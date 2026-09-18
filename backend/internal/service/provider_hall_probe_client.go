package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ProviderHallProbeTimeout bounds one probe request end to end.
const ProviderHallProbeTimeout = 90 * time.Second

// Sample error codes.
const (
	ProviderHallErrTransportRefused = "transport_refused"
	ProviderHallErrRedirect         = "redirect_refused"
	ProviderHallErrTimeout          = "timeout"
	ProviderHallErrStreamIncomplete = "stream_incomplete"
	ProviderHallErrUncertainTimeout = "uncertain_timeout"
	ProviderHallErrResponseLost     = "response_lost"
	ProviderHallErrUpstream         = "upstream_error"
)

// ProviderHallAllowLoopback reports whether probes may target a loopback
// gateway origin (local end-to-end runs only).
func ProviderHallAllowLoopback() bool {
	return os.Getenv("PROVIDER_HALL_ALLOW_LOOPBACK") == "1"
}

// ProviderHallIsLoopbackOrigin reports whether the origin's host is loopback.
func ProviderHallIsLoopbackOrigin(origin string) bool {
	u, err := url.Parse(strings.TrimSpace(origin))
	if err != nil {
		return false
	}
	host := strings.ToLower(u.Hostname())
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback() || ip.IsUnspecified()
	}
	return false
}

// ProviderHallProtocolPath maps a protocol to the gateway path.
func ProviderHallProtocolPath(protocol string) (string, bool) {
	switch protocol {
	case "responses":
		return "/v1/responses", true
	case "chat_completions":
		return "/v1/chat/completions", true
	case "messages":
		return "/v1/messages", true
	}
	return "", false
}

// ProviderHallBuildRequestBody renders the streaming request for one case.
func ProviderHallBuildRequestBody(protocol, model string, tc ProviderHallTestCase) ([]byte, error) {
	maxTokens := tc.MaxTokens
	if maxTokens <= 0 {
		maxTokens = ProviderHallProbeMaxTokens
	}
	var body map[string]any
	switch protocol {
	case "responses":
		body = map[string]any{
			"model":             model,
			"input":             []any{map[string]any{"role": "user", "content": []any{map[string]any{"type": "input_text", "text": tc.Prompt}}}},
			"stream":            true,
			"store":             false,
			"max_output_tokens": maxTokens,
		}
		if tc.Tool {
			body["tools"] = []any{map[string]any{"type": "function", "name": ProviderHallToolName, "parameters": ProviderHallToolSchema}}
			body["tool_choice"] = map[string]any{"type": "function", "name": ProviderHallToolName}
		}
	case "chat_completions":
		body = map[string]any{
			"model":          model,
			"messages":       []any{map[string]any{"role": "user", "content": tc.Prompt}},
			"stream":         true,
			"max_tokens":     maxTokens,
			"stream_options": map[string]any{"include_usage": true},
		}
		if tc.Tool {
			body["tools"] = []any{map[string]any{"type": "function", "function": map[string]any{"name": ProviderHallToolName, "parameters": ProviderHallToolSchema}}}
			body["tool_choice"] = map[string]any{"type": "function", "function": map[string]any{"name": ProviderHallToolName}}
		}
	case "messages":
		body = map[string]any{
			"model":      model,
			"max_tokens": maxTokens,
			"stream":     true,
			"messages":   []any{map[string]any{"role": "user", "content": tc.Prompt}},
		}
		if tc.Tool {
			body["tools"] = []any{map[string]any{"name": ProviderHallToolName, "input_schema": ProviderHallToolSchema}}
			body["tool_choice"] = map[string]any{"type": "tool", "name": ProviderHallToolName}
		}
	default:
		return nil, fmt.Errorf("provider hall: unsupported protocol %q", protocol)
	}
	return json.Marshal(body)
}

type providerHallSSEEvent struct {
	event string
	data  string
}

// providerHallReadSSE yields events; a "data: [DONE]" is passed through.
func providerHallReadSSE(r io.Reader, fn func(ev providerHallSSEEvent) bool) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	var cur providerHallSSEEvent
	var data []string
	flush := func() bool {
		if len(data) == 0 && cur.event == "" {
			return true
		}
		cur.data = strings.Join(data, "\n")
		keep := fn(cur)
		cur, data = providerHallSSEEvent{}, nil
		return keep
	}
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == "":
			if !flush() {
				return nil
			}
		case strings.HasPrefix(line, ":"):
		case strings.HasPrefix(line, "event:"):
			cur.event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		case strings.HasPrefix(line, "data:"):
			data = append(data, strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	flush()
	return nil
}

