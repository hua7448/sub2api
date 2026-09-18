package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/redis/go-redis/v9"
)

const (
	oauthTokenKeyPrefix       = "oauth:token:"
	oauthRefreshLockKeyPrefix = "oauth:refresh_lock:"
)

// Only the acquisition that still owns the lease may delete it. A late release
// after TTL expiry must leave a newer worker's refresh lock intact.
var oauthRefreshLockReleaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0
`)

type geminiTokenCache struct {
	rdb *redis.Client
}

func NewGeminiTokenCache(rdb *redis.Client) service.GeminiTokenCache {
	return &geminiTokenCache{rdb: rdb}
}

func (c *geminiTokenCache) GetAccessToken(ctx context.Context, cacheKey string) (string, error) {
	key := fmt.Sprintf("%s%s", oauthTokenKeyPrefix, cacheKey)
	return c.rdb.Get(ctx, key).Result()
}

func (c *geminiTokenCache) SetAccessToken(ctx context.Context, cacheKey string, token string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", oauthTokenKeyPrefix, cacheKey)
	return c.rdb.Set(ctx, key, token, ttl).Err()
}

func (c *geminiTokenCache) DeleteAccessToken(ctx context.Context, cacheKey string) error {
	key := fmt.Sprintf("%s%s", oauthTokenKeyPrefix, cacheKey)
	return c.rdb.Del(ctx, key).Err()
}

func (c *geminiTokenCache) AcquireRefreshLock(ctx context.Context, cacheKey, owner string, ttl time.Duration) (bool, error) {
	if owner == "" {
		return false, fmt.Errorf("refresh lock owner is required")
	}
	key := fmt.Sprintf("%s%s", oauthRefreshLockKeyPrefix, cacheKey)
	return c.rdb.SetNX(ctx, key, owner, ttl).Result()
}

func (c *geminiTokenCache) ReleaseRefreshLock(ctx context.Context, cacheKey, owner string) error {
	if owner == "" {
		return fmt.Errorf("refresh lock owner is required")
	}
	key := fmt.Sprintf("%s%s", oauthRefreshLockKeyPrefix, cacheKey)
	return oauthRefreshLockReleaseScript.Run(ctx, c.rdb, []string{key}, owner).Err()
}
