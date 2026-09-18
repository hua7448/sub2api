package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrProviderHallJobInFlight       = infraerrors.Conflict("PROVIDER_HALL_JOB_IN_FLIGHT", "job has samples in flight and cannot be cancelled")
	ErrProviderHallJobNotCancellable = infraerrors.Conflict("PROVIDER_HALL_JOB_NOT_CANCELLABLE", "job already finished")
	ErrProviderHallTasksDisabled     = infraerrors.BadRequest("PROVIDER_HALL_TASKS_DISABLED", "provider hall tasks are disabled")
	ErrProviderHallBudgetExhausted   = infraerrors.BadRequest("PROVIDER_HALL_BUDGET_EXHAUSTED", "provider hall daily budget is exhausted")
	ErrProviderHallTargetDisabled    = infraerrors.BadRequest("PROVIDER_HALL_TARGET_DISABLED", "target is not enabled")
)

const (
	providerHallRunnerLockKey       = "provider-hall-scheduler"
	providerHallRunnerLockTTL       = time.Minute
	providerHallRunnerTick          = 15 * time.Second
	providerHallRunnerLease         = 120 * time.Second
	providerHallRunnerLeaseRenew    = 30 * time.Second
	providerHallRunnerStopWait      = 90 * time.Second
	providerHallRunnerMaxInflight   = 2
	providerHallRunnerBacklogAfter  = 120 * time.Second
	providerHallRunnerUncertainWait = 10 * time.Minute
	providerHallRunnerProbeJitter   = 30 * time.Second
	providerHallRunnerVerifyJitter  = 5 * time.Minute
	providerHallRunnerMaxJobsTick   = 4
	providerHallRunnerWaitPoll      = 2 * time.Second
	providerHallRunnerWaitMax       = 60 * time.Second
)

// ProviderHallRunner schedules, claims and dispatches probe and verification
// jobs on the primary instance. Every paid request goes through the local
// gateway with the target's dedicated key; nothing is ever resent.
type ProviderHallRunner struct {
	jobs      ProviderHallJobRepository
	cfgRepo   ProviderHallRepository
	db        *sql.DB
	lockCache LeaderLockCache
	cfg       *config.Config
	client    *ProviderHallProbeClient
	version   string
	owner     string
	now       func() time.Time
	rng       *rand.Rand
	rngMu     sync.Mutex
	log       *zap.Logger

	tick         time.Duration
	lease        time.Duration
	backlogAfter time.Duration
	maxInflight  int
	stopWait     time.Duration

	startOnce sync.Once
	stopOnce  sync.Once
	stopCh    chan struct{}
	done      chan struct{}
	inflight  sync.WaitGroup
	stopping  chan struct{}
}

// NewProviderHallRunner builds the runner. buildVersion is sent as the
// User-Agent suffix so gateway logs can distinguish probe traffic.
func NewProviderHallRunner(jobs ProviderHallJobRepository, cfgRepo ProviderHallRepository, db *sql.DB, lockCache LeaderLockCache, cfg *config.Config, buildVersion string) *ProviderHallRunner {
	return &ProviderHallRunner{
		jobs:         jobs,
		cfgRepo:      cfgRepo,
		db:           db,
		lockCache:    lockCache,
		cfg:          cfg,
		client:       NewProviderHallProbeClient(ProviderHallProbeTimeout),
		version:      buildVersion,
		owner:        ProviderHallNodeID() + "/" + uuid.NewString()[:8],
		now:          time.Now,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
		log:          logger.L().With(zap.String("component", "service.provider_hall_runner")),
		tick:         providerHallRunnerTick,
		lease:        providerHallRunnerLease,
		backlogAfter: providerHallRunnerBacklogAfter,
		maxInflight:  providerHallRunnerMaxInflight,
		stopWait:     providerHallRunnerStopWait,
		stopCh:       make(chan struct{}),
		done:         make(chan struct{}),
		stopping:     make(chan struct{}),
	}
}

// SetClock overrides time for tests.
func (r *ProviderHallRunner) SetClock(now func() time.Time) {
	if r != nil && now != nil {
		r.now = now
		r.client.now = now
	}
}

