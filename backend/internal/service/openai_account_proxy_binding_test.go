package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type openAIAccountProxyBindingRepo struct {
	ProxyRepository
	proxy  *Proxy
	err    error
	calls  int
	lastID int64
}

func (r *openAIAccountProxyBindingRepo) GetByID(_ context.Context, id int64) (*Proxy, error) {
	r.calls++
	r.lastID = id
	return r.proxy, r.err
}

type openAIAccountProxyBindingTokenCache struct {
	*stubQuotaTokenCache
	calls int
}

func (c *openAIAccountProxyBindingTokenCache) GetAccessToken(ctx context.Context, key string) (string, error) {
	c.calls++
	return c.stubQuotaTokenCache.GetAccessToken(ctx, key)
}

func openAIAccountProxyBindingAccount() *Account {
	proxyID := int64(42)
	return &Account{ID: 71, Platform: PlatformOpenAI, Type: AccountTypeOAuth, ProxyID: &proxyID,
		Credentials: map[string]any{"access_token": "access-binding-test", "refresh_token": "refresh-binding-test", "chatgpt_account_id": "account-binding-test"}}
}
func openAIAccountProxyBindingProxy() *Proxy {
	return &Proxy{ID: 42, Protocol: "http", Host: "proxy.example", Port: 8080, Username: "binding-user", Password: "binding-secret"}
}

func TestOpenAIAccountProxyBinding_RefreshRejectsUnavailableBeforeAnyUpstream(t *testing.T) {
	for _, credentialMode := range []string{"refresh", "access_only", "personal_access_token"} {
		for _, scenario := range []string{"missing_repository", "lookup_failure", "lookup_missing", "invalid_protocol", "empty_host", "invalid_port", "mismatched_relation"} {
			t.Run(credentialMode+"/"+scenario, func(t *testing.T) {
				account := openAIAccountProxyBindingAccount()
				var repo ProxyRepository
				switch scenario {
				case "lookup_failure":
					repo = &openAIAccountProxyBindingRepo{err: errors.New("postgres failed binding-secret")}
				case "lookup_missing":
					repo = &openAIAccountProxyBindingRepo{}
				case "invalid_protocol":
					p := openAIAccountProxyBindingProxy()
					p.Protocol = "file"
					account.Proxy = p
					repo = &openAIAccountProxyBindingRepo{proxy: p}
				case "empty_host":
					p := openAIAccountProxyBindingProxy()
					p.Host = ""
					account.Proxy = p
					repo = &openAIAccountProxyBindingRepo{proxy: p}
				case "invalid_port":
					p := openAIAccountProxyBindingProxy()
					p.Port = 0
					account.Proxy = p
					repo = &openAIAccountProxyBindingRepo{proxy: p}
				case "mismatched_relation":
					p := openAIAccountProxyBindingProxy()
					p.ID = 43
					account.Proxy = p
				}
				client := &openaiOAuthClientRefreshStub{resp: &openai.TokenResponse{AccessToken: "new-binding-token", ExpiresIn: 3600}}
				service := NewOpenAIOAuthService(repo, client)
				defer service.Stop()
				privacyCalls := 0
				service.SetPrivacyClientFactory(func(string) (*req.Client, error) { privacyCalls++; return nil, errors.New("blocked test transport") })
				var patCalls atomic.Int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { patCalls.Add(1); w.WriteHeader(http.StatusUnauthorized) }))
				defer server.Close()
				originalPATURL := openAICodexPATWhoamiURL
				openAICodexPATWhoamiURL = server.URL
				defer func() { openAICodexPATWhoamiURL = originalPATURL }()
				if credentialMode != "refresh" {
					delete(account.Credentials, "refresh_token")
				}
				if credentialMode == "personal_access_token" {
					account.Credentials["access_token"] = "at-binding-test"
					account.Credentials["auth_mode"] = OpenAIAuthModePersonalAccessToken
				}
				_, err := service.RefreshAccountToken(context.Background(), account)
				assert.Zero(t, atomic.LoadInt32(&client.refreshCalls), "bound-proxy failure must stop token refresh")
				assert.Zero(t, privacyCalls, "bound-proxy failure must stop enrichment")
				assert.Zero(t, patCalls.Load(), "bound-proxy failure must stop PAT validation")
				require.Error(t, err)
				require.NotContains(t, err.Error(), "binding-secret")
			})
		}
	}
}

