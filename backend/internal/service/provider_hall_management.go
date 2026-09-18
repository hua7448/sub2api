package service

import (
	"context"
	"strings"
	"time"
)

type ProviderHallModelRefresh struct {
	AccountID int64                        `json:"account_id"`
	Success   bool                         `json:"success"`
	Error     string                       `json:"error,omitempty"`
	Models    []ProviderHallModelCandidate `json:"models"`
}

func (s *ProviderHallService) RefreshModels(ctx context.Context, groupID int64) ([]ProviderHallModelRefresh, error) {
	if s.accounts == nil || s.modelFetcher == nil {
		return nil, ErrProviderHallNotReady
	}
	if _, err := s.GetGroup(ctx, groupID); err != nil {
		return nil, err
	}
	accounts, err := s.accounts.ListByGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	g, err := s.groups.GetByIDLite(ctx, groupID)
	if err != nil {
		return nil, err
	}
	out := []ProviderHallModelRefresh{}
	for _, account := range accounts {
		if account.Platform != PlatformOpenAI {
			continue
		}
		item := ProviderHallModelRefresh{AccountID: account.ID, Models: []ProviderHallModelCandidate{}}
		requestCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		models, e := s.modelFetcher.FetchUpstreamSupportedModels(requestCtx, &account)
		cancel()
		if e != nil {
			item.Error = "upstream_models_unavailable"
		} else {
			item.Success = true
			for _, upstream := range models {
				if strings.ContainsAny(upstream, "*?\r\n") {
					continue
				}
				public := []string{}
				for model, mapped := range account.GetModelMapping() {
					if mapped == upstream && !strings.ContainsAny(model, "*?") {
						public = append(public, model)
					}
				}
				if len(public) == 0 {
					public = append(public, upstream)
				}
				for _, model := range public {
					for _, protocol := range []string{"responses", "chat_completions", "messages"} {
						routedModel := model
						available := g.IsActive() && account.IsActive() && (protocol != "messages" || g.AllowMessagesDispatch)
						if g.Platform == PlatformComposite {
							route, e := NewCompositeRouteResolver(s.modelRoutes).Resolve(ctx, groupID, model, protocol)
							if e != nil {
								return nil, e
							}
							available = available && route.Matched && route.TargetPlatform == PlatformOpenAI
							routedModel = route.UpstreamModel
						}
						available = available && account.IsModelSupported(routedModel) && account.GetMappedModel(routedModel) == upstream
						reason := ""
						if !available {
							reason = "route_or_account_unavailable"
						}
						item.Models = append(item.Models, ProviderHallModelCandidate{Model: model, UpstreamModel: upstream, Protocol: protocol, Source: "upstream", AccountID: account.ID, Available: available, Reason: reason})
					}
				}
			}
		}
		out = append(out, item)
		if ctx.Err() != nil {
			break
		}
	}
	return out, nil
}

type providerHallExpectedVersionKey struct{}
type ProviderHallExpectedVersions struct{ Target, Profile int64 }

func ProviderHallWithExpectedVersions(ctx context.Context, target, profile int64) context.Context {
	return context.WithValue(ctx, providerHallExpectedVersionKey{}, ProviderHallExpectedVersions{target, profile})
}

type ProviderHallGroupSummary struct {
	KeyStatuses        map[string]int           `json:"key_statuses"`
	LatestProbe        *ProviderHallAdminResult `json:"latest_probe"`
	LatestVerification *ProviderHallAdminResult `json:"latest_verification"`
	ProviderHallGroup
	Name               string               `json:"name"`
	Platform           string               `json:"platform"`
	Status             string               `json:"status"`
	Supported          bool                 `json:"supported"`
	Reason             string               `json:"reason"`
	Targets            []ProviderHallTarget `json:"targets"`
	Models             []string             `json:"models"`
	EffectiveProfileID *int64               `json:"effective_profile_id"`
}

type ProviderHallAdminResult struct {
	JobID   int64     `json:"job_id"`
	Status  string    `json:"status"`
	Verdict string    `json:"verdict"`
	At      time.Time `json:"at"`
}

type ProviderHallGroupFilter struct {
	Search, Platform, Listed, Sort string
	Page, PageSize                 int
}
type ProviderHallGroupPage struct {
	Items    []ProviderHallGroupSummary `json:"items"`
	Total    int                        `json:"total"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"page_size"`
}
type ProviderHallProbeKeyOption struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Status         string `json:"status"`
	Registered     bool   `json:"registered"`
	GroupID        int64  `json:"group_id"`
	OperatorUserID int64  `json:"operator_user_id"`
}
type ProviderHallModelCandidate struct {
	Model         string `json:"model"`
	Protocol      string `json:"protocol"`
	Source        string `json:"source"`
	AccountID     int64  `json:"account_id"`
	UpstreamModel string `json:"upstream_model"`
	Available     bool   `json:"available"`
	Reason        string `json:"reason"`
}
type ProviderHallManagementRepository interface {
	ListAdminGroups(context.Context, ProviderHallGroupFilter) (*ProviderHallGroupPage, error)
	ListProbeKeys(context.Context, int64) ([]ProviderHallProbeKeyOption, error)
	EnsureProbeKey(context.Context, int64, int64, int64) (*ProviderHallProbeKeyOption, error)
	ListModelCandidates(context.Context, int64) ([]ProviderHallModelCandidate, error)
}

func (s *ProviderHallService) management() (ProviderHallManagementRepository, error) {
	r, ok := s.repo.(ProviderHallManagementRepository)
	if !ok {
		return nil, ErrProviderHallNotReady
	}
	return r, nil
}

func (s *ProviderHallService) ListAdminGroups(ctx context.Context, filter ProviderHallGroupFilter) (*ProviderHallGroupPage, error) {
	r, err := s.management()
	if err != nil {
		return nil, err
	}
	return r.ListAdminGroups(ctx, filter)
}
func (s *ProviderHallService) ListModelCandidates(ctx context.Context, groupID int64) ([]ProviderHallModelCandidate, error) {
	r, err := s.management()
	if err != nil {
		return nil, err
	}
	return r.ListModelCandidates(ctx, groupID)
}
func (s *ProviderHallService) ListProbeKeys(ctx context.Context, groupID int64) ([]ProviderHallProbeKeyOption, error) {
	r, err := s.management()
	if err != nil {
		return nil, err
	}
	return r.ListProbeKeys(ctx, groupID)
}
func (s *ProviderHallService) EnsureProbeKey(ctx context.Context, groupID, profileID, actorID int64) (*ProviderHallProbeKeyOption, error) {
	if groupID < 1 || profileID < 1 || actorID < 1 {
		return nil, providerHallInvalid("profile_id")
	}
	r, err := s.management()
	if err != nil {
		return nil, err
	}
	return r.EnsureProbeKey(ctx, groupID, profileID, actorID)
}
