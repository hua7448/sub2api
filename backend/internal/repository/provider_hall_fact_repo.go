package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

// providerHallFactRepository writes request facts with raw SQL. Fact tables
// deliberately have no Ent schema: they are append-heavy and never edited by
// hand through the admin UI.
type providerHallFactRepository struct{ db *sql.DB }

func NewProviderHallFactRepository(db *sql.DB) service.ProviderHallFactRepository {
	return &providerHallFactRepository{db: db}
}

func (r *providerHallFactRepository) RegisterEpoch(ctx context.Context, nodeID, buildVersion string, algorithmVersion int, now time.Time) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()
	// A restart without a clean exit is a crash: everything after the previous
	// epoch's confirmation point is a conservative gap.
	rows, err := tx.QueryContext(ctx, `SELECT id, COALESCE(confirmed_at, registered_at), heartbeat_at FROM provider_hall_collector_epochs WHERE node_id = $1 AND exited_at IS NULL FOR UPDATE`, nodeID)
	if err != nil {
		return 0, err
	}
	type stale struct {
		id                     int64
		confirmedAt, heartbeat time.Time
	}
	var stales []stale
	for rows.Next() {
		var s stale
		if err := rows.Scan(&s.id, &s.confirmedAt, &s.heartbeat); err != nil {
			_ = rows.Close()
			return 0, err
		}
		stales = append(stales, s)
	}
	_ = rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	for _, s := range stales {
		if _, err := tx.ExecContext(ctx, `UPDATE provider_hall_collector_epochs SET exited_at = $2, exit_reason = 'crash' WHERE id = $1`, s.id, now); err != nil {
			return 0, err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO provider_hall_coverage_gaps (node_id, epoch_id, scope, started_at, ended_at, reason) VALUES ($1, $2, 'collection', $3, $4, 'crash')`, nodeID, s.id, s.confirmedAt, now); err != nil {
			return 0, err
		}
	}
	var id int64
	if err := tx.QueryRowContext(ctx, `INSERT INTO provider_hall_collector_epochs (node_id, build_version, algorithm_version, registered_at, heartbeat_at) VALUES ($1, $2, $3, $4, $4) RETURNING id`, nodeID, buildVersion, algorithmVersion, now).Scan(&id); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *providerHallFactRepository) Heartbeat(ctx context.Context, epochID int64, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE provider_hall_collector_epochs SET heartbeat_at = $2 WHERE id = $1 AND exited_at IS NULL`, epochID, now)
	return err
}

func (r *providerHallFactRepository) MarkEpochExited(ctx context.Context, epochID int64, reason string, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE provider_hall_collector_epochs SET exited_at = $2, exit_reason = $3, heartbeat_at = $2 WHERE id = $1 AND exited_at IS NULL`, epochID, now, reason)
	return err
}

func (r *providerHallFactRepository) OpenGap(ctx context.Context, gap service.ProviderHallGap) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `INSERT INTO provider_hall_coverage_gaps (node_id, epoch_id, scope, started_at, reason) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		gap.NodeID, gap.EpochID, gap.Scope, gap.StartedAt, gap.Reason).Scan(&id)
	if err == nil {
		_, _ = r.db.ExecContext(ctx, `UPDATE provider_hall_collector_epochs SET overflowed = true WHERE id = $1`, gap.EpochID)
	}
	return id, err
}

