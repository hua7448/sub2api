package repository

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/service"
	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/bandwidth"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/proxy"
)

type OpenAIOAuthServiceSuite struct {
	suite.Suite
	ctx      context.Context
	srv      *httptest.Server
	svc      *openaiOAuthService
	received chan url.Values
}

type tlsClientStub struct {
	do func(req *fhttp.Request) (*fhttp.Response, error)
}

func (s *tlsClientStub) GetCookies(u *url.URL) []*fhttp.Cookie             { return nil }
func (s *tlsClientStub) SetCookies(u *url.URL, cookies []*fhttp.Cookie)    {}
func (s *tlsClientStub) SetCookieJar(jar fhttp.CookieJar)                  {}
func (s *tlsClientStub) GetCookieJar() fhttp.CookieJar                     { return nil }
func (s *tlsClientStub) SetProxy(proxyUrl string) error                    { return nil }
func (s *tlsClientStub) GetProxy() string                                  { return "" }
func (s *tlsClientStub) SetFollowRedirect(followRedirect bool)             {}
func (s *tlsClientStub) GetFollowRedirect() bool                           { return false }
func (s *tlsClientStub) CloseIdleConnections()                             {}
func (s *tlsClientStub) Do(req *fhttp.Request) (*fhttp.Response, error)    { return s.do(req) }
func (s *tlsClientStub) Get(url string) (resp *fhttp.Response, err error)  { return nil, nil }
func (s *tlsClientStub) Head(url string) (resp *fhttp.Response, err error) { return nil, nil }
func (s *tlsClientStub) Post(url, contentType string, body io.Reader) (resp *fhttp.Response, err error) {
	return nil, nil
}
func (s *tlsClientStub) GetBandwidthTracker() bandwidth.BandwidthTracker { return nil }
func (s *tlsClientStub) GetDialer() proxy.ContextDialer                  { return nil }
func (s *tlsClientStub) GetTLSDialer() tls_client.TLSDialerFunc          { return nil }

func (s *OpenAIOAuthServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.received = make(chan url.Values, 1)
}

func (s *OpenAIOAuthServiceSuite) TearDownTest() {
	if s.srv != nil {
		s.srv.Close()
		s.srv = nil
	}
}

func (s *OpenAIOAuthServiceSuite) setupServer(handler http.HandlerFunc) {
	s.srv = newLocalTestServer(s.T(), handler)
	s.svc = &openaiOAuthService{tokenURL: s.srv.URL}
}

func (s *OpenAIOAuthServiceSuite) TestExchangeCode_DefaultRedirectURI() {
	errCh := make(chan string, 1)
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			errCh <- "method mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := r.ParseForm(); err != nil {
			errCh <- "ParseForm failed"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.PostForm.Get("grant_type"); got != "authorization_code" {
			errCh <- "grant_type mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.PostForm.Get("client_id"); got != openai.ClientID {
			errCh <- "client_id mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.PostForm.Get("code"); got != "code" {
			errCh <- "code mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.PostForm.Get("redirect_uri"); got != openai.DefaultRedirectURI {
			errCh <- "redirect_uri mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.PostForm.Get("code_verifier"); got != "ver" {
			errCh <- "code_verifier mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		wantUA, wantOriginator := service.CodexCanonicalAuthIdentity()
		if got := r.Header.Get("User-Agent"); got != wantUA {
			errCh <- "user-agent mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.Header.Get("originator"); got != wantOriginator {
			errCh <- "originator mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"at","refresh_token":"rt","token_type":"bearer","expires_in":3600}`)
	}))

	resp, err := s.svc.ExchangeCode(s.ctx, "code", "ver", "", "", "")
	require.NoError(s.T(), err, "ExchangeCode")
	select {
	case msg := <-errCh:
		require.Fail(s.T(), msg)
	default:
	}
	require.Equal(s.T(), "at", resp.AccessToken)
	require.Equal(s.T(), "rt", resp.RefreshToken)
}

