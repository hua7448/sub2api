package repository

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

type capacityQueryCapture struct {
	db       *sql.DB
	captured *string
}

func (c capacityQueryCapture) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c capacityQueryCapture) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	*c.captured = query
	return c.db.QueryContext(ctx, query, args...)
}

func TestListSchedulableCapacityKeepStatusActiveBypassesTempWindow(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	mock.ExpectQuery("SELECT").
		WithArgs(sqlmock.AnyArg(), service.StatusActive, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_id", "account_id", "concurrency", "extra",
			"session_window_start", "session_window_end", "session_window_status",
		}))

	var captured string
	repo := newAccountRepositoryWithSQL(nil, capacityQueryCapture{db: db, captured: &captured}, nil)
	rows, err := repo.ListSchedulableCapacityByGroupIDs(context.Background(), []int64{7})

	require.NoError(t, err)
	require.Empty(t, rows)
	require.NoError(t, mock.ExpectationsWereMet())
	normalized := normalizeSQLWhitespace(captured)
	require.Contains(t, normalized,
		"OR COALESCE(a.extra, '{}'::jsonb) @> '{\"keep_status_active\":true}'::jsonb")
	require.GreaterOrEqual(t, strings.Count(normalized,
		"COALESCE(a.extra, '{}'::jsonb) @> '{\"keep_status_active\":true}'::jsonb"), 4,
		"temp, expiry, overload, and rate-limit windows must all be bypassed")
}
