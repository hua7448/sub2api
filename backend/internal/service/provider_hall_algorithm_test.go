//go:build unit

package service

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// Matrix A: fixed samples for the pure algorithms. Values are exact so a
// formula change is a visible diff, never a tolerance drift.
func TestProviderHallExactLatencyFixedSamples(t *testing.T) {
	// 1000, 2000, ..., 19000 plus one 100000 outlier: N=20, fast95 keeps 19
	// samples (mean 10000), P90 rank = 18 -> 18000.
	counts := map[int64]int64{100000: 1}
	for ms := int64(1000); ms <= 19000; ms += 1000 {
		counts[ms] = 1
	}
	got := ProviderHallExactLatency(counts)
	require.Equal(t, int64(20), got.Count)
	require.NotNil(t, got.Fast95MeanMS)
	require.True(t, got.Fast95MeanMS.Equal(decimal.NewFromInt(10000)), got.Fast95MeanMS.String())
	require.NotNil(t, got.P90MS)
	require.Equal(t, int64(18000), *got.P90MS)

	t.Run("n_zero", func(t *testing.T) {
		got := ProviderHallExactLatency(nil)
		require.Zero(t, got.Count)
		require.Nil(t, got.Fast95MeanMS)
		require.Nil(t, got.P90MS)
		got = ProviderHallExactLatency(map[int64]int64{0: 5, -3: 2, 500: 0})
		require.Zero(t, got.Count, "non-positive latency or count never contributes")
	})
	t.Run("n_one", func(t *testing.T) {
		got := ProviderHallExactLatency(map[int64]int64{750: 1})
		require.Equal(t, int64(1), got.Count)
		require.True(t, got.Fast95MeanMS.Equal(decimal.NewFromInt(750)))
		require.Equal(t, int64(750), *got.P90MS)
	})
	t.Run("n_ten_drops_nothing", func(t *testing.T) {
		// N=10: fast95 = 10 - 10/20 = 10 samples, P90 rank = 10 - 1 = 9.
		counts := map[int64]int64{}
		for ms := int64(100); ms <= 1000; ms += 100 {
			counts[ms] = 1
		}
		got := ProviderHallExactLatency(counts)
		require.Equal(t, int64(10), got.Count)
		require.True(t, got.Fast95MeanMS.Equal(decimal.NewFromInt(550)), got.Fast95MeanMS.String())
		require.Equal(t, int64(900), *got.P90MS)
	})
	t.Run("n_twenty_drops_one", func(t *testing.T) {
		counts := map[int64]int64{}
		for ms := int64(100); ms <= 2000; ms += 100 {
			counts[ms] = 1
		}
		got := ProviderHallExactLatency(counts)
		require.Equal(t, int64(20), got.Count)
		// mean of 100..1900 = 1000
		require.True(t, got.Fast95MeanMS.Equal(decimal.NewFromInt(1000)), got.Fast95MeanMS.String())
		require.Equal(t, int64(1800), *got.P90MS)
	})
	t.Run("duplicate_boundary_values", func(t *testing.T) {
		// 19 samples at 800 and 1 at 1100: fast95 drops one of the slowest
		// (the 1100), P90 rank 18 falls inside the 800 run.
		got := ProviderHallExactLatency(map[int64]int64{800: 19, 1100: 1})
		require.Equal(t, int64(20), got.Count)
		require.True(t, got.Fast95MeanMS.Equal(decimal.NewFromInt(800)), got.Fast95MeanMS.String())
		require.Equal(t, int64(800), *got.P90MS)
		// 1 at 800, 19 at 1100: drop one 1100; mean = (800 + 18*1100)/19.
		got = ProviderHallExactLatency(map[int64]int64{800: 1, 1100: 19})
		want := decimal.NewFromInt(800 + 18*1100).Div(decimal.NewFromInt(19))
		require.True(t, got.Fast95MeanMS.Equal(want), got.Fast95MeanMS.String())
		require.Equal(t, int64(1100), *got.P90MS)
	})
	t.Run("input_map_not_mutated", func(t *testing.T) {
		counts := map[int64]int64{800: 2, 1100: 3}
		_ = ProviderHallExactLatency(counts)
		require.Equal(t, map[int64]int64{800: 2, 1100: 3}, counts)
	})
}

