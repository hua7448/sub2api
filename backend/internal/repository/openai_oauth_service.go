package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/Wei-Shaw/sub2api/internal/util/httputil"
	fhttp "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
)

const legacyOpenAITokenURL = "https://public.xyhelper.cn/oauth/token"
const openAIAccountCheckURL = "https://chatgpt.com/backend-api/accounts/check/v4-2023-04-27"

// NewOpenAIOAuthClient creates a new OpenAI OAuth client
func NewOpenAIOAuthClient() service.OpenAIOAuthClient {
	return &openaiOAuthService{
		tokenURL:        openai.TokenURL,
		legacyTokenURL:  legacyOpenAITokenURL,
		accountCheckURL: openAIAccountCheckURL,
	}
}

type openaiOAuthService struct {
	tokenURL         string
	legacyTokenURL   string
	accountCheckURL  string
	tlsClientFactory func(proxyURL string) (tls_client.HttpClient, error)
}

func (s *openaiOAuthService) ExchangeCode(ctx context.Context, code, codeVerifier, redirectURI, proxyURL, clientID string) (*openai.TokenResponse, error) {
	if redirectURI == "" {
		redirectURI = openai.DefaultRedirectURI
	}
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		clientID = openai.ClientID
	}

	return s.exchangeCodeWithClientID(ctx, code, codeVerifier, redirectURI, proxyURL, clientID)
}

func (s *openaiOAuthService) exchangeCodeWithClientID(ctx context.Context, code, codeVerifier, redirectURI, proxyURL, clientID string) (*openai.TokenResponse, error) {
	client, err := createOpenAIReqClient(proxyURL)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_CLIENT_INIT_FAILED", "create HTTP client: %v", err)
	}

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", clientID)
	formData.Set("code", code)
	formData.Set("redirect_uri", redirectURI)
	formData.Set("code_verifier", codeVerifier)

	var tokenResp openai.TokenResponse

	authUA, authOriginator := service.CodexCanonicalAuthIdentity()
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("User-Agent", authUA).
		SetHeader("originator", authOriginator).
		SetFormDataFromValues(formData).
		SetSuccessResult(&tokenResp).
		Post(s.tokenURL)

	if err != nil {
		if shouldReturnOpenAINoProxyHint(ctx, proxyURL, err) {
			return nil, newOpenAINoProxyHintError(err)
		}
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_REQUEST_FAILED", "request failed: %v", err)
	}

	if !resp.IsSuccessState() {
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_TOKEN_EXCHANGE_FAILED", "token exchange failed: status %d, body: %s", resp.StatusCode, resp.String())
	}

	return &tokenResp, nil
}

func (s *openaiOAuthService) RefreshToken(ctx context.Context, refreshToken, proxyURL string) (*openai.TokenResponse, error) {
	return s.RefreshTokenWithClientID(ctx, refreshToken, proxyURL, "")
}

func (s *openaiOAuthService) RefreshTokenWithClientID(ctx context.Context, refreshToken, proxyURL string, clientID string) (*openai.TokenResponse, error) {
	if isOpenAIAccessTokenInput(refreshToken) {
		return &openai.TokenResponse{
			AccessToken: strings.TrimSpace(refreshToken),
			TokenType:   "Bearer",
			PassThrough: true,
		}, nil
	}

	// 调用方应始终传入正确的 client_id；为兼容旧数据，未指定时默认使用 OpenAI ClientID
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		clientID = openai.ClientID
	}
	// Refresh tokens are sent only to the official endpoint. Token prefixes do
	// not authorize sharing account credentials with a third-party service.
	return s.refreshTokenWithClientID(ctx, refreshToken, proxyURL, clientID)
}

func (s *openaiOAuthService) refreshTokenWithClientID(ctx context.Context, refreshToken, proxyURL, clientID string) (*openai.TokenResponse, error) {
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", refreshToken)
	formData.Set("client_id", clientID)
	formData.Set("scope", openai.RefreshScopes)

	authUA, authOriginator := service.CodexAuthIdentityFromContext(ctx)
	headers := http.Header{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"User-Agent":   {authUA},
		"Originator":   {authOriginator},
	}
	return requestOpenAIRefreshToken(ctx, s.tokenURL, proxyURL, refreshToken, strings.NewReader(formData.Encode()), headers)
}

