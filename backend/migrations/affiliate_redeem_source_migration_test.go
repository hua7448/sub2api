package migrations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAffiliateRedeemSourceMigrationTracksAndDeduplicatesRedeemRewards(t *testing.T) {
	sqlBytes, err := FS.ReadFile("185_affiliate_redeem_code_sources.sql")
	require.NoError(t, err)
	sql := strings.Join(strings.Fields(string(sqlBytes)), " ")

	require.Contains(t, sql, "ADD COLUMN IF NOT EXISTS source_redeem_code_id BIGINT")
	require.Contains(t, sql, "CREATE UNIQUE INDEX IF NOT EXISTS idx_user_affiliate_ledger_redeem_reward_unique")
	require.Contains(t, sql, "action IN ('accrue', 'accrue_days')")
	require.Contains(t, sql, "rc.used_by = ual.source_user_id")
}
