package repository

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/require"
)

func TestGetSharedReqClient_DoesNotShareResponseCookiesBetweenAccounts(t *testing.T) {
	sharedReqClients = sync.Map{}
	requests := make(chan http.Header, 2)
	server := newLocalTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- r.Header.Clone()
		http.SetCookie(w, &http.Cookie{Name: "upstream_session", Value: "account-a-session", Path: "/"})
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	opts := reqClientOptions{Timeout: time.Second}
	accountAClient, err := getSharedReqClient(opts)
	require.NoError(t, err)
	accountBClient, err := getSharedReqClient(opts)
	require.NoError(t, err)
	require.Same(t, accountAClient, accountBClient)

	for _, account := range []struct {
		client *req.Client
		token  string
	}{
		{accountAClient, "account-a-token"},
		{accountBClient, "account-b-token"},
	} {
		response, requestErr := account.client.R().
			SetHeader("Authorization", "Bearer "+account.token).
			Get(server.URL)
		require.NoError(t, requestErr)
		require.Equal(t, http.StatusNoContent, response.StatusCode)
		header := <-requests
		require.Equal(t, "Bearer "+account.token, header.Get("Authorization"))
		require.Empty(t, header.Values("Cookie"), "a shared client must not replay an upstream account cookie")
	}
}

func TestGetSharedReqClient_ExplicitCookiesRemainRequestScoped(t *testing.T) {
	for _, tc := range []struct {
		name  string
		apply func(*req.Request) *req.Request
	}{
		{
			name: "request cookie",
			apply: func(r *req.Request) *req.Request {
				return r.SetCookies(&http.Cookie{Name: "sessionKey", Value: "explicit-session"})
			},
		},
		{
			name: "cookie header",
			apply: func(r *req.Request) *req.Request {
				return r.SetHeader("Cookie", "sessionKey=explicit-session")
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sharedReqClients = sync.Map{}
			requests := make(chan http.Header, 2)
			server := newLocalTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests <- r.Header.Clone()
				w.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()
			client, err := getSharedReqClient(reqClientOptions{Timeout: time.Second})
			require.NoError(t, err)

			response, err := tc.apply(client.R()).Get(server.URL)
			require.NoError(t, err)
			require.Equal(t, http.StatusNoContent, response.StatusCode)
			require.Equal(t, "sessionKey=explicit-session", (<-requests).Get("Cookie"))

			response, err = client.R().Get(server.URL)
			require.NoError(t, err)
			require.Equal(t, http.StatusNoContent, response.StatusCode)
			require.Empty(t, (<-requests).Values("Cookie"))
		})
	}
}
