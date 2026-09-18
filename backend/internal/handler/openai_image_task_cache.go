package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	studioImageTaskKeyHeader = "X-Sub2API-Image-Task-Key"
	studioImageTaskCacheTTL  = 30 * time.Minute
	studioImageTaskMaxBody   = 32 << 20
)

// studioImageTaskReplayBlockedHeaders lists headers that must not be replayed
// from a cached image task back to a different client. They either expose the
// original request's identifiers or are connection-specific.
var studioImageTaskReplayBlockedHeaders = map[string]struct{}{
	"content-length":            {},
	"date":                      {},
	"server":                    {},
	"connection":                {},
	"keep-alive":                {},
	"transfer-encoding":         {},
	"set-cookie":                {},
	"x-request-id":              {},
	"x-client-request-id":       {},
	"x-sub2api-image-task-status": {},
}

type studioImageTask struct {
	requestHash string
	done        chan struct{}
	startedAt   time.Time

	mu        sync.Mutex
	status    int
	header    http.Header
	body      []byte
	truncated bool
}

var studioImageTasks sync.Map

type studioImageTaskCaptureWriter struct {
	gin.ResponseWriter
	buf       bytes.Buffer
	truncated bool
}

func (w *studioImageTaskCaptureWriter) Write(b []byte) (int, error) {
	w.capture(b)
	return w.ResponseWriter.Write(b)
}

func (w *studioImageTaskCaptureWriter) WriteString(s string) (int, error) {
	w.capture([]byte(s))
	return w.ResponseWriter.WriteString(s)
}

func (w *studioImageTaskCaptureWriter) capture(b []byte) {
	if w.truncated || len(b) == 0 {
		return
	}
	remaining := studioImageTaskMaxBody - w.buf.Len()
	if remaining <= 0 {
		w.truncated = true
		return
	}
	if len(b) > remaining {
		_, _ = w.buf.Write(b[:remaining])
		w.truncated = true
		return
	}
	_, _ = w.buf.Write(b)
}

func handleStudioImageTaskRequest(c *gin.Context, userID, apiKeyID int64, parsed *service.OpenAIImagesRequest, body []byte) (bool, func()) {
	if c == nil || c.Request == nil || parsed == nil || parsed.Stream {
		return false, nil
	}
	taskKey := normalizeStudioImageTaskKey(c.GetHeader(studioImageTaskKeyHeader))
	if taskKey == "" {
		return false, nil
	}

	scope := fmt.Sprintf("%d:%d:%s", userID, apiKeyID, taskKey)
	requestHash := studioImageTaskRequestHash(parsed.Endpoint, parsed.ContentType, body)
	task, owner, conflict := acquireStudioImageTask(scope, requestHash)
	if conflict {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"type":    "conflict_error",
				"message": "Image task key is already used for a different request",
			},
		})
		return true, nil
	}
	if !owner {
		serveStudioImageTask(c, task)
		return true, nil
	}

	originalWriter := c.Writer
	captureWriter := &studioImageTaskCaptureWriter{ResponseWriter: originalWriter}
	c.Writer = captureWriter
	c.Request = c.Request.WithContext(context.WithoutCancel(c.Request.Context()))

	finish := func() {
		task.complete(captureWriter)
		if c.Writer == captureWriter {
			c.Writer = originalWriter
		}
	}
	return false, finish
}

func acquireStudioImageTask(scope, requestHash string) (*studioImageTask, bool, bool) {
	task := &studioImageTask{
		requestHash: requestHash,
		done:        make(chan struct{}),
		startedAt:   time.Now(),
	}
	scheduleStudioImageTaskCleanup := func(t *studioImageTask) {
		time.AfterFunc(studioImageTaskCacheTTL, func() {
			if current, ok := studioImageTasks.Load(scope); ok && current == t {
				studioImageTasks.Delete(scope)
			}
		})
	}
	actual, loaded := studioImageTasks.LoadOrStore(scope, task)
	if !loaded {
		scheduleStudioImageTaskCleanup(task)
		return task, true, false
	}
	existing, ok := actual.(*studioImageTask)
	if !ok || existing == nil {
		studioImageTasks.Store(scope, task)
		scheduleStudioImageTaskCleanup(task)
		return task, true, false
	}
	existing.mu.Lock()
	conflict := existing.requestHash != requestHash
	existing.mu.Unlock()
	return existing, false, conflict
}

func serveStudioImageTask(c *gin.Context, task *studioImageTask) {
	if c == nil || c.Request == nil || task == nil {
		return
	}
	select {
	case <-task.done:
	case <-c.Request.Context().Done():
		return
	}

	status, header, body, truncated := task.snapshot()
	if truncated {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"type":    "api_error",
				"message": "Cached image task response is too large to replay",
			},
		})
		return
	}
	for key, values := range header {
		if _, blocked := studioImageTaskReplayBlockedHeaders[strings.ToLower(key)]; blocked {
			continue
		}
		c.Writer.Header().Del(key)
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	c.Header("X-Sub2API-Image-Task-Status", "replayed")
	contentType := strings.TrimSpace(header.Get("Content-Type"))
	if contentType == "" {
		contentType = "application/json"
	}
	if status == 0 {
		status = http.StatusOK
	}
	c.Data(status, contentType, body)
}

func (t *studioImageTask) complete(w *studioImageTaskCaptureWriter) {
	if t == nil || w == nil {
		return
	}
	t.mu.Lock()
	t.status = w.Status()
	if t.status == 0 {
		t.status = http.StatusOK
	}
	t.header = w.Header().Clone()
	t.body = append(t.body[:0], w.buf.Bytes()...)
	t.truncated = w.truncated
	t.mu.Unlock()
	close(t.done)
}

func (t *studioImageTask) snapshot() (int, http.Header, []byte, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.status, t.header.Clone(), append([]byte(nil), t.body...), t.truncated
}

func normalizeStudioImageTaskKey(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if len(raw) > 256 {
		sum := sha256.Sum256([]byte(raw))
		raw = raw[:128] + "." + hex.EncodeToString(sum[:8])
	}
	var b strings.Builder
	b.Grow(len(raw))
	for _, r := range raw {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.' || r == ':':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}

func studioImageTaskRequestHash(endpoint, contentType string, body []byte) string {
	h := sha256.New()
	_, _ = io.WriteString(h, strings.TrimSpace(endpoint))
	_, _ = io.WriteString(h, "\n")
	normalizedContentType := strings.ToLower(strings.TrimSpace(contentType))
	if strings.HasPrefix(normalizedContentType, "multipart/form-data") {
		// Browser-generated multipart boundaries change across refresh recovery
		// attempts even when the logical image task is the same. Use the body
		// length as a coarse-grained fingerprint so that two clearly different
		// uploads sharing a task key are still detected as a conflict.
		_, _ = io.WriteString(h, "multipart/form-data\n")
		_, _ = fmt.Fprintf(h, "len=%d", len(body))
		return hex.EncodeToString(h.Sum(nil))
	}
	_, _ = io.WriteString(h, normalizedContentType)
	_, _ = io.WriteString(h, "\n")
	_, _ = h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}
