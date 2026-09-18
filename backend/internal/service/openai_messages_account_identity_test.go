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
)

func TestForwardAsAnthropic_AccountCodexIdentityIsolation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, forceCLI := range []bool{false, true} {
		for _, stream := range []bool{false, true} {
			t.Run(fmt.Sprintf("force_cli=%t/stream=%t", forceCLI, stream), func(t *testing.T) {
				upstream := &httpUpstreamRecorder{}
				svc := &OpenAIGatewayService{
					httpUpstream: upstream,
					cfg: &config.Config{
						Gateway:  config.GatewayConfig{ForceCodexCLI: forceCLI},
						Security: config.SecurityConfig{URLAllowlist: config.URLAllowlistConfig{Enabled: false}},
					},
				}
				mac := &Account{
					ID: 11, Name: "mac-account", Platform: PlatformOpenAI, Type: AccountTypeOAuth, Concurrency: 3,
					Credentials: map[string]any{
						"access_token": "mac-token", "chatgpt_account_id": "mac-chatgpt-account",
						"user_agent": "codex_vscode/0.125.0 (Mac OS X 14.0; arm64) vscode",
					},
				}
				defaultAccount := &Account{
					ID: 12, Name: "default-account", Platform: PlatformOpenAI, Type: AccountTypeOAuth, Concurrency: 5,
					Credentials: map[string]any{
						"access_token": "default-token", "chatgpt_account_id": "default-chatgpt-account",
					},
				}
				canonical := resolveCodexOutboundIdentity("")
				// Reuse one service for A -> B -> A to catch cross-account identity or credential state.
				for index, account := range []*Account{mac, defaultAccount, mac} {
					upstream.resp = openAICompatSSECompletedResponse(fmt.Sprintf("resp_identity_%d", index), "gpt-5.4")
					body := []byte(fmt.Sprintf(`{"model":"claude-sonnet-4-5","max_tokens":16,"messages":[{"role":"user","content":"hello"}],"stream":%t}`, stream))
					c, _ := gin.CreateTestContext(httptest.NewRecorder())
					c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewReader(body))
					c.Request.Header.Set("Content-Type", "application/json")
					c.Request.Header.Set("User-Agent", "untrusted-client/1.0")
					c.Request.Header.Set("originator", "untrusted-client")
					c.Request.Header.Set("Authorization", "Bearer downstream-token")
					c.Request.Header.Set("chatgpt-account-id", "downstream-account")

					result, err := svc.ForwardAsAnthropic(context.Background(), c, account, body, "", "gpt-5.4")
					require.NoError(t, err)
					require.NotNil(t, result)
					require.Len(t, upstream.requests, index+1)
					request := upstream.requests[index]
					wantUA, wantOriginator := canonical.userAgent, canonical.originator
					if account == mac && !forceCLI {
						wantUA = "codex_vscode/" + canonical.version + " (Mac OS X 14.0; arm64) vscode"
						wantOriginator = "codex_vscode"
					}
					require.Equal(t, wantUA, request.Header.Get("User-Agent"))
					require.Equal(t, wantOriginator, request.Header.Get("originator"))
					require.Equal(t, canonical.version, request.Header.Get("version"))
					require.Equal(t, "responses=experimental", request.Header.Get("OpenAI-Beta"))
					require.Equal(t, "Bearer "+account.GetCredential("access_token"), request.Header.Get("Authorization"))
					require.Equal(t, account.GetCredential("chatgpt_account_id"), request.Header.Get("chatgpt-account-id"))
				}
				require.Equal(t, "codex_vscode/0.125.0 (Mac OS X 14.0; arm64) vscode", mac.GetCredential("user_agent"))
				require.Empty(t, defaultAccount.GetCredential("user_agent"))
				require.Equal(t, 3, mac.Concurrency)
				require.Equal(t, 5, defaultAccount.Concurrency)
			})
		}
	}
}
