package admin

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type ProviderHallHandler struct {
	service    *service.ProviderHallService
	runner     providerHallJobRunner
	jobs       providerHallJobReads
	queue      providerHallQueueSource
	aggregator providerHallAggregatorSource
	admin      service.ProviderHallAdminRepository
	now        func() time.Time
}

// NewProviderHallHandler takes the concrete runtime pieces wire provides. Any
// of them may be nil (api-only instance, partial build); the task endpoints
// then answer with the not-ready business error instead of panicking.
func NewProviderHallHandler(svc *service.ProviderHallService, runner *service.ProviderHallRunner, jobs service.ProviderHallJobRepository, collector *service.ProviderHallCollector, aggregator *service.ProviderHallAggregator, adminRepo service.ProviderHallAdminRepository) *ProviderHallHandler {
	h := &ProviderHallHandler{service: svc, admin: adminRepo, now: time.Now}
	if runner != nil {
		h.runner = runner
	}
	if jobs != nil {
		h.jobs = jobs
	}
	if collector != nil {
		h.queue = collector
	}
	if aggregator != nil {
		h.aggregator = aggregator
	}
	return h
}

func providerHallAdmin(c *gin.Context) (int64, bool) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok || subject.UserID < 1 {
		response.Unauthorized(c, "User not authenticated")
		return 0, false
	}
	role, ok := middleware.GetUserRoleFromContext(c)
	if !ok || role != service.RoleAdmin {
		response.Forbidden(c, "Administrator access required")
		return 0, false
	}
	return subject.UserID, true
}

func providerHallID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		response.BadRequest(c, "Invalid ID")
		return 0, false
	}
	return id, true
}

func providerHallBody(c *gin.Context, target any) bool {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 64<<10)
	d := json.NewDecoder(c.Request.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(target); err != nil {
		response.BadRequest(c, "Invalid provider hall request body")
		return false
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		response.BadRequest(c, "Expected one JSON object")
		return false
	}
	return true
}

// providerHallConfigView adds the build's readiness flags to the config so the
// admin page can enable exactly the switches this deployment supports.
func (h *ProviderHallHandler) providerHallConfigView(cfg *service.ProviderHallConfig) any {
	if cfg == nil {
		return nil
	}
	ready := h.service.Readiness()
	return dto.ProviderHallConfigView{
		ProviderHallConfig: cfg,
		Readiness:          dto.ProviderHallReadiness{Collection: ready.Collection, Tasks: ready.Tasks, Display: ready.Display},
	}
}

func (h *ProviderHallHandler) GetConfig(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	result, err := h.service.GetConfig(c.Request.Context())
	if !response.ErrorFrom(c, err) {
		response.Success(c, h.providerHallConfigView(result))
	}
}

func (h *ProviderHallHandler) UpdateConfig(c *gin.Context) {
	actor, ok := providerHallAdmin(c)
	if !ok {
		return
	}
	var input dto.ProviderHallConfigInput
	if !providerHallBody(c, &input) {
		return
	}
	if input.Version == nil {
		response.BadRequest(c, "version is required")
		return
	}
	result, err := h.service.UpdateConfig(c.Request.Context(), service.ProviderHallConfig{
		ProviderHallVersion: service.ProviderHallVersion{Version: *input.Version},
		CollectionEnabled:   input.CollectionEnabled, DisplayEnabled: input.DisplayEnabled, TasksEnabled: input.TasksEnabled,
		DefaultModel: input.DefaultModel, DefaultProtocol: input.DefaultProtocol, DefaultRange: input.DefaultRange,
		GatewayOrigin: input.GatewayOrigin, OperatorUserID: input.OperatorUserID, DailyBudget: input.DailyBudget, ExpectedNodes: input.ExpectedNodes,
	}, actor)
	if !response.ErrorFrom(c, err) {
		response.Success(c, h.providerHallConfigView(result))
	}
}

