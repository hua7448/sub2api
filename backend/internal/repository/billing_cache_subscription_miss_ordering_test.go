//go:build unit

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBillingCacheSubscriptionMissOrderingSurvivesResetRetry(t *testing.T) {
	cache, _ := newMiniRedisCache(t)
	ctx := context.Background()
	const userID, groupID int64 = 107, 207
	dataKey := billingSubKey(userID, groupID)
	orderingKey := billingSubResetBarrierKey(userID, groupID)
	cutoff := time.Date(2026, 8, 1, 5, 43, 21, 300000000, time.UTC)
	t1 := cutoff.Add(time.Microsecond)
	t2 := cutoff.Add(2 * time.Microsecond)
	t3 := cutoff.Add(3 * time.Microsecond)
	t4 := cutoff.Add(4 * time.Microsecond)

	require.NoError(t, cache.SetSubscriptionResetBarrierAndInvalidate(ctx, userID, groupID, cutoff))

	// The post-reset billing write arrives while the data cache is missing.
	require.NoError(t, cache.UpdateSubscriptionUsageAt(ctx, userID, groupID, 5, t2))
	exists, err := cache.rdb.Exists(ctx, dataKey).Result()
	require.NoError(t, err)
	require.Zero(t, exists)
	highwater, err := cache.rdb.HGet(ctx, orderingKey, "write_highwater").Int64()
	require.NoError(t, err)
	require.Equal(t, t2.UnixMicro(), highwater)

	for _, staleVersion := range []int64{cutoff.UnixMicro(), t1.UnixMicro()} {
		require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(staleVersion, 0)))
		exists, err = cache.rdb.Exists(ctx, dataKey).Result()
		require.NoError(t, err)
		require.Zero(t, exists, "snapshot older than a cache-miss billing write must be rejected")
	}

	// A DB snapshot at the missing-cache write highwater can populate data.
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(t2.UnixMicro(), 5)))
	got, err := cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 5, got.DailyUsage, 0.000001)
	require.Equal(t, t2.UnixMicro(), got.Version)
	highwater, err = cache.rdb.HGet(ctx, orderingKey, "write_highwater").Int64()
	require.NoError(t, err)
	require.Equal(t, t2.UnixMicro(), highwater)

	// A still newer DB snapshot advances ordering highwater.
	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(t3.UnixMicro(), 5.5)))
	got, err = cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.Equal(t, t3.UnixMicro(), got.Version)
	highwater, err = cache.rdb.HGet(ctx, orderingKey, "write_highwater").Int64()
	require.NoError(t, err)

	// The delayed reset retry deletes data but must retain the newer highwater.
	require.NoError(t, cache.SetSubscriptionResetBarrierAndInvalidate(ctx, userID, groupID, cutoff))
	highwater, err = cache.rdb.HGet(ctx, orderingKey, "write_highwater").Int64()
	require.NoError(t, err)
	require.Equal(t, t3.UnixMicro(), highwater)

	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(t2.UnixMicro(), 5)))
	exists, err = cache.rdb.Exists(ctx, dataKey).Result()
	require.NoError(t, err)
	require.Zero(t, exists, "retry DEL must not let an older snapshot recreate data")

	require.NoError(t, cache.SetSubscriptionCache(ctx, userID, groupID, barrierTestSubscription(t4.UnixMicro(), 6)))
	got, err = cache.GetSubscriptionCache(ctx, userID, groupID)
	require.NoError(t, err)
	require.InDelta(t, 6, got.DailyUsage, 0.000001)
	require.Equal(t, t4.UnixMicro(), got.Version)
}
