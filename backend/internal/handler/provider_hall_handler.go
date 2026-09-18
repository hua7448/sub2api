package handler

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// ProviderHallUserHandler serves the user-facing hall (batch B5). Every
// response is built from service results through the dto contract; nothing
// from the request context other than the user ID reaches the service.
type ProviderHallUserHandler struct {
	query *service.ProviderHallQueryService
}

func NewProviderHallUserHandler(query *service.ProviderHallQueryService) *ProviderHallUserHandler {
	return &ProviderHallUserHandler{query: query}
}

// DisplayGuard hides the routes entirely (404) while display is off. The
// service re-checks the switch on every call as well.
func (h *ProviderHallUserHandler) DisplayGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h == nil || h.query == nil || !h.query.DisplayEnabled(c.Request.Context()) {
			response.ErrorFrom(c, service.ErrProviderHallNotFound)
			c.Abort()
			return
		}
		c.Next()
	}
}

func providerHallUser(c *gin.Context) (int64, bool) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok || subject.UserID < 1 {
		response.Unauthorized(c, "user not found in context")
		return 0, false
	}
	return subject.UserID, true
}

func providerHallInvalidQuery(c *gin.Context, field string) {
	response.ErrorFrom(c, service.ErrProviderHallInvalidQuery.WithMetadata(map[string]string{"field": field}))
}

func providerHallIntQuery(c *gin.Context, name string) (int, bool) {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return 0, true
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 1 {
		providerHallInvalidQuery(c, name)
		return 0, false
	}
	return v, true
}

func providerHallGroupID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		response.ErrorFrom(c, service.ErrProviderHallNotFound)
		return 0, false
	}
	return id, true
}

// List GET /api/v1/provider-hall
func (h *ProviderHallUserHandler) List(c *gin.Context) {
	userID, ok := providerHallUser(c)
	if !ok {
		return
	}
	q := service.ProviderHallListQuery{
		Range:      strings.TrimSpace(c.Query("range")),
		Model:      strings.TrimSpace(c.Query("model")),
		Protocol:   strings.TrimSpace(c.Query("protocol")),
		Search:     c.Query("search"),
		SnapshotID: strings.TrimSpace(c.Query("snapshot_id")),
	}
	sortRules, ok := service.ParseProviderHallSort(c.Query("sort"))
	if !ok {
		providerHallInvalidQuery(c, "sort")
		return
	}
	q.Sort = sortRules
	if q.Page, ok = providerHallIntQuery(c, "page"); !ok {
		return
	}
	if q.PageSize, ok = providerHallIntQuery(c, "page_size"); !ok {
		return
	}
	result, err := h.query.List(c.Request.Context(), userID, q)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, providerHallListDTO(result))
}

// GetGroup GET /api/v1/provider-hall/groups/:id
func (h *ProviderHallUserHandler) GetGroup(c *gin.Context) {
	userID, ok := providerHallUser(c)
	if !ok {
		return
	}
	groupID, ok := providerHallGroupID(c)
	if !ok {
		return
	}
	result, err := h.query.GetGroup(c.Request.Context(), userID, groupID, strings.TrimSpace(c.Query("range")))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, providerHallDetailDTO(result))
}

// ListVerifications GET /api/v1/provider-hall/groups/:id/verifications
func (h *ProviderHallUserHandler) ListVerifications(c *gin.Context) {
	userID, ok := providerHallUser(c)
	if !ok {
		return
	}
	groupID, ok := providerHallGroupID(c)
	if !ok {
		return
	}
	var profileID *int64
	if raw := strings.TrimSpace(c.Query("profile_id")); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || id < 1 {
			providerHallInvalidQuery(c, "profile_id")
			return
		}
		profileID = &id
	}
	page, ok := providerHallIntQuery(c, "page")
	if !ok {
		return
	}
	pageSize, ok := providerHallIntQuery(c, "page_size")
	if !ok {
		return
	}
	items, total, err := h.query.ListVerifications(c.Request.Context(), userID, groupID, profileID, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	out := dto.ProviderHallVerificationListResponse{Items: make([]dto.ProviderHallVerificationReport, 0, len(items)), Pagination: dto.ProviderHallPagination{Page: page, PageSize: pageSize, Total: total}}
	for _, it := range items {
		out.Items = append(out.Items, providerHallReportDTO(it))
	}
	response.Success(c, out)
}

