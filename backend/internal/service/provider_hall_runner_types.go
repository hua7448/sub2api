package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Job, sample, report and spend types shared by the runner (B4), the admin
// API (B6) and the user API (B5). Nothing here carries key material.

type ProviderHallJobKind string

const (
	ProviderHallJobProbe        ProviderHallJobKind = "probe"
	ProviderHallJobVerification ProviderHallJobKind = "verification"
)

type ProviderHallJobStatus string

const (
	ProviderHallJobQueued    ProviderHallJobStatus = "queued"
	ProviderHallJobRunning   ProviderHallJobStatus = "running"
	ProviderHallJobSucceeded ProviderHallJobStatus = "succeeded"
	ProviderHallJobFailed    ProviderHallJobStatus = "failed"
	ProviderHallJobCancelled ProviderHallJobStatus = "cancelled"
	ProviderHallJobUnknown   ProviderHallJobStatus = "unknown"
)

type ProviderHallSampleStatus string

const (
	ProviderHallSamplePrepared   ProviderHallSampleStatus = "prepared"
	ProviderHallSampleDispatched ProviderHallSampleStatus = "dispatched"
	ProviderHallSampleReceived   ProviderHallSampleStatus = "received"
	ProviderHallSampleUncertain  ProviderHallSampleStatus = "uncertain"
)

// Sample results.
const (
	ProviderHallResultPassed        = "passed"
	ProviderHallResultFailed        = "failed"
	ProviderHallResultError         = "error"
	ProviderHallResultModelMismatch = "model_mismatch"
)

// Job error / cancel codes.
const (
	ProviderHallJobCodeTargetDisabled   = "target_disabled"
	ProviderHallJobCodeTasksDisabled    = "tasks_disabled"
	ProviderHallJobCodeConfigChanged    = "config_changed"
	ProviderHallJobCodeProbeKeyInvalid  = "probe_key_invalid"
	ProviderHallJobCodeBillingBacklog   = "billing_backlog"
	ProviderHallJobCodeBudgetExhausted  = "budget_exhausted"
	ProviderHallJobCodeGatewayOrigin    = "gateway_origin_invalid"
	ProviderHallJobCodeSampleUncertain  = "sample_uncertain"
	ProviderHallJobCodeProbeFailed      = "probe_failed"
	ProviderHallJobCodeLeaseLost        = "lease_lost"
	ProviderHallJobCodeCancelledByAdmin = "cancelled_by_admin"
	ProviderHallJobCodeRunnerStopped    = "runner_stopped"
)

// ProviderHallJobProfile is the profile as it was when the job was created.
type ProviderHallJobProfile struct {
	ID            int64    `json:"id"`
	Version       int64    `json:"version"`
	Model         string   `json:"model"`
	Protocol      string   `json:"protocol"`
	SupportsTools bool     `json:"supports_tools"`
	OutputLimit   int      `json:"output_limit"`
	ModelAliases  []string `json:"model_aliases"`
}

// ProviderHallJobSnapshot is stored as config_snapshot. Dispatch re-validates
// the live rows against it and cancels on any version drift.
type ProviderHallJobSnapshot struct {
	Profile        ProviderHallJobProfile `json:"profile"`
	TargetVersion  int64                  `json:"target_version"`
	ProbeKeyID     int64                  `json:"probe_key_id"`
	OperatorUserID int64                  `json:"operator_user_id"`
	GatewayOrigin  string                 `json:"gateway_origin"`
}