// Owner is the lease owner identity of this runner.
func (r *ProviderHallRunner) Owner() string {
	if r == nil {
		return ""
	}
	return r.owner
}

// WaitInflight blocks until every job goroutine started by RunOnce finished.
func (r *ProviderHallRunner) WaitInflight() {
	if r != nil {
		r.inflight.Wait()
	}
}

// Ready reports whether the task runtime is implemented in this build.
func (r *ProviderHallRunner) Ready() bool { return r != nil && r.jobs != nil && r.cfgRepo != nil }

// JobControl exposes the immediate cancel/stale hooks for the config service.
func (r *ProviderHallRunner) JobControl() ProviderHallJobControl {
	if r == nil || r.jobs == nil {
		return nil
	}
	return r.jobs
}

func (r *ProviderHallRunner) Start() {
	if r == nil || !r.Ready() {
		return
	}
	r.startOnce.Do(func() { go r.loop() })
}

// Stop halts scheduling and waits up to 90s for in-flight dispatches, so a
// sample is never left dispatched without its result being recorded when the
// process can still do so.
func (r *ProviderHallRunner) Stop() {
	if r == nil {
		return
	}
	r.stopOnce.Do(func() {
		close(r.stopCh)
		waited := make(chan struct{})
		go func() {
			r.inflight.Wait()
			close(waited)
		}()
		select {
		case <-waited:
		case <-time.After(r.stopWait):
			r.log.Warn("provider_hall.runner_stop_timeout")
		}
		select {
		case <-r.done:
		case <-time.After(2 * time.Second):
		}
	})
}

func (r *ProviderHallRunner) loop() {
	defer close(r.done)
	ticker := time.NewTicker(r.tick)
	defer ticker.Stop()
	r.safeRunOnce()
	for {
		select {
		case <-ticker.C:
			r.safeRunOnce()
		case <-r.stopCh:
			return
		}
	}
}

func (r *ProviderHallRunner) safeRunOnce() {
	defer func() {
		if rec := recover(); rec != nil {
			r.log.Error("provider_hall.runner_panic", zap.Any("panic", rec))
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), r.tick*4)
	defer cancel()
	r.RunOnce(ctx)
}

func (r *ProviderHallRunner) stopped() bool {
	select {
	case <-r.stopCh:
		return true
	default:
		return false
	}
}

// RunOnce performs one scheduler tick under the leader lock. Exported for
// tests and for admin "run now" tooling.
func (r *ProviderHallRunner) RunOnce(ctx context.Context) {
	if r == nil || !r.Ready() {
		return
	}
	release, acquired := tryAcquireSingletonLeaderLock(ctx, r.lockCache, r.db, providerHallRunnerLockKey, r.owner, providerHallRunnerLockTTL)
	if !acquired {
		return
	}
	if release != nil {
		defer release()
	}
	now := r.now().UTC()
	cfg, err := r.cfgRepo.GetConfig(ctx)
	if err != nil {
		r.log.Warn("provider_hall.runner_config_failed", zap.Error(err))
		return
	}
	// Bookkeeping runs even when tasks are off: outstanding samples still
	// need their facts and bills reconciled.
	if _, err := r.jobs.ReconcileSamples(ctx, now, providerHallRunnerUncertainWait); err != nil {
		r.log.Warn("provider_hall.runner_reconcile_failed", zap.Error(err))
	}
	if _, _, err := r.jobs.RecoverExpiredLeases(ctx, now); err != nil {
		r.log.Warn("provider_hall.runner_recover_failed", zap.Error(err))
	}
	if _, err := r.jobs.MarkDriftedVerificationsStale(ctx); err != nil {
		r.log.Warn("provider_hall.runner_stale_failed", zap.Error(err))
	}
	r.finalize(ctx, now)

	targets, err := r.jobs.ListSchedulableTargets(ctx)
	if err != nil {
		r.log.Warn("provider_hall.runner_targets_failed", zap.Error(err))
		return
	}
	active := make([]int64, 0, len(targets))
	for _, t := range targets {
		if cfg.TasksEnabled && ProviderHallAutoSchedule(cfg.AutoScheduleEnabled) && t.AutoScheduleEnabled && t.Enabled && t.Listed && t.ProbeKeyID != nil {
			active = append(active, t.TargetID)
		}
	}
	if _, err := r.jobs.CancelQueuedNotIn(ctx, active, ProviderHallJobCodeTargetDisabled, now); err != nil {
		r.log.Warn("provider_hall.runner_cancel_failed", zap.Error(err))
	}
	if _, err := r.jobs.CancelQueuedDrifted(ctx, now); err != nil {
		r.log.Warn("provider_hall.runner_cancel_drift_failed", zap.Error(err))
	}
	if cfg.TasksEnabled && ProviderHallAutoSchedule(cfg.AutoScheduleEnabled) {
		r.schedule(ctx, cfg, targets, now)
	}
	r.claimAndRun(ctx, cfg, now)
}

