//go:build unit

package service

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// Fakes for the hall read model. Everything is in memory and deterministic.

type hallCfgRepo struct {
	ProviderHallRepository
	cfg      ProviderHallConfig
	listings map[int64]ProviderHallGroup
	profiles []ProviderHallProfile
}

func (r *hallCfgRepo) GetConfig(context.Context) (*ProviderHallConfig, error) {
	c := r.cfg
	return &c, nil
}
func (r *hallCfgRepo) GetGroup(_ context.Context, id int64) (*ProviderHallGroup, error) {
	g, ok := r.listings[id]
	if !ok {
		return nil, ErrProviderHallNotFound
	}
	return &g, nil
}
func (r *hallCfgRepo) ListProfiles(context.Context) ([]ProviderHallProfile, error) {
	return r.profiles, nil
}

type hallAccess struct {
	mu     sync.Mutex
	groups map[int64][]Group
	calls  int
}

func (a *hallAccess) GetAvailableGroups(_ context.Context, userID int64) ([]Group, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.calls++
	return append([]Group(nil), a.groups[userID]...), nil
}

type hallRead struct {
	mu      sync.Mutex
	groups  []ProviderHallListedGroup
	targets []ProviderHallReadTarget
	window  *time.Time
	algo    int
	snaps   []ProviderHallSnapshotRow
	series  []ProviderHallSeriesPoint
	reports []ProviderHallVerification
	probes  []ProviderHallGroupProbePoint
	calls   map[string]int
}

func (r *hallRead) count(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.calls == nil {
		r.calls = map[string]int{}
	}
	r.calls[name]++
}
func (r *hallRead) ListListedGroups(context.Context) ([]ProviderHallListedGroup, error) {
	r.count("groups")
	return r.groups, nil
}
func (r *hallRead) ListEnabledTargets(context.Context) ([]ProviderHallReadTarget, error) {
	r.count("targets")
	return r.targets, nil
}
func (r *hallRead) LatestWindow(context.Context) (*time.Time, int, error) {
	r.count("window")
	return r.window, r.algo, nil
}
func (r *hallRead) LoadSnapshots(_ context.Context, at time.Time, _ int) ([]ProviderHallSnapshotRow, error) {
	r.count("snapshots")
	var out []ProviderHallSnapshotRow
	for _, s := range r.snaps {
		if s.WindowEnd.Equal(at) {
			out = append(out, s)
		}
	}
	return out, nil
}
func (r *hallRead) LoadSeries(context.Context, []time.Time, int) ([]ProviderHallSeriesPoint, error) {
	r.count("series")
	return r.series, nil
}
func (r *hallRead) LatestVerifications(context.Context) ([]ProviderHallVerification, error) {
	r.count("reports")
	return r.reports, nil
}
func (r *hallRead) ProbeSeriesAll(context.Context, time.Time, time.Time, int) ([]ProviderHallGroupProbePoint, error) {
	r.count("probe_series")
	return r.probes, nil
}

type hallRates struct {
	mu    sync.Mutex
	rates map[int64]map[int64]float64 // user → group → multiplier
}

func (r *hallRates) ResolveUserGroupRateMultiplier(_ context.Context, userID, groupID int64, def float64) float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	if m, ok := r.rates[userID][groupID]; ok {
		return m
	}
	return def
}

type hallPricing struct{ byModel map[string]*ResolvedPricing }

func (p *hallPricing) Resolve(_ context.Context, in PricingInput) *ResolvedPricing {
	return p.byModel[in.Model]
}

type hallProbes struct {
	latest  map[int64]ProviderHallProbeHealth
	recent  []ProviderHallProbeHealth
	series  []ProviderHallProbePoint
	reports []ProviderHallVerification
}

func (p *hallProbes) LatestProbePerTarget(context.Context, time.Time) (map[int64]ProviderHallProbeHealth, error) {
	return p.latest, nil
}
func (p *hallProbes) RecentProbeSamples(context.Context, int64, *int64, time.Time, int) ([]ProviderHallProbeHealth, error) {
	return p.recent, nil
}
func (p *hallProbes) ProbeSeries(context.Context, ProviderHallProbeSeriesQuery) ([]ProviderHallProbePoint, error) {
	return p.series, nil
}
func (p *hallProbes) ListVerifications(_ context.Context, groupID int64, profileID *int64, page, pageSize int) ([]ProviderHallVerification, int, error) {
	var out []ProviderHallVerification
	for _, v := range p.reports {
		if v.GroupID == groupID && (profileID == nil || v.ProfileID == *profileID) {
			out = append(out, v)
		}
	}
	return out, len(out), nil
}

