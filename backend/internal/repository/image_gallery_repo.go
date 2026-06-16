package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type imageGalleryRepository struct {
	db *sql.DB
}

func NewImageGalleryRepository(db *sql.DB) service.ImageGalleryRepository {
	return &imageGalleryRepository{db: db}
}

func (r *imageGalleryRepository) CreateJob(ctx context.Context, job *service.ImageGenerationJob) error {
	now := time.Now()
	if len(job.Params) == 0 {
		job.Params = json.RawMessage(`{}`)
	}
	err := r.db.QueryRowContext(ctx, `
INSERT INTO image_generation_jobs (user_id, api_key_id, status, model, prompt, params, error, started_at, completed_at, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$10)
RETURNING id, created_at, updated_at
`, job.UserID, job.APIKeyID, job.Status, job.Model, job.Prompt, []byte(job.Params), job.Error, job.StartedAt, job.CompletedAt, now).
		Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)
	return err
}

func (r *imageGalleryRepository) UpdateJob(ctx context.Context, job *service.ImageGenerationJob) error {
	job.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
UPDATE image_generation_jobs
SET status=$2, error=$3, completed_at=$4, updated_at=$5
WHERE id=$1
`, job.ID, job.Status, job.Error, job.CompletedAt, job.UpdatedAt)
	return err
}

func (r *imageGalleryRepository) GetJob(ctx context.Context, userID, jobID int64) (*service.ImageGenerationJob, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, user_id, api_key_id, status, model, prompt, params, error, started_at, completed_at, created_at, updated_at
FROM image_generation_jobs
WHERE id=$1 AND user_id=$2
`, jobID, userID)
	return scanImageJob(row)
}

func (r *imageGalleryRepository) CreateAsset(ctx context.Context, asset *service.ImageAsset) error {
	now := time.Now()
	if len(asset.Params) == 0 {
		asset.Params = json.RawMessage(`{}`)
	}
	err := r.db.QueryRowContext(ctx, `
INSERT INTO image_assets (
  user_id, job_id, usage_log_id, api_key_id, path, thumb_path, sha256, mime_type, width, height, size_bytes,
  visibility, review_status, hidden, prompt_snapshot, model, params, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$18)
RETURNING id, created_at, updated_at
`, asset.UserID, asset.JobID, asset.UsageLogID, asset.APIKeyID, asset.Path, asset.ThumbPath, asset.SHA256, asset.MimeType, asset.Width, asset.Height, asset.SizeBytes, asset.Visibility, asset.ReviewStatus, asset.Hidden, asset.Prompt, asset.Model, []byte(asset.Params), now).
		Scan(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt)
	return err
}

func (r *imageGalleryRepository) GetAsset(ctx context.Context, id int64) (*service.ImageAsset, error) {
	row := r.db.QueryRowContext(ctx, imageAssetSelect(`a.id=$1 AND a.deleted_at IS NULL`, false), id)
	return scanImageAsset(row)
}

func (r *imageGalleryRepository) ListHistory(ctx context.Context, userID int64, params pagination.PaginationParams) ([]service.ImageAsset, *pagination.PaginationResult, error) {
	params = normalizePagination(params)
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM image_assets WHERE user_id=$1 AND deleted_at IS NULL`, userID).Scan(&total); err != nil {
		return nil, nil, err
	}
	rows, err := r.db.QueryContext(ctx, imageAssetSelect(`a.user_id=$1 AND a.deleted_at IS NULL`, false)+` ORDER BY a.created_at DESC LIMIT $2 OFFSET $3`, userID, params.Limit(), params.Offset())
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	items, err := scanImageAssets(rows)
	return items, paginationResult(total, params), err
}

func (r *imageGalleryRepository) ListPublic(ctx context.Context, userID int64, params pagination.PaginationParams) ([]service.ImageAsset, *pagination.PaginationResult, error) {
	params = normalizePagination(params)
	where := `visibility='public' AND review_status='approved' AND hidden=FALSE AND deleted_at IS NULL`
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM image_assets WHERE `+where).Scan(&total); err != nil {
		return nil, nil, err
	}
	rows, err := r.db.QueryContext(ctx, imageAssetSelect(`a.`+where, true)+` ORDER BY a.created_at DESC LIMIT $2 OFFSET $3`, userID, params.Limit(), params.Offset())
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	items, err := scanImageAssets(rows)
	return items, paginationResult(total, params), err
}

