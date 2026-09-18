package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const imageDownloadTestURL = "https://" + ImageDownloadAllowedHost + "/images/generated.png?signature=test"

type imageDownloadRoundTripFunc func(*http.Request) (*http.Response, error)

func (f imageDownloadRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type countingImageBody struct {
	reader *bytes.Reader
	read   int
	closed bool
}

func newCountingImageBody(data []byte) *countingImageBody {
	return &countingImageBody{reader: bytes.NewReader(data)}
}

func (b *countingImageBody) Read(p []byte) (int, error) {
	n, err := b.reader.Read(p)
	b.read += n
	return n, err
}

func (b *countingImageBody) Close() error {
	b.closed = true
	return nil
}

func testPNG(size int) []byte {
	if size < 8 {
		size = 8
	}
	data := make([]byte, size)
	copy(data, []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a})
	return data
}

func TestImageDownloadServiceStreamsAllowedImage(t *testing.T) {
	payload := testPNG(4096)
	upstreamBody := newCountingImageBody(payload)
	client := &http.Client{Transport: imageDownloadRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, ImageDownloadAllowedHost, req.URL.Host)
		require.Equal(t, "identity", req.Header.Get("Accept-Encoding"))
		require.Equal(t, "Sub2API-Image-Download/1.0", req.Header.Get("User-Agent"))
		return &http.Response{
			StatusCode:    http.StatusOK,
			Header:        http.Header{"Content-Type": []string{"image/png"}},
			Body:          upstreamBody,
			ContentLength: int64(len(payload)),
			Request:       req,
		}, nil
	})}
	svc := newImageDownloadServiceWithClient(client)

	download, err := svc.Fetch(context.Background(), imageDownloadTestURL)
	require.NoError(t, err)
	require.Equal(t, int64(len(payload)), download.ContentLength)
	require.Equal(t, "image/png", download.ContentType)
	require.Equal(t, imageDownloadSniffBytes, upstreamBody.read, "Fetch must not buffer the complete image")

	got, err := io.ReadAll(download.Body)
	require.NoError(t, err)
	require.Equal(t, payload, got)
	require.NoError(t, download.Close())
	require.True(t, upstreamBody.closed)
}

func TestImageDownloadServiceRejectsUntrustedURLsBeforeRequest(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: imageDownloadRoundTripFunc(func(*http.Request) (*http.Response, error) {
		calls++
		return nil, errors.New("must not be called")
	})}
	svc := newImageDownloadServiceWithClient(client)

	for _, rawURL := range []string{
		"http://" + ImageDownloadAllowedHost + "/image.png",
		"https://example.com/image.png",
		"https://" + ImageDownloadAllowedHost + ".evil.example/image.png",
		"https://127.0.0.1/image.png",
		"https://user@" + ImageDownloadAllowedHost + "/image.png",
		"https://" + ImageDownloadAllowedHost + ":8443/image.png",
		"https://" + ImageDownloadAllowedHost + "/image.png#fragment",
	} {
		t.Run(rawURL, func(t *testing.T) {
			download, err := svc.Fetch(context.Background(), rawURL)
			require.Nil(t, download)
			require.ErrorIs(t, err, ErrImageDownloadInvalidURL)
		})
	}
	require.Zero(t, calls)
}

func TestImageDownloadServiceRejectsRedirectToUntrustedHost(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: imageDownloadRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusFound,
			Header:     http.Header{"Location": []string{"https://127.0.0.1/private.png"}},
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	})}
	svc := newImageDownloadServiceWithClient(client)

	download, err := svc.Fetch(context.Background(), imageDownloadTestURL)
	require.Nil(t, download)
	require.ErrorIs(t, err, ErrImageDownloadUnavailable)
	require.Equal(t, 1, calls, "blocked redirect must not be dialed")
}

