package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestAccountFromServiceShallow_RedactsSensitiveCredentials(t *testing.T) {
	src := &service.Account{
		ID:       42,
		Name:     "demo",
		Platform: "anthropic",
		Type:     "oauth",
		Credentials: map[string]any{
			"access_token":  "at-secret",
			"refresh_token": "rt-secret",
			"id_token":      "id-secret",
			"api_key":       "sk-secret",
			"base_url":      "https://api.example.com",
			"model_mapping": map[string]any{"foo": "bar"},
		},
	}

	got := AccountFromServiceShallow(src)
	require.NotNil(t, got)

	// 敏感键不在 Credentials 里
	require.NotContains(t, got.Credentials, "access_token")
	require.NotContains(t, got.Credentials, "refresh_token")
	require.NotContains(t, got.Credentials, "id_token")
	require.NotContains(t, got.Credentials, "api_key")
	// 非敏感键保留
	require.Equal(t, "https://api.example.com", got.Credentials["base_url"])
	require.Equal(t, map[string]any{"foo": "bar"}, got.Credentials["model_mapping"])

	// 状态 map 标记敏感键存在
	require.True(t, got.CredentialsStatus["has_access_token"])
	require.True(t, got.CredentialsStatus["has_refresh_token"])
	require.True(t, got.CredentialsStatus["has_id_token"])
	require.True(t, got.CredentialsStatus["has_api_key"])

	// JSON 序列化校验：响应体里不会出现敏感子串
	raw, err := json.Marshal(got)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "rt-secret")
	require.NotContains(t, string(raw), "at-secret")
	require.NotContains(t, string(raw), "sk-secret")
	require.NotContains(t, string(raw), "id-secret")
	// 状态标识应序列化进 JSON
	require.Contains(t, string(raw), "credentials_status")
	require.Contains(t, string(raw), "has_refresh_token")

	// 原始 service.Account 不应被改动
	require.Equal(t, "rt-secret", src.Credentials["refresh_token"])
}

func TestAccountFromServiceShallow_RedactsOllamaCloudManagedExtra(t *testing.T) {
	snapshot := map[string]any{
		"status":          service.OllamaCloudUsageStatusOK,
		"last_attempt_at": "2026-07-22T12:00:00Z",
		"next_refresh_at": "2026-07-22T13:00:00Z",
		"data":            map[string]any{"plan": "Pro"},
	}
	src := &service.Account{
		ID: 9, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
		Credentials: map[string]any{"base_url": "https://ollama.com", "api_key": "secret-key"},
		Extra: map[string]any{
			service.OllamaCloudUsageSessionExtraKey:     "ciphertext-secret",
			service.OllamaCloudUsageAutoRefreshExtraKey: true,
			service.OllamaCloudUsageSnapshotExtraKey:    snapshot,
			"ordinary":                                  "kept",
		},
	}

	got := AccountFromServiceShallow(src)
	require.NotContains(t, got.Extra, service.OllamaCloudUsageSessionExtraKey)
	require.NotContains(t, got.Extra, service.OllamaCloudUsageAutoRefreshExtraKey)
	require.NotContains(t, got.Extra, service.OllamaCloudUsageSnapshotExtraKey)
	require.Equal(t, "kept", got.Extra["ordinary"])
	require.NotNil(t, got.OllamaCloudUsage)
	require.True(t, got.OllamaCloudUsage.Configured)
	require.True(t, got.OllamaCloudUsage.AutoRefreshEnabled)
	require.Equal(t, "Pro", got.OllamaCloudUsage.Snapshot.Data.Plan)

	raw, err := json.Marshal(got)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "ciphertext-secret")
	require.NotContains(t, string(raw), "secret-key")
	require.Contains(t, src.Extra, service.OllamaCloudUsageSessionExtraKey)
}