func (r *imageGalleryRepository) ListAdminAssets(ctx context.Context, status string, params pagination.PaginationParams) ([]service.ImageAsset, *pagination.PaginationResult, error) {
	params = normalizePagination(params)
	where := `deleted_at IS NULL`
	args := []any{}
	if status = strings.TrimSpace(status); status != "" {
		where += ` AND review_status=$1`
		args = append(args, status)
	}
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM image_assets WHERE `+where, args...).Scan(&total); err != nil {
		return nil, nil, err
	}
	args = append(args, params.Limit(), params.Offset())
	limitParam := len(args) - 1
	offsetParam := len(args)
	rows, err := r.db.QueryContext(ctx, imageAssetSelect(`a.`+where, false)+fmt.Sprintf(` ORDER BY a.created_at DESC LIMIT $%d OFFSET $%d`, limitParam, offsetParam), args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	items, err := scanImageAssets(rows)
	return items, paginationResult(total, params), err
}

func (r *imageGalleryRepository) UpdateAsset(ctx context.Context, asset *service.ImageAsset) error {
	asset.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
UPDATE image_assets
SET visibility=$2, review_status=$3, hidden=$4, updated_at=$5
WHERE id=$1 AND deleted_at IS NULL
`, asset.ID, asset.Visibility, asset.ReviewStatus, asset.Hidden, asset.UpdatedAt)
	return err
}

func (r *imageGalleryRepository) SoftDeleteAsset(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE image_assets SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *imageGalleryRepository) IsFavorite(ctx context.Context, userID, assetID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM image_gallery_favorites WHERE user_id=$1 AND asset_id=$2)`, userID, assetID).Scan(&exists)
	return exists, err
}

func (r *imageGalleryRepository) SetFavorite(ctx context.Context, userID, assetID int64, favorite bool) error {
	if favorite {
		_, err := r.db.ExecContext(ctx, `INSERT INTO image_gallery_favorites (user_id, asset_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, userID, assetID)
		return err
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM image_gallery_favorites WHERE user_id=$1 AND asset_id=$2`, userID, assetID)
	return err
}

func (r *imageGalleryRepository) ListTemplates(ctx context.Context, includeDisabled bool, query, category string, params pagination.PaginationParams) ([]service.ImageGalleryTemplate, *pagination.PaginationResult, error) {
	params = normalizePagination(params)
	where := []string{"1=1"}
	args := []any{}
	if !includeDisabled {
		where = append(where, "enabled=TRUE")
	}
	if q := strings.TrimSpace(query); q != "" {
		args = append(args, "%"+q+"%")
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR prompt ILIKE $%d)", len(args), len(args)))
	}
	if c := strings.TrimSpace(category); c != "" {
		args = append(args, c)
		where = append(where, fmt.Sprintf("category=$%d", len(args)))
	}
	whereSQL := strings.Join(where, " AND ")
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM image_templates WHERE `+whereSQL, args...).Scan(&total); err != nil {
		return nil, nil, err
	}
	args = append(args, params.Limit(), params.Offset())
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT id, title, prompt, category, tags, source, source_url, author, license, variables, enabled, featured, created_at, updated_at
FROM image_templates
WHERE %s
ORDER BY featured DESC, created_at DESC
LIMIT $%d OFFSET $%d
`, whereSQL, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	items := []service.ImageGalleryTemplate{}
	for rows.Next() {
		item, err := scanTemplateRows(rows)
		if err != nil {
			return nil, nil, err
		}
		items = append(items, *item)
	}
	return items, paginationResult(total, params), rows.Err()
}

func (r *imageGalleryRepository) UpdateTemplate(ctx context.Context, id int64, req service.ImageGalleryAdminTemplateUpdate) (*service.ImageGalleryTemplate, error) {
	current, err := r.getTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		current.Title = strings.TrimSpace(*req.Title)
	}
	if req.Prompt != nil {
		current.Prompt = strings.TrimSpace(*req.Prompt)
	}
	if req.Category != nil {
		current.Category = strings.TrimSpace(*req.Category)
	}
	if req.Tags != nil {
		current.Tags = *req.Tags
	}
	if req.Enabled != nil {
		current.Enabled = *req.Enabled
	}
	if req.Featured != nil {
		current.Featured = *req.Featured
	}
	if req.SourceURL != nil {
		current.SourceURL = strings.TrimSpace(*req.SourceURL)
	}
	if req.Author != nil {
		current.Author = strings.TrimSpace(*req.Author)
	}
	if req.License != nil {
		current.License = strings.TrimSpace(*req.License)
	}
	tags, _ := json.Marshal(current.Tags)
	err = r.db.QueryRowContext(ctx, `
