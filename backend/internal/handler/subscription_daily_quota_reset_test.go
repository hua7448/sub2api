//go:build unit

package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/repository"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

func newDailyQuotaResetHandlerSQLMock(t *testing.T) (*SubscriptionHandler, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	client := dbent.NewClient(dbent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	t.Cleanup(func() { _ = client.Close() })
	repo := repository.NewUserSubscriptionRepository(client)
	svc := service.NewSubscriptionService(nil, repo, nil, client, nil)
	return NewSubscriptionHandler(svc), mock
}

func dailyQuotaResetRouter(h *SubscriptionHandler, userID int64) *gin.Engine {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: userID})
		c.Next()
	})
	router.POST("/api/v1/subscriptions/:id/reset-daily-quota", h.ResetDailyQuota)
	return router
}

func dailyQuotaResetRequest(subscriptionID, key, body string) *http.Request {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/subscriptions/"+subscriptionID+"/reset-daily-quota",
		bytes.NewBufferString(body),
	)
	req.Header.Set("Content-Type", "application/json")
	if key != "" {
		req.Header.Set("Idempotency-Key", key)
	}
	return req
}

func TestResetDailyQuotaRequiresConfirmationAndIdempotencyKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := dailyQuotaResetRouter(NewSubscriptionHandler(nil), 77)

	missingConfirm := httptest.NewRecorder()
	router.ServeHTTP(missingConfirm, dailyQuotaResetRequest("501", "daily-reset-501", `{"confirm":false}`))
	require.Equal(t, http.StatusBadRequest, missingConfirm.Code)
	require.Contains(t, missingConfirm.Body.String(), "confirm must be true")

	missingKey := httptest.NewRecorder()
	router.ServeHTTP(missingKey, dailyQuotaResetRequest("501", "", `{"confirm":true}`))
	require.Equal(t, http.StatusBadRequest, missingKey.Code)
	require.Contains(t, missingKey.Body.String(), "IDEMPOTENCY_KEY_REQUIRED")
}

