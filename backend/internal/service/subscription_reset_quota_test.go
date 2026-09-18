//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/stretchr/testify/require"
)

// resetQuotaUserSubRepoStub 支持 GetByID、ResetUsageWindows，
// 其余方法继承 userSubRepoNoop（panic）。
type resetQuotaUserSubRepoStub struct {
	userSubRepoNoop

	sub *UserSubscription

	resetDailyCalled   bool
	resetWeeklyCalled  bool
	resetMonthlyCalled bool
	resetDailyErr      error
	resetWeeklyErr     error
	resetMonthlyErr    error
	dailyStart         time.Time
	periodicStart      time.Time
	resetCutoff        time.Time
	refreshErr         error
	getByIDCalls       int
}

func (r *resetQuotaUserSubRepoStub) GetByID(_ context.Context, id int64) (*UserSubscription, error) {
	r.getByIDCalls++
	if r.getByIDCalls > 1 && r.refreshErr != nil {
		return nil, r.refreshErr
	}
	if r.sub == nil || r.sub.ID != id {
		return nil, ErrSubscriptionNotFound
	}
	cp := *r.sub
	return &cp, nil
}

func (r *resetQuotaUserSubRepoStub) ResetUsageWindows(_ context.Context, _ int64, resetDaily, resetWeekly, resetMonthly bool, dailyStart, periodicStart time.Time) (time.Time, error) {
	r.resetDailyCalled = resetDaily
	r.resetWeeklyCalled = resetWeekly
	r.resetMonthlyCalled = resetMonthly
	r.dailyStart = dailyStart
	r.periodicStart = periodicStart
	if resetDaily && r.resetDailyErr != nil {
		return time.Time{}, r.resetDailyErr
	}
	if resetWeekly && r.resetWeeklyErr != nil {
		return time.Time{}, r.resetWeeklyErr
	}
	if resetMonthly && r.resetMonthlyErr != nil {
		return time.Time{}, r.resetMonthlyErr
	}
	if r.sub == nil {
		return periodicStart, nil
	}
	if resetDaily {
		r.sub.DailyUsageUSD = 0
		r.sub.DailyWindowStart = &dailyStart
	}
	if resetWeekly {
		r.sub.WeeklyUsageUSD = 0
		r.sub.WeeklyWindowStart = &periodicStart
	}
	if resetMonthly {
		r.sub.MonthlyUsageUSD = 0
		r.sub.MonthlyWindowStart = &periodicStart
	}
	cutoff := r.resetCutoff
	if cutoff.IsZero() {
		cutoff = periodicStart
	}
	r.sub.UpdatedAt = cutoff
	return cutoff, nil
}

func (r *resetQuotaUserSubRepoStub) ResetDailyUsage(_ context.Context, _ int64, _ *time.Time, windowStart time.Time) error {
	r.resetDailyCalled = true
	if r.resetDailyErr == nil && r.sub != nil {
		r.sub.DailyUsageUSD = 0
		r.sub.DailyWindowStart = &windowStart
	}
	return r.resetDailyErr
}

func (r *resetQuotaUserSubRepoStub) ResetWeeklyUsage(_ context.Context, _ int64, _ *time.Time, _ time.Time) error {
	r.resetWeeklyCalled = true
	return r.resetWeeklyErr
}

func (r *resetQuotaUserSubRepoStub) ResetMonthlyUsage(_ context.Context, _ int64, _ *time.Time, _ time.Time) error {
	r.resetMonthlyCalled = true
	return r.resetMonthlyErr
}

func newResetQuotaSvc(stub *resetQuotaUserSubRepoStub) *SubscriptionService {
	return NewSubscriptionService(groupRepoNoop{}, stub, nil, nil, nil)
}

