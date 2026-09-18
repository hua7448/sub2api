//go:build unit

package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSetAuditExtraAllowsUserDailyQuotaResetSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)

	SetAuditExtra(ctx, map[string]any{
		"subscription_id": int64(501),
		"reset_at":        "2026-08-05T04:50:00Z",
		"ignored_secret":  "must-not-be-recorded",
	})

	raw, exists := ctx.Get(auditCtxKeyExtra)
	require.True(t, exists)
	extra, ok := raw.(map[string]any)
	require.True(t, ok)
	require.Equal(t, int64(501), extra["subscription_id"])
	require.Equal(t, "2026-08-05T04:50:00Z", extra["reset_at"])
	require.NotContains(t, extra, "ignored_secret")
}
