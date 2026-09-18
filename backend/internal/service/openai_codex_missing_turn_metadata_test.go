package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	coderws "github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func requireGeneratedCodexTurnMetadata(t *testing.T, body []byte) map[string]any {
	t.Helper()
	raw := gjson.GetBytes(body, "client_metadata.x-codex-turn-metadata")
	require.Equal(t, gjson.String, raw.Type, "missing turn metadata must be generated as a JSON string")
	var metadata map[string]any
	require.NoError(t, json.Unmarshal([]byte(raw.String()), &metadata))
	for _, pair := range [][2]string{
		{"installation_id", "x-codex-installation-id"},
		{"session_id", "session_id"},
		{"thread_id", "thread_id"},
		{"turn_id", "turn_id"},
		{"window_id", "x-codex-window-id"},
	} {
		value := gjson.GetBytes(body, "client_metadata."+pair[1]).String()
		require.NotEmpty(t, value, "missing flat identity %s", pair[1])
		require.Equal(t, value, metadata[pair[0]], "flat and embedded %s must agree", pair[0])
	}
	timestamp, ok := metadata["turn_started_at_unix_ms"].(float64)
	require.True(t, ok, "turn timestamp must be a JSON number")
	require.Greater(t, timestamp, float64(0))
	return metadata
}

func requireGeneratedCodexTurnHeader(t *testing.T, headers http.Header, metadata map[string]any) {
	t.Helper()
	var headerMetadata map[string]any
	require.NotEmpty(t, headers.Get(openAIWSTurnMetadataHeader), "missing HTTP/WS handshake metadata must be generated")
	require.NoError(t, json.Unmarshal([]byte(headers.Get(openAIWSTurnMetadataHeader)), &headerMetadata))
	require.Equal(t, metadata, headerMetadata)
	require.Equal(t, metadata["installation_id"], headers.Get("x-codex-installation-id"))
	require.Equal(t, metadata["session_id"], headers.Get("session-id"))
	require.Equal(t, metadata["thread_id"], headers.Get("thread-id"))
	require.Equal(t, metadata["window_id"], headers.Get("x-codex-window-id"))
}

func TestCodexMissingTurnMetadataForwardPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, transport := range []string{"http", "http_passthrough", "websocket"} {
		for _, mode := range []string{"off", "device", "session", "full"} {
			t.Run(transport+"/"+mode, func(t *testing.T) {
				var accountMetadata []map[string]any
				for accountIndex, seed := range []string{
					"11111111-1111-4111-8111-111111111111",
					"22222222-2222-4222-8222-222222222222",
				} {
					useWS := transport == "websocket"
					cfg := &config.Config{}
					cfg.Security.URLAllowlist.Enabled = false
					cfg.Gateway.OpenAIWS.Enabled = useWS
					cfg.Gateway.OpenAIWS.OAuthEnabled = useWS
					cfg.Gateway.OpenAIWS.ResponsesWebsocketsV2 = useWS
					cfg.Gateway.OpenAIWS.MaxConnsPerAccount = 3
					cfg.Gateway.OpenAIWS.MaxIdlePerAccount = 1
					account := newTestOAuthAccount(int64(1600+accountIndex), map[string]any{
						codexFingerprintModeExtraKey: mode, codexFingerprintSeedExtraKey: seed,
						"openai_oauth_passthrough":        transport == "http_passthrough",
						"responses_websockets_v2_enabled": useWS,
					})
					account.Concurrency = 3
					account.Status = StatusActive
					account.Schedulable = true
					account.Credentials = map[string]any{"access_token": fmt.Sprintf("token-%d", accountIndex), "chatgpt_account_id": fmt.Sprintf("account-%d", accountIndex)}
					upstream := &httpUpstreamRecorder{responses: []*http.Response{
						openAICompatSSECompletedResponse("resp_missing_1", "gpt-5.1"),
						openAICompatSSECompletedResponse("resp_missing_2", "gpt-5.1"),
					}}
					conn := &openAIWSCaptureConn{events: [][]byte{
						[]byte(`{"type":"response.completed","response":{"id":"resp_missing_1","model":"gpt-5.1","usage":{"input_tokens":1,"output_tokens":1}}}`),
						[]byte(`{"type":"response.completed","response":{"id":"resp_missing_2","model":"gpt-5.1","usage":{"input_tokens":1,"output_tokens":1}}}`),
					}}
					dialer := &openAIWSCaptureDialer{conn: conn}
					pool := newOpenAIWSConnPool(cfg)
					defer pool.Close()
					pool.setClientDialerForTest(dialer)
					svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstream, cache: &stubGatewayCache{}, openaiWSResolver: NewOpenAIWSProtocolResolver(cfg), toolCorrector: NewCodexToolCorrector(), openaiWSPool: pool}
					body := []byte(`{"model":"gpt-5.1","instructions":"Help the user.","stream":false,"input":[{"type":"message","role":"user","content":"hello"}]}`)
					var turns []map[string]any
					for turn := 0; turn < 2; turn++ {
						c, _ := gin.CreateTestContext(httptest.NewRecorder())
						c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
						c.Request.Header.Set("User-Agent", CodexCanonicalUserAgent())
						c.Request.Header.Set("session-id", "same-real-client-session")
						require.Empty(t, c.Request.Header.Get(openAIWSTurnMetadataHeader))
						_, err := svc.Forward(context.Background(), c, account, body)
						require.NoError(t, err)
						var outbound []byte
						var headers http.Header
						if useWS {
							require.Empty(t, upstream.requests, "WS must not fall back to HTTP")
							require.Len(t, conn.writes, turn+1)
							outbound = []byte(requestToJSONString(conn.writes[turn]))
							headers = dialer.lastHeaders
						} else {
							require.Len(t, upstream.requests, turn+1)
							outbound = upstream.bodies[turn]
							headers = upstream.requests[turn].Header
						}
						if mode == "off" || mode == "device" {
							require.Empty(t, headers.Get(openAIWSTurnMetadataHeader))
							require.False(t, gjson.GetBytes(outbound, "client_metadata.x-codex-turn-metadata").Exists())
							continue
						}
						metadata := requireGeneratedCodexTurnMetadata(t, outbound)
						// An existing WS handshake is not retransmitted on later turns.
						if !useWS || turn == 0 {
							requireGeneratedCodexTurnHeader(t, headers, metadata)
						}
						turns = append(turns, metadata)
					}
					if useWS {
						require.Equal(t, 1, dialer.DialCount(), "metadata synthesis must retain WS connection reuse")
					}
					require.Equal(t, 3, account.Concurrency)
					if len(turns) == 2 {
						require.NotEqual(t, turns[0]["turn_id"], turns[1]["turn_id"])
						for _, key := range []string{"installation_id", "session_id", "thread_id", "window_id"} {
							require.Equal(t, turns[0][key], turns[1][key])
						}
						accountMetadata = append(accountMetadata, turns[0])
					}
				}
				if len(accountMetadata) == 2 {
					for _, key := range []string{"installation_id", "session_id", "thread_id", "turn_id"} {
						require.NotEqual(t, accountMetadata[0][key], accountMetadata[1][key], "%s must be isolated by account", key)
					}
				}
			})
		}
	}
}

