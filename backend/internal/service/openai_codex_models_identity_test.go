package service

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestFetchCodexModelsManifestOAuthUserAgentMatchesNegotiatedVersion(t *testing.T) {
	const canonicalVersion = "0.146.0"
	const accountUA = "codex-tui/0.153.3 (Mac OS 26.5.1; arm64) iTerm.app/3.6.11 (codex-tui; 0.153.3)"
	const manifestBody = `{"models":[{"slug":"gpt-5.5","display_name":"GPT-5.5","unknown":{"keep":true}}]}`
	SetCodexCanonicalUserAgentResolver(func() string { return buildCodexCLIUserAgent(canonicalVersion) })
	t.Cleanup(func() { SetCodexCanonicalUserAgentResolver(nil) })

	for _, forceCLI := range []bool{false, true} {
		for _, test := range []struct {
			name, queryVersion, headerVersion string
		}{
			{name: "newer", queryVersion: "0.153.3", headerVersion: "0.153.3"},
			{name: "supported_older", queryVersion: "0.144.0", headerVersion: "0.144.0"},
			{name: "below_minimum", queryVersion: "0.1.0", headerVersion: canonicalVersion},
			{name: "invalid", queryVersion: "not-a-version", headerVersion: canonicalVersion},
		} {
			t.Run(fmt.Sprintf("force_cli=%t/%s", forceCLI, test.name), func(t *testing.T) {
				var received *http.Request
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
					received = request.Clone(context.Background())
					w.Header().Set("ETag", `W/"account-manifest"`)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(manifestBody))
				}))
				defer server.Close()
				originalURL := chatgptCodexModelsURL
				chatgptCodexModelsURL = server.URL
				defer func() { chatgptCodexModelsURL = originalURL }()
				account := newCodexModelsTestAccount()
				account.Credentials["user_agent"] = accountUA
				svc := &OpenAIGatewayService{cfg: &config.Config{Gateway: config.GatewayConfig{ForceCodexCLI: forceCLI}}}

				manifest, err := svc.FetchCodexModelsManifest(context.Background(), account, test.queryVersion, `W/"previous"`)
				require.NoError(t, err)
				require.NotNil(t, received)
				require.Equal(t, test.queryVersion, received.URL.Query().Get("client_version"))
				require.Equal(t, test.headerVersion, received.Header.Get("Version"))
				wantUA := "codex-tui/" + test.headerVersion + " (Mac OS 26.5.1; arm64) iTerm.app/3.6.11 (codex-tui; " + test.headerVersion + ")"
				if forceCLI {
					wantUA = buildCodexCLIUserAgent(test.headerVersion)
				}
				require.Equal(t, wantUA, received.Header.Get("User-Agent"))
				require.Equal(t, "codex-tui", received.Header.Get("Originator"))
				require.Equal(t, "Bearer test-access-token", received.Header.Get("Authorization"))
				require.Equal(t, "acc-123", received.Header.Get("chatgpt-account-id"))
				require.Equal(t, `W/"previous"`, received.Header.Get("If-None-Match"))
				require.Equal(t, manifestBody, string(manifest.Body))
				require.Equal(t, `W/"account-manifest"`, manifest.ETag)
				require.Equal(t, accountUA, account.GetOpenAIUserAgent())
			})
		}
	}
}
