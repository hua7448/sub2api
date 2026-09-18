package service

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

const (
	resetTodayUsageTimezone       = subscriptionDailyCalendarTimezone
	resetTodayUsageAdvisoryLockID = int64(0x7375623272747501)
	resetTodayUsageCacheWorkers   = 16
	resetTodayUsageCacheBatchTTL  = 15 * time.Second
)

// ResetTodayUsageResult is the durable operation summary returned to the admin API.
type ResetTodayUsageResult struct {
	CutoffAt                  time.Time `json:"cutoff_at"`
	CapturedAt                time.Time `json:"captured_at"`
	Subscriptions             int       `json:"subscriptions"`
	TodayWindowSubscriptions  int       `json:"today_window_subscriptions"`
	ResetSubscriptions        int       `json:"reset_subscriptions"`
	ResetAndDeductedUSD       float64   `json:"reset_and_deducted_usd"`
	StaleDailyClearedUSD      float64   `json:"stale_daily_cleared_usd"`
	CacheInvalidations        int       `json:"cache_invalidations"`
	CacheInvalidationFailures int       `json:"cache_invalidation_failures"`
}

type resetTodayUsageTarget struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}

type resetTodayUsageAdjustment struct {
	ID        int64
	DeductUSD decimal.Decimal
}

type resetTodayUsageOperation struct {
	ID      int64
	Result  *ResetTodayUsageResult
	Targets []resetTodayUsageTarget
}

// AdminResetTodayUsage resets the current Shanghai calendar day's subscription
// usage while preserving all quota-window start timestamps. The operation hash
// provides a second, transactional idempotency boundary in addition to the HTTP
// idempotency coordinator.
func (s *SubscriptionService) AdminResetTodayUsage(ctx context.Context, actorUserID int64, idempotencyKeyHash string) (*ResetTodayUsageResult, error) {
	if s == nil || s.entClient == nil {
		return nil, fmt.Errorf("subscription reset requires a database client")
	}
	if !validResetTodayUsageKeyHash(idempotencyKeyHash) {
		return nil, ErrIdempotencyKeyInvalid
	}

	op, err := s.resetTodayUsageTransaction(ctx, actorUserID, idempotencyKeyHash)
	if err != nil {
		return nil, err
	}

	succeeded, failedTargets := s.invalidateResetTodayUsageTargetBatch(op.Targets, op.Result.CutoffAt, true)
	s.scheduleResetTodayUsageCacheRetries(failedTargets, op.Result.CutoffAt)
	op.Result.CacheInvalidations = succeeded
	op.Result.CacheInvalidationFailures = len(failedTargets)
	s.persistResetTodayUsageCacheSummary(op.ID, op.Result)
	return op.Result, nil
}

