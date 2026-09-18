package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestForceCodexImageGenerationBridge(t *testing.T) {
	extra := map[string]any{
		"keep":                                        "value",
		"codex_image_generation_bridge":               false,
		"codex_image_generation_bridge_enabled":       false,
		"codex_image_generation_explicit_tool_policy": "strip",
	}

	got := ForceCodexImageGenerationBridge(extra)

	require.Equal(t, true, got[featureKeyCodexImageGenerationBridge])
	require.Equal(t, "value", got["keep"])
	require.NotContains(t, got, featureKeyCodexImageGenerationBridgeLegacy)
	require.NotContains(t, got, featureKeyCodexImageGenerationExplicitToolPolicy)
	require.Equal(t, false, extra["codex_image_generation_bridge"], "input map must not be mutated")
}
