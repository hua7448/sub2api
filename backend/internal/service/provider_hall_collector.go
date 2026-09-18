package service

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// ProviderHallTargetRef is the collector's and runner's view of one enabled
// target joined with its profile.
type ProviderHallTargetRef struct {
	TargetID                    int64
	GroupID                     int64
	ProfileID                   int64
	Model                       string
	Protocol                    string
	Enabled                     bool
	ProbeKeyID                  *int64
	ProbeIntervalSeconds        int
	VerificationIntervalSeconds int
	TargetVersion               int64
	ProfileVersion              int64
	SupportsTools               bool
	OutputLimit                 int
	ModelAliases                []string
}

type ProviderHallGap struct {
	NodeID    string
	EpochID   int64
	Scope     string // collection | billing
	StartedAt time.Time
	Reason    string
}

// ProviderHallBillingEvent links the asynchronous billing outcome to a fact.
type ProviderHallBillingEvent struct {
	TraceID          uuid.UUID
	BillingRequestID string
	APIKeyID         int64
	Fingerprint      string
	Status           string // applied | duplicate | failed | uncertain | not_applicable
	Mode             string
	IsSubscription   bool
	Multiplier       float64
	InputBaseCost    decimal.Decimal
	TotalBaseCost    decimal.Decimal
	ActualCost       decimal.Decimal
	BilledAt         time.Time
}

type ProviderHallBillingSink interface {
	RecordBilling(ProviderHallBillingEvent)
}

// ProviderHallFactBatch is persisted atomically. ConfirmedAt is nil when the
// watermark must not advance (overflow, missing barrier).
type ProviderHallFactBatch struct {
	EpochID     int64
	Starts      []ProviderHallRequestRow
	Finishes    []ProviderHallRequestRow
	Billings    []ProviderHallBillingEvent
	MaxSeq      uint64
	ConfirmedAt *time.Time
	Overflowed  bool
}

type ProviderHallFactRepository interface {
	RegisterEpoch(ctx context.Context, nodeID, buildVersion string, algorithmVersion int, now time.Time) (int64, error)
	Heartbeat(ctx context.Context, epochID int64, now time.Time) error
	MarkEpochExited(ctx context.Context, epochID int64, reason string, now time.Time) error
	WriteBatch(ctx context.Context, batch ProviderHallFactBatch) error
	OpenGap(ctx context.Context, gap ProviderHallGap) (int64, error)
	CloseGap(ctx context.Context, id int64, endedAt time.Time) error
	ListEnabledTargets(ctx context.Context) ([]ProviderHallTargetRef, error)
	ListProbeKeys(ctx context.Context) (map[int64]time.Time, error)
}

type providerHallEventKind uint8

const (
	providerHallEventStart providerHallEventKind = iota + 1
	providerHallEventFinish
	providerHallEventBilling
	providerHallEventBarrier
)

type providerHallEvent struct {
	kind    providerHallEventKind
	seq     uint64
	at      time.Time
	row     ProviderHallRequestRow
	billing ProviderHallBillingEvent
}

type providerHallTargetKey struct {
	groupID  int64
	protocol string
	model    string
}

type providerHallTargetIndex struct {
	enabled   bool
	targets   map[providerHallTargetKey]ProviderHallTargetRef
	groups    map[int64]bool
	probeKeys map[int64]time.Time
}

const (
	providerHallQueueSize         = 8192
	providerHallBatchSize         = 500
	providerHallFlushInterval     = 250 * time.Millisecond
	providerHallBarrierInterval   = time.Second
	providerHallBarrierSlack      = time.Second
	providerHallHeartbeatInterval = 10 * time.Second
	providerHallRefreshInterval   = 30 * time.Second
	providerHallDBTimeout         = 10 * time.Second
	providerHallStopTimeout       = 5 * time.Second
)