func (r *ProviderHallRunner) jitter(max time.Duration, signed bool) time.Duration {
	r.rngMu.Lock()
	defer r.rngMu.Unlock()
	d := time.Duration(r.rng.Int63n(int64(max) + 1))
	if signed && r.rng.Intn(2) == 0 {
		return -d
	}
	return d
}

// schedule enqueues the current slot of every active target. Missed slots are
// never backfilled.
func (r *ProviderHallRunner) schedule(ctx context.Context, cfg *ProviderHallConfig, targets []ProviderHallSchedulableTarget, now time.Time) {
	if cfg.OperatorUserID == nil || cfg.GatewayOrigin == "" {
		return
	}
	for _, t := range targets {
		if !t.AutoScheduleEnabled || !t.Enabled || !t.Listed || t.ProbeKeyID == nil {
			continue
		}
		snapshot := ProviderHallJobSnapshot{Profile: t.Profile, TargetVersion: t.TargetVersion, ProbeKeyID: *t.ProbeKeyID, OperatorUserID: *cfg.OperatorUserID, GatewayOrigin: cfg.GatewayOrigin}
		for _, spec := range []struct {
			kind     ProviderHallJobKind
			interval int
			jitter   time.Duration
			signed   bool
		}{
			{ProviderHallJobProbe, t.ProbeIntervalSeconds, providerHallRunnerProbeJitter, true},
			{ProviderHallJobVerification, t.VerificationIntervalSeconds, providerHallRunnerVerifyJitter, false},
		} {
			if spec.interval <= 0 {
				continue
			}
			slot := now.Truncate(time.Duration(spec.interval) * time.Second)
			notBefore := slot.Add(r.jitter(spec.jitter, spec.signed))
			if notBefore.Before(slot) && spec.signed {
				// Negative jitter on the current slot would be in the past; keep it
				// but never earlier than the slot start minus jitter.
				notBefore = slot
			}
			if _, _, err := r.jobs.EnqueueSlot(ctx, ProviderHallEnqueueInput{Kind: spec.kind, TargetID: t.TargetID, GroupID: t.GroupID, ProfileID: t.ProfileID,
				Snapshot: snapshot, SlotAt: &slot, NotBefore: &notBefore}); err != nil {
				r.log.Warn("provider_hall.runner_enqueue_failed", zap.Int64("target_id", t.TargetID), zap.String("kind", string(spec.kind)), zap.Error(err))
			}
		}
	}
}

