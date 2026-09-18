package service

import (
	"time"

	"github.com/shopspring/decimal"
)

// Metric states and reason codes shared by the user API (B5) and admin views.
// The state machine is fixed: disabled → incomplete → insufficient → stale → ok.
const (
	ProviderHallMetricOK            = "ok"
	ProviderHallMetricInsufficient  = "insufficient"
	ProviderHallMetricIncomplete    = "incomplete"
	ProviderHallMetricStale         = "stale"
	ProviderHallMetricDisabled      = "disabled"
	ProviderHallMetricNotApplicable = "not_applicable"

	ProviderHallReasonCollectionDisabled = "collection_disabled"
	ProviderHallReasonTargetDisabled     = "target_disabled"
	ProviderHallReasonGroupUnlisted      = "group_unlisted"
	ProviderHallReasonSnapshotMissing    = "snapshot_missing"
	ProviderHallReasonSamplesBelow       = "samples_below_20"
	ProviderHallReasonNoUsage            = "no_usage"
	ProviderHallReasonBillsBelow         = "bills_below_200"
	ProviderHallReasonStale              = "snapshot_stale"

	ProviderHallCoverageComplete        = "complete"
	ProviderHallCoverageCollectionGap   = "collection_gap"
	ProviderHallCoverageNodeUnconfirmed = "node_unconfirmed"
	ProviderHallCoverageVersionMismatch = "version_mismatch"
	ProviderHallBillingCoverageComplete = "complete"
	ProviderHallBillingCoverageGap      = "billing_gap"
	ProviderHallBillingCoveragePending  = "pending"

	// ProviderHallMinQualitySamples is the minimum number of samples before a
	// quality metric (TTFT, success rate, cache rate) is published.
	ProviderHallMinQualitySamples = 20
	// ProviderHallMinBilledSamples is the minimum number of confirmed bills
	// before a historical price is published.
	ProviderHallMinBilledSamples = 200
	// ProviderHallSnapshotWindow is the metric window length.
	ProviderHallSnapshotWindow = time.Hour
	// ProviderHallSnapshotStaleAfter marks a snapshot stale once its window end
	// is older than this.
	ProviderHallSnapshotStaleAfter = 5 * time.Minute
)

// ProviderHallMetric mirrors dto.ProviderHallMetric without importing the dto
// package (dto imports service). Only ok/stale carry a value.
type ProviderHallMetric[T any] struct {
	Value       *T
	State       string
	ReasonCode  string
	WindowStart *time.Time
	WindowEnd   *time.Time
	ComputedAt  *time.Time
}

// ProviderHallSnapshotRow is one provider_hall_snapshots row. profile_id 0 is
// the merged row over the group's enabled profiles.
type ProviderHallSnapshotRow struct {
	ID                  int64
	GroupID             int64
	ProfileID           int64
	WindowEnd           time.Time
	Tier                int
	AlgorithmVersion    int
	TTFTSampleCount     int64
	TTFTFast95MeanMs    *decimal.Decimal
	TTFTP90Ms           *int64
	SuccessCount        int64
	FailedCount         int64
	Submissions         int64
	SuccessRate         *decimal.Decimal
	UsageSuccessCount   int64
	InputTokens         int64
	CacheReadTokens     int64
	CacheCreationTokens int64
	CacheRate           *decimal.Decimal
	BilledCount         int64
	BilledInputTokens   int64
	BilledInputCost     decimal.Decimal
	HistoricalPrice     *decimal.Decimal
	SubscriptionShare   *decimal.Decimal
	Coverage            string
	BillingCoverage     string
	CoverageReason      string
	ComputedAt          time.Time
	ComputedVersion     int
}

// ProviderHallMetricInput carries everything the mapping needs. Snapshot is
// nil when no row exists for the requested window.
type ProviderHallMetricInput struct {
	Snapshot          *ProviderHallSnapshotRow
	Now               time.Time
	CollectionEnabled bool
	TargetEnabled     bool
	GroupListed       bool
}

