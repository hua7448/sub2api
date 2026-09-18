package service

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/shopspring/decimal"
)

// User-facing read model of the provider hall (batch B5). The service never
// exposes account IDs, credentials, upstream addresses, request bodies or
// absolute real-traffic counts; sample sizes surface only as metric states.

var ErrProviderHallInvalidQuery = infraerrors.BadRequest("PROVIDER_HALL_INVALID_QUERY", "invalid provider hall query")

const (
	ProviderHallDefaultPageSize = 50
	ProviderHallMaxPageSize     = 100
	ProviderHallMaxSearchRunes  = 100
	ProviderHallMaxReportPage   = 50
	// ProviderHallBaseCacheTTL bounds how long user-independent base data is
	// reused between requests. Personalised responses are never cached.
	ProviderHallBaseCacheTTL = 30 * time.Second
	// ProviderHallProbeFreshness is the number of probe intervals a probe
	// result stays authoritative for the health dot.
	ProviderHallProbeFreshness = 3
	// ProviderHallE2EWindow / MinSamples define the detail availability metric.
	ProviderHallE2EWindow     = 6 * time.Hour
	ProviderHallE2EMinSamples = 12

	ProviderHallHealthUp      = "up"
	ProviderHallHealthDown    = "down"
	ProviderHallHealthUnknown = "unknown"

	ProviderHallReasonSamplesBelow12 = "samples_below_12"
	ProviderHallReasonNoProbe        = "no_probe"
)

// ProviderHallReadTarget is an enabled target of a listed group joined with
// its profile and reference prices.
type ProviderHallReadTarget struct {
	TargetID             int64
	GroupID              int64
	ProfileID            int64
	Model                string
	Protocol             string
	SupportsTools        bool
	ProbeIntervalSeconds int
	Reference            ProviderHallReference
}

// ProviderHallListedGroup is a listed hall group with its original name.
type ProviderHallListedGroup struct {
	ProviderHallGroup
	GroupName string
}

// ProviderHallSeriesPoint is the slim projection of a tier-5 snapshot used
// by sparklines and trends.
type ProviderHallSeriesPoint struct {
	GroupID     int64
	ProfileID   int64
	WindowEnd   time.Time
	SampleCount int64
	Fast95Ms    *float64
	Coverage    string
}

// ProviderHallGroupProbePoint is one probe trend bucket of one group.
type ProviderHallGroupProbePoint struct {
	GroupID int64
	ProviderHallProbePoint
}

// ProviderHallReadRepository serves the hall with a fixed number of batched
// statements per request. It never reads usage_logs or provider_hall_requests.
type ProviderHallReadRepository interface {
	ListListedGroups(ctx context.Context) ([]ProviderHallListedGroup, error)
	ListEnabledTargets(ctx context.Context) ([]ProviderHallReadTarget, error)
	LatestWindow(ctx context.Context) (windowEnd *time.Time, algorithmVersion int, err error)
	LoadSnapshots(ctx context.Context, windowEnd time.Time, algorithmVersion int) ([]ProviderHallSnapshotRow, error)
	LoadSeries(ctx context.Context, windowEnds []time.Time, algorithmVersion int) ([]ProviderHallSeriesPoint, error)
	LatestVerifications(ctx context.Context) ([]ProviderHallVerification, error)
	ProbeSeriesAll(ctx context.Context, start, end time.Time, bucketSeconds int) ([]ProviderHallGroupProbePoint, error)
}

// ProviderHallProbeReader is the slice of the job repository the hall reads.
type ProviderHallProbeReader interface {
	LatestProbePerTarget(ctx context.Context, since time.Time) (map[int64]ProviderHallProbeHealth, error)
	RecentProbeSamples(ctx context.Context, groupID int64, profileID *int64, since time.Time, limit int) ([]ProviderHallProbeHealth, error)
	ProbeSeries(ctx context.Context, q ProviderHallProbeSeriesQuery) ([]ProviderHallProbePoint, error)
	ListVerifications(ctx context.Context, groupID int64, profileID *int64, page, pageSize int) ([]ProviderHallVerification, int, error)
}

type providerHallRateResolver interface {
	ResolveUserGroupRateMultiplier(ctx context.Context, userID, groupID int64, groupDefaultMultiplier float64) float64
}

type providerHallPricingResolver interface {
	Resolve(ctx context.Context, input PricingInput) *ResolvedPricing
}

// Result types. The handler converts them into dto.* one to one.

type ProviderHallProfileRef struct {
	ProfileID int64
	Model     string
	Protocol  string
}

type ProviderHallModelHealthResult struct {
	ProviderHallProfileRef
	Status    string
	CheckedAt *time.Time
}

type ProviderHallHealthResult struct {
	Status    string
	CheckedAt *time.Time
	Models    []ProviderHallModelHealthResult
}

type ProviderHallSparkPointResult struct {
	T          time.Time
	RealTTFTMs *float64
	ProbeMs    *float64
	Gap        bool
}

type ProviderHallSparklineResult struct {
	Range         string
	BucketSeconds int
	Points        []ProviderHallSparkPointResult
}

