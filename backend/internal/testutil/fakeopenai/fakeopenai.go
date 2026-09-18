// Package fakeopenai is a shared fake OpenAI-compatible upstream for Go tests
// and the provider hall end-to-end script. It serves the Responses, Chat
// Completions and Anthropic Messages endpoints, streaming (SSE) or not, and
// answers the provider hall verification prompts correctly by default so a
// verification job passes unless a knob says otherwise.
//
// The package has no build tag: unit and integration tests alike may import
// it. It depends only on the standard library.
package fakeopenai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// FaultHeader carries per-request fault injection:
// "delay=<ms>;status=<code>;truncate=1;first_token_delay=<ms>". Fields are
// optional and override the server Options for that request only.
const FaultHeader = "X-Fake-Upstream"

// ToolName is the forced tool every protocol answers with when the request
// declares it.
const ToolName = "hall_echo"

// Options configure default behaviour. The zero value serves correct answers
// with no delay.
type Options struct {
	// ModelAlias maps a requested model to the model name echoed back
	// (e.g. {"gpt-5": "gpt-5-2026-01-01"}). Unmapped models are echoed as-is.
	ModelAlias map[string]string
	// MismatchModel, when non-empty, is echoed for every request regardless
	// of what was asked (model verification "suspected" path).
	MismatchModel string

	// Usage reported in the terminal event. Zero values fall back to
	// InputTokens=12, OutputTokens=len(answer tokens), CacheReadTokens=0.
	InputTokens     int
	OutputTokens    int
	CacheReadTokens int

	// WrongArithmetic answers "Compute a op b." with the wrong integer.
	WrongArithmetic bool
	// WrongJSON answers the token/length prompt with an extra key.
	WrongJSON bool
	// WrongTool answers the forced tool call with a wrong code.
	WrongTool bool

	// Delay is applied before the response headers are sent.
	Delay time.Duration
	// FirstTokenDelay is applied after headers, before the first content delta.
	FirstTokenDelay time.Duration
	// ChunkDelay is applied between successive stream events.
	ChunkDelay time.Duration
	// Status429Every makes every Nth request (1-based) answer 429.
	Status429Every int
	// TruncateStream ends every stream after the first content delta without
	// a terminal event (simulates an upstream that died mid-stream).
	TruncateStream bool
	// RequireAuth rejects requests without Authorization/x-api-key with 401.
	RequireAuth bool
}

// Request is one recorded inbound request.
type Request struct {
	Method     string
	Path       string
	Header     http.Header
	Body       map[string]any
	Raw        []byte
	Protocol   string // responses | chat_completions | messages | ""
	ReceivedAt time.Time
}

// Server is the fake upstream. Its URL is Server.URL (from httptest.Server).
type Server struct {
	*httptest.Server
	mu       sync.Mutex
	opts     Options
	requests []Request
	counter  atomic.Int64
}

// New starts a fake upstream with the given options. Close it with Close.
func New(opts Options) *Server {
	s := &Server{opts: opts}
	s.Server = httptest.NewServer(http.HandlerFunc(s.handle))
	return s
}

// NewTLS starts the fake upstream over TLS (self-signed; use Server.Client()).
func NewTLS(opts Options) *Server {
	s := &Server{opts: opts}
	s.Server = httptest.NewTLSServer(http.HandlerFunc(s.handle))
	return s
}

// Handler returns the raw handler for callers that manage their own listener.
func (s *Server) Handler() http.Handler { return http.HandlerFunc(s.handle) }

// Options returns a copy of the current defaults.
func (s *Server) Options() Options {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.opts
}

// SetOptions replaces the defaults for subsequent requests.
func (s *Server) SetOptions(opts Options) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.opts = opts
}

// Requests returns a copy of every recorded request in arrival order.
func (s *Server) Requests() []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Request(nil), s.requests...)
}

// Count returns how many requests have been recorded since the last Reset.
func (s *Server) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.requests)
}

// Reset clears recorded requests and the 429 rotation counter.
func (s *Server) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = nil
	s.counter.Store(0)
}

// ---------------------------------------------------------------------------

