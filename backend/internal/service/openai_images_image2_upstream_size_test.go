package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestRewriteOpenAIImagesRequest_GPTImage2PreservesExactSizeAndImage2Fields(t *testing.T) {
	body := []byte(`{"model":"gpt-image-2","prompt":"draw","size":"3840x2160","quality":"hd","type":"url","reference_images":"https://example.com/reference.jpg"}`)

	rewritten, contentType, err := rewriteOpenAIImagesRequest(body, "application/json", "gpt-image-2", "3840x2160")
	require.NoError(t, err)
	require.Equal(t, "application/json", contentType)
	require.Equal(t, "3840x2160", gjson.GetBytes(rewritten, "size").String())
	require.Equal(t, "hd", gjson.GetBytes(rewritten, "quality").String())
	require.Equal(t, "url", gjson.GetBytes(rewritten, "type").String())
	require.Equal(t, "https://example.com/reference.jpg", gjson.GetBytes(rewritten, "reference_images").String())
}

func TestRewriteOpenAIImagesRequest_LegacyModelStillNormalizesCustomSize(t *testing.T) {
	body := []byte(`{"model":"gpt-image-1","prompt":"draw","size":"3840x2160"}`)

	rewritten, _, err := rewriteOpenAIImagesRequest(body, "application/json", "gpt-image-1", "3840x2160")
	require.NoError(t, err)
	require.Equal(t, "1536x1024", gjson.GetBytes(rewritten, "size").String())
}
