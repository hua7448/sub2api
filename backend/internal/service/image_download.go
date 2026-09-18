package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ImageDownloadAllowedHost = "pre-signed-firefly-prod.s3-accelerate.amazonaws.com"
	ImageDownloadMaxBytes    = 64 << 20

	imageDownloadSniffBytes   = 512
	imageDownloadMaxRedirects = 5
)

var (
	ErrImageDownloadInvalidURL    = errors.New("invalid image download URL")
	ErrImageDownloadUnavailable   = errors.New("image URL has expired or is unavailable")
	ErrImageDownloadInvalidLength = errors.New("upstream image has an invalid content length")
	ErrImageDownloadTooLarge      = errors.New("upstream image exceeds the 64 MiB download limit")
	ErrImageDownloadInvalidType   = errors.New("upstream response is not a supported image")
)

type ImageDownload struct {
	Body          io.ReadCloser
	ContentType   string
	ContentLength int64
}

type ImageDownloadService struct {
	client *http.Client
}

func NewImageDownloadService() *ImageDownloadService {
	return newImageDownloadServiceWithClient(nil)
}

func newImageDownloadServiceWithClient(client *http.Client) *ImageDownloadService {
	if client == nil {
		client = newImageDownloadHTTPClient(5 * time.Minute)
	} else {
		clone := *client
		client = &clone
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= imageDownloadMaxRedirects {
			return errors.New("too many image download redirects")
		}
		return validateImageDownloadURL(req.URL)
	}
	return &ImageDownloadService{client: client}
}

func (s *ImageDownloadService) Fetch(ctx context.Context, rawURL string) (*ImageDownload, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || validateImageDownloadURL(parsed) != nil {
		return nil, ErrImageDownloadInvalidURL
	}
	if s == nil || s.client == nil {
		return nil, ErrImageDownloadUnavailable
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, ErrImageDownloadInvalidURL
	}
	req.Header.Set("Accept", "image/webp,image/png,image/jpeg,image/gif,image/*;q=0.8")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", "Sub2API-Image-Download/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return nil, ErrImageDownloadUnavailable
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		_ = resp.Body.Close()
		return nil, ErrImageDownloadUnavailable
	}
	if resp.ContentLength <= 0 {
		_ = resp.Body.Close()
		return nil, ErrImageDownloadInvalidLength
	}
	if resp.ContentLength > ImageDownloadMaxBytes {
		_ = resp.Body.Close()
		return nil, ErrImageDownloadTooLarge
	}

	prefix := make([]byte, min(int64(imageDownloadSniffBytes), resp.ContentLength))
	if _, err := io.ReadFull(resp.Body, prefix); err != nil {
		_ = resp.Body.Close()
		return nil, ErrImageDownloadUnavailable
	}
	contentType, err := validatedImageContentType(resp.Header.Get("Content-Type"), prefix)
	if err != nil {
		_ = resp.Body.Close()
		return nil, err
	}

	remaining := resp.ContentLength - int64(len(prefix))
	body := &prefixedImageReadCloser{
		Reader: io.MultiReader(bytes.NewReader(prefix), io.LimitReader(resp.Body, remaining)),
		closer: resp.Body,
	}
	return &ImageDownload{
		Body:          body,
		ContentType:   contentType,
		ContentLength: resp.ContentLength,
	}, nil
}

func validateImageDownloadURL(parsed *url.URL) error {
	if parsed == nil ||
		!strings.EqualFold(parsed.Scheme, "https") ||
		parsed.User != nil ||
		parsed.Port() != "" ||
		!strings.EqualFold(parsed.Hostname(), ImageDownloadAllowedHost) {
		return ErrImageDownloadInvalidURL
	}
	if parsed.Opaque != "" || parsed.Fragment != "" {
		return ErrImageDownloadInvalidURL
	}
	return nil
}

func validatedImageContentType(headerValue string, prefix []byte) (string, error) {
	headerType := ""
	if strings.TrimSpace(headerValue) != "" {
		parsedType, _, err := mime.ParseMediaType(strings.TrimSpace(headerValue))
		if err != nil {
			return "", ErrImageDownloadInvalidType
		}
		headerType = canonicalImageContentType(parsedType)
		if headerType != "application/octet-stream" && !isSupportedImageContentType(headerType) {
			return "", ErrImageDownloadInvalidType
		}
	}

	detectedType, _, err := mime.ParseMediaType(http.DetectContentType(prefix))
	if err != nil {
		return "", ErrImageDownloadInvalidType
	}
	detectedType = canonicalImageContentType(detectedType)
	if !isSupportedImageContentType(detectedType) {
		return "", ErrImageDownloadInvalidType
	}
	if isSupportedImageContentType(headerType) && headerType != detectedType {
		return "", ErrImageDownloadInvalidType
	}
	return detectedType, nil
}

func isSupportedImageContentType(contentType string) bool {
	switch canonicalImageContentType(contentType) {
	case "image/png", "image/jpeg", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

func canonicalImageContentType(contentType string) string {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpg", "image/pjpeg":
		return "image/jpeg"
	default:
		return strings.ToLower(strings.TrimSpace(contentType))
	}
}

type prefixedImageReadCloser struct {
	io.Reader
	closer io.Closer
}

func (r *prefixedImageReadCloser) Close() error {
	if r == nil || r.closer == nil {
		return nil
	}
	return r.closer.Close()
}

func (d *ImageDownload) Close() error {
	if d == nil || d.Body == nil {
		return nil
	}
	return d.Body.Close()
}

func (d *ImageDownload) Valid() bool {
	return d != nil &&
		d.Body != nil &&
		d.ContentLength > 0 &&
		d.ContentLength <= ImageDownloadMaxBytes &&
		isSupportedImageContentType(d.ContentType)
}
