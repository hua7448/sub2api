package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUserSubscriptionResetTimesUseCalendarDayForDailyLegacyAnchor(t *testing.T) {
	location := time.FixedZone("UTC+8", 8*60*60)
	startsAt := time.Date(2026, 7, 29, 14, 52, 1, 0, location)
	legacyWindowStart := startOfDay(startsAt)
	sub := &UserSubscription{
		StartsAt:           startsAt,
		ExpiresAt:          startsAt.Add(90 * 24 * time.Hour),
		DailyWindowStart:   &legacyWindowStart,
		WeeklyWindowStart:  &legacyWindowStart,
		MonthlyWindowStart: &legacyWindowStart,
	}

	require.Equal(t, subscriptionDailyCalendarStart(legacyWindowStart).AddDate(0, 0, 1), *sub.DailyResetTime())
	require.Equal(t, startsAt.Add(7*24*time.Hour), *sub.WeeklyResetTime())
	require.Equal(t, startsAt.Add(30*24*time.Hour), *sub.MonthlyResetTime())
}

func TestUserSubscriptionResetTimesUseCalendarDayForDailyManualAnchor(t *testing.T) {
	startsAt := time.Date(2026, 7, 29, 14, 52, 1, 0, time.UTC)
	manualWindowStart := startsAt.Add(-2 * time.Hour)
	sub := &UserSubscription{
		StartsAt:           startsAt,
		ExpiresAt:          startsAt.Add(90 * 24 * time.Hour),
		DailyWindowStart:   &manualWindowStart,
		WeeklyWindowStart:  &manualWindowStart,
		MonthlyWindowStart: &manualWindowStart,
	}

	require.Equal(t, subscriptionDailyCalendarStart(manualWindowStart).AddDate(0, 0, 1), *sub.DailyResetTime())
	require.Equal(t, manualWindowStart.Add(7*24*time.Hour), *sub.WeeklyResetTime())
	require.Equal(t, manualWindowStart.Add(30*24*time.Hour), *sub.MonthlyResetTime())
}

func TestUserSubscriptionResetTimesDoNotOutliveSubscription(t *testing.T) {
	startsAt := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	finalWindowStart := startsAt.Add(29 * 24 * time.Hour)
	expiresAt := subscriptionDailyCalendarStart(finalWindowStart).AddDate(0, 0, 1).Add(-time.Hour)
	sub := &UserSubscription{
		StartsAt:           startsAt,
		ExpiresAt:          expiresAt,
		DailyWindowStart:   &finalWindowStart,
		WeeklyWindowStart:  &finalWindowStart,
		MonthlyWindowStart: &finalWindowStart,
	}

	require.False(t, sub.HasOneTimeDailyQuota())
	require.Nil(t, sub.DailyResetTime())
	require.Nil(t, sub.WeeklyResetTime())
	require.Nil(t, sub.MonthlyResetTime())
}

func TestCalculateProgressUsesExpiryWhenNoFurtherWindowFits(t *testing.T) {
	startsAt := time.Date(2090, 1, 1, 10, 0, 0, 0, time.UTC)
	finalWindowStart := startsAt.Add(29 * 24 * time.Hour)
	expiresAt := subscriptionDailyCalendarStart(finalWindowStart).AddDate(0, 0, 1).Add(-time.Hour)
	limit := 100.0
	sub := &UserSubscription{
		StartsAt:           startsAt,
		ExpiresAt:          expiresAt,
		DailyWindowStart:   &finalWindowStart,
		WeeklyWindowStart:  &finalWindowStart,
		MonthlyWindowStart: &finalWindowStart,
	}

	progress := newTestSubscriptionService().calculateProgress(sub, &Group{
		DailyLimitUSD:   &limit,
		WeeklyLimitUSD:  &limit,
		MonthlyLimitUSD: &limit,
	})

	require.Equal(t, expiresAt, progress.Daily.ResetsAt)
	require.Equal(t, expiresAt, progress.Weekly.ResetsAt)
	require.Equal(t, expiresAt, progress.Monthly.ResetsAt)
}

func TestUserSubscriptionResetTimesReturnNilWithoutWindow(t *testing.T) {
	sub := &UserSubscription{}

	require.Nil(t, sub.DailyResetTime())
	require.Nil(t, sub.WeeklyResetTime())
	require.Nil(t, sub.MonthlyResetTime())
}

func TestCalculateProgressUsesEffectiveLegacyResetTimes(t *testing.T) {
	startsAt := time.Date(2026, 7, 29, 14, 52, 1, 0, time.UTC)
	legacyWindowStart := startOfDay(startsAt)
	limit := 200.0
	sub := &UserSubscription{
		StartsAt:           startsAt,
		ExpiresAt:          startsAt.Add(90 * 24 * time.Hour),
		WeeklyWindowStart:  &legacyWindowStart,
		MonthlyWindowStart: &legacyWindowStart,
	}

	progress := newTestSubscriptionService().calculateProgress(sub, &Group{WeeklyLimitUSD: &limit, MonthlyLimitUSD: &limit})

	require.Equal(t, startsAt.Add(7*24*time.Hour), progress.Weekly.ResetsAt)
	require.Equal(t, startsAt.Add(30*24*time.Hour), progress.Monthly.ResetsAt)
}
