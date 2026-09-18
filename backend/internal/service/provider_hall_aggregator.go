package service

import (
	"context"
	"database/sql"
	"sort"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	providerHallAggregatorLockKey      = "provider-hall-aggregator"
	providerHallAggregatorTick         = time.Minute
	providerHallAggregatorRunTimeout   = 55 * time.Second
	providerHallAggregatorLag          = 45 * time.Second
	providerHallAggregatorMaxMinutes   = 240
	providerHallAggregatorReevalWindow = 65 * time.Minute
	providerHallAggregatorPruneEvery   = 10
	providerHallAggregatorLostAfter    = 45 * time.Second
	providerHallReconcileBatch         = 500

	ProviderHallRetentionRequests  = 7 * 24 * time.Hour
	ProviderHallRetentionMetrics   = 7 * 24 * time.Hour
	ProviderHallRetentionTier1     = 48 * time.Hour
	ProviderHallRetentionTier5     = 35 * 24 * time.Hour
	ProviderHallRetentionDirty     = 7 * 24 * time.Hour
	ProviderHallRetentionGaps      = 35 * 24 * time.Hour
	ProviderHallRetentionEpochs    = 35 * 24 * time.Hour
	ProviderHallRetentionJobs      = 90 * 24 * time.Hour
	ProviderHallReconcileMinAge    = 2 * time.Minute
	ProviderHallReconcileFailAfter = 10 * time.Minute
	// ProviderHallReconcileGiveUp turns rows that never found a ledger entry
	// into failed so retention can eventually release them.
	ProviderHallReconcileGiveUp = 24 * time.Hour
)

type ProviderHallAggregatorState struct {
	AlgorithmVersion int
	LastWindowEnd    *time.Time
	LastRunAt        *time.Time
	LastError        string
}

// ProviderHallNodeStatus is the latest epoch of one expected node.
type ProviderHallNodeStatus struct {
	NodeID           string
	EpochID          int64
	AlgorithmVersion int
	ConfirmedAt      *time.Time
	Exited           bool
	Present          bool
}

type ProviderHallDirtyBucket struct {
	ID        int64
	GroupID   int64
	ProfileID int64
	Minute    time.Time
}

// ProviderHallGapRange is a coverage gap projected onto the timeline.
type ProviderHallGapRange struct {
	Scope     string
	StartedAt time.Time
	EndedAt   *time.Time
}

type ProviderHallRecomputeInput struct {
	Now              time.Time
	T                time.Time
	Minutes          []time.Time
	NewWindowEnd     time.Time
	DirtyIDs         []int64
	ReevalFrom       time.Time
	ReevalTo         time.Time
	ExpectedNodes    []string
	Targets          []ProviderHallTargetRef
	AlgorithmVersion int
}

type ProviderHallRecomputeResult struct {
	MinutesRecomputed  int
	SnapshotsPublished int
	CoverageFlipped    int // minutes whose coverage changed during re-evaluation
}

type ProviderHallReconcileResult struct {
	Applied, Uncertain, Failed, Examined int
}

type ProviderHallPruneResult struct {
	Deleted map[string]int64
}

// ProviderHallAggregationRepository is implemented over raw *sql.DB in the
// repository package. Recompute performs steps 5–7 of the plan in one
// ReadCommitted transaction so a publish is atomic.
type ProviderHallAggregationRepository interface {
	GetState(ctx context.Context) (*ProviderHallAggregatorState, error)
	RecordRun(ctx context.Context, now time.Time, lastError string) error
	MarkLostEpochs(ctx context.Context, now, staleBefore time.Time) (int, error)
	ListDirty(ctx context.Context, limit int) ([]ProviderHallDirtyBucket, error)
	Recompute(ctx context.Context, in ProviderHallRecomputeInput) (*ProviderHallRecomputeResult, error)
	Reconcile(ctx context.Context, now time.Time, limit int) (*ProviderHallReconcileResult, error)
	Prune(ctx context.Context, now time.Time) (*ProviderHallPruneResult, error)
}

