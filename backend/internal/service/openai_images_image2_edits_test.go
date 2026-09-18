package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

func TestParseOpenAIImagesRequest_JSONEditTypeURLAliasesResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"gpt-image-2","prompt":"edit","size":"3840x2160","images":[{"image_url":"https://example.com/source.png"}],"type":"url"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = req

	parsed, err := (&OpenAIGatewayService{}).ParseOpenAIImagesRequest(ctx, body)
	require.NoError(t, err)
	require.Equal(t, "url", parsed.ResponseFormat)
	require.Equal(t, "3840x2160", parsed.Size)
	require.Equal(t, []string{"https://example.com/source.png"}, parsed.InputImageURLs)
	require.Equal(t, OpenAIImagesCapabilityExactSize, parsed.RequiredCapability)
}

func TestOpenAIGatewayServiceForwardImages_APIKeyJSONEditPreservesExactSizeAndNestedImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"gpt-image-2","prompt":"edit","size":"3840x2160","images":[{"image_url":"https://example.com/source.png"}],"type":"url"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"created":1710000010,"data":[{"url":"https://image-result.example/exact.png","size":"3840x2160"}]}`)),
	}}
	svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
	parsed, err := svc.ParseOpenAIImagesRequest(ctx, body)
	require.NoError(t, err)

	account := &Account{ID: 99, Name: "openai-apikey", Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Credentials: map[string]any{
		"api_key": "test-api-key", "base_url": "https://image-upstream.example/v1/",
	}, Extra: map[string]any{OpenAIImageExactSizeSupportedExtraKey: true}}
	result, err := svc.ForwardImages(context.Background(), ctx, account, body, parsed, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "https://image-upstream.example/v1/images/edits", upstream.lastReq.URL.String())
	require.Equal(t, "gpt-image-2", gjson.GetBytes(upstream.lastBody, "model").String())
	require.Equal(t, "3840x2160", gjson.GetBytes(upstream.lastBody, "size").String())
	require.Equal(t, "url", gjson.GetBytes(upstream.lastBody, "type").String())
	require.Equal(t, "https://example.com/source.png", gjson.GetBytes(upstream.lastBody, "images.0.image_url").String())
	require.Equal(t, "https://image-result.example/exact.png", gjson.Get(rec.Body.String(), "data.0.url").String())
	require.Equal(t, "3840x2160", gjson.Get(rec.Body.String(), "data.0.size").String())
}

func TestOpenAIGatewayServiceForwardImages_OAuthJSONEditExactSizeFailsBeforeUpstream(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"gpt-image-2","prompt":"edit","size":"3840x2160","images":[{"image_url":"https://example.com/source.png"}],"type":"url"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = req

	upstream := &httpUpstreamRecorder{}
	svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
	parsed, err := svc.ParseOpenAIImagesRequest(ctx, body)
	require.NoError(t, err)

	account := &Account{ID: 100, Name: "openai-oauth", Platform: PlatformOpenAI, Type: AccountTypeOAuth}
	result, err := svc.ForwardImages(context.Background(), ctx, account, body, parsed, "")
	require.ErrorContains(t, err, "exact image size requires an opted-in API key image account")
	require.Nil(t, result)
	require.Nil(t, upstream.lastReq)
}

func TestOpenAIGatewayServiceForwardImages_OAuthJSONEditStandardSizeUsesTypeURLAlias(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"gpt-image-2","prompt":"edit","size":"1024x1024","images":[{"image_url":"https://example.com/source.png"}],"type":"url"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
		Body: io.NopCloser(strings.NewReader(
			"data: {\"type\":\"response.created\",\"response\":{\"created_at\":1710000011,\"tools\":[{\"type\":\"image_generation\",\"model\":\"gpt-image-2\",\"size\":\"1024x1024\"}]}}\n\n" +
				"data: {\"type\":\"response.completed\",\"response\":{\"created_at\":1710000011,\"usage\":{\"input_tokens\":1,\"output_tokens\":2},\"tools\":[{\"type\":\"image_generation\",\"model\":\"gpt-image-2\",\"size\":\"1024x1024\"}],\"output\":[{\"type\":\"image_generation_call\",\"result\":\"ZWRpdGVk\",\"output_format\":\"png\"}]}}\n\n" +
				"data: [DONE]\n\n",
		)),
	}}
	svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
	parsed, err := svc.ParseOpenAIImagesRequest(ctx, body)
	require.NoError(t, err)
	require.Equal(t, "url", parsed.ResponseFormat)
	require.Equal(t, OpenAIImagesCapabilityNative, parsed.RequiredCapability)

	account := &Account{ID: 101, Name: "openai-oauth", Platform: PlatformOpenAI, Type: AccountTypeOAuth, Credentials: map[string]any{
		"access_token": "token-123",
	}}
	result, err := svc.ForwardImages(context.Background(), ctx, account, body, parsed, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "edit", gjson.GetBytes(upstream.lastBody, "tools.0.action").String())
	require.Equal(t, "1024x1024", gjson.GetBytes(upstream.lastBody, "tools.0.size").String())
	require.Equal(t, "https://example.com/source.png", gjson.GetBytes(upstream.lastBody, "input.0.content.1.image_url").String())
	require.Equal(t, "data:image/png;base64,ZWRpdGVk", gjson.Get(rec.Body.String(), "data.0.url").String())
}

func TestParseOpenAIImagesRequest_ResponseFormatWinsTypeAlias(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"gpt-image-2","prompt":"edit","size":"1024x1024","images":[{"image_url":"https://example.com/source.png"}],"response_format":"b64_json","type":"url"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = req

	parsed, err := (&OpenAIGatewayService{}).ParseOpenAIImagesRequest(ctx, body)
	require.NoError(t, err)
	require.Equal(t, "b64_json", parsed.ResponseFormat)
}

func TestOpenAIGatewayServiceForwardImages_ExactSizeRejectsLegacyModelMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"gpt-image-2","prompt":"edit","size":"3840x2160","images":[{"image_url":"https://example.com/source.png"}],"type":"url"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = req

	upstream := &httpUpstreamRecorder{}
	svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
	parsed, err := svc.ParseOpenAIImagesRequest(ctx, body)
	require.NoError(t, err)

	account := &Account{ID: 102, Name: "mapped-apikey", Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Credentials: map[string]any{
		"api_key": "test-api-key",
	}, Extra: map[string]any{OpenAIImageExactSizeSupportedExtraKey: true}}
	result, err := svc.ForwardImages(context.Background(), ctx, account, body, parsed, "gpt-image-1")
	require.ErrorContains(t, err, "exact image size requires upstream model gpt-image-2")
	require.Nil(t, result)
	require.Nil(t, upstream.lastReq)
}

func TestOpenAIGatewayServiceForwardImages_UnmarkedAPIKeyExactSizeFailsBeforeUpstream(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"gpt-image-2","prompt":"edit","size":"3840x2160","images":[{"image_url":"https://example.com/source.png"}],"type":"url"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/images/edits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = req

	upstream := &httpUpstreamRecorder{}
	svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
	parsed, err := svc.ParseOpenAIImagesRequest(ctx, body)
	require.NoError(t, err)

	account := &Account{ID: 104, Name: "unmarked-apikey", Platform: PlatformOpenAI, Type: AccountTypeAPIKey}
	result, err := svc.ForwardImages(context.Background(), ctx, account, body, parsed, "")
	require.ErrorContains(t, err, "exact image size requires an opted-in API key image account")
	require.Nil(t, result)
	require.Nil(t, upstream.lastReq)
}

func TestOpenAIGatewayServiceForwardImages_APIKeyExactSizeStreamingPreservesRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"gpt-image-2","prompt":"draw","size":"3840x2160","stream":true,"response_format":"b64_json"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/images/generations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"created":1710000012,"data":[{"b64_json":"ZmluYWw=","size":"3840x2160"}]}`)),
	}}
	svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
	parsed, err := svc.ParseOpenAIImagesRequest(ctx, body)
	require.NoError(t, err)
	require.Equal(t, OpenAIImagesCapabilityExactSize, parsed.RequiredCapability)

	account := &Account{ID: 103, Name: "stream-apikey", Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Credentials: map[string]any{
		"api_key": "test-api-key", "base_url": "https://image-upstream.example/v1/",
	}, Extra: map[string]any{OpenAIImageExactSizeSupportedExtraKey: true}}
	result, err := svc.ForwardImages(context.Background(), ctx, account, body, parsed, "")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.Stream)
	require.Equal(t, "3840x2160", gjson.GetBytes(upstream.lastBody, "size").String())
	require.Equal(t, "3840x2160", gjson.Get(rec.Body.String(), "data.0.size").String())
}

func TestOpenAIImagesRequestResolveRequiredCapabilityForEffectiveModel(t *testing.T) {
	t.Run("channel mapping to Image2 requires exact-size account", func(t *testing.T) {
		req := &OpenAIImagesRequest{Model: "public-image-alias", Size: "3840x2160", ExplicitModel: true, ExplicitSize: true}
		require.NoError(t, req.ResolveRequiredCapabilityForEffectiveModel("gpt-image-2"))
		require.Equal(t, OpenAIImagesCapabilityExactSize, req.RequiredCapability)
	})

	t.Run("Image2 exact size cannot map to legacy model", func(t *testing.T) {
		req := &OpenAIImagesRequest{Model: "gpt-image-2", Size: "3840x2160", ExplicitModel: true, ExplicitSize: true}
		err := req.ResolveRequiredCapabilityForEffectiveModel("gpt-image-1")
		require.ErrorContains(t, err, "exact image size requires upstream model gpt-image-2")
	})
}
