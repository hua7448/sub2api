//go:build unit

package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const hallResponsesSSE = `event: response.created
data: {"type":"response.created","response":{"id":"resp_1","model":"gpt-x"}}

event: response.output_text.delta
data: {"type":"response.output_text.delta","delta":"12"}

event: response.output_text.delta
data: {"type":"response.output_text.delta","delta":"3456"}

event: response.completed
data: {"type":"response.completed","response":{"id":"resp_1","model":"gpt-x-2025","status":"completed","output":[{"type":"message"}],"usage":{"input_tokens":21,"output_tokens":4}}}

`

const hallResponsesToolSSE = `event: response.output_item.added
data: {"type":"response.output_item.added","item":{"type":"function_call","name":"hall_echo"}}

event: response.function_call_arguments.delta
data: {"type":"response.function_call_arguments.delta","delta":"{\"code\""}

event: response.completed
data: {"type":"response.completed","response":{"model":"gpt-x","output":[{"type":"function_call","name":"hall_echo","arguments":"{\"code\":\"zz\"}"}],"usage":{"input_tokens":30,"output_tokens":9}}}

`

const hallChatSSE = `data: {"id":"c1","model":"gpt-x","choices":[{"index":0,"delta":{"role":"assistant","content":""}}]}

data: {"id":"c1","model":"gpt-x","choices":[{"index":0,"delta":{"content":"{\"token\":"}}]}

data: {"id":"c1","model":"gpt-x","choices":[{"index":0,"delta":{"content":"\"abc12345\",\"length\":8}"}}]}

data: {"id":"c1","model":"gpt-x","choices":[],"usage":{"prompt_tokens":40,"completion_tokens":12}}

data: [DONE]

`

const hallChatToolSSE = `data: {"id":"c2","model":"gpt-x","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"hall_echo","arguments":""}}]}}]}

data: {"id":"c2","model":"gpt-x","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"{\"code\":"}}]}}]}

data: {"id":"c2","model":"gpt-x","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"function":{"arguments":"\"zz\"}"}}]}}]}

data: [DONE]

`

const hallMessagesSSE = `event: message_start
data: {"type":"message_start","message":{"id":"m1","model":"claude-x","usage":{"input_tokens":17,"output_tokens":1}}}

event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"OK"}}

event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":3}}

event: message_stop
data: {"type":"message_stop"}

`

const hallMessagesToolSSE = `event: message_start
data: {"type":"message_start","message":{"model":"claude-x","usage":{"input_tokens":17}}}

event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"tu_1","name":"hall_echo","input":{}}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{\"code\": \"z"}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"z\"}"}}

event: message_delta
data: {"type":"message_delta","usage":{"output_tokens":6}}

event: message_stop
data: {"type":"message_stop"}

`

func hallClock(start time.Time, steps ...time.Duration) func() time.Time {
	i := 0
	return func() time.Time {
		if i >= len(steps) {
			return start.Add(steps[len(steps)-1])
		}
		d := steps[i]
		i++
		return start.Add(d)
	}
}

