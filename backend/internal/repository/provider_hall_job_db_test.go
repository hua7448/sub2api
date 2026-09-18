//go:build integration || providerhall_localdb

package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func init() {
	registerProviderHallLocalDBContract("b4_jobs", providerHallJobDatabaseContracts)
}

// hallGateway is a fake local gateway: it answers every model request with a
// terminal Responses stream that echoes the requested model, and records the
// task headers it saw.
type hallGateway struct {
	*httptest.Server
	mu     sync.Mutex
	tasks  []string
	hits   atomic.Int32
	status atomic.Int32
}

func newHallGateway(t *testing.T) *hallGateway {
	t.Helper()
	g := &hallGateway{}
	g.status.Store(http.StatusOK)
	g.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g.hits.Add(1)
		body, _ := io.ReadAll(r.Body)
		var req struct {
			Model string `json:"model"`
		}
		_ = json.Unmarshal(body, &req)
		g.mu.Lock()
		g.tasks = append(g.tasks, r.Header.Get(service.ProviderHallTaskHeader))
		g.mu.Unlock()
		if status := int(g.status.Load()); status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("X-Client-Request-ID", r.Header.Get("X-Client-Request-ID"))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "event: response.output_text.delta\ndata: {\"type\":\"response.output_text.delta\",\"delta\":\"OK\"}\n\n")
		fmt.Fprintf(w, "event: response.completed\ndata: {\"type\":\"response.completed\",\"response\":{\"model\":%q,\"usage\":{\"input_tokens\":12,\"output_tokens\":2}}}\n\n", req.Model)
	}))
	t.Cleanup(g.Server.Close)
	return g
}

type providerHallJobFixture struct {
	*providerHallTargetFixture
	jobs    service.ProviderHallJobRepository
	gateway *hallGateway
	target  service.ProviderHallTarget
	key     int64
}

func newProviderHallJobFixture(t *testing.T, db *sql.DB) *providerHallJobFixture {
	t.Helper()
	t.Setenv("PROVIDER_HALL_ALLOW_LOOPBACK", "1")
	ctx := context.Background()
	f := &providerHallJobFixture{providerHallTargetFixture: newProviderHallTargetFixture(t, db), jobs: NewProviderHallJobRepository(db), gateway: newHallGateway(t)}
	key := f.providerHallTargetFixture.key(t)
	f.key = key.ID
	listing, err := f.repo.GetGroup(ctx, f.group.ID)
	require.NoError(t, err)
	listing.Listed, listing.UpdatedBy = true, &f.user.ID
	savedListing, err := f.repo.SaveGroup(ctx, *listing)
	require.NoError(t, err)
	input := f.input(key.ID)
	input.Version = savedListing.Version
	// Daily slots keep the scheduler's own jobs out of the way of the
	// hand-built scenarios (the slot only changes at midnight UTC).
	input.Items[0].ProbeIntervalSeconds, input.Items[0].VerificationIntervalSeconds = 86400, 86400
	auto := true
	input.Items[0].AutoScheduleEnabled = &auto
	set, err := f.svc.SaveTargets(ctx, input, f.user.ID)
	require.NoError(t, err)
	f.target = set.Items[0]
	cfg, err := f.repo.GetConfig(ctx)
	require.NoError(t, err)
	cfg.CollectionEnabled, cfg.TasksEnabled = true, true
	cfg.AutoScheduleEnabled = &auto
	cfg.GatewayOrigin, cfg.DailyBudget, cfg.UpdatedBy = f.gateway.URL, "0", &f.user.ID
	_, err = f.repo.UpdateConfig(ctx, *cfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := db.ExecContext(ctx, `DELETE FROM provider_hall_spend WHERE api_key_id = $1 OR job_id IN (SELECT id FROM provider_hall_jobs WHERE group_id = $2)`, f.key, f.group.ID)
		require.NoError(t, err)
		for _, q := range []string{
			`DELETE FROM provider_hall_requests WHERE group_id = $1`,
			`DELETE FROM provider_hall_jobs WHERE group_id = $1`,
		} {
			_, err := db.ExecContext(ctx, q, f.group.ID)
			require.NoError(t, err)
		}
	})
	return f
}

// blockSlots pre-creates cancelled rows for today's slots so RunOnce cannot
// enqueue (and send) scheduler jobs inside a scenario that counts requests.
func (f *providerHallJobFixture) blockSlots(t *testing.T) {
	t.Helper()
	for _, kind := range []string{"probe", "verification"} {
		_, err := f.db.ExecContext(context.Background(), `
INSERT INTO provider_hall_jobs (kind, target_id, group_id, profile_id, config_snapshot, slot_at, status, error_code, finished_at)
VALUES ($1, $2, $3, $4, '{}', date_trunc('day', now() AT TIME ZONE 'UTC') AT TIME ZONE 'UTC', 'cancelled', 'test_blocked', now())
ON CONFLICT (target_id, kind, slot_at) WHERE slot_at IS NOT NULL DO NOTHING`, kind, f.target.ID, f.group.ID, f.profiles[0].ID)
		require.NoError(t, err)
	}
}

