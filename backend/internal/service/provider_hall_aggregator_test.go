//go:build unit

package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type providerHallAggRepoFake struct {
	mu         sync.Mutex
	state      ProviderHallAggregatorState
	dirty      []ProviderHallDirtyBucket
	recomputes []ProviderHallRecomputeInput
	reconciles int
	prunes     int
	lost       int
	runs       []string
}

func (f *providerHallAggRepoFake) GetState(context.Context) (*ProviderHallAggregatorState, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	st := f.state
	return &st, nil
}
func (f *providerHallAggRepoFake) RecordRun(_ context.Context, _ time.Time, msg string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.runs = append(f.runs, msg)
	return nil
}
func (f *providerHallAggRepoFake) MarkLostEpochs(context.Context, time.Time, time.Time) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lost++
	return 0, nil
}
func (f *providerHallAggRepoFake) ListDirty(_ context.Context, limit int) ([]ProviderHallDirtyBucket, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.dirty) > limit {
		return append([]ProviderHallDirtyBucket(nil), f.dirty[:limit]...), nil
	}
	return append([]ProviderHallDirtyBucket(nil), f.dirty...), nil
}
func (f *providerHallAggRepoFake) Recompute(_ context.Context, in ProviderHallRecomputeInput) (*ProviderHallRecomputeResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.recomputes = append(f.recomputes, in)
	f.state.LastWindowEnd = &in.NewWindowEnd
	remaining := f.dirty[:0]
	consumed := map[int64]bool{}
	for _, id := range in.DirtyIDs {
		consumed[id] = true
	}
	for _, d := range f.dirty {
		if !consumed[d.ID] {
			remaining = append(remaining, d)
		}
	}
	f.dirty = remaining
	return &ProviderHallRecomputeResult{MinutesRecomputed: len(in.Minutes)}, nil
}
func (f *providerHallAggRepoFake) Reconcile(context.Context, time.Time, int) (*ProviderHallReconcileResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.reconciles++
	return &ProviderHallReconcileResult{}, nil
}
func (f *providerHallAggRepoFake) Prune(context.Context, time.Time) (*ProviderHallPruneResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.prunes++
	return &ProviderHallPruneResult{}, nil
}

type providerHallCfgRepoStub struct {
	ProviderHallRepository
	cfg ProviderHallConfig
}

func (s providerHallCfgRepoStub) GetConfig(context.Context) (*ProviderHallConfig, error) {
	cfg := s.cfg
	return &cfg, nil
}

func newTestHallAggregator(t *testing.T, enabled bool, now time.Time) (*ProviderHallAggregator, *providerHallAggRepoFake) {
	t.Helper()
	repo := &providerHallAggRepoFake{}
	facts := &providerHallFactRepoFake{targets: []ProviderHallTargetRef{{TargetID: 1, GroupID: 10, ProfileID: 20, Enabled: true}}}
	cfg := providerHallCfgRepoStub{cfg: ProviderHallConfig{CollectionEnabled: enabled, ExpectedNodes: []string{"n1"}}}
	a := NewProviderHallAggregator(repo, facts, cfg, nil, nil, nil)
	a.now = func() time.Time { return now }
	a.acquire = func(context.Context) (func(), bool) { return func() {}, true }
	require.True(t, a.Ready())
	return a, repo
}

func TestProviderHallPlanMinutes(t *testing.T) {
	base := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	t.Run("first_run_has_no_backfill", func(t *testing.T) {
		p := providerHallPlanMinutes(base, base, nil, 240)
		require.Empty(t, p.Minutes)
		require.Equal(t, base, p.NewWindowEnd)
		require.False(t, p.Truncated)
	})
	t.Run("range_plus_dirty_dedup_and_watermark", func(t *testing.T) {
		dirty := []ProviderHallDirtyBucket{{ID: 1, Minute: base.Add(-30 * time.Minute)}, {ID: 2, Minute: base.Add(2 * time.Minute)}, {ID: 3, Minute: base.Add(10 * time.Minute)}}
		p := providerHallPlanMinutes(base, base.Add(3*time.Minute), dirty, 240)
		require.Equal(t, []time.Time{base.Add(-30 * time.Minute), base, base.Add(time.Minute), base.Add(2 * time.Minute)}, p.Minutes)
		require.Equal(t, base.Add(3*time.Minute), p.NewWindowEnd)
		require.Equal(t, []int64{1, 2}, p.DirtyIDs, "a dirty minute at or after T is left for a later tick")
		require.False(t, p.Truncated)
	})
	t.Run("truncation_keeps_watermark_contiguous", func(t *testing.T) {
		W := base.Add(-300 * time.Minute)
		dirty := []ProviderHallDirtyBucket{{ID: 9, Minute: base.Add(-400 * time.Minute)}}
		p := providerHallPlanMinutes(W, base, dirty, 240)
		require.True(t, p.Truncated)
		require.Len(t, p.Minutes, 240)
		require.Equal(t, base.Add(-400*time.Minute), p.Minutes[0])
		// 239 minutes of the range were covered: W .. W+239m, so the watermark
		// stops at W+239m and the next tick continues from there.
		require.Equal(t, W.Add(239*time.Minute), p.NewWindowEnd)
		require.Equal(t, []int64{9}, p.DirtyIDs)
	})
	t.Run("dirty_outside_prefix_does_not_advance", func(t *testing.T) {
		W := base.Add(-241 * time.Minute)
		p := providerHallPlanMinutes(W, base, nil, 240)
		require.True(t, p.Truncated)
		require.Equal(t, W.Add(240*time.Minute), p.NewWindowEnd)
	})
}

