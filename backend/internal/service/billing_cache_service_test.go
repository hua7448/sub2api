package service

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type billingCacheWorkerStub struct {
	balanceUpdates            int64
	subscriptionUpdates       int64
	subscriptionInvalidations int64
	subscriptionCache         *SubscriptionCacheData
	subscriptionWriteAt       int64
	subscriptionCacheErr      error
	lastSetSubscription       *SubscriptionCacheData
}

func (b *billingCacheWorkerStub) GetUserBalance(ctx context.Context, userID int64) (float64, error) {
	return 0, errors.New("not implemented")
}

func (b *billingCacheWorkerStub) SetUserBalance(ctx context.Context, userID int64, balance float64) error {
	atomic.AddInt64(&b.balanceUpdates, 1)
	return nil
}

func (b *billingCacheWorkerStub) DeductUserBalance(ctx context.Context, userID int64, amount float64) error {
	atomic.AddInt64(&b.balanceUpdates, 1)
	return nil
}

func (b *billingCacheWorkerStub) InvalidateUserBalance(ctx context.Context, userID int64) error {
	return nil
}

func (b *billingCacheWorkerStub) GetSubscriptionCache(ctx context.Context, userID, groupID int64) (*SubscriptionCacheData, error) {
	if b.subscriptionCache != nil || b.subscriptionCacheErr != nil {
		return b.subscriptionCache, b.subscriptionCacheErr
	}
	return nil, errors.New("not implemented")
}

func (b *billingCacheWorkerStub) SetSubscriptionCache(ctx context.Context, userID, groupID int64, data *SubscriptionCacheData) error {
	atomic.AddInt64(&b.subscriptionUpdates, 1)
	b.lastSetSubscription = data
	return nil
}

func (b *billingCacheWorkerStub) UpdateSubscriptionUsage(ctx context.Context, userID, groupID int64, cost float64) error {
	atomic.AddInt64(&b.subscriptionUpdates, 1)
	return nil
}
func (b *billingCacheWorkerStub) UpdateSubscriptionUsageAt(ctx context.Context, userID, groupID int64, cost float64, writeAt time.Time) error {
	atomic.StoreInt64(&b.subscriptionWriteAt, writeAt.UnixNano())
	return b.UpdateSubscriptionUsage(ctx, userID, groupID, cost)
}

func (b *billingCacheWorkerStub) InvalidateSubscriptionCache(ctx context.Context, userID, groupID int64) error {
	atomic.AddInt64(&b.subscriptionInvalidations, 1)
	return nil
}

func (b *billingCacheWorkerStub) GetAPIKeyRateLimit(ctx context.Context, keyID int64) (*APIKeyRateLimitCacheData, error) {
	return nil, errors.New("not implemented")
}

func (b *billingCacheWorkerStub) SetAPIKeyRateLimit(ctx context.Context, keyID int64, data *APIKeyRateLimitCacheData) error {
	return nil
}

func (b *billingCacheWorkerStub) UpdateAPIKeyRateLimitUsage(ctx context.Context, keyID int64, cost float64) error {
	return nil
}

func (b *billingCacheWorkerStub) InvalidateAPIKeyRateLimit(ctx context.Context, keyID int64) error {
	return nil
}

func (b *billingCacheWorkerStub) GetUserPlatformQuotaCache(ctx context.Context, userID int64, platform string) (*UserPlatformQuotaCacheEntry, bool, error) {
	return nil, false, nil
}

func (b *billingCacheWorkerStub) SetUserPlatformQuotaCache(ctx context.Context, userID int64, platform string, entry *UserPlatformQuotaCacheEntry, ttl time.Duration) error {
	return nil
}

func (b *billingCacheWorkerStub) DeleteUserPlatformQuotaCache(ctx context.Context, userID int64, platform string) error {
	return nil
}

func (b *billingCacheWorkerStub) IncrUserPlatformQuotaUsageCache(ctx context.Context, userID int64, platform string, cost float64, ttl time.Duration, markDirty bool) error {
	return nil
}

func (b *billingCacheWorkerStub) PopDirtyUserPlatformQuotaKeys(ctx context.Context, n int) ([]UserPlatformQuotaKey, error) {
	return nil, nil
}

func (b *billingCacheWorkerStub) ReaddDirtyUserPlatformQuotaKeys(ctx context.Context, keys []UserPlatformQuotaKey) error {
	return nil
}

func (b *billingCacheWorkerStub) BatchGetUserPlatformQuotaCache(ctx context.Context, keys []UserPlatformQuotaKey) ([]*UserPlatformQuotaCacheEntry, error) {
	return nil, nil
}

type billingCacheNormalizerStub struct {
	*billingCacheWorkerStub
	applied        bool
	normalizeErr   error
	normalizeCalls int64
	normalized     *SubscriptionCacheData
}

func (b *billingCacheNormalizerStub) NormalizeSubscriptionCache(_ context.Context, _, _ int64, data *SubscriptionCacheData) (bool, error) {
	atomic.AddInt64(&b.normalizeCalls, 1)
	b.normalized = data
	return b.applied, b.normalizeErr
}

type billingCacheRoundTripStub struct {
	billingCacheWorkerStub
}

