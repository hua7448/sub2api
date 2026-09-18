//go:build unit

package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

type providerHallTaskCfgRepo struct {
	service.ProviderHallRepository
	cfg service.ProviderHallConfig
}

func (r *providerHallTaskCfgRepo) GetConfig(context.Context) (*service.ProviderHallConfig, error) {
	cfg := r.cfg
	return &cfg, nil
}

type providerHallRunnerFake struct {
	enqueued []struct {
		kind           service.ProviderHallJobKind
		group, profile int64
		key            string
		actor          int64
	}
	reused     bool
	enqueueErr error
	cancelErr  error
	cancelled  []int64
}

func (f *providerHallRunnerFake) EnqueueManual(_ context.Context, kind service.ProviderHallJobKind, groupID, profileID int64, key string, actorID int64) (*service.ProviderHallJob, bool, error) {
	f.enqueued = append(f.enqueued, struct {
		kind           service.ProviderHallJobKind
		group, profile int64
		key            string
		actor          int64
	}{kind, groupID, profileID, key, actorID})
	if f.enqueueErr != nil {
		return nil, false, f.enqueueErr
	}
	return &service.ProviderHallJob{ID: 42, Kind: kind, GroupID: groupID, ProfileID: profileID, Status: service.ProviderHallJobQueued}, f.reused, nil
}

func (f *providerHallRunnerFake) CancelJob(_ context.Context, id int64) error {
	f.cancelled = append(f.cancelled, id)
	return f.cancelErr
}

type providerHallJobReadsFake struct {
	jobs      []service.ProviderHallJob
	filter    service.ProviderHallJobFilter
	page      int
	pageSize  int
	samples   []service.ProviderHallSample
	report    *service.ProviderHallVerification
	counts    service.ProviderHallJobCounts
	spend     service.ProviderHallSpendSummary
	verifs    []service.ProviderHallVerification
	verifArgs struct {
		group   int64
		profile *int64
	}
}

func (f *providerHallJobReadsFake) ListJobs(_ context.Context, filter service.ProviderHallJobFilter, page, pageSize int) ([]service.ProviderHallJob, int, error) {
	f.filter, f.page, f.pageSize = filter, page, pageSize
	return f.jobs, len(f.jobs), nil
}
func (f *providerHallJobReadsFake) GetJob(_ context.Context, id int64) (*service.ProviderHallJob, error) {
	for i := range f.jobs {
		if f.jobs[i].ID == id {
			return &f.jobs[i], nil
		}
	}
	return nil, service.ErrProviderHallNotFound
}
func (f *providerHallJobReadsFake) ListSamples(context.Context, int64) ([]service.ProviderHallSample, error) {
	return f.samples, nil
}
func (f *providerHallJobReadsFake) GetVerification(context.Context, int64) (*service.ProviderHallVerification, error) {
	return f.report, nil
}
func (f *providerHallJobReadsFake) CountJobs(context.Context, time.Time) (service.ProviderHallJobCounts, error) {
	return f.counts, nil
}
func (f *providerHallJobReadsFake) SumSpend(_ context.Context, day string) (service.ProviderHallSpendSummary, error) {
	f.spend.Day = day
	return f.spend, nil
}
func (f *providerHallJobReadsFake) ListVerifications(_ context.Context, groupID int64, profileID *int64, page, pageSize int) ([]service.ProviderHallVerification, int, error) {
	f.verifArgs.group, f.verifArgs.profile, f.page, f.pageSize = groupID, profileID, page, pageSize
	return f.verifs, len(f.verifs), nil
}

type providerHallAdminRepoFake struct {
	epochs []service.ProviderHallNodeEpoch
	gaps   []service.ProviderHallOpenGap
	state  service.ProviderHallAggregatorState
	dirty  int
	rec    service.ProviderHallReconciliationCounts
}