func TestProviderHallSuccessRateFixedSamples(t *testing.T) {
	rate := ProviderHallSuccessRate(2, 1, 3)
	require.NotNil(t, rate)
	require.True(t, rate.Equal(decimal.NewFromInt(2).Div(decimal.NewFromInt(3))), rate.String())
	// Retries: 2 successes over 5 submissions, denominator = max(5, 2+1).
	rate = ProviderHallSuccessRate(2, 1, 5)
	require.True(t, rate.Equal(decimal.NewFromInt(2).Div(decimal.NewFromInt(5))), rate.String())
	// Submissions undercount (no-account failures): denominator = success+failed.
	rate = ProviderHallSuccessRate(2, 2, 0)
	require.True(t, rate.Equal(decimal.NewFromFloat(0.5)))
	require.Nil(t, ProviderHallSuccessRate(0, 0, 0))
	require.Nil(t, ProviderHallSuccessRate(-1, 0, 0))
	require.Nil(t, ProviderHallSuccessRate(0, 0, -1))
	require.True(t, ProviderHallSuccessRate(0, 3, 3).IsZero())
	require.True(t, ProviderHallSuccessRate(3, 0, 3).Equal(decimal.NewFromInt(1)))
}

func TestProviderHallCacheRateFixedSamples(t *testing.T) {
	// 800 ordinary + 200 cache read + 0 creation -> 0.2
	rate := ProviderHallCacheRate(800, 200, 0)
	require.NotNil(t, rate)
	require.True(t, rate.Equal(decimal.NewFromFloat(0.2)), rate.String())
	// Cache creation dilutes the rate: 600 + 200 + 200 -> 0.2
	rate = ProviderHallCacheRate(600, 200, 200)
	require.True(t, rate.Equal(decimal.NewFromFloat(0.2)), rate.String())
	require.Nil(t, ProviderHallCacheRate(0, 0, 0), "no input tokens: no rate")
	require.Nil(t, ProviderHallCacheRate(-1, 0, 0))
	require.Nil(t, ProviderHallCacheRate(0, -1, 0))
	require.Nil(t, ProviderHallCacheRate(0, 0, -1))
	require.True(t, ProviderHallCacheRate(0, 10, 0).Equal(decimal.NewFromInt(1)))
	require.True(t, ProviderHallCacheRate(10, 0, 0).IsZero())
}

func TestProviderHallInputCostAndPrice(t *testing.T) {
	d := func(s string) decimal.Decimal { return decimal.RequireFromString(s) }
	// $1 per million: actual 0.001 for 1000 input tokens with all cost in input.
	cost := ProviderHallInputCost(d("0.001"), d("0.001"), d("0.001"))
	require.NotNil(t, cost)
	price := ProviderHallInputPrice(*cost, 1000)
	require.NotNil(t, price)
	require.True(t, price.Equal(decimal.NewFromInt(1)), price.String())

	// Input share: actual 3, input base 1, total base 4 -> input cost 0.75.
	cost = ProviderHallInputCost(d("3"), d("1"), d("4"))
	require.True(t, cost.Equal(d("0.75")), cost.String())
	// Multiplier applied through actual: base 1/4 at 2x -> 1.5 input cost.
	cost = ProviderHallInputCost(d("8"), d("1"), d("4"))
	require.True(t, cost.Equal(d("2")), cost.String())
	// Legitimate zero bill.
	cost = ProviderHallInputCost(decimal.Zero, d("1"), d("4"))
	require.NotNil(t, cost)
	require.True(t, cost.IsZero())
	// Rejections: negative, zero total base, input greater than total.
	require.Nil(t, ProviderHallInputCost(d("-1"), d("1"), d("4")))
	require.Nil(t, ProviderHallInputCost(d("1"), d("-1"), d("4")))
	require.Nil(t, ProviderHallInputCost(d("1"), d("1"), decimal.Zero))
	require.Nil(t, ProviderHallInputCost(d("1"), d("5"), d("4")))

	// Price rejects zero/negative tokens and negative cost; keeps precision.
	require.Nil(t, ProviderHallInputPrice(d("1"), 0))
	require.Nil(t, ProviderHallInputPrice(d("1"), -5))
	require.Nil(t, ProviderHallInputPrice(d("-1"), 5))
	price = ProviderHallInputPrice(d("0.0000000123"), 3)
	require.True(t, price.Equal(d("0.0000000123").Mul(decimal.NewFromInt(1_000_000)).Div(decimal.NewFromInt(3))), price.String())
	require.True(t, ProviderHallInputPrice(decimal.Zero, 7).IsZero())
}

func TestProviderHallBudgetDay(t *testing.T) {
	// 2026-09-11 16:00 UTC is 2026-09-12 00:00 in Asia/Shanghai.
	require.Equal(t, "2026-09-12", ProviderHallBudgetDay(time.Date(2026, 9, 11, 16, 0, 0, 0, time.UTC)))
	require.Equal(t, "2026-09-11", ProviderHallBudgetDay(time.Date(2026, 9, 11, 15, 59, 59, 0, time.UTC)))
	// Input zone never matters.
	ny := time.FixedZone("America/New_York", -4*60*60)
	require.Equal(t, "2026-09-12", ProviderHallBudgetDay(time.Date(2026, 9, 11, 12, 0, 0, 0, ny)))
}
