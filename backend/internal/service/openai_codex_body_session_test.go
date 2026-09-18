package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestCodexFingerprintBodySessionPrecedence(t *testing.T) {
	account := newTestOAuthAccount(8801, map[string]any{codexFingerprintModeExtraKey: "session"})
	seed, ok := codexFingerprintSeed(account.Extra)
	require.True(t, ok)
	for _, tc := range []struct{ name, header, turnHeader, body, want string }{
		{"explicit header", "header-session", "", `{"prompt_cache_key":"cache-session","client_metadata":{"session_id":"body-session"}}`, "header-session"},
		{"header metadata", "", `{"session_id":"turn-session"}`, `{"prompt_cache_key":"cache-session"}`, "turn-session"},
		{"body session", "", "", `{"prompt_cache_key":"explicit-cache","client_metadata":{"session_id":"body-session"}}`, "body-session"},
		{"body alias", "", "", `{"client_metadata":{"session-id":"alias-session"}}`, "alias-session"},
		{"embedded metadata", "", "", `{"client_metadata":{"x-codex-turn-metadata":"{\"session_id\":\"embedded-session\"}"}}`, "embedded-session"},
		{"cache fallback", "", "", `{"prompt_cache_key":"cache-session"}`, "cache-session"},
		{"invalid identity types", "", "null", `{"prompt_cache_key":42,"client_metadata":{"session_id":12}}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := http.Header{}
			h.Set("session-id", tc.header)
			h.Set("x-codex-turn-metadata", tc.turnHeader)
			before := h.Clone()
			ids := resolveCodexFingerprintIDsFromRequest(account, h, []byte(tc.body))
			require.NotNil(t, ids)
			want := resolveConvergedThreadID(seed, tc.want)
			if want == "" {
				want = ids.sessionID
			}
			require.Equal(t, want, ids.threadID)
			require.Equal(t, before, h, "identity lookup must not mutate client headers")
		})
	}
	for _, mode := range []string{"off", "device", "full"} {
		account := newTestOAuthAccount(8802, map[string]any{codexFingerprintModeExtraKey: mode})
		a := resolveCodexFingerprintIDsFromRequest(account, nil, []byte(`{"prompt_cache_key":"a"}`))
		b := resolveCodexFingerprintIDsFromRequest(account, nil, []byte(`{"prompt_cache_key":"b"}`))
		if mode == "off" {
			require.Nil(t, a)
			require.Nil(t, b)
			continue
		}
		require.Equal(t, a.threadID, b.threadID, "body fallback must preserve mode boundaries")
	}
}

func TestCodexBodyOnlySessionsAcrossForwardTransports(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, transport := range []string{"http", "http_passthrough", "websocket"} {
		t.Run(transport, func(t *testing.T) {
			cfg := &config.Config{}
			useWS := transport == "websocket"
			cfg.Gateway.OpenAIWS.Enabled = useWS
			cfg.Gateway.OpenAIWS.OAuthEnabled = useWS
			cfg.Gateway.OpenAIWS.ResponsesWebsocketsV2 = useWS
			cfg.Gateway.OpenAIWS.MaxConnsPerAccount = 3
			cfg.Gateway.OpenAIWS.MaxIdlePerAccount = 1
			account := newTestOAuthAccount(8803, map[string]any{codexFingerprintModeExtraKey: "session", "openai_oauth_passthrough": transport == "http_passthrough", "responses_websockets_v2_enabled": useWS})
			account.Concurrency = 3
			account.Credentials = map[string]any{"access_token": "test-token", "chatgpt_account_id": "test-account"}
			var firstThread string
			for _, session := range []string{"conversation-a", "conversation-b", "conversation-a"} {
				upstream := &httpUpstreamRecorder{responses: []*http.Response{openAICompatSSECompletedResponse("resp_session", "gpt-5.1")}}
				conn := &openAIWSCaptureConn{events: [][]byte{[]byte(`{"type":"response.completed","response":{"id":"resp_session_ws","model":"gpt-5.1","usage":{"input_tokens":1,"output_tokens":1}}}`)}}
				dialer := &openAIWSCaptureDialer{conn: conn}
				pool := newOpenAIWSConnPool(cfg)
				pool.setClientDialerForTest(dialer)
				svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstream, cache: &stubGatewayCache{}, openaiWSResolver: NewOpenAIWSProtocolResolver(cfg), toolCorrector: NewCodexToolCorrector(), openaiWSPool: pool}
				body, err := json.Marshal(map[string]any{"model": "gpt-5.1", "instructions": "Help the user.", "stream": false, "prompt_cache_key": session, "client_metadata": map[string]any{"x-codex-turn-metadata": `{"permissions":{"network":false}}`}, "input": []map[string]string{{"type": "input_text", "text": "hello"}}})
				require.NoError(t, err)
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))
				c.Request.Header.Set("User-Agent", CodexCanonicalUserAgent())
				c.Request.Header.Set("x-codex-turn-metadata", `{"sandbox":"client-sandbox"}`)
				_, err = svc.Forward(context.Background(), c, account, body)
				pool.Close()
				require.NoError(t, err)
				var headers http.Header
				var payload []byte
				if useWS {
					headers = dialer.lastHeaders
					payload, err = json.Marshal(conn.lastWrite)
					require.NoError(t, err)
				} else {
					require.NotNil(t, upstream.lastReq)
					headers = upstream.lastReq.Header
					payload = upstream.lastBody
				}
				require.Equal(t, "client-sandbox", gjson.Get(gjson.GetBytes(payload, "client_metadata.x-codex-turn-metadata").String(), "sandbox").String())
				require.Equal(t, "false", gjson.Get(headers.Get("x-codex-turn-metadata"), "permissions.network").Raw)
				thread := headers.Get("thread-id")
				require.NotEmpty(t, thread)
				require.Equal(t, thread, gjson.GetBytes(payload, "client_metadata.thread_id").String())
				require.Equal(t, session, gjson.GetBytes(payload, "prompt_cache_key").String(), "caller-owned cache keys must remain unchanged")
				if firstThread == "" {
					firstThread = thread
				} else if session == "conversation-a" {
					require.Equal(t, firstThread, thread)
				} else {
					require.NotEqual(t, firstThread, thread, "different body-only conversations must not share one thread")
				}
			}
		})
	}
}