func (f *providerHallAdminRepoFake) ListLatestEpochs(context.Context) ([]service.ProviderHallNodeEpoch, error) {
	return f.epochs, nil
}
func (f *providerHallAdminRepoFake) ListOpenGaps(context.Context) ([]service.ProviderHallOpenGap, error) {
	return f.gaps, nil
}
func (f *providerHallAdminRepoFake) GetAggregatorState(context.Context) (*service.ProviderHallAggregatorState, error) {
	st := f.state
	return &st, nil
}
func (f *providerHallAdminRepoFake) CountDirty(context.Context) (int, error) { return f.dirty, nil }
func (f *providerHallAdminRepoFake) CountReconciliation(context.Context, time.Time) (service.ProviderHallReconciliationCounts, error) {
	return f.rec, nil
}

type providerHallQueueFake struct{}

func (providerHallQueueFake) QueueStats() (int, int, uint64, bool) { return 3, 8192, 1, false }

type providerHallAggFake struct {
	at  time.Time
	err string
}

func (f providerHallAggFake) Status() (time.Time, string) { return f.at, f.err }

func providerHallTaskRequest(method, path, body string, params gin.Params, role string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var reader *strings.Reader
	if body != "" {
		reader = strings.NewReader(body)
	} else {
		reader = strings.NewReader("")
	}
	c.Request = httptest.NewRequest(method, path, reader)
	c.Params = params
	if role != "" {
		c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 9})
		c.Set(string(middleware.ContextKeyUserRole), role)
	}
	return w, c
}

func providerHallTaskHandler(runner *providerHallRunnerFake, jobs *providerHallJobReadsFake, cfg service.ProviderHallConfig, admin *providerHallAdminRepoFake) *ProviderHallHandler {
	svc := service.NewProviderHallService(&providerHallTaskCfgRepo{cfg: cfg}, nil, nil, nil)
	h := &ProviderHallHandler{service: svc, now: func() time.Time { return time.Date(2026, 9, 12, 3, 0, 0, 0, time.UTC) }}
	if runner != nil {
		h.runner = runner
	}
	if jobs != nil {
		h.jobs = jobs
	}
	if admin != nil {
		h.admin = admin
	}
	return h
}

func TestProviderHallAdminHandler_TaskAuthentication(t *testing.T) {
	h := NewProviderHallHandler(nil, nil, nil, nil, nil, nil)
	for _, method := range []gin.HandlerFunc{h.EnqueueProbe, h.EnqueueVerification, h.ListJobs, h.GetJob, h.CancelJob, h.ListGroupVerifications, h.Health} {
		for _, role := range []string{"", service.RoleUser} {
			w, c := providerHallTaskRequest(http.MethodGet, "/", "", gin.Params{{Key: "id", Value: "1"}}, role)
			method(c)
			if role == "" {
				require.Equal(t, 401, w.Code)
			} else {
				require.Equal(t, 403, w.Code)
			}
		}
	}
}