func (r *providerHallFactRepository) CloseGap(ctx context.Context, id int64, endedAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE provider_hall_coverage_gaps SET ended_at = GREATEST($2, started_at) WHERE id = $1 AND ended_at IS NULL`, id, endedAt)
	return err
}

func (r *providerHallFactRepository) ListEnabledTargets(ctx context.Context) ([]service.ProviderHallTargetRef, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT t.id, t.group_id, t.profile_id, p.model, p.protocol, t.enabled, t.probe_key_id,
       t.probe_interval_seconds, t.verification_interval_seconds, t.version, p.version,
       p.supports_tools, p.output_limit, p.model_aliases
FROM provider_hall_targets t
JOIN provider_hall_profiles p ON p.id = t.profile_id
WHERE t.enabled
ORDER BY t.group_id, t.profile_id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallTargetRef
	for rows.Next() {
		var t service.ProviderHallTargetRef
		var probeKey sql.NullInt64
		var aliases []byte
		if err := rows.Scan(&t.TargetID, &t.GroupID, &t.ProfileID, &t.Model, &t.Protocol, &t.Enabled, &probeKey,
			&t.ProbeIntervalSeconds, &t.VerificationIntervalSeconds, &t.TargetVersion, &t.ProfileVersion,
			&t.SupportsTools, &t.OutputLimit, &aliases); err != nil {
			return nil, err
		}
		if probeKey.Valid {
			id := probeKey.Int64
			t.ProbeKeyID = &id
		}
		if len(aliases) > 0 {
			_ = json.Unmarshal(aliases, &t.ModelAliases)
		}
		if t.ModelAliases == nil {
			t.ModelAliases = []string{}
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *providerHallFactRepository) ListProbeKeys(ctx context.Context) (map[int64]time.Time, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT api_key_id, registered_at FROM provider_hall_probe_keys`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := map[int64]time.Time{}
	for rows.Next() {
		var id int64
		var at time.Time
		if err := rows.Scan(&id, &at); err != nil {
			return nil, err
		}
		out[id] = at.UTC()
	}
	return out, rows.Err()
}

const providerHallInsertRequestSQL = `
INSERT INTO provider_hall_requests (
    trace_id, node_id, epoch_id, last_seq, group_id, profile_id, api_key_id, protocol,
    requested_model, upstream_model, response_model, platform, source, sample_id, stream,
    started_at, first_content_at, ended_at, ttft_ms, submissions, outcome, exclusion_reason,
    usage_known, input_tokens, output_tokens, cache_read_tokens, cache_creation_tokens, algorithm_version, updated_at
) VALUES ($1::uuid, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, now())`

const providerHallFinishConflictSQL = `
ON CONFLICT (trace_id) DO UPDATE SET
    last_seq = GREATEST(provider_hall_requests.last_seq, EXCLUDED.last_seq),
    upstream_model = CASE WHEN EXCLUDED.upstream_model <> '' THEN EXCLUDED.upstream_model ELSE provider_hall_requests.upstream_model END,
    response_model = CASE WHEN EXCLUDED.response_model <> '' THEN EXCLUDED.response_model ELSE provider_hall_requests.response_model END,
    platform = CASE WHEN EXCLUDED.platform <> '' THEN EXCLUDED.platform ELSE provider_hall_requests.platform END,
    first_content_at = EXCLUDED.first_content_at, ended_at = EXCLUDED.ended_at, ttft_ms = EXCLUDED.ttft_ms,
    submissions = EXCLUDED.submissions, outcome = EXCLUDED.outcome, exclusion_reason = EXCLUDED.exclusion_reason,
    usage_known = EXCLUDED.usage_known, input_tokens = EXCLUDED.input_tokens, output_tokens = EXCLUDED.output_tokens,
    cache_read_tokens = EXCLUDED.cache_read_tokens, cache_creation_tokens = EXCLUDED.cache_creation_tokens,
    updated_at = now()
WHERE provider_hall_requests.outcome = 'pending'`

const providerHallBillingUpdateSQL = `
UPDATE provider_hall_requests SET
    billing_request_id = $2, billing_api_key_id = $3, request_fingerprint = $4, billing_status = $5,
    billing_mode = $6, is_subscription = $7, rate_multiplier = $8, input_base_cost = $9::numeric,
    total_base_cost = $10::numeric, actual_cost = $11::numeric, billed_at = $12, updated_at = now()
WHERE trace_id = $1::uuid AND billing_status IN ('pending', 'uncertain')`

