//go:build unit

package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type userDailyQuotaResetRepoStub struct {
	userSubRepoNoop
	result             *UserDailyQuotaResetResult
	err                error
	lastUserID         int64
	lastSubscriptionID int64
	lastOperationKey   string
}

func (r *userDailyQuotaResetRepoStub) ResetDailyQuotaForUser(_ context.Context, userID, subscriptionID int64, operationKeyHash string) (*UserDailyQuotaResetResult, error) {
	r.lastUserID = userID
	r.lastSubscriptionID = subscriptionID
	r.lastOperationKey = operationKeyHash
	if r.err != nil {
		return nil, r.err
	}
	result := *r.result
	return &result, nil
}

type userDailyQuotaResetCacheStub struct {
	billingCacheWorkerStub

	mu              sync.Mutex
	barrierCalls    int
	publishCalls    int
	barrierFailures int
	cutoffs         []time.Time
}

type userDailyQuotaFinalFailureCacheStub struct {
	billingCacheWorkerStub

	mu      sync.Mutex
	started chan struct{}
	release chan struct{}
	once    sync.Once
	cutoffs []time.Time
}

type subscriptionQuotaGenerationFailureCacheStub struct {
	billingCacheWorkerStub

	mu      sync.Mutex
	started chan struct{}
	release chan struct{}
	once    sync.Once
	cutoffs []time.Time
}

func (c *subscriptionQuotaGenerationFailureCacheStub) SetSubscriptionResetBarrierAndInvalidate(_ context.Context, _, _ int64, cutoff time.Time) error {
	c.mu.Lock()
	c.cutoffs = append(c.cutoffs, cutoff)
	call := len(c.cutoffs)
	c.mu.Unlock()
	if call == 1 {
		c.once.Do(func() { close(c.started) })
		<-c.release
	}
	return errors.New("redis barrier unavailable")
}

func (c *subscriptionQuotaGenerationFailureCacheStub) PublishSubscriptionCacheInvalidation(context.Context, string) error {
	return nil
}

func (c *subscriptionQuotaGenerationFailureCacheStub) SubscribeSubscriptionCacheInvalidation(context.Context, func(string)) error {
	return nil
}

func (c *subscriptionQuotaGenerationFailureCacheStub) recordedCutoffs() []time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]time.Time(nil), c.cutoffs...)
}

func (c *userDailyQuotaFinalFailureCacheStub) SetSubscriptionResetBarrierAndInvalidate(_ context.Context, _, _ int64, cutoff time.Time) error {
	c.mu.Lock()
	c.cutoffs = append(c.cutoffs, cutoff)
	call := len(c.cutoffs)
	c.mu.Unlock()
	if call == 1 {
		c.once.Do(func() { close(c.started) })
		<-c.release
		return errors.New("redis barrier unavailable")
	}
	return nil
}

func (c *userDailyQuotaFinalFailureCacheStub) PublishSubscriptionCacheInvalidation(context.Context, string) error {
	return nil
}

func (c *userDailyQuotaFinalFailureCacheStub) SubscribeSubscriptionCacheInvalidation(context.Context, func(string)) error {
	return nil
}

func (c *userDailyQuotaFinalFailureCacheStub) recordedCutoffs() []time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]time.Time(nil), c.cutoffs...)
}

func (c *userDailyQuotaResetCacheStub) SetSubscriptionResetBarrierAndInvalidate(_ context.Context, _, _ int64, cutoff time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.barrierCalls++
	c.cutoffs = append(c.cutoffs, cutoff)
	if c.barrierFailures > 0 {
		c.barrierFailures--
		return errors.New("redis barrier unavailable")
	}
	return nil
}

func (c *userDailyQuotaResetCacheStub) PublishSubscriptionCacheInvalidation(context.Context, string) error {
	c.mu.Lock()
	c.publishCalls++
	c.mu.Unlock()
	return nil
}

func (c *userDailyQuotaResetCacheStub) SubscribeSubscriptionCacheInvalidation(context.Context, func(string)) error {
	return nil
}

func (c *userDailyQuotaResetCacheStub) calls() (int, int, []time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.barrierCalls, c.publishCalls, append([]time.Time(nil), c.cutoffs...)
}

