package repository

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestOpenAIRefreshRejectsCredentialRedirects(t *testing.T) {
	for _, status := range []int{301, 302, 303, 307, 308} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			var destinationCalls atomic.Int64
			destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				destinationCalls.Add(1)
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{"access_token":"redirected-token"}`)
			}))
			defer destination.Close()
			const secret = "RT-do-not-forward-or-log"
			official := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, destination.URL+"/?echo="+secret, status)
			}))
			defer official.Close()
			svc := &openaiOAuthService{tokenURL: official.URL, legacyTokenURL: destination.URL}
			result, err := svc.RefreshTokenWithClientID(t.Context(), secret, "", openai.ClientID)
			require.ErrorContains(t, err, "status "+strconv.Itoa(status))
			require.NotContains(t, err.Error(), secret)
			require.Nil(t, result)
			require.Zero(t, destinationCalls.Load(), "a redirect must not change the credential recipient")

			// The shared req client keeps its original redirect policy for other,
			// non-refresh calls. Refresh hardening must not mutate the global pool.
			shared, err := createOpenAIReqClient("")
			require.NoError(t, err)
			resp, err := shared.R().SetContext(t.Context()).Get(official.URL)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			require.Equal(t, int64(1), destinationCalls.Load())
		})
	}
}

func TestOpenAIRefreshKeepsConfiguredProxy(t *testing.T) {
	var proxyCalls atomic.Int64
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyCalls.Add(1)
		require.Equal(t, "http://official.invalid/oauth/token", r.RequestURI)
		require.NoError(t, r.ParseForm())
		require.Equal(t, "xy-opaque-refresh", r.Form.Get("refresh_token"))
		require.Equal(t, "custom-client-id", r.Form.Get("client_id"))
		require.Equal(t, openai.RefreshScopes, r.Form.Get("scope"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"proxied-access-token","refresh_token":"rotated","expires_in":3600}`)
	}))
	defer proxy.Close()
	svc := &openaiOAuthService{tokenURL: "http://official.invalid/oauth/token"}
	result, err := svc.RefreshTokenWithClientID(t.Context(), "xy-opaque-refresh", proxy.URL, "custom-client-id")
	require.NoError(t, err)
	require.Equal(t, "proxied-access-token", result.AccessToken)
	require.Equal(t, "rotated", result.RefreshToken)
	require.Equal(t, int64(1), proxyCalls.Load())
}

func TestOpenAIRefreshFailuresDoNotExposeProviderResponse(t *testing.T) {
	const secret = "RT-sensitive-refresh-value"
	for _, test := range []struct {
		name   string
		status int
		body   string
	}{
		{"error echo", http.StatusBadRequest, `{"refresh_token":"` + secret + `"}`},
		{"malformed JSON", http.StatusOK, `{"access_token":` + secret},
		{"missing access token", http.StatusOK, `{"refresh_token":"` + secret + `"}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(test.status)
				_, _ = io.WriteString(w, test.body)
			}))
			defer server.Close()
			svc := &openaiOAuthService{tokenURL: server.URL}
			result, err := svc.RefreshTokenWithClientID(t.Context(), secret, "", openai.ClientID)
			require.Error(t, err)
			require.Nil(t, result)
			require.NotContains(t, err.Error(), secret)
			require.NotContains(t, err.Error(), strings.TrimSpace(test.body))
		})
	}
}

func TestOpenAIRefreshUsesAccountAuthIdentity(t *testing.T) {
	service.SetCodexCanonicalUserAgentResolver(func() string { return "codex_cli_rs/0.200.1 (Linux; x86_64)" })
	t.Cleanup(func() { service.SetCodexCanonicalUserAgentResolver(nil) })
	captured := make(chan http.Header, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured <- r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"access-token","expires_in":3600}`)
	}))
	defer server.Close()
	svc := &openaiOAuthService{tokenURL: server.URL}
	for _, test := range []struct {
		input, wantUA, wantOriginator string
	}{
		{"codex_cli_rs/0.100.0 (Mac OS X 14.0; arm64)", "codex_cli_rs/0.200.1 (Mac OS X 14.0; arm64)", "codex_cli_rs"},
		{"codex_vscode/0.101.0 (Windows; x86_64)", "codex_vscode/0.200.1 (Windows; x86_64)", "codex_vscode"},
		{"", "codex_cli_rs/0.200.1 (Linux; x86_64)", "codex_cli_rs"},
	} {
		ctx := service.WithCodexAccountAuthIdentity(t.Context(), test.input)
		_, err := svc.RefreshTokenWithClientID(ctx, "RT-private-token", "", openai.ClientID)
		require.NoError(t, err)
		headers := <-captured
		require.Equal(t, test.wantUA, headers.Get("User-Agent"))
		require.Equal(t, test.wantOriginator, headers.Get("originator"))
		require.Empty(t, headers.Get("version"), "the authentication endpoint must not receive the inference-only version header")
		require.NotContains(t, headers.Get("User-Agent"), "RT-private-token")
	}

	ctx := service.WithCodexAccountAuthIdentity(t.Context(), "codex_vscode/0.100.0 (Windows; x86_64)")
	_, err := svc.ExchangeCode(ctx, "code", "verifier", "", "", openai.ClientID)
	require.NoError(t, err)
	headers := <-captured
	canonicalUA, canonicalOriginator := service.CodexCanonicalAuthIdentity()
	require.Equal(t, canonicalUA, headers.Get("User-Agent"), "initial code exchange has no selected account and keeps canonical metadata")
	require.Equal(t, canonicalOriginator, headers.Get("originator"))
}