func TestProviderHallAdminHandler_Probes(t *testing.T) {
	for _, tc := range []struct {
		name   string
		body   string
		runner *providerHallRunnerFake
		status int
		reason string
		reused bool
	}{
		{"accepted", `{"profile_id":7,"idempotency_key":"click-1"}`, &providerHallRunnerFake{}, 202, "", false},
		{"reused", `{"profile_id":7,"idempotency_key":"click-1"}`, &providerHallRunnerFake{reused: true}, 202, "", true},
		{"tasks_disabled", `{"profile_id":7}`, &providerHallRunnerFake{enqueueErr: service.ErrProviderHallTasksDisabled}, 400, "PROVIDER_HALL_TASKS_DISABLED", false},
		{"budget_exhausted", `{"profile_id":7}`, &providerHallRunnerFake{enqueueErr: service.ErrProviderHallBudgetExhausted}, 400, "PROVIDER_HALL_BUDGET_EXHAUSTED", false},
		{"target_disabled", `{"profile_id":7}`, &providerHallRunnerFake{enqueueErr: service.ErrProviderHallTargetDisabled}, 400, "PROVIDER_HALL_TARGET_DISABLED", false},
		{"missing_profile", `{"idempotency_key":"x"}`, &providerHallRunnerFake{}, 400, "", false},
		{"long_key", `{"profile_id":7,"idempotency_key":"` + strings.Repeat("k", 129) + `"}`, &providerHallRunnerFake{}, 400, "", false},
		{"unknown_field", `{"profile_id":7,"requested_by":1}`, &providerHallRunnerFake{}, 400, "", false},
		{"runner_absent", `{"profile_id":7}`, nil, 409, "PROVIDER_HALL_NOT_READY", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := providerHallTaskHandler(tc.runner, nil, service.ProviderHallConfig{}, nil)
			w, c := providerHallTaskRequest(http.MethodPost, "/groups/5/probes", tc.body, gin.Params{{Key: "id", Value: "5"}}, service.RoleAdmin)
			h.EnqueueProbe(c)
			require.Equal(t, tc.status, w.Code, w.Body.String())
			if tc.reason != "" {
				require.Contains(t, w.Body.String(), `"reason":"`+tc.reason+`"`)
			}
			if tc.status == 202 {
				var envelope struct {
					Data struct {
						JobID  int64  `json:"job_id"`
						Status string `json:"status"`
						Reused bool   `json:"reused"`
					} `json:"data"`
				}
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &envelope))
				require.Equal(t, int64(42), envelope.Data.JobID)
				require.Equal(t, "queued", envelope.Data.Status)
				require.Equal(t, tc.reused, envelope.Data.Reused)
				require.Len(t, tc.runner.enqueued, 1)
				require.Equal(t, service.ProviderHallJobProbe, tc.runner.enqueued[0].kind)
				require.Equal(t, int64(5), tc.runner.enqueued[0].group)
				require.Equal(t, int64(7), tc.runner.enqueued[0].profile)
				require.Equal(t, "click-1", tc.runner.enqueued[0].key)
				require.Equal(t, int64(9), tc.runner.enqueued[0].actor)
			} else if tc.runner != nil && tc.runner.enqueueErr == nil {
				require.Empty(t, tc.runner.enqueued, "invalid input must not reach the runner")
			}
		})
	}
}

func TestProviderHallAdminHandler_Verifications(t *testing.T) {
	runner := &providerHallRunnerFake{}
	h := providerHallTaskHandler(runner, nil, service.ProviderHallConfig{}, nil)
	w, c := providerHallTaskRequest(http.MethodPost, "/groups/5/verifications", `{"profile_id":8,"idempotency_key":"v-1"}`, gin.Params{{Key: "id", Value: "5"}}, service.RoleAdmin)
	h.EnqueueVerification(c)
	require.Equal(t, 202, w.Code, w.Body.String())
	require.Equal(t, service.ProviderHallJobVerification, runner.enqueued[0].kind)

	// Group reports include unlisted groups and pass the optional profile filter.
	jobs := &providerHallJobReadsFake{verifs: []service.ProviderHallVerification{{JobID: 3, GroupID: 5, ProfileID: 8, Verdict: "passed", ExecutionStatus: "completed", CompletedAt: time.Now(), ExpiresAt: time.Now().Add(time.Hour), Summary: service.ProviderHallVerificationSummary{}}}}
	h = providerHallTaskHandler(nil, jobs, service.ProviderHallConfig{}, nil)
	w, c = providerHallTaskRequest(http.MethodGet, "/groups/5/verifications?profile_id=8&page=2&page_size=10", "", gin.Params{{Key: "id", Value: "5"}}, service.RoleAdmin)
	h.ListGroupVerifications(c)
	require.Equal(t, 200, w.Code, w.Body.String())
	require.Equal(t, int64(5), jobs.verifArgs.group)
	require.NotNil(t, jobs.verifArgs.profile)
	require.Equal(t, int64(8), *jobs.verifArgs.profile)
	require.Equal(t, 2, jobs.page)
	require.Equal(t, 10, jobs.pageSize)
	require.Contains(t, w.Body.String(), `"verdict":"passed"`)
	require.Contains(t, w.Body.String(), `"total":1`)
	for _, query := range []string{"profile_id=0", "page_size=51", "page=0"} {
		w, c = providerHallTaskRequest(http.MethodGet, "/groups/5/verifications?"+query, "", gin.Params{{Key: "id", Value: "5"}}, service.RoleAdmin)
		h.ListGroupVerifications(c)
		require.Equal(t, 400, w.Code, query)
	}
}

