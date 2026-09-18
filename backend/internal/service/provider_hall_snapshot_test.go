//go:build unit

package service

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func hallSnapshotFixture(windowEnd time.Time) *ProviderHallSnapshotRow {
	mean := decimal.RequireFromString("10000")
	p90 := int64(18000)
	rate := decimal.RequireFromString("0.95")
	cache := decimal.RequireFromString("0.25")
	price := decimal.RequireFromString("2.5")
	return &ProviderHallSnapshotRow{
		GroupID: 1, ProfileID: 2, WindowEnd: windowEnd, Tier: 5, AlgorithmVersion: 1,
		TTFTSampleCount: 20, TTFTFast95MeanMs: &mean, TTFTP90Ms: &p90,
		SuccessCount: 19, FailedCount: 1, Submissions: 21, SuccessRate: &rate,
		UsageSuccessCount: 20, InputTokens: 4000, CacheReadTokens: 1000, CacheRate: &cache,
		BilledCount: 200, BilledInputTokens: 40000, BilledInputCost: decimal.RequireFromString("0.1"), HistoricalPrice: &price,
		Coverage: ProviderHallCoverageComplete, BillingCoverage: ProviderHallBillingCoverageComplete, ComputedAt: windowEnd,
	}
}

func TestProviderHallMetricStatePriority(t *testing.T) {
	now := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	fresh := hallSnapshotFixture(now.Add(-time.Minute))
	cases := []struct {
		name          string
		mutate        func(in *ProviderHallMetricInput)
		state, reason string
		priceState    string
		priceReason   string
	}{
		{"ok", nil, ProviderHallMetricOK, "", ProviderHallMetricOK, ""},
		{"collection_disabled_beats_everything", func(in *ProviderHallMetricInput) {
			in.CollectionEnabled = false
			in.Snapshot.Coverage = ProviderHallCoverageCollectionGap
			in.Snapshot.TTFTSampleCount = 1
		}, ProviderHallMetricDisabled, ProviderHallReasonCollectionDisabled, ProviderHallMetricDisabled, ProviderHallReasonCollectionDisabled},
		{"target_disabled", func(in *ProviderHallMetricInput) { in.TargetEnabled = false }, ProviderHallMetricDisabled, ProviderHallReasonTargetDisabled, ProviderHallMetricDisabled, ProviderHallReasonTargetDisabled},
		{"group_unlisted", func(in *ProviderHallMetricInput) { in.GroupListed = false }, ProviderHallMetricDisabled, ProviderHallReasonGroupUnlisted, ProviderHallMetricDisabled, ProviderHallReasonGroupUnlisted},
		{"gap_beats_insufficient_and_stale", func(in *ProviderHallMetricInput) {
			in.Snapshot.Coverage = ProviderHallCoverageCollectionGap
			in.Snapshot.TTFTSampleCount = 1
			in.Snapshot.WindowEnd = now.Add(-time.Hour)
		}, ProviderHallMetricIncomplete, ProviderHallCoverageCollectionGap, ProviderHallMetricIncomplete, ProviderHallCoverageCollectionGap},
		{"node_unconfirmed", func(in *ProviderHallMetricInput) { in.Snapshot.Coverage = ProviderHallCoverageNodeUnconfirmed }, ProviderHallMetricIncomplete, ProviderHallCoverageNodeUnconfirmed, ProviderHallMetricIncomplete, ProviderHallCoverageNodeUnconfirmed},
		{"billing_gap_only_affects_price", func(in *ProviderHallMetricInput) { in.Snapshot.BillingCoverage = ProviderHallBillingCoverageGap }, ProviderHallMetricOK, "", ProviderHallMetricIncomplete, ProviderHallBillingCoverageGap},
		{"billing_pending_only_affects_price", func(in *ProviderHallMetricInput) { in.Snapshot.BillingCoverage = ProviderHallBillingCoveragePending }, ProviderHallMetricOK, "", ProviderHallMetricIncomplete, "billing_pending"},
		{"insufficient_beats_stale", func(in *ProviderHallMetricInput) {
			in.Snapshot.TTFTSampleCount = 19
			in.Snapshot.BilledCount = 199
			in.Snapshot.WindowEnd = now.Add(-time.Hour)
		}, ProviderHallMetricInsufficient, ProviderHallReasonSamplesBelow, ProviderHallMetricInsufficient, ProviderHallReasonBillsBelow},
		{"stale", func(in *ProviderHallMetricInput) { in.Snapshot.WindowEnd = now.Add(-6 * time.Minute) }, ProviderHallMetricStale, ProviderHallReasonStale, ProviderHallMetricStale, ProviderHallReasonStale},
		{"five_minutes_is_not_stale", func(in *ProviderHallMetricInput) { in.Snapshot.WindowEnd = now.Add(-5 * time.Minute) }, ProviderHallMetricOK, "", ProviderHallMetricOK, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			snap := *fresh
			in := ProviderHallMetricInput{Snapshot: &snap, Now: now, CollectionEnabled: true, TargetEnabled: true, GroupListed: true}
			if tc.mutate != nil {
				tc.mutate(&in)
			}
			fast95, p90 := ProviderHallTTFTMetric(in)
			require.Equal(t, tc.state, fast95.State)
			require.Equal(t, tc.reason, fast95.ReasonCode)
			require.Equal(t, tc.state, p90.State)
			price := ProviderHallPriceMetric(in)
			require.Equal(t, tc.priceState, price.State)
			require.Equal(t, tc.priceReason, price.ReasonCode)
			hasValue := tc.state == ProviderHallMetricOK || tc.state == ProviderHallMetricStale
			require.Equal(t, hasValue, fast95.Value != nil, "only ok/stale carry values")
			require.Equal(t, hasValue, p90.Value != nil)
			if hasValue {
				require.Equal(t, 10000.0, *fast95.Value)
				require.Equal(t, int64(18000), *p90.Value)
				require.NotNil(t, fast95.WindowStart)
				require.Equal(t, in.Snapshot.WindowEnd.Add(-time.Hour), *fast95.WindowStart)
			}
			if tc.priceState == ProviderHallMetricOK {
				require.Equal(t, "2.5000000000", *price.Value)
			}
		})
	}

	t.Run("snapshot_missing", func(t *testing.T) {
		in := ProviderHallMetricInput{Now: now, CollectionEnabled: true, TargetEnabled: true, GroupListed: true}
		m := ProviderHallSuccessMetric(in)
		require.Equal(t, ProviderHallMetricInsufficient, m.State)
		require.Equal(t, ProviderHallReasonSnapshotMissing, m.ReasonCode)
		require.Nil(t, m.WindowEnd)
	})
	t.Run("success_and_cache_denominators", func(t *testing.T) {
		snap := *fresh
		in := ProviderHallMetricInput{Snapshot: &snap, Now: now, CollectionEnabled: true, TargetEnabled: true, GroupListed: true}
		snap.SuccessCount, snap.FailedCount = 10, 9
		require.Equal(t, ProviderHallReasonSamplesBelow, ProviderHallSuccessMetric(in).ReasonCode)
		snap.SuccessCount = 11
		require.Equal(t, ProviderHallMetricOK, ProviderHallSuccessMetric(in).State)
		snap.UsageSuccessCount = 19
		require.Equal(t, ProviderHallReasonSamplesBelow, ProviderHallCacheMetric(in).ReasonCode)
		snap.UsageSuccessCount, snap.CacheRate = 20, nil
		require.Equal(t, ProviderHallReasonNoUsage, ProviderHallCacheMetric(in).ReasonCode)
	})
}