func validResetTodayUsageKeyHash(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func (s *SubscriptionService) resetTodayUsageTransaction(ctx context.Context, actorUserID int64, keyHash string) (_ *resetTodayUsageOperation, retErr error) {
	tx, err := s.entClient.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin reset-today transaction: %w", err)
	}
	defer func() {
		if retErr != nil {
			_ = tx.Rollback()
		}
	}()
	client := tx.Client()

	operationID, claimed, err := claimResetTodayUsageOperation(ctx, client, actorUserID, keyHash)
	if err != nil {
		return nil, err
	}
	if !claimed {
		op, err := loadResetTodayUsageOperation(ctx, client, actorUserID, keyHash)
		if err != nil {
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit reset-today replay: %w", err)
		}
		return op, nil
	}

	if _, err := client.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1)`, resetTodayUsageAdvisoryLockID); err != nil {
		return nil, fmt.Errorf("lock reset-today operation: %w", err)
	}
	if err := lockActiveSubscriptions(ctx, client); err != nil {
		return nil, err
	}

	wallTime, cutoff, err := resetTodayUsageDatabaseTime(ctx, client)
	if err != nil {
		return nil, err
	}
	dayStart, err := resetTodayUsageDayStart(wallTime)
	if err != nil {
		return nil, err
	}
	if err := backfillLegacyResetTodayUsageOperations(ctx, client, wallTime, dayStart); err != nil {
		return nil, err
	}

	result, targets, adjustments, err := captureResetTodayUsageTargets(ctx, client, wallTime, dayStart)
	if err != nil {
		return nil, err
	}
	result.CutoffAt = cutoff
	result.CapturedAt = wallTime

	if len(adjustments) > 0 {
		ids, deductions := resetTodayUsageAdjustmentArrays(adjustments)
		// A Shanghai calendar day's ledger may span a weekly or monthly window
		// rollover. Cap each window independently so only usage still present in
		// that window is deducted, while a longer-lived window receives the full
		// eligible daily deduction.
		updated, err := client.ExecContext(ctx, `
			WITH adjustments AS MATERIALIZED (
				SELECT id, deduct_usd
				FROM unnest($4::bigint[], $5::numeric[]) AS adjustment(id, deduct_usd)
			)
			UPDATE user_subscriptions AS us
			SET daily_usage_usd = 0,
				weekly_usage_usd = GREATEST(
					us.weekly_usage_usd - LEAST(
						adjustments.deduct_usd,
						GREATEST(us.weekly_usage_usd, 0)
					),
					0
				),
				monthly_usage_usd = GREATEST(
					us.monthly_usage_usd - LEAST(
						adjustments.deduct_usd,
						GREATEST(us.monthly_usage_usd, 0)
					),
					0
				),
				daily_quota_reset_operations = COALESCE((
					SELECT jsonb_agg(
						CASE
							WHEN operation.value ?& ARRAY['reset_at', 'cleared_daily_window_start', 'cleared_daily_usage_usd']
								AND NOT operation.value ? 'admin_reset_consumed_operation_id'
								AND jsonb_typeof(operation.value->'reset_at') = 'string'
								AND jsonb_typeof(operation.value->'cleared_daily_window_start') = 'string'
								AND jsonb_typeof(operation.value->'cleared_daily_usage_usd') = 'number'
								AND (operation.value->>'reset_at')::timestamptz >= $1
								AND (operation.value->>'reset_at')::timestamptz <= $2
								AND (operation.value->>'cleared_daily_window_start')::timestamptz >= $1
								AND (operation.value->>'cleared_daily_window_start')::timestamptz <= $2
							THEN operation.value || jsonb_build_object(
								'admin_reset_consumed_operation_id', $6::bigint,
								'admin_reset_consumed_at', $2
							)
							ELSE operation.value
						END
						ORDER BY operation.position
					)
					FROM jsonb_array_elements(
						CASE
							WHEN jsonb_typeof(us.daily_quota_reset_operations) = 'array'
							THEN us.daily_quota_reset_operations
							ELSE '[]'::jsonb
						END
					) WITH ORDINALITY AS operation(value, position)
				), '[]'::jsonb),
				updated_at = $3
			FROM adjustments
			WHERE us.id = adjustments.id
		`, dayStart.UTC(), wallTime, cutoff, pq.Array(ids), pq.Array(deductions), operationID)
		if err != nil {
			return nil, fmt.Errorf("update reset-today subscriptions: %w", err)
		}
		affected, err := updated.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("count reset-today subscriptions: %w", err)
		}
		if affected != int64(len(adjustments)) {
			return nil, fmt.Errorf("reset-today updated %d subscriptions, expected %d", affected, len(adjustments))
		}
	}

	if err := completeResetTodayUsageOperation(ctx, client, operationID, result, targets); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit reset-today transaction: %w", err)
	}
	return &resetTodayUsageOperation{ID: operationID, Result: result, Targets: targets}, nil
}

type resetTodayUsageDB interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func resetTodayUsageDatabaseTime(ctx context.Context, client resetTodayUsageDB) (time.Time, time.Time, error) {
	rows, err := client.QueryContext(ctx, `
		WITH reset_clock AS MATERIALIZED (
			SELECT clock_timestamp() AS captured_at
		)
		SELECT
			captured_at,
			GREATEST(
				captured_at,
				COALESCE((
					SELECT MAX(updated_at) + interval '1 microsecond'
					FROM user_subscriptions
					WHERE deleted_at IS NULL AND status = $1
				), captured_at)
			) AS cutoff_at
		FROM reset_clock
	`, SubscriptionStatusActive)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("read reset-today database clock: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("read reset-today database clock row: %w", err)
		}
		return time.Time{}, time.Time{}, fmt.Errorf("read reset-today database clock: %w", sql.ErrNoRows)
	}
	var capturedAt, cutoffAt time.Time
	if err := rows.Scan(&capturedAt, &cutoffAt); err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("scan reset-today database clock: %w", err)
	}
	return capturedAt.UTC(), cutoffAt.UTC(), nil
}

func resetTodayUsageDayStart(cutoff time.Time) (time.Time, error) {
	return subscriptionDailyCalendarStart(cutoff), nil
}

func claimResetTodayUsageOperation(ctx context.Context, client resetTodayUsageDB, actorUserID int64, keyHash string) (int64, bool, error) {
	rows, err := client.QueryContext(ctx, `
		INSERT INTO subscription_quota_reset_operations (
			actor_user_id, idempotency_key_hash, status
		) VALUES ($1, $2, 'processing')
		ON CONFLICT (actor_user_id, idempotency_key_hash) DO NOTHING
		RETURNING id
	`, actorUserID, keyHash)
	if err != nil {
		return 0, false, fmt.Errorf("claim reset-today operation: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return 0, false, fmt.Errorf("read reset-today operation claim: %w", err)
		}
		return 0, false, nil
	}
	var operationID int64
	if err := rows.Scan(&operationID); err != nil {
		return 0, false, fmt.Errorf("scan reset-today operation claim: %w", err)
	}
	return operationID, true, nil
}

func loadResetTodayUsageOperation(ctx context.Context, client resetTodayUsageDB, actorUserID int64, keyHash string) (*resetTodayUsageOperation, error) {
	rows, err := client.QueryContext(ctx, `
		SELECT id, status, summary, target_pairs
		FROM subscription_quota_reset_operations
		WHERE actor_user_id = $1 AND idempotency_key_hash = $2
	`, actorUserID, keyHash)
	if err != nil {
		return nil, fmt.Errorf("load reset-today operation: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("read reset-today operation: %w", err)
		}
		return nil, fmt.Errorf("reset-today operation disappeared after conflict")
	}

	var (
		operationID int64
		status      string
		summaryJSON []byte
		targetsJSON []byte
	)
	if err := rows.Scan(&operationID, &status, &summaryJSON, &targetsJSON); err != nil {
		return nil, fmt.Errorf("scan reset-today operation: %w", err)
	}
	if status != "succeeded" {
		return nil, fmt.Errorf("reset-today operation has unexpected status %q", status)
	}
	result := new(ResetTodayUsageResult)
	if err := json.Unmarshal(summaryJSON, result); err != nil {
		return nil, fmt.Errorf("decode reset-today operation summary: %w", err)
	}
	var targets []resetTodayUsageTarget
	if err := json.Unmarshal(targetsJSON, &targets); err != nil {
		return nil, fmt.Errorf("decode reset-today operation targets: %w", err)
	}
	return &resetTodayUsageOperation{ID: operationID, Result: result, Targets: targets}, nil
}

func lockActiveSubscriptions(ctx context.Context, client resetTodayUsageDB) error {
	rows, err := client.QueryContext(ctx, `
		SELECT id
		FROM user_subscriptions
		WHERE deleted_at IS NULL AND status = $1
		ORDER BY id
		FOR UPDATE
	`, SubscriptionStatusActive)
	if err != nil {
		return fmt.Errorf("lock active subscriptions for reset-today: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("scan active subscription lock: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("read active subscription locks: %w", err)
	}
	return nil
}

func resetTodayUsageWindowIsToday(windowStart sql.NullTime, dayStart, cutoff time.Time) bool {
	return windowStart.Valid &&
		!windowStart.Time.Before(dayStart) &&
		!windowStart.Time.After(cutoff)
}

// backfillLegacyResetTodayUsageOperations bridges self-reset records written by
// binaries that predate the cleared_* ledger fields. It reconstructs only usage
// rows already durable before each reset_at. A concurrently delayed usage-log
// commit cannot be attributed retrospectively, so this compatibility path does
// not sleep while holding every active subscription row lock. PostgreSQL 15 has
// no safe timestamp-input probe; non-string values are counted as malformed,
// while invalid timestamp strings deliberately abort the transaction on cast.
func backfillLegacyResetTodayUsageOperations(ctx context.Context, client resetTodayUsageDB, cutoff, dayStart time.Time) error {
	rows, err := client.QueryContext(ctx, `
		WITH subscription_operations AS MATERIALIZED (
			SELECT
				us.id AS subscription_id,
				us.user_id,
				us.group_id,
				us.starts_at,
				us.weekly_window_start,
				us.monthly_window_start,
				operation.value,
				operation.position
			FROM user_subscriptions AS us
			CROSS JOIN LATERAL jsonb_array_elements(
				CASE
					WHEN jsonb_typeof(us.daily_quota_reset_operations) = 'array'
					THEN us.daily_quota_reset_operations
					ELSE '[]'::jsonb
				END
			) WITH ORDINALITY AS operation(value, position)
			WHERE us.deleted_at IS NULL
				AND us.status = $1
				AND us.starts_at <= $2
				AND us.expires_at > $2
		),
		successful_admin_resets AS MATERIALIZED (
			SELECT cutoff_at, target_pairs
			FROM subscription_quota_reset_operations
			WHERE status = 'succeeded'
				AND cutoff_at IS NOT NULL
				AND cutoff_at < $2
		),
		parsed_self_reset_operations AS MATERIALIZED (
			SELECT
				operation.*,
				CASE
					WHEN NOT operation.value ? 'admin_reset_consumed_operation_id'
						AND jsonb_typeof(operation.value->'reset_at') = 'string'
					THEN (operation.value->>'reset_at')::timestamptz
				END AS reset_at,
				CASE
					WHEN NOT operation.value ? 'admin_reset_consumed_operation_id'
						AND jsonb_typeof(operation.value->'cleared_daily_window_start') = 'string'
					THEN (operation.value->>'cleared_daily_window_start')::timestamptz
				END AS cleared_window_start,
				CASE
					WHEN NOT operation.value ? 'admin_reset_consumed_operation_id'
						AND jsonb_typeof(operation.value->'replay_expires_at') = 'string'
					THEN (operation.value->>'replay_expires_at')::timestamptz
				END AS replay_expires_at
			FROM subscription_operations AS operation
		),
		valid_self_reset_operations AS MATERIALIZED (
			SELECT parsed.*
			FROM parsed_self_reset_operations AS parsed
			WHERE parsed.value ?& ARRAY[
					'idempotency_key_hash', 'subscription_id', 'group_id', 'reset_at',
					'daily_usage_usd', 'weekly_usage_usd', 'monthly_usage_usd',
					'daily_quota_reset_available_at', 'replay_expires_at'
				]
				AND jsonb_typeof(parsed.value->'idempotency_key_hash') = 'string'
				AND parsed.value->>'idempotency_key_hash' ~ '^[0-9a-f]{64}$'
				AND jsonb_typeof(parsed.value->'subscription_id') = 'number'
				AND jsonb_typeof(parsed.value->'group_id') = 'number'
				AND parsed.reset_at IS NOT NULL
				AND jsonb_typeof(parsed.value->'daily_usage_usd') = 'number'
				AND (parsed.value->>'daily_usage_usd')::numeric = 0
				AND jsonb_typeof(parsed.value->'weekly_usage_usd') = 'number'
				AND jsonb_typeof(parsed.value->'monthly_usage_usd') = 'number'
				AND jsonb_typeof(parsed.value->'daily_quota_reset_available_at') IN ('string', 'null')
				AND jsonb_typeof(parsed.value->'replay_expires_at') = 'string'
				AND parsed.replay_expires_at >= parsed.reset_at
				AND parsed.value->'subscription_id' = to_jsonb(parsed.subscription_id)
				AND parsed.value->'group_id' = to_jsonb(parsed.group_id)
		),
		malformed_self_reset_operations AS MATERIALIZED (
			SELECT parsed.subscription_id, parsed.position
			FROM parsed_self_reset_operations AS parsed
			WHERE NOT parsed.value ? 'admin_reset_consumed_operation_id'
				AND (
					jsonb_typeof(parsed.value->'reset_at') IS DISTINCT FROM 'string'
					OR (
						parsed.reset_at >= $3
						AND parsed.reset_at <= $2
						AND (
							(
								NOT (parsed.value ?| ARRAY[
									'cleared_daily_window_start', 'cleared_daily_usage_usd'
								])
								AND NOT EXISTS (
									SELECT 1
									FROM valid_self_reset_operations AS valid
									WHERE valid.subscription_id = parsed.subscription_id
										AND valid.position = parsed.position
								)
							)
							OR (
								parsed.value ?| ARRAY[
									'cleared_daily_window_start', 'cleared_daily_usage_usd'
								]
								AND NOT (
									parsed.value ?& ARRAY[
										'cleared_daily_window_start', 'cleared_daily_usage_usd'
									]
									AND parsed.cleared_window_start IS NOT NULL
									AND jsonb_typeof(parsed.value->'cleared_daily_usage_usd') = 'number'
									AND EXISTS (
										SELECT 1
										FROM valid_self_reset_operations AS valid
										WHERE valid.subscription_id = parsed.subscription_id
											AND valid.position = parsed.position
									)
								)
							)
						)
					)
				)
		),
		legacy_operations AS MATERIALIZED (
			SELECT valid.*
			FROM valid_self_reset_operations AS valid
			WHERE NOT (valid.value ?| ARRAY[
					'cleared_daily_window_start', 'cleared_daily_usage_usd'
				])
				AND NOT valid.value ? 'admin_reset_consumed_operation_id'
				AND valid.reset_at >= $3
				AND valid.reset_at <= $2
		),
		reconstruction_windows AS MATERIALIZED (
			SELECT
				legacy.subscription_id,
				legacy.user_id,
				legacy.group_id,
				legacy.position,
				legacy.reset_at,
				legacy.weekly_window_start,
				legacy.monthly_window_start,
				(legacy.value->>'weekly_usage_usd')::numeric AS stored_weekly_usage_usd,
				(legacy.value->>'monthly_usage_usd')::numeric AS stored_monthly_usage_usd,
				GREATEST(
					$3::timestamptz,
					legacy.starts_at,
					COALESCE((
						SELECT MAX(previous.reset_at)
						FROM valid_self_reset_operations AS previous
						WHERE previous.subscription_id = legacy.subscription_id
							AND previous.position < legacy.position
							AND previous.reset_at < legacy.reset_at
					), $3::timestamptz),
					COALESCE((
						SELECT MAX(admin_operation.cutoff_at)
						FROM successful_admin_resets AS admin_operation
						WHERE admin_operation.cutoff_at >= $3
							AND admin_operation.cutoff_at < legacy.reset_at
							AND (
								CASE
									WHEN jsonb_typeof(admin_operation.target_pairs) = 'array'
									THEN admin_operation.target_pairs
									ELSE '[]'::jsonb
								END
							) @> jsonb_build_array(jsonb_build_object(
								'user_id', legacy.user_id,
								'group_id', legacy.group_id
							))
					), $3::timestamptz)
				) AS cleared_window_start
			FROM legacy_operations AS legacy
		),
		legacy_proof_windows AS MATERIALIZED (
			SELECT
				reconstruction.*,
				CASE
					WHEN reconstruction.weekly_window_start IS NOT NULL
						AND reconstruction.weekly_window_start <= reconstruction.reset_at
						AND NOT EXISTS (
							SELECT 1
							FROM successful_admin_resets AS admin_operation
							WHERE admin_operation.cutoff_at >= reconstruction.weekly_window_start
								AND admin_operation.cutoff_at < reconstruction.reset_at
								AND (
									CASE
										WHEN jsonb_typeof(admin_operation.target_pairs) = 'array'
										THEN admin_operation.target_pairs
										ELSE '[]'::jsonb
									END
								) @> jsonb_build_array(jsonb_build_object(
									'user_id', reconstruction.user_id,
									'group_id', reconstruction.group_id
								))
						)
					THEN reconstruction.weekly_window_start
				END AS weekly_proof_start,
				CASE
					WHEN reconstruction.monthly_window_start IS NOT NULL
						AND reconstruction.monthly_window_start <= reconstruction.reset_at
						AND NOT EXISTS (
							SELECT 1
							FROM successful_admin_resets AS admin_operation
							WHERE admin_operation.cutoff_at >= reconstruction.monthly_window_start
								AND admin_operation.cutoff_at < reconstruction.reset_at
								AND (
									CASE
										WHEN jsonb_typeof(admin_operation.target_pairs) = 'array'
										THEN admin_operation.target_pairs
										ELSE '[]'::jsonb
									END
								) @> jsonb_build_array(jsonb_build_object(
									'user_id', reconstruction.user_id,
									'group_id', reconstruction.group_id
								))
						)
					THEN reconstruction.monthly_window_start
				END AS monthly_proof_start
			FROM reconstruction_windows AS reconstruction
		),
		legacy_usage_evidence AS MATERIALIZED (
			SELECT proof_window.*, usage_evidence.*
			FROM legacy_proof_windows AS proof_window
			CROSS JOIN LATERAL (
				SELECT
					COALESCE(SUM(GREATEST(usage_log.actual_cost, 0)) FILTER (
						WHERE usage_log.created_at >= proof_window.cleared_window_start
					), 0) AS subwindow_raw10_usd,
					COALESCE(SUM(ROUND(GREATEST(usage_log.actual_cost, 0), $5)) FILTER (
						WHERE usage_log.created_at >= proof_window.cleared_window_start
					), 0) AS subwindow_per_row_round8_usd,
					COALESCE(SUM(GREATEST(usage_log.actual_cost, 0)) FILTER (
						WHERE proof_window.weekly_proof_start IS NOT NULL
							AND usage_log.created_at >= proof_window.weekly_proof_start
					), 0) AS weekly_raw10_usd,
					COALESCE(SUM(ROUND(GREATEST(usage_log.actual_cost, 0), $5)) FILTER (
						WHERE proof_window.weekly_proof_start IS NOT NULL
							AND usage_log.created_at >= proof_window.weekly_proof_start
					), 0) AS weekly_per_row_round8_usd,
					COALESCE(SUM(GREATEST(usage_log.actual_cost, 0)) FILTER (
						WHERE proof_window.monthly_proof_start IS NOT NULL
							AND usage_log.created_at >= proof_window.monthly_proof_start
					), 0) AS monthly_raw10_usd,
					COALESCE(SUM(ROUND(GREATEST(usage_log.actual_cost, 0), $5)) FILTER (
						WHERE proof_window.monthly_proof_start IS NOT NULL
							AND usage_log.created_at >= proof_window.monthly_proof_start
					), 0) AS monthly_per_row_round8_usd
				FROM usage_logs AS usage_log
				WHERE usage_log.subscription_id = proof_window.subscription_id
					AND usage_log.billing_type = $4
					AND usage_log.created_at >= LEAST(
						proof_window.cleared_window_start,
						COALESCE(proof_window.weekly_proof_start, proof_window.cleared_window_start),
						COALESCE(proof_window.monthly_proof_start, proof_window.cleared_window_start)
					)
					AND usage_log.created_at < proof_window.reset_at
			) AS usage_evidence
		),
		legacy_usage_proofs AS MATERIALIZED (
			SELECT
				evidence.*,
				(
					(evidence.weekly_proof_start IS NOT NULL
						AND evidence.weekly_raw10_usd = evidence.stored_weekly_usage_usd)
					OR (evidence.monthly_proof_start IS NOT NULL
						AND evidence.monthly_raw10_usd = evidence.stored_monthly_usage_usd)
				) AS raw10_proven,
				(
					(evidence.weekly_proof_start IS NOT NULL
						AND evidence.weekly_per_row_round8_usd = evidence.stored_weekly_usage_usd)
					OR (evidence.monthly_proof_start IS NOT NULL
						AND evidence.monthly_per_row_round8_usd = evidence.stored_monthly_usage_usd)
				) AS per_row_round8_proven
			FROM legacy_usage_evidence AS evidence
		),
		unproven_legacy_operations AS MATERIALIZED (
			SELECT proof.subscription_id, proof.position
			FROM legacy_usage_proofs AS proof
			WHERE (NOT proof.raw10_proven AND NOT proof.per_row_round8_proven)
				OR (
					proof.raw10_proven
					AND proof.per_row_round8_proven
					AND proof.subwindow_raw10_usd <> proof.subwindow_per_row_round8_usd
				)
		),
		reconstructed_operations AS MATERIALIZED (
			SELECT
				proof.subscription_id,
				proof.position,
				proof.cleared_window_start,
				CASE
					WHEN proof.raw10_proven THEN proof.subwindow_raw10_usd
					ELSE proof.subwindow_per_row_round8_usd
				END AS cleared_usage_usd,
				CASE
					WHEN proof.raw10_proven AND proof.per_row_round8_proven
					THEN 'period_snapshot_proof=raw10+per_row_round8/subwindow_equal'
					WHEN proof.raw10_proven
					THEN 'period_snapshot_proof=raw10'
					ELSE 'period_snapshot_proof=per_row_round8'
				END AS reconstruction_source
			FROM legacy_usage_proofs AS proof
			WHERE proof.raw10_proven <> proof.per_row_round8_proven
				OR (
					proof.raw10_proven
					AND proof.per_row_round8_proven
					AND proof.subwindow_raw10_usd = proof.subwindow_per_row_round8_usd
				)
		),
		rebuilt_subscriptions AS MATERIALIZED (
			SELECT
				operation.subscription_id,
				jsonb_agg(
					CASE
						WHEN reconstructed.position IS NOT NULL
						THEN operation.value || jsonb_build_object(
							'cleared_daily_window_start', reconstructed.cleared_window_start,
							'cleared_daily_usage_usd', reconstructed.cleared_usage_usd,
							'legacy_cleared_usage_reconstructed_at', $2,
							'legacy_cleared_usage_reconstruction_source',
								reconstructed.reconstruction_source
						)
						ELSE operation.value
					END
					ORDER BY operation.position
				) AS operations
			FROM subscription_operations AS operation
			LEFT JOIN reconstructed_operations AS reconstructed
				ON reconstructed.subscription_id = operation.subscription_id
				AND reconstructed.position = operation.position
			WHERE EXISTS (
				SELECT 1
				FROM reconstructed_operations AS candidate
				WHERE candidate.subscription_id = operation.subscription_id
			)
			GROUP BY operation.subscription_id
		),
		updated_subscriptions AS (
			UPDATE user_subscriptions AS us
			SET daily_quota_reset_operations = rebuilt.operations
			FROM rebuilt_subscriptions AS rebuilt
			WHERE us.id = rebuilt.subscription_id
				AND NOT EXISTS (SELECT 1 FROM malformed_self_reset_operations)
				AND NOT EXISTS (SELECT 1 FROM unproven_legacy_operations)
			RETURNING us.id
		)
		SELECT
			(SELECT COUNT(*) FROM malformed_self_reset_operations),
			(SELECT COUNT(*) FROM unproven_legacy_operations),
			(SELECT COUNT(*) FROM updated_subscriptions)
	`, SubscriptionStatusActive, cutoff, dayStart.UTC(), BillingTypeSubscription, UsageBillingMonetaryScale)
	if err != nil {
		return fmt.Errorf("backfill legacy reset-today operations: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return fmt.Errorf("read legacy reset-today backfill result: %w", err)
		}
		return fmt.Errorf("read legacy reset-today backfill result: %w", sql.ErrNoRows)
	}
	var malformed, unproven, updated int64
	if err := rows.Scan(&malformed, &unproven, &updated); err != nil {
		return fmt.Errorf("scan legacy reset-today backfill result: %w", err)
	}
	if malformed > 0 {
		return fmt.Errorf("refuse reset-today with %d malformed current-day self-reset operations", malformed)
	}
	if unproven > 0 {
		return fmt.Errorf("refuse reset-today with %d unproven or ambiguous legacy self-reset operations", unproven)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("read legacy reset-today backfill rows: %w", err)
	}
	return nil
}

func captureResetTodayUsageTargets(ctx context.Context, client resetTodayUsageDB, cutoff, dayStart time.Time) (*ResetTodayUsageResult, []resetTodayUsageTarget, []resetTodayUsageAdjustment, error) {
	rows, err := client.QueryContext(ctx, `
		SELECT id, user_id, group_id, daily_window_start, daily_usage_usd::text,
			COALESCE((
				SELECT SUM(GREATEST((operation.value->>'cleared_daily_usage_usd')::numeric, 0))
				FROM jsonb_array_elements(
					CASE
						WHEN jsonb_typeof(user_subscriptions.daily_quota_reset_operations) = 'array'
						THEN user_subscriptions.daily_quota_reset_operations
						ELSE '[]'::jsonb
					END
				) AS operation(value)
				WHERE operation.value ?& ARRAY['reset_at', 'cleared_daily_window_start', 'cleared_daily_usage_usd']
					AND NOT operation.value ? 'admin_reset_consumed_operation_id'
					AND jsonb_typeof(operation.value->'reset_at') = 'string'
					AND jsonb_typeof(operation.value->'cleared_daily_window_start') = 'string'
					AND jsonb_typeof(operation.value->'cleared_daily_usage_usd') = 'number'
					AND (operation.value->>'reset_at')::timestamptz >= $3
					AND (operation.value->>'reset_at')::timestamptz <= $2
					AND (operation.value->>'cleared_daily_window_start')::timestamptz >= $3
					AND (operation.value->>'cleared_daily_window_start')::timestamptz <= $2
			), 0)::text AS unconsumed_reset_usage_usd
		FROM user_subscriptions
		WHERE deleted_at IS NULL
			AND status = $1
			AND starts_at <= $2
			AND expires_at > $2
		ORDER BY id
		FOR UPDATE
	`, SubscriptionStatusActive, cutoff, dayStart.UTC())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("capture reset-today subscriptions: %w", err)
	}
	defer rows.Close()

	result := new(ResetTodayUsageResult)
	targets := make([]resetTodayUsageTarget, 0)
	adjustments := make([]resetTodayUsageAdjustment, 0)
	seenPairs := make(map[resetTodayUsageTarget]struct{})
	resetAmount := decimal.Zero
	staleAmount := decimal.Zero

	for rows.Next() {
		var (
			id               int64
			userID           int64
			groupID          int64
			dailyWindowStart sql.NullTime
			dailyUsageRaw    string
			resetUsageRaw    string
		)
		if err := rows.Scan(&id, &userID, &groupID, &dailyWindowStart, &dailyUsageRaw, &resetUsageRaw); err != nil {
			return nil, nil, nil, fmt.Errorf("scan reset-today subscription: %w", err)
		}
		dailyUsage, err := decimal.NewFromString(dailyUsageRaw)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("parse subscription %d daily usage: %w", id, err)
		}
		positiveUsage := dailyUsage
		if positiveUsage.IsNegative() {
			positiveUsage = decimal.Zero
		}
		resetUsage, err := decimal.NewFromString(resetUsageRaw)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("parse subscription %d unconsumed reset usage: %w", id, err)
		}
		if resetUsage.IsNegative() {
			resetUsage = decimal.Zero
		}

		result.Subscriptions++
		pair := resetTodayUsageTarget{UserID: userID, GroupID: groupID}
		if _, exists := seenPairs[pair]; !exists {
			seenPairs[pair] = struct{}{}
			targets = append(targets, pair)
		}

		isTodayWindow := resetTodayUsageWindowIsToday(dailyWindowStart, dayStart, cutoff)
		deductUsage := resetUsage
		if isTodayWindow {
			result.TodayWindowSubscriptions++
			deductUsage = deductUsage.Add(positiveUsage)
		} else if positiveUsage.IsPositive() {
			staleAmount = staleAmount.Add(positiveUsage)
		}
		if deductUsage.IsPositive() {
			result.ResetSubscriptions++
			resetAmount = resetAmount.Add(deductUsage)
		}
		adjustments = append(adjustments, resetTodayUsageAdjustment{ID: id, DeductUSD: deductUsage})
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, fmt.Errorf("read reset-today subscriptions: %w", err)
	}
	result.ResetAndDeductedUSD, _ = resetAmount.Float64()
	result.StaleDailyClearedUSD, _ = staleAmount.Float64()
	return result, targets, adjustments, nil
}

func resetTodayUsageAdjustmentArrays(adjustments []resetTodayUsageAdjustment) ([]int64, []string) {
	ids := make([]int64, 0, len(adjustments))
	deductions := make([]string, 0, len(adjustments))
	for _, adjustment := range adjustments {
		ids = append(ids, adjustment.ID)
		deductions = append(deductions, adjustment.DeductUSD.String())
	}
	return ids, deductions
}

func completeResetTodayUsageOperation(ctx context.Context, client resetTodayUsageDB, operationID int64, result *ResetTodayUsageResult, targets []resetTodayUsageTarget) error {
	summaryJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("encode reset-today operation summary: %w", err)
	}
	targetsJSON, err := json.Marshal(targets)
	if err != nil {
		return fmt.Errorf("encode reset-today operation targets: %w", err)
	}
	updated, err := client.ExecContext(ctx, `
		UPDATE subscription_quota_reset_operations
		SET status = 'succeeded',
			cutoff_at = $2,
			captured_at = $3,
			summary = $4::jsonb,
			target_pairs = $5::jsonb,
			updated_at = NOW()
		WHERE id = $1 AND status = 'processing'
	`, operationID, result.CutoffAt, result.CapturedAt, string(summaryJSON), string(targetsJSON))
	if err != nil {
		return fmt.Errorf("complete reset-today operation: %w", err)
	}
	affected, err := updated.RowsAffected()
	if err != nil {
		return fmt.Errorf("count completed reset-today operation: %w", err)
	}
	if affected != 1 {
		return fmt.Errorf("complete reset-today operation affected %d rows", affected)
	}
	return nil
}

func (s *SubscriptionService) invalidateResetTodayUsageTargetBatch(targets []resetTodayUsageTarget, cutoff time.Time, logFailures bool) (int, []resetTodayUsageTarget) {
	return s.invalidateResetTodayUsageTargetBatchWithTimeouts(
		targets,
		cutoff,
		logFailures,
		resetTodayUsageCacheBatchTTL,
		subscriptionQuotaCacheTimeout,
	)
}

func (s *SubscriptionService) invalidateResetTodayUsageTargetBatchWithTimeouts(
	targets []resetTodayUsageTarget,
	cutoff time.Time,
	logFailures bool,
	batchTimeout time.Duration,
	targetTimeout time.Duration,
) (int, []resetTodayUsageTarget) {
	if len(targets) == 0 {
		return 0, nil
	}

	// Clear every local L1 entry before any remote operation can block.
	for _, target := range targets {
		s.InvalidateSubCache(target.UserID, target.GroupID)
	}
	if s.subCacheL1 != nil {
		s.subCacheL1.Wait()
	}

	type invalidationResult struct {
		target resetTodayUsageTarget
		err    error
	}
	workerCount := min(resetTodayUsageCacheWorkers, len(targets))
	jobs := make(chan resetTodayUsageTarget, len(targets))
	results := make(chan invalidationResult, len(targets))
	batchCtx, batchCancel := context.WithTimeout(context.Background(), batchTimeout)
	defer batchCancel()
	var workers sync.WaitGroup
	workers.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer workers.Done()
			for target := range jobs {
				if err := batchCtx.Err(); err != nil {
					results <- invalidationResult{target: target, err: fmt.Errorf("reset-today cache batch deadline: %w", err)}
					continue
				}
				ctx, cancel := context.WithTimeout(batchCtx, targetTimeout)
				err := s.invalidateSubscriptionResetRemoteCaches(ctx, target.UserID, target.GroupID, cutoff)
				cancel()
				results <- invalidationResult{target: target, err: err}
			}
		}()
	}
	for _, target := range targets {
		jobs <- target
	}
	close(jobs)
	workers.Wait()
	close(results)

	succeeded := 0
	failedTargets := make([]resetTodayUsageTarget, 0)
	for result := range results {
		if result.err == nil {
			succeeded++
			continue
		}
		failedTargets = append(failedTargets, result.target)
		if logFailures {
			log.Printf("Warning: reset-today cache invalidation failed for user %d group %d: %v", result.target.UserID, result.target.GroupID, result.err)
		}
	}
	return succeeded, failedTargets
}

func (s *SubscriptionService) scheduleResetTodayUsageCacheRetries(targets []resetTodayUsageTarget, cutoff time.Time) {
	if len(targets) == 0 || s.billingCacheService == nil {
		return
	}
	retryTargets := append([]resetTodayUsageTarget(nil), targets...)
	retryDelays := append([]time.Duration(nil), subscriptionQuotaCacheRetryDelays...)
	go func() {
		previousDelay := time.Duration(0)
		for _, delay := range retryDelays {
			timer := time.NewTimer(delay - previousDelay)
			<-timer.C
			previousDelay = delay
			_, failedTargets := s.invalidateResetTodayUsageTargetBatch(retryTargets, cutoff, false)
			if len(failedTargets) == 0 {
				return
			}
			log.Printf("Warning: reset-today follow-up cache invalidation had %d failures at delay %s", len(failedTargets), delay)
			retryTargets = failedTargets
		}
	}()
}

func (s *SubscriptionService) persistResetTodayUsageCacheSummary(operationID int64, result *ResetTodayUsageResult) {
	if s == nil || s.entClient == nil || result == nil {
		return
	}
	summaryJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("Warning: encode reset-today cache summary failed: %v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := s.entClient.ExecContext(ctx, `
		UPDATE subscription_quota_reset_operations
		SET summary = $2::jsonb, updated_at = NOW()
		WHERE id = $1 AND status = 'succeeded'
	`, operationID, string(summaryJSON)); err != nil {
		log.Printf("Warning: persist reset-today cache summary failed for operation %d: %v", operationID, err)
	}
}