func TestProviderHallAggregatorRunComputesTAndW(t *testing.T) {
	now := time.Date(2026, 9, 12, 10, 3, 50, 0, time.UTC) // T = 10:03:05 → 10:03
	a, repo := newTestHallAggregator(t, true, now)
	a.RunOnce()
	require.Len(t, repo.recomputes, 1)
	in := repo.recomputes[0]
	T := time.Date(2026, 9, 12, 10, 3, 0, 0, time.UTC)
	require.Equal(t, T, in.T)
	require.Equal(t, T, in.NewWindowEnd, "first run starts at T without backfill")
	require.Empty(t, in.Minutes)
	require.Equal(t, T.Add(-65*time.Minute), in.ReevalFrom)
	require.Equal(t, []string{"n1"}, in.ExpectedNodes)
	require.Len(t, in.Targets, 1)
	require.Equal(t, 1, repo.reconciles)
	require.Equal(t, 1, repo.lost)
	require.Equal(t, 1, repo.prunes, "tick 1 prunes")
	require.Equal(t, []string{""}, repo.runs)

	// Second tick three minutes later: W advanced to the previous T.
	a.now = func() time.Time { return now.Add(3 * time.Minute) }
	a.RunOnce()
	require.Len(t, repo.recomputes, 2)
	in = repo.recomputes[1]
	require.Equal(t, []time.Time{T, T.Add(time.Minute), T.Add(2 * time.Minute)}, in.Minutes)
	require.Equal(t, T.Add(3*time.Minute), in.NewWindowEnd)
	require.Equal(t, 1, repo.prunes, "tick 2 does not prune")
}

func TestProviderHallAggregatorSkipsWithoutLock(t *testing.T) {
	now := time.Date(2026, 9, 12, 10, 3, 50, 0, time.UTC)
	a, repo := newTestHallAggregator(t, true, now)
	a.acquire = func(context.Context) (func(), bool) { return nil, false }
	a.RunOnce()
	require.Empty(t, repo.recomputes)
	require.Empty(t, repo.runs)
}

func TestProviderHallAggregatorDisabledOnlyPrunes(t *testing.T) {
	now := time.Date(2026, 9, 12, 10, 3, 50, 0, time.UTC)
	a, repo := newTestHallAggregator(t, false, now)
	for range 11 {
		a.RunOnce()
	}
	require.Empty(t, repo.recomputes)
	require.Zero(t, repo.reconciles)
	require.Equal(t, 2, repo.prunes, "ticks 1 and 11")
}

func TestProviderHallAggregatorTruncationKicks(t *testing.T) {
	now := time.Date(2026, 9, 12, 10, 3, 50, 0, time.UTC)
	a, repo := newTestHallAggregator(t, true, now)
	old := now.Add(-400 * time.Minute).Truncate(time.Minute)
	repo.state.LastWindowEnd = &old
	a.RunOnce()
	require.Len(t, repo.recomputes, 1)
	require.Len(t, repo.recomputes[0].Minutes, 240)
	select {
	case <-a.kickCh:
	default:
		t.Fatal("truncated replay must kick the loop")
	}
	a.RunOnce()
	require.Len(t, repo.recomputes, 2)
	require.Equal(t, old.Add(240*time.Minute), repo.recomputes[1].Minutes[0])
	require.Len(t, repo.recomputes[1].Minutes, 400-240)
	select {
	case <-a.kickCh:
		t.Fatal("complete replay must not kick")
	default:
	}
}