func (s *OpenAIOAuthServiceSuite) TestExchangeCode_UsesOfficialTokenEndpointOnly() {
	var officialCalls int
	official := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		officialCalls++
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		require.Equal(s.T(), "authorization_code", r.PostForm.Get("grant_type"))
		require.Equal(s.T(), "custom-client", r.PostForm.Get("client_id"))
		require.Equal(s.T(), "code", r.PostForm.Get("code"))
		require.Equal(s.T(), "http://localhost:9999/cb", r.PostForm.Get("redirect_uri"))
		require.Equal(s.T(), "ver", r.PostForm.Get("code_verifier"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"official-at","refresh_token":"official-rt","token_type":"bearer","expires_in":3600}`)
	}))
	defer official.Close()

	var legacyCalls int
	legacy := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		legacyCalls++
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer legacy.Close()

	s.svc = &openaiOAuthService{
		tokenURL:       official.URL,
		legacyTokenURL: legacy.URL,
	}

	resp, err := s.svc.ExchangeCode(s.ctx, "code", "ver", "http://localhost:9999/cb", "", "custom-client")
	require.NoError(s.T(), err, "ExchangeCode")
	require.Equal(s.T(), "official-at", resp.AccessToken)
	require.Equal(s.T(), "official-rt", resp.RefreshToken)
	require.Equal(s.T(), 1, officialCalls)
	require.Equal(s.T(), 0, legacyCalls)
}

func (s *OpenAIOAuthServiceSuite) TestExchangeCode_DoesNotUseLegacyWhenOfficialFails() {
	var legacyCalls int
	legacy := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		legacyCalls++
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"legacy-at","refresh_token":"legacy-rt","token_type":"bearer","expires_in":3600}`)
	}))
	defer legacy.Close()

	var officialCalls int
	official := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		officialCalls++
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":{"code":"token_exchange_user_error"}}`)
	}))
	defer official.Close()

	s.svc = &openaiOAuthService{
		tokenURL:       official.URL,
		legacyTokenURL: legacy.URL,
	}

	_, err := s.svc.ExchangeCode(s.ctx, "code", "ver", openai.DefaultRedirectURI, "", "")
	require.Error(s.T(), err, "ExchangeCode")
	require.ErrorContains(s.T(), err, "token exchange failed: status 400")
	require.Equal(s.T(), 0, legacyCalls)
	require.Equal(s.T(), 1, officialCalls)
}

func (s *OpenAIOAuthServiceSuite) TestRefreshToken_FormFields() {
	errCh := make(chan string, 1)
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			errCh <- "ParseForm failed"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.PostForm.Get("grant_type"); got != "refresh_token" {
			errCh <- "grant_type mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.PostForm.Get("refresh_token"); got != "rt" {
			errCh <- "refresh_token mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.PostForm.Get("client_id"); got != openai.ClientID {
			errCh <- "client_id mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.PostForm.Get("scope"); got != openai.RefreshScopes {
			errCh <- "scope mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		wantUA, wantOriginator := service.CodexCanonicalAuthIdentity()
		if got := r.Header.Get("User-Agent"); got != wantUA {
			errCh <- "user-agent mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if got := r.Header.Get("originator"); got != wantOriginator {
			errCh <- "originator mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"at2","refresh_token":"rt2","token_type":"bearer","expires_in":3600}`)
	}))

	resp, err := s.svc.RefreshToken(s.ctx, "rt", "")
	require.NoError(s.T(), err, "RefreshToken")
	select {
	case msg := <-errCh:
		require.Fail(s.T(), msg)
	default:
	}
	require.Equal(s.T(), "at2", resp.AccessToken)
	require.Equal(s.T(), "rt2", resp.RefreshToken)
}

// TestRefreshToken_DefaultsToOpenAIClientID 验证未指定 client_id 时默认使用 OpenAI ClientID，
// 且只发送一次请求（不再盲猜多个 client_id）。
func (s *OpenAIOAuthServiceSuite) TestRefreshToken_DefaultsToOpenAIClientID() {
	var seenClientIDs []string
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		clientID := r.PostForm.Get("client_id")
		seenClientIDs = append(seenClientIDs, clientID)
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"at","refresh_token":"rt","token_type":"bearer","expires_in":3600}`)
	}))

	resp, err := s.svc.RefreshToken(s.ctx, "rt", "")
	require.NoError(s.T(), err, "RefreshToken")
	require.Equal(s.T(), "at", resp.AccessToken)
	// 只发送了一次请求，使用默认的 OpenAI ClientID
	require.Equal(s.T(), []string{openai.ClientID}, seenClientIDs)
}

func (s *OpenAIOAuthServiceSuite) TestRefreshToken_UseProvidedClientID() {
	const customClientID = "custom-client-id"
	var seenClientIDs []string
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		clientID := r.PostForm.Get("client_id")
		seenClientIDs = append(seenClientIDs, clientID)
		if clientID != customClientID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"at-custom","refresh_token":"rt-custom","token_type":"bearer","expires_in":3600}`)
	}))

	resp, err := s.svc.RefreshTokenWithClientID(s.ctx, "rt", "", customClientID)
	require.NoError(s.T(), err, "RefreshTokenWithClientID")
	require.Equal(s.T(), "at-custom", resp.AccessToken)
	require.Equal(s.T(), "rt-custom", resp.RefreshToken)
	require.Equal(s.T(), []string{customClientID}, seenClientIDs)
}

