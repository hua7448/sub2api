package service

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageDownloadPublicIPPolicyRejectsPrivateAndSpecialUse(t *testing.T) {
	blocked := []string{
		"0.0.0.1",
		"10.0.0.1",
		"100.64.0.1",
		"127.0.0.1",
		"169.254.169.254",
		"172.16.0.1",
		"192.0.0.1",
		"192.0.2.1",
		"192.168.0.1",
		"198.18.0.1",
		"198.51.100.1",
		"203.0.113.1",
		"224.0.0.1",
		"240.0.0.1",
		"::1",
		"64:ff9b::1",
		"100::1",
		"100:0:0:1::1",
		"2001:db8::1",
		"2002::1",
		"3fff::1",
		"5f00::1",
		"fc00::1",
		"fe80::1",
		"ff00::1",
	}
	for _, rawIP := range blocked {
		t.Run(rawIP, func(t *testing.T) {
			require.False(t, isImageDownloadPublicIP(net.ParseIP(rawIP)))
		})
	}

	require.True(t, isImageDownloadPublicIP(net.ParseIP("8.8.8.8")))
	require.True(t, isImageDownloadPublicIP(net.ParseIP("2606:4700:4700::1111")))
}

func TestImageDownloadSafeDialRejectsPrivateDNSAnswerBeforeDial(t *testing.T) {
	lookupCalls := 0
	dialCalls := 0
	lookup := func(_ context.Context, host string) ([]net.IPAddr, error) {
		lookupCalls++
		require.Equal(t, ImageDownloadAllowedHost, host)
		return []net.IPAddr{{IP: net.ParseIP("127.0.0.1")}}, nil
	}
	dial := func(context.Context, string, string) (net.Conn, error) {
		dialCalls++
		return nil, errors.New("must not dial")
	}

	conn, err := imageDownloadSafeDialContextWith(
		context.Background(),
		"tcp",
		net.JoinHostPort(ImageDownloadAllowedHost, "443"),
		lookup,
		dial,
	)

	require.Nil(t, conn)
	require.Error(t, err)
	require.Contains(t, err.Error(), "blocked by image download SSRF policy")
	require.Equal(t, 1, lookupCalls)
	require.Zero(t, dialCalls)
}

func TestImageDownloadSafeDialRejectsMixedPublicPrivateAnswers(t *testing.T) {
	dialCalls := 0
	lookup := func(context.Context, string) ([]net.IPAddr, error) {
		return []net.IPAddr{
			{IP: net.ParseIP("8.8.8.8")},
			{IP: net.ParseIP("10.0.0.1")},
		}, nil
	}
	dial := func(context.Context, string, string) (net.Conn, error) {
		dialCalls++
		return nil, errors.New("must not dial")
	}

	conn, err := imageDownloadSafeDialContextWith(
		context.Background(),
		"tcp",
		net.JoinHostPort(ImageDownloadAllowedHost, "443"),
		lookup,
		dial,
	)

	require.Nil(t, conn)
	require.Error(t, err)
	require.Zero(t, dialCalls)
}
