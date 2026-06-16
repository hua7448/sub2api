package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type ImageGalleryHandler struct {
	service             *service.ImageGalleryService
	openAIGateway       *OpenAIGatewayHandler
	subscriptionService *service.SubscriptionService
}

func NewImageGalleryHandler(service *service.ImageGalleryService, openAIGateway *OpenAIGatewayHandler, subscriptionService *service.SubscriptionService) *ImageGalleryHandler {
	return &ImageGalleryHandler{service: service, openAIGateway: openAIGateway, subscriptionService: subscriptionService}
}

func (h *ImageGalleryHandler) Settings(c *gin.Context) {
	settings, err := h.service.PublicSettings(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, settings)
}

func (h *ImageGalleryHandler) EligibleKeys(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	keys, err := h.service.EligibleKeys(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, keys)
}

func (h *ImageGalleryHandler) ProxyImagesGenerations(c *gin.Context) {
	h.proxyOpenAI(c, "/v1/images/generations", h.openAIGateway.Images)
}

func (h *ImageGalleryHandler) ProxyImagesEdits(c *gin.Context) {
	h.proxyOpenAI(c, "/v1/images/edits", h.openAIGateway.Images)
}

func (h *ImageGalleryHandler) ProxyResponses(c *gin.Context) {
	h.proxyOpenAI(c, "/v1/responses", h.openAIGateway.Responses)
}

func (h *ImageGalleryHandler) proxyOpenAI(c *gin.Context, path string, handler func(*gin.Context)) {
	if h.openAIGateway == nil {
		response.Error(c, 500, "OpenAI gateway is unavailable")
		return
	}
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	apiKeyID, ok := parseProxyAPIKeyID(c)
	if !ok {
		return
	}
	apiKey, err := h.service.ValidateProxyAPIKey(c.Request.Context(), subject.UserID, apiKeyID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	subscription, ok := h.proxySubscription(c, apiKey)
	if !ok {
		return
	}

	c.Request.URL.Path = path
	c.Request.RequestURI = path
	c.Request.Header.Del("Authorization")
	c.Request.Header.Del("Proxy-Authorization")
	c.Request.Header.Del("X-Api-Key")
	c.Request.Header.Del("X-Goog-Api-Key")
	middleware2.SetAuthenticatedAPIKeyContext(c, apiKey, subscription)
	handler(c)
}

func (h *ImageGalleryHandler) proxySubscription(c *gin.Context, apiKey *service.APIKey) (*service.UserSubscription, bool) {
	if apiKey == nil || apiKey.Group == nil || apiKey.User == nil || !apiKey.Group.IsSubscriptionType() {
		return nil, true
	}
	if h.subscriptionService == nil {
		response.Error(c, 403, "No active subscription found for this group")
		return nil, false
	}
	subscription, err := h.subscriptionService.GetActiveSubscription(c.Request.Context(), apiKey.User.ID, apiKey.Group.ID)
	if err != nil {
		response.Error(c, 403, "No active subscription found for this group")
		return nil, false
	}
	needsMaintenance, err := h.subscriptionService.ValidateAndCheckLimits(subscription, apiKey.Group)
	if err != nil {
		status := 403
		if errors.Is(err, service.ErrDailyLimitExceeded) ||
			errors.Is(err, service.ErrWeeklyLimitExceeded) ||
			errors.Is(err, service.ErrMonthlyLimitExceeded) {
			status = 429
		}
		response.Error(c, status, err.Error())
		return nil, false
	}
	if needsMaintenance {
		maintenanceCopy := *subscription
		h.subscriptionService.DoWindowMaintenance(&maintenanceCopy)
	}
	return subscription, true
}

func parseProxyAPIKeyID(c *gin.Context) (int64, bool) {
	raw := strings.TrimSpace(c.GetHeader("X-Sub2API-Key-ID"))
	if raw == "" {
		response.BadRequest(c, "X-Sub2API-Key-ID header is required")
		return 0, false
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid X-Sub2API-Key-ID header")
		return 0, false
	}
	return id, true
}

func (h *ImageGalleryHandler) Templates(c *gin.Context) {
	params := galleryPagination(c)
	items, p, err := h.service.Templates(c.Request.Context(), false, c.Query("q"), c.Query("category"), params)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, p.Total, p.Page, p.PageSize)
}

func (h *ImageGalleryHandler) Generate(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	var req service.ImageGalleryGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	result, err := h.service.Generate(c.Request.Context(), subject.UserID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

func (h *ImageGalleryHandler) Job(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	job, err := h.service.Job(c.Request.Context(), subject.UserID, id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, job)
}

func (h *ImageGalleryHandler) History(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	items, p, err := h.service.History(c.Request.Context(), subject.UserID, galleryPagination(c))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, p.Total, p.Page, p.PageSize)
}

func (h *ImageGalleryHandler) UpdateItem(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req service.ImageGalleryUpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	item, err := h.service.UpdateItem(c.Request.Context(), subject.UserID, id, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

func (h *ImageGalleryHandler) DeleteItem(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.service.DeleteItem(c.Request.Context(), subject.UserID, id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

func (h *ImageGalleryHandler) Publish(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	item, err := h.service.Publish(c.Request.Context(), subject.UserID, id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

func (h *ImageGalleryHandler) Public(c *gin.Context) {
	subject, _ := middleware2.GetAuthSubjectFromContext(c)
	items, p, err := h.service.Public(c.Request.Context(), subject.UserID, galleryPagination(c))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, p.Total, p.Page, p.PageSize)
}

func (h *ImageGalleryHandler) Favorite(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.service.Favorite(c.Request.Context(), subject.UserID, id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "ok"})
}

func (h *ImageGalleryHandler) Asset(c *gin.Context) {
	subject, _ := middleware2.GetAuthSubjectFromContext(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	path, mimeType, err := h.service.OpenAsset(c.Request.Context(), subject.UserID, id, false)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if c.Query("download") == "1" {
		c.Header("Content-Disposition", "attachment")
	}
	c.Header("Content-Type", mimeType)
	c.File(path)
}

func galleryPagination(c *gin.Context) pagination.PaginationParams {
	page, pageSize := response.ParsePagination(c)
	return pagination.PaginationParams{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}
}

func parseIDParam(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid ID")
		return 0, false
	}
	return id, true
}