var hallNow = time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)

func hallTokenPricing(input, cache float64) *ResolvedPricing {
	return &ResolvedPricing{Mode: BillingModeToken, BasePricing: &ModelPricing{InputPricePerToken: input, CacheReadPricePerToken: cache}}
}

func hallSnapshot(group, profile int64, at time.Time, mutate func(*ProviderHallSnapshotRow)) ProviderHallSnapshotRow {
	fast := decimal.NewFromInt(900)
	p90 := int64(1500)
	rate := decimal.NewFromFloat(0.95)
	cache := decimal.NewFromFloat(0.4)
	price := decimal.NewFromFloat(1.25)
	share := decimal.Zero
	s := ProviderHallSnapshotRow{
		GroupID: group, ProfileID: profile, WindowEnd: at, Tier: 5, AlgorithmVersion: 1,
		TTFTSampleCount: 50, TTFTFast95MeanMs: &fast, TTFTP90Ms: &p90,
		SuccessCount: 95, FailedCount: 5, Submissions: 100, SuccessRate: &rate,
		UsageSuccessCount: 90, InputTokens: 1000, CacheReadTokens: 400, CacheRate: &cache,
		BilledCount: 300, BilledInputTokens: 900, BilledInputCost: decimal.NewFromFloat(0.001125), HistoricalPrice: &price, SubscriptionShare: &share,
		Coverage: ProviderHallCoverageComplete, BillingCoverage: ProviderHallBillingCoverageComplete, ComputedAt: at,
	}
	if mutate != nil {
		mutate(&s)
	}
	return s
}

// hallFixture builds a hall with three listed groups (1 standard, 2 exclusive,
// 3 subscription), one unlisted (4) and one inactive (5).
type hallFixture struct {
	svc    *ProviderHallQueryService
	cfg    *hallCfgRepo
	access *hallAccess
	read   *hallRead
	rates  *hallRates
	probes *hallProbes
	window time.Time
}

