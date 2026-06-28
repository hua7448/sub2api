package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

const (
	ModelPricingBoardProviderAnthropic = PlatformAnthropic
	ModelPricingBoardProviderOpenAI    = PlatformOpenAI
)

var modelPricingBoardBasePlatforms = map[string]struct{}{
	ModelPricingBoardProviderAnthropic: {},
	ModelPricingBoardProviderOpenAI:    {},
}

// ModelPricingBoardItem is one row in the user-facing model pricing board.
// Prices are stored per token in USD; the frontend formats them as USD / 1M tokens.
type ModelPricingBoardItem struct {
	ModelID                string
	Platform               string
	GroupID                int64
	GroupName              string
	ChannelID              int64
	ChannelName            string
	RateMultiplier         float64
	UserRateOverridden     bool
	SiteInputPrice         *float64
	SiteOutputPrice        *float64
	SiteCacheReadPrice     *float64
	OfficialInputPrice     *float64
	OfficialOutputPrice    *float64
	OfficialCacheReadPrice *float64
	InputSavings           *float64
	OutputSavings          *float64
	CacheReadSavings       *float64
}

type modelPricingBoardCandidate struct {
	modelID            string
	platform           string
	group              AvailableGroupRef
	channelID          int64
	channelName        string
	rateMultiplier     float64
	userRateOverridden bool
	pricing            *ChannelModelPricing
	official           *LiteLLMModelPricing
}

type modelPricingBoardGroupProvider interface {
	GetAvailableGroups(ctx context.Context, userID int64) ([]Group, error)
}

// ModelPricingBoardService builds the current user's lowest effective model prices
// from the same channel/group data used by the user-facing available-channels page.
type ModelPricingBoardService struct {
	channelService    *ChannelService
	groupProvider     modelPricingBoardGroupProvider
	userGroupRateRepo UserGroupRateRepository
	pricingService    *PricingService
	billingService    *BillingService
}

func NewModelPricingBoardService(
	channelService *ChannelService,
	groupProvider modelPricingBoardGroupProvider,
	userGroupRateRepo UserGroupRateRepository,
	pricingService *PricingService,
	billingService *BillingService,
) *ModelPricingBoardService {
	return &ModelPricingBoardService{
		channelService:    channelService,
		groupProvider:     groupProvider,
		userGroupRateRepo: userGroupRateRepo,
		pricingService:    pricingService,
		billingService:    billingService,
	}
}

// BuildUserBoard returns the lowest effective token-model price per platform+model
// for the given user. Only active channels and user-accessible active groups are considered.
func (s *ModelPricingBoardService) BuildUserBoard(ctx context.Context, userID int64) ([]ModelPricingBoardItem, error) {
	if s == nil || s.channelService == nil || s.groupProvider == nil {
		return []ModelPricingBoardItem{}, nil
	}

	userGroups, err := s.groupProvider.GetAvailableGroups(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get available groups: %w", err)
	}
	allowedGroups := make(map[int64]struct{}, len(userGroups))
	for i := range userGroups {
		g := userGroups[i]
		if isModelPricingBoardBasePlatform(g.Platform) {
			allowedGroups[g.ID] = struct{}{}
		}
	}
	if len(allowedGroups) == 0 {
		return []ModelPricingBoardItem{}, nil
	}

	userRates := map[int64]float64{}
	if s.userGroupRateRepo != nil {
		rates, err := s.userGroupRateRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("get user group rates: %w", err)
		}
		userRates = rates
	}

	channels, err := s.channelService.ListAvailable(ctx)
	if err != nil {
		return nil, fmt.Errorf("list available channels: %w", err)
	}

	bestByModel := make(map[string]ModelPricingBoardItem)
	for _, ch := range channels {
		if ch.Status != StatusActive {
			continue
		}
		visibleGroupsByPlatform := visiblePricingBoardGroups(ch.Groups, allowedGroups)
		if len(visibleGroupsByPlatform) == 0 {
			continue
		}
		for _, model := range ch.SupportedModels {
			if !isModelPricingBoardPlatform(model.Platform, model.Name) {
				continue
			}
			if model.Pricing == nil {
				continue
			}
			mode := model.Pricing.BillingMode
			if mode == "" {
				mode = BillingModeToken
			}
			if mode != BillingModeToken {
				continue
			}
			groups := groupsForPricingBoardModel(visibleGroupsByPlatform, model.Platform, model.Name)
			if len(groups) == 0 {
				continue
			}
			official := s.lookupOfficialPricing(model.Name)
			for _, group := range groups {
				rate, overridden := effectiveGroupRate(group, userRates)
				candidate := modelPricingBoardCandidate{
					modelID:            model.Name,
					platform:           model.Platform,
					group:              group,
					channelID:          ch.ID,
					channelName:        ch.Name,
					rateMultiplier:     rate,
					userRateOverridden: overridden,
					pricing:            model.Pricing,
					official:           official,
				}
				item := candidate.toItem()
				key := boardModelKey(item.Platform, item.ModelID)
				if existing, ok := bestByModel[key]; !ok || boardItemLess(item, existing) {
					bestByModel[key] = item
				}
			}
		}
	}

	out := make([]ModelPricingBoardItem, 0, len(bestByModel))
	for _, item := range bestByModel {
		out = append(out, item)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Platform != out[j].Platform {
			return out[i].Platform < out[j].Platform
		}
		return strings.ToLower(out[i].ModelID) < strings.ToLower(out[j].ModelID)
	})
	return out, nil
}

