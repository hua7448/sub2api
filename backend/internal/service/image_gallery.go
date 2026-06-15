package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

const (
	SettingKeyImageGalleryEnabled               = "image_gallery_enabled"
	SettingKeyImageGalleryPublicEnabled         = "image_gallery_public_enabled"
	SettingKeyImageGalleryPublishRequiresReview = "image_gallery_publish_requires_review"
	SettingKeyImageGalleryStorageBackend        = "image_gallery_storage_backend"
	SettingKeyImageGalleryStoragePath           = "image_gallery_storage_path"
	SettingKeyImageGalleryMaxUploadMB           = "image_gallery_max_upload_mb"
	SettingKeyImageGalleryUserQuotaMB           = "image_gallery_user_quota_mb"
	SettingKeyImageGalleryRetentionDays         = "image_gallery_retention_days"
	SettingKeyImageGalleryAllowedModels         = "image_gallery_allowed_models"
	SettingKeyImageGalleryDefaultModel          = "image_gallery_default_model"
	SettingKeyImageGalleryAllowedSizes          = "image_gallery_allowed_sizes"
	SettingKeyImageGalleryAllowedQuality        = "image_gallery_allowed_quality"
	SettingKeyImageGalleryAllowedOutputFormats  = "image_gallery_allowed_output_formats"
	SettingKeyImageGalleryMaxN                  = "image_gallery_max_n"
	SettingKeyImageGalleryTemplatesEnabled      = "image_gallery_templates_enabled"
	SettingKeyImageGalleryTemplateImportEnabled = "image_gallery_template_import_enabled"
)

const (
	ImageGalleryJobPending        = "pending"
	ImageGalleryJobRunning        = "running"
	ImageGalleryJobSucceeded      = "succeeded"
	ImageGalleryJobFailed         = "failed"
	ImageGalleryJobSaveFailed     = "save_failed"
	ImageGalleryVisibilityPrivate = "private"
	ImageGalleryVisibilityPublic  = "public"
	ImageGalleryReviewNone        = "none"
	ImageGalleryReviewPending     = "pending"
	ImageGalleryReviewApproved    = "approved"
	ImageGalleryReviewRejected    = "rejected"
)

var (
	ErrImageGalleryDisabled       = infraerrors.Forbidden("IMAGE_GALLERY_DISABLED", "image gallery is disabled")
	ErrImageGalleryPublicDisabled = infraerrors.Forbidden("IMAGE_GALLERY_PUBLIC_DISABLED", "public image gallery is disabled")
	ErrImageGalleryInvalidParam   = infraerrors.BadRequest("IMAGE_GALLERY_INVALID_PARAM", "invalid image gallery parameter")
	ErrImageGalleryNotFound       = infraerrors.NotFound("IMAGE_GALLERY_NOT_FOUND", "image gallery item not found")
	ErrImageGalleryForbidden      = infraerrors.Forbidden("IMAGE_GALLERY_FORBIDDEN", "image gallery access denied")
)

type ImageGallerySettings struct {
	Enabled               bool     `json:"enabled"`
	PublicEnabled         bool     `json:"public_enabled"`
	PublishRequiresReview bool     `json:"publish_requires_review"`
	StorageBackend        string   `json:"storage_backend"`
	StoragePath           string   `json:"storage_path"`
	MaxUploadMB           int      `json:"max_upload_mb"`
	UserQuotaMB           int      `json:"user_quota_mb"`
	RetentionDays         int      `json:"retention_days"`
	AllowedModels         []string `json:"allowed_models"`
	DefaultModel          string   `json:"default_model"`
	AllowedSizes          []string `json:"allowed_sizes"`
	AllowedQuality        []string `json:"allowed_quality"`
	AllowedOutputFormats  []string `json:"allowed_output_formats"`
	MaxN                  int      `json:"max_n"`
	TemplatesEnabled      bool     `json:"templates_enabled"`
	TemplateImportEnabled bool     `json:"template_import_enabled"`
}

