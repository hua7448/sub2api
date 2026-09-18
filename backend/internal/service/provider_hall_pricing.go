package service

import (
	"time"

	"github.com/shopspring/decimal"
)

// Quote and predicted-rate rules of the user hall (batch B5). Every amount is
// a decimal with ten places; float64 inputs from the pricing catalogue are
// converted exactly once, at the boundary.

const (
	ProviderHallUnitUSDPerMillion   = "usd_per_million"
	ProviderHallUnitQuotaPerMillion = "quota_per_million"

	ProviderHallReasonPricingUnavailable = "pricing_unavailable"
	ProviderHallReasonTieredPricing      = "tiered_pricing"
	ProviderHallReasonZeroPrice          = "zero_price"
	ProviderHallReasonQuoteUnavailable   = "quote_unavailable"
	ProviderHallReasonCacheCreation      = "cache_creation_present"
	ProviderHallReasonCacheRateMissing   = "cache_rate_unavailable"
	ProviderHallReasonReferenceMissing   = "reference_missing"
	ProviderHallReasonReferenceZero      = "reference_zero"
	ProviderHallReasonReferenceExpired   = "reference_expired"
	ProviderHallReasonMixedBilling       = "mixed_billing"

	// ProviderHallReferenceMaxAge is how long an admin-confirmed reference
	// price stays current before the predicted rate turns stale.
	ProviderHallReferenceMaxAge = 7 * 24 * time.Hour
	providerHallPriceScale      = 10
)

var providerHallMillion = decimal.NewFromInt(1_000_000)

// ProviderHallQuote is the user's current quoted price on one profile.
type ProviderHallQuote struct {
	InputPrice string
	CachePrice string
	Unit       string
	Applicable bool
	ReasonCode string
	// input and cache keep the exact decimals for the predicted rate.
	input, cache decimal.Decimal
}

// ProviderHallQuoteFromPricing builds the quote from resolved pricing and the
// user's effective multiplier. Only flat token pricing with a positive input
// price is quotable; tiered ladders, per-request and media modes are not.
func ProviderHallQuoteFromPricing(rp *ResolvedPricing, multiplier float64, unit string) ProviderHallQuote {
	q := ProviderHallQuote{Unit: unit, InputPrice: "", CachePrice: ""}
	switch {
	case rp == nil || rp.BasePricing == nil:
		q.ReasonCode = ProviderHallReasonPricingUnavailable
		return q
	case rp.Mode != "" && rp.Mode != BillingModeToken:
		q.ReasonCode = ProviderHallReasonPricingUnavailable
		return q
	case len(rp.Intervals) > 0:
		q.ReasonCode = ProviderHallReasonTieredPricing
		return q
	case !(rp.BasePricing.InputPricePerToken > 0):
		q.ReasonCode = ProviderHallReasonZeroPrice
		return q
	}
	m := decimal.NewFromFloat(multiplier)
	if m.IsNegative() {
		m = decimal.Zero
	}
	q.input = decimal.NewFromFloat(rp.BasePricing.InputPricePerToken).Mul(providerHallMillion).Mul(m)
	q.cache = decimal.NewFromFloat(rp.BasePricing.CacheReadPricePerToken).Mul(providerHallMillion).Mul(m)
	if q.cache.IsNegative() {
		q.cache = decimal.Zero
	}
	q.InputPrice = q.input.StringFixed(providerHallPriceScale)
	q.CachePrice = q.cache.StringFixed(providerHallPriceScale)
	q.Applicable = true
	return q
}

// ProviderHallReference is the admin-maintained benchmark of one profile.
type ProviderHallReference struct {
	InputPrice  *decimal.Decimal
	CachePrice  *decimal.Decimal
	CacheRate   *decimal.Decimal
	ConfirmedAt *time.Time
}

// ProviderHallPredictedRate compares the user's effective blended input price
// against the reference blended price:
//
//	m × (P×(1−h) + C×h) / (Pr×(1−hr) + Cr×hr)
//
// P/C are the unmultiplied quoted prices per million, h the observed cache
// rate, Pr/Cr/hr the reference. It returns nil when the denominator is zero.
func ProviderHallPredictedRate(m, P, C, h, Pr, Cr, hr decimal.Decimal) *decimal.Decimal {
	one := decimal.NewFromInt(1)
	numerator := P.Mul(one.Sub(h)).Add(C.Mul(h)).Mul(m)
	denominator := Pr.Mul(one.Sub(hr)).Add(Cr.Mul(hr))
	if !denominator.IsPositive() || numerator.IsNegative() {
		return nil
	}
	v := numerator.Div(denominator)
	return &v
}

