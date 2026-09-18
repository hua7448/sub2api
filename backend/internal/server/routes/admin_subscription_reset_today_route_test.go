//go:build unit

package routes

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Wei-Shaw/sub2api/internal/handler"
	adminhandler "github.com/Wei-Shaw/sub2api/internal/handler/admin"
)

func TestRegisterSubscriptionRoutesIncludesResetTodayUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	adminGroup := router.Group("/api/v1/admin")
	handlers := &handler.Handlers{Admin: &handler.AdminHandlers{
		Subscription: adminhandler.NewSubscriptionHandler(nil),
	}}

	registerSubscriptionRoutes(adminGroup, handlers)

	for _, route := range router.Routes() {
		if route.Method == http.MethodPost && route.Path == "/api/v1/admin/subscriptions/reset-today-usage" {
			return
		}
	}
	t.Fatal("POST /api/v1/admin/subscriptions/reset-today-usage was not registered")
}

func TestRegisterSubscriptionRoutesIncludesResetQuota(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	adminGroup := router.Group("/api/v1/admin")
	handlers := &handler.Handlers{Admin: &handler.AdminHandlers{
		Subscription: adminhandler.NewSubscriptionHandler(nil),
	}}

	registerSubscriptionRoutes(adminGroup, handlers)

	for _, route := range router.Routes() {
		if route.Method == http.MethodPost && route.Path == "/api/v1/admin/subscriptions/:id/reset-quota" {
			return
		}
	}
	t.Fatal("POST /api/v1/admin/subscriptions/:id/reset-quota was not registered")
}