// ProviderHallCollector turns tracker events into fact rows. It runs on every
// instance that serves gateway traffic. Any failure only degrades coverage
// bookkeeping; it never blocks or fails a real request.
type ProviderHallCollector struct {
	facts   ProviderHallFactRepository
	cfgRepo ProviderHallRepository
	nodeID  string
	build   string

	epochID    atomic.Int64
	queue      chan providerHallEvent
	seq        atomic.Uint64
	dropped    atomic.Uint64
	overflowed atomic.Bool
	index      atomic.Pointer[providerHallTargetIndex]

	openGapID          int64
	droppedAtLastFlush uint64
	lastConfirmed      time.Time
	registeredAt       time.Time

	startOnce sync.Once
	stopOnce  sync.Once
	stopCh    chan struct{}
	done      chan struct{}
	now       func() time.Time

	flushInterval     time.Duration
	barrierInterval   time.Duration
	heartbeatInterval time.Duration
	refreshInterval   time.Duration
	batchSize         int
	log               *zap.Logger
}

func NewProviderHallCollector(facts ProviderHallFactRepository, cfgRepo ProviderHallRepository, nodeID, buildVersion string) *ProviderHallCollector {
	return newProviderHallCollector(facts, cfgRepo, nodeID, buildVersion, providerHallQueueSize)
}

func newProviderHallCollector(facts ProviderHallFactRepository, cfgRepo ProviderHallRepository, nodeID, buildVersion string, queueSize int) *ProviderHallCollector {
	if queueSize < 1 {
		queueSize = providerHallQueueSize
	}
	return &ProviderHallCollector{
		facts:             facts,
		cfgRepo:           cfgRepo,
		nodeID:            nodeID,
		build:             buildVersion,
		queue:             make(chan providerHallEvent, queueSize),
		stopCh:            make(chan struct{}),
		done:              make(chan struct{}),
		now:               time.Now,
		flushInterval:     providerHallFlushInterval,
		barrierInterval:   providerHallBarrierInterval,
		heartbeatInterval: providerHallHeartbeatInterval,
		refreshInterval:   providerHallRefreshInterval,
		batchSize:         providerHallBatchSize,
		log:               logger.L().With(zap.String("component", "service.provider_hall_collector")),
	}
}

// ProviderHallNodeID identifies this process for expected_nodes matching.
func ProviderHallNodeID() string {
	if v := strings.TrimSpace(os.Getenv("INSTANCE_ID")); v != "" {
		return v
	}
	if host, err := os.Hostname(); err == nil && strings.TrimSpace(host) != "" {
		return strings.TrimSpace(host)
	}
	return "unknown"
}

func (c *ProviderHallCollector) clock() time.Time {
	if c == nil || c.now == nil {
		return time.Now()
	}
	return c.now()
}

func (c *ProviderHallCollector) NodeID() string {
	if c == nil {
		return ""
	}
	return c.nodeID
}

// QueueStats is reported by the admin health endpoint.
func (c *ProviderHallCollector) QueueStats() (depth, capacity int, dropped uint64, overflowed bool) {
	if c == nil {
		return 0, 0, 0, false
	}
	return len(c.queue), cap(c.queue), c.dropped.Load(), c.overflowed.Load()
}

func (c *ProviderHallCollector) Start() {
	if c == nil {
		return
	}
	c.startOnce.Do(func() { go c.run() })
}

func (c *ProviderHallCollector) Stop() {
	if c == nil {
		return
	}
	c.stopOnce.Do(func() {
		close(c.stopCh)
		select {
		case <-c.done:
		case <-time.After(providerHallStopTimeout + 2*time.Second):
			c.log.Warn("provider_hall.collector_stop_timeout")
		}
	})
}

// RefreshNow reloads the target index synchronously.
func (c *ProviderHallCollector) RefreshNow() {
	if c != nil {
		c.refresh()
	}
}

// FlushNow drains queued events synchronously. It is meant for tests and
// callers that did not Start the background writer; never call it while the
// writer goroutine runs.
func (c *ProviderHallCollector) FlushNow() {
	if c == nil {
		return
	}
	var pending []providerHallEvent
drain:
	for {
		select {
		case ev := <-c.queue:
			pending = append(pending, ev)
		default:
			break drain
		}
	}
	c.flush(pending)
}

