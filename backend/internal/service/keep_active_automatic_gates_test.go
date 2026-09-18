//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/stretchr/testify/require"
)

// keepActiveGeminiUsageRepo returns a fully consumed local quota while also
// recording whether a precheck queried usage at all. Protected accounts must
// bypass the query, not merely ignore its result.
type keepActiveGeminiUsageRepo struct {
	usageBatchLogRepoStub
	calls int
}

func (r *keepActiveGeminiUsageRepo) GetModelStatsWithFilters(
	_ context.Context,
	_ time.Time,
	_ time.Time,
	_ int64,
	_ int64,
	_ int64,
	_ int64,
	_ *int16,
	_ *bool,
	_ *int8,
) ([]usagestats.ModelStat, error) {
	r.calls++
	return []usagestats.ModelStat{{Model: "gemini-3-pro", Requests: 50}}, nil
}

func keepActiveGeminiAccount(id int64) *Account {
	return &Account{
		ID:       id,
		Platform: PlatformGemini,
		Type:     AccountTypeAPIKey,
		Extra: map[string]any{
			KeepStatusActiveExtraKey: true,
		},
	}
}

func TestRateLimitService_PreCheckUsage_KeepActiveBypassesLocalGeminiQuota(t *testing.T) {
	usage := &keepActiveGeminiUsageRepo{}
	svc := NewRateLimitService(nil, usage, &config.Config{}, NewGeminiQuotaService(&config.Config{}, nil), nil)

	ok, err := svc.PreCheckUsage(context.Background(), keepActiveGeminiAccount(7101), "gemini-3-pro")
	require.NoError(t, err)
	require.True(t, ok)
	require.Zero(t, usage.calls, "protected accounts must not query or apply local RPD/RPM prechecks")
}

func TestRateLimitService_PreCheckUsageBatch_KeepActiveBypassesLocalGeminiQuota(t *testing.T) {
	usage := &keepActiveGeminiUsageRepo{}
	svc := NewRateLimitService(nil, usage, &config.Config{}, NewGeminiQuotaService(&config.Config{}, nil), nil)
	account := keepActiveGeminiAccount(7102)

	result, err := svc.PreCheckUsageBatch(context.Background(), []*Account{account}, "gemini-3-pro")
	require.NoError(t, err)
	require.True(t, result[account.ID])
	require.Zero(t, usage.calls, "protected accounts must not enter the batch quota query")
}

func TestRateLimitService_CNProtectedAccountSkipsAutomaticBalancePenalty(t *testing.T) {
	account := &Account{
		ID:       7103,
		Platform: PlatformKimi,
		Extra:    map[string]any{KeepStatusActiveExtraKey: true},
	}
	ratelimit := &RateLimitService{}

	require.NotPanics(t, func() {
		ratelimit.handleCNProviderInsufficientBalance(context.Background(), account, "insufficient balance")
	})
	require.True(t, ratelimit.applyCNProviderReactive429(context.Background(), account, nil, []byte(`{"error":"insufficient balance"}`)))

	checker := &CNProviderBalanceCheckService{}
	require.Equal(t, cnBalanceNoChange, checker.checkOne(context.Background(), account, 10))
}

func TestRateLimitService_GetTempUnschedStatusHandlesNilRepositoryAccount(t *testing.T) {
	repo := &rateLimitClearRepoStub{getByIDAccount: nil}
	svc := NewRateLimitService(repo, nil, &config.Config{}, nil, nil)

	state, err := svc.GetTempUnschedStatus(context.Background(), 7104)
	require.Nil(t, state)
	require.ErrorIs(t, err, ErrAccountNotFound)
}

func TestAntigravityInternal500Penalty_KeepActiveSkipsAutomaticStateChanges(t *testing.T) {
	repo := &internal500AccountRepoStub{}
	cache := &mockInternal500Cache{incrementCount: 3}
	svc := &AntigravityGatewayService{accountRepo: repo, internal500Cache: cache}
	account := &Account{
		ID:    7105,
		Name:  "protected-antigravity",
		Extra: map[string]any{KeepStatusActiveExtraKey: true},
	}

	svc.applyInternal500Penalty(context.Background(), "[test]", account, 3)
	svc.handleInternal500RetryExhausted(context.Background(), "[test]", account)

	require.Empty(t, repo.tempUnschedCalls)
	require.Empty(t, repo.setErrorCalls)
	require.Empty(t, cache.incrementCalls, "protected accounts must not even advance the automatic penalty counter")
}
