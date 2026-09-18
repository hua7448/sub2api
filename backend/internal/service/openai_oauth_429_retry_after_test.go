//go:build unit

package service

import (
	"context"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestOpenAIOAuth429RetryAfterPersistsAndBlocksUntilDeadline(t *testing.T) {
	for _, format := range []string{"seconds", "http-date"} {
		t.Run(format, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx := context.Background()
				primary := Account{ID: 42901, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 0}
				backup := Account{ID: 42902, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 5}
				repo := &oauth429RetryAfterRepo{}
				rl := NewRateLimitService(repo, nil, nil, nil, nil)
				gateway := &OpenAIGatewayService{
					accountRepo: schedulerTestOpenAIAccountRepo{accounts: []Account{primary, backup}},
					cfg:         &config.Config{}, rateLimitService: rl,
				}
				rl.SetAccountRuntimeBlocker(gateway)
				deadline := time.Now().Add(time.Minute)
				headers := make(http.Header)
				headers.Set("Retry-After", "60")
				if format == "http-date" {
					headers.Set("Retry-After", deadline.UTC().Format(http.TimeFormat))
				}

				require.False(t, gateway.handleOpenAIAccountUpstreamError(ctx, &primary, http.StatusTooManyRequests, headers, nil))
				require.Equal(t, primary.ID, repo.rateLimitedID)
				require.True(t, deadline.Equal(repo.rateLimitedAt), "reset = %v, want %v", repo.rateLimitedAt, deadline)
				blockedUntil, ok := gateway.openaiAccountRuntimeBlockUntil.Load(primary.ID)
				require.True(t, ok)
				require.Equal(t, repo.rateLimitedAt, blockedUntil, "the DB update and immediate scheduler gate must share the exact response deadline")

				// Cross the old five-second fallback without waiting in wall time.
				// The actual account selector must still choose the healthy backup.
				time.Sleep(6 * time.Second)
				selected, err := gateway.SelectAccountForModelWithExclusions(ctx, nil, "", "gpt-5.1", nil)
				require.NoError(t, err)
				require.Equal(t, backup.ID, selected.ID)

				// A newly loaded account on another process must also remain blocked,
				// independently of this gateway's in-memory cooldown map.
				persisted := primary
				persisted.RateLimitResetAt = &repo.rateLimitedAt
				require.True(t, persisted.IsRateLimited())
				require.False(t, persisted.IsSchedulable())
				freshGateway := &OpenAIGatewayService{cfg: &config.Config{}, accountRepo: schedulerTestOpenAIAccountRepo{accounts: []Account{persisted, backup}}}
				selected, err = freshGateway.SelectAccountForModelWithExclusions(ctx, nil, "", "gpt-5.1", nil)
				require.NoError(t, err)
				require.Equal(t, backup.ID, selected.ID)

				time.Sleep(55 * time.Second)
				selected, err = gateway.SelectAccountForModelWithExclusions(ctx, nil, "", "gpt-5.1", nil)
				require.NoError(t, err)
				require.Equal(t, primary.ID, selected.ID)
				require.False(t, persisted.IsRateLimited())
			})
		})
	}
}

func TestOpenAIOAuth429RetryAfterFastPathWithoutRateLimitService(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		gateway := &OpenAIGatewayService{}
		account := &Account{ID: 42903, Platform: PlatformOpenAI, Type: AccountTypeOAuth}
		headers := make(http.Header)
		headers.Set("Retry-After", "60")
		gateway.handleOpenAIAccountUpstreamError(context.Background(), account, http.StatusTooManyRequests, headers, nil)
		time.Sleep(6 * time.Second)
		require.True(t, gateway.isOpenAIAccountRuntimeBlocked(account))
		time.Sleep(55 * time.Second)
		require.False(t, gateway.isOpenAIAccountRuntimeBlocked(account))
	})
}

type oauth429RetryAfterRepo struct {
	openAI429SnapshotRepo
	extra      map[string]any
	clearCalls int
}

func (r *oauth429RetryAfterRepo) UpdateExtra(_ context.Context, _ int64, updates map[string]any) error {
	if r.extra == nil {
		r.extra = make(map[string]any)
	}
	for key, value := range updates {
		r.extra[key] = value
	}
	return nil
}

func (r *oauth429RetryAfterRepo) ClearRateLimit(_ context.Context, _ int64) error {
	r.clearCalls++
	return nil
}