func TestImageDownloadServiceAllowsRedirectOnTrustedHost(t *testing.T) {
	payload := testPNG(1024)
	calls := 0
	client := &http.Client{Transport: imageDownloadRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		if req.URL.Path == "/images/generated.png" {
			return &http.Response{
				StatusCode: http.StatusFound,
				Header:     http.Header{"Location": []string{"https://" + ImageDownloadAllowedHost + "/images/final.png"}},
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			}, nil
		}
		return &http.Response{
			StatusCode:    http.StatusOK,
			Header:        http.Header{"Content-Type": []string{"application/octet-stream"}},
			Body:          io.NopCloser(bytes.NewReader(payload)),
			ContentLength: int64(len(payload)),
			Request:       req,
		}, nil
	})}
	svc := newImageDownloadServiceWithClient(client)

	download, err := svc.Fetch(context.Background(), imageDownloadTestURL)
	require.NoError(t, err)
	defer func() { _ = download.Close() }()
	require.Equal(t, 2, calls)
	require.Equal(t, "image/png", download.ContentType)
}

func TestImageDownloadServiceRejectsMissingAndOversizedContentLength(t *testing.T) {
	for _, tc := range []struct {
		name    string
		length  int64
		wantErr error
	}{
		{name: "missing", length: -1, wantErr: ErrImageDownloadInvalidLength},
		{name: "zero", length: 0, wantErr: ErrImageDownloadInvalidLength},
		{name: "oversized", length: ImageDownloadMaxBytes + 1, wantErr: ErrImageDownloadTooLarge},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := newCountingImageBody(testPNG(512))
			client := &http.Client{Transport: imageDownloadRoundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        http.Header{"Content-Type": []string{"image/png"}},
					Body:          body,
					ContentLength: tc.length,
					Request:       req,
				}, nil
			})}
			download, err := newImageDownloadServiceWithClient(client).Fetch(context.Background(), imageDownloadTestURL)
			require.Nil(t, download)
			require.ErrorIs(t, err, tc.wantErr)
			require.True(t, body.closed)
			require.Zero(t, body.read)
		})
	}
}

func TestImageDownloadServiceReturnsUnavailableForExpiredAndUpstreamFailures(t *testing.T) {
	t.Run("expired", func(t *testing.T) {
		client := &http.Client{Transport: imageDownloadRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode:    http.StatusForbidden,
				Header:        http.Header{"Content-Type": []string{"application/xml"}},
				Body:          io.NopCloser(strings.NewReader("<Error>ExpiredToken</Error>")),
				ContentLength: 27,
				Request:       req,
			}, nil
		})}
		download, err := newImageDownloadServiceWithClient(client).Fetch(context.Background(), imageDownloadTestURL)
		require.Nil(t, download)
		require.ErrorIs(t, err, ErrImageDownloadUnavailable)
	})

	t.Run("network error", func(t *testing.T) {
		client := &http.Client{Transport: imageDownloadRoundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("upstream reset")
		})}
		download, err := newImageDownloadServiceWithClient(client).Fetch(context.Background(), imageDownloadTestURL)
		require.Nil(t, download)
		require.ErrorIs(t, err, ErrImageDownloadUnavailable)
	})
}

func TestImageDownloadServiceRejectsInvalidOrMismatchedContentType(t *testing.T) {
	for _, tc := range []struct {
		name        string
		contentType string
		payload     []byte
	}{
		{name: "html", contentType: "text/html", payload: []byte("<html>not an image</html>")},
		{name: "mismatched", contentType: "image/jpeg", payload: testPNG(512)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &http.Client{Transport: imageDownloadRoundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode:    http.StatusOK,
					Header:        http.Header{"Content-Type": []string{tc.contentType}},
					Body:          io.NopCloser(bytes.NewReader(tc.payload)),
					ContentLength: int64(len(tc.payload)),
					Request:       req,
				}, nil
			})}
			download, err := newImageDownloadServiceWithClient(client).Fetch(context.Background(), imageDownloadTestURL)
			require.Nil(t, download)
			require.ErrorIs(t, err, ErrImageDownloadInvalidType)
		})
	}
}