func (b *billingCacheRoundTripStub) SetSubscriptionCache(_ context.Context, _, _ int64, data *SubscriptionCacheData) error {
	atomic.AddInt64(&b.subscriptionUpdates, 1)
	copied := *data
	copied.CacheRevision = 1
	b.subscriptionCache = &copied
	b.lastSetSubscription = &copied
	return nil
}

type billingCacheSubscriptionRepoStub struct {
	userSubRepoNoop
	sub *UserSubscription
}

func (b *billingCacheSubscriptionRepoStub) GetActiveByUserIDAndGroupID(context.Context, int64, int64) (*UserSubscription, error) {
	return b.sub, nil
}

func TestBillingCacheServiceQueueHighLoad(t *testing.T) {
	cache := &billingCacheWorkerStub{}
	svc := NewBillingCacheService(cache, nil, nil, nil, nil, nil, &config.Config{}, nil)
	t.Cleanup(svc.Stop)

	start := time.Now()
	for i := 0; i < cacheWriteBufferSize*2; i++ {
		svc.QueueDeductBalance(1, 1)
	}
	require.Less(t, time.Since(start), 2*time.Second)

	svc.QueueUpdateSubscriptionUsage(1, 2, 1.5)

	require.Eventually(t, func() bool {
		return atomic.LoadInt64(&cache.balanceUpdates) > 0
	}, 2*time.Second, 10*time.Millisecond)

	require.Eventually(t, func() bool {
		return atomic.LoadInt64(&cache.subscriptionUpdates) > 0
	}, 2*time.Second, 10*time.Millisecond)
}

func TestBillingCacheServiceEnqueueAfterStopReturnsFalse(t *testing.T) {
	cache := &billingCacheWorkerStub{}
	svc := NewBillingCacheService(cache, nil, nil, nil, nil, nil, &config.Config{}, nil)
	svc.Stop()

	enqueued := svc.enqueueCacheWrite(cacheWriteTask{
		kind:   cacheWriteDeductBalance,
		userID: 1,
		amount: 1,
	})
	require.False(t, enqueued)
}

func TestBillingCacheServiceSubscriptionEligibilityResetsExpiredWeeklyWindow(t *testing.T) {
	weeklyLimit := 450.0
	oldWeeklyStart := time.Now().Add(-8 * 24 * time.Hour)
	baseCache := &billingCacheWorkerStub{
		subscriptionCache: &SubscriptionCacheData{
			Status:            SubscriptionStatusActive,
			StartsAt:          oldWeeklyStart,
			ExpiresAt:         time.Now().Add(24 * time.Hour),
			WeeklyWindowStart: &oldWeeklyStart,
			WeeklyUsage:       weeklyLimit,
			Version:           1,
			SnapshotVersion:   1,
			CacheRevision:     7,
		},
	}
	cache := &billingCacheNormalizerStub{
		billingCacheWorkerStub: baseCache,
		applied:                true,
	}
	svc := NewBillingCacheService(cache, nil, nil, nil, nil, nil, &config.Config{}, nil)
	t.Cleanup(svc.Stop)

	err := svc.checkSubscriptionEligibility(context.Background(), 1184, &Group{
		ID:               11,
		SubscriptionType: SubscriptionTypeSubscription,
		WeeklyLimitUSD:   &weeklyLimit,
	}, nil)

	require.NoError(t, err)
	require.Equal(t, int64(1), atomic.LoadInt64(&cache.normalizeCalls))
	require.NotNil(t, cache.normalized)
	require.Zero(t, cache.normalized.WeeklyUsage)
	require.Equal(t, int64(1), cache.normalized.Version)
	require.Equal(t, int64(1), cache.normalized.SnapshotVersion)
	require.NotNil(t, cache.normalized.WeeklyWindowStart)
	require.Zero(t, atomic.LoadInt64(&cache.subscriptionInvalidations))
	require.Nil(t, cache.lastSetSubscription)
}

func TestBillingCacheServiceSubscriptionEligibilityMigratesLegacyCacheFromSubscriptionSnapshot(t *testing.T) {
	weeklyLimit := 450.0
	now := time.Now()
	currentWeeklyStart := startOfDay(now)
	cache := &billingCacheWorkerStub{
		subscriptionCache: &SubscriptionCacheData{
			Status:      SubscriptionStatusActive,
			ExpiresAt:   now.Add(24 * time.Hour),
			WeeklyUsage: weeklyLimit,
			Version:     1,
		},
	}
	svc := NewBillingCacheService(cache, nil, nil, nil, nil, nil, &config.Config{}, nil)
	t.Cleanup(svc.Stop)

	err := svc.checkSubscriptionEligibility(context.Background(), 1184, &Group{
		ID:               11,
		SubscriptionType: SubscriptionTypeSubscription,
		WeeklyLimitUSD:   &weeklyLimit,
	}, &UserSubscription{
		ID:                666,
		UserID:            1184,
		GroupID:           11,
		Status:            SubscriptionStatusActive,
		StartsAt:          now.Add(-24 * time.Hour),
		ExpiresAt:         now.Add(24 * time.Hour),
		WeeklyWindowStart: &currentWeeklyStart,
		WeeklyUsageUSD:    0,
		UpdatedAt:         now,
	})

	require.NoError(t, err)
	require.Nil(t, cache.lastSetSubscription)
	require.Equal(t, int64(1), atomic.LoadInt64(&cache.subscriptionInvalidations))
}
