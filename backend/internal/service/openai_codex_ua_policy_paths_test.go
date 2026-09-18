//go:build unit

package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const codexUAPolicyAccountUA = "codex-tui/0.153.3 (Mac OS 26.5.1; arm64) iTerm.app/3.6.11 (codex-tui; 0.153.3)"
const codexUAPolicyVersion = "0.146.0"
const codexUAPolicyCanonicalUA = "codex-tui/0.146.0 (Ubuntu 22.4.0; x86_64) xterm-256color"
const codexUAPolicyNormalizedAccountUA = "codex-tui/0.146.0 (Mac OS 26.5.1; arm64) iTerm.app/3.6.11 (codex-tui; 0.146.0)"

// Exercise independently constructed outbound requests: a successful inference
// identity must also be used for administrative probes and the models manifest.
func TestCodexUAPolicyModelsAndProbesMatchInference(t *testing.T) {
	gin.SetMode(gin.TestMode)
	originalEnforcement := codexIdentityEnforcement.Load()
	SetCodexCanonicalUserAgentResolver(func() string { return codexUAPolicyCanonicalUA })
	t.Cleanup(func() {
		SetCodexCanonicalUserAgentResolver(nil)
		SetCodexIdentityEnforcementEnabled(originalEnforcement)
	})
	for _, policy := range []struct {
		name    string
		cfg     *config.Config
		enforce bool
	}{
		{name: "no_config", enforce: true},
		{name: "account_profile", cfg: &config.Config{}, enforce: true},
		{name: "force_global", cfg: &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: true}}, enforce: true},
		{name: "force_global_with_pairing_only", cfg: &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: true}}, enforce: false},
	} {
		t.Run(policy.name, func(t *testing.T) {
			SetCodexIdentityEnforcementEnabled(policy.enforce)
			for _, profile := range []struct{ name, ua string }{
				{name: "cpa", ua: codexUAPolicyAccountUA},
				{name: "global_default"},
			} {
				t.Run(profile.name, func(t *testing.T) {
					account := newCodexModelsTestAccount()
					account.Concurrency = 3
					account.Credentials["user_agent"] = profile.ua
					wantUA := codexUAPolicyCanonicalUA
					if profile.ua != "" && (policy.cfg == nil || !policy.cfg.Gateway.ForceCodexCLI) {
						wantUA = codexUAPolicyNormalizedAccountUA
					}
					gateway := &OpenAIGatewayService{cfg: policy.cfg}
					checkHeaders := func(t *testing.T, headers http.Header) {
						t.Helper()
						require.Equal(t, wantUA, headers.Get("User-Agent"))
						require.Equal(t, "codex-tui", headers.Get("Originator"))
						if policy.enforce || headers.Get("Version") != "" {
							require.Equal(t, codexUAPolicyVersion, headers.Get("Version"))
						}
						require.Equal(t, "Bearer test-access-token", headers.Get("Authorization"))
						require.Equal(t, "acc-123", headers.Get("chatgpt-account-id"))
						require.Equal(t, profile.ua, account.GetOpenAIUserAgent())
						require.Equal(t, 3, account.Concurrency)
					}
					t.Run("responses_http", func(t *testing.T) {
						c, _ := newTestContext()
						c.Request.URL.Path = "/openai/v1/responses"
						c.Request.Header.Set("User-Agent", "different-local-client/1.0")
						request, err := gateway.buildUpstreamRequest(t.Context(), c, account, []byte(`{"model":"gpt-5.4","input":[]}`), "test-access-token", false, "", false)
						require.NoError(t, err)
						checkHeaders(t, request.Header)
					})
					t.Run("responses_ws", func(t *testing.T) {
						c, _ := newTestContext()
						c.Request.Header.Set("User-Agent", "different-local-client/1.0")
						headers, _, err := gateway.buildOpenAIWSHeaders(t.Context(), c, account, "test-access-token", OpenAIWSProtocolDecision{}, false, "", "", "", "gpt-5.4", "")
						require.NoError(t, err)
						checkHeaders(t, headers)
					})
					t.Run("models", func(t *testing.T) {
						var headers http.Header
						server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
							headers = request.Header.Clone()
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write([]byte(`{"models":[{"slug":"gpt-5.4"}]}`))
						}))
						defer server.Close()
						originalURL := chatgptCodexModelsURL
						chatgptCodexModelsURL = server.URL
						defer func() { chatgptCodexModelsURL = originalURL }()
						manifest, err := gateway.FetchCodexModelsManifest(t.Context(), account, "", "")
						require.NoError(t, err)
						require.NotNil(t, manifest)
						checkHeaders(t, headers)
					})
					for _, probe := range []struct {
						name     string
						response string
						run      func(*AccountTestService, *gin.Context, *Account) error
					}{
						{name: "normal_probe", response: "data: {\"type\":\"response.completed\"}\n\n", run: func(s *AccountTestService, c *gin.Context, a *Account) error {
							return s.testOpenAIAccountConnection(c, a, "gpt-5.4", "", "")
						}},
						{name: "compact_probe", response: compactProbeSSESuccessBody, run: func(s *AccountTestService, c *gin.Context, a *Account) error {
							return s.testOpenAICompactConnection(c, a, "gpt-5.4")
						}},
						{name: "image_probe", response: "data: {\"type\":\"response.output_item.done\",\"item\":{\"type\":\"image_generation_call\",\"result\":\"aGVsbG8=\",\"output_format\":\"png\"}}\n\ndata: {\"type\":\"response.completed\",\"response\":{\"output\":[]}}\n\n", run: func(s *AccountTestService, c *gin.Context, a *Account) error {
							return s.testOpenAIImageOAuth(c, context.Background(), a, "gpt-image-2", "draw a cat")
						}},
					} {
						t.Run(probe.name, func(t *testing.T) {
							upstream := &httpUpstreamRecorder{resp: &http.Response{
								StatusCode: http.StatusOK,
								Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
								Body:       io.NopCloser(strings.NewReader(probe.response)),
							}}
							service := &AccountTestService{cfg: policy.cfg, httpUpstream: upstream}
							c, _ := newTestContext()
							require.NoError(t, probe.run(service, c, account))
							require.NotNil(t, upstream.lastReq)
							checkHeaders(t, upstream.lastReq.Header)
						})
					}
				})
			}
		})
	}
}

