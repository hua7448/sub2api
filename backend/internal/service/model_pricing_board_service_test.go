//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type stubModelPricingBoardGroupProvider struct {
	groups []Group
	err    error
}

func (s *stubModelPricingBoardGroupProvider) GetAvailableGroups(context.Context, int64) ([]Group, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.groups, nil
}

type stubModelPricingBoardRateRepo struct {
	UserGroupRateRepository
	rates map[int64]float64
	err   error
}

func (s *stubModelPricingBoardRateRepo) GetByUserID(context.Context, int64) (map[int64]float64, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.rates == nil {
		return map[int64]float64{}, nil
	}
	return s.rates, nil
}

func newModelPricingBoardServiceForTest(
	channels []Channel,
	activeGroups []Group,
	userGroups []Group,
	userRates map[int64]float64,
	pricing *PricingService,
) *ModelPricingBoardService {
	channelRepo := &mockChannelRepository{
		listAllFn: func(context.Context) ([]Channel, error) {
			return channels, nil
		},
	}
	channelSvc := NewChannelService(channelRepo, &stubGroupRepoForAvailable{activeGroups: activeGroups}, nil, pricing)
	return NewModelPricingBoardService(
		channelSvc,
		&stubModelPricingBoardGroupProvider{groups: userGroups},
		&stubModelPricingBoardRateRepo{rates: userRates},
		pricing,
	)
}

func TestModelPricingBoard_NoAccessibleGroupsReturnsEmpty(t *testing.T) {
	svc := newModelPricingBoardServiceForTest(nil, nil, nil, nil, nil)

	items, err := svc.BuildUserBoard(context.Background(), 100)
	require.NoError(t, err)
	require.Empty(t, items)
}

func TestModelPricingBoard_SelectsLowestEffectiveSitePrice(t *testing.T) {
	inputA, outputA := 2e-6, 1e-5
	inputB, outputB := 3e-6, 1.2e-5
	channels := []Channel{
		{
			ID:       1,
			Name:     "expensive",
			Status:   StatusActive,
			GroupIDs: []int64{1},
			ModelPricing: []ChannelModelPricing{{
				Platform:    PlatformOpenAI,
				Models:      []string{"gpt-5"},
				BillingMode: BillingModeToken,
				InputPrice:  &inputA,
				OutputPrice: &outputA,
			}},
		},
		{
			ID:       2,
			Name:     "cheaper",
			Status:   StatusActive,
			GroupIDs: []int64{2},
			ModelPricing: []ChannelModelPricing{{
				Platform:    PlatformOpenAI,
				Models:      []string{"gpt-5"},
				BillingMode: BillingModeToken,
				InputPrice:  &inputB,
				OutputPrice: &outputB,
			}},
		},
	}
	activeGroups := []Group{
		{ID: 1, Name: "standard", Platform: PlatformOpenAI, Status: StatusActive, RateMultiplier: 1.0},
		{ID: 2, Name: "vip", Platform: PlatformOpenAI, Status: StatusActive, RateMultiplier: 0.5},
	}
	svc := newModelPricingBoardServiceForTest(channels, activeGroups, activeGroups, nil, nil)

	items, err := svc.BuildUserBoard(context.Background(), 100)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, "gpt-5", items[0].ModelID)
	require.Equal(t, "cheaper", items[0].ChannelName)
	require.Equal(t, "vip", items[0].GroupName)
	require.NotNil(t, items[0].SiteInputPrice)
	require.InDelta(t, 1.5e-6, *items[0].SiteInputPrice, 1e-12)
	require.NotNil(t, items[0].SiteOutputPrice)
	require.InDelta(t, 6e-6, *items[0].SiteOutputPrice, 1e-12)
}

func TestModelPricingBoard_UserRateOverridesGroupDefault(t *testing.T) {
	input, output := 2e-6, 1e-5
	channels := []Channel{{
		ID:       1,
		Name:     "main",
		Status:   StatusActive,
		GroupIDs: []int64{10},
		ModelPricing: []ChannelModelPricing{{
			Platform:    PlatformAnthropic,
			Models:      []string{"claude-sonnet-4"},
			BillingMode: BillingModeToken,
			InputPrice:  &input,
			OutputPrice: &output,
		}},
	}}
	groups := []Group{{ID: 10, Name: "pro", Platform: PlatformAnthropic, Status: StatusActive, RateMultiplier: 1.2}}
	svc := newModelPricingBoardServiceForTest(channels, groups, groups, map[int64]float64{10: 0.7}, nil)

	items, err := svc.BuildUserBoard(context.Background(), 100)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.True(t, items[0].UserRateOverridden)
	require.Equal(t, 0.7, items[0].RateMultiplier)
	require.NotNil(t, items[0].SiteInputPrice)
	require.InDelta(t, 1.4e-6, *items[0].SiteInputPrice, 1e-12)
}

func TestModelPricingBoard_OfficialMissingKeepsSitePriceAndNullSavings(t *testing.T) {
	input, output := 2e-6, 1e-5
	channels := []Channel{{
		ID:       1,
		Name:     "main",
		Status:   StatusActive,
		GroupIDs: []int64{10},
		ModelPricing: []ChannelModelPricing{{
			Platform:    PlatformOpenAI,
			Models:      []string{"unknown-model"},
			BillingMode: BillingModeToken,
			InputPrice:  &input,
			OutputPrice: &output,
		}},
	}}
	groups := []Group{{ID: 10, Name: "pro", Platform: PlatformOpenAI, Status: StatusActive, RateMultiplier: 1.0}}
	svc := newModelPricingBoardServiceForTest(channels, groups, groups, nil, nil)

	items, err := svc.BuildUserBoard(context.Background(), 100)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.NotNil(t, items[0].SiteInputPrice)
	require.Nil(t, items[0].OfficialInputPrice)
	require.Nil(t, items[0].InputSavings)
}