func TestOpenAIRefreshPreservesSafeFailureCodesWithoutTokenEcho(t *testing.T) {
	const secret = "RT-sensitive-refresh-error-echo"
	for _, code := range []string{
		"invalid_grant", "invalid_refresh_token", "token_expired", "app_session_terminated",
		"refresh_token_reused", "refresh_token_invalidated", "invalid_client", "unauthorized_client",
		"access_denied", "missing_project_id", "invalid_scope", "entitlement_denied",
	} {
		for _, shape := range []string{"string", "code", "type"} {
			t.Run(code+"/"+shape, func(t *testing.T) {
				var errorValue any = code
				if shape != "string" {
					errorValue = map[string]string{shape: code, "message": secret}
				}
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					_ = json.NewEncoder(w).Encode(map[string]any{"error": errorValue, "echo": secret})
				}))
				defer server.Close()
				svc := &openaiOAuthService{tokenURL: server.URL}
				result, err := svc.RefreshTokenWithClientID(t.Context(), secret, "", openai.ClientID)
				require.Nil(t, result)
				require.ErrorContains(t, err, "status 401, code: "+code)
				require.NotContains(t, err.Error(), secret)
				require.NotContains(t, err.Error(), `"message":`)
			})
		}
	}
}

func TestOpenAIRefreshDropsUnknownAndOversizedErrors(t *testing.T) {
	const secret = "RT-sensitive-refresh-error-echo"
	for _, test := range []struct {
		name, body string
	}{
		{"unknown string", `{"error":"` + secret + `"}`},
		{"unknown object code", `{"error":{"code":"` + secret + `","message":"invalid_grant"}}`},
		{"unknown object type", `{"error":{"type":"` + secret + `"}}`},
		{"known code substring is not allowed", `{"error":"invalid_grant: ` + secret + `"}`},
		{"oversized known error", `{"error":"invalid_grant","echo":"` + strings.Repeat("x", 64<<10) + secret + `"}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = io.WriteString(w, test.body)
			}))
			defer server.Close()
			svc := &openaiOAuthService{tokenURL: server.URL}
			result, err := svc.RefreshTokenWithClientID(t.Context(), secret, "", openai.ClientID)
			require.Nil(t, result)
			require.ErrorContains(t, err, "status 400")
			require.NotContains(t, err.Error(), "code:")
			require.NotContains(t, err.Error(), "invalid_grant")
			require.NotContains(t, err.Error(), secret)
		})
	}
}