// OAuth refresh must choose the same profile as inference even when invoked
// with a context inherited from a different account or a previous refresh.
func TestCodexUAPolicyOAuthRefreshMatchesInference(t *testing.T) {
	SetCodexCanonicalUserAgentResolver(func() string { return codexUAPolicyCanonicalUA })
	t.Cleanup(func() { SetCodexCanonicalUserAgentResolver(nil) })
	for _, policy := range []struct {
		name string
		cfg  *config.Config
	}{
		{name: "no_config"},
		{name: "account_profile", cfg: &config.Config{}},
		{name: "force_global", cfg: &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: true}}},
	} {
		t.Run(policy.name, func(t *testing.T) {
			client := &openaiOAuthClientRefreshStub{resp: &openai.TokenResponse{AccessToken: "new-token", ExpiresIn: 3600}}
			service := NewOpenAIOAuthService(nil, client)
			service.SetConfig(policy.cfg)
			t.Cleanup(service.Stop)
			parent := WithCodexAccountAuthIdentity(t.Context(), "codex_vscode/0.155.0 (Windows; x86_64)")
			for _, profile := range []struct{ name, ua, wantUA, wantOriginator string }{
				{name: "account_a", ua: codexUAPolicyAccountUA, wantUA: codexUAPolicyNormalizedAccountUA, wantOriginator: "codex-tui"},
				{name: "account_b", ua: "codex_vscode/0.153.3 (Windows; x86_64)", wantUA: "codex_vscode/0.146.0 (Windows; x86_64)", wantOriginator: "codex_vscode"},
				{name: "no_account_profile", wantUA: codexUAPolicyCanonicalUA, wantOriginator: "codex-tui"},
				{name: "account_a_again", ua: codexUAPolicyAccountUA, wantUA: codexUAPolicyNormalizedAccountUA, wantOriginator: "codex-tui"},
			} {
				t.Run(profile.name, func(t *testing.T) {
					account := &Account{ID: 801, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Concurrency: 3,
						Credentials: map[string]any{"refresh_token": "RT-test-account-token", "client_id": "test-client", "user_agent": profile.ua},
					}
					info, err := service.RefreshAccountToken(parent, account)
					require.NoError(t, err)
					wantUA, wantOriginator := profile.wantUA, profile.wantOriginator
					if policy.cfg != nil && policy.cfg.Gateway.ForceCodexCLI {
						wantUA, wantOriginator = codexUAPolicyCanonicalUA, "codex-tui"
					}
					require.Equal(t, wantUA, client.lastUserAgent)
					require.Equal(t, wantOriginator, client.lastOriginator)
					require.Equal(t, "test-client", client.lastClientID)
					require.Equal(t, "RT-test-account-token", info.RefreshToken)
					require.Equal(t, profile.ua, account.GetOpenAIUserAgent())
					require.Equal(t, 3, account.Concurrency)
				})
			}
		})
	}
}

