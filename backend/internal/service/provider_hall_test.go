//go:build unit

package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

type providerHallRepoStub struct {
	ProviderHallRepository
	config  ProviderHallConfig
	listing ProviderHallGroup
	writes  int
}

func (r *providerHallRepoStub) GetConfig(context.Context) (*ProviderHallConfig, error) {
	return &r.config, nil
}
func (r *providerHallRepoStub) UpdateConfig(_ context.Context, cfg ProviderHallConfig) (*ProviderHallConfig, error) {
	if cfg.Version != r.config.Version {
		return nil, ErrProviderHallConflict
	}
	r.writes++
	cfg.Version++
	r.config = cfg
	return &cfg, nil
}
func (r *providerHallRepoStub) GetGroup(context.Context, int64) (*ProviderHallGroup, error) {
	return &r.listing, nil
}
func (r *providerHallRepoStub) SaveGroup(_ context.Context, group ProviderHallGroup) (*ProviderHallGroup, error) {
	r.writes++
	r.listing = group
	return &group, nil
}
func (r *providerHallRepoStub) SaveProfile(_ context.Context, p ProviderHallProfile) (*ProviderHallProfile, error) {
	r.writes++
	return &p, nil
}

type providerHallUsersStub struct {
	UserRepository
	user User
}

func (r *providerHallUsersStub) GetByID(context.Context, int64) (*User, error) { return &r.user, nil }

type providerHallGroupsStub struct {
	GroupRepository
	groups []Group
}

func (r *providerHallGroupsStub) ListActive(context.Context) ([]Group, error) { return r.groups, nil }
func (r *providerHallGroupsStub) GetByIDLite(_ context.Context, id int64) (*Group, error) {
	for _, group := range r.groups {
		if group.ID == id {
			return &group, nil
		}
	}
	return nil, ErrProviderHallNotFound
}

func TestProviderHallGroupChineseLengthLimits(t *testing.T) {
	for _, tc := range []struct {
		name        string
		nameLen     int
		description int
		valid       bool
	}{
		{"previous_byte_boundary", 34, 667, true},
		{"character_limit", 100, 2000, true},
		{"name_over_limit", 101, 2000, false},
		{"description_over_limit", 100, 2001, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &providerHallRepoStub{}
			groups := &providerHallGroupsStub{groups: []Group{{ID: 1, Platform: PlatformOpenAI, Status: StatusActive}}}
			svc := NewProviderHallService(repo, groups, nil, nil)
			input := ProviderHallGroup{GroupID: 1, DisplayName: strings.Repeat("\u4e2d", tc.nameLen), Description: strings.Repeat("\u6587", tc.description)}
			result, err := svc.SaveGroup(context.Background(), input, 1)
			if tc.valid {
				require.NoError(t, err)
				require.Equal(t, input.DisplayName, result.DisplayName)
				require.Equal(t, input.Description, result.Description)
				require.Equal(t, 1, repo.writes)
			} else {
				require.Error(t, err)
				require.Zero(t, repo.writes)
			}
		})
	}
}

type providerHallSubscriptionsStub struct {
	UserSubscriptionRepository
	subscriptions []UserSubscription
}

func (r *providerHallSubscriptionsStub) ListActiveByUserID(context.Context, int64) ([]UserSubscription, error) {
	return r.subscriptions, nil
}

