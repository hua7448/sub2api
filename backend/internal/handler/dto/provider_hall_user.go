package dto

import (
	"encoding/json"
	"time"
)

// User-facing provider hall contract (batch B5). Field names are the JSON
// contract shared with the frontend (src/types/providerHall.ts). Nothing here
// may ever carry account IDs, credentials, upstream addresses, request bodies
// or absolute real-traffic counts; sample sizes are only expressed through
// Metric.state / reason_code.

// ProviderHallProfileRef identifies one model profile without exposing its
// reference pricing.
type ProviderHallProfileRef struct {
	ProfileID int64  `json:"profile_id"`
	Model     string `json:"model"`
	Protocol  string `json:"protocol"`
}

// ProviderHallQuote is the current quoted price for the user on the group's
// default profile. Prices are decimal strings per million tokens.
type ProviderHallQuote struct {
	InputPrice string `json:"input_price"`
	CachePrice string `json:"cache_price"`
	Unit       string `json:"unit"` // usd_per_million | quota_per_million
	Applicable bool   `json:"applicable"`
	ReasonCode string `json:"reason_code"` // pricing_unavailable | tiered_pricing | zero_price | ""
}

type ProviderHallModelHealth struct {
	ProfileID int64      `json:"profile_id"`
	Model     string     `json:"model"`
	Protocol  string     `json:"protocol"`
	Status    string     `json:"status"` // up | down | unknown
	CheckedAt *time.Time `json:"checked_at"`
}

type ProviderHallHealth struct {
	Status    string                    `json:"status"` // up | down | unknown
	CheckedAt *time.Time                `json:"checked_at"`
	Models    []ProviderHallModelHealth `json:"models"`
}

type ProviderHallRowMetrics struct {
	TTFTFast95Ms ProviderHallMetric[float64] `json:"ttft_fast95_ms"`
	TTFTP90Ms    ProviderHallMetric[int64]   `json:"ttft_p90_ms"`
	CacheRate    ProviderHallMetric[float64] `json:"cache_rate"`
	SuccessRate  ProviderHallMetric[float64] `json:"success_rate"`
}

type ProviderHallSparkPoint struct {
	T          time.Time `json:"t"`
	RealTTFTMs *float64  `json:"real_ttft_ms"`
	ProbeMs    *float64  `json:"probe_ms"`
	Gap        bool      `json:"gap"`
}

type ProviderHallSparkline struct {
	Range         string                   `json:"range"`
	BucketSeconds int                      `json:"bucket_seconds"`
	Points        []ProviderHallSparkPoint `json:"points"`
}

type ProviderHallVerificationBadge struct {
	Verdict     *string    `json:"verdict"` // passed | failed | suspected | insufficient
	ReasonCode  string     `json:"reason_code"`
	CompletedAt *time.Time `json:"completed_at"`
	Expired     bool       `json:"expired"`
	ReportID    *int64     `json:"report_id"`
}

type ProviderHallRow struct {
	GroupID         int64                         `json:"group_id"`
	Name            string                        `json:"name"`
	Description     string                        `json:"description"`
	DisplayOrder    int                           `json:"display_order"`
	Rate            string                        `json:"rate"` // effective multiplier, decimal string
	Quote           ProviderHallQuote             `json:"quote"`
	HistoricalPrice ProviderHallMetric[string]    `json:"historical_price"`
	PredictedRate   ProviderHallMetric[string]    `json:"predicted_rate"`
	Health          ProviderHallHealth            `json:"health"`
	Metrics         ProviderHallRowMetrics        `json:"metrics"`
	Sparkline       ProviderHallSparkline         `json:"sparkline"`
	Verification    ProviderHallVerificationBadge `json:"verification"`
	DefaultProfile  ProviderHallProfileRef        `json:"default_profile"`
}

type ProviderHallCatalog struct {
	Models  []ProviderHallProfileRef `json:"models"`
	Default ProviderHallProfileRef   `json:"default"`
}

type ProviderHallSummary struct {
	Available int `json:"available"`
	Listed    int `json:"listed"`
	Abnormal  int `json:"abnormal"`
	Verified  int `json:"verified"`
}

type ProviderHallPagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type ProviderHallListResponse struct {
	SnapshotID    string                 `json:"snapshot_id"`
	DataThrough   *time.Time             `json:"data_through"`
	MetricVersion int                    `json:"metric_version"`
	PricingAt     time.Time              `json:"pricing_at"`
	Range         string                 `json:"range"`
	Catalog       ProviderHallCatalog    `json:"catalog"`
	Summary       ProviderHallSummary    `json:"summary"`
	Items         []ProviderHallRow      `json:"items"`
	Pagination    ProviderHallPagination `json:"pagination"`
}

type ProviderHallProbeTokens struct {
	Input  int       `json:"input"`
	Output int       `json:"output"`
	At     time.Time `json:"at"`
}

type ProviderHallDetailMetrics struct {
	E2EAvailability6h ProviderHallMetric[float64] `json:"e2e_availability_6h"`
	TPS               ProviderHallMetric[float64] `json:"tps"`
	ProbeTTFTMs       ProviderHallMetric[int64]   `json:"probe_ttft_ms"`
	ProbeTotalMs      ProviderHallMetric[int64]   `json:"probe_total_ms"`
	P90Ms             ProviderHallMetric[int64]   `json:"p90_ms"`
	ProbeTokens       *ProviderHallProbeTokens    `json:"probe_tokens"`
	LastProbeAt       *time.Time                  `json:"last_probe_at"`
}

type ProviderHallTrendPoint struct {
	T            time.Time `json:"t"`
	RealFast95Ms *float64  `json:"real_fast95_ms"`
	ProbeTotalMs *float64  `json:"probe_total_ms"`
	ProbeTTFTMs  *float64  `json:"probe_ttft_ms"`
	Gap          bool      `json:"gap"`
	ProbeFailed  int       `json:"probe_failed"`
}

type ProviderHallTrend struct {
	Range         string                   `json:"range"`
	BucketSeconds int                      `json:"bucket_seconds"`
	Points        []ProviderHallTrendPoint `json:"points"`
}

type ProviderHallProfileStatus struct {
	ProviderHallProfileRef
	Health       ProviderHallModelHealth       `json:"health"`
	Verification ProviderHallVerificationBadge `json:"verification"`
}

type ProviderHallDetailResponse struct {
	ProviderHallRow
	Detail   ProviderHallDetailMetrics   `json:"detail"`
	Trend    ProviderHallTrend           `json:"trend"`
	Profiles []ProviderHallProfileStatus `json:"profiles"`
}

type ProviderHallVerificationReport struct {
	ReportID        int64           `json:"report_id"`
	ProfileID       int64           `json:"profile_id"`
	Model           string          `json:"model"`
	Protocol        string          `json:"protocol"`
	Verdict         string          `json:"verdict"`
	ExecutionStatus string          `json:"execution_status"`
	ReasonCode      string          `json:"reason_code"`
	CompletedAt     time.Time       `json:"completed_at"`
	ExpiresAt       time.Time       `json:"expires_at"`
	Expired         bool            `json:"expired"`
	Summary         json.RawMessage `json:"summary"`
}

type ProviderHallVerificationListResponse struct {
	Items      []ProviderHallVerificationReport `json:"items"`
	Pagination ProviderHallPagination           `json:"pagination"`
}
