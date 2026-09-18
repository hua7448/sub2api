//go:build unit

package service

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

type providerHallFactRepoFake struct {
	mu         sync.Mutex
	batches    []ProviderHallFactBatch
	gaps       []ProviderHallGap
	closed     map[int64]time.Time
	exited     []string
	heartbeats int
	registered int
	targets    []ProviderHallTargetRef
	probeKeys  map[int64]time.Time
	writeErr   error
	nextGap    int64
}

func (f *providerHallFactRepoFake) RegisterEpoch(context.Context, string, string, int, time.Time) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.registered++
	return 7, nil
}
func (f *providerHallFactRepoFake) Heartbeat(context.Context, int64, time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.heartbeats++
	return nil
}
func (f *providerHallFactRepoFake) MarkEpochExited(_ context.Context, _ int64, reason string, _ time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.exited = append(f.exited, reason)
	return nil
}
func (f *providerHallFactRepoFake) WriteBatch(_ context.Context, b ProviderHallFactBatch) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.writeErr != nil {
		return f.writeErr
	}
	f.batches = append(f.batches, b)
	return nil
}
func (f *providerHallFactRepoFake) OpenGap(_ context.Context, g ProviderHallGap) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nextGap++
	f.gaps = append(f.gaps, g)
	return f.nextGap, nil
}
func (f *providerHallFactRepoFake) CloseGap(_ context.Context, id int64, at time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed == nil {
		f.closed = map[int64]time.Time{}
	}
	f.closed[id] = at
	return nil
}
func (f *providerHallFactRepoFake) ListEnabledTargets(context.Context) ([]ProviderHallTargetRef, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]ProviderHallTargetRef(nil), f.targets...), nil
}
func (f *providerHallFactRepoFake) ListProbeKeys(context.Context) (map[int64]time.Time, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := map[int64]time.Time{}
	for k, v := range f.probeKeys {
		out[k] = v
	}
	return out, nil
}
func (f *providerHallFactRepoFake) snapshot() ([]ProviderHallFactBatch, []ProviderHallGap, map[int64]time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	closed := map[int64]time.Time{}
	for k, v := range f.closed {
		closed[k] = v
	}
	return append([]ProviderHallFactBatch(nil), f.batches...), append([]ProviderHallGap(nil), f.gaps...), closed
}

var providerHallTestEpoch = time.Date(2026, 9, 11, 8, 0, 0, 0, time.UTC)

func newTestHallCollector(t *testing.T, queue int, enabled bool) (*ProviderHallCollector, *providerHallFactRepoFake) {
	t.Helper()
	facts := &providerHallFactRepoFake{
		targets: []ProviderHallTargetRef{
			{TargetID: 1, GroupID: 10, ProfileID: 100, Model: "gpt-test", Protocol: "responses", Enabled: true},
			{TargetID: 2, GroupID: 10, ProfileID: 101, Model: "gpt-test", Protocol: "chat_completions", Enabled: true},
		},
		probeKeys: map[int64]time.Time{99: providerHallTestEpoch.Add(-time.Hour)},
	}
	cfg := &providerHallRepoStub{config: ProviderHallConfig{CollectionEnabled: enabled}}
	c := newProviderHallCollector(facts, cfg, "node-a", "test", queue)
	c.now = func() time.Time { return providerHallTestEpoch }
	c.epochID.Store(7)
	c.registeredAt = providerHallTestEpoch.Add(-time.Minute)
	c.lastConfirmed = c.registeredAt
	c.RefreshNow()
	return c, facts
}

func hallBegin(c *ProviderHallCollector, in ProviderHallBeginInput) (*ProviderHallRequestTracker, context.Context) {
	if in.StartedAt.IsZero() {
		in.StartedAt = providerHallTestEpoch
	}
	return c.Begin(context.Background(), in)
}

func TestProviderHallTaskHeaderRoundTrip(t *testing.T) {
	ref := ProviderHallTaskRef{Kind: ProviderHallSourceVerification, JobID: 12, SampleID: 34, TraceID: uuid.New()}
	parsed, ok := ParseProviderHallTaskHeader(FormatProviderHallTaskHeader(ref))
	require.True(t, ok)
	require.Equal(t, ref, parsed)
	for _, bad := range []string{"", "probe:1:2", "user:1:2:" + uuid.NewString(), "probe:0:2:" + uuid.NewString(), "probe:1:2:not-a-uuid", "probe:x:2:" + uuid.NewString()} {
		_, ok := ParseProviderHallTaskHeader(bad)
		require.False(t, ok, bad)
	}
}

