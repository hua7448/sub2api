//go:build unit

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func dailyQuotaResetEligibilitySubscription(now time.Time) UserSubscription {
	dailyLimit := 20.0
	weeklyStart := now.Add(-24 * time.Hour)
	return UserSubscription{
		StartsAt:          now.Add(-10 * 24 * time.Hour),
		ExpiresAt:         now.Add(30 * 24 * time.Hour),
		Status:            SubscriptionStatusActive,
		WeeklyWindowStart: &weeklyStart,
		Group:             &Group{DailyLimitUSD: &dailyLimit},
	}
}

func TestDailyQuotaResetAvailabilityAllowsUnusedWeeklyOpportunity(t *testing.T) {
	now := time.Date(2026, 8, 5, 12, 50, 0, 0, time.UTC)
	sub := dailyQuotaResetEligibilitySubscription(now)

	sub.SetDailyQuotaResetAvailabilityAt(now)

	require.True(t, sub.DailyQuotaResetAvailable)
	require.Nil(t, sub.DailyQuotaResetAvailableAt)
}

func TestDailyQuotaResetAvailabilityBlocksCurrentWeeklyWindow(t *testing.T) {
	now := time.Date(2026, 8, 5, 12, 50, 0, 0, time.UTC)
	sub := dailyQuotaResetEligibilitySubscription(now)
	marker := *sub.WeeklyWindowStart
	sub.DailyQuotaResetWeekStart = &marker

	sub.SetDailyQuotaResetAvailabilityAt(now)

	require.False(t, sub.DailyQuotaResetAvailable)
	require.Equal(t, marker.Add(7*24*time.Hour), *sub.DailyQuotaResetAvailableAt)
}

func TestDailyQuotaResetAvailabilityReopensInNextWeeklyWindow(t *testing.T) {
	now := time.Date(2026, 8, 13, 12, 50, 0, 0, time.UTC)
	sub := dailyQuotaResetEligibilitySubscription(now)
	weeklyStart := now.Add(-8 * 24 * time.Hour)
	sub.WeeklyWindowStart = &weeklyStart
	marker := weeklyStart
	sub.DailyQuotaResetWeekStart = &marker
	subs := []UserSubscription{sub}

	normalizeExpiredWindowsAt(subs, now)

	require.True(t, subs[0].DailyQuotaResetAvailable)
	require.Nil(t, subs[0].DailyQuotaResetAvailableAt)
	require.Nil(t, subs[0].WeeklyWindowStart)
}

func TestDailyQuotaResetAvailabilityDoesNotAdvertiseAfterExpiry(t *testing.T) {
	now := time.Date(2026, 8, 5, 12, 50, 0, 0, time.UTC)
	sub := dailyQuotaResetEligibilitySubscription(now)
	marker := *sub.WeeklyWindowStart
	sub.DailyQuotaResetWeekStart = &marker
	sub.ExpiresAt = marker.Add(6 * 24 * time.Hour)

	sub.SetDailyQuotaResetAvailabilityAt(now)

	require.False(t, sub.DailyQuotaResetAvailable)
	require.Nil(t, sub.DailyQuotaResetAvailableAt)
}

func TestDailyQuotaResetAvailabilityUsesLegacySubscriptionAnchor(t *testing.T) {
	startsAt := time.Date(2026, 8, 5, 10, 0, 0, 0, time.UTC)
	legacyStart := time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)
	now := startsAt.Add(24 * time.Hour)
	dailyLimit := 20.0
	marker := startsAt
	sub := UserSubscription{
		StartsAt:                 startsAt,
		ExpiresAt:                startsAt.Add(30 * 24 * time.Hour),
		Status:                   SubscriptionStatusActive,
		WeeklyWindowStart:        &legacyStart,
		DailyQuotaResetWeekStart: &marker,
		Group:                    &Group{DailyLimitUSD: &dailyLimit},
	}

	sub.SetDailyQuotaResetAvailabilityAt(now)

	require.False(t, sub.DailyQuotaResetAvailable)
	require.Equal(t, startsAt.Add(7*24*time.Hour), *sub.DailyQuotaResetAvailableAt)
}

func TestDailyQuotaResetAvailabilityRejectsFutureMarker(t *testing.T) {
	now := time.Date(2026, 8, 5, 12, 50, 0, 0, time.UTC)
	sub := dailyQuotaResetEligibilitySubscription(now)
	futureMarker := sub.WeeklyWindowStart.Add(7 * 24 * time.Hour)
	sub.DailyQuotaResetWeekStart = &futureMarker

	sub.SetDailyQuotaResetAvailabilityAt(now)

	require.False(t, sub.DailyQuotaResetAvailable)
	require.Nil(t, sub.DailyQuotaResetAvailableAt)
}

func TestDailyQuotaResetAvailabilityFailsClosedForFutureWeeklyAnchor(t *testing.T) {
	now := time.Date(2026, 8, 5, 12, 50, 0, 0, time.UTC)
	sub := dailyQuotaResetEligibilitySubscription(now)
	futureAnchor := now.Add(24 * time.Hour)
	sub.WeeklyWindowStart = &futureAnchor
	marker := futureAnchor
	sub.DailyQuotaResetWeekStart = &marker

	sub.SetDailyQuotaResetAvailabilityAt(now)

	require.False(t, sub.DailyQuotaResetAvailable)
	require.Nil(t, sub.DailyQuotaResetAvailableAt)
}

func TestDailyQuotaResetAvailabilityUsesAsiaShanghaiLegacyAnchor(t *testing.T) {
	location, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)
	startsAt := time.Date(2026, 8, 5, 10, 0, 0, 0, location)
	legacyStart := time.Date(2026, 8, 5, 0, 0, 0, 0, location)
	now := startsAt.Add(24 * time.Hour)
	dailyLimit := 20.0
	marker := startsAt
	sub := UserSubscription{
		StartsAt:                 startsAt,
		ExpiresAt:                startsAt.Add(30 * 24 * time.Hour),
		Status:                   SubscriptionStatusActive,
		WeeklyWindowStart:        &legacyStart,
		DailyQuotaResetWeekStart: &marker,
		Group:                    &Group{DailyLimitUSD: &dailyLimit},
	}

	sub.SetDailyQuotaResetAvailabilityAt(now)

	require.False(t, sub.DailyQuotaResetAvailable)
	require.Equal(t, startsAt.Add(7*24*time.Hour), *sub.DailyQuotaResetAvailableAt)
	require.Equal(t, 10, sub.DailyQuotaResetAvailableAt.Hour())
}

func TestDailyQuotaResetAvailabilityUsesFixedHoursAcrossDST(t *testing.T) {
	location, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	weeklyStart := time.Date(2026, 3, 7, 12, 0, 0, 0, location)
	now := weeklyStart.Add(24 * time.Hour)
	dailyLimit := 20.0
	marker := weeklyStart
	sub := UserSubscription{
		StartsAt:                 weeklyStart.Add(-7 * 24 * time.Hour),
		ExpiresAt:                weeklyStart.Add(30 * 24 * time.Hour),
		Status:                   SubscriptionStatusActive,
		WeeklyWindowStart:        &weeklyStart,
		DailyQuotaResetWeekStart: &marker,
		Group:                    &Group{DailyLimitUSD: &dailyLimit},
	}

	sub.SetDailyQuotaResetAvailabilityAt(now)

	require.False(t, sub.DailyQuotaResetAvailable)
	require.Equal(t, 168*time.Hour, sub.DailyQuotaResetAvailableAt.Sub(weeklyStart))
	require.Equal(t, 13, sub.DailyQuotaResetAvailableAt.Hour())
}
