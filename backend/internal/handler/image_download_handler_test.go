package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type imageDownloadFetcherFunc func(context.Context, string) (*service.ImageDownload, error)

func (f imageDownloadFetcherFunc) Fetch(ctx context.Context, rawURL string) (*service.ImageDownload, error) {
	return f(ctx, rawURL)
}

func imageDownloadHandlerRequest(t *testing.T, h *ImageDownloadHandler, body any) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{ID: 123})
		c.Next()
	})
	router.POST("/v1/images/download", h.Download)
	payload, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/v1/images/download", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func TestImageDownloadHandlerStreamsAttachment(t *testing.T) {
	payload := append([]byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}, bytes.Repeat([]byte{0x01}, 2048)...)
	const rawURL = "https://" + service.ImageDownloadAllowedHost + "/images/result.png?signature=test"
	h := &ImageDownloadHandler{download: imageDownloadFetcherFunc(func(_ context.Context, gotURL string) (*service.ImageDownload, error) {
		require.Equal(t, rawURL, gotURL)
		return &service.ImageDownload{
			Body:          io.NopCloser(bytes.NewReader(payload)),
			ContentType:   "image/png",
			ContentLength: int64(len(payload)),
		}, nil
	}), limiter: newImageDownloadConcurrencyGate()}

	recorder := imageDownloadHandlerRequest(t, h, map[string]string{
		"url":      rawURL,
		"filename": "../../report\r\n.exe",
	})

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, payload, recorder.Body.Bytes())
	require.Equal(t, "private, no-store", recorder.Header().Get("Cache-Control"))
	require.Equal(t, "image/png", recorder.Header().Get("Content-Type"))
	require.Equal(t, "no", recorder.Header().Get("X-Accel-Buffering"))
	require.Equal(t, "nosniff", recorder.Header().Get("X-Content-Type-Options"))
	require.Equal(t, "attachment; filename=report__.png", recorder.Header().Get("Content-Disposition"))
	require.Equal(t, len(payload), recorder.Body.Len())
}

func TestImageDownloadHandlerReturnsJSONBeforeStreamingOnFailure(t *testing.T) {
	for _, tc := range []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "invalid URL", err: service.ErrImageDownloadInvalidURL, wantStatus: http.StatusBadRequest, wantCode: "invalid_image_url"},
		{name: "missing length", err: service.ErrImageDownloadInvalidLength, wantStatus: http.StatusBadGateway, wantCode: "invalid_upstream_response"},
		{name: "oversized", err: service.ErrImageDownloadTooLarge, wantStatus: http.StatusBadGateway, wantCode: "image_too_large"},
		{name: "expired", err: service.ErrImageDownloadUnavailable, wantStatus: http.StatusBadGateway, wantCode: "image_unavailable"},
		{name: "unknown upstream", err: errors.New("reset"), wantStatus: http.StatusBadGateway, wantCode: "image_unavailable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := &ImageDownloadHandler{download: imageDownloadFetcherFunc(func(context.Context, string) (*service.ImageDownload, error) {
				return nil, tc.err
			}), limiter: newImageDownloadConcurrencyGate()}
			recorder := imageDownloadHandlerRequest(t, h, map[string]string{
				"url": "https://" + service.ImageDownloadAllowedHost + "/image.png",
			})
			require.Equal(t, tc.wantStatus, recorder.Code)
			require.Equal(t, "no-store", recorder.Header().Get("Cache-Control"))
			require.Empty(t, recorder.Header().Get("Content-Disposition"))
			require.Contains(t, recorder.Body.String(), tc.wantCode)
		})
	}
}

func TestImageDownloadHandlerRejectsMissingURL(t *testing.T) {
	h := &ImageDownloadHandler{download: imageDownloadFetcherFunc(func(context.Context, string) (*service.ImageDownload, error) {
		t.Fatal("fetch must not be called")
		return nil, nil
	}), limiter: newImageDownloadConcurrencyGate()}
	recorder := imageDownloadHandlerRequest(t, h, map[string]string{"filename": "image.png"})
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "invalid_request_error")
}

func TestImageDownloadHandlerRejectsPerKeyConcurrencyOverflow(t *testing.T) {
	gate := newImageDownloadConcurrencyGate()
	releaseOne, acquired := gate.TryAcquire(123)
	require.True(t, acquired)
	releaseTwo, acquired := gate.TryAcquire(123)
	require.True(t, acquired)
	defer releaseOne()
	defer releaseTwo()

	h := &ImageDownloadHandler{
		download: imageDownloadFetcherFunc(func(context.Context, string) (*service.ImageDownload, error) {
			t.Fatal("fetch must not be called when the API key is at its download limit")
			return nil, nil
		}),
		limiter: gate,
	}
	recorder := imageDownloadHandlerRequest(t, h, map[string]string{
		"url": "https://" + service.ImageDownloadAllowedHost + "/image.png",
	})

	require.Equal(t, http.StatusTooManyRequests, recorder.Code)
	require.Contains(t, recorder.Body.String(), "rate_limit_error")
}

