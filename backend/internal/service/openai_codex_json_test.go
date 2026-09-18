package service

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalCodexASCIIJSONMatchesOfficialUnicodeEscaping(t *testing.T) {
	value := map[string]any{
		"workspaces": map[string]any{
			"/tmp/東京": map[string]any{
				"label": "Agentlarım",
				"emoji": "🚀",
			},
		},
		"literal": "<>&",
	}

	encoded, err := marshalCodexASCIIJSON(value)
	require.NoError(t, err)
	require.Equal(t,
		`{"literal":"<>&","workspaces":{"/tmp/\u6771\u4eac":{"emoji":"\ud83d\ude80","label":"Agentlar\u0131m"}}}`,
		string(encoded),
	)
	for _, b := range encoded {
		require.Less(t, b, byte(utf8RuneSelf), "Codex metadata must be ASCII-only")
	}

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.Equal(t, value, decoded)
}

func TestCodexFingerprintMetadataRewritesEscapeUnicodeAndAlignAliases(t *testing.T) {
	account := newTestOAuthAccount(7701, map[string]any{codexFingerprintModeExtraKey: "session"})
	ids := resolveCodexFingerprintIDs(account, "client-session-unicode", codexFingerprintSession)
	require.NotNil(t, ids)

	turnMetadata := `{"installation_id":"old-install","session_id":"old-session","session-id":"old-session-alias","x-codex-session-id":"old-x-session-alias","thread_id":"old-thread","turn_id":"old-turn","window_id":"old-thread:0","x-codex-installation-id":"old-install-alias","thread-id":"old-thread-alias","x-client-request-id":"old-request-alias","x-codex-window-id":"old-window-alias","turn-id":"old-turn-alias","sandbox":"东京🚀"}`
	h := make(http.Header)
	h.Set("X-Codex-Turn-Metadata", turnMetadata)
	rewriteCodexTurnMetadataFields(h, map[string]any{
		"installation_id": ids.installationID,
		"session_id":      ids.sessionID,
		"thread_id":       ids.threadID,
		"turn_id":         ids.turnID,
		"window_id":       ids.windowID,
	})
	headerRaw := h.Get("x-codex-turn-metadata")
	require.NotEmpty(t, headerRaw)
	require.True(t, isASCIIBytes([]byte(headerRaw)))
	require.Contains(t, headerRaw, `\u4e1c\u4eac\ud83d\ude80`)

	var headerMeta map[string]any
	require.NoError(t, json.Unmarshal([]byte(headerRaw), &headerMeta))
	require.Equal(t, ids.installationID, headerMeta["installation_id"])
	require.Equal(t, ids.installationID, headerMeta["x-codex-installation-id"])
	require.Equal(t, ids.sessionID, headerMeta["session_id"])
	require.Equal(t, ids.sessionID, headerMeta["session-id"])
	require.Equal(t, ids.sessionID, headerMeta["x-codex-session-id"])
	require.Equal(t, ids.threadID, headerMeta["thread_id"])
	require.Equal(t, ids.threadID, headerMeta["thread-id"])
	require.Equal(t, ids.threadID, headerMeta["x-client-request-id"])
	require.Equal(t, ids.windowID, headerMeta["window_id"])
	require.Equal(t, ids.windowID, headerMeta["x-codex-window-id"])
	require.Equal(t, ids.turnID, headerMeta["turn_id"])
	require.Equal(t, ids.turnID, headerMeta["turn-id"])

	body := map[string]any{
		"client_metadata": map[string]any{
			"installation_id":         "old-flat-install",
			"x-codex-installation-id": "old-flat-install-alias",
			"session_id":              "old-flat-session",
			"session-id":              "old-flat-session-alias",
			"x-codex-session-id":      "old-flat-x-session-alias",
			"thread_id":               "old-flat-thread",
			"thread-id":               "old-flat-thread-alias",
			"x-client-request-id":     "old-flat-request-alias",
			"window_id":               "old-flat-window",
			"x-codex-window-id":       "old-flat-window-alias",
			"turn-id":                 "old-flat-turn-alias",
			"x-codex-turn-metadata":   turnMetadata,
		},
	}
	rawInput, err := json.Marshal(body)
	require.NoError(t, err)
	require.True(t, applyCodexFingerprintClientMetadata(body, ids))
	cm, ok := body["client_metadata"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, ids.installationID, cm["installation_id"])
	require.Equal(t, ids.installationID, cm["x-codex-installation-id"])
	require.Equal(t, ids.sessionID, cm["session_id"])
	require.Equal(t, ids.sessionID, cm["session-id"])
	require.Equal(t, ids.sessionID, cm["x-codex-session-id"])
	require.Equal(t, ids.threadID, cm["thread_id"])
	require.Equal(t, ids.threadID, cm["thread-id"])
	require.Equal(t, ids.threadID, cm["x-client-request-id"])
	require.Equal(t, ids.windowID, cm["window_id"])
	require.Equal(t, ids.windowID, cm["x-codex-window-id"])
	require.Equal(t, ids.turnID, cm["turn_id"])
	require.Equal(t, ids.turnID, cm["turn-id"])
	embedded, ok := cm["x-codex-turn-metadata"].(string)
	require.True(t, ok)
	require.True(t, isASCIIBytes([]byte(embedded)))
	require.Contains(t, embedded, `\u4e1c\u4eac\ud83d\ude80`)
	var embeddedMeta map[string]any
	require.NoError(t, json.Unmarshal([]byte(embedded), &embeddedMeta))
	require.Equal(t, ids.turnID, embeddedMeta["turn_id"])
	require.Equal(t, ids.turnID, embeddedMeta["turn-id"])

	rawBody, changed, err := applyCodexFingerprintClientMetadataRaw(rawInput, cloneCodexFingerprintIDsForTest(ids))
	require.NoError(t, err)
	require.True(t, changed)
	require.True(t, isASCIIBytes(rawBody))
	var rawDecoded map[string]any
	require.NoError(t, json.Unmarshal(rawBody, &rawDecoded))
	rawCM, ok := rawDecoded["client_metadata"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, cm, rawCM)
}

