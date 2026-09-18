package admin

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *ProviderHallHandler) PreflightConfig(c *gin.Context) {
	actor, ok := providerHallAdmin(c)
	if !ok {
		return
	}
	var i dto.ProviderHallConfigInput
	if !providerHallBody(c, &i) {
		return
	}
	if i.Version == nil {
		response.BadRequest(c, "version is required")
		return
	}
	_, err := h.service.ValidateConfig(c.Request.Context(), service.ProviderHallConfig{ProviderHallVersion: service.ProviderHallVersion{Version: *i.Version}, CollectionEnabled: i.CollectionEnabled, DisplayEnabled: i.DisplayEnabled, TasksEnabled: i.TasksEnabled, AutoScheduleEnabled: i.AutoScheduleEnabled, DefaultModel: i.DefaultModel, DefaultProtocol: i.DefaultProtocol, DefaultRange: i.DefaultRange, GatewayOrigin: i.GatewayOrigin, OperatorUserID: i.OperatorUserID, DailyBudget: i.DailyBudget, ExpectedNodes: i.ExpectedNodes}, actor)
	if response.ErrorFrom(c, err) {
		return
	}
	current, err := h.service.GetConfig(c.Request.Context())
	if response.ErrorFrom(c, err) {
		return
	}
	if current.Version != *i.Version {
		response.ErrorFrom(c, service.ErrProviderHallConflict)
		return
	}
	affected := []service.ProviderHallTarget{}
	changed := (current.OperatorUserID == nil) != (i.OperatorUserID == nil) || current.OperatorUserID != nil && i.OperatorUserID != nil && *current.OperatorUserID != *i.OperatorUserID
	if changed {
		for page := 1; ; page++ {
			gs, e := h.service.ListAdminGroups(c.Request.Context(), service.ProviderHallGroupFilter{Page: page, PageSize: 100})
			if response.ErrorFrom(c, e) {
				return
			}
			for _, g := range gs.Items {
				for _, t := range g.Targets {
					if t.Enabled {
						affected = append(affected, t)
					}
				}
			}
			if page*100 >= gs.Total {
				break
			}
		}
	}
	response.Success(c, gin.H{"valid": true, "disabled_targets": affected})
}

func (h *ProviderHallHandler) CheckGateway(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	var i struct {
		Origin string `json:"origin"`
	}
	if !providerHallBody(c, &i) {
		return
	}
	status, err := service.ProviderHallCheckGateway(c.Request.Context(), i.Origin)
	if !response.ErrorFrom(c, err) {
		response.Success(c, gin.H{"status": status, "healthy": status == 200})
	}
}

func (h *ProviderHallHandler) ListAdminGroups(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("page_size"))
	if page > 1000000 || len(c.Query("search")) > 200 || c.Query("listed") != "" && c.Query("listed") != "true" && c.Query("listed") != "false" {
		response.BadRequest(c, "invalid group filter")
		return
	}
	out, err := h.service.ListAdminGroups(c.Request.Context(), service.ProviderHallGroupFilter{Search: c.Query("search"), Platform: c.Query("platform"), Listed: c.Query("listed"), Sort: c.Query("sort"), Page: page, PageSize: size})
	if !response.ErrorFrom(c, err) {
		response.Success(c, out)
	}
}
func (h *ProviderHallHandler) ListGroupModels(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	out, err := h.service.ListModelCandidates(c.Request.Context(), id)
	if !response.ErrorFrom(c, err) {
		response.Success(c, out)
	}
}

// ListAllProfileCandidates serves the profile editor: one merged candidate per
// model and protocol across every group, so the admin picks a model once
// instead of re-picking a source group.
func (h *ProviderHallHandler) ListAllProfileCandidates(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	out, err := h.service.ListAllModelCandidates(c.Request.Context())
	if !response.ErrorFrom(c, err) {
		response.Success(c, out)
	}
}

func (h *ProviderHallHandler) DeleteProfile(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	if err := h.service.DeleteProfile(c.Request.Context(), id); !response.ErrorFrom(c, err) {
		response.Success(c, map[string]int64{"profile_id": id})
	}
}

