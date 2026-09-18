package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func countCompactionTriggers(body []byte) int {
	count := 0
	input := gjson.GetBytes(body, "input")
	input.ForEach(func(_, item gjson.Result) bool {
		if item.Get("type").String() == "compaction_trigger" {
			count++
		}
		return true
	})
	return count
}

func TestNormalizeOpenAILegacyCompactRemoteV2Body(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantInputLen int
	}{
		{name: "array", body: `{"model":"gpt-5.4","input":[{"type":"message","role":"user","content":"hello"}]}`, wantInputLen: 2},
		{name: "string", body: `{"model":"gpt-5.4","input":"hello"}`, wantInputLen: 2},
		{name: "object", body: `{"model":"gpt-5.4","input":{"type":"message","role":"user","content":"hello"}}`, wantInputLen: 2},
		{name: "existing trigger", body: `{"model":"gpt-5.4","stream":true,"store":false,"input":[{"type":"message","role":"user","content":"hello"},{"type":"compaction_trigger"}]}`, wantInputLen: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized, _, err := normalizeOpenAILegacyCompactRemoteV2Body([]byte(tt.body))
			require.NoError(t, err)
			require.True(t, gjson.GetBytes(normalized, "stream").Bool())
			require.Equal(t, gjson.False, gjson.GetBytes(normalized, "store").Type)
			require.Len(t, gjson.GetBytes(normalized, "input").Array(), tt.wantInputLen)
			require.Equal(t, 1, countCompactionTriggers(normalized))
		})
	}

	_, _, err := normalizeOpenAILegacyCompactRemoteV2Body([]byte(`{"model":"gpt-5.4"}`))
	require.ErrorContains(t, err, "requires input")
}

func TestShouldAdaptOpenAILegacyCompactToRemoteV2IsExactAndOAuthOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	oauth := &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth}
	apiKey := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKey}

	newContext := func(path string) *gin.Context {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, path, nil)
		return c
	}

	require.True(t, shouldAdaptOpenAILegacyCompactToRemoteV2(newContext("/v1/responses/compact"), oauth))
	require.True(t, shouldAdaptOpenAILegacyCompactToRemoteV2(newContext("/responses/compact/"), oauth))
	require.False(t, shouldAdaptOpenAILegacyCompactToRemoteV2(newContext("/v1/responses/compact/detail"), oauth))
	require.False(t, shouldAdaptOpenAILegacyCompactToRemoteV2(newContext("/v1/responses/compact"), apiKey))
}

func TestOpenAIGatewayServiceForwardOAuthLegacyCompactUsesRemoteV2AndReturnsUnaryJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses/compact", bytes.NewReader(nil))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("User-Agent", "codex-tui/0.147.0")
	c.Request.Header.Set("x-codex-beta-features", "existing_feature, remote_compaction_v2")
	SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)

	upstreamSSE := strings.Join([]string{
		`event: response.output_item.done`,
		`data: {"type":"response.output_item.done","output_index":0,"item":{"id":"cmp_legacy","type":"compaction","status":"completed","encrypted_content":"legacy-v2"}}`,
		"",
		`event: response.completed`,
		`data: {"type":"response.completed","response":{"id":"resp_legacy_v2","status":"completed","model":"gpt-5.4","output":[],"usage":{"input_tokens":7,"output_tokens":2,"total_tokens":9}}}`,
		"",
	}, "\n")
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}, "x-request-id": []string{"rid-legacy-v2"}},
		Body:       io.NopCloser(strings.NewReader(upstreamSSE)),
	}}
	cfg := &config.Config{}
	cfg.Security.URLAllowlist.Enabled = false
	svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstream, toolCorrector: NewCodexToolCorrector()}
	account := &Account{
		ID: 2019, Name: "oauth-legacy-compact", Platform: PlatformOpenAI, Type: AccountTypeOAuth, Concurrency: 1,
		Credentials: map[string]any{"access_token": "oauth-token", "chatgpt_account_id": "chatgpt-acc"},
		Status:      StatusActive, Schedulable: true,
	}
	body := []byte(`{"model":"gpt-5.4","instructions":"compact-test","input":[{"type":"message","role":"user","content":"hello"}]}`)

	result, err := svc.Forward(context.Background(), c, account, body)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.Stream)
	require.Equal(t, "resp_legacy_v2", result.ResponseID)
	require.Equal(t, 7, result.Usage.InputTokens)
	require.Equal(t, 2, result.Usage.OutputTokens)

	require.Equal(t, chatgptCodexURL, upstream.lastReq.URL.String())
	require.Equal(t, "text/event-stream", upstream.lastReq.Header.Get("Accept"))
	require.Equal(t, "existing_feature,remote_compaction_v2", upstream.lastReq.Header.Get("x-codex-beta-features"))
	require.True(t, gjson.GetBytes(upstream.lastBody, "stream").Bool())
	require.Equal(t, gjson.False, gjson.GetBytes(upstream.lastBody, "store").Type)
	require.Equal(t, 1, countCompactionTriggers(upstream.lastBody))

	require.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	require.NotContains(t, rec.Body.String(), "data:")
	require.NotContains(t, rec.Body.String(), "event:")
	require.Equal(t, "resp_legacy_v2", gjson.Get(rec.Body.String(), "id").String())
	require.Equal(t, "compaction", gjson.Get(rec.Body.String(), "output.0.type").String())
	require.Equal(t, "legacy-v2", gjson.Get(rec.Body.String(), "output.0.encrypted_content").String())
}