type fault struct {
	delay           time.Duration
	firstTokenDelay time.Duration
	chunkDelay      time.Duration
	status          int
	truncate        bool
}

func parseFault(base Options, header string) fault {
	f := fault{delay: base.Delay, firstTokenDelay: base.FirstTokenDelay, chunkDelay: base.ChunkDelay, truncate: base.TruncateStream}
	for _, part := range strings.Split(header, ";") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		key, value := strings.ToLower(strings.TrimSpace(kv[0])), strings.TrimSpace(kv[1])
		switch key {
		case "delay":
			if ms, err := strconv.Atoi(value); err == nil && ms >= 0 {
				f.delay = time.Duration(ms) * time.Millisecond
			}
		case "first_token_delay", "ttft":
			if ms, err := strconv.Atoi(value); err == nil && ms >= 0 {
				f.firstTokenDelay = time.Duration(ms) * time.Millisecond
			}
		case "chunk_delay":
			if ms, err := strconv.Atoi(value); err == nil && ms >= 0 {
				f.chunkDelay = time.Duration(ms) * time.Millisecond
			}
		case "status":
			if code, err := strconv.Atoi(value); err == nil && code >= 100 && code <= 599 {
				f.status = code
			}
		case "truncate":
			f.truncate = value == "1" || strings.EqualFold(value, "true")
		}
	}
	return f
}

func protocolForPath(path string) string {
	p := strings.TrimRight(strings.ToLower(path), "/")
	switch {
	case strings.HasSuffix(p, "/responses"):
		return "responses"
	case strings.HasSuffix(p, "/chat/completions"):
		return "chat_completions"
	case strings.HasSuffix(p, "/messages"):
		return "messages"
	}
	return ""
}

var (
	arithPattern = regexp.MustCompile(`Compute\s+(-?\d+)\s*([+\-*xX×])\s*(-?\d+)`)
	jsonPattern  = regexp.MustCompile(`"token"\s+set to the string\s+"([^"]+)"\s+and\s+"length"\s+set to the integer\s+(\d+)`)
	toolPattern  = regexp.MustCompile(`with the code\s+"([^"]+)"`)
)

// answer is the protocol-neutral reply for one request.
type answer struct {
	text     string // empty when toolArgs is set
	toolArgs string // JSON arguments for hall_echo, empty when text answer
	model    string
	stream   bool
	input    int
	output   int
	cached   int
}

func (s *Server) buildAnswer(opts Options, protocol string, body map[string]any) answer {
	a := answer{}
	model, _ := body["model"].(string)
	a.model = model
	if alias, ok := opts.ModelAlias[model]; ok && alias != "" {
		a.model = alias
	}
	if opts.MismatchModel != "" {
		a.model = opts.MismatchModel
	}
	a.stream, _ = body["stream"].(bool)
	prompt := extractPrompt(protocol, body)
	forcedTool := wantsTool(protocol, body)

	switch {
	case forcedTool:
		code := "unknown"
		if m := toolPattern.FindStringSubmatch(prompt); m != nil {
			code = m[1]
		}
		if opts.WrongTool {
			code = "wrong-" + code
		}
		args, _ := json.Marshal(map[string]string{"code": code})
		a.toolArgs = string(args)
	case arithPattern.MatchString(prompt):
		m := arithPattern.FindStringSubmatch(prompt)
		x, _ := strconv.ParseInt(m[1], 10, 64)
		y, _ := strconv.ParseInt(m[3], 10, 64)
		var v int64
		switch m[2] {
		case "+":
			v = x + y
		case "-":
			v = x - y
		default:
			v = x * y
		}
		if opts.WrongArithmetic {
			v++
		}
		a.text = strconv.FormatInt(v, 10)
	case jsonPattern.MatchString(prompt):
		m := jsonPattern.FindStringSubmatch(prompt)
		n, _ := strconv.Atoi(m[2])
		if opts.WrongJSON {
			a.text = fmt.Sprintf(`{"token":%q,"length":%d,"extra":true}`, m[1], n)
		} else {
			a.text = fmt.Sprintf(`{"token":%q,"length":%d}`, m[1], n)
		}
	default:
		a.text = "OK"
	}
	a.input = opts.InputTokens
	if a.input <= 0 {
		a.input = 12
	}
	a.output = opts.OutputTokens
	if a.output <= 0 {
		a.output = len(strings.Fields(a.text)) + len(a.toolArgs)/4
		if a.output < 1 {
			a.output = 1
		}
	}
	a.cached = opts.CacheReadTokens
	if a.cached < 0 {
		a.cached = 0
	}
	return a
}

