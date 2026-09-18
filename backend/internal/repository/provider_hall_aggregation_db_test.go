//go:build integration || providerhall_localdb

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	registerProviderHallLocalDBContract("b3_aggregation", providerHallAggregationDatabaseContracts)
}

// hallReq describes one fact row for the aggregation contracts.
type hallReq struct {
	group, profile int64
	endedAt        time.Time
	outcome        string
	ttft           int
	stream         bool
	submissions    int
	usage          bool
	input, read    int64
	billing        string // pending|applied|...
	subscription   bool
	requestID      string
	apiKeyID       int64
	startedAt      time.Time
	source         string
}

func insertHallReq(t *testing.T, db *sql.DB, r hallReq) uuid.UUID {
	t.Helper()
	id := uuid.New()
	started := r.startedAt
	if started.IsZero() {
		started = r.endedAt.Add(-2 * time.Second)
	}
	if r.submissions == 0 && r.outcome != "excluded" {
		r.submissions = 1
	}
	if r.billing == "" {
		r.billing = "pending"
	}
	if r.source == "" {
		r.source = "user"
	}
	if r.apiKeyID == 0 {
		r.apiKeyID = 1
	}
	var ttft any
	if r.ttft > 0 {
		ttft = r.ttft
	}
	var profile any
	if r.profile != 0 {
		profile = r.profile
	}
	var ended any
	if !r.endedAt.IsZero() {
		ended = r.endedAt
	}
	var reqID any
	if r.requestID != "" {
		reqID = r.requestID
	}
	var mode, actual, inputCost, totalCost any
	if r.billing == "applied" {
		mode, actual, inputCost, totalCost = "token", "0.004", "0.001", "0.002"
	}
	_, err := db.Exec(`INSERT INTO provider_hall_requests (trace_id, node_id, epoch_id, last_seq, group_id, profile_id, api_key_id, protocol, requested_model, source, stream,
		started_at, ended_at, ttft_ms, submissions, outcome, usage_known, input_tokens, output_tokens, cache_read_tokens, cache_creation_tokens,
		billing_request_id, billing_api_key_id, billing_status, billing_mode, is_subscription, actual_cost, input_base_cost, total_base_cost, algorithm_version)
		VALUES ($1, 'n1', 1, 1, $2, $3, $4, 'responses', 'gpt-hall', $5, $6, $7, $8, $9, $10, $11, $12, $13, 5, $14, 0, $15, $4, $16, $17, $18, $19::numeric, $20::numeric, $21::numeric, $22)`,
		id, r.group, profile, r.apiKeyID, r.source, r.stream, started, ended, ttft, r.submissions, r.outcome, r.usage, r.input, r.read,
		reqID, r.billing, mode, r.subscription, actual, inputCost, totalCost, service.ProviderHallAlgorithmVersion)
	require.NoError(t, err)
	return id
}

type hallSnap struct {
	sampleCount, success, failed, submissions, billed int64
	fast95, p90, price, successRate                   sql.NullString
	coverage, billingCoverage, reason                 string
	computedVersion                                   int
}

func readHallSnap(t *testing.T, db *sql.DB, group, profile int64, windowEnd time.Time) (hallSnap, bool) {
	t.Helper()
	var s hallSnap
	err := db.QueryRow(`SELECT ttft_sample_count, success_count, failed_count, submissions, billed_count, ttft_fast95_mean_ms::text, ttft_p90_ms::text, historical_price::text, success_rate::text, coverage, billing_coverage, coverage_reason, computed_version
		FROM provider_hall_snapshots WHERE group_id=$1 AND profile_id=$2 AND window_end=$3 AND algorithm_version=$4`, group, profile, windowEnd, service.ProviderHallAlgorithmVersion).
		Scan(&s.sampleCount, &s.success, &s.failed, &s.submissions, &s.billed, &s.fast95, &s.p90, &s.price, &s.successRate, &s.coverage, &s.billingCoverage, &s.reason, &s.computedVersion)
	if err == sql.ErrNoRows {
		return s, false
	}
	require.NoError(t, err)
	return s, true
}