func providerHallRequestArgs(row service.ProviderHallRequestRow) []any {
	var ttft any
	if row.TTFTMs != nil {
		ttft = *row.TTFTMs
	}
	toNull := func(v *int64) any {
		if v == nil {
			return nil
		}
		return *v
	}
	var first, ended any
	if row.FirstContentAt != nil {
		first = *row.FirstContentAt
	}
	if row.EndedAt != nil {
		ended = *row.EndedAt
	}
	outcome := row.Outcome
	if outcome == "" {
		outcome = service.ProviderHallOutcomePending
	}
	algo := row.AlgorithmVersion
	if algo == 0 {
		algo = service.ProviderHallAlgorithmVersion
	}
	return []any{
		row.TraceID.String(), row.NodeID, row.EpochID, int64(row.Seq), row.GroupID, toNull(row.ProfileID), row.APIKeyID, row.Protocol,
		row.RequestedModel, row.UpstreamModel, row.ResponseModel, row.Platform, string(row.Source), toNull(row.SampleID), row.Stream,
		row.StartedAt, first, ended, ttft, row.Submissions, string(outcome), string(row.ExclusionReason),
		row.UsageKnown, toNull(row.InputTokens), toNull(row.OutputTokens), toNull(row.CacheReadTokens), toNull(row.CacheCreation), algo,
	}
}

func providerHallDecimalArg(d decimal.Decimal) string { return d.StringFixed(10) }

func (r *providerHallFactRepository) WriteBatch(ctx context.Context, batch service.ProviderHallFactBatch) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	touched := make([]string, 0, len(batch.Finishes)+len(batch.Billings))
	if len(batch.Starts) > 0 {
		stmt, err := tx.PrepareContext(ctx, providerHallInsertRequestSQL+` ON CONFLICT (trace_id) DO NOTHING`)
		if err != nil {
			return fmt.Errorf("prepare start: %w", err)
		}
		for _, row := range batch.Starts {
			if _, err := stmt.ExecContext(ctx, providerHallRequestArgs(row)...); err != nil {
				_ = stmt.Close()
				return fmt.Errorf("insert start %s: %w", row.TraceID, err)
			}
		}
		_ = stmt.Close()
	}
	if len(batch.Finishes) > 0 {
		stmt, err := tx.PrepareContext(ctx, providerHallInsertRequestSQL+providerHallFinishConflictSQL)
		if err != nil {
			return fmt.Errorf("prepare finish: %w", err)
		}
		for _, row := range batch.Finishes {
			if _, err := stmt.ExecContext(ctx, providerHallRequestArgs(row)...); err != nil {
				_ = stmt.Close()
				return fmt.Errorf("upsert finish %s: %w", row.TraceID, err)
			}
			touched = append(touched, row.TraceID.String())
		}
		_ = stmt.Close()
	}
	if len(batch.Billings) > 0 {
		stmt, err := tx.PrepareContext(ctx, providerHallBillingUpdateSQL)
		if err != nil {
			return fmt.Errorf("prepare billing: %w", err)
		}
		for _, ev := range batch.Billings {
			if _, err := stmt.ExecContext(ctx, ev.TraceID.String(), ev.BillingRequestID, ev.APIKeyID, ev.Fingerprint, ev.Status,
				ev.Mode, ev.IsSubscription, ev.Multiplier, providerHallDecimalArg(ev.InputBaseCost), providerHallDecimalArg(ev.TotalBaseCost),
				providerHallDecimalArg(ev.ActualCost), ev.BilledAt); err != nil {
				_ = stmt.Close()
				return fmt.Errorf("update billing %s: %w", ev.TraceID, err)
			}
			touched = append(touched, ev.TraceID.String())
		}
		_ = stmt.Close()
		if err := providerHallWriteSpend(ctx, tx, batch.Billings); err != nil {
			return err
		}
	}
	if len(touched) > 0 {
		// Fact changes and dirty marks commit together so the aggregator can
		// never miss a late finish or a late bill.
		if _, err := tx.ExecContext(ctx, `
INSERT INTO provider_hall_dirty_buckets (group_id, profile_id, minute, reason)
SELECT group_id, profile_id, date_trunc('minute', ended_at), 'collector'
FROM provider_hall_requests
WHERE trace_id = ANY($1::uuid[]) AND ended_at IS NOT NULL AND source = 'user' AND profile_id IS NOT NULL
ON CONFLICT (group_id, profile_id, minute) DO NOTHING`, pq.Array(touched)); err != nil {
			return fmt.Errorf("mark dirty: %w", err)
		}
	}
	if batch.EpochID > 0 {
		if batch.ConfirmedAt != nil {
			if _, err := tx.ExecContext(ctx, `UPDATE provider_hall_collector_epochs SET persisted_seq = GREATEST(persisted_seq, $2), confirmed_at = GREATEST(COALESCE(confirmed_at, $3), $3), overflowed = false, heartbeat_at = now() WHERE id = $1`, batch.EpochID, int64(batch.MaxSeq), *batch.ConfirmedAt); err != nil {
				return fmt.Errorf("advance watermark: %w", err)
			}
		} else if _, err := tx.ExecContext(ctx, `UPDATE provider_hall_collector_epochs SET persisted_seq = GREATEST(persisted_seq, $2), overflowed = $3, heartbeat_at = now() WHERE id = $1`, batch.EpochID, int64(batch.MaxSeq), batch.Overflowed); err != nil {
			return fmt.Errorf("record seq: %w", err)
		}
	}
	return tx.Commit()
}