// extractPrompt returns the last user text in the request.
func extractPrompt(protocol string, body map[string]any) string {
	switch protocol {
	case "responses":
		switch input := body["input"].(type) {
		case string:
			return input
		case []any:
			var last string
			for _, item := range input {
				m, ok := item.(map[string]any)
				if !ok {
					continue
				}
				if role, _ := m["role"].(string); role != "" && role != "user" {
					continue
				}
				switch content := m["content"].(type) {
				case string:
					last = content
				case []any:
					for _, part := range content {
						pm, ok := part.(map[string]any)
						if !ok {
							continue
						}
						if text, _ := pm["text"].(string); text != "" {
							last = text
						}
					}
				}
			}
			return last
		}
	case "chat_completions", "messages":
		msgs, _ := body["messages"].([]any)
		var last string
		for _, item := range msgs {
			m, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if role, _ := m["role"].(string); role != "user" {
				continue
			}
			switch content := m["content"].(type) {
			case string:
				last = content
			case []any:
				for _, part := range content {
					pm, ok := part.(map[string]any)
					if !ok {
						continue
					}
					if text, _ := pm["text"].(string); text != "" {
						last = text
					}
				}
			}
		}
		return last
	}
	return ""
}

// wantsTool reports whether the request declares hall_echo and forces it.
func wantsTool(protocol string, body map[string]any) bool {
	tools, _ := body["tools"].([]any)
	declared := false
	for _, t := range tools {
		m, ok := t.(map[string]any)
		if !ok {
			continue
		}
		name, _ := m["name"].(string)
		if fn, ok := m["function"].(map[string]any); ok && name == "" {
			name, _ = fn["name"].(string)
		}
		if name == ToolName {
			declared = true
		}
	}
	if !declared {
		return false
	}
	switch choice := body["tool_choice"].(type) {
	case string:
		return choice == "required" || choice == "any" || choice == ToolName
	case map[string]any:
		if name, _ := choice["name"].(string); name == ToolName {
			return true
		}
		if fn, ok := choice["function"].(map[string]any); ok {
			name, _ := fn["name"].(string)
			return name == ToolName
		}
		typ, _ := choice["type"].(string)
		return typ == "required" || typ == "any"
	}
	// Declared but not forced: still call it, the verification prompt asks for it.
	return protocol != "" && strings.Contains(extractPrompt(protocol, body), ToolName)
}

// ---------------------------------------------------------------------------

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(io.LimitReader(r.Body, 8<<20))
	var body map[string]any
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &body)
	}
	protocol := protocolForPath(r.URL.Path)
	rec := Request{Method: r.Method, Path: r.URL.Path, Header: r.Header.Clone(), Body: body, Raw: raw, Protocol: protocol, ReceivedAt: time.Now()}
	s.mu.Lock()
	s.requests = append(s.requests, rec)
	opts := s.opts
	s.mu.Unlock()
	n := s.counter.Add(1)

	if opts.RequireAuth && r.Header.Get("Authorization") == "" && r.Header.Get("x-api-key") == "" {
		writeError(w, protocol, http.StatusUnauthorized, "missing_api_key", "no credentials")
		return
	}
	if protocol == "" || r.Method != http.MethodPost {
		writeError(w, "", http.StatusNotFound, "not_found", "unknown endpoint "+r.URL.Path)
		return
	}
	f := parseFault(opts, r.Header.Get(FaultHeader))
	if f.delay > 0 {
		time.Sleep(f.delay)
	}
	if f.status == 0 && opts.Status429Every > 0 && n%int64(opts.Status429Every) == 0 {
		f.status = http.StatusTooManyRequests
	}
	if f.status != 0 && (f.status < 200 || f.status >= 300) {
		if f.status == http.StatusTooManyRequests {
			w.Header().Set("Retry-After", "1")
		}
		writeError(w, protocol, f.status, "fake_upstream_fault", fmt.Sprintf("injected status %d", f.status))
		return
	}
	a := s.buildAnswer(opts, protocol, body)
	id := fmt.Sprintf("fake-%s-%d", protocol, n)
	w.Header().Set("X-Request-Id", id)
	if !a.stream {
		writeJSON(w, http.StatusOK, nonStreamBody(protocol, id, a))
		return
	}
	st := &streamWriter{w: w, flusher: flusherOf(w), chunkDelay: f.chunkDelay}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	st.flush()
	switch protocol {
	case "responses":
		streamResponses(st, id, a, f)
	case "chat_completions":
		streamChat(st, id, a, f, body)
	case "messages":
		streamMessages(st, id, a, f)
	}
}

