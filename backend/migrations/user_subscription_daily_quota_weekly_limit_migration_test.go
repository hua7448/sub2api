package migrations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserSubscriptionDailyQuotaWeeklyLimitMigration(t *testing.T) {
	content, err := FS.ReadFile("194_user_subscription_daily_quota_weekly_limit.sql")
	require.NoError(t, err)

	sql := strings.ToUpper(string(content))
	require.Contains(t, sql, "ALTER TABLE USER_SUBSCRIPTIONS")
	require.Contains(t, sql, "ADD COLUMN IF NOT EXISTS DAILY_QUOTA_RESET_WEEK_START TIMESTAMPTZ")
	require.Contains(t, sql, "ADD COLUMN IF NOT EXISTS DAILY_QUOTA_RESET_OPERATIONS JSONB NOT NULL DEFAULT '[]'::JSONB")
	require.NotContains(t, sql, "IDEMPOTENCY_KEY VARCHAR")
	require.NotContains(t, sql, "IDEMPOTENCY_KEY TEXT")
	require.NotContains(t, sql, "DAILY_QUOTA_RESET_WEEK_START TIMESTAMPTZ NOT NULL")
}
