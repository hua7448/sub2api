//go:build perf

package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Wei-Shaw/sub2api/ent"
	_ "github.com/Wei-Shaw/sub2api/ent/runtime"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// Provider hall performance acceptance (plan matrix H). Opt-in via the perf
// build tag; owns an isolated PostgreSQL cluster like the localdb suite.
//
//	PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin \
//	PROVIDER_HALL_PG_PORT=15495 go test -tags=perf ./internal/repository -run ProviderHallPerf -v -timeout 30m
const (
	perfGroupCount   = 100
	perfRequestRows  = 1_000_000
	perfRequestDays  = 8 // one day beyond retention so prune has work
	perfSnapshotDays = 30
	perfTier1Hours   = 48
	perfConcurrency  = 20
)

func perfEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func perfPercentile(d []time.Duration, p float64) time.Duration {
	if len(d) == 0 {
		return 0
	}
	c := append([]time.Duration(nil), d...)
	sort.Slice(c, func(i, j int) bool { return c[i] < c[j] })
	idx := int(float64(len(c)-1)*p + 0.5)
	if idx >= len(c) {
		idx = len(c) - 1
	}
	return c[idx]
}

// --- statement logging connector (perf copy of the b5_read helper) ---------

type perfStatement struct {
	sql  string
	args []driver.NamedValue
}

type perfStatementLog struct {
	mu    sync.Mutex
	items []perfStatement
}

func (l *perfStatementLog) add(query string, args []driver.NamedValue) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = append(l.items, perfStatement{sql: query, args: append([]driver.NamedValue(nil), args...)})
}

func (l *perfStatementLog) reset() []perfStatement {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := l.items
	l.items = nil
	return out
}

type perfConnector struct {
	base driver.Connector
	log  *perfStatementLog
}

func (c perfConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.base.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &perfConn{Conn: conn, log: c.log}, nil
}
func (c perfConnector) Driver() driver.Driver { return c.base.Driver() }

type perfConn struct {
	driver.Conn
	log *perfStatementLog
}

func (c *perfConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	c.log.add(query, args)
	if q, ok := c.Conn.(driver.QueryerContext); ok {
		return q.QueryContext(ctx, query, args)
	}
	return nil, driver.ErrSkip
}
func (c *perfConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	c.log.add(query, args)
	if e, ok := c.Conn.(driver.ExecerContext); ok {
		return e.ExecContext(ctx, query, args)
	}
	return nil, driver.ErrSkip
}
func (c *perfConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	c.log.add(query, nil)
	if p, ok := c.Conn.(driver.ConnPrepareContext); ok {
		return p.PrepareContext(ctx, query)
	}
	return c.Conn.Prepare(query)
}
func (c *perfConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if b, ok := c.Conn.(driver.ConnBeginTx); ok {
		return b.BeginTx(ctx, opts)
	}
	return c.Conn.Begin()
}
func (c *perfConn) Ping(ctx context.Context) error {
	if p, ok := c.Conn.(driver.Pinger); ok {
		return p.Ping(ctx)
	}
	return nil
}

type perfAccess struct{ groups []service.Group }

func (a perfAccess) GetAvailableGroups(context.Context, int64) ([]service.Group, error) {
	return a.groups, nil
}

type perfRates struct{}

func (perfRates) ResolveUserGroupRateMultiplier(_ context.Context, _, _ int64, def float64) float64 {
	return def
}

type perfPricing struct{}

func (perfPricing) Resolve(context.Context, service.PricingInput) *service.ResolvedPricing {
	return &service.ResolvedPricing{Mode: service.BillingModeToken, BasePricing: &service.ModelPricing{InputPricePerToken: 0.000002, CacheReadPricePerToken: 0.0000002}}
}