// Begin registers a request as a candidate sample. It returns (nil, ctx) when
// collection is off or the request cannot belong to any enabled target.
func (c *ProviderHallCollector) Begin(ctx context.Context, in ProviderHallBeginInput) (*ProviderHallRequestTracker, context.Context) {
	if c == nil || ctx == nil {
		return nil, ctx
	}
	idx := c.index.Load()
	if idx == nil || !idx.enabled || in.GroupID < 1 || in.APIKeyID < 1 {
		return nil, ctx
	}
	if in.StartedAt.IsZero() {
		in.StartedAt = c.clock()
	}
	registeredAt, isProbeKey := idx.probeKeys[in.APIKeyID]
	isProbeKey = isProbeKey && !registeredAt.After(in.StartedAt)
	if !idx.groups[in.GroupID] && !isProbeKey {
		return nil, ctx
	}
	row := ProviderHallRequestRow{
		TraceID:          uuid.New(),
		NodeID:           c.nodeID,
		EpochID:          c.epochID.Load(),
		GroupID:          in.GroupID,
		APIKeyID:         in.APIKeyID,
		Protocol:         in.Protocol,
		RequestedModel:   in.RequestedModel,
		Source:           ProviderHallSourceUser,
		Stream:           in.Stream,
		StartedAt:        in.StartedAt,
		Outcome:          ProviderHallOutcomePending,
		AlgorithmVersion: ProviderHallAlgorithmVersion,
	}
	if target, ok := idx.targets[providerHallTargetKey{groupID: in.GroupID, protocol: in.Protocol, model: in.RequestedModel}]; ok {
		profileID := target.ProfileID
		row.ProfileID = &profileID
	}
	if header := strings.TrimSpace(in.TaskHeader); header != "" {
		ref, ok := ParseProviderHallTaskHeader(header)
		if ok && isProbeKey {
			row.Source = ref.Kind
			sampleID := ref.SampleID
			row.SampleID = &sampleID
			row.TraceID = ref.TraceID
		} else {
			row.ExclusionReason = ProviderHallExclusionTaskHeaderMismatch
		}
	} else if isProbeKey {
		row.Source = ProviderHallSourceProbe
	}
	t := &ProviderHallRequestTracker{sink: c, ctx: ctx, row: row}
	c.enqueue(providerHallEvent{kind: providerHallEventStart, row: row})
	return t, context.WithValue(ctx, providerHallTrackerCtxKey{}, t)
}

// RecordBilling implements ProviderHallBillingSink.
func (c *ProviderHallCollector) RecordBilling(ev ProviderHallBillingEvent) {
	if c == nil || ev.TraceID == uuid.Nil {
		return
	}
	if ev.BilledAt.IsZero() {
		ev.BilledAt = c.clock()
	}
	c.enqueue(providerHallEvent{kind: providerHallEventBilling, billing: ev})
}

func (c *ProviderHallCollector) enqueue(ev providerHallEvent) {
	if c == nil {
		return
	}
	ev.seq = c.seq.Add(1)
	if ev.at.IsZero() {
		ev.at = c.clock()
	}
	select {
	case c.queue <- ev:
	default:
		c.dropped.Add(1)
		c.overflowed.Store(true)
	}
}

func (c *ProviderHallCollector) dbContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), providerHallDBTimeout)
}

func (c *ProviderHallCollector) run() {
	defer close(c.done)
	defer func() {
		if r := recover(); r != nil {
			c.log.Error("provider_hall.collector_panic", zap.Any("panic", r))
		}
	}()
	if !c.register() {
		return
	}
	c.refresh()

	flushTicker := time.NewTicker(c.flushInterval)
	barrierTicker := time.NewTicker(c.barrierInterval)
	heartbeatTicker := time.NewTicker(c.heartbeatInterval)
	refreshTicker := time.NewTicker(c.refreshInterval)
	defer flushTicker.Stop()
	defer barrierTicker.Stop()
	defer heartbeatTicker.Stop()
	defer refreshTicker.Stop()

	var pending []providerHallEvent
	for {
		select {
		case ev := <-c.queue:
			pending = append(pending, ev)
			if len(pending) >= c.batchSize {
				pending = c.flush(pending)
			}
		case <-flushTicker.C:
			pending = c.flush(pending)
		case <-barrierTicker.C:
			c.enqueue(providerHallEvent{kind: providerHallEventBarrier, at: c.clock().Add(-providerHallBarrierSlack)})
		case <-heartbeatTicker.C:
			c.heartbeat()
		case <-refreshTicker.C:
			c.refresh()
		case <-c.stopCh:
			c.shutdown(pending)
			return
		}
	}
}