func TestOpenAIOAuth429RetryAfterSurvivesNonExhaustedQuotaSnapshot(t *testing.T) {
	for _, quotaReset := range []string{"60", "3600"} {
		t.Run(quotaReset, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				repo := &oauth429RetryAfterRepo{}
				rl := NewRateLimitService(repo, nil, nil, nil, nil)
				account := &Account{ID: 42904, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true}
				headers := make(http.Header)
				headers.Set("Retry-After", "60")
				headers.Set("x-codex-primary-used-percent", "30")
				headers.Set("x-codex-primary-reset-after-seconds", quotaReset)
				headers.Set("x-codex-primary-window-minutes", "300")
				headers.Set("x-codex-secondary-used-percent", "40")
				headers.Set("x-codex-secondary-reset-after-seconds", "604800")
				headers.Set("x-codex-secondary-window-minutes", "10080")
				rl.handle429(context.Background(), account, headers, nil)

				persisted := *account
				persisted.RateLimitResetAt = &repo.rateLimitedAt
				persisted.Extra = repo.extra
				require.True(t, persisted.IsRateLimited(), "available quota must not override an explicit request-rate cooldown")
				require.False(t, persisted.IsSchedulable())
				usage := &AccountUsageService{accountRepo: repo}
				usage.clearOpenAIRateLimitIfCodexSnapshotRecovered(context.Background(), &persisted)
				require.Zero(t, repo.clearCalls, "quota refresh must preserve Retry-After until its deadline")
				require.NotNil(t, persisted.RateLimitResetAt)
				time.Sleep(61 * time.Second)
				require.False(t, persisted.IsRateLimited())
				require.True(t, persisted.IsSchedulable())
			})
		})
	}
}

func TestOpenAIOAuth429RetryAfterDeadlinePrecedenceAndInvalidValues(t *testing.T) {
	for _, tc := range []struct {
		name         string
		retryAfter   string
		usedPercent  string
		codexReset   string
		body         string
		wantCooldown time.Duration
	}{
		{name: "retry shorter than exhausted codex", retryAfter: "60", usedPercent: "100", codexReset: "3600", wantCooldown: time.Hour},
		{name: "retry longer than exhausted codex", retryAfter: "3600", usedPercent: "100", codexReset: "60", wantCooldown: time.Hour},
		{name: "retry shorter than body", retryAfter: "60", body: `{"error":{"type":"usage_limit_reached","resets_in_seconds":120}}`, wantCooldown: 2 * time.Minute},
		{name: "retry longer than body", retryAfter: "120", body: `{"error":{"type":"rate_limit_exceeded","resets_in_seconds":60}}`, wantCooldown: 2 * time.Minute},
		{name: "body longer than exhausted codex", retryAfter: "60", usedPercent: "100", codexReset: "3600", body: `{"error":{"type":"usage_limit_reached","resets_in_seconds":7200}}`, wantCooldown: 2 * time.Hour},
		{name: "ignore echoed unexhausted window", retryAfter: "60", usedPercent: "30", codexReset: "3600", body: `{"error":{"type":"rate_limit_exceeded","resets_in_seconds":3600}}`, wantCooldown: time.Minute},
		{name: "unexhausted window without retry retains fallback", usedPercent: "30", codexReset: "3600", body: `{"error":{"type":"rate_limit_exceeded","resets_in_seconds":3600}}`, wantCooldown: 5 * time.Second},
		{name: "fractional seconds", retryAfter: "0.25", wantCooldown: 250 * time.Millisecond},
		{name: "zero", retryAfter: "0", wantCooldown: 5 * time.Second},
		{name: "negative", retryAfter: "-60", wantCooldown: 5 * time.Second},
		{name: "NaN", retryAfter: "NaN", wantCooldown: 5 * time.Second},
		{name: "infinity", retryAfter: "+Inf", wantCooldown: 5 * time.Second},
		{name: "overflow", retryAfter: "9223372036854775807", wantCooldown: 5 * time.Second},
		{name: "malformed", retryAfter: "later", wantCooldown: 5 * time.Second},
		{name: "expired date", retryAfter: "Fri, 31 Dec 1999 23:59:59 GMT", wantCooldown: 5 * time.Second},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				repo := &oauth429RetryAfterRepo{}
				rl := NewRateLimitService(repo, nil, nil, nil, nil)
				gateway := &OpenAIGatewayService{rateLimitService: rl}
				rl.SetAccountRuntimeBlocker(gateway)
				account := &Account{ID: 42905, Platform: PlatformOpenAI, Type: AccountTypeOAuth}
				headers := make(http.Header)
				headers.Set("Retry-After", tc.retryAfter)
				if tc.usedPercent != "" {
					headers.Set("x-codex-primary-used-percent", tc.usedPercent)
					headers.Set("x-codex-primary-reset-after-seconds", tc.codexReset)
					headers.Set("x-codex-primary-window-minutes", "300")
					headers.Set("x-codex-secondary-used-percent", "40")
					headers.Set("x-codex-secondary-reset-after-seconds", "604800")
					headers.Set("x-codex-secondary-window-minutes", "10080")
				}
				wantReset := time.Now().Add(tc.wantCooldown)
				gateway.handleOpenAIAccountUpstreamError(context.Background(), account, http.StatusTooManyRequests, headers, []byte(tc.body))
				require.True(t, wantReset.Equal(repo.rateLimitedAt), "reset = %v, want %v", repo.rateLimitedAt, wantReset)
				blockedUntil, ok := gateway.openaiAccountRuntimeBlockUntil.Load(account.ID)
				require.True(t, ok)
				require.Equal(t, repo.rateLimitedAt, blockedUntil)
			})
		})
	}
}

