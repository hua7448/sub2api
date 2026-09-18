//go:build unit

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBillingCacheNormalizeCASRejectsOutOfOrderIncrementMutation(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 111, 211
	baseSnapshot := time.Date(2026, 8, 1, 5, 43, 21, 700000000, time.UTC)
	earlierWrite := baseSnapshot.Add(time.Microsecond)
	laterWrite := baseSnapshot.Add(2 * time.Microsecond)

	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(baseSnapshot.UnixMicro(), 0)))
	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 2, laterWrite))
	derived, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	staleRevision := derived.CacheRevision
	derived.DailyUsage = 0
	derived.WeeklyUsage = 0
	derived.MonthlyUsage = 0

	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 1, earlierWrite))
	applied, err := cache.NormalizeSubscriptionCache(ctx, userID, groupID, derived)
	require.NoError(t, err)
	require.False(t, applied, "an out-of-order increment must invalidate the derived CAS")

	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 3, got.DailyUsage, 0.000001)
	require.Equal(t, laterWrite.UnixMicro(), got.Version)
	require.Equal(t, staleRevision+1, got.CacheRevision)
}

func TestBillingCacheNormalizeCASRejectsEqualVersionDatabaseSnapshot(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 112, 212
	baseSnapshot := time.Date(2026, 8, 1, 5, 43, 21, 800000000, time.UTC)
	writeAt := baseSnapshot.Add(time.Microsecond)

	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(baseSnapshot.UnixMicro(), 0)))
	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 2, writeAt))
	derived, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	derived.DailyUsage = 0
	derived.WeeklyUsage = 0
	derived.MonthlyUsage = 0

	snapshot := barrierTestSubscription(writeAt.UnixMicro(), 7)
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, snapshot))
	applied, err := cache.NormalizeSubscriptionCache(ctx, userID, groupID, derived)
	require.NoError(t, err)
	require.False(t, applied, "a same-version DB snapshot must invalidate the derived CAS")

	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 7, got.DailyUsage, 0.000001)
	require.Equal(t, writeAt.UnixMicro(), got.SnapshotVersion)
	require.Greater(t, got.CacheRevision, derived.CacheRevision)
}

func TestBillingCacheNormalizeCASSuccessfullyAppliesAndAdvancesRevision(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 113, 213
	versionAt := time.Date(2026, 8, 1, 5, 43, 21, 900000000, time.UTC)
	oldWindow := versionAt.Add(-8 * 24 * time.Hour)
	newWindow := versionAt

	snapshot := barrierTestSubscription(versionAt.UnixMicro(), 9)
	snapshot.WeeklyWindowStart = &oldWindow
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, snapshot))
	derived, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	initialRevision := derived.CacheRevision
	derived.WeeklyUsage = 0
	derived.WeeklyWindowStart = &newWindow

	applied, err := cache.NormalizeSubscriptionCache(ctx, userID, groupID, derived)
	require.NoError(t, err)
	require.True(t, applied)

	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.Zero(t, got.WeeklyUsage)
	require.NotNil(t, got.WeeklyWindowStart)
	require.Equal(t, newWindow.Unix(), got.WeeklyWindowStart.Unix())
	require.Equal(t, versionAt.UnixMicro(), got.Version)
	require.Equal(t, versionAt.UnixMicro(), got.SnapshotVersion)
	require.Equal(t, initialRevision+1, got.CacheRevision)
}

func TestBillingCacheNormalizeCASDoesNotCreateMissingOrLegacyHash(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 114, 214
	versionAt := time.Date(2026, 8, 1, 5, 43, 22, 0, time.UTC)
	desired := barrierTestSubscription(versionAt.UnixMicro(), 0)
	desired.CacheRevision = 1

	applied, err := cache.NormalizeSubscriptionCache(ctx, userID, groupID, desired)
	require.NoError(t, err)
	require.False(t, applied)
	exists, err := cache.rdb.Exists(ctx, billingSubKey(userID, groupID)).Result()
	require.NoError(t, err)
	require.Zero(t, exists)

	legacy := barrierTestSubscription(versionAt.UnixMicro(), 4)
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, legacy))
	require.NoError(t, cache.rdb.HDel(ctx, billingSubKey(userID, groupID), subFieldCacheRevision).Err())
	cachedLegacy, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.Zero(t, cachedLegacy.CacheRevision)
	cachedLegacy.DailyUsage = 0

	applied, err = cache.NormalizeSubscriptionCache(ctx, userID, groupID, cachedLegacy)
	require.NoError(t, err)
	require.False(t, applied)
	dailyUsage, err := cache.rdb.HGet(ctx, billingSubKey(userID, groupID), subFieldDailyUsage).Float64()
	require.NoError(t, err)
	require.InDelta(t, 4, dailyUsage, 0.000001)
}