func TestOpenAIAccountProxyBinding_QuotaRejectsUnavailableBeforeTokenOrClient(t *testing.T) {
	for _, operation := range []string{"query", "reset"} {
		for _, scenario := range []string{"missing_repository", "lookup_failure", "lookup_missing", "invalid_port"} {
			t.Run(operation+"/"+scenario, func(t *testing.T) {
				account := openAIAccountProxyBindingAccount()
				repo := &stubQuotaAccountRepo{accounts: map[int64]*Account{account.ID: account}}
				cache := &openAIAccountProxyBindingTokenCache{stubQuotaTokenCache: &stubQuotaTokenCache{tokens: map[string]string{OpenAITokenCacheKey(account): "cached-binding-token"}}}
				var proxyRepo ProxyRepository
				switch scenario {
				case "lookup_failure":
					proxyRepo = &openAIAccountProxyBindingRepo{err: errors.New("postgres failed binding-secret")}
				case "lookup_missing":
					proxyRepo = &openAIAccountProxyBindingRepo{}
				case "invalid_port":
					p := openAIAccountProxyBindingProxy()
					p.Port = 65536
					account.Proxy = p
				}
				factoryCalls := 0
				service := NewOpenAIQuotaService(repo, proxyRepo, NewOpenAITokenProvider(repo, cache, nil), func(string) (*req.Client, error) { factoryCalls++; return nil, errors.New("blocked test transport") })
				var err error
				if operation == "query" {
					_, err = service.QueryUsage(context.Background(), account.ID)
				} else {
					_, err = service.ResetCredit(context.Background(), account.ID)
				}
				assert.Zero(t, cache.calls, "resolve bound egress before token acquisition can initiate refresh")
				assert.Zero(t, factoryCalls, "unavailable proxy must not construct a default-egress client")
				require.Error(t, err)
				require.NotContains(t, err.Error(), "binding-secret")
			})
		}
	}
}

func TestOpenAIAccountProxyBinding_RefreshPreservesSelectedProxy(t *testing.T) {
	for _, scenario := range []string{"hydrated", "lookup", "repair_mismatched_relation", "unbound"} {
		t.Run(scenario, func(t *testing.T) {
			account := openAIAccountProxyBindingAccount()
			proxy := openAIAccountProxyBindingProxy()
			repo := &openAIAccountProxyBindingRepo{proxy: proxy}
			wantURL := proxy.URL()
			wantLookups := 1
			switch scenario {
			case "hydrated":
				account.Proxy = proxy
				wantLookups = 0
			case "repair_mismatched_relation":
				stale := *proxy
				stale.ID = 43
				stale.Host = "wrong.example"
				account.Proxy = &stale
			case "unbound":
				account.ProxyID = nil
				account.Proxy = proxy
				wantURL = ""
				wantLookups = 0
			}
			client := &openaiOAuthClientRefreshStub{resp: &openai.TokenResponse{AccessToken: "new-binding-token", ExpiresIn: 3600}}
			service := NewOpenAIOAuthService(repo, client)
			defer service.Stop()
			service.SetPrivacyClientFactory(func(got string) (*req.Client, error) {
				require.Equal(t, wantURL, got)
				return nil, errors.New("stop enrichment transport")
			})
			info, err := service.RefreshAccountToken(context.Background(), account)
			require.NoError(t, err)
			require.Equal(t, "new-binding-token", info.AccessToken)
			require.Equal(t, wantURL, client.lastProxyURL)
			require.Equal(t, wantLookups, repo.calls)
		})
	}
}

func TestOpenAIAccountProxyBinding_QuotaUsesResolvedParentProxy(t *testing.T) {
	for _, scenario := range []string{"hydrated", "lookup", "shadow_parent", "unbound"} {
		t.Run(scenario, func(t *testing.T) {
			account := openAIAccountProxyBindingAccount()
			proxy := openAIAccountProxyBindingProxy()
			proxyRepo := &openAIAccountProxyBindingRepo{proxy: proxy}
			repo := &stubQuotaAccountRepo{accounts: map[int64]*Account{account.ID: account}}
			requestedID := account.ID
			wantURL := proxy.URL()
			wantLookups := 1
			switch scenario {
			case "hydrated":
				account.Proxy = proxy
				wantLookups = 0
			case "shadow_parent":
				account.Proxy = proxy
				wantLookups = 0
				parentID := account.ID
				wrongProxy := *proxy
				wrongProxy.ID = 99
				wrongProxy.Host = "shadow-wrong.example"
				shadow := &Account{ID: 72, Platform: PlatformOpenAI, Type: AccountTypeOAuth, ParentAccountID: &parentID, QuotaDimension: QuotaDimensionSpark, ProxyID: &wrongProxy.ID, Proxy: &wrongProxy}
				repo.accounts[shadow.ID] = shadow
				requestedID = shadow.ID
			case "unbound":
				account.ProxyID = nil
				account.Proxy = proxy
				wantURL = ""
				wantLookups = 0
			}
			cache := &stubQuotaTokenCache{tokens: map[string]string{OpenAITokenCacheKey(account): "cached-binding-token"}}
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				require.Equal(t, "Bearer cached-binding-token", r.Header.Get("authorization"))
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(`{"plan_type":"plus","rate_limit_reset_credits":{"available_count":0}}`))
			}))
			defer server.Close()
			factory := newQuotaRedirectingFactory(server)
			service := NewOpenAIQuotaService(repo, proxyRepo, NewOpenAITokenProvider(repo, cache, nil), func(got string) (*req.Client, error) { require.Equal(t, wantURL, got); return factory(got) })
			usage, err := service.QueryUsage(context.Background(), requestedID)
			require.NoError(t, err)
			require.Equal(t, "plus", usage.PlanType)
			require.Positive(t, requests)
			require.Equal(t, wantLookups, proxyRepo.calls)
		})
	}
}

