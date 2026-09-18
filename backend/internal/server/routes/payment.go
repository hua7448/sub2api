package routes

import (
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/handler/admin"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterPaymentRoutes registers all payment-related routes:
// user-facing endpoints, webhook endpoints, and admin endpoints.
func RegisterPaymentRoutes(
	v1 *gin.RouterGroup,
	paymentHandler *handler.PaymentHandler,
	webhookHandler *handler.PaymentWebhookHandler,
	adminPaymentHandler *admin.PaymentHandler,
	jwtAuth middleware.JWTAuthMiddleware,
	adminAuth middleware.AdminAuthMiddleware,
	auditLog middleware.AuditLogMiddleware,
	settingService *service.SettingService,
	panelRateLimiter *middleware.PanelRateLimiter,
) {
	disabled := func(c *gin.Context) {
		response.ErrorWithDetails(c, http.StatusServiceUnavailable, "Payment functionality is temporarily unavailable.", "PAYMENT_DISABLED", nil)
	}

	// --- User-facing payment endpoints (authenticated) ---
	authenticated := v1.Group("/payment")
	authenticated.Use(gin.HandlerFunc(jwtAuth))
	authenticated.Use(middleware.BackendModeUserGuard(settingService))
	// 面板全局按用户限流
	authenticated.Use(panelRateLimiter.Global())
	{
		authenticated.GET("/config", disabled)
		authenticated.GET("/checkout-info", disabled)
		authenticated.GET("/plans", disabled)
		authenticated.GET("/limits", disabled)

		orders := authenticated.Group("/orders")
		{
			orders.POST("", disabled)
			orders.POST("/verify", disabled)
			orders.GET("/my", disabled)
			orders.GET("/:id", disabled)
			orders.POST("/:id/cancel", disabled)
			orders.POST("/:id/refund-request", disabled)
			orders.GET("/refund-eligible-providers", disabled)
		}
	}

	// --- Public payment endpoints (no auth) ---
	// Signed resume-token recovery is the preferred public lookup path.
	// The legacy anonymous out_trade_no verify endpoint remains available as a
	// persisted-state compatibility path for staggered upgrades.
	public := v1.Group("/payment/public")
	{
		public.POST("/orders/verify", disabled)
		public.POST("/orders/resolve", disabled)
	}

	// --- Webhook endpoints (no auth) ---
	webhook := v1.Group("/payment/webhook")
	{
		// EasyPay sends GET callbacks with query params
		webhook.GET("/easypay", disabled)
		webhook.POST("/easypay", disabled)
		webhook.POST("/alipay", disabled)
		webhook.POST("/wxpay", disabled)
		webhook.POST("/stripe", disabled)
		webhook.POST("/airwallex", disabled)
	}

	// --- Admin payment endpoints (admin auth) ---
	adminGroup := v1.Group("/admin/payment")
	adminGroup.Use(gin.HandlerFunc(adminAuth))
	adminGroup.Use(gin.HandlerFunc(auditLog))
	adminGroup.Use(middleware.AdminComplianceGuard(settingService))
	{
		// Dashboard
		adminGroup.GET("/dashboard", disabled)

		// Config
		adminGroup.GET("/config", disabled)
		adminGroup.PUT("/config", disabled)

		// Orders
		adminOrders := adminGroup.Group("/orders")
		{
			adminOrders.GET("", disabled)
			adminOrders.GET("/:id", disabled)
			adminOrders.POST("/:id/cancel", disabled)
			adminOrders.POST("/:id/retry", disabled)
			adminOrders.POST("/:id/refund", disabled)
			adminOrders.POST("/:id/refund/query", disabled)
		}

		// Subscription Plans
		plans := adminGroup.Group("/plans")
		{
			plans.GET("", disabled)
			plans.POST("", disabled)
			plans.PUT("/:id", disabled)
			plans.DELETE("/:id", disabled)
		}

		// Provider Instances
		providers := adminGroup.Group("/providers")
		{
			providers.GET("", disabled)
			providers.POST("", disabled)
			providers.PUT("/:id", disabled)
			providers.DELETE("/:id", disabled)
		}
	}
}