// DTO conversion. Kept in the handler because service cannot import dto.

func providerHallMetricDTO[T any](m service.ProviderHallMetric[T]) dto.ProviderHallMetric[T] {
	return dto.ProviderHallMetric[T]{Value: m.Value, State: dto.ProviderHallMetricState(m.State), ReasonCode: m.ReasonCode, WindowStart: m.WindowStart, WindowEnd: m.WindowEnd, ComputedAt: m.ComputedAt}
}

func providerHallRefDTO(r service.ProviderHallProfileRef) dto.ProviderHallProfileRef {
	return dto.ProviderHallProfileRef{ProfileID: r.ProfileID, Model: r.Model, Protocol: r.Protocol}
}

func providerHallModelHealthDTO(m service.ProviderHallModelHealthResult) dto.ProviderHallModelHealth {
	return dto.ProviderHallModelHealth{ProfileID: m.ProfileID, Model: m.Model, Protocol: m.Protocol, Status: m.Status, CheckedAt: m.CheckedAt}
}

func providerHallBadgeDTO(b service.ProviderHallVerificationBadgeResult) dto.ProviderHallVerificationBadge {
	return dto.ProviderHallVerificationBadge{Verdict: b.Verdict, ReasonCode: b.ReasonCode, CompletedAt: b.CompletedAt, Expired: b.Expired, ReportID: b.ReportID}
}

func providerHallRowDTO(r *service.ProviderHallRowResult) dto.ProviderHallRow {
	models := make([]dto.ProviderHallModelHealth, 0, len(r.Health.Models))
	for _, m := range r.Health.Models {
		models = append(models, providerHallModelHealthDTO(m))
	}
	points := make([]dto.ProviderHallSparkPoint, 0, len(r.Sparkline.Points))
	for _, p := range r.Sparkline.Points {
		points = append(points, dto.ProviderHallSparkPoint{T: p.T, RealTTFTMs: p.RealTTFTMs, ProbeMs: p.ProbeMs, Gap: p.Gap})
	}
	return dto.ProviderHallRow{
		GroupID:      r.GroupID,
		Name:         r.Name,
		Description:  r.Description,
		DisplayOrder: r.DisplayOrder,
		Rate:         r.Rate,
		Quote: dto.ProviderHallQuote{
			InputPrice: r.Quote.InputPrice, CachePrice: r.Quote.CachePrice, Unit: r.Quote.Unit,
			Applicable: r.Quote.Applicable, ReasonCode: r.Quote.ReasonCode,
		},
		HistoricalPrice: providerHallMetricDTO(r.HistoricalPrice),
		PredictedRate:   providerHallMetricDTO(r.PredictedRate),
		Health:          dto.ProviderHallHealth{Status: r.Health.Status, CheckedAt: r.Health.CheckedAt, Models: models},
		Metrics: dto.ProviderHallRowMetrics{
			TTFTFast95Ms: providerHallMetricDTO(r.TTFTFast95Ms),
			TTFTP90Ms:    providerHallMetricDTO(r.TTFTP90Ms),
			CacheRate:    providerHallMetricDTO(r.CacheRate),
			SuccessRate:  providerHallMetricDTO(r.SuccessRate),
		},
		Sparkline:      dto.ProviderHallSparkline{Range: r.Sparkline.Range, BucketSeconds: r.Sparkline.BucketSeconds, Points: points},
		Verification:   providerHallBadgeDTO(r.Verification),
		DefaultProfile: providerHallRefDTO(r.DefaultProfile),
	}
}

