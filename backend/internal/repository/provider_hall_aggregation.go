package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

// providerHallAggregationRepository owns the aggregation tables with raw SQL.
// Everything runs on the primary aggregator; no request path touches it.
type providerHallAggregationRepository struct{ db *sql.DB }

func NewProviderHallAggregationRepository(db *sql.DB) service.ProviderHallAggregationRepository {
	return &providerHallAggregationRepository{db: db}
}

const providerHallPruneBatch = 5000

func (r *providerHallAggregationRepository) GetState(ctx context.Context) (*service.ProviderHallAggregatorState, error) {
	st := &service.ProviderHallAggregatorState{}
	var lastWindow, lastRun sql.NullTime
	err := r.db.QueryRowContext(ctx, `SELECT algorithm_version, last_window_end, last_run_at, last_error FROM provider_hall_aggregator_state WHERE id = 1`).
		Scan(&st.AlgorithmVersion, &lastWindow, &lastRun, &st.LastError)
	if err != nil {
		return nil, err
	}
	if lastWindow.Valid {
		t := lastWindow.Time.UTC()
		st.LastWindowEnd = &t
	}
	if lastRun.Valid {
		t := lastRun.Time.UTC()
		st.LastRunAt = &t
	}
	return st, nil
}

func (r *providerHallAggregationRepository) RecordRun(ctx context.Context, now time.Time, lastError string) error {
	if len(lastError) > 2000 {
		lastError = lastError[:2000]
	}
	_, err := r.db.ExecContext(ctx, `UPDATE provider_hall_aggregator_state SET last_run_at = $1, last_error = $2, updated_at = now() WHERE id = 1`, now, lastError)
	return err
}