func TestProviderHallBuildSnapshotMatrixA(t *testing.T) {
	// Matrix A: 1000..19000 plus one 100000 outlier → fast95 mean 10000, P90 18000.
	latency := map[int64]int64{}
	for ms := int64(1000); ms <= 19000; ms += 1000 {
		latency[ms] = 1
	}
	latency[100000] = 1
	end := time.Date(2026, 9, 12, 10, 5, 0, 0, time.UTC)
	minutes := []ProviderHallMinuteRow{
		{Minute: end.Add(-2 * time.Minute), SuccessCount: 10, Submissions: 11, TTFTSampleCount: 10, UsageSuccessCount: 10, InputTokens: 1000, CacheReadTokens: 200, CacheCreationTokens: 100, BilledCount: 10, BilledInputTokens: 1000, BilledInputCost: decimal.RequireFromString("0.002"), Coverage: ProviderHallCoverageComplete, BillingCoverage: ProviderHallBillingCoverageComplete},
		{Minute: end.Add(-time.Minute), SuccessCount: 10, FailedCount: 2, Submissions: 13, TTFTSampleCount: 10, UsageSuccessCount: 10, InputTokens: 1000, CacheReadTokens: 200, BilledCount: 10, BilledSubscriptionCount: 5, BilledInputTokens: 1000, BilledInputCost: decimal.RequireFromString("0.002"), Coverage: ProviderHallCoverageComplete, BillingCoverage: ProviderHallBillingCoveragePending},
	}
	s := ProviderHallBuildSnapshot(7, 0, end, minutes, latency, 1, end)
	require.Equal(t, 5, s.Tier)
	require.Equal(t, int64(20), s.TTFTSampleCount)
	require.Equal(t, "10000", s.TTFTFast95MeanMs.String())
	require.Equal(t, int64(18000), *s.TTFTP90Ms)
	require.Equal(t, int64(20), s.SuccessCount)
	require.Equal(t, int64(2), s.FailedCount)
	require.Equal(t, int64(24), s.Submissions)
	require.Equal(t, "0.8333333333", s.SuccessRate.StringFixed(10))
	require.Equal(t, "0.2", s.CacheRate.String())
	require.Equal(t, int64(20), s.BilledCount)
	require.Equal(t, "2", s.HistoricalPrice.String(), "$0.004 over 2000 tokens = $2 per million")
	require.Equal(t, "0.25", s.SubscriptionShare.String())
	require.Equal(t, ProviderHallCoverageComplete, s.Coverage)
	require.Equal(t, ProviderHallBillingCoveragePending, s.BillingCoverage)
	require.Equal(t, ProviderHallBillingCoveragePending, s.CoverageReason)

	empty := ProviderHallBuildSnapshot(7, 3, end.Add(time.Minute), nil, nil, 1, end)
	require.Equal(t, 1, empty.Tier)
	require.Nil(t, empty.TTFTFast95MeanMs)
	require.Nil(t, empty.SuccessRate)
	require.Nil(t, empty.CacheRate)
	require.Nil(t, empty.HistoricalPrice)
	require.Nil(t, empty.SubscriptionShare)
	require.Equal(t, ProviderHallCoverageComplete, empty.Coverage)

	worst := ProviderHallBuildSnapshot(7, 0, end, []ProviderHallMinuteRow{{Coverage: ProviderHallCoverageVersionMismatch}, {Coverage: ProviderHallCoverageCollectionGap}, {Coverage: ProviderHallCoverageNodeUnconfirmed}}, nil, 1, end)
	require.Equal(t, ProviderHallCoverageCollectionGap, worst.Coverage)
	require.Equal(t, ProviderHallCoverageCollectionGap, worst.CoverageReason)
}

func TestProviderHallTrendBucket(t *testing.T) {
	for _, tc := range []struct {
		key            string
		bucket, points int
	}{{"6h", 300, 72}, {"24h", 900, 96}, {"7d", 3600, 168}, {"30d", 21600, 120}} {
		b, p, ok := ProviderHallTrendBucket(tc.key)
		require.True(t, ok)
		require.Equal(t, tc.bucket, b)
		require.Equal(t, tc.points, p)
		require.Equal(t, tc.bucket*tc.points, map[string]int{"6h": 6 * 3600, "24h": 24 * 3600, "7d": 7 * 86400, "30d": 30 * 86400}[tc.key])
	}
	_, _, ok := ProviderHallTrendBucket("1h")
	require.False(t, ok)
}
