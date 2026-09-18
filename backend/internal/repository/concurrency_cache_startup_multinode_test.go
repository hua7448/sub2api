package repository

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestConcurrencyStartupCleanupPreservesOtherNodesActiveSlots(t *testing.T) {
	server := miniredis.RunT(t)
	ctx := t.Context()
	newNode := func() *concurrencyCache {
		client := redis.NewClient(&redis.Options{Addr: server.Addr()})
		t.Cleanup(func() { _ = client.Close() })
		return NewConcurrencyCache(client, 15, 90).(*concurrencyCache)
	}
	nodeA, nodeB, nodeC := newNode(), newNode(), newNode()
	const accountID, userID = int64(8101), int64(8102)
	acquired, err := nodeA.AcquireAccountSlot(ctx, accountID, 1, "node-a-request")
	require.NoError(t, err)
	require.True(t, acquired)
	acquired, err = nodeA.AcquireUserSlot(ctx, userID, 1, "node-a-request")
	require.NoError(t, err)
	require.True(t, acquired)
	waiting, err := nodeA.IncrementAccountWaitCount(ctx, accountID, 10)
	require.NoError(t, err)
	require.True(t, waiting)
	waiting, err = nodeA.IncrementWaitCount(ctx, userID, 10)
	require.NoError(t, err)
	require.True(t, waiting)

	// There is deliberately no legacy sweep marker: even the first startup
	// against a shared Redis must preserve another live node's waiting count.
	require.NoError(t, nodeB.CleanupStaleProcessSlots(ctx, "node-b-"))
	acquired, err = nodeC.AcquireAccountSlot(ctx, accountID, 1, "node-c-request")
	require.NoError(t, err)
	require.False(t, acquired, "starting node B must not free node A's account slot")
	acquired, err = nodeC.AcquireUserSlot(ctx, userID, 1, "node-c-request")
	require.NoError(t, err)
	require.False(t, acquired, "starting node B must not free node A's user slot")
	for _, key := range []string{accountWaitKey(accountID), waitQueueKey(userID)} {
		count, err := nodeA.rdb.Get(ctx, key).Int()
		require.NoError(t, err)
		require.Equal(t, 1, count)
	}

	require.NoError(t, nodeA.ReleaseAccountSlot(ctx, accountID, "node-a-request"))
	require.NoError(t, nodeA.ReleaseUserSlot(ctx, userID, "node-a-request"))
	acquired, err = nodeC.AcquireAccountSlot(ctx, accountID, 1, "node-c-request")
	require.NoError(t, err)
	require.True(t, acquired, "node C may enter after the real holder releases")
	acquired, err = nodeC.AcquireUserSlot(ctx, userID, 1, "node-c-request")
	require.NoError(t, err)
	require.True(t, acquired)
}

func TestConcurrencyStartupCleanupReapsOnlyExpiredSlots(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	cache := NewConcurrencyCache(client, 15, 90).(*concurrencyCache)
	ctx := t.Context()
	now, err := cache.redisUnixSeconds(ctx)
	require.NoError(t, err)
	for _, spec := range []slotIndexSpec{accountSlotIndex, userSlotIndex} {
		t.Run(spec.indexKey, func(t *testing.T) {
			const mixedID, expiredID, waitingID = int64(8201), int64(8202), int64(8203)
			require.NoError(t, client.ZAdd(ctx, spec.slotKey(mixedID),
				redis.Z{Score: float64(now), Member: "other-node-active"},
				redis.Z{Score: float64(now - int64(cache.slotTTLSeconds)), Member: "other-node-expired"},
				redis.Z{Score: float64(now - int64(cache.slotTTLSeconds)), Member: "restarting-node-expired"},
			).Err())
			require.NoError(t, client.ZAdd(ctx, spec.slotKey(expiredID),
				redis.Z{Score: float64(now - int64(cache.slotTTLSeconds) - 1), Member: "other-node-expired"},
			).Err())
			for _, id := range []int64{mixedID, expiredID, waitingID} {
				require.NoError(t, client.ZAdd(ctx, spec.indexKey, redis.Z{
					Score: float64(now - 1), Member: strconv.FormatInt(id, 10),
				}).Err())
			}
			require.NoError(t, client.Set(ctx, spec.waitKey(waitingID), 2, time.Minute).Err())
			require.NoError(t, cache.CleanupStaleProcessSlots(ctx, "restarting-node-"))
			members, err := client.ZRange(ctx, spec.slotKey(mixedID), 0, -1).Result()
			require.NoError(t, err)
			require.Equal(t, []string{"other-node-active"}, members)
			exists, err := client.Exists(ctx, spec.slotKey(expiredID)).Result()
			require.NoError(t, err)
			require.Zero(t, exists)
			_, err = client.ZScore(ctx, spec.indexKey, strconv.FormatInt(expiredID, 10)).Result()
			require.ErrorIs(t, err, redis.Nil)
			for _, id := range []int64{mixedID, waitingID} {
				score, err := client.ZScore(ctx, spec.indexKey, strconv.FormatInt(id, 10)).Result()
				require.NoError(t, err)
				require.Greater(t, score, float64(now), fmt.Sprintf("live load %d must stay indexed", id))
			}
		})
	}
}

func TestConcurrencyStartupCleanupBoundsLegacyWaitKeysWithoutDeletingLiveCounts(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	cache := NewConcurrencyCache(client, 15, 90).(*concurrencyCache)
	ctx := t.Context()
	legacyKey, activeKey := accountWaitKey(8301), waitQueueKey(8302)
	require.NoError(t, client.Set(ctx, legacyKey, 2, 0).Err())
	require.NoError(t, client.Set(ctx, activeKey, 3, 30*time.Second).Err())
	require.NoError(t, cache.CleanupStaleProcessSlots(ctx, "node-b-"))
	for key, want := range map[string]int{legacyKey: 2, activeKey: 3} {
		got, err := client.Get(ctx, key).Int()
		require.NoError(t, err)
		require.Equal(t, want, got)
	}
	require.Equal(t, 90*time.Second, server.TTL(legacyKey))
	require.Equal(t, 30*time.Second, server.TTL(activeKey))
	server.FastForward(31 * time.Second)
	require.False(t, server.Exists(activeKey))
	require.NoError(t, cache.CleanupStaleProcessSlots(ctx, "node-c-"))
	require.Equal(t, 59*time.Second, server.TTL(legacyKey), "later startup must not renew the legacy waiting lifetime")
	server.FastForward(60 * time.Second)
	require.False(t, server.Exists(legacyKey))
}
