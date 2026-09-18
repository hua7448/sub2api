//go:build unit

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestApplyUsageBillingEffectsUsesStrictlyMonotonicSubscriptionDatabaseWriteTime(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	mock.ExpectBegin()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)

	previousWriteAt := time.Date(2026, 8, 1, 5, 43, 21, 456789000, time.UTC)
	// Simulate clock_timestamp() landing in the same database microsecond.
	writeAt := previousWriteAt.Add(time.Microsecond)
	mock.ExpectQuery(`(?s)UPDATE user_subscriptions us.*updated_at = GREATEST\(clock_timestamp\(\), us.updated_at \+ interval '1 microsecond'\).*RETURNING us.updated_at`).
		WithArgs(2.5, int64(88)).
		WillReturnRows(sqlmock.NewRows([]string{"updated_at"}).AddRow(writeAt))
	mock.ExpectCommit()

	result := &service.UsageBillingApplyResult{Applied: true}
	err = (&usageBillingRepository{}).applyUsageBillingEffects(ctx, tx, &service.UsageBillingCommand{
		SubscriptionID:   ptrInt64ForUsageWriteAtTest(88),
		SubscriptionCost: 2.5,
	}, result)
	require.NoError(t, err)
	require.NotNil(t, result.SubscriptionUsageWriteAt)
	require.Equal(t, writeAt, *result.SubscriptionUsageWriteAt)
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func ptrInt64ForUsageWriteAtTest(value int64) *int64 {
	return &value
}