func (c *ProviderHallCollector) register() bool {
	backoff := time.Second
	for {
		ctx, cancel := c.dbContext()
		now := c.clock()
		id, err := c.facts.RegisterEpoch(ctx, c.nodeID, c.build, ProviderHallAlgorithmVersion, now)
		cancel()
		if err == nil {
			c.epochID.Store(id)
			c.registeredAt = now
			c.lastConfirmed = now
			return true
		}
		c.log.Warn("provider_hall.collector_register_failed", zap.Error(err))
		select {
		case <-c.stopCh:
			return false
		case <-time.After(backoff):
		}
		if backoff < 30*time.Second {
			backoff *= 2
		}
	}
}

func (c *ProviderHallCollector) heartbeat() {
	ctx, cancel := c.dbContext()
	defer cancel()
	if err := c.facts.Heartbeat(ctx, c.epochID.Load(), c.clock()); err != nil {
		c.log.Warn("provider_hall.collector_heartbeat_failed", zap.Error(err))
	}
}

// refresh reloads the target index. Errors keep the previous index so a
// transient database failure never flips collection on or off.
func (c *ProviderHallCollector) refresh() {
	ctx, cancel := c.dbContext()
	defer cancel()
	cfg, err := c.cfgRepo.GetConfig(ctx)
	if err != nil {
		c.log.Warn("provider_hall.collector_refresh_config_failed", zap.Error(err))
		return
	}
	idx := &providerHallTargetIndex{
		enabled:   cfg.CollectionEnabled,
		targets:   map[providerHallTargetKey]ProviderHallTargetRef{},
		groups:    map[int64]bool{},
		probeKeys: map[int64]time.Time{},
	}
	if cfg.CollectionEnabled {
		targets, err := c.facts.ListEnabledTargets(ctx)
		if err != nil {
			c.log.Warn("provider_hall.collector_refresh_targets_failed", zap.Error(err))
			return
		}
		for _, t := range targets {
			if !t.Enabled {
				continue
			}
			idx.targets[providerHallTargetKey{groupID: t.GroupID, protocol: t.Protocol, model: t.Model}] = t
			idx.groups[t.GroupID] = true
		}
		keys, err := c.facts.ListProbeKeys(ctx)
		if err != nil {
			c.log.Warn("provider_hall.collector_refresh_probe_keys_failed", zap.Error(err))
			return
		}
		idx.probeKeys = keys
	}
	c.index.Store(idx)
}

// flush persists pending events. It returns the events that must be retried.
func (c *ProviderHallCollector) flush(pending []providerHallEvent) []providerHallEvent {
	if len(pending) == 0 {
		c.maintainGap(nil)
		return pending
	}
	batch := ProviderHallFactBatch{EpochID: c.epochID.Load()}
	var barrier *time.Time
	for _, ev := range pending {
		if ev.seq > batch.MaxSeq {
			batch.MaxSeq = ev.seq
		}
		switch ev.kind {
		case providerHallEventStart:
			row := ev.row
			row.Seq = ev.seq
			batch.Starts = append(batch.Starts, row)
		case providerHallEventFinish:
			row := ev.row
			row.Seq = ev.seq
			batch.Finishes = append(batch.Finishes, row)
		case providerHallEventBilling:
			batch.Billings = append(batch.Billings, ev.billing)
		case providerHallEventBarrier:
			at := ev.at
			barrier = &at
		}
	}
	advance := c.maintainGap(barrier)
	if advance && barrier != nil {
		batch.ConfirmedAt = barrier
	}
	batch.Overflowed = c.overflowed.Load()

	ctx, cancel := c.dbContext()
	err := c.facts.WriteBatch(ctx, batch)
	cancel()
	if err != nil {
		c.log.Warn("provider_hall.collector_write_failed", zap.Error(err), zap.Int("events", len(pending)))
		if len(pending) > 4*c.batchSize {
			// Bound memory: shed the oldest half and mark the loss like a queue overflow.
			c.dropped.Add(uint64(len(pending) / 2))
			c.overflowed.Store(true)
			return append([]providerHallEvent(nil), pending[len(pending)/2:]...)
		}
		return pending
	}
	if batch.ConfirmedAt != nil {
		c.lastConfirmed = *batch.ConfirmedAt
	}
	return pending[:0]
}