func (h *ProviderHallHandler) RefreshGroupModels(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	out, err := h.service.RefreshModels(c.Request.Context(), id)
	if !response.ErrorFrom(c, err) {
		response.Success(c, out)
	}
}
func (h *ProviderHallHandler) ListProbeKeys(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	out, err := h.service.ListProbeKeys(c.Request.Context(), id)
	if !response.ErrorFrom(c, err) {
		response.Success(c, out)
	}
}
func (h *ProviderHallHandler) EnsureProbeKey(c *gin.Context) {
	actor, ok := providerHallAdmin(c)
	if !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	var input struct {
		ProfileID int64 `json:"profile_id"`
	}
	if !providerHallBody(c, &input) {
		return
	}
	if input.ProfileID < 1 {
		response.BadRequest(c, "profile_id is required")
		return
	}
	out, err := h.service.EnsureProbeKey(c.Request.Context(), id, input.ProfileID, actor)
	if !response.ErrorFrom(c, err) {
		response.Success(c, out)
	}
}

type providerHallSettingsInput struct {
	dto.ProviderHallGroupInput
	Items *[]dto.ProviderHallTargetInput `json:"items"`
}

func (h *ProviderHallHandler) SaveSettings(c *gin.Context)   { h.saveSettings(c, false) }
func (h *ProviderHallHandler) PreflightGroup(c *gin.Context) { h.saveSettings(c, true) }
func (h *ProviderHallHandler) saveSettings(c *gin.Context, preflight bool) {
	actor, ok := providerHallAdmin(c)
	if !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	var input providerHallSettingsInput
	if !providerHallBody(c, &input) {
		return
	}
	if input.Version == nil || input.Items == nil {
		response.BadRequest(c, "version and items are required")
		return
	}
	out, err := h.service.SaveTargets(c.Request.Context(), providerHallSettings(id, input, preflight), actor)
	if !response.ErrorFrom(c, err) {
		response.Success(c, out)
	}
}
func providerHallSettings(id int64, input providerHallSettingsInput, preflight bool) service.ProviderHallTargetSet {
	items := []service.ProviderHallTarget{}
	for _, i := range *input.Items {
		items = append(items, service.ProviderHallTarget{ProfileID: i.ProfileID, ProbeKeyID: i.ProbeKeyID, Enabled: i.Enabled, AutoScheduleEnabled: i.AutoScheduleEnabled, ProbeIntervalSeconds: i.ProbeIntervalSeconds, VerificationIntervalSeconds: i.VerificationIntervalSeconds})
	}
	return service.ProviderHallTargetSet{ProviderHallVersion: service.ProviderHallVersion{Version: *input.Version}, GroupID: id, Items: items, Preflight: preflight,
		Listing: &service.ProviderHallGroup{GroupID: id, Listed: input.Listed, DisplayName: input.DisplayName, Description: input.Description, DisplayOrder: input.DisplayOrder}}
}

func (h *ProviderHallHandler) BatchGroups(c *gin.Context) {
	actor, ok := providerHallAdmin(c)
	if !ok {
		return
	}
	var input struct {
		Preview bool `json:"preview"`
		Groups  []struct {
			ID int64 `json:"id"`
			providerHallSettingsInput
		} `json:"groups"`
	}
	// A batch can contain up to 100 explicit target sets.
	if !providerHallBodyLimit(c, &input, 4<<20) {
		return
	}
	if len(input.Groups) == 0 || len(input.Groups) > 100 {
		response.BadRequest(c, "select 1 to 100 groups")
		return
	}
	seen := map[int64]bool{}
	for _, g := range input.Groups {
		if g.ID < 1 || g.Version == nil || g.Items == nil || seen[g.ID] {
			response.BadRequest(c, "unique groups, versions and explicit targets are required")
			return
		}
		seen[g.ID] = true
	}
	type result struct {
		ID       int64                          `json:"id"`
		Success  bool                           `json:"success"`
		Error    string                         `json:"error,omitempty"`
		Reason   string                         `json:"reason,omitempty"`
		Metadata map[string]string              `json:"metadata,omitempty"`
		Settings *service.ProviderHallTargetSet `json:"settings,omitempty"`
	}
	results := []result{}
	for _, g := range input.Groups {
		out, err := h.service.SaveTargets(c.Request.Context(), providerHallSettings(g.ID, g.providerHallSettingsInput, input.Preview), actor)
		item := result{ID: g.ID, Success: err == nil, Settings: out}
		if err != nil {
			item.Error = infraerrors.Message(err)
			item.Reason = infraerrors.Reason(err)
			item.Metadata = infraerrors.FromError(err).Metadata
		}
		results = append(results, item)
	}
	response.Success(c, results)
}
