package repository

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

func newUserDailyQuotaResetRepo(t *testing.T, capturedSQL *string) (*userSubscriptionRepository, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(captureEntQueryMatcher{actual: capturedSQL}))
	require.NoError(t, err)
	driver := entsql.OpenDB(dialect.Postgres, db)
	client := dbent.NewClient(dbent.Driver(driver))
	t.Cleanup(func() { _ = client.Close() })
	return &userSubscriptionRepository{client: client}, mock
}

func TestResetDailyQuotaForUserUsesAtomicOwnedEffectiveUpdate(t *testing.T) {
	var capturedSQL string
	repo, mock := newUserDailyQuotaResetRepo(t, &capturedSQL)
	resetAt := time.Date(2026, 8, 5, 4, 12, 13, 456000000, time.UTC)
	availableAt := resetAt.Add(7 * 24 * time.Hour)

	mock.ExpectQuery("reset daily quota").
		WithArgs(int64(501), int64(77), service.SubscriptionStatusActive, "operation-hash", int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second)).
		WillReturnRows(sqlmock.NewRows([]string{
			"did_reset", "operation_replayed", "id", "group_id", "updated_at", "daily_usage_usd", "weekly_usage_usd", "monthly_usage_usd", "available_at",
		}).AddRow(true, false, int64(501), int64(19), resetAt, 0.0, 12.5, 31.75, availableAt))

	result, err := repo.ResetDailyQuotaForUser(context.Background(), 77, 501, "operation-hash")

	require.NoError(t, err)
	require.Equal(t, int64(501), result.SubscriptionID)
	require.Equal(t, int64(19), result.GroupID)
	require.Equal(t, resetAt, result.ResetAt)
	require.Zero(t, result.DailyUsageUSD)
	require.Equal(t, 12.5, result.WeeklyUsageUSD)
	require.Equal(t, 31.75, result.MonthlyUsageUSD)
	require.False(t, result.DailyQuotaResetAvailable)
	require.Equal(t, availableAt, *result.DailyQuotaResetAvailableAt)
	require.False(t, result.OperationReplayed)
	require.NoError(t, mock.ExpectationsWereMet())

	sqlText := strings.ToUpper(normalizeSQLWhitespace(capturedSQL))
	require.Contains(t, sqlText, "SET DAILY_USAGE_USD = 0")
	require.Contains(t, sqlText, "DAILY_WINDOW_START = GREATEST")
	require.Contains(t, sqlText, "UPDATED_AT = GREATEST")
	require.Contains(t, sqlText, "US.ID = $1")
	require.Contains(t, sqlText, "US.USER_ID = $2")
	require.Contains(t, sqlText, "OWNED_SUBSCRIPTION.SUBSCRIPTION_DELETED_AT IS NULL")
	require.Contains(t, sqlText, "OWNED_SUBSCRIPTION.STATUS = $3")
	require.Contains(t, sqlText, "OWNED_SUBSCRIPTION.STARTS_AT <= OWNED_SUBSCRIPTION.CAPTURED_AT")
	require.Contains(t, sqlText, "OWNED_SUBSCRIPTION.EXPIRES_AT > OWNED_SUBSCRIPTION.CAPTURED_AT")
	require.Contains(t, sqlText, "FOR UPDATE OF US")
	require.Contains(t, sqlText, "DAILY_QUOTA_RESET_WEEK_START = CANDIDATE.CURRENT_WEEK_START")
	require.Contains(t, sqlText, "DAILY_QUOTA_RESET_OPERATIONS = COALESCE")
	require.Contains(t, sqlText, "JSONB_BUILD_ARRAY")
	require.Contains(t, sqlText, "'IDEMPOTENCY_KEY_HASH', $4")
	require.Contains(t, sqlText, "'CLEARED_DAILY_WINDOW_START', COALESCE(US.DAILY_WINDOW_START, DATE_TRUNC('DAY', CANDIDATE.CAPTURED_AT AT TIME ZONE 'ASIA/SHANGHAI') AT TIME ZONE 'ASIA/SHANGHAI')")
	require.Contains(t, sqlText, "'CLEARED_DAILY_USAGE_USD', GREATEST(US.DAILY_USAGE_USD, 0)")
	require.Contains(t, sqlText, "'REPLAY_EXPIRES_AT'")
	require.Contains(t, sqlText, "OPERATION_REPLAY AS MATERIALIZED")
	require.Less(t, strings.Index(sqlText, "OPERATION_REPLAY AS MATERIALIZED"), strings.Index(sqlText, "CANDIDATE AS MATERIALIZED"))
	require.Contains(t, sqlText, "DAILY_QUOTA_RESET_WEEK_START IS NULL")
	require.Contains(t, sqlText, "DAILY_QUOTA_RESET_WEEK_START < CANDIDATE.CURRENT_WEEK_START")
	require.NotContains(t, sqlText, "DAILY_QUOTA_RESET_WEEK_START IS DISTINCT FROM")
	require.Contains(t, sqlText, "WEEKLY_WINDOW_START = COALESCE")
	require.Contains(t, sqlText, "MONTHLY_WINDOW_START = COALESCE")
	require.Contains(t, sqlText, "INTERVAL '168 HOURS'")
	require.NotContains(t, sqlText, "INTERVAL '7 DAYS'")
	require.Contains(t, sqlText, "DATE_TRUNC('DAY', OWNED_SUBSCRIPTION.STARTS_AT)")
	require.NotContains(t, sqlText, "TIME ZONE 'UTC'")
	require.Contains(t, sqlText, "OWNED_SUBSCRIPTION.DAILY_LIMIT_USD > 0")
	require.NotContains(t, sqlText, "SET WEEKLY_USAGE_USD")
	require.NotContains(t, sqlText, "SET MONTHLY_USAGE_USD")
	require.NotContains(t, sqlText, "UPDATE USAGE_LOGS")
}