func TestProviderHallAggregatorStartStop(t *testing.T) {
	now := time.Date(2026, 9, 12, 10, 3, 50, 0, time.UTC)
	a, repo := newTestHallAggregator(t, true, now)
	a.tick = 20 * time.Millisecond
	a.Start()
	require.Eventually(t, func() bool {
		repo.mu.Lock()
		defer repo.mu.Unlock()
		return len(repo.recomputes) >= 2
	}, 2*time.Second, 5*time.Millisecond)
	a.Stop()
	a.Stop()
	last, errMsg := a.Status()
	require.Equal(t, now, last)
	require.Empty(t, errMsg)
}

func TestProviderHallMinuteCoverage(t *testing.T) {
	m := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	confirmed := m.Add(time.Minute)
	late := m.Add(30 * time.Second)
	closedBefore := m
	okNode := ProviderHallNodeStatus{NodeID: "n1", Present: true, ConfirmedAt: &confirmed, AlgorithmVersion: 1}
	cases := []struct {
		name  string
		gaps  []ProviderHallGapRange
		nodes []ProviderHallNodeStatus
		want  string
	}{
		{"complete", nil, []ProviderHallNodeStatus{okNode}, ProviderHallCoverageComplete},
		{"open_gap", []ProviderHallGapRange{{Scope: "collection", StartedAt: m.Add(-time.Hour)}}, []ProviderHallNodeStatus{okNode}, ProviderHallCoverageCollectionGap},
		{"gap_closed_at_minute_start_does_not_overlap", []ProviderHallGapRange{{Scope: "collection", StartedAt: m.Add(-time.Hour), EndedAt: &closedBefore}}, []ProviderHallNodeStatus{okNode}, ProviderHallCoverageComplete},
		{"billing_gap_is_not_collection", []ProviderHallGapRange{{Scope: "billing", StartedAt: m.Add(-time.Hour)}}, []ProviderHallNodeStatus{okNode}, ProviderHallCoverageComplete},
		{"gap_beats_unconfirmed", []ProviderHallGapRange{{Scope: "collection", StartedAt: m}}, []ProviderHallNodeStatus{{NodeID: "ghost"}}, ProviderHallCoverageCollectionGap},
		{"absent_node", nil, []ProviderHallNodeStatus{okNode, {NodeID: "ghost"}}, ProviderHallCoverageNodeUnconfirmed},
		{"late_watermark", nil, []ProviderHallNodeStatus{{NodeID: "n1", Present: true, ConfirmedAt: &late, AlgorithmVersion: 1}}, ProviderHallCoverageNodeUnconfirmed},
		{"unconfirmed_beats_version", nil, []ProviderHallNodeStatus{{NodeID: "n1", Present: true, ConfirmedAt: &confirmed, AlgorithmVersion: 2}, {NodeID: "ghost"}}, ProviderHallCoverageNodeUnconfirmed},
		{"version_mismatch", nil, []ProviderHallNodeStatus{{NodeID: "n1", Present: true, ConfirmedAt: &confirmed, AlgorithmVersion: 2}}, ProviderHallCoverageVersionMismatch},
		{"no_expected_nodes", nil, nil, ProviderHallCoverageComplete},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, ProviderHallMinuteCoverage(m, tc.gaps, tc.nodes, 1))
		})
	}
	require.Equal(t, ProviderHallBillingCoverageGap, ProviderHallMinuteBillingCoverage(m, []ProviderHallGapRange{{Scope: "billing", StartedAt: m}}, 0, 0))
	require.Equal(t, ProviderHallBillingCoveragePending, ProviderHallMinuteBillingCoverage(m, nil, 0, 1))
	require.Equal(t, ProviderHallBillingCoverageComplete, ProviderHallMinuteBillingCoverage(m, nil, 0, 0))
}

func TestProviderHallSnapshotTimes(t *testing.T) {
	T := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	times := ProviderHallSnapshotTimes(T, []time.Time{T.Add(-90 * time.Minute)}, T)
	require.Len(t, times, 61, "T plus the 60 windows that contain the minute")
	require.Equal(t, T.Add(-89*time.Minute), times[0])
	require.Equal(t, T, times[len(times)-1])
	times = ProviderHallSnapshotTimes(T, []time.Time{T.Add(-5 * time.Minute)}, T)
	require.Len(t, times, 5, "windows past T are not published yet")
	times = ProviderHallSnapshotTimes(T, []time.Time{T.Add(-40 * 24 * time.Hour)}, T)
	require.Len(t, times, 1, "windows beyond tier-5 retention are skipped")
}