// MarkLostEpochs closes epochs whose heartbeat stopped without a clean exit
// and records the unconfirmed tail as a collection gap.
func (r *providerHallAggregationRepository) MarkLostEpochs(ctx context.Context, now, staleBefore time.Time) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.QueryContext(ctx, `SELECT id, node_id, COALESCE(confirmed_at, registered_at) FROM provider_hall_collector_epochs WHERE exited_at IS NULL AND heartbeat_at < $1 FOR UPDATE SKIP LOCKED`, staleBefore)
	if err != nil {
		return 0, err
	}
	type lost struct {
		id        int64
		node      string
		confirmed time.Time
	}
	var lostEpochs []lost
	for rows.Next() {
		var l lost
		if err := rows.Scan(&l.id, &l.node, &l.confirmed); err != nil {
			_ = rows.Close()
			return 0, err
		}
		lostEpochs = append(lostEpochs, l)
	}
	_ = rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	for _, l := range lostEpochs {
		if _, err := tx.ExecContext(ctx, `UPDATE provider_hall_collector_epochs SET exited_at = $2, exit_reason = 'lost' WHERE id = $1`, l.id, now); err != nil {
			return 0, err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO provider_hall_coverage_gaps (node_id, epoch_id, scope, started_at, ended_at, reason) VALUES ($1, $2, 'collection', $3, $4, 'missed_heartbeat')`, l.node, l.id, l.confirmed, now); err != nil {
			return 0, err
		}
	}
	return len(lostEpochs), tx.Commit()
}

func (r *providerHallAggregationRepository) ListDirty(ctx context.Context, limit int) ([]service.ProviderHallDirtyBucket, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, group_id, profile_id, minute FROM provider_hall_dirty_buckets ORDER BY id LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallDirtyBucket
	for rows.Next() {
		var d service.ProviderHallDirtyBucket
		if err := rows.Scan(&d.ID, &d.GroupID, &d.ProfileID, &d.Minute); err != nil {
			return nil, err
		}
		d.Minute = d.Minute.UTC()
		out = append(out, d)
	}
	return out, rows.Err()
}

type providerHallQuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func providerHallNodeStatuses(ctx context.Context, q providerHallQuerier, expected []string) ([]service.ProviderHallNodeStatus, error) {
	if len(expected) == 0 {
		// No roster configured: every node that is currently alive is expected.
		rows, err := q.QueryContext(ctx, `SELECT DISTINCT node_id FROM provider_hall_collector_epochs WHERE exited_at IS NULL`)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var node string
			if err := rows.Scan(&node); err != nil {
				_ = rows.Close()
				return nil, err
			}
			expected = append(expected, node)
		}
		_ = rows.Close()
		if err := rows.Err(); err != nil {
			return nil, err
		}
		sort.Strings(expected)
	}
	out := make([]service.ProviderHallNodeStatus, 0, len(expected))
	for _, node := range expected {
		st := service.ProviderHallNodeStatus{NodeID: node}
		var confirmed, exited sql.NullTime
		err := q.QueryRowContext(ctx, `SELECT id, algorithm_version, confirmed_at, exited_at FROM provider_hall_collector_epochs WHERE node_id = $1 ORDER BY registered_at DESC, id DESC LIMIT 1`, node).
			Scan(&st.EpochID, &st.AlgorithmVersion, &confirmed, &exited)
		switch {
		case err == sql.ErrNoRows:
		case err != nil:
			return nil, err
		default:
			st.Present = true
			if confirmed.Valid {
				t := confirmed.Time.UTC()
				st.ConfirmedAt = &t
			}
			st.Exited = exited.Valid
		}
		out = append(out, st)
	}
	return out, nil
}

func providerHallLoadGaps(ctx context.Context, q providerHallQuerier, lo, hi time.Time) ([]service.ProviderHallGapRange, error) {
	rows, err := q.QueryContext(ctx, `SELECT scope, started_at, ended_at FROM provider_hall_coverage_gaps WHERE started_at < $2 AND (ended_at IS NULL OR ended_at > $1)`, lo, hi)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallGapRange
	for rows.Next() {
		var g service.ProviderHallGapRange
		var ended sql.NullTime
		if err := rows.Scan(&g.Scope, &g.StartedAt, &ended); err != nil {
			return nil, err
		}
		g.StartedAt = g.StartedAt.UTC()
		if ended.Valid {
			t := ended.Time.UTC()
			g.EndedAt = &t
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

const providerHallMinuteMetricsSQL = `
INSERT INTO provider_hall_metrics_1m (
    group_id, profile_id, minute, algorithm_version,
    success_count, failed_count, submissions, excluded_count, ttft_sample_count, usage_success_count,
    input_tokens, cache_read_tokens, cache_creation_tokens, output_tokens,
    billed_count, billed_subscription_count, billed_input_tokens, billed_input_cost,
    billing_pending_count, billing_uncertain_count, computed_at)
SELECT group_id, profile_id, date_trunc('minute', ended_at) AS minute, $2::smallint,
    count(*) FILTER (WHERE outcome = 'success'),
    count(*) FILTER (WHERE outcome = 'failed'),
    COALESCE(sum(submissions), 0),
    count(*) FILTER (WHERE outcome = 'excluded'),
    count(*) FILTER (WHERE outcome = 'success' AND stream AND ttft_ms IS NOT NULL),
    count(*) FILTER (WHERE outcome = 'success' AND usage_known),
    COALESCE(sum(input_tokens) FILTER (WHERE outcome = 'success' AND usage_known), 0),
    COALESCE(sum(cache_read_tokens) FILTER (WHERE outcome = 'success' AND usage_known), 0),
    COALESCE(sum(cache_creation_tokens) FILTER (WHERE outcome = 'success' AND usage_known), 0),
    COALESCE(sum(output_tokens) FILTER (WHERE outcome = 'success' AND usage_known), 0),
    count(*) FILTER (WHERE billed),
    count(*) FILTER (WHERE billed AND is_subscription),
    COALESCE(sum(input_tokens) FILTER (WHERE billed), 0),
    COALESCE(sum(CASE WHEN total_base_cost > 0 AND input_base_cost <= total_base_cost
                      THEN actual_cost * input_base_cost / total_base_cost ELSE 0 END) FILTER (WHERE billed), 0),
    count(*) FILTER (WHERE outcome = 'success' AND billing_status = 'pending'),
    count(*) FILTER (WHERE outcome = 'success' AND billing_status = 'uncertain'),
    $4
FROM (
    SELECT r.*, (outcome = 'success' AND billing_status = 'applied' AND billing_mode = 'token'
                 AND usage_known AND actual_cost IS NOT NULL AND actual_cost >= 0) AS billed
    FROM provider_hall_requests r
    WHERE source = 'user' AND profile_id IS NOT NULL AND ended_at IS NOT NULL AND algorithm_version = $2
      AND date_trunc('minute', ended_at) = ANY($1::timestamptz[])
      AND started_at >= $3 AND started_at < $5
) f
GROUP BY group_id, profile_id, date_trunc('minute', ended_at)`

const providerHallMinuteLatencySQL = `
INSERT INTO provider_hall_latency_counts_1m (group_id, profile_id, minute, algorithm_version, ttft_ms, count)
SELECT group_id, profile_id, date_trunc('minute', ended_at), $2::smallint, ttft_ms, count(*)
FROM provider_hall_requests
WHERE source = 'user' AND profile_id IS NOT NULL AND ended_at IS NOT NULL AND algorithm_version = $2
  AND outcome = 'success' AND stream AND ttft_ms IS NOT NULL AND ttft_ms > 0
  AND date_trunc('minute', ended_at) = ANY($1::timestamptz[])
  AND started_at >= $3 AND started_at < $4
GROUP BY group_id, profile_id, date_trunc('minute', ended_at), ttft_ms`

// Recompute performs plan steps 5–7 atomically: rebuild the requested
// minutes, re-evaluate coverage of the trailing window, publish every
// dependent snapshot, consume dirty marks and advance the watermark.
func (r *providerHallAggregationRepository) Recompute(ctx context.Context, in service.ProviderHallRecomputeInput) (result *service.ProviderHallRecomputeResult, err error) {
	result = &service.ProviderHallRecomputeResult{}
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	minutes := make([]time.Time, 0, len(in.Minutes))
	for _, m := range in.Minutes {
		minutes = append(minutes, m.UTC().Truncate(time.Minute))
	}
	sort.Slice(minutes, func(i, j int) bool { return minutes[i].Before(minutes[j]) })
	lo, hi := in.ReevalFrom.UTC(), in.T.UTC().Add(time.Minute)
	if len(minutes) > 0 && minutes[0].Before(lo) {
		lo = minutes[0]
	}
	if len(minutes) > 0 && minutes[len(minutes)-1].Add(time.Minute).After(hi) {
		hi = minutes[len(minutes)-1].Add(time.Minute)
	}
	nodes, err := providerHallNodeStatuses(ctx, tx, in.ExpectedNodes)
	if err != nil {
		return nil, err
	}
	gaps, err := providerHallLoadGaps(ctx, tx, lo, hi)
	if err != nil {
		return nil, err
	}
	if len(minutes) > 0 {
		if err = r.recomputeMinutes(ctx, tx, in, minutes, nodes, gaps); err != nil {
			return nil, fmt.Errorf("recompute minutes: %w", err)
		}
		result.MinutesRecomputed = len(minutes)
	}
	flipped, err := r.reevalCoverage(ctx, tx, in, minutes, nodes, gaps)
	if err != nil {
		return nil, fmt.Errorf("reevaluate coverage: %w", err)
	}
	result.CoverageFlipped = len(flipped)
	times := service.ProviderHallSnapshotTimes(in.T, append(append([]time.Time{}, minutes...), flipped...), in.Now)
	published, err := r.publishSnapshots(ctx, tx, in, times)
	if err != nil {
		return nil, fmt.Errorf("publish snapshots: %w", err)
	}
	result.SnapshotsPublished = published
	if len(in.DirtyIDs) > 0 {
		if _, err = tx.ExecContext(ctx, `DELETE FROM provider_hall_dirty_buckets WHERE id = ANY($1::bigint[])`, pq.Array(in.DirtyIDs)); err != nil {
			return nil, err
		}
	}
	if _, err = tx.ExecContext(ctx, `UPDATE provider_hall_aggregator_state SET last_window_end = $1, algorithm_version = $2, updated_at = now() WHERE id = 1`, in.NewWindowEnd.UTC(), in.AlgorithmVersion); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *providerHallAggregationRepository) recomputeMinutes(ctx context.Context, tx *sql.Tx, in service.ProviderHallRecomputeInput, minutes []time.Time, nodes []service.ProviderHallNodeStatus, gaps []service.ProviderHallGapRange) error {
	arr := pq.Array(minutes)
	for _, table := range []string{"provider_hall_latency_counts_1m", "provider_hall_metrics_1m"} {
		if _, err := tx.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s WHERE algorithm_version = $2 AND minute = ANY($1::timestamptz[])`, table), arr, in.AlgorithmVersion); err != nil {
			return err
		}
	}
	// A request can only end in minute m if it started within the last two
	// hours of it; the bound lets the started_at index prune the scan.
	startLo, startHi := minutes[0].Add(-2*time.Hour), minutes[len(minutes)-1].Add(time.Minute)
	if _, err := tx.ExecContext(ctx, providerHallMinuteMetricsSQL, arr, in.AlgorithmVersion, startLo, in.Now, startHi); err != nil {
		return fmt.Errorf("metrics: %w", err)
	}
	if _, err := tx.ExecContext(ctx, providerHallMinuteLatencySQL, arr, in.AlgorithmVersion, startLo, startHi); err != nil {
		return fmt.Errorf("latency: %w", err)
	}
	// Every enabled target gets a row per minute so coverage is recorded even
	// when no request landed.
	if len(in.Targets) > 0 {
		groups, profiles := make([]int64, 0, len(in.Targets)), make([]int64, 0, len(in.Targets))
		for _, t := range in.Targets {
			if t.Enabled {
				groups, profiles = append(groups, t.GroupID), append(profiles, t.ProfileID)
			}
		}
		if len(groups) > 0 {
			if _, err := tx.ExecContext(ctx, `
INSERT INTO provider_hall_metrics_1m (group_id, profile_id, minute, algorithm_version, computed_at)
SELECT t.g, t.p, m.minute, $3::smallint, $4
FROM unnest($1::bigint[], $2::bigint[]) AS t(g, p)
CROSS JOIN unnest($5::timestamptz[]) AS m(minute)
ON CONFLICT (group_id, profile_id, minute, algorithm_version) DO NOTHING`, pq.Array(groups), pq.Array(profiles), in.AlgorithmVersion, in.Now, arr); err != nil {
				return fmt.Errorf("placeholders: %w", err)
			}
		}
	}
	coverages := make([]string, 0, len(minutes))
	billingGaps := make([]bool, 0, len(minutes))
	for _, m := range minutes {
		coverages = append(coverages, service.ProviderHallMinuteCoverage(m, gaps, nodes, in.AlgorithmVersion))
		billingGaps = append(billingGaps, service.ProviderHallMinuteBillingCoverage(m, gaps, 0, 0) == service.ProviderHallBillingCoverageGap)
	}
	if _, err := tx.ExecContext(ctx, `
UPDATE provider_hall_metrics_1m AS x SET
    coverage = c.coverage,
    billing_coverage = CASE WHEN c.billing_gap THEN 'billing_gap'
                            WHEN x.billing_pending_count + x.billing_uncertain_count > 0 THEN 'pending'
                            ELSE 'complete' END
FROM unnest($1::timestamptz[], $2::text[], $3::boolean[]) AS c(minute, coverage, billing_gap)
WHERE x.minute = c.minute AND x.algorithm_version = $4`, arr, pq.Array(coverages), pq.Array(billingGaps), in.AlgorithmVersion); err != nil {
		return fmt.Errorf("coverage: %w", err)
	}
	return nil
}

