package repository

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	_ "github.com/Wei-Shaw/sub2api/ent/runtime"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

func TestAccountRepository_SetErrorKeepStatusActiveGuardSkipsOutboxOnNoop(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	client := dbent.NewClient(dbent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	t.Cleanup(func() { _ = client.Close() })

	mock.ExpectExec(`(?s)` + regexp.QuoteMeta("UPDATE \"accounts\"") +
		`.*` + regexp.QuoteMeta("NOT (COALESCE(\"accounts\".\"extra\", '{}'::jsonb) @> '{\"keep_status_active\":true}'::jsonb)")).
		WillReturnResult(sqlmock.NewResult(0, 0))

	repo := newAccountRepositoryWithSQL(client, db, nil)
	err = repo.SetError(context.Background(), 42, "automatic upstream error")

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet(), "guarded no-op must not write scheduler outbox")
}

func TestKeepStatusActivePredicatesUseExactJSONBoolean(t *testing.T) {
	var enabledSQL string
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(captureEntQueryMatcher{actual: &enabledSQL}))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	client := dbent.NewClient(dbent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	t.Cleanup(func() { _ = client.Close() })

	mock.ExpectQuery("keep status active predicate").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	_, err = client.Account.Query().
		Where(keepStatusActiveEnabledPredicate()).
		Select("id").
		All(context.Background())
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	require.Contains(t, normalizeSQLWhitespace(enabledSQL),
		`COALESCE("accounts"."extra", '{}'::jsonb) @> '{"keep_status_active":true}'::jsonb`)
	require.NotContains(t, normalizeSQLWhitespace(enabledSQL), "->>", "JSON string values must not enable Keep Active")
}

func TestAccountRepository_BulkKeepStatusActiveAtomicallyClearsSchedulingBarriersAndWritesOutbox(t *testing.T) {
	exec := &recordingSQLExecutor{result: rowsAffectedResult(2)}
	repo := newAccountRepositoryWithSQL(nil, exec, nil)
	manualSchedulable := false

	rows, err := repo.BulkUpdate(context.Background(), []int64{41, 42}, service.AccountBulkUpdate{
		Extra:       map[string]any{service.KeepStatusActiveExtraKey: true},
		Schedulable: &manualSchedulable,
	})

	require.NoError(t, err)
	require.Equal(t, int64(2), rows)
	require.Len(t, exec.execQueries, 2, "the durable reset and one bulk outbox event are required")
	normalized := normalizeSQLWhitespace(exec.execQueries[0])
	require.Contains(t, normalized, "status = $1")
	require.Contains(t, normalized, "error_message = ''")
	require.Contains(t, normalized, "rate_limited_at = NULL")
	require.Contains(t, normalized, "rate_limit_reset_at = NULL")
	require.Contains(t, normalized, "overload_until = NULL")
	require.Contains(t, normalized, "temp_unschedulable_until = NULL")
	require.Contains(t, normalized, "temp_unschedulable_reason = NULL")
	require.Contains(t, normalized, "schedulable = $2", "manual schedulable=false must remain effective")
	require.Contains(t, normalized, "- 'antigravity_quota_scopes' - 'model_rate_limits'")
	require.Equal(t, service.StatusActive, exec.execArgs[0][0])
	require.Equal(t, false, exec.execArgs[0][1])

	var extra map[string]any
	require.NoError(t, json.Unmarshal(exec.execArgs[0][2].([]byte), &extra))
	require.Equal(t, true, extra[service.KeepStatusActiveExtraKey])
	require.Contains(t, strings.ToLower(exec.execQueries[1]), "insert into scheduler_outbox")
	require.Equal(t, service.SchedulerOutboxEventAccountBulkChanged, exec.execArgs[1][0])
}

func TestAccountRepository_BulkKeepStatusActiveFalseDoesNotClearExistingState(t *testing.T) {
	exec := &recordingSQLExecutor{result: rowsAffectedResult(1)}
	repo := newAccountRepositoryWithSQL(nil, exec, nil)

	_, err := repo.BulkUpdate(context.Background(), []int64{41}, service.AccountBulkUpdate{
		Extra: map[string]any{service.KeepStatusActiveExtraKey: false},
	})

	require.NoError(t, err)
	require.Len(t, exec.execQueries, 2)
	normalized := normalizeSQLWhitespace(exec.execQueries[0])
	for _, clause := range []string{
		"error_message = ''",
		"rate_limited_at = NULL",
		"rate_limit_reset_at = NULL",
		"overload_until = NULL",
		"temp_unschedulable_until = NULL",
		"temp_unschedulable_reason = NULL",
	} {
		require.NotContains(t, normalized, clause)
	}
	require.NotContains(t, normalized, "- 'antigravity_quota_scopes'")
	require.NotContains(t, normalized, "- 'model_rate_limits'")
}

