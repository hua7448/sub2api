package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/ent/predicate"
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
	"github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/ent/usersubscription"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type userSubscriptionRepository struct {
	client *dbent.Client
}

func NewUserSubscriptionRepository(client *dbent.Client) service.UserSubscriptionRepository {
	return &userSubscriptionRepository{client: client}
}

func (r *userSubscriptionRepository) Create(ctx context.Context, sub *service.UserSubscription) error {
	if sub == nil {
		return service.ErrSubscriptionNilInput
	}

	client := clientFromContext(ctx, r.client)
	builder := client.UserSubscription.Create().
		SetUserID(sub.UserID).
		SetGroupID(sub.GroupID).
		SetExpiresAt(sub.ExpiresAt).
		SetNillableDailyWindowStart(sub.DailyWindowStart).
		SetNillableWeeklyWindowStart(sub.WeeklyWindowStart).
		SetNillableMonthlyWindowStart(sub.MonthlyWindowStart).
		SetDailyUsageUsd(sub.DailyUsageUSD).
		SetWeeklyUsageUsd(sub.WeeklyUsageUSD).
		SetMonthlyUsageUsd(sub.MonthlyUsageUSD).
		SetNillableAssignedBy(sub.AssignedBy)

	if sub.StartsAt.IsZero() {
		builder.SetStartsAt(time.Now())
	} else {
		builder.SetStartsAt(sub.StartsAt)
	}
	if sub.Status != "" {
		builder.SetStatus(sub.Status)
	}
	if !sub.AssignedAt.IsZero() {
		builder.SetAssignedAt(sub.AssignedAt)
	}
	// Keep compatibility with historical behavior: always store notes as a string value.
	builder.SetNotes(sub.Notes)

	created, err := builder.Save(ctx)
	if err == nil {
		applyUserSubscriptionEntityToService(sub, created)
	}
	return translatePersistenceError(err, nil, service.ErrSubscriptionAlreadyExists)
}

func (r *userSubscriptionRepository) GetByID(ctx context.Context, id int64) (*service.UserSubscription, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.UserSubscription.Query().
		Where(usersubscription.IDEQ(id)).
		WithUser().
		WithGroup().
		WithAssignedByUser().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
	}
	return userSubscriptionEntityToService(m), nil
}

func (r *userSubscriptionRepository) GetByIDForUpdate(ctx context.Context, id int64) (*service.UserSubscription, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.UserSubscription.Query().
		Where(usersubscription.IDEQ(id)).
		ForUpdate().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
	}
	return userSubscriptionEntityToService(m), nil
}

func (r *userSubscriptionRepository) GetByIDIncludeDeleted(ctx context.Context, id int64) (*service.UserSubscription, error) {
	client := clientFromContext(ctx, r.client)
	queryCtx := mixins.SkipSoftDelete(ctx)
	m, err := client.UserSubscription.Query().
		Where(usersubscription.IDEQ(id)).
		WithUser().
		WithGroup().
		WithAssignedByUser().
		Only(queryCtx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
	}
	return userSubscriptionEntityToServicePreserveStatus(m), nil
}

func (r *userSubscriptionRepository) GetByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*service.UserSubscription, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.UserSubscription.Query().
		Where(usersubscription.UserIDEQ(userID), usersubscription.GroupIDEQ(groupID)).
		WithGroup().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
	}
	return userSubscriptionEntityToService(m), nil
}

func (r *userSubscriptionRepository) GetActiveByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (*service.UserSubscription, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.UserSubscription.Query().
		Where(
			usersubscription.UserIDEQ(userID),
			usersubscription.GroupIDEQ(groupID),
			usersubscription.StatusEQ(service.SubscriptionStatusActive),
			usersubscription.ExpiresAtGT(time.Now()),
		).
		WithGroup().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
	}
	return userSubscriptionEntityToService(m), nil
}

