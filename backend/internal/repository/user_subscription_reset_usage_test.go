//go:build unit

package repository

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestResetUsageWindowsReturnsStrictlyMonotonicDatabaseCutoff(t *testing.T) {
	var capturedSQL string
	repo, mock := newUserDailyQuotaResetRepo(t, &capturedSQL)
	dailyStart := time.Date(2026, 8, 19, 0, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	periodicStart := time.Date(2026, 8, 19, 7, 31, 2, 123456000, time.FixedZone("UTC+8", 8*60*60))
	databaseCutoff := time.Date(2026, 8, 18, 23, 31, 2, 654322000, time.UTC)

	mock.ExpectQuery("reset usage windows").
		WithArgs(int64(501), true, false, true, dailyStart, periodicStart).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(databaseCutoff))

	cutoff, err := repo.ResetUsageWindows(context.Background(), 501, true, false, true, dailyStart, periodicStart)

	require.NoError(t, err)
	require.Equal(t, databaseCutoff, cutoff)
	require.NoError(t, mock.ExpectationsWereMet())

	sqlText := strings.ToUpper(normalizeSQLWhitespace(capturedSQL))
	require.Contains(t, sqlText, "UPDATE USER_SUBSCRIPTIONS AS US")
	require.Contains(t, sqlText, "DAILY_USAGE_USD = CASE WHEN $2::BOOLEAN THEN 0")
	require.Contains(t, sqlText, "WEEKLY_USAGE_USD = CASE WHEN $3::BOOLEAN THEN 0")
	require.Contains(t, sqlText, "MONTHLY_USAGE_USD = CASE WHEN $4::BOOLEAN THEN 0")
	require.Contains(t, sqlText, "UPDATED_AT = GREATEST(CLOCK_TIMESTAMP(), US.UPDATED_AT + INTERVAL '1 MICROSECOND')")
	require.Contains(t, sqlText, "RETURNING US.UPDATED_AT")
	require.NotContains(t, sqlText, "UPDATED_AT = $")
}

func TestResetUsageWindowsMissingSubscriptionReturnsNotFound(t *testing.T) {
	var capturedSQL string
	repo, mock := newUserDailyQuotaResetRepo(t, &capturedSQL)
	now := time.Date(2026, 8, 19, 1, 2, 3, 0, time.UTC)
	mock.ExpectQuery("reset usage windows").
		WithArgs(int64(999), true, false, false, now, now).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}))

	cutoff, err := repo.ResetUsageWindows(context.Background(), 999, true, false, false, now, now)

	require.True(t, cutoff.IsZero())
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