func TestResetDailyQuotaForUserNullWindowWritesShanghaiNativeLedgerWindow(t *testing.T) {
	var capturedSQL string
	repo, mock := newUserDailyQuotaResetRepo(t, &capturedSQL)
	resetAt := time.Date(2026, 8, 18, 23, 31, 2, 654322000, time.UTC)

	mock.ExpectQuery("reset daily quota with null window").
		WithArgs(int64(501), int64(77), service.SubscriptionStatusActive, "null-window-operation", int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second)).
		WillReturnRows(sqlmock.NewRows([]string{
			"did_reset", "operation_replayed", "id", "group_id", "updated_at", "daily_usage_usd", "weekly_usage_usd", "monthly_usage_usd", "available_at",
		}).AddRow(true, false, int64(501), int64(19), resetAt, 0.0, 0.0, 0.0, nil))

	result, err := repo.ResetDailyQuotaForUser(context.Background(), 77, 501, "null-window-operation")

	require.NoError(t, err)
	require.Equal(t, resetAt, result.ResetAt)
	require.NoError(t, mock.ExpectationsWereMet())

	sqlText := strings.ToUpper(normalizeSQLWhitespace(capturedSQL))
	require.Contains(t, sqlText, "'CLEARED_DAILY_WINDOW_START', COALESCE(US.DAILY_WINDOW_START, DATE_TRUNC('DAY', CANDIDATE.CAPTURED_AT AT TIME ZONE 'ASIA/SHANGHAI') AT TIME ZONE 'ASIA/SHANGHAI')")
	require.NotContains(t, sqlText, "'CLEARED_DAILY_WINDOW_START', US.DAILY_WINDOW_START,")
}