func (r *userSubscriptionRepository) Update(ctx context.Context, sub *service.UserSubscription) error {
	if sub == nil {
		return service.ErrSubscriptionNilInput
	}

	client := clientFromContext(ctx, r.client)
	builder := client.UserSubscription.UpdateOneID(sub.ID).
		SetUserID(sub.UserID).
		SetGroupID(sub.GroupID).
		SetStartsAt(sub.StartsAt).
		SetExpiresAt(sub.ExpiresAt).
		SetStatus(sub.Status).
		SetNillableDailyWindowStart(sub.DailyWindowStart).
		SetNillableWeeklyWindowStart(sub.WeeklyWindowStart).
		SetNillableMonthlyWindowStart(sub.MonthlyWindowStart).
		SetDailyUsageUsd(sub.DailyUsageUSD).
		SetWeeklyUsageUsd(sub.WeeklyUsageUSD).
		SetMonthlyUsageUsd(sub.MonthlyUsageUSD).
		SetNillableAssignedBy(sub.AssignedBy).
		SetAssignedAt(sub.AssignedAt).
		SetNotes(sub.Notes)

	updated, err := builder.Save(ctx)
	if err == nil {
		applyUserSubscriptionEntityToService(sub, updated)
		return nil
	}
	return translatePersistenceError(err, service.ErrSubscriptionNotFound, service.ErrSubscriptionAlreadyExists)
}

func (r *userSubscriptionRepository) Delete(ctx context.Context, id int64) error {
	// Match GORM semantics: deleting a missing row is not an error.
	client := clientFromContext(ctx, r.client)
	_, err := client.UserSubscription.Delete().Where(usersubscription.IDEQ(id)).Exec(ctx)
	return err
}

func (r *userSubscriptionRepository) Restore(ctx context.Context, subscriptionID int64, restoredStatus string) (*service.UserSubscription, error) {
	client := clientFromContext(ctx, r.client)
	queryCtx := mixins.SkipSoftDelete(ctx)
	_, err := client.UserSubscription.UpdateOneID(subscriptionID).
		SetStatus(restoredStatus).
		ClearDeletedAt().
		SetUpdatedAt(time.Now()).
		Save(queryCtx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrSubscriptionNotFound, service.ErrSubscriptionRestoreConflict)
	}
	return r.GetByID(ctx, subscriptionID)
}

