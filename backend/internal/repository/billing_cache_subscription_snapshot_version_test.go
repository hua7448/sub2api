//go:build unit

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBillingCacheSubscriptionSnapshotSkipsDelayedIncrementAtSameVersion(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 108, 208
	cutoff := time.Date(2026, 8, 1, 5, 43, 21, 400000000, time.UTC)
	writeAt := cutoff.Add(2 * time.Microsecond)

	require.NoError(t, cache.SetSubscriptionResetBarrierAndInvalidate(ctx, userID, groupID, cutoff))
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(writeAt.UnixMicro(), 5)))

	// The DB snapshot already includes this billing write.
	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 5, writeAt))
	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 5, got.DailyUsage, 0.000001)
	require.InDelta(t, 5, got.WeeklyUsage, 0.000001)
	require.InDelta(t, 5, got.MonthlyUsage, 0.000001)
	require.Equal(t, writeAt.UnixMicro(), got.Version)
	require.Equal(t, writeAt.UnixMicro(), got.SnapshotVersion)
}

func TestBillingCacheSubscriptionSnapshotBaselineSkipsIncludedAndCountsNewerWrite(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 109, 209
	cutoff := time.Date(2026, 8, 1, 5, 43, 21, 500000000, time.UTC)
	snapshotAt := cutoff.Add(time.Microsecond)
	newerWrite := cutoff.Add(2 * time.Microsecond)

	require.NoError(t, cache.SetSubscriptionResetBarrierAndInvalidate(ctx, userID, groupID, cutoff))
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(snapshotAt.UnixMicro(), 3)))

	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 3, snapshotAt))
	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 2, newerWrite))

	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 5, got.DailyUsage, 0.000001)
	require.Equal(t, newerWrite.UnixMicro(), got.Version)
	require.Equal(t, snapshotAt.UnixMicro(), got.SnapshotVersion)
}

func TestBillingCacheSubscriptionNormalizePreservesSnapshotBaselineForOutOfOrderWrite(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 110, 210
	baseSnapshot := time.Date(2026, 8, 1, 5, 43, 21, 600000000, time.UTC)
	earlierWrite := baseSnapshot.Add(time.Microsecond)
	laterWrite := baseSnapshot.Add(2 * time.Microsecond)

	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(baseSnapshot.UnixMicro(), 0)))
	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 2, laterWrite))

	normalized, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.Equal(t, laterWrite.UnixMicro(), normalized.Version)
	require.Equal(t, baseSnapshot.UnixMicro(), normalized.SnapshotVersion)
	initialRevision := normalized.CacheRevision
	require.Positive(t, initialRevision)

	// Simulate local window normalization writing the derived cache state back.
	applied, err := cache.NormalizeSubscriptionCache(ctx, userID, groupID, normalized)
	require.NoError(t, err)
	require.True(t, applied)
	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 1, earlierWrite))

	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 3, got.DailyUsage, 0.000001)
	require.Equal(t, laterWrite.UnixMicro(), got.Version)
	require.Equal(t, baseSnapshot.UnixMicro(), got.SnapshotVersion)
	require.Equal(t, initialRevision+2, got.CacheRevision)
}
