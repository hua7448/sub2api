package service

import (
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

type ProviderHallLatency struct {
	Count        int64
	Fast95MeanMS *decimal.Decimal
	P90MS        *int64
}

// ProviderHallExactLatency merges exact millisecond frequencies. The caller
// supplies only successful streaming samples, never coarse histogram buckets.
func ProviderHallExactLatency(counts map[int64]int64) ProviderHallLatency {
	result := ProviderHallLatency{}
	values := make([]int64, 0, len(counts))
	for ms, count := range counts {
		if ms > 0 && count > 0 {
			result.Count += count
			values = append(values, ms)
		}
	}
	if result.Count == 0 {
		return result
	}
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	// ceil(95*N/100) and ceil(90*N/100), without float rounding.
	fastCount := result.Count - result.Count/20
	p90Rank := result.Count - result.Count/10
	remaining, cumulative := fastCount, int64(0)
	sum := decimal.Zero
	for _, ms := range values {
		count := counts[ms]
		take := min(remaining, count)
		sum = sum.Add(decimal.NewFromInt(ms).Mul(decimal.NewFromInt(take)))
		remaining -= take
		cumulative += count
		if result.P90MS == nil && cumulative >= p90Rank {
			value := ms
			result.P90MS = &value
		}
	}
	mean := sum.Div(decimal.NewFromInt(fastCount))
	result.Fast95MeanMS = &mean
	return result
}

func ProviderHallCacheRate(ordinaryInput, cacheRead, cacheCreation int64) *decimal.Decimal {
	if ordinaryInput < 0 || cacheRead < 0 || cacheCreation < 0 {
		return nil
	}
	total := decimal.NewFromInt(ordinaryInput).Add(decimal.NewFromInt(cacheRead)).Add(decimal.NewFromInt(cacheCreation))
	if total.IsZero() {
		return nil
	}
	rate := decimal.NewFromInt(cacheRead).Div(total)
	return &rate
}

func ProviderHallSuccessRate(success, failed, submissions int64) *decimal.Decimal {
	if success < 0 || failed < 0 || submissions < 0 {
		return nil
	}
	denominator := decimal.Max(decimal.NewFromInt(submissions), decimal.NewFromInt(success).Add(decimal.NewFromInt(failed)))
	if denominator.IsZero() {
		return nil
	}
	rate := decimal.NewFromInt(success).Div(denominator)
	return &rate
}

// ProviderHallInputCost accepts only confirmed bills selected by the collector.
// A zero actual bill is legitimate; missing/uncertain bills must never call it.
func ProviderHallInputCost(actual, inputBase, totalBase decimal.Decimal) *decimal.Decimal {
	if actual.IsNegative() || inputBase.IsNegative() || !totalBase.IsPositive() || inputBase.GreaterThan(totalBase) {
		return nil
	}
	cost := actual.Mul(inputBase).Div(totalBase)
	return &cost
}

func ProviderHallInputPrice(inputCost decimal.Decimal, inputTokens int64) *decimal.Decimal {
	if inputTokens <= 0 || inputCost.IsNegative() {
		return nil
	}
	price := inputCost.Mul(decimal.NewFromInt(1_000_000)).Div(decimal.NewFromInt(inputTokens))
	return &price
}

func ProviderHallBudgetDay(dispatchedAt time.Time) string {
	// Asia/Shanghai uses UTC+8 throughout the supported (post-2026) task dates.
	return dispatchedAt.In(time.FixedZone("Asia/Shanghai", 8*60*60)).Format("2006-01-02")
}