func newHallFixture(t *testing.T) *hallFixture {
	t.Helper()
	window := hallNow.Add(-2 * time.Minute)
	ref := func(in, cache, rate string, age time.Duration) ProviderHallReference {
		i, _ := decimal.NewFromString(in)
		c, _ := decimal.NewFromString(cache)
		r, _ := decimal.NewFromString(rate)
		at := hallNow.Add(-age)
		return ProviderHallReference{InputPrice: &i, CachePrice: &c, CacheRate: &r, ConfirmedAt: &at}
	}
	cfg := &hallCfgRepo{
		cfg: ProviderHallConfig{CollectionEnabled: true, DisplayEnabled: true, DefaultModel: "gpt-5", DefaultProtocol: "responses", DefaultRange: "6h"},
		listings: map[int64]ProviderHallGroup{
			1: {GroupID: 1, Listed: true, DisplayName: "Alpha", DisplayOrder: 2},
			2: {GroupID: 2, Listed: true, DisplayName: "Beta", Description: "fast lane", DisplayOrder: 1},
			3: {GroupID: 3, Listed: true, DisplayName: "", DisplayOrder: 3},
			4: {GroupID: 4, Listed: false, DisplayName: "Hidden"},
			5: {GroupID: 5, Listed: true, DisplayName: "Dead"},
		},
		profiles: []ProviderHallProfile{{ID: 10, Model: "gpt-5", Protocol: "responses"}, {ID: 11, Model: "gpt-5-mini", Protocol: "chat_completions"}},
	}
	groups := []Group{
		{ID: 1, Name: "alpha-raw", Platform: PlatformOpenAI, Status: StatusActive, RateMultiplier: 1},
		{ID: 2, Name: "beta-raw", Platform: PlatformOpenAI, Status: StatusActive, RateMultiplier: 2, IsExclusive: true},
		{ID: 3, Name: "gamma-raw", Platform: PlatformComposite, Status: StatusActive, RateMultiplier: 1, SubscriptionType: SubscriptionTypeSubscription},
		{ID: 4, Name: "hidden-raw", Platform: PlatformOpenAI, Status: StatusActive, RateMultiplier: 1},
		{ID: 5, Name: "dead-raw", Platform: PlatformOpenAI, Status: "disabled", RateMultiplier: 1},
		{ID: 6, Name: "anthropic-raw", Platform: "anthropic", Status: StatusActive, RateMultiplier: 1},
	}
	access := &hallAccess{groups: map[int64][]Group{
		100: groups,                            // sees everything the permission API returns
		101: {groups[0], groups[3], groups[4]}, // no exclusive, no subscription
		102: {groups[1]},
		103: {groups[1]},
	}}
	read := &hallRead{
		groups: []ProviderHallListedGroup{
			{ProviderHallGroup: cfg.listings[2], GroupName: "beta-raw"},
			{ProviderHallGroup: cfg.listings[1], GroupName: "alpha-raw"},
			{ProviderHallGroup: cfg.listings[3], GroupName: "gamma-raw"},
			{ProviderHallGroup: cfg.listings[5], GroupName: "dead-raw"},
		},
		targets: []ProviderHallReadTarget{
			{TargetID: 1, GroupID: 1, ProfileID: 10, Model: "gpt-5", Protocol: "responses", ProbeIntervalSeconds: 300, Reference: ref("1.0", "0.1", "0.5", time.Hour)},
			{TargetID: 2, GroupID: 1, ProfileID: 11, Model: "gpt-5-mini", Protocol: "chat_completions", ProbeIntervalSeconds: 300},
			{TargetID: 3, GroupID: 2, ProfileID: 10, Model: "gpt-5", Protocol: "responses", ProbeIntervalSeconds: 300, Reference: ref("1.0", "0.1", "0.5", 8*24*time.Hour)},
			{TargetID: 4, GroupID: 3, ProfileID: 11, Model: "gpt-5-mini", Protocol: "chat_completions", ProbeIntervalSeconds: 300},
		},
		window: &window,
		algo:   1,
		snaps: []ProviderHallSnapshotRow{
			hallSnapshot(1, 0, window, nil),
			hallSnapshot(1, 10, window, nil),
			hallSnapshot(2, 0, window, func(s *ProviderHallSnapshotRow) { v := decimal.NewFromInt(400); s.TTFTFast95MeanMs = &v }),
			hallSnapshot(2, 10, window, func(s *ProviderHallSnapshotRow) { s.BilledCount = 10 }),
			hallSnapshot(3, 0, window, func(s *ProviderHallSnapshotRow) { s.TTFTSampleCount = 5 }),
			hallSnapshot(3, 11, window, func(s *ProviderHallSnapshotRow) { one := decimal.NewFromInt(1); s.SubscriptionShare = &one }),
		},
		reports: []ProviderHallVerification{
			{JobID: 900, GroupID: 1, ProfileID: 10, Verdict: "passed", CompletedAt: hallNow.Add(-time.Hour), ExpiresAt: hallNow.Add(47 * time.Hour)},
			{JobID: 901, GroupID: 2, ProfileID: 10, Verdict: "passed", CompletedAt: hallNow.Add(-72 * time.Hour), ExpiresAt: hallNow.Add(-24 * time.Hour)},
		},
	}
	rates := &hallRates{rates: map[int64]map[int64]float64{102: {2: 0.5}, 103: {2: 3}}}
	pricing := &hallPricing{byModel: map[string]*ResolvedPricing{
		"gpt-5":      hallTokenPricing(0.000001, 0.0000001), // $1 / $0.1 per million
		"gpt-5-mini": {Mode: BillingModeToken, BasePricing: &ModelPricing{InputPricePerToken: 0.0000002}, Intervals: []PricingInterval{{}}},
	}}
	probes := &hallProbes{latest: map[int64]ProviderHallProbeHealth{
		1: {TargetID: 1, Result: ProviderHallResultPassed, ReceivedAt: hallNow.Add(-time.Minute)},
		2: {TargetID: 2, Result: ProviderHallResultFailed, ReceivedAt: hallNow.Add(-time.Minute)},
		3: {TargetID: 3, Result: ProviderHallResultPassed, ReceivedAt: hallNow.Add(-time.Minute)},
		4: {TargetID: 4, Result: ProviderHallResultPassed, ReceivedAt: hallNow.Add(-2 * time.Hour)}, // older than 3 intervals
	}}
	hall := NewProviderHallService(cfg, nil, nil, nil)
	hall.access = access
	svc := NewProviderHallQueryService(read, cfg, hall, access, rates, pricing, probes)
	svc.now = func() time.Time { return hallNow }
	return &hallFixture{svc: svc, cfg: cfg, access: access, read: read, rates: rates, probes: probes, window: window}
}

func hallIDs(items []*ProviderHallRowResult) []int64 {
	out := make([]int64, 0, len(items))
	for _, it := range items {
		out = append(out, it.GroupID)
	}
	return out
}

