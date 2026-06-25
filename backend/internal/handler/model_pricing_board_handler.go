package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type ModelPricingBoardHandler struct {
	boardService   *service.ModelPricingBoardService
	settingService *service.SettingService
}

func NewModelPricingBoardHandler(
	boardService *service.ModelPricingBoardService,
	settingService *service.SettingService,
) *ModelPricingBoardHandler {
	return &ModelPricingBoardHandler{
		boardService:   boardService,
		settingService: settingService,
	}
}

type modelPricingBoardItemResponse struct {
	ModelID                string   `json:"model_id"`
	Platform               string   `json:"platform"`
	GroupID                int64    `json:"group_id"`
	GroupName              string   `json:"group_name"`
	ChannelID              int64    `json:"channel_id"`
	ChannelName            string   `json:"channel_name"`
	RateMultiplier         float64  `json:"rate_multiplier"`
	UserRateOverridden     bool     `json:"user_rate_overridden"`
	SiteInputPrice         *float64 `json:"site_input_price"`
	SiteOutputPrice        *float64 `json:"site_output_price"`
	SiteCacheReadPrice     *float64 `json:"site_cache_read_price"`
	OfficialInputPrice     *float64 `json:"official_input_price"`
	OfficialOutputPrice    *float64 `json:"official_output_price"`
	OfficialCacheReadPrice *float64 `json:"official_cache_read_price"`
	InputSavings           *float64 `json:"input_savings"`
	OutputSavings          *float64 `json:"output_savings"`
	CacheReadSavings       *float64 `json:"cache_read_savings"`
}

type modelPricingBoardResponse struct {
	Items []modelPricingBoardItemResponse `json:"items"`
}

func (h *ModelPricingBoardHandler) featureEnabled(c *gin.Context) bool {
	if h.settingService == nil {
		return false
	}
	return h.settingService.GetModelPricingBoardRuntime(c.Request.Context()).Enabled
}

func toModelPricingBoardResponse(items []service.ModelPricingBoardItem) modelPricingBoardResponse {
	out := make([]modelPricingBoardItemResponse, 0, len(items))
	for _, item := range items {
		out = append(out, modelPricingBoardItemResponse{
			ModelID:                item.ModelID,
			Platform:               item.Platform,
			GroupID:                item.GroupID,
			GroupName:              item.GroupName,
			ChannelID:              item.ChannelID,
			ChannelName:            item.ChannelName,
			RateMultiplier:         item.RateMultiplier,
			UserRateOverridden:     item.UserRateOverridden,
			SiteInputPrice:         item.SiteInputPrice,
			SiteOutputPrice:        item.SiteOutputPrice,
			SiteCacheReadPrice:     item.SiteCacheReadPrice,
			OfficialInputPrice:     item.OfficialInputPrice,
			OfficialOutputPrice:    item.OfficialOutputPrice,
			OfficialCacheReadPrice: item.OfficialCacheReadPrice,
			InputSavings:           item.InputSavings,
			OutputSavings:          item.OutputSavings,
			CacheReadSavings:       item.CacheReadSavings,
		})
	}
	return modelPricingBoardResponse{Items: out}
}

// Board GET /api/v1/model-pricing/board
func (h *ModelPricingBoardHandler) Board(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}
	if !h.featureEnabled(c) {
		response.Success(c, modelPricingBoardResponse{Items: []modelPricingBoardItemResponse{}})
		return
	}
	items, err := h.boardService.BuildUserBoard(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, toModelPricingBoardResponse(items))
}