func TestResetUserDailyQuotaCacheFailureDoesNotFailAndRetries(t *testing.T) {
	previousDelays := append([]time.Duration(nil), subscriptionQuotaCacheRetryDelays...)
	subscriptionQuotaCacheRetryDelays = []time.Duration{time.Millisecond, 2 * time.Millisecond}
	t.Cleanup(func() { subscriptionQuotaCacheRetryDelays = previousDelays })

	location := time.FixedZone("UTC+8", 8*60*60)
	resetAt := time.Date(2026, 8, 5, 12, 50, 0, 123000000, location)
	repo := &userDailyQuotaResetRepoStub{result: &UserDailyQuotaResetResult{
		SubscriptionID:  501,
		GroupID:         19,
		ResetAt:         resetAt,
		DailyUsageUSD:   0,
		WeeklyUsageUSD:  12.5,
		MonthlyUsageUSD: 31.75,
	}}
	cache := &userDailyQuotaResetCacheStub{barrierFailures: 1}
	svc := NewSubscriptionService(nil, repo, &BillingCacheService{cache: cache}, nil, nil)

	result, err := svc.ResetUserDailyQuota(context.Background(), 77, 501, "operation-hash")

	require.NoError(t, err)
	require.Equal(t, int64(77), repo.lastUserID)
	require.Equal(t, int64(501), repo.lastSubscriptionID)
	require.Equal(t, resetAt.UTC(), result.ResetAt)
	require.Equal(t, "operation-hash", repo.lastOperationKey)
	require.Zero(t, result.DailyUsageUSD)
	require.Equal(t, 12.5, result.WeeklyUsageUSD)
	require.Equal(t, 31.75, result.MonthlyUsageUSD)

	require.Eventually(t, func() bool {
		barriers, publishes, _ := cache.calls()
		return barriers >= 2 && publishes >= 2
	}, time.Second, time.Millisecond)
	_, _, cutoffs := cache.calls()
	for _, cutoff := range cutoffs {
		require.Equal(t, resetAt.UTC(), cutoff)
	}
}

func TestResetUserDailyQuotaCacheRetriesThroughLaterRecovery(t *testing.T) {
	previousDelays := append([]time.Duration(nil), subscriptionQuotaCacheRetryDelays...)
	subscriptionQuotaCacheRetryDelays = []time.Duration{
		time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		4 * time.Millisecond,
		5 * time.Millisecond,
	}
	t.Cleanup(func() { subscriptionQuotaCacheRetryDelays = previousDelays })

	resetAt := time.Date(2026, 8, 5, 12, 50, 0, 0, time.UTC)
	repo := &userDailyQuotaResetRepoStub{result: &UserDailyQuotaResetResult{
		SubscriptionID:  501,
		GroupID:         19,
		ResetAt:         resetAt,
		DailyUsageUSD:   0,
		WeeklyUsageUSD:  12.5,
		MonthlyUsageUSD: 31.75,
	}}
	cache := &userDailyQuotaResetCacheStub{barrierFailures: 5}
	svc := NewSubscriptionService(nil, repo, &BillingCacheService{cache: cache}, nil, nil)

	result, err := svc.ResetUserDailyQuota(context.Background(), 77, 501, "operation-hash")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Eventually(t, func() bool {
		barriers, publishes, _ := cache.calls()
		return barriers == 6 && publishes == 6
	}, time.Second, time.Millisecond)
	_, _, cutoffs := cache.calls()
	require.Len(t, cutoffs, 6)
	for _, cutoff := range cutoffs {
		require.Equal(t, resetAt, cutoff)
	}
}

func TestUserDailyQuotaCacheRetriesCoalesceByUserGroup(t *testing.T) {
	previousDelays := append([]time.Duration(nil), subscriptionQuotaCacheRetryDelays...)
	subscriptionQuotaCacheRetryDelays = []time.Duration{20 * time.Millisecond, 40 * time.Millisecond}
	t.Cleanup(func() { subscriptionQuotaCacheRetryDelays = previousDelays })

	cache := &userDailyQuotaResetCacheStub{barrierFailures: 100}
	svc := NewSubscriptionService(nil, &userDailyQuotaResetRepoStub{}, &BillingCacheService{cache: cache}, nil, nil)
	baseCutoff := time.Date(2026, 8, 5, 12, 50, 0, 0, time.UTC)
	latestCutoff := baseCutoff.Add(99 * time.Microsecond)

	for i := 0; i < 100; i++ {
		svc.scheduleUserDailyQuotaCacheRetries(77, 19, baseCutoff.Add(time.Duration(i)*time.Microsecond))
	}

	svc.dailyQuotaCacheRetryMu.Lock()
	require.Len(t, svc.dailyQuotaCacheRetries, 1)
	state := svc.dailyQuotaCacheRetries[subCacheKey(77, 19)]
	require.NotNil(t, state)
	require.Equal(t, latestCutoff, state.resetAt)
	svc.dailyQuotaCacheRetryMu.Unlock()

	require.Eventually(t, func() bool {
		barriers, publishes, cutoffs := cache.calls()
		if barriers < 2 || barriers > 4 || publishes != barriers || len(cutoffs) != barriers {
			return false
		}
		for _, cutoff := range cutoffs {
			if !cutoff.Equal(latestCutoff) {
				return false
			}
		}
		svc.dailyQuotaCacheRetryMu.Lock()
		defer svc.dailyQuotaCacheRetryMu.Unlock()
		return len(svc.dailyQuotaCacheRetries) == 0
	}, time.Second, time.Millisecond)
}