func TestCodexMissingTurnMetadataWebSocketIngress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, ingress := range []string{"native", "passthrough"} {
		for _, mode := range []string{"off", "device", "session", "full"} {
			t.Run(ingress+"/"+mode, func(t *testing.T) {
				var accountMetadata []map[string]any
				for accountIndex, seed := range []string{
					"11111111-1111-4111-8111-111111111111",
					"22222222-2222-4222-8222-222222222222",
				} {
					cfg := &config.Config{}
					cfg.Security.URLAllowlist.Enabled = false
					cfg.Security.URLAllowlist.AllowInsecureHTTP = true
					cfg.Gateway.OpenAIWS.Enabled = true
					cfg.Gateway.OpenAIWS.OAuthEnabled = true
					cfg.Gateway.OpenAIWS.ResponsesWebsocketsV2 = true
					cfg.Gateway.OpenAIWS.MaxConnsPerAccount = 3
					cfg.Gateway.OpenAIWS.MaxIdlePerAccount = 1
					cfg.Gateway.OpenAIWS.QueueLimitPerConn = 8
					cfg.Gateway.OpenAIWS.DialTimeoutSeconds = 3
					cfg.Gateway.OpenAIWS.ReadTimeoutSeconds = 3
					cfg.Gateway.OpenAIWS.WriteTimeoutSeconds = 3
					account := newTestOAuthAccount(int64(1700+accountIndex), map[string]any{
						codexFingerprintModeExtraKey: mode, codexFingerprintSeedExtraKey: seed,
						"responses_websockets_v2_enabled": true,
					})
					account.Concurrency = 3
					account.Status = StatusActive
					account.Schedulable = true
					account.Credentials = map[string]any{"access_token": fmt.Sprintf("token-%d", accountIndex), "chatgpt_account_id": fmt.Sprintf("account-%d", accountIndex)}
					captureConn := &openAIWSCaptureConn{events: [][]byte{
						[]byte(`{"type":"response.completed","response":{"id":"resp_missing_ingress_1","model":"gpt-5.1","usage":{"input_tokens":1,"output_tokens":1}}}`),
						[]byte(`{"type":"response.completed","response":{"id":"resp_missing_ingress_2","model":"gpt-5.1","usage":{"input_tokens":1,"output_tokens":1}}}`),
					}}
					dialer := &openAIWSCaptureDialer{conn: captureConn}
					pool := newOpenAIWSConnPool(cfg)
					defer pool.Close()
					pool.setClientDialerForTest(dialer)
					upstream := &httpUpstreamRecorder{}
					svc := &OpenAIGatewayService{cfg: cfg, httpUpstream: upstream, cache: &stubGatewayCache{}, openaiWSResolver: NewOpenAIWSProtocolResolver(cfg), toolCorrector: NewCodexToolCorrector(), openaiWSPool: pool}
					if ingress == "passthrough" {
						cfg.Gateway.OpenAIWS.ModeRouterV2Enabled = true
						cfg.Gateway.OpenAIWS.IngressModeDefault = OpenAIWSIngressModeCtxPool
						account.Extra["openai_oauth_responses_websockets_v2_mode"] = OpenAIWSIngressModePassthrough
						svc.openaiWSPassthroughDialer = dialer
						// Keep the relay from consuming the next event before the next
						// response.create frame is sent by the real downstream client.
						captureConn.readDelays = []time.Duration{50 * time.Millisecond, 250 * time.Millisecond}
					}

					serverErrCh := make(chan error, 1)
					wsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						conn, err := coderws.Accept(w, r, &coderws.AcceptOptions{CompressionMode: coderws.CompressionContextTakeover})
						if err != nil {
							serverErrCh <- err
							return
						}
						defer conn.CloseNow()
						ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
						ginCtx.Request = r.Clone(r.Context())
						readCtx, cancelRead := context.WithTimeout(r.Context(), 3*time.Second)
						_, firstMessage, readErr := conn.Read(readCtx)
						cancelRead()
						if readErr != nil {
							serverErrCh <- readErr
							return
						}
						serverErrCh <- svc.ProxyResponsesWebSocketFromClient(r.Context(), ginCtx, conn, account, account.GetCredential("access_token"), firstMessage, nil)
					}))
					defer wsServer.Close()
					clientHeaders := make(http.Header)
					clientHeaders.Set("User-Agent", CodexCanonicalUserAgent())
					clientHeaders.Set("session-id", "same-real-client-session")
					require.Empty(t, clientHeaders.Get(openAIWSTurnMetadataHeader))
					dialCtx, cancelDial := context.WithTimeout(context.Background(), 3*time.Second)
					clientConn, _, err := coderws.Dial(dialCtx, "ws"+strings.TrimPrefix(wsServer.URL, "http"), &coderws.DialOptions{HTTPHeader: clientHeaders})
					cancelDial()
					require.NoError(t, err)
					defer clientConn.CloseNow()
					for turn, payload := range []string{
						`{"type":"response.create","model":"gpt-5.1","stream":false,"input":[{"type":"message","role":"user","content":"hello"}]}`,
						`{"type":"response.create","model":"gpt-5.1","stream":false,"input":[{"type":"message","role":"user","content":"continue"}],"previous_response_id":"resp_missing_ingress_1"}`,
					} {
						writeCtx, cancelWrite := context.WithTimeout(context.Background(), 3*time.Second)
						err := clientConn.Write(writeCtx, coderws.MessageText, []byte(payload))
						cancelWrite()
						require.NoError(t, err)
						readCtx, cancelRead := context.WithTimeout(context.Background(), 3*time.Second)
						msgType, event, err := clientConn.Read(readCtx)
						cancelRead()
						require.NoError(t, err)
						require.Equal(t, coderws.MessageText, msgType)
						require.Equal(t, fmt.Sprintf("resp_missing_ingress_%d", turn+1), gjson.GetBytes(event, "response.id").String())
					}
					// A relay may already close after its final fake upstream event.
					_ = clientConn.Close(coderws.StatusNormalClosure, "done")
					select {
					case serverErr := <-serverErrCh:
						if ingress == "passthrough" && serverErr != nil {
							require.True(t, strings.Contains(serverErr.Error(), "StatusNormalClosure") || strings.Contains(serverErr.Error(), "connection closed"), serverErr)
						} else {
							require.NoError(t, serverErr)
						}
					case <-time.After(5 * time.Second):
						t.Fatal("websocket ingress did not finish")
					}
					require.Equal(t, 1, dialer.DialCount(), "two turns must reuse the same upstream connection")
					require.Empty(t, upstream.requests, "native WS must not fall back to HTTP")
					require.Equal(t, 3, account.Concurrency)
					require.Len(t, captureConn.writes, 2)
					var turns []map[string]any
					for turn, request := range captureConn.writes {
						body := []byte(requestToJSONString(request))
						if mode == "off" || mode == "device" {
							require.Empty(t, dialer.lastHeaders.Get(openAIWSTurnMetadataHeader))
							require.False(t, gjson.GetBytes(body, "client_metadata.x-codex-turn-metadata").Exists())
							continue
						}
						metadata := requireGeneratedCodexTurnMetadata(t, body)
						if turn == 0 {
							requireGeneratedCodexTurnHeader(t, dialer.lastHeaders, metadata)
						}
						turns = append(turns, metadata)
					}
					if len(turns) == 2 {
						require.NotEqual(t, turns[0]["turn_id"], turns[1]["turn_id"])
						for _, key := range []string{"installation_id", "session_id", "thread_id", "window_id"} {
							require.Equal(t, turns[0][key], turns[1][key])
						}
						accountMetadata = append(accountMetadata, turns[0])
					}
				}
				if len(accountMetadata) == 2 {
					for _, key := range []string{"installation_id", "session_id", "thread_id", "turn_id"} {
						require.NotEqual(t, accountMetadata[0][key], accountMetadata[1][key], "%s must be isolated by account", key)
					}
				}
			})
		}
	}
}
