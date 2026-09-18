package service

import (
	"context"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var (
	ErrProviderHallOperator    = infraerrors.BadRequest("PROVIDER_HALL_OPERATOR_NOT_CONFIGURED", "an active operator with group access is required")
	ErrProviderHallProbeKey    = infraerrors.BadRequest("PROVIDER_HALL_PROBE_KEY_INVALID", "probe key must be active, dedicated to the operator and bound to this group")
	ErrProviderHallTargetRoute = infraerrors.BadRequest("PROVIDER_HALL_TARGET_ROUTE_UNSUPPORTED", "target model and protocol must route to OpenAI")
)

const ProviderHallMaxTargetsPerGroup = 100

type ProviderHallTarget struct {
	ProviderHallVersion
	ID                          int64  `json:"id"`
	GroupID                     int64  `json:"group_id"`
	ProfileID                   int64  `json:"profile_id"`
	ProbeKeyID                  *int64 `json:"probe_key_id"`
	Enabled                     bool   `json:"enabled"`
	AutoScheduleEnabled         *bool  `json:"auto_schedule_enabled"`
	ProbeIntervalSeconds        int    `json:"probe_interval_seconds"`
	VerificationIntervalSeconds int    `json:"verification_interval_seconds"`
}

// The set uses the hall group's version, shared with name/listing writes.
// Omitted targets are retained disabled so their IDs remain stable for reports.
type ProviderHallTargetSet struct {
	ProviderHallVersion
	GroupID   int64                `json:"group_id"`
	Items     []ProviderHallTarget `json:"items"`
	Listing   *ProviderHallGroup   `json:"listing,omitempty"`
	Preflight bool                 `json:"-"`
}

func (s *ProviderHallService) GetTargets(ctx context.Context, groupID int64) (*ProviderHallTargetSet, error) {
	if groupID < 1 {
		return nil, providerHallInvalid("group_id")
	}
	return s.repo.GetTargets(ctx, groupID)
}

func (s *ProviderHallService) SaveTargets(ctx context.Context, input ProviderHallTargetSet, actorID int64) (*ProviderHallTargetSet, error) {
	if input.GroupID < 1 || input.Version < 0 || actorID < 1 {
		return nil, providerHallInvalid("version")
	}
	if len(input.Items) > ProviderHallMaxTargetsPerGroup {
		return nil, providerHallInvalid("items")
	}
	if input.Listing != nil {
		g := input.Listing
		if utf8.RuneCountInString(g.DisplayName) > 100 || utf8.RuneCountInString(g.Description) > 2000 || g.DisplayOrder < -2147483648 || g.DisplayOrder > 2147483647 {
			return nil, providerHallInvalid("group")
		}
		g.DisplayName = strings.TrimSpace(g.DisplayName)
	}
	seen := map[int64]bool{}
	items := make([]ProviderHallTarget, 0, len(input.Items))
	for _, item := range input.Items {
		if item.ProfileID < 1 || seen[item.ProfileID] {
			return nil, providerHallInvalid("profile_id")
		}
		seen[item.ProfileID] = true
		if item.ProbeIntervalSeconds == 0 {
			item.ProbeIntervalSeconds = 300
		}
		if item.VerificationIntervalSeconds == 0 {
			item.VerificationIntervalSeconds = 86400
		}
		if item.ProbeIntervalSeconds < 60 || item.ProbeIntervalSeconds > 86400 {
			return nil, providerHallInvalid("probe_interval_seconds")
		}
		if item.VerificationIntervalSeconds < 3600 || item.VerificationIntervalSeconds > 604800 {
			return nil, providerHallInvalid("verification_interval_seconds")
		}
		if item.ProbeKeyID != nil && *item.ProbeKeyID < 1 || item.Enabled && item.ProbeKeyID == nil {
			return nil, ErrProviderHallProbeKey
		}
		// Copy writable fields only; IDs and audit metadata are server-owned.
		items = append(items, ProviderHallTarget{GroupID: input.GroupID, ProfileID: item.ProfileID,
			ProbeKeyID: item.ProbeKeyID, Enabled: item.Enabled, AutoScheduleEnabled: item.AutoScheduleEnabled, ProbeIntervalSeconds: item.ProbeIntervalSeconds,
			VerificationIntervalSeconds: item.VerificationIntervalSeconds})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ProfileID < items[j].ProfileID })
	input.Items, input.UpdatedBy = items, &actorID
	var previous map[int64]ProviderHallTarget
	if s.jobs != nil && !input.Preflight {
		before, err := s.repo.GetTargets(ctx, input.GroupID)
		if err != nil {
			return nil, err
		}
		previous = make(map[int64]ProviderHallTarget, len(before.Items))
		for _, t := range before.Items {
			previous[t.ID] = t
		}
	}
	saved, err := s.repo.SaveTargets(ctx, input)
	if err != nil {
		return nil, err
	}
	// Queued jobs of a disabled or changed target must not send. The runner
	// also cancels drifted/disabled jobs each tick and re-checks versions at
	// dispatch, so this hook only shortens the window.
	if s.jobs != nil && !input.Preflight {
		now := time.Now().UTC()
		for _, t := range saved.Items {
			old, existed := previous[t.ID]
			switch {
			case !t.Enabled:
				_, _ = s.jobs.CancelQueuedForTarget(ctx, t.ID, ProviderHallJobCodeTargetDisabled, now)
			case existed && old.Version != t.Version:
				_, _ = s.jobs.CancelQueuedForTarget(ctx, t.ID, ProviderHallJobCodeConfigChanged, now)
			}
		}
	}
	return saved, nil
}

// ValidateProviderHallTargetBinding is also usable by the runner's dispatch
// recheck. Configuration-time validation never authorizes a later paid request.
func ValidateProviderHallTargetBinding(group *Group, profile *ProviderHallProfile, key *APIKey, operator *User, hasSubscription bool, route CompositeRouteDecision, now time.Time) error {
	if group == nil || !group.IsActive() {
		return ErrProviderHallTargetRoute
	}
	if group.Platform != PlatformOpenAI && group.Platform != PlatformComposite {
		return ErrProviderHallTargetRoute
	}
	if profile == nil || !providerHallProtocol(profile.Protocol) {
		return ErrProviderHallTargetRoute
	}
	if group.Platform == PlatformComposite && (!route.Matched || route.TargetPlatform != PlatformOpenAI) {
		return ErrProviderHallTargetRoute
	}
	if profile.Protocol == "messages" && !group.AllowMessagesDispatch {
		return ErrProviderHallTargetRoute
	}
	if operator == nil || !operator.IsActive() || operator.DeletedAt != nil {
		return ErrProviderHallOperator
	}
	if group.IsSubscriptionType() {
		if !hasSubscription {
			return ErrProviderHallOperator
		}
	} else if !operator.CanBindGroup(group.ID, group.IsExclusive) {
		return ErrProviderHallOperator
	}
	if key == nil || !key.IsActive() || key.UserID != operator.ID || key.GroupID == nil || *key.GroupID != group.ID ||
		key.ExpiresAt != nil && !key.ExpiresAt.After(now) || key.IsQuotaExhausted() {
		return ErrProviderHallProbeKey
	}
	return nil
}

// New configurations default to manual-only; nil is used by legacy write clients.
func ProviderHallAutoSchedule(value *bool) bool { return value != nil && *value }