func TestProviderHallParseStreamFixtures(t *testing.T) {
	start := time.Date(2026, 9, 12, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name, protocol, fixture string
		text, model             string
		in, out                 int
		tool                    *ProviderHallToolCall
	}{
		{"responses_text", "responses", hallResponsesSSE, "123456", "gpt-x-2025", 21, 4, nil},
		{"responses_tool", "responses", hallResponsesToolSSE, "", "gpt-x", 30, 9, &ProviderHallToolCall{Name: "hall_echo", Arguments: `{"code":"zz"}`}},
		{"chat_text", "chat_completions", hallChatSSE, `{"token":"abc12345","length":8}`, "gpt-x", 40, 12, nil},
		{"chat_tool", "chat_completions", hallChatToolSSE, "", "gpt-x", 0, 0, &ProviderHallToolCall{Name: "hall_echo", Arguments: `{"code":"zz"}`}},
		{"messages_text", "messages", hallMessagesSSE, "OK", "claude-x", 17, 3, nil},
		{"messages_tool", "messages", hallMessagesToolSSE, "", "claude-x", 17, 6, &ProviderHallToolCall{Name: "hall_echo", Arguments: `{"code": "zz"}`}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// First content lands at +120ms, the end at +900ms.
			resp, err := ProviderHallParseStream(tc.protocol, strings.NewReader(tc.fixture), start, hallClock(start, 120*time.Millisecond, 900*time.Millisecond))
			require.NoError(t, err)
			require.True(t, resp.Terminal)
			require.Empty(t, resp.ErrorCode)
			require.Equal(t, tc.text, resp.Text)
			require.Equal(t, tc.model, resp.Model)
			require.NotNil(t, resp.TTFTMs)
			require.Equal(t, 120, *resp.TTFTMs, "TTFT is the first non-empty content delta")
			require.NotNil(t, resp.TotalMs)
			require.Equal(t, 900, *resp.TotalMs)
			require.Equal(t, 780, *ProviderHallGenerationMs(resp.TotalMs, resp.TTFTMs))
			if tc.in > 0 {
				require.Equal(t, tc.in, *resp.InputTokens)
				require.Equal(t, tc.out, *resp.OutputTokens)
				require.InDelta(t, float64(tc.out)/0.78, *ProviderHallTPS(resp.OutputTokens, ProviderHallGenerationMs(resp.TotalMs, resp.TTFTMs)), 0.01)
			} else {
				require.Nil(t, resp.InputTokens)
			}
			if tc.tool != nil {
				require.Len(t, resp.ToolCalls, 1)
				require.Equal(t, *tc.tool, resp.ToolCalls[0])
			} else {
				require.Empty(t, resp.ToolCalls)
			}
		})
	}
	t.Run("incomplete_stream", func(t *testing.T) {
		cut := strings.SplitAfter(hallResponsesSSE, "3456\"}\n\n")[0]
		resp, err := ProviderHallParseStream("responses", strings.NewReader(cut), start, hallClock(start, 10*time.Millisecond))
		require.NoError(t, err)
		require.False(t, resp.Terminal)
		require.Equal(t, ProviderHallErrStreamIncomplete, resp.ErrorCode)
		require.Equal(t, "123456", resp.Text)
	})
	t.Run("failed_response", func(t *testing.T) {
		body := "event: response.failed\ndata: {\"type\":\"response.failed\",\"response\":{\"model\":\"gpt-x\",\"error\":{\"code\":\"rate_limited\"}}}\n\n"
		resp, err := ProviderHallParseStream("responses", strings.NewReader(body), start, hallClock(start, 5*time.Millisecond))
		require.NoError(t, err)
		require.True(t, resp.Terminal)
		require.Equal(t, "upstream_rate_limited", resp.ErrorCode)
		require.Nil(t, resp.TTFTMs, "an error stream has no first content")
	})
}

func TestProviderHallBuildRequestBody(t *testing.T) {
	tc := ProviderHallTestCase{Prompt: "hi", MaxTokens: 77, Tool: true}
	for protocol, want := range map[string][]string{
		"responses":        {`"max_output_tokens":77`, `"store":false`, `"tool_choice":{"name":"hall_echo","type":"function"}`, `"input_text"`},
		"chat_completions": {`"max_tokens":77`, `"include_usage":true`, `"tool_choice":{"function":{"name":"hall_echo"},"type":"function"}`},
		"messages":         {`"max_tokens":77`, `"tool_choice":{"name":"hall_echo","type":"tool"}`, `"input_schema"`},
	} {
		body, err := ProviderHallBuildRequestBody(protocol, "gpt-x", tc)
		require.NoError(t, err, protocol)
		require.Contains(t, string(body), `"stream":true`)
		require.Contains(t, string(body), `"model":"gpt-x"`)
		for _, w := range want {
			require.Contains(t, string(body), w, protocol)
		}
		var obj map[string]any
		require.NoError(t, json.Unmarshal(body, &obj))
	}
	_, err := ProviderHallBuildRequestBody("ws", "gpt-x", tc)
	require.Error(t, err)
}