func TestAccountFromServiceShallow_NilCredentialsOmitsStatus(t *testing.T) {
	src := &service.Account{ID: 1, Name: "n", Platform: "anthropic", Type: "oauth"}
	got := AccountFromServiceShallow(src)
	require.NotNil(t, got)
	require.Nil(t, got.Credentials)
	require.Nil(t, got.CredentialsStatus)
}

func TestAccountFromServiceShallow_KeepActiveSuppressesAutomaticUnavailableState(t *testing.T) {
	future := time.Now().Add(10 * time.Minute)
	src := &service.Account{
		ID:                      77,
		Status:                  service.StatusActive,
		Schedulable:             true,
		RateLimitedAt:           &future,
		RateLimitResetAt:        &future,
		OverloadUntil:           &future,
		TempUnschedulableUntil:  &future,
		TempUnschedulableReason: "automatic upstream cooldown",
		Extra: map[string]any{
			service.KeepStatusActiveExtraKey: true,
			"model_rate_limits": map[string]any{
				"gpt-5": map[string]any{"rate_limit_reset_at": future.Format(time.RFC3339)},
			},
			"antigravity_quota_scopes": map[string]any{"gemini": true},
			"ordinary":                 "kept",
		},
	}

	got := AccountFromServiceShallow(src)
	require.NotNil(t, got)
	require.True(t, got.KeepStatusActive)
	require.Nil(t, got.RateLimitedAt)
	require.Nil(t, got.RateLimitResetAt)
	require.Nil(t, got.OverloadUntil)
	require.Nil(t, got.TempUnschedulableUntil)
	require.Empty(t, got.TempUnschedulableReason)
	require.NotContains(t, got.Extra, "model_rate_limits")
	require.NotContains(t, got.Extra, "antigravity_quota_scopes")
	require.Equal(t, "kept", got.Extra["ordinary"])

	// Mapping is presentation-only: durable audit state remains untouched.
	require.NotNil(t, src.OverloadUntil)
	require.NotNil(t, src.TempUnschedulableUntil)
	require.Contains(t, src.Extra, "model_rate_limits")
	require.Contains(t, src.Extra, "antigravity_quota_scopes")
}

func TestAccountFromServiceShallow_OrdinaryAccountKeepsActiveAutomaticState(t *testing.T) {
	past := time.Now().Add(-10 * time.Minute)
	src := &service.Account{
		RateLimitedAt:           &past,
		RateLimitResetAt:        &past,
		OverloadUntil:           &past,
		TempUnschedulableUntil:  &past,
		TempUnschedulableReason: "automatic upstream cooldown",
		Extra: map[string]any{
			"model_rate_limits":        map[string]any{"gpt-5": map[string]any{"rate_limit_reset_at": past.Format(time.RFC3339)}},
			"antigravity_quota_scopes": map[string]any{"gemini": true},
		},
	}

	got := AccountFromServiceShallow(src)
	require.Nil(t, got.RateLimitedAt)
	require.Nil(t, got.RateLimitResetAt)
	require.Equal(t, &past, got.OverloadUntil)
	require.Equal(t, &past, got.TempUnschedulableUntil)
	require.Equal(t, src.TempUnschedulableReason, got.TempUnschedulableReason)
	require.Contains(t, got.Extra, "model_rate_limits")
	require.Contains(t, got.Extra, "antigravity_quota_scopes")
}

func TestAccountFromServiceShallow_KeepActiveMalformedValueFailsClosed(t *testing.T) {
	future := time.Now().Add(10 * time.Minute)
	src := &service.Account{
		OverloadUntil:          &future,
		TempUnschedulableUntil: &future,
		Extra: map[string]any{
			service.KeepStatusActiveExtraKey: "true",
			"model_rate_limits":              map[string]any{"gpt-5": map[string]any{"rate_limit_reset_at": future.Format(time.RFC3339)}},
		},
	}

	got := AccountFromServiceShallow(src)
	require.False(t, got.KeepStatusActive)
	require.Equal(t, &future, got.OverloadUntil)
	require.Equal(t, &future, got.TempUnschedulableUntil)
	require.Contains(t, got.Extra, "model_rate_limits")
}