func TestProviderHallQueryPermissionsAndVisibility(t *testing.T) {
	f := newHallFixture(t)
	ctx := context.Background()

	res, err := f.svc.List(ctx, 100, ProviderHallListQuery{})
	require.NoError(t, err)
	// Unlisted (4), inactive (5) and non-OpenAI (6) never appear; default order is display_order.
	require.Equal(t, []int64{2, 1, 3}, hallIDs(res.Items))
	require.Equal(t, 3, res.Summary.Listed)
	require.Equal(t, 3, res.Total)
	require.Equal(t, providerHallSnapshotIDFor(f.window), res.SnapshotID)
	require.Equal(t, f.window, *res.DataThrough)
	require.Equal(t, "6h", res.Range)
	require.Len(t, res.Catalog.Models, 2)
	require.Equal(t, int64(10), res.Catalog.Default.ProfileID)

	// User without exclusive / subscription access sees only group 1.
	res, err = f.svc.List(ctx, 101, ProviderHallListQuery{})
	require.NoError(t, err)
	require.Equal(t, []int64{1}, hallIDs(res.Items))
	// Detail and reports of an invisible group are 404, even though it is listed.
	_, err = f.svc.GetGroup(ctx, 101, 2, "")
	require.ErrorIs(t, err, ErrProviderHallNotFound)
	_, _, err = f.svc.ListVerifications(ctx, 101, 2, nil, 1, 10)
	require.ErrorIs(t, err, ErrProviderHallNotFound)
	// Unlisted group is 404 even for a user who may bind it.
	_, err = f.svc.GetGroup(ctx, 100, 4, "")
	require.ErrorIs(t, err, ErrProviderHallNotFound)
	// Unknown user sees nothing but no error.
	res, err = f.svc.List(ctx, 999, ProviderHallListQuery{})
	require.NoError(t, err)
	require.Empty(t, res.Items)

	// Display switch off: everything is 404.
	f.cfg.cfg.DisplayEnabled = false
	_, err = f.svc.List(ctx, 100, ProviderHallListQuery{})
	require.ErrorIs(t, err, ErrProviderHallNotFound)
	_, err = f.svc.GetGroup(ctx, 100, 1, "")
	require.ErrorIs(t, err, ErrProviderHallNotFound)
	require.False(t, f.svc.DisplayEnabled(ctx))
}

