package service

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodexTurnMetadataCompletionMissingAndInvalid(t *testing.T) {
	for _, mode := range []codexFingerprintMode{codexFingerprintSession, codexFingerprintFull} {
		for _, tc := range []struct {
			name  string
			value any
		}{
			{"absent", nil}, {"empty", ""}, {"whitespace", " \t "},
			{"null", "null"}, {"array", "[]"}, {"number", "42"},
			{"malformed", "{"}, {"non-string", 42},
		} {
			t.Run(string(mode)+"/"+tc.name, func(t *testing.T) {
				account := newTestOAuthAccount(8701, map[string]any{codexFingerprintModeExtraKey: string(mode)})
				ids := resolveCodexFingerprintIDs(account, "client-thread", mode)
				require.NotNil(t, ids)
				h := make(http.Header)
				cm := map[string]any{"custom": "preserved"}
				if tc.value != nil {
					cm["x-codex-turn-metadata"] = tc.value
					if raw, ok := tc.value.(string); ok {
						h.Set("x-codex-turn-metadata", raw)
					}
				}
				body := map[string]any{"model": "gpt-test", "input": []any{}, "client_metadata": cm}
				raw, err := json.Marshal(body)
				require.NoError(t, err)
				applyCodexFingerprintHeaders(h, ids)
				require.True(t, applyCodexFingerprintClientMetadata(body, ids))
				patched, changed, err := applyCodexFingerprintClientMetadataRaw(raw, cloneCodexFingerprintIDsForTest(ids))
				require.NoError(t, err)
				require.True(t, changed)
				var decoded map[string]any
				require.NoError(t, json.Unmarshal(patched, &decoded))
				require.Equal(t, body, decoded)
				require.Equal(t, "preserved", cm["custom"])
				var metadata map[string]any
				require.NoError(t, json.Unmarshal([]byte(h.Get("x-codex-turn-metadata")), &metadata))
				require.Equal(t, map[string]any{
					"installation_id": ids.installationID, "session_id": ids.sessionID,
					"thread_id": ids.threadID, "turn_id": ids.turnID, "window_id": ids.windowID,
					"turn_started_at_unix_ms": float64(ids.turnStartedAtUnixMs),
				}, metadata)
				require.Equal(t, h.Get("x-codex-turn-metadata"), cm["x-codex-turn-metadata"])
			})
		}
	}
}

func TestCodexTurnMetadataCompletionPreservesOpaqueFields(t *testing.T) {
	ids := resolveCodexFingerprintIDs(newTestOAuthAccount(8702, map[string]any{codexFingerprintModeExtraKey: "full"}), "client", codexFingerprintFull)
	original := `{"sandbox":{"mode":"workspace-write","paths":["/tmp/東京"]},"thread_source":"client-provided","permissions":{"network":false},"session-id":"stale","turn-id":"stale"}`
	h := make(http.Header)
	h.Set("x-codex-turn-metadata", original)
	cm := map[string]any{"x-codex-turn-metadata": original}
	applyCodexFingerprintHeaders(h, ids)
	require.True(t, applyCodexFingerprintToClientMetadataMap(cm, ids))
	require.Equal(t, h.Get("x-codex-turn-metadata"), cm["x-codex-turn-metadata"])
	require.True(t, isASCIIBytes([]byte(h.Get("x-codex-turn-metadata"))))
	var before, after map[string]any
	require.NoError(t, json.Unmarshal([]byte(original), &before))
	require.NoError(t, json.Unmarshal([]byte(h.Get("x-codex-turn-metadata")), &after))
	for _, key := range []string{"sandbox", "thread_source", "permissions"} {
		require.Equal(t, before[key], after[key])
	}
	require.Equal(t, ids.sessionID, after["session-id"])
	require.Equal(t, ids.turnID, after["turn-id"])
}

func TestCodexTurnMetadataCompletionKeepsModeAndConversationBoundaries(t *testing.T) {
	for _, mode := range []codexFingerprintMode{codexFingerprintOff, codexFingerprintDevice, codexFingerprintSession, codexFingerprintFull} {
		t.Run(string(mode), func(t *testing.T) {
			account := newTestOAuthAccount(8703, map[string]any{codexFingerprintModeExtraKey: string(mode)})
			ids := resolveCodexFingerprintIDsFromRequest(account, nil)
			h := make(http.Header)
			h.Set("conversation_id", "client-conversation")
			h.Set("conversation-id", "client-alias")
			applyCodexFingerprintHeaders(h, ids)
			cm := map[string]any{}
			applyCodexFingerprintToClientMetadataMap(cm, ids)
			if mode == codexFingerprintOff || mode == codexFingerprintDevice {
				require.Empty(t, h.Get("x-codex-turn-metadata"))
				require.NotContains(t, cm, "x-codex-turn-metadata")
				require.Equal(t, "client-conversation", h.Get("conversation_id"))
				require.Equal(t, "client-alias", h.Get("conversation-id"))
			} else {
				require.Equal(t, ids.sessionID, h.Get("conversation_id"))
				require.Equal(t, ids.sessionID, h.Get("conversation-id"))
			}
			absent := make(http.Header)
			applyCodexFingerprintHeaders(absent, ids)
			require.Empty(t, absent.Get("conversation_id"))
			require.Empty(t, absent.Get("conversation-id"))
			require.Empty(t, absent.Get("x-codex-turn-state"))
		})
	}
}
