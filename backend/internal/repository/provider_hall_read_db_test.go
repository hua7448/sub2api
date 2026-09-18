//go:build integration || providerhall_localdb

package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func init() {
	registerProviderHallLocalDBContract("b5_read", providerHallReadDatabaseContracts)
}

// hallStatementLog wraps a pq connection and records every statement so the
// contract can assert a fixed statement count and EXPLAIN each read.
type hallStatement struct {
	sql  string
	args []driver.NamedValue
}

type hallStatementLog struct {
	mu    sync.Mutex
	items []hallStatement
}

func (l *hallStatementLog) add(query string, args []driver.NamedValue) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.items = append(l.items, hallStatement{sql: query, args: append([]driver.NamedValue(nil), args...)})
}

func (l *hallStatementLog) reset() []hallStatement {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := l.items
	l.items = nil
	return out
}

type hallCountingConnector struct {
	base driver.Connector
	log  *hallStatementLog
}

func (c hallCountingConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.base.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &hallCountingConn{Conn: conn, log: c.log}, nil
}
func (c hallCountingConnector) Driver() driver.Driver { return c.base.Driver() }

type hallCountingConn struct {
	driver.Conn
	log *hallStatementLog
}

func (c *hallCountingConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	c.log.add(query, args)
	if q, ok := c.Conn.(driver.QueryerContext); ok {
		return q.QueryContext(ctx, query, args)
	}
	return nil, driver.ErrSkip
}
func (c *hallCountingConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	c.log.add(query, args)
	if e, ok := c.Conn.(driver.ExecerContext); ok {
		return e.ExecContext(ctx, query, args)
	}
	return nil, driver.ErrSkip
}
func (c *hallCountingConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	c.log.add(query, nil)
	if p, ok := c.Conn.(driver.ConnPrepareContext); ok {
		return p.PrepareContext(ctx, query)
	}
	return c.Conn.Prepare(query)
}
func (c *hallCountingConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if b, ok := c.Conn.(driver.ConnBeginTx); ok {
		return b.BeginTx(ctx, opts)
	}
	return c.Conn.Begin()
}
func (c *hallCountingConn) Ping(ctx context.Context) error {
	if p, ok := c.Conn.(driver.Pinger); ok {
		return p.Ping(ctx)
	}
	return nil
}

// hallCountingDB reopens the harness database through a statement-logging
// connector. The DSN is rebuilt from the live session (socket dir + port).
func hallCountingDB(t *testing.T, db *sql.DB) (*sql.DB, *hallStatementLog) {
	t.Helper()
	var dir, port, name string
	require.NoError(t, db.QueryRow(`SHOW unix_socket_directories`).Scan(&dir))
	require.NoError(t, db.QueryRow(`SHOW port`).Scan(&port))
	require.NoError(t, db.QueryRow(`SELECT current_database()`).Scan(&name))
	dir = strings.Split(dir, ",")[0]
	connector, err := pq.NewConnector(fmt.Sprintf("host=%s port=%s user=hall_test dbname=%s sslmode=disable", strings.TrimSpace(dir), port, name))
	require.NoError(t, err)
	log := &hallStatementLog{}
	counted := sql.OpenDB(hallCountingConnector{base: connector, log: log})
	t.Cleanup(func() { _ = counted.Close() })
	return counted, log
}

type hallReadAccess struct{ groups []service.Group }

func (a hallReadAccess) GetAvailableGroups(context.Context, int64) ([]service.Group, error) {
	return a.groups, nil
}

type hallReadRates struct{}

func (hallReadRates) ResolveUserGroupRateMultiplier(_ context.Context, _, _ int64, def float64) float64 {
	return def
}

type hallReadPricing struct{}

func (hallReadPricing) Resolve(context.Context, service.PricingInput) *service.ResolvedPricing {
	return &service.ResolvedPricing{Mode: service.BillingModeToken, BasePricing: &service.ModelPricing{InputPricePerToken: 0.000002, CacheReadPricePerToken: 0.0000002}}
}