UPDATE image_templates
SET title=$2, prompt=$3, category=$4, tags=$5, source_url=$6, author=$7, license=$8, enabled=$9, featured=$10, updated_at=NOW()
WHERE id=$1
RETURNING updated_at
`, id, current.Title, current.Prompt, current.Category, tags, current.SourceURL, current.Author, current.License, current.Enabled, current.Featured).Scan(&current.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return current, nil
}

func (r *imageGalleryRepository) ImportTemplates(ctx context.Context, req service.ImageGalleryImportTemplatesRequest) (int, int, error) {
	imported := 0
	skipped := 0
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()
	for _, item := range req.Templates {
		if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.Prompt) == "" {
			skipped++
			continue
		}
		tags, _ := json.Marshal(item.Tags)
		variables := item.Variables
		if len(variables) == 0 {
			variables = json.RawMessage(`{}`)
		}
		source := firstRepoNonEmpty(item.Source, req.Source, "manual")
		sourceURL := firstRepoNonEmpty(item.SourceURL, req.SourceURL)
		_, err := tx.ExecContext(ctx, `
INSERT INTO image_templates (title, prompt, category, tags, source, source_url, author, license, variables, enabled, featured)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`, item.Title, item.Prompt, item.Category, tags, source, sourceURL, item.Author, item.License, []byte(variables), item.Enabled, item.Featured)
		if err != nil {
			return imported, skipped, err
		}
		imported++
	}
	_, _ = tx.ExecContext(ctx, `INSERT INTO image_template_import_runs (source, source_url, status, imported_count, skipped_count) VALUES ($1,$2,'completed',$3,$4)`, firstRepoNonEmpty(req.Source, "manual"), req.SourceURL, imported, skipped)
	return imported, skipped, tx.Commit()
}

func (r *imageGalleryRepository) FindLatestUsageLogID(ctx context.Context, userID, apiKeyID int64, since time.Time) (*int64, *float64, error) {
	var id int64
	var cost sql.NullFloat64
	err := r.db.QueryRowContext(ctx, `
SELECT id, actual_cost
FROM usage_logs
WHERE user_id=$1 AND api_key_id=$2 AND image_count > 0 AND created_at >= $3
ORDER BY created_at DESC
LIMIT 1
`, userID, apiKeyID, since).Scan(&id, &cost)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	var costPtr *float64
	if cost.Valid {
		costPtr = &cost.Float64
	}
	return &id, costPtr, nil
}

func (r *imageGalleryRepository) StorageBytesByUser(ctx context.Context, userID int64) (int64, error) {
	var total sql.NullInt64
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(size_bytes),0) FROM image_assets WHERE user_id=$1 AND deleted_at IS NULL`, userID).Scan(&total)
	return total.Int64, err
}

