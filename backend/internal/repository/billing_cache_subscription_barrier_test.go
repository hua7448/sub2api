//go:build unit

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func barrierTestSubscription(version int64, dailyUsage float64) *service.SubscriptionCacheData {
	return &service.SubscriptionCacheData{
		Status:          service.SubscriptionStatusActive,
		StartsAt:        time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		ExpiresAt:       time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		DailyUsage:      dailyUsage,
		WeeklyUsage:     dailyUsage,
		MonthlyUsage:    dailyUsage,
		SnapshotVersion: version,
		Version:         version,
	}
}

func TestBillingCacheSubscriptionBarrierRejectsStaleSnapshotAndAllowsFresh(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 101, 201
	cutoff := time.Date(2026, 8, 1, 5, 43, 21, 123456000, time.UTC)

	stale := barrierTestSubscription(cutoff.Add(-time.Microsecond).UnixMicro(), 9)
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, stale))
	require.NoError(t, cache.SetSubscriptionResetBarrierAndInvalidate(ctx, userID, groupID, cutoff))

	_, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.ErrorIs(t, err, redis.Nil)

	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, stale))
	_, err = cache.GetSubscriptionCache(ctx, userID, groupID)
	require.ErrorIs(t, err, redis.Nil, "snapshot older than reset cutoff must not recreate cache")

	fresh := barrierTestSubscription(cutoff.UnixMicro(), 0)
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, fresh))
	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.Equal(t, cutoff.UnixMicro(), got.Version)
	require.Zero(t, got.DailyUsage)

	ttl, err := cache.rdb.TTL(ctx, billingSubResetBarrierKey(userID, groupID)).Result()
	require.NoError(t, err)
	require.GreaterOrEqual(t, ttl, 10*time.Minute)
	require.LessOrEqual(t, ttl, billingSubResetBarrierTTL)
}

func TestBillingCacheSubscriptionBarrierRejectsOldIncrementAndAllowsNewIncrement(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 102, 202
	cutoff := time.Date(2026, 8, 1, 5, 43, 21, 654321000, time.UTC)

	require.NoError(t, cache.SetSubscriptionResetBarrierAndInvalidate(ctx, userID, groupID, cutoff))
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(cutoff.UnixMicro(), 0)))

	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 4.5, cutoff.Add(-time.Microsecond)))
	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 4.5, cutoff))

	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.Zero(t, got.DailyUsage, "write timestamps at or before reset cutoff must be rejected")

	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 1.25, cutoff.Add(time.Microsecond)))
	got, err = cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 1.25, got.DailyUsage, 0.000001)
	require.InDelta(t, 1.25, got.WeeklyUsage, 0.000001)
	require.InDelta(t, 1.25, got.MonthlyUsage, 0.000001)
}

func TestBillingCacheSubscriptionBarrierIsScopedToUserGroupPair(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	cutoff := time.Date(2026, 8, 1, 5, 43, 21, 0, time.UTC)

	require.NoError(t, cache.SetSubscriptionResetBarrierAndInvalidate(ctx, 103, 203, cutoff))
	olderSnapshot := barrierTestSubscription(cutoff.Add(-time.Hour).UnixMicro(), 2)
	require.NoError(t, cache.SetSubscriptionCache(ctx, 104, 204, olderSnapshot))

	got, err := cache.GetSubscriptionCache(ctx, 104, 204)
	require.NoError(t, err)
	require.InDelta(t, 2, got.DailyUsage, 0.000001)
}
