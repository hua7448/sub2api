package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
)

// registerProviderHallUserRoutes mounts the user hall. The display guard
// answers 404 while the hall is switched off; the query endpoints are heavy
// reads and share the stricter per-user limiter.
func registerProviderHallUserRoutes(authenticated *gin.RouterGroup, h *handler.Handlers, panelRateLimiter *middleware.PanelRateLimiter) {
	if h == nil || h.ProviderHall == nil {
		return
	}
	hall := authenticated.Group("/provider-hall")
	hall.Use(panelRateLimiter.Heavy(), h.ProviderHall.DisplayGuard())
	{
		hall.GET("", h.ProviderHall.List)
		hall.GET("/groups/:id", h.ProviderHall.GetGroup)
		hall.GET("/groups/:id/verifications", h.ProviderHall.ListVerifications)
	}
}
