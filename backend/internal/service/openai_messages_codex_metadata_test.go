package service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestForwardAsAnthropic_CodexMetadataUsesOneIdentityAcrossBodyAndHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, mode := range []string{"session", "full"} {
		t.Run(mode, func(t *testing.T) {
			upstream := &httpUpstreamRecorder{}
			svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
			accountA := newTestOAuthAccount(4411, map[string]any{codexFingerprintModeExtraKey: mode})
			accountB := newTestOAuthAccount(4412, map[string]any{
				codexFingerprintModeExtraKey: mode,
				codexFingerprintSeedExtraKey: "44444444-4444-4444-8444-444444444444",
			})
			for _, account := range []*Account{accountA, accountB} {
				account.Concurrency = 3
				account.Credentials = map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-account"}
			}
			const cacheKey = "stable-compat-cache"
			var firstTurnID, firstSessionID string
			for index, account := range []*Account{accountA, accountB, accountA} {
				upstream.resp = openAICompatSSECompletedResponse(fmt.Sprintf("resp_messages_metadata_%d", index), "gpt-5.4")
				if index == 0 {
					upstream.resp.Header.Set("x-codex-turn-state", "account-a-turn-state")
				}
				body := []byte(`{"model":"claude-sonnet-4-5","max_tokens":16,"messages":[{"role":"user","content":"hello"}],"stream":false}`)
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
				c.Request.Header.Set("Content-Type", "application/json")
				// Real Messages callers need not send any Codex identity headers.
				// Full also converges a later caller that sends a different session.
				if mode == "full" && index == 2 {
					c.Request.Header.Set("session-id", "another-client-session")
					c.Request.Header.Set("conversation_id", "another-client-conversation")
				}

				result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, cacheKey, "gpt-5.4")
				require.NoError(t, err)
				require.NotNil(t, result)
				seed := requireValidCodexFingerprintSeed(t, account.Extra)
				wantInstall := resolveConvergedInstallationID(account, seed)
				wantSession := resolveConvergedSessionID(seed)
				wantThread := resolveConvergedThreadID(seed, cacheKey)
				if mode == "full" {
					wantThread = wantSession
				}
				requireCodexFingerprintHeaderBodyParity(t, upstream.lastReq.Header, upstream.lastBody, wantInstall, wantSession, wantThread)
				require.True(t, gjson.GetBytes(upstream.lastBody, "client_metadata.x-codex-turn-metadata").Exists(), "Messages must populate missing embedded turn metadata")
				require.False(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").Exists(), "OAuth Messages cache identity stays header-only")
				require.False(t, gjson.GetBytes(upstream.lastBody, "previous_response_id").Exists())
				turnID := gjson.GetBytes(upstream.lastBody, "client_metadata.turn_id").String()
				switch index {
				case 0:
					firstTurnID, firstSessionID = turnID, wantSession
					require.Empty(t, upstream.lastReq.Header.Get("x-codex-turn-state"))
				case 1:
					require.NotEqual(t, firstSessionID, wantSession)
					require.Empty(t, upstream.lastReq.Header.Get("x-codex-turn-state"), "another account must not receive account A's cached state")
				case 2:
					require.Equal(t, firstSessionID, wantSession)
					require.NotEqual(t, firstTurnID, turnID)
					require.Equal(t, "account-a-turn-state", upstream.lastReq.Header.Get("x-codex-turn-state"), "internal cache key and account isolation must survive metadata convergence")
				}
				if mode == "full" && index == 2 {
					require.Equal(t, wantSession, upstream.lastReq.Header.Get("conversation_id"))
				} else {
					require.Empty(t, upstream.lastReq.Header.Get("conversation_id"), "do not invent the optional Messages conversation header")
				}
				require.Equal(t, 3, account.Concurrency)
			}
		})
	}
}

func TestForwardAsAnthropic_CodexMetadataOffAndDevicePreserveLegacySession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, mode := range []string{"off", "device"} {
		t.Run(mode, func(t *testing.T) {
			upstream := &httpUpstreamRecorder{resp: openAICompatSSECompletedResponse("resp_messages_legacy", "gpt-5.4")}
			svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
			account := newTestOAuthAccount(4413, map[string]any{
				codexFingerprintModeExtraKey: mode,
				codexFingerprintSeedExtraKey: testCodexFingerprintSeed,
			})
			account.Credentials = map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-account"}
			body := []byte(`{"model":"claude-sonnet-4-5","max_tokens":16,"messages":[{"role":"user","content":"hello"}],"stream":false}`)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
			// A previous attempt's staged identity must not leak after switching mode.
			stageCodexFingerprintIDs(c, resolveCodexFingerprintIDs(account, "stale-session", codexFingerprintFull))
			_, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "stable-compat-cache", "gpt-5.4")
			require.NoError(t, err)
			require.Equal(t, generateSessionUUID(isolateOpenAISessionID(0, "stable-compat-cache")), upstream.lastReq.Header.Get("session_id"))
			require.Empty(t, upstream.lastReq.Header.Get("thread-id"))
			require.False(t, gjson.GetBytes(upstream.lastBody, "client_metadata.session_id").Exists())
			require.False(t, gjson.GetBytes(upstream.lastBody, "client_metadata.turn_id").Exists())
			require.False(t, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").Exists())
			if mode == "off" {
				require.Empty(t, upstream.lastReq.Header.Get("x-codex-installation-id"))
				require.False(t, gjson.GetBytes(upstream.lastBody, "client_metadata").Exists())
			} else {
				installationID := resolveConvergedInstallationID(account, requireValidCodexFingerprintSeed(t, account.Extra))
				require.Equal(t, installationID, upstream.lastReq.Header.Get("x-codex-installation-id"))
				require.Equal(t, installationID, gjson.GetBytes(upstream.lastBody, "client_metadata.x-codex-installation-id").String())
			}
		})
	}
}