func (r *userSubscriptionRepository) ListByUserID(ctx context.Context, userID int64) ([]service.UserSubscription, error) {
	client := clientFromContext(ctx, r.client)
	subs, err := client.UserSubscription.Query().
		Where(usersubscription.UserIDEQ(userID)).
		WithGroup().
		Order(dbent.Desc(usersubscription.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return userSubscriptionEntitiesToService(subs), nil
}

func (r *userSubscriptionRepository) ListActiveByUserID(ctx context.Context, userID int64) ([]service.UserSubscription, error) {
	client := clientFromContext(ctx, r.client)
	subs, err := client.UserSubscription.Query().
		Where(
			usersubscription.UserIDEQ(userID),
			usersubscription.StatusEQ(service.SubscriptionStatusActive),
			usersubscription.ExpiresAtGT(time.Now()),
		).
		WithGroup().
		Order(dbent.Desc(usersubscription.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return userSubscriptionEntitiesToService(subs), nil
}

func (r *userSubscriptionRepository) ListByGroupID(ctx context.Context, groupID int64, params pagination.PaginationParams) ([]service.UserSubscription, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.UserSubscription.Query().Where(usersubscription.GroupIDEQ(groupID))

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	subs, err := q.
		WithUser().
		WithGroup().
		Order(dbent.Desc(usersubscription.FieldCreatedAt)).
		Offset(params.Offset()).
		Limit(params.Limit()).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	return userSubscriptionEntitiesToService(subs), paginationResultFromTotal(int64(total), params), nil
}

func (r *userSubscriptionRepository) List(ctx context.Context, params pagination.PaginationParams, userID, groupID *int64, status, platform, sortBy, sortOrder string) ([]service.UserSubscription, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.UserSubscription.Query()
	includeSoftDeleted := status == "" || status == service.SubscriptionStatusRevoked
	if userID != nil {
		q = q.Where(usersubscription.UserIDEQ(*userID))
	}
	if groupID != nil {
		q = q.Where(usersubscription.GroupIDEQ(*groupID))
	}
	if platform != "" {
		groupPredicates := []predicate.Group{group.PlatformEQ(platform)}
		if includeSoftDeleted {
			groupPredicates = append(groupPredicates, group.DeletedAtIsNil())
		}
		q = q.Where(usersubscription.HasGroupWith(groupPredicates...))
	}

	// Status filtering with real-time expiration check
	now := time.Now()
	switch status {
	case service.SubscriptionStatusActive:
		// Active: status is active AND not yet expired
		q = q.Where(
			usersubscription.StatusEQ(service.SubscriptionStatusActive),
			usersubscription.ExpiresAtGT(now),
		)
	case service.SubscriptionStatusExpired:
		// Expired: status is expired OR (status is active but already expired)
		q = q.Where(
			usersubscription.Or(
				usersubscription.StatusEQ(service.SubscriptionStatusExpired),
				usersubscription.And(
					usersubscription.StatusEQ(service.SubscriptionStatusActive),
					usersubscription.ExpiresAtLTE(now),
				),
			),
		)
	case service.SubscriptionStatusRevoked:
		// Revoked is a DTO/API display state backed by user_subscriptions.deleted_at.
		q = q.Where(usersubscription.DeletedAtNotNil())
	case "":
		// No filter. Use SkipSoftDelete below so admin "all status" includes revoked history.
	default:
		// Other persisted status.
		q = q.Where(usersubscription.StatusEQ(status))
	}

	queryCtx := ctx
	if includeSoftDeleted {
		queryCtx = mixins.SkipSoftDelete(ctx)
	}

	total, err := q.Clone().Count(queryCtx)
	if err != nil {
		return nil, nil, err
	}

	if !includeSoftDeleted {
		q = q.WithUser().WithGroup().WithAssignedByUser()
	}

	// Determine sort field
	var field string
	switch sortBy {
	case "expires_at":
		field = usersubscription.FieldExpiresAt
	case "status":
		field = usersubscription.FieldStatus
	default:
		field = usersubscription.FieldCreatedAt
	}

	// Determine sort order (default: desc)
	if sortOrder == "asc" && sortBy != "" {
		q = q.Order(dbent.Asc(field))
	} else {
		q = q.Order(dbent.Desc(field))
	}

	subs, err := q.
		Offset(params.Offset()).
		Limit(params.Limit()).
		All(queryCtx)
	if err != nil {
		return nil, nil, err
	}

	result := userSubscriptionEntitiesToService(subs)
	if includeSoftDeleted {
		if err := r.attachUserSubscriptionRelations(ctx, result); err != nil {
			return nil, nil, err
		}
	}

	return result, paginationResultFromTotal(int64(total), params), nil
}

func (r *userSubscriptionRepository) ExistsByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (bool, error) {
	client := clientFromContext(ctx, r.client)
	return client.UserSubscription.Query().
		Where(usersubscription.UserIDEQ(userID), usersubscription.GroupIDEQ(groupID)).
		Exist(ctx)
}

func (r *userSubscriptionRepository) ExistsActiveByUserIDAndGroupID(ctx context.Context, userID, groupID int64) (bool, error) {
	return r.ExistsByUserIDAndGroupID(ctx, userID, groupID)
}

func (r *userSubscriptionRepository) ExtendExpiry(ctx context.Context, subscriptionID int64, newExpiresAt time.Time) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.UserSubscription.UpdateOneID(subscriptionID).
		SetExpiresAt(newExpiresAt).
		Save(ctx)
	return translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
}

func (r *userSubscriptionRepository) UpdateStatus(ctx context.Context, subscriptionID int64, status string) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.UserSubscription.UpdateOneID(subscriptionID).
		SetStatus(status).
		Save(ctx)
	return translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
}

func (r *userSubscriptionRepository) UpdateNotes(ctx context.Context, subscriptionID int64, notes string) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.UserSubscription.UpdateOneID(subscriptionID).
		SetNotes(notes).
		Save(ctx)
	return translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
}

func (r *userSubscriptionRepository) ActivateWindows(ctx context.Context, id int64, dailyStart, periodicStart time.Time) error {
	client := clientFromContext(ctx, r.client)
	n, err := client.UserSubscription.Update().
		Where(
			usersubscription.IDEQ(id),
			usersubscription.DailyWindowStartIsNil(),
			usersubscription.WeeklyWindowStartIsNil(),
			usersubscription.MonthlyWindowStartIsNil(),
		).
		SetDailyWindowStart(dailyStart).
		SetWeeklyWindowStart(periodicStart).
		SetMonthlyWindowStart(periodicStart).
		Save(ctx)
	return r.translateConditionalWindowReset(ctx, client, id, n, err)
}

// ResetDailyQuotaForUser atomically verifies ownership and effectiveness, then
// starts a fresh daily window without changing weekly or monthly usage.
func (r *userSubscriptionRepository) ResetDailyQuotaForUser(ctx context.Context, userID, subscriptionID int64, operationKeyHash string) (*service.UserDailyQuotaResetResult, error) {
	client := clientFromContext(ctx, r.client)
	rows, err := client.QueryContext(ctx, `
		WITH reset_clock AS MATERIALIZED (
			SELECT clock_timestamp() AS captured_at
		),
		owned_subscription AS MATERIALIZED (
			SELECT us.id,
				us.group_id,
				us.deleted_at AS subscription_deleted_at,
				us.status,
				us.starts_at,
				us.expires_at,
				us.weekly_window_start,
				us.daily_quota_reset_week_start,
				CASE
					WHEN jsonb_typeof(us.daily_quota_reset_operations) = 'array'
					THEN us.daily_quota_reset_operations
					ELSE '[]'::jsonb
				END AS daily_quota_reset_operations,
				us.daily_usage_usd::float8 AS daily_usage_usd,
				us.weekly_usage_usd::float8 AS weekly_usage_usd,
				us.monthly_usage_usd::float8 AS monthly_usage_usd,
				g.deleted_at AS group_deleted_at,
				g.daily_limit_usd,
				reset_clock.captured_at
			FROM user_subscriptions AS us
			LEFT JOIN groups AS g ON g.id = us.group_id
			CROSS JOIN reset_clock
			WHERE us.id = $1
				AND us.user_id = $2
			FOR UPDATE OF us
		),
		operation_replay AS MATERIALIZED (
			SELECT
				(operation.value->>'subscription_id')::bigint AS id,
				(operation.value->>'group_id')::bigint AS group_id,
				(operation.value->>'reset_at')::timestamptz AS reset_at,
				(operation.value->>'daily_usage_usd')::float8 AS daily_usage_usd,
				(operation.value->>'weekly_usage_usd')::float8 AS weekly_usage_usd,
				(operation.value->>'monthly_usage_usd')::float8 AS monthly_usage_usd,
				(operation.value->>'daily_quota_reset_available_at')::timestamptz AS available_at
			FROM owned_subscription
			CROSS JOIN LATERAL jsonb_array_elements(
				owned_subscription.daily_quota_reset_operations
			) AS operation(value)
			WHERE operation.value->>'idempotency_key_hash' = $4
				AND operation.value ?& ARRAY[
					'subscription_id', 'group_id', 'reset_at', 'daily_usage_usd',
					'weekly_usage_usd', 'monthly_usage_usd', 'replay_expires_at'
				]
				AND (operation.value->>'replay_expires_at')::timestamptz
					> owned_subscription.captured_at
			ORDER BY (operation.value->>'reset_at')::timestamptz DESC
			LIMIT 1
		),
		candidate AS MATERIALIZED (
			SELECT owned_subscription.*,
				CASE
					WHEN owned_subscription.weekly_window_start IS NULL
					THEN owned_subscription.captured_at
					ELSE weekly_anchor.anchor_at
						+ GREATEST(
							FLOOR(EXTRACT(EPOCH FROM (owned_subscription.captured_at - weekly_anchor.anchor_at)) / 604800),
							0
						)::bigint * interval '168 hours'
				END AS current_week_start
			FROM owned_subscription
			CROSS JOIN LATERAL (
				SELECT CASE
					WHEN owned_subscription.weekly_window_start =
						date_trunc('day', owned_subscription.starts_at)
						AND date_trunc('day', owned_subscription.starts_at) < owned_subscription.starts_at
					THEN owned_subscription.starts_at
					ELSE owned_subscription.weekly_window_start
				END AS anchor_at
			) AS weekly_anchor
			WHERE NOT EXISTS (SELECT 1 FROM operation_replay)
				AND owned_subscription.subscription_deleted_at IS NULL
				AND owned_subscription.status = $3
				AND owned_subscription.starts_at <= owned_subscription.captured_at
				AND owned_subscription.expires_at > owned_subscription.captured_at
				AND owned_subscription.group_deleted_at IS NULL
				AND owned_subscription.daily_limit_usd > 0
		),
		reset AS (
			UPDATE user_subscriptions AS us
			SET daily_usage_usd = 0,
				daily_window_start = GREATEST(candidate.captured_at, us.updated_at + interval '1 microsecond'),
				weekly_window_start = COALESCE(us.weekly_window_start, candidate.current_week_start),
				monthly_window_start = COALESCE(us.monthly_window_start, candidate.captured_at),
				daily_quota_reset_week_start = candidate.current_week_start,
				daily_quota_reset_operations = COALESCE((
					SELECT jsonb_agg(operation.value ORDER BY (operation.value->>'reset_at')::timestamptz)
					FROM jsonb_array_elements(candidate.daily_quota_reset_operations) AS operation(value)
					WHERE operation.value ? 'replay_expires_at'
						AND (operation.value->>'replay_expires_at')::timestamptz > candidate.captured_at
				), '[]'::jsonb) || jsonb_build_array(
					jsonb_build_object(
						'idempotency_key_hash', $4,
						'subscription_id', us.id,
						'group_id', us.group_id,
						'reset_at', GREATEST(candidate.captured_at, us.updated_at + interval '1 microsecond'),
						'cleared_daily_window_start', COALESCE(us.daily_window_start, date_trunc('day', candidate.captured_at AT TIME ZONE 'Asia/Shanghai') AT TIME ZONE 'Asia/Shanghai'),
						'cleared_daily_usage_usd', GREATEST(us.daily_usage_usd, 0),
						'daily_usage_usd', 0,
						'weekly_usage_usd', us.weekly_usage_usd,
						'monthly_usage_usd', us.monthly_usage_usd,
						'daily_quota_reset_available_at', CASE
							WHEN candidate.current_week_start + interval '168 hours' < candidate.expires_at
							THEN candidate.current_week_start + interval '168 hours'
						END,
						'replay_expires_at', candidate.captured_at + ($5::bigint * interval '1 second')
					)
				),
				updated_at = GREATEST(candidate.captured_at, us.updated_at + interval '1 microsecond')
			FROM candidate
			WHERE us.id = candidate.id
				AND (
					us.daily_quota_reset_week_start IS NULL
					OR us.daily_quota_reset_week_start < candidate.current_week_start
				)
			RETURNING us.id, us.group_id, us.updated_at,
				us.daily_usage_usd::float8,
				us.weekly_usage_usd::float8,
				us.monthly_usage_usd::float8
		)
		SELECT TRUE, FALSE, reset.id, reset.group_id, reset.updated_at,
			reset.daily_usage_usd, reset.weekly_usage_usd, reset.monthly_usage_usd,
			CASE
				WHEN candidate.current_week_start + interval '168 hours' < candidate.expires_at
				THEN candidate.current_week_start + interval '168 hours'
			END AS available_at
		FROM reset
		JOIN candidate ON candidate.id = reset.id
		UNION ALL
		SELECT TRUE, TRUE, operation_replay.id, operation_replay.group_id, operation_replay.reset_at,
			operation_replay.daily_usage_usd, operation_replay.weekly_usage_usd, operation_replay.monthly_usage_usd,
			operation_replay.available_at
		FROM operation_replay
		UNION ALL
		SELECT FALSE, FALSE, candidate.id, candidate.group_id, candidate.captured_at,
			candidate.daily_usage_usd, candidate.weekly_usage_usd, candidate.monthly_usage_usd,
			CASE
				WHEN candidate.daily_quota_reset_week_start > candidate.current_week_start THEN NULL
				WHEN candidate.current_week_start + interval '168 hours' < candidate.expires_at
				THEN candidate.current_week_start + interval '168 hours'
			END AS available_at
		FROM candidate
		WHERE NOT EXISTS (SELECT 1 FROM reset)
			AND NOT EXISTS (SELECT 1 FROM operation_replay)
	`,
		subscriptionID,
		userID,
		service.SubscriptionStatusActive,
		operationKeyHash,
		int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := new(service.UserDailyQuotaResetResult)
	var didReset bool
	var operationReplayed bool
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, service.ErrSubscriptionNotFound
	}
	if err := rows.Scan(
		&didReset,
		&operationReplayed,
		&result.SubscriptionID,
		&result.GroupID,
		&result.ResetAt,
		&result.DailyUsageUSD,
		&result.WeeklyUsageUSD,
		&result.MonthlyUsageUSD,
		&result.DailyQuotaResetAvailableAt,
	); err != nil {
		return nil, err
	}
	if !didReset {
		return nil, service.NewDailyQuotaResetWeeklyLimitError(result.DailyQuotaResetAvailableAt)
	}
	result.DailyQuotaResetAvailable = false
	result.OperationReplayed = operationReplayed
	return result, nil
}

func (r *userSubscriptionRepository) ResetUsageWindows(ctx context.Context, id int64, resetDaily, resetWeekly, resetMonthly bool, dailyStart, periodicStart time.Time) (time.Time, error) {
	client := clientFromContext(ctx, r.client)
	rows, err := client.QueryContext(ctx, `
		UPDATE user_subscriptions AS us
		SET daily_usage_usd = CASE WHEN $2::boolean THEN 0 ELSE us.daily_usage_usd END,
			daily_window_start = CASE WHEN $2::boolean THEN $5::timestamptz ELSE us.daily_window_start END,
			weekly_usage_usd = CASE WHEN $3::boolean THEN 0 ELSE us.weekly_usage_usd END,
			weekly_window_start = CASE WHEN $3::boolean THEN $6::timestamptz ELSE us.weekly_window_start END,
			monthly_usage_usd = CASE WHEN $4::boolean THEN 0 ELSE us.monthly_usage_usd END,
			monthly_window_start = CASE WHEN $4::boolean THEN $6::timestamptz ELSE us.monthly_window_start END,
			updated_at = GREATEST(clock_timestamp(), us.updated_at + interval '1 microsecond')
		WHERE us.id = $1
			AND us.deleted_at IS NULL
		RETURNING us.updated_at
	`, id, resetDaily, resetWeekly, resetMonthly, dailyStart, periodicStart)
	if err != nil {
		return time.Time{}, translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return time.Time{}, err
		}
		return time.Time{}, translatePersistenceError(sql.ErrNoRows, service.ErrSubscriptionNotFound, nil)
	}
	var cutoff time.Time
	if err := rows.Scan(&cutoff); err != nil {
		return time.Time{}, err
	}
	if rows.Next() {
		return time.Time{}, fmt.Errorf("reset usage windows returned more than one cutoff row")
	}
	if err := rows.Err(); err != nil {
		return time.Time{}, err
	}
	if err := rows.Close(); err != nil {
		return time.Time{}, err
	}
	return cutoff, nil
}

