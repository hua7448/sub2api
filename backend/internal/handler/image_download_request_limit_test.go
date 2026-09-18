package handler

import (
	"context"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestImageDownloadHandlerRejectsOversizedRequestBeforeFetch(t *testing.T) {
	h := &ImageDownloadHandler{
		download: imageDownloadFetcherFunc(func(context.Context, string) (*service.ImageDownload, error) {
			t.Fatal("fetch must not be called for an oversized request")
			return nil, nil
		}),
		limiter: newImageDownloadConcurrencyGate(),
	}

	recorder := imageDownloadHandlerRequest(t, h, map[string]string{
		"url": "https://" + service.ImageDownloadAllowedHost + "/" + strings.Repeat("a", imageDownloadRequestMaxBytes),
	})

	require.Equal(t, 400, recorder.Code)
	require.Contains(t, recorder.Body.String(), "invalid_request_error")
}
