package repository

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

type groupCountQueryCapture struct {
	db       *sql.DB
	captured *string
}

func (c groupCountQueryCapture) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c groupCountQueryCapture) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	*c.captured = query
	return c.db.QueryContext(ctx, query, args...)
}

func TestGroupAccountCountsClassifyKeepActiveTempWindowAsAvailable(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	mock.ExpectQuery("SELECT").
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"total", "active", "rate_limited"}).AddRow(1, 1, 0))

	var captured string
	repo := newGroupRepositoryWithSQL(nil, groupCountQueryCapture{db: db, captured: &captured})
	total, active, err := repo.GetAccountCount(context.Background(), 7)

	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, active)
	require.NoError(t, mock.ExpectationsWereMet())
	normalized := normalizeSQLWhitespace(captured)
	require.Contains(t, normalized,
		"OR COALESCE(a.extra, '{}'::jsonb) @> '{\"keep_status_active\":true}'::jsonb")
	require.Contains(t, normalized,
		"NOT (COALESCE(a.extra, '{}'::jsonb) @> '{\"keep_status_active\":true}'::jsonb)")
	require.GreaterOrEqual(t, strings.Count(normalized,
		"COALESCE(a.extra, '{}'::jsonb) @> '{\"keep_status_active\":true}'::jsonb"), 4,
		"available count must bypass expiry, rate-limit, overload, and temp windows")
}