func TestAdminResetQuota_ResetBoth(t *testing.T) {
	useUTCGlobalAndShanghaiSubscriptionCalendar(t)
	stub := &resetQuotaUserSubRepoStub{
		sub: &UserSubscription{ID: 1, UserID: 10, GroupID: 20},
	}
	svc := newResetQuotaSvc(stub)
	resetAt := time.Date(2026, 7, 1, 10, 37, 42, 123, time.UTC)
	svc.now = func() time.Time { return resetAt }

	result, err := svc.AdminResetQuota(context.Background(), 1, true, true, false)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, stub.resetDailyCalled, "应调用 ResetDailyUsage")
	require.True(t, stub.resetWeeklyCalled, "应调用 ResetWeeklyUsage")
	require.False(t, stub.resetMonthlyCalled, "不应调用 ResetMonthlyUsage")
	// 手动重置后日窗口锚定当天 0 点（保持 0 点刷新节奏），周窗口锚定重置时刻。
	require.Equal(t, subscriptionDailyCalendarStart(resetAt), stub.dailyStart)
	require.Equal(t, resetAt, stub.periodicStart)
	require.Equal(t, subscriptionDailyCalendarStart(resetAt), *result.DailyWindowStart)
	require.Equal(t, resetAt, *result.WeeklyWindowStart)
}

func TestAdminResetQuota_ResetDailyOnly(t *testing.T) {
	stub := &resetQuotaUserSubRepoStub{
		sub: &UserSubscription{ID: 2, UserID: 10, GroupID: 20},
	}
	svc := newResetQuotaSvc(stub)

	result, err := svc.AdminResetQuota(context.Background(), 2, true, false, false)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, stub.resetDailyCalled, "应调用 ResetDailyUsage")
	require.False(t, stub.resetWeeklyCalled, "不应调用 ResetWeeklyUsage")
	require.False(t, stub.resetMonthlyCalled, "不应调用 ResetMonthlyUsage")
}

func TestAdminResetQuota_ResetWeeklyOnly(t *testing.T) {
	stub := &resetQuotaUserSubRepoStub{
		sub: &UserSubscription{ID: 3, UserID: 10, GroupID: 20},
	}
	svc := newResetQuotaSvc(stub)

	result, err := svc.AdminResetQuota(context.Background(), 3, false, true, false)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, stub.resetDailyCalled, "不应调用 ResetDailyUsage")
	require.True(t, stub.resetWeeklyCalled, "应调用 ResetWeeklyUsage")
	require.False(t, stub.resetMonthlyCalled, "不应调用 ResetMonthlyUsage")
}

func TestAdminResetQuota_BothFalseReturnsError(t *testing.T) {
	stub := &resetQuotaUserSubRepoStub{
		sub: &UserSubscription{ID: 7, UserID: 10, GroupID: 20},
	}
	svc := newResetQuotaSvc(stub)

	_, err := svc.AdminResetQuota(context.Background(), 7, false, false, false)

	require.ErrorIs(t, err, ErrInvalidInput)
	require.False(t, stub.resetDailyCalled)
	require.False(t, stub.resetWeeklyCalled)
	require.False(t, stub.resetMonthlyCalled)
}

func TestAdminResetQuota_SubscriptionNotFound(t *testing.T) {
	stub := &resetQuotaUserSubRepoStub{sub: nil}
	svc := newResetQuotaSvc(stub)

	_, err := svc.AdminResetQuota(context.Background(), 999, true, true, true)

	require.ErrorIs(t, err, ErrSubscriptionNotFound)
	require.False(t, stub.resetDailyCalled)
	require.False(t, stub.resetWeeklyCalled)
	require.False(t, stub.resetMonthlyCalled)
}

func TestAdminResetQuota_ResetDailyUsageError(t *testing.T) {
	dbErr := errors.New("db error")
	stub := &resetQuotaUserSubRepoStub{
		sub:           &UserSubscription{ID: 4, UserID: 10, GroupID: 20},
		resetDailyErr: dbErr,
	}
	svc := newResetQuotaSvc(stub)

	_, err := svc.AdminResetQuota(context.Background(), 4, true, true, false)

	require.ErrorIs(t, err, dbErr)
	require.True(t, stub.resetDailyCalled)
	require.True(t, stub.resetWeeklyCalled, "原子重置应在一次调用中提交所选窗口")
}