// requestOpenAIRefreshToken shares the configured transport, while keeping the
// redirect policy local to this credential-bearing request. In particular a
// 307/308 must not replay the refresh token to a different endpoint.
func requestOpenAIRefreshToken(ctx context.Context, endpoint, proxyURL, refreshToken string, body io.Reader, headers http.Header) (*openai.TokenResponse, error) {
	shared, err := createOpenAIReqClient(proxyURL)
	if err != nil {
		return nil, infraerrors.New(http.StatusBadGateway, "OPENAI_OAUTH_CLIENT_INIT_FAILED", "create refresh HTTP client failed")
	}
	client := *shared.GetClient()
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, infraerrors.New(http.StatusBadGateway, "OPENAI_OAUTH_REQUEST_FAILED", "create token refresh request failed")
	}
	request.Header = headers
	resp, err := client.Do(request)
	if err != nil {
		// Preserve useful connection errors without echoing a token if it was
		// accidentally included in a request URL or a transport error.
		message := err.Error()
		if refreshToken != "" {
			message = strings.ReplaceAll(message, refreshToken, "[redacted]")
			message = strings.ReplaceAll(message, url.QueryEscape(refreshToken), "[redacted]")
		}
		safeErr := errors.New(message)
		if shouldReturnOpenAINoProxyHint(ctx, proxyURL, err) {
			return nil, newOpenAINoProxyHintError(safeErr)
		}
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_REQUEST_FAILED", "request failed: %v", safeErr)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Preserve only known machine-readable codes: invalid_grant drives
		// cross-worker race recovery and other codes stop futile retries.
		// Provider messages and unknown values may echo submitted credentials.
		if code := safeOpenAIRefreshErrorCode(resp.Body); code != "" {
			return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_TOKEN_REFRESH_FAILED", "token refresh failed: status %d, code: %s", resp.StatusCode, code)
		}
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_TOKEN_REFRESH_FAILED", "token refresh failed: status %d", resp.StatusCode)
	}
	var tokenResp openai.TokenResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&tokenResp); err != nil {
		return nil, infraerrors.New(http.StatusBadGateway, "OPENAI_OAUTH_TOKEN_REFRESH_FAILED", "invalid token refresh response")
	}
	if strings.TrimSpace(tokenResp.AccessToken) == "" {
		return nil, infraerrors.New(http.StatusBadGateway, "OPENAI_OAUTH_TOKEN_REFRESH_FAILED", "token refresh response missing access_token")
	}
	return &tokenResp, nil
}

// safeOpenAIRefreshErrorCode never returns provider-controlled descriptions or
// arbitrary error strings. Only fixed codes consumed by the refresh recovery,
// retry, and provider-configuration classifiers may leave this boundary.
func safeOpenAIRefreshErrorCode(body io.Reader) string {
	const maxErrorBodyBytes = 64 << 10
	encoded, err := io.ReadAll(io.LimitReader(body, maxErrorBodyBytes+1))
	if err != nil || len(encoded) > maxErrorBodyBytes {
		return ""
	}
	var envelope struct {
		Error json.RawMessage `json:"error"`
	}
	if json.Unmarshal(encoded, &envelope) != nil {
		return ""
	}
	var value string
	if json.Unmarshal(envelope.Error, &value) == nil {
		return allowlistedOpenAIRefreshErrorCode(value)
	}
	var detail struct {
		Code string `json:"code"`
		Type string `json:"type"`
	}
	if json.Unmarshal(envelope.Error, &detail) != nil {
		return ""
	}
	if code := allowlistedOpenAIRefreshErrorCode(detail.Code); code != "" {
		return code
	}
	return allowlistedOpenAIRefreshErrorCode(detail.Type)
}

func allowlistedOpenAIRefreshErrorCode(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "invalid_grant":
		return "invalid_grant"
	case "invalid_refresh_token":
		return "invalid_refresh_token"
	case "token_expired":
		return "token_expired"
	case "app_session_terminated":
		return "app_session_terminated"
	case "refresh_token_reused":
		return "refresh_token_reused"
	case "refresh_token_invalidated":
		return "refresh_token_invalidated"
	case "invalid_client":
		return "invalid_client"
	case "unauthorized_client":
		return "unauthorized_client"
	case "access_denied":
		return "access_denied"
	case "missing_project_id":
		return "missing_project_id"
	case "invalid_scope", "unknown_scope", "unknown scope":
		return "invalid_scope"
	case "entitlement_denied":
		return "entitlement_denied"
	default:
		return ""
	}
}

// refreshTokenWithLegacyService retains the historical adapter for review and
// compatibility work. Automatic refresh intentionally never calls this adapter.
func (s *openaiOAuthService) refreshTokenWithLegacyService(ctx context.Context, refreshToken, proxyURL, clientID string) (*openai.TokenResponse, error) {
	payload := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": strings.TrimSpace(refreshToken),
	}
	if trimmedClientID := strings.TrimSpace(clientID); trimmedClientID != "" {
		payload["client_id"] = trimmedClientID
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, infraerrors.New(http.StatusBadGateway, "OPENAI_OAUTH_REQUEST_FAILED", "create token refresh request failed")
	}
	headers := http.Header{
		"Content-Type": {"application/json"},
		"Accept":       {"*/*"},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
	}
	return requestOpenAIRefreshToken(ctx, s.legacyTokenURL, proxyURL, refreshToken, bytes.NewReader(body), headers)
}