type ProviderHallVerificationBadgeResult struct {
	Verdict     *string
	ReasonCode  string
	CompletedAt *time.Time
	Expired     bool
	ReportID    *int64
}

type ProviderHallRowResult struct {
	GroupID         int64
	Name            string
	Description     string
	DisplayOrder    int
	Rate            string
	RateValue       decimal.Decimal
	Quote           ProviderHallQuote
	HistoricalPrice ProviderHallMetric[string]
	HistoricalUnit  string
	PredictedRate   ProviderHallMetric[string]
	Health          ProviderHallHealthResult
	TTFTFast95Ms    ProviderHallMetric[float64]
	TTFTP90Ms       ProviderHallMetric[int64]
	CacheRate       ProviderHallMetric[float64]
	SuccessRate     ProviderHallMetric[float64]
	Sparkline       ProviderHallSparklineResult
	Verification    ProviderHallVerificationBadgeResult
	DefaultProfile  ProviderHallProfileRef

	group   Group
	targets []ProviderHallReadTarget
}

func (r *ProviderHallRowResult) sortKey() providerHallSortKey {
	return providerHallSortKey{
		groupID:         r.GroupID,
		displayOrder:    r.DisplayOrder,
		rate:            r.RateValue,
		historicalPrice: providerHallSortDecimal(r.HistoricalPrice, providerHallParseDecimal),
		predictedRate:   providerHallSortDecimal(r.PredictedRate, providerHallParseDecimal),
		ttftFast95:      providerHallSortFloat(r.TTFTFast95Ms),
		cacheRate:       providerHallSortFloat(r.CacheRate),
		successRate:     providerHallSortFloat(r.SuccessRate),
	}
}

type ProviderHallListQuery struct {
	Range      string
	Model      string
	Protocol   string
	Search     string
	Sort       []ProviderHallSortRule
	Page       int
	PageSize   int
	SnapshotID string
}

type ProviderHallCatalogResult struct {
	Models  []ProviderHallProfileRef
	Default ProviderHallProfileRef
}

type ProviderHallSummaryResult struct {
	Available int
	Listed    int
	Abnormal  int
	Verified  int
}

type ProviderHallListResult struct {
	SnapshotID    string
	DataThrough   *time.Time
	MetricVersion int
	PricingAt     time.Time
	Range         string
	Catalog       ProviderHallCatalogResult
	Summary       ProviderHallSummaryResult
	Items         []*ProviderHallRowResult
	Page          int
	PageSize      int
	Total         int
}

type ProviderHallProbeTokensResult struct {
	Input  int
	Output int
	At     time.Time
}

type ProviderHallTrendPointResult struct {
	T            time.Time
	RealFast95Ms *float64
	ProbeTotalMs *float64
	ProbeTTFTMs  *float64
	Gap          bool
	ProbeFailed  int
}

type ProviderHallProfileStatusResult struct {
	ProviderHallProfileRef
	Health       ProviderHallModelHealthResult
	Verification ProviderHallVerificationBadgeResult
}

type ProviderHallDetailResult struct {
	Row               *ProviderHallRowResult
	E2EAvailability6h ProviderHallMetric[float64]
	TPS               ProviderHallMetric[float64]
	ProbeTTFTMs       ProviderHallMetric[int64]
	ProbeTotalMs      ProviderHallMetric[int64]
	P90Ms             ProviderHallMetric[int64]
	ProbeTokens       *ProviderHallProbeTokensResult
	LastProbeAt       *time.Time
	TrendRange        string
	TrendBucket       int
	Trend             []ProviderHallTrendPointResult
	Profiles          []ProviderHallProfileStatusResult
}

type ProviderHallReportResult struct {
	ProviderHallVerification
	Model    string
	Protocol string
	Expired  bool
}

// providerHallBase is the user-independent slice of one request, cached for
// ProviderHallBaseCacheTTL keyed by (window end, range).
type providerHallBase struct {
	loadedAt  time.Time
	T         time.Time
	algo      int
	rangeKey  string
	bucket    int
	points    int
	groups    []ProviderHallListedGroup
	targets   map[int64][]ProviderHallReadTarget // by group
	snapshots map[[2]int64]*ProviderHallSnapshotRow
	series    map[int64]map[time.Time]ProviderHallSeriesPoint // (group, profile 0)
	probes    map[int64]ProviderHallProbeHealth               // by target
	reports   map[[2]int64]*ProviderHallVerification
	probeSpk  map[int64]map[time.Time]ProviderHallProbePoint // by group, bucket start
	ends      []time.Time
}

type providerHallSnapshotKey = [2]int64

type ProviderHallQueryService struct {
	read    ProviderHallReadRepository
	cfgRepo ProviderHallRepository
	hall    *ProviderHallService
	access  ProviderHallGroupAccess
	rates   providerHallRateResolver
	pricing providerHallPricingResolver
	probes  ProviderHallProbeReader
	now     func() time.Time

	// loadMu serialises base loads so concurrent misses share one fetch;
	// cache entries are immutable after publication.
	loadMu sync.Mutex
	cache  map[string]*providerHallBase
}