// ProviderHallAggregator rolls request facts into minute metrics and hourly
// snapshots on the primary instance. It runs once a minute under a leader
// lock; every step is idempotent so a crashed tick is simply redone.
type ProviderHallAggregator struct {
	repo      ProviderHallAggregationRepository
	facts     ProviderHallFactRepository
	cfgRepo   ProviderHallRepository
	db        *sql.DB
	lockCache LeaderLockCache
	cfg       *config.Config

	instanceID string
	now        func() time.Time
	acquire    func(ctx context.Context) (func(), bool)
	tick       time.Duration
	runTimeout time.Duration

	startOnce sync.Once
	stopOnce  sync.Once
	stopCh    chan struct{}
	kickCh    chan struct{}
	done      chan struct{}
	mu        sync.Mutex
	ticks     int
	lastRun   time.Time
	lastErr   string
	log       *zap.Logger
}

func NewProviderHallAggregator(repo ProviderHallAggregationRepository, facts ProviderHallFactRepository, cfgRepo ProviderHallRepository, db *sql.DB, lockCache LeaderLockCache, cfg *config.Config) *ProviderHallAggregator {
	a := &ProviderHallAggregator{
		repo:       repo,
		facts:      facts,
		cfgRepo:    cfgRepo,
		db:         db,
		lockCache:  lockCache,
		cfg:        cfg,
		instanceID: uuid.NewString(),
		now:        time.Now,
		tick:       providerHallAggregatorTick,
		runTimeout: providerHallAggregatorRunTimeout,
		stopCh:     make(chan struct{}),
		kickCh:     make(chan struct{}, 1),
		done:       make(chan struct{}),
		log:        logger.L().With(zap.String("component", "service.provider_hall_aggregator")),
	}
	a.acquire = func(ctx context.Context) (func(), bool) {
		return tryAcquireSingletonLeaderLock(ctx, a.lockCache, a.db, providerHallAggregatorLockKey, a.instanceID, 2*time.Minute)
	}
	return a
}

// Ready reports whether the aggregation runtime is implemented and wired.
func (a *ProviderHallAggregator) Ready() bool {
	return a != nil && a.repo != nil && a.facts != nil && a.cfgRepo != nil
}

func (a *ProviderHallAggregator) Start() {
	if a == nil || !a.Ready() {
		return
	}
	a.startOnce.Do(func() { go a.loop() })
}

func (a *ProviderHallAggregator) Stop() {
	if a == nil {
		return
	}
	a.stopOnce.Do(func() {
		close(a.stopCh)
		if a.Ready() {
			select {
			case <-a.done:
			case <-time.After(providerHallAggregatorRunTimeout + 5*time.Second):
				a.log.Warn("provider_hall.aggregator_stop_timeout")
			}
		}
	})
}

// Kick wakes the loop early (settings change, truncated dirty replay).
func (a *ProviderHallAggregator) Kick() {
	if a == nil {
		return
	}
	select {
	case a.kickCh <- struct{}{}:
	default:
	}
}