func TestProviderHallConfigurationValidation(t *testing.T) {
	valid := ProviderHallConfig{ProviderHallVersion: ProviderHallVersion{Version: 1}, DefaultProtocol: "responses", DefaultRange: "6h", DailyBudget: "0.00000001"}
	for _, tc := range []struct {
		name string
		edit func(*ProviderHallConfig)
	}{
		{"version", func(c *ProviderHallConfig) { c.Version = 0 }},
		{"negative_budget", func(c *ProviderHallConfig) { c.DailyBudget = "-1" }},
		{"precision", func(c *ProviderHallConfig) { c.DailyBudget = "0.000000001" }},
		{"overflow", func(c *ProviderHallConfig) { c.DailyBudget = "1000000000000" }},
		{"exponent", func(c *ProviderHallConfig) { c.DailyBudget = "1e-8" }},
		{"protocol", func(c *ProviderHallConfig) { c.DefaultProtocol = "websocket" }},
		{"range", func(c *ProviderHallConfig) { c.DefaultRange = "1h" }},
		{"origin_path", func(c *ProviderHallConfig) { c.GatewayOrigin = "https://panel.example/v1" }},
		{"origin_credentials", func(c *ProviderHallConfig) { c.GatewayOrigin = "https://user:secret@panel.example" }},
		{"origin_query", func(c *ProviderHallConfig) { c.GatewayOrigin = "https://panel.example?" }},
		{"origin_fragment", func(c *ProviderHallConfig) { c.GatewayOrigin = "https://panel.example#" }},
		{"origin_http", func(c *ProviderHallConfig) { c.GatewayOrigin = "http://panel.example" }},
		{"duplicate_node", func(c *ProviderHallConfig) { c.ExpectedNodes = []string{"api-1", "api-1"} }},
		{"invalid_node", func(c *ProviderHallConfig) { c.ExpectedNodes = []string{"api 1"} }},
		{"collection_not_ready", func(c *ProviderHallConfig) { c.CollectionEnabled = true }},
		{"display_not_ready", func(c *ProviderHallConfig) { c.DisplayEnabled = true }},
		{"tasks_not_ready", func(c *ProviderHallConfig) { c.TasksEnabled = true }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &providerHallRepoStub{config: valid}
			svc := &ProviderHallService{repo: repo}
			cfg := valid
			tc.edit(&cfg)
			_, err := svc.UpdateConfig(context.Background(), cfg, 4)
			require.Error(t, err)
			require.Zero(t, repo.writes)
		})
	}
	repo := &providerHallRepoStub{config: valid}
	svc := &ProviderHallService{repo: repo}
	saved := 0
	svc.SetConfigSavedCallback(func() { saved++ })
	updated, err := svc.UpdateConfig(context.Background(), valid, 4)
	require.NoError(t, err)
	require.Equal(t, 1, saved, "a successful save must invalidate derived caches")
	require.Equal(t, int64(2), updated.Version)
	require.Equal(t, int64(4), *updated.UpdatedBy)
	require.NotNil(t, updated.ExpectedNodes)
	_, err = svc.UpdateConfig(context.Background(), valid, 4)
	require.ErrorIs(t, err, ErrProviderHallConflict)
	require.Equal(t, 1, saved, "a rejected save must not fire the callback")
}

func TestProviderHallConfigChangeRefreshesPublicSettings(t *testing.T) {
	svc := &SettingService{}
	invalidated := 0
	svc.SetOnUpdateCallback(func() { invalidated++ })
	reader := &providerHallDisplayReaderStub{enabled: false}
	svc.SetProviderHallDisplayReader(reader)
	ctx := context.Background()
	require.False(t, svc.providerHallDisplayEnabled(ctx))
	reader.enabled = true
	require.False(t, svc.providerHallDisplayEnabled(ctx), "cached for 30s without a notification")
	svc.NotifyProviderHallConfigChanged()
	require.True(t, svc.providerHallDisplayEnabled(ctx))
	require.Equal(t, 1, invalidated, "the injected HTML cache must be invalidated")
}

type providerHallDisplayReaderStub struct{ enabled bool }

func (r *providerHallDisplayReaderStub) DisplayEnabled(context.Context) bool { return r.enabled }

