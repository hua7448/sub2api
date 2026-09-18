package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type codexCompatFingerprintUpstream struct {
	httpUpstreamRecorder
	concurrency int
}

func (u *codexCompatFingerprintUpstream) Do(req *http.Request, proxyURL string, accountID int64, concurrency int) (*http.Response, error) {
	u.concurrency = concurrency
	return u.httpUpstreamRecorder.Do(req, proxyURL, accountID, concurrency)
}

func (u *codexCompatFingerprintUpstream) DoWithTLS(req *http.Request, proxyURL string, accountID int64, concurrency int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.Do(req, proxyURL, accountID, concurrency)
}

func newCodexCompatFingerprintUpstream() *codexCompatFingerprintUpstream {
	return &codexCompatFingerprintUpstream{httpUpstreamRecorder: httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"error":{"type":"invalid_request_error","message":"test stops after outbound capture"}}`)),
	}}}
}

func TestCodexCompatibilityFingerprintChatCompletions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, mode := range []string{"off", "device", "session", "full", "api_key"} {
		t.Run(mode, func(t *testing.T) {
			body := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"hello"}],"stream":false}`)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Request.Header.Set("session-id", "caller-thread")
			c.Request.Header.Set("x-codex-turn-metadata", `{"sandbox":"client-provided"}`)
			c.Set("api_key", &APIKey{ID: 91})
			account := newTestOAuthAccount(6201, map[string]any{codexFingerprintModeExtraKey: mode})
			account.Concurrency = 3
			account.Credentials = map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "account-test"}
			if mode == "api_key" {
				account.Type = AccountTypeAPIKey
				account.Credentials = map[string]any{"api_key": "sk-test"}
				account.Extra = map[string]any{"openai_responses_supported": true, codexFingerprintModeExtraKey: "full", codexFingerprintSeedExtraKey: testCodexFingerprintSeed}
			}
			// Re-entering this compatibility path must not reuse a prior attempt's IDs.
			previous := newTestOAuthAccount(account.ID, map[string]any{codexFingerprintModeExtraKey: "full"})
			stageCodexFingerprintIDs(c, resolveCodexFingerprintIDsFromRequest(previous, c.Request.Header))
			upstream := newCodexCompatFingerprintUpstream()
			svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
			_, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, "explicit-cache-key", "")
			require.Error(t, err)
			require.NotNil(t, upstream.lastReq)
			require.Equal(t, 3, upstream.concurrency)
			require.Equal(t, "explicit-cache-key", gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String())
			header := upstream.lastReq.Header
			cm := gjson.GetBytes(upstream.lastBody, "client_metadata")
			if mode == "off" || mode == "api_key" {
				require.Empty(t, header.Get("x-codex-installation-id"))
				require.False(t, cm.Exists())
				require.Equal(t, generateSessionUUID(isolateOpenAISessionID(91, "explicit-cache-key")), header.Get("session_id"))
				return
			}
			ids := stagedCodexFingerprintIDs(c, account)
			require.NotNil(t, ids)
			require.Equal(t, ids.installationID, header.Get("x-codex-installation-id"))
			require.Equal(t, ids.installationID, cm.Get("x-codex-installation-id").String())
			if mode == "device" {
				require.False(t, cm.Get("thread_id").Exists())
				require.Equal(t, generateSessionUUID(isolateOpenAISessionID(91, "explicit-cache-key")), header.Get("session_id"))
				return
			}
			require.Equal(t, ids.sessionID, header.Get("session_id"), "post-build compatibility session must not overwrite the converged identity")
			require.Equal(t, ids.sessionID, cm.Get("session_id").String())
			require.Equal(t, ids.threadID, header.Get("thread-id"))
			require.Equal(t, ids.threadID, cm.Get("thread_id").String())
			require.Equal(t, ids.turnID, cm.Get("turn_id").String())
			var metadata map[string]any
			require.NoError(t, json.Unmarshal([]byte(header.Get("x-codex-turn-metadata")), &metadata))
			require.Equal(t, ids.turnID, metadata["turn_id"])
			require.Equal(t, "client-provided", metadata["sandbox"])
			if mode == "full" {
				require.Equal(t, ids.sessionID, ids.threadID)
			}
		})
	}
}

func TestCodexCompatibilityFingerprintChatCompletionsWithoutSessionHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, mode := range []string{"session", "full"} {
		t.Run(mode, func(t *testing.T) {
			account := newTestOAuthAccount(6203, map[string]any{codexFingerprintModeExtraKey: mode})
			account.Concurrency = 3
			account.Credentials = map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "account-test"}
			threads := make([]string, 0, 3)
			turns := make(map[string]bool)
			for _, cacheKey := range []string{"cache-A", "cache-B", "cache-A"} {
				body := []byte(`{"model":"gpt-5.4","messages":[{"role":"user","content":"hello"}],"stream":false}`)
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
				c.Request.Header.Set("Content-Type", "application/json")
				upstream := newCodexCompatFingerprintUpstream()
				svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
				_, err := svc.ForwardAsChatCompletions(context.Background(), c, account, body, cacheKey, "")
				require.Error(t, err)
				require.NotNil(t, upstream.lastReq)
				require.Equal(t, 3, upstream.concurrency)
				require.Equal(t, cacheKey, gjson.GetBytes(upstream.lastBody, "prompt_cache_key").String(), "fallback identifies the thread without changing the cache key")
				header := upstream.lastReq.Header
				threadID := header.Get("thread-id")
				require.NotEmpty(t, threadID)
				require.Equal(t, threadID, gjson.GetBytes(upstream.lastBody, "client_metadata.thread_id").String())
				turnID := gjson.Get(header.Get("x-codex-turn-metadata"), "turn_id").String()
				require.NotEmpty(t, turnID)
				require.Equal(t, turnID, gjson.GetBytes(upstream.lastBody, "client_metadata.turn_id").String())
				threads = append(threads, threadID)
				turns[turnID] = true
				if mode == "full" {
					require.Equal(t, header.Get("session_id"), threadID)
				}
			}
			require.Equal(t, threads[0], threads[2], "returning to A must retain its thread")
			if mode == "session" {
				require.NotEqual(t, threads[0], threads[1], "A and B must have different session-mode threads")
			} else {
				require.Equal(t, threads[0], threads[1], "full retains the account's unified thread")
			}
			require.Len(t, turns, 3, "each request must have a fresh turn ID")
		})
	}
}