func TestProviderHallAdminHandler_Jobs(t *testing.T) {
	jobs := &providerHallJobReadsFake{
		jobs:    []service.ProviderHallJob{{ID: 11, Kind: service.ProviderHallJobVerification, GroupID: 5, ProfileID: 7, Status: service.ProviderHallJobRunning, Snapshot: service.ProviderHallJobSnapshot{ProbeKeyID: 3, GatewayOrigin: "https://gw.example"}}},
		samples: []service.ProviderHallSample{{ID: 1, JobID: 11, TestID: "arith", Seq: 1, Status: service.ProviderHallSampleReceived, Detail: json.RawMessage(`{"expected":"42"}`)}},
		report:  &service.ProviderHallVerification{JobID: 11, Verdict: "insufficient", ExecutionStatus: "partial"},
	}
	h := providerHallTaskHandler(nil, jobs, service.ProviderHallConfig{}, nil)
	w, c := providerHallTaskRequest(http.MethodGet, "/jobs?status=running&kind=verification&group_id=5&page=1&page_size=25", "", nil, service.RoleAdmin)
	h.ListJobs(c)
	require.Equal(t, 200, w.Code, w.Body.String())
	require.Equal(t, service.ProviderHallJobFilter{Status: service.ProviderHallJobRunning, Kind: service.ProviderHallJobVerification, GroupID: 5}, jobs.filter)
	require.Equal(t, 25, jobs.pageSize)
	require.Contains(t, w.Body.String(), `"total":1`)
	body := w.Body.String()
	for _, secret := range []string{"sk-", "Authorization", "probe_key\"", "api_key\""} {
		require.NotContains(t, body, secret)
	}
	for _, query := range []string{"status=bogus", "kind=nope", "group_id=-1", "page_size=101"} {
		w, c = providerHallTaskRequest(http.MethodGet, "/jobs?"+query, "", nil, service.RoleAdmin)
		h.ListJobs(c)
		require.Equal(t, 400, w.Code, query)
	}

	w, c = providerHallTaskRequest(http.MethodGet, "/jobs/11", "", gin.Params{{Key: "id", Value: "11"}}, service.RoleAdmin)
	h.GetJob(c)
	require.Equal(t, 200, w.Code, w.Body.String())
	var detail struct {
		Data struct {
			Job          service.ProviderHallJob           `json:"job"`
			Samples      []service.ProviderHallSample      `json:"samples"`
			Verification *service.ProviderHallVerification `json:"verification"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &detail))
	require.Equal(t, int64(11), detail.Data.Job.ID)
	require.Len(t, detail.Data.Samples, 1)
	require.NotNil(t, detail.Data.Verification)
	require.Equal(t, "insufficient", detail.Data.Verification.Verdict)

	w, c = providerHallTaskRequest(http.MethodGet, "/jobs/12", "", gin.Params{{Key: "id", Value: "12"}}, service.RoleAdmin)
	h.GetJob(c)
	require.Equal(t, 404, w.Code)

	// Without a job repository (api-only instance) listing is empty, not an error.
	h = providerHallTaskHandler(nil, nil, service.ProviderHallConfig{}, nil)
	w, c = providerHallTaskRequest(http.MethodGet, "/jobs", "", nil, service.RoleAdmin)
	h.ListJobs(c)
	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), `"items":[]`)
}

func TestProviderHallAdminHandler_Cancel(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
		reason string
	}{
		{"cancelled", nil, 200, ""},
		{"in_flight", service.ErrProviderHallJobInFlight, 409, "PROVIDER_HALL_JOB_IN_FLIGHT"},
		{"finished", service.ErrProviderHallJobNotCancellable, 409, "PROVIDER_HALL_JOB_NOT_CANCELLABLE"},
		{"missing", service.ErrProviderHallNotFound, 404, "PROVIDER_HALL_NOT_FOUND"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runner := &providerHallRunnerFake{cancelErr: tc.err}
			h := providerHallTaskHandler(runner, nil, service.ProviderHallConfig{}, nil)
			w, c := providerHallTaskRequest(http.MethodPost, "/jobs/11/cancel", "", gin.Params{{Key: "id", Value: "11"}}, service.RoleAdmin)
			h.CancelJob(c)
			require.Equal(t, tc.status, w.Code, w.Body.String())
			require.Equal(t, []int64{11}, runner.cancelled)
			if tc.reason != "" {
				require.Contains(t, w.Body.String(), `"reason":"`+tc.reason+`"`)
			} else {
				require.Contains(t, w.Body.String(), `"status":"cancelled"`)
			}
		})
	}
	h := providerHallTaskHandler(nil, nil, service.ProviderHallConfig{}, nil)
	w, c := providerHallTaskRequest(http.MethodPost, "/jobs/abc/cancel", "", gin.Params{{Key: "id", Value: "abc"}}, service.RoleAdmin)
	h.CancelJob(c)
	require.Equal(t, 400, w.Code)
}

func TestProviderHallAdminHandler_Health(t *testing.T) {
	now := time.Date(2026, 9, 12, 3, 0, 0, 0, time.UTC)
	confirmed := now.Add(-30 * time.Second)
	watermark := now.Add(-2 * time.Minute)
	exited := now.Add(-time.Hour)
	admin := &providerHallAdminRepoFake{
		epochs: []service.ProviderHallNodeEpoch{
			{NodeID: "node-a", EpochID: 1, BuildVersion: "v1.2", HeartbeatAt: now.Add(-5 * time.Second), ConfirmedAt: &confirmed, PersistedSeq: 99},
			{NodeID: "node-b", EpochID: 2, BuildVersion: "v1.2", HeartbeatAt: now.Add(-2 * time.Minute), Overflowed: true},
			{NodeID: "node-c", EpochID: 3, BuildVersion: "v1.1", HeartbeatAt: exited, ExitedAt: &exited, ExitReason: "shutdown"},
		},
		gaps:  []service.ProviderHallOpenGap{{ID: 5, NodeID: "node-b", EpochID: 2, Scope: "collection", StartedAt: now.Add(-90 * time.Second), Reason: "queue_overflow"}},
		state: service.ProviderHallAggregatorState{LastWindowEnd: &watermark, LastRunAt: &confirmed, LastError: ""},
		dirty: 4,
		rec:   service.ProviderHallReconciliationCounts{Pending: 2, Uncertain: 1, Failed24h: 3},
	}
	jobs := &providerHallJobReadsFake{
		counts: service.ProviderHallJobCounts{Queued: 1, Running: 1, Unknown: 0, Failed24h: 2},
		spend:  service.ProviderHallSpendSummary{Confirmed: decimal.RequireFromString("12.5"), Uncertain: decimal.RequireFromString("0.25"), InFlight: 1},
	}
	cfg := service.ProviderHallConfig{CollectionEnabled: true, TasksEnabled: true, DailyBudget: "10", ExpectedNodes: []string{"node-a", "node-b", "node-z"}}
	h := providerHallTaskHandler(nil, jobs, cfg, admin)
	h.queue = providerHallQueueFake{}
	h.aggregator = providerHallAggFake{at: now.Add(-10 * time.Second), err: "lock miss"}
	w, c := providerHallTaskRequest(http.MethodGet, "/health", "", nil, service.RoleAdmin)
	h.Health(c)
	require.Equal(t, 200, w.Code, w.Body.String())
	var envelope struct {
		Data struct {
			Collection struct {
				Enabled bool `json:"enabled"`
				Nodes   []struct {
					NodeID string `json:"node_id"`
					Lost   bool   `json:"lost"`
				} `json:"nodes"`
				MissingExpected []string `json:"missing_expected"`
				OpenGaps        []struct {
					Reason string `json:"reason"`
				} `json:"open_gaps"`
				LocalQueue struct {
					Depth    int    `json:"depth"`
					Capacity int    `json:"capacity"`
					Dropped  uint64 `json:"dropped"`
				} `json:"local_queue"`
			} `json:"collection"`
			Aggregator struct {
				Watermark  *time.Time `json:"watermark"`
				LagSeconds *int64     `json:"lag_seconds"`
				DirtyCount int        `json:"dirty_count"`
				LastRunAt  *time.Time `json:"last_run_at"`
				LastError  string     `json:"last_error"`
			} `json:"aggregator"`
			Reconciliation struct {
				Pending   int `json:"pending"`
				Uncertain int `json:"uncertain"`
				Failed24h int `json:"failed_24h"`
			} `json:"reconciliation"`
			Budget struct {
				Day            string `json:"day"`
				Budget         string `json:"budget"`
				ConfirmedSpend string `json:"confirmed_spend"`
				UncertainSpend string `json:"uncertain_spend"`
				InFlight       int    `json:"in_flight"`
				PausedReason   string `json:"paused_reason"`
			} `json:"budget"`
			Jobs struct {
				Queued    int `json:"queued"`
				Running   int `json:"running"`
				Failed24h int `json:"failed_24h"`
			} `json:"jobs"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &envelope))
	d := envelope.Data
	require.True(t, d.Collection.Enabled)
	require.Len(t, d.Collection.Nodes, 3)
	require.False(t, d.Collection.Nodes[0].Lost)
	require.True(t, d.Collection.Nodes[1].Lost, "missed heartbeat shows as lost before the aggregator closes the epoch")
	require.False(t, d.Collection.Nodes[2].Lost, "a clean exit is not a loss")
	require.Equal(t, []string{"node-b", "node-z"}, d.Collection.MissingExpected)
	require.Len(t, d.Collection.OpenGaps, 1)
	require.Equal(t, "queue_overflow", d.Collection.OpenGaps[0].Reason)
	require.Equal(t, 3, d.Collection.LocalQueue.Depth)
	require.Equal(t, 8192, d.Collection.LocalQueue.Capacity)
	require.Equal(t, uint64(1), d.Collection.LocalQueue.Dropped)
	require.NotNil(t, d.Aggregator.Watermark)
	require.NotNil(t, d.Aggregator.LagSeconds)
	require.Equal(t, int64(120), *d.Aggregator.LagSeconds)
	require.Equal(t, 4, d.Aggregator.DirtyCount)
	require.Equal(t, "lock miss", d.Aggregator.LastError, "fresher in-process status wins")
	require.Equal(t, 2, d.Reconciliation.Pending)
	require.Equal(t, 1, d.Reconciliation.Uncertain)
	require.Equal(t, 3, d.Reconciliation.Failed24h)
	require.Equal(t, "2026-09-12", d.Budget.Day, "Shanghai day of 03:00 UTC")
	require.Equal(t, "10", d.Budget.Budget)
	require.Equal(t, "12.5", d.Budget.ConfirmedSpend)
	require.Equal(t, "0.25", d.Budget.UncertainSpend)
	require.Equal(t, 1, d.Budget.InFlight)
	require.Equal(t, "budget_exhausted", d.Budget.PausedReason)
	require.Equal(t, 1, d.Jobs.Queued)
	require.Equal(t, 2, d.Jobs.Failed24h)
	require.NotContains(t, w.Body.String(), "probe_key")

	// Tasks disabled takes precedence over budget; a bare build still answers.
	cfg.TasksEnabled = false
	h = providerHallTaskHandler(nil, nil, cfg, nil)
	w, c = providerHallTaskRequest(http.MethodGet, "/health", "", nil, service.RoleAdmin)
	h.Health(c)
	require.Equal(t, 200, w.Code, w.Body.String())
	require.Contains(t, w.Body.String(), `"paused_reason":"tasks_disabled"`)
	require.Contains(t, w.Body.String(), `"nodes":[]`)
	require.Contains(t, w.Body.String(), `"missing_expected":["node-a","node-b","node-z"]`)
}