func TestProviderHallAuthorizationRechecksSubscriptionAndRevocation(t *testing.T) {
	ctx := context.Background()
	users := &providerHallUsersStub{user: User{ID: 1, Status: StatusActive}}
	groups := &providerHallGroupsStub{groups: []Group{
		{ID: 1, Platform: PlatformOpenAI, Status: StatusActive},
		{ID: 2, Platform: PlatformOpenAI, Status: StatusActive, IsExclusive: true},
		{ID: 3, Platform: PlatformComposite, Status: StatusActive, SubscriptionType: SubscriptionTypeSubscription},
	}}
	subs := &providerHallSubscriptionsStub{}
	keys := &APIKeyService{userRepo: users, groupRepo: groups, userSubRepo: subs}
	repo := &providerHallRepoStub{config: ProviderHallConfig{DisplayEnabled: true}, listing: ProviderHallGroup{Listed: true}}
	svc := NewProviderHallService(repo, groups, users, keys)
	_, err := svc.AuthorizeUserGroup(ctx, 1, 1)
	require.NoError(t, err)
	_, err = svc.AuthorizeUserGroup(ctx, 1, 2)
	require.ErrorIs(t, err, ErrProviderHallNotFound)
	_, err = svc.AuthorizeUserGroup(ctx, 1, 3)
	require.ErrorIs(t, err, ErrProviderHallNotFound)
	users.user.AllowedGroups = []int64{2}
	_, err = svc.AuthorizeUserGroup(ctx, 1, 2)
	require.NoError(t, err)
	users.user.AllowedGroups = nil
	_, err = svc.AuthorizeUserGroup(ctx, 1, 2)
	require.ErrorIs(t, err, ErrProviderHallNotFound)
	subs.subscriptions = []UserSubscription{{GroupID: 3}}
	_, err = svc.AuthorizeUserGroup(ctx, 1, 3)
	require.NoError(t, err)
	subs.subscriptions = nil
	_, err = svc.AuthorizeUserGroup(ctx, 1, 3)
	require.ErrorIs(t, err, ErrProviderHallNotFound)
	repo.listing.Listed = false
	_, err = svc.AuthorizeUserGroup(ctx, 1, 1)
	require.ErrorIs(t, err, ErrProviderHallNotFound)
	repo.listing.Listed = true
	repo.config.DisplayEnabled = false
	_, err = svc.AuthorizeUserGroup(ctx, 1, 1)
	require.ErrorIs(t, err, ErrProviderHallNotFound)
}

func TestProviderHallProfileReferenceValidation(t *testing.T) {
	p := ProviderHallProfile{Model: "gpt-test", Protocol: "responses", OutputLimit: 256}
	repo := &providerHallRepoStub{}
	svc := &ProviderHallService{repo: repo}
	_, err := svc.SaveProfile(context.Background(), p, 1)
	require.NoError(t, err)
	zero := "0"
	p.ReferenceInputPrice = &zero
	_, err = svc.SaveProfile(context.Background(), p, 1)
	require.Error(t, err)
	confirmed := time.Now().Add(-time.Hour)
	p.ReferenceCachePrice, p.ReferenceCacheRate, p.ReferenceConfirmedAt = &zero, &zero, &confirmed
	_, err = svc.SaveProfile(context.Background(), p, 1)
	require.NoError(t, err)
	tooHigh := "1.0000000001"
	p.ReferenceCacheRate = &tooHigh
	_, err = svc.SaveProfile(context.Background(), p, 1)
	require.Error(t, err)
	p.ReferenceCacheRate = &zero
	p.ModelAliases = []string{"gpt-*"}
	_, err = svc.SaveProfile(context.Background(), p, 1)
	require.Error(t, err)
}

func TestProviderHallAlgorithmFixedExamples(t *testing.T) {
	counts := map[int64]int64{100000: 1}
	for n := int64(1000); n <= 19000; n += 1000 {
		counts[n] = 1
	}
	latency := ProviderHallExactLatency(counts)
	require.Equal(t, int64(20), latency.Count)
	require.Equal(t, "10000", latency.Fast95MeanMS.String())
	require.Equal(t, int64(18000), *latency.P90MS)
	require.True(t, ProviderHallCacheRate(300, 800, 0).Equal(decimal.NewFromInt(800).Div(decimal.NewFromInt(1100))))
	require.True(t, ProviderHallSuccessRate(8, 2, 12).Equal(decimal.NewFromInt(2).Div(decimal.NewFromInt(3))))
	require.Equal(t, "1", ProviderHallInputPrice(decimal.RequireFromString("0.2"), 200000).String())
	require.Equal(t, "0.05", ProviderHallInputCost(decimal.RequireFromString("0.2"), decimal.NewFromInt(1), decimal.NewFromInt(4)).String())
	require.Equal(t, "0", ProviderHallInputCost(decimal.Zero, decimal.NewFromInt(1), decimal.NewFromInt(4)).String())
	require.Nil(t, ProviderHallInputPrice(decimal.Zero, 0))
	require.Nil(t, ProviderHallInputCost(decimal.Zero, decimal.Zero, decimal.Zero))
	require.Nil(t, ProviderHallCacheRate(0, 0, 0))
	require.True(t, ProviderHallCacheRate(10, 0, 0).IsZero())
	require.Equal(t, "0.5", ProviderHallCacheRate(10, 20, 10).String())
	require.Nil(t, ProviderHallSuccessRate(0, 0, 0))
	require.True(t, ProviderHallSuccessRate(0, 1, 0).IsZero())
	for _, n := range []int64{0, 1, 10, 20} {
		result := ProviderHallExactLatency(map[int64]int64{15: n, 0: 100, -1: 100, 99: -2})
		require.Equal(t, n, result.Count)
		if n == 0 {
			require.Nil(t, result.Fast95MeanMS)
			require.Nil(t, result.P90MS)
		} else {
			require.Equal(t, "15", result.Fast95MeanMS.String())
			require.Equal(t, int64(15), *result.P90MS)
		}
	}
	require.Equal(t, "2026-09-10", ProviderHallBudgetDay(time.Date(2026, 9, 10, 15, 59, 59, 0, time.UTC)))
	require.Equal(t, "2026-09-11", ProviderHallBudgetDay(time.Date(2026, 9, 10, 16, 0, 0, 0, time.UTC)))
}

