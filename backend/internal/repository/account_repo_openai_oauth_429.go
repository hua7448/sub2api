package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

var _ service.OpenAIOAuth429RateLimitWriter = (*accountRepository)(nil)

// SetOpenAIOAuth429RateLimited binds Retry-After provenance to the exact persisted
// limit in one row update. A separate UpdateExtra followed by SetRateLimited can
// interleave with another response and incorrectly label its cooldown deadline.
func (r *accountRepository) SetOpenAIOAuth429RateLimited(ctx context.Context, id int64, resetAt, retryAfterUntil time.Time) error {
	if id <= 0 || resetAt.IsZero() || retryAfterUntil.IsZero() || retryAfterUntil.After(resetAt) {
		return fmt.Errorf("invalid OpenAI OAuth 429 deadline")
	}
	resetAt = resetAt.UTC().Truncate(time.Microsecond)
	retryAfterUntil = retryAfterUntil.UTC().Truncate(time.Microsecond)
	payload, err := json.Marshal(map[string]any{
		service.OpenAIOAuth429RetryAfterExtraKey: map[string]string{
			"reset_at":          resetAt.UTC().Format(time.RFC3339Nano),
			"retry_after_until": retryAfterUntil.UTC().Format(time.RFC3339Nano),
		},
	})
	if err != nil {
		return err
	}
	result, err := r.sql.ExecContext(ctx, `UPDATE accounts
		SET rate_limited_at = $1, rate_limit_reset_at = $2,
			extra = COALESCE(extra, '{}'::jsonb) || $3::jsonb, updated_at = $1
		WHERE id = $4 AND deleted_at IS NULL
			AND platform = 'openai' AND type = 'oauth' AND parent_account_id IS NULL
			AND NOT (COALESCE(extra, '{}'::jsonb) @> '{"keep_status_active":true}'::jsonb)`,
		time.Now(), resetAt, payload, id)
	if err != nil {
		return err
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if updated == 0 {
		return nil
	}
	// Preserve the existing SetRateLimited publication and snapshot behavior.
	if err := enqueueSchedulerOutbox(ctx, r.sql, service.SchedulerOutboxEventAccountChanged, &id, nil, nil); err != nil {
		logger.LegacyPrintf("repository.account", "[SchedulerOutbox] enqueue OAuth 429 rate limit failed: account=%d err=%v", id, err)
	}
	r.syncSchedulerAccountSnapshot(ctx, id)
	return nil
}