func TestProviderHallProbeClientSendsHeadersThroughGateway(t *testing.T) {
	t.Setenv("PROVIDER_HALL_ALLOW_LOOPBACK", "1")
	var seenTask, seenAuth, seenAccept, seenUA, seenClientID, seenPath string
	var hits atomic.Int32
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		seenPath = r.URL.Path
		seenTask = r.Header.Get(ProviderHallTaskHeader)
		seenAuth = r.Header.Get("Authorization")
		seenAccept = r.Header.Get("Accept")
		seenUA = r.Header.Get("User-Agent")
		seenClientID = r.Header.Get("X-Client-Request-ID")
		body, _ := io.ReadAll(r.Body)
		require.Contains(t, string(body), `"stream":true`)
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("X-Client-Request-ID", "gw-"+seenClientID)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, hallChatSSE)
	}))
	defer gateway.Close()
	trace := uuid.New()
	client := NewProviderHallProbeClient(5 * time.Second)
	req := ProviderHallProbeRequest{Origin: gateway.URL, Protocol: "chat_completions", Model: "gpt-x", Key: "sk-probe",
		Task: ProviderHallTaskRef{Kind: ProviderHallSourceProbe, JobID: 3, SampleID: 9, TraceID: trace}, Case: ProviderHallTestCase{Prompt: "p", MaxTokens: 5}, Version: "1.2.3"}
	res := client.Send(context.Background(), req)
	require.False(t, res.Uncertain)
	require.Empty(t, res.ErrorCode)
	require.Equal(t, http.StatusOK, res.HTTPStatus)
	require.Equal(t, "gw-"+trace.String(), res.ClientRequestID, "gateway echo is recorded for reconciliation")
	require.NotNil(t, res.Response)
	require.True(t, res.Response.Terminal)
	require.Equal(t, int32(1), hits.Load())
	require.Equal(t, "/v1/chat/completions", seenPath)
	require.Equal(t, "probe:3:9:"+trace.String(), seenTask)
	require.Equal(t, "Bearer sk-probe", seenAuth)
	require.Equal(t, "text/event-stream", seenAccept)
	require.Equal(t, "sub2api-provider-hall/1.2.3", seenUA)
	require.Equal(t, trace.String(), seenClientID)

	t.Run("non_2xx_is_received_error", func(t *testing.T) {
		bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusTooManyRequests) }))
		defer bad.Close()
		req := req
		req.Origin = bad.URL
		res := client.Send(context.Background(), req)
		require.False(t, res.Uncertain)
		require.Equal(t, "http_429", res.ErrorCode)
	})
	t.Run("stream_cut_after_headers_is_uncertain", func(t *testing.T) {
		cut := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, "data: {\"id\":\"c\",\"model\":\"gpt-x\",\"choices\":[{\"delta\":{\"content\":\"x\"}}]}\n\n")
		}))
		defer cut.Close()
		req := req
		req.Origin = cut.URL
		res := client.Send(context.Background(), req)
		require.True(t, res.Uncertain, "the gateway admitted the request; never resend")
		require.Equal(t, ProviderHallErrStreamIncomplete, res.ErrorCode)
	})
	t.Run("connection_refused_is_not_uncertain", func(t *testing.T) {
		closed := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		origin := closed.URL
		closed.Close()
		req := req
		req.Origin = origin
		res := client.Send(context.Background(), req)
		require.False(t, res.Uncertain)
		require.Equal(t, ProviderHallErrTransportRefused, res.ErrorCode)
	})
}

func TestProviderHallProbeClientRefusesRedirectsAndLoopback(t *testing.T) {
	t.Setenv("PROVIDER_HALL_ALLOW_LOOPBACK", "1")
	var upstreamHits atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { upstreamHits.Add(1) }))
	defer target.Close()
	redirecting := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+r.URL.Path, http.StatusTemporaryRedirect)
	}))
	defer redirecting.Close()
	client := NewProviderHallProbeClient(5 * time.Second)
	req := ProviderHallProbeRequest{Origin: redirecting.URL, Protocol: "responses", Model: "gpt-x", Key: "k",
		Task: ProviderHallTaskRef{Kind: ProviderHallSourceProbe, JobID: 1, SampleID: 1, TraceID: uuid.New()}, Case: ProviderHallTestCase{Prompt: "p"}}
	res := client.Send(context.Background(), req)
	require.Equal(t, ProviderHallErrRedirect, res.ErrorCode)
	require.False(t, res.Uncertain)
	require.Equal(t, int32(0), upstreamHits.Load(), "the redirect target never receives the key")

	t.Setenv("PROVIDER_HALL_ALLOW_LOOPBACK", "")
	_, err := client.NewRequest(context.Background(), req)
	require.Error(t, err, "loopback origins need PROVIDER_HALL_ALLOW_LOOPBACK=1")
	req.Origin = "http://gateway.example.com"
	_, err = client.NewRequest(context.Background(), req)
	require.Error(t, err, "non-loopback origins must be https")
	req.Origin = "https://gateway.example.com"
	httpReq, err := client.NewRequest(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "https://gateway.example.com/v1/responses", httpReq.URL.String())
	require.True(t, ProviderHallIsLoopbackOrigin("http://localhost:8080"))
	require.True(t, ProviderHallIsLoopbackOrigin("https://[::1]"))
	require.False(t, ProviderHallIsLoopbackOrigin("https://10.0.0.1"))
}