type ProviderHallJob struct {
	GroupName      string                  `json:"group_name"`
	Source         string                  `json:"source"`
	ID             int64                   `json:"id"`
	Kind           ProviderHallJobKind     `json:"kind"`
	TargetID       int64                   `json:"target_id"`
	GroupID        int64                   `json:"group_id"`
	ProfileID      int64                   `json:"profile_id"`
	Snapshot       ProviderHallJobSnapshot `json:"config_snapshot"`
	SlotAt         *time.Time              `json:"slot_at"`
	NotBefore      *time.Time              `json:"not_before"`
	IdempotencyKey *string                 `json:"idempotency_key"`
	RequestedBy    *int64                  `json:"requested_by"`
	Status         ProviderHallJobStatus   `json:"status"`
	LeaseOwner     *string                 `json:"lease_owner"`
	LeaseUntil     *time.Time              `json:"lease_until"`
	Attempts       int                     `json:"attempts"`
	BudgetDay      *string                 `json:"budget_day"`
	ErrorCode      string                  `json:"error_code"`
	ErrorMessage   string                  `json:"error_message"`
	CreatedAt      time.Time               `json:"created_at"`
	StartedAt      *time.Time              `json:"started_at"`
	FinishedAt     *time.Time              `json:"finished_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

type ProviderHallSample struct {
	ID              int64                    `json:"id"`
	JobID           int64                    `json:"job_id"`
	TestID          string                   `json:"test_id"`
	Seq             int                      `json:"seq"`
	TraceID         uuid.UUID                `json:"trace_id"`
	Status          ProviderHallSampleStatus `json:"status"`
	PreparedAt      time.Time                `json:"prepared_at"`
	DispatchedAt    *time.Time               `json:"dispatched_at"`
	ReceivedAt      *time.Time               `json:"received_at"`
	ClientRequestID *string                  `json:"client_request_id"`
	HTTPStatus      *int                     `json:"http_status"`
	TTFTMs          *int                     `json:"ttft_ms"`
	TotalMs         *int                     `json:"total_ms"`
	GenerationMs    *int                     `json:"generation_ms"`
	InputTokens     *int                     `json:"input_tokens"`
	OutputTokens    *int                     `json:"output_tokens"`
	ResponseModel   *string                  `json:"response_model"`
	Result          *string                  `json:"result"`
	ErrorCode       string                   `json:"error_code"`
	Detail          json.RawMessage          `json:"detail"`
	BillingStatus   string                   `json:"billing_status"`
	ActualCost      *decimal.Decimal         `json:"actual_cost"`
}

type ProviderHallVerification struct {
	JobID           int64                           `json:"job_id"`
	GroupID         int64                           `json:"group_id"`
	ProfileID       int64                           `json:"profile_id"`
	TargetID        int64                           `json:"target_id"`
	Verdict         string                          `json:"verdict"`
	ExecutionStatus string                          `json:"execution_status"`
	ReasonCode      string                          `json:"reason_code"`
	Summary         ProviderHallVerificationSummary `json:"summary"`
	ProfileVersion  int64                           `json:"profile_version"`
	TargetVersion   int64                           `json:"target_version"`
	CompletedAt     time.Time                       `json:"completed_at"`
	ExpiresAt       time.Time                       `json:"expires_at"`
	Stale           bool                            `json:"stale"`
	CreatedAt       time.Time                       `json:"created_at"`
}

// Expired reports whether the report may no longer be shown as current.
func (v *ProviderHallVerification) Expired(now time.Time) bool {
	return v == nil || v.Stale || !now.Before(v.ExpiresAt)
}

// ProviderHallProbeHealth is the latest received probe sample of one target.
type ProviderHallProbeHealth struct {
	TargetID     int64      `json:"target_id"`
	GroupID      int64      `json:"group_id"`
	ProfileID    int64      `json:"profile_id"`
	JobID        int64      `json:"job_id"`
	SampleID     int64      `json:"sample_id"`
	Result       string     `json:"result"`
	ErrorCode    string     `json:"error_code"`
	ReceivedAt   time.Time  `json:"received_at"`
	TTFTMs       *int       `json:"ttft_ms"`
	TotalMs      *int       `json:"total_ms"`
	GenerationMs *int       `json:"generation_ms"`
	InputTokens  *int       `json:"input_tokens"`
	OutputTokens *int       `json:"output_tokens"`
	DispatchedAt *time.Time `json:"dispatched_at"`
}

// ProviderHallProbePoint is one trend bucket of probe samples.
type ProviderHallProbePoint struct {
	BucketStart time.Time `json:"bucket_start"`
	Total       int       `json:"total"`
	Passed      int       `json:"passed"`
	Failed      int       `json:"failed"`
	AvgTotalMs  *float64  `json:"avg_total_ms"`
	AvgTTFTMs   *float64  `json:"avg_ttft_ms"`
}

type ProviderHallProbeSeriesQuery struct {
	GroupID       int64
	ProfileID     *int64 // nil = every enabled profile of the group
	Start, End    time.Time
	BucketSeconds int
}

type ProviderHallEnqueueInput struct {
	Kind           ProviderHallJobKind
	TargetID       int64
	GroupID        int64
	ProfileID      int64
	Snapshot       ProviderHallJobSnapshot
	SlotAt         *time.Time
	NotBefore      *time.Time
	IdempotencyKey *string
	RequestedBy    *int64
}

type ProviderHallJobFilter struct {
	GroupName string
	Model     string
	Source    string
	ProfileID int64
	Status    ProviderHallJobStatus
	Kind      ProviderHallJobKind
	GroupID   int64
}

// ProviderHallSchedulableTarget is an enabled target of a listed group with
// its live profile, as seen by the scheduler.
type ProviderHallSchedulableTarget struct {
	TargetID                    int64
	GroupID                     int64
	ProfileID                   int64
	ProbeKeyID                  *int64
	TargetVersion               int64
	ProbeIntervalSeconds        int
	VerificationIntervalSeconds int
	AutoScheduleEnabled         bool
	Listed                      bool
	Enabled                     bool
	Profile                     ProviderHallJobProfile
}

// ProviderHallDispatchDecision is the outcome of the pre-dispatch check.
type ProviderHallDispatchDecision string

const (
	ProviderHallDispatchGo      ProviderHallDispatchDecision = "dispatch"
	ProviderHallDispatchCancel  ProviderHallDispatchDecision = "cancel"
	ProviderHallDispatchFail    ProviderHallDispatchDecision = "fail"
	ProviderHallDispatchRequeue ProviderHallDispatchDecision = "requeue"
	ProviderHallDispatchWait    ProviderHallDispatchDecision = "wait"
)

type ProviderHallDispatchOutcome struct {
	Decision ProviderHallDispatchDecision
	Code     string
	Message  string
	// Key is the raw probe key, present only for Decision == dispatch. It is
	// used for one Authorization header and never stored or logged.
	Key       string
	BudgetDay string
	Origin    string
}

type ProviderHallDispatchInput struct {
	Job          *ProviderHallJob
	SampleID     int64
	Now          time.Time
	TasksEnabled bool
	// MaxInflight is the global dispatched-sample cap (2).
	MaxInflight int
	// BacklogAfter is how long a pending/uncertain bill may block dispatch.
	BacklogAfter time.Duration
}

type ProviderHallSampleReceipt struct {
	SampleID        int64
	ReceivedAt      time.Time
	ClientRequestID string
	HTTPStatus      int
	TTFTMs          *int
	TotalMs         *int
	GenerationMs    *int
	InputTokens     *int
	OutputTokens    *int
	ResponseModel   string
	Result          string
	ErrorCode       string
	Detail          map[string]any
}

type ProviderHallJobCounts struct {
	Queued    int `json:"queued"`
	Running   int `json:"running"`
	Unknown   int `json:"unknown"`
	Failed24h int `json:"failed_24h"`
}

type ProviderHallSpendSummary struct {
	Day       string          `json:"day"`
	Confirmed decimal.Decimal `json:"confirmed"`
	Uncertain decimal.Decimal `json:"uncertain"`
	InFlight  int             `json:"in_flight"`
}

// ProviderHallJobRepository persists jobs, samples, reports and spend.
type ProviderHallJobRepository interface {
	ListSchedulableTargets(ctx context.Context) ([]ProviderHallSchedulableTarget, error)
	EnqueueSlot(ctx context.Context, in ProviderHallEnqueueInput) (int64, bool, error)
	EnqueueManual(ctx context.Context, in ProviderHallEnqueueInput) (*ProviderHallJob, bool, error)
	CancelQueuedNotIn(ctx context.Context, activeTargetIDs []int64, reason string, now time.Time) (int64, error)
	CancelQueuedForTarget(ctx context.Context, targetID int64, reason string, now time.Time) (int64, error)
	CancelQueuedDrifted(ctx context.Context, now time.Time) (int64, error)
	RecoverExpiredLeases(ctx context.Context, now time.Time) (requeued, unknown int64, err error)
	Claim(ctx context.Context, owner string, now time.Time, lease time.Duration) (*ProviderHallJob, error)
	RenewLease(ctx context.Context, jobID int64, owner string, until time.Time) (bool, error)
	Requeue(ctx context.Context, jobID int64, notBefore time.Time, code string) error
	Complete(ctx context.Context, jobID int64, status ProviderHallJobStatus, code, message string, now time.Time) error
	MarkUnknown(ctx context.Context, jobID int64, code string, now time.Time) error
	Cancel(ctx context.Context, jobID int64, reason string, now time.Time) error
	CreateSamples(ctx context.Context, jobID int64, cases []ProviderHallTestCase) error
	ListSamples(ctx context.Context, jobID int64) ([]ProviderHallSample, error)
	TryDispatch(ctx context.Context, in ProviderHallDispatchInput) (ProviderHallDispatchOutcome, error)
	MarkSampleReceived(ctx context.Context, r ProviderHallSampleReceipt) error
	MarkSampleUncertain(ctx context.Context, sampleID int64, clientRequestID, code string, now time.Time) error
	ReconcileSamples(ctx context.Context, now time.Time, uncertainTimeout time.Duration) (int64, error)
	ListInflight(ctx context.Context) ([]ProviderHallSample, error)
	ListFinalizableJobs(ctx context.Context, now time.Time) ([]ProviderHallJob, error)
	SaveVerification(ctx context.Context, v ProviderHallVerification) error
	MarkVerificationsStale(ctx context.Context, profileID int64) (int64, error)
	MarkDriftedVerificationsStale(ctx context.Context) (int64, error)
	SumSpend(ctx context.Context, day string) (ProviderHallSpendSummary, error)
	ListJobs(ctx context.Context, f ProviderHallJobFilter, page, pageSize int) ([]ProviderHallJob, int, error)
	GetJob(ctx context.Context, id int64) (*ProviderHallJob, error)
	GetVerification(ctx context.Context, jobID int64) (*ProviderHallVerification, error)
	CountJobs(ctx context.Context, now time.Time) (ProviderHallJobCounts, error)
	LatestProbePerTarget(ctx context.Context, since time.Time) (map[int64]ProviderHallProbeHealth, error)
	ProbeSeries(ctx context.Context, q ProviderHallProbeSeriesQuery) ([]ProviderHallProbePoint, error)
	RecentProbeSamples(ctx context.Context, groupID int64, profileID *int64, since time.Time, limit int) ([]ProviderHallProbeHealth, error)
	LatestVerification(ctx context.Context, groupID, profileID int64) (*ProviderHallVerification, error)
	ListVerifications(ctx context.Context, groupID int64, profileID *int64, page, pageSize int) ([]ProviderHallVerification, int, error)
}

// ProviderHallJobControl is the small surface the configuration service uses
// to react to admin changes immediately; the runner also converges on its own.
type ProviderHallJobControl interface {
	CancelQueuedForTarget(ctx context.Context, targetID int64, reason string, now time.Time) (int64, error)
	MarkVerificationsStale(ctx context.Context, profileID int64) (int64, error)
}
