package dto

import (
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestUserSubscriptionMappersExposeEffectiveWindowResetTimes(t *testing.T) {
	originalTimezone := timezone.Name()
	require.NoError(t, timezone.Init("UTC"))
	t.Cleanup(func() {
		require.NoError(t, timezone.Init(originalTimezone))
	})
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)

	startsAt := time.Date(2026, 7, 29, 14, 52, 1, 0, time.UTC)
	legacyWindowStart := time.Date(2026, 7, 29, 0, 0, 0, 0, time.UTC)
	dailyResetAvailableAt := startsAt.Add(7 * 24 * time.Hour)
	sub := &service.UserSubscription{
		StartsAt:                   startsAt,
		ExpiresAt:                  startsAt.Add(90 * 24 * time.Hour),
		DailyWindowStart:           &legacyWindowStart,
		WeeklyWindowStart:          &legacyWindowStart,
		MonthlyWindowStart:         &legacyWindowStart,
		DailyQuotaResetAvailable:   false,
		DailyQuotaResetAvailableAt: &dailyResetAvailableAt,
	}

	userDTO := UserSubscriptionFromService(sub)
	adminDTO := UserSubscriptionFromServiceAdmin(sub)

	require.NotNil(t, userDTO)
	require.NotNil(t, adminDTO)
	require.Equal(t, time.Date(2026, 7, 30, 0, 0, 0, 0, shanghai), *userDTO.DailyWindowResetsAt)
	require.Equal(t, startsAt.Add(7*24*time.Hour), *userDTO.WeeklyWindowResetsAt)
	require.Equal(t, startsAt.Add(30*24*time.Hour), *userDTO.MonthlyWindowResetsAt)
	require.False(t, userDTO.DailyQuotaResetAvailable)
	require.Equal(t, dailyResetAvailableAt, *userDTO.DailyQuotaResetAvailableAt)
	require.Equal(t, userDTO.DailyWindowResetsAt, adminDTO.DailyWindowResetsAt)
	require.Equal(t, userDTO.WeeklyWindowResetsAt, adminDTO.WeeklyWindowResetsAt)
	require.Equal(t, userDTO.MonthlyWindowResetsAt, adminDTO.MonthlyWindowResetsAt)
	require.Equal(t, userDTO.DailyQuotaResetAvailable, adminDTO.DailyQuotaResetAvailable)
	require.Equal(t, userDTO.DailyQuotaResetAvailableAt, adminDTO.DailyQuotaResetAvailableAt)
}

func TestUserSubscriptionMapperReturnsNullResetTimesWithoutWindows(t *testing.T) {
	out := UserSubscriptionFromService(&service.UserSubscription{})

	require.NotNil(t, out)
	require.Nil(t, out.DailyWindowResetsAt)
	require.Nil(t, out.WeeklyWindowResetsAt)
	require.Nil(t, out.MonthlyWindowResetsAt)
}