func TestProviderHallPerf(t *testing.T) {
	bin := os.Getenv("PROVIDER_HALL_PG_BIN")
	require.NotEmpty(t, bin, "set PROVIDER_HALL_PG_BIN to the PostgreSQL bin directory")
	port := perfEnv("PROVIDER_HALL_PG_PORT", "15495")
	dir, err := os.MkdirTemp("/tmp", "provider-hall-perf-")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	data := filepath.Join(dir, "data")
	run := func(name string, args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		out, err := exec.CommandContext(ctx, filepath.Join(bin, name), args...).CombinedOutput()
		require.NoError(t, err, "%s: %s", name, out)
	}
	run("initdb", "-D", data, "-U", "hall_test", "--auth=trust", "--no-locale", "-E", "UTF8")
	run("pg_ctl", "-D", data, "-l", filepath.Join(dir, "postgres.log"), "-o",
		fmt.Sprintf("-h '' -k %s -p %s -c shared_buffers=512MB -c work_mem=64MB -c maintenance_work_mem=256MB -c synchronous_commit=off", dir, port), "-w", "start")
	t.Cleanup(func() { run("pg_ctl", "-D", data, "-m", "immediate", "-w", "stop") })
	dsn := func(name string) string {
		return fmt.Sprintf("host=%s port=%s user=hall_test dbname=%s sslmode=disable", dir, port, name)
	}
	connect := func(name string) *sql.DB {
		db, err := sql.Open("postgres", dsn(name))
		require.NoError(t, err)
		db.SetMaxOpenConns(64)
		t.Cleanup(func() { _ = db.Close() })
		return db
	}
	admin := connect("postgres")
	_, err = admin.Exec("CREATE DATABASE hall_perf")
	require.NoError(t, err)
	db := connect("hall_perf")
	ctx := context.Background()
	require.NoError(t, ApplyMigrations(ctx, db))

	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	cfgRepo := NewProviderHallRepository(client)
	now := time.Now().UTC()
	// Newest published window for the read tests: a tier-5 boundary in the past.
	T := now.Add(-2 * time.Minute).Truncate(5 * time.Minute)

	// --- fixtures: profiles, groups, listings, targets -----------------------
	protocols := []string{"responses", "chat_completions"}
	profiles := make([]*service.ProviderHallProfile, 0, 2)
	for i, protocol := range protocols {
		p, err := cfgRepo.SaveProfile(ctx, service.ProviderHallProfile{Model: fmt.Sprintf("hall-perf-%d", i), Protocol: protocol, OutputLimit: 256, ModelAliases: []string{}})
		require.NoError(t, err)
		profiles = append(profiles, p)
	}
	groupIDs := make([]int64, 0, perfGroupCount)
	groups := make([]service.Group, 0, perfGroupCount)
	for i := 0; i < perfGroupCount; i++ {
		g, err := client.Group.Create().SetName(fmt.Sprintf("hall-perf-%03d", i)).SetPlatform(service.PlatformOpenAI).SetRateMultiplier(1.5).Save(ctx)
		require.NoError(t, err)
		groupIDs = append(groupIDs, g.ID)
		groups = append(groups, service.Group{ID: g.ID, Name: g.Name, Platform: service.PlatformOpenAI, Status: service.StatusActive, RateMultiplier: 1.5})
		_, err = db.ExecContext(ctx, `INSERT INTO provider_hall_groups (group_id, listed, display_name, display_order) VALUES ($1, true, $2, $3)`, g.ID, fmt.Sprintf("Perf %03d", i), perfGroupCount-i)
		require.NoError(t, err)
		keyID := int64(900_000_000) + g.ID
		_, err = db.ExecContext(ctx, `INSERT INTO provider_hall_probe_keys (api_key_id, operator_user_id, group_id, registered_by) VALUES ($1, 1, $2, 1)`, keyID, g.ID)
		require.NoError(t, err)
		for _, p := range profiles {
			_, err = db.ExecContext(ctx, `INSERT INTO provider_hall_targets (group_id, profile_id, probe_key_id, enabled) VALUES ($1, $2, $3, true)`, g.ID, p.ID, keyID)
			require.NoError(t, err)
		}
	}
	cfg, err := cfgRepo.GetConfig(ctx)
	require.NoError(t, err)
	cfg.DisplayEnabled, cfg.CollectionEnabled = true, true
	cfg.DefaultModel, cfg.DefaultProtocol = profiles[0].Model, "responses"
	_, err = cfgRepo.UpdateConfig(ctx, *cfg)
	require.NoError(t, err)

	// --- snapshots: 30 days tier 5 + 48 h tier 1, merged row and default profile
	started := time.Now()
	for _, profileID := range []int64{0, profiles[0].ID} {
		_, err = db.ExecContext(ctx, `
INSERT INTO provider_hall_snapshots (group_id, profile_id, window_end, tier, algorithm_version, ttft_sample_count, ttft_fast95_mean_ms, ttft_p90_ms,
    success_count, failed_count, submissions, success_rate, usage_success_count, input_tokens, cache_read_tokens, cache_creation_tokens, cache_rate,
    billed_count, billed_input_tokens, billed_input_cost, historical_price, subscription_share, coverage, billing_coverage)
SELECT g.id, $3, w.window_end, 5, 1, 40, 500 + g.ord * 10 + (extract(epoch from w.window_end)::bigint % 97), 900,
       95, 5, 100, 0.95, 50, 1000, 300, 0, 0.3, 250, 1000, 0.0015, 1.5, 0, 'complete', 'complete'
FROM unnest($1::bigint[]) WITH ORDINALITY AS g(id, ord),
     (SELECT $2::timestamptz - (i * interval '5 minutes') AS window_end FROM generate_series(0, $4) AS i) AS w`,
			pq.Array(groupIDs), T, profileID, perfSnapshotDays*24*12)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, `
INSERT INTO provider_hall_snapshots (group_id, profile_id, window_end, tier, algorithm_version, ttft_sample_count, ttft_fast95_mean_ms, ttft_p90_ms,
    success_count, failed_count, submissions, success_rate, usage_success_count, input_tokens, cache_read_tokens, cache_creation_tokens, cache_rate,
    billed_count, billed_input_tokens, billed_input_cost, historical_price, subscription_share, coverage, billing_coverage)
SELECT g.id, $3, w.window_end, 1, 1, 40, 500 + g.ord * 10, 900,
       95, 5, 100, 0.95, 50, 1000, 300, 0, 0.3, 250, 1000, 0.0015, 1.5, 0, 'complete', 'complete'
FROM unnest($1::bigint[]) WITH ORDINALITY AS g(id, ord),
     (SELECT $2::timestamptz - (i * interval '1 minute') AS window_end FROM generate_series(1, $4) AS i WHERE i % 5 <> 0) AS w`,
			pq.Array(groupIDs), T, profileID, perfTier1Hours*60)
		require.NoError(t, err)
	}
	var snapshotRows int64
	require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_snapshots`).Scan(&snapshotRows))
	t.Logf("seeded %d snapshot rows in %s", snapshotRows, time.Since(started).Round(time.Millisecond))
	_, err = db.ExecContext(ctx, `UPDATE provider_hall_aggregator_state SET last_window_end = $1, algorithm_version = 1 WHERE id = 1`, T)
	require.NoError(t, err)

	// --- requests: 1,000,000 user facts over 8 days via COPY ------------------
	started = time.Now()
	factRepo := NewProviderHallFactRepository(db)
	epochID, err := factRepo.RegisterEpoch(ctx, "perf-node", "perf", service.ProviderHallAlgorithmVersion, now.Add(-perfRequestDays*24*time.Hour))
	require.NoError(t, err)
	confirmed := now.Add(-30 * time.Second)
	require.NoError(t, factRepo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epochID, MaxSeq: perfRequestRows, ConfirmedAt: &confirmed}))
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("provider_hall_requests",
		"trace_id", "node_id", "epoch_id", "last_seq", "group_id", "profile_id", "api_key_id", "protocol",
		"requested_model", "source", "stream", "started_at", "first_content_at", "ended_at", "ttft_ms", "submissions", "outcome",
		"usage_known", "input_tokens", "output_tokens", "cache_read_tokens", "cache_creation_tokens",
		"billing_request_id", "billing_api_key_id", "billing_status", "billing_mode", "is_subscription", "rate_multiplier",
		"input_base_cost", "total_base_cost", "actual_cost", "billed_at", "algorithm_version"))
	require.NoError(t, err)
	rng := rand.New(rand.NewSource(20260912))
	spanSeconds := int64(perfRequestDays * 24 * 3600)
	origin := now.Add(-time.Duration(spanSeconds) * time.Second)
	for i := 0; i < perfRequestRows; i++ {
		g := groupIDs[rng.Intn(perfGroupCount)]
		pi := rng.Intn(len(profiles))
		start := origin.Add(time.Duration(rng.Int63n(spanSeconds-60)) * time.Second).Add(time.Duration(rng.Intn(1000)) * time.Millisecond)
		ttft := 100 + rng.Intn(1900)
		total := ttft + 100 + rng.Intn(2500)
		first := start.Add(time.Duration(ttft) * time.Millisecond)
		ended := start.Add(time.Duration(total) * time.Millisecond)
		outcome := "success"
		var ttftArg, firstArg any = ttft, first
		if rng.Intn(20) == 0 {
			outcome, ttftArg, firstArg = "failed", nil, nil
		}
		in, out, cached := int64(500+rng.Intn(2500)), int64(50+rng.Intn(400)), int64(rng.Intn(400))
		inputCost := float64(in) * 0.000002
		totalCost := inputCost + float64(out)*0.000008
		_, err = stmt.ExecContext(ctx,
			uuid.New().String(), "perf-node", epochID, int64(i+1), g, profiles[pi].ID, int64(1), protocols[pi],
			profiles[pi].Model, "user", true, start, firstArg, ended, ttftArg, 1, outcome,
			outcome == "success", in, out, cached, int64(0),
			fmt.Sprintf("client:%d", i), int64(1), "applied", "token", false, 1.5,
			fmt.Sprintf("%.10f", inputCost), fmt.Sprintf("%.10f", totalCost), fmt.Sprintf("%.8f", totalCost*1.5), ended.Add(time.Second), 1)
		require.NoError(t, err)
	}
	_, err = stmt.ExecContext(ctx)
	require.NoError(t, err)
	require.NoError(t, stmt.Close())
	require.NoError(t, tx.Commit())
	var requestRows int64
	require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_requests`).Scan(&requestRows))
	require.Equal(t, int64(perfRequestRows), requestRows)
	t.Logf("copied %d request rows in %s", requestRows, time.Since(started).Round(time.Millisecond))
	_, err = db.ExecContext(ctx, `ANALYZE`)
	require.NoError(t, err)

	newService := func(base *sql.DB) *service.ProviderHallQueryService {
		c := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, base)))
		cfgR := NewProviderHallRepository(c)
		hall := service.NewProviderHallService(cfgR, nil, nil, nil)
		access := perfAccess{groups: groups}
		hall.SetGroupAccessForTest(access)
		svc := service.NewProviderHallQueryService(NewProviderHallReadRepository(base), cfgR, hall, access, perfRates{}, perfPricing{}, NewProviderHallJobRepository(base))
		svc.SetClockForTest(func() time.Time { return T.Add(2 * time.Minute) })
		return svc
	}

	t.Run("read_plans_avoid_fact_tables", func(t *testing.T) {
		connector, err := pq.NewConnector(dsn("hall_perf"))
		require.NoError(t, err)
		log := &perfStatementLog{}
		counted := sql.OpenDB(perfConnector{base: connector, log: log})
		defer func() { _ = counted.Close() }()
		svc := newService(counted)
		log.reset()
		res, err := svc.List(ctx, 1, service.ProviderHallListQuery{Range: "30d", PageSize: 100})
		require.NoError(t, err)
		require.Len(t, res.Items, perfGroupCount)
		listStatements := log.reset()
		_, err = svc.GetGroup(ctx, 1, groupIDs[7], "7d")
		require.NoError(t, err)
		detailStatements := log.reset()
		t.Logf("cold list statements: %d; detail statements: %d", len(listStatements), len(detailStatements))
		require.LessOrEqual(t, len(listStatements), 12)
		explained := 0
		for _, st := range append(listStatements, detailStatements...) {
			lower := strings.ToLower(st.sql)
			require.NotContains(t, lower, "usage_logs", st.sql)
			require.NotContains(t, lower, "provider_hall_requests", st.sql)
			if !strings.HasPrefix(strings.TrimSpace(lower), "select") {
				continue
			}
			args := make([]any, 0, len(st.args))
			for _, a := range st.args {
				args = append(args, a.Value)
			}
			rows, err := db.QueryContext(ctx, "EXPLAIN (FORMAT JSON) "+st.sql, args...)
			require.NoError(t, err, st.sql)
			var plan strings.Builder
			for rows.Next() {
				var line string
				require.NoError(t, rows.Scan(&line))
				plan.WriteString(line)
			}
			_ = rows.Close()
			require.NotContains(t, plan.String(), "usage_logs", st.sql)
			require.NotContains(t, plan.String(), "provider_hall_requests", st.sql)
			explained++
		}
		t.Logf("explained %d SELECT statements; none touch usage_logs or provider_hall_requests", explained)
	})

	t.Run("list_concurrency_p95", func(t *testing.T) {
		measure := func(svc *service.ProviderHallQueryService, rounds int) []time.Duration {
			var mu sync.Mutex
			var wg sync.WaitGroup
			var samples []time.Duration
			for r := 0; r < rounds; r++ {
				for w := 0; w < perfConcurrency; w++ {
					wg.Add(1)
					go func(user int64) {
						defer wg.Done()
						t0 := time.Now()
						res, err := svc.List(ctx, user, service.ProviderHallListQuery{Range: "30d", PageSize: 50})
						d := time.Since(t0)
						mu.Lock()
						defer mu.Unlock()
						if err == nil && len(res.Items) != 50 {
							err = fmt.Errorf("unexpected page size %d", len(res.Items))
						}
						require.NoError(t, err)
						samples = append(samples, d)
					}(int64(w + 1))
				}
				wg.Wait()
			}
			return samples
		}
		cold := measure(newService(db), 1)
		coldP95 := perfPercentile(cold, 0.95)
		warmSvc := newService(db)
		_, err := warmSvc.List(ctx, 1, service.ProviderHallListQuery{Range: "30d", PageSize: 50})
		require.NoError(t, err)
		warm := measure(warmSvc, 5)
		warmP95 := perfPercentile(warm, 0.95)
		t.Logf("list P95: cold %s (n=%d, max %s) / warm %s (n=%d, max %s)", coldP95.Round(time.Millisecond), len(cold), perfPercentile(cold, 1).Round(time.Millisecond), warmP95.Round(time.Millisecond), len(warm), perfPercentile(warm, 1).Round(time.Millisecond))
		require.LessOrEqual(t, coldP95, 2*time.Second, "cold list P95")
		require.LessOrEqual(t, warmP95, 500*time.Millisecond, "warm list P95")

		var detail []time.Duration
		for i := 0; i < 20; i++ {
			t0 := time.Now()
			_, err := warmSvc.GetGroup(ctx, 1, groupIDs[i], "7d")
			require.NoError(t, err)
			detail = append(detail, time.Since(t0))
		}
		t.Logf("detail P95 (sequential, 20 groups): %s", perfPercentile(detail, 0.95).Round(time.Millisecond))
	})

	t.Run("aggregator_tick_and_prune", func(t *testing.T) {
		aggRepo := NewProviderHallAggregationRepository(db)
		targets, err := factRepo.ListEnabledTargets(ctx)
		require.NoError(t, err)
		require.Len(t, targets, perfGroupCount*len(profiles))
		aggNow := time.Now().UTC()
		aggT := aggNow.Add(-45 * time.Second).Truncate(time.Minute)
		aggW := aggT.Add(-240 * time.Minute)
		_, err = db.ExecContext(ctx, `UPDATE provider_hall_aggregator_state SET last_window_end = $1 WHERE id = 1`, aggW)
		require.NoError(t, err)

		t0 := time.Now()
		rec, err := aggRepo.Reconcile(ctx, aggNow, 500)
		require.NoError(t, err)
		reconcileDur := time.Since(t0)

		minutes := make([]time.Time, 0, 240)
		for m := aggW; m.Before(aggT); m = m.Add(time.Minute) {
			minutes = append(minutes, m)
		}
		t0 = time.Now()
		res, err := aggRepo.Recompute(ctx, service.ProviderHallRecomputeInput{
			Now: aggNow, T: aggT, Minutes: minutes, NewWindowEnd: aggT,
			ReevalFrom: aggT.Add(-65 * time.Minute), ReevalTo: aggT,
			Targets: targets, AlgorithmVersion: service.ProviderHallAlgorithmVersion,
		})
		require.NoError(t, err)
		recomputeDur := time.Since(t0)
		t.Logf("reconcile examined=%d in %s; recompute 240 minutes: %d minutes, %d snapshots in %s", rec.Examined, reconcileDur.Round(time.Millisecond), res.MinutesRecomputed, res.SnapshotsPublished, recomputeDur.Round(time.Millisecond))
		require.Less(t, recomputeDur, 55*time.Second, "240-minute recompute must fit the tick budget")

		t0 = time.Now()
		pruned, err := aggRepo.Prune(ctx, aggNow)
		require.NoError(t, err)
		pruneDur := time.Since(t0)
		t.Logf("prune in %s: %v", pruneDur.Round(time.Millisecond), pruned.Deleted)
		require.Greater(t, pruned.Deleted["requests"], int64(0), "day 8 facts fall out of the 7-day retention")

		// Full tick through the aggregator (lock, node maintenance, reconcile,
		// recompute) against a fresh 240-minute backlog.
		_, err = db.ExecContext(ctx, `UPDATE provider_hall_aggregator_state SET last_window_end = $1 WHERE id = 1`, aggW)
		require.NoError(t, err)
		agg := service.NewProviderHallAggregator(aggRepo, factRepo, cfgRepo, db, nil, nil)
		t0 = time.Now()
		agg.RunOnce()
		tickDur := time.Since(t0)
		_, lastErr := agg.Status()
		require.Empty(t, lastErr)
		var lastWindow time.Time
		require.NoError(t, db.QueryRowContext(ctx, `SELECT last_window_end FROM provider_hall_aggregator_state WHERE id = 1`).Scan(&lastWindow))
		t.Logf("aggregator RunOnce (240-minute backlog incl. prune on tick 1): %s, watermark now %s", tickDur.Round(time.Millisecond), lastWindow.UTC().Format(time.RFC3339))
		require.Less(t, tickDur, 55*time.Second)
	})
}
