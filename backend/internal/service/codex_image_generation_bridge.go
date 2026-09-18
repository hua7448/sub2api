package service

import "strings"

const (
	featureKeyCodexImageGenerationBridge             = "codex_image_generation_bridge"
	featureKeyCodexImageGenerationBridgeLegacy       = "codex_image_generation_bridge_enabled"
	featureKeyCodexImageGenerationExplicitToolPolicy = "codex_image_generation_explicit_tool_policy"

	codexImageGenerationExplicitToolPolicyAllow = "allow"
	codexImageGenerationExplicitToolPolicyStrip = "strip"
)

// ForceCodexImageGenerationBridge returns a copy of extra configured for the
// account-level "force inject" mode. Legacy and blocking fields are removed so
// they cannot override the canonical setting.
func ForceCodexImageGenerationBridge(extra map[string]any) map[string]any {
	out := make(map[string]any, len(extra)+1)
	for key, value := range extra {
		out[key] = value
	}
	out[featureKeyCodexImageGenerationBridge] = true
	delete(out, featureKeyCodexImageGenerationBridgeLegacy)
	delete(out, featureKeyCodexImageGenerationExplicitToolPolicy)
	return out
}

func boolOverridePtr(v bool) *bool {
	return &v
}

func boolOverrideFromMap(values map[string]any, keys ...string) *bool {
	if values == nil {
		return nil
	}
	for _, key := range keys {
		if v, ok := values[key].(bool); ok {
			return boolOverridePtr(v)
		}
	}
	return nil
}

func stringOverrideFromMap(values map[string]any, keys ...string) (string, bool) {
	if values == nil {
		return "", false
	}
	for _, key := range keys {
		if v, ok := values[key].(string); ok {
			return v, true
		}
	}
	return "", false
}

func normalizeCodexImageGenerationExplicitToolPolicy(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case codexImageGenerationExplicitToolPolicyStrip, "remove", "drop":
		return codexImageGenerationExplicitToolPolicyStrip
	default:
		return codexImageGenerationExplicitToolPolicyAllow
	}
}

func platformBoolOverride(values map[string]any, key string, platform string) *bool {
	if values == nil {
		return nil
	}
	if v, ok := values[key].(bool); ok {
		return boolOverridePtr(v)
	}
	raw, ok := values[key].(map[string]any)
	if !ok {
		return nil
	}
	platform = strings.TrimSpace(platform)
	if platform == "" {
		return nil
	}
	if v, ok := raw[platform].(bool); ok {
		return boolOverridePtr(v)
	}
	return nil
}

// CodexImageGenerationBridgeOverride returns the channel-level override for Codex
// image_generation bridge injection. Nil means follow the global/account policy.
func (c *Channel) CodexImageGenerationBridgeOverride(platform string) *bool {
	if c == nil {
		return nil
	}
	return platformBoolOverride(c.FeaturesConfig, featureKeyCodexImageGenerationBridge, platform)
}

// CodexImageGenerationBridgeOverride returns the account-level override for Codex
// image_generation bridge injection. Nil means follow the channel/global policy.
func (a *Account) CodexImageGenerationBridgeOverride() *bool {
	if a == nil || a.Platform != PlatformOpenAI || a.Extra == nil {
		return nil
	}
	if override := boolOverrideFromMap(a.Extra, featureKeyCodexImageGenerationBridge, featureKeyCodexImageGenerationBridgeLegacy); override != nil {
		return override
	}
	openaiConfig, _ := a.Extra[PlatformOpenAI].(map[string]any)
	return boolOverrideFromMap(openaiConfig, featureKeyCodexImageGenerationBridge, featureKeyCodexImageGenerationBridgeLegacy)
}

// CodexImageGenerationExplicitToolPolicy returns the account-level policy for
// client-provided Codex /responses image_generation tools. Unknown or unset
// values default to allow to preserve existing behavior.
func (a *Account) CodexImageGenerationExplicitToolPolicy() string {
	if a == nil || a.Platform != PlatformOpenAI || a.Extra == nil {
		return codexImageGenerationExplicitToolPolicyAllow
	}
	if policy, ok := stringOverrideFromMap(a.Extra, featureKeyCodexImageGenerationExplicitToolPolicy); ok {
		return normalizeCodexImageGenerationExplicitToolPolicy(policy)
	}
	openaiConfig, _ := a.Extra[PlatformOpenAI].(map[string]any)
	if policy, ok := stringOverrideFromMap(openaiConfig, featureKeyCodexImageGenerationExplicitToolPolicy); ok {
		return normalizeCodexImageGenerationExplicitToolPolicy(policy)
	}
	return codexImageGenerationExplicitToolPolicyAllow
}
