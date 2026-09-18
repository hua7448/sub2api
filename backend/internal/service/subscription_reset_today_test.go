//go:build unit

package service

import (
	"context"
	stdsql "database/sql"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/DATA-DOG/go-sqlmock"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/dgraph-io/ristretto"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func newResetTodaySQLMockService(t *testing.T) (*SubscriptionService, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	driver := entsql.OpenDB(dialect.Postgres, db)
	client := dbent.NewClient(dbent.Driver(driver))
	t.Cleanup(func() { _ = client.Close() })
	return &SubscriptionService{entClient: client}, mock
}

func resetTodayTestDayStart(t *testing.T) time.Time {
	t.Helper()
	dayStart, err := resetTodayUsageDayStart(resetTodayTestCapturedAt())
	require.NoError(t, err)
	return dayStart
}

func resetTodayTestCapturedAt() time.Time {
	return time.Date(2026, 8, 1, 5, 43, 21, 123456000, time.UTC)
}

func resetTodayTestCutoff() time.Time {
	return resetTodayTestCapturedAt().Add(time.Microsecond)
}

type resetTodayFailingCache struct {
	billingCacheWorkerStub
	mu          sync.Mutex
	invalidated int
	published   int
	cutoffs     []time.Time
}

type resetTodaySelectiveRetryCache struct {
	billingCacheWorkerStub

	mu         sync.Mutex
	failUserID int64
	calls      map[int64]int
	callTimes  map[int64][]time.Time
}

func (c *resetTodaySelectiveRetryCache) SetSubscriptionResetBarrierAndInvalidate(_ context.Context, userID, _ int64, _ time.Time) error {
	c.mu.Lock()
	if c.calls == nil {
		c.calls = make(map[int64]int)
	}
	c.calls[userID]++
	if c.callTimes == nil {
		c.callTimes = make(map[int64][]time.Time)
	}
	c.callTimes[userID] = append(c.callTimes[userID], time.Now())
	c.mu.Unlock()
	if userID == c.failUserID {
		return errors.New("redis barrier unavailable")
	}
	return nil
}

func (c *resetTodaySelectiveRetryCache) callSnapshot(userID int64) (int, []time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.calls[userID], append([]time.Time(nil), c.callTimes[userID]...)
}

type resetTodaySlowTargetCache struct {
	billingCacheWorkerStub

	slowStarted sync.Once
	fastStarted sync.Once
	slowReady   chan struct{}
	fastReady   chan struct{}
	releaseSlow chan struct{}
}

func (c *resetTodaySlowTargetCache) SetSubscriptionResetBarrierAndInvalidate(ctx context.Context, userID, _ int64, _ time.Time) error {
	switch userID {
	case 11:
		c.slowStarted.Do(func() { close(c.slowReady) })
		select {
		case <-c.releaseSlow:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	case 12:
		c.fastStarted.Do(func() { close(c.fastReady) })
	}
	return nil
}

type resetTodayBlockingCache struct {
	billingCacheWorkerStub

	mu        sync.Mutex
	active    int
	maxActive int
	calls     int
}

func (c *resetTodayBlockingCache) SetSubscriptionResetBarrierAndInvalidate(ctx context.Context, _, _ int64, _ time.Time) error {
	c.mu.Lock()
	c.active++
	c.calls++
	if c.active > c.maxActive {
		c.maxActive = c.active
	}
	c.mu.Unlock()
	<-ctx.Done()
	c.mu.Lock()
	c.active--
	c.mu.Unlock()
	return ctx.Err()
}

func (c *resetTodayBlockingCache) stats() (int, int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.calls, c.maxActive
}

func (c *resetTodayFailingCache) InvalidateSubscriptionCache(context.Context, int64, int64) error {
	c.mu.Lock()
	c.invalidated++
	c.mu.Unlock()
	return errors.New("redis delete unavailable")
}

func (c *resetTodayFailingCache) SetSubscriptionResetBarrierAndInvalidate(_ context.Context, _ int64, _ int64, cutoff time.Time) error {
	c.mu.Lock()
	c.invalidated++
	c.cutoffs = append(c.cutoffs, cutoff)
	c.mu.Unlock()
	return errors.New("redis barrier unavailable")
}

func (c *resetTodayFailingCache) PublishSubscriptionCacheInvalidation(context.Context, string) error {
	c.mu.Lock()
	c.published++
	c.mu.Unlock()
	return nil
}

func (c *resetTodayFailingCache) SubscribeSubscriptionCacheInvalidation(context.Context, func(string)) error {
	return nil
}

func (c *resetTodayFailingCache) calls() (int, int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.invalidated, c.published
}

func (c *resetTodayFailingCache) barrierCutoffs() []time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]time.Time(nil), c.cutoffs...)
}

