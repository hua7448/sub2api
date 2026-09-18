package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestOpenAIGatewayService_Forward_AccountProfileAndFullAcrossHTTPAndWS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const accountUA = "codex-tui/0.153.3 (Mac OS 26.5.1; arm64) iTerm.app/3.6.11 (codex-tui; 0.153.3)"
	const canonicalVersion = "0.146.0"
	SetCodexCanonicalUserAgentResolver(func() string { return buildCodexCLIUserAgent(canonicalVersion) })
	t.Cleanup(func() { SetCodexCanonicalUserAgentResolver(nil) })

	for _, forceCLI := range []bool{false, true} {
		t.Run(fmt.Sprintf("force_cli=%t", forceCLI), func(t *testing.T) {
			var fullThreadID string
			for _, useWS := range []bool{false, true} {
				t.Run(fmt.Sprintf("websocket=%t", useWS), func(t *testing.T) {
					cfg := &config.Config{}
					cfg.Security.URLAllowlist.Enabled = false
					cfg.Gateway.ForceCodexCLI = forceCLI
					cfg.Gateway.OpenAIWS.Enabled = useWS
					cfg.Gateway.OpenAIWS.OAuthEnabled = useWS
					cfg.Gateway.OpenAIWS.ResponsesWebsocketsV2 = useWS
					cfg.Gateway.OpenAIWS.MaxConnsPerAccount = 10
					cfg.Gateway.OpenAIWS.MaxIdlePerAccount = 1
					cfg.Gateway.OpenAIWS.DynamicMaxConnsByAccountConcurrencyEnabled = true
					account := newTestOAuthAccount(1400, map[string]any{
						codexFingerprintModeExtraKey:      "full",
						"responses_websockets_v2_enabled": useWS,
					})
					account.Concurrency = 3
					account.Status = StatusActive
					account.Schedulable = true
					account.Credentials = map[string]any{
						"access_token": "profile-token", "chatgpt_account_id": "profile-account",
						"user_agent": accountUA,
					}
					var headers http.Header
					var payload map[string]any
					httpCalls := 0
					upstream := &codexModelsHTTPUpstreamStub{do: func(request *http.Request, _ string, id int64, concurrency int) (*http.Response, error) {
						httpCalls++
						require.Equal(t, account.ID, id)
						require.Equal(t, 3, concurrency)
						headers = request.Header.Clone()
						raw, err := io.ReadAll(request.Body)
						require.NoError(t, err)
						require.NoError(t, json.Unmarshal(raw, &payload))
						return openAICompatSSECompletedResponse("resp_profile_http", "gpt-5.1"), nil
					}}
					captureConn := &openAIWSCaptureConn{events: [][]byte{
						[]byte(`{"type":"response.completed","response":{"id":"resp_profile_ws","model":"gpt-5.1","usage":{"input_tokens":1,"output_tokens":1}}}`),
					}}
					captureDialer := &openAIWSCaptureDialer{conn: captureConn}
					pool := newOpenAIWSConnPool(cfg)
					defer pool.Close()
					pool.setClientDialerForTest(captureDialer)
					svc := &OpenAIGatewayService{
						cfg: cfg, httpUpstream: upstream, cache: &stubGatewayCache{},
						openaiWSResolver: NewOpenAIWSProtocolResolver(cfg), toolCorrector: NewCodexToolCorrector(), openaiWSPool: pool,
					}
					originalSession := fmt.Sprintf("real-session-ws-%t", useWS)
					body := []byte(fmt.Sprintf(`{"model":"gpt-5.1","stream":false,"input":[{"type":"input_text","text":"hello"}],"client_metadata":{"session_id":%q,"thread_id":"original-thread"}}`, originalSession))
					c, _ := gin.CreateTestContext(httptest.NewRecorder())
					c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewReader(body))
					c.Request.Header.Set("User-Agent", "downstream-client/1.0")
					c.Request.Header.Set("session-id", originalSession)
					c.Request.Header.Set("thread-id", "original-thread")

					result, err := svc.Forward(context.Background(), c, account, body)
					require.NoError(t, err)
					require.NotNil(t, result)
					if useWS {
						require.Zero(t, httpCalls, "WS must not silently fall back to HTTP")
						require.Equal(t, 1, captureDialer.DialCount())
						headers = captureDialer.lastHeaders
						payload = captureConn.lastWrite
						require.Equal(t, 3, pool.effectiveMaxConnsByAccount(account))
					} else {
						require.Equal(t, 1, httpCalls)
						require.Zero(t, captureDialer.DialCount())
					}
					wantUA := "codex-tui/" + canonicalVersion + " (Mac OS 26.5.1; arm64) iTerm.app/3.6.11 (codex-tui; " + canonicalVersion + ")"
					if forceCLI {
						wantUA = CodexCanonicalUserAgent()
					}
					require.Equal(t, wantUA, headers.Get("User-Agent"))
					require.Equal(t, "codex-tui", headers.Get("originator"))
					require.Equal(t, canonicalVersion, headers.Get("version"))
					require.Equal(t, "Bearer profile-token", headers.Get("Authorization"))
					require.Equal(t, "profile-account", headers.Get("chatgpt-account-id"))
					threadID := headers.Get("thread-id")
					require.NotEmpty(t, threadID)
					require.NotEqual(t, "original-thread", threadID)
					require.Equal(t, threadID, headers.Get("session-id"))
					if fullThreadID == "" {
						fullThreadID = threadID
					} else {
						require.Equal(t, fullThreadID, threadID, "full must keep the same account thread across original sessions and transports")
					}
					metadata, ok := payload["client_metadata"].(map[string]any)
					require.True(t, ok)
					require.Equal(t, threadID, metadata["thread_id"])
					require.Equal(t, threadID, metadata["session_id"])
					require.Equal(t, headers.Get("x-codex-installation-id"), metadata["x-codex-installation-id"])
					require.Equal(t, 3, account.Concurrency)
					require.Equal(t, accountUA, account.GetOpenAIUserAgent())
				})
			}
		})
	}
}
