//go:build unit

package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProviderHallGatewayCheckFixedPathNoCredentialsOrRedirect(t *testing.T) {
	t.Setenv("PROVIDER_HALL_ALLOW_LOOPBACK", "1")
	var hits atomic.Int32
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		require.Equal(t, "/health", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		require.Empty(t, r.Header.Get("Authorization"))
		http.Redirect(w, r, "/v1/responses", http.StatusFound)
	}))
	defer s.Close()
	_, err := ProviderHallCheckGateway(context.Background(), s.URL)
	require.Error(t, err)
	require.EqualValues(t, 1, hits.Load())
	for _, origin := range []string{s.URL + "/v1/responses", s.URL + "?path=other", "http://user:pass@localhost", s.URL + "#fragment"} {
		_, err := ProviderHallCheckGateway(context.Background(), origin)
		require.Error(t, err)
	}
	require.EqualValues(t, 1, hits.Load())
	t.Setenv("PROVIDER_HALL_ALLOW_LOOPBACK", "")
	_, err = ProviderHallCheckGateway(context.Background(), s.URL)
	require.Error(t, err)
	require.EqualValues(t, 1, hits.Load())
}
