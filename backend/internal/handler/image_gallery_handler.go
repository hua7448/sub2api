package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
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
	h.proxyOpenAI(c, "/v1/images/generations", h.openAIGateway.Images, h.validateImagesGenerationProxyRequest)
}

func (h *ImageGalleryHandler) ProxyImagesEdits(c *gin.Context) {
	h.proxyOpenAI(c, "/v1/images/edits", h.openAIGateway.Images, h.validateImagesEditProxyRequest)
}

func (h *ImageGalleryHandler) ProxyResponses(c *gin.Context) {
	h.proxyOpenAI(c, "/v1/responses", h.openAIGateway.Responses, h.validateResponsesProxyRequest)
}

func (h *ImageGalleryHandler) proxyOpenAI(c *gin.Context, path string, handler func(*gin.Context), guard func(*gin.Context, service.ImageGallerySettings) bool) {
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
	if guard != nil {
		settings, err := h.service.Settings(c.Request.Context())
		if err != nil {
			response.ErrorFrom(c, err)
			return
		}
		if !guard(c, settings) {
			return
		}
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

type imageGalleryProxyParams struct {
	Model        string
	Size         string
	Quality      string
	OutputFormat string
	N            int
}

type imageGalleryResponsesProxyInspection struct {
	Model            string
	HasImageTool     bool
	HasAgentTool     bool
	HasWebSearchTool bool
	ImageToolParams  []imageGalleryProxyParams
}

func (h *ImageGalleryHandler) validateImagesGenerationProxyRequest(c *gin.Context, settings service.ImageGallerySettings) bool {
	body, ok := readAndRestoreImageGalleryProxyBody(c)
	if !ok {
		return false
	}
	params, err := parseImageGalleryJSONProxyParams(body)
	if err != nil {
		response.ErrorFrom(c, err)
		return false
	}
	if err := validateImageGalleryProxyParams(params, settings, settings.AllowedModels, true); err != nil {
		response.ErrorFrom(c, err)
		return false
	}
	return true
}

func (h *ImageGalleryHandler) validateImagesEditProxyRequest(c *gin.Context, settings service.ImageGallerySettings) bool {
	body, ok := readAndRestoreImageGalleryProxyBody(c)
	if !ok {
		return false
	}
	params, imageCount, err := parseImageGalleryMultipartProxyParams(body, c.GetHeader("Content-Type"), settings)
	if err != nil {
		response.ErrorFrom(c, err)
		return false
	}
	if imageCount == 0 {
		response.ErrorFrom(c, infraerrors.BadRequest("IMAGE_GALLERY_IMAGE_REQUIRED", "image edit requires at least one reference image"))
		return false
	}
	if imageCount > 4 {
		response.ErrorFrom(c, infraerrors.BadRequest("IMAGE_GALLERY_TOO_MANY_REFERENCE_IMAGES", "too many reference images"))
		return false
	}
	if err := validateImageGalleryProxyParams(params, settings, settings.AllowedModels, true); err != nil {
		response.ErrorFrom(c, err)
		return false
	}
	return true
}

func (h *ImageGalleryHandler) validateResponsesProxyRequest(c *gin.Context, settings service.ImageGallerySettings) bool {
	body, ok := readAndRestoreImageGalleryProxyBody(c)
	if !ok {
		return false
	}
	inspection, err := inspectImageGalleryResponsesProxyBody(body)
	if err != nil {
		response.ErrorFrom(c, err)
		return false
	}
	if !inspection.HasImageTool {
		response.ErrorFrom(c, infraerrors.BadRequest("IMAGE_GALLERY_RESPONSES_IMAGE_TOOL_REQUIRED", "responses proxy requires the image_generation tool"))
		return false
	}
	if !settings.AgentEnabled && inspection.HasAgentTool {
		response.ErrorFrom(c, infraerrors.Forbidden("IMAGE_GALLERY_AGENT_DISABLED", "image playground agent mode is disabled"))
		return false
	}
	if !settings.AgentWebSearchEnabled && inspection.HasWebSearchTool {
		response.ErrorFrom(c, infraerrors.Forbidden("IMAGE_GALLERY_AGENT_WEB_SEARCH_DISABLED", "agent web search is disabled"))
		return false
	}
	params := imageGalleryProxyParams{Model: inspection.Model, N: 1}
	if err := validateImageGalleryProxyParams(params, settings, settings.AllowedAgentModels, false); err != nil {
		response.ErrorFrom(c, err)
		return false
	}
	for _, imageParams := range inspection.ImageToolParams {
		if err := validateImageGalleryProxyImageParams(imageParams, settings); err != nil {
			response.ErrorFrom(c, err)
			return false
		}
	}
	return true
}

func readAndRestoreImageGalleryProxyBody(c *gin.Context) ([]byte, bool) {
	if c.Request.Body == nil {
		return nil, true
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to read request body")
		return nil, false
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	c.Request.ContentLength = int64(len(body))
	return body, true
}

func parseImageGalleryJSONProxyParams(body []byte) (imageGalleryProxyParams, error) {
	var payload map[string]any
	if len(body) == 0 || json.Unmarshal(body, &payload) != nil {
		return imageGalleryProxyParams{}, infraerrors.BadRequest("IMAGE_GALLERY_INVALID_PROXY_REQUEST", "invalid image proxy request body")
	}
	return imageGalleryProxyParams{
		Model:        imageGalleryStringField(payload, "model"),
		Size:         imageGalleryStringField(payload, "size"),
		Quality:      imageGalleryStringField(payload, "quality"),
		OutputFormat: strings.ToLower(imageGalleryStringField(payload, "output_format")),
		N:            imageGalleryIntField(payload, "n", 1),
	}, nil
}

func parseImageGalleryMultipartProxyParams(body []byte, contentType string, settings service.ImageGallerySettings) (imageGalleryProxyParams, int, error) {
	mediaType, mediaParams, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.HasPrefix(strings.ToLower(mediaType), "multipart/") || mediaParams["boundary"] == "" {
		return imageGalleryProxyParams{}, 0, infraerrors.BadRequest("IMAGE_GALLERY_INVALID_PROXY_REQUEST", "invalid multipart image proxy request")
	}
	maxMemory := int64(settings.MaxUploadMB)
	if maxMemory <= 0 {
		maxMemory = 20
	}
	maxMemory = maxMemory * 1024 * 1024 * 6
	form, err := multipart.NewReader(bytes.NewReader(body), mediaParams["boundary"]).ReadForm(maxMemory)
	if err != nil {
		return imageGalleryProxyParams{}, 0, infraerrors.BadRequest("IMAGE_GALLERY_INVALID_PROXY_REQUEST", "invalid multipart image proxy request")
	}
	defer form.RemoveAll()

	maxFileBytes := int64(settings.MaxUploadMB)
	if maxFileBytes <= 0 {
		maxFileBytes = 20
	}
	maxFileBytes *= 1024 * 1024
	imageCount := 0
	for field, files := range form.File {
		switch field {
		case "image", "image[]":
			imageCount += len(files)
		}
		for _, file := range files {
			if file != nil && file.Size > maxFileBytes {
				return imageGalleryProxyParams{}, 0, service.ErrImageGalleryUploadTooLarge
			}
		}
	}
	params := imageGalleryProxyParams{
		Model:        firstImageGalleryFormValue(form.Value, "model"),
		Size:         firstImageGalleryFormValue(form.Value, "size"),
		Quality:      firstImageGalleryFormValue(form.Value, "quality"),
		OutputFormat: strings.ToLower(firstImageGalleryFormValue(form.Value, "output_format")),
		N:            imageGalleryIntString(firstImageGalleryFormValue(form.Value, "n"), 1),
	}
	return params, imageCount, nil
}

func inspectImageGalleryResponsesProxyBody(body []byte) (imageGalleryResponsesProxyInspection, error) {
	var payload map[string]any
	if len(body) == 0 || json.Unmarshal(body, &payload) != nil {
		return imageGalleryResponsesProxyInspection{}, infraerrors.BadRequest("IMAGE_GALLERY_INVALID_PROXY_REQUEST", "invalid responses proxy request body")
	}
	tools, ok := payload["tools"].([]any)
	if !ok || len(tools) == 0 {
		return imageGalleryResponsesProxyInspection{}, infraerrors.BadRequest("IMAGE_GALLERY_RESPONSES_IMAGE_TOOL_REQUIRED", "responses proxy requires the image_generation tool")
	}
	inspection := imageGalleryResponsesProxyInspection{
		Model: imageGalleryStringField(payload, "model"),
	}
	for _, rawTool := range tools {
		tool, ok := rawTool.(map[string]any)
		if !ok {
			return imageGalleryResponsesProxyInspection{}, infraerrors.BadRequest("IMAGE_GALLERY_INVALID_PROXY_REQUEST", "invalid responses tool definition")
		}
		toolType := strings.ToLower(imageGalleryStringField(tool, "type"))
		toolName := strings.ToLower(imageGalleryStringField(tool, "name"))
		if strings.Contains(toolType, "web_search") || strings.Contains(toolName, "web_search") {
			inspection.HasWebSearchTool = true
			inspection.HasAgentTool = true
		}
		if toolType == "image_generation" {
			inspection.HasImageTool = true
			inspection.ImageToolParams = append(inspection.ImageToolParams, imageGalleryProxyParams{
				Size:         imageGalleryStringField(tool, "size"),
				Quality:      imageGalleryStringField(tool, "quality"),
				OutputFormat: strings.ToLower(imageGalleryStringField(tool, "output_format")),
				N:            1,
			})
			continue
		}
		if toolType != "" {
			inspection.HasAgentTool = true
		}
		if toolType == "function" || toolName == "generate_image_batch" || toolName == "continue_generation" {
			inspection.HasAgentTool = true
		}
	}
	return inspection, nil
}

func validateImageGalleryProxyParams(params imageGalleryProxyParams, settings service.ImageGallerySettings, allowedModels []string, imageModel bool) error {
	model := strings.TrimSpace(params.Model)
	if model == "" {
		return infraerrors.BadRequest("IMAGE_GALLERY_MODEL_REQUIRED", "image proxy model is required")
	}
	if !imageGalleryContainsString(allowedModels, model) {
		if imageModel {
			return infraerrors.BadRequest("IMAGE_GALLERY_MODEL_NOT_ALLOWED", "image model is not allowed")
		}
		return infraerrors.BadRequest("IMAGE_GALLERY_AGENT_MODEL_NOT_ALLOWED", "agent model is not allowed")
	}
	if params.N <= 0 {
		params.N = 1
	}
	if params.N > settings.MaxN {
		return infraerrors.BadRequest("IMAGE_GALLERY_N_TOO_LARGE", "image count exceeds the configured limit")
	}
	return validateImageGalleryProxyImageParams(params, settings)
}

func validateImageGalleryProxyImageParams(params imageGalleryProxyParams, settings service.ImageGallerySettings) error {
	if params.Size != "" && !imageGalleryContainsString(settings.AllowedSizes, params.Size) {
		return infraerrors.BadRequest("IMAGE_GALLERY_SIZE_NOT_ALLOWED", "image size is not allowed")
	}
	if params.Quality != "" && !imageGalleryContainsString(settings.AllowedQuality, params.Quality) {
		return infraerrors.BadRequest("IMAGE_GALLERY_QUALITY_NOT_ALLOWED", "image quality is not allowed")
	}
	if params.OutputFormat != "" && !imageGalleryContainsString(settings.AllowedOutputFormats, params.OutputFormat) {
		return infraerrors.BadRequest("IMAGE_GALLERY_FORMAT_NOT_ALLOWED", "image output format is not allowed")
	}
	return nil
}

func imageGalleryStringField(payload map[string]any, key string) string {
	value, _ := payload[key].(string)
	return strings.TrimSpace(value)
}

func imageGalleryIntField(payload map[string]any, key string, fallback int) int {
	switch value := payload[key].(type) {
	case float64:
		return int(value)
	case int:
		return value
	case json.Number:
		out, err := strconv.Atoi(value.String())
		if err == nil {
			return out
		}
	case string:
		return imageGalleryIntString(value, fallback)
	}
	return fallback
}

func imageGalleryIntString(value string, fallback int) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	out, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return out
}

func firstImageGalleryFormValue(values map[string][]string, key string) string {
	for _, value := range values[key] {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func imageGalleryContainsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
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

func (h *ImageGalleryHandler) CreateGenerationJob(c *gin.Context) {
	h.createJob(c, false)
}

func (h *ImageGalleryHandler) CreateEditJob(c *gin.Context) {
	h.createJob(c, true)
}

func (h *ImageGalleryHandler) createJob(c *gin.Context, expectEdit bool) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	req, ok := h.bindJobRequest(c)
	if !ok {
		return
	}
	if expectEdit && len(req.ReferenceImages) == 0 {
		response.BadRequest(c, "image is required for image edit jobs")
		return
	}
	if !expectEdit {
		req.ReferenceImages = nil
		req.MaskImage = ""
	}
	result, err := h.service.CreateJob(c.Request.Context(), subject.UserID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Accepted(c, result)
}

func (h *ImageGalleryHandler) bindJobRequest(c *gin.Context) (service.ImageGalleryGenerateRequest, bool) {
	contentType := strings.ToLower(c.GetHeader("Content-Type"))
	if strings.HasPrefix(contentType, "multipart/form-data") {
		req, err := h.bindMultipartJobRequest(c)
		if err != nil {
			writeImageGalleryBindError(c, err)
			return service.ImageGalleryGenerateRequest{}, false
		}
		return req, true
	}
	var req service.ImageGalleryGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeImageGalleryBindError(c, err)
		return req, false
	}
	return req, true
}

func (h *ImageGalleryHandler) bindMultipartJobRequest(c *gin.Context) (service.ImageGalleryGenerateRequest, error) {
	settings, err := h.service.PublicSettings(c.Request.Context())
	if err != nil {
		return service.ImageGalleryGenerateRequest{}, err
	}
	maxBytes := int64(settings.MaxUploadMB) * 1024 * 1024
	if maxBytes <= 0 {
		maxBytes = 20 << 20
	}
	reader, err := c.Request.MultipartReader()
	if err != nil {
		return service.ImageGalleryGenerateRequest{}, err
	}

	var req service.ImageGalleryGenerateRequest
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return service.ImageGalleryGenerateRequest{}, err
		}
		name := part.FormName()
		if name == "" {
			continue
		}
		if part.FileName() == "" {
			valueBytes, err := io.ReadAll(io.LimitReader(part, 1<<20))
			if err != nil {
				return service.ImageGalleryGenerateRequest{}, err
			}
			assignImageGalleryJobField(&req, name, string(valueBytes))
			continue
		}
		data, err := io.ReadAll(io.LimitReader(part, maxBytes+1))
		if err != nil {
			return service.ImageGalleryGenerateRequest{}, err
		}
		if int64(len(data)) > maxBytes {
			return service.ImageGalleryGenerateRequest{}, service.ErrImageGalleryUploadTooLarge
		}
		if len(data) == 0 {
			return service.ImageGalleryGenerateRequest{}, infraerrors.BadRequest("IMAGE_GALLERY_INVALID_UPLOAD", "empty image upload")
		}
		contentType := part.Header.Get("Content-Type")
		if strings.TrimSpace(contentType) == "" {
			contentType = http.DetectContentType(data)
		}
		dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(data))
		if name == "mask" {
			req.MaskImage = dataURL
		} else if name == "image" || name == "image[]" || strings.HasPrefix(name, "image[") || strings.HasPrefix(name, "images") {
			req.ReferenceImages = append(req.ReferenceImages, dataURL)
		}
	}
	return req, nil
}