func (r *imageGalleryRepository) DeletedOrExpiredAssets(ctx context.Context, cutoff *time.Time, limit int) ([]service.ImageAsset, error) {
	where := `deleted_at IS NOT NULL`
	args := []any{}
	if cutoff != nil {
		args = append(args, *cutoff)
		where = `deleted_at IS NOT NULL OR a.created_at < $1`
	}
	args = append(args, limit)
	rows, err := r.db.QueryContext(ctx, imageAssetSelect(`a.`+where, false)+fmt.Sprintf(` ORDER BY a.created_at ASC LIMIT $%d`, len(args)), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanImageAssets(rows)
}

func (r *imageGalleryRepository) HardDeleteAssetRecord(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM image_assets WHERE id=$1`, id)
	return err
}

func (r *imageGalleryRepository) getTemplate(ctx context.Context, id int64) (*service.ImageGalleryTemplate, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, title, prompt, category, tags, source, source_url, author, license, variables, enabled, featured, created_at, updated_at
FROM image_templates WHERE id=$1
`, id)
	return scanTemplate(row)
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanImageJob(row rowScanner) (*service.ImageGenerationJob, error) {
	var job service.ImageGenerationJob
	var params []byte
	err := row.Scan(&job.ID, &job.UserID, &job.APIKeyID, &job.Status, &job.Model, &job.Prompt, &params, &job.Error, &job.StartedAt, &job.CompletedAt, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrImageGalleryNotFound
		}
		return nil, err
	}
	job.Params = json.RawMessage(params)
	return &job, nil
}

func imageAssetSelect(where string, includeFavorite bool) string {
	fav := `FALSE AS favorited`
	if includeFavorite {
		fav = `EXISTS(SELECT 1 FROM image_gallery_favorites f WHERE f.asset_id=a.id AND f.user_id=$1) AS favorited`
	}
	return `
SELECT a.id, a.user_id, a.job_id, a.usage_log_id, a.api_key_id, a.path, a.thumb_path, a.sha256, a.mime_type,
       a.width, a.height, a.size_bytes, a.visibility, a.review_status, a.hidden, a.prompt_snapshot, a.model, a.params,
       ul.actual_cost, ` + fav + `, a.created_at, a.updated_at
FROM image_assets a
LEFT JOIN usage_logs ul ON ul.id = a.usage_log_id
WHERE ` + where
}

func scanImageAsset(row rowScanner) (*service.ImageAsset, error) {
	var item service.ImageAsset
	var params []byte
	var cost sql.NullFloat64
	err := row.Scan(&item.ID, &item.UserID, &item.JobID, &item.UsageLogID, &item.APIKeyID, &item.Path, &item.ThumbPath, &item.SHA256, &item.MimeType, &item.Width, &item.Height, &item.SizeBytes, &item.Visibility, &item.ReviewStatus, &item.Hidden, &item.Prompt, &item.Model, &params, &cost, &item.Favorited, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrImageGalleryNotFound
		}
		return nil, err
	}
	item.Params = json.RawMessage(params)
	if cost.Valid {
		item.ActualCost = &cost.Float64
	}
	return &item, nil
}

func scanImageAssets(rows *sql.Rows) ([]service.ImageAsset, error) {
	items := []service.ImageAsset{}
	for rows.Next() {
		item, err := scanImageAsset(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, rows.Err()
}

func scanTemplate(row rowScanner) (*service.ImageGalleryTemplate, error) {
	return scanTemplateRows(row)
}

func scanTemplateRows(row rowScanner) (*service.ImageGalleryTemplate, error) {
	var item service.ImageGalleryTemplate
	var tags []byte
	var variables []byte
	err := row.Scan(&item.ID, &item.Title, &item.Prompt, &item.Category, &tags, &item.Source, &item.SourceURL, &item.Author, &item.License, &variables, &item.Enabled, &item.Featured, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrImageGalleryNotFound
		}
		return nil, err
	}
	_ = json.Unmarshal(tags, &item.Tags)
	item.Variables = json.RawMessage(variables)
	return &item, nil
}

func normalizePagination(params pagination.PaginationParams) pagination.PaginationParams {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	return params
}

func paginationResult(total int64, params pagination.PaginationParams) *pagination.PaginationResult {
	pages := int(math.Ceil(float64(total) / float64(params.Limit())))
	if pages < 1 {
		pages = 1
	}
	return &pagination.PaginationResult{Total: total, Page: params.Page, PageSize: params.Limit(), Pages: pages}
}

func firstRepoNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
