package service

import (
	"context"
	"net"
	"net/http"
	"net/netip"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/servertiming"
)

const (
	imageDownloadDialTimeout           = 10 * time.Second
	imageDownloadDialKeepAlive         = 30 * time.Second
	imageDownloadIdleConnTimeout       = 30 * time.Second
	imageDownloadTLSHandshakeTimeout   = 10 * time.Second
	imageDownloadResponseHeaderTimeout = 30 * time.Second
)

var imageDownloadBlockedPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.31.196.0/24"),
	netip.MustParsePrefix("192.52.193.0/24"),
	netip.MustParsePrefix("192.88.99.0/24"),
	netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("192.175.48.0/24"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("224.0.0.0/4"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("64:ff9b::/96"),
	netip.MustParsePrefix("64:ff9b:1::/48"),
	netip.MustParsePrefix("100::/64"),
	netip.MustParsePrefix("100:0:0:1::/64"),
	netip.MustParsePrefix("2001::/23"),
	netip.MustParsePrefix("2001:20::/28"),
	netip.MustParsePrefix("2001:db8::/32"),
	netip.MustParsePrefix("2002::/16"),
	netip.MustParsePrefix("2620:4f:8000::/48"),
	netip.MustParsePrefix("3fff::/20"),
	netip.MustParsePrefix("5f00::/16"),
	netip.MustParsePrefix("fc00::/7"),
	netip.MustParsePrefix("fe80::/10"),
	netip.MustParsePrefix("ff00::/8"),
}

var imageDownloadDialer = &net.Dialer{
	Timeout:   imageDownloadDialTimeout,
	KeepAlive: imageDownloadDialKeepAlive,
}

type imageDownloadLookupIPFunc func(context.Context, string) ([]net.IPAddr, error)
type imageDownloadDialIPFunc func(context.Context, string, string) (net.Conn, error)

func newImageDownloadHTTPClient(timeout time.Duration) *http.Client {
	transport := &http.Transport{
		DialContext:           imageDownloadSafeDialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          16,
		IdleConnTimeout:       imageDownloadIdleConnTimeout,
		TLSHandshakeTimeout:   imageDownloadTLSHandshakeTimeout,
		ResponseHeaderTimeout: imageDownloadResponseHeaderTimeout,
	}
	return &http.Client{
		Timeout:   timeout,
		Transport: servertiming.WrapRoundTripper(transport),
	}
}

func imageDownloadSafeDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return imageDownloadSafeDialContextWith(
		ctx,
		network,
		address,
		net.DefaultResolver.LookupIPAddr,
		imageDownloadDialer.DialContext,
	)
}

func imageDownloadSafeDialContextWith(
	ctx context.Context,
	network string,
	address string,
	lookup imageDownloadLookupIPFunc,
	dial imageDownloadDialIPFunc,
) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	var addresses []net.IPAddr
	if ip := net.ParseIP(host); ip != nil {
		addresses = []net.IPAddr{{IP: ip}}
	} else {
		addresses, err = lookup(ctx, host)
		if err != nil {
			return nil, err
		}
	}
	if len(addresses) == 0 {
		return nil, &net.AddrError{Err: "no addresses for host", Addr: host}
	}

	// Reject the whole resolution set if any answer is special-use. This keeps
	// mixed public/private DNS answers from becoming a rebinding bypass.
	for _, resolved := range addresses {
		if !isImageDownloadPublicIP(resolved.IP) {
			return nil, &net.AddrError{Err: "blocked by image download SSRF policy", Addr: resolved.IP.String()}
		}
	}

	var lastErr error
	for _, resolved := range addresses {
		conn, dialErr := dial(ctx, network, net.JoinHostPort(resolved.IP.String(), port))
		if dialErr == nil {
			return conn, nil
		}
		lastErr = dialErr
	}
	if lastErr == nil {
		lastErr = &net.AddrError{Err: "no usable addresses", Addr: host}
	}
	return nil, lastErr
}

func isImageDownloadPublicIP(ip net.IP) bool {
	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return false
	}
	addr = addr.Unmap()
	if !addr.IsValid() || !addr.IsGlobalUnicast() || addr.IsPrivate() {
		return false
	}
	for _, prefix := range imageDownloadBlockedPrefixes {
		if prefix.Contains(addr) {
			return false
		}
	}
	return true
}