func TestResetDailyQuotaForUserHidesMissingIneligibleOrForeignSubscription(t *testing.T) {
	var capturedSQL string
	repo, mock := newUserDailyQuotaResetRepo(t, &capturedSQL)
	mock.ExpectQuery("reset daily quota").
		WithArgs(int64(999), int64(77), service.SubscriptionStatusActive, "operation-hash", int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second)).
		WillReturnRows(sqlmock.NewRows([]string{
			"did_reset", "operation_replayed", "id", "group_id", "updated_at", "daily_usage_usd", "weekly_usage_usd", "monthly_usage_usd", "available_at",
		}))

	result, err := repo.ResetDailyQuotaForUser(context.Background(), 77, 999, "operation-hash")

	require.Nil(t, result)
	require.True(t, errors.Is(err, service.ErrSubscriptionNotFound))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetDailyQuotaForUserReturnsWeeklyLimitWithAvailability(t *testing.T) {
	var capturedSQL string
	repo, mock := newUserDailyQuotaResetRepo(t, &capturedSQL)
	resetAt := time.Date(2026, 8, 5, 4, 12, 13, 0, time.UTC)
	availableAt := resetAt.Add(7 * 24 * time.Hour)

	mock.ExpectQuery("reset daily quota").
		WithArgs(int64(501), int64(77), service.SubscriptionStatusActive, "new-operation-hash", int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second)).
		WillReturnRows(sqlmock.NewRows([]string{
			"did_reset", "operation_replayed", "id", "group_id", "updated_at", "daily_usage_usd", "weekly_usage_usd", "monthly_usage_usd", "available_at",
		}).AddRow(false, false, int64(501), int64(19), resetAt, 2.5, 12.5, 31.75, availableAt))

	result, err := repo.ResetDailyQuotaForUser(context.Background(), 77, 501, "new-operation-hash")

	require.Nil(t, result)
	var appErr *infraerrors.ApplicationError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, "DAILY_QUOTA_RESET_WEEKLY_LIMIT", appErr.Reason)
	require.Equal(t, availableAt.Format(time.RFC3339Nano), appErr.Metadata["available_at"])
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetDailyQuotaForUserReplaysDurableWeeklyOperation(t *testing.T) {
	var capturedSQL string
	repo, mock := newUserDailyQuotaResetRepo(t, &capturedSQL)
	resetAt := time.Date(2026, 8, 5, 4, 12, 13, 0, time.UTC)
	availableAt := resetAt.Add(7 * 24 * time.Hour)

	mock.ExpectQuery("reset daily quota").
		WithArgs(int64(501), int64(77), service.SubscriptionStatusActive, "operation-hash", int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second)).
		WillReturnRows(sqlmock.NewRows([]string{
			"did_reset", "operation_replayed", "id", "group_id", "updated_at", "daily_usage_usd", "weekly_usage_usd", "monthly_usage_usd", "available_at",
		}).AddRow(true, true, int64(501), int64(19), resetAt, 0.0, 12.5, 31.75, availableAt))

	result, err := repo.ResetDailyQuotaForUser(context.Background(), 77, 501, "operation-hash")

	require.NoError(t, err)
	require.True(t, result.OperationReplayed)
	require.Equal(t, resetAt, result.ResetAt)
	require.Zero(t, result.DailyUsageUSD)
	require.Equal(t, 12.5, result.WeeklyUsageUSD)
	require.Equal(t, 31.75, result.MonthlyUsageUSD)
	require.Equal(t, availableAt, *result.DailyQuotaResetAvailableAt)
	require.NoError(t, mock.ExpectationsWereMet())

	sqlText := strings.ToUpper(normalizeSQLWhitespace(capturedSQL))
	require.Contains(t, sqlText, "OWNED_SUBSCRIPTION")
	require.Contains(t, sqlText, "OPERATION.VALUE->>'IDEMPOTENCY_KEY_HASH' = $4")
	require.Contains(t, sqlText, "OPERATION.VALUE->>'REPLAY_EXPIRES_AT'")
	require.Contains(t, sqlText, "FROM OPERATION_REPLAY")
}