func NewProviderHallQueryService(read ProviderHallReadRepository, cfgRepo ProviderHallRepository, hall *ProviderHallService, access ProviderHallGroupAccess, rates providerHallRateResolver, pricing providerHallPricingResolver, probes ProviderHallProbeReader) *ProviderHallQueryService {
	return &ProviderHallQueryService{
		read: read, cfgRepo: cfgRepo, hall: hall, access: access, rates: rates, pricing: pricing, probes: probes,
		now:   time.Now,
		cache: map[string]*providerHallBase{},
	}
}

// SetClockForTest pins the clock; tests in other packages use it.
func (s *ProviderHallQueryService) SetClockForTest(now func() time.Time) {
	if s != nil && now != nil {
		s.now = now
	}
}

// DisplayEnabled re-reads the switch; the route guard and every method use
// it so a flipped switch takes effect on the next request.
func (s *ProviderHallQueryService) DisplayEnabled(ctx context.Context) bool {
	if s == nil || s.cfgRepo == nil {
		return false
	}
	cfg, err := s.cfgRepo.GetConfig(ctx)
	return err == nil && cfg != nil && cfg.DisplayEnabled
}

func providerHallSnapshotIDFor(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return strconv.FormatInt(t.UTC().Unix(), 10)
}

// ParseProviderHallSnapshotID accepts the value published as snapshot_id.
func ParseProviderHallSnapshotID(raw string) (time.Time, bool) {
	sec, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || sec <= 0 || sec%60 != 0 {
		return time.Time{}, false
	}
	return time.Unix(sec, 0).UTC(), true
}

func providerHallRangeOK(r string) bool {
	_, _, ok := ProviderHallTrendBucket(r)
	return ok
}