// The PAT branch of the same account refresh uses a different auth endpoint.
// Its identity profile must still match, without inference-only version headers.
func TestCodexUAPolicyPATRefreshAndDirectValidation(t *testing.T) {
	SetCodexCanonicalUserAgentResolver(func() string { return codexUAPolicyCanonicalUA })
	t.Cleanup(func() { SetCodexCanonicalUserAgentResolver(nil) })
	for _, test := range []struct {
		name                          string
		forceGlobal, directValidation bool
		wantUA                        string
	}{
		{name: "account_profile", wantUA: codexUAPolicyNormalizedAccountUA},
		{name: "force_global", forceGlobal: true, wantUA: codexUAPolicyCanonicalUA},
		{name: "direct_validation_without_account", directValidation: true, wantUA: codexUAPolicyCanonicalUA},
	} {
		t.Run(test.name, func(t *testing.T) {
			var headers http.Header
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				headers = request.Header.Clone()
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"email":"user@example.test","chatgpt_user_id":"user-test","chatgpt_account_id":"account-test","chatgpt_plan_type":"plus","chatgpt_account_is_fedramp":false}`))
			}))
			defer server.Close()
			originalURL := openAICodexPATWhoamiURL
			openAICodexPATWhoamiURL = server.URL
			defer func() { openAICodexPATWhoamiURL = originalURL }()
			service := NewOpenAIOAuthService(nil, nil)
			service.SetConfig(&config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: test.forceGlobal}})
			t.Cleanup(service.Stop)
			account := &Account{ID: 802, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Concurrency: 3,
				Credentials: map[string]any{"access_token": "at-test-token", "auth_mode": OpenAIAuthModePersonalAccessToken, "user_agent": codexUAPolicyAccountUA},
			}
			var info *OpenAITokenInfo
			var err error
			if test.directValidation {
				info, err = service.ValidateCodexPersonalAccessToken(t.Context(), "at-test-token", "")
			} else {
				inherited := WithCodexAccountAuthIdentity(t.Context(), "codex_vscode/0.155.0 (Windows; x86_64)")
				info, err = service.RefreshAccountToken(inherited, account)
			}
			require.NoError(t, err)
			require.NotNil(t, headers)
			require.Equal(t, test.wantUA, headers.Get("User-Agent"))
			require.Equal(t, "codex-tui", headers.Get("Originator"))
			require.Empty(t, headers.Get("Version"), "auth whoami must retain its UA/originator-only protocol")
			require.Equal(t, "Bearer at-test-token", headers.Get("Authorization"))
			require.Equal(t, OpenAIAuthModePersonalAccessToken, info.AuthMode)
			require.Equal(t, "account-test", info.ChatGPTAccountID)
			require.Equal(t, codexUAPolicyAccountUA, account.GetOpenAIUserAgent())
			require.Equal(t, 3, account.Concurrency)
		})
	}
}
