package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestCORSPreflightAllowsImageDownloadAuthorization(t *testing.T) {
	middleware := CORS(config.CORSConfig{
		AllowedOrigins:   []string{"https://panel.sharesai.xyz"},
		AllowCredentials: true,
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodOptions, "/v1/images/download", nil)
	c.Request.Header.Set("Origin", "https://panel.sharesai.xyz")
	c.Request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	c.Request.Header.Set("Access-Control-Request-Headers", "authorization,content-type")

	middleware(c)

	require.Equal(t, http.StatusNoContent, w.Code)
	require.Equal(t, "https://panel.sharesai.xyz", w.Header().Get("Access-Control-Allow-Origin"))
	require.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), http.MethodPost)
	require.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	require.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}
