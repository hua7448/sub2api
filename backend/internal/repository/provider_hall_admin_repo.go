package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type providerHallAdminRepository struct{ db *sql.DB }

func NewProviderHallAdminRepository(db *sql.DB) service.ProviderHallAdminRepository {
	return &providerHallAdminRepository{db: db}
}

func (r *providerHallAdminRepository) ListLatestEpochs(ctx context.Context) ([]service.ProviderHallNodeEpoch, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT DISTINCT ON (node_id) node_id, id, build_version, algorithm_version, registered_at, heartbeat_at, confirmed_at, persisted_seq, overflowed, exited_at, exit_reason
FROM provider_hall_collector_epochs
ORDER BY node_id, registered_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallNodeEpoch
	for rows.Next() {
		var e service.ProviderHallNodeEpoch
		var confirmed, exited sql.NullTime
		if err := rows.Scan(&e.NodeID, &e.EpochID, &e.BuildVersion, &e.AlgorithmVersion, &e.RegisteredAt, &e.HeartbeatAt, &confirmed, &e.PersistedSeq, &e.Overflowed, &exited, &e.ExitReason); err != nil {
			return nil, err
		}
		e.RegisteredAt, e.HeartbeatAt = e.RegisteredAt.UTC(), e.HeartbeatAt.UTC()
		if confirmed.Valid {
			t := confirmed.Time.UTC()
			e.ConfirmedAt = &t
		}
		if exited.Valid {
			t := exited.Time.UTC()
			e.ExitedAt = &t
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *providerHallAdminRepository) ListOpenGaps(ctx context.Context) ([]service.ProviderHallOpenGap, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, node_id, epoch_id, scope, started_at, reason FROM provider_hall_coverage_gaps WHERE ended_at IS NULL ORDER BY started_at, id LIMIT 200`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallOpenGap
	for rows.Next() {
		var g service.ProviderHallOpenGap
		if err := rows.Scan(&g.ID, &g.NodeID, &g.EpochID, &g.Scope, &g.StartedAt, &g.Reason); err != nil {
			return nil, err
		}
		g.StartedAt = g.StartedAt.UTC()
		out = append(out, g)
	}
	return out, rows.Err()
}

func (r *providerHallAdminRepository) GetAggregatorState(ctx context.Context) (*service.ProviderHallAggregatorState, error) {
	st := &service.ProviderHallAggregatorState{}
	var lastWindow, lastRun sql.NullTime
	if err := r.db.QueryRowContext(ctx, `SELECT algorithm_version, last_window_end, last_run_at, last_error FROM provider_hall_aggregator_state WHERE id = 1`).
		Scan(&st.AlgorithmVersion, &lastWindow, &lastRun, &st.LastError); err != nil {
		return nil, err
	}
	if lastWindow.Valid {
		t := lastWindow.Time.UTC()
		st.LastWindowEnd = &t
	}
	if lastRun.Valid {
		t := lastRun.Time.UTC()
		st.LastRunAt = &t
	}
	return st, nil
}

func (r *providerHallAdminRepository) CountDirty(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_dirty_buckets`).Scan(&n)
	return n, err
}

func (r *providerHallAdminRepository) CountReconciliation(ctx context.Context, now time.Time) (service.ProviderHallReconciliationCounts, error) {
	var c service.ProviderHallReconciliationCounts
	err := r.db.QueryRowContext(ctx, `
SELECT count(*) FILTER (WHERE billing_status = 'pending' AND outcome = 'success'),
       count(*) FILTER (WHERE billing_status = 'uncertain' AND outcome = 'success'),
       count(*) FILTER (WHERE billing_status = 'failed' AND updated_at >= $1)
FROM provider_hall_requests
WHERE started_at >= $2`, now.Add(-24*time.Hour).UTC(), now.Add(-7*24*time.Hour).UTC()).Scan(&c.Pending, &c.Uncertain, &c.Failed24h)
	return c, err
}