func (r *userSubscriptionRepository) ResetDailyUsage(ctx context.Context, id int64, expectedWindowStart *time.Time, newWindowStart time.Time) error {
	client := clientFromContext(ctx, r.client)
	query := client.UserSubscription.Update().Where(usersubscription.IDEQ(id))
	if expectedWindowStart == nil {
		query = query.Where(usersubscription.DailyWindowStartIsNil())
	} else {
		query = query.Where(usersubscription.DailyWindowStartEQ(*expectedWindowStart))
	}
	n, err := query.
		SetDailyUsageUsd(0).
		SetDailyWindowStart(newWindowStart).
		Save(ctx)
	return r.translateConditionalWindowReset(ctx, client, id, n, err)
}

func (r *userSubscriptionRepository) ResetWeeklyUsage(ctx context.Context, id int64, expectedWindowStart *time.Time, newWindowStart time.Time) error {
	client := clientFromContext(ctx, r.client)
	query := client.UserSubscription.Update().Where(usersubscription.IDEQ(id))
	if expectedWindowStart == nil {
		query = query.Where(usersubscription.WeeklyWindowStartIsNil())
	} else {
		query = query.Where(usersubscription.WeeklyWindowStartEQ(*expectedWindowStart))
	}
	n, err := query.
		SetWeeklyUsageUsd(0).
		SetWeeklyWindowStart(newWindowStart).
		Save(ctx)
	return r.translateConditionalWindowReset(ctx, client, id, n, err)
}

