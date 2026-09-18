package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestAccountsLoadBatchUsesReadReplicaWithoutMutatingIt(t *testing.T) {
	ctx := context.Background()
	primaryServer := miniredis.RunT(t)
	readServer := miniredis.RunT(t)
	now := time.Now().UTC().Truncate(time.Second)
	readServer.SetTime(now)

	primary := redis.NewClient(&redis.Options{Addr: primaryServer.Addr()})
	readReplica := redis.NewClient(&redis.Options{Addr: readServer.Addr()})
	t.Cleanup(func() {
		_ = primary.Close()
		_ = readReplica.Close()
	})

	cache := NewConcurrencyCacheWithReadReplica(primary, readReplica, 15, 900)
	accountID := int64(41)
	regularKey := accountSlotKey(accountID)
	liveKey := liveAccountSlotKey(accountID)
	waitKey := accountWaitKey(accountID)
	require.NoError(t, readReplica.ZAdd(ctx, regularKey,
		redis.Z{Score: float64(now.Add(-16 * time.Minute).Unix()), Member: "expired"},
		redis.Z{Score: float64(now.Unix()), Member: "active"},
	).Err())
	require.NoError(t, readReplica.ZAdd(ctx, liveKey,
		redis.Z{Score: float64(now.Unix()), Member: "live"},
	).Err())
	require.NoError(t, readReplica.Set(ctx, waitKey, 2, time.Minute).Err())

	loadMap, err := cache.GetAccountsLoadBatch(ctx, []service.AccountWithConcurrency{{ID: accountID, MaxConcurrency: 10}})
	require.NoError(t, err)
	require.Equal(t, 2, loadMap[accountID].CurrentConcurrency)
	require.Equal(t, 2, loadMap[accountID].WaitingCount)
	require.Equal(t, 40, loadMap[accountID].LoadRate)

	// Advisory reads must not clean the replica. Authoritative acquisition on the
	// primary remains responsible for mutation and expiration cleanup.
	remaining, err := readReplica.ZCard(ctx, regularKey).Result()
	require.NoError(t, err)
	require.EqualValues(t, 2, remaining)
	primaryKeys, err := primary.DBSize(ctx).Result()
	require.NoError(t, err)
	require.Zero(t, primaryKeys)
}

func TestAccountsLoadBatchFallsBackToPrimaryWhenReadReplicaFails(t *testing.T) {
	ctx := context.Background()
	primaryServer := miniredis.RunT(t)
	readServer := miniredis.RunT(t)
	primaryServer.SetTime(time.Now().UTC().Truncate(time.Second))

	primary := redis.NewClient(&redis.Options{Addr: primaryServer.Addr()})
	readReplica := redis.NewClient(&redis.Options{Addr: readServer.Addr()})
	t.Cleanup(func() { _ = primary.Close() })

	accountID := int64(42)
	require.NoError(t, primary.ZAdd(ctx, accountSlotKey(accountID), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: "primary-active",
	}).Err())
	require.NoError(t, readReplica.Close())

	cache := NewConcurrencyCacheWithReadReplica(primary, readReplica, 15, 900)
	loadMap, err := cache.GetAccountsLoadBatch(ctx, []service.AccountWithConcurrency{{ID: accountID, MaxConcurrency: 4}})
	require.NoError(t, err)
	require.Equal(t, 1, loadMap[accountID].CurrentConcurrency)
	require.Equal(t, 25, loadMap[accountID].LoadRate)
}