const providerHallSpendUpsertSQL = `
INSERT INTO provider_hall_spend (billing_request_id, api_key_id, job_id, sample_id, budget_day, actual_cost, status, confirmed_at)
SELECT $2, $3, s.job_id, s.id,
       COALESCE(j.budget_day, ($6::timestamptz AT TIME ZONE 'Asia/Shanghai')::date),
       $4::numeric, $5::varchar, CASE WHEN $5::varchar = 'confirmed' THEN $6::timestamptz END
FROM provider_hall_samples s
JOIN provider_hall_jobs j ON j.id = s.job_id
WHERE s.trace_id = $1::uuid
ON CONFLICT (billing_request_id, api_key_id) DO UPDATE SET
    status = EXCLUDED.status, actual_cost = EXCLUDED.actual_cost, confirmed_at = EXCLUDED.confirmed_at,
    job_id = EXCLUDED.job_id, sample_id = EXCLUDED.sample_id
WHERE provider_hall_spend.status <> 'confirmed'`

const providerHallSampleBillingMirrorSQL = `
UPDATE provider_hall_samples SET billing_status = $2, actual_cost = $3::numeric
WHERE trace_id = $1::uuid AND billing_status <> 'confirmed'`

// providerHallWriteSpend mirrors probe/verification bills into the spend
// ledger and the sample row (batch B4). Requests of ordinary users are only
// billed to themselves and never appear here; a trace without a sample row
// matches nothing.
func providerHallWriteSpend(ctx context.Context, tx *sql.Tx, billings []service.ProviderHallBillingEvent) error {
	for _, ev := range billings {
		spendStatus, sampleStatus, writeSpend, ok := service.ProviderHallSpendStatusForBilling(ev.Status)
		if !ok {
			continue
		}
		if writeSpend && ev.BillingRequestID != "" {
			if _, err := tx.ExecContext(ctx, providerHallSpendUpsertSQL, ev.TraceID.String(), ev.BillingRequestID, ev.APIKeyID,
				providerHallDecimalArg(ev.ActualCost), spendStatus, ev.BilledAt); err != nil {
				return fmt.Errorf("upsert spend %s: %w", ev.TraceID, err)
			}
		}
		var cost any
		if sampleStatus == service.ProviderHallSampleBillingConfirmed {
			cost = providerHallDecimalArg(ev.ActualCost)
		}
		if _, err := tx.ExecContext(ctx, providerHallSampleBillingMirrorSQL, ev.TraceID.String(), sampleStatus, cost); err != nil {
			return fmt.Errorf("mirror sample billing %s: %w", ev.TraceID, err)
		}
	}
	return nil
}
