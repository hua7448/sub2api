//go:build integration || providerhall_localdb

package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func init() { registerProviderHallLocalDBContract("b6_admin", providerHallAdminDatabaseContracts) }

// providerHallAdminDatabaseContracts covers the health reads: latest epoch per
// node, open gaps only, aggregator state, dirty count and reconciliation counts.
func providerHallAdminDatabaseContracts(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	repo := NewProviderHallAdminRepository(db)
	facts := NewProviderHallFactRepository(db)
	base := time.Date(2026, 9, 12, 5, 0, 0, 0, time.UTC)

	t.Run("admin_health_reads", func(t *testing.T) {
		old, err := facts.RegisterEpoch(ctx, "admin-node", "v0", 1, base.Add(-2*time.Hour))
		require.NoError(t, err)
		require.NoError(t, facts.MarkEpochExited(ctx, old, "shutdown", base.Add(-time.Hour)))
		current, err := facts.RegisterEpoch(ctx, "admin-node", "v1", 1, base)
		require.NoError(t, err)
		epochs, err := repo.ListLatestEpochs(ctx)
		require.NoError(t, err)
		var found *service.ProviderHallNodeEpoch
		for i := range epochs {
			if epochs[i].NodeID == "admin-node" {
				found = &epochs[i]
			}
		}
		require.NotNil(t, found)
		require.Equal(t, current, found.EpochID, "only the latest epoch per node is reported")
		require.Equal(t, "v1", found.BuildVersion)
		require.Nil(t, found.ExitedAt)

		openID, err := facts.OpenGap(ctx, service.ProviderHallGap{NodeID: "admin-node", EpochID: current, Scope: "collection", StartedAt: base, Reason: "queue_overflow"})
		require.NoError(t, err)
		closedID, err := facts.OpenGap(ctx, service.ProviderHallGap{NodeID: "admin-node", EpochID: current, Scope: "billing", StartedAt: base, Reason: "db_unavailable"})
		require.NoError(t, err)
		require.NoError(t, facts.CloseGap(ctx, closedID, base.Add(time.Minute)))
		gaps, err := repo.ListOpenGaps(ctx)
		require.NoError(t, err)
		var ids []int64
		for _, g := range gaps {
			if g.NodeID == "admin-node" {
				ids = append(ids, g.ID)
			}
		}
		require.Equal(t, []int64{openID}, ids, "closed gaps are excluded")

		state, err := repo.GetAggregatorState(ctx)
		require.NoError(t, err)
		require.Equal(t, 1, state.AlgorithmVersion)

		profile := int64(1)
		started := base.Add(-time.Hour)
		ended := started.Add(2 * time.Second)
		rows := make([]service.ProviderHallRequestRow, 0, 4)
		for _, status := range []string{"pending", "uncertain", "failed", "applied"} {
			row := service.ProviderHallRequestRow{TraceID: uuid.New(), NodeID: "admin-node", EpochID: current, GroupID: 990, ProfileID: &profile, APIKeyID: 1, Protocol: "responses",
				RequestedModel: "m", Source: service.ProviderHallSourceUser, StartedAt: started, EndedAt: &ended, Outcome: service.ProviderHallOutcomeSuccess, Submissions: 1, AlgorithmVersion: 1}
			rows = append(rows, row)
			_ = status
		}
		require.NoError(t, facts.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: current, Starts: rows, Finishes: rows, MaxSeq: 4}))
		for i, status := range []string{"pending", "uncertain", "failed", "applied"} {
			_, err := db.ExecContext(ctx, `UPDATE provider_hall_requests SET billing_status = $2, updated_at = $3 WHERE trace_id = $1::uuid`, rows[i].TraceID.String(), status, base)
			require.NoError(t, err)
		}
		dirty, err := repo.CountDirty(ctx)
		require.NoError(t, err)
		require.GreaterOrEqual(t, dirty, 1)
		rec, err := repo.CountReconciliation(ctx, base)
		require.NoError(t, err)
		require.GreaterOrEqual(t, rec.Pending, 1)
		require.GreaterOrEqual(t, rec.Uncertain, 1)
		require.GreaterOrEqual(t, rec.Failed24h, 1)
		stale, err := repo.CountReconciliation(ctx, base.Add(8*24*time.Hour))
		require.NoError(t, err)
		require.Zero(t, stale.Failed24h, "a week-old failure is outside the 24h window")
		require.NoError(t, facts.MarkEpochExited(ctx, current, "shutdown", base.Add(time.Hour)))
	})
}