func assignImageGalleryJobField(req *service.ImageGalleryGenerateRequest, name, value string) {
	value = strings.TrimSpace(value)
	switch name {
	case "api_key_id":
		if id, err := strconv.ParseInt(value, 10, 64); err == nil {
			req.APIKeyID = id
		}
	case "client_task_id":
		req.ClientTaskID = value
	case "prompt":
		req.Prompt = value
	case "model":
		req.Model = value
	case "size":
		req.Size = value
	case "quality":
		req.Quality = value
	case "n":
		if n, err := strconv.Atoi(value); err == nil {
			req.N = n
		}
	case "output_format":
		req.OutputFormat = value
	}
}

func writeImageGalleryBindError(c *gin.Context, err error) {
	if maxErr, ok := extractMaxBytesError(err); ok {
		response.Error(c, http.StatusRequestEntityTooLarge, buildBodyTooLargeMessage(maxErr.Limit))
		return
	}
	if infraerrors.Code(err) == http.StatusRequestEntityTooLarge {
		response.ErrorFrom(c, err)
		return
	}
	response.BadRequest(c, "Invalid request: "+err.Error())
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
	job, err := h.service.JobStatus(c.Request.Context(), subject.UserID, id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, job)
}

func (h *ImageGalleryHandler) JobByClientTask(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	clientTaskID := strings.TrimSpace(c.Param("clientTaskId"))
	job, err := h.service.JobByClientTaskID(c.Request.Context(), subject.UserID, clientTaskID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, job)
}

func (h *ImageGalleryHandler) CancelJob(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	result, err := h.service.CancelJob(c.Request.Context(), subject.UserID, id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
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