func TestImageDownloadConcurrencyGateBoundsTotalAndReleases(t *testing.T) {
	gate := newImageDownloadConcurrencyGate()
	releases := make([]func(), 0, imageDownloadMaxConcurrent)
	for i := 0; i < imageDownloadMaxConcurrent; i++ {
		release, acquired := gate.TryAcquire(int64(i + 1))
		require.True(t, acquired)
		releases = append(releases, release)
	}
	_, acquired := gate.TryAcquire(999)
	require.False(t, acquired)

	releases[0]()
	release, acquired := gate.TryAcquire(999)
	require.True(t, acquired)
	release()
	for _, release := range releases[1:] {
		release()
	}

	gate.mu.Lock()
	defer gate.mu.Unlock()
	require.Zero(t, gate.active)
	require.Empty(t, gate.byAPIKey)
}

func TestImageDownloadConcurrencyGateRateLimitsAndResets(t *testing.T) {
	t.Run("per key", func(t *testing.T) {
		gate := newImageDownloadConcurrencyGate()
		for i := 0; i < imageDownloadMaxPerKeyMinute; i++ {
			release, acquired := gate.TryAcquire(123)
			require.True(t, acquired)
			release()
		}
		_, acquired := gate.TryAcquire(123)
		require.False(t, acquired)

		gate.mu.Lock()
		gate.windowAt = time.Now().Add(-imageDownloadRateWindow)
		gate.mu.Unlock()
		release, acquired := gate.TryAcquire(123)
		require.True(t, acquired)
		release()
	})

	t.Run("process total", func(t *testing.T) {
		gate := newImageDownloadConcurrencyGate()
		for i := 0; i < imageDownloadMaxTotalMinute; i++ {
			apiKeyID := int64(i/imageDownloadMaxPerKeyMinute + 1)
			release, acquired := gate.TryAcquire(apiKeyID)
			require.True(t, acquired)
			release()
		}
		_, acquired := gate.TryAcquire(999)
		require.False(t, acquired)
	})
}

func TestImageDownloadConcurrencyGateParallelAcquireAndIdempotentRelease(t *testing.T) {
	const workers = 64
	gate := newImageDownloadConcurrencyGate()
	unblock := make(chan struct{})
	var attempted sync.WaitGroup
	var completed sync.WaitGroup
	var acquired atomic.Int64
	attempted.Add(workers)
	completed.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer completed.Done()
			release, ok := gate.TryAcquire(123)
			attempted.Done()
			if !ok {
				return
			}
			acquired.Add(1)
			<-unblock
			release()
			release()
		}()
	}

	attempted.Wait()
	require.Equal(t, int64(imageDownloadMaxPerAPIKey), acquired.Load())
	close(unblock)
	completed.Wait()

	gate.mu.Lock()
	defer gate.mu.Unlock()
	require.Zero(t, gate.active)
	require.Empty(t, gate.byAPIKey)
}

func TestImageDownloadHandlerRecordsTruncatedStream(t *testing.T) {
	payload := append([]byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a}, bytes.Repeat([]byte{0x01}, 32)...)
	h := &ImageDownloadHandler{
		download: imageDownloadFetcherFunc(func(context.Context, string) (*service.ImageDownload, error) {
			return &service.ImageDownload{
				Body:          io.NopCloser(bytes.NewReader(payload)),
				ContentType:   "image/png",
				ContentLength: int64(len(payload) + 1),
			}, nil
		}),
		limiter: newImageDownloadConcurrencyGate(),
	}

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	body, err := json.Marshal(map[string]string{"url": "https://" + service.ImageDownloadAllowedHost + "/image.png"})
	require.NoError(t, err)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/images/download", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{ID: 123})

	h.Download(c)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.NotEmpty(t, c.Errors)
	require.Contains(t, c.Errors.String(), "wrote")
}

func TestSanitizeImageDownloadFilename(t *testing.T) {
	for _, tc := range []struct {
		name        string
		input       string
		contentType string
		want        string
	}{
		{name: "default", contentType: "image/png", want: "generated-image.png"},
		{name: "path traversal", input: "../../final.exe", contentType: "image/png", want: "final.png"},
		{name: "windows path", input: `..\..\final.jpg`, contentType: "image/jpeg", want: "final.jpg"},
		{name: "replace extension", input: "photo.png", contentType: "image/webp", want: "photo.webp"},
		{name: "control characters", input: "photo\r\n.png", contentType: "image/png", want: "photo__.png"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, sanitizeImageDownloadFilename(tc.input, tc.contentType))
		})
	}
}
