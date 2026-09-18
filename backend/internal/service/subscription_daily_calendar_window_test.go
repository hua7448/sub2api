package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/stretchr/testify/require"
)

type subscriptionGetByIDSnapshotRepo struct {
	userSubRepoNoop
	snapshot *UserSubscription
}

func (r *subscriptionGetByIDSnapshotRepo) GetByID(_ context.Context, _ int64) (*UserSubscription, error) {
	return r.snapshot, nil
}

func useUTCGlobalAndShanghaiSubscriptionCalendar(t *testing.T) *time.Location {
	t.Helper()
	originalTimezone := timezone.Name()
	require.NoError(t, timezone.Init("UTC"))
	require.Equal(t, "UTC", timezone.Name())
	t.Cleanup(func() {
		require.NoError(t, timezone.Init(originalTimezone))
	})
	return subscriptionDailyCalendarLocation
}

func TestSubscriptionDailyCalendarBoundaryUsesShanghaiMidnight(t *testing.T) {
	location := useUTCGlobalAndShanghaiSubscriptionCalendar(t)
	windowStart := time.Date(2026, 8, 6, 15, 0, 0, 0, location)
	sub := &UserSubscription{
		StartsAt:         time.Date(2026, 8, 1, 12, 0, 0, 0, location),
		ExpiresAt:        time.Date(2026, 8, 20, 12, 0, 0, 0, location),
		DailyWindowStart: &windowStart,
	}

	beforeMidnightUTC := time.Date(2026, 8, 6, 15, 59, 59, 0, time.UTC)
	atMidnightUTC := time.Date(2026, 8, 6, 16, 0, 0, 0, time.UTC)
	require.False(t, sub.NeedsDailyResetAt(beforeMidnightUTC))
	require.True(t, sub.NeedsDailyResetAt(atMidnightUTC))

	newWindowStart, ok := sub.automaticDailyWindowStartAt(atMidnightUTC)
	require.True(t, ok)
	require.Equal(t, time.Date(2026, 8, 7, 0, 0, 0, 0, location), newWindowStart)
}

func TestSubscriptionDailyResetTimeUsesNextCalendarMidnightAfterManualReset(t *testing.T) {
	location := useUTCGlobalAndShanghaiSubscriptionCalendar(t)
	manualResetAt := time.Date(2026, 8, 6, 15, 0, 0, 0, location)
	sub := &UserSubscription{
		StartsAt:         time.Date(2026, 8, 1, 12, 0, 0, 0, location),
		ExpiresAt:        time.Date(2026, 8, 20, 12, 0, 0, 0, location),
		DailyWindowStart: &manualResetAt,
	}

	require.Equal(t, time.Date(2026, 8, 7, 0, 0, 0, 0, location), *sub.DailyResetTime())
}

func TestCheckAndResetWindowsAdvancesDailyAcrossMultipleDaysOnly(t *testing.T) {
	location := useUTCGlobalAndShanghaiSubscriptionCalendar(t)
	now := time.Date(2026, 8, 9, 12, 0, 0, 0, location)
	dailyStart := time.Date(2026, 8, 6, 15, 0, 0, 0, location)
	weeklyStart := time.Date(2026, 8, 8, 12, 0, 0, 0, location)
	monthlyStart := time.Date(2026, 8, 1, 12, 0, 0, 0, location)
	repo := &dailyResetTrackingUserSubRepo{}
	svc := NewSubscriptionService(groupRepoNoop{}, repo, nil, nil, nil)
	svc.now = func() time.Time { return now }
	sub := &UserSubscription{
		ID:                 1,
		UserID:             10,
		GroupID:            20,
		StartsAt:           time.Date(2026, 8, 1, 12, 0, 0, 0, location),
		ExpiresAt:          time.Date(2026, 8, 20, 12, 0, 0, 0, location),
		DailyWindowStart:   &dailyStart,
		WeeklyWindowStart:  &weeklyStart,
		MonthlyWindowStart: &monthlyStart,
		DailyUsageUSD:      9,
		WeeklyUsageUSD:     19,
		MonthlyUsageUSD:    29,
	}

	require.NoError(t, svc.CheckAndResetWindows(context.Background(), sub))
	require.True(t, repo.resetDailyCalled)
	require.Equal(t, time.Date(2026, 8, 9, 0, 0, 0, 0, location), repo.resetAt)
	require.Zero(t, sub.DailyUsageUSD)
	require.Equal(t, 19.0, sub.WeeklyUsageUSD)
	require.Equal(t, 29.0, sub.MonthlyUsageUSD)
	require.Equal(t, weeklyStart, *sub.WeeklyWindowStart)
	require.Equal(t, monthlyStart, *sub.MonthlyWindowStart)
}