func TestProviderHallCollectorBeginFiltering(t *testing.T) {
	c, _ := newTestHallCollector(t, 64, true)
	base := ProviderHallBeginInput{APIKeyID: 5, GroupID: 10, Protocol: "responses", RequestedModel: "gpt-test", Stream: true}

	tracker, ctx := hallBegin(c, base)
	require.NotNil(t, tracker)
	require.Same(t, tracker, ProviderHallTrackerFromContext(ctx))
	require.Equal(t, ProviderHallSourceUser, tracker.Source())
	require.NotNil(t, tracker.row.ProfileID)
	require.Equal(t, int64(100), *tracker.row.ProfileID)

	other := base
	other.RequestedModel = "gpt-other"
	tracker, _ = hallBegin(c, other)
	require.NotNil(t, tracker, "unknown model on a hall group is still recorded for diagnostics")
	require.Nil(t, tracker.row.ProfileID)

	unknownGroup := base
	unknownGroup.GroupID = 11
	tracker, _ = hallBegin(c, unknownGroup)
	require.Nil(t, tracker, "groups without enabled targets are not collected")

	probe := unknownGroup
	probe.APIKeyID = 99
	tracker, _ = hallBegin(c, probe)
	require.NotNil(t, tracker, "probe keys are always recorded so spend can be reconciled")
	require.Equal(t, ProviderHallSourceProbe, tracker.Source())

	tooEarly := probe
	tooEarly.StartedAt = providerHallTestEpoch.Add(-2 * time.Hour)
	tracker, _ = hallBegin(c, tooEarly)
	require.Nil(t, tracker, "traffic before key registration is normal traffic")

	ref := ProviderHallTaskRef{Kind: ProviderHallSourceVerification, JobID: 1, SampleID: 2, TraceID: uuid.New()}
	task := base
	task.APIKeyID = 99
	task.TaskHeader = FormatProviderHallTaskHeader(ref)
	tracker, _ = hallBegin(c, task)
	require.NotNil(t, tracker)
	require.Equal(t, ProviderHallSourceVerification, tracker.Source())
	require.Equal(t, ref.TraceID, tracker.TraceID())
	require.Equal(t, int64(2), *tracker.row.SampleID)

	forged := base
	forged.TaskHeader = task.TaskHeader
	tracker, _ = hallBegin(c, forged)
	require.NotNil(t, tracker)
	require.Equal(t, ProviderHallSourceUser, tracker.Source())
	require.NotEqual(t, ref.TraceID, tracker.TraceID())
	tracker.RecordSubmission()
	tracker.Attempt(providerHallTestEpoch, &OpenAIForwardResult{}, nil)
	tracker.Finish()
	c.FlushNow()
	batches, _, _ := c.facts.(*providerHallFactRepoFake).snapshot()
	last := batches[len(batches)-1].Finishes
	require.Len(t, last, 1)
	require.Equal(t, ProviderHallOutcomeExcluded, last[0].Outcome)
	require.Equal(t, ProviderHallExclusionTaskHeaderMismatch, last[0].ExclusionReason)

	disabled, _ := newTestHallCollector(t, 64, false)
	tracker, _ = hallBegin(disabled, base)
	require.Nil(t, tracker)
	var nilCollector *ProviderHallCollector
	tracker, ctx = nilCollector.Begin(context.Background(), base)
	require.Nil(t, tracker)
	require.NotNil(t, ctx)
}

