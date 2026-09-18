//go:build unit

package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/stretchr/testify/require"
)

func TestResetTodayUsageWindowClassificationAtShanghaiMidnight(t *testing.T) {
	originalTimezone := timezone.Name()
	require.NoError(t, timezone.Init("UTC"))
	require.Equal(t, "UTC", timezone.Name())
	t.Cleanup(func() {
		require.NoError(t, timezone.Init(originalTimezone))
	})
	location := subscriptionDailyCalendarLocation

	cutoff := time.Date(2026, 8, 2, 12, 50, 0, 0, location)
	dayStart, err := resetTodayUsageDayStart(cutoff)
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 8, 2, 0, 0, 0, 0, location), dayStart)

	previousDayLastSecond := time.Date(2026, 8, 1, 23, 59, 59, 0, location)
	currentDayMidnight := time.Date(2026, 8, 2, 0, 0, 0, 0, location)

	require.False(t, resetTodayUsageWindowIsToday(sql.NullTime{Time: previousDayLastSecond, Valid: true}, dayStart, cutoff))
	require.True(t, resetTodayUsageWindowIsToday(sql.NullTime{Time: currentDayMidnight, Valid: true}, dayStart, cutoff))
}