func providerHallAggregationDatabaseContracts(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	repo := NewProviderHallAggregationRepository(db)
	now := time.Now().UTC().Truncate(time.Minute)
	// Facts have no foreign keys: synthetic group/profile identities keep the
	// contract independent of the admin fixtures.
	group, profileA, profileB := int64(910_000)+now.Unix()%1000, int64(920_001), int64(920_002)
	targets := []service.ProviderHallTargetRef{
		{TargetID: 1, GroupID: group, ProfileID: profileA, Enabled: true},
		{TargetID: 2, GroupID: group, ProfileID: profileB, Enabled: true},
	}
	var epochID int64
	require.NoError(t, db.QueryRow(`INSERT INTO provider_hall_collector_epochs (node_id, algorithm_version, registered_at, heartbeat_at, confirmed_at) VALUES ('n1', $1, $2, $3, $3) RETURNING id`,
		service.ProviderHallAlgorithmVersion, now.Add(-24*time.Hour), now.Add(time.Hour)).Scan(&epochID))
	t.Cleanup(func() {
		for _, q := range []string{
			`DELETE FROM provider_hall_requests WHERE group_id=$1`,
			`DELETE FROM provider_hall_metrics_1m WHERE group_id=$1`,
			`DELETE FROM provider_hall_latency_counts_1m WHERE group_id=$1`,
			`DELETE FROM provider_hall_snapshots WHERE group_id=$1`,
			`DELETE FROM provider_hall_dirty_buckets WHERE group_id=$1`,
		} {
			_, _ = db.Exec(q, group)
		}
		_, _ = db.Exec(`DELETE FROM provider_hall_coverage_gaps WHERE node_id IN ('n1','lost-node')`)
		_, _ = db.Exec(`DELETE FROM provider_hall_collector_epochs WHERE node_id IN ('n1','lost-node')`)
		_, _ = db.Exec(`UPDATE provider_hall_aggregator_state SET last_window_end = NULL, last_error = '' WHERE id = 1`)
	})
	recompute := func(t *testing.T, T time.Time, minutes []time.Time, dirtyIDs []int64, expected []string) *service.ProviderHallRecomputeResult {
		t.Helper()
		res, err := repo.Recompute(ctx, service.ProviderHallRecomputeInput{
			Now: T.Add(45 * time.Second), T: T, Minutes: minutes, NewWindowEnd: T, DirtyIDs: dirtyIDs,
			ReevalFrom: T.Add(-65 * time.Minute), ReevalTo: T, ExpectedNodes: expected, Targets: targets,
			AlgorithmVersion: service.ProviderHallAlgorithmVersion,
		})
		require.NoError(t, err)
		return res
	}

	m := now.Add(-3 * time.Hour)
	t.Run("matrix_a_samples_to_snapshot", func(t *testing.T) {
		for ms := 1000; ms <= 19000; ms += 1000 {
			insertHallReq(t, db, hallReq{group: group, profile: profileA, endedAt: m.Add(10 * time.Second), outcome: "success", ttft: ms, stream: true, usage: true, input: 1000, read: 100, billing: "applied"})
		}
		insertHallReq(t, db, hallReq{group: group, profile: profileA, endedAt: m.Add(20 * time.Second), outcome: "success", ttft: 100000, stream: true, usage: true, input: 1000, read: 100, billing: "applied", subscription: true})
		// Non-streaming successes and excluded rows never enter the latency set.
		insertHallReq(t, db, hallReq{group: group, profile: profileA, endedAt: m.Add(30 * time.Second), outcome: "success", ttft: 5, stream: false, usage: true, input: 10})
		insertHallReq(t, db, hallReq{group: group, profile: profileA, endedAt: m.Add(30 * time.Second), outcome: "excluded", submissions: 0})
		// A retried failure counts twice in submissions and once as failed.
		insertHallReq(t, db, hallReq{group: group, profile: profileB, endedAt: m.Add(40 * time.Second), outcome: "failed", submissions: 2})
		// A different minute is untouched by this recompute.
		insertHallReq(t, db, hallReq{group: group, profile: profileB, endedAt: m.Add(2 * time.Minute), outcome: "success", stream: true, ttft: 700})

		T := m.Add(time.Minute)
		res := recompute(t, T, []time.Time{m}, nil, []string{"n1"})
		require.Equal(t, 1, res.MinutesRecomputed)
		require.Equal(t, 3, res.SnapshotsPublished, "merged row plus one per enabled profile")

		var placeholder int
		require.NoError(t, db.QueryRow(`SELECT count(*) FROM provider_hall_metrics_1m WHERE group_id=$1 AND minute=$2`, group, m).Scan(&placeholder))
		require.Equal(t, 2, placeholder, "every enabled target has a minute row")

		a, ok := readHallSnap(t, db, group, profileA, T)
		require.True(t, ok)
		require.Equal(t, int64(20), a.sampleCount)
		require.Equal(t, "10000.000", a.fast95.String)
		require.Equal(t, "18000", a.p90.String)
		require.Equal(t, int64(21), a.success)
		require.Equal(t, int64(21), a.submissions)
		require.Equal(t, int64(20), a.billed)
		require.Equal(t, "2.0000000000", a.price.String, "0.004*0.001/0.002 per bill over 1000 tokens → $2 per million")
		require.Equal(t, "1.0000000000", a.successRate.String)
		require.Equal(t, service.ProviderHallCoverageComplete, a.coverage)
		require.Equal(t, service.ProviderHallBillingCoveragePending, a.billingCoverage, "the unbilled non-stream success is still pending")
		require.Equal(t, 1, a.computedVersion)

		var share string
		require.NoError(t, db.QueryRow(`SELECT subscription_share::text FROM provider_hall_snapshots WHERE group_id=$1 AND profile_id=$2 AND window_end=$3`, group, profileA, T).Scan(&share))
		require.Equal(t, "0.0500000000", share)

		b, ok := readHallSnap(t, db, group, profileB, T)
		require.True(t, ok)
		require.Equal(t, int64(0), b.success)
		require.Equal(t, int64(1), b.failed)
		require.Equal(t, int64(2), b.submissions)
		require.Equal(t, "0.0000000000", b.successRate.String)

		merged, ok := readHallSnap(t, db, group, 0, T)
		require.True(t, ok)
		require.Equal(t, int64(20), merged.sampleCount)
		require.Equal(t, "10000.000", merged.fast95.String)
		require.Equal(t, int64(21), merged.success)
		require.Equal(t, int64(1), merged.failed)
		require.Equal(t, int64(23), merged.submissions)
		require.Equal(t, "0.9130434783", merged.successRate.String)

		_, ok = readHallSnap(t, db, group, profileB, m.Add(3*time.Minute))
		require.False(t, ok, "minutes outside M are not computed")
		var tier int
		require.NoError(t, db.QueryRow(`SELECT tier FROM provider_hall_snapshots WHERE group_id=$1 AND profile_id=0 AND window_end=$2`, group, T).Scan(&tier))
		require.Equal(t, service.ProviderHallSnapshotTier(T), tier)
		var wm time.Time
		require.NoError(t, db.QueryRow(`SELECT last_window_end FROM provider_hall_aggregator_state WHERE id=1`).Scan(&wm))
		require.Equal(t, T, wm.UTC())
	})

	t.Run("late_terminal_marks_dirty_and_recomputes_sixty_snapshots", func(t *testing.T) {
		// A finish that lands after the watermark passed: the collector marks the
		// bucket dirty in the same transaction.
		insertHallReq(t, db, hallReq{group: group, profile: profileA, endedAt: m.Add(50 * time.Second), outcome: "failed"})
		_, err := db.Exec(`INSERT INTO provider_hall_dirty_buckets (group_id, profile_id, minute, reason) VALUES ($1, $2, $3, 'collector')`, group, profileA, m)
		require.NoError(t, err)
		dirty, err := repo.ListDirty(ctx, 240)
		require.NoError(t, err)
		var ids []int64
		for _, d := range dirty {
			if d.GroupID == group {
				ids = append(ids, d.ID)
				require.Equal(t, m, d.Minute)
			}
		}
		require.Len(t, ids, 1)
		T2 := m.Add(90 * time.Minute)
		res := recompute(t, T2, []time.Time{m}, ids, []string{"n1"})
		require.Equal(t, 3*61, res.SnapshotsPublished, "60 windows containing m plus T2, three rows each")
		a, ok := readHallSnap(t, db, group, profileA, m.Add(time.Minute))
		require.True(t, ok)
		require.Equal(t, int64(1), a.failed)
		require.Equal(t, 2, a.computedVersion, "republished snapshot bumps its version")
		last, ok := readHallSnap(t, db, group, profileA, m.Add(60*time.Minute))
		require.True(t, ok)
		require.Equal(t, int64(20), last.sampleCount, "the window [m, m+60m) still contains the minute")
		_, ok = readHallSnap(t, db, group, profileA, m.Add(61*time.Minute))
		require.False(t, ok, "m+61m no longer contains m")
		var remaining int
		require.NoError(t, db.QueryRow(`SELECT count(*) FROM provider_hall_dirty_buckets WHERE group_id=$1`, group).Scan(&remaining))
		require.Zero(t, remaining)
	})

	t.Run("collection_gap_makes_window_incomplete_until_closed", func(t *testing.T) {
		m2 := now.Add(-2 * time.Hour)
		var gapID int64
		require.NoError(t, db.QueryRow(`INSERT INTO provider_hall_coverage_gaps (node_id, epoch_id, scope, started_at, reason) VALUES ('n1', $1, 'collection', $2, 'queue_overflow') RETURNING id`, epochID, m2.Add(30*time.Second)).Scan(&gapID))
		T := m2.Add(time.Minute)
		recompute(t, T, []time.Time{m2}, nil, []string{"n1"})
		s, ok := readHallSnap(t, db, group, 0, T)
		require.True(t, ok)
		require.Equal(t, service.ProviderHallCoverageCollectionGap, s.coverage)
		require.Equal(t, service.ProviderHallCoverageCollectionGap, s.reason)

		// Closing the gap exactly at the minute start removes the overlap; the
		// trailing-window re-evaluation flips the minute and republishes.
		_, err := db.Exec(`UPDATE provider_hall_coverage_gaps SET ended_at = started_at WHERE id = $1`, gapID)
		require.NoError(t, err)
		_, err = db.Exec(`UPDATE provider_hall_coverage_gaps SET started_at = $2, ended_at = $2 WHERE id = $1`, gapID, m2)
		require.NoError(t, err)
		res := recompute(t, T.Add(time.Minute), nil, nil, []string{"n1"})
		require.Equal(t, 1, res.CoverageFlipped)
		s, ok = readHallSnap(t, db, group, 0, T)
		require.True(t, ok)
		require.Equal(t, service.ProviderHallCoverageComplete, s.coverage)
		require.Equal(t, 2, s.computedVersion)
	})

	t.Run("expected_node_absent_blocks_complete", func(t *testing.T) {
		m3 := now.Add(-100 * time.Minute)
		T := m3.Add(time.Minute)
		recompute(t, T, []time.Time{m3}, nil, []string{"n1", "ghost"})
		s, ok := readHallSnap(t, db, group, profileA, T)
		require.True(t, ok)
		require.Equal(t, service.ProviderHallCoverageNodeUnconfirmed, s.coverage)
		// The wrong roster also degraded the older minute inside the trailing
		// re-evaluation window; correcting it flips both minutes back in place.
		m2 := now.Add(-2 * time.Hour)
		s2, ok := readHallSnap(t, db, group, 0, m2.Add(time.Minute))
		require.True(t, ok)
		require.Equal(t, service.ProviderHallCoverageNodeUnconfirmed, s2.coverage)
		res := recompute(t, T.Add(time.Minute), nil, nil, []string{"n1"})
		require.Equal(t, 2, res.CoverageFlipped)
		s, _ = readHallSnap(t, db, group, profileA, T)
		require.Equal(t, service.ProviderHallCoverageComplete, s.coverage)
		s2, _ = readHallSnap(t, db, group, 0, m2.Add(time.Minute))
		require.Equal(t, service.ProviderHallCoverageComplete, s2.coverage)
		// An empty roster expects every live node. Other contracts in this
		// database leave stale epochs behind; retire them so only n1 is alive.
		_, err := db.Exec(`UPDATE provider_hall_collector_epochs SET exited_at = $2, exit_reason = 'shutdown' WHERE exited_at IS NULL AND id <> $1`, epochID, now)
		require.NoError(t, err)
		m4 := now.Add(-99 * time.Minute)
		recompute(t, m4.Add(time.Minute), []time.Time{m4}, nil, nil)
		s, _ = readHallSnap(t, db, group, profileA, m4.Add(time.Minute))
		require.Equal(t, service.ProviderHallCoverageComplete, s.coverage)
	})

	t.Run("lost_epoch_becomes_gap", func(t *testing.T) {
		var lostID int64
		require.NoError(t, db.QueryRow(`INSERT INTO provider_hall_collector_epochs (node_id, algorithm_version, registered_at, heartbeat_at, confirmed_at) VALUES ('lost-node', 1, $1, $2, $2) RETURNING id`, now.Add(-time.Hour), now.Add(-2*time.Minute)).Scan(&lostID))
		n, err := repo.MarkLostEpochs(ctx, now, now.Add(-45*time.Second))
		require.NoError(t, err)
		require.Equal(t, 1, n)
		var reason string
		var exited sql.NullTime
		require.NoError(t, db.QueryRow(`SELECT exit_reason, exited_at FROM provider_hall_collector_epochs WHERE id=$1`, lostID).Scan(&reason, &exited))
		require.Equal(t, "lost", reason)
		require.True(t, exited.Valid)
		var gapReason string
		var started, ended time.Time
		require.NoError(t, db.QueryRow(`SELECT reason, started_at, ended_at FROM provider_hall_coverage_gaps WHERE epoch_id=$1`, lostID).Scan(&gapReason, &started, &ended))
		require.Equal(t, "missed_heartbeat", gapReason)
		require.Equal(t, now.Add(-2*time.Minute), started.UTC())
		require.Equal(t, now, ended.UTC())
		n, err = repo.MarkLostEpochs(ctx, now, now.Add(-45*time.Second))
		require.NoError(t, err)
		require.Zero(t, n, "idempotent")
	})

	t.Run("reconcile_reads_ledgers_only", func(t *testing.T) {
		client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
		user, err := client.User.Create().SetEmail(fmt.Sprintf("hall-agg-%d@example.test", time.Now().UnixNano())).SetPasswordHash("test-only").Save(ctx)
		require.NoError(t, err)
		key, err := client.APIKey.Create().SetUserID(user.ID).SetName("hall-agg").SetKey(fmt.Sprintf("hall-agg-%d", time.Now().UnixNano())).Save(ctx)
		require.NoError(t, err)
		account, err := client.Account.Create().SetName(fmt.Sprintf("hall-agg-%d", time.Now().UnixNano())).SetPlatform(service.PlatformOpenAI).SetType(service.AccountTypeAPIKey).SetStatus(service.StatusActive).SetCredentials(map[string]any{"api_key": "sk-test"}).Save(ctx)
		require.NoError(t, err)
		t.Cleanup(func() {
			_, _ = db.Exec(`DELETE FROM usage_logs WHERE api_key_id=$1`, key.ID)
			_, _ = db.Exec(`DELETE FROM usage_billing_dedup WHERE api_key_id=$1`, key.ID)
			_, _ = db.Exec(`DELETE FROM accounts WHERE id=$1`, account.ID)
			_, _ = db.Exec(`DELETE FROM users WHERE id=$1`, user.ID)
		})
		started := now.Add(-5 * time.Minute)
		ended := started.Add(time.Second)
		mkReq := func(reqID string, age time.Duration, status string) uuid.UUID {
			return insertHallReq(t, db, hallReq{group: group, profile: profileA, startedAt: now.Add(-age), endedAt: now.Add(-age).Add(time.Second), outcome: "success", stream: true, ttft: 900, usage: true, input: 100, requestID: reqID, apiKeyID: key.ID, billing: status})
		}
		bothID := mkReq("client:both", 5*time.Minute, "pending")
		_, err = db.Exec(`INSERT INTO usage_billing_dedup (request_id, api_key_id, request_fingerprint) VALUES ('client:both', $1, 'fp')`, key.ID)
		require.NoError(t, err)
		_, err = db.Exec(`INSERT INTO usage_logs (user_id, api_key_id, account_id, request_id, model, input_cost, total_cost, actual_cost) VALUES ($1, $2, $3, 'client:both', 'gpt-hall', 0.0015, 0.0030, 0.0060)`, user.ID, key.ID, account.ID)
		require.NoError(t, err)
		dedupOnlyID := mkReq("client:dedup", 5*time.Minute, "pending")
		_, err = db.Exec(`INSERT INTO usage_billing_dedup (request_id, api_key_id, request_fingerprint) VALUES ('client:dedup', $1, 'fp')`, key.ID)
		require.NoError(t, err)
		noneYoungID := mkReq("client:none-young", 5*time.Minute, "pending")
		noneOldID := mkReq("client:none-old", 11*time.Minute, "pending")
		noKeyOldID := mkReq("", 11*time.Minute, "pending")
		noKeyYoungID := mkReq("", 3*time.Minute, "pending")
		tooFreshID := mkReq("", time.Minute, "pending")
		staleUncertainID := mkReq("client:stale", 25*time.Hour, "uncertain")
		_, err = db.Exec(`INSERT INTO usage_billing_dedup (request_id, api_key_id, request_fingerprint) VALUES ('client:stale', $1, 'fp')`, key.ID)
		require.NoError(t, err)
		_ = ended

		res, err := repo.Reconcile(ctx, now, 500)
		require.NoError(t, err)
		require.Equal(t, 1, res.Applied)
		require.Equal(t, 1, res.Uncertain)
		// Earlier subtests left two keyless successes older than ten minutes.
		require.GreaterOrEqual(t, res.Failed, 3)
		require.Equal(t, res.Examined-2, res.Applied+res.Uncertain+res.Failed, "only the two young rows stay pending")

		status := func(id uuid.UUID) (string, sql.NullString, sql.NullString) {
			var s string
			var actual, mode sql.NullString
			require.NoError(t, db.QueryRow(`SELECT billing_status, actual_cost::text, billing_mode FROM provider_hall_requests WHERE trace_id=$1`, id).Scan(&s, &actual, &mode))
			return s, actual, mode
		}
		s, actual, mode := status(bothID)
		require.Equal(t, "applied", s)
		require.Equal(t, "0.00600000", actual.String)
		require.Equal(t, "token", mode.String)
		s, _, _ = status(dedupOnlyID)
		require.Equal(t, "uncertain", s)
		s, _, _ = status(noneYoungID)
		require.Equal(t, "pending", s)
		s, _, _ = status(noneOldID)
		require.Equal(t, "failed", s)
		s, _, _ = status(noKeyOldID)
		require.Equal(t, "failed", s)
		s, _, _ = status(noKeyYoungID)
		require.Equal(t, "pending", s)
		s, _, _ = status(tooFreshID)
		require.Equal(t, "pending", s, "younger than two minutes is not examined")
		s, _, _ = status(staleUncertainID)
		require.Equal(t, "failed", s, "a day without a usage log gives up")

		var dirtyCount int
		require.NoError(t, db.QueryRow(`SELECT count(*) FROM provider_hall_dirty_buckets WHERE group_id=$1 AND reason='reconcile' AND minute >= $2`, group, now.Add(-30*time.Minute)).Scan(&dirtyCount))
		require.Equal(t, 2, dirtyCount, "four changes in the last half hour share two minutes; the bucket key is unique")
		_, err = db.Exec(`DELETE FROM provider_hall_dirty_buckets WHERE group_id=$1`, group)
		require.NoError(t, err)
	})

	t.Run("retention_keeps_pending_reconciliation", func(t *testing.T) {
		old := now.Add(-8 * 24 * time.Hour)
		keepID := insertHallReq(t, db, hallReq{group: group, profile: profileA, startedAt: old, endedAt: old.Add(time.Second), outcome: "success", billing: "pending", requestID: "client:keep"})
		dropID := insertHallReq(t, db, hallReq{group: group, profile: profileA, startedAt: old, endedAt: old.Add(time.Second), outcome: "success", billing: "applied"})
		failedID := insertHallReq(t, db, hallReq{group: group, profile: profileA, startedAt: old, endedAt: old.Add(time.Second), outcome: "failed"})
		_, err := db.Exec(`INSERT INTO provider_hall_metrics_1m (group_id, profile_id, minute, algorithm_version) VALUES ($1, $2, $3, 1)`, group, profileA, old)
		require.NoError(t, err)
		_, err = db.Exec(`INSERT INTO provider_hall_snapshots (group_id, profile_id, window_end, tier, algorithm_version, ttft_sample_count, success_count, failed_count, submissions, usage_success_count, input_tokens, cache_read_tokens, cache_creation_tokens, billed_count, billed_input_tokens, billed_input_cost, coverage, billing_coverage)
			VALUES ($1, $2, $3, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 'complete', 'complete'), ($1, $2, $4, 5, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 'complete', 'complete'), ($1, $2, $5, 5, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 'complete', 'complete')`,
			group, profileA, now.Add(-49*time.Hour).Truncate(5*time.Minute).Add(time.Minute), now.Add(-49*time.Hour).Truncate(5*time.Minute), now.Add(-36*24*time.Hour).Truncate(5*time.Minute))
		require.NoError(t, err)
		res, err := repo.Prune(ctx, now)
		require.NoError(t, err)
		require.GreaterOrEqual(t, res.Deleted["requests"], int64(2))
		count := func(q string, args ...any) int {
			var n int
			require.NoError(t, db.QueryRow(q, args...).Scan(&n))
			return n
		}
		require.Equal(t, 1, count(`SELECT count(*) FROM provider_hall_requests WHERE trace_id=$1`, keepID))
		require.Equal(t, 0, count(`SELECT count(*) FROM provider_hall_requests WHERE trace_id=$1`, dropID))
		require.Equal(t, 0, count(`SELECT count(*) FROM provider_hall_requests WHERE trace_id=$1`, failedID))
		require.Equal(t, 0, count(`SELECT count(*) FROM provider_hall_metrics_1m WHERE group_id=$1 AND minute=$2`, group, old))
		require.Equal(t, 0, count(`SELECT count(*) FROM provider_hall_snapshots WHERE group_id=$1 AND tier=1 AND window_end < $2`, group, now.Add(-48*time.Hour)))
		require.Equal(t, 1, count(`SELECT count(*) FROM provider_hall_snapshots WHERE group_id=$1 AND tier=5 AND window_end=$2`, group, now.Add(-49*time.Hour).Truncate(5*time.Minute)), "tier 5 survives 48h")
		require.Equal(t, 0, count(`SELECT count(*) FROM provider_hall_snapshots WHERE group_id=$1 AND tier=5 AND window_end < $2`, group, now.Add(-35*24*time.Hour)))
		require.Equal(t, 3, count(`SELECT count(*) FROM provider_hall_snapshots WHERE group_id=$1 AND window_end=$2`, group, m.Add(time.Minute)), "recent snapshots untouched")
	})

	t.Run("state_round_trip", func(t *testing.T) {
		require.NoError(t, repo.RecordRun(ctx, now, "boom"))
		st, err := repo.GetState(ctx)
		require.NoError(t, err)
		require.Equal(t, "boom", st.LastError)
		require.NotNil(t, st.LastRunAt)
		require.NotNil(t, st.LastWindowEnd)
		require.NoError(t, repo.RecordRun(ctx, now, ""))
	})
}