// reevalCoverage re-checks coverage of the trailing window (minutes not just
// recomputed). It returns every minute whose coverage changed so the
// dependent snapshots are republished in both directions: a late node
// confirmation completes them, a roster correction can also degrade them.
func (r *providerHallAggregationRepository) reevalCoverage(ctx context.Context, tx *sql.Tx, in service.ProviderHallRecomputeInput, recomputed []time.Time, nodes []service.ProviderHallNodeStatus, gaps []service.ProviderHallGapRange) ([]time.Time, error) {
	rows, err := tx.QueryContext(ctx, `SELECT minute, coverage FROM provider_hall_metrics_1m WHERE algorithm_version = $1 AND minute >= $2 AND minute < $3 AND NOT (minute = ANY($4::timestamptz[])) GROUP BY minute, coverage`, in.AlgorithmVersion, in.ReevalFrom.UTC(), in.ReevalTo.UTC(), pq.Array(recomputed))
	if err != nil {
		return nil, err
	}
	type pair struct {
		minute   time.Time
		coverage string
	}
	var pairs []pair
	for rows.Next() {
		var p pair
		if err := rows.Scan(&p.minute, &p.coverage); err != nil {
			_ = rows.Close()
			return nil, err
		}
		p.minute = p.minute.UTC()
		pairs = append(pairs, p)
	}
	_ = rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var flipped []time.Time
	seen := map[time.Time]bool{}
	for _, p := range pairs {
		want := service.ProviderHallMinuteCoverage(p.minute, gaps, nodes, in.AlgorithmVersion)
		if want == p.coverage {
			continue
		}
		if _, err := tx.ExecContext(ctx, `UPDATE provider_hall_metrics_1m SET coverage = $3 WHERE algorithm_version = $1 AND minute = $2 AND coverage = $4`, in.AlgorithmVersion, p.minute, want, p.coverage); err != nil {
			return nil, err
		}
		if !seen[p.minute] {
			seen[p.minute] = true
			flipped = append(flipped, p.minute)
		}
	}
	return flipped, nil
}