func (r *userSubscriptionRepository) ResetMonthlyUsage(ctx context.Context, id int64, expectedWindowStart *time.Time, newWindowStart time.Time) error {
	client := clientFromContext(ctx, r.client)
	query := client.UserSubscription.Update().Where(usersubscription.IDEQ(id))
	if expectedWindowStart == nil {
		query = query.Where(usersubscription.MonthlyWindowStartIsNil())
	} else {
		query = query.Where(usersubscription.MonthlyWindowStartEQ(*expectedWindowStart))
	}
	n, err := query.
		SetMonthlyUsageUsd(0).
		SetMonthlyWindowStart(newWindowStart).
		Save(ctx)
	return r.translateConditionalWindowReset(ctx, client, id, n, err)
}

func (r *userSubscriptionRepository) translateConditionalWindowReset(ctx context.Context, client *dbent.Client, id int64, affected int, err error) error {
	if err != nil {
		return translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
	}
	if affected > 0 {
		return nil
	}

	// A stale reset is an expected no-op: another request already advanced the
	// window. Preserve not-found semantics for callers that target a missing row.
	exists, err := client.UserSubscription.Query().Where(usersubscription.IDEQ(id)).Exist(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrSubscriptionNotFound, nil)
	}
	if !exists {
		return service.ErrSubscriptionNotFound
	}
	return nil
}

