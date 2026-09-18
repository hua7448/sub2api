package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestGatewayRoutesImageDownloadPathIsRegisteredWithoutRootAlias(t *testing.T) {
	router := newGatewayRoutesTestRouter()
	registered := make(map[string]bool)
	for _, route := range router.Routes() {
		registered[route.Method+" "+route.Path] = true
	}

	require.True(t, registered["POST /v1/images/download"])
	require.False(t, registered["POST /images/download"])
}

func TestGatewayRoutesImageDownloadRequiresAPIKeyAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	authCalls := 0
	RegisterGatewayRoutes(
		router,
		&handler.Handlers{
			Gateway:       &handler.GatewayHandler{},
			OpenAIGateway: &handler.OpenAIGatewayHandler{},
			AsyncImage:    handler.NewAsyncImageHandler(nil, nil),
		},
		servermiddleware.APIKeyAuthMiddleware(func(c *gin.Context) {
			authCalls++
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
		}),
		nil,
		nil,
		nil,
		nil,
		nil,
		&config.Config{Gateway: config.GatewayConfig{MaxBodySize: 1024, TextMaxBodySize: 1024}},
	)

	req := httptest.NewRequest(http.MethodPost, "/v1/images/download", strings.NewReader(`{"url":"https://pre-signed-firefly-prod.s3-accelerate.amazonaws.com/image.png"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)
	require.Equal(t, 1, authCalls)
}
