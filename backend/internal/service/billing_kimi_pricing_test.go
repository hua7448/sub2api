package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestBillingServiceKimiK3FallbackPricing(t *testing.T) {
	svc := NewBillingService(&config.Config{}, nil)

	tests := []struct {
		model string
	}{
		{model: "kimi-k3"},
		{model: "kimi-3"},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			pricing, err := svc.GetModelPricing(tt.model)
			require.NoError(t, err)
			require.NotNil(t, pricing)
			require.InDelta(t, 2.80e-6, pricing.InputPricePerToken, 1e-12)
			require.InDelta(t, 14e-6, pricing.OutputPricePerToken, 1e-12)
			require.InDelta(t, 0.28e-6, pricing.CacheReadPricePerToken, 1e-12)
		})
	}
}
