//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type providerHallTargetRepoStub struct {
	ProviderHallRepository
	saved *ProviderHallTargetSet
}

func (r *providerHallTargetRepoStub) SaveTargets(_ context.Context, input ProviderHallTargetSet) (*ProviderHallTargetSet, error) {
	r.saved = &input
	return &input, nil
}

func TestProviderHallTargetSetValidation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input ProviderHallTargetSet
	}{
		{"group", ProviderHallTargetSet{GroupID: -1}},
		{"version", ProviderHallTargetSet{GroupID: 1, ProviderHallVersion: ProviderHallVersion{Version: -1}}},
		{"duplicate_profile", ProviderHallTargetSet{GroupID: 1, Items: []ProviderHallTarget{{ProfileID: 2}, {ProfileID: 2}}}},
		{"missing_key", ProviderHallTargetSet{GroupID: 1, Items: []ProviderHallTarget{{ProfileID: 2, Enabled: true}}}},
		{"interval_low", ProviderHallTargetSet{GroupID: 1, Items: []ProviderHallTarget{{ProfileID: 2, ProbeIntervalSeconds: 59}}}},
		{"interval_high", ProviderHallTargetSet{GroupID: 1, Items: []ProviderHallTarget{{ProfileID: 2, VerificationIntervalSeconds: 604801}}}},
		{"too_many", ProviderHallTargetSet{GroupID: 1, Items: make([]ProviderHallTarget, 101)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &providerHallTargetRepoStub{}
			svc := NewProviderHallService(repo, nil, nil, nil)
			_, err := svc.SaveTargets(context.Background(), tc.input, 7)
			require.Error(t, err)
			require.Nil(t, repo.saved)
		})
	}
	repo := &providerHallTargetRepoStub{}
	svc := NewProviderHallService(repo, nil, nil, nil)
	input := ProviderHallTargetSet{GroupID: 4, Items: []ProviderHallTarget{{ID: 99, GroupID: 123, ProfileID: 8}, {ProfileID: 3}}}
	result, err := svc.SaveTargets(context.Background(), input, 7)
	require.NoError(t, err)
	require.Equal(t, int64(7), *result.UpdatedBy)
	require.Equal(t, int64(3), result.Items[0].ProfileID)
	require.Equal(t, 300, result.Items[0].ProbeIntervalSeconds)
	require.Equal(t, 86400, result.Items[0].VerificationIntervalSeconds)
	require.Zero(t, result.Items[1].ID)
	require.Equal(t, int64(4), result.Items[1].GroupID)
	require.Equal(t, int64(99), input.Items[0].ID, "normalization must not mutate caller input")
}

func TestProviderHallTargetBinding(t *testing.T) {
	now := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name string
		edit func(*Group, *ProviderHallProfile, *APIKey, *User, *bool, *CompositeRouteDecision)
		want error
	}{
		{"public", func(*Group, *ProviderHallProfile, *APIKey, *User, *bool, *CompositeRouteDecision) {}, nil},
		{"other_owner", func(_ *Group, _ *ProviderHallProfile, k *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			k.UserID = 99
		}, ErrProviderHallProbeKey},
		{"wrong_group", func(_ *Group, _ *ProviderHallProfile, k *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			k.GroupID = new(int64(99))
		}, ErrProviderHallProbeKey},
		{"no_group", func(_ *Group, _ *ProviderHallProfile, k *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			k.GroupID = nil
		}, ErrProviderHallProbeKey},
		{"disabled", func(_ *Group, _ *ProviderHallProfile, k *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			k.Status = StatusAPIKeyDisabled
		}, ErrProviderHallProbeKey},
		{"expired_exactly_now", func(_ *Group, _ *ProviderHallProfile, k *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			k.ExpiresAt = &now
		}, ErrProviderHallProbeKey},
		{"quota_exhausted", func(_ *Group, _ *ProviderHallProfile, k *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			k.Quota, k.QuotaUsed = 1, 1
		}, ErrProviderHallProbeKey},
		{"exclusive", func(g *Group, _ *ProviderHallProfile, _ *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			g.IsExclusive = true
		}, ErrProviderHallOperator},
		{"exclusive_allowed", func(g *Group, _ *ProviderHallProfile, _ *APIKey, u *User, _ *bool, _ *CompositeRouteDecision) {
			g.IsExclusive = true
			u.AllowedGroups = []int64{g.ID}
		}, nil},
		{"subscription_expired", func(g *Group, _ *ProviderHallProfile, _ *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			g.SubscriptionType = SubscriptionTypeSubscription
		}, ErrProviderHallOperator},
		{"subscription_active", func(g *Group, _ *ProviderHallProfile, _ *APIKey, _ *User, sub *bool, _ *CompositeRouteDecision) {
			g.SubscriptionType = SubscriptionTypeSubscription
			*sub = true
		}, nil},
		{"operator_disabled", func(_ *Group, _ *ProviderHallProfile, _ *APIKey, u *User, _ *bool, _ *CompositeRouteDecision) {
			u.Status = "disabled"
		}, ErrProviderHallOperator},
		{"operator_deleted", func(_ *Group, _ *ProviderHallProfile, _ *APIKey, u *User, _ *bool, _ *CompositeRouteDecision) {
			u.DeletedAt = &now
		}, ErrProviderHallOperator},
		{"messages_disabled", func(_ *Group, p *ProviderHallProfile, _ *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			p.Protocol = "messages"
		}, ErrProviderHallTargetRoute},
		{"messages_enabled", func(g *Group, p *ProviderHallProfile, _ *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			p.Protocol = "messages"
			g.AllowMessagesDispatch = true
		}, nil},
		{"composite_openai", func(g *Group, _ *ProviderHallProfile, _ *APIKey, _ *User, _ *bool, r *CompositeRouteDecision) {
			g.Platform = PlatformComposite
			r.Matched = true
			r.TargetPlatform = PlatformOpenAI
		}, nil},
		{"composite_other_provider", func(g *Group, _ *ProviderHallProfile, _ *APIKey, _ *User, _ *bool, r *CompositeRouteDecision) {
			g.Platform = PlatformComposite
			r.Matched = true
			r.TargetPlatform = PlatformDeepseek
		}, ErrProviderHallTargetRoute},
		{"composite_unresolved", func(g *Group, _ *ProviderHallProfile, _ *APIKey, _ *User, _ *bool, _ *CompositeRouteDecision) {
			g.Platform = PlatformComposite
		}, ErrProviderHallTargetRoute},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := Group{ID: 2, Platform: PlatformOpenAI, Status: StatusActive}
			p := ProviderHallProfile{Model: "gpt-test", Protocol: "responses"}
			k := APIKey{ID: 3, UserID: 1, GroupID: &g.ID, Status: StatusActive}
			u := User{ID: 1, Status: StatusActive}
			sub, route := false, CompositeRouteDecision{}
			tc.edit(&g, &p, &k, &u, &sub, &route)
			err := ValidateProviderHallTargetBinding(&g, &p, &k, &u, sub, route, now)
			if tc.want == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.want)
			}
		})
	}
}
