package resources

import _ "embed"

// modelPricingFallbackJSON is bundled into release binaries so hot updates that
// replace only /app/sub2api still carry SmartAPI fallback pricing additions.
//
//go:embed model-pricing/model_prices_and_context_window.json
var modelPricingFallbackJSON []byte

func ModelPricingFallbackJSON() []byte {
	return modelPricingFallbackJSON
}