// ProviderHallParseStream consumes one streamed response body. start is the
// moment the request was sent; TTFT is the first content delta.
func ProviderHallParseStream(protocol string, body io.Reader, start time.Time, now func() time.Time) (*ProviderHallProbeResponse, error) {
	if now == nil {
		now = time.Now
	}
	resp := &ProviderHallProbeResponse{}
	firstContent := func() {
		if resp.TTFTMs == nil {
			ms := int(now().Sub(start) / time.Millisecond)
			if ms < 1 {
				ms = 1
			}
			resp.TTFTMs = &ms
		}
	}
	var text strings.Builder
	toolArgs := map[int]*ProviderHallToolCall{}
	toolOrder := []int{}
	var parseErr error
	handle := func(ev providerHallSSEEvent) bool {
		if strings.TrimSpace(ev.data) == "[DONE]" {
			if protocol == "chat_completions" {
				resp.Terminal = true
			}
			return false
		}
		var payload map[string]json.RawMessage
		if err := json.Unmarshal([]byte(ev.data), &payload); err != nil {
			return true
		}
		eventType := ev.event
		if eventType == "" {
			eventType = providerHallJSONString(payload["type"])
		}
		switch protocol {
		case "responses":
			switch eventType {
			case "response.output_text.delta":
				if delta := providerHallJSONString(payload["delta"]); delta != "" {
					firstContent()
					text.WriteString(delta)
				}
			case "response.output_item.added":
				var item struct {
					Type string `json:"type"`
					Name string `json:"name"`
				}
				_ = json.Unmarshal(payload["item"], &item)
				if item.Type == "function_call" {
					firstContent()
				}
			case "response.completed", "response.incomplete", "response.failed":
				var response struct {
					Model  string `json:"model"`
					Status string `json:"status"`
					Output []struct {
						Type      string `json:"type"`
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"output"`
					Usage struct {
						InputTokens  *int `json:"input_tokens"`
						OutputTokens *int `json:"output_tokens"`
					} `json:"usage"`
					Error *struct {
						Code string `json:"code"`
					} `json:"error"`
				}
				_ = json.Unmarshal(payload["response"], &response)
				resp.Model = response.Model
				resp.InputTokens, resp.OutputTokens = response.Usage.InputTokens, response.Usage.OutputTokens
				for _, out := range response.Output {
					if out.Type == "function_call" {
						resp.ToolCalls = append(resp.ToolCalls, ProviderHallToolCall{Name: out.Name, Arguments: out.Arguments})
					}
				}
				resp.Terminal = true
				if eventType != "response.completed" {
					resp.ErrorCode = ProviderHallErrUpstream
					if response.Error != nil && response.Error.Code != "" {
						resp.ErrorCode = "upstream_" + response.Error.Code
					}
				}
				return false
			case "error":
				resp.ErrorCode = ProviderHallErrUpstream
				resp.Terminal = true
				return false
			}
		case "chat_completions":
			if model := providerHallJSONString(payload["model"]); model != "" {
				resp.Model = model
			}
			if payload["error"] != nil {
				resp.ErrorCode = ProviderHallErrUpstream
				resp.Terminal = true
				return false
			}
			var chunk struct {
				Choices []struct {
					Delta struct {
						Content   string `json:"content"`
						ToolCalls []struct {
							Index    int `json:"index"`
							Function struct {
								Name      string `json:"name"`
								Arguments string `json:"arguments"`
							} `json:"function"`
						} `json:"tool_calls"`
					} `json:"delta"`
				} `json:"choices"`
				Usage *struct {
					PromptTokens     *int `json:"prompt_tokens"`
					CompletionTokens *int `json:"completion_tokens"`
				} `json:"usage"`
			}
			_ = json.Unmarshal([]byte(ev.data), &chunk)
			for _, choice := range chunk.Choices {
				if choice.Delta.Content != "" {
					firstContent()
					text.WriteString(choice.Delta.Content)
				}
				for _, call := range choice.Delta.ToolCalls {
					firstContent()
					entry := toolArgs[call.Index]
					if entry == nil {
						entry = &ProviderHallToolCall{}
						toolArgs[call.Index] = entry
						toolOrder = append(toolOrder, call.Index)
					}
					if call.Function.Name != "" {
						entry.Name = call.Function.Name
					}
					entry.Arguments += call.Function.Arguments
				}
			}
			if chunk.Usage != nil {
				resp.InputTokens, resp.OutputTokens = chunk.Usage.PromptTokens, chunk.Usage.CompletionTokens
			}
		case "messages":
			switch eventType {
			case "message_start":
				var msg struct {
					Message struct {
						Model string `json:"model"`
						Usage struct {
							InputTokens *int `json:"input_tokens"`
						} `json:"usage"`
					} `json:"message"`
				}
				_ = json.Unmarshal([]byte(ev.data), &msg)
				resp.Model = msg.Message.Model
				resp.InputTokens = msg.Message.Usage.InputTokens
			case "content_block_start":
				var block struct {
					Index        int `json:"index"`
					ContentBlock struct {
						Type string `json:"type"`
						Name string `json:"name"`
					} `json:"content_block"`
				}
				_ = json.Unmarshal([]byte(ev.data), &block)
				if block.ContentBlock.Type == "tool_use" {
					firstContent()
					entry := &ProviderHallToolCall{Name: block.ContentBlock.Name}
					toolArgs[block.Index] = entry
					toolOrder = append(toolOrder, block.Index)
				}
			case "content_block_delta":
				var delta struct {
					Index int `json:"index"`
					Delta struct {
						Type        string `json:"type"`
						Text        string `json:"text"`
						PartialJSON string `json:"partial_json"`
					} `json:"delta"`
				}
				_ = json.Unmarshal([]byte(ev.data), &delta)
				switch delta.Delta.Type {
				case "text_delta":
					if delta.Delta.Text != "" {
						firstContent()
						text.WriteString(delta.Delta.Text)
					}
				case "input_json_delta":
					firstContent()
					if entry := toolArgs[delta.Index]; entry != nil {
						entry.Arguments += delta.Delta.PartialJSON
					}
				}
			case "message_delta":
				var d struct {
					Usage struct {
						OutputTokens *int `json:"output_tokens"`
					} `json:"usage"`
				}
				_ = json.Unmarshal([]byte(ev.data), &d)
				if d.Usage.OutputTokens != nil {
					resp.OutputTokens = d.Usage.OutputTokens
				}
			case "message_stop":
				resp.Terminal = true
				return false
			case "error":
				resp.ErrorCode = ProviderHallErrUpstream
				resp.Terminal = true
				return false
			}
		}
		return true
	}
	parseErr = providerHallReadSSE(body, handle)
	for _, idx := range toolOrder {
		if entry := toolArgs[idx]; entry != nil {
			if entry.Arguments == "" {
				entry.Arguments = "{}"
			}
			resp.ToolCalls = append(resp.ToolCalls, *entry)
		}
	}
	resp.Text = text.String()
	total := int(now().Sub(start) / time.Millisecond)
	if total < 1 {
		total = 1
	}
	resp.TotalMs = &total
	if parseErr != nil {
		return resp, parseErr
	}
	if !resp.Terminal && resp.ErrorCode == "" {
		resp.ErrorCode = ProviderHallErrStreamIncomplete
	}
	return resp, nil
}

func providerHallJSONString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return ""
	}
	return s
}

// ProviderHallGenerationMs derives generation time from total and TTFT.
func ProviderHallGenerationMs(totalMs, ttftMs *int) *int {
	if totalMs == nil || ttftMs == nil || *totalMs < *ttftMs {
		return nil
	}
	g := *totalMs - *ttftMs
	return &g
}

// ProviderHallTPS computes output tokens per second of generation.
func ProviderHallTPS(outputTokens, generationMs *int) *float64 {
	if outputTokens == nil || generationMs == nil || *generationMs <= 0 || *outputTokens < 0 {
		return nil
	}
	v := float64(*outputTokens) / float64(*generationMs) * 1000
	return &v
}

// ProviderHallProbeRequest is one outbound sample.
type ProviderHallProbeRequest struct {
	Origin   string
	Protocol string
	Model    string
	Key      string
	Task     ProviderHallTaskRef
	Case     ProviderHallTestCase
	Version  string
}

// ProviderHallProbeResult classifies the transport outcome.
type ProviderHallProbeResult struct {
	// Uncertain means the request may have been admitted by the gateway; the
	// runner must never resend it.
	Uncertain       bool
	ErrorCode       string
	HTTPStatus      int
	ClientRequestID string
	Response        *ProviderHallProbeResponse
}

// ProviderHallProbeClient sends samples through the local gateway.
type ProviderHallProbeClient struct {
	client *http.Client
	now    func() time.Time
}

func NewProviderHallProbeClient(timeout time.Duration) *ProviderHallProbeClient {
	if timeout <= 0 {
		timeout = ProviderHallProbeTimeout
	}
	return &ProviderHallProbeClient{
		client: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return errors.New("provider hall: redirects are not followed")
			},
		},
		now: time.Now,
	}
}

// NewRequest builds the HTTP request; exported so tests can assert headers.
func (c *ProviderHallProbeClient) NewRequest(ctx context.Context, in ProviderHallProbeRequest) (*http.Request, error) {
	path, ok := ProviderHallProtocolPath(in.Protocol)
	if !ok {
		return nil, fmt.Errorf("provider hall: unsupported protocol %q", in.Protocol)
	}
	origin := strings.TrimRight(strings.TrimSpace(in.Origin), "/")
	u, err := url.Parse(origin)
	if err != nil || origin == "" || u.Host == "" || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return nil, fmt.Errorf("provider hall: invalid gateway origin")
	}
	if ProviderHallIsLoopbackOrigin(origin) {
		if !ProviderHallAllowLoopback() {
			return nil, fmt.Errorf("provider hall: loopback gateway origin requires PROVIDER_HALL_ALLOW_LOOPBACK=1")
		}
	} else if u.Scheme != "https" {
		return nil, fmt.Errorf("provider hall: gateway origin must use https")
	}
	body, err := ProviderHallBuildRequestBody(in.Protocol, in.Model, in.Case)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, origin+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+in.Key)
	req.Header.Set(ProviderHallTaskHeader, FormatProviderHallTaskHeader(in.Task))
	req.Header.Set("X-Client-Request-ID", in.Task.TraceID.String())
	version := strings.TrimSpace(in.Version)
	if version == "" {
		version = "dev"
	}
	req.Header.Set("User-Agent", "sub2api-provider-hall/"+version)
	if in.Protocol == "messages" {
		req.Header.Set("x-api-key", in.Key)
		req.Header.Set("anthropic-version", "2023-06-01")
	}
	return req, nil
}

// Send performs the request and classifies the result.
func (c *ProviderHallProbeClient) Send(ctx context.Context, in ProviderHallProbeRequest) ProviderHallProbeResult {
	req, err := c.NewRequest(ctx, in)
	if err != nil {
		return ProviderHallProbeResult{ErrorCode: "request_invalid"}
	}
	start := c.now()
	res, err := c.client.Do(req)
	if err != nil {
		code := providerHallTransportCode(err)
		// A redirect is answered by something that is not the gateway's model
		// endpoint; a refused dial never reached it. Everything else may have
		// been admitted and must wait for reconciliation.
		return ProviderHallProbeResult{Uncertain: code != ProviderHallErrRedirect && !providerHallDefinitelyNotAdmitted(err), ErrorCode: code}
	}
	defer func() { _ = res.Body.Close() }()
	out := ProviderHallProbeResult{HTTPStatus: res.StatusCode, ClientRequestID: strings.TrimSpace(res.Header.Get("X-Client-Request-ID"))}
	if out.ClientRequestID == "" {
		out.ClientRequestID = in.Task.TraceID.String()
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 64*1024))
		out.ErrorCode = fmt.Sprintf("http_%d", res.StatusCode)
		return out
	}
	parsed, parseErr := ProviderHallParseStream(in.Protocol, res.Body, start, c.now)
	out.Response = parsed
	if parseErr != nil || (parsed != nil && !parsed.Terminal) {
		// Headers arrived: the gateway admitted the request, so the bill may
		// still land. Wait for reconciliation instead of resending.
		out.Uncertain = true
		out.ErrorCode = ProviderHallErrStreamIncomplete
		if parseErr != nil && errors.Is(parseErr, context.DeadlineExceeded) {
			out.ErrorCode = ProviderHallErrTimeout
		}
		return out
	}
	if parsed.ErrorCode != "" {
		out.ErrorCode = parsed.ErrorCode
	}
	return out
}

func providerHallTransportCode(err error) string {
	switch {
	case err == nil:
		return ""
	case strings.Contains(err.Error(), "redirects are not followed"):
		return ProviderHallErrRedirect
	case errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err):
		return ProviderHallErrTimeout
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return ProviderHallErrTimeout
	}
	if providerHallDefinitelyNotAdmitted(err) {
		return ProviderHallErrTransportRefused
	}
	return "transport_error"
}

// providerHallDefinitelyNotAdmitted is true for failures that happen before
// any byte reached the gateway (DNS, connection refused).
func providerHallDefinitelyNotAdmitted(err error) bool {
	if err == nil {
		return false
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) && opErr.Op == "dial" {
		return true
	}
	if strings.Contains(err.Error(), "redirects are not followed") {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "connection refused") || strings.Contains(msg, "no such host")
}

var _ = uuid.Nil
