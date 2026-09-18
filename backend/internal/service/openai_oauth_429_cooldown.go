package service

import (
	"context"
	"net/http"
	"time"
)

// OpenAIOAuth429RateLimitWriter persists an explicit request-rate deadline and
// its provenance together. The production account repository implements this
// narrow capability; splitting these writes could mislabel a concurrent 429.
type OpenAIOAuth429RateLimitWriter interface {
	SetOpenAIOAuth429RateLimited(ctx context.Context, accountID int64, resetAt, retryAfterUntil time.Time) error
}

const OpenAIOAuth429RetryAfterExtraKey = "openai_oauth_429_retry_after"

type openAIOAuth429CooldownContextKey struct{}

type openAIOAuth429CooldownDecision struct {
	accountID    int64
	resetAt      *time.Time
	retryAfterAt *time.Time
}

// Resolve once per response so persistence and the immediate scheduling gate
// agree even when parsing relative headers or when writing account state is slow.
func openAIOAuth429CooldownDecisionFromContext(ctx context.Context, account *Account, headers http.Header, body []byte) openAIOAuth429CooldownDecision {
	if decision, ok := ctx.Value(openAIOAuth429CooldownContextKey{}).(openAIOAuth429CooldownDecision); ok && decision.accountID == account.ID {
		return decision
	}
	return resolveOpenAIOAuth429CooldownDecision(account.ID, headers, body, time.Now())
}

func resolveOpenAIOAuth429CooldownDecision(accountID int64, headers http.Header, body []byte, now time.Time) openAIOAuth429CooldownDecision {
	decision := openAIOAuth429CooldownDecision{accountID: accountID}
	consider := func(resetAt *time.Time) {
		if resetAt == nil {
			return
		}
		// Match PostgreSQL timestamp precision in both the immediate block and
		// persisted source metadata, without depending on driver rounding.
		normalized := resetAt.UTC().Truncate(time.Microsecond)
		if normalized.After(now) && (decision.resetAt == nil || normalized.After(*decision.resetAt)) {
			decision.resetAt = &normalized
		}
	}

	codexResetAt, hasCodexHeaders := calculateOpenAI429ResetTimeFromHeadersAt(headers, now)
	consider(codexResetAt)
	if resetUnix := parseOpenAIRateLimitResetTimeAt(body, now); resetUnix != nil {
		resetAt := time.Unix(*resetUnix, 0)
		// Some upstreams echo a quota-window reset in the body while both
		// Codex windows still have quota. Preserve the existing suppression;
		// Retry-After remains independent evidence of a request rate limit.
		if !hasCodexHeaders || !openAICodexHeadersBelowLimitResetMatches(headers, resetAt, now) {
			consider(&resetAt)
		}
	}
	if retryAfterAt := parseRetryAfterResetTime(headers, now); retryAfterAt != nil && retryAfterAt.After(now) {
		normalized := retryAfterAt.UTC().Truncate(time.Microsecond)
		if normalized.After(now) {
			decision.retryAfterAt = &normalized
			consider(&normalized)
		}
	}
	return decision
}

// Match the persisted generation before overriding legacy quota-window recovery.
// PostgreSQL timestamps have microsecond precision; expired or mismatched source
// metadata must not protect a later, unrelated limit or undo an admin clear.
func hasActiveOpenAIOAuth429RetryAfter(account *Account, now time.Time) bool {
	if account == nil || !account.IsOpenAIOAuth() || account.IsShadow() || account.RateLimitResetAt == nil || !account.RateLimitResetAt.After(now) {
		return false
	}
	source, ok := account.Extra[OpenAIOAuth429RetryAfterExtraKey].(map[string]any)
	if !ok {
		return false
	}
	resetRaw, resetOK := source["reset_at"].(string)
	retryRaw, retryOK := source["retry_after_until"].(string)
	if !resetOK || !retryOK {
		return false
	}
	resetAt, resetErr := time.Parse(time.RFC3339Nano, resetRaw)
	retryAfterAt, retryErr := time.Parse(time.RFC3339Nano, retryRaw)
	if resetErr != nil || retryErr != nil || !retryAfterAt.After(now) || retryAfterAt.After(resetAt) {
		return false
	}
	return account.RateLimitResetAt.Truncate(time.Microsecond).Equal(resetAt.Truncate(time.Microsecond))
}
