//go:build unit

package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSetAuditExtraAllowsResetTodayOperationSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	fields := map[string]any{
		"cutoff_at":                   "2026-08-01T04:50:00Z",
		"captured_at":                 "2026-08-01T04:50:01Z",
		"subscriptions":               669,
		"today_window_subscriptions":  640,
		"reset_subscriptions":         620,
		"reset_and_deducted_usd":      7290.01,
		"stale_daily_cleared_usd":     12.34,
		"cache_invalidations":         669,
		"cache_invalidation_failures": 0,
	}

	SetAuditExtra(ctx, fields)

	value, exists := ctx.Get(auditCtxKeyExtra)
	require.True(t, exists)
	require.Equal(t, fields, value)
}