func TestAdminResetQuota_ResetWeeklyUsageError(t *testing.T) {
	dbErr := errors.New("db error")
	stub := &resetQuotaUserSubRepoStub{
		sub:            &UserSubscription{ID: 5, UserID: 10, GroupID: 20},
		resetWeeklyErr: dbErr,
	}
	svc := newResetQuotaSvc(stub)

	_, err := svc.AdminResetQuota(context.Background(), 5, false, true, false)

	require.ErrorIs(t, err, dbErr)
	require.True(t, stub.resetWeeklyCalled)
}

func TestAdminResetQuota_ResetMonthlyOnly(t *testing.T) {
	stub := &resetQuotaUserSubRepoStub{
		sub: &UserSubscription{ID: 8, UserID: 10, GroupID: 20},
	}
	svc := newResetQuotaSvc(stub)

	result, err := svc.AdminResetQuota(context.Background(), 8, false, false, true)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, stub.resetDailyCalled, "不应调用 ResetDailyUsage")
	require.False(t, stub.resetWeeklyCalled, "不应调用 ResetWeeklyUsage")
	require.True(t, stub.resetMonthlyCalled, "应调用 ResetMonthlyUsage")
}

func TestAdminResetQuota_BeforeStartsAtSameDayPreservesAutomaticBoundary(t *testing.T) {
	startsAt := time.Date(2026, 7, 1, 15, 0, 0, 0, time.UTC)
	resetAt := time.Date(2026, 7, 1, 10, 37, 42, 123, time.UTC)
	stub := &resetQuotaUserSubRepoStub{
		sub: &UserSubscription{
			ID:        10,
			UserID:    10,
			GroupID:   20,
			StartsAt:  startsAt,
			ExpiresAt: startsAt.Add(45 * 24 * time.Hour),
		},
	}
	svc := newResetQuotaSvc(stub)
	svc.now = func() time.Time { return resetAt }

	result, err := svc.AdminResetQuota(context.Background(), 10, false, false, true)

	require.NoError(t, err)
	require.Equal(t, resetAt, *result.MonthlyWindowStart)
	boundary, ok := result.automaticWindowStartAt(result.MonthlyWindowStart, 30*24*time.Hour, resetAt.Add(30*24*time.Hour))
	require.True(t, ok)
	require.Equal(t, resetAt.Add(30*24*time.Hour), boundary)
}

func TestAdminResetQuota_ResetMonthlyUsageError(t *testing.T) {
	dbErr := errors.New("db error")
	stub := &resetQuotaUserSubRepoStub{
		sub:             &UserSubscription{ID: 9, UserID: 10, GroupID: 20},
		resetMonthlyErr: dbErr,
	}
	svc := newResetQuotaSvc(stub)

	_, err := svc.AdminResetQuota(context.Background(), 9, false, false, true)

	require.ErrorIs(t, err, dbErr)
	require.True(t, stub.resetMonthlyCalled)
}

func TestAdminResetQuota_ReturnsRefreshedSub(t *testing.T) {
	stub := &resetQuotaUserSubRepoStub{
		sub: &UserSubscription{
			ID:            6,
			UserID:        10,
			GroupID:       20,
			DailyUsageUSD: 99.9,
		},
	}

	svc := newResetQuotaSvc(stub)
	result, err := svc.AdminResetQuota(context.Background(), 6, true, false, false)

	require.NoError(t, err)
	// ResetUsageWindows stub 会将 sub.DailyUsageUSD 归零，
	// 服务应返回第二次 GetByID 的刷新值而非初始的 99.9
	require.Equal(t, float64(0), result.DailyUsageUSD, "返回的订阅应反映已归零的用量")
	require.True(t, stub.resetDailyCalled)
}

