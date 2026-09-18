//go:build unit

package service

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestBillingCacheServiceNormalizeCASConflictInvalidatesCache(t *testing.T) {
	weeklyLimit := 450.0
	oldWeeklyStart := time.Now().Add(-8 * 24 * time.Hour)
	baseCache := &billingCacheWorkerStub{
		subscriptionCache: &SubscriptionCacheData{
			Status:            SubscriptionStatusActive,
			StartsAt:          oldWeeklyStart,
			ExpiresAt:         time.Now().Add(24 * time.Hour),
			WeeklyWindowStart: &oldWeeklyStart,
			WeeklyUsage:       weeklyLimit,
			Version:           10,
			SnapshotVersion:   8,
			CacheRevision:     4,
		},
	}
	cache := &billingCacheNormalizerStub{
		billingCacheWorkerStub: baseCache,
		applied:                false,
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
	require.Equal(t, int64(1), atomic.LoadInt64(&cache.subscriptionInvalidations))
	require.Nil(t, cache.lastSetSubscription)
}

func TestBillingCacheServiceSubscriptionMissSetsAndReadsRevisionSynchronously(t *testing.T) {
	now := time.Now().UTC()
	cache := &billingCacheRoundTripStub{}
	repo := &billingCacheSubscriptionRepoStub{
		sub: &UserSubscription{
			ID:             77,
			UserID:         1184,
			GroupID:        11,
			Status:         SubscriptionStatusActive,
			StartsAt:       now.Add(-time.Hour),
			ExpiresAt:      now.Add(24 * time.Hour),
			DailyUsageUSD:  1.25,
			WeeklyUsageUSD: 2.5,
			UpdatedAt:      now,
		},
	}
	svc := NewBillingCacheService(cache, nil, repo, nil, nil, nil, &config.Config{}, nil)
	t.Cleanup(svc.Stop)

	data, err := svc.GetSubscriptionStatus(context.Background(), 1184, 11)

	require.NoError(t, err)
	require.True(t, data.cacheBacked)
	require.Equal(t, int64(1), data.CacheRevision)
	require.Equal(t, now.UnixMicro(), data.Version)
	require.Equal(t, now.UnixMicro(), data.SnapshotVersion)
	require.InDelta(t, 1.25, data.DailyUsage, 0.000001)
	require.Equal(t, int64(1), atomic.LoadInt64(&cache.subscriptionUpdates))
}
