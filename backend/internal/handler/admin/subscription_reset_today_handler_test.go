//go:build unit

package admin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestResetTodayUsageRequiresExplicitConfirmation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewSubscriptionHandler(nil)
	router := gin.New()
	router.POST("/api/v1/admin/subscriptions/reset-today-usage", handler.ResetTodayUsage)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/subscriptions/reset-today-usage", bytes.NewBufferString(`{"confirm":false}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "reset-at-2026-08-01T12:50:00+08:00")
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "confirm must be true")
}

func TestResetTodayUsageRequiresIdempotencyKeyEvenInObserveOnlyMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewSubscriptionHandler(nil)
	router := gin.New()
	router.POST("/api/v1/admin/subscriptions/reset-today-usage", handler.ResetTodayUsage)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/subscriptions/reset-today-usage", bytes.NewBufferString(`{"confirm":true}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "IDEMPOTENCY_KEY_REQUIRED")
}

func TestResetTodayUsageSetsExplicitAuditActionBeforeValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewSubscriptionHandler(nil)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/subscriptions/reset-today-usage", bytes.NewBufferString(`{"confirm":false}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.ResetTodayUsage(ctx)

	action, exists := ctx.Get("audit_action")
	require.True(t, exists)
	require.Equal(t, service.AuditActionSubscriptionResetTodayUsage, action)
}
