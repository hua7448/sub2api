//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type accountRepoStubForClearAccountError struct {
	mockAccountRepoForGemini
	account                  *Account
	clearErrorCalls          int
	clearRateLimitCalls      int
	clearAntigravityCalls    int
	clearModelRateLimitCalls int
	clearTempUnschedCalls    int
	updateExtraCalls         int
	lastExtraUpdates         map[string]any
}

func (r *accountRepoStubForClearAccountError) GetByID(ctx context.Context, id int64) (*Account, error) {
	return r.account, nil
}

func (r *accountRepoStubForClearAccountError) ClearError(ctx context.Context, id int64) error {
	r.clearErrorCalls++
	r.account.Status = StatusActive
	r.account.ErrorMessage = ""
	return nil
}

func (r *accountRepoStubForClearAccountError) ClearRateLimit(ctx context.Context, id int64) error {
	r.clearRateLimitCalls++
	r.account.RateLimitedAt = nil
	r.account.RateLimitResetAt = nil
	return nil
}

func (r *accountRepoStubForClearAccountError) ClearAntigravityQuotaScopes(ctx context.Context, id int64) error {
	r.clearAntigravityCalls++
	return nil
}

func (r *accountRepoStubForClearAccountError) ClearModelRateLimits(ctx context.Context, id int64) error {
	r.clearModelRateLimitCalls++
	return nil
}

func (r *accountRepoStubForClearAccountError) ClearTempUnschedulable(ctx context.Context, id int64) error {
	r.clearTempUnschedCalls++
	r.account.TempUnschedulableUntil = nil
	r.account.TempUnschedulableReason = ""
	return nil
}

func (r *accountRepoStubForClearAccountError) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	r.updateExtraCalls++
	r.lastExtraUpdates = updates
	if r.account.Extra == nil {
		r.account.Extra = make(map[string]any)
	}
	for key, value := range updates {
		r.account.Extra[key] = value
	}
	return nil
}

func TestAdminService_ClearAccountError_AlsoClearsRecoverableRuntimeState(t *testing.T) {
	until := time.Now().Add(10 * time.Minute)
	resetAt := time.Now().Add(5 * time.Minute)
	repo := &accountRepoStubForClearAccountError{
		account: &Account{
			ID:                      31,
			Platform:                PlatformOpenAI,
			Type:                    AccountTypeOAuth,
			Status:                  StatusError,
			ErrorMessage:            "refresh failed",
			RateLimitResetAt:        &resetAt,
			TempUnschedulableUntil:  &until,
			TempUnschedulableReason: "missing refresh token",
		},
	}
	blocker := &runtimeBlockRecorder{}
	svc := &adminServiceImpl{accountRepo: repo, runtimeBlocker: blocker}

	updated, err := svc.ClearAccountError(context.Background(), 31)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, 1, repo.clearErrorCalls)
	require.Equal(t, 1, repo.clearRateLimitCalls)
	require.Equal(t, 1, repo.clearAntigravityCalls)
	require.Equal(t, 1, repo.clearModelRateLimitCalls)
	require.Equal(t, 1, repo.clearTempUnschedCalls)
	require.Nil(t, updated.RateLimitResetAt)
	require.Nil(t, updated.TempUnschedulableUntil)
	require.Empty(t, updated.TempUnschedulableReason)
	require.Equal(t, []int64{31}, blocker.clearedIDs)
}

func TestAdminService_SetAccountKeepStatusActive_EnablePersistsFlagAndClearsRuntimeState(t *testing.T) {
	until := time.Now().Add(10 * time.Minute)
	resetAt := time.Now().Add(5 * time.Minute)
	repo := &accountRepoStubForClearAccountError{
		account: &Account{
			ID:                      41,
			Platform:                PlatformOpenAI,
			Type:                    AccountTypeAPIKey,
			Status:                  StatusError,
			ErrorMessage:            "transport failed",
			RateLimitResetAt:        &resetAt,
			TempUnschedulableUntil:  &until,
			TempUnschedulableReason: "connection refused",
			Extra:                   map[string]any{},
		},
	}
	blocker := &runtimeBlockRecorder{}
	svc := &adminServiceImpl{accountRepo: repo, runtimeBlocker: blocker}

	updated, err := svc.SetAccountKeepStatusActive(context.Background(), 41, true)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, 1, repo.updateExtraCalls)
	require.Equal(t, true, repo.lastExtraUpdates[KeepStatusActiveExtraKey])
	require.True(t, updated.KeepStatusActive())
	require.Equal(t, StatusActive, updated.Status)
	require.Empty(t, updated.ErrorMessage)
	require.Nil(t, updated.RateLimitResetAt)
	require.Nil(t, updated.TempUnschedulableUntil)
	require.Equal(t, []int64{41}, blocker.clearedIDs)
}

func TestAdminService_SetAccountKeepStatusActive_DisableOnlyPersistsFlag(t *testing.T) {
	repo := &accountRepoStubForClearAccountError{
		account: &Account{
			ID:       42,
			Status:   StatusActive,
			Extra:    map[string]any{KeepStatusActiveExtraKey: true},
			Platform: PlatformOpenAI,
		},
	}
	blocker := &runtimeBlockRecorder{}
	svc := &adminServiceImpl{accountRepo: repo, runtimeBlocker: blocker}

	updated, err := svc.SetAccountKeepStatusActive(context.Background(), 42, false)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, false, repo.lastExtraUpdates[KeepStatusActiveExtraKey])
	require.False(t, updated.KeepStatusActive())
	require.Zero(t, repo.clearErrorCalls)
	require.Empty(t, blocker.clearedIDs)
}
