package service

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestBillingCacheServiceQueueUpdateSubscriptionUsageAtPassesTimestamp(t *testing.T) {
	cache := &billingCacheWorkerStub{}
	svc := NewBillingCacheService(cache, nil, nil, nil, nil, nil, &config.Config{}, nil)
	t.Cleanup(svc.Stop)

	writeAt := time.Date(2026, 8, 1, 5, 43, 21, 123456000, time.UTC)
	svc.QueueUpdateSubscriptionUsageAt(11, 22, 1.25, writeAt)

	require.Eventually(t, func() bool {
		return atomic.LoadInt64(&cache.subscriptionWriteAt) == writeAt.UnixNano()
	}, 2*time.Second, 10*time.Millisecond)
}
