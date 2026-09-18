package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestOpenAIImagesScheduleLatencyMs(t *testing.T) {
	ttft := 123
	require.Same(t, &ttft, openAIImagesScheduleLatencyMs(&service.OpenAIForwardResult{
		Duration:     5 * time.Second,
		FirstTokenMs: &ttft,
	}))

	latency := openAIImagesScheduleLatencyMs(&service.OpenAIForwardResult{
		Duration: 3*time.Minute + 250*time.Millisecond,
	})
	require.NotNil(t, latency)
	require.Equal(t, 180250, *latency)

	require.Nil(t, openAIImagesScheduleLatencyMs(nil))
	require.Nil(t, openAIImagesScheduleLatencyMs(&service.OpenAIForwardResult{}))
}

func TestOpenAIImagesRequestSessionHashUsesStudioTaskKeyFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := &OpenAIGatewayHandler{gatewayService: &service.OpenAIGatewayService{}}
	body := []byte(`{"model":"gpt-image-2","prompt":"draw a cat"}`)

	t.Run("empty without explicit session or studio task", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)

		require.Empty(t, handler.imageRequestSessionHash(c, body))
	})

	t.Run("studio task key becomes per-image scheduling seed", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)
		c.Request.Header.Set(studioImageTaskKeyHeader, "batch:3")

		got := handler.imageRequestSessionHash(c, body)
		require.Equal(t, fmt.Sprintf("%016x", xxhash.Sum64String("openai-images-task:batch:3")), got)
	})

	t.Run("explicit session still wins", func(t *testing.T) {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/generations", nil)
		c.Request.Header.Set("session_id", "explicit-session")
		c.Request.Header.Set(studioImageTaskKeyHeader, "batch:3")

		got := handler.imageRequestSessionHash(c, body)
		require.Equal(t, fmt.Sprintf("%016x", xxhash.Sum64String("explicit-session")), got)
	})
}
