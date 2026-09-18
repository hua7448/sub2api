package service

import (
	"context"
	"time"
)

// Admin-only read models for the provider hall health endpoint (batch B6).
// They describe collector nodes, coverage gaps and reconciliation backlog; no
// key material, request bodies or per-user data appear here.

// ProviderHallNodeEpoch is the latest collector epoch of one node.
type ProviderHallNodeEpoch struct {
	NodeID           string
	EpochID          int64
	BuildVersion     string
	AlgorithmVersion int
	RegisteredAt     time.Time
	HeartbeatAt      time.Time
	ConfirmedAt      *time.Time
	PersistedSeq     int64
	Overflowed       bool
	ExitedAt         *time.Time
	ExitReason       string
}

// ProviderHallOpenGap is a coverage gap that has not been closed yet.
type ProviderHallOpenGap struct {
	ID        int64
	NodeID    string
	EpochID   int64
	Scope     string
	StartedAt time.Time
	Reason    string
}

// ProviderHallReconciliationCounts summarises request facts whose billing is
// still unresolved or failed recently.
type ProviderHallReconciliationCounts struct {
	Pending   int
	Uncertain int
	Failed24h int
}

// ProviderHallAdminRepository serves the admin health view with raw SQL reads.
type ProviderHallAdminRepository interface {
	ListLatestEpochs(ctx context.Context) ([]ProviderHallNodeEpoch, error)
	ListOpenGaps(ctx context.Context) ([]ProviderHallOpenGap, error)
	GetAggregatorState(ctx context.Context) (*ProviderHallAggregatorState, error)
	CountDirty(ctx context.Context) (int, error)
	CountReconciliation(ctx context.Context, now time.Time) (ProviderHallReconciliationCounts, error)
}

// ProviderHallNodeLostAfter mirrors the aggregator's missed-heartbeat rule for
// display purposes: a node that has not beaten in this long is shown as lost
// even before the aggregator closes its epoch.
const ProviderHallNodeLostAfter = 45 * time.Second