func TestProviderHallQueryRowContents(t *testing.T) {
	f := newHallFixture(t)
	ctx := context.Background()
	res, err := f.svc.List(ctx, 100, ProviderHallListQuery{})
	require.NoError(t, err)
	byID := map[int64]*ProviderHallRowResult{}
	for _, it := range res.Items {
		byID[it.GroupID] = it
	}
	g1 := byID[1]
	require.Equal(t, "Alpha", g1.Name)
	require.Equal(t, int64(10), g1.DefaultProfile.ProfileID, "config default wins when enabled")
	require.Equal(t, "1", g1.Rate)
	require.True(t, g1.Quote.Applicable)
	require.Equal(t, "1.0000000000", g1.Quote.InputPrice)
	require.Equal(t, "0.1000000000", g1.Quote.CachePrice)
	require.Equal(t, ProviderHallUnitUSDPerMillion, g1.Quote.Unit)
	require.Equal(t, ProviderHallMetricOK, g1.TTFTFast95Ms.State)
	require.Equal(t, 900.0, *g1.TTFTFast95Ms.Value)
	require.Equal(t, int64(1500), *g1.TTFTP90Ms.Value)
	require.Equal(t, 0.95, *g1.SuccessRate.Value)
	require.Equal(t, 0.4, *g1.CacheRate.Value)
	require.Equal(t, ProviderHallMetricOK, g1.HistoricalPrice.State)
	require.Equal(t, "1.2500000000", *g1.HistoricalPrice.Value)
	require.Equal(t, ProviderHallUnitUSDPerMillion, g1.HistoricalUnit)
	// predicted = (1×0.6 + 0.1×0.4) / (1×0.5 + 0.1×0.5) = 0.64 / 0.55
	require.Equal(t, ProviderHallMetricOK, g1.PredictedRate.State)
	require.Equal(t, "1.1636", *g1.PredictedRate.Value)
	require.Equal(t, ProviderHallHealthDown, g1.Health.Status, "one failing model marks the group down")
	require.Len(t, g1.Health.Models, 2)
	require.Equal(t, ProviderHallHealthUp, g1.Health.Models[0].Status)
	require.Equal(t, ProviderHallHealthDown, g1.Health.Models[1].Status)
	require.NotNil(t, g1.Verification.Verdict)
	require.Equal(t, "passed", *g1.Verification.Verdict)
	require.False(t, g1.Verification.Expired)
	require.Equal(t, int64(900), *g1.Verification.ReportID)
	require.Len(t, g1.Sparkline.Points, 72)
	require.Equal(t, 300, g1.Sparkline.BucketSeconds)

	g2 := byID[2]
	require.Equal(t, "2", g2.Rate, "group default multiplier when the user has no override")
	require.Equal(t, "2.0000000000", g2.Quote.InputPrice)
	require.Equal(t, ProviderHallMetricInsufficient, g2.HistoricalPrice.State)
	require.Equal(t, ProviderHallReasonBillsBelow, g2.HistoricalPrice.ReasonCode)
	require.Equal(t, ProviderHallMetricStale, g2.PredictedRate.State, "reference older than 7 days")
	require.Equal(t, ProviderHallReasonReferenceExpired, g2.PredictedRate.ReasonCode)
	require.True(t, g2.Verification.Expired)
	require.Equal(t, ProviderHallHealthUp, g2.Health.Status)

	g3 := byID[3]
	require.Equal(t, "gamma-raw", g3.Name, "empty display name falls back to the group name")
	require.Equal(t, ProviderHallUnitQuotaPerMillion, g3.Quote.Unit)
	require.False(t, g3.Quote.Applicable)
	require.Equal(t, ProviderHallReasonTieredPricing, g3.Quote.ReasonCode)
	require.Equal(t, ProviderHallMetricNotApplicable, g3.PredictedRate.State)
	require.Equal(t, ProviderHallReasonQuoteUnavailable, g3.PredictedRate.ReasonCode)
	require.Equal(t, ProviderHallMetricInsufficient, g3.TTFTFast95Ms.State)
	require.Equal(t, ProviderHallUnitQuotaPerMillion, g3.HistoricalUnit)
	require.Equal(t, ProviderHallHealthUnknown, g3.Health.Status, "probe older than three intervals is unknown")
	require.Nil(t, g3.Verification.Verdict)

	require.Equal(t, ProviderHallSummaryResult{Available: 1, Listed: 3, Abnormal: 1, Verified: 1}, res.Summary)
}