func TestUserDailyQuotaCacheRetryKeepsCutoffJoiningFinalFailure(t *testing.T) {
	previousDelays := append([]time.Duration(nil), subscriptionQuotaCacheRetryDelays...)
	subscriptionQuotaCacheRetryDelays = []time.Duration{time.Millisecond}
	t.Cleanup(func() { subscriptionQuotaCacheRetryDelays = previousDelays })

	cache := &userDailyQuotaFinalFailureCacheStub{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	svc := NewSubscriptionService(nil, &userDailyQuotaResetRepoStub{}, &BillingCacheService{cache: cache}, nil, nil)
	firstCutoff := time.Date(2026, 8, 5, 12, 50, 0, 0, time.UTC)
	latestCutoff := firstCutoff.Add(time.Minute)

	svc.scheduleUserDailyQuotaCacheRetries(77, 19, firstCutoff)
	select {
	case <-cache.started:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for final retry attempt")
	}
	svc.scheduleUserDailyQuotaCacheRetries(77, 19, latestCutoff)
	close(cache.release)

	require.Eventually(t, func() bool {
		cutoffs := cache.recordedCutoffs()
		if len(cutoffs) != 2 || !cutoffs[0].Equal(firstCutoff) || !cutoffs[1].Equal(latestCutoff) {
			return false
		}
		svc.dailyQuotaCacheRetryMu.Lock()
		defer svc.dailyQuotaCacheRetryMu.Unlock()
		return len(svc.dailyQuotaCacheRetries) == 0
	}, time.Second, time.Millisecond)
}

func TestSubscriptionQuotaCacheRetryNewGenerationGetsFullCycleAfterLastFailure(t *testing.T) {
	previousDelays := append([]time.Duration(nil), subscriptionQuotaCacheRetryDelays...)
	subscriptionQuotaCacheRetryDelays = []time.Duration{time.Millisecond, 2 * time.Millisecond}
	t.Cleanup(func() { subscriptionQuotaCacheRetryDelays = previousDelays })

	cache := &subscriptionQuotaGenerationFailureCacheStub{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	svc := NewSubscriptionService(nil, &userDailyQuotaResetRepoStub{}, &BillingCacheService{cache: cache}, nil, nil)
	firstCutoff := time.Date(2026, 8, 5, 12, 50, 0, 0, time.UTC)
	latestCutoff := firstCutoff.Add(time.Minute)

	svc.scheduleSubscriptionResetCacheRetries(77, 19, firstCutoff)
	select {
	case <-cache.started:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first retry attempt")
	}
	svc.scheduleSubscriptionResetCacheRetries(77, 19, latestCutoff)
	close(cache.release)

	require.Eventually(t, func() bool {
		cutoffs := cache.recordedCutoffs()
		if len(cutoffs) != 4 {
			return false
		}
		svc.dailyQuotaCacheRetryMu.Lock()
		defer svc.dailyQuotaCacheRetryMu.Unlock()
		return len(svc.dailyQuotaCacheRetries) == 0
	}, time.Second, time.Millisecond)
	cutoffs := cache.recordedCutoffs()
	require.Equal(t, firstCutoff, cutoffs[0])
	for _, cutoff := range cutoffs[1:] {
		require.Equal(t, latestCutoff, cutoff)
	}
}

func TestResetUserDailyQuotaRepositoryFailureSkipsCaches(t *testing.T) {
	repo := &userDailyQuotaResetRepoStub{err: ErrSubscriptionNotFound}
	cache := &userDailyQuotaResetCacheStub{}
	svc := NewSubscriptionService(nil, repo, &BillingCacheService{cache: cache}, nil, nil)

	result, err := svc.ResetUserDailyQuota(context.Background(), 77, 999, "operation-hash")

	require.Nil(t, result)
	require.ErrorIs(t, err, ErrSubscriptionNotFound)
	barriers, publishes, _ := cache.calls()
	require.Zero(t, barriers)
	require.Zero(t, publishes)
}
