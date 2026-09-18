//go:build unit

package repository

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestSchedulerAccountReadReplicaHitMissAndFailureFallback(t *testing.T) {
	ctx := context.Background()
	primaryServer := miniredis.RunT(t)
	readServer := miniredis.RunT(t)
	primary := redis.NewClient(&redis.Options{Addr: primaryServer.Addr()})
	readReplica := redis.NewClient(&redis.Options{Addr: readServer.Addr()})
	t.Cleanup(func() {
		_ = primary.Close()
		_ = readReplica.Close()
	})

	cache := newSchedulerCacheWithReadReplica(primary, readReplica, 128, 256).(*schedulerCache)
	primaryWriter := newSchedulerCacheWithChunkSizes(primary, 128, 256).(*schedulerCache)
	readWriter := newSchedulerCacheWithChunkSizes(readReplica, 128, 256).(*schedulerCache)

	primaryAccount := &service.Account{ID: 501, Name: "primary", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey}
	require.NoError(t, primaryWriter.SetAccount(ctx, primaryAccount))

	// Replica miss falls back to the authoritative cache.
	account, err := cache.GetAccount(ctx, primaryAccount.ID)
	require.NoError(t, err)
	require.Equal(t, "primary", account.Name)

	readAccount := *primaryAccount
	readAccount.Name = "replica"
	require.NoError(t, readWriter.SetAccount(ctx, &readAccount))
	account, err = cache.GetAccount(ctx, primaryAccount.ID)
	require.NoError(t, err)
	require.Equal(t, "replica", account.Name)

	// All cache writes still target the authoritative Redis.
	writeOnly := &service.Account{ID: 502, Name: "write-primary", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey}
	require.NoError(t, cache.SetAccount(ctx, writeOnly))
	primaryExists, err := primary.Exists(ctx, schedulerAccountKey("502")).Result()
	require.NoError(t, err)
	readExists, err := readReplica.Exists(ctx, schedulerAccountKey("502")).Result()
	require.NoError(t, err)
	require.EqualValues(t, 1, primaryExists)
	require.Zero(t, readExists)

	require.NoError(t, readReplica.Close())
	account, err = cache.GetAccount(ctx, primaryAccount.ID)
	require.NoError(t, err)
	require.Equal(t, "primary", account.Name)
}

func TestSchedulerSnapshotReadReplicaMissFallsBackThenHitWins(t *testing.T) {
	ctx := context.Background()
	primaryServer := miniredis.RunT(t)
	readServer := miniredis.RunT(t)
	primary := redis.NewClient(&redis.Options{Addr: primaryServer.Addr()})
	readReplica := redis.NewClient(&redis.Options{Addr: readServer.Addr()})
	t.Cleanup(func() {
		_ = primary.Close()
		_ = readReplica.Close()
	})

	cache := newSchedulerCacheWithReadReplica(primary, readReplica, 128, 256).(*schedulerCache)
	primaryWriter := newSchedulerCacheWithChunkSizes(primary, 128, 256).(*schedulerCache)
	readWriter := newSchedulerCacheWithChunkSizes(readReplica, 128, 256).(*schedulerCache)
	bucket := service.SchedulerBucket{GroupID: 77, Platform: service.PlatformOpenAI, Mode: service.SchedulerModeSingle}

	primaryToken, err := primaryWriter.CaptureBucketWriteToken(ctx, bucket)
	require.NoError(t, err)
	require.NoError(t, primaryWriter.SetSnapshot(ctx, bucket, primaryToken, []service.Account{{
		ID: 601, Name: "primary-snapshot", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
	}}))

	accounts, hit, err := cache.GetSnapshot(ctx, bucket)
	require.NoError(t, err)
	require.True(t, hit)
	require.Equal(t, "primary-snapshot", accounts[0].Name)

	readToken, err := readWriter.CaptureBucketWriteToken(ctx, bucket)
	require.NoError(t, err)
	require.NoError(t, readWriter.SetSnapshot(ctx, bucket, readToken, []service.Account{{
		ID: 601, Name: "replica-snapshot", Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
	}}))
	accounts, hit, err = cache.GetSnapshot(ctx, bucket)
	require.NoError(t, err)
	require.True(t, hit)
	require.Equal(t, "replica-snapshot", accounts[0].Name)
}
