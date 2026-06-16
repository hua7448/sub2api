package admin

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type ImageGalleryHandler struct {
	service *service.ImageGalleryService
}

func NewImageGalleryHandler(service *service.ImageGalleryService) *ImageGalleryHandler {
	return &ImageGalleryHandler{service: service}
}

func (h *ImageGalleryHandler) GetSettings(c *gin.Context) {
	settings, err := h.service.Settings(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, settings)
}

func (h *ImageGalleryHandler) UpdateSettings(c *gin.Context) {
	var req service.ImageGallerySettings
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	settings, err := h.service.UpdateSettings(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, settings)
}

func (h *ImageGalleryHandler) Items(c *gin.Context) {
	items, p, err := h.service.AdminItems(c.Request.Context(), c.Query("review_status"), adminGalleryPagination(c))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, p.Total, p.Page, p.PageSize)
}

func (h *ImageGalleryHandler) Approve(c *gin.Context) {
	h.setReview(c, service.ImageGalleryReviewApproved)
}

func (h *ImageGalleryHandler) Reject(c *gin.Context) {
	h.setReview(c, service.ImageGalleryReviewRejected)
}

func (h *ImageGalleryHandler) Hide(c *gin.Context) {
	h.setReview(c, "hidden")
}

func (h *ImageGalleryHandler) Delete(c *gin.Context) {
	id, ok := adminParseID(c)
	if !ok {
		return
	}
	if err := h.service.AdminDelete(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

func (h *ImageGalleryHandler) Templates(c *gin.Context) {
	items, p, err := h.service.Templates(c.Request.Context(), true, c.Query("q"), c.Query("category"), adminGalleryPagination(c))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, p.Total, p.Page, p.PageSize)
}

func (h *ImageGalleryHandler) ImportTemplates(c *gin.Context) {
	var req service.ImageGalleryImportTemplatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	imported, skipped, err := h.service.AdminImportTemplates(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"imported": imported, "skipped": skipped})
}

func (h *ImageGalleryHandler) UpdateTemplate(c *gin.Context) {
	id, ok := adminParseID(c)
	if !ok {
		return
	}
	var req service.ImageGalleryAdminTemplateUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	item, err := h.service.AdminUpdateTemplate(c.Request.Context(), id, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

func (h *ImageGalleryHandler) StorageStats(c *gin.Context) {
	stats, err := h.service.StorageStats(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, stats)
}

func (h *ImageGalleryHandler) Cleanup(c *gin.Context) {
	result, err := h.service.Cleanup(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

func (h *ImageGalleryHandler) setReview(c *gin.Context, status string) {
	id, ok := adminParseID(c)
	if !ok {
		return
	}
	item, err := h.service.AdminSetReview(c.Request.Context(), id, status)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, item)
}

func adminGalleryPagination(c *gin.Context) pagination.PaginationParams {
	page, pageSize := response.ParsePagination(c)
	return pagination.PaginationParams{Page: page, PageSize: pageSize, SortBy: "created_at", SortOrder: "desc"}
}

func adminParseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid ID")
		return 0, false
	}
	return id, true
}
