//go:build unit

package service

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBulkCodexFingerprintModeValidatesAndPreservesIncrementalUpdates(t *testing.T) {
	for _, mode := range []string{"full", "off", "device", "session"} {
		t.Run(mode, func(t *testing.T) {
			account := &Account{
				ID: 1, Platform: PlatformOpenAI, Type: AccountTypeOAuth,
				Extra: map[string]any{codexFingerprintSeedExtraKey: testCodexFingerprintSeed, "unrelated": "keep"},
			}
			repo := &accountRepoStubForBulkUpdate{getByIDsAccounts: []*Account{account}}
			concurrency := 3
			result, err := (&adminServiceImpl{accountRepo: repo}).BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
				AccountIDs: []int64{1}, Concurrency: &concurrency,
				Extra: map[string]any{codexFingerprintModeExtraKey: mode, "another_setting": true},
			})
			require.NoError(t, err)
			require.True(t, repo.getByIDsCalled, "fingerprint-only writes must validate actual target accounts")
			require.Equal(t, 1, repo.bulkUpdateCalls)
			require.Equal(t, map[string]any{codexFingerprintModeExtraKey: mode, "another_setting": true}, repo.lastBulkUpdate.Extra)
			require.Equal(t, mode != "off", repo.lastBulkUpdate.EnsureCodexFingerprintSeed)
			require.Equal(t, &concurrency, repo.lastBulkUpdate.Concurrency)
			require.NotContains(t, repo.lastBulkUpdate.Extra, codexFingerprintSeedExtraKey, "the existing seed must remain owned by the atomic repository merge")
			require.Equal(t, testCodexFingerprintSeed, account.Extra[codexFingerprintSeedExtraKey])
			require.Equal(t, "keep", account.Extra["unrelated"])
			require.Equal(t, []int64{1}, result.SuccessIDs)
			require.Empty(t, result.FailedIDs)
			require.Equal(t, []BulkUpdateAccountResult{{AccountID: 1, Success: true}}, result.Results)
		})
	}
}

func TestBulkCodexFingerprintModeRejectsInvalidValuesBeforeWrites(t *testing.T) {
	for _, mode := range []any{nil, true, 3, "", "invalid", "FULL", " full ", []string{"full"}} {
		t.Run(stringifyBulkCodexModeForTest(mode), func(t *testing.T) {
			repo := &accountRepoStubForBulkUpdate{}
			result, err := (&adminServiceImpl{accountRepo: repo}).BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
				AccountIDs: []int64{1}, Extra: map[string]any{codexFingerprintModeExtraKey: mode},
			})
			require.Nil(t, result)
			requireApplicationErrorReason(t, err, "OPENAI_CODEX_FINGERPRINT_MODE_INVALID")
			require.Zero(t, repo.bulkUpdateCalls)
		})
	}
}

func stringifyBulkCodexModeForTest(value any) string {
	if value == nil {
		return "null"
	}
	if text, ok := value.(string); ok {
		return "string:" + text
	}
	return "non-string"
}

func TestBulkCodexFingerprintModeRejectsInvalidTargetsBeforeWrites(t *testing.T) {
	for _, other := range []*Account{
		nil,
		{ID: 2, Platform: PlatformAnthropic, Type: AccountTypeOAuth},
		{ID: 2, Platform: PlatformOpenAI, Type: AccountTypeAPIKey},
		{ID: 2, Platform: PlatformOpenAI, Type: AccountTypeSetupToken},
		{ID: 2, Platform: PlatformGemini, Type: AccountTypeOAuth},
	} {
		for _, mode := range []string{"off", "full"} {
			repo := &accountRepoStubForBulkUpdate{getByIDsAccounts: []*Account{{ID: 1, Platform: PlatformOpenAI, Type: AccountTypeOAuth}, other}}
			result, err := (&adminServiceImpl{accountRepo: repo}).BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
				AccountIDs: []int64{1, 2}, Extra: map[string]any{codexFingerprintModeExtraKey: mode},
			})
			require.Nil(t, result)
			requireApplicationErrorReason(t, err, "OPENAI_BULK_TARGET_INVALID")
			require.Zero(t, repo.bulkUpdateCalls, "mixed or missing targets must reject the complete mutation before any write")
		}
	}
}

