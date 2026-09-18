//go:build unit

package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	coderws "github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type boundEgressUpstream struct{ calls atomic.Int32 }

func (u *boundEgressUpstream) Do(*http.Request, string, int64, int) (*http.Response, error) {
	u.calls.Add(1)
	return nil, errors.New("test upstream reached")
}
func (u *boundEgressUpstream) DoWithTLS(r *http.Request, p string, a int64, c int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.Do(r, p, a, c)
}

type boundEgressDialer struct{ calls atomic.Int32 }

func (d *boundEgressDialer) Dial(context.Context, string, http.Header, string) (openAIWSClientConn, int, http.Header, error) {
	d.calls.Add(1)
	return nil, 0, nil, errors.New("test dial reached")
}

type boundEgressAccountRepo struct {
	AccountRepository
	account *Account
	calls   int
}

func (r *boundEgressAccountRepo) GetByID(context.Context, int64) (*Account, error) {
	r.calls++
	return r.account, nil
}
func (r *boundEgressAccountRepo) UpdateExtra(context.Context, int64, map[string]any) error {
	return nil
}

func boundEgressTestAccount() *Account {
	p := int64(75)
	return &Account{ID: 8751, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 3, ProxyID: &p, Credentials: map[string]any{"access_token": "unit-token", "chatgpt_account_id": "unit-account", "chatgpt_account_is_fedramp": false}, Extra: map[string]any{codexFingerprintModeExtraKey: "full", codexFingerprintSeedExtraKey: testCodexFingerprintSeed}}
}

func TestOpenAIBoundEgressRejectsMissingProxyBeforeHTTP(t *testing.T) {
	for _, path := range []string{"responses", "passthrough", "messages", "alpha", "live", "live_sideband", "probe", "compact_probe"} {
		t.Run(path, func(t *testing.T) {
			account := boundEgressTestAccount()
			upstream := &boundEgressUpstream{}
			svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream, cache: &stubGatewayCache{}, toolCorrector: NewCodexToolCorrector()}
			c, _ := newTestContext()
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
			var err error
			switch path {
			case "responses", "passthrough":
				account.Extra["openai_oauth_passthrough"] = path == "passthrough"
				_, err = svc.Forward(context.Background(), c, account, []byte(`{"model":"gpt-5.4","input":"hello","stream":false}`))
			case "messages":
				_, err = svc.ForwardAsAnthropic(context.Background(), c, account, []byte(`{"model":"gpt-5.4","max_tokens":128,"messages":[{"role":"user","content":"hello"}],"stream":false}`), "", "")
			case "alpha":
				_, err = svc.ForwardAlphaSearch(context.Background(), c, account, []byte(`{"id":"unit-search","model":"gpt-5.4","commands":{"search_query":[{"q":"test"}]}}`))
			case "live":
				_, err = svc.createUpstreamLiveCall(context.Background(), account, &LiveCallRequest{SDP: "v=offer\r\n", Session: json.RawMessage(`{"model":"gpt-live-test"}`)}, `{"v":1}`)
			case "live_sideband":
				_, err = svc.liveSidebandHeaders(context.Background(), account, &LiveCallRecord{CallID: "unit-call"})
			case "probe", "compact_probe":
				repo := &boundEgressAccountRepo{account: account}
				tester := &AccountTestService{accountRepo: repo, httpUpstream: upstream}
				mode := ""
				if path == "compact_probe" {
					mode = AccountTestModeCompact
				}
				err = tester.TestAccountConnection(c, account.ID, "gpt-5.4", "", mode)
				require.Equal(t, 1, repo.calls, "guard must reuse the loaded account")
			}
			require.Error(t, err)
			require.Zero(t, upstream.calls.Load(), "missing bound proxy must never reach upstream")
			require.Contains(t, strings.ToLower(err.Error()), "proxy")
			require.Equal(t, 3, account.Concurrency)
		})
	}
}

