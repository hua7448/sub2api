package repository

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestSetOpenAIOAuth429RateLimitedAtomicallyBindsDeadlineAndSource(t *testing.T) {
	exec := &recordingSQLExecutor{result: rowsAffectedResult(1)}
	repo := newAccountRepositoryWithSQL(nil, exec, nil)
	resetAt := time.Date(2026, 9, 6, 12, 0, 0, 123456789, time.UTC)
	retryUntil := resetAt.Add(-time.Minute)

	require.NoError(t, repo.SetOpenAIOAuth429RateLimited(context.Background(), 42, resetAt, retryUntil))
	require.Len(t, exec.execQueries, 2, "one account mutation followed by the existing outbox publication")
	query := normalizeSQLWhitespace(exec.execQueries[0])
	require.Contains(t, query, "SET rate_limited_at = $1, rate_limit_reset_at = $2, extra = COALESCE(extra, '{}'::jsonb) || $3::jsonb")
	require.Contains(t, query, "id = $4 AND deleted_at IS NULL")
	require.Contains(t, query, "platform = 'openai' AND type = 'oauth' AND parent_account_id IS NULL")
	require.Contains(t, query, `NOT (COALESCE(extra, '{}'::jsonb) @> '{"keep_status_active":true}'::jsonb)`)
	require.Equal(t, int64(42), exec.execArgs[0][3])
	storedReset, ok := exec.execArgs[0][1].(time.Time)
	require.True(t, ok)
	require.Equal(t, resetAt.Truncate(time.Microsecond), storedReset)
	payload, ok := exec.execArgs[0][2].([]byte)
	require.True(t, ok)
	var extra map[string]any
	require.NoError(t, json.Unmarshal(payload, &extra))
	require.Len(t, extra, 1, "only Retry-After provenance is merged; credentials and unrelated extra are untouched")
	source := extra[service.OpenAIOAuth429RetryAfterExtraKey].(map[string]any)
	require.Equal(t, storedReset.Format(time.RFC3339Nano), source["reset_at"])
	require.Equal(t, retryUntil.Truncate(time.Microsecond).Format(time.RFC3339Nano), source["retry_after_until"])
	require.Contains(t, exec.execQueries[1], "INSERT INTO scheduler_outbox")
	require.Equal(t, service.SchedulerOutboxEventAccountChanged, exec.execArgs[1][0])
}

func TestSetOpenAIOAuth429RateLimitedNoopAndFailureDoNotPublish(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
	}{
		{name: "guarded or missing account"},
		{name: "account update failed", err: errors.New("write failed")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			exec := &recordingSQLExecutor{result: rowsAffectedResult(0), err: tc.err}
			repo := newAccountRepositoryWithSQL(nil, exec, nil)
			deadline := time.Now().Add(time.Minute)
			err := repo.SetOpenAIOAuth429RateLimited(context.Background(), 42, deadline, deadline)
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.err)
			}
			require.Len(t, exec.execQueries, 1, "no separate metadata mutation and no outbox on a failed or guarded write")
		})
	}
}

func TestSetOpenAIOAuth429RateLimitedRejectsInvalidDeadlineWithoutWrites(t *testing.T) {
	exec := &recordingSQLExecutor{result: rowsAffectedResult(1)}
	repo := newAccountRepositoryWithSQL(nil, exec, nil)
	deadline := time.Now().Add(time.Minute)
	require.Error(t, repo.SetOpenAIOAuth429RateLimited(context.Background(), 42, deadline, deadline.Add(time.Second)))
	require.Error(t, repo.SetOpenAIOAuth429RateLimited(context.Background(), 42, deadline, time.Time{}))
	require.Error(t, repo.SetOpenAIOAuth429RateLimited(context.Background(), 0, deadline, deadline))
	require.Empty(t, exec.execQueries)
}