func flusherOf(w http.ResponseWriter) http.Flusher {
	if f, ok := w.(http.Flusher); ok {
		return f
	}
	return nil
}

type streamWriter struct {
	w          http.ResponseWriter
	flusher    http.Flusher
	chunkDelay time.Duration
	events     int
}

func (st *streamWriter) flush() {
	if st.flusher != nil {
		st.flusher.Flush()
	}
}

// event writes one SSE event. Responses and Messages include the event name;
// Chat Completions uses bare data lines.
func (st *streamWriter) event(name string, payload any) {
	if st.events > 0 && st.chunkDelay > 0 {
		time.Sleep(st.chunkDelay)
	}
	st.events++
	data, _ := json.Marshal(payload)
	if name != "" {
		_, _ = fmt.Fprintf(st.w, "event: %s\n", name)
	}
	_, _ = fmt.Fprintf(st.w, "data: %s\n\n", data)
	st.flush()
}

func (st *streamWriter) raw(line string) {
	st.events++
	_, _ = fmt.Fprintf(st.w, "data: %s\n\n", line)
	st.flush()
}

// splitText yields text deltas of a few characters so TTFT and generation are
// distinguishable.
func splitText(text string) []string {
	if text == "" {
		return nil
	}
	const width = 6
	var out []string
	for len(text) > width {
		out = append(out, text[:width])
		text = text[width:]
	}
	return append(out, text)
}

func responsesUsage(a answer) map[string]any {
	return map[string]any{
		"input_tokens":          a.input,
		"input_tokens_details":  map[string]any{"cached_tokens": a.cached},
		"output_tokens":         a.output,
		"output_tokens_details": map[string]any{"reasoning_tokens": 0},
		"total_tokens":          a.input + a.output,
	}
}

func chatUsage(a answer) map[string]any {
	return map[string]any{
		"prompt_tokens":         a.input,
		"completion_tokens":     a.output,
		"total_tokens":          a.input + a.output,
		"prompt_tokens_details": map[string]any{"cached_tokens": a.cached},
	}
}

func messagesUsage(a answer, output int) map[string]any {
	return map[string]any{
		"input_tokens":                a.input,
		"cache_read_input_tokens":     a.cached,
		"cache_creation_input_tokens": 0,
		"output_tokens":               output,
	}
}

func responsesOutput(a answer, callID string) []any {
	if a.toolArgs != "" {
		return []any{map[string]any{"type": "function_call", "id": "fc_" + callID, "call_id": "call_" + callID, "name": ToolName, "arguments": a.toolArgs, "status": "completed"}}
	}
	return []any{map[string]any{"type": "message", "id": "msg_" + callID, "role": "assistant", "status": "completed",
		"content": []any{map[string]any{"type": "output_text", "text": a.text, "annotations": []any{}}}}}
}