func visiblePricingBoardGroups(groups []AvailableGroupRef, allowed map[int64]struct{}) map[string][]AvailableGroupRef {
	out := make(map[string][]AvailableGroupRef, 2)
	for _, group := range groups {
		if _, ok := allowed[group.ID]; !ok {
			continue
		}
		if !isModelPricingBoardBasePlatform(group.Platform) {
			continue
		}
		out[group.Platform] = append(out[group.Platform], group)
	}
	return out
}

func groupsForPricingBoardModel(groupsByPlatform map[string][]AvailableGroupRef, platform, modelID string) []AvailableGroupRef {
	if groups := groupsByPlatform[platform]; len(groups) > 0 {
		return groups
	}
	if !isDomesticModelID(modelID) {
		return nil
	}
	var out []AvailableGroupRef
	seen := make(map[int64]struct{})
	for _, groups := range groupsByPlatform {
		for _, group := range groups {
			if _, ok := seen[group.ID]; ok {
				continue
			}
			seen[group.ID] = struct{}{}
			out = append(out, group)
		}
	}
	return out
}

func isModelPricingBoardBasePlatform(platform string) bool {
	_, ok := modelPricingBoardBasePlatforms[platform]
	return ok
}

func isModelPricingBoardPlatform(platform, modelID string) bool {
	if isModelPricingBoardBasePlatform(platform) {
		return true
	}
	return isDomesticModelID(modelID)
}

func isDomesticModelID(modelID string) bool {
	model := strings.ToLower(strings.TrimSpace(modelID))
	return strings.Contains(model, "deepseek") ||
		strings.Contains(model, "glm") ||
		strings.Contains(model, "chatglm") ||
		strings.Contains(model, "kimi") ||
		strings.Contains(model, "moonshot") ||
		strings.Contains(model, "minimax") ||
		strings.Contains(model, "doubao")
}

func effectiveGroupRate(group AvailableGroupRef, userRates map[int64]float64) (float64, bool) {
	if rate, ok := userRates[group.ID]; ok {
		return rate, true
	}
	return group.RateMultiplier, false
}

func (s *ModelPricingBoardService) lookupOfficialPricing(modelID string) *LiteLLMModelPricing {
	if s == nil || s.pricingService == nil {
		return s.lookupFallbackOfficialPricing(modelID)
	}
	if pricing := s.pricingService.GetModelPricing(modelID); pricing != nil {
		return pricing
	}
	return s.lookupFallbackOfficialPricing(modelID)
}

