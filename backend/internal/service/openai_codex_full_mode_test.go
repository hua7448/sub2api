package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// These assertions check captured outbound traffic on the HTTP and WS paths.
// Embedded turn metadata is optional, but must agree when the caller supplied it.
func requireCodexFingerprintBodyIDs(t *testing.T, body []byte, installationID, sessionID, threadID string) {
	t.Helper()
	require.Equal(t, installationID, gjson.GetBytes(body, "client_metadata.x-codex-installation-id").String())
	require.Equal(t, sessionID, gjson.GetBytes(body, "client_metadata.session_id").String())
	require.Equal(t, threadID, gjson.GetBytes(body, "client_metadata.thread_id").String())
	require.Equal(t, threadID+":0", gjson.GetBytes(body, "client_metadata.x-codex-window-id").String())
	turnID := gjson.GetBytes(body, "client_metadata.turn_id").String()
	require.NotEmpty(t, turnID)
	if embedded := gjson.GetBytes(body, "client_metadata.x-codex-turn-metadata"); embedded.Exists() {
		metadata := embedded.String()
		require.Equal(t, installationID, gjson.Get(metadata, "installation_id").String())
		require.Equal(t, sessionID, gjson.Get(metadata, "session_id").String())
		require.Equal(t, threadID, gjson.Get(metadata, "thread_id").String())
		require.Equal(t, threadID+":0", gjson.Get(metadata, "window_id").String())
		require.Equal(t, turnID, gjson.Get(metadata, "turn_id").String())
		require.NotZero(t, gjson.Get(metadata, "turn_started_at_unix_ms").Int())
	}
}

func requireCodexFingerprintHeaderBodyParity(t *testing.T, headers http.Header, body []byte, installationID, sessionID, threadID string) {
	t.Helper()
	requireCodexFingerprintBodyIDs(t, body, installationID, sessionID, threadID)
	require.Equal(t, installationID, headers.Get("x-codex-installation-id"))
	require.Equal(t, sessionID, headers.Get("session-id"))
	require.Equal(t, sessionID, headers.Get("session_id"))
	require.Equal(t, threadID, headers.Get("thread-id"))
	require.Equal(t, threadID, headers.Get("x-client-request-id"))
	require.Equal(t, threadID+":0", headers.Get("x-codex-window-id"))
	metadata := headers.Get("x-codex-turn-metadata")
	require.Equal(t, installationID, gjson.Get(metadata, "installation_id").String())
	require.Equal(t, sessionID, gjson.Get(metadata, "session_id").String())
	require.Equal(t, threadID, gjson.Get(metadata, "thread_id").String())
	require.Equal(t, threadID+":0", gjson.Get(metadata, "window_id").String())
	require.Equal(t, gjson.GetBytes(body, "client_metadata.turn_id").String(), gjson.Get(metadata, "turn_id").String())
	if embedded := gjson.GetBytes(body, "client_metadata.x-codex-turn-metadata"); embedded.Exists() {
		require.Equal(t, gjson.Get(embedded.String(), "turn_started_at_unix_ms").Int(), gjson.Get(metadata, "turn_started_at_unix_ms").Int())
	}
}

func TestCodexFullModeCreatedAccountsHaveIndependentStableIdentities(t *testing.T) {
	repo := &upstreamBillingProbeAccountRepo{}
	svc := &adminServiceImpl{accountRepo: repo}
	createAccount := func(name string) *Account {
		account, err := svc.CreateAccount(context.Background(), &CreateAccountInput{
			Name:                 name,
			Platform:             PlatformOpenAI,
			Type:                 AccountTypeOAuth,
			SkipDefaultGroupBind: true,
			Extra: map[string]any{
				codexFingerprintModeExtraKey: "full",
				codexFingerprintSeedExtraKey: userSuppliedCodexFingerprintSeed,
			},
		})
		require.NoError(t, err)
		return account
	}

	firstAccount := createAccount("first-full-mode-account")
	secondAccount := createAccount("second-full-mode-account")
	firstSeed := requireValidCodexFingerprintSeed(t, firstAccount.Extra)
	secondSeed := requireValidCodexFingerprintSeed(t, secondAccount.Extra)
	require.NotEqual(t, firstSeed, secondSeed)
	require.NotEqual(t, userSuppliedCodexFingerprintSeed, firstSeed)
	require.NotEqual(t, userSuppliedCodexFingerprintSeed, secondSeed)
	require.NotEqual(t, firstAccount.ID, secondAccount.ID)

	clientHeaders := make(http.Header)
	clientHeaders.Set("session-id", "same-client-session")
	firstIDs := resolveCodexFingerprintIDsFromRequest(firstAccount, clientHeaders)
	secondIDs := resolveCodexFingerprintIDsFromRequest(secondAccount, clientHeaders)
	require.NotNil(t, firstIDs)
	require.NotNil(t, secondIDs)
	require.Equal(t, firstIDs.sessionID, firstIDs.threadID)
	require.Equal(t, secondIDs.sessionID, secondIDs.threadID)
	require.NotEqual(t, firstIDs.installationID, secondIDs.installationID)
	require.NotEqual(t, firstIDs.sessionID, secondIDs.sessionID)
	require.NotEqual(t, firstIDs.threadID, secondIDs.threadID)

	// A new caller session and a fresh account snapshot retain the same persisted
	// account identity; each turn remains distinct for protocol correctness.
	reloaded, err := repo.GetByID(context.Background(), firstAccount.ID)
	require.NoError(t, err)
	clientHeaders.Set("session-id", "another-client-session")
	nextIDs := resolveCodexFingerprintIDsFromRequest(reloaded, clientHeaders)
	require.NotNil(t, nextIDs)
	require.Equal(t, firstIDs.installationID, nextIDs.installationID)
	require.Equal(t, firstIDs.sessionID, nextIDs.sessionID)
	require.Equal(t, firstIDs.threadID, nextIDs.threadID)
	require.Equal(t, firstIDs.windowID, nextIDs.windowID)
	require.NotEmpty(t, firstIDs.turnID)
	require.NotEmpty(t, nextIDs.turnID)
	require.NotEqual(t, firstIDs.turnID, nextIDs.turnID)
	require.Equal(t, firstSeed, requireValidCodexFingerprintSeed(t, reloaded.Extra))
}