// EnqueueManual creates an on-demand job for the admin API. It validates the
// switch, target and key state and reuses an existing job for the same key.
func (r *ProviderHallRunner) EnqueueManual(ctx context.Context, kind ProviderHallJobKind, groupID, profileID int64, idempotencyKey string, actorID int64) (*ProviderHallJob, bool, error) {
	if r == nil || !r.Ready() {
		return nil, false, providerHallNotReady("tasks_enabled")
	}
	if kind != ProviderHallJobProbe && kind != ProviderHallJobVerification {
		return nil, false, providerHallInvalid("kind")
	}
	cfg, err := r.cfgRepo.GetConfig(ctx)
	if err != nil {
		return nil, false, err
	}
	if !cfg.TasksEnabled || cfg.OperatorUserID == nil || cfg.GatewayOrigin == "" {
		return nil, false, ErrProviderHallTasksDisabled
	}
	targets, err := r.jobs.ListSchedulableTargets(ctx)
	if err != nil {
		return nil, false, err
	}
	var target *ProviderHallSchedulableTarget
	for i := range targets {
		if targets[i].GroupID == groupID && targets[i].ProfileID == profileID {
			target = &targets[i]
		}
	}
	if target == nil || !target.Enabled || target.ProbeKeyID == nil {
		return nil, false, ErrProviderHallTargetDisabled
	}
	if expected, ok := ctx.Value(providerHallExpectedVersionKey{}).(ProviderHallExpectedVersions); ok &&
		(expected.Target != target.TargetVersion || expected.Profile != target.Profile.Version) {
		return nil, false, ErrProviderHallConflict
	}
	now := r.now().UTC()
	spend, err := r.jobs.SumSpend(ctx, ProviderHallBudgetDay(now))
	if err != nil {
		return nil, false, err
	}
	if ProviderHallBudgetExhausted(cfg.DailyBudget, spend.Confirmed) {
		return nil, false, ErrProviderHallBudgetExhausted
	}
	in := ProviderHallEnqueueInput{Kind: kind, TargetID: target.TargetID, GroupID: groupID, ProfileID: profileID, NotBefore: &now,
		Snapshot: ProviderHallJobSnapshot{Profile: target.Profile, TargetVersion: target.TargetVersion, ProbeKeyID: *target.ProbeKeyID, OperatorUserID: *cfg.OperatorUserID, GatewayOrigin: cfg.GatewayOrigin}}
	if key := strings.TrimSpace(idempotencyKey); key != "" {
		in.IdempotencyKey = &key
	}
	if actorID > 0 {
		in.RequestedBy = &actorID
	}
	return r.jobs.EnqueueManual(ctx, in)
}

// CancelJob is the admin cancel: queued or not-yet-dispatched running jobs.
func (r *ProviderHallRunner) CancelJob(ctx context.Context, jobID int64) error {
	if r == nil || !r.Ready() {
		return providerHallNotReady("tasks_enabled")
	}
	return r.jobs.Cancel(ctx, jobID, ProviderHallJobCodeCancelledByAdmin, r.now().UTC())
}

func (r *ProviderHallRunner) claimAndRun(ctx context.Context, cfg *ProviderHallConfig, now time.Time) {
	for i := 0; i < providerHallRunnerMaxJobsTick; i++ {
		if r.stopped() {
			return
		}
		job, err := r.jobs.Claim(ctx, r.owner, r.now().UTC(), r.lease)
		if err != nil {
			r.log.Warn("provider_hall.runner_claim_failed", zap.Error(err))
			return
		}
		if job == nil {
			return
		}
		r.inflight.Add(1)
		go func(job *ProviderHallJob) {
			defer r.inflight.Done()
			defer func() {
				if rec := recover(); rec != nil {
					r.log.Error("provider_hall.runner_job_panic", zap.Int64("job_id", job.ID), zap.Any("panic", rec))
				}
			}()
			r.runJob(context.Background(), cfg, job)
		}(job)
	}
}

// RunJob executes one claimed job to completion; exported for tests.
func (r *ProviderHallRunner) RunJob(ctx context.Context, job *ProviderHallJob) {
	cfg, err := r.cfgRepo.GetConfig(ctx)
	if err != nil {
		r.log.Warn("provider_hall.runner_config_failed", zap.Error(err))
		return
	}
	r.runJob(ctx, cfg, job)
}