func (s *OpenAIOAuthServiceSuite) TestRefreshToken_AccessTokenInputPassesThroughWithoutTokenExchange() {
	var officialCalls int
	official := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		officialCalls++
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer official.Close()

	var legacyCalls int
	legacy := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		legacyCalls++
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer legacy.Close()

	s.svc = &openaiOAuthService{
		tokenURL:       official.URL,
		legacyTokenURL: legacy.URL,
	}

	resp, err := s.svc.RefreshTokenWithClientID(s.ctx, " eyJhbGciOiJub25lIn0.payload.sig ", "", openai.ClientID)
	require.NoError(s.T(), err, "RefreshTokenWithClientID")
	require.Equal(s.T(), "eyJhbGciOiJub25lIn0.payload.sig", resp.AccessToken)
	require.Empty(s.T(), resp.RefreshToken)
	require.True(s.T(), resp.PassThrough)
	require.Equal(s.T(), 0, officialCalls)
	require.Equal(s.T(), 0, legacyCalls)

	resp, err = s.svc.RefreshTokenWithClientID(s.ctx, "at-project-token", "", openai.ClientID)
	require.NoError(s.T(), err, "RefreshTokenWithClientID")
	require.Equal(s.T(), "at-project-token", resp.AccessToken)
	require.True(s.T(), resp.PassThrough)
	require.Equal(s.T(), 0, officialCalls)
	require.Equal(s.T(), 0, legacyCalls)
}

func (s *OpenAIOAuthServiceSuite) TestRefreshToken_DefaultCallsOfficialOnlyForAllPrefixes() {
	for _, token := range []string{"RT-project-one", "rt-official", "xy-legacy", "XY-project", "ordinary-refresh-token"} {
		s.Run(token, func() {
			var officialCalls, legacyCalls int
			official := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				officialCalls++
				require.NoError(s.T(), r.ParseForm())
				require.Equal(s.T(), token, r.Form.Get("refresh_token"))
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{"access_token":"official-at","refresh_token":"official-rt","expires_in":3600}`)
			}))
			defer official.Close()
			legacy := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				legacyCalls++
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, `{"access_token":"legacy-at"}`)
			}))
			defer legacy.Close()
			svc := &openaiOAuthService{tokenURL: official.URL, legacyTokenURL: legacy.URL}
			resp, err := svc.RefreshTokenWithClientID(s.ctx, token, "", openai.ClientID)
			require.NoError(s.T(), err)
			require.Equal(s.T(), "official-at", resp.AccessToken)
			require.Equal(s.T(), 1, officialCalls)
			require.Zero(s.T(), legacyCalls, "token prefixes must not disclose credentials to another provider")
		})
	}
}

func (s *OpenAIOAuthServiceSuite) TestRefreshToken_OfficialFailureNeverFallsBackToLegacy() {
	var officialCalls, legacyCalls int
	official := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		officialCalls++
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":"invalid_grant","echo":"RT-sensitive-token"}`)
	}))
	defer official.Close()
	legacy := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		legacyCalls++
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"legacy-at"}`)
	}))
	defer legacy.Close()
	svc := &openaiOAuthService{tokenURL: official.URL, legacyTokenURL: legacy.URL}
	resp, err := svc.RefreshTokenWithClientID(s.ctx, "RT-sensitive-token", "", openai.ClientID)
	require.ErrorContains(s.T(), err, "status 400")
	require.NotContains(s.T(), err.Error(), "RT-sensitive-token")
	require.Nil(s.T(), resp)
	require.Equal(s.T(), 1, officialCalls)
	require.Zero(s.T(), legacyCalls)
}

func (s *OpenAIOAuthServiceSuite) TestFetchAccountInfo_FallsBackToTLSClientOnCloudflareChallenge() {
	challengeBody := `<!DOCTYPE html><html><head><title>Just a moment...</title></head><body><script>window._cf_chl_opt={};</script></body></html>`
	primary := newLocalTestServer(s.T(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, challengeBody)
	}))
	defer primary.Close()

	s.svc = &openaiOAuthService{
		accountCheckURL: primary.URL,
		tlsClientFactory: func(proxyURL string) (tls_client.HttpClient, error) {
			return &tlsClientStub{
				do: func(req *fhttp.Request) (*fhttp.Response, error) {
					return &fhttp.Response{
						StatusCode: http.StatusOK,
						Header:     fhttp.Header{"Content-Type": []string{"application/json"}},
						Body: io.NopCloser(strings.NewReader(
							`{"accounts":{"default":{"account":{"account_id":"default-account","plan_type":"plus"}},"team":{"account":{"account_id":"team-account","plan_type":"team"}}}}`,
						)),
					}, nil
				},
			}, nil
		},
	}

	info, err := s.svc.FetchAccountInfo(s.ctx, "access-token", "")
	require.NoError(s.T(), err)
	require.Len(s.T(), info.Accounts, 1)
	require.Equal(s.T(), "team-account", info.Accounts[0].AccountID)
	require.Equal(s.T(), "team", info.Accounts[0].PlanType)
}

func (s *OpenAIOAuthServiceSuite) TestNonSuccessStatus_IncludesBody() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "bad")
	}))

	_, err := s.svc.ExchangeCode(s.ctx, "code", "ver", openai.DefaultRedirectURI, "", "")
	require.Error(s.T(), err)
	require.ErrorContains(s.T(), err, "status 400")
	require.ErrorContains(s.T(), err, "bad")
}

func (s *OpenAIOAuthServiceSuite) TestRequestError_ClosedServer() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	s.srv.Close()

	_, err := s.svc.ExchangeCode(s.ctx, "code", "ver", openai.DefaultRedirectURI, "", "")
	require.Error(s.T(), err)
	require.ErrorContains(s.T(), err, "request failed")
}

func (s *OpenAIOAuthServiceSuite) TestExchangeCode_RequestErrorWithoutProxyReturnsProxyHint() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	s.srv.Close()

	_, err := s.svc.ExchangeCode(s.ctx, "code", "ver", openai.DefaultRedirectURI, "", "")

	require.Error(s.T(), err)
	require.Equal(s.T(), "OPENAI_OAUTH_PROXY_REQUIRED", infraerrors.Reason(err))
	require.Contains(s.T(), infraerrors.Message(err), "no proxy is configured")
}

func (s *OpenAIOAuthServiceSuite) TestContextCancel() {
	started := make(chan struct{})
	block := make(chan struct{})
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-block
	}))

	ctx, cancel := context.WithCancel(s.ctx)

	done := make(chan error, 1)
	go func() {
		_, err := s.svc.ExchangeCode(ctx, "code", "ver", openai.DefaultRedirectURI, "", "")
		done <- err
	}()

	<-started
	cancel()
	close(block)

	err := <-done
	require.Error(s.T(), err)
}

func (s *OpenAIOAuthServiceSuite) TestExchangeCode_UsesProvidedRedirectURI() {
	want := "http://localhost:9999/cb"
	errCh := make(chan string, 1)
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if got := r.PostForm.Get("redirect_uri"); got != want {
			errCh <- "redirect_uri mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"at","token_type":"bearer","expires_in":1}`)
	}))

	_, err := s.svc.ExchangeCode(s.ctx, "code", "ver", want, "", "")
	require.NoError(s.T(), err, "ExchangeCode")
	select {
	case msg := <-errCh:
		require.Fail(s.T(), msg)
	default:
	}
}

