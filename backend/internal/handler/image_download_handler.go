package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	imageDownloadCopyBufferBytes = 32 << 10
	imageDownloadRequestMaxBytes = 8 << 10
	imageDownloadMaxConcurrent   = 16
	imageDownloadMaxPerAPIKey    = 2
	imageDownloadMaxPerKeyMinute = 30
	imageDownloadMaxTotalMinute  = 240
	imageDownloadRateWindow      = time.Minute
)

type imageDownloadFetcher interface {
	Fetch(ctx context.Context, rawURL string) (*service.ImageDownload, error)
}

type ImageDownloadHandler struct {
	download imageDownloadFetcher
	limiter  *imageDownloadConcurrencyGate
}

type imageDownloadRequest struct {
	URL      string `json:"url" binding:"required"`
	Filename string `json:"filename"`
}

type imageDownloadConcurrencyGate struct {
	mu       sync.Mutex
	active   int
	byAPIKey map[int64]int
	windowAt time.Time
	requests int
	byKeyRPM map[int64]int
}

func NewImageDownloadHandler() *ImageDownloadHandler {
	return &ImageDownloadHandler{
		download: service.NewImageDownloadService(),
		limiter:  newImageDownloadConcurrencyGate(),
	}
}

func newImageDownloadConcurrencyGate() *imageDownloadConcurrencyGate {
	return &imageDownloadConcurrencyGate{
		byAPIKey: make(map[int64]int),
		byKeyRPM: make(map[int64]int),
	}
}

func (g *imageDownloadConcurrencyGate) TryAcquire(apiKeyID int64) (func(), bool) {
	if g == nil || apiKeyID <= 0 {
		return nil, false
	}
	g.mu.Lock()
	now := time.Now()
	if g.windowAt.IsZero() || now.Sub(g.windowAt) >= imageDownloadRateWindow {
		g.windowAt = now
		g.requests = 0
		clear(g.byKeyRPM)
	}
	if g.active >= imageDownloadMaxConcurrent ||
		g.byAPIKey[apiKeyID] >= imageDownloadMaxPerAPIKey ||
		g.requests >= imageDownloadMaxTotalMinute ||
		g.byKeyRPM[apiKeyID] >= imageDownloadMaxPerKeyMinute {
		g.mu.Unlock()
		return nil, false
	}
	g.active++
	g.byAPIKey[apiKeyID]++
	g.requests++
	g.byKeyRPM[apiKeyID]++
	g.mu.Unlock()

	var once sync.Once
	return func() {
		once.Do(func() {
			g.mu.Lock()
			if g.active > 0 {
				g.active--
			}
			if remaining := g.byAPIKey[apiKeyID] - 1; remaining > 0 {
				g.byAPIKey[apiKeyID] = remaining
			} else {
				delete(g.byAPIKey, apiKeyID)
			}
			g.mu.Unlock()
		})
	}, true
}

func (h *ImageDownloadHandler) Download(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, imageDownloadRequestMaxBytes)
	var req imageDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		imageDownloadJSONError(c, http.StatusBadRequest, "invalid_request_error", "url is required")
		return
	}
	if h == nil || h.download == nil {
		imageDownloadJSONError(c, http.StatusServiceUnavailable, "api_error", "image download is unavailable")
		return
	}
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		imageDownloadJSONError(c, http.StatusUnauthorized, "authentication_error", "API key authentication is required")
		return
	}
	if h.limiter == nil {
		imageDownloadJSONError(c, http.StatusServiceUnavailable, "api_error", "image download is unavailable")
		return
	}
	release, acquired := h.limiter.TryAcquire(apiKey.ID)
	if !acquired {
		imageDownloadJSONError(c, http.StatusTooManyRequests, "rate_limit_error", "too many concurrent image downloads")
		return
	}
	defer release()

	download, err := h.download.Fetch(c.Request.Context(), req.URL)
	if err != nil {
		imageDownloadError(c, err)
		return
	}
	defer func() { _ = download.Close() }()
	if !download.Valid() {
		imageDownloadJSONError(c, http.StatusBadGateway, "invalid_upstream_response", "upstream returned an invalid image")
		return
	}

	filename := sanitizeImageDownloadFilename(req.Filename, download.ContentType)
	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": filename})
	c.Header("Cache-Control", "private, no-store")
	c.Header("Content-Disposition", disposition)
	c.Header("Content-Length", strconv.FormatInt(download.ContentLength, 10))
	c.Header("Content-Type", download.ContentType)
	c.Header("X-Accel-Buffering", "no")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Status(http.StatusOK)
	buffer := make([]byte, imageDownloadCopyBufferBytes)
	written, copyErr := io.CopyBuffer(c.Writer, download.Body, buffer)
	if copyErr != nil {
		_ = c.Error(fmt.Errorf("stream image download: %w", copyErr))
		return
	}
	if written != download.ContentLength {
		_ = c.Error(fmt.Errorf("stream image download: wrote %d of %d bytes", written, download.ContentLength))
	}
}

func imageDownloadError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrImageDownloadInvalidURL):
		imageDownloadJSONError(c, http.StatusBadRequest, "invalid_image_url", "url must be an HTTPS image URL from the image generation service")
	case errors.Is(err, service.ErrImageDownloadTooLarge):
		imageDownloadJSONError(c, http.StatusBadGateway, "image_too_large", "upstream image exceeds the 64 MiB download limit")
	case errors.Is(err, service.ErrImageDownloadInvalidLength):
		imageDownloadJSONError(c, http.StatusBadGateway, "invalid_upstream_response", "upstream image did not provide a valid content length")
	case errors.Is(err, service.ErrImageDownloadInvalidType):
		imageDownloadJSONError(c, http.StatusBadGateway, "invalid_upstream_response", "upstream response is not a supported image")
	default:
		imageDownloadJSONError(c, http.StatusBadGateway, "image_unavailable", "image URL has expired or is unavailable")
	}
}

func imageDownloadJSONError(c *gin.Context, status int, code, message string) {
	c.Header("Cache-Control", "no-store")
	c.JSON(status, gin.H{"error": gin.H{"type": code, "code": code, "message": message}})
}

func sanitizeImageDownloadFilename(input, contentType string) string {
	input = strings.ReplaceAll(strings.TrimSpace(input), `\`, "/")
	base := path.Base(input)
	if base == "." || base == "/" || base == "" {
		base = "generated-image"
	}
	runes := make([]rune, 0, len(base))
	for _, r := range base {
		if unicode.IsControl(r) || strings.ContainsRune(`<>:"/\|?*`, r) {
			r = '_'
		}
		runes = append(runes, r)
		if len(runes) == 100 {
			break
		}
	}
	base = strings.Trim(strings.TrimSpace(string(runes)), ".")
	if base == "" {
		base = "generated-image"
	}
	extension := imageDownloadExtension(contentType)
	if current := path.Ext(base); current != "" {
		base = strings.TrimSuffix(base, current)
	}
	base = strings.TrimRight(base, ". ")
	if base == "" {
		base = "generated-image"
	}
	return base + extension
}

func imageDownloadExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".png"
	}
}