func (r *ProviderHallRunner) runJob(ctx context.Context, cfg *ProviderHallConfig, job *ProviderHallJob) {
	jobCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	renewStop := make(chan struct{})
	defer close(renewStop)
	go r.renewLease(jobCtx, job, renewStop, cancel)

	if job.Kind == ProviderHallJobProbe || job.Kind == ProviderHallJobVerification {
		var cases []ProviderHallTestCase
		if job.Kind == ProviderHallJobProbe {
			cases = ProviderHallProbeSuite(job.Snapshot.Profile)
		} else {
			r.rngMu.Lock()
			cases = ProviderHallBuildSuite(job.Snapshot.Profile, r.rng)
			r.rngMu.Unlock()
		}
		if err := r.jobs.CreateSamples(jobCtx, job.ID, cases); err != nil {
			r.log.Warn("provider_hall.runner_samples_failed", zap.Int64("job_id", job.ID), zap.Error(err))
			return
		}
	}
	samples, err := r.jobs.ListSamples(jobCtx, job.ID)
	if err != nil {
		r.log.Warn("provider_hall.runner_list_samples_failed", zap.Int64("job_id", job.ID), zap.Error(err))
		return
	}
	for i := range samples {
		sample := &samples[i]
		if sample.Status != ProviderHallSamplePrepared {
			continue
		}
		tc, ok := providerHallCaseFromDetail(sample.Detail)
		if !ok {
			_ = r.jobs.Complete(jobCtx, job.ID, ProviderHallJobFailed, "sample_corrupt", "sample case missing", r.now().UTC())
			return
		}
		outcome, done := r.dispatchWithChecks(jobCtx, cfg, job, sample)
		if done {
			return
		}
		r.send(jobCtx, job, sample, tc, outcome)
		if jobCtx.Err() != nil {
			return
		}
	}
	r.completeJob(jobCtx, job)
}

func providerHallCaseFromDetail(detail json.RawMessage) (ProviderHallTestCase, bool) {
	var wrapper struct {
		Case *ProviderHallTestCase `json:"case"`
	}
	if err := json.Unmarshal(detail, &wrapper); err != nil || wrapper.Case == nil {
		return ProviderHallTestCase{}, false
	}
	return *wrapper.Case, true
}

// dispatchWithChecks loops until the sample is dispatched or the job leaves
// the running state. done=true means the job was cancelled/failed/requeued.
func (r *ProviderHallRunner) dispatchWithChecks(ctx context.Context, cfg *ProviderHallConfig, job *ProviderHallJob, sample *ProviderHallSample) (ProviderHallDispatchOutcome, bool) {
	waitStart := r.now()
	for {
		if ctx.Err() != nil {
			return ProviderHallDispatchOutcome{}, true
		}
		if r.stopped() {
			_ = r.jobs.Requeue(ctx, job.ID, r.now().UTC(), ProviderHallJobCodeRunnerStopped)
			return ProviderHallDispatchOutcome{}, true
		}
		now := r.now().UTC()
		outcome, err := r.jobs.TryDispatch(ctx, ProviderHallDispatchInput{Job: job, SampleID: sample.ID, Now: now, TasksEnabled: cfg.TasksEnabled, MaxInflight: r.maxInflight, BacklogAfter: r.backlogAfter})
		if err != nil {
			r.log.Warn("provider_hall.runner_dispatch_failed", zap.Int64("job_id", job.ID), zap.Error(err))
			_ = r.jobs.Requeue(ctx, job.ID, now.Add(time.Minute), "dispatch_error")
			return ProviderHallDispatchOutcome{}, true
		}
		switch outcome.Decision {
		case ProviderHallDispatchGo:
			if job.BudgetDay == nil {
				day := outcome.BudgetDay
				job.BudgetDay = &day
			}
			sample.Status = ProviderHallSampleDispatched
			return outcome, false
		case ProviderHallDispatchCancel:
			_ = r.jobs.Complete(ctx, job.ID, ProviderHallJobCancelled, outcome.Code, outcome.Message, now)
			return outcome, true
		case ProviderHallDispatchFail:
			_ = r.jobs.Complete(ctx, job.ID, ProviderHallJobFailed, outcome.Code, outcome.Message, now)
			return outcome, true
		case ProviderHallDispatchRequeue:
			_ = r.jobs.Requeue(ctx, job.ID, now.Add(time.Minute), outcome.Code)
			return outcome, true
		default: // wait for concurrency
			if r.now().Sub(waitStart) > providerHallRunnerWaitMax {
				// Give the slot back instead of holding a lease indefinitely.
				_ = r.jobs.Requeue(ctx, job.ID, r.now().UTC().Add(30*time.Second), "concurrency_wait")
				return outcome, true
			}
			select {
			case <-ctx.Done():
				return outcome, true
			case <-r.stopCh:
				_ = r.jobs.Requeue(ctx, job.ID, r.now().UTC(), ProviderHallJobCodeRunnerStopped)
				return outcome, true
			case <-time.After(providerHallRunnerWaitPoll):
			}
		}
	}
}

