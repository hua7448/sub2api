package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

// providerHallReadRepository serves the user hall. Every method is one
// statement over the aggregate tables; usage_logs and provider_hall_requests
// are never touched here, so the list cost does not grow with traffic.
type providerHallReadRepository struct{ db *sql.DB }

func NewProviderHallReadRepository(db *sql.DB) service.ProviderHallReadRepository {
	return &providerHallReadRepository{db: db}
}

func (r *providerHallReadRepository) ListListedGroups(ctx context.Context) ([]service.ProviderHallListedGroup, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT h.group_id, h.version, h.listed, h.display_name, h.description, h.display_order, h.updated_at, g.name
FROM provider_hall_groups h JOIN groups g ON g.id = h.group_id AND g.deleted_at IS NULL
WHERE h.listed
ORDER BY h.display_order, h.group_id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []service.ProviderHallListedGroup{}
	for rows.Next() {
		var g service.ProviderHallListedGroup
		if err := rows.Scan(&g.GroupID, &g.Version, &g.Listed, &g.DisplayName, &g.Description, &g.DisplayOrder, &g.UpdatedAt, &g.GroupName); err != nil {
			return nil, err
		}
		g.UpdatedAt = g.UpdatedAt.UTC()
		out = append(out, g)
	}
	return out, rows.Err()
}

func providerHallScanDecimal(v sql.NullString) *decimal.Decimal {
	if !v.Valid {
		return nil
	}
	d, err := decimal.NewFromString(v.String)
	if err != nil {
		return nil
	}
	return &d
}

