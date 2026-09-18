//go:build unit

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type codexAuxiliaryMetadataUpstream struct {
	httpUpstreamRecorder
	lastAccountConcurrency int
}

func (u *codexAuxiliaryMetadataUpstream) Do(req *http.Request, proxy string, accountID int64, concurrency int) (*http.Response, error) {
	u.lastAccountConcurrency = concurrency
	return u.httpUpstreamRecorder.Do(req, proxy, accountID, concurrency)
}
func (u *codexAuxiliaryMetadataUpstream) DoWithTLS(req *http.Request, proxy string, accountID int64, concurrency int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.Do(req, proxy, accountID, concurrency)
}

func TestCodexTurnMetadataAuxiliaryPaths(t *testing.T) {
	for _, path := range []string{"probe", "compact_probe", "search_pat", "search_standalone"} {
		for _, mode := range []string{"off", "device", "session", "full"} {
			t.Run(path+"/"+mode, func(t *testing.T) {
				account := newTestOAuthAccount(8801, map[string]any{codexFingerprintModeExtraKey: mode})
				account.Concurrency = 3
				account.Credentials = map[string]any{"access_token": "test-token", "chatgpt_account_id": "test-account", "chatgpt_account_is_fedramp": false}
				if path == "search_pat" {
					account.Credentials["auth_mode"] = OpenAIAuthModePersonalAccessToken
				}
				var first map[string]any
				for turn := 0; turn < 2; turn++ {
					response := "data: {\"type\":\"response.completed\"}\n\n"
					switch path {
					case "compact_probe":
						response = compactProbeSSESuccessBody
					case "search_pat":
						response = alphaSearchResponsesSSE("search result")
					case "search_standalone":
						response = `{"output":"search result"}`
					}
					upstream := &codexAuxiliaryMetadataUpstream{httpUpstreamRecorder: httpUpstreamRecorder{resp: &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(response))}}}
					c, _ := newTestContext()
					c.Request.Header.Set("session-id", fmt.Sprintf("client-%d", turn))
					alphaBody := []byte(`{"id":"search-id","model":"gpt-5.4","commands":{"search_query":[{"q":"test"}]}}`)
					var err error
					if strings.HasPrefix(path, "search") {
						svc := &OpenAIGatewayService{cfg: &config.Config{}, httpUpstream: upstream}
						_, err = svc.ForwardAlphaSearch(context.Background(), c, account, alphaBody)
					} else {
						svc := &AccountTestService{httpUpstream: upstream}
						if path == "compact_probe" {
							err = svc.testOpenAICompactConnection(c, account, "gpt-5.4")
						} else {
							err = svc.testOpenAIAccountConnection(c, account, "gpt-5.4", "", "")
						}
					}
					require.NoError(t, err)
					require.NotNil(t, upstream.lastReq)
					require.Equal(t, 3, upstream.lastAccountConcurrency)
					h := upstream.lastReq.Header
					if path == "search_standalone" {
						require.JSONEq(t, string(alphaBody), string(upstream.lastBody))
						for _, key := range []string{"OpenAI-Beta", "session_id", "session-id", "conversation_id", "thread-id", "x-codex-window-id", "x-codex-turn-state"} {
							require.Empty(t, h.Get(key))
						}
					}
					if mode == "off" || mode == "device" {
						require.Empty(t, h.Get("x-codex-turn-metadata"))
						require.False(t, gjson.GetBytes(upstream.lastBody, "client_metadata.x-codex-turn-metadata").Exists())
						continue
					}
					var metadata map[string]any
					require.NoError(t, json.Unmarshal([]byte(h.Get("x-codex-turn-metadata")), &metadata))
					require.Len(t, metadata, 6)
					require.NotEmpty(t, metadata["turn_id"])
					require.Greater(t, metadata["turn_started_at_unix_ms"].(float64), float64(0))
					if path != "search_standalone" {
						require.Equal(t, h.Get("x-codex-turn-metadata"), gjson.GetBytes(upstream.lastBody, "client_metadata.x-codex-turn-metadata").String())
						requireCodexFingerprintHeaderBodyParity(t, h, upstream.lastBody, metadata["installation_id"].(string), metadata["session_id"].(string), metadata["thread_id"].(string))
						if conversation := h.Get("conversation_id"); conversation != "" {
							require.Equal(t, metadata["session_id"], conversation)
						}
					}
					if first != nil {
						require.Equal(t, first["installation_id"], metadata["installation_id"])
						require.Equal(t, first["session_id"], metadata["session_id"])
						require.NotEqual(t, first["turn_id"], metadata["turn_id"])
						if mode == "full" {
							require.Equal(t, first["thread_id"], metadata["thread_id"])
						}
					}
					first = metadata
				}
			})
		}
	}
}