type ImageGalleryTemplate struct {
	ID        int64           `json:"id"`
	Title     string          `json:"title"`
	Prompt    string          `json:"prompt"`
	Category  string          `json:"category"`
	Tags      []string        `json:"tags"`
	Source    string          `json:"source"`
	SourceURL string          `json:"source_url"`
	Author    string          `json:"author"`
	License   string          `json:"license"`
	Variables json.RawMessage `json:"variables"`
	Enabled   bool            `json:"enabled"`
	Featured  bool            `json:"featured"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ImageGenerationJob struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	APIKeyID    *int64          `json:"api_key_id,omitempty"`
	Status      string          `json:"status"`
	Model       string          `json:"model"`
	Prompt      string          `json:"prompt"`
	Params      json.RawMessage `json:"params"`
	Error       string          `json:"error"`
	StartedAt   *time.Time      `json:"started_at,omitempty"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ImageAsset struct {
	ID           int64           `json:"id"`
	UserID       int64           `json:"user_id"`
	JobID        *int64          `json:"job_id,omitempty"`
	UsageLogID   *int64          `json:"usage_log_id,omitempty"`
	APIKeyID     *int64          `json:"api_key_id,omitempty"`
	Path         string          `json:"-"`
	ThumbPath    string          `json:"-"`
	URL          string          `json:"url"`
	DownloadURL  string          `json:"download_url"`
	SHA256       string          `json:"sha256"`
	MimeType     string          `json:"mime_type"`
	Width        int             `json:"width"`
	Height       int             `json:"height"`
	SizeBytes    int64           `json:"size_bytes"`
	Visibility   string          `json:"visibility"`
	ReviewStatus string          `json:"review_status"`
	Hidden       bool            `json:"hidden"`
	Prompt       string          `json:"prompt"`
	Model        string          `json:"model"`
	Params       json.RawMessage `json:"params"`
	ActualCost   *float64        `json:"actual_cost,omitempty"`
	Favorited    bool            `json:"favorited"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type ImageGalleryEligibleKey struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	GroupName      string  `json:"group_name"`
	Status         string  `json:"status"`
	QuotaRemaining float64 `json:"quota_remaining"`
}

type ImageGalleryGenerateRequest struct {
	APIKeyID        int64    `json:"api_key_id"`
	Prompt          string   `json:"prompt"`
	Model           string   `json:"model"`
	Size            string   `json:"size"`
	Quality         string   `json:"quality"`
	N               int      `json:"n"`
	OutputFormat    string   `json:"output_format"`
	ReferenceImages []string `json:"reference_images"`
	MaskImage       string   `json:"mask_image"`
}

type ImageGalleryGenerateResult struct {
	Job    *ImageGenerationJob `json:"job"`
	Assets []ImageAsset        `json:"assets"`
	Raw    json.RawMessage     `json:"raw,omitempty"`
}

type ImageGalleryUpdateItemRequest struct {
	Visibility *string `json:"visibility"`
	Hidden     *bool   `json:"hidden"`
}

type ImageGalleryAdminTemplateUpdate struct {
	Title     *string   `json:"title"`
	Prompt    *string   `json:"prompt"`
	Category  *string   `json:"category"`
	Tags      *[]string `json:"tags"`
	Enabled   *bool     `json:"enabled"`
	Featured  *bool     `json:"featured"`
	SourceURL *string   `json:"source_url"`
	Author    *string   `json:"author"`
	License   *string   `json:"license"`
}

type ImageGalleryImportTemplatesRequest struct {
	Templates []ImageGalleryTemplate `json:"templates"`
	Source    string                 `json:"source"`
	SourceURL string                 `json:"source_url"`
}

type ImageGalleryStorageStats struct {
	StoragePath string `json:"storage_path"`
	FileCount   int64  `json:"file_count"`
	TotalBytes  int64  `json:"total_bytes"`
}

type ImageGalleryCleanupResult struct {
	DeletedFiles int64 `json:"deleted_files"`
	DeletedBytes int64 `json:"deleted_bytes"`
}

type ImageGalleryRepository interface {
	CreateJob(ctx context.Context, job *ImageGenerationJob) error
	UpdateJob(ctx context.Context, job *ImageGenerationJob) error
	GetJob(ctx context.Context, userID, jobID int64) (*ImageGenerationJob, error)
	CreateAsset(ctx context.Context, asset *ImageAsset) error
	GetAsset(ctx context.Context, id int64) (*ImageAsset, error)
	ListHistory(ctx context.Context, userID int64, params pagination.PaginationParams) ([]ImageAsset, *pagination.PaginationResult, error)
	ListPublic(ctx context.Context, userID int64, params pagination.PaginationParams) ([]ImageAsset, *pagination.PaginationResult, error)
	ListAdminAssets(ctx context.Context, status string, params pagination.PaginationParams) ([]ImageAsset, *pagination.PaginationResult, error)
	UpdateAsset(ctx context.Context, asset *ImageAsset) error
	SoftDeleteAsset(ctx context.Context, id int64) error
	IsFavorite(ctx context.Context, userID, assetID int64) (bool, error)
	SetFavorite(ctx context.Context, userID, assetID int64, favorite bool) error
	ListTemplates(ctx context.Context, includeDisabled bool, query, category string, params pagination.PaginationParams) ([]ImageGalleryTemplate, *pagination.PaginationResult, error)
	UpdateTemplate(ctx context.Context, id int64, req ImageGalleryAdminTemplateUpdate) (*ImageGalleryTemplate, error)
	ImportTemplates(ctx context.Context, req ImageGalleryImportTemplatesRequest) (int, int, error)
	FindLatestUsageLogID(ctx context.Context, userID, apiKeyID int64, since time.Time) (*int64, *float64, error)
	StorageBytesByUser(ctx context.Context, userID int64) (int64, error)
	DeletedOrExpiredAssets(ctx context.Context, cutoff *time.Time, limit int) ([]ImageAsset, error)
	HardDeleteAssetRecord(ctx context.Context, id int64) error
}

type ImageGatewayClient interface {
	Generate(ctx context.Context, apiKey string, request ImageGatewayRequest) ([]byte, int, http.Header, error)
}

type ImageGatewayRequest struct {
	Payload        map[string]any
	ReferenceFiles []ImageGatewayFile
	MaskFile       *ImageGatewayFile
}

type ImageGatewayFile struct {
	Name        string
	ContentType string
	Data        []byte
}

type LocalImageGatewayClient struct {
	BaseURL string
	Client  *http.Client
}

func (c *LocalImageGatewayClient) Generate(ctx context.Context, apiKey string, request ImageGatewayRequest) ([]byte, int, http.Header, error) {
	body, contentType, endpoint, err := request.Encode()
	if err != nil {
		return nil, 0, nil, err
	}
	baseURL := strings.TrimRight(c.BaseURL, "/")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:3000"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, 0, nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := c.Client
	if client == nil {
		client = &http.Client{Timeout: 180 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, nil, err
	}
	defer resp.Body.Close()
	respBody, readErr := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
	if readErr != nil {
		return nil, resp.StatusCode, resp.Header, readErr
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBody, resp.StatusCode, resp.Header, fmt.Errorf("image gateway returned HTTP %d", resp.StatusCode)
	}
	return respBody, resp.StatusCode, resp.Header, nil
}

func (r ImageGatewayRequest) Encode() ([]byte, string, string, error) {
	if len(r.ReferenceFiles) == 0 && r.MaskFile == nil {
		body, err := json.Marshal(r.Payload)
		return body, "application/json", "/v1/images/generations", err
	}
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	keys := make([]string, 0, len(r.Payload))
	for key := range r.Payload {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if r.Payload[key] == nil {
			continue
		}
		if err := writer.WriteField(key, fmt.Sprint(r.Payload[key])); err != nil {
			_ = writer.Close()
			return nil, "", "", err
		}
	}
	for i, file := range r.ReferenceFiles {
		fieldName := "image"
		if i > 0 {
			fieldName = fmt.Sprintf("image[%d]", i)
		}
		if err := writeImageGatewayFile(writer, fieldName, file); err != nil {
			_ = writer.Close()
			return nil, "", "", err
		}
	}
	if r.MaskFile != nil {
		if err := writeImageGatewayFile(writer, "mask", *r.MaskFile); err != nil {
			_ = writer.Close()
			return nil, "", "", err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), "/v1/images/edits", nil
}

type ImageGalleryService struct {
	repo          ImageGalleryRepository
	settingRepo   SettingRepository
	apiKeyService *APIKeyService
	gatewayClient ImageGatewayClient
}

func NewImageGalleryService(repo ImageGalleryRepository, settingRepo SettingRepository, apiKeyService *APIKeyService) *ImageGalleryService {
	return &ImageGalleryService{
		repo:          repo,
		settingRepo:   settingRepo,
		apiKeyService: apiKeyService,
		gatewayClient: &LocalImageGatewayClient{},
	}
}

func (s *ImageGalleryService) SetGatewayClient(client ImageGatewayClient) {
	if client != nil {
		s.gatewayClient = client
	}
}

func (s *ImageGalleryService) Settings(ctx context.Context) (ImageGallerySettings, error) {
	values, err := s.settingRepo.GetMultiple(ctx, imageGallerySettingKeys())
	if err != nil {
		return DefaultImageGallerySettings(), err
	}
	return parseImageGallerySettings(values), nil
}

func (s *ImageGalleryService) UpdateSettings(ctx context.Context, settings ImageGallerySettings) (ImageGallerySettings, error) {
	normalized := normalizeImageGallerySettings(settings)
	values := map[string]string{
		SettingKeyImageGalleryEnabled:               formatBool(normalized.Enabled),
		SettingKeyImageGalleryPublicEnabled:         formatBool(normalized.PublicEnabled),
		SettingKeyImageGalleryPublishRequiresReview: formatBool(normalized.PublishRequiresReview),
		SettingKeyImageGalleryStorageBackend:        normalized.StorageBackend,
		SettingKeyImageGalleryStoragePath:           normalized.StoragePath,
		SettingKeyImageGalleryMaxUploadMB:           fmt.Sprintf("%d", normalized.MaxUploadMB),
		SettingKeyImageGalleryUserQuotaMB:           fmt.Sprintf("%d", normalized.UserQuotaMB),
		SettingKeyImageGalleryRetentionDays:         fmt.Sprintf("%d", normalized.RetentionDays),
		SettingKeyImageGalleryAllowedModels:         mustJSON(normalized.AllowedModels),
		SettingKeyImageGalleryDefaultModel:          normalized.DefaultModel,
		SettingKeyImageGalleryAllowedSizes:          mustJSON(normalized.AllowedSizes),
		SettingKeyImageGalleryAllowedQuality:        mustJSON(normalized.AllowedQuality),
		SettingKeyImageGalleryAllowedOutputFormats:  mustJSON(normalized.AllowedOutputFormats),
		SettingKeyImageGalleryMaxN:                  fmt.Sprintf("%d", normalized.MaxN),
		SettingKeyImageGalleryTemplatesEnabled:      formatBool(normalized.TemplatesEnabled),
		SettingKeyImageGalleryTemplateImportEnabled: formatBool(normalized.TemplateImportEnabled),
	}
	if err := s.settingRepo.SetMultiple(ctx, values); err != nil {
		return normalized, err
	}
	return normalized, nil
}

func (s *ImageGalleryService) PublicSettings(ctx context.Context) (ImageGallerySettings, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return settings, err
	}
	settings.StoragePath = ""
	settings.StorageBackend = "local"
	return settings, nil
}

func (s *ImageGalleryService) EligibleKeys(ctx context.Context, userID int64) ([]ImageGalleryEligibleKey, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return nil, err
	}
	if !settings.Enabled {
		return nil, ErrImageGalleryDisabled
	}
	keys, _, err := s.apiKeyService.List(ctx, userID, pagination.PaginationParams{Page: 1, PageSize: 100, SortBy: "created_at", SortOrder: "desc"}, APIKeyListFilters{})
	if err != nil {
		return nil, err
	}
	out := make([]ImageGalleryEligibleKey, 0, len(keys))
	for i := range keys {
		k := keys[i]
		if !k.IsActive() || k.IsExpired() || k.IsQuotaExhausted() || !GroupAllowsImageGeneration(k.Group) {
			continue
		}
		groupName := ""
		if k.Group != nil {
			groupName = k.Group.Name
		}
		out = append(out, ImageGalleryEligibleKey{
			ID:             k.ID,
			Name:           k.Name,
			GroupName:      groupName,
			Status:         k.Status,
			QuotaRemaining: k.GetQuotaRemaining(),
		})
	}
	return out, nil
}

func (s *ImageGalleryService) Generate(ctx context.Context, userID int64, req ImageGalleryGenerateRequest) (*ImageGalleryGenerateResult, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return nil, err
	}
	if !settings.Enabled {
		return nil, ErrImageGalleryDisabled
	}
	req = normalizeGenerateRequest(req, settings)
	if err := validateGenerateRequest(req, settings); err != nil {
		return nil, err
	}
	apiKey, err := s.apiKeyService.GetByID(ctx, req.APIKeyID)
	if err != nil {
		return nil, err
	}
	if apiKey.UserID != userID {
		return nil, ErrImageGalleryForbidden
	}
	if apiKey.Key == "" {
		return nil, infraerrors.BadRequest("IMAGE_GALLERY_KEY_UNAVAILABLE", "api key secret is unavailable")
	}
	if !apiKey.IsActive() || apiKey.IsExpired() || apiKey.IsQuotaExhausted() {
		return nil, infraerrors.Forbidden("IMAGE_GALLERY_KEY_UNAVAILABLE", "api key is unavailable")
	}
	if !GroupAllowsImageGeneration(apiKey.Group) {
		return nil, infraerrors.Forbidden("IMAGE_GALLERY_IMAGE_NOT_ALLOWED", ImageGenerationPermissionMessage())
	}
	usedBytes, err := s.repo.StorageBytesByUser(ctx, userID)
	if err == nil && settings.UserQuotaMB > 0 && usedBytes >= int64(settings.UserQuotaMB)*1024*1024 {
		return nil, infraerrors.Forbidden("IMAGE_GALLERY_QUOTA_EXCEEDED", "image gallery storage quota exceeded")
	}
	referenceFiles, maskFile, err := decodeGenerateUploads(req, int64(settings.MaxUploadMB)*1024*1024)
	if err != nil {
		return nil, err
	}

	params := generateParamsJSON(req)
	job := &ImageGenerationJob{
		UserID:   userID,
		APIKeyID: &req.APIKeyID,
		Status:   ImageGalleryJobRunning,
		Model:    req.Model,
		Prompt:   req.Prompt,
		Params:   params,
	}
	now := time.Now()
	job.StartedAt = &now
	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, err
	}

	payload := map[string]any{
		"model":  req.Model,
		"prompt": req.Prompt,
		"n":      req.N,
	}
	if req.Size != "" {
		payload["size"] = req.Size
	}
	if req.Quality != "" {
		payload["quality"] = req.Quality
	}
	if req.OutputFormat != "" {
		payload["output_format"] = req.OutputFormat
	}
	started := time.Now()
	body, _, _, genErr := s.gatewayClient.Generate(ctx, apiKey.Key, ImageGatewayRequest{
		Payload:        payload,
		ReferenceFiles: referenceFiles,
		MaskFile:       maskFile,
	})
	if genErr != nil {
		job.Status = ImageGalleryJobFailed
		job.Error = truncateImageGalleryString(string(body), 2000)
		if job.Error == "" {
			job.Error = genErr.Error()
		}
		completed := time.Now()
		job.CompletedAt = &completed
		_ = s.repo.UpdateJob(ctx, job)
		return &ImageGalleryGenerateResult{Job: job, Raw: json.RawMessage(body)}, genErr
	}

	images, parseErr := extractGeneratedImages(body)
	if parseErr != nil {
		job.Status = ImageGalleryJobSaveFailed
		job.Error = parseErr.Error()
		completed := time.Now()
		job.CompletedAt = &completed
		_ = s.repo.UpdateJob(ctx, job)
		return &ImageGalleryGenerateResult{Job: job, Raw: json.RawMessage(body)}, parseErr
	}
	usageID, actualCost, _ := s.repo.FindLatestUsageLogID(ctx, userID, req.APIKeyID, started.Add(-2*time.Second))
	assets := make([]ImageAsset, 0, len(images))
	for _, img := range images {
		saved, saveErr := s.saveImage(ctx, settings, userID, job.ID, req.APIKeyID, usageID, actualCost, req, params, img)
		if saveErr != nil {
			job.Status = ImageGalleryJobSaveFailed
			job.Error = saveErr.Error()
			continue
		}
		assets = append(assets, *saved)
	}
	if len(assets) == 0 {
		if job.Error == "" {
			job.Error = "no image assets saved"
		}
		job.Status = ImageGalleryJobSaveFailed
	} else if job.Status != ImageGalleryJobSaveFailed {
		job.Status = ImageGalleryJobSucceeded
	}
	completed := time.Now()
	job.CompletedAt = &completed
	_ = s.repo.UpdateJob(ctx, job)
	return &ImageGalleryGenerateResult{Job: job, Assets: assets, Raw: json.RawMessage(body)}, nil
}

func (s *ImageGalleryService) saveImage(ctx context.Context, settings ImageGallerySettings, userID, jobID, apiKeyID int64, usageID *int64, actualCost *float64, req ImageGalleryGenerateRequest, params json.RawMessage, img generatedImage) (*ImageAsset, error) {
	data := img.Data
	if len(data) == 0 && img.URL != "" {
		downloaded, err := downloadImage(ctx, img.URL, int64(settings.MaxUploadMB)*1024*1024)
		if err != nil {
			return nil, err
		}
		data = downloaded
	}
	if len(data) == 0 {
		return nil, errors.New("empty image data")
	}
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])
	contentType := http.DetectContentType(data)
	ext := extensionForContentType(contentType, req.OutputFormat)
	now := time.Now()
	rel := filepath.ToSlash(filepath.Join("originals", now.Format("2006"), now.Format("01"), fmt.Sprintf("%d", userID), hash+ext))
	full := filepath.Join(settings.StoragePath, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(full, data, 0o644); err != nil {
		return nil, err
	}
	width, height := decodeImageSize(data)
	asset := &ImageAsset{
		UserID:       userID,
		JobID:        &jobID,
		UsageLogID:   usageID,
		APIKeyID:     &apiKeyID,
		Path:         rel,
		SHA256:       hash,
		MimeType:     contentType,
		Width:        width,
		Height:       height,
		SizeBytes:    int64(len(data)),
		Visibility:   ImageGalleryVisibilityPrivate,
		ReviewStatus: ImageGalleryReviewNone,
		Prompt:       req.Prompt,
		Model:        req.Model,
		Params:       params,
		ActualCost:   actualCost,
	}
	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		_ = os.Remove(full)
		return nil, err
	}
	return decorateAsset(asset), nil
}

func (s *ImageGalleryService) Job(ctx context.Context, userID, jobID int64) (*ImageGenerationJob, error) {
	return s.repo.GetJob(ctx, userID, jobID)
}

func (s *ImageGalleryService) History(ctx context.Context, userID int64, params pagination.PaginationParams) ([]ImageAsset, *pagination.PaginationResult, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return nil, nil, err
	}
	if !settings.Enabled {
		return nil, nil, ErrImageGalleryDisabled
	}
	items, p, err := s.repo.ListHistory(ctx, userID, params)
	decorateAssets(items)
	return items, p, err
}

func (s *ImageGalleryService) Public(ctx context.Context, userID int64, params pagination.PaginationParams) ([]ImageAsset, *pagination.PaginationResult, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return nil, nil, err
	}
	if !settings.Enabled {
		return nil, nil, ErrImageGalleryDisabled
	}
	if !settings.PublicEnabled {
		return nil, nil, ErrImageGalleryPublicDisabled
	}
	items, p, err := s.repo.ListPublic(ctx, userID, params)
	decorateAssets(items)
	return items, p, err
}

func (s *ImageGalleryService) UpdateItem(ctx context.Context, userID, id int64, req ImageGalleryUpdateItemRequest) (*ImageAsset, error) {
	asset, err := s.repo.GetAsset(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset.UserID != userID {
		return nil, ErrImageGalleryNotFound
	}
	if req.Visibility != nil {
		visibility := strings.TrimSpace(*req.Visibility)
		if visibility != ImageGalleryVisibilityPrivate && visibility != ImageGalleryVisibilityPublic {
			return nil, ErrImageGalleryInvalidParam
		}
		asset.Visibility = visibility
	}
	if req.Hidden != nil {
		asset.Hidden = *req.Hidden
	}
	if err := s.repo.UpdateAsset(ctx, asset); err != nil {
		return nil, err
	}
	return decorateAsset(asset), nil
}

func (s *ImageGalleryService) DeleteItem(ctx context.Context, userID, id int64) error {
	asset, err := s.repo.GetAsset(ctx, id)
	if err != nil {
		return err
	}
	if asset.UserID != userID {
		return ErrImageGalleryNotFound
	}
	return s.repo.SoftDeleteAsset(ctx, id)
}

func (s *ImageGalleryService) Publish(ctx context.Context, userID, id int64) (*ImageAsset, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return nil, err
	}
	if !settings.PublicEnabled {
		return nil, ErrImageGalleryPublicDisabled
	}
	asset, err := s.repo.GetAsset(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset.UserID != userID {
		return nil, ErrImageGalleryNotFound
	}
	asset.Visibility = ImageGalleryVisibilityPublic
	if settings.PublishRequiresReview {
		asset.ReviewStatus = ImageGalleryReviewPending
	} else {
		asset.ReviewStatus = ImageGalleryReviewApproved
	}
	if err := s.repo.UpdateAsset(ctx, asset); err != nil {
		return nil, err
	}
	return decorateAsset(asset), nil
}

func (s *ImageGalleryService) Favorite(ctx context.Context, userID, assetID int64) error {
	asset, err := s.repo.GetAsset(ctx, assetID)
	if err != nil {
		return err
	}
	if asset.Visibility != ImageGalleryVisibilityPublic || asset.ReviewStatus != ImageGalleryReviewApproved || asset.Hidden {
		return ErrImageGalleryNotFound
	}
	favorited, err := s.repo.IsFavorite(ctx, userID, assetID)
	if err != nil {
		return err
	}
	return s.repo.SetFavorite(ctx, userID, assetID, !favorited)
}

func (s *ImageGalleryService) Templates(ctx context.Context, includeDisabled bool, query, category string, params pagination.PaginationParams) ([]ImageGalleryTemplate, *pagination.PaginationResult, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return nil, nil, err
	}
	if !settings.TemplatesEnabled && !includeDisabled {
		return []ImageGalleryTemplate{}, &pagination.PaginationResult{Page: params.Page, PageSize: params.PageSize, Total: 0, Pages: 1}, nil
	}
	return s.repo.ListTemplates(ctx, includeDisabled, query, category, params)
}

func (s *ImageGalleryService) AdminItems(ctx context.Context, status string, params pagination.PaginationParams) ([]ImageAsset, *pagination.PaginationResult, error) {
	items, p, err := s.repo.ListAdminAssets(ctx, status, params)
	decorateAssets(items)
	return items, p, err
}

func (s *ImageGalleryService) AdminSetReview(ctx context.Context, id int64, status string) (*ImageAsset, error) {
	asset, err := s.repo.GetAsset(ctx, id)
	if err != nil {
		return nil, err
	}
	switch status {
	case ImageGalleryReviewApproved, ImageGalleryReviewRejected:
		asset.ReviewStatus = status
		asset.Visibility = ImageGalleryVisibilityPublic
	case "hidden":
		asset.Hidden = true
	default:
		return nil, ErrImageGalleryInvalidParam
	}
	if err := s.repo.UpdateAsset(ctx, asset); err != nil {
		return nil, err
	}
	return decorateAsset(asset), nil
}

func (s *ImageGalleryService) AdminDelete(ctx context.Context, id int64) error {
	return s.repo.SoftDeleteAsset(ctx, id)
}

func (s *ImageGalleryService) AdminUpdateTemplate(ctx context.Context, id int64, req ImageGalleryAdminTemplateUpdate) (*ImageGalleryTemplate, error) {
	return s.repo.UpdateTemplate(ctx, id, req)
}

func (s *ImageGalleryService) AdminImportTemplates(ctx context.Context, req ImageGalleryImportTemplatesRequest) (int, int, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return 0, 0, err
	}
	if !settings.TemplateImportEnabled {
		return 0, 0, infraerrors.Forbidden("IMAGE_GALLERY_TEMPLATE_IMPORT_DISABLED", "template import is disabled")
	}
	return s.repo.ImportTemplates(ctx, req)
}

func (s *ImageGalleryService) StorageStats(ctx context.Context) (ImageGalleryStorageStats, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return ImageGalleryStorageStats{}, err
	}
	stats := ImageGalleryStorageStats{StoragePath: settings.StoragePath}
	_ = filepath.WalkDir(settings.StoragePath, func(path string, d os.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return nil
		}
		info, statErr := d.Info()
		if statErr != nil {
			return nil
		}
		stats.FileCount++
		stats.TotalBytes += info.Size()
		return nil
	})
	return stats, nil
}

func (s *ImageGalleryService) Cleanup(ctx context.Context) (ImageGalleryCleanupResult, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return ImageGalleryCleanupResult{}, err
	}
	var cutoff *time.Time
	if settings.RetentionDays > 0 {
		t := time.Now().AddDate(0, 0, -settings.RetentionDays)
		cutoff = &t
	}
	items, err := s.repo.DeletedOrExpiredAssets(ctx, cutoff, 500)
	if err != nil {
		return ImageGalleryCleanupResult{}, err
	}
	var result ImageGalleryCleanupResult
	for _, item := range items {
		for _, rel := range []string{item.Path, item.ThumbPath} {
			if rel == "" {
				continue
			}
			full, ok := safeJoin(settings.StoragePath, rel)
			if !ok {
				continue
			}
			info, statErr := os.Stat(full)
			if statErr == nil {
				if rmErr := os.Remove(full); rmErr == nil {
					result.DeletedFiles++
					result.DeletedBytes += info.Size()
				}
			}
		}
		_ = s.repo.HardDeleteAssetRecord(ctx, item.ID)
	}
	return result, nil
}

func (s *ImageGalleryService) OpenAsset(ctx context.Context, userID int64, id int64, requireOwner bool) (string, string, error) {
	settings, err := s.Settings(ctx)
	if err != nil {
		return "", "", err
	}
	asset, err := s.repo.GetAsset(ctx, id)
	if err != nil {
		return "", "", err
	}
	allowed := asset.Visibility == ImageGalleryVisibilityPublic && asset.ReviewStatus == ImageGalleryReviewApproved && !asset.Hidden
	if requireOwner || !allowed {
		allowed = asset.UserID == userID
	}
	if !allowed {
		return "", "", ErrImageGalleryNotFound
	}
	full, ok := safeJoin(settings.StoragePath, asset.Path)
	if !ok {
		return "", "", ErrImageGalleryNotFound
	}
	return full, asset.MimeType, nil
}

func DefaultImageGallerySettings() ImageGallerySettings {
	return ImageGallerySettings{
		Enabled:               false,
		PublicEnabled:         false,
		PublishRequiresReview: true,
		StorageBackend:        "local",
		StoragePath:           "data/image-gallery",
		MaxUploadMB:           20,
		UserQuotaMB:           1024,
		RetentionDays:         0,
		AllowedModels:         []string{"gpt-image-2"},
		DefaultModel:          "gpt-image-2",
		AllowedSizes:          []string{"1024x1024", "1024x1536", "1536x1024", "auto"},
		AllowedQuality:        []string{"auto", "low", "medium", "high"},
		AllowedOutputFormats:  []string{"png", "jpeg", "webp"},
		MaxN:                  4,
		TemplatesEnabled:      true,
		TemplateImportEnabled: false,
	}
}

func imageGallerySettingKeys() []string {
	return []string{
		SettingKeyImageGalleryEnabled,
		SettingKeyImageGalleryPublicEnabled,
		SettingKeyImageGalleryPublishRequiresReview,
		SettingKeyImageGalleryStorageBackend,
		SettingKeyImageGalleryStoragePath,
		SettingKeyImageGalleryMaxUploadMB,
		SettingKeyImageGalleryUserQuotaMB,
		SettingKeyImageGalleryRetentionDays,
		SettingKeyImageGalleryAllowedModels,
		SettingKeyImageGalleryDefaultModel,
		SettingKeyImageGalleryAllowedSizes,
		SettingKeyImageGalleryAllowedQuality,
		SettingKeyImageGalleryAllowedOutputFormats,
		SettingKeyImageGalleryMaxN,
		SettingKeyImageGalleryTemplatesEnabled,
		SettingKeyImageGalleryTemplateImportEnabled,
	}
}

func parseImageGallerySettings(values map[string]string) ImageGallerySettings {
	defaults := DefaultImageGallerySettings()
	defaults.Enabled = parseBoolSetting(values[SettingKeyImageGalleryEnabled], defaults.Enabled)
	defaults.PublicEnabled = parseBoolSetting(values[SettingKeyImageGalleryPublicEnabled], defaults.PublicEnabled)
	defaults.PublishRequiresReview = parseBoolSetting(values[SettingKeyImageGalleryPublishRequiresReview], defaults.PublishRequiresReview)
	defaults.StorageBackend = firstImageGalleryNonEmpty(values[SettingKeyImageGalleryStorageBackend], defaults.StorageBackend)
	defaults.StoragePath = firstImageGalleryNonEmpty(values[SettingKeyImageGalleryStoragePath], defaults.StoragePath)
	defaults.MaxUploadMB = parseIntSetting(values[SettingKeyImageGalleryMaxUploadMB], defaults.MaxUploadMB)
	defaults.UserQuotaMB = parseIntSetting(values[SettingKeyImageGalleryUserQuotaMB], defaults.UserQuotaMB)
	defaults.RetentionDays = parseIntSetting(values[SettingKeyImageGalleryRetentionDays], defaults.RetentionDays)
	defaults.AllowedModels = parseStringSlice(values[SettingKeyImageGalleryAllowedModels], defaults.AllowedModels)
	defaults.DefaultModel = firstImageGalleryNonEmpty(values[SettingKeyImageGalleryDefaultModel], defaults.DefaultModel)
	defaults.AllowedSizes = parseStringSlice(values[SettingKeyImageGalleryAllowedSizes], defaults.AllowedSizes)
	defaults.AllowedQuality = parseStringSlice(values[SettingKeyImageGalleryAllowedQuality], defaults.AllowedQuality)
	defaults.AllowedOutputFormats = parseStringSlice(values[SettingKeyImageGalleryAllowedOutputFormats], defaults.AllowedOutputFormats)
	defaults.MaxN = parseIntSetting(values[SettingKeyImageGalleryMaxN], defaults.MaxN)
	defaults.TemplatesEnabled = parseBoolSetting(values[SettingKeyImageGalleryTemplatesEnabled], defaults.TemplatesEnabled)
	defaults.TemplateImportEnabled = parseBoolSetting(values[SettingKeyImageGalleryTemplateImportEnabled], defaults.TemplateImportEnabled)
	return normalizeImageGallerySettings(defaults)
}

func normalizeImageGallerySettings(settings ImageGallerySettings) ImageGallerySettings {
	defaults := DefaultImageGallerySettings()
	if settings.StorageBackend == "" {
		settings.StorageBackend = defaults.StorageBackend
	}
	if settings.StoragePath == "" {
		settings.StoragePath = defaults.StoragePath
	}
	if settings.MaxUploadMB <= 0 {
		settings.MaxUploadMB = defaults.MaxUploadMB
	}
	if settings.UserQuotaMB < 0 {
		settings.UserQuotaMB = defaults.UserQuotaMB
	}
	if settings.RetentionDays < 0 {
		settings.RetentionDays = 0
	}
	if len(settings.AllowedModels) == 0 {
		settings.AllowedModels = defaults.AllowedModels
	}
	if settings.DefaultModel == "" {
		settings.DefaultModel = settings.AllowedModels[0]
	}
	if len(settings.AllowedSizes) == 0 {
		settings.AllowedSizes = defaults.AllowedSizes
	}
	if len(settings.AllowedQuality) == 0 {
		settings.AllowedQuality = defaults.AllowedQuality
	}
	if len(settings.AllowedOutputFormats) == 0 {
		settings.AllowedOutputFormats = defaults.AllowedOutputFormats
	}
	if settings.MaxN <= 0 {
		settings.MaxN = defaults.MaxN
	}
	settings.AllowedModels = cleanStringSlice(settings.AllowedModels)
	settings.AllowedSizes = cleanStringSlice(settings.AllowedSizes)
	settings.AllowedQuality = cleanStringSlice(settings.AllowedQuality)
	settings.AllowedOutputFormats = cleanStringSlice(settings.AllowedOutputFormats)
	return settings
}

func normalizeGenerateRequest(req ImageGalleryGenerateRequest, settings ImageGallerySettings) ImageGalleryGenerateRequest {
	req.Prompt = strings.TrimSpace(req.Prompt)
	req.Model = firstImageGalleryNonEmpty(strings.TrimSpace(req.Model), settings.DefaultModel)
	req.Size = strings.TrimSpace(req.Size)
	req.Quality = strings.TrimSpace(req.Quality)
	req.OutputFormat = strings.ToLower(strings.TrimSpace(req.OutputFormat))
	cleanRefs := make([]string, 0, len(req.ReferenceImages))
	for _, image := range req.ReferenceImages {
		if image = strings.TrimSpace(image); image != "" {
			cleanRefs = append(cleanRefs, image)
		}
	}
	req.ReferenceImages = cleanRefs
	req.MaskImage = strings.TrimSpace(req.MaskImage)
	if req.N <= 0 {
		req.N = 1
	}
	return req
}

func validateGenerateRequest(req ImageGalleryGenerateRequest, settings ImageGallerySettings) error {
	if req.APIKeyID <= 0 || req.Prompt == "" {
		return ErrImageGalleryInvalidParam
	}
	if !containsString(settings.AllowedModels, req.Model) {
		return infraerrors.BadRequest("IMAGE_GALLERY_MODEL_NOT_ALLOWED", "image model is not allowed")
	}
	if req.Size != "" && !containsString(settings.AllowedSizes, req.Size) {
		return infraerrors.BadRequest("IMAGE_GALLERY_SIZE_NOT_ALLOWED", "image size is not allowed")
	}
	if req.Quality != "" && !containsString(settings.AllowedQuality, req.Quality) {
		return infraerrors.BadRequest("IMAGE_GALLERY_QUALITY_NOT_ALLOWED", "image quality is not allowed")
	}
	if req.OutputFormat != "" && !containsString(settings.AllowedOutputFormats, req.OutputFormat) {
		return infraerrors.BadRequest("IMAGE_GALLERY_FORMAT_NOT_ALLOWED", "image output format is not allowed")
	}
	if req.N > settings.MaxN {
		return infraerrors.BadRequest("IMAGE_GALLERY_N_TOO_LARGE", "image count exceeds the configured limit")
	}
	if len(req.ReferenceImages) > 4 {
		return infraerrors.BadRequest("IMAGE_GALLERY_TOO_MANY_REFERENCE_IMAGES", "too many reference images")
	}
	if req.MaskImage != "" && len(req.ReferenceImages) == 0 {
		return infraerrors.BadRequest("IMAGE_GALLERY_MASK_REQUIRES_REFERENCE", "mask requires a reference image")
	}
	return nil
}

type generatedImage struct {
	Data []byte
	URL  string
}

func extractGeneratedImages(body []byte) ([]generatedImage, error) {
	var payload struct {
		Data []struct {
			B64JSON string `json:"b64_json"`
			URL     string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	out := make([]generatedImage, 0, len(payload.Data))
	for _, item := range payload.Data {
		if item.B64JSON != "" {
			data, err := base64.StdEncoding.DecodeString(item.B64JSON)
			if err != nil {
				return nil, err
			}
			out = append(out, generatedImage{Data: data})
			continue
		}
		if strings.TrimSpace(item.URL) != "" {
			out = append(out, generatedImage{URL: strings.TrimSpace(item.URL)})
		}
	}
	if len(out) == 0 {
		return nil, errors.New("image response did not contain image data")
	}
	return out, nil
}

func decodeGenerateUploads(req ImageGalleryGenerateRequest, maxBytes int64) ([]ImageGatewayFile, *ImageGatewayFile, error) {
	if maxBytes <= 0 {
		maxBytes = 20 << 20
	}
	files := make([]ImageGatewayFile, 0, len(req.ReferenceImages))
	for i, raw := range req.ReferenceImages {
		file, err := decodeImageGatewayFile(raw, fmt.Sprintf("reference-%d.png", i+1), maxBytes)
		if err != nil {
			return nil, nil, err
		}
		files = append(files, file)
	}
	var mask *ImageGatewayFile
	if req.MaskImage != "" {
		file, err := decodeImageGatewayFile(req.MaskImage, "mask.png", maxBytes)
		if err != nil {
			return nil, nil, err
		}
		mask = &file
	}
	return files, mask, nil
}

func decodeImageGatewayFile(raw, name string, maxBytes int64) (ImageGatewayFile, error) {
	contentType := ""
	dataPart := strings.TrimSpace(raw)
	if strings.HasPrefix(strings.ToLower(dataPart), "data:") {
		header, body, ok := strings.Cut(dataPart, ",")
		if !ok {
			return ImageGatewayFile{}, infraerrors.BadRequest("IMAGE_GALLERY_INVALID_UPLOAD", "invalid image upload")
		}
		meta := strings.TrimPrefix(header, "data:")
		if semi := strings.Index(meta, ";"); semi >= 0 {
			contentType = meta[:semi]
		} else {
			contentType = meta
		}
		dataPart = body
	}
	data, err := base64.StdEncoding.DecodeString(dataPart)
	if err != nil {
		return ImageGatewayFile{}, infraerrors.BadRequest("IMAGE_GALLERY_INVALID_UPLOAD", "invalid image upload")
	}
	if len(data) == 0 {
		return ImageGatewayFile{}, infraerrors.BadRequest("IMAGE_GALLERY_INVALID_UPLOAD", "empty image upload")
	}
	if int64(len(data)) > maxBytes {
		return ImageGatewayFile{}, infraerrors.BadRequest("IMAGE_GALLERY_UPLOAD_TOO_LARGE", "image upload exceeds configured limit")
	}
	detected := http.DetectContentType(data)
	if contentType == "" || contentType == "application/octet-stream" {
		contentType = detected
	}
	if !strings.HasPrefix(strings.ToLower(contentType), "image/") || !strings.HasPrefix(strings.ToLower(detected), "image/") {
		return ImageGatewayFile{}, infraerrors.BadRequest("IMAGE_GALLERY_INVALID_UPLOAD", "upload must be an image")
	}
	return ImageGatewayFile{Name: name, ContentType: contentType, Data: data}, nil
}

func writeImageGatewayFile(writer *multipart.Writer, fieldName string, file ImageGatewayFile) error {
	contentType := strings.TrimSpace(file.ContentType)
	if contentType == "" {
		contentType = http.DetectContentType(file.Data)
	}
	fileName := strings.TrimSpace(file.Name)
	if fileName == "" {
		fileName = fieldName + extensionForContentType(contentType, "")
	}
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldName), escapeQuotes(fileName)))
	header.Set("Content-Type", contentType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}
	_, err = part.Write(file.Data)
	return err
}

func downloadImage(ctx context.Context, imageURL string, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		maxBytes = 20 << 20
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("download image returned HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("downloaded image exceeds %d bytes", maxBytes)
	}
	return data, nil
}

func decodeImageSize(data []byte) (int, int) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0
	}
	return cfg.Width, cfg.Height
}

func extensionForContentType(contentType, preferred string) string {
	preferred = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(preferred)), ".")
	switch preferred {
	case "png", "jpeg", "jpg", "webp":
		if preferred == "jpeg" {
			return ".jpg"
		}
		return "." + preferred
	}
	exts, _ := mime.ExtensionsByType(contentType)
	if len(exts) > 0 {
		return exts[0]
	}
	return ".png"
}

func decorateAsset(asset *ImageAsset) *ImageAsset {
	if asset == nil {
		return nil
	}
	asset.URL = fmt.Sprintf("/api/v1/image-gallery/assets/%d", asset.ID)
	asset.DownloadURL = asset.URL + "?download=1"
	return asset
}

func decorateAssets(items []ImageAsset) {
	for i := range items {
		decorateAsset(&items[i])
	}
}

func generateParamsJSON(req ImageGalleryGenerateRequest) json.RawMessage {
	body, _ := json.Marshal(map[string]any{
		"size":                  req.Size,
		"quality":               req.Quality,
		"n":                     req.N,
		"output_format":         req.OutputFormat,
		"reference_image_count": len(req.ReferenceImages),
		"has_mask":              req.MaskImage != "",
	})
	return body
}

func escapeQuotes(value string) string {
	return strings.ReplaceAll(value, `"`, `\"`)
}

func safeJoin(root, rel string) (string, bool) {
	if filepath.IsAbs(rel) || strings.Contains(rel, "..") {
		return "", false
	}
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", false
	}
	full, err := filepath.Abs(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return "", false
	}
	return full, full == rootAbs || strings.HasPrefix(full, rootAbs+string(os.PathSeparator))
}

func parseBoolSetting(value string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return fallback
	}
}

func parseIntSetting(value string, fallback int) int {
	var out int
	if _, err := fmt.Sscanf(strings.TrimSpace(value), "%d", &out); err != nil {
		return fallback
	}
	return out
}

func parseStringSlice(value string, fallback []string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	var out []string
	if err := json.Unmarshal([]byte(value), &out); err == nil {
		return cleanStringSlice(out)
	}
	parts := strings.Split(value, ",")
	return cleanStringSlice(parts)
}

func cleanStringSlice(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, item := range in {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	sort.Strings(out)
	return out
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func mustJSON(v any) string {
	body, _ := json.Marshal(v)
	return string(body)
}

func formatBool(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func firstImageGalleryNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func truncateImageGalleryString(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max]
}