func isOpenAIAccessTokenInput(token string) bool {
	trimmed := strings.TrimSpace(token)
	lower := strings.ToLower(trimmed)
	return strings.HasPrefix(trimmed, "eyJ") ||
		strings.HasPrefix(lower, "at") ||
		strings.HasPrefix(lower, "sess-")
}

func (s *openaiOAuthService) FetchAccountInfo(ctx context.Context, accessToken, proxyURL string) (*service.OpenAIAccountInfoResult, error) {
	reqResult, reqErr := s.fetchAccountInfoWithReq(ctx, accessToken, proxyURL)
	if reqErr == nil && reqResult != nil && len(reqResult.Accounts) > 0 {
		return reqResult, nil
	}

	tlsResult, tlsErr := s.fetchAccountInfoWithTLSClient(ctx, accessToken, proxyURL)
	if tlsErr == nil && tlsResult != nil && len(tlsResult.Accounts) > 0 {
		return tlsResult, nil
	}

	if reqErr == nil && reqResult != nil {
		return reqResult, nil
	}
	if tlsErr == nil && tlsResult != nil {
		return tlsResult, nil
	}
	if reqErr != nil && tlsErr != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_ACCOUNT_INFO_FAILED", "req client failed: %v; tls client failed: %v", reqErr, tlsErr)
	}
	if reqErr != nil {
		return nil, reqErr
	}
	return nil, tlsErr
}

func (s *openaiOAuthService) fetchAccountInfoWithReq(ctx context.Context, accessToken, proxyURL string) (*service.OpenAIAccountInfoResult, error) {
	client, err := createOpenAIReqClient(proxyURL)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_CLIENT_INIT_FAILED", "create HTTP client: %v", err)
	}

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+strings.TrimSpace(accessToken)).
		SetHeader("Accept", "*/*").
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Content-Type", "application/json").
		SetHeader("Host", "chatgpt.com").
		SetHeader("Origin", "https://chatgpt.com").
		SetHeader("Referer", "https://chatgpt.com/").
		SetHeader("Oai-Device-Id", uuid.NewString()).
		SetHeader("Oai-Language", "en-US").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		Get(s.accountCheckURL)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_ACCOUNT_INFO_REQUEST_FAILED", "request failed: %v", err)
	}

	return parseOpenAIAccountInfoResponse(resp.StatusCode, resp.Header, []byte(resp.String()))
}

func (s *openaiOAuthService) fetchAccountInfoWithTLSClient(ctx context.Context, accessToken, proxyURL string) (*service.OpenAIAccountInfoResult, error) {
	clientFactory := s.tlsClientFactory
	if clientFactory == nil {
		clientFactory = newOpenAIAccountCheckTLSClient
	}
	client, err := clientFactory(proxyURL)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_ACCOUNT_INFO_TLS_CLIENT_INIT_FAILED", "create tls client: %v", err)
	}

	req, err := fhttp.NewRequestWithContext(ctx, fhttp.MethodGet, s.accountCheckURL, nil)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusInternalServerError, "OPENAI_OAUTH_ACCOUNT_INFO_REQUEST_BUILD_FAILED", "build request failed: %v", err)
	}
	setOpenAIAccountCheckFHeaders(req.Header, accessToken)
	req.Host = "chatgpt.com"

	resp, err := client.Do(req)
	if err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_ACCOUNT_INFO_TLS_REQUEST_FAILED", "request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	return parseOpenAIAccountInfoResponse(resp.StatusCode, cloneFHTTPHeader(resp.Header), body)
}

func createOpenAIReqClient(proxyURL string) (*req.Client, error) {
	return getSharedReqClient(reqClientOptions{
		ProxyURL: proxyURL,
		Timeout:  120 * time.Second,
	})
}

func newOpenAIAccountCheckTLSClient(proxyURL string) (tls_client.HttpClient, error) {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(120),
		tls_client.WithClientProfile(profiles.Chrome_131),
		tls_client.WithRandomTLSExtensionOrder(),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithForceHttp1(),
	}
	if trimmed := strings.TrimSpace(proxyURL); trimmed != "" {
		options = append(options, tls_client.WithProxyUrl(trimmed))
	}
	return tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
}

