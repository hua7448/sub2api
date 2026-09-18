package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	httppool "github.com/Wei-Shaw/sub2api/internal/pkg/httpclient"
	"github.com/stretchr/testify/require"
)

type auxiliaryIdentityRoundTripper func(*http.Request) (*http.Response, error)

func (f auxiliaryIdentityRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func captureCodexSnapshotProbe(t *testing.T, receive func(*http.Request)) {
	t.Helper()
	client, err := httppool.GetClient(httppool.Options{Timeout: 15 * time.Second, ResponseHeaderTimeout: 10 * time.Second})
	require.NoError(t, err)
	original := client.Transport
	client.Transport = auxiliaryIdentityRoundTripper(func(req *http.Request) (*http.Response, error) {
		receive(req)
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("data: [DONE]\n\n")), Request: req}, nil
	})
	t.Cleanup(func() { client.Transport = original })
}

func TestCodexBoundProxyAuxiliaryRequestsStopBeforeSending(t *testing.T) {
	for _, path := range []string{"models", "usage_probe"} {
		t.Run(path, func(t *testing.T) {
			account := newCodexModelsTestAccount()
			proxyID := int64(904)
			account.ProxyID = &proxyID
			account.Proxy = nil
			calls := 0
			var err error
			if path == "models" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					calls++
					_, _ = io.WriteString(w, `{"models":[{"slug":"gpt-5.4"}]}`)
				}))
				defer server.Close()
				original := chatgptCodexModelsURL
				chatgptCodexModelsURL = server.URL
				t.Cleanup(func() { chatgptCodexModelsURL = original })
				_, err = (&OpenAIGatewayService{}).FetchCodexModelsManifest(context.Background(), account, "", "")
			} else {
				captureCodexSnapshotProbe(t, func(*http.Request) { calls++ })
				_, err = (&AccountUsageService{}).probeOpenAICodexSnapshot(context.Background(), account)
			}
			require.Error(t, err)
			require.Zero(t, calls, "a missing bound proxy must stop auxiliary network traffic")
		})
	}
}

func TestCodexUAPolicyPassiveProbeUsesWiredGatewayPolicy(t *testing.T) {
	const accountUA = "codex-tui/0.153.3 (Mac OS 26.5.1; arm64) iTerm.app/3.6.11 (codex-tui; 0.153.3)"
	const version = "0.153.3"
	SetCodexCanonicalUserAgentResolver(func() string { return buildCodexCLIUserAgent(version) })
	previous := codexIdentityEnforcement.Load()
	SetCodexIdentityEnforcementEnabled(true)
	t.Cleanup(func() { SetCodexCanonicalUserAgentResolver(nil); SetCodexIdentityEnforcementEnabled(previous) })
	for _, force := range []bool{false, true} {
		name := "account_profile"
		if force {
			name = "forced_global"
		}
		t.Run(name, func(t *testing.T) {
			account := newCodexModelsTestAccount()
			account.Credentials["user_agent"] = accountUA
			cfg := &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: force}}
			service := ProvideAccountUsageService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &OpenAIGatewayService{cfg: cfg})
			var headers http.Header
			captureCodexSnapshotProbe(t, func(req *http.Request) { headers = req.Header.Clone() })
			_, err := service.probeOpenAICodexSnapshot(context.Background(), account)
			require.NoError(t, err)
			require.NotNil(t, headers)
			wantUA := accountUA
			if force {
				wantUA = buildCodexCLIUserAgent(version)
			}
			require.Equal(t, wantUA, headers.Get("User-Agent"))
			require.Equal(t, version, headers.Get("Version"))
			require.Equal(t, "codex-tui", headers.Get("Originator"))
			require.Equal(t, accountUA, account.GetOpenAIUserAgent(), "effective policy does not rewrite saved profile")
		})
	}
}