func TestOpenAIBoundEgressGetAccessTokenPreservesShadowAndOtherProviders(t *testing.T) {
	t.Run("missing child proxy stops before parent lookup", func(t *testing.T) {
		account := boundEgressTestAccount()
		parentID := int64(91)
		account.ParentAccountID = &parentID
		repo := &boundEgressAccountRepo{account: boundEgressTestAccount()}
		svc := &OpenAIGatewayService{accountRepo: repo}
		_, _, err := svc.GetAccessToken(context.Background(), account)
		require.Error(t, err)
		require.Zero(t, repo.calls)
	})
	t.Run("valid child egress retains parent credentials", func(t *testing.T) {
		account := boundEgressTestAccount()
		parentID := int64(91)
		account.ParentAccountID = &parentID
		account.Proxy = &Proxy{ID: *account.ProxyID, Protocol: "http", Host: "127.0.0.1", Port: 18080}
		parent := boundEgressTestAccount()
		parent.ID = parentID
		parent.ProxyID = nil
		parent.Credentials["access_token"] = "parent-token"
		repo := &boundEgressAccountRepo{account: parent}
		svc := &OpenAIGatewayService{accountRepo: repo}
		token, _, err := svc.GetAccessToken(context.Background(), account)
		require.NoError(t, err)
		require.Equal(t, "parent-token", token)
		require.Equal(t, 1, repo.calls)
		require.Equal(t, int64(75), *account.ProxyID)
		require.Nil(t, parent.Proxy)
	})
	for _, kind := range []string{"unbound_oauth", "apikey", "grok_oauth"} {
		t.Run(kind, func(t *testing.T) {
			account := boundEgressTestAccount()
			switch kind {
			case "unbound_oauth":
				account.ProxyID = nil
			case "apikey":
				account.Type = AccountTypeAPIKey
				account.Credentials["api_key"] = "unit-key"
			case "grok_oauth":
				account.Platform = PlatformGrok
				account.Credentials["sso"] = "unit-sso"
			}
			_, _, err := (&OpenAIGatewayService{}).GetAccessToken(context.Background(), account)
			require.NoError(t, err)
		})
	}
}

func TestOpenAIBoundEgressRejectsMissingProxyAtWebSocketIngress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, mode := range []string{OpenAIWSIngressModeCtxPool, OpenAIWSIngressModePassthrough, OpenAIWSIngressModeHTTPBridge} {
		t.Run(mode, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.Gateway.OpenAIWS.Enabled = true
			cfg.Gateway.OpenAIWS.OAuthEnabled = true
			cfg.Gateway.OpenAIWS.ResponsesWebsocketsV2 = true
			cfg.Gateway.OpenAIWS.ModeRouterV2Enabled = true
			cfg.Gateway.OpenAIWS.MaxConnsPerAccount = 3
			cfg.Gateway.OpenAIWS.DialTimeoutSeconds = 1
			account := boundEgressTestAccount()
			account.Extra["openai_oauth_responses_websockets_v2_mode"] = mode
			upstream := &boundEgressUpstream{}
			dialer := &boundEgressDialer{}
			pool := newOpenAIWSConnPool(cfg)
			t.Cleanup(pool.Close)
			pool.setClientDialerForTest(dialer)
			svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstream, cache: &stubGatewayCache{}, toolCorrector: NewCodexToolCorrector(), openaiWSResolver: NewOpenAIWSProtocolResolver(cfg), openaiWSPool: pool, openaiWSPassthroughDialer: dialer}
			outcome := make(chan error, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, err := coderws.Accept(w, r, nil)
				if err != nil {
					outcome <- err
					return
				}
				defer conn.CloseNow()
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				c.Request = r
				outcome <- svc.ProxyResponsesWebSocketFromClient(r.Context(), c, conn, account, "unit-token", []byte(`{"type":"response.create","model":"gpt-5.4","input":"hello"}`), nil)
			}))
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			conn, _, err := coderws.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http"), nil)
			require.NoError(t, err)
			defer conn.CloseNow()
			select {
			case err = <-outcome:
			case <-ctx.Done():
				t.Fatal("gateway did not reject missing proxy")
			}
			require.Error(t, err)
			require.Zero(t, dialer.calls.Load())
			require.Zero(t, upstream.calls.Load())
			require.Contains(t, strings.ToLower(err.Error()), "proxy")
			require.Equal(t, 3, account.Concurrency)
		})
	}
}

func TestOpenAIBoundEgressAgentTaskRegistrationDoesNotUseDefaultEgress(t *testing.T) {
	_, privateKey := newTestAgentIdentityKey(t)
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		_, _ = w.Write([]byte(`{"task_id":"unit-task"}`))
	}))
	defer server.Close()
	oldBase := openAIAgentIdentityAuthAPIBaseURL
	openAIAgentIdentityAuthAPIBaseURL = server.URL
	t.Cleanup(func() { openAIAgentIdentityAuthAPIBaseURL = oldBase })
	account := boundEgressTestAccount()
	account.Credentials[openAIAuthModeCredentialKey] = OpenAIAuthModeAgentIdentity
	account.Credentials["agent_private_key"] = privateKey
	account.Credentials["agent_runtime_id"] = "runtime-test"
	taskID, err := registerAgentIdentityTask(context.Background(), account)
	require.Error(t, err)
	require.Empty(t, taskID)
	require.Zero(t, calls.Load())
	require.Contains(t, strings.ToLower(err.Error()), "proxy")
}