func (h *ProviderHallHandler) GetGroup(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	result, err := h.service.GetGroup(c.Request.Context(), id)
	if !response.ErrorFrom(c, err) {
		response.Success(c, result)
	}
}

func (h *ProviderHallHandler) UpdateGroup(c *gin.Context) {
	actor, ok := providerHallAdmin(c)
	if !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	var input dto.ProviderHallGroupInput
	if !providerHallBody(c, &input) {
		return
	}
	if input.Version == nil {
		response.BadRequest(c, "version is required (0 for first configuration)")
		return
	}
	result, err := h.service.SaveGroup(c.Request.Context(), service.ProviderHallGroup{
		ProviderHallVersion: service.ProviderHallVersion{Version: *input.Version},
		GroupID:             id, Listed: input.Listed, DisplayName: input.DisplayName, Description: input.Description, DisplayOrder: input.DisplayOrder,
	}, actor)
	if !response.ErrorFrom(c, err) {
		response.Success(c, result)
	}
}

func (h *ProviderHallHandler) ListProfiles(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	result, err := h.service.ListProfiles(c.Request.Context())
	if !response.ErrorFrom(c, err) {
		response.Success(c, result)
	}
}

func (h *ProviderHallHandler) GetTargets(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	result, err := h.service.GetTargets(c.Request.Context(), id)
	if !response.ErrorFrom(c, err) {
		response.Success(c, result)
	}
}

func (h *ProviderHallHandler) UpdateTargets(c *gin.Context) {
	actor, ok := providerHallAdmin(c)
	if !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	var input dto.ProviderHallTargetSetInput
	if !providerHallBody(c, &input) {
		return
	}
	if input.Version == nil || input.Items == nil {
		response.BadRequest(c, "version and items are required; use an empty items array to disable all targets")
		return
	}
	items := make([]service.ProviderHallTarget, 0, len(*input.Items))
	for _, item := range *input.Items {
		items = append(items, service.ProviderHallTarget{ProfileID: item.ProfileID, ProbeKeyID: item.ProbeKeyID, Enabled: item.Enabled,
			ProbeIntervalSeconds: item.ProbeIntervalSeconds, VerificationIntervalSeconds: item.VerificationIntervalSeconds})
	}
	result, err := h.service.SaveTargets(c.Request.Context(), service.ProviderHallTargetSet{
		ProviderHallVersion: service.ProviderHallVersion{Version: *input.Version}, GroupID: id, Items: items,
	}, actor)
	if !response.ErrorFrom(c, err) {
		response.Success(c, result)
	}
}

func (h *ProviderHallHandler) CreateProfile(c *gin.Context) { h.saveProfile(c, true) }
func (h *ProviderHallHandler) UpdateProfile(c *gin.Context) { h.saveProfile(c, false) }

func (h *ProviderHallHandler) saveProfile(c *gin.Context, create bool) {
	actor, ok := providerHallAdmin(c)
	if !ok {
		return
	}
	var id int64
	if !create {
		id, ok = providerHallID(c)
		if !ok {
			return
		}
	}
	var input dto.ProviderHallProfileInput
	if !providerHallBody(c, &input) {
		return
	}
	if input.Version == nil {
		response.BadRequest(c, "version is required (0 to create)")
		return
	}
	result, err := h.service.SaveProfile(c.Request.Context(), service.ProviderHallProfile{
		ProviderHallVersion: service.ProviderHallVersion{Version: *input.Version},
		ID:                  id, Model: input.Model, Protocol: input.Protocol, SupportsTools: input.SupportsTools, OutputLimit: input.OutputLimit,
		ModelAliases: input.ModelAliases, ReferenceInputPrice: input.ReferenceInputPrice, ReferenceCachePrice: input.ReferenceCachePrice,
		ReferenceCacheRate: input.ReferenceCacheRate, ReferenceConfirmedAt: input.ReferenceConfirmedAt,
	}, actor)
	if response.ErrorFrom(c, err) {
		return
	}
	if create {
		response.Created(c, result)
	} else {
		response.Success(c, result)
	}
}