func TestNormalizeExpiredWindowsAdvancesDailyResponseCopy(t *testing.T) {
	location := useUTCGlobalAndShanghaiSubscriptionCalendar(t)
	now := time.Date(2026, 8, 9, 12, 0, 0, 0, location)
	dailyStart := time.Date(2026, 8, 6, 15, 0, 0, 0, location)
	weeklyStart := time.Date(2026, 8, 8, 12, 0, 0, 0, location)
	monthlyStart := time.Date(2026, 8, 1, 12, 0, 0, 0, location)
	subs := []UserSubscription{{
		StartsAt:           time.Date(2026, 8, 1, 12, 0, 0, 0, location),
		ExpiresAt:          time.Date(2026, 8, 20, 12, 0, 0, 0, location),
		DailyWindowStart:   &dailyStart,
		WeeklyWindowStart:  &weeklyStart,
		MonthlyWindowStart: &monthlyStart,
		DailyUsageUSD:      9,
		WeeklyUsageUSD:     19,
		MonthlyUsageUSD:    29,
	}}

	normalizeExpiredWindowsAt(subs, now)

	require.Equal(t, time.Date(2026, 8, 9, 0, 0, 0, 0, location), *subs[0].DailyWindowStart)
	require.Zero(t, subs[0].DailyUsageUSD)
	require.Equal(t, weeklyStart, *subs[0].WeeklyWindowStart)
	require.Equal(t, monthlyStart, *subs[0].MonthlyWindowStart)
	require.Equal(t, time.Date(2026, 8, 10, 0, 0, 0, 0, location), *subs[0].DailyResetTime())
}

func TestValidateAndCheckLimitsUsesDailyCalendarBoundary(t *testing.T) {
	location := useUTCGlobalAndShanghaiSubscriptionCalendar(t)
	now := time.Date(2026, 8, 7, 0, 0, 0, 0, location)
	dailyStart := time.Date(2026, 8, 6, 15, 0, 0, 0, location)
	limit := 10.0
	sub := &UserSubscription{
		Status:           SubscriptionStatusActive,
		StartsAt:         time.Date(2026, 8, 1, 12, 0, 0, 0, location),
		ExpiresAt:        time.Date(2026, 8, 20, 12, 0, 0, 0, location),
		DailyWindowStart: &dailyStart,
		DailyUsageUSD:    11,
	}
	svc := NewSubscriptionService(groupRepoNoop{}, userSubRepoNoop{}, nil, nil, nil)
	svc.now = func() time.Time { return now }

	needsMaintenance, err := svc.ValidateAndCheckLimits(sub, &Group{DailyLimitUSD: &limit})

	require.NoError(t, err)
	require.True(t, needsMaintenance)
	require.Zero(t, sub.DailyUsageUSD)
}

