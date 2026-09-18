package service

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestCodexMetadataPreservesBothCarriersAndPrecision(t *testing.T) {
	for _, mode := range []string{"session", "full"} {
		t.Run(mode, func(t *testing.T) {
			account := newTestOAuthAccount(8901, map[string]any{codexFingerprintModeExtraKey: mode})
			headers := http.Header{}
			headers.Set("x-codex-turn-metadata", `{"sandbox":{"mode":"workspace-write","paths":["/tmp/東京"]},"thread_source":"header-source","session_id":"client-session","prompt_cache_key":"client-session","conversation_id":"stale","precise":9007199254740993}`)
			original := headers.Clone()
			ids := resolveCodexFingerprintIDsFromRequest(account, headers)
			body := map[string]any{"prompt_cache_key": "independent-cache", "client_metadata": map[string]any{"x-codex-turn-metadata": `{"permissions":{"network":false},"thread_source":"body-source","session_id":"client-session","prompt_cache_key":"client-session","note":"keep client-session literal"}`}}
			require.True(t, applyCodexFingerprintClientMetadata(body, ids))
			applyCodexFingerprintHeaders(headers, ids)
			cm := body["client_metadata"].(map[string]any)
			bodyMeta := cm["x-codex-turn-metadata"].(string)
			headerMeta := headers.Get("x-codex-turn-metadata")
			for _, raw := range []string{bodyMeta, headerMeta} {
				require.Equal(t, "workspace-write", gjson.Get(raw, "sandbox.mode").String())
				require.Equal(t, "/tmp/東京", gjson.Get(raw, "sandbox.paths.0").String())
				require.Equal(t, "9007199254740993", gjson.Get(raw, "precise").Raw)
				require.Equal(t, "false", gjson.Get(raw, "permissions.network").Raw)
				require.Equal(t, ids.sessionID, gjson.Get(raw, "session_id").String())
				require.Equal(t, ids.sessionID, gjson.Get(raw, "conversation_id").String())
				require.Equal(t, ids.sessionID, gjson.Get(raw, "prompt_cache_key").String())
				require.Equal(t, "keep client-session literal", gjson.Get(raw, "note").String())
				require.True(t, isASCIIBytes([]byte(raw)))
			}
			require.Equal(t, "header-source", gjson.Get(headerMeta, "thread_source").String())
			require.Equal(t, "body-source", gjson.Get(bodyMeta, "thread_source").String())
			require.Equal(t, "independent-cache", body["prompt_cache_key"])
			require.Contains(t, original.Get("x-codex-turn-metadata"), "stale")
		})
	}
}

func TestCodexMetadataMergeModeAndTurnIsolation(t *testing.T) {
	for _, mode := range []string{"off", "device", "session", "full"} {
		account := newTestOAuthAccount(8902, map[string]any{codexFingerprintModeExtraKey: mode})
		h := http.Header{}
		h.Set("x-codex-turn-metadata", `{"sandbox":"caller-value"}`)
		ids := resolveCodexFingerprintIDsFromRequest(account, h)
		cm := map[string]any{}
		applyCodexFingerprintToClientMetadataMap(cm, ids)
		if mode == "off" || mode == "device" {
			require.NotContains(t, cm, "x-codex-turn-metadata")
			continue
		}
		require.Equal(t, "caller-value", gjson.Get(cm["x-codex-turn-metadata"].(string), "sandbox").String())
		first := map[string]any{"x-codex-turn-metadata": `{"body_only_first_turn":true}`}
		applyCodexFingerprintToClientMetadataMap(first, ids)
		next := *ids
		nextBody := map[string]any{}
		applyCodexFingerprintToClientMetadataMap(nextBody, &next)
		out := http.Header{}
		applyCodexFingerprintHeaders(out, &next)
		require.False(t, gjson.Get(out.Get("x-codex-turn-metadata"), "body_only_first_turn").Exists(), "a new turn must not inherit the previous body's metadata")
	}
}

func TestCodexMetadataCacheAliasesAndMalformedFallback(t *testing.T) {
	fields := map[string]any{"session_id": "account-session"}
	for _, raw := range []string{`{"session_id":"old","prompt_cache_key":"independent"}`, `{"prompt_cache_key":"explicit-only"}`} {
		h := http.Header{}
		h.Set("x-codex-turn-metadata", raw)
		rewriteCodexTurnMetadataFields(h, fields)
		require.Equal(t, gjson.Get(raw, "prompt_cache_key").String(), gjson.Get(h.Get("x-codex-turn-metadata"), "prompt_cache_key").String())
	}
	for _, raw := range []string{"", "{", "null", "[]", "{} {}"} {
		merged := mergeCodexFingerprintTurnMetadata(raw, `{"sandbox":"existing-only"}`, nil)
		require.True(t, json.Valid([]byte(merged)))
		require.Equal(t, "existing-only", gjson.Get(merged, "sandbox").String())
	}
	account := newTestOAuthAccount(8903, map[string]any{codexFingerprintModeExtraKey: "session"})
	ids := resolveCodexFingerprintIDsFromRequest(account, nil)
	h := http.Header{}
	applyCodexFingerprintHeaders(h, ids)
	require.False(t, gjson.Get(h.Get("x-codex-turn-metadata"), "sandbox").Exists(), "do not invent absent environment metadata")
	require.False(t, gjson.Get(h.Get("x-codex-turn-metadata"), "prompt_cache_key").Exists())
}

func TestCodexMetadataFallbackCacheBelongsToOriginalCarrier(t *testing.T) {
	fields := map[string]any{"session_id": "account-session"}
	for _, pair := range [][2]string{
		{`{"session_id":"body-session"}`, `{"session_id":"header-session","prompt_cache_key":"header-session"}`},
		{`{"session_id":"header-session"}`, `{"session_id":"body-session","prompt_cache_key":"body-session"}`},
	} {
		merged := mergeCodexFingerprintTurnMetadata(pair[0], pair[1], fields)
		require.Equal(t, "account-session", gjson.Get(merged, "prompt_cache_key").String())
	}
}