func (r *providerHallReadRepository) ListEnabledTargets(ctx context.Context) ([]service.ProviderHallReadTarget, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT t.id, t.group_id, t.profile_id, p.model, p.protocol, p.supports_tools, t.probe_interval_seconds,
       p.reference_input_price::text, p.reference_cache_price::text, p.reference_cache_rate::text, p.reference_confirmed_at
FROM provider_hall_targets t
JOIN provider_hall_profiles p ON p.id = t.profile_id
JOIN provider_hall_groups h ON h.group_id = t.group_id AND h.listed
WHERE t.enabled
ORDER BY t.group_id, t.profile_id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []service.ProviderHallReadTarget{}
	for rows.Next() {
		var t service.ProviderHallReadTarget
		var in, cache, rate sql.NullString
		var confirmed sql.NullTime
		if err := rows.Scan(&t.TargetID, &t.GroupID, &t.ProfileID, &t.Model, &t.Protocol, &t.SupportsTools, &t.ProbeIntervalSeconds, &in, &cache, &rate, &confirmed); err != nil {
			return nil, err
		}
		t.Reference = service.ProviderHallReference{InputPrice: providerHallScanDecimal(in), CachePrice: providerHallScanDecimal(cache), CacheRate: providerHallScanDecimal(rate)}
		if confirmed.Valid {
			at := confirmed.Time.UTC()
			t.Reference.ConfirmedAt = &at
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *providerHallReadRepository) LatestWindow(ctx context.Context) (*time.Time, int, error) {
	var last sql.NullTime
	var algo int
	err := r.db.QueryRowContext(ctx, `SELECT last_window_end, algorithm_version FROM provider_hall_aggregator_state WHERE id = 1`).Scan(&last, &algo)
	if err == sql.ErrNoRows {
		return nil, service.ProviderHallAlgorithmVersion, nil
	}
	if err != nil {
		return nil, 0, err
	}
	if !last.Valid {
		return nil, algo, nil
	}
	t := last.Time.UTC()
	return &t, algo, nil
}

const providerHallSnapshotColumns = `id, group_id, profile_id, window_end, tier, algorithm_version,
       ttft_sample_count, ttft_fast95_mean_ms::text, ttft_p90_ms,
       success_count, failed_count, submissions, success_rate::text,
       usage_success_count, input_tokens, cache_read_tokens, cache_creation_tokens, cache_rate::text,
       billed_count, billed_input_tokens, billed_input_cost::text, historical_price::text, subscription_share::text,
       coverage, billing_coverage, coverage_reason, computed_at, computed_version`

func providerHallScanSnapshot(rows *sql.Rows) (service.ProviderHallSnapshotRow, error) {
	var s service.ProviderHallSnapshotRow
	var fast95, successRate, cacheRate, billedCost, price, share sql.NullString
	var p90 sql.NullInt64
	if err := rows.Scan(&s.ID, &s.GroupID, &s.ProfileID, &s.WindowEnd, &s.Tier, &s.AlgorithmVersion,
		&s.TTFTSampleCount, &fast95, &p90,
		&s.SuccessCount, &s.FailedCount, &s.Submissions, &successRate,
		&s.UsageSuccessCount, &s.InputTokens, &s.CacheReadTokens, &s.CacheCreationTokens, &cacheRate,
		&s.BilledCount, &s.BilledInputTokens, &billedCost, &price, &share,
		&s.Coverage, &s.BillingCoverage, &s.CoverageReason, &s.ComputedAt, &s.ComputedVersion); err != nil {
		return s, err
	}
	s.WindowEnd, s.ComputedAt = s.WindowEnd.UTC(), s.ComputedAt.UTC()
	s.TTFTFast95MeanMs, s.SuccessRate, s.CacheRate = providerHallScanDecimal(fast95), providerHallScanDecimal(successRate), providerHallScanDecimal(cacheRate)
	s.HistoricalPrice, s.SubscriptionShare = providerHallScanDecimal(price), providerHallScanDecimal(share)
	if d := providerHallScanDecimal(billedCost); d != nil {
		s.BilledInputCost = *d
	}
	if p90.Valid {
		v := p90.Int64
		s.TTFTP90Ms = &v
	}
	return s, nil
}

func (r *providerHallReadRepository) LoadSnapshots(ctx context.Context, windowEnd time.Time, algorithmVersion int) ([]service.ProviderHallSnapshotRow, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+providerHallSnapshotColumns+`
FROM provider_hall_snapshots WHERE window_end = $1 AND algorithm_version = $2`, windowEnd.UTC(), algorithmVersion)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []service.ProviderHallSnapshotRow{}
	for rows.Next() {
		s, err := providerHallScanSnapshot(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *providerHallReadRepository) LoadSeries(ctx context.Context, windowEnds []time.Time, algorithmVersion int) ([]service.ProviderHallSeriesPoint, error) {
	if len(windowEnds) == 0 {
		return []service.ProviderHallSeriesPoint{}, nil
	}
	ends := make([]time.Time, len(windowEnds))
	for i, t := range windowEnds {
		ends[i] = t.UTC()
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT group_id, profile_id, window_end, ttft_sample_count, ttft_fast95_mean_ms::float8, coverage
FROM provider_hall_snapshots
WHERE profile_id = 0 AND tier = 5 AND algorithm_version = $2 AND window_end = ANY($1::timestamptz[])`, pq.Array(ends), algorithmVersion)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []service.ProviderHallSeriesPoint{}
	for rows.Next() {
		var p service.ProviderHallSeriesPoint
		var fast95 sql.NullFloat64
		if err := rows.Scan(&p.GroupID, &p.ProfileID, &p.WindowEnd, &p.SampleCount, &fast95, &p.Coverage); err != nil {
			return nil, err
		}
		p.WindowEnd = p.WindowEnd.UTC()
		if fast95.Valid {
			v := fast95.Float64
			p.Fast95Ms = &v
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// providerHallQualified prefixes every column of a comma-separated list.
func providerHallQualified(alias, columns string) string {
	parts := strings.Split(columns, ",")
	for i, c := range parts {
		parts[i] = alias + "." + strings.TrimSpace(c)
	}
	return strings.Join(parts, ", ")
}

func (r *providerHallReadRepository) LatestVerifications(ctx context.Context) ([]service.ProviderHallVerification, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT DISTINCT ON (v.group_id, v.profile_id) `+providerHallQualified("v", providerHallVerificationColumns)+`
FROM provider_hall_verifications v
JOIN provider_hall_groups h ON h.group_id = v.group_id AND h.listed
ORDER BY v.group_id, v.profile_id, v.completed_at DESC, v.job_id DESC`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []service.ProviderHallVerification{}
	for rows.Next() {
		v, err := providerHallScanVerification(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *v)
	}
	return out, rows.Err()
}

func (r *providerHallReadRepository) ProbeSeriesAll(ctx context.Context, start, end time.Time, bucketSeconds int) ([]service.ProviderHallGroupProbePoint, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT j.group_id, date_bin(make_interval(secs => $1), s.received_at, $2::timestamptz) AS bucket,
       count(*), count(*) FILTER (WHERE s.result = 'passed'), count(*) FILTER (WHERE s.result <> 'passed' OR s.result IS NULL),
       avg(s.total_ms) FILTER (WHERE s.result = 'passed'), avg(s.ttft_ms) FILTER (WHERE s.result = 'passed')
FROM provider_hall_samples s JOIN provider_hall_jobs j ON j.id = s.job_id
WHERE j.kind = 'probe' AND s.status = 'received' AND s.received_at >= $2 AND s.received_at < $3
GROUP BY j.group_id, bucket ORDER BY j.group_id, bucket`, bucketSeconds, start.UTC(), end.UTC())
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []service.ProviderHallGroupProbePoint{}
	for rows.Next() {
		var p service.ProviderHallGroupProbePoint
		var avgTotal, avgTTFT sql.NullFloat64
		if err := rows.Scan(&p.GroupID, &p.BucketStart, &p.Total, &p.Passed, &p.Failed, &avgTotal, &avgTTFT); err != nil {
			return nil, err
		}
		p.BucketStart = p.BucketStart.UTC()
		if avgTotal.Valid {
			p.AvgTotalMs = &avgTotal.Float64
		}
		if avgTTFT.Valid {
			p.AvgTTFTMs = &avgTTFT.Float64
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