func TestSubscriptionDailyCalendarResetDoesNotOutliveExpiry(t *testing.T) {
	location := useUTCGlobalAndShanghaiSubscriptionCalendar(t)
	dailyStart := time.Date(2026, 8, 6, 15, 0, 0, 0, location)
	expiresAt := time.Date(2026, 8, 7, 0, 0, 0, 0, location)
	sub := &UserSubscription{
		StartsAt:         time.Date(2026, 8, 1, 12, 0, 0, 0, location),
		ExpiresAt:        expiresAt,
		DailyWindowStart: &dailyStart,
	}

	_, ok := sub.automaticDailyWindowStartAt(expiresAt)
	require.False(t, ok)
	require.Nil(t, sub.DailyResetTime())
}

func TestCalculateProgressUsesDailyCalendarMidnight(t *testing.T) {
	location := useUTCGlobalAndShanghaiSubscriptionCalendar(t)
	now := time.Date(2090, 8, 9, 12, 0, 0, 0, location)
	dailyStart := time.Date(2090, 8, 6, 15, 0, 0, 0, location)
	limit := 10.0
	sub := &UserSubscription{
		StartsAt:         time.Date(2090, 8, 1, 12, 0, 0, 0, location),
		ExpiresAt:        time.Date(2090, 8, 20, 12, 0, 0, 0, location),
		DailyWindowStart: &dailyStart,
		DailyUsageUSD:    4,
	}
	svc := newTestSubscriptionService()
	svc.now = func() time.Time { return now }

	progress := svc.calculateProgress(sub, &Group{DailyLimitUSD: &limit})

	require.NotNil(t, progress.Daily)
	require.Equal(t, time.Date(2090, 8, 9, 0, 0, 0, 0, location), progress.Daily.WindowStart)
	require.Zero(t, progress.Daily.UsedUSD)
	require.Equal(t, time.Date(2090, 8, 10, 0, 0, 0, 0, location), progress.Daily.ResetsAt)
	require.Equal(t, dailyStart, *sub.DailyWindowStart, "progress normalization must not mutate the repository snapshot")
	require.Equal(t, 4.0, sub.DailyUsageUSD)
}

func TestGetByIDNormalizesDailyResponseCopyWithoutMutatingRepositorySnapshot(t *testing.T) {
	location := useUTCGlobalAndShanghaiSubscriptionCalendar(t)
	now := time.Date(2026, 8, 9, 12, 0, 0, 0, location)
	dailyStart := time.Date(2026, 8, 6, 15, 0, 0, 0, location)
	snapshot := &UserSubscription{
		ID:                       42,
		StartsAt:                 time.Date(2026, 8, 1, 12, 0, 0, 0, location),
		ExpiresAt:                time.Date(2026, 8, 20, 12, 0, 0, 0, location),
		DailyWindowStart:         &dailyStart,
		DailyUsageUSD:            9,
		WeeklyUsageUSD:           19,
		MonthlyUsageUSD:          29,
		DailyQuotaResetAvailable: true,
	}
	repo := &subscriptionGetByIDSnapshotRepo{snapshot: snapshot}
	svc := NewSubscriptionService(groupRepoNoop{}, repo, nil, nil, nil)
	svc.now = func() time.Time { return now }

	got, err := svc.GetByID(context.Background(), snapshot.ID)

	require.NoError(t, err)
	require.NotSame(t, snapshot, got)
	require.Equal(t, time.Date(2026, 8, 9, 0, 0, 0, 0, location), *got.DailyWindowStart)
	require.Zero(t, got.DailyUsageUSD)
	require.Equal(t, time.Date(2026, 8, 10, 0, 0, 0, 0, location), *got.DailyResetTime())
	require.True(t, got.DailyQuotaResetAvailable, "GetByID must not recalculate reset availability")
	require.Equal(t, 19.0, got.WeeklyUsageUSD)
	require.Equal(t, 29.0, got.MonthlyUsageUSD)

	require.Equal(t, dailyStart, *snapshot.DailyWindowStart)
	require.Equal(t, 9.0, snapshot.DailyUsageUSD)
	require.True(t, snapshot.DailyQuotaResetAvailable)
}