func TestBulkCodexFingerprintModeMatchesRuntimePATEligibility(t *testing.T) {
	account := &Account{
		ID: 1, Platform: PlatformOpenAI, Type: AccountTypeOAuth,
		Credentials: map[string]any{"auth_mode": OpenAIAuthModePersonalAccessToken},
	}
	require.True(t, account.IsOpenAIOAuth())
	repo := &accountRepoStubForBulkUpdate{getByIDsAccounts: []*Account{account}}
	result, err := (&adminServiceImpl{accountRepo: repo}).BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
		AccountIDs: []int64{1}, Extra: map[string]any{codexFingerprintModeExtraKey: "full"},
	})
	require.NoError(t, err)
	require.Equal(t, 1, result.Success)
	require.True(t, repo.lastBulkUpdate.EnsureCodexFingerprintSeed)
}

func TestBulkCodexFingerprintUserAgentUpdatesOnlySelectedCredential(t *testing.T) {
	for _, ua := range []string{
		"",
		"codex-tui/0.153.3 (Mac OS 26.5.1; arm64) iTerm.app/3.6.11 (codex-tui; 0.153.3)",
		"  codex_cli_rs/0.146.0 (Windows 11; x86_64) WindowsTerminal  ",
	} {
		t.Run(ua, func(t *testing.T) {
			account := &Account{
				ID: 1, Platform: PlatformOpenAI, Type: AccountTypeOAuth,
				Credentials: map[string]any{"access_token": "keep-access", "refresh_token": "keep-refresh", "user_agent": "old-profile"},
				Extra:       map[string]any{codexFingerprintSeedExtraKey: testCodexFingerprintSeed},
			}
			repo := &accountRepoStubForBulkUpdate{getByIDsAccounts: []*Account{account}}
			result, err := (&adminServiceImpl{accountRepo: repo}).BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
				AccountIDs: []int64{1}, Extra: map[string]any{codexFingerprintModeExtraKey: "full"},
				Credentials: map[string]any{"user_agent": ua},
			})
			require.NoError(t, err)
			require.Equal(t, 1, result.Success)
			require.Equal(t, map[string]any{"user_agent": strings.TrimSpace(ua)}, repo.lastBulkUpdate.Credentials)
			require.Nil(t, repo.lastBulkUpdate.Concurrency, "selecting a fingerprint must not change concurrency")
			require.NotContains(t, repo.lastBulkUpdate.Extra, codexFingerprintSeedExtraKey)
			require.Equal(t, "keep-access", account.Credentials["access_token"])
			require.Equal(t, "keep-refresh", account.Credentials["refresh_token"])
			require.Equal(t, testCodexFingerprintSeed, account.Extra[codexFingerprintSeedExtraKey])
		})
	}
}

func TestBulkCodexFingerprintUserAgentRejectsInvalidBeforeWrites(t *testing.T) {
	for _, ua := range []any{nil, true, 3, "Mozilla/5.0", "codex-tui/invalid", "codex-tui/0.153.3\r\nInjected: value", "codex-tui/0.153.3\t", "codex-tui/0.153.3 中文", "codex-tui/0.153.3 " + strings.Repeat("x", 1024)} {
		t.Run(stringifyBulkCodexModeForTest(ua), func(t *testing.T) {
			repo := &accountRepoStubForBulkUpdate{}
			result, err := (&adminServiceImpl{accountRepo: repo}).BulkUpdateAccounts(context.Background(), &BulkUpdateAccountsInput{
				AccountIDs: []int64{1}, Extra: map[string]any{codexFingerprintModeExtraKey: "full"},
				Credentials: map[string]any{"user_agent": ua},
			})
			require.Nil(t, result)
			requireApplicationErrorReason(t, err, "OPENAI_CODEX_USER_AGENT_INVALID")
			require.Zero(t, repo.bulkUpdateCalls)
		})
	}
}