// ProviderHallPredictedRateInput carries the row-level inputs.
type ProviderHallPredictedRateInput struct {
	Quote     ProviderHallQuote
	Snapshot  *ProviderHallSnapshotRow // default-profile snapshot
	CacheRate ProviderHallMetric[float64]
	Reference ProviderHallReference
	Now       time.Time
}

// ProviderHallPredictedRateMetric applies the fixed rule set:
// not_applicable when the quote is unavailable, cache creation tokens are
// present, the cache rate is not published, the reference is missing or its
// blended price is zero; stale when the reference is older than seven days.
func ProviderHallPredictedRateMetric(in ProviderHallPredictedRateInput) ProviderHallMetric[string] {
	m := ProviderHallMetric[string]{State: ProviderHallMetricNotApplicable}
	providerHallMetricWindow(&m, in.Snapshot)
	switch {
	case !in.Quote.Applicable:
		m.ReasonCode = ProviderHallReasonQuoteUnavailable
		return m
	case in.Snapshot != nil && in.Snapshot.CacheCreationTokens > 0:
		m.ReasonCode = ProviderHallReasonCacheCreation
		return m
	case in.CacheRate.Value == nil || (in.CacheRate.State != ProviderHallMetricOK && in.CacheRate.State != ProviderHallMetricStale):
		m.ReasonCode = ProviderHallReasonCacheRateMissing
		return m
	case in.Reference.InputPrice == nil || in.Reference.CachePrice == nil || in.Reference.CacheRate == nil || in.Reference.ConfirmedAt == nil:
		m.ReasonCode = ProviderHallReasonReferenceMissing
		return m
	}
	// The multiplier is already folded into the quote, so m = 1 here.
	h := decimal.NewFromFloat(*in.CacheRate.Value)
	rate := ProviderHallPredictedRate(decimal.NewFromInt(1), in.Quote.input, in.Quote.cache, h, *in.Reference.InputPrice, *in.Reference.CachePrice, *in.Reference.CacheRate)
	if rate == nil {
		m.ReasonCode = ProviderHallReasonReferenceZero
		return m
	}
	v := rate.StringFixed(4)
	m.Value = &v
	switch {
	case in.Now.Sub(*in.Reference.ConfirmedAt) > ProviderHallReferenceMaxAge:
		m.State, m.ReasonCode = ProviderHallMetricStale, ProviderHallReasonReferenceExpired
	case in.CacheRate.State == ProviderHallMetricStale:
		m.State, m.ReasonCode = ProviderHallMetricStale, ProviderHallReasonStale
	default:
		m.State, m.ReasonCode = ProviderHallMetricOK, ""
	}
	return m
}

// ProviderHallHistoricalPriceMetric adds the unit rule to the price metric:
// bills that were all subscription quota or all balance are publishable;
// a mix has no single unit and is not applicable.
func ProviderHallHistoricalPriceMetric(in ProviderHallMetricInput) (ProviderHallMetric[string], string) {
	m := ProviderHallPriceMetric(in)
	if m.Value == nil {
		return m, ""
	}
	share := in.Snapshot.SubscriptionShare
	switch {
	case share == nil || share.IsZero():
		return m, ProviderHallUnitUSDPerMillion
	case share.Equal(decimal.NewFromInt(1)):
		return m, ProviderHallUnitQuotaPerMillion
	}
	m.Value = nil
	m.State, m.ReasonCode = ProviderHallMetricNotApplicable, ProviderHallReasonMixedBilling
	return m, ""
}

// ProviderHallEffectiveMultiplier is the user's group multiplier times the
// peak factor at the pricing instant, rendered as a short decimal string.
func ProviderHallEffectiveMultiplier(userMultiplier, peak float64) (float64, string) {
	if userMultiplier < 0 {
		userMultiplier = 0
	}
	if peak < 0 {
		peak = 0
	}
	m := decimal.NewFromFloat(userMultiplier).Mul(decimal.NewFromFloat(peak))
	f, _ := m.Float64()
	return f, m.String()
}