// maintainGap opens a coverage gap when events were dropped and closes it once
// the queue has drained without further drops. It returns whether the
// watermark may advance on this flush.
func (c *ProviderHallCollector) maintainGap(barrier *time.Time) bool {
	dropped := c.dropped.Load()
	defer func() { c.droppedAtLastFlush = dropped }()
	if c.overflowed.Load() && c.openGapID == 0 {
		start := c.lastConfirmed
		if start.IsZero() {
			start = c.registeredAt
		}
		ctx, cancel := c.dbContext()
		id, err := c.facts.OpenGap(ctx, ProviderHallGap{NodeID: c.nodeID, EpochID: c.epochID.Load(), Scope: "collection", StartedAt: start, Reason: "queue_overflow"})
		cancel()
		if err != nil {
			c.log.Warn("provider_hall.collector_open_gap_failed", zap.Error(err))
			return false
		}
		c.openGapID = id
		c.log.Warn("provider_hall.collector_queue_overflow", zap.Int64("gap_id", id), zap.Uint64("dropped", dropped))
		return false
	}
	if c.openGapID != 0 {
		if barrier == nil || dropped != c.droppedAtLastFlush || len(c.queue) >= cap(c.queue)/2 {
			return false
		}
		ctx, cancel := c.dbContext()
		err := c.facts.CloseGap(ctx, c.openGapID, *barrier)
		cancel()
		if err != nil {
			c.log.Warn("provider_hall.collector_close_gap_failed", zap.Error(err))
			return false
		}
		c.log.Info("provider_hall.collector_gap_closed", zap.Int64("gap_id", c.openGapID))
		c.openGapID = 0
		c.overflowed.Store(false)
		return true
	}
	return !c.overflowed.Load()
}

func (c *ProviderHallCollector) shutdown(pending []providerHallEvent) {
	deadline := time.Now().Add(providerHallStopTimeout)
drain:
	for {
		select {
		case ev := <-c.queue:
			pending = append(pending, ev)
		default:
			break drain
		}
	}
	for len(pending) > 0 && time.Now().Before(deadline) {
		remaining := c.flush(pending)
		if len(remaining) == len(pending) {
			time.Sleep(200 * time.Millisecond)
		}
		pending = remaining
	}
	ctx, cancel := context.WithTimeout(context.Background(), providerHallStopTimeout)
	defer cancel()
	now := c.clock()
	if len(pending) > 0 {
		start := c.lastConfirmed
		if start.IsZero() {
			start = c.registeredAt
		}
		if _, err := c.facts.OpenGap(ctx, ProviderHallGap{NodeID: c.nodeID, EpochID: c.epochID.Load(), Scope: "collection", StartedAt: start, Reason: "node_exit"}); err != nil {
			c.log.Warn("provider_hall.collector_exit_gap_failed", zap.Error(err))
		}
	}
	if c.openGapID != 0 {
		if err := c.facts.CloseGap(ctx, c.openGapID, now); err != nil {
			c.log.Warn("provider_hall.collector_close_gap_failed", zap.Error(err))
		}
	}
	if err := c.facts.MarkEpochExited(ctx, c.epochID.Load(), "shutdown", now); err != nil && !errors.Is(err, context.DeadlineExceeded) {
		c.log.Warn("provider_hall.collector_exit_failed", zap.Error(err))
	}
}