func TestCodexFingerprintPromptCacheKeyCapturesSessionAliases(t *testing.T) {
	account := newTestOAuthAccount(7702, map[string]any{codexFingerprintModeExtraKey: "session"})
	baseIDs := resolveCodexFingerprintIDs(account, "client-session-alias", codexFingerprintSession)
	require.NotNil(t, baseIDs)

	tests := []struct {
		name       string
		body       []byte
		wantKey    string
		wantRawKey string
	}{
		{
			name:       "canonical wins over aliases",
			body:       []byte(`{"prompt_cache_key":"canonical-session","client_metadata":{"session_id":"canonical-session","session-id":"alias-session","x-codex-session-id":"x-alias-session"}}`),
			wantKey:    baseIDs.sessionID,
			wantRawKey: baseIDs.sessionID,
		},
		{
			name:       "hyphen alias fallback",
			body:       []byte(`{"prompt_cache_key":"alias-session","client_metadata":{"session-id":"alias-session"}}`),
			wantKey:    baseIDs.sessionID,
			wantRawKey: baseIDs.sessionID,
		},
		{
			name:       "x-codex alias fallback",
			body:       []byte(`{"prompt_cache_key":"x-alias-session","client_metadata":{"x-codex-session-id":"x-alias-session"}}`),
			wantKey:    baseIDs.sessionID,
			wantRawKey: baseIDs.sessionID,
		},
		{
			name:       "explicit cache override is preserved",
			body:       []byte(`{"prompt_cache_key":"explicit-cache","client_metadata":{"session-id":"alias-session"}}`),
			wantKey:    "explicit-cache",
			wantRawKey: "explicit-cache",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mapBody map[string]any
			require.NoError(t, json.Unmarshal(tt.body, &mapBody))
			mapIDs := cloneCodexFingerprintIDsForTest(baseIDs)
			require.True(t, applyCodexFingerprintClientMetadata(mapBody, mapIDs))
			require.Equal(t, tt.wantKey, mapBody["prompt_cache_key"])

			rawIDs := cloneCodexFingerprintIDsForTest(baseIDs)
			rawBody, changed, err := applyCodexFingerprintClientMetadataRaw(tt.body, rawIDs)
			require.NoError(t, err)
			require.True(t, changed)
			var rawDecoded map[string]any
			require.NoError(t, json.Unmarshal(rawBody, &rawDecoded))
			require.Equal(t, tt.wantRawKey, rawDecoded["prompt_cache_key"])
		})
	}
}

func isASCIIBytes(value []byte) bool {
	for _, b := range value {
		if b >= utf8RuneSelf {
			return false
		}
	}
	return true
}

const utf8RuneSelf = 0x80