func providerHallMetricWindow[T any](m *ProviderHallMetric[T], s *ProviderHallSnapshotRow) {
	if s == nil {
		return
	}
	start, end, computed := s.WindowEnd.Add(-ProviderHallSnapshotWindow), s.WindowEnd, s.ComputedAt
	m.WindowStart, m.WindowEnd, m.ComputedAt = &start, &end, &computed
}

// providerHallMetricGate applies the shared prefix of the state machine
// (disabled → incomplete → snapshot missing). It returns ok=false with the
// metric already filled when the caller must stop.
func providerHallMetricGate[T any](in ProviderHallMetricInput, billingAware bool) (ProviderHallMetric[T], bool) {
	m := ProviderHallMetric[T]{State: ProviderHallMetricOK}
	providerHallMetricWindow(&m, in.Snapshot)
	switch {
	case !in.CollectionEnabled:
		m.State, m.ReasonCode = ProviderHallMetricDisabled, ProviderHallReasonCollectionDisabled
	case !in.TargetEnabled:
		m.State, m.ReasonCode = ProviderHallMetricDisabled, ProviderHallReasonTargetDisabled
	case !in.GroupListed:
		m.State, m.ReasonCode = ProviderHallMetricDisabled, ProviderHallReasonGroupUnlisted
	}
	if m.State != ProviderHallMetricOK {
		return m, false
	}
	if in.Snapshot == nil {
		m.State, m.ReasonCode = ProviderHallMetricInsufficient, ProviderHallReasonSnapshotMissing
		return m, false
	}
	if in.Snapshot.Coverage != ProviderHallCoverageComplete {
		m.State, m.ReasonCode = ProviderHallMetricIncomplete, in.Snapshot.Coverage
		return m, false
	}
	if billingAware && in.Snapshot.BillingCoverage != ProviderHallBillingCoverageComplete {
		m.State, m.ReasonCode = ProviderHallMetricIncomplete, "billing_"+in.Snapshot.BillingCoverage
		if in.Snapshot.BillingCoverage == ProviderHallBillingCoverageGap {
			m.ReasonCode = ProviderHallBillingCoverageGap
		}
		return m, false
	}
	return m, true
}

func providerHallMetricFinish[T any](m *ProviderHallMetric[T], in ProviderHallMetricInput, value *T) {
	m.Value = value
	if in.Now.Sub(in.Snapshot.WindowEnd) > ProviderHallSnapshotStaleAfter {
		m.State, m.ReasonCode = ProviderHallMetricStale, ProviderHallReasonStale
		return
	}
	m.State, m.ReasonCode = ProviderHallMetricOK, ""
}

// ProviderHallTTFTMetric maps the merged latency figures. Both metrics share
// the sample-count gate so fast95 and P90 never disagree on availability.
func ProviderHallTTFTMetric(in ProviderHallMetricInput) (fast95 ProviderHallMetric[float64], p90 ProviderHallMetric[int64]) {
	fast95, ok := providerHallMetricGate[float64](in, false)
	p90, _ = providerHallMetricGate[int64](in, false)
	if !ok {
		return fast95, p90
	}
	s := in.Snapshot
	if s.TTFTSampleCount < ProviderHallMinQualitySamples || s.TTFTFast95MeanMs == nil || s.TTFTP90Ms == nil {
		fast95.State, fast95.ReasonCode = ProviderHallMetricInsufficient, ProviderHallReasonSamplesBelow
		p90.State, p90.ReasonCode = ProviderHallMetricInsufficient, ProviderHallReasonSamplesBelow
		return fast95, p90
	}
	mean, _ := s.TTFTFast95MeanMs.Round(3).Float64()
	p90v := *s.TTFTP90Ms
	providerHallMetricFinish(&fast95, in, &mean)
	providerHallMetricFinish(&p90, in, &p90v)
	return fast95, p90
}