func (f *providerHallJobFixture) runner(t *testing.T) *service.ProviderHallRunner {
	t.Helper()
	r := service.NewProviderHallRunner(f.jobs, f.repo, f.db, nil, nil, "test")
	require.True(t, r.Ready())
	return r
}

func (f *providerHallJobFixture) snapshot(t *testing.T) service.ProviderHallJobSnapshot {
	t.Helper()
	cfg, err := f.repo.GetConfig(context.Background())
	require.NoError(t, err)
	targets, err := f.jobs.ListSchedulableTargets(context.Background())
	require.NoError(t, err)
	for _, target := range targets {
		if target.TargetID == f.target.ID {
			require.True(t, target.Enabled)
			require.True(t, target.Listed)
			return service.ProviderHallJobSnapshot{Profile: target.Profile, TargetVersion: target.TargetVersion, ProbeKeyID: *target.ProbeKeyID, OperatorUserID: *cfg.OperatorUserID, GatewayOrigin: cfg.GatewayOrigin}
		}
	}
	t.Fatalf("target %d not schedulable", f.target.ID)
	return service.ProviderHallJobSnapshot{}
}

func (f *providerHallJobFixture) enqueue(t *testing.T, kind service.ProviderHallJobKind, idem string) *service.ProviderHallJob {
	t.Helper()
	in := service.ProviderHallEnqueueInput{Kind: kind, TargetID: f.target.ID, GroupID: f.group.ID, ProfileID: f.profiles[0].ID, Snapshot: f.snapshot(t)}
	if idem != "" {
		in.IdempotencyKey = &idem
	}
	job, _, err := f.jobs.EnqueueManual(context.Background(), in)
	require.NoError(t, err)
	return job
}

func (f *providerHallJobFixture) job(t *testing.T, id int64) *service.ProviderHallJob {
	t.Helper()
	job, err := f.jobs.GetJob(context.Background(), id)
	require.NoError(t, err)
	return job
}

func (f *providerHallJobFixture) samples(t *testing.T, id int64) []service.ProviderHallSample {
	t.Helper()
	samples, err := f.jobs.ListSamples(context.Background(), id)
	require.NoError(t, err)
	return samples
}

func (f *providerHallJobFixture) insertRequest(t *testing.T, trace uuid.UUID, outcome, billingStatus string, cost string) {
	t.Helper()
	now := time.Now().UTC()
	_, err := f.db.ExecContext(context.Background(), `
INSERT INTO provider_hall_requests (trace_id, node_id, epoch_id, last_seq, group_id, profile_id, api_key_id, protocol, requested_model, source, started_at, ended_at, ttft_ms, submissions, outcome, usage_known, input_tokens, output_tokens,
    billing_request_id, billing_api_key_id, billing_status, billing_mode, actual_cost, billed_at)
VALUES ($1, 'test', 1, 1, $2, $3, $4, 'responses', $5, 'probe', $6, $6, 150, 1, $7, true, 12, 2, $8, $4, $9, 'token', $10::numeric, $6)`,
		trace.String(), f.group.ID, f.profiles[0].ID, f.key, f.profiles[0].Model, now, outcome, "client:"+trace.String(), billingStatus, cost)
	require.NoError(t, err)
}

