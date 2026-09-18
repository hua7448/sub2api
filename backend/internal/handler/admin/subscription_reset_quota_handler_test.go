//go:build unit

package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type resetQuotaSelectorRepo struct {
	service.UserSubscriptionRepository

	sub                     *service.UserSubscription
	resetDaily, resetWeekly bool
	resetMonthly            bool
}

func (r *resetQuotaSelectorRepo) GetByID(_ context.Context, id int64) (*service.UserSubscription, error) {
	if r.sub == nil || r.sub.ID != id {
		return nil, service.ErrSubscriptionNotFound
	}
	copy := *r.sub
	return &copy, nil
}

func (r *resetQuotaSelectorRepo) ResetUsageWindows(_ context.Context, _ int64, resetDaily, resetWeekly, resetMonthly bool, dailyStart, periodicStart time.Time) (time.Time, error) {
	r.resetDaily = resetDaily
	r.resetWeekly = resetWeekly
	r.resetMonthly = resetMonthly
	if resetDaily {
		r.sub.DailyUsageUSD = 0
		r.sub.DailyWindowStart = &dailyStart
	}
	if resetWeekly {
		r.sub.WeeklyUsageUSD = 0
		r.sub.WeeklyWindowStart = &periodicStart
	}
	if resetMonthly {
		r.sub.MonthlyUsageUSD = 0
		r.sub.MonthlyWindowStart = &periodicStart
	}
	return periodicStart, nil
}

func TestResetSubscriptionQuotaRequestDecodesDailyWeeklyMonthlySelectors(t *testing.T) {
	var request ResetSubscriptionQuotaRequest
	require.NoError(t, json.Unmarshal([]byte(`{"daily":true,"weekly":true,"monthly":true}`), &request))
	require.True(t, request.Daily)
	require.True(t, request.Weekly)
	require.True(t, request.Monthly)
}

func TestResetQuotaForwardsDailyWeeklyMonthlySelectors(t *testing.T) {
	tests := []struct {
		name                               string
		body                               string
		wantDaily, wantWeekly, wantMonthly bool
	}{
		{name: "daily", body: `{"daily":true}`, wantDaily: true},
		{name: "weekly", body: `{"weekly":true}`, wantWeekly: true},
		{name: "monthly", body: `{"monthly":true}`, wantMonthly: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &resetQuotaSelectorRepo{sub: &service.UserSubscription{
				ID:              42,
				UserID:          7,
				GroupID:         8,
				StartsAt:        time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC),
				ExpiresAt:       time.Date(2090, 8, 1, 12, 0, 0, 0, time.UTC),
				Status:          service.SubscriptionStatusActive,
				DailyUsageUSD:   1,
				WeeklyUsageUSD:  2,
				MonthlyUsageUSD: 3,
			}}
			subscriptionService := service.NewSubscriptionService(nil, repo, nil, nil, nil)
			t.Cleanup(subscriptionService.Stop)
			handler := NewSubscriptionHandler(subscriptionService)
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/api/v1/admin/subscriptions/:id/reset-quota", handler.ResetQuota)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/subscriptions/42/reset-quota", bytes.NewBufferString(test.body))
			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
			require.Equal(t, test.wantDaily, repo.resetDaily)
			require.Equal(t, test.wantWeekly, repo.resetWeekly)
			require.Equal(t, test.wantMonthly, repo.resetMonthly)
		})
	}
}

func TestResetQuotaRequiresAtLeastOneSelectedWindow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewSubscriptionHandler(nil)
	router := gin.New()
	router.POST("/api/v1/admin/subscriptions/:id/reset-quota", handler.ResetQuota)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/subscriptions/42/reset-quota", bytes.NewBufferString(`{}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "At least one of 'daily', 'weekly', or 'monthly' must be true")
}