type providerHallLatencyRow struct {
	profileID int64
	minute    time.Time
	ttftMs    int64
	count     int64
}

func (r *providerHallAggregationRepository) publishSnapshots(ctx context.Context, tx *sql.Tx, in service.ProviderHallRecomputeInput, times []time.Time) (int, error) {
	if len(times) == 0 || len(in.Targets) == 0 {
		return 0, nil
	}
	byGroup := map[int64][]int64{}
	for _, t := range in.Targets {
		if t.Enabled {
			byGroup[t.GroupID] = append(byGroup[t.GroupID], t.ProfileID)
		}
	}
	groupIDs := make([]int64, 0, len(byGroup))
	for g := range byGroup {
		groupIDs = append(groupIDs, g)
	}
	sort.Slice(groupIDs, func(i, j int) bool { return groupIDs[i] < groupIDs[j] })
	rangeLo, rangeHi := times[0].Add(-service.ProviderHallSnapshotWindow), times[len(times)-1]
	tier1Floor := in.Now.Add(-service.ProviderHallRetentionTier1)
	published := 0
	for _, g := range groupIDs {
		profiles := byGroup[g]
		sort.Slice(profiles, func(i, j int) bool { return profiles[i] < profiles[j] })
		metricRows, err := r.loadMinuteRows(ctx, tx, g, profiles, rangeLo, rangeHi, in.AlgorithmVersion)
		if err != nil {
			return published, err
		}
		latencyRows, err := r.loadLatencyRows(ctx, tx, g, profiles, rangeLo, rangeHi, in.AlgorithmVersion)
		if err != nil {
			return published, err
		}
		for _, t := range times {
			if service.ProviderHallSnapshotTier(t) == 1 && t.Before(tier1Floor) {
				continue
			}
			wlo := t.Add(-service.ProviderHallSnapshotWindow)
			inWindow := func(m time.Time) bool { return !m.Before(wlo) && m.Before(t) }
			for _, p := range append([]int64{0}, profiles...) {
				var minutes []service.ProviderHallMinuteRow
				latency := map[int64]int64{}
				for _, row := range metricRows {
					if (p == 0 || row.ProfileID == p) && inWindow(row.Minute) {
						minutes = append(minutes, row)
					}
				}
				for _, row := range latencyRows {
					if (p == 0 || row.profileID == p) && inWindow(row.minute) {
						latency[row.ttftMs] += row.count
					}
				}
				snap := service.ProviderHallBuildSnapshot(g, p, t, minutes, latency, in.AlgorithmVersion, in.Now)
				if err := providerHallUpsertSnapshot(ctx, tx, snap); err != nil {
					return published, err
				}
				published++
			}
		}
	}
	return published, nil
}