func ProviderHallSuccessMetric(in ProviderHallMetricInput) ProviderHallMetric[float64] {
	m, ok := providerHallMetricGate[float64](in, false)
	if !ok {
		return m
	}
	s := in.Snapshot
	if s.SuccessCount+s.FailedCount < ProviderHallMinQualitySamples || s.SuccessRate == nil {
		m.State, m.ReasonCode = ProviderHallMetricInsufficient, ProviderHallReasonSamplesBelow
		return m
	}
	v, _ := s.SuccessRate.Round(10).Float64()
	providerHallMetricFinish(&m, in, &v)
	return m
}

func ProviderHallCacheMetric(in ProviderHallMetricInput) ProviderHallMetric[float64] {
	m, ok := providerHallMetricGate[float64](in, false)
	if !ok {
		return m
	}
	s := in.Snapshot
	switch {
	case s.UsageSuccessCount < ProviderHallMinQualitySamples:
		m.State, m.ReasonCode = ProviderHallMetricInsufficient, ProviderHallReasonSamplesBelow
		return m
	case s.CacheRate == nil:
		m.State, m.ReasonCode = ProviderHallMetricInsufficient, ProviderHallReasonNoUsage
		return m
	}
	v, _ := s.CacheRate.Round(10).Float64()
	providerHallMetricFinish(&m, in, &v)
	return m
}

// ProviderHallPriceMetric maps the historical input price (per million input
// tokens) as a decimal string. It is billing-coverage aware.
func ProviderHallPriceMetric(in ProviderHallMetricInput) ProviderHallMetric[string] {
	m, ok := providerHallMetricGate[string](in, true)
	if !ok {
		return m
	}
	s := in.Snapshot
	switch {
	case s.BilledCount < ProviderHallMinBilledSamples:
		m.State, m.ReasonCode = ProviderHallMetricInsufficient, ProviderHallReasonBillsBelow
		return m
	case s.HistoricalPrice == nil:
		m.State, m.ReasonCode = ProviderHallMetricInsufficient, ProviderHallReasonNoUsage
		return m
	}
	v := s.HistoricalPrice.StringFixed(10)
	providerHallMetricFinish(&m, in, &v)
	return m
}

// ProviderHallTrendBucket returns the bucket size and point count for a
// display range. Real curves read tier-5 snapshots at bucket ends.
func ProviderHallTrendBucket(rangeKey string) (bucketSeconds, points int, ok bool) {
	switch rangeKey {
	case "6h":
		return 300, 72, true
	case "24h":
		return 900, 96, true
	case "7d":
		return 3600, 168, true
	case "30d":
		return 21600, 120, true
	}
	return 0, 0, false
}

// ProviderHallSnapshotTier is 5 when the window end sits on a 5-minute mark.
func ProviderHallSnapshotTier(windowEnd time.Time) int {
	if windowEnd.UTC().Minute()%5 == 0 {
		return 5
	}
	return 1
}

// ProviderHallCoverageWorst merges per-minute coverage into a window value.
func ProviderHallCoverageWorst(values ...string) string {
	rank := map[string]int{ProviderHallCoverageComplete: 0, ProviderHallCoverageVersionMismatch: 1, ProviderHallCoverageNodeUnconfirmed: 2, ProviderHallCoverageCollectionGap: 3}
	worst := ProviderHallCoverageComplete
	for _, v := range values {
		if rank[v] > rank[worst] {
			worst = v
		}
	}
	return worst
}

func ProviderHallBillingCoverageWorst(values ...string) string {
	rank := map[string]int{ProviderHallBillingCoverageComplete: 0, ProviderHallBillingCoveragePending: 1, ProviderHallBillingCoverageGap: 2}
	worst := ProviderHallBillingCoverageComplete
	for _, v := range values {
		if rank[v] > rank[worst] {
			worst = v
		}
	}
	return worst
}