func providerHallReadDatabaseContracts(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	cfgRepo := NewProviderHallRepository(client)
	const groupCount = 100
	T := time.Date(2026, 9, 12, 12, 0, 0, 0, time.UTC)
	stamp := time.Now().UnixNano()

	// Two profiles, 100 listed groups each with two enabled targets, tier-5
	// snapshots for 30 days (120 six-hour buckets) plus the T window.
	profiles := make([]*service.ProviderHallProfile, 0, 2)
	for i, protocol := range []string{"responses", "chat_completions"} {
		p, err := cfgRepo.SaveProfile(ctx, service.ProviderHallProfile{Model: fmt.Sprintf("hall-read-%d-%d", stamp, i), Protocol: protocol, OutputLimit: 256, ModelAliases: []string{}})
		require.NoError(t, err)
		profiles = append(profiles, p)
	}
	groupIDs := make([]int64, 0, groupCount)
	groups := make([]service.Group, 0, groupCount)
	for i := 0; i < groupCount; i++ {
		g, err := client.Group.Create().SetName(fmt.Sprintf("hall-read-%d-%d", stamp, i)).SetPlatform(service.PlatformOpenAI).SetRateMultiplier(1.5).Save(ctx)
		require.NoError(t, err)
		groupIDs = append(groupIDs, g.ID)
		groups = append(groups, service.Group{ID: g.ID, Name: g.Name, Platform: service.PlatformOpenAI, Status: service.StatusActive, RateMultiplier: 1.5})
		_, err = db.ExecContext(ctx, `INSERT INTO provider_hall_groups (group_id, listed, display_name, display_order) VALUES ($1, true, $2, $3)`, g.ID, fmt.Sprintf("Read %03d", i), groupCount-i)
		require.NoError(t, err)
		// Enabled targets require a registered probe key; the registry has no
		// FK to api_keys, so a synthetic id per group is enough here.
		keyID := int64(900_000_000) + g.ID
		_, err = db.ExecContext(ctx, `INSERT INTO provider_hall_probe_keys (api_key_id, operator_user_id, group_id, registered_by) VALUES ($1, 1, $2, 1)`, keyID, g.ID)
		require.NoError(t, err)
		for _, p := range profiles {
			_, err = db.ExecContext(ctx, `INSERT INTO provider_hall_targets (group_id, profile_id, probe_key_id, enabled) VALUES ($1, $2, $3, true)`, g.ID, p.ID, keyID)
			require.NoError(t, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM provider_hall_snapshots WHERE group_id = ANY($1)`, pq.Array(groupIDs))
		_, _ = db.ExecContext(ctx, `DELETE FROM groups WHERE id = ANY($1)`, pq.Array(groupIDs))
		_, _ = db.ExecContext(ctx, `DELETE FROM provider_hall_probe_keys WHERE group_id = ANY($1)`, pq.Array(groupIDs))
		for _, p := range profiles {
			_, _ = db.ExecContext(ctx, `DELETE FROM provider_hall_profiles WHERE id=$1`, p.ID)
		}
		_, _ = db.ExecContext(ctx, `UPDATE provider_hall_aggregator_state SET last_window_end = NULL WHERE id = 1`)
	})
	// Merged rows: 121 windows per group (T and 120 six-hour ends); the
	// default-profile row only at T. fast95 encodes the group index so the
	// sort can be verified.
	_, err := db.ExecContext(ctx, `
INSERT INTO provider_hall_snapshots (group_id, profile_id, window_end, tier, algorithm_version, ttft_sample_count, ttft_fast95_mean_ms, ttft_p90_ms,
    success_count, failed_count, submissions, success_rate, usage_success_count, input_tokens, cache_read_tokens, cache_creation_tokens, cache_rate,
    billed_count, billed_input_tokens, billed_input_cost, historical_price, subscription_share, coverage, billing_coverage)
SELECT g.id, 0, w.window_end, 5, 1, 40, 500 + g.ord * 10, 900, 95, 5, 100, 0.95, 50, 1000, 300, 0, 0.3,
       250, 1000, 0.0015, 1.5, 0, 'complete', 'complete'
FROM unnest($1::bigint[]) WITH ORDINALITY AS g(id, ord),
     (SELECT $2::timestamptz - (i * interval '6 hours') AS window_end FROM generate_series(0, 120) AS i) AS w`, pq.Array(groupIDs), T)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `
INSERT INTO provider_hall_snapshots (group_id, profile_id, window_end, tier, algorithm_version, ttft_sample_count, ttft_fast95_mean_ms, ttft_p90_ms,
    success_count, failed_count, submissions, success_rate, usage_success_count, input_tokens, cache_read_tokens, cache_creation_tokens, cache_rate,
    billed_count, billed_input_tokens, billed_input_cost, historical_price, subscription_share, coverage, billing_coverage)
SELECT g.id, $3, $2::timestamptz, 5, 1, 40, 600, 900, 95, 5, 100, 0.95, 50, 1000, 300, 0, 0.3,
       250, 1000, 0.0015, 1.5, 0, 'complete', 'complete'
FROM unnest($1::bigint[]) AS g(id)`, pq.Array(groupIDs), T, profiles[0].ID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `UPDATE provider_hall_aggregator_state SET last_window_end = $1, algorithm_version = 1 WHERE id = 1`, T)
	require.NoError(t, err)
	original, err := cfgRepo.GetConfig(ctx)
	require.NoError(t, err)
	cfg := *original
	cfg.DisplayEnabled, cfg.CollectionEnabled = true, true
	cfg.DefaultModel, cfg.DefaultProtocol = profiles[0].Model, "responses"
	// Readiness gates are enforced by the service, not the repository.
	_, err = cfgRepo.UpdateConfig(ctx, cfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		current, err := cfgRepo.GetConfig(ctx)
		require.NoError(t, err)
		original.Version = current.Version
		_, _ = cfgRepo.UpdateConfig(ctx, *original)
	})

	counted, log := hallCountingDB(t, db)
	countedClient := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, counted)))
	countedCfg := NewProviderHallRepository(countedClient)
	newService := func(access hallReadAccess) *service.ProviderHallQueryService {
		hall := service.NewProviderHallService(countedCfg, nil, nil, nil)
		hall.SetGroupAccessForTest(access)
		svc := service.NewProviderHallQueryService(NewProviderHallReadRepository(counted), countedCfg, hall, access, hallReadRates{}, hallReadPricing{}, NewProviderHallJobRepository(counted))
		svc.SetClockForTest(func() time.Time { return T.Add(2 * time.Minute) })
		return svc
	}

	t.Run("read_list_values_and_sort", func(t *testing.T) {
		svc := newService(hallReadAccess{groups: groups})
		rules, ok := service.ParseProviderHallSort("ttft_fast95:asc")
		require.True(t, ok)
		res, err := svc.List(ctx, 1, service.ProviderHallListQuery{Range: "30d", Sort: rules, PageSize: 100})
		require.NoError(t, err)
		require.Equal(t, groupCount, res.Total)
		require.Len(t, res.Items, groupCount)
		require.Equal(t, groupIDs[0], res.Items[0].GroupID, "lowest fast95 first")
		require.Equal(t, groupIDs[groupCount-1], res.Items[groupCount-1].GroupID)
		first := res.Items[0]
		require.Equal(t, "Read 000", first.Name)
		require.Equal(t, 510.0, *first.TTFTFast95Ms.Value)
		require.Equal(t, service.ProviderHallMetricOK, first.TTFTFast95Ms.State)
		require.Equal(t, 0.3, *first.CacheRate.Value)
		require.Equal(t, "1.5000000000", *first.HistoricalPrice.Value)
		require.Equal(t, "1.5", first.Rate)
		require.Equal(t, "3.0000000000", first.Quote.InputPrice)
		require.Equal(t, profiles[0].ID, first.DefaultProfile.ProfileID)
		require.Len(t, first.Sparkline.Points, 120)
		require.NotNil(t, first.Sparkline.Points[119].RealTTFTMs, "the newest bucket end is T")
		require.False(t, first.Sparkline.Points[119].Gap)
		require.NotNil(t, first.Sparkline.Points[0].RealTTFTMs, "30 days of tier-5 rows are present")
		require.Equal(t, service.ProviderHallHealthUnknown, first.Health.Status, "no probes seeded")
		require.Len(t, first.Health.Models, 2)
		require.Len(t, res.Catalog.Models, 2)
		require.Equal(t, T, *res.DataThrough)

		// Default order is display_order: the last created group sorts first.
		res, err = svc.List(ctx, 1, service.ProviderHallListQuery{Range: "6h", PageSize: 5})
		require.NoError(t, err)
		require.Equal(t, groupIDs[groupCount-1], res.Items[0].GroupID)
		require.Len(t, res.Items[0].Sparkline.Points, 72)
		require.NotNil(t, res.Items[0].Sparkline.Points[71].RealTTFTMs)
		require.Nil(t, res.Items[0].Sparkline.Points[70].RealTTFTMs, "only six-hour ends were seeded")
		require.True(t, res.Items[0].Sparkline.Points[70].Gap)

		// Detail on the counted path.
		d, err := svc.GetGroup(ctx, 1, groupIDs[3], "7d")
		require.NoError(t, err)
		require.Equal(t, int64(900), *d.P90Ms.Value)
		require.Len(t, d.Trend, 168)
		require.Len(t, d.Profiles, 2)
		require.Equal(t, service.ProviderHallReasonSamplesBelow12, d.E2EAvailability6h.ReasonCode)
		reports, total, err := svc.ListVerifications(ctx, 1, groupIDs[3], nil, 1, 10)
		require.NoError(t, err)
		require.Zero(t, total)
		require.Empty(t, reports)
	})

	t.Run("read_statement_count_is_fixed_and_plans_avoid_facts", func(t *testing.T) {
		countList := func(access hallReadAccess, rangeKey string) []hallStatement {
			svc := newService(access) // fresh service ⇒ cold base cache
			log.reset()
			res, err := svc.List(ctx, 1, service.ProviderHallListQuery{Range: rangeKey, PageSize: 100})
			require.NoError(t, err)
			require.Len(t, res.Items, len(access.groups))
			return log.reset()
		}
		all := countList(hallReadAccess{groups: groups}, "30d")
		few := countList(hallReadAccess{groups: groups[:5]}, "30d")
		t.Logf("cold list statements: %d (100 groups) vs %d (5 groups)", len(all), len(few))
		require.Equal(t, len(all), len(few), "statement count must not depend on the number of visible groups")
		require.LessOrEqual(t, len(all), 12, "cold list issues a bounded, fixed number of statements: %d", len(all))
		for _, st := range all {
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
			var plan strings.Builder
			rows, err := db.QueryContext(ctx, "EXPLAIN (FORMAT JSON) "+st.sql, args...)
			require.NoError(t, err, st.sql)
			for rows.Next() {
				var line string
				require.NoError(t, rows.Scan(&line))
				plan.WriteString(line)
			}
			_ = rows.Close()
			require.NotContains(t, plan.String(), "usage_logs", st.sql)
			require.NotContains(t, plan.String(), "provider_hall_requests", st.sql)
		}
		// A warm cache serves the same list without touching the aggregate
		// tables again; only the switch, permissions and window are re-read.
		svc := newService(hallReadAccess{groups: groups})
		_, err := svc.List(ctx, 1, service.ProviderHallListQuery{Range: "24h"})
		require.NoError(t, err)
		log.reset()
		_, err = svc.List(ctx, 2, service.ProviderHallListQuery{Range: "24h"})
		require.NoError(t, err)
		warm := log.reset()
		for _, st := range warm {
			require.NotContains(t, strings.ToLower(st.sql), "provider_hall_snapshots", "warm request re-read snapshots: %s", st.sql)
		}
		require.Less(t, len(warm), len(all))
	})
}
