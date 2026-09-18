//go:build unit

package routes

import (
	"net/http"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/gin-gonic/gin"
)

func TestRegisterUserSubscriptionRoutesIncludesResetDailyQuota(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	v1 := router.Group("/api/v1")
	handlers := &handler.Handlers{
		Subscription: handler.NewSubscriptionHandler(nil),
	}

	registerUserSubscriptionRoutes(v1, handlers, nil)

	for _, route := range router.Routes() {
		if route.Method == http.MethodPost && route.Path == "/api/v1/subscriptions/:id/reset-daily-quota" {
			return
		}
	}
	t.Fatal("POST /api/v1/subscriptions/:id/reset-daily-quota was not registered")
}