// ProviderHallMinuteRow is one provider_hall_metrics_1m row.
type ProviderHallMinuteRow struct {
	GroupID                 int64
	ProfileID               int64
	Minute                  time.Time
	SuccessCount            int64
	FailedCount             int64
	Submissions             int64
	ExcludedCount           int64
	TTFTSampleCount         int64
	UsageSuccessCount       int64
	InputTokens             int64
	CacheReadTokens         int64
	CacheCreationTokens     int64
	OutputTokens            int64
	BilledCount             int64
	BilledSubscriptionCount int64
	BilledInputTokens       int64
	BilledInputCost         decimal.Decimal
	BillingPendingCount     int64
	BillingUncertainCount   int64
	Coverage                string
	BillingCoverage         string
}

// ProviderHallBuildSnapshot folds the minute rows of one window and the
// merged exact latency frequencies into a snapshot. Minutes absent from the
// window carry no expectation and do not degrade coverage.
func ProviderHallBuildSnapshot(groupID, profileID int64, windowEnd time.Time, minutes []ProviderHallMinuteRow, latency map[int64]int64, algorithmVersion int, computedAt time.Time) ProviderHallSnapshotRow {
	s := ProviderHallSnapshotRow{
		GroupID:          groupID,
		ProfileID:        profileID,
		WindowEnd:        windowEnd,
		Tier:             ProviderHallSnapshotTier(windowEnd),
		AlgorithmVersion: algorithmVersion,
		BilledInputCost:  decimal.Zero,
		Coverage:         ProviderHallCoverageComplete,
		BillingCoverage:  ProviderHallBillingCoverageComplete,
		ComputedAt:       computedAt,
		ComputedVersion:  1,
	}
	var billedSubscription int64
	coverages := make([]string, 0, len(minutes))
	billingCoverages := make([]string, 0, len(minutes))
	for _, m := range minutes {
		s.SuccessCount += m.SuccessCount
		s.FailedCount += m.FailedCount
		s.Submissions += m.Submissions
		s.TTFTSampleCount += m.TTFTSampleCount
		s.UsageSuccessCount += m.UsageSuccessCount
		s.InputTokens += m.InputTokens
		s.CacheReadTokens += m.CacheReadTokens
		s.CacheCreationTokens += m.CacheCreationTokens
		s.BilledCount += m.BilledCount
		billedSubscription += m.BilledSubscriptionCount
		s.BilledInputTokens += m.BilledInputTokens
		s.BilledInputCost = s.BilledInputCost.Add(m.BilledInputCost)
		coverages = append(coverages, m.Coverage)
		billingCoverages = append(billingCoverages, m.BillingCoverage)
	}
	lat := ProviderHallExactLatency(latency)
	if lat.Count > 0 {
		s.TTFTSampleCount = lat.Count
		s.TTFTFast95MeanMs = lat.Fast95MeanMS
		s.TTFTP90Ms = lat.P90MS
	}
	s.SuccessRate = ProviderHallSuccessRate(s.SuccessCount, s.FailedCount, s.Submissions)
	s.CacheRate = ProviderHallCacheRate(s.InputTokens-s.CacheReadTokens-s.CacheCreationTokens, s.CacheReadTokens, s.CacheCreationTokens)
	if s.BilledCount > 0 {
		s.HistoricalPrice = ProviderHallInputPrice(s.BilledInputCost, s.BilledInputTokens)
		share := decimal.NewFromInt(billedSubscription).Div(decimal.NewFromInt(s.BilledCount))
		s.SubscriptionShare = &share
	}
	s.Coverage = ProviderHallCoverageWorst(coverages...)
	s.BillingCoverage = ProviderHallBillingCoverageWorst(billingCoverages...)
	if s.Coverage != ProviderHallCoverageComplete {
		s.CoverageReason = s.Coverage
	} else if s.BillingCoverage != ProviderHallBillingCoverageComplete {
		s.CoverageReason = s.BillingCoverage
	}
	return s
}
