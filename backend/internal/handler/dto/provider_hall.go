package dto

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type ProviderHallConfigInput struct {
	Version             *int64   `json:"version"`
	CollectionEnabled   bool     `json:"collection_enabled"`
	DisplayEnabled      bool     `json:"display_enabled"`
	TasksEnabled        bool     `json:"tasks_enabled"`
	AutoScheduleEnabled *bool    `json:"auto_schedule_enabled"`
	DefaultModel        string   `json:"default_model"`
	DefaultProtocol     string   `json:"default_protocol"`
	DefaultRange        string   `json:"default_range"`
	GatewayOrigin       string   `json:"gateway_origin"`
	OperatorUserID      *int64   `json:"operator_user_id"`
	DailyBudget         string   `json:"daily_budget"`
	ExpectedNodes       []string `json:"expected_nodes"`
}

// ProviderHallReadiness tells the admin page which runtime switches this
// build can actually turn on.
type ProviderHallReadiness struct {
	Collection bool `json:"collection"`
	Tasks      bool `json:"tasks"`
	Display    bool `json:"display"`
}

// ProviderHallConfigView is the GET/PUT config response: the stored config
// plus readiness.
type ProviderHallConfigView struct {
	*service.ProviderHallConfig
	Readiness ProviderHallReadiness `json:"readiness"`
}

type ProviderHallGroupInput struct {
	Version      *int64 `json:"version"`
	Listed       bool   `json:"listed"`
	DisplayName  string `json:"display_name"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
}

type ProviderHallTargetInput struct {
	ProfileID                   int64  `json:"profile_id"`
	ProbeKeyID                  *int64 `json:"probe_key_id"`
	Enabled                     bool   `json:"enabled"`
	AutoScheduleEnabled         *bool  `json:"auto_schedule_enabled"`
	ProbeIntervalSeconds        int    `json:"probe_interval_seconds"`
	VerificationIntervalSeconds int    `json:"verification_interval_seconds"`
}

type ProviderHallTargetSetInput struct {
	Version *int64                     `json:"version"`
	Items   *[]ProviderHallTargetInput `json:"items"`
}

type ProviderHallProfileInput struct {
	Version              *int64     `json:"version"`
	Model                string     `json:"model"`
	Protocol             string     `json:"protocol"`
	SupportsTools        bool       `json:"supports_tools"`
	OutputLimit          int        `json:"output_limit"`
	ModelAliases         []string   `json:"model_aliases"`
	ReferenceInputPrice  *string    `json:"reference_input_price"`
	ReferenceCachePrice  *string    `json:"reference_cache_price"`
	ReferenceCacheRate   *string    `json:"reference_cache_rate"`
	ReferenceConfirmedAt *time.Time `json:"reference_confirmed_at"`
}

type ProviderHallMetricState string

const (
	ProviderHallMetricOK            ProviderHallMetricState = "ok"
	ProviderHallMetricInsufficient  ProviderHallMetricState = "insufficient"
	ProviderHallMetricIncomplete    ProviderHallMetricState = "incomplete"
	ProviderHallMetricStale         ProviderHallMetricState = "stale"
	ProviderHallMetricDisabled      ProviderHallMetricState = "disabled"
	ProviderHallMetricNotApplicable ProviderHallMetricState = "not_applicable"
)

// Unavailable values and unpublished timestamps remain JSON null, never zero.
type ProviderHallMetric[T any] struct {
	Value       *T                      `json:"value"`
	State       ProviderHallMetricState `json:"state"`
	ReasonCode  string                  `json:"reason_code"`
	WindowStart *time.Time              `json:"window_start"`
	WindowEnd   *time.Time              `json:"window_end"`
	ComputedAt  *time.Time              `json:"computed_at"`
}

// Admin task API (batch B6). Enqueue bodies and health read model; job,
// sample and report rows reuse the service types (no key material).

type ProviderHallEnqueueInput struct {
	TargetVersion  *int64 `json:"target_version"`
	ProfileVersion *int64 `json:"profile_version"`
	ProfileID      int64  `json:"profile_id"`
	IdempotencyKey string `json:"idempotency_key"`
}

type ProviderHallEnqueueResponse struct {
	JobID  int64  `json:"job_id"`
	Status string `json:"status"`
	Reused bool   `json:"reused"`
}

type ProviderHallHealthNode struct {
	NodeID       string     `json:"node_id"`
	EpochID      int64      `json:"epoch_id"`
	Version      string     `json:"version"`
	HeartbeatAt  time.Time  `json:"heartbeat_at"`
	ConfirmedAt  *time.Time `json:"confirmed_at"`
	PersistedSeq int64      `json:"persisted_seq"`
	Overflowed   bool       `json:"overflowed"`
	Lost         bool       `json:"lost"`
}

type ProviderHallHealthGap struct {
	ID        int64     `json:"id"`
	NodeID    string    `json:"node_id"`
	EpochID   int64     `json:"epoch_id"`
	Scope     string    `json:"scope"`
	StartedAt time.Time `json:"started_at"`
	Reason    string    `json:"reason"`
}

type ProviderHallHealthQueue struct {
	Depth    int    `json:"depth"`
	Capacity int    `json:"capacity"`
	Dropped  uint64 `json:"dropped"`
}

type ProviderHallHealthCollection struct {
	Enabled         bool                     `json:"enabled"`
	Nodes           []ProviderHallHealthNode `json:"nodes"`
	MissingExpected []string                 `json:"missing_expected"`
	OpenGaps        []ProviderHallHealthGap  `json:"open_gaps"`
	LocalQueue      ProviderHallHealthQueue  `json:"local_queue"`
}

type ProviderHallHealthAggregator struct {
	Watermark  *time.Time `json:"watermark"`
	LagSeconds *int64     `json:"lag_seconds"`
	DirtyCount int        `json:"dirty_count"`
	LastRunAt  *time.Time `json:"last_run_at"`
	LastError  string     `json:"last_error"`
}

type ProviderHallHealthReconciliation struct {
	Pending   int `json:"pending"`
	Uncertain int `json:"uncertain"`
	Failed24h int `json:"failed_24h"`
}

type ProviderHallHealthBudget struct {
	Day            string `json:"day"`
	Budget         string `json:"budget"`
	ConfirmedSpend string `json:"confirmed_spend"`
	UncertainSpend string `json:"uncertain_spend"`
	InFlight       int    `json:"in_flight"`
	PausedReason   string `json:"paused_reason"`
}

type ProviderHallHealthJobs struct {
	Queued    int `json:"queued"`
	Running   int `json:"running"`
	Unknown   int `json:"unknown"`
	Failed24h int `json:"failed_24h"`
}

type ProviderHallAdminHealth struct {
	TasksEnabled        bool                             `json:"tasks_enabled"`
	AutoScheduleEnabled bool                             `json:"auto_schedule_enabled"`
	GeneratedAt         time.Time                        `json:"generated_at"`
	Collection          ProviderHallHealthCollection     `json:"collection"`
	Aggregator          ProviderHallHealthAggregator     `json:"aggregator"`
	Reconciliation      ProviderHallHealthReconciliation `json:"reconciliation"`
	Budget              ProviderHallHealthBudget         `json:"budget"`
	Jobs                ProviderHallHealthJobs           `json:"jobs"`
}