func providerHallListDTO(r *service.ProviderHallListResult) dto.ProviderHallListResponse {
	out := dto.ProviderHallListResponse{
		SnapshotID:    r.SnapshotID,
		DataThrough:   r.DataThrough,
		MetricVersion: r.MetricVersion,
		PricingAt:     r.PricingAt,
		Range:         r.Range,
		Catalog:       dto.ProviderHallCatalog{Models: make([]dto.ProviderHallProfileRef, 0, len(r.Catalog.Models)), Default: providerHallRefDTO(r.Catalog.Default)},
		Summary:       dto.ProviderHallSummary{Available: r.Summary.Available, Listed: r.Summary.Listed, Abnormal: r.Summary.Abnormal, Verified: r.Summary.Verified},
		Items:         make([]dto.ProviderHallRow, 0, len(r.Items)),
		Pagination:    dto.ProviderHallPagination{Page: r.Page, PageSize: r.PageSize, Total: r.Total},
	}
	for _, ref := range r.Catalog.Models {
		out.Catalog.Models = append(out.Catalog.Models, providerHallRefDTO(ref))
	}
	for _, row := range r.Items {
		out.Items = append(out.Items, providerHallRowDTO(row))
	}
	return out
}

func providerHallDetailDTO(d *service.ProviderHallDetailResult) dto.ProviderHallDetailResponse {
	out := dto.ProviderHallDetailResponse{
		ProviderHallRow: providerHallRowDTO(d.Row),
		Detail: dto.ProviderHallDetailMetrics{
			E2EAvailability6h: providerHallMetricDTO(d.E2EAvailability6h),
			TPS:               providerHallMetricDTO(d.TPS),
			ProbeTTFTMs:       providerHallMetricDTO(d.ProbeTTFTMs),
			ProbeTotalMs:      providerHallMetricDTO(d.ProbeTotalMs),
			P90Ms:             providerHallMetricDTO(d.P90Ms),
			LastProbeAt:       d.LastProbeAt,
		},
		Trend:    dto.ProviderHallTrend{Range: d.TrendRange, BucketSeconds: d.TrendBucket, Points: make([]dto.ProviderHallTrendPoint, 0, len(d.Trend))},
		Profiles: make([]dto.ProviderHallProfileStatus, 0, len(d.Profiles)),
	}
	if d.ProbeTokens != nil {
		out.Detail.ProbeTokens = &dto.ProviderHallProbeTokens{Input: d.ProbeTokens.Input, Output: d.ProbeTokens.Output, At: d.ProbeTokens.At}
	}
	for _, p := range d.Trend {
		out.Trend.Points = append(out.Trend.Points, dto.ProviderHallTrendPoint{T: p.T, RealFast95Ms: p.RealFast95Ms, ProbeTotalMs: p.ProbeTotalMs, ProbeTTFTMs: p.ProbeTTFTMs, Gap: p.Gap, ProbeFailed: p.ProbeFailed})
	}
	for _, p := range d.Profiles {
		out.Profiles = append(out.Profiles, dto.ProviderHallProfileStatus{ProviderHallProfileRef: providerHallRefDTO(p.ProviderHallProfileRef), Health: providerHallModelHealthDTO(p.Health), Verification: providerHallBadgeDTO(p.Verification)})
	}
	return out
}

func providerHallReportDTO(r service.ProviderHallReportResult) dto.ProviderHallVerificationReport {
	summary, err := json.Marshal(r.Summary)
	if err != nil {
		summary = []byte("{}")
	}
	return dto.ProviderHallVerificationReport{
		ReportID: r.JobID, ProfileID: r.ProfileID, Model: r.Model, Protocol: r.Protocol,
		Verdict: r.Verdict, ExecutionStatus: r.ExecutionStatus, ReasonCode: r.ReasonCode,
		CompletedAt: r.CompletedAt, ExpiresAt: r.ExpiresAt, Expired: r.Expired, Summary: summary,
	}
}