func TestAdminResetQuotaUsesDatabaseCutoffBarrierPublishesAndInvalidatesLocalL1(t *testing.T) {
	appTime := time.Date(2026, 8, 19, 7, 40, 1, 987654321, time.UTC)
	databaseCutoff := time.Date(2026, 8, 19, 7, 39, 55, 123457000, time.UTC)
	stub := &resetQuotaUserSubRepoStub{
		sub: &UserSubscription{
			ID:              51,
			UserID:          77,
			GroupID:         19,
			DailyUsageUSD:   2.5,
			WeeklyUsageUSD:  12.5,
			MonthlyUsageUSD: 31.75,
		},
		resetCutoff: databaseCutoff,
	}
	redisCache := &userDailyQuotaResetCacheStub{}
	svc := NewSubscriptionService(nil, stub, &BillingCacheService{cache: redisCache}, nil, nil)
	svc.now = func() time.Time { return appTime }

	l1, err := ristretto.NewCache(&ristretto.Config{NumCounters: 100, MaxCost: 100, BufferItems: 64})
	require.NoError(t, err)
	t.Cleanup(l1.Close)
	svc.subCacheL1 = l1
	cacheKey := subCacheKey(77, 19)
	require.True(t, l1.Set(cacheKey, &UserSubscription{ID: 51, DailyUsageUSD: 99}, 1))
	l1.Wait()
	_, found := l1.Get(cacheKey)
	require.True(t, found)

	result, err := svc.AdminResetQuota(context.Background(), 51, true, false, false)

	require.NoError(t, err)
	require.Equal(t, databaseCutoff, result.UpdatedAt)
	require.Zero(t, result.DailyUsageUSD)
	_, found = l1.Get(cacheKey)
	require.False(t, found, "local L1 entry must be synchronously invalidated before returning")
	barriers, publishes, cutoffs := redisCache.calls()
	require.Equal(t, 1, barriers)
	require.Equal(t, 1, publishes, "other nodes must receive the L1 invalidation event")
	require.Equal(t, []time.Time{databaseCutoff}, cutoffs)
	require.NotEqual(t, appTime, cutoffs[0], "the cache barrier must use the database RETURNING cutoff")
}

func TestAdminResetQuotaCacheFailureDoesNotFailCommittedResetAndRetries(t *testing.T) {
	previousDelays := append([]time.Duration(nil), subscriptionQuotaCacheRetryDelays...)
	subscriptionQuotaCacheRetryDelays = []time.Duration{time.Millisecond, 2 * time.Millisecond}
	t.Cleanup(func() { subscriptionQuotaCacheRetryDelays = previousDelays })

	databaseCutoff := time.Date(2026, 8, 19, 7, 45, 2, 222223000, time.UTC)
	stub := &resetQuotaUserSubRepoStub{
		sub:         &UserSubscription{ID: 52, UserID: 77, GroupID: 19, WeeklyUsageUSD: 12.5},
		resetCutoff: databaseCutoff,
	}
	redisCache := &userDailyQuotaResetCacheStub{barrierFailures: 1}
	svc := NewSubscriptionService(nil, stub, &BillingCacheService{cache: redisCache}, nil, nil)

	result, err := svc.AdminResetQuota(context.Background(), 52, false, true, false)

	require.NoError(t, err, "a committed reset must not be exposed as a retryable write failure")
	require.NotNil(t, result)
	require.Zero(t, result.WeeklyUsageUSD)
	require.Eventually(t, func() bool {
		barriers, publishes, cutoffs := redisCache.calls()
		return barriers >= 2 && publishes >= 2 && len(cutoffs) >= 2
	}, time.Second, time.Millisecond)
	_, _, cutoffs := redisCache.calls()
	for _, cutoff := range cutoffs {
		require.Equal(t, databaseCutoff, cutoff)
	}
}

func TestAdminResetQuotaRefreshFailureReturnsCommittedSnapshot(t *testing.T) {
	databaseCutoff := time.Date(2026, 8, 19, 7, 50, 3, 333334000, time.UTC)
	stub := &resetQuotaUserSubRepoStub{
		sub: &UserSubscription{
			ID:              53,
			UserID:          77,
			GroupID:         19,
			MonthlyUsageUSD: 31.75,
		},
		resetCutoff: databaseCutoff,
		refreshErr:  errors.New("post-commit read unavailable"),
	}
	svc := newResetQuotaSvc(stub)

	result, err := svc.AdminResetQuota(context.Background(), 53, false, false, true)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Zero(t, result.MonthlyUsageUSD)
	require.Equal(t, databaseCutoff, result.UpdatedAt)
	require.Equal(t, 2, stub.getByIDCalls)
}
