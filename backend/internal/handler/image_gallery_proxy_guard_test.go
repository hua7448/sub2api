package handler

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestImageGalleryResponsesProxyGuard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	baseSettings := service.DefaultImageGallerySettings()

	tests := []struct {
		name     string
		body     string
		mutate   func(*service.ImageGallerySettings)
		wantCode int
	}{
		{
			name: "allows standard responses image generation when agent disabled",
			body: `{"model":"gpt-5.5","input":"draw","tools":[{"type":"image_generation","size":"1024x1024","quality":"auto","output_format":"png"}],"tool_choice":"required"}`,
			mutate: func(settings *service.ImageGallerySettings) {
				settings.AgentEnabled = false
			},
			wantCode: http.StatusOK,
		},
		{
			name: "rejects agent function tools when agent disabled",
			body: `{"model":"gpt-5.5","input":"draw","tools":[{"type":"image_generation","size":"1024x1024","quality":"auto","output_format":"png"},{"type":"function","name":"generate_image_batch"}]}`,
			mutate: func(settings *service.ImageGallerySettings) {
				settings.AgentEnabled = false
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "rejects web search when disabled",
			body: `{"model":"gpt-5.5","input":"draw","tools":[{"type":"image_generation","size":"1024x1024","quality":"auto","output_format":"png"},{"type":"web_search"}]}`,
			mutate: func(settings *service.ImageGallerySettings) {
				settings.AgentEnabled = true
				settings.AgentWebSearchEnabled = false
			},
			wantCode: http.StatusForbidden,
		},
		{
			name:     "rejects disallowed responses model",
			body:     `{"model":"gpt-image-2","input":"draw","tools":[{"type":"image_generation","size":"1024x1024","quality":"auto","output_format":"png"}]}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "rejects disallowed image tool size",
			body:     `{"model":"gpt-5.5","input":"draw","tools":[{"type":"image_generation","size":"2048x2048","quality":"auto","output_format":"png"}]}`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := baseSettings
			if tt.mutate != nil {
				tt.mutate(&settings)
			}
			c, w := newImageGalleryProxyGuardContext(http.MethodPost, "/api/v1/image-playground/proxy/responses", "application/json", strings.NewReader(tt.body))
			h := &ImageGalleryHandler{}

			ok := h.validateResponsesProxyRequest(c, settings)

			if tt.wantCode == http.StatusOK {
				require.True(t, ok)
				require.Equal(t, http.StatusOK, w.Code)
				return
			}
			require.False(t, ok)
			require.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestImageGalleryImagesProxyGuard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	settings := service.DefaultImageGallerySettings()
	h := &ImageGalleryHandler{}

	t.Run("rejects disallowed image model", func(t *testing.T) {
		body := `{"model":"gpt-5.5","prompt":"draw","size":"1024x1024","quality":"auto","output_format":"png","n":1}`
		c, w := newImageGalleryProxyGuardContext(http.MethodPost, "/api/v1/image-playground/proxy/images/generations", "application/json", strings.NewReader(body))

		ok := h.validateImagesGenerationProxyRequest(c, settings)

		require.False(t, ok)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("rejects excessive image count", func(t *testing.T) {
		body := `{"model":"gpt-image-2","prompt":"draw","size":"1024x1024","quality":"auto","output_format":"png","n":99}`
		c, w := newImageGalleryProxyGuardContext(http.MethodPost, "/api/v1/image-playground/proxy/images/generations", "application/json", strings.NewReader(body))

		ok := h.validateImagesGenerationProxyRequest(c, settings)

		require.False(t, ok)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("allows multipart edit within limits", func(t *testing.T) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		require.NoError(t, writer.WriteField("model", "gpt-image-2"))
		require.NoError(t, writer.WriteField("prompt", "edit"))
		require.NoError(t, writer.WriteField("size", "1024x1024"))
		require.NoError(t, writer.WriteField("quality", "auto"))
		require.NoError(t, writer.WriteField("output_format", "png"))
		file, err := writer.CreateFormFile("image[]", "image.png")
		require.NoError(t, err)
		_, err = file.Write([]byte("fake-image"))
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		c, w := newImageGalleryProxyGuardContext(http.MethodPost, "/api/v1/image-playground/proxy/images/edits", writer.FormDataContentType(), &body)

		ok := h.validateImagesEditProxyRequest(c, settings)

		require.True(t, ok)
		require.Equal(t, http.StatusOK, w.Code)
	})
}

func newImageGalleryProxyGuardContext(method, path, contentType string, body io.Reader) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, body)
	c.Request.Header.Set("Content-Type", contentType)
	return c, w
}
