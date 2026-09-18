package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/gin-gonic/gin"
)

func registerProviderHallRoutes(admin *gin.RouterGroup, h *handler.Handlers) {
	hall := admin.Group("/provider-hall")
	hall.GET("/config", h.Admin.ProviderHall.GetConfig)
	hall.PUT("/config", h.Admin.ProviderHall.UpdateConfig)
	hall.POST("/config/preflight", h.Admin.ProviderHall.PreflightConfig)
	hall.POST("/config/check-gateway", h.Admin.ProviderHall.CheckGateway)
	hall.GET("/groups", h.Admin.ProviderHall.ListAdminGroups)
	hall.POST("/groups/batch", h.Admin.ProviderHall.BatchGroups)
	hall.GET("/groups/:id/models", h.Admin.ProviderHall.ListGroupModels)
	hall.POST("/groups/:id/models/refresh", h.Admin.ProviderHall.RefreshGroupModels)
	hall.GET("/groups/:id/probe-keys", h.Admin.ProviderHall.ListProbeKeys)
	hall.POST("/groups/:id/probe-keys", h.Admin.ProviderHall.EnsureProbeKey)
	hall.PUT("/groups/:id/settings", h.Admin.ProviderHall.SaveSettings)
	hall.POST("/groups/:id/preflight", h.Admin.ProviderHall.PreflightGroup)
	hall.GET("/groups/:id", h.Admin.ProviderHall.GetGroup)
	hall.PUT("/groups/:id", h.Admin.ProviderHall.UpdateGroup)
	hall.GET("/groups/:id/targets", h.Admin.ProviderHall.GetTargets)
	hall.PUT("/groups/:id/targets", h.Admin.ProviderHall.UpdateTargets)
	hall.GET("/profiles", h.Admin.ProviderHall.ListProfiles)
	hall.POST("/profiles", h.Admin.ProviderHall.CreateProfile)
	hall.PUT("/profiles/:id", h.Admin.ProviderHall.UpdateProfile)
	hall.DELETE("/profiles/:id", h.Admin.ProviderHall.DeleteProfile)
	hall.GET("/profile-candidates", h.Admin.ProviderHall.ListAllProfileCandidates)
	// Task and operations endpoints (batch B6).
	hall.POST("/groups/:id/probes", h.Admin.ProviderHall.EnqueueProbe)
	hall.POST("/groups/:id/verifications", h.Admin.ProviderHall.EnqueueVerification)
	hall.GET("/groups/:id/verifications", h.Admin.ProviderHall.ListGroupVerifications)
	hall.GET("/jobs", h.Admin.ProviderHall.ListJobs)
	hall.GET("/jobs/:id", h.Admin.ProviderHall.GetJob)
	hall.POST("/jobs/:id/cancel", h.Admin.ProviderHall.CancelJob)
	hall.GET("/health", h.Admin.ProviderHall.Health)
}