func (r *providerHallAggregationRepository) loadMinuteRows(ctx context.Context, q providerHallQuerier, groupID int64, profiles []int64, lo, hi time.Time, algo int) ([]service.ProviderHallMinuteRow, error) {
	rows, err := q.QueryContext(ctx, `
SELECT group_id, profile_id, minute, success_count, failed_count, submissions, excluded_count, ttft_sample_count, usage_success_count,
       input_tokens, cache_read_tokens, cache_creation_tokens, output_tokens,
       billed_count, billed_subscription_count, billed_input_tokens, billed_input_cost::text,
       billing_pending_count, billing_uncertain_count, coverage, billing_coverage
FROM provider_hall_metrics_1m
WHERE group_id = $1 AND profile_id = ANY($2::bigint[]) AND algorithm_version = $3 AND minute >= $4 AND minute < $5`, groupID, pq.Array(profiles), algo, lo, hi)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallMinuteRow
	for rows.Next() {
		var m service.ProviderHallMinuteRow
		var cost string
		if err := rows.Scan(&m.GroupID, &m.ProfileID, &m.Minute, &m.SuccessCount, &m.FailedCount, &m.Submissions, &m.ExcludedCount, &m.TTFTSampleCount, &m.UsageSuccessCount,
			&m.InputTokens, &m.CacheReadTokens, &m.CacheCreationTokens, &m.OutputTokens,
			&m.BilledCount, &m.BilledSubscriptionCount, &m.BilledInputTokens, &cost,
			&m.BillingPendingCount, &m.BillingUncertainCount, &m.Coverage, &m.BillingCoverage); err != nil {
			return nil, err
		}
		m.Minute = m.Minute.UTC()
		if m.BilledInputCost, err = decimal.NewFromString(cost); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *providerHallAggregationRepository) loadLatencyRows(ctx context.Context, q providerHallQuerier, groupID int64, profiles []int64, lo, hi time.Time, algo int) ([]providerHallLatencyRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT profile_id, minute, ttft_ms, count FROM provider_hall_latency_counts_1m WHERE group_id = $1 AND profile_id = ANY($2::bigint[]) AND algorithm_version = $3 AND minute >= $4 AND minute < $5`, groupID, pq.Array(profiles), algo, lo, hi)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []providerHallLatencyRow
	for rows.Next() {
		var l providerHallLatencyRow
		if err := rows.Scan(&l.profileID, &l.minute, &l.ttftMs, &l.count); err != nil {
			return nil, err
		}
		l.minute = l.minute.UTC()
		out = append(out, l)
	}
	return out, rows.Err()
}

func providerHallNullDecimal(d *decimal.Decimal, places int32) any {
	if d == nil {
		return nil
	}
	return d.StringFixed(places)
}

func providerHallUpsertSnapshot(ctx context.Context, q providerHallQuerier, s service.ProviderHallSnapshotRow) error {
	var p90 any
	if s.TTFTP90Ms != nil {
		p90 = *s.TTFTP90Ms
	}
	_, err := q.ExecContext(ctx, `
INSERT INTO provider_hall_snapshots (
    group_id, profile_id, window_end, tier, algorithm_version,
    ttft_sample_count, ttft_fast95_mean_ms, ttft_p90_ms,
    success_count, failed_count, submissions, success_rate,
    usage_success_count, input_tokens, cache_read_tokens, cache_creation_tokens, cache_rate,
    billed_count, billed_input_tokens, billed_input_cost, historical_price, subscription_share,
    coverage, billing_coverage, coverage_reason, computed_at, computed_version)
VALUES ($1, $2, $3, $4, $5, $6, $7::numeric, $8, $9, $10, $11, $12::numeric, $13, $14, $15, $16, $17::numeric,
        $18, $19, $20::numeric, $21::numeric, $22::numeric, $23, $24, $25, $26, 1)
ON CONFLICT (group_id, profile_id, window_end, algorithm_version) DO UPDATE SET
    tier = EXCLUDED.tier,
    ttft_sample_count = EXCLUDED.ttft_sample_count, ttft_fast95_mean_ms = EXCLUDED.ttft_fast95_mean_ms, ttft_p90_ms = EXCLUDED.ttft_p90_ms,
    success_count = EXCLUDED.success_count, failed_count = EXCLUDED.failed_count, submissions = EXCLUDED.submissions, success_rate = EXCLUDED.success_rate,
    usage_success_count = EXCLUDED.usage_success_count, input_tokens = EXCLUDED.input_tokens, cache_read_tokens = EXCLUDED.cache_read_tokens,
    cache_creation_tokens = EXCLUDED.cache_creation_tokens, cache_rate = EXCLUDED.cache_rate,
    billed_count = EXCLUDED.billed_count, billed_input_tokens = EXCLUDED.billed_input_tokens, billed_input_cost = EXCLUDED.billed_input_cost,
    historical_price = EXCLUDED.historical_price, subscription_share = EXCLUDED.subscription_share,
    coverage = EXCLUDED.coverage, billing_coverage = EXCLUDED.billing_coverage, coverage_reason = EXCLUDED.coverage_reason,
    computed_at = EXCLUDED.computed_at, computed_version = provider_hall_snapshots.computed_version + 1`,
		s.GroupID, s.ProfileID, s.WindowEnd.UTC(), s.Tier, s.AlgorithmVersion,
		s.TTFTSampleCount, providerHallNullDecimal(s.TTFTFast95MeanMs, 3), p90,
		s.SuccessCount, s.FailedCount, s.Submissions, providerHallNullDecimal(s.SuccessRate, 10),
		s.UsageSuccessCount, s.InputTokens, s.CacheReadTokens, s.CacheCreationTokens, providerHallNullDecimal(s.CacheRate, 10),
		s.BilledCount, s.BilledInputTokens, s.BilledInputCost.StringFixed(10), providerHallNullDecimal(s.HistoricalPrice, 10), providerHallNullDecimal(s.SubscriptionShare, 10),
		s.Coverage, s.BillingCoverage, s.CoverageReason, s.ComputedAt.UTC())
	return err
}

// Reconcile settles successful requests whose bill has not reported back. It
// only reads the billing ledgers; it never retries a charge.
func (r *providerHallAggregationRepository) Reconcile(ctx context.Context, now time.Time, limit int) (result *service.ProviderHallReconcileResult, err error) {
	result = &service.ProviderHallReconcileResult{}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	rows, err := tx.QueryContext(ctx, `
SELECT trace_id, source, group_id, profile_id, ended_at, started_at, COALESCE(billing_request_id, ''), billing_api_key_id, api_key_id, billing_status
FROM provider_hall_requests
WHERE outcome = 'success' AND billing_status IN ('pending', 'uncertain') AND started_at > $1 AND started_at < $2
ORDER BY started_at LIMIT $3 FOR UPDATE SKIP LOCKED`, now.Add(-service.ProviderHallRetentionRequests), now.Add(-service.ProviderHallReconcileMinAge), limit)
	if err != nil {
		return nil, err
	}
	type cand struct {
		traceID, source, requestID, status string
		groupID, apiKeyID                  int64
		profileID, billingKeyID            sql.NullInt64
		endedAt                            sql.NullTime
		startedAt                          time.Time
	}
	var cands []cand
	for rows.Next() {
		var c cand
		if err = rows.Scan(&c.traceID, &c.source, &c.groupID, &c.profileID, &c.endedAt, &c.startedAt, &c.requestID, &c.billingKeyID, &c.apiKeyID, &c.status); err != nil {
			_ = rows.Close()
			return nil, err
		}
		cands = append(cands, c)
	}
	_ = rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	result.Examined = len(cands)
	for _, c := range cands {
		age := now.Sub(c.startedAt)
		next := ""
		var actual, input, total sql.NullString
		if c.requestID != "" {
			keyID := c.apiKeyID
			if c.billingKeyID.Valid {
				keyID = c.billingKeyID.Int64
			}
			var dedup bool
			if err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM usage_billing_dedup WHERE request_id = $1 AND api_key_id = $2)`, c.requestID, keyID).Scan(&dedup); err != nil {
				return nil, err
			}
			logErr := tx.QueryRowContext(ctx, `SELECT actual_cost::text, input_cost::text, total_cost::text FROM usage_logs WHERE request_id = $1 AND api_key_id = $2 ORDER BY created_at DESC LIMIT 1`, c.requestID, keyID).Scan(&actual, &input, &total)
			hasLog := logErr == nil
			if logErr != nil && logErr != sql.ErrNoRows {
				err = logErr
				return nil, err
			}
			switch {
			case dedup && hasLog:
				next = "applied"
			case dedup && age > service.ProviderHallReconcileGiveUp:
				next = "failed"
			case dedup:
				next = "uncertain"
			case age > service.ProviderHallReconcileFailAfter:
				next = "failed"
			}
		} else if age > service.ProviderHallReconcileFailAfter {
			next = "failed"
		}
		if next == "" || next == c.status {
			continue
		}
		if next == "applied" {
			_, err = tx.ExecContext(ctx, `UPDATE provider_hall_requests SET billing_status = 'applied', actual_cost = $2::numeric, input_base_cost = $3::numeric, total_base_cost = $4::numeric,
				billing_mode = COALESCE(billing_mode, 'token'), billed_at = COALESCE(billed_at, $5), updated_at = now() WHERE trace_id = $1::uuid`,
				c.traceID, actual.String, input.String, total.String, now)
			result.Applied++
		} else {
			_, err = tx.ExecContext(ctx, `UPDATE provider_hall_requests SET billing_status = $2, updated_at = now() WHERE trace_id = $1::uuid`, c.traceID, next)
			if next == "failed" {
				result.Failed++
			} else {
				result.Uncertain++
			}
		}
		if err != nil {
			return nil, err
		}
		if c.source == "user" && c.profileID.Valid && c.endedAt.Valid {
			if _, err = tx.ExecContext(ctx, `INSERT INTO provider_hall_dirty_buckets (group_id, profile_id, minute, reason) VALUES ($1, $2, date_trunc('minute', $3::timestamptz), 'reconcile') ON CONFLICT (group_id, profile_id, minute) DO NOTHING`, c.groupID, c.profileID.Int64, c.endedAt.Time); err != nil {
				return nil, err
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func providerHallDeleteBatches(ctx context.Context, db *sql.DB, query string, args ...any) (int64, error) {
	var total int64
	for {
		res, err := db.ExecContext(ctx, query, args...)
		if err != nil {
			return total, err
		}
		n, _ := res.RowsAffected()
		total += n
		if n < providerHallPruneBatch {
			return total, nil
		}
		if ctx.Err() != nil {
			return total, ctx.Err()
		}
	}
}

// Prune applies the retention rules in bounded batches. Requests awaiting
// reconciliation and samples of unfinished jobs are never deleted.
func (r *providerHallAggregationRepository) Prune(ctx context.Context, now time.Time) (*service.ProviderHallPruneResult, error) {
	res := &service.ProviderHallPruneResult{Deleted: map[string]int64{}}
	var hasSamples bool
	if err := r.db.QueryRowContext(ctx, `SELECT to_regclass('provider_hall_samples') IS NOT NULL`).Scan(&hasSamples); err != nil {
		return nil, err
	}
	sampleGuard := ""
	if hasSamples {
		sampleGuard = ` AND (sample_id IS NULL OR NOT EXISTS (SELECT 1 FROM provider_hall_samples s WHERE s.id = provider_hall_requests.sample_id AND s.status IN ('dispatched', 'uncertain')))`
	}
	type rule struct {
		name  string
		query string
		args  []any
	}
	rules := []rule{
		{"requests", `DELETE FROM provider_hall_requests WHERE trace_id IN (SELECT trace_id FROM provider_hall_requests WHERE started_at < $1 AND NOT (outcome = 'success' AND billing_status IN ('pending', 'uncertain'))` + sampleGuard + ` LIMIT $2)`, []any{now.Add(-service.ProviderHallRetentionRequests), providerHallPruneBatch}},
		{"metrics_1m", `DELETE FROM provider_hall_metrics_1m WHERE ctid = ANY(ARRAY(SELECT ctid FROM provider_hall_metrics_1m WHERE minute < $1 LIMIT $2))`, []any{now.Add(-service.ProviderHallRetentionMetrics), providerHallPruneBatch}},
		{"latency_counts_1m", `DELETE FROM provider_hall_latency_counts_1m WHERE ctid = ANY(ARRAY(SELECT ctid FROM provider_hall_latency_counts_1m WHERE minute < $1 LIMIT $2))`, []any{now.Add(-service.ProviderHallRetentionMetrics), providerHallPruneBatch}},
		{"snapshots_tier1", `DELETE FROM provider_hall_snapshots WHERE id IN (SELECT id FROM provider_hall_snapshots WHERE tier = 1 AND window_end < $1 ORDER BY id LIMIT $2)`, []any{now.Add(-service.ProviderHallRetentionTier1), providerHallPruneBatch}},
		{"snapshots_tier5", `DELETE FROM provider_hall_snapshots WHERE id IN (SELECT id FROM provider_hall_snapshots WHERE tier = 5 AND window_end < $1 ORDER BY id LIMIT $2)`, []any{now.Add(-service.ProviderHallRetentionTier5), providerHallPruneBatch}},
		{"dirty_buckets", `DELETE FROM provider_hall_dirty_buckets WHERE id IN (SELECT id FROM provider_hall_dirty_buckets WHERE created_at < $1 ORDER BY id LIMIT $2)`, []any{now.Add(-service.ProviderHallRetentionDirty), providerHallPruneBatch}},
		{"coverage_gaps", `DELETE FROM provider_hall_coverage_gaps WHERE id IN (SELECT id FROM provider_hall_coverage_gaps WHERE ended_at IS NOT NULL AND ended_at < $1 ORDER BY id LIMIT $2)`, []any{now.Add(-service.ProviderHallRetentionGaps), providerHallPruneBatch}},
		{"collector_epochs", `DELETE FROM provider_hall_collector_epochs WHERE id IN (SELECT id FROM provider_hall_collector_epochs WHERE exited_at IS NOT NULL AND exited_at < $1 ORDER BY id LIMIT $2)`, []any{now.Add(-service.ProviderHallRetentionEpochs), providerHallPruneBatch}},
	}
	if hasSamples {
		rules = append(rules,
			rule{"jobs", `DELETE FROM provider_hall_jobs WHERE id IN (SELECT id FROM provider_hall_jobs j WHERE j.created_at < $1 AND NOT EXISTS (SELECT 1 FROM provider_hall_samples s WHERE s.job_id = j.id AND s.status IN ('dispatched', 'uncertain')) ORDER BY id LIMIT $2)`, []any{now.Add(-service.ProviderHallRetentionJobs), providerHallPruneBatch}},
			rule{"spend", `DELETE FROM provider_hall_spend WHERE ctid = ANY(ARRAY(SELECT ctid FROM provider_hall_spend WHERE created_at < $1 LIMIT $2))`, []any{now.Add(-service.ProviderHallRetentionJobs), providerHallPruneBatch}},
		)
	}
	for _, rule := range rules {
		n, err := providerHallDeleteBatches(ctx, r.db, rule.query, rule.args...)
		res.Deleted[rule.name] = n
		if err != nil {
			return res, fmt.Errorf("prune %s: %w", rule.name, err)
		}
	}
	return res, nil
}