func TestAccountRepository_AutomaticHealthWritersGuardKeepStatusActive(t *testing.T) {
	for _, tt := range []struct {
		name   string
		mutate func(context.Context, *accountRepository) error
	}{
		{"rate limit", func(ctx context.Context, r *accountRepository) error {
			return r.SetRateLimited(ctx, 42, time.Now().Add(time.Hour))
		}},
		{"rate limit later", func(ctx context.Context, r *accountRepository) error {
			return r.SetRateLimitedIfLater(ctx, 42, time.Now().Add(time.Hour))
		}},
		{"overload", func(ctx context.Context, r *accountRepository) error {
			return r.SetOverloaded(ctx, 42, time.Now().Add(time.Hour))
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var captured string
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(captureEntQueryMatcher{actual: &captured}))
			require.NoError(t, err)
			t.Cleanup(func() { _ = db.Close() })
			client := dbent.NewClient(dbent.Driver(entsql.OpenDB(dialect.Postgres, db)))
			t.Cleanup(func() { _ = client.Close() })
			mock.ExpectExec("automatic health writer").WillReturnResult(sqlmock.NewResult(0, 0))
			repo := newAccountRepositoryWithSQL(client, db, nil)
			require.NoError(t, tt.mutate(context.Background(), repo))
			require.NoError(t, mock.ExpectationsWereMet())
			require.Contains(t, normalizeSQLWhitespace(captured),
				`NOT (COALESCE("accounts"."extra", '{}'::jsonb) @> '{"keep_status_active":true}'::jsonb)`)
		})
	}
}

func TestAccountRepository_SetModelRateLimitGuardsKeepStatusActiveAndSkipsOutboxOnNoop(t *testing.T) {
	var captured string
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(captureEntQueryMatcher{actual: &captured}))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	client := dbent.NewClient(dbent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	t.Cleanup(func() { _ = client.Close() })
	mock.ExpectExec("model rate limit").WillReturnResult(sqlmock.NewResult(0, 0))
	repo := newAccountRepositoryWithSQL(client, db, nil)
	require.NoError(t, repo.SetModelRateLimit(context.Background(), 42, "model", time.Now().Add(time.Hour)))
	require.NoError(t, mock.ExpectationsWereMet())
	require.Contains(t, normalizeSQLWhitespace(captured),
		"NOT (COALESCE(extra, '{}'::jsonb) @> '{\"keep_status_active\":true}'::jsonb)")
	require.NotContains(t, captured, "scheduler_outbox")
}

func TestSchedulablePredicatesBypassAllAutomaticWindowsForKeepStatusActive(t *testing.T) {
	var captured string
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(captureEntQueryMatcher{actual: &captured}))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	client := dbent.NewClient(dbent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	t.Cleanup(func() { _ = client.Close() })
	mock.ExpectQuery("schedulable query").WillReturnRows(sqlmock.NewRows([]string{"id"}))
	_, err = client.Account.Query().Where(
		tempUnschedulablePredicate(), notExpiredPredicate(time.Now()),
		overloadAvailablePredicate(time.Now()), rateLimitAvailablePredicate(time.Now()),
	).Select("id").All(context.Background())
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	normalized := normalizeSQLWhitespace(captured)
	require.GreaterOrEqual(t, strings.Count(normalized,
		`COALESCE("accounts"."extra", '{}'::jsonb) @> '{"keep_status_active":true}'::jsonb`), 4)
}

func TestKeepStatusActiveSQLLiteralUsesJSONBooleanNotString(t *testing.T) {
	enabled := keepStatusActiveEnabledSQL("extra")
	require.Equal(t, `COALESCE(extra, '{}'::jsonb) @> '{"keep_status_active":true}'::jsonb`, enabled)
	require.NotContains(t, enabled, "->>")
	require.NotContains(t, enabled, `\"true\"`, "the JSON value must be boolean true")
	require.Equal(t, `NOT (COALESCE(extra, '{}'::jsonb) @> '{"keep_status_active":true}'::jsonb)`, "NOT ("+enabled+")")
}