func providerHallJobDatabaseContracts(t *testing.T, db *sql.DB) {
	ctx := context.Background()

	t.Run("jobs_schedule_claim_once_and_dispatch_through_gateway", func(t *testing.T) {
		f := newProviderHallJobFixture(t, db)
		r := f.runner(t)
		r.RunOnce(ctx)
		r.WaitInflight()
		jobs, total, err := f.jobs.ListJobs(ctx, service.ProviderHallJobFilter{GroupID: f.group.ID}, 1, 50)
		require.NoError(t, err)
		require.Equal(t, 2, total, "one probe and one verification slot per enabled target")
		kinds := map[service.ProviderHallJobKind]service.ProviderHallJob{}
		for _, j := range jobs {
			kinds[j.Kind] = j
			require.NotNil(t, j.SlotAt)
			require.NotNil(t, j.NotBefore)
		}
		require.Contains(t, kinds, service.ProviderHallJobProbe)
		require.Contains(t, kinds, service.ProviderHallJobVerification)
		// Re-running the same tick never creates a second job for the slot.
		r.RunOnce(ctx)
		r.WaitInflight()
		_, total, err = f.jobs.ListJobs(ctx, service.ProviderHallJobFilter{GroupID: f.group.ID}, 1, 50)
		require.NoError(t, err)
		require.Equal(t, 2, total)

		// Manual probe with two competing claimers: exactly one wins.
		job := f.enqueue(t, service.ProviderHallJobProbe, "")
		var wg sync.WaitGroup
		claimed := make(chan *service.ProviderHallJob, 2)
		for i := range 2 {
			wg.Go(func() {
				got, err := f.jobs.Claim(ctx, fmt.Sprintf("claimer-%d", i), time.Now().UTC(), 2*time.Minute)
				require.NoError(t, err)
				if got != nil && got.ID == job.ID {
					claimed <- got
				}
			})
		}
		wg.Wait()
		close(claimed)
		var winners []*service.ProviderHallJob
		for j := range claimed {
			winners = append(winners, j)
		}
		require.Len(t, winners, 1, "FOR UPDATE SKIP LOCKED hands the job to one claimer")
		require.Equal(t, service.ProviderHallJobRunning, winners[0].Status)
		before := f.gateway.hits.Load()
		r.RunJob(ctx, winners[0])
		require.Equal(t, before+1, f.gateway.hits.Load(), "a probe job sends exactly one request")
		done := f.job(t, job.ID)
		require.Equal(t, service.ProviderHallJobSucceeded, done.Status)
		require.NotNil(t, done.BudgetDay)
		samples := f.samples(t, job.ID)
		require.Len(t, samples, 1)
		require.Equal(t, service.ProviderHallSampleReceived, samples[0].Status)
		require.Equal(t, "passed", *samples[0].Result)
		require.Equal(t, f.profiles[0].Model, *samples[0].ResponseModel)
		require.Equal(t, 12, *samples[0].InputTokens)
		require.NotNil(t, samples[0].TTFTMs)
		require.Equal(t, samples[0].TraceID.String(), *samples[0].ClientRequestID)
		f.gateway.mu.Lock()
		require.Equal(t, fmt.Sprintf("probe:%d:%d:%s", job.ID, samples[0].ID, samples[0].TraceID), f.gateway.tasks[len(f.gateway.tasks)-1])
		f.gateway.mu.Unlock()
		health, err := f.jobs.LatestProbePerTarget(ctx, time.Now().Add(-time.Hour))
		require.NoError(t, err)
		require.Equal(t, "passed", health[f.target.ID].Result)
		series, err := f.jobs.ProbeSeries(ctx, service.ProviderHallProbeSeriesQuery{GroupID: f.group.ID, Start: time.Now().Add(-time.Hour), End: time.Now().Add(time.Hour), BucketSeconds: 300})
		require.NoError(t, err)
		require.NotEmpty(t, series)
	})

	t.Run("verification_job_scores_and_reports", func(t *testing.T) {
		f := newProviderHallJobFixture(t, db)
		f.blockSlots(t)
		r := f.runner(t)
		job := f.enqueue(t, service.ProviderHallJobVerification, "")
		claimed, err := f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), 2*time.Minute)
		require.NoError(t, err)
		require.Equal(t, job.ID, claimed.ID)
		r.RunJob(ctx, claimed)
		samples := f.samples(t, job.ID)
		require.Len(t, samples, 6, "arith and json suites; no tool suite without supports_tools")
		report, err := f.jobs.GetVerification(ctx, job.ID)
		require.NoError(t, err)
		require.NotNil(t, report)
		// The fake gateway answers "OK" to everything: every case fails, so the
		// verdict is failed with the arithmetic suite as reason.
		require.Equal(t, service.ProviderHallVerdictFailed, report.Verdict)
		require.Equal(t, "arith_failed", report.ReasonCode)
		require.Equal(t, service.ProviderHallExecutionCompleted, report.ExecutionStatus)
		require.Equal(t, 3, report.Summary.Arithmetic.Failed)
		require.Equal(t, 6, report.Summary.Model.Matched)
		require.Nil(t, report.Summary.Tool)
		require.WithinDuration(t, report.CompletedAt.Add(48*time.Hour), report.ExpiresAt, time.Second)
		require.Equal(t, service.ProviderHallJobSucceeded, f.job(t, job.ID).Status)
		latest, err := f.jobs.LatestVerification(ctx, f.group.ID, f.profiles[0].ID)
		require.NoError(t, err)
		require.Equal(t, job.ID, latest.JobID)
		list, total, err := f.jobs.ListVerifications(ctx, f.group.ID, nil, 1, 10)
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Len(t, list, 1)

		// Changing the profile identity stales the report (service hook and
		// runner fingerprint check both converge on the same result).
		f.svc.SetJobControl(f.jobs)
		profile := *f.profiles[0]
		profile.SupportsTools = true
		saved, err := f.svc.SaveProfile(ctx, profile, f.user.ID)
		require.NoError(t, err)
		f.profiles[0] = saved
		latest, err = f.jobs.LatestVerification(ctx, f.group.ID, f.profiles[0].ID)
		require.NoError(t, err)
		require.True(t, latest.Stale)
		require.True(t, latest.Expired(time.Now()))
	})

	t.Run("prepared_crash_requeues_dispatched_crash_never_resends", func(t *testing.T) {
		f := newProviderHallJobFixture(t, db)
		f.blockSlots(t)
		r := f.runner(t)
		// Prepared: lease expires, job returns to the queue and is claimable again.
		job := f.enqueue(t, service.ProviderHallJobProbe, "")
		claimed, err := f.jobs.Claim(ctx, "crashed-node", time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		require.Equal(t, job.ID, claimed.ID)
		require.NoError(t, f.jobs.CreateSamples(ctx, job.ID, service.ProviderHallProbeSuite(job.Snapshot.Profile)))
		requeued, unknown, err := f.jobs.RecoverExpiredLeases(ctx, time.Now().Add(2*time.Minute).UTC())
		require.NoError(t, err)
		require.Equal(t, int64(1), requeued)
		require.Zero(t, unknown)
		again := f.job(t, job.ID)
		require.Equal(t, service.ProviderHallJobQueued, again.Status)
		require.Nil(t, again.LeaseOwner)
		claimed, err = f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		require.Equal(t, job.ID, claimed.ID)
		require.Equal(t, 2, claimed.Attempts)

		// Dispatched: the sample was marked dispatched (commit before send) and
		// the node died. Nobody resends; reconciliation fills the result.
		samples := f.samples(t, job.ID)
		out, err := f.jobs.TryDispatch(ctx, service.ProviderHallDispatchInput{Job: claimed, SampleID: samples[0].ID, Now: time.Now().UTC(), TasksEnabled: true})
		require.NoError(t, err)
		require.Equal(t, service.ProviderHallDispatchGo, out.Decision)
		require.NotEmpty(t, out.Key)
		require.Equal(t, service.ProviderHallSampleDispatched, f.samples(t, job.ID)[0].Status)
		requeued, unknown, err = f.jobs.RecoverExpiredLeases(ctx, time.Now().Add(2*time.Minute).UTC())
		require.NoError(t, err)
		require.Zero(t, requeued)
		require.Equal(t, int64(1), unknown)
		require.Equal(t, service.ProviderHallJobUnknown, f.job(t, job.ID).Status)
		before := f.gateway.hits.Load()
		r.RunOnce(ctx)
		r.WaitInflight()
		require.Equal(t, before, f.gateway.hits.Load(), "an unknown job is only checked, never resent")
		require.Equal(t, service.ProviderHallJobUnknown, f.job(t, job.ID).Status)
		// The gateway's fact arrives (the request had actually completed).
		f.insertRequest(t, samples[0].TraceID, "success", "applied", "0.00300000")
		require.NoError(t, f.jobs.MarkSampleUncertain(ctx, samples[0].ID, "", service.ProviderHallErrTimeout, time.Now()))
		n, err := f.jobs.ReconcileSamples(ctx, time.Now().UTC(), 10*time.Minute)
		require.NoError(t, err)
		require.Positive(t, n)
		sample := f.samples(t, job.ID)[0]
		require.Equal(t, service.ProviderHallSampleReceived, sample.Status)
		require.Equal(t, "passed", *sample.Result)
		require.Equal(t, 150, *sample.TTFTMs)
		require.Equal(t, "confirmed", sample.BillingStatus)
		require.True(t, sample.ActualCost.Equal(decimal.RequireFromString("0.003")))
		r.RunOnce(ctx)
		r.WaitInflight()
		require.Equal(t, service.ProviderHallJobSucceeded, f.job(t, job.ID).Status, "finalized without any new request")
		require.Equal(t, before, f.gateway.hits.Load())
		spend, err := f.jobs.SumSpend(ctx, service.ProviderHallBudgetDay(time.Now()))
		require.NoError(t, err)
		require.True(t, spend.Confirmed.Equal(decimal.RequireFromString("0.003")), spend.Confirmed.String())

		// An uncertain sample with no fact at all times out after 10 minutes.
		job2 := f.enqueue(t, service.ProviderHallJobProbe, "")
		claimed2, err := f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		require.Equal(t, job2.ID, claimed2.ID)
		require.NoError(t, f.jobs.CreateSamples(ctx, job2.ID, service.ProviderHallProbeSuite(job2.Snapshot.Profile)))
		s2 := f.samples(t, job2.ID)[0]
		_, err = f.jobs.TryDispatch(ctx, service.ProviderHallDispatchInput{Job: claimed2, SampleID: s2.ID, Now: time.Now().Add(-11 * time.Minute).UTC(), TasksEnabled: true})
		require.NoError(t, err)
		require.NoError(t, f.jobs.MarkSampleUncertain(ctx, s2.ID, "", service.ProviderHallErrTimeout, time.Now()))
		_, err = f.jobs.ReconcileSamples(ctx, time.Now().UTC(), 10*time.Minute)
		require.NoError(t, err)
		s2 = f.samples(t, job2.ID)[0]
		require.Equal(t, service.ProviderHallSampleReceived, s2.Status)
		require.Equal(t, service.ProviderHallErrUncertainTimeout, s2.ErrorCode)
		require.Equal(t, "unbilled", s2.BillingStatus)

		// A definitive 4xx answer carries no bill: it is closed as unbilled after a
		// short grace (well under the backlog threshold), so one rejected sample
		// cannot pause every dispatch for the whole uncertain timeout.
		require.NoError(t, f.jobs.Complete(ctx, job2.ID, service.ProviderHallJobFailed, service.ProviderHallErrUncertainTimeout, "", time.Now().UTC()))
		job3 := f.enqueue(t, service.ProviderHallJobProbe, "")
		claimed3, err := f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		require.NotNil(t, claimed3)
		require.Equal(t, job3.ID, claimed3.ID)
		require.NoError(t, f.jobs.CreateSamples(ctx, job3.ID, service.ProviderHallProbeSuite(job3.Snapshot.Profile)))
		s3 := f.samples(t, job3.ID)[0]
		_, err = f.jobs.TryDispatch(ctx, service.ProviderHallDispatchInput{Job: claimed3, SampleID: s3.ID, Now: time.Now().Add(-3 * time.Minute).UTC(), TasksEnabled: true})
		require.NoError(t, err)
		require.NoError(t, f.jobs.MarkSampleReceived(ctx, service.ProviderHallSampleReceipt{SampleID: s3.ID, ReceivedAt: time.Now().Add(-2 * time.Minute).UTC(), HTTPStatus: 403, Result: "error", ErrorCode: "http_403"}))
		require.Equal(t, "pending", f.samples(t, job3.ID)[0].BillingStatus)
		_, err = f.jobs.ReconcileSamples(ctx, time.Now().Add(-100*time.Second).UTC(), 10*time.Minute)
		require.NoError(t, err)
		require.Equal(t, "pending", f.samples(t, job3.ID)[0].BillingStatus, "inside the grace a late bill may still arrive")
		_, err = f.jobs.ReconcileSamples(ctx, time.Now().UTC(), 10*time.Minute)
		require.NoError(t, err)
		require.Equal(t, "unbilled", f.samples(t, job3.ID)[0].BillingStatus, "definitive 4xx closes well before the uncertain timeout")
	})

	t.Run("idempotency_key_reuse_and_admin_cancel", func(t *testing.T) {
		f := newProviderHallJobFixture(t, db)
		f.blockSlots(t)
		r := f.runner(t)
		first, reused, err := r.EnqueueManual(ctx, service.ProviderHallJobProbe, f.group.ID, f.profiles[0].ID, "idem-1", f.user.ID)
		require.NoError(t, err)
		require.False(t, reused)
		require.Equal(t, f.user.ID, *first.RequestedBy)
		second, reused, err := r.EnqueueManual(ctx, service.ProviderHallJobProbe, f.group.ID, f.profiles[0].ID, "idem-1", f.user.ID)
		require.NoError(t, err)
		require.True(t, reused)
		require.Equal(t, first.ID, second.ID)
		_, _, err = r.EnqueueManual(ctx, service.ProviderHallJobProbe, f.group.ID, f.profiles[1].ID, "", f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallTargetDisabled, "profile without an enabled target")
		require.NoError(t, r.CancelJob(ctx, first.ID))
		require.Equal(t, service.ProviderHallJobCancelled, f.job(t, first.ID).Status)
		require.ErrorIs(t, r.CancelJob(ctx, first.ID), service.ErrProviderHallJobNotCancellable)
		// Running with a dispatched sample cannot be cancelled.
		job := f.enqueue(t, service.ProviderHallJobProbe, "")
		claimed, err := f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		require.NoError(t, f.jobs.CreateSamples(ctx, job.ID, service.ProviderHallProbeSuite(job.Snapshot.Profile)))
		require.NoError(t, r.CancelJob(ctx, job.ID), "running with only prepared samples is cancellable")
		job = f.enqueue(t, service.ProviderHallJobProbe, "")
		claimed, err = f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		require.Equal(t, job.ID, claimed.ID)
		require.NoError(t, f.jobs.CreateSamples(ctx, job.ID, service.ProviderHallProbeSuite(job.Snapshot.Profile)))
		_, err = f.jobs.TryDispatch(ctx, service.ProviderHallDispatchInput{Job: claimed, SampleID: f.samples(t, job.ID)[0].ID, Now: time.Now().UTC(), TasksEnabled: true})
		require.NoError(t, err)
		require.ErrorIs(t, r.CancelJob(ctx, job.ID), service.ErrProviderHallJobInFlight)
		counts, err := f.jobs.CountJobs(ctx, time.Now())
		require.NoError(t, err)
		require.Equal(t, 1, counts.Running)
	})

	t.Run("disable_cancels_unsent_and_config_drift_cancels_at_dispatch", func(t *testing.T) {
		f := newProviderHallJobFixture(t, db)
		f.blockSlots(t)
		r := f.runner(t)
		f.svc.SetJobControl(f.jobs)
		queued := f.enqueue(t, service.ProviderHallJobProbe, "")
		set, err := f.repo.GetTargets(ctx, f.group.ID)
		require.NoError(t, err)
		set.Items[0].Enabled = false
		_, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.NoError(t, err)
		cancelled := f.job(t, queued.ID)
		require.Equal(t, service.ProviderHallJobCancelled, cancelled.Status)
		require.Equal(t, service.ProviderHallJobCodeTargetDisabled, cancelled.ErrorCode)
		before := f.gateway.hits.Load()
		r.RunOnce(ctx)
		r.WaitInflight()
		require.Equal(t, before, f.gateway.hits.Load())
		_, total, err := f.jobs.ListJobs(ctx, service.ProviderHallJobFilter{GroupID: f.group.ID, Status: service.ProviderHallJobQueued}, 1, 10)
		require.NoError(t, err)
		require.Zero(t, total, "a disabled target schedules nothing")

		// Re-enable, queue a job, then bump the target version: the stale
		// snapshot is cancelled at dispatch time.
		set, err = f.repo.GetTargets(ctx, f.group.ID)
		require.NoError(t, err)
		set.Items[0].Enabled = true
		set, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.NoError(t, err)
		f.target = set.Items[0]
		job := f.enqueue(t, service.ProviderHallJobProbe, "")
		claimed, err := f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		require.Equal(t, job.ID, claimed.ID)
		set.Items[0].ProbeIntervalSeconds = 600
		_, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.NoError(t, err)
		r.RunJob(ctx, claimed)
		require.Equal(t, before, f.gateway.hits.Load(), "no request after a config change")
		drifted := f.job(t, job.ID)
		require.Equal(t, service.ProviderHallJobCancelled, drifted.Status)
		require.Equal(t, service.ProviderHallJobCodeConfigChanged, drifted.ErrorCode)
		// Queued jobs with a drifted snapshot are cancelled by the tick too.
		stale := f.enqueue(t, service.ProviderHallJobProbe, "")
		_, err = f.db.ExecContext(ctx, `UPDATE provider_hall_jobs SET config_snapshot = jsonb_set(config_snapshot, '{target_version}', '999') WHERE id = $1`, stale.ID)
		require.NoError(t, err)
		n, err := f.jobs.CancelQueuedDrifted(ctx, time.Now().UTC())
		require.NoError(t, err)
		require.Equal(t, int64(1), n)
	})

	t.Run("invalid_key_never_sends", func(t *testing.T) {
		f := newProviderHallJobFixture(t, db)
		f.blockSlots(t)
		r := f.runner(t)
		job := f.enqueue(t, service.ProviderHallJobProbe, "")
		_, err := f.db.ExecContext(ctx, `UPDATE api_keys SET status = 'inactive' WHERE id = $1`, f.key)
		require.NoError(t, err)
		claimed, err := f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		before := f.gateway.hits.Load()
		r.RunJob(ctx, claimed)
		require.Equal(t, before, f.gateway.hits.Load())
		failed := f.job(t, job.ID)
		require.Equal(t, service.ProviderHallJobFailed, failed.Status)
		require.Equal(t, service.ProviderHallJobCodeProbeKeyInvalid, failed.ErrorCode)
		require.Equal(t, service.ProviderHallSamplePrepared, f.samples(t, job.ID)[0].Status, "the sample was never dispatched")
	})

	t.Run("budget_backlog_concurrency_and_shanghai_day", func(t *testing.T) {
		f := newProviderHallJobFixture(t, db)
		f.blockSlots(t)
		r := f.runner(t)
		cfg, err := f.repo.GetConfig(ctx)
		require.NoError(t, err)
		cfg.DailyBudget, cfg.UpdatedBy = "1.00000000", &f.user.ID
		_, err = f.repo.UpdateConfig(ctx, *cfg)
		require.NoError(t, err)

		// 16:30Z on 2026-09-12 is already 2026-09-13 in Asia/Shanghai.
		dispatchAt := time.Date(2026, 9, 12, 16, 30, 0, 0, time.UTC)
		job := f.enqueue(t, service.ProviderHallJobVerification, "")
		claimed, err := f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		require.NoError(t, f.jobs.CreateSamples(ctx, job.ID, service.ProviderHallBuildSuite(job.Snapshot.Profile, nil)))
		samples := f.samples(t, job.ID)
		out, err := f.jobs.TryDispatch(ctx, service.ProviderHallDispatchInput{Job: claimed, SampleID: samples[0].ID, Now: dispatchAt, TasksEnabled: true})
		require.NoError(t, err)
		require.Equal(t, service.ProviderHallDispatchGo, out.Decision)
		require.Equal(t, "2026-09-13", out.BudgetDay)
		require.Equal(t, "2026-09-13", *f.job(t, job.ID).BudgetDay)
		claimed.BudgetDay = f.job(t, job.ID).BudgetDay

		// Per-target concurrency: the second sample waits while the first is in flight.
		out, err = f.jobs.TryDispatch(ctx, service.ProviderHallDispatchInput{Job: claimed, SampleID: samples[1].ID, Now: dispatchAt, TasksEnabled: true})
		require.NoError(t, err)
		require.Equal(t, service.ProviderHallDispatchWait, out.Decision)
		// The in-flight bill lands after the budget is spent: it is still recorded.
		_, err = f.db.ExecContext(ctx, `INSERT INTO provider_hall_spend (billing_request_id, api_key_id, budget_day, actual_cost, status, confirmed_at) VALUES ('client:other', $1, '2026-09-13', 1.5, 'confirmed', now())`, f.key)
		require.NoError(t, err)
		f.insertRequest(t, samples[0].TraceID, "success", "pending", "0")
		facts := NewProviderHallFactRepository(f.db)
		require.NoError(t, facts.WriteBatch(ctx, service.ProviderHallFactBatch{Billings: []service.ProviderHallBillingEvent{{
			TraceID: samples[0].TraceID, BillingRequestID: "client:" + samples[0].TraceID.String(), APIKeyID: f.key, Status: "applied", Mode: "token",
			ActualCost: decimal.RequireFromString("0.25"), BilledAt: dispatchAt.Add(time.Minute)}}}))
		spend, err := f.jobs.SumSpend(ctx, "2026-09-13")
		require.NoError(t, err)
		require.True(t, spend.Confirmed.Equal(decimal.RequireFromString("1.75")), spend.Confirmed.String())
		require.Equal(t, 1, spend.InFlight)
		require.Equal(t, "confirmed", f.samples(t, job.ID)[0].BillingStatus)
		var spendDay string
		require.NoError(t, f.db.QueryRowContext(ctx, `SELECT budget_day::text FROM provider_hall_spend WHERE sample_id = $1`, samples[0].ID).Scan(&spendDay))
		require.Equal(t, "2026-09-13", spendDay, "spend inherits the job's day even after midnight")
		require.NoError(t, f.jobs.MarkSampleReceived(ctx, service.ProviderHallSampleReceipt{SampleID: samples[0].ID, ReceivedAt: dispatchAt.Add(time.Second), HTTPStatus: 200, Result: "failed", Detail: map[string]any{}}))

		// Budget reached: the next sample cancels the job; manual enqueue refuses.
		r.SetClock(func() time.Time { return dispatchAt.Add(2 * time.Minute) })
		out, err = f.jobs.TryDispatch(ctx, service.ProviderHallDispatchInput{Job: claimed, SampleID: samples[1].ID, Now: dispatchAt.Add(2 * time.Minute), TasksEnabled: true})
		require.NoError(t, err)
		require.Equal(t, service.ProviderHallDispatchCancel, out.Decision)
		require.Equal(t, service.ProviderHallJobCodeBudgetExhausted, out.Code)
		_, _, err = r.EnqueueManual(ctx, service.ProviderHallJobProbe, f.group.ID, f.profiles[0].ID, "", f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallBudgetExhausted)
		_, err = f.db.ExecContext(ctx, `DELETE FROM provider_hall_spend WHERE billing_request_id = 'client:other'`)
		require.NoError(t, err)

		// Billing backlog: a bill pending for more than 120s pauses dispatch.
		_, err = f.db.ExecContext(ctx, `UPDATE provider_hall_samples SET billing_status = 'pending', dispatched_at = $2 WHERE id = $1`, samples[0].ID, dispatchAt)
		require.NoError(t, err)
		out, err = f.jobs.TryDispatch(ctx, service.ProviderHallDispatchInput{Job: claimed, SampleID: samples[1].ID, Now: dispatchAt.Add(3 * time.Minute), TasksEnabled: true})
		require.NoError(t, err)
		require.Equal(t, service.ProviderHallDispatchRequeue, out.Decision)
		require.Equal(t, service.ProviderHallJobCodeBillingBacklog, out.Code)
		out, err = f.jobs.TryDispatch(ctx, service.ProviderHallDispatchInput{Job: claimed, SampleID: samples[1].ID, Now: dispatchAt.Add(time.Minute), TasksEnabled: true})
		require.NoError(t, err)
		require.Equal(t, service.ProviderHallDispatchGo, out.Decision, "a young pending bill does not block")

		// Global concurrency: two dispatched samples anywhere block a third target.
		other := f.enqueue(t, service.ProviderHallJobProbe, "")
		_, err = f.db.ExecContext(ctx, `UPDATE provider_hall_jobs SET status = 'running', target_id = target_id + 1000000 WHERE id = $1`, other.ID)
		require.NoError(t, err)
		require.NoError(t, f.jobs.CreateSamples(ctx, other.ID, service.ProviderHallProbeSuite(other.Snapshot.Profile)))
		otherSample := f.samples(t, other.ID)[0]
		_, err = f.db.ExecContext(ctx, `UPDATE provider_hall_samples SET status = 'dispatched', dispatched_at = $2 WHERE id = $1`, otherSample.ID, dispatchAt.Add(time.Minute))
		require.NoError(t, err)
		third := f.enqueue(t, service.ProviderHallJobProbe, "")
		_, err = f.db.ExecContext(ctx, `UPDATE provider_hall_jobs SET status = 'running', target_id = target_id + 2000000 WHERE id = $1`, third.ID)
		require.NoError(t, err)
		require.NoError(t, f.jobs.CreateSamples(ctx, third.ID, service.ProviderHallProbeSuite(third.Snapshot.Profile)))
		thirdJob := f.job(t, third.ID)
		_, err = f.db.ExecContext(ctx, `UPDATE provider_hall_samples SET billing_status = 'unbilled' WHERE job_id = $1`, job.ID)
		require.NoError(t, err)
		var inflight int
		require.NoError(t, f.db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_samples WHERE status = 'dispatched'`).Scan(&inflight))
		require.Equal(t, 2, inflight)
		// thirdJob's target no longer exists: the check order puts target
		// existence before concurrency, so use the real target for the count.
		_, err = f.db.ExecContext(ctx, `UPDATE provider_hall_jobs SET target_id = $2 WHERE id = $1`, third.ID, f.target.ID)
		require.NoError(t, err)
		thirdJob.TargetID = f.target.ID
		out, err = f.jobs.TryDispatch(ctx, service.ProviderHallDispatchInput{Job: thirdJob, SampleID: f.samples(t, third.ID)[0].ID, Now: dispatchAt.Add(time.Minute), TasksEnabled: true})
		require.NoError(t, err)
		require.Equal(t, service.ProviderHallDispatchWait, out.Decision)
	})

	t.Run("tasks_disabled_cancels_at_dispatch_and_loopback_rule", func(t *testing.T) {
		f := newProviderHallJobFixture(t, db)
		f.blockSlots(t)
		r := f.runner(t)
		job := f.enqueue(t, service.ProviderHallJobProbe, "")
		claimed, err := f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		cfg, err := f.repo.GetConfig(ctx)
		require.NoError(t, err)
		cfg.TasksEnabled, cfg.UpdatedBy = false, &f.user.ID
		_, err = f.repo.UpdateConfig(ctx, *cfg)
		require.NoError(t, err)
		before := f.gateway.hits.Load()
		r.RunJob(ctx, claimed)
		require.Equal(t, before, f.gateway.hits.Load())
		require.Equal(t, service.ProviderHallJobCodeTasksDisabled, f.job(t, job.ID).ErrorCode)
		require.Equal(t, service.ProviderHallJobCancelled, f.job(t, job.ID).Status)

		// The runner refuses a loopback origin without the explicit override.
		t.Setenv("PROVIDER_HALL_ALLOW_LOOPBACK", "")
		_, err = f.db.ExecContext(ctx, `UPDATE provider_hall_config SET tasks_enabled = true WHERE id = 1`)
		require.NoError(t, err)
		job2 := f.enqueue(t, service.ProviderHallJobProbe, "")
		claimed2, err := f.jobs.Claim(ctx, r.Owner(), time.Now().UTC(), time.Minute)
		require.NoError(t, err)
		require.Equal(t, job2.ID, claimed2.ID)
		r.RunJob(ctx, claimed2)
		require.Equal(t, before, f.gateway.hits.Load(), "loopback origin blocked without PROVIDER_HALL_ALLOW_LOOPBACK")
		s := f.samples(t, job2.ID)[0]
		require.Equal(t, service.ProviderHallSampleReceived, s.Status)
		require.Equal(t, "request_invalid", s.ErrorCode)
	})
}
