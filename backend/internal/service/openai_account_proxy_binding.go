package service

import (
	"context"
	"net/http"
	"net/netip"
	"strings"
	"unicode"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/proxyurl"
)

// resolveOpenAIAccountProxyURL keeps an explicit account binding from silently
// becoming the default egress when a relation or repository lookup is missing.
// Callers pass the account that owns the request egress: quota/refresh may use
// the credential parent, while inference preserves its selected account egress.
// An unbound account retains its existing default-egress behavior. Expiry and
// administrator-configured fallback policy remain the repository's concern.
func resolveOpenAIAccountProxyURL(ctx context.Context, account *Account, proxyRepo ProxyRepository) (string, error) {
	unavailable := func(reason string) (string, error) {
		// Do not expose repository errors or proxy URLs: either may contain secrets.
		return "", infraerrors.Newf(http.StatusBadGateway, "OPENAI_ACCOUNT_PROXY_UNAVAILABLE", "bound OpenAI account proxy is unavailable: %s", reason)
	}
	if account == nil {
		return unavailable("account is missing")
	}
	if account.ProxyID == nil {
		return "", nil
	}
	if *account.ProxyID <= 0 {
		return unavailable("invalid proxy binding")
	}
	proxy := account.Proxy
	if proxy == nil || proxy.ID != *account.ProxyID {
		if proxyRepo == nil {
			return unavailable("bound proxy is not loaded")
		}
		var err error
		proxy, err = proxyRepo.GetByID(ctx, *account.ProxyID)
		if err != nil {
			return unavailable("proxy lookup failed")
		}
		if proxy == nil || proxy.ID != *account.ProxyID {
			return unavailable("bound proxy was not found")
		}
	}
	if proxy.Port < 1 || proxy.Port > 65535 {
		return unavailable("invalid proxy port")
	}
	// Proxy.URL uses net.JoinHostPort; reject malformed stored hosts before it can
	// turn an empty host or embedded port into an unintended endpoint.
	if proxy.Host == "" || strings.ContainsAny(proxy.Host, "/\\?#@[]") || strings.IndexFunc(proxy.Host, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }) >= 0 {
		return unavailable("invalid proxy host")
	}
	if strings.Contains(proxy.Host, ":") {
		if _, err := netip.ParseAddr(proxy.Host); err != nil {
			return unavailable("invalid proxy host")
		}
	} else if strings.Contains(proxy.Host, "%") {
		return unavailable("invalid proxy host")
	}
	resolvedURL, parsed, err := proxyurl.Parse(proxy.URL())
	if err != nil || parsed == nil {
		return unavailable("invalid proxy configuration")
	}
	return resolvedURL, nil
}