func TestOpenAIOAuth429RetryAfterRemainsExplicitWhenFallbackDisabled(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		settingRepo := newMockSettingRepo()
		settingRepo.data[SettingKeyRateLimit429CooldownSettings] = `{"enabled":false,"cooldown_seconds":12}`
		repo := &oauth429RetryAfterRepo{}
		rl := NewRateLimitService(repo, nil, nil, nil, nil)
		rl.SetSettingService(NewSettingService(settingRepo, &config.Config{}))
		gateway := &OpenAIGatewayService{rateLimitService: rl}
		rl.SetAccountRuntimeBlocker(gateway)
		account := &Account{ID: 42906, Platform: PlatformOpenAI, Type: AccountTypeOAuth}
		headers := make(http.Header)
		headers.Set("Retry-After", "60")
		deadline := time.Now().Add(time.Minute)
		gateway.handleOpenAIAccountUpstreamError(context.Background(), account, http.StatusTooManyRequests, headers, nil)
		require.True(t, deadline.Equal(repo.rateLimitedAt))
		time.Sleep(6 * time.Second)
		require.True(t, gateway.isOpenAIAccountRuntimeBlocked(account))
	})
}

func (r *oauth429RetryAfterRepo) SetOpenAIOAuth429RateLimited(ctx context.Context, id int64, resetAt, retryAfterUntil time.Time) error {
	if err := r.UpdateExtra(ctx, id, map[string]any{
		OpenAIOAuth429RetryAfterExtraKey: map[string]any{
			"reset_at":          resetAt.UTC().Format(time.RFC3339Nano),
			"retry_after_until": retryAfterUntil.UTC().Format(time.RFC3339Nano),
		},
	}); err != nil {
		return err
	}
	return r.SetRateLimited(ctx, id, resetAt)
}

func TestOpenAIOAuth429RetryAfterSourceCannotProtectUnrelatedOrExpiredLimits(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*Account, map[string]any, time.Time)
	}{
		{name: "missing", mutate: func(a *Account, _ map[string]any, _ time.Time) { delete(a.Extra, OpenAIOAuth429RetryAfterExtraKey) }},
		{name: "different reset", mutate: func(_ *Account, source map[string]any, now time.Time) {
			source["reset_at"] = now.Add(2 * time.Minute).Format(time.RFC3339Nano)
		}},
		{name: "expired retry", mutate: func(_ *Account, source map[string]any, now time.Time) {
			source["retry_after_until"] = now.Add(-time.Second).Format(time.RFC3339Nano)
		}},
		{name: "malformed", mutate: func(_ *Account, source map[string]any, _ time.Time) { source["retry_after_until"] = "invalid" }},
		{name: "retry after effective deadline", mutate: func(_ *Account, source map[string]any, now time.Time) {
			source["retry_after_until"] = now.Add(2 * time.Minute).Format(time.RFC3339Nano)
		}},
		{name: "admin cleared", mutate: func(a *Account, _ map[string]any, _ time.Time) { a.RateLimitResetAt = nil }},
		{name: "keep active", mutate: func(a *Account, _ map[string]any, _ time.Time) { a.Extra[KeepStatusActiveExtraKey] = true }},
		{name: "API key", mutate: func(a *Account, _ map[string]any, _ time.Time) { a.Type = AccountTypeAPIKey }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			now := time.Now()
			resetAt := now.Add(time.Minute)
			source := map[string]any{"reset_at": resetAt.Format(time.RFC3339Nano), "retry_after_until": resetAt.Format(time.RFC3339Nano)}
			account := &Account{
				ID: 42907, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true, RateLimitResetAt: &resetAt,
				Extra: map[string]any{
					OpenAIOAuth429RetryAfterExtraKey: source,
					"codex_5h_used_percent":          30.0, "codex_7d_used_percent": 40.0,
					"codex_5h_reset_at": resetAt.Format(time.RFC3339Nano),
				},
			}
			tc.mutate(account, source, now)
			require.False(t, account.IsRateLimited(), "the legacy quota-window exception still applies without a matching active Retry-After")
		})
	}
}

func TestOpenAIOAuth429RetryAfterMissingAtomicWriterKeepsImmediateBlock(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Older/custom repositories cannot silently replace the atomic writer
		// with separate writes. Persistence failure must leave the fast path safe.
		repo := &openAI429SnapshotRepo{}
		rl := NewRateLimitService(repo, nil, nil, nil, nil)
		gateway := &OpenAIGatewayService{rateLimitService: rl}
		rl.SetAccountRuntimeBlocker(gateway)
		account := &Account{ID: 42908, Platform: PlatformOpenAI, Type: AccountTypeOAuth}
		headers := make(http.Header)
		headers.Set("Retry-After", "60")
		gateway.handleOpenAIAccountUpstreamError(context.Background(), account, http.StatusTooManyRequests, headers, nil)
		require.Zero(t, repo.rateLimitedID, "unsupported atomic persistence must not fall back to unsafe split writes")
		time.Sleep(6 * time.Second)
		require.True(t, gateway.isOpenAIAccountRuntimeBlocked(account))
	})
}
