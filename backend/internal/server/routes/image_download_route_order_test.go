package routes

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGatewayImageDownloadRouteRunsBeforeCompositeBodyParsing(t *testing.T) {
	source, err := os.ReadFile("gateway.go")
	require.NoError(t, err)
	text := string(source)

	authIndex := strings.Index(text, "gateway.Use(gin.HandlerFunc(apiKeyAuth))")
	downloadIndex := strings.Index(text, `gateway.POST("/images/download", imageDownloadHandler.Download)`)
	compositeIndex := strings.Index(text, "gateway.Use(compositeTarget)")

	require.NotEqual(t, -1, authIndex)
	require.NotEqual(t, -1, downloadIndex)
	require.NotEqual(t, -1, compositeIndex)
	require.Less(t, authIndex, downloadIndex, "download must remain API-key authenticated")
	require.Less(t, downloadIndex, compositeIndex, "8 KiB handler limit must run before composite body parsing")
}
