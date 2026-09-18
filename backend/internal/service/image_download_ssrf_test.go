package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageDownloadServiceDefaultClientBlocksPrivateDialAddress(t *testing.T) {
	svc := NewImageDownloadService()
	require.NotNil(t, svc)
	require.NotNil(t, svc.client)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://127.0.0.1/image.png", nil)
	require.NoError(t, err)
	resp, err := svc.client.Do(req)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}

	require.Error(t, err)
	require.Contains(t, err.Error(), "blocked by image download SSRF policy")
}

func TestImageDownloadServiceStopsAfterRedirectLimit(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: imageDownloadRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusFound,
			Header: http.Header{
				"Location": []string{"https://" + ImageDownloadAllowedHost + req.URL.Path + "/next"},
			},
			Body:    http.NoBody,
			Request: req,
		}, nil
	})}
	svc := newImageDownloadServiceWithClient(client)

	download, err := svc.Fetch(context.Background(), imageDownloadTestURL)
	require.Nil(t, download)
	require.ErrorIs(t, err, ErrImageDownloadUnavailable)
	require.Equal(t, imageDownloadMaxRedirects, calls)
}
