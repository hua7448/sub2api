//go:build unit

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAccountIsSchedulable_QuotaExceeded(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		account *Account
		want    bool
	}{
		{
			name: "apikey daily quota exceeded",
			account: &Account{
				Status:      StatusActive,
				Schedulable: true,
				Type:        AccountTypeAPIKey,
				Extra: map[string]any{
					"quota_daily_limit": 10.0,
					"quota_daily_used":  10.0,
					"quota_daily_start": now.Add(-1 * time.Hour).Format(time.RFC3339),
				},
			},
			want: false,
		},
		{
			name: "apikey weekly quota exceeded",
			account: &Account{
				Status:      StatusActive,
				Schedulable: true,
				Type:        AccountTypeAPIKey,
				Extra: map[string]any{
					"quota_weekly_limit": 50.0,
					"quota_weekly_used":  50.0,
					"quota_weekly_start": now.Add(-2 * 24 * time.Hour).Format(time.RFC3339),
				},
			},
			want: false,
		},
		{
			name: "apikey total quota exceeded",
			account: &Account{
				Status:      StatusActive,
				Schedulable: true,
				Type:        AccountTypeAPIKey,
				Extra: map[string]any{
					"quota_limit": 100.0,
					"quota_used":  100.0,
				},
			},
			want: false,
		},
		{
			name: "apikey quota not exceeded",
			account: &Account{
				Status:      StatusActive,
				Schedulable: true,
				Type:        AccountTypeAPIKey,
				Extra: map[string]any{
					"quota_daily_limit": 10.0,
					"quota_daily_used":  5.0,
					"quota_daily_start": now.Add(-1 * time.Hour).Format(time.RFC3339),
				},
			},
			want: true,
		},
		{
			name: "apikey expired daily period restores schedulable",
			account: &Account{
				Status:      StatusActive,
				Schedulable: true,
				Type:        AccountTypeAPIKey,
				Extra: map[string]any{
					"quota_daily_limit": 10.0,
					"quota_daily_used":  10.0,
					"quota_daily_start": now.Add(-25 * time.Hour).Format(time.RFC3339),
				},
			},
			want: true,
		},
		{
			name: "oauth ignores quota exceeded",
			account: &Account{
				Status:      StatusActive,
				Schedulable: true,
				Type:        AccountTypeOAuth,
				Extra: map[string]any{
					"quota_daily_limit": 10.0,
					"quota_daily_used":  10.0,
					"quota_daily_start": now.Add(-1 * time.Hour).Format(time.RFC3339),
				},
			},
			want: true,
		},
		{
			name: "bedrock quota exceeded",
			account: &Account{
				Status:      StatusActive,
				Schedulable: true,
				Type:        AccountTypeBedrock,
				Extra: map[string]any{
					"quota_limit": 200.0,
					"quota_used":  200.0,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.account.IsSchedulable())
		})
	}
}

func TestOpenAICodexWindowResetDoesNotCountAsRuntimeRateLimit(t *testing.T) {
	now := time.Now()
	resetAt := now.Add(7 * 24 * time.Hour)
	account := &Account{
		Platform:         PlatformOpenAI,
		Status:           StatusActive,
		Schedulable:      true,
		Type:             AccountTypeOAuth,
		RateLimitResetAt: &resetAt,
		RateLimitedAt:    &now,
		Extra: map[string]any{
			"codex_5h_used_percent": 99.0,
			"codex_7d_used_percent": 31.0,
			"codex_5h_reset_at":     now.Add(5 * time.Hour).Format(time.RFC3339),
			"codex_7d_reset_at":     resetAt.Add(30 * time.Second).Format(time.RFC3339),
		},
	}

	require.False(t, account.IsRateLimited())
	require.True(t, account.IsSchedulable())
}

func TestOpenAICodexExhaustedWindowStillCountsAsRuntimeRateLimit(t *testing.T) {
	now := time.Now()
	resetAt := now.Add(7 * 24 * time.Hour)
	account := &Account{
		Platform:         PlatformOpenAI,
		Status:           StatusActive,
		Schedulable:      true,
		Type:             AccountTypeOAuth,
		RateLimitResetAt: &resetAt,
		RateLimitedAt:    &now,
		Extra: map[string]any{
			"codex_5h_used_percent": 1.0,
			"codex_7d_used_percent": 100.0,
			"codex_5h_reset_at":     now.Add(5 * time.Hour).Format(time.RFC3339),
			"codex_7d_reset_at":     resetAt.Format(time.RFC3339),
		},
	}

	require.True(t, account.IsRateLimited())
	require.False(t, account.IsSchedulable())
}