// Status is reported by the admin health endpoint.
func (a *ProviderHallAggregator) Status() (lastRun time.Time, lastErr string) {
	if a == nil {
		return time.Time{}, ""
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.lastRun, a.lastErr
}

func (a *ProviderHallAggregator) loop() {
	defer close(a.done)
	started := false
	for {
		if started {
			timer := time.NewTimer(a.tick)
			select {
			case <-timer.C:
			case <-a.kickCh:
				timer.Stop()
			case <-a.stopCh:
				timer.Stop()
				return
			}
		}
		started = true
		select {
		case <-a.stopCh:
			return
		default:
		}
		a.runOnce()
	}
}

// RunOnce executes one aggregation tick. Exported for tests and the e2e
// harness; production uses the loop.
func (a *ProviderHallAggregator) RunOnce() { a.runOnce() }

func (a *ProviderHallAggregator) runOnce() {
	defer func() {
		if r := recover(); r != nil {
			a.log.Error("provider_hall.aggregator_panic", zap.Any("panic", r))
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), a.runTimeout)
	defer cancel()
	release, ok := a.acquire(ctx)
	if !ok {
		return
	}
	if release != nil {
		defer release()
	}
	a.mu.Lock()
	a.ticks++
	tick := a.ticks
	a.mu.Unlock()
	kick, err := a.run(ctx, tick)
	now := a.now()
	msg := ""
	if err != nil {
		msg = err.Error()
		a.log.Warn("provider_hall.aggregator_run_failed", zap.Error(err))
	}
	a.mu.Lock()
	a.lastRun, a.lastErr = now, msg
	a.mu.Unlock()
	if rerr := a.repo.RecordRun(ctx, now, msg); rerr != nil {
		a.log.Warn("provider_hall.aggregator_record_run_failed", zap.Error(rerr))
	}
	if kick {
		a.Kick()
	}
}

func (a *ProviderHallAggregator) run(ctx context.Context, tick int) (kick bool, err error) {
	cfg, err := a.cfgRepo.GetConfig(ctx)
	if err != nil {
		return false, err
	}
	now := a.now().UTC()
	if !cfg.CollectionEnabled {
		if tick%providerHallAggregatorPruneEvery == 1 {
			_, err = a.repo.Prune(ctx, now)
		}
		return false, err
	}
	// Step 3: nodes that stopped heart-beating are lost; their tail is a gap.
	if _, err := a.repo.MarkLostEpochs(ctx, now, now.Add(-providerHallAggregatorLostAfter)); err != nil {
		return false, err
	}
	// Step 8 runs before the recompute so freshly reconciled bills are counted.
	if _, err := a.repo.Reconcile(ctx, now, providerHallReconcileBatch); err != nil {
		return false, err
	}
	state, err := a.repo.GetState(ctx)
	if err != nil {
		return false, err
	}
	T := now.Add(-providerHallAggregatorLag).Truncate(time.Minute)
	W := T
	if state.LastWindowEnd != nil && !state.LastWindowEnd.IsZero() {
		W = state.LastWindowEnd.UTC().Truncate(time.Minute)
	}
	// Never rebuild history older than the metrics retention.
	if floor := T.Add(-ProviderHallRetentionMetrics); W.Before(floor) {
		W = floor
	}
	dirty, err := a.repo.ListDirty(ctx, providerHallAggregatorMaxMinutes)
	if err != nil {
		return false, err
	}
	plan := providerHallPlanMinutes(W, T, dirty, providerHallAggregatorMaxMinutes)
	targets, err := a.facts.ListEnabledTargets(ctx)
	if err != nil {
		return false, err
	}
	in := ProviderHallRecomputeInput{
		Now:              now,
		T:                T,
		Minutes:          plan.Minutes,
		NewWindowEnd:     plan.NewWindowEnd,
		DirtyIDs:         plan.DirtyIDs,
		ReevalFrom:       T.Add(-providerHallAggregatorReevalWindow),
		ReevalTo:         T,
		ExpectedNodes:    cfg.ExpectedNodes,
		Targets:          targets,
		AlgorithmVersion: ProviderHallAlgorithmVersion,
	}
	if _, err := a.repo.Recompute(ctx, in); err != nil {
		return false, err
	}
	if tick%providerHallAggregatorPruneEvery == 1 {
		if _, err := a.repo.Prune(ctx, now); err != nil {
			return plan.Truncated, err
		}
	}
	return plan.Truncated, nil
}

type providerHallMinutePlan struct {
	Minutes      []time.Time
	NewWindowEnd time.Time
	DirtyIDs     []int64
	Truncated    bool
}

// providerHallPlanMinutes builds M = [W,T) ∪ dirty, capped at limit minutes.
// The watermark only advances over a contiguous prefix of [W,T) that was
// actually recomputed; dirty rows are consumed only when their minute is in M.
func providerHallPlanMinutes(W, T time.Time, dirty []ProviderHallDirtyBucket, limit int) providerHallMinutePlan {
	W, T = W.UTC().Truncate(time.Minute), T.UTC().Truncate(time.Minute)
	set := map[time.Time]bool{}
	for m := W; m.Before(T); m = m.Add(time.Minute) {
		set[m] = true
	}
	for _, d := range dirty {
		m := d.Minute.UTC().Truncate(time.Minute)
		if m.Before(T) {
			set[m] = true
		}
	}
	all := make([]time.Time, 0, len(set))
	for m := range set {
		all = append(all, m)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Before(all[j]) })
	plan := providerHallMinutePlan{}
	if limit > 0 && len(all) > limit {
		all = all[:limit]
		plan.Truncated = true
	}
	plan.Minutes = all
	chosen := map[time.Time]bool{}
	for _, m := range all {
		chosen[m] = true
	}
	newW := W
	for newW.Before(T) && chosen[newW] {
		newW = newW.Add(time.Minute)
	}
	plan.NewWindowEnd = newW
	for _, d := range dirty {
		m := d.Minute.UTC().Truncate(time.Minute)
		if chosen[m] {
			plan.DirtyIDs = append(plan.DirtyIDs, d.ID)
		}
	}
	return plan
}