func TestProviderHallQueryExclusiveRatesNeverLeakAcrossUsers(t *testing.T) {
	f := newHallFixture(t)
	ctx := context.Background()
	var wg sync.WaitGroup
	errs := make(chan error, 200)
	for i := 0; i < 50; i++ {
		for _, tc := range []struct {
			user  int64
			rate  string
			price string
		}{{102, "0.5", "0.5000000000"}, {103, "3", "3.0000000000"}} {
			wg.Add(1)
			go func() {
				defer wg.Done()
				res, err := f.svc.List(ctx, tc.user, ProviderHallListQuery{})
				if err != nil {
					errs <- err
					return
				}
				if len(res.Items) != 1 || res.Items[0].Rate != tc.rate || res.Items[0].Quote.InputPrice != tc.price {
					errs <- fmt.Errorf("user %d saw %+v", tc.user, res.Items)
				}
			}()
		}
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
	// The user-independent base was loaded once and reused within the TTL.
	f.read.mu.Lock()
	defer f.read.mu.Unlock()
	require.Equal(t, 1, f.read.calls["snapshots"])
	require.Equal(t, 1, f.read.calls["groups"])
}

func TestProviderHallQuerySearchFilterAndPaging(t *testing.T) {
	f := newHallFixture(t)
	ctx := context.Background()
	res, err := f.svc.List(ctx, 100, ProviderHallListQuery{Search: "FAST"})
	require.NoError(t, err)
	require.Equal(t, []int64{2}, hallIDs(res.Items), "description search is case-insensitive")
	res, err = f.svc.List(ctx, 100, ProviderHallListQuery{Search: "gamma"})
	require.NoError(t, err)
	require.Equal(t, []int64{3}, hallIDs(res.Items), "original group name is searchable")
	res, err = f.svc.List(ctx, 100, ProviderHallListQuery{Model: "gpt-5-mini"})
	require.NoError(t, err)
	require.Equal(t, []int64{1, 3}, hallIDs(res.Items))
	res, err = f.svc.List(ctx, 100, ProviderHallListQuery{Model: "gpt-5-mini", Protocol: "responses"})
	require.NoError(t, err)
	require.Empty(t, res.Items)
	require.Len(t, res.Catalog.Models, 2, "catalog is independent of the filter")
	res, err = f.svc.List(ctx, 100, ProviderHallListQuery{Page: 2, PageSize: 2})
	require.NoError(t, err)
	require.Equal(t, []int64{3}, hallIDs(res.Items))
	require.Equal(t, 3, res.Total)
	require.Equal(t, 3, res.Summary.Listed, "summary covers the filtered set, not the page")

	for _, q := range []ProviderHallListQuery{{Range: "1h"}, {PageSize: 101}, {Protocol: "grpc"}, {Search: string(make([]rune, 101))}, {SnapshotID: "abc"}} {
		_, err := f.svc.List(ctx, 100, q)
		require.ErrorIs(t, err, ErrProviderHallInvalidQuery, "%+v", q)
	}
}

func TestProviderHallQuerySortMissingValuesLast(t *testing.T) {
	f := newHallFixture(t)
	ctx := context.Background()
	// Group 2 has the lowest fast95 (400), group 1 900, group 3 unpublished.
	rules, ok := ParseProviderHallSort("ttft_fast95:asc")
	require.True(t, ok)
	res, err := f.svc.List(ctx, 100, ProviderHallListQuery{Sort: rules})
	require.NoError(t, err)
	require.Equal(t, []int64{2, 1, 3}, hallIDs(res.Items))
	rules, _ = ParseProviderHallSort("ttft_fast95:desc")
	res, err = f.svc.List(ctx, 100, ProviderHallListQuery{Sort: rules})
	require.NoError(t, err)
	require.Equal(t, []int64{1, 2, 3}, hallIDs(res.Items), "missing stays last even descending")
	// Three levels: historical price (only 1 and 3 published, both 1.25) → rate desc → group id.
	rules, _ = ParseProviderHallSort("historical_price:asc,rate:desc,display_order")
	res, err = f.svc.List(ctx, 100, ProviderHallListQuery{Sort: rules})
	require.NoError(t, err)
	require.Equal(t, []int64{1, 3, 2}, hallIDs(res.Items))
	rules, _ = ParseProviderHallSort("predicted_rate:desc")
	res, err = f.svc.List(ctx, 100, ProviderHallListQuery{Sort: rules})
	require.NoError(t, err)
	require.Equal(t, []int64{2, 1, 3}, hallIDs(res.Items), "stale values still sort (group 2 quotes 2×); not_applicable is last")

	for _, bad := range []string{"name", "rate:up", "rate,rate", "rate,ttft_fast95,cache_rate,success_rate"} {
		_, ok := ParseProviderHallSort(bad)
		require.False(t, ok, bad)
	}
}

func TestProviderHallQueryDetailAndReports(t *testing.T) {
	f := newHallFixture(t)
	ctx := context.Background()
	in, out, gen, ttft, total := 12, 30, 600, 400, 1000
	f.probes.recent = []ProviderHallProbeHealth{
		{TargetID: 1, Result: ProviderHallResultFailed, ReceivedAt: hallNow.Add(-time.Minute), InputTokens: &in, OutputTokens: &out},
		{TargetID: 1, Result: ProviderHallResultPassed, ReceivedAt: hallNow.Add(-6 * time.Minute), InputTokens: &in, OutputTokens: &out, GenerationMs: &gen, TTFTMs: &ttft, TotalMs: &total},
	}
	for i := 0; i < 10; i++ {
		f.probes.recent = append(f.probes.recent, ProviderHallProbeHealth{TargetID: 1, Result: ProviderHallResultPassed, ReceivedAt: hallNow.Add(-time.Duration(10+i) * time.Minute)})
	}
	d, err := f.svc.GetGroup(ctx, 100, 1, "")
	require.NoError(t, err)
	require.Equal(t, int64(1), d.Row.GroupID)
	require.Equal(t, ProviderHallMetricOK, d.E2EAvailability6h.State)
	require.InDelta(t, 11.0/12.0, *d.E2EAvailability6h.Value, 1e-9)
	require.Equal(t, ProviderHallMetricStale, d.TPS.State, "latest passed probe is six minutes old")
	require.Equal(t, 50.0, *d.TPS.Value)
	require.Equal(t, int64(400), *d.ProbeTTFTMs.Value)
	require.Equal(t, int64(1000), *d.ProbeTotalMs.Value)
	require.Equal(t, int64(1500), *d.P90Ms.Value)
	require.Equal(t, ProviderHallProbeTokensResult{Input: 12, Output: 30, At: hallNow.Add(-time.Minute)}, *d.ProbeTokens)
	require.Equal(t, hallNow.Add(-time.Minute), *d.LastProbeAt)
	require.Len(t, d.Trend, 72)
	require.Len(t, d.Profiles, 2)
	require.Equal(t, "passed", *d.Profiles[0].Verification.Verdict)
	require.Nil(t, d.Profiles[1].Verification.Verdict)

	// Fewer than twelve samples → insufficient; no passed probe → no TPS.
	f.probes.recent = f.probes.recent[:1]
	d, err = f.svc.GetGroup(ctx, 100, 1, "24h")
	require.NoError(t, err)
	require.Equal(t, ProviderHallReasonSamplesBelow12, d.E2EAvailability6h.ReasonCode)
	require.Equal(t, ProviderHallReasonNoProbe, d.TPS.ReasonCode)
	require.Equal(t, 900, d.TrendBucket)
	require.Len(t, d.Trend, 96)
	_, err = f.svc.GetGroup(ctx, 100, 1, "2h")
	require.ErrorIs(t, err, ErrProviderHallInvalidQuery)

	f.probes.reports = f.read.reports
	items, total2, err := f.svc.ListVerifications(ctx, 100, 1, nil, 1, 10)
	require.NoError(t, err)
	require.Equal(t, 1, total2)
	require.Equal(t, "gpt-5", items[0].Model)
	require.Equal(t, "responses", items[0].Protocol)
	require.False(t, items[0].Expired)
	pid := int64(11)
	items, _, err = f.svc.ListVerifications(ctx, 100, 1, &pid, 1, 10)
	require.NoError(t, err)
	require.Empty(t, items)
	_, _, err = f.svc.ListVerifications(ctx, 100, 1, nil, 1, 51)
	require.ErrorIs(t, err, ErrProviderHallInvalidQuery)
}

func TestProviderHallQueryNoWindowYet(t *testing.T) {
	f := newHallFixture(t)
	f.read.window = nil
	res, err := f.svc.List(context.Background(), 100, ProviderHallListQuery{})
	require.NoError(t, err)
	require.Equal(t, "", res.SnapshotID)
	require.Nil(t, res.DataThrough)
	require.Len(t, res.Items, 3)
	require.Equal(t, ProviderHallReasonSnapshotMissing, res.Items[0].TTFTFast95Ms.ReasonCode)
	require.Empty(t, res.Items[0].Sparkline.Points)
	require.Zero(t, f.read.calls["snapshots"])
}

func TestProviderHallPredictedRateFixedCases(t *testing.T) {
	one, tenth, half := decimal.NewFromInt(1), decimal.NewFromFloat(0.1), decimal.NewFromFloat(0.5)
	fresh := hallNow.Add(-time.Hour)
	old := hallNow.Add(-8 * 24 * time.Hour)
	ref := ProviderHallReference{InputPrice: &one, CachePrice: &tenth, CacheRate: &half, ConfirmedAt: &fresh}
	quote := ProviderHallQuoteFromPricing(hallTokenPricing(0.000001, 0.0000001), 1, ProviderHallUnitUSDPerMillion)
	cache := 0.4
	okCache := ProviderHallMetric[float64]{Value: &cache, State: ProviderHallMetricOK}
	snap := hallSnapshot(1, 10, hallNow, nil)
	cases := []struct {
		name   string
		in     ProviderHallPredictedRateInput
		state  string
		reason string
		value  string
	}{
		{"ok", ProviderHallPredictedRateInput{Quote: quote, Snapshot: &snap, CacheRate: okCache, Reference: ref, Now: hallNow}, ProviderHallMetricOK, "", "1.1636"},
		{"tiered_quote", ProviderHallPredictedRateInput{Quote: ProviderHallQuoteFromPricing(&ResolvedPricing{Mode: BillingModeToken, BasePricing: &ModelPricing{InputPricePerToken: 1}, Intervals: []PricingInterval{{}}}, 1, ""), Snapshot: &snap, CacheRate: okCache, Reference: ref, Now: hallNow}, ProviderHallMetricNotApplicable, ProviderHallReasonQuoteUnavailable, ""},
		{"cache_creation", ProviderHallPredictedRateInput{Quote: quote, Snapshot: func() *ProviderHallSnapshotRow { s := snap; s.CacheCreationTokens = 1; return &s }(), CacheRate: okCache, Reference: ref, Now: hallNow}, ProviderHallMetricNotApplicable, ProviderHallReasonCacheCreation, ""},
		{"cache_rate_insufficient", ProviderHallPredictedRateInput{Quote: quote, Snapshot: &snap, CacheRate: ProviderHallMetric[float64]{State: ProviderHallMetricInsufficient}, Reference: ref, Now: hallNow}, ProviderHallMetricNotApplicable, ProviderHallReasonCacheRateMissing, ""},
		{"reference_missing", ProviderHallPredictedRateInput{Quote: quote, Snapshot: &snap, CacheRate: okCache, Reference: ProviderHallReference{}, Now: hallNow}, ProviderHallMetricNotApplicable, ProviderHallReasonReferenceMissing, ""},
		{"reference_zero", ProviderHallPredictedRateInput{Quote: quote, Snapshot: &snap, CacheRate: okCache, Reference: ProviderHallReference{InputPrice: &decimal.Zero, CachePrice: &decimal.Zero, CacheRate: &half, ConfirmedAt: &fresh}, Now: hallNow}, ProviderHallMetricNotApplicable, ProviderHallReasonReferenceZero, ""},
		{"reference_expired", ProviderHallPredictedRateInput{Quote: quote, Snapshot: &snap, CacheRate: okCache, Reference: ProviderHallReference{InputPrice: &one, CachePrice: &tenth, CacheRate: &half, ConfirmedAt: &old}, Now: hallNow}, ProviderHallMetricStale, ProviderHallReasonReferenceExpired, "1.1636"},
		{"cache_rate_stale", ProviderHallPredictedRateInput{Quote: quote, Snapshot: &snap, CacheRate: ProviderHallMetric[float64]{Value: &cache, State: ProviderHallMetricStale}, Reference: ref, Now: hallNow}, ProviderHallMetricStale, ProviderHallReasonStale, "1.1636"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := ProviderHallPredictedRateMetric(tc.in)
			require.Equal(t, tc.state, m.State)
			require.Equal(t, tc.reason, m.ReasonCode)
			if tc.value == "" {
				require.Nil(t, m.Value)
			} else {
				require.Equal(t, tc.value, *m.Value)
			}
		})
	}
	// Multiplier folds into the quote: doubling the multiplier doubles the rate.
	doubled := ProviderHallQuoteFromPricing(hallTokenPricing(0.000001, 0.0000001), 2, ProviderHallUnitUSDPerMillion)
	require.Equal(t, "2.0000000000", doubled.InputPrice)
	m := ProviderHallPredictedRateMetric(ProviderHallPredictedRateInput{Quote: doubled, Snapshot: &snap, CacheRate: okCache, Reference: ref, Now: hallNow})
	require.Equal(t, "2.3273", *m.Value)
	// Quote reasons.
	require.Equal(t, ProviderHallReasonPricingUnavailable, ProviderHallQuoteFromPricing(nil, 1, "").ReasonCode)
	require.Equal(t, ProviderHallReasonPricingUnavailable, ProviderHallQuoteFromPricing(&ResolvedPricing{Mode: BillingModePerRequest, BasePricing: &ModelPricing{InputPricePerToken: 1}}, 1, "").ReasonCode)
	require.Equal(t, ProviderHallReasonZeroPrice, ProviderHallQuoteFromPricing(hallTokenPricing(0, 0), 1, "").ReasonCode)
	// Historical unit rules.
	mixed := hallSnapshot(1, 10, hallNow, func(s *ProviderHallSnapshotRow) { h := decimal.NewFromFloat(0.5); s.SubscriptionShare = &h })
	hm, unit := ProviderHallHistoricalPriceMetric(ProviderHallMetricInput{Snapshot: &mixed, Now: hallNow, CollectionEnabled: true, TargetEnabled: true, GroupListed: true})
	require.Equal(t, ProviderHallMetricNotApplicable, hm.State)
	require.Equal(t, ProviderHallReasonMixedBilling, hm.ReasonCode)
	require.Equal(t, "", unit)
	// Snapshot id round trip and rejection of non-minute values.
	id := providerHallSnapshotIDFor(hallNow)
	parsed, ok := ParseProviderHallSnapshotID(id)
	require.True(t, ok)
	require.Equal(t, hallNow, parsed)
	_, ok = ParseProviderHallSnapshotID("1757671230")
	require.False(t, ok)
	// Effective multiplier keeps a short exact decimal.
	_, s := ProviderHallEffectiveMultiplier(1.1, 1.5)
	require.Equal(t, "1.65", s)
}
