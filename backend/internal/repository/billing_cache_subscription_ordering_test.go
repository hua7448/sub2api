//go:build unit

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBillingCacheSubscriptionSnapshotCannotOverwriteNewerIncrement(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 105, 205
	cutoff := time.Date(2026, 8, 1, 5, 43, 21, 100000000, time.UTC)
	t1 := cutoff.Add(time.Microsecond)
	t2 := cutoff.Add(2 * time.Microsecond)
	t3 := cutoff.Add(3 * time.Microsecond)

	require.NoError(t, cache.SetSubscriptionResetBarrierAndInvalidate(ctx, userID, groupID, cutoff))
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(cutoff.UnixMicro(), 0)))

	// Equal versions remain writable for local window normalization.
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(cutoff.UnixMicro(), 0.25)))
	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 0.25, got.DailyUsage, 0.000001)

	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(cutoff.UnixMicro(), 0)))
	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 1.5, t2))

	for _, staleVersion := range []int64{cutoff.UnixMicro(), t1.UnixMicro()} {
		require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(staleVersion, 0)))
		got, err = cache.GetSubscriptionCache(ctx, userID, groupID)
		require.NoError(t, err)
		require.InDelta(t, 1.5, got.DailyUsage, 0.000001)
		require.Equal(t, t2.UnixMicro(), got.Version)
	}

	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(t3.UnixMicro(), 7)))
	got, err = cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 7, got.DailyUsage, 0.000001)
	require.Equal(t, t3.UnixMicro(), got.Version)
}

func TestBillingCacheSubscriptionPostResetIncrementsAccumulateOutOfOrder(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 106, 206
	cutoff := time.Date(2026, 8, 1, 5, 43, 21, 200000000, time.UTC)
	earlierWrite := cutoff.Add(2 * time.Microsecond)
	laterWrite := cutoff.Add(4 * time.Microsecond)

	require.NoError(t, cache.SetSubscriptionResetBarrierAndInvalidate(ctx, userID, groupID, cutoff))
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(cutoff.UnixMicro(), 0)))
	// Simulate a pre-deployment cache hash without snapshot_version.
	require.NoError(t, cache.rdb.HDel(ctx, billingSubKey(userID, groupID), subFieldSnapshotVersion).Err())

	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 2.25, laterWrite))
	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 1.75, earlierWrite))

	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 4, got.DailyUsage, 0.000001)
	require.InDelta(t, 4, got.WeeklyUsage, 0.000001)
	require.InDelta(t, 4, got.MonthlyUsage, 0.000001)
	require.Equal(t, laterWrite.UnixMicro(), got.Version)
	highwater, err := cache.rdb.HGet(ctx, billingSubResetBarrierKey(userID, groupID), "write_highwater").Int64()
	require.NoError(t, err)
	require.Equal(t, laterWrite.UnixMicro(), highwater)
	snapshotVersion, err := cache.rdb.HGet(ctx, billingSubKey(userID, groupID), subFieldSnapshotVersion).Int64()
	require.NoError(t, err)
	require.Equal(t, cutoff.UnixMicro(), snapshotVersion)
}