func TestProviderHallTrackerOutcomes(t *testing.T) {
	failover := &UpstreamFailoverError{StatusCode: 429}
	ttft := 120
	for _, tc := range []struct {
		name      string
		drive     func(*ProviderHallRequestTracker, context.CancelFunc)
		outcome   ProviderHallOutcome
		exclusion ProviderHallExclusion
		subs      int
		ttft      *int
		usage     bool
	}{
		{"success_first_attempt", func(tr *ProviderHallRequestTracker, _ context.CancelFunc) {
			tr.RecordSubmission()
			tr.Attempt(providerHallTestEpoch.Add(50*time.Millisecond), &OpenAIForwardResult{FirstTokenMs: &ttft, Usage: OpenAIUsage{InputTokens: 10, OutputTokens: 2, CacheReadInputTokens: 4}}, nil)
		}, ProviderHallOutcomeSuccess, ProviderHallExclusionNone, 1, intPtr(170), true},
		{"retry_then_success_counts_both", func(tr *ProviderHallRequestTracker, _ context.CancelFunc) {
			tr.RecordSubmission()
			tr.Attempt(providerHallTestEpoch, nil, failover)
			tr.RecordSubmission()
			tr.Attempt(providerHallTestEpoch, &OpenAIForwardResult{}, nil)
		}, ProviderHallOutcomeSuccess, ProviderHallExclusionNone, 2, nil, true},
		{"all_attempts_failed", func(tr *ProviderHallRequestTracker, _ context.CancelFunc) {
			tr.RecordSubmission()
			tr.Attempt(providerHallTestEpoch, nil, failover)
			tr.RecordSubmission()
			tr.Attempt(providerHallTestEpoch, nil, failover)
		}, ProviderHallOutcomeFailed, ProviderHallExclusionNone, 2, nil, false},
		{"stream_interrupted_after_200", func(tr *ProviderHallRequestTracker, _ context.CancelFunc) {
			tr.RecordSubmission()
			tr.Attempt(providerHallTestEpoch, &OpenAIForwardResult{FirstTokenMs: &ttft, Usage: OpenAIUsage{InputTokens: 3}}, errors.New("stream usage incomplete: missing terminal event"))
		}, ProviderHallOutcomeFailed, ProviderHallExclusionNone, 1, intPtr(120), true},
		{"no_account", func(tr *ProviderHallRequestTracker, _ context.CancelFunc) { tr.NoAccount() },
			ProviderHallOutcomeFailed, ProviderHallExclusionNone, 0, nil, false},
		{"client_cancelled", func(tr *ProviderHallRequestTracker, cancel context.CancelFunc) {
			tr.RecordSubmission()
			cancel()
			tr.Attempt(providerHallTestEpoch, &OpenAIForwardResult{ClientDisconnect: true}, errors.New("client gone"))
		}, ProviderHallOutcomeExcluded, ProviderHallExclusionClientCancel, 1, nil, true},
		{"policy_reject", func(tr *ProviderHallRequestTracker, _ context.CancelFunc) {
			tr.Reject(ProviderHallExclusionPolicyReject)
		},
			ProviderHallOutcomeExcluded, ProviderHallExclusionPolicyReject, 0, nil, false},
		{"early_exit", func(*ProviderHallRequestTracker, context.CancelFunc) {},
			ProviderHallOutcomeExcluded, ProviderHallExclusionEarlyExit, 0, nil, false},
		{"success_wins_over_late_cancel", func(tr *ProviderHallRequestTracker, cancel context.CancelFunc) {
			tr.RecordSubmission()
			tr.Attempt(providerHallTestEpoch, &OpenAIForwardResult{}, nil)
			cancel()
		}, ProviderHallOutcomeSuccess, ProviderHallExclusionNone, 1, nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, facts := newTestHallCollector(t, 64, true)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			tracker, _ := c.Begin(ctx, ProviderHallBeginInput{APIKeyID: 5, GroupID: 10, Protocol: "responses", RequestedModel: "gpt-test", StartedAt: providerHallTestEpoch})
			require.NotNil(t, tracker)
			tc.drive(tracker, cancel)
			tracker.Finish()
			tracker.Finish()
			tracker.RecordSubmission() // ignored after finish
			c.FlushNow()
			batches, _, _ := facts.snapshot()
			require.Len(t, batches, 1)
			require.Len(t, batches[0].Starts, 1)
			require.Len(t, batches[0].Finishes, 1)
			row := batches[0].Finishes[0]
			require.Equal(t, tc.outcome, row.Outcome)
			require.Equal(t, tc.exclusion, row.ExclusionReason)
			require.Equal(t, tc.subs, row.Submissions)
			require.Equal(t, tc.usage, row.UsageKnown)
			require.NotNil(t, row.EndedAt)
			if tc.ttft == nil {
				require.Nil(t, row.TTFTMs)
			} else {
				require.NotNil(t, row.TTFTMs)
				require.Equal(t, *tc.ttft, *row.TTFTMs)
			}
		})
	}
}

