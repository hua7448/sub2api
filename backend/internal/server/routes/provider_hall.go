package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/gin-gonic/gin"
)

func registerProviderHallRoutes(admin *gin.RouterGroup, h *handler.Handlers) {
	hall := admin.Group("/provider-hall")
	hall.GET("/config", h.Admin.ProviderHall.GetConfig)
	hall.PUT("/config", h.Admin.ProviderHall.UpdateConfig)
	hall.GET("/groups/:id", h.Admin.ProviderHall.GetGroup)
	hall.PUT("/groups/:id", h.Admin.ProviderHall.UpdateGroup)
	hall.GET("/groups/:id/targets", h.Admin.ProviderHall.GetTargets)
	hall.PUT("/groups/:id/targets", h.Admin.ProviderHall.UpdateTargets)
	hall.GET("/profiles", h.Admin.ProviderHall.ListProfiles)
	hall.POST("/profiles", h.Admin.ProviderHall.CreateProfile)
	hall.PUT("/profiles/:id", h.Admin.ProviderHall.UpdateProfile)
	// Task and operations endpoints (batch B6).
	hall.POST("/groups/:id/probes", h.Admin.ProviderHall.EnqueueProbe)
	hall.POST("/groups/:id/verifications", h.Admin.ProviderHall.EnqueueVerification)
	hall.GET("/groups/:id/verifications", h.Admin.ProviderHall.ListGroupVerifications)
	hall.GET("/jobs", h.Admin.ProviderHall.ListJobs)
	hall.GET("/jobs/:id", h.Admin.ProviderHall.GetJob)
	hall.POST("/jobs/:id/cancel", h.Admin.ProviderHall.CancelJob)
	hall.GET("/health", h.Admin.ProviderHall.Health)
}
