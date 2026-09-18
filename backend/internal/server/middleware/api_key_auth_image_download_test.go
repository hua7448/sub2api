//go:build unit

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestIsImageDownloadReadRequiresExactMethodAndPath(t *testing.T) {
	require.True(t, isImageDownloadRead(http.MethodPost, "/v1/images/download"))
	require.False(t, isImageDownloadRead(http.MethodGet, "/v1/images/download"))
	require.False(t, isImageDownloadRead(http.MethodPost, "/images/download"))
	require.False(t, isImageDownloadRead(http.MethodPost, "/v1/images/download/extra"))
}

func TestAPIKeyAuthImageDownloadSkipsBillingChecks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	group := &service.Group{
		ID:               42,
		Name:             "subscription",
		Status:           service.StatusActive,
		Hydrated:         true,
		SubscriptionType: service.SubscriptionTypeSubscription,
	}
	user := &service.User{ID: 7, Role: service.RoleUser, Status: service.StatusActive, Balance: 0, Concurrency: 3}
	expiredAt := time.Now().Add(-time.Hour)
	apiKey := &service.APIKey{
		ID: 100, UserID: user.ID, Key: "image-download-auth-only", Status: service.StatusAPIKeyQuotaExhausted,
		User: user, GroupID: &group.ID, Group: group, Quota: 1, QuotaUsed: 1, ExpiresAt: &expiredAt,
	}
	subscriptionCalls := 0
	apiKeyRepo := &stubApiKeyRepo{getByKey: func(context.Context, string) (*service.APIKey, error) {
		clone := *apiKey
		return &clone, nil
	}}
	subscriptionRepo := &stubUserSubscriptionRepo{getActive: func(context.Context, int64, int64) (*service.UserSubscription, error) {
		subscriptionCalls++
		return nil, service.ErrSubscriptionNotFound
	}}
	cfg := &config.Config{RunMode: config.RunModeStandard}
	apiKeyService := service.NewAPIKeyService(apiKeyRepo, nil, nil, nil, nil, nil, cfg)
	subscriptionService := service.NewSubscriptionService(nil, subscriptionRepo, nil, nil, cfg)
	t.Cleanup(subscriptionService.Stop)

	router := gin.New()
	router.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(apiKeyService, subscriptionService, cfg)))
	router.POST("/v1/images/download", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/images/download", nil)
	req.Header.Set("x-api-key", apiKey.Key)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 1, subscriptionCalls)
}