func (s *ModelPricingBoardService) lookupFallbackOfficialPricing(modelID string) *LiteLLMModelPricing {
	if s == nil || s.billingService == nil {
		return nil
	}
	pricing, err := s.billingService.GetModelPricing(modelID)
	if err != nil || pricing == nil {
		return nil
	}
	return &LiteLLMModelPricing{
		InputCostPerToken:               pricing.InputPricePerToken,
		OutputCostPerToken:              pricing.OutputPricePerToken,
		CacheReadInputTokenCost:         pricing.CacheReadPricePerToken,
		SupportsPromptCaching:           pricing.CacheReadPricePerToken > 0,
		LongContextInputTokenThreshold:  pricing.LongContextInputThreshold,
		LongContextInputCostMultiplier:  pricing.LongContextInputMultiplier,
		LongContextOutputCostMultiplier: pricing.LongContextOutputMultiplier,
		Mode:                            "chat",
	}
}

func (c modelPricingBoardCandidate) toItem() ModelPricingBoardItem {
	item := ModelPricingBoardItem{
		ModelID:            c.modelID,
		Platform:           c.platform,
		GroupID:            c.group.ID,
		GroupName:          c.group.Name,
		ChannelID:          c.channelID,
		ChannelName:        c.channelName,
		RateMultiplier:     c.rateMultiplier,
		UserRateOverridden: c.userRateOverridden,
	}
	if c.pricing != nil {
		item.SiteInputPrice = scaledPricePtr(c.pricing.InputPrice, c.rateMultiplier)
		item.SiteOutputPrice = scaledPricePtr(c.pricing.OutputPrice, c.rateMultiplier)
		item.SiteCacheReadPrice = scaledPricePtr(c.pricing.CacheReadPrice, c.rateMultiplier)
	}
	if c.official != nil {
		item.OfficialInputPrice = nonZeroPtr(c.official.InputCostPerToken)
		item.OfficialOutputPrice = nonZeroPtr(c.official.OutputCostPerToken)
		item.OfficialCacheReadPrice = nonZeroPtr(c.official.CacheReadInputTokenCost)
	}
	item.InputSavings = savingsRatio(item.OfficialInputPrice, item.SiteInputPrice)
	item.OutputSavings = savingsRatio(item.OfficialOutputPrice, item.SiteOutputPrice)
	item.CacheReadSavings = savingsRatio(item.OfficialCacheReadPrice, item.SiteCacheReadPrice)
	return item
}

func scaledPricePtr(price *float64, multiplier float64) *float64 {
	if price == nil {
		return nil
	}
	v := *price * multiplier
	return &v
}

func savingsRatio(official, site *float64) *float64 {
	if official == nil || site == nil || *official <= 0 {
		return nil
	}
	v := (*official - *site) / *official
	return &v
}

func boardModelKey(platform, modelID string) string {
	return platform + "\x00" + strings.ToLower(modelID)
}

func boardItemLess(a, b ModelPricingBoardItem) bool {
	scoreA, okA := comparableBoardPrice(a)
	scoreB, okB := comparableBoardPrice(b)
	if okA && okB && scoreA != scoreB {
		return scoreA < scoreB
	}
	if okA != okB {
		return okA
	}
	if a.RateMultiplier != b.RateMultiplier {
		return a.RateMultiplier < b.RateMultiplier
	}
	if a.ChannelName != b.ChannelName {
		return strings.ToLower(a.ChannelName) < strings.ToLower(b.ChannelName)
	}
	return strings.ToLower(a.GroupName) < strings.ToLower(b.GroupName)
}

func comparableBoardPrice(item ModelPricingBoardItem) (float64, bool) {
	if item.SiteInputPrice != nil || item.SiteOutputPrice != nil {
		var total float64
		if item.SiteInputPrice != nil {
			total += *item.SiteInputPrice
		}
		if item.SiteOutputPrice != nil {
			total += *item.SiteOutputPrice
		}
		return total, true
	}
	if item.SiteCacheReadPrice != nil {
		return *item.SiteCacheReadPrice, true
	}
	return 0, false
}