// IncrementUsage 原子性地累加订阅用量。
// 限额检查已在请求前由 BillingCacheService.CheckBillingEligibility 完成，
// 此处仅负责记录实际消费，确保消费数据的完整性。
func (r *userSubscriptionRepository) IncrementUsage(ctx context.Context, id int64, costUSD float64) error {
	const updateSQL = `
		UPDATE user_subscriptions us
		SET
			daily_usage_usd = us.daily_usage_usd + $1,
			weekly_usage_usd = us.weekly_usage_usd + $1,
			monthly_usage_usd = us.monthly_usage_usd + $1,
			updated_at = NOW()
		FROM groups g
		WHERE us.id = $2
			AND us.deleted_at IS NULL
			AND us.group_id = g.id
			AND g.deleted_at IS NULL
	`

	client := clientFromContext(ctx, r.client)
	result, err := client.ExecContext(ctx, updateSQL, costUSD, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected > 0 {
		return nil
	}

	// affected == 0：订阅不存在或已删除
	return service.ErrSubscriptionNotFound
}

func (r *userSubscriptionRepository) BatchUpdateExpiredStatus(ctx context.Context) (int64, error) {
	client := clientFromContext(ctx, r.client)
	n, err := client.UserSubscription.Update().
		Where(
			usersubscription.StatusEQ(service.SubscriptionStatusActive),
			usersubscription.ExpiresAtLTE(time.Now()),
		).
		SetStatus(service.SubscriptionStatusExpired).
		Save(ctx)
	return int64(n), err
}

// Extra repository helpers (currently used only by integration tests).

func (r *userSubscriptionRepository) ListExpired(ctx context.Context) ([]service.UserSubscription, error) {
	client := clientFromContext(ctx, r.client)
	subs, err := client.UserSubscription.Query().
		Where(
			usersubscription.StatusEQ(service.SubscriptionStatusActive),
			usersubscription.ExpiresAtLTE(time.Now()),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return userSubscriptionEntitiesToService(subs), nil
}

func (r *userSubscriptionRepository) CountByGroupID(ctx context.Context, groupID int64) (int64, error) {
	client := clientFromContext(ctx, r.client)
	count, err := client.UserSubscription.Query().Where(usersubscription.GroupIDEQ(groupID)).Count(ctx)
	return int64(count), err
}

func (r *userSubscriptionRepository) CountActiveByGroupID(ctx context.Context, groupID int64) (int64, error) {
	client := clientFromContext(ctx, r.client)
	count, err := client.UserSubscription.Query().
		Where(
			usersubscription.GroupIDEQ(groupID),
			usersubscription.StatusEQ(service.SubscriptionStatusActive),
			usersubscription.ExpiresAtGT(time.Now()),
		).
		Count(ctx)
	return int64(count), err
}

func (r *userSubscriptionRepository) DeleteByGroupID(ctx context.Context, groupID int64) (int64, error) {
	client := clientFromContext(ctx, r.client)
	n, err := client.UserSubscription.Delete().Where(usersubscription.GroupIDEQ(groupID)).Exec(ctx)
	return int64(n), err
}

func (r *userSubscriptionRepository) attachUserSubscriptionRelations(ctx context.Context, subs []service.UserSubscription) error {
	if len(subs) == 0 {
		return nil
	}

	userIDs := make([]int64, 0, len(subs))
	groupIDs := make([]int64, 0, len(subs))
	assignedByIDs := make([]int64, 0, len(subs))
	for i := range subs {
		userIDs = append(userIDs, subs[i].UserID)
		groupIDs = append(groupIDs, subs[i].GroupID)
		if subs[i].AssignedBy != nil {
			assignedByIDs = append(assignedByIDs, *subs[i].AssignedBy)
		}
	}

	client := clientFromContext(ctx, r.client)
	users, err := client.User.Query().Where(user.IDIn(uniqueInt64s(userIDs)...)).All(ctx)
	if err != nil {
		return err
	}
	userByID := make(map[int64]*service.User, len(users))
	for _, u := range users {
		userByID[u.ID] = userEntityToService(u)
	}

	groups, err := client.Group.Query().Where(group.IDIn(uniqueInt64s(groupIDs)...)).All(ctx)
	if err != nil {
		return err
	}
	groupByID := make(map[int64]*service.Group, len(groups))
	for _, g := range groups {
		groupByID[g.ID] = groupEntityToService(g)
	}

	assignedByID := map[int64]*service.User{}
	if len(assignedByIDs) > 0 {
		assignedUsers, err := client.User.Query().Where(user.IDIn(uniqueInt64s(assignedByIDs)...)).All(ctx)
		if err != nil {
			return err
		}
		assignedByID = make(map[int64]*service.User, len(assignedUsers))
		for _, u := range assignedUsers {
			assignedByID[u.ID] = userEntityToService(u)
		}
	}

	for i := range subs {
		subs[i].User = userByID[subs[i].UserID]
		subs[i].Group = groupByID[subs[i].GroupID]
		if subs[i].AssignedBy != nil {
			subs[i].AssignedByUser = assignedByID[*subs[i].AssignedBy]
		}
	}
	return nil
}

func uniqueInt64s(values []int64) []int64 {
	seen := make(map[int64]struct{}, len(values))
	out := make([]int64, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func userSubscriptionEntityToService(m *dbent.UserSubscription) *service.UserSubscription {
	return userSubscriptionEntityToServiceWithStatusMapping(m, true)
}

func userSubscriptionEntityToServicePreserveStatus(m *dbent.UserSubscription) *service.UserSubscription {
	return userSubscriptionEntityToServiceWithStatusMapping(m, false)
}

func userSubscriptionEntityToServiceWithStatusMapping(m *dbent.UserSubscription, mapDeletedToRevoked bool) *service.UserSubscription {
	if m == nil {
		return nil
	}
	status := m.Status
	if mapDeletedToRevoked && m.DeletedAt != nil {
		status = service.SubscriptionStatusRevoked
	}
	out := &service.UserSubscription{
		ID:                       m.ID,
		UserID:                   m.UserID,
		GroupID:                  m.GroupID,
		StartsAt:                 m.StartsAt,
		ExpiresAt:                m.ExpiresAt,
		Status:                   status,
		DailyWindowStart:         m.DailyWindowStart,
		WeeklyWindowStart:        m.WeeklyWindowStart,
		MonthlyWindowStart:       m.MonthlyWindowStart,
		DailyQuotaResetWeekStart: m.DailyQuotaResetWeekStart,
		DailyUsageUSD:            m.DailyUsageUsd,
		WeeklyUsageUSD:           m.WeeklyUsageUsd,
		MonthlyUsageUSD:          m.MonthlyUsageUsd,
		AssignedBy:               m.AssignedBy,
		AssignedAt:               m.AssignedAt,
		Notes:                    derefString(m.Notes),
		CreatedAt:                m.CreatedAt,
		UpdatedAt:                m.UpdatedAt,
		DeletedAt:                m.DeletedAt,
	}
	if m.Edges.User != nil {
		out.User = userEntityToService(m.Edges.User)
	}
	if m.Edges.Group != nil {
		out.Group = groupEntityToService(m.Edges.Group)
	}
	if m.Edges.AssignedByUser != nil {
		out.AssignedByUser = userEntityToService(m.Edges.AssignedByUser)
	}
	return out
}

func userSubscriptionEntitiesToService(models []*dbent.UserSubscription) []service.UserSubscription {
	out := make([]service.UserSubscription, 0, len(models))
	for i := range models {
		if s := userSubscriptionEntityToService(models[i]); s != nil {
			out = append(out, *s)
		}
	}
	return out
}

func applyUserSubscriptionEntityToService(dst *service.UserSubscription, src *dbent.UserSubscription) {
	if dst == nil || src == nil {
		return
	}
	dst.ID = src.ID
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}