type providerHallCountingUpstream struct {
	calls atomic.Int32
}

func (u *providerHallCountingUpstream) Do(*http.Request, string, int64, int) (*http.Response, error) {
	u.calls.Add(1)
	return nil, errors.New("not sent")
}
func (u *providerHallCountingUpstream) DoWithTLS(*http.Request, string, int64, int, *tlsfingerprint.Profile) (*http.Response, error) {
	u.calls.Add(1)
	return nil, errors.New("not sent")
}

func TestProviderHallSubmissionsCountedAtUpstreamBoundary(t *testing.T) {
	c, facts := newTestHallCollector(t, 64, true)
	tracker, ctx := hallBegin(c, ProviderHallBeginInput{APIKeyID: 5, GroupID: 10, Protocol: "responses", RequestedModel: "gpt-test"})
	base := &providerHallCountingUpstream{}
	upstream := WrapProviderHallUpstream(base)
	require.Equal(t, upstream, WrapProviderHallUpstream(upstream), "wrapping is idempotent")
	detached := context.WithoutCancel(ctx)
	for _, target := range []struct {
		method, url string
	}{
		{http.MethodPost, "https://api.openai.com/v1/responses"},
		{http.MethodPost, "https://chatgpt.com/backend-api/codex/responses"},
		{http.MethodPost, "https://api.openai.com/v1/chat/completions/"},
		{http.MethodPost, "https://api.anthropic.com/v1/messages"},
		{http.MethodPost, "https://auth.openai.com/oauth/token"},
		{http.MethodGet, "https://api.openai.com/v1/models"},
		{http.MethodPost, "https://api.openai.com/v1/responses/compact"},
	} {
		req, err := http.NewRequestWithContext(detached, target.method, target.url, nil)
		require.NoError(t, err)
		_, _ = upstream.Do(req, "", 1, 1)
		_, _ = upstream.DoWithTLS(req, "", 1, 1, nil)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://api.openai.com/v1/responses", nil)
	require.NoError(t, err)
	_, _ = upstream.Do(req, "", 1, 1)
	require.Equal(t, int32(15), base.calls.Load(), "every call reaches the real upstream")
	tracker.Attempt(providerHallTestEpoch, &OpenAIForwardResult{}, nil)
	tracker.Finish()
	c.FlushNow()
	batches, _, _ := facts.snapshot()
	require.Equal(t, 8, batches[0].Finishes[0].Submissions, "four model endpoints x two transports; token, models and compact are not submissions")
}

func TestProviderHallTrackerNilSafe(t *testing.T) {
	var tracker *ProviderHallRequestTracker
	tracker.RecordSubmission()
	tracker.SelectAccount(&Account{Platform: PlatformOpenAI})
	tracker.Attempt(time.Now(), nil, nil)
	tracker.Reject(ProviderHallExclusionPolicyReject)
	tracker.NoAccount()
	tracker.Finish()
	require.Equal(t, uuid.Nil, tracker.TraceID())
	require.Empty(t, tracker.Source())
	require.Nil(t, ProviderHallTrackerFromContext(context.Background()))
	require.Nil(t, ProviderHallTrackerFromContext(nil))
	var c *ProviderHallCollector
	c.RecordBilling(ProviderHallBillingEvent{TraceID: uuid.New()})
	c.Start()
	c.Stop()
	c.RefreshNow()
	c.FlushNow()
	depth, capacity, dropped, overflowed := c.QueueStats()
	require.Zero(t, depth+capacity)
	require.Zero(t, dropped)
	require.False(t, overflowed)
}

func TestProviderHallCollectorBillingBeforeAndAfterFinish(t *testing.T) {
	c, facts := newTestHallCollector(t, 64, true)
	tracker, _ := hallBegin(c, ProviderHallBeginInput{APIKeyID: 5, GroupID: 10, Protocol: "responses", RequestedModel: "gpt-test"})
	tracker.SelectAccount(&Account{Platform: PlatformOpenAI})
	tracker.RecordSubmission()
	tracker.Attempt(providerHallTestEpoch, &OpenAIForwardResult{UpstreamModel: "gpt-test-2026", Usage: OpenAIUsage{InputTokens: 12}}, nil)
	// Billing can land on the worker pool before the handler returns.
	c.RecordBilling(ProviderHallBillingEvent{TraceID: tracker.TraceID(), BillingRequestID: "client:abc", APIKeyID: 5, Status: "applied", Mode: "token", ActualCost: decimal.RequireFromString("0.5")})
	tracker.Finish()
	c.RecordBilling(ProviderHallBillingEvent{TraceID: uuid.Nil}) // ignored
	c.enqueue(providerHallEvent{kind: providerHallEventBarrier, at: providerHallTestEpoch.Add(-time.Second)})
	c.FlushNow()
	batches, gaps, _ := facts.snapshot()
	require.Empty(t, gaps)
	require.Len(t, batches, 1)
	b := batches[0]
	require.Equal(t, int64(7), b.EpochID)
	require.Len(t, b.Starts, 1)
	require.Len(t, b.Billings, 1)
	require.Len(t, b.Finishes, 1)
	require.Equal(t, PlatformOpenAI, b.Finishes[0].Platform)
	require.Equal(t, "gpt-test-2026", b.Finishes[0].UpstreamModel)
	require.Equal(t, "client:abc", b.Billings[0].BillingRequestID)
	require.NotNil(t, b.ConfirmedAt)
	require.Equal(t, providerHallTestEpoch.Add(-time.Second), *b.ConfirmedAt)
	require.Greater(t, b.MaxSeq, uint64(0))
	require.False(t, b.Overflowed)
	require.Equal(t, providerHallTestEpoch.Add(-time.Second), c.lastConfirmed)
}

func TestProviderHallCollectorOverflowOpensAndClosesGap(t *testing.T) {
	c, facts := newTestHallCollector(t, 4, true)
	var trackers []*ProviderHallRequestTracker
	for i := 0; i < 6; i++ {
		tracker, _ := hallBegin(c, ProviderHallBeginInput{APIKeyID: 5, GroupID: 10, Protocol: "responses", RequestedModel: "gpt-test"})
		require.NotNil(t, tracker)
		trackers = append(trackers, tracker)
	}
	_, capacity, dropped, overflowed := c.QueueStats()
	require.Equal(t, 4, capacity)
	require.Equal(t, uint64(2), dropped)
	require.True(t, overflowed)

	c.enqueue(providerHallEvent{kind: providerHallEventBarrier, at: providerHallTestEpoch}) // dropped too: queue is full
	c.FlushNow()
	batches, gaps, closed := facts.snapshot()
	require.Len(t, gaps, 1)
	require.Equal(t, "queue_overflow", gaps[0].Reason)
	require.Equal(t, "collection", gaps[0].Scope)
	require.Equal(t, c.registeredAt, gaps[0].StartedAt)
	require.Empty(t, closed)
	require.Len(t, batches, 1)
	require.Nil(t, batches[0].ConfirmedAt, "watermark must not advance while events were lost")
	require.True(t, batches[0].Overflowed)

	// Drain without new drops: the next barrier closes the gap and resumes.
	for _, tracker := range trackers[:1] {
		tracker.Finish()
	}
	c.FlushNow()
	_, _, closed = facts.snapshot()
	require.Empty(t, closed, "a flush without a barrier cannot close the gap")
	closeAt := providerHallTestEpoch.Add(5 * time.Second)
	c.enqueue(providerHallEvent{kind: providerHallEventBarrier, at: closeAt})
	c.FlushNow()
	batches, gaps, closed = facts.snapshot()
	require.Len(t, gaps, 1)
	require.Equal(t, closeAt, closed[1])
	require.False(t, c.overflowed.Load())
	require.NotNil(t, batches[len(batches)-1].ConfirmedAt)
	require.Equal(t, closeAt, *batches[len(batches)-1].ConfirmedAt)

	// A second overflow opens a second gap from the last confirmed point.
	for i := 0; i < 8; i++ {
		trackers[1].sink.enqueue(providerHallEvent{kind: providerHallEventStart, row: trackers[1].row})
	}
	c.FlushNow()
	_, gaps, _ = facts.snapshot()
	require.Len(t, gaps, 2)
	require.Equal(t, closeAt, gaps[1].StartedAt)
}

func TestProviderHallCollectorWriteFailureRetainsAndBounds(t *testing.T) {
	c, facts := newTestHallCollector(t, 64, true)
	c.batchSize = 2
	facts.writeErr = errors.New("db down")
	tracker, _ := hallBegin(c, ProviderHallBeginInput{APIKeyID: 5, GroupID: 10, Protocol: "responses", RequestedModel: "gpt-test"})
	tracker.Finish()
	pending := c.flush([]providerHallEvent{{kind: providerHallEventStart, row: tracker.row, seq: 1}, {kind: providerHallEventFinish, row: tracker.row, seq: 2}})
	require.Len(t, pending, 2, "a failed write keeps the batch for retry")
	require.False(t, c.overflowed.Load())
	big := make([]providerHallEvent, 0, 10)
	for i := 0; i < 10; i++ {
		big = append(big, providerHallEvent{kind: providerHallEventStart, row: tracker.row, seq: uint64(i)})
	}
	pending = c.flush(big)
	require.Len(t, pending, 5, "unbounded retention is shed like an overflow")
	require.True(t, c.overflowed.Load())
	facts.writeErr = nil
	require.Empty(t, c.flush(pending))
	_, gaps, _ := facts.snapshot()
	require.Len(t, gaps, 1)
}

func TestProviderHallCollectorStartStopUnderLoad(t *testing.T) {
	facts := &providerHallFactRepoFake{
		targets:   []ProviderHallTargetRef{{TargetID: 1, GroupID: 10, ProfileID: 100, Model: "gpt-test", Protocol: "responses", Enabled: true}},
		probeKeys: map[int64]time.Time{},
	}
	cfg := &providerHallRepoStub{config: ProviderHallConfig{CollectionEnabled: true}}
	c := newProviderHallCollector(facts, cfg, "node-b", "test", 4096)
	c.flushInterval = 5 * time.Millisecond
	c.barrierInterval = 5 * time.Millisecond
	c.heartbeatInterval = 10 * time.Millisecond
	c.refreshInterval = 10 * time.Millisecond
	c.Start()
	c.Start()
	require.Eventually(t, func() bool { return c.index.Load() != nil }, time.Second, time.Millisecond)
	var wg sync.WaitGroup
	const workers, perWorker = 8, 25
	for w := 0; w < workers; w++ {
		wg.Go(func() {
			for i := 0; i < perWorker; i++ {
				tracker, _ := c.Begin(context.Background(), ProviderHallBeginInput{APIKeyID: 5, GroupID: 10, Protocol: "responses", RequestedModel: "gpt-test", StartedAt: time.Now()})
				tracker.RecordSubmission()
				tracker.Attempt(time.Now(), &OpenAIForwardResult{}, nil)
				c.RecordBilling(ProviderHallBillingEvent{TraceID: tracker.TraceID(), Status: "applied"})
				tracker.Finish()
			}
		})
	}
	wg.Wait()
	require.Eventually(t, func() bool {
		batches, _, _ := facts.snapshot()
		for _, b := range batches {
			if b.ConfirmedAt != nil {
				return true
			}
		}
		return false
	}, 2*time.Second, 5*time.Millisecond, "barriers advance the watermark while running")
	c.Stop()
	c.Stop()
	batches, gaps, _ := facts.snapshot()
	require.Empty(t, gaps)
	starts, finishes, billings := 0, 0, 0
	confirmed := false
	for _, b := range batches {
		starts += len(b.Starts)
		finishes += len(b.Finishes)
		billings += len(b.Billings)
		confirmed = confirmed || b.ConfirmedAt != nil
	}
	require.Equal(t, workers*perWorker, starts)
	require.Equal(t, workers*perWorker, finishes)
	require.Equal(t, workers*perWorker, billings)
	require.True(t, confirmed)
	require.Equal(t, []string{"shutdown"}, facts.exited)
	require.Equal(t, 1, facts.registered)
	_, _, dropped, _ := c.QueueStats()
	require.Zero(t, dropped)
}