func setOpenAIAccountCheckHeaders(headers http.Header, accessToken string) {
	headers.Set("Authorization", "Bearer "+strings.TrimSpace(accessToken))
	headers.Set("Accept", "*/*")
	headers.Set("Accept-Language", "en-US,en;q=0.9")
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "https://chatgpt.com")
	headers.Set("Referer", "https://chatgpt.com/")
	headers.Set("Oai-Device-Id", uuid.NewString())
	headers.Set("Oai-Language", "en-US")
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
}

func setOpenAIAccountCheckFHeaders(headers fhttp.Header, accessToken string) {
	headers.Set("Authorization", "Bearer "+strings.TrimSpace(accessToken))
	headers.Set("Accept", "*/*")
	headers.Set("Accept-Language", "en-US,en;q=0.9")
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "https://chatgpt.com")
	headers.Set("Referer", "https://chatgpt.com/")
	headers.Set("Oai-Device-Id", uuid.NewString())
	headers.Set("Oai-Language", "en-US")
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
}

func parseOpenAIAccountInfoResponse(statusCode int, headers http.Header, body []byte) (*service.OpenAIAccountInfoResult, error) {
	if isOpenAICloudflareChallengeResponse(statusCode, headers, body) {
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_ACCOUNT_INFO_CF_CHALLENGE", "account info blocked by Cloudflare challenge: %s", httputil.FormatCloudflareChallengeMessage("challenge detected", headers, body))
	}
	if statusCode != http.StatusOK {
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_ACCOUNT_INFO_FAILED", "account info fetch failed: status %d, body: %s", statusCode, httputil.TruncateBody(body, 512))
	}

	type accountPayload struct {
		Account struct {
			AccountID string `json:"account_id"`
			PlanType  string `json:"plan_type"`
		} `json:"account"`
	}
	var result struct {
		Accounts map[string]accountPayload `json:"accounts"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, infraerrors.Newf(http.StatusBadGateway, "OPENAI_OAUTH_ACCOUNT_INFO_PARSE_FAILED", "parse account info response failed: %v; body: %s", err, httputil.TruncateBody(body, 512))
	}

	accounts := make([]service.OpenAIAccountInfo, 0, len(result.Accounts))
	for key, item := range result.Accounts {
		if strings.EqualFold(strings.TrimSpace(key), "default") {
			continue
		}
		accounts = append(accounts, service.OpenAIAccountInfo{
			AccountID:  strings.TrimSpace(item.Account.AccountID),
			PlanType:   strings.TrimSpace(item.Account.PlanType),
			AccountKey: key,
		})
	}
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].AccountKey < accounts[j].AccountKey
	})

	return &service.OpenAIAccountInfoResult{
		Accounts:   accounts,
		TotalCount: len(accounts),
	}, nil
}

func isOpenAICloudflareChallengeResponse(statusCode int, headers http.Header, body []byte) bool {
	if httputil.IsCloudflareChallengeResponse(statusCode, headers, body) {
		return true
	}

	preview := strings.ToLower(strings.TrimSpace(string(bytes.TrimSpace(body))))
	if preview == "" {
		return false
	}
	if strings.Contains(preview, "window._cf_chl_opt") ||
		strings.Contains(preview, "just a moment") ||
		strings.Contains(preview, "enable javascript and cookies to continue") ||
		strings.Contains(preview, "challenge-platform") {
		return true
	}
	contentType := strings.ToLower(strings.TrimSpace(headers.Get("content-type")))
	if strings.Contains(contentType, "text/html") &&
		(strings.Contains(preview, "<html") || strings.Contains(preview, "<!doctype html")) &&
		strings.Contains(preview, "cloudflare") {
		return true
	}
	return false
}

func cloneFHTTPHeader(headers fhttp.Header) http.Header {
	cloned := make(http.Header, len(headers))
	for k, values := range headers {
		copied := make([]string, len(values))
		copy(copied, values)
		cloned[k] = copied
	}
	return cloned
}
func shouldReturnOpenAINoProxyHint(ctx context.Context, proxyURL string, err error) bool {
	if strings.TrimSpace(proxyURL) != "" || err == nil {
		return false
	}
	if ctx != nil && ctx.Err() != nil {
		return false
	}
	return !errors.Is(err, context.Canceled)
}

func newOpenAINoProxyHintError(cause error) error {
	return infraerrors.New(
		http.StatusBadGateway,
		"OPENAI_OAUTH_PROXY_REQUIRED",
		"OpenAI OAuth request failed: no proxy is configured and this server could not reach OpenAI directly. Select a proxy that can access OpenAI, then retry; if the authorization code has expired, regenerate the authorization URL.",
	).WithCause(cause)
}