func TestResetDailyQuotaSetsFixedAuditActionBeforeValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewSubscriptionHandler(nil)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 77})
	ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}
	ctx.Request = dailyQuotaResetRequest("invalid", "daily-reset-invalid", `{"confirm":true}`)

	h.ResetDailyQuota(ctx)

	action, exists := ctx.Get("audit_action")
	require.True(t, exists)
	require.Equal(t, service.AuditActionUserSubscriptionResetDailyQuota, action)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestResetDailyQuotaUsesJWTUserAndReplaysIdempotentResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, mock := newDailyQuotaResetHandlerSQLMock(t)
	router := dailyQuotaResetRouter(h, 77)

	previousCoordinator := service.DefaultIdempotencyCoordinator()
	cfg := service.DefaultIdempotencyConfig()
	cfg.ObserveOnly = false
	service.SetDefaultIdempotencyCoordinator(service.NewIdempotencyCoordinator(newUserMemoryIdempotencyRepoStub(), cfg))
	t.Cleanup(func() { service.SetDefaultIdempotencyCoordinator(previousCoordinator) })

	resetAt := time.Date(2026, 8, 5, 4, 50, 0, 123000000, time.UTC)
	availableAt := resetAt.Add(7 * 24 * time.Hour)
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*owned_subscription.daily_limit_usd > 0.*UPDATE user_subscriptions AS us.*RETURNING`).
		WithArgs(int64(501), int64(77), service.SubscriptionStatusActive, service.HashIdempotencyKey("user:77\x00daily-reset-501"), int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second)).
		WillReturnRows(sqlmock.NewRows([]string{
			"did_reset", "operation_replayed", "id", "group_id", "updated_at", "daily_usage_usd", "weekly_usage_usd", "monthly_usage_usd", "available_at",
		}).AddRow(true, false, int64(501), int64(19), resetAt, 0.0, 12.5, 31.75, availableAt))

	first := httptest.NewRecorder()
	router.ServeHTTP(first, dailyQuotaResetRequest("501", "daily-reset-501", `{"confirm":true}`))
	require.Equal(t, http.StatusOK, first.Code)
	require.Contains(t, first.Body.String(), `"subscription_id":501`)
	require.Contains(t, first.Body.String(), `"daily_usage_usd":0`)
	require.Contains(t, first.Body.String(), `"weekly_usage_usd":12.5`)
	require.NotContains(t, first.Body.String(), `"group_id"`)
	require.Contains(t, first.Body.String(), `"daily_quota_reset_available":false`)
	require.Contains(t, first.Body.String(), `"daily_quota_reset_available_at"`)

	replay := httptest.NewRecorder()
	router.ServeHTTP(replay, dailyQuotaResetRequest("501", "daily-reset-501", `{"confirm":true}`))
	require.Equal(t, http.StatusOK, replay.Code)
	require.Equal(t, "true", replay.Header().Get("X-Idempotency-Replayed"))

	conflict := httptest.NewRecorder()
	router.ServeHTTP(conflict, dailyQuotaResetRequest("502", "daily-reset-501", `{"confirm":true}`))
	require.Equal(t, http.StatusConflict, conflict.Code)
	require.Contains(t, conflict.Body.String(), "IDEMPOTENCY_KEY_CONFLICT")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetDailyQuotaMarksDurableDatabaseReplay(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, mock := newDailyQuotaResetHandlerSQLMock(t)
	router := dailyQuotaResetRouter(h, 77)

	previousCoordinator := service.DefaultIdempotencyCoordinator()
	cfg := service.DefaultIdempotencyConfig()
	cfg.ObserveOnly = false
	service.SetDefaultIdempotencyCoordinator(service.NewIdempotencyCoordinator(newUserMemoryIdempotencyRepoStub(), cfg))
	t.Cleanup(func() { service.SetDefaultIdempotencyCoordinator(previousCoordinator) })

	resetAt := time.Date(2026, 8, 5, 4, 50, 0, 123000000, time.UTC)
	availableAt := resetAt.Add(7 * 24 * time.Hour)
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*operation_replay AS MATERIALIZED.*UPDATE user_subscriptions AS us.*RETURNING`).
		WithArgs(
			int64(501),
			int64(77),
			service.SubscriptionStatusActive,
			service.HashIdempotencyKey("user:77\x00recovered-reset-501"),
			int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second),
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"did_reset", "operation_replayed", "id", "group_id", "updated_at", "daily_usage_usd", "weekly_usage_usd", "monthly_usage_usd", "available_at",
		}).AddRow(true, true, int64(501), int64(19), resetAt, 0.0, 12.5, 31.75, availableAt))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, dailyQuotaResetRequest("501", "recovered-reset-501", `{"confirm":true}`))

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "true", response.Header().Get("X-Idempotency-Replayed"))
	require.Contains(t, response.Body.String(), `"reset_at":"`+resetAt.Format(time.RFC3339Nano)+`"`)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetDailyQuotaIdempotencyKeyIsIsolatedPerUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, mock := newDailyQuotaResetHandlerSQLMock(t)
	firstUserRouter := dailyQuotaResetRouter(h, 77)
	secondUserRouter := dailyQuotaResetRouter(h, 88)

	previousCoordinator := service.DefaultIdempotencyCoordinator()
	cfg := service.DefaultIdempotencyConfig()
	cfg.ObserveOnly = false
	service.SetDefaultIdempotencyCoordinator(service.NewIdempotencyCoordinator(newUserMemoryIdempotencyRepoStub(), cfg))
	t.Cleanup(func() { service.SetDefaultIdempotencyCoordinator(previousCoordinator) })

	resetAt := time.Date(2026, 8, 5, 4, 50, 0, 0, time.UTC)
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*UPDATE user_subscriptions AS us.*RETURNING`).
		WithArgs(int64(501), int64(77), service.SubscriptionStatusActive, service.HashIdempotencyKey("user:77\x00shared-client-generated-key"), int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second)).
		WillReturnRows(sqlmock.NewRows([]string{
			"did_reset", "operation_replayed", "id", "group_id", "updated_at", "daily_usage_usd", "weekly_usage_usd", "monthly_usage_usd", "available_at",
		}).AddRow(true, false, int64(501), int64(19), resetAt, 0.0, 12.5, 31.75, resetAt.Add(7*24*time.Hour)))
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*UPDATE user_subscriptions AS us.*RETURNING`).
		WithArgs(int64(502), int64(88), service.SubscriptionStatusActive, service.HashIdempotencyKey("user:88\x00shared-client-generated-key"), int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second)).
		WillReturnRows(sqlmock.NewRows([]string{
			"did_reset", "operation_replayed", "id", "group_id", "updated_at", "daily_usage_usd", "weekly_usage_usd", "monthly_usage_usd", "available_at",
		}).AddRow(true, false, int64(502), int64(20), resetAt, 0.0, 22.5, 41.75, resetAt.Add(7*24*time.Hour)))

	const sharedExternalKey = "shared-client-generated-key"
	first := httptest.NewRecorder()
	firstUserRouter.ServeHTTP(
		first,
		dailyQuotaResetRequest("501", sharedExternalKey, `{"confirm":true}`),
	)
	require.Equal(t, http.StatusOK, first.Code)
	require.Empty(t, first.Header().Get("X-Idempotency-Replayed"))

	second := httptest.NewRecorder()
	secondUserRouter.ServeHTTP(
		second,
		dailyQuotaResetRequest("502", sharedExternalKey, `{"confirm":true}`),
	)
	require.Equal(t, http.StatusOK, second.Code)
	require.Empty(t, second.Header().Get("X-Idempotency-Replayed"))
	require.Contains(t, second.Body.String(), `"subscription_id":502`)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetDailyQuotaReturnsWeeklyLimitForNewKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, mock := newDailyQuotaResetHandlerSQLMock(t)
	router := dailyQuotaResetRouter(h, 77)

	previousCoordinator := service.DefaultIdempotencyCoordinator()
	cfg := service.DefaultIdempotencyConfig()
	cfg.ObserveOnly = false
	service.SetDefaultIdempotencyCoordinator(service.NewIdempotencyCoordinator(newUserMemoryIdempotencyRepoStub(), cfg))
	t.Cleanup(func() { service.SetDefaultIdempotencyCoordinator(previousCoordinator) })

	resetAt := time.Date(2026, 8, 5, 4, 50, 0, 0, time.UTC)
	availableAt := resetAt.Add(7 * 24 * time.Hour)
	mock.ExpectQuery(`(?s)WITH reset_clock AS MATERIALIZED.*UPDATE user_subscriptions AS us.*RETURNING`).
		WithArgs(int64(501), int64(77), service.SubscriptionStatusActive, service.HashIdempotencyKey("user:77\x00new-key-this-week"), int64(service.UserDailyQuotaResetIdempotencyTTL/time.Second)).
		WillReturnRows(sqlmock.NewRows([]string{
			"did_reset", "operation_replayed", "id", "group_id", "updated_at", "daily_usage_usd", "weekly_usage_usd", "monthly_usage_usd", "available_at",
		}).AddRow(false, false, int64(501), int64(19), resetAt, 2.5, 12.5, 31.75, availableAt))

	response := httptest.NewRecorder()
	router.ServeHTTP(response, dailyQuotaResetRequest("501", "new-key-this-week", `{"confirm":true}`))

	require.Equal(t, http.StatusConflict, response.Code)
	require.Contains(t, response.Body.String(), `"reason":"DAILY_QUOTA_RESET_WEEKLY_LIMIT"`)
	require.Contains(t, response.Body.String(), `"available_at":"`+availableAt.Format(time.RFC3339Nano)+`"`)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetDailyQuotaIdempotencyTTLSpansWeeklyWindow(t *testing.T) {
	require.Equal(t, 8*24*time.Hour, userDailyQuotaResetIdempotencyTTL)
}