func streamResponses(st *streamWriter, id string, a answer, f fault) {
	base := map[string]any{"id": "resp_" + id, "object": "response", "model": a.model, "status": "in_progress", "output": []any{}}
	st.event("response.created", map[string]any{"type": "response.created", "sequence_number": 0, "response": base})
	st.event("response.in_progress", map[string]any{"type": "response.in_progress", "sequence_number": 1, "response": base})
	if f.firstTokenDelay > 0 {
		time.Sleep(f.firstTokenDelay)
	}
	if a.toolArgs != "" {
		item := map[string]any{"type": "function_call", "id": "fc_" + id, "call_id": "call_" + id, "name": ToolName, "arguments": "", "status": "in_progress"}
		st.event("response.output_item.added", map[string]any{"type": "response.output_item.added", "output_index": 0, "item": item})
		if f.truncate {
			return
		}
		for _, delta := range splitText(a.toolArgs) {
			st.event("response.function_call_arguments.delta", map[string]any{"type": "response.function_call_arguments.delta", "item_id": "fc_" + id, "output_index": 0, "delta": delta})
		}
		st.event("response.function_call_arguments.done", map[string]any{"type": "response.function_call_arguments.done", "item_id": "fc_" + id, "output_index": 0, "arguments": a.toolArgs})
		done := responsesOutput(a, id)[0]
		st.event("response.output_item.done", map[string]any{"type": "response.output_item.done", "output_index": 0, "item": done})
	} else {
		item := map[string]any{"type": "message", "id": "msg_" + id, "role": "assistant", "status": "in_progress", "content": []any{}}
		st.event("response.output_item.added", map[string]any{"type": "response.output_item.added", "output_index": 0, "item": item})
		st.event("response.content_part.added", map[string]any{"type": "response.content_part.added", "item_id": "msg_" + id, "output_index": 0, "content_index": 0, "part": map[string]any{"type": "output_text", "text": "", "annotations": []any{}}})
		for i, delta := range splitText(a.text) {
			st.event("response.output_text.delta", map[string]any{"type": "response.output_text.delta", "item_id": "msg_" + id, "output_index": 0, "content_index": 0, "delta": delta})
			if f.truncate && i == 0 {
				return
			}
		}
		st.event("response.output_text.done", map[string]any{"type": "response.output_text.done", "item_id": "msg_" + id, "output_index": 0, "content_index": 0, "text": a.text})
		st.event("response.content_part.done", map[string]any{"type": "response.content_part.done", "item_id": "msg_" + id, "output_index": 0, "content_index": 0, "part": map[string]any{"type": "output_text", "text": a.text, "annotations": []any{}}})
		done := responsesOutput(a, id)[0]
		st.event("response.output_item.done", map[string]any{"type": "response.output_item.done", "output_index": 0, "item": done})
	}
	final := map[string]any{"id": "resp_" + id, "object": "response", "model": a.model, "status": "completed", "output": responsesOutput(a, id), "usage": responsesUsage(a)}
	st.event("response.completed", map[string]any{"type": "response.completed", "response": final})
}

func chatChunk(id, model string, delta map[string]any, finish any) map[string]any {
	return map[string]any{
		"id": "chatcmpl-" + id, "object": "chat.completion.chunk", "created": time.Now().Unix(), "model": model,
		"choices": []any{map[string]any{"index": 0, "delta": delta, "finish_reason": finish}},
	}
}

func streamChat(st *streamWriter, id string, a answer, f fault, body map[string]any) {
	includeUsage := false
	if so, ok := body["stream_options"].(map[string]any); ok {
		includeUsage, _ = so["include_usage"].(bool)
	}
	st.event("", chatChunk(id, a.model, map[string]any{"role": "assistant", "content": ""}, nil))
	if f.firstTokenDelay > 0 {
		time.Sleep(f.firstTokenDelay)
	}
	finish := "stop"
	if a.toolArgs != "" {
		finish = "tool_calls"
		st.event("", chatChunk(id, a.model, map[string]any{"tool_calls": []any{map[string]any{"index": 0, "id": "call_" + id, "type": "function", "function": map[string]any{"name": ToolName, "arguments": ""}}}}, nil))
		if f.truncate {
			return
		}
		for _, delta := range splitText(a.toolArgs) {
			st.event("", chatChunk(id, a.model, map[string]any{"tool_calls": []any{map[string]any{"index": 0, "function": map[string]any{"arguments": delta}}}}, nil))
		}
	} else {
		for i, delta := range splitText(a.text) {
			st.event("", chatChunk(id, a.model, map[string]any{"content": delta}, nil))
			if f.truncate && i == 0 {
				return
			}
		}
	}
	st.event("", chatChunk(id, a.model, map[string]any{}, finish))
	if includeUsage {
		st.event("", map[string]any{"id": "chatcmpl-" + id, "object": "chat.completion.chunk", "created": time.Now().Unix(), "model": a.model, "choices": []any{}, "usage": chatUsage(a)})
	}
	st.raw("[DONE]")
}