func TestAdminResetTodayUsage_AtomicSummaryAndCacheFailureStillSucceeds(t *testing.T) {
	svc, mock := newResetTodaySQLMockService(t)
	cache := &resetTodayFailingCache{}
	svc.billingCacheService = &BillingCacheService{cache: cache}

	keyHash := HashIdempotencyKey("reset-2026-08-01")
	dayStart := resetTodayTestDayStart(t)
	staleStart := dayStart.Add(-24 * time.Hour)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO subscription_quota_reset_operations`).
		WithArgs(int64(7), keyHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(77)))
	mock.ExpectExec(`SELECT pg_advisory_xact_lock`).
		WithArgs(resetTodayUsageAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`(?s)SELECT id\s+FROM user_subscriptions.*FOR UPDATE`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(101)).AddRow(int64(102)))
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*MAX\(updated_at\).*interval '1 microsecond'.*FROM reset_clock`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"captured_at", "cutoff_at"}).
			AddRow(resetTodayTestCapturedAt(), resetTodayTestCutoff()))
	mock.ExpectQuery(`(?s)WITH subscription_operations AS MATERIALIZED.*reconstructed_operations AS MATERIALIZED.*updated_subscriptions AS.*UPDATE user_subscriptions AS us.*SELECT`).
		WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC(), BillingTypeSubscription, UsageBillingMonetaryScale).
		WillReturnRows(sqlmock.NewRows([]string{"malformed", "unproven", "updated"}).AddRow(int64(0), int64(0), int64(1)))
	mock.ExpectQuery(`SELECT id, user_id, group_id, daily_window_start, daily_usage_usd::text`).
		WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "group_id", "daily_window_start", "daily_usage_usd", "unconsumed_reset_usage_usd"}).
			AddRow(int64(101), int64(11), int64(21), dayStart, "3.0000000000", "10.0000000000").
			AddRow(int64(102), int64(12), int64(22), staleStart, "7.5000000000", "0"))
	mock.ExpectExec(`(?s)UPDATE user_subscriptions.*weekly_usage_usd = GREATEST.*monthly_usage_usd = GREATEST`).
		WithArgs(dayStart.UTC(), resetTodayTestCapturedAt(), resetTodayTestCutoff(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(77)).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec(`(?s)UPDATE subscription_quota_reset_operations\s+SET status = 'succeeded'`).
		WithArgs(int64(77), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectExec(`(?s)UPDATE subscription_quota_reset_operations\s+SET summary =`).
		WithArgs(int64(77), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := svc.AdminResetTodayUsage(context.Background(), 7, keyHash)

	require.NoError(t, err)
	require.Equal(t, resetTodayTestCutoff(), result.CutoffAt)
	require.Equal(t, 2, result.Subscriptions)
	require.Equal(t, resetTodayTestCapturedAt(), result.CapturedAt)
	require.Equal(t, 1, result.TodayWindowSubscriptions)
	require.Equal(t, 1, result.ResetSubscriptions)
	require.InDelta(t, 13, result.ResetAndDeductedUSD, 0.0000000001)
	require.InDelta(t, 7.5, result.StaleDailyClearedUSD, 0.0000000001)
	require.Zero(t, result.CacheInvalidations)
	require.Equal(t, 2, result.CacheInvalidationFailures)
	invalidated, published := cache.calls()
	require.Equal(t, 2, invalidated, "every target must receive Redis DEL")
	require.Equal(t, 2, published, "PUBLISH must still run when DEL fails")
	require.NoError(t, mock.ExpectationsWereMet())
	cutoffs := cache.barrierCutoffs()
	require.Len(t, cutoffs, 2)
	require.Equal(t, resetTodayTestCutoff(), cutoffs[0])
	require.Equal(t, resetTodayTestCutoff(), cutoffs[1])
}

func TestResetTodayCacheRetryUsesOnlyFailedSubsetAndRunsPastCacheTTL(t *testing.T) {
	previousDelays := append([]time.Duration(nil), subscriptionQuotaCacheRetryDelays...)
	subscriptionQuotaCacheRetryDelays = []time.Duration{time.Millisecond, 3 * time.Millisecond, 8 * time.Millisecond}
	t.Cleanup(func() { subscriptionQuotaCacheRetryDelays = previousDelays })

	cache := &resetTodaySelectiveRetryCache{failUserID: 11}
	svc := &SubscriptionService{billingCacheService: &BillingCacheService{cache: cache}}
	cutoff := resetTodayTestCutoff()
	targets := []resetTodayUsageTarget{{UserID: 11, GroupID: 21}, {UserID: 12, GroupID: 22}}

	succeeded, failedTargets := svc.invalidateResetTodayUsageTargetBatchWithTimeouts(targets, cutoff, false, time.Second, 100*time.Millisecond)
	require.Equal(t, 1, succeeded)
	require.Equal(t, []resetTodayUsageTarget{{UserID: 11, GroupID: 21}}, failedTargets)

	retryStarted := time.Now()
	svc.scheduleResetTodayUsageCacheRetries(failedTargets, cutoff)
	require.Eventually(t, func() bool {
		failedCalls, _ := cache.callSnapshot(11)
		successCalls, _ := cache.callSnapshot(12)
		return failedCalls == 4 && successCalls == 1
	}, time.Second, time.Millisecond)
	_, failedCallTimes := cache.callSnapshot(11)
	require.Len(t, failedCallTimes, 4)
	require.GreaterOrEqual(t, failedCallTimes[3].Sub(retryStarted), 5*time.Millisecond, "persistent failure must still be retried after the compressed cache TTL")
}

func TestResetTodayCacheBatchClearsAllL1BeforeSlowRemoteAndDoesNotStarveFastTarget(t *testing.T) {
	cache := &resetTodaySlowTargetCache{
		slowReady:   make(chan struct{}),
		fastReady:   make(chan struct{}),
		releaseSlow: make(chan struct{}),
	}
	svc := &SubscriptionService{billingCacheService: &BillingCacheService{cache: cache}}
	l1, err := ristretto.NewCache(&ristretto.Config{NumCounters: 100, MaxCost: 100, BufferItems: 64})
	require.NoError(t, err)
	t.Cleanup(l1.Close)
	svc.subCacheL1 = l1
	for _, target := range []resetTodayUsageTarget{{UserID: 11, GroupID: 21}, {UserID: 12, GroupID: 22}} {
		require.True(t, l1.Set(subCacheKey(target.UserID, target.GroupID), &UserSubscription{ID: target.UserID}, 1))
	}
	l1.Wait()

	type batchResult struct {
		succeeded int
		failed    []resetTodayUsageTarget
	}
	done := make(chan batchResult, 1)
	go func() {
		succeeded, failed := svc.invalidateResetTodayUsageTargetBatchWithTimeouts(
			[]resetTodayUsageTarget{{UserID: 11, GroupID: 21}, {UserID: 12, GroupID: 22}},
			resetTodayTestCutoff(),
			false,
			time.Second,
			500*time.Millisecond,
		)
		done <- batchResult{succeeded: succeeded, failed: failed}
	}()

	select {
	case <-cache.slowReady:
	case <-time.After(time.Second):
		t.Fatal("slow target did not start")
	}
	select {
	case <-cache.fastReady:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("fast target was starved by the slow target")
	}
	_, foundSlow := l1.Get(subCacheKey(11, 21))
	_, foundFast := l1.Get(subCacheKey(12, 22))
	require.False(t, foundSlow)
	require.False(t, foundFast, "every local L1 key must be gone before remote invalidation waits")
	close(cache.releaseSlow)
	result := <-done
	require.Equal(t, 2, result.succeeded)
	require.Empty(t, result.failed)
}

func TestResetTodayCacheBatchHasBoundedConcurrencyAndSharedDeadline(t *testing.T) {
	cache := &resetTodayBlockingCache{}
	svc := &SubscriptionService{billingCacheService: &BillingCacheService{cache: cache}}
	targets := make([]resetTodayUsageTarget, 65)
	for i := range targets {
		targets[i] = resetTodayUsageTarget{UserID: int64(i + 1), GroupID: 21}
	}

	started := time.Now()
	succeeded, failedTargets := svc.invalidateResetTodayUsageTargetBatchWithTimeouts(
		targets,
		resetTodayTestCutoff(),
		false,
		40*time.Millisecond,
		time.Second,
	)
	elapsed := time.Since(started)

	require.Zero(t, succeeded)
	require.Len(t, failedTargets, len(targets), "deadline must classify queued targets as failed for retry")
	require.Less(t, elapsed, 250*time.Millisecond, "batch must be bounded by one shared deadline, not N/worker waves")
	calls, maxActive := cache.stats()
	require.LessOrEqual(t, maxActive, resetTodayUsageCacheWorkers)
	require.LessOrEqual(t, calls, resetTodayUsageCacheWorkers, "targets left after the shared deadline must fail without another remote call")
	require.Contains(t, failedTargets, targets[len(targets)-1], "targets beyond the first worker wave must still be returned for retry")
}

func TestResetTodayUsageTransactionAcceptsNativeSelfResetFromInitiallyNullDailyWindow(t *testing.T) {
	svc, mock := newResetTodaySQLMockService(t)
	keyHash := HashIdempotencyKey("native-self-reset")
	dayStart := resetTodayTestDayStart(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO subscription_quota_reset_operations`).
		WithArgs(int64(7), keyHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(78)))
	mock.ExpectExec(`SELECT pg_advisory_xact_lock`).
		WithArgs(resetTodayUsageAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`(?s)SELECT id\s+FROM user_subscriptions.*FOR UPDATE`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(101)))
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*MAX\(updated_at\).*interval '1 microsecond'.*FROM reset_clock`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"captured_at", "cutoff_at"}).
			AddRow(resetTodayTestCapturedAt(), resetTodayTestCutoff()))
	// The user-reset repository fills an initially NULL daily window with the
	// current Shanghai midnight. That native operation must not be classified as
	// malformed or sent through legacy reconstruction.
	mock.ExpectQuery(`(?s)legacy_operations AS MATERIALIZED.*cleared_daily_window_start.*cleared_daily_usage_usd.*updated_subscriptions AS.*SELECT`).
		WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC(), BillingTypeSubscription, UsageBillingMonetaryScale).
		WillReturnRows(sqlmock.NewRows([]string{"malformed", "unproven", "updated"}).AddRow(int64(0), int64(0), int64(0)))
	// An initially unactivated window clears no prior usage; usage accrued after
	// the self-reset is still reset normally by the administrator operation.
	mock.ExpectQuery(`(?s)SELECT id, user_id, group_id, daily_window_start, daily_usage_usd::text.*cleared_daily_usage_usd`).
		WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "group_id", "daily_window_start", "daily_usage_usd", "unconsumed_reset_usage_usd"}).
			AddRow(int64(101), int64(11), int64(21), dayStart, "3.0000000000", "0"))
	mock.ExpectExec(`(?s)UPDATE user_subscriptions.*admin_reset_consumed_operation_id.*weekly_usage_usd = GREATEST|(?s)UPDATE user_subscriptions.*weekly_usage_usd = GREATEST.*admin_reset_consumed_operation_id`).
		WithArgs(dayStart.UTC(), resetTodayTestCapturedAt(), resetTodayTestCutoff(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(78)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)UPDATE subscription_quota_reset_operations\s+SET status = 'succeeded'`).
		WithArgs(int64(78), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	op, err := svc.resetTodayUsageTransaction(context.Background(), 7, keyHash)

	require.NoError(t, err)
	require.NotNil(t, op)
	require.Equal(t, 1, op.Result.ResetSubscriptions)
	require.InDelta(t, 3, op.Result.ResetAndDeductedUSD, 0.0000000001)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetTodayUsageTransaction_ReplaysDurableOperationWithoutSecondReset(t *testing.T) {
	svc, mock := newResetTodaySQLMockService(t)
	keyHash := HashIdempotencyKey("durable-replay")
	cutoff := time.Date(2026, 8, 1, 4, 50, 0, 0, time.UTC)
	summary := `{"cutoff_at":"2026-08-01T04:50:00Z","captured_at":"2026-08-01T04:50:01Z","subscriptions":2,"today_window_subscriptions":1,"reset_subscriptions":1,"reset_and_deducted_usd":3.25,"stale_daily_cleared_usd":7.5,"cache_invalidations":2,"cache_invalidation_failures":0}`
	targets := `[{"user_id":11,"group_id":21},{"user_id":12,"group_id":22}]`

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO subscription_quota_reset_operations`).
		WithArgs(int64(7), keyHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT id, status, summary, target_pairs`).
		WithArgs(int64(7), keyHash).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "summary", "target_pairs"}).
			AddRow(int64(77), "succeeded", []byte(summary), []byte(targets)))
	mock.ExpectCommit()

	op, err := svc.resetTodayUsageTransaction(context.Background(), 7, keyHash)

	require.NoError(t, err)
	require.Equal(t, int64(77), op.ID)
	require.Equal(t, cutoff, op.Result.CutoffAt)
	require.Equal(t, 2, op.Result.Subscriptions)
	require.Len(t, op.Targets, 2)
	require.NoError(t, mock.ExpectationsWereMet(), "replay must not issue quota UPDATE statements")
}

func TestResetTodayUsageTransaction_RollsBackAllChangesOnUpdateFailure(t *testing.T) {
	svc, mock := newResetTodaySQLMockService(t)
	keyHash := HashIdempotencyKey("rollback-on-update-error")
	dayStart := resetTodayTestDayStart(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO subscription_quota_reset_operations`).
		WithArgs(int64(7), keyHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(88)))
	mock.ExpectExec(`SELECT pg_advisory_xact_lock`).
		WithArgs(resetTodayUsageAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`(?s)SELECT id\s+FROM user_subscriptions.*FOR UPDATE`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(101)))
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*MAX\(updated_at\).*interval '1 microsecond'.*FROM reset_clock`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"captured_at", "cutoff_at"}).
			AddRow(resetTodayTestCapturedAt(), resetTodayTestCutoff()))
	mock.ExpectQuery(`(?s)WITH subscription_operations AS MATERIALIZED.*reconstructed_operations AS MATERIALIZED.*updated_subscriptions AS.*UPDATE user_subscriptions AS us.*SELECT`).
		WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC(), BillingTypeSubscription, UsageBillingMonetaryScale).
		WillReturnRows(sqlmock.NewRows([]string{"malformed", "unproven", "updated"}).AddRow(int64(0), int64(0), int64(1)))
	mock.ExpectQuery(`SELECT id, user_id, group_id, daily_window_start, daily_usage_usd::text`).
		WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "group_id", "daily_window_start", "daily_usage_usd", "unconsumed_reset_usage_usd"}).
			AddRow(int64(101), int64(11), int64(21), dayStart, "3.2500000000", "10.0000000000"))
	mock.ExpectExec(`UPDATE user_subscriptions`).
		WithArgs(dayStart.UTC(), resetTodayTestCapturedAt(), resetTodayTestCutoff(), sqlmock.AnyArg(), sqlmock.AnyArg(), int64(88)).
		WillReturnError(errors.New("update failed"))
	mock.ExpectRollback()

	_, err := svc.resetTodayUsageTransaction(context.Background(), 7, keyHash)

	require.ErrorContains(t, err, "update reset-today subscriptions")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetTodayUsageTransaction_RollsBackOnMalformedNativeOperation(t *testing.T) {
	svc, mock := newResetTodaySQLMockService(t)
	keyHash := HashIdempotencyKey("rollback-on-malformed-native-operation")
	dayStart := resetTodayTestDayStart(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO subscription_quota_reset_operations`).
		WithArgs(int64(7), keyHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(89)))
	mock.ExpectExec(`SELECT pg_advisory_xact_lock`).
		WithArgs(resetTodayUsageAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`(?s)SELECT id\s+FROM user_subscriptions.*FOR UPDATE`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(101)))
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*MAX\(updated_at\).*interval '1 microsecond'.*FROM reset_clock`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"captured_at", "cutoff_at"}).
			AddRow(resetTodayTestCapturedAt(), resetTodayTestCutoff()))
	mock.ExpectQuery(`(?s)WITH subscription_operations AS MATERIALIZED.*malformed_self_reset_operations AS MATERIALIZED.*updated_subscriptions AS.*UPDATE user_subscriptions AS us.*SELECT`).
		WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC(), BillingTypeSubscription, UsageBillingMonetaryScale).
		WillReturnRows(sqlmock.NewRows([]string{"malformed", "unproven", "updated"}).AddRow(int64(1), int64(0), int64(0)))
	mock.ExpectRollback()

	_, err := svc.resetTodayUsageTransaction(context.Background(), 7, keyHash)

	require.ErrorContains(t, err, "malformed current-day self-reset operations")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetTodayUsageTransaction_RollsBackOnMalformedLegacyOperation(t *testing.T) {
	svc, mock := newResetTodaySQLMockService(t)
	keyHash := HashIdempotencyKey("rollback-on-malformed-legacy-operation")
	dayStart := resetTodayTestDayStart(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO subscription_quota_reset_operations`).
		WithArgs(int64(7), keyHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(90)))
	mock.ExpectExec(`SELECT pg_advisory_xact_lock`).
		WithArgs(resetTodayUsageAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`(?s)SELECT id\s+FROM user_subscriptions.*FOR UPDATE`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(101)))
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*MAX\(updated_at\).*interval '1 microsecond'.*FROM reset_clock`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"captured_at", "cutoff_at"}).
			AddRow(resetTodayTestCapturedAt(), resetTodayTestCutoff()))
	mock.ExpectQuery(`(?s)WITH subscription_operations AS MATERIALIZED.*valid_self_reset_operations AS MATERIALIZED.*malformed_self_reset_operations AS MATERIALIZED.*updated_subscriptions AS.*UPDATE user_subscriptions AS us.*SELECT`).
		WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC(), BillingTypeSubscription, UsageBillingMonetaryScale).
		WillReturnRows(sqlmock.NewRows([]string{"malformed", "unproven", "updated"}).AddRow(int64(1), int64(0), int64(0)))
	mock.ExpectRollback()

	_, err := svc.resetTodayUsageTransaction(context.Background(), 7, keyHash)

	require.ErrorContains(t, err, "malformed current-day self-reset operations")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetTodayUsageTransaction_RollsBackOnInvalidLegacyTimestamp(t *testing.T) {
	svc, mock := newResetTodaySQLMockService(t)
	keyHash := HashIdempotencyKey("rollback-on-invalid-legacy-timestamp")
	dayStart := resetTodayTestDayStart(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO subscription_quota_reset_operations`).
		WithArgs(int64(7), keyHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(91)))
	mock.ExpectExec(`SELECT pg_advisory_xact_lock`).
		WithArgs(resetTodayUsageAdvisoryLockID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`(?s)SELECT id\s+FROM user_subscriptions.*FOR UPDATE`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(101)))
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*MAX\(updated_at\).*interval '1 microsecond'.*FROM reset_clock`).
		WithArgs(SubscriptionStatusActive).
		WillReturnRows(sqlmock.NewRows([]string{"captured_at", "cutoff_at"}).
			AddRow(resetTodayTestCapturedAt(), resetTodayTestCutoff()))
	mock.ExpectQuery(`(?s)WITH subscription_operations AS MATERIALIZED.*parsed_self_reset_operations AS MATERIALIZED.*SELECT`).
		WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC(), BillingTypeSubscription, UsageBillingMonetaryScale).
		WillReturnError(errors.New("invalid input syntax for type timestamp with time zone"))
	mock.ExpectRollback()

	_, err := svc.resetTodayUsageTransaction(context.Background(), 7, keyHash)

	require.ErrorContains(t, err, "backfill legacy reset-today operations")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetTodayUsageTransaction_RollsBackOnUnprovenOrAmbiguousLegacyProof(t *testing.T) {
	for index, name := range []string{
		"mixed_with_unproven_legacy_operation",
		"both_proofs_with_ambiguous_subwindow",
	} {
		t.Run(name, func(t *testing.T) {
			svc, mock := newResetTodaySQLMockService(t)
			keyHash := HashIdempotencyKey(name)
			dayStart := resetTodayTestDayStart(t)

			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO subscription_quota_reset_operations`).
				WithArgs(int64(7), keyHash).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(100 + index)))
			mock.ExpectExec(`SELECT pg_advisory_xact_lock`).
				WithArgs(resetTodayUsageAdvisoryLockID).
				WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectQuery(`(?s)SELECT id\s+FROM user_subscriptions.*FOR UPDATE`).
				WithArgs(SubscriptionStatusActive).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(101)))
			mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*MAX\(updated_at\).*interval '1 microsecond'.*FROM reset_clock`).
				WithArgs(SubscriptionStatusActive).
				WillReturnRows(sqlmock.NewRows([]string{"captured_at", "cutoff_at"}).
					AddRow(resetTodayTestCapturedAt(), resetTodayTestCutoff()))
			mock.ExpectQuery(`(?s)legacy_usage_proofs AS MATERIALIZED.*unproven_legacy_operations AS MATERIALIZED.*updated_subscriptions AS.*SELECT`).
				WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC(), BillingTypeSubscription, UsageBillingMonetaryScale).
				WillReturnRows(sqlmock.NewRows([]string{"malformed", "unproven", "updated"}).AddRow(int64(0), int64(1), int64(0)))
			mock.ExpectRollback()

			_, err := svc.resetTodayUsageTransaction(context.Background(), 7, keyHash)

			require.ErrorContains(t, err, "unproven or ambiguous legacy self-reset operations")
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBackfillLegacyResetTodayUsageOperations_ExecutesTransitionReconstruction(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	dayStart := resetTodayTestDayStart(t)

	mock.ExpectQuery(`(?s)legacy_operations AS MATERIALIZED.*reconstruction_windows AS MATERIALIZED.*reconstructed_operations AS MATERIALIZED.*updated_subscriptions AS.*UPDATE user_subscriptions AS us.*SELECT`).
		WithArgs(SubscriptionStatusActive, resetTodayTestCapturedAt(), dayStart.UTC(), BillingTypeSubscription, UsageBillingMonetaryScale).
		WillReturnRows(sqlmock.NewRows([]string{"malformed", "unproven", "updated"}).AddRow(int64(0), int64(0), int64(1)))

	err = backfillLegacyResetTodayUsageOperations(context.Background(), db, resetTodayTestCapturedAt(), dayStart)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestLegacyResetTodayUsageBackfillSQLGuardsAndBoundaries(t *testing.T) {
	source, err := os.ReadFile("subscription_reset_today.go")
	require.NoError(t, err)
	text := string(source)

	// Native records already have both cleared fields and must never be
	// reconstructed. Likewise, a legacy record consumed by an earlier admin
	// operation must not be charged a second time.
	require.Contains(t, text, `WHERE NOT (valid.value ?| ARRAY[`)
	require.Contains(t, text, `'cleared_daily_window_start', 'cleared_daily_usage_usd'`)
	require.Contains(t, text, "AND NOT valid.value ? 'admin_reset_consumed_operation_id'")

	// Only the complete legacy schema is eligible; the original array order is
	// retained when cleared_* is added to a qualifying operation.
	require.Contains(t, text, `'idempotency_key_hash', 'subscription_id', 'group_id', 'reset_at',`)
	require.Contains(t, text, `AND jsonb_typeof(parsed.value->'subscription_id') = 'number'`)
	require.NotContains(t, text, `pg_input_is_valid`)
	require.Contains(t, text, `AND jsonb_typeof(operation.value->'reset_at') = 'string'`)
	require.Contains(t, text, `THEN (operation.value->>'reset_at')::timestamptz`)
	require.Contains(t, text, `THEN (operation.value->>'replay_expires_at')::timestamptz`)
	require.Contains(t, text, `AND parsed.value->>'idempotency_key_hash' ~ '^[0-9a-f]{64}$'`)
	require.Contains(t, text, `AND (parsed.value->>'daily_usage_usd')::numeric = 0`)
	require.Contains(t, text, `AND jsonb_typeof(parsed.value->'daily_quota_reset_available_at') IN ('string', 'null')`)
	require.Contains(t, text, `AND jsonb_typeof(parsed.value->'replay_expires_at') = 'string'`)
	require.Contains(t, text, `AND parsed.replay_expires_at >= parsed.reset_at`)
	require.Contains(t, text, `AND parsed.value->'subscription_id' = to_jsonb(parsed.subscription_id)`)
	require.Contains(t, text, `AND parsed.value->'group_id' = to_jsonb(parsed.group_id)`)
	require.Contains(t, text, `'cleared_daily_window_start', reconstructed.cleared_window_start`)
	require.Contains(t, text, `'cleared_daily_usage_usd', reconstructed.cleared_usage_usd`)
	require.Contains(t, text, `'legacy_cleared_usage_reconstructed_at', $2`)
	require.Contains(t, text, `reconstructed.reconstruction_source`)
	require.Contains(t, text, `ORDER BY operation.position`)

	// The reconstruction window begins after every boundary that could have
	// already cleared usage, and the correlated ledger scan matches the existing
	// (subscription_id, created_at) index prefix/range.
	require.Contains(t, text, `GREATEST(
					$3::timestamptz,
					legacy.starts_at`)
	require.Contains(t, text, `SELECT MAX(previous.reset_at)`)
	require.Contains(t, text, `AND previous.position < legacy.position`)
	require.Contains(t, text, `SELECT MAX(admin_operation.cutoff_at)`)
	require.Contains(t, text, `WHERE status = 'succeeded'`)
	require.Contains(t, text, `AND admin_operation.cutoff_at < legacy.reset_at`)
	require.Contains(t, text, `) @> jsonb_build_array(jsonb_build_object(`)
	require.Contains(t, text, `'user_id', legacy.user_id`)
	require.Contains(t, text, `'group_id', legacy.group_id`)
	require.NotContains(t, text, `FROM reconstruction_windows AS window`)
	require.Contains(t, text, `FROM reconstruction_windows AS reconstruction`)
	require.Contains(t, text, `WHERE usage_log.subscription_id = proof_window.subscription_id`)
	require.Contains(t, text, `AND usage_log.billing_type = $4`)
	require.Equal(t, 1, strings.Count(text, `FROM usage_logs AS usage_log`))
	require.Contains(t, text, `SUM(GREATEST(usage_log.actual_cost, 0)) FILTER (`)
	require.Contains(t, text, `SUM(ROUND(GREATEST(usage_log.actual_cost, 0), $5)) FILTER (`)
	require.Contains(t, text, `AND usage_log.created_at >= LEAST(`)
	require.Contains(t, text, `AND usage_log.created_at < proof_window.reset_at`)
	require.Contains(t, text, `evidence.weekly_raw10_usd = evidence.stored_weekly_usage_usd`)
	require.Contains(t, text, `evidence.monthly_raw10_usd = evidence.stored_monthly_usage_usd`)
	require.Contains(t, text, `evidence.weekly_per_row_round8_usd = evidence.stored_weekly_usage_usd`)
	require.Contains(t, text, `evidence.monthly_per_row_round8_usd = evidence.stored_monthly_usage_usd`)
	require.Contains(t, text, `admin_operation.cutoff_at >= reconstruction.weekly_window_start`)
	require.Contains(t, text, `admin_operation.cutoff_at >= reconstruction.monthly_window_start`)
	require.Contains(t, text, `unproven_legacy_operations AS MATERIALIZED (`)
	require.Contains(t, text, `proof.subwindow_raw10_usd <> proof.subwindow_per_row_round8_usd`)
	require.Contains(t, text, `AND NOT EXISTS (SELECT 1 FROM unproven_legacy_operations)`)
	require.Contains(t, text, `unproven or ambiguous legacy self-reset operations`)

	// PostgreSQL 15 has no pg_input_is_valid. String timestamps are cast directly,
	// so invalid input aborts the SQL statement and the surrounding transaction.
	// Non-string reset_at values, malformed legacy shapes, and partial native
	// ledgers are counted explicitly. The update runs only when that count is zero.
	require.Contains(t, text, `malformed_self_reset_operations AS MATERIALIZED (`)
	require.Contains(t, text, `jsonb_typeof(parsed.value->'reset_at') IS DISTINCT FROM 'string'`)
	require.Contains(t, text, `NOT (parsed.value ?| ARRAY[`)
	require.Contains(t, text, `FROM valid_self_reset_operations AS valid`)
	require.Contains(t, text, `AND valid.position = parsed.position`)
	require.Contains(t, text, `parsed.value ?| ARRAY[`)
	require.Contains(t, text, `AND parsed.cleared_window_start IS NOT NULL`)
	require.Contains(t, text, `AND jsonb_typeof(parsed.value->'cleared_daily_usage_usd') = 'number'`)
	require.Contains(t, text, `AND EXISTS (`)
	require.Contains(t, text, `AND NOT EXISTS (SELECT 1 FROM malformed_self_reset_operations)`)
	require.Contains(t, text, `refuse reset-today with %d malformed current-day self-reset operations`)
}

func TestLegacyResetTodayUsageProofSQLSupportsPureRaw10History(t *testing.T) {
	source, err := os.ReadFile("subscription_reset_today.go")
	require.NoError(t, err)
	text := string(source)

	require.Contains(t, text, `WHEN proof.raw10_proven THEN proof.subwindow_raw10_usd`)
	require.Contains(t, text, `WHEN proof.raw10_proven
					THEN 'period_snapshot_proof=raw10'`)
}

func TestLegacyResetTodayUsageProofSQLSupportsPurePerRowRound8History(t *testing.T) {
	source, err := os.ReadFile("subscription_reset_today.go")
	require.NoError(t, err)
	text := string(source)

	require.Contains(t, text, `ELSE proof.subwindow_per_row_round8_usd`)
	require.Contains(t, text, `ELSE 'period_snapshot_proof=per_row_round8'`)
	require.Contains(t, text, `THEN 'period_snapshot_proof=raw10+per_row_round8/subwindow_equal'`)
}

func TestResetTodayUsageSQLPreservesWindowStartsAndClampsUsage(t *testing.T) {
	source, err := os.ReadFile("subscription_reset_today.go")
	require.NoError(t, err)
	text := string(source)
	require.Contains(t, text, "weekly_usage_usd = GREATEST(")
	require.Contains(t, text, "monthly_usage_usd = GREATEST(")
	require.Contains(t, text, "GREATEST(us.weekly_usage_usd, 0)")
	require.Contains(t, text, "GREATEST(us.monthly_usage_usd, 0)")
	require.Contains(t, text, "NOT operation.value ? 'admin_reset_consumed_operation_id'")
	require.Contains(t, text, "'admin_reset_consumed_operation_id', $6::bigint")
	require.Contains(t, text, "'admin_reset_consumed_at', $2")
	require.Contains(t, text, "(operation.value->>'reset_at')::timestamptz >= $3")
	require.Contains(t, text, "(operation.value->>'cleared_daily_window_start')::timestamptz >= $3")
	require.Contains(t, text, "(operation.value->>'cleared_daily_window_start')::timestamptz <= $2")
	require.NotContains(t, text, "SET daily_window_start =")
	require.NotContains(t, text, "SET weekly_window_start =")
	require.NotContains(t, text, "SET monthly_window_start =")
}

func TestResetTodayUsageSQLCapsEachPeriodWindowDeductionIndependently(t *testing.T) {
	source, err := os.ReadFile("subscription_reset_today.go")
	require.NoError(t, err)
	text := string(source)
	require.Contains(t, text, `us.weekly_usage_usd - LEAST(
						adjustments.deduct_usd,
						GREATEST(us.weekly_usage_usd, 0)`)
	require.Contains(t, text, `us.monthly_usage_usd - LEAST(
						adjustments.deduct_usd,
						GREATEST(us.monthly_usage_usd, 0)`)

	capDeduction := func(todayLedger, currentWindowUsage float64) float64 {
		if todayLedger < 0 {
			todayLedger = 0
		}
		if currentWindowUsage < 0 {
			currentWindowUsage = 0
		}
		if todayLedger < currentWindowUsage {
			return todayLedger
		}
		return currentWindowUsage
	}

	const todayLedger = 10 + 3
	tests := []struct {
		name                  string
		weeklyUsage           float64
		monthlyUsage          float64
		expectedWeeklyDeduct  float64
		expectedMonthlyDeduct float64
	}{
		{
			name:                  "weekly rolled over after ledger ten and new usage three",
			weeklyUsage:           3,
			monthlyUsage:          23,
			expectedWeeklyDeduct:  3,
			expectedMonthlyDeduct: 13,
		},
		{
			name:                  "monthly rolled over while weekly still contains the full day",
			weeklyUsage:           23,
			monthlyUsage:          3,
			expectedWeeklyDeduct:  13,
			expectedMonthlyDeduct: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedWeeklyDeduct, capDeduction(todayLedger, tt.weeklyUsage))
			require.Equal(t, tt.expectedMonthlyDeduct, capDeduction(todayLedger, tt.monthlyUsage))
		})
	}
}

func TestCaptureResetTodayUsageTargetsCombinesOnlyEligibleTodayUsage(t *testing.T) {
	dayStart := resetTodayTestDayStart(t)
	capturedAt := resetTodayTestCapturedAt()

	tests := []struct {
		name                 string
		windowStart          time.Time
		dailyUsage           string
		unconsumedResetUsage string
		expectedDeduction    string
		expectedTodayWindows int
		expectedResetCount   int
		expectedResetTotal   float64
		expectedStaleTotal   float64
	}{
		{
			name:                 "self reset ten plus three current usage deducts thirteen",
			windowStart:          dayStart,
			dailyUsage:           "3",
			unconsumedResetUsage: "10",
			expectedDeduction:    "13",
			expectedTodayWindows: 1,
			expectedResetCount:   1,
			expectedResetTotal:   13,
		},
		{
			name:                 "second admin reset does not deduct consumed operation again",
			windowStart:          dayStart,
			dailyUsage:           "0",
			unconsumedResetUsage: "0",
			expectedDeduction:    "0",
			expectedTodayWindows: 1,
		},
		{
			name:                 "usage after admin reset deducts only the new amount",
			windowStart:          dayStart,
			dailyUsage:           "2.75",
			unconsumedResetUsage: "0",
			expectedDeduction:    "2.75",
			expectedTodayWindows: 1,
			expectedResetCount:   1,
			expectedResetTotal:   2.75,
		},
		{
			name:                 "stale daily window is cleared without weekly monthly deduction",
			windowStart:          dayStart.Add(-24 * time.Hour),
			dailyUsage:           "7.5",
			unconsumedResetUsage: "0",
			expectedDeduction:    "0",
			expectedStaleTotal:   7.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() { _ = db.Close() }()

			mock.ExpectQuery(`SELECT id, user_id, group_id, daily_window_start, daily_usage_usd::text`).
				WithArgs(SubscriptionStatusActive, capturedAt, dayStart.UTC()).
				WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "group_id", "daily_window_start", "daily_usage_usd", "unconsumed_reset_usage_usd"}).
					AddRow(int64(101), int64(11), int64(21), tt.windowStart, tt.dailyUsage, tt.unconsumedResetUsage))

			result, targets, adjustments, err := captureResetTodayUsageTargets(context.Background(), db, capturedAt, dayStart)

			require.NoError(t, err)
			require.Equal(t, 1, result.Subscriptions)
			require.Equal(t, tt.expectedTodayWindows, result.TodayWindowSubscriptions)
			require.Equal(t, tt.expectedResetCount, result.ResetSubscriptions)
			require.InDelta(t, tt.expectedResetTotal, result.ResetAndDeductedUSD, 0.0000000001)
			require.InDelta(t, tt.expectedStaleTotal, result.StaleDailyClearedUSD, 0.0000000001)
			require.Equal(t, []resetTodayUsageTarget{{UserID: 11, GroupID: 21}}, targets)
			require.Len(t, adjustments, 1)
			require.Equal(t, int64(101), adjustments[0].ID)
			require.Equal(t, tt.expectedDeduction, adjustments[0].DeductUSD.String())
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestClaimResetTodayUsageOperation_ConcurrentDuplicatesHaveSingleOwner(t *testing.T) {
	db, err := stdsql.Open("sqlite", "file:reset-today-claim?mode=memory&cache=shared")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	// A single shared connection makes the assertion deterministic while the
	// goroutines still contend at the operation boundary used by the handler.
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`
		CREATE TABLE subscription_quota_reset_operations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			actor_user_id INTEGER NOT NULL,
			idempotency_key_hash TEXT NOT NULL,
			status TEXT NOT NULL,
			UNIQUE(actor_user_id, idempotency_key_hash)
		)
	`)
	require.NoError(t, err)

	const callers = 8
	keyHash := HashIdempotencyKey("concurrent-duplicate")
	start := make(chan struct{})
	type claimResult struct {
		claimed bool
		err     error
	}
	results := make(chan claimResult, callers)
	for i := 0; i < callers; i++ {
		go func() {
			<-start
			_, claimed, claimErr := claimResetTodayUsageOperation(context.Background(), db, 7, keyHash)
			results <- claimResult{claimed: claimed, err: claimErr}
		}()
	}
	close(start)

	owners := 0
	for i := 0; i < callers; i++ {
		result := <-results
		require.NoError(t, result.err)
		if result.claimed {
			owners++
		}
	}
	require.Equal(t, 1, owners)

	var rows int
	require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM subscription_quota_reset_operations`).Scan(&rows))
	require.Equal(t, 1, rows)
}