func (r *ProviderHallRunner) send(ctx context.Context, job *ProviderHallJob, sample *ProviderHallSample, tc ProviderHallTestCase, outcome ProviderHallDispatchOutcome) {
	req := ProviderHallProbeRequest{
		Origin:   outcome.Origin,
		Protocol: job.Snapshot.Profile.Protocol,
		Model:    job.Snapshot.Profile.Model,
		Key:      outcome.Key,
		Task:     ProviderHallTaskRef{Kind: ProviderHallSource(job.Kind), JobID: job.ID, SampleID: sample.ID, TraceID: sample.TraceID},
		Case:     tc,
		Version:  r.version,
	}
	// The send itself must not be cut short by the job context: once a sample
	// is dispatched its result is worth waiting for.
	sendCtx, cancel := context.WithTimeout(context.Background(), ProviderHallProbeTimeout+5*time.Second)
	defer cancel()
	result := r.client.Send(sendCtx, req)
	now := r.now().UTC()
	if result.Uncertain {
		if err := r.jobs.MarkSampleUncertain(ctx, sample.ID, result.ClientRequestID, result.ErrorCode, now); err != nil {
			r.log.Warn("provider_hall.runner_mark_uncertain_failed", zap.Int64("sample_id", sample.ID), zap.Error(err))
		}
		sample.Status = ProviderHallSampleUncertain
		return
	}
	receipt := ProviderHallSampleReceipt{SampleID: sample.ID, ReceivedAt: now, ClientRequestID: result.ClientRequestID, HTTPStatus: result.HTTPStatus, ErrorCode: result.ErrorCode}
	if result.Response != nil {
		resp := result.Response
		receipt.TTFTMs, receipt.TotalMs = resp.TTFTMs, resp.TotalMs
		receipt.GenerationMs = ProviderHallGenerationMs(resp.TotalMs, resp.TTFTMs)
		receipt.InputTokens, receipt.OutputTokens = resp.InputTokens, resp.OutputTokens
		receipt.ResponseModel = resp.Model
		receipt.Result, receipt.Detail = ProviderHallEvaluate(tc, resp)
		if receipt.Result == ProviderHallResultError && receipt.ErrorCode == "" {
			if code, _ := receipt.Detail["error"].(string); code != "" {
				receipt.ErrorCode = code
			}
		}
	} else {
		receipt.Result = ProviderHallResultError
		receipt.Detail = map[string]any{"test_id": tc.TestID, "seq": tc.Seq, "error": result.ErrorCode}
	}
	if receipt.Result == ProviderHallResultPassed && job.Kind == ProviderHallJobVerification && ProviderHallModelMatch(job.Snapshot.Profile, receipt.ResponseModel) == "mismatched" {
		// Keep the case result; the verdict counts mismatches separately.
		receipt.Detail["model_mismatch"] = true
	}
	if err := r.jobs.MarkSampleReceived(ctx, receipt); err != nil {
		r.log.Warn("provider_hall.runner_mark_received_failed", zap.Int64("sample_id", sample.ID), zap.Error(err))
	}
	sample.Status = ProviderHallSampleReceived
	sample.Result = &receipt.Result
	sample.ResponseModel = &receipt.ResponseModel
}

func (r *ProviderHallRunner) renewLease(ctx context.Context, job *ProviderHallJob, stop <-chan struct{}, cancel context.CancelFunc) {
	ticker := time.NewTicker(providerHallRunnerLeaseRenew)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			ok, err := r.jobs.RenewLease(ctx, job.ID, r.owner, r.now().UTC().Add(r.lease))
			if err == nil && !ok {
				r.log.Warn("provider_hall.runner_lease_lost", zap.Int64("job_id", job.ID))
				cancel()
				return
			}
		}
	}
}