// ProviderHallMinuteCoverage evaluates one minute [m, m+1m) against the
// collection gaps and expected node states, in the fixed priority
// collection_gap > node_unconfirmed > version_mismatch > complete.
func ProviderHallMinuteCoverage(m time.Time, gaps []ProviderHallGapRange, nodes []ProviderHallNodeStatus, algorithmVersion int) string {
	end := m.Add(time.Minute)
	for _, g := range gaps {
		if g.Scope != "collection" {
			continue
		}
		if g.StartedAt.Before(end) && (g.EndedAt == nil || g.EndedAt.After(m)) {
			return ProviderHallCoverageCollectionGap
		}
	}
	mismatch := false
	for _, n := range nodes {
		if !n.Present || n.ConfirmedAt == nil || n.ConfirmedAt.Before(end) {
			return ProviderHallCoverageNodeUnconfirmed
		}
		if n.AlgorithmVersion != algorithmVersion {
			mismatch = true
		}
	}
	if mismatch {
		return ProviderHallCoverageVersionMismatch
	}
	return ProviderHallCoverageComplete
}

// ProviderHallMinuteBillingCoverage: a billing gap dominates; otherwise any
// successful request still waiting for its bill makes the minute pending.
func ProviderHallMinuteBillingCoverage(m time.Time, gaps []ProviderHallGapRange, pending, uncertain int64) string {
	end := m.Add(time.Minute)
	for _, g := range gaps {
		if g.Scope != "billing" {
			continue
		}
		if g.StartedAt.Before(end) && (g.EndedAt == nil || g.EndedAt.After(m)) {
			return ProviderHallBillingCoverageGap
		}
	}
	if pending > 0 || uncertain > 0 {
		return ProviderHallBillingCoveragePending
	}
	return ProviderHallBillingCoverageComplete
}

// ProviderHallSnapshotTimes returns S = {T} ∪ {m+k, k∈[1,60]} for m ∈ M,
// limited to t ≤ T and t within the tier-5 retention.
func ProviderHallSnapshotTimes(T time.Time, minutes []time.Time, now time.Time) []time.Time {
	T = T.UTC().Truncate(time.Minute)
	floor := now.Add(-ProviderHallRetentionTier5)
	set := map[time.Time]bool{T: true}
	for _, m := range minutes {
		m = m.UTC().Truncate(time.Minute)
		for k := 1; k <= 60; k++ {
			t := m.Add(time.Duration(k) * time.Minute)
			if t.After(T) {
				break
			}
			if t.Before(floor) {
				continue
			}
			set[t] = true
		}
	}
	out := make([]time.Time, 0, len(set))
	for t := range set {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Before(out[j]) })
	return out
}