func (s *OpenAIOAuthServiceSuite) TestExchangeCode_UseProvidedClientID() {
	wantClientID := "custom-exchange-client-id"
	errCh := make(chan string, 1)
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if got := r.PostForm.Get("client_id"); got != wantClientID {
			errCh <- "client_id mismatch"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"at","token_type":"bearer","expires_in":1}`)
	}))

	_, err := s.svc.ExchangeCode(s.ctx, "code", "ver", openai.DefaultRedirectURI, "", wantClientID)
	require.NoError(s.T(), err, "ExchangeCode")
	select {
	case msg := <-errCh:
		require.Fail(s.T(), msg)
	default:
	}
}

func (s *OpenAIOAuthServiceSuite) TestTokenURL_CanBeOverriddenWithQuery() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		s.received <- r.PostForm
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"access_token":"at","token_type":"bearer","expires_in":1}`)
	}))
	s.svc.tokenURL = s.srv.URL + "?x=1"

	_, err := s.svc.ExchangeCode(s.ctx, "code", "ver", openai.DefaultRedirectURI, "", "")
	require.NoError(s.T(), err, "ExchangeCode")
	select {
	case <-s.received:
	default:
		require.Fail(s.T(), "expected server to receive request")
	}
}

func (s *OpenAIOAuthServiceSuite) TestExchangeCode_SuccessButInvalidJSON() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "not-valid-json")
	}))

	_, err := s.svc.ExchangeCode(s.ctx, "code", "ver", openai.DefaultRedirectURI, "", "")
	require.Error(s.T(), err, "expected error for invalid JSON response")
}

func (s *OpenAIOAuthServiceSuite) TestRefreshToken_NonSuccessStatus() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, "unauthorized")
	}))

	_, err := s.svc.RefreshToken(s.ctx, "rt", "")
	require.Error(s.T(), err, "expected error for non-2xx status")
	require.ErrorContains(s.T(), err, "status 401")
}

func TestNewOpenAIOAuthClient_DefaultTokenURL(t *testing.T) {
	client := NewOpenAIOAuthClient()
	svc, ok := client.(*openaiOAuthService)
	require.True(t, ok)
	require.Equal(t, openai.TokenURL, svc.tokenURL)
}

func TestOpenAIOAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(OpenAIOAuthServiceSuite))
}