func streamMessages(st *streamWriter, id string, a answer, f fault) {
	st.event("message_start", map[string]any{"type": "message_start", "message": map[string]any{
		"id": "msg_" + id, "type": "message", "role": "assistant", "model": a.model, "content": []any{}, "stop_reason": nil, "stop_sequence": nil,
		"usage": messagesUsage(a, 0)}})
	st.event("ping", map[string]any{"type": "ping"})
	if f.firstTokenDelay > 0 {
		time.Sleep(f.firstTokenDelay)
	}
	stop := "end_turn"
	if a.toolArgs != "" {
		stop = "tool_use"
		st.event("content_block_start", map[string]any{"type": "content_block_start", "index": 0, "content_block": map[string]any{"type": "tool_use", "id": "toolu_" + id, "name": ToolName, "input": map[string]any{}}})
		if f.truncate {
			return
		}
		for _, delta := range splitText(a.toolArgs) {
			st.event("content_block_delta", map[string]any{"type": "content_block_delta", "index": 0, "delta": map[string]any{"type": "input_json_delta", "partial_json": delta}})
		}
	} else {
		st.event("content_block_start", map[string]any{"type": "content_block_start", "index": 0, "content_block": map[string]any{"type": "text", "text": ""}})
		for i, delta := range splitText(a.text) {
			st.event("content_block_delta", map[string]any{"type": "content_block_delta", "index": 0, "delta": map[string]any{"type": "text_delta", "text": delta}})
			if f.truncate && i == 0 {
				return
			}
		}
	}
	st.event("content_block_stop", map[string]any{"type": "content_block_stop", "index": 0})
	st.event("message_delta", map[string]any{"type": "message_delta", "delta": map[string]any{"stop_reason": stop, "stop_sequence": nil}, "usage": map[string]any{"output_tokens": a.output}})
	st.event("message_stop", map[string]any{"type": "message_stop"})
}

func nonStreamBody(protocol, id string, a answer) map[string]any {
	switch protocol {
	case "responses":
		return map[string]any{"id": "resp_" + id, "object": "response", "created_at": time.Now().Unix(), "status": "completed", "model": a.model, "output": responsesOutput(a, id), "usage": responsesUsage(a)}
	case "chat_completions":
		msg := map[string]any{"role": "assistant", "content": a.text}
		finish := "stop"
		if a.toolArgs != "" {
			msg["content"] = nil
			msg["tool_calls"] = []any{map[string]any{"id": "call_" + id, "type": "function", "function": map[string]any{"name": ToolName, "arguments": a.toolArgs}}}
			finish = "tool_calls"
		}
		return map[string]any{"id": "chatcmpl-" + id, "object": "chat.completion", "created": time.Now().Unix(), "model": a.model,
			"choices": []any{map[string]any{"index": 0, "message": msg, "finish_reason": finish}}, "usage": chatUsage(a)}
	case "messages":
		var content []any
		stop := "end_turn"
		if a.toolArgs != "" {
			var input map[string]any
			_ = json.Unmarshal([]byte(a.toolArgs), &input)
			content = []any{map[string]any{"type": "tool_use", "id": "toolu_" + id, "name": ToolName, "input": input}}
			stop = "tool_use"
		} else {
			content = []any{map[string]any{"type": "text", "text": a.text}}
		}
		return map[string]any{"id": "msg_" + id, "type": "message", "role": "assistant", "model": a.model, "content": content, "stop_reason": stop, "stop_sequence": nil, "usage": messagesUsage(a, a.output)}
	}
	return map[string]any{}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, protocol string, status int, code, message string) {
	if protocol == "messages" {
		writeJSON(w, status, map[string]any{"type": "error", "error": map[string]any{"type": code, "message": message}})
		return
	}
	writeJSON(w, status, map[string]any{"error": map[string]any{"type": code, "code": code, "message": message}})
}