// completeJob scores a job whose samples all reached a terminal state. Jobs
// with uncertain samples become unknown and are finalized by a later tick.
func (r *ProviderHallRunner) completeJob(ctx context.Context, job *ProviderHallJob) {
	samples, err := r.jobs.ListSamples(ctx, job.ID)
	if err != nil {
		r.log.Warn("provider_hall.runner_list_samples_failed", zap.Int64("job_id", job.ID), zap.Error(err))
		return
	}
	now := r.now().UTC()
	for _, s := range samples {
		if s.Status == ProviderHallSampleUncertain || s.Status == ProviderHallSampleDispatched {
			_ = r.jobs.MarkUnknown(ctx, job.ID, ProviderHallJobCodeSampleUncertain, now)
			return
		}
	}
	r.score(ctx, job, samples, now)
}

func (r *ProviderHallRunner) score(ctx context.Context, job *ProviderHallJob, samples []ProviderHallSample, now time.Time) {
	switch job.Kind {
	case ProviderHallJobProbe:
		passed := false
		code := ProviderHallJobCodeProbeFailed
		for _, s := range samples {
			if s.Result != nil && *s.Result == ProviderHallResultPassed {
				passed = true
			} else if s.ErrorCode != "" {
				code = s.ErrorCode
			}
		}
		if passed {
			_ = r.jobs.Complete(ctx, job.ID, ProviderHallJobSucceeded, "", "", now)
		} else {
			_ = r.jobs.Complete(ctx, job.ID, ProviderHallJobFailed, code, "probe sample did not pass", now)
		}
	case ProviderHallJobVerification:
		results := make([]ProviderHallSampleResult, 0, len(samples))
		for _, s := range samples {
			res := ProviderHallSampleResult{TestID: s.TestID, Seq: s.Seq, Status: s.Status}
			if s.Result != nil {
				res.Result = *s.Result
			}
			if s.ResponseModel != nil {
				res.ResponseModel = *s.ResponseModel
			}
			results = append(results, res)
		}
		verdict, execution, reason, summary := ProviderHallVerdict(results, job.Snapshot.Profile)
		if err := r.jobs.SaveVerification(ctx, ProviderHallVerification{JobID: job.ID, GroupID: job.GroupID, ProfileID: job.ProfileID, TargetID: job.TargetID,
			Verdict: verdict, ExecutionStatus: execution, ReasonCode: reason, Summary: summary, ProfileVersion: job.Snapshot.Profile.Version, TargetVersion: job.Snapshot.TargetVersion,
			CompletedAt: now, ExpiresAt: now.Add(ProviderHallVerificationTTL)}); err != nil {
			r.log.Warn("provider_hall.runner_save_verification_failed", zap.Int64("job_id", job.ID), zap.Error(err))
			_ = r.jobs.Complete(ctx, job.ID, ProviderHallJobFailed, "report_write_failed", err.Error(), now)
			return
		}
		status := ProviderHallJobSucceeded
		if execution == ProviderHallExecutionError {
			status = ProviderHallJobFailed
		}
		_ = r.jobs.Complete(ctx, job.ID, status, reason, fmt.Sprintf("verdict=%s", verdict), now)
	}
}

// finalize scores jobs that were left unknown (uncertain samples that have
// since been reconciled) or whose lease lapsed with every sample terminal.
func (r *ProviderHallRunner) finalize(ctx context.Context, now time.Time) {
	jobs, err := r.jobs.ListFinalizableJobs(ctx, now)
	if err != nil {
		r.log.Warn("provider_hall.runner_finalize_list_failed", zap.Error(err))
		return
	}
	for i := range jobs {
		job := &jobs[i]
		samples, err := r.jobs.ListSamples(ctx, job.ID)
		if err != nil {
			continue
		}
		r.score(ctx, job, samples, now)
	}
}