// NormalizeListQuery applies defaults and validates. The default range comes
// from the hall configuration.
func (s *ProviderHallQueryService) NormalizeListQuery(ctx context.Context, q ProviderHallListQuery) (ProviderHallListQuery, error) {
	if q.Range == "" {
		cfg, err := s.cfgRepo.GetConfig(ctx)
		if err != nil {
			return q, err
		}
		q.Range = cfg.DefaultRange
	}
	if !providerHallRangeOK(q.Range) {
		return q, ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": "range"})
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize == 0 {
		q.PageSize = ProviderHallDefaultPageSize
	}
	if q.PageSize < 1 || q.PageSize > ProviderHallMaxPageSize {
		return q, ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": "page_size"})
	}
	q.Search = strings.TrimSpace(q.Search)
	if utf8.RuneCountInString(q.Search) > ProviderHallMaxSearchRunes {
		return q, ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": "search"})
	}
	if q.Protocol != "" && !providerHallProtocol(q.Protocol) {
		return q, ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": "protocol"})
	}
	if utf8.RuneCountInString(q.Model) > 200 {
		return q, ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": "model"})
	}
	if len(q.Sort) == 0 {
		q.Sort = []ProviderHallSortRule{{Field: "display_order"}}
	}
	if q.SnapshotID != "" {
		if _, ok := ParseProviderHallSnapshotID(q.SnapshotID); !ok {
			return q, ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": "snapshot_id"})
		}
	}
	return q, nil
}

// permittedGroups is the strict first step: groups the user may bind, active,
// on an OpenAI-capable platform. Expired subscriptions and revoked exclusive
// grants are already excluded by GetAvailableGroups.
func (s *ProviderHallQueryService) permittedGroups(ctx context.Context, userID int64) (map[int64]Group, error) {
	if userID < 1 {
		return map[int64]Group{}, nil
	}
	groups, err := s.access.GetAvailableGroups(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make(map[int64]Group, len(groups))
	for _, g := range groups {
		if !g.IsActive() || (g.Platform != PlatformOpenAI && g.Platform != PlatformComposite) {
			continue
		}
		out[g.ID] = g
	}
	return out, nil
}

// resolveWindow picks T: the requested snapshot or the aggregator watermark.
func (s *ProviderHallQueryService) resolveWindow(ctx context.Context, snapshotID string) (time.Time, int, error) {
	latest, algo, err := s.read.LatestWindow(ctx)
	if err != nil {
		return time.Time{}, 0, err
	}
	if algo == 0 {
		algo = ProviderHallAlgorithmVersion
	}
	if snapshotID != "" {
		t, ok := ParseProviderHallSnapshotID(snapshotID)
		if !ok {
			return time.Time{}, 0, ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": "snapshot_id"})
		}
		return t, algo, nil
	}
	if latest == nil {
		return time.Time{}, algo, nil
	}
	return latest.UTC(), algo, nil
}

func providerHallBucketEnds(T time.Time, bucketSeconds, points int) []time.Time {
	if T.IsZero() || bucketSeconds <= 0 || points <= 0 {
		return nil
	}
	bucket := time.Duration(bucketSeconds) * time.Second
	end := T.UTC().Truncate(bucket)
	ends := make([]time.Time, points)
	for i := range ends {
		ends[i] = end.Add(-time.Duration(points-1-i) * bucket)
	}
	return ends
}

// loadBase fetches every user-independent input with a fixed set of
// statements and caches it briefly.
func (s *ProviderHallQueryService) loadBase(ctx context.Context, T time.Time, algo int, rangeKey string) (*providerHallBase, error) {
	key := providerHallSnapshotIDFor(T) + "|" + rangeKey + "|" + strconv.Itoa(algo)
	now := s.now()
	s.loadMu.Lock()
	defer s.loadMu.Unlock()
	if b, ok := s.cache[key]; ok && now.Sub(b.loadedAt) < ProviderHallBaseCacheTTL {
		return b, nil
	}

	bucket, points, _ := ProviderHallTrendBucket(rangeKey)
	b := &providerHallBase{
		loadedAt: now, T: T, algo: algo, rangeKey: rangeKey, bucket: bucket, points: points,
		targets:   map[int64][]ProviderHallReadTarget{},
		snapshots: map[providerHallSnapshotKey]*ProviderHallSnapshotRow{},
		series:    map[int64]map[time.Time]ProviderHallSeriesPoint{},
		reports:   map[providerHallSnapshotKey]*ProviderHallVerification{},
		probeSpk:  map[int64]map[time.Time]ProviderHallProbePoint{},
	}
	var err error
	if b.groups, err = s.read.ListListedGroups(ctx); err != nil {
		return nil, err
	}
	targets, err := s.read.ListEnabledTargets(ctx)
	if err != nil {
		return nil, err
	}
	maxInterval := 0
	for _, t := range targets {
		b.targets[t.GroupID] = append(b.targets[t.GroupID], t)
		if t.ProbeIntervalSeconds > maxInterval {
			maxInterval = t.ProbeIntervalSeconds
		}
	}
	if !T.IsZero() {
		rows, err := s.read.LoadSnapshots(ctx, T, algo)
		if err != nil {
			return nil, err
		}
		for i := range rows {
			row := rows[i]
			b.snapshots[providerHallSnapshotKey{row.GroupID, row.ProfileID}] = &row
		}
		b.ends = providerHallBucketEnds(T, bucket, points)
		series, err := s.read.LoadSeries(ctx, b.ends, algo)
		if err != nil {
			return nil, err
		}
		for _, p := range series {
			if p.ProfileID != 0 {
				continue
			}
			if b.series[p.GroupID] == nil {
				b.series[p.GroupID] = map[time.Time]ProviderHallSeriesPoint{}
			}
			b.series[p.GroupID][p.WindowEnd.UTC()] = p
		}
		if len(b.ends) > 0 {
			start := b.ends[0].Add(-time.Duration(bucket) * time.Second)
			spk, err := s.read.ProbeSeriesAll(ctx, start, b.ends[len(b.ends)-1], bucket)
			if err != nil {
				return nil, err
			}
			for _, p := range spk {
				if b.probeSpk[p.GroupID] == nil {
					b.probeSpk[p.GroupID] = map[time.Time]ProviderHallProbePoint{}
				}
				b.probeSpk[p.GroupID][p.BucketStart.UTC()] = p.ProviderHallProbePoint
			}
		}
	}
	if s.probes != nil {
		since := now.Add(-time.Duration(ProviderHallProbeFreshness*max(maxInterval, 300)) * time.Second)
		if b.probes, err = s.probes.LatestProbePerTarget(ctx, since); err != nil {
			return nil, err
		}
	}
	reports, err := s.read.LatestVerifications(ctx)
	if err != nil {
		return nil, err
	}
	for i := range reports {
		v := reports[i]
		b.reports[providerHallSnapshotKey{v.GroupID, v.ProfileID}] = &v
	}
	s.cache[key] = b
	// Drop stale entries so the map cannot grow with every distinct T.
	for k, v := range s.cache {
		if now.Sub(v.loadedAt) >= ProviderHallBaseCacheTTL {
			delete(s.cache, k)
		}
	}
	return b, nil
}

func providerHallDefaultTarget(cfg *ProviderHallConfig, targets []ProviderHallReadTarget) *ProviderHallReadTarget {
	var best *ProviderHallReadTarget
	for i := range targets {
		t := &targets[i]
		if cfg != nil && t.Model == cfg.DefaultModel && t.Protocol == cfg.DefaultProtocol {
			return t
		}
		if best == nil || t.ProfileID < best.ProfileID {
			best = t
		}
	}
	return best
}

func providerHallRef(t *ProviderHallReadTarget) ProviderHallProfileRef {
	if t == nil {
		return ProviderHallProfileRef{}
	}
	return ProviderHallProfileRef{ProfileID: t.ProfileID, Model: t.Model, Protocol: t.Protocol}
}

func (s *ProviderHallQueryService) targetHealth(b *providerHallBase, t ProviderHallReadTarget, now time.Time) ProviderHallModelHealthResult {
	h := ProviderHallModelHealthResult{ProviderHallProfileRef: providerHallRef(&t), Status: ProviderHallHealthUnknown}
	p, ok := b.probes[t.TargetID]
	if !ok {
		return h
	}
	interval := t.ProbeIntervalSeconds
	if interval <= 0 {
		interval = 300
	}
	if now.Sub(p.ReceivedAt) > time.Duration(ProviderHallProbeFreshness*interval)*time.Second {
		return h
	}
	at := p.ReceivedAt
	h.CheckedAt = &at
	if p.Result == ProviderHallResultPassed {
		h.Status = ProviderHallHealthUp
	} else {
		h.Status = ProviderHallHealthDown
	}
	return h
}

func providerHallGroupHealth(models []ProviderHallModelHealthResult) ProviderHallHealthResult {
	out := ProviderHallHealthResult{Status: ProviderHallHealthUnknown, Models: models}
	if out.Models == nil {
		out.Models = []ProviderHallModelHealthResult{}
	}
	down, unknown := false, false
	for _, m := range models {
		switch m.Status {
		case ProviderHallHealthDown:
			down = true
		case ProviderHallHealthUnknown:
			unknown = true
		}
		if m.CheckedAt != nil && (out.CheckedAt == nil || m.CheckedAt.After(*out.CheckedAt)) {
			at := *m.CheckedAt
			out.CheckedAt = &at
		}
	}
	switch {
	case len(models) == 0:
		out.Status = ProviderHallHealthUnknown
	case down:
		out.Status = ProviderHallHealthDown
	case unknown:
		out.Status = ProviderHallHealthUnknown
	default:
		out.Status = ProviderHallHealthUp
	}
	return out
}

func providerHallBadge(v *ProviderHallVerification, now time.Time) ProviderHallVerificationBadgeResult {
	if v == nil {
		return ProviderHallVerificationBadgeResult{}
	}
	verdict, id, at := v.Verdict, v.JobID, v.CompletedAt
	return ProviderHallVerificationBadgeResult{Verdict: &verdict, ReasonCode: v.ReasonCode, CompletedAt: &at, Expired: v.Expired(now), ReportID: &id}
}

func (s *ProviderHallQueryService) sparkline(b *providerHallBase, groupID int64) ProviderHallSparklineResult {
	out := ProviderHallSparklineResult{Range: b.rangeKey, BucketSeconds: b.bucket, Points: []ProviderHallSparkPointResult{}}
	real := b.series[groupID]
	probe := b.probeSpk[groupID]
	for _, end := range b.ends {
		pt := ProviderHallSparkPointResult{T: end, Gap: true}
		if p, ok := real[end]; ok && p.Coverage == ProviderHallCoverageComplete {
			pt.Gap = false
			if p.SampleCount >= ProviderHallMinQualitySamples && p.Fast95Ms != nil {
				v := *p.Fast95Ms
				pt.RealTTFTMs = &v
			}
		}
		if pp, ok := probe[end.Add(-time.Duration(b.bucket)*time.Second)]; ok && pp.AvgTotalMs != nil {
			v := *pp.AvgTotalMs
			pt.ProbeMs = &v
		}
		out.Points = append(out.Points, pt)
	}
	return out
}

// buildRow assembles one group's row from the base and the user's pricing.
func (s *ProviderHallQueryService) buildRow(ctx context.Context, cfg *ProviderHallConfig, b *providerHallBase, listing ProviderHallListedGroup, group Group, userID int64, now time.Time) *ProviderHallRowResult {
	targets := b.targets[listing.GroupID]
	def := providerHallDefaultTarget(cfg, targets)
	row := &ProviderHallRowResult{
		GroupID:        listing.GroupID,
		Name:           listing.DisplayName,
		Description:    listing.Description,
		DisplayOrder:   listing.DisplayOrder,
		DefaultProfile: providerHallRef(def),
		group:          group,
		targets:        targets,
	}
	if row.Name == "" {
		row.Name = listing.GroupName
	}
	// Merged metrics (profile 0) for TTFT / P90 / success rate.
	merged := ProviderHallMetricInput{Snapshot: b.snapshots[providerHallSnapshotKey{listing.GroupID, 0}], Now: now, CollectionEnabled: cfg.CollectionEnabled, TargetEnabled: len(targets) > 0, GroupListed: listing.Listed}
	row.TTFTFast95Ms, row.TTFTP90Ms = ProviderHallTTFTMetric(merged)
	row.SuccessRate = ProviderHallSuccessMetric(merged)
	// Default-profile metrics for cache rate / historical price / prediction.
	var defSnap *ProviderHallSnapshotRow
	if def != nil {
		defSnap = b.snapshots[providerHallSnapshotKey{listing.GroupID, def.ProfileID}]
	}
	defIn := ProviderHallMetricInput{Snapshot: defSnap, Now: now, CollectionEnabled: cfg.CollectionEnabled, TargetEnabled: def != nil, GroupListed: listing.Listed}
	row.CacheRate = ProviderHallCacheMetric(defIn)
	row.HistoricalPrice, row.HistoricalUnit = ProviderHallHistoricalPriceMetric(defIn)

	// Pricing overlay: multiplier at the pricing instant, then the quote.
	userMult := group.RateMultiplier
	if s.rates != nil {
		userMult = s.rates.ResolveUserGroupRateMultiplier(ctx, userID, group.ID, group.RateMultiplier)
	}
	mult, rateStr := ProviderHallEffectiveMultiplier(userMult, group.PeakMultiplierAt(now))
	row.Rate = rateStr
	row.RateValue, _ = decimal.NewFromString(rateStr)
	unit := ProviderHallUnitUSDPerMillion
	if group.IsSubscriptionType() {
		unit = ProviderHallUnitQuotaPerMillion
	}
	if def != nil && s.pricing != nil {
		gid := group.ID
		g := group
		row.Quote = ProviderHallQuoteFromPricing(s.pricing.Resolve(ctx, PricingInput{Model: def.Model, GroupID: &gid, Group: &g}), mult, unit)
	} else {
		row.Quote = ProviderHallQuote{Unit: unit, ReasonCode: ProviderHallReasonPricingUnavailable}
	}
	var ref ProviderHallReference
	if def != nil {
		ref = def.Reference
	}
	row.PredictedRate = ProviderHallPredictedRateMetric(ProviderHallPredictedRateInput{Quote: row.Quote, Snapshot: defSnap, CacheRate: row.CacheRate, Reference: ref, Now: now})

	models := make([]ProviderHallModelHealthResult, 0, len(targets))
	for _, t := range targets {
		models = append(models, s.targetHealth(b, t, now))
	}
	row.Health = providerHallGroupHealth(models)
	if def != nil {
		row.Verification = providerHallBadge(b.reports[providerHallSnapshotKey{listing.GroupID, def.ProfileID}], now)
	}
	row.Sparkline = s.sparkline(b, listing.GroupID)
	return row
}

func providerHallMatchesSearch(listing ProviderHallListedGroup, search string) bool {
	if search == "" {
		return true
	}
	needle := strings.ToLower(search)
	for _, hay := range []string{listing.DisplayName, listing.Description, listing.GroupName} {
		if strings.Contains(strings.ToLower(hay), needle) {
			return true
		}
	}
	return false
}

func providerHallMatchesFilter(targets []ProviderHallReadTarget, model, protocol string) bool {
	if model == "" && protocol == "" {
		return true
	}
	for _, t := range targets {
		if (model == "" || t.Model == model) && (protocol == "" || t.Protocol == protocol) {
			return true
		}
	}
	return false
}

// List runs the fixed pipeline: display switch → permissions → listing →
// search/filter → window → batch reads → pricing → sort → page → summary.
func (s *ProviderHallQueryService) List(ctx context.Context, userID int64, q ProviderHallListQuery) (*ProviderHallListResult, error) {
	cfg, err := s.cfgRepo.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !cfg.DisplayEnabled {
		return nil, ErrProviderHallNotFound
	}
	if q, err = s.NormalizeListQuery(ctx, q); err != nil {
		return nil, err
	}
	permitted, err := s.permittedGroups(ctx, userID)
	if err != nil {
		return nil, err
	}
	T, algo, err := s.resolveWindow(ctx, q.SnapshotID)
	if err != nil {
		return nil, err
	}
	base, err := s.loadBase(ctx, T, algo, q.Range)
	if err != nil {
		return nil, err
	}
	now := s.now()
	result := &ProviderHallListResult{
		SnapshotID: providerHallSnapshotIDFor(T), MetricVersion: algo, PricingAt: now, Range: q.Range,
		Items: []*ProviderHallRowResult{}, Page: q.Page, PageSize: q.PageSize,
	}
	if !T.IsZero() {
		t := T
		result.DataThrough = &t
	}
	// Catalog covers every profile the user could see, before search/filter.
	catalog := map[int64]ProviderHallProfileRef{}
	rows := make([]*ProviderHallRowResult, 0, len(base.groups))
	for _, listing := range base.groups {
		group, ok := permitted[listing.GroupID]
		if !ok || !listing.Listed {
			continue
		}
		for _, t := range base.targets[listing.GroupID] {
			catalog[t.ProfileID] = providerHallRef(&t)
		}
		if !providerHallMatchesSearch(listing, q.Search) || !providerHallMatchesFilter(base.targets[listing.GroupID], q.Model, q.Protocol) {
			continue
		}
		rows = append(rows, s.buildRow(ctx, cfg, base, listing, group, userID, now))
	}
	result.Catalog.Models = make([]ProviderHallProfileRef, 0, len(catalog))
	for _, ref := range catalog {
		result.Catalog.Models = append(result.Catalog.Models, ref)
	}
	sort.Slice(result.Catalog.Models, func(i, j int) bool {
		a, b := result.Catalog.Models[i], result.Catalog.Models[j]
		if a.Model != b.Model {
			return a.Model < b.Model
		}
		if a.Protocol != b.Protocol {
			return a.Protocol < b.Protocol
		}
		return a.ProfileID < b.ProfileID
	})
	result.Catalog.Default = ProviderHallProfileRef{Model: cfg.DefaultModel, Protocol: cfg.DefaultProtocol}
	for _, ref := range result.Catalog.Models {
		if ref.Model == cfg.DefaultModel && ref.Protocol == cfg.DefaultProtocol {
			result.Catalog.Default = ref
			break
		}
	}
	providerHallSortRows(rows, q.Sort)
	for _, row := range rows {
		result.Summary.Listed++
		switch row.Health.Status {
		case ProviderHallHealthUp:
			result.Summary.Available++
		case ProviderHallHealthDown:
			result.Summary.Abnormal++
		}
		if row.Verification.Verdict != nil && *row.Verification.Verdict == "passed" && !row.Verification.Expired {
			result.Summary.Verified++
		}
	}
	result.Total = len(rows)
	start := (q.Page - 1) * q.PageSize
	if start < len(rows) {
		end := min(start+q.PageSize, len(rows))
		result.Items = rows[start:end]
	}
	return result, nil
}

// authorizedRow resolves one group the user may see and builds its row.
func (s *ProviderHallQueryService) authorizedRow(ctx context.Context, userID, groupID int64, rangeKey string) (*ProviderHallConfig, *providerHallBase, *ProviderHallRowResult, error) {
	listing, err := s.hall.AuthorizeUserGroup(ctx, userID, groupID)
	if err != nil {
		return nil, nil, nil, err
	}
	cfg, err := s.cfgRepo.GetConfig(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	permitted, err := s.permittedGroups(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	group, ok := permitted[groupID]
	if !ok {
		return nil, nil, nil, ErrProviderHallNotFound
	}
	T, algo, err := s.resolveWindow(ctx, "")
	if err != nil {
		return nil, nil, nil, err
	}
	base, err := s.loadBase(ctx, T, algo, rangeKey)
	if err != nil {
		return nil, nil, nil, err
	}
	var found *ProviderHallListedGroup
	for i := range base.groups {
		if base.groups[i].GroupID == groupID && base.groups[i].Listed {
			found = &base.groups[i]
			break
		}
	}
	if found == nil {
		// Listed a moment ago per AuthorizeUserGroup but not in the cached base:
		// build from the listing we just authorised.
		found = &ProviderHallListedGroup{ProviderHallGroup: *listing, GroupName: group.Name}
	}
	return cfg, base, s.buildRow(ctx, cfg, base, *found, group, userID, s.now()), nil
}

func providerHallProbeMetric[T any](value *T, at *time.Time, now time.Time, reason string) ProviderHallMetric[T] {
	m := ProviderHallMetric[T]{State: ProviderHallMetricInsufficient, ReasonCode: reason}
	if value == nil {
		return m
	}
	m.Value, m.State, m.ReasonCode = value, ProviderHallMetricOK, ""
	if at != nil {
		t := *at
		m.ComputedAt, m.WindowEnd = &t, &t
		if now.Sub(t) > ProviderHallSnapshotStaleAfter {
			m.State, m.ReasonCode = ProviderHallMetricStale, ProviderHallReasonStale
		}
	}
	return m
}

// GetGroup returns the row plus probe-derived detail metrics, the trend and
// per-profile status. Unknown or unauthorised groups are 404.
func (s *ProviderHallQueryService) GetGroup(ctx context.Context, userID, groupID int64, rangeKey string) (*ProviderHallDetailResult, error) {
	if rangeKey == "" {
		cfg, err := s.cfgRepo.GetConfig(ctx)
		if err != nil {
			return nil, err
		}
		rangeKey = cfg.DefaultRange
	}
	if !providerHallRangeOK(rangeKey) {
		return nil, ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": "range"})
	}
	cfg, base, row, err := s.authorizedRow(ctx, userID, groupID, rangeKey)
	if err != nil {
		return nil, err
	}
	now := s.now()
	out := &ProviderHallDetailResult{Row: row, TrendRange: rangeKey, TrendBucket: base.bucket, Trend: []ProviderHallTrendPointResult{}, Profiles: []ProviderHallProfileStatusResult{}}

	// Probe-derived detail. Absolute counts stay internal.
	e2e := ProviderHallMetric[float64]{State: ProviderHallMetricInsufficient, ReasonCode: ProviderHallReasonSamplesBelow12}
	tps := ProviderHallMetric[float64]{State: ProviderHallMetricInsufficient, ReasonCode: ProviderHallReasonNoProbe}
	ttft := ProviderHallMetric[int64]{State: ProviderHallMetricInsufficient, ReasonCode: ProviderHallReasonNoProbe}
	total := ProviderHallMetric[int64]{State: ProviderHallMetricInsufficient, ReasonCode: ProviderHallReasonNoProbe}
	if s.probes != nil {
		samples, err := s.probes.RecentProbeSamples(ctx, groupID, nil, now.Add(-ProviderHallE2EWindow), 1000)
		if err != nil {
			return nil, err
		}
		passed := 0
		for _, sm := range samples {
			if sm.Result == ProviderHallResultPassed {
				passed++
			}
		}
		if len(samples) >= ProviderHallE2EMinSamples {
			v := float64(passed) / float64(len(samples))
			e2e = providerHallProbeMetric(&v, &samples[0].ReceivedAt, now, "")
			e2e.State, e2e.ReasonCode = ProviderHallMetricOK, ""
			start := now.Add(-ProviderHallE2EWindow)
			e2e.WindowStart, e2e.WindowEnd = &start, &now
		}
		if len(samples) > 0 {
			latest := samples[0]
			at := latest.ReceivedAt
			out.LastProbeAt = &at
			if latest.InputTokens != nil && latest.OutputTokens != nil {
				out.ProbeTokens = &ProviderHallProbeTokensResult{Input: *latest.InputTokens, Output: *latest.OutputTokens, At: at}
			}
		}
		for _, sm := range samples {
			if sm.Result != ProviderHallResultPassed {
				continue
			}
			at := sm.ReceivedAt
			if v := ProviderHallTPS(sm.OutputTokens, sm.GenerationMs); v != nil {
				tps = providerHallProbeMetric(v, &at, now, "")
			}
			if sm.TTFTMs != nil {
				v := int64(*sm.TTFTMs)
				ttft = providerHallProbeMetric(&v, &at, now, "")
			}
			if sm.TotalMs != nil {
				v := int64(*sm.TotalMs)
				total = providerHallProbeMetric(&v, &at, now, "")
			}
			break
		}
	}
	out.E2EAvailability6h, out.TPS, out.ProbeTTFTMs, out.ProbeTotalMs = e2e, tps, ttft, total
	var defSnap *ProviderHallSnapshotRow
	if row.DefaultProfile.ProfileID > 0 {
		defSnap = base.snapshots[providerHallSnapshotKey{groupID, row.DefaultProfile.ProfileID}]
	}
	_, out.P90Ms = ProviderHallTTFTMetric(ProviderHallMetricInput{Snapshot: defSnap, Now: now, CollectionEnabled: cfg.CollectionEnabled, TargetEnabled: row.DefaultProfile.ProfileID > 0, GroupListed: true})

	// Trend: real fast95 from the merged tier-5 snapshots, probe from samples.
	if len(base.ends) > 0 && s.probes != nil {
		start := base.ends[0].Add(-time.Duration(base.bucket) * time.Second)
		points, err := s.probes.ProbeSeries(ctx, ProviderHallProbeSeriesQuery{GroupID: groupID, Start: start, End: base.ends[len(base.ends)-1], BucketSeconds: base.bucket})
		if err != nil {
			return nil, err
		}
		byStart := make(map[time.Time]ProviderHallProbePoint, len(points))
		for _, p := range points {
			byStart[p.BucketStart.UTC()] = p
		}
		real := base.series[groupID]
		for _, end := range base.ends {
			pt := ProviderHallTrendPointResult{T: end, Gap: true}
			if p, ok := real[end]; ok && p.Coverage == ProviderHallCoverageComplete {
				pt.Gap = false
				if p.SampleCount >= ProviderHallMinQualitySamples && p.Fast95Ms != nil {
					v := *p.Fast95Ms
					pt.RealFast95Ms = &v
				}
			}
			if pp, ok := byStart[end.Add(-time.Duration(base.bucket)*time.Second)]; ok {
				pt.ProbeTotalMs, pt.ProbeTTFTMs, pt.ProbeFailed = pp.AvgTotalMs, pp.AvgTTFTMs, pp.Failed
			}
			out.Trend = append(out.Trend, pt)
		}
	}
	for _, t := range row.targets {
		out.Profiles = append(out.Profiles, ProviderHallProfileStatusResult{
			ProviderHallProfileRef: providerHallRef(&t),
			Health:                 s.targetHealth(base, t, now),
			Verification:           providerHallBadge(base.reports[providerHallSnapshotKey{groupID, t.ProfileID}], now),
		})
	}
	return out, nil
}

// ListVerifications pages the group's reports (default profile unless one is
// requested). Authorisation is identical to the detail endpoint.
func (s *ProviderHallQueryService) ListVerifications(ctx context.Context, userID, groupID int64, profileID *int64, page, pageSize int) ([]ProviderHallReportResult, int, error) {
	if _, err := s.hall.AuthorizeUserGroup(ctx, userID, groupID); err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	if pageSize < 1 || pageSize > ProviderHallMaxReportPage {
		return nil, 0, ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": "page_size"})
	}
	if profileID != nil && *profileID < 1 {
		return nil, 0, ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": "profile_id"})
	}
	cfg, err := s.cfgRepo.GetConfig(ctx)
	if err != nil {
		return nil, 0, err
	}
	targets, err := s.read.ListEnabledTargets(ctx)
	if err != nil {
		return nil, 0, err
	}
	var own []ProviderHallReadTarget
	for _, t := range targets {
		if t.GroupID == groupID {
			own = append(own, t)
		}
	}
	if profileID == nil {
		def := providerHallDefaultTarget(cfg, own)
		if def == nil {
			return []ProviderHallReportResult{}, 0, nil
		}
		id := def.ProfileID
		profileID = &id
	}
	profiles, err := s.cfgRepo.ListProfiles(ctx)
	if err != nil {
		return nil, 0, err
	}
	byID := make(map[int64]ProviderHallProfile, len(profiles))
	for _, p := range profiles {
		byID[p.ID] = p
	}
	if s.probes == nil {
		return []ProviderHallReportResult{}, 0, nil
	}
	items, total, err := s.probes.ListVerifications(ctx, groupID, profileID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	now := s.now()
	out := make([]ProviderHallReportResult, 0, len(items))
	for i := range items {
		v := items[i]
		p := byID[v.ProfileID]
		out = append(out, ProviderHallReportResult{ProviderHallVerification: v, Model: p.Model, Protocol: p.Protocol, Expired: v.Expired(now)})
	}
	return out, total, nil
}