func TestProviderHallReadinessGatesEachSwitch(t *testing.T) {
	operator := int64(9)
	valid := ProviderHallConfig{ProviderHallVersion: ProviderHallVersion{Version: 1}, DefaultProtocol: "responses", DefaultRange: "6h", DailyBudget: "1", GatewayOrigin: "https://panel.example", OperatorUserID: &operator}
	users := &providerHallUsersStub{user: User{ID: operator, Status: StatusActive}}
	for _, tc := range []struct {
		name  string
		ready ProviderHallReadiness
		edit  func(*ProviderHallConfig)
		ok    bool
		code  string
	}{
		{"collection_ready", ProviderHallReadiness{Collection: true}, func(c *ProviderHallConfig) { c.CollectionEnabled = true }, true, ""},
		{"collection_not_ready", ProviderHallReadiness{}, func(c *ProviderHallConfig) { c.CollectionEnabled = true }, false, "PROVIDER_HALL_NOT_READY"},
		{"tasks_not_ready", ProviderHallReadiness{Collection: true}, func(c *ProviderHallConfig) { c.CollectionEnabled, c.TasksEnabled = true, true }, false, "PROVIDER_HALL_NOT_READY"},
		{"tasks_requires_collection", ProviderHallReadiness{Collection: true, Tasks: true}, func(c *ProviderHallConfig) { c.TasksEnabled = true }, false, "PROVIDER_HALL_INVALID_CONFIG"},
		{"tasks_requires_origin", ProviderHallReadiness{Collection: true, Tasks: true}, func(c *ProviderHallConfig) { c.CollectionEnabled, c.TasksEnabled, c.GatewayOrigin = true, true, "" }, false, "PROVIDER_HALL_INVALID_CONFIG"},
		{"tasks_requires_operator", ProviderHallReadiness{Collection: true, Tasks: true}, func(c *ProviderHallConfig) { c.CollectionEnabled, c.TasksEnabled, c.OperatorUserID = true, true, nil }, false, "PROVIDER_HALL_INVALID_CONFIG"},
		{"tasks_ready", ProviderHallReadiness{Collection: true, Tasks: true}, func(c *ProviderHallConfig) { c.CollectionEnabled, c.TasksEnabled = true, true }, true, ""},
		{"display_not_ready", ProviderHallReadiness{Collection: true}, func(c *ProviderHallConfig) { c.CollectionEnabled, c.DisplayEnabled = true, true }, false, "PROVIDER_HALL_NOT_READY"},
		{"display_requires_collection", ProviderHallReadiness{Collection: true, Display: true}, func(c *ProviderHallConfig) { c.DisplayEnabled = true }, false, "PROVIDER_HALL_INVALID_CONFIG"},
		{"display_ready", ProviderHallReadiness{Collection: true, Display: true}, func(c *ProviderHallConfig) { c.CollectionEnabled, c.DisplayEnabled = true, true }, true, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &providerHallRepoStub{config: valid}
			svc := &ProviderHallService{repo: repo, users: users}
			svc.SetReadiness(tc.ready)
			require.Equal(t, tc.ready, svc.Readiness())
			cfg := valid
			tc.edit(&cfg)
			_, err := svc.UpdateConfig(context.Background(), cfg, 4)
			if tc.ok {
				require.NoError(t, err)
				require.Equal(t, 1, repo.writes)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.code)
			require.Zero(t, repo.writes)
		})
	}
}
