package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/stretchr/testify/require"
)

func TestCodexAccountAuthIdentityUsesCurrentVersionAndSeparateContexts(t *testing.T) {
	version := "0.200.1"
	SetCodexCanonicalUserAgentResolver(func() string { return "codex_cli_rs/" + version + " (Linux; x86_64)" })
	t.Cleanup(func() { SetCodexCanonicalUserAgentResolver(nil) })
	ctxA := WithCodexAccountAuthIdentity(t.Context(), "codex_cli_rs/0.100.0 (Mac OS X 14.0; arm64)")
	ctxB := WithCodexAccountAuthIdentity(t.Context(), "codex_vscode/0.101.0 (Windows; x86_64)")
	uaA, originatorA := CodexAuthIdentityFromContext(ctxA)
	uaB, originatorB := CodexAuthIdentityFromContext(ctxB)
	require.Equal(t, "codex_cli_rs/0.200.1 (Mac OS X 14.0; arm64)", uaA)
	require.Equal(t, "codex_cli_rs", originatorA)
	require.Equal(t, "codex_vscode/0.200.1 (Windows; x86_64)", uaB)
	require.Equal(t, "codex_vscode", originatorB)
	canonicalUA, canonicalOriginator := CodexCanonicalAuthIdentity()
	defaultUA, defaultOriginator := CodexAuthIdentityFromContext(t.Context())
	require.Equal(t, canonicalUA, defaultUA)
	require.Equal(t, canonicalOriginator, defaultOriginator)
	invalidUA, invalidOriginator := CodexAuthIdentityFromContext(WithCodexAccountAuthIdentity(t.Context(), "unrecognized-client/1.0"))
	require.Equal(t, canonicalUA, invalidUA)
	require.Equal(t, canonicalOriginator, invalidOriginator)

	version = "0.201.2"
	updatedUA, _ := CodexAuthIdentityFromContext(ctxA)
	require.Equal(t, "0.201.2", openai.CodexUserAgentVersion(updatedUA), "an existing account setting must follow version synchronization")
}

type openAIRefreshIdentityProxyRepo struct {
	ProxyRepository
	proxy *Proxy
}

func (s *openAIRefreshIdentityProxyRepo) GetByID(context.Context, int64) (*Proxy, error) {
	return s.proxy, nil
}

func TestOpenAIOAuthRefreshAccountCarriesItsIdentityAndProxy(t *testing.T) {
	SetCodexCanonicalUserAgentResolver(func() string { return "codex_cli_rs/0.200.1 (Linux; x86_64)" })
	t.Cleanup(func() { SetCodexCanonicalUserAgentResolver(nil) })
	proxy := &Proxy{ID: 7, Protocol: "http", Host: "proxy.example", Port: 8080}
	client := &openaiOAuthClientRefreshStub{resp: &openai.TokenResponse{AccessToken: "access-token", ExpiresIn: 3600}}
	svc := NewOpenAIOAuthService(&openAIRefreshIdentityProxyRepo{proxy: proxy}, client)
	t.Cleanup(svc.Stop)
	// A previous request's metadata must not override the selected account.
	parent := WithCodexAccountAuthIdentity(t.Context(), "codex_vscode/0.123.0 (Other; arm64)")
	for _, userAgent := range []string{
		"codex_cli_rs/0.100.0 (Mac OS X 14.0; arm64)",
		"codex_vscode/0.101.0 (Windows; x86_64)",
		"",
	} {
		account := &Account{
			ID: 42, Platform: PlatformOpenAI, Type: AccountTypeOAuth, ProxyID: &proxy.ID,
			Credentials: map[string]any{"refresh_token": "RT-private-account-token", "client_id": "custom-client", "user_agent": userAgent},
		}
		info, err := svc.RefreshAccountToken(parent, account)
		require.NoError(t, err)
		want := resolveCodexOutboundIdentity(userAgent)
		require.Equal(t, want.userAgent, client.lastUserAgent)
		require.Equal(t, want.originator, client.lastOriginator)
		require.Equal(t, proxy.URL(), client.lastProxyURL)
		require.Equal(t, "custom-client", client.lastClientID)
		require.Equal(t, "RT-private-account-token", info.RefreshToken, "an absent rotated token must preserve the input token")
		require.NotContains(t, client.lastUserAgent, "RT-private-account-token")
	}
}