func TestResolveOpenAIAccountProxyURL_RejectsInvalidBoundConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Account, *openAIAccountProxyBindingRepo)
	}{
		{"zero_binding", func(a *Account, r *openAIAccountProxyBindingRepo) { *a.ProxyID = 0 }},
		{"negative_binding", func(a *Account, r *openAIAccountProxyBindingRepo) { *a.ProxyID = -4 }},
		{"lookup_error", func(a *Account, r *openAIAccountProxyBindingRepo) {
			a.Proxy = nil
			r.err = errors.New("postgres binding-user binding-secret")
		}},
		{"missing_result", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy = nil; r.proxy = nil }},
		{"wrong_result_id", func(a *Account, r *openAIAccountProxyBindingRepo) {
			a.Proxy = nil
			p := *r.proxy
			p.ID = 43
			r.proxy = &p
		}},
		{"empty_scheme", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Protocol = "" }},
		{"unknown_scheme", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Protocol = "ftp" }},
		{"empty_host", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Host = "" }},
		{"host_whitespace", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Host = " proxy.example" }},
		{"host_control", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Host = "proxy.example\r\nother" }},
		{"host_embedded_port", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Host = "proxy.example:9999" }},
		{"host_url", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Host = "http://proxy.example" }},
		{"host_path", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Host = "proxy.example/path" }},
		{"host_credentials", func(a *Account, r *openAIAccountProxyBindingRepo) {
			a.Proxy.Host = "binding-user:binding-secret@proxy.example"
		}},
		{"host_escape", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Host = "proxy%ZZ.example" }},
		{"zero_port", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Port = 0 }},
		{"negative_port", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Port = -1 }},
		{"oversized_port", func(a *Account, r *openAIAccountProxyBindingRepo) { a.Proxy.Port = 65536 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := openAIAccountProxyBindingAccount()
			account.Proxy = openAIAccountProxyBindingProxy()
			repo := &openAIAccountProxyBindingRepo{proxy: account.Proxy}
			tt.mutate(account, repo)
			got, err := resolveOpenAIAccountProxyURL(context.Background(), account, repo)
			require.Error(t, err)
			require.Empty(t, got)
			require.NotContains(t, err.Error(), "binding-user")
			require.NotContains(t, err.Error(), "binding-secret")
			require.NotContains(t, err.Error(), "proxy.example")
		})
	}
}

func TestResolveOpenAIAccountProxyURL_PreservesSupportedBindings(t *testing.T) {
	for _, protocol := range []string{"http", "https", "socks5", "socks5h"} {
		for _, host := range []string{"proxy.example", "127.0.0.1", "2001:db8::1", "fe80::1%eth0"} {
			t.Run(protocol+"/"+host, func(t *testing.T) {
				account := openAIAccountProxyBindingAccount()
				account.Proxy = openAIAccountProxyBindingProxy()
				account.Proxy.Protocol = protocol
				account.Proxy.Host = host
				// Hydrated bindings do not depend on repository availability.
				repo := &openAIAccountProxyBindingRepo{err: errors.New("database unavailable")}
				got, err := resolveOpenAIAccountProxyURL(context.Background(), account, repo)
				require.NoError(t, err)
				require.Zero(t, repo.calls)
				want := account.Proxy.URL()
				if protocol == "socks5" {
					want = "socks5h" + want[len("socks5"):]
				}
				require.Equal(t, want, got)
			})
		}
	}
}

func TestOpenAIAccountProxyBinding_QuotaRejectsMissingParentProxy(t *testing.T) {
	parent := openAIAccountProxyBindingAccount()
	parentID := parent.ID
	shadowProxy := openAIAccountProxyBindingProxy()
	shadowProxy.ID = 99
	shadow := &Account{ID: 72, Platform: PlatformOpenAI, Type: AccountTypeOAuth, ParentAccountID: &parentID, QuotaDimension: QuotaDimensionSpark, ProxyID: &shadowProxy.ID, Proxy: shadowProxy}
	repo := &stubQuotaAccountRepo{accounts: map[int64]*Account{parent.ID: parent, shadow.ID: shadow}}
	cache := &openAIAccountProxyBindingTokenCache{stubQuotaTokenCache: &stubQuotaTokenCache{tokens: map[string]string{OpenAITokenCacheKey(parent): "cached-binding-token"}}}
	factoryCalls := 0
	service := NewOpenAIQuotaService(repo, nil, NewOpenAITokenProvider(repo, cache, nil), func(string) (*req.Client, error) { factoryCalls++; return nil, errors.New("blocked test transport") })
	_, err := service.QueryUsage(context.Background(), shadow.ID)
	require.Error(t, err)
	require.Zero(t, cache.calls)
	require.Zero(t, factoryCalls)
}
