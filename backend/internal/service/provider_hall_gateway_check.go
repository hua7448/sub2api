package service

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Only a fixed, unauthenticated health endpoint is allowed here.
func ProviderHallCheckGateway(ctx context.Context, origin string) (int, error) {
	origin = strings.TrimSuffix(strings.TrimSpace(origin), "/")
	u, err := url.Parse(origin)
	if err != nil || len(origin) > 2048 || u.Hostname() == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Opaque != "" || strings.Contains(origin, "#") {
		return 0, providerHallInvalid("gateway_origin")
	}
	loopback := ProviderHallIsLoopbackOrigin(origin)
	if loopback && !ProviderHallAllowLoopback() || !loopback && u.Scheme != "https" || loopback && u.Scheme != "https" && u.Scheme != "http" {
		return 0, providerHallInvalid("gateway_origin")
	}
	client := NewProviderHallProbeClient(5 * time.Second)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, origin+"/health", nil)
	if err != nil {
		return 0, providerHallInvalid("gateway_origin")
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return 0, providerHallInvalid("gateway_connection")
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	return resp.StatusCode, nil
}
