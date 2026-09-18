package handler

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// SetProviderHallCollector attaches the provider hall request collector. The
// handler works unchanged when it is nil.
func (h *OpenAIGatewayHandler) SetProviderHallCollector(c *service.ProviderHallCollector) {
	if h != nil {
		h.providerHall = c
	}
}

// providerHallBegin registers the request with the collector once auth, body
// validation, admission and pricing context are done. The task header is
// removed from the inbound request unconditionally so no upstream builder can
// ever forward it, even when collection is disabled.
func (h *OpenAIGatewayHandler) providerHallBegin(c *gin.Context, apiKey *service.APIKey, protocol, reqModel string, stream bool, startedAt time.Time) *service.ProviderHallRequestTracker {
	if c == nil || c.Request == nil {
		return nil
	}
	taskHeader := c.GetHeader(service.ProviderHallTaskHeader)
	c.Request.Header.Del(service.ProviderHallTaskHeader)
	if h == nil || h.providerHall == nil || apiKey == nil || apiKey.GroupID == nil {
		return nil
	}
	tracker, ctx := h.providerHall.Begin(c.Request.Context(), service.ProviderHallBeginInput{
		APIKeyID:       apiKey.ID,
		GroupID:        *apiKey.GroupID,
		Protocol:       protocol,
		RequestedModel: reqModel,
		Stream:         stream,
		StartedAt:      startedAt,
		TaskHeader:     taskHeader,
	})
	if tracker == nil {
		return nil
	}
	c.Request = c.Request.WithContext(ctx)
	return tracker
}
