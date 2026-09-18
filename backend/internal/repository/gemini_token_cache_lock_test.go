package repository

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func newOAuthRefreshLockTestCache(t *testing.T) (*geminiTokenCache, *miniredis.Miniredis) {
	t.Helper()
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	return &geminiTokenCache{rdb: client}, server
}

func TestOAuthRefreshLockExpiredOwnerCannotReleaseNewLease(t *testing.T) {
	cache, server := newOAuthRefreshLockTestCache(t)
	ctx := context.Background()
	const cacheKey = "claude:791"
	const ttl = 30 * time.Second
	key := oauthRefreshLockKeyPrefix + cacheKey

	acquired, err := cache.AcquireRefreshLock(ctx, cacheKey, "owner-A", ttl)
	require.NoError(t, err)
	require.True(t, acquired)
	require.Equal(t, ttl, server.TTL(key))
	server.FastForward(ttl + time.Second)
	acquired, err = cache.AcquireRefreshLock(ctx, cacheKey, "owner-B", ttl)
	require.NoError(t, err)
	require.True(t, acquired)

	for range 2 {
		require.NoError(t, cache.ReleaseRefreshLock(ctx, cacheKey, "owner-A"))
		owner, err := server.Get(key)
		require.NoError(t, err)
		require.Equal(t, "owner-B", owner, "late or repeated release must preserve the new owner's lease")
		require.Equal(t, ttl, server.TTL(key), "late release must not alter the new lease TTL")
	}
	acquired, err = cache.AcquireRefreshLock(ctx, cacheKey, "owner-C", ttl)
	require.NoError(t, err)
	require.False(t, acquired, "third worker must remain excluded while B owns the lease")
}

func TestOAuthRefreshLockOnlyOwnerCanRelease(t *testing.T) {
	cache, server := newOAuthRefreshLockTestCache(t)
	ctx := context.Background()
	const cacheKey = "openai:792"
	key := oauthRefreshLockKeyPrefix + cacheKey
	acquired, err := cache.AcquireRefreshLock(ctx, cacheKey, "owner-A", time.Minute)
	require.NoError(t, err)
	require.True(t, acquired)
	acquired, err = cache.AcquireRefreshLock(ctx, cacheKey, "owner-B", time.Minute)
	require.NoError(t, err)
	require.False(t, acquired)
	require.NoError(t, cache.ReleaseRefreshLock(ctx, cacheKey, "owner-B"))
	require.True(t, server.Exists(key))
	require.NoError(t, cache.ReleaseRefreshLock(ctx, cacheKey, "owner-A"))
	require.False(t, server.Exists(key))
	acquired, err = cache.AcquireRefreshLock(ctx, cacheKey, "owner-B", time.Minute)
	require.NoError(t, err)
	require.True(t, acquired)
	require.NoError(t, cache.ReleaseRefreshLock(ctx, cacheKey, "owner-A"))
	owner, err := server.Get(key)
	require.NoError(t, err)
	require.Equal(t, "owner-B", owner)
}

func TestOAuthRefreshLockRejectsEmptyOwner(t *testing.T) {
	cache, server := newOAuthRefreshLockTestCache(t)
	ctx := context.Background()
	const cacheKey = "gemini:793"
	acquired, err := cache.AcquireRefreshLock(ctx, cacheKey, "", time.Minute)
	require.Error(t, err)
	require.False(t, acquired)
	require.False(t, server.Exists(oauthRefreshLockKeyPrefix+cacheKey))
	acquired, err = cache.AcquireRefreshLock(ctx, cacheKey, "owner-A", time.Minute)
	require.NoError(t, err)
	require.True(t, acquired)
	require.Error(t, cache.ReleaseRefreshLock(ctx, cacheKey, ""))
	require.True(t, server.Exists(oauthRefreshLockKeyPrefix+cacheKey))
}
