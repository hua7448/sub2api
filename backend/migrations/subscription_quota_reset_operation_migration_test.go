package migrations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubscriptionQuotaResetOperationMigrationProvidesDurableIdempotencyAudit(t *testing.T) {
	content, err := FS.ReadFile("186_subscription_quota_reset_operations.sql")
	require.NoError(t, err)
	sql := strings.Join(strings.Fields(string(content)), " ")

	require.Contains(t, sql, "CREATE TABLE IF NOT EXISTS subscription_quota_reset_operations")
	require.Contains(t, sql, "actor_user_id BIGINT NOT NULL")
	require.Contains(t, sql, "idempotency_key_hash VARCHAR(64) NOT NULL")
	require.Contains(t, sql, "summary JSONB NOT NULL")
	require.Contains(t, sql, "target_pairs JSONB NOT NULL")
	require.Contains(t, sql, "ON subscription_quota_reset_operations (actor_user_id, idempotency_key_hash)")
	require.NotContains(t, sql, "idempotency_key VARCHAR")
}
