//go:build unit

package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// --- collector fakes -------------------------------------------------------

type hallFactRepoFake struct {
	mu       sync.Mutex
	batches  []service.ProviderHallFactBatch
	targets  []service.ProviderHallTargetRef
	probeKey map[int64]time.Time
}

func (f *hallFactRepoFake) RegisterEpoch(context.Context, string, string, int, time.Time) (int64, error) {
	return 1, nil
}
func (f *hallFactRepoFake) Heartbeat(context.Context, int64, time.Time) error { return nil }
func (f *hallFactRepoFake) MarkEpochExited(context.Context, int64, string, time.Time) error {
	return nil
}
func (f *hallFactRepoFake) WriteBatch(_ context.Context, b service.ProviderHallFactBatch) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.batches = append(f.batches, b)
	return nil
}
func (f *hallFactRepoFake) OpenGap(context.Context, service.ProviderHallGap) (int64, error) {
	return 1, nil
}
func (f *hallFactRepoFake) CloseGap(context.Context, int64, time.Time) error { return nil }
func (f *hallFactRepoFake) ListEnabledTargets(context.Context) ([]service.ProviderHallTargetRef, error) {
	return f.targets, nil
}
func (f *hallFactRepoFake) ListProbeKeys(context.Context) (map[int64]time.Time, error) {
	return f.probeKey, nil
}

// rows flattens every finish row written so far, keyed by trace.
func (f *hallFactRepoFake) rows() (starts, finishes []service.ProviderHallRequestRow, billings []service.ProviderHallBillingEvent) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, b := range f.batches {
		starts = append(starts, b.Starts...)
		finishes = append(finishes, b.Finishes...)
		billings = append(billings, b.Billings...)
	}
	return starts, finishes, billings
}

type hallCfgRepoFake struct {
	service.ProviderHallRepository
	enabled bool
}

func (f *hallCfgRepoFake) GetConfig(context.Context) (*service.ProviderHallConfig, error) {
	return &service.ProviderHallConfig{CollectionEnabled: f.enabled, DefaultProtocol: "responses", DefaultRange: "6h"}, nil
}

// --- upstream fake ---------------------------------------------------------

type hallUpstreamCall struct {
	accountID int64
	path      string
	header    http.Header
	body      []byte
}

type hallUpstreamFake struct {
	mu      sync.Mutex
	calls   []hallUpstreamCall
	respond func(call hallUpstreamCall, n int) *http.Response
}

func (u *hallUpstreamFake) do(req *http.Request, accountID int64) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
	}
	call := hallUpstreamCall{accountID: accountID, path: req.URL.Path, header: req.Header.Clone(), body: body}
	u.mu.Lock()
	u.calls = append(u.calls, call)
	n := len(u.calls)
	u.mu.Unlock()
	return u.respond(call, n), nil
}

func (u *hallUpstreamFake) Do(req *http.Request, _ string, accountID int64, _ int) (*http.Response, error) {
	return u.do(req, accountID)
}

func (u *hallUpstreamFake) DoWithTLS(req *http.Request, _ string, accountID int64, _ int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.do(req, accountID)
}

func (u *hallUpstreamFake) snapshot() []hallUpstreamCall {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]hallUpstreamCall(nil), u.calls...)
}

func hallJSON(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func hallSSE(events ...string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
		Body:       io.NopCloser(strings.NewReader(strings.Join(events, "\n\n") + "\n\n")),
	}
}

// hallAnswer serves a non-stream JSON body to non-stream upstream calls and
// the SSE fixture to streaming ones. Chat Completions and Messages entrances
// always stream from the upstream Responses endpoint, whatever the client asked.
func hallAnswer(call hallUpstreamCall, _ int) *http.Response {
	if strings.Contains(string(call.body), `"stream":true`) {
		return hallSSE(hallResponsesSSE...)
	}
	return hallJSON(http.StatusOK, hallResponsesOK)
}

const hallResponsesOK = `{"id":"resp_ok","object":"response","model":"gpt-hall","status":"completed","output":[{"type":"message","id":"msg_1","status":"completed","role":"assistant","content":[{"type":"output_text","text":"OK","annotations":[]}]}],"usage":{"input_tokens":11,"output_tokens":2,"total_tokens":13,"input_tokens_details":{"cached_tokens":4}}}`

var hallResponsesSSE = []string{
	"event: response.created\ndata: " + `{"type":"response.created","response":{"id":"resp_sse","object":"response","model":"gpt-hall","status":"in_progress"}}`,
	"event: response.output_item.added\ndata: " + `{"type":"response.output_item.added","output_index":0,"item":{"type":"message","id":"msg_1","status":"in_progress","role":"assistant","content":[]}}`,
	"event: response.output_text.delta\ndata: " + `{"type":"response.output_text.delta","item_id":"msg_1","output_index":0,"content_index":0,"delta":"OK"}`,
	"event: response.completed\ndata: " + `{"type":"response.completed","response":{"id":"resp_sse","object":"response","model":"gpt-hall","status":"completed","output":[{"type":"message","id":"msg_1","status":"completed","role":"assistant","content":[{"type":"output_text","text":"OK","annotations":[]}]}],"usage":{"input_tokens":11,"output_tokens":2,"total_tokens":13}}}`,
}

// --- harness ----------------------------------------------------------------

type hallHarness struct {
	h         *OpenAIGatewayHandler
	collector *service.ProviderHallCollector
	facts     *hallFactRepoFake
	upstream  *hallUpstreamFake
	groupID   int64
}

const (
	hallGroupID  = int64(7301)
	hallUserKey  = int64(8301)
	hallProbeKey = int64(8399)
)

func newHallHarness(t *testing.T, accounts []service.Account, respond func(hallUpstreamCall, int) *http.Response) *hallHarness {
	t.Helper()
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{RunMode: config.RunModeSimple}
	cfg.Default.RateMultiplier = 1
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Gateway.MaxAccountSwitches = 3

	upstream := &hallUpstreamFake{respond: respond}
	billingCacheSvc := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	t.Cleanup(billingCacheSvc.Stop)
	gatewaySvc := service.NewOpenAIGatewayService(
		&openAIWSFailoverHandlerAccountRepoStub{accounts: accounts},
		nil, nil, nil, nil, nil, nil, cfg, nil, nil,
		service.NewBillingService(cfg, nil), nil, billingCacheSvc, upstream, &service.DeferredService{},
		nil, nil, nil, nil, nil, nil, nil,
	)
	facts := &hallFactRepoFake{
		targets: []service.ProviderHallTargetRef{
			{TargetID: 1, GroupID: hallGroupID, ProfileID: 100, Model: "gpt-hall", Protocol: "responses", Enabled: true},
			{TargetID: 2, GroupID: hallGroupID, ProfileID: 101, Model: "gpt-hall", Protocol: "chat_completions", Enabled: true},
			{TargetID: 3, GroupID: hallGroupID, ProfileID: 102, Model: "gpt-hall", Protocol: "messages", Enabled: true},
		},
		probeKey: map[int64]time.Time{hallProbeKey: time.Now().Add(-time.Hour)},
	}
	collector := service.NewProviderHallCollector(facts, &hallCfgRepoFake{enabled: true}, "node-test", "test")
	collector.RefreshNow()
	gatewaySvc.SetProviderHallBillingSink(collector)

	h := NewOpenAIGatewayHandler(gatewaySvc, service.NewConcurrencyService(nil), billingCacheSvc,
		service.NewAPIKeyService(nil, nil, nil, nil, nil, nil, cfg), nil, nil, nil, nil, cfg)
	h.SetProviderHallCollector(collector)
	return &hallHarness{h: h, collector: collector, facts: facts, upstream: upstream, groupID: hallGroupID}
}

func hallAccount(id int64, platform string) service.Account {
	return service.Account{
		ID: id, Name: "acct", Platform: platform,
		Type: service.AccountTypeAPIKey, Status: service.StatusActive, Schedulable: true, Priority: int(id % 100),
		Credentials: map[string]any{"api_key": "sk-test", "base_url": "https://api.example.test"},
		Extra:       map[string]any{"openai_passthrough": true},
	}
}

type hallRequestOpts struct {
	keyID   int64
	header  map[string]string
	cancel  bool
	entry   string // responses | chat_completions | messages
	stream  bool
	groupID int64
}

func (hh *hallHarness) do(t *testing.T, opts hallRequestOpts) *httptest.ResponseRecorder {
	t.Helper()
	if opts.keyID == 0 {
		opts.keyID = hallUserKey
	}
	if opts.groupID == 0 {
		opts.groupID = hh.groupID
	}
	stream := "false"
	if opts.stream {
		stream = "true"
	}
	var path, body string
	switch opts.entry {
	case "chat_completions":
		path, body = "/v1/chat/completions", `{"model":"gpt-hall","messages":[{"role":"user","content":"hi"}],"stream":`+stream+`}`
	case "messages":
		path, body = "/v1/messages", `{"model":"gpt-hall","max_tokens":16,"messages":[{"role":"user","content":"hi"}],"stream":`+stream+`}`
	default:
		opts.entry = "responses"
		path, body = "/openai/v1/responses", `{"model":"gpt-hall","input":"hi","stream":`+stream+`}`
	}
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range opts.header {
		req.Header.Set(k, v)
	}
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()
	if opts.cancel {
		cancel()
	}
	c.Request = req.WithContext(ctx)
	groupID := opts.groupID
	c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
		ID: opts.keyID, GroupID: &groupID,
		User:  &service.User{ID: 1703, Status: service.StatusActive},
		Group: &service.Group{ID: groupID, Platform: service.PlatformOpenAI, Status: service.StatusActive, AllowMessagesDispatch: true},
	})
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 1703, Concurrency: 0})
	switch opts.entry {
	case "chat_completions":
		hh.h.ChatCompletions(c)
	case "messages":
		hh.h.Messages(c)
	default:
		hh.h.Responses(c)
	}
	// RecordUsage may run on a detached goroutine; give it a moment to report.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		hh.collector.FlushNow()
		_, _, billings := hh.facts.rows()
		if len(billings) > 0 || rec.Code >= 400 || opts.cancel {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	hh.collector.FlushNow()
	return rec
}

func (hh *hallHarness) singleFinish(t *testing.T) service.ProviderHallRequestRow {
	t.Helper()
	starts, finishes, _ := hh.facts.rows()
	require.Len(t, starts, 1, "exactly one start row")
	require.Len(t, finishes, 1, "exactly one terminal row")
	require.Equal(t, starts[0].TraceID, finishes[0].TraceID)
	require.Equal(t, service.ProviderHallOutcomePending, starts[0].Outcome)
	require.NotEqual(t, service.ProviderHallOutcomePending, finishes[0].Outcome)
	require.NotNil(t, finishes[0].EndedAt)
	return finishes[0]
}

// --- tests -------------------------------------------------------------------

func TestProviderHallHooks_EntrancesNonStreamSuccess(t *testing.T) {
	for _, entry := range []string{"responses", "chat_completions", "messages"} {
		t.Run(entry, func(t *testing.T) {
			hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI)}, hallAnswer)
			rec := hh.do(t, hallRequestOpts{entry: entry})
			require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
			row := hh.singleFinish(t)
			require.Equal(t, service.ProviderHallOutcomeSuccess, row.Outcome)
			require.Equal(t, service.ProviderHallExclusionNone, row.ExclusionReason)
			require.Equal(t, 1, row.Submissions)
			require.Equal(t, entry, row.Protocol)
			require.Equal(t, hallGroupID, row.GroupID)
			require.Equal(t, hallUserKey, row.APIKeyID)
			require.Equal(t, service.PlatformOpenAI, row.Platform)
			require.Equal(t, service.ProviderHallSourceUser, row.Source)
			require.NotNil(t, row.ProfileID, "entrance protocol maps to its own profile")
			require.True(t, row.UsageKnown)
			require.NotNil(t, row.InputTokens)
			require.Equal(t, int64(11), *row.InputTokens)
			require.Equal(t, "gpt-hall", row.RequestedModel)
			require.False(t, row.Stream, "client asked for a non-stream answer")
			if entry == "responses" {
				require.Nil(t, row.TTFTMs, "non-stream passthrough carries no first-content time")
			}
			_, _, billings := hh.facts.rows()
			require.Len(t, billings, 1, "billing outcome links back to the trace")
			require.Equal(t, row.TraceID, billings[0].TraceID)
			require.Equal(t, "not_applicable", billings[0].Status, "simple run mode never bills")
			for _, call := range hh.upstream.snapshot() {
				require.Empty(t, call.header.Values(service.ProviderHallTaskHeader))
			}
		})
	}
}

func TestProviderHallHooks_ResponsesSSERecordsFirstContent(t *testing.T) {
	hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI)}, func(hallUpstreamCall, int) *http.Response {
		return hallSSE(hallResponsesSSE...)
	})
	rec := hh.do(t, hallRequestOpts{entry: "responses", stream: true})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), "response.completed")
	row := hh.singleFinish(t)
	require.Equal(t, service.ProviderHallOutcomeSuccess, row.Outcome)
	require.Equal(t, 1, row.Submissions)
	require.True(t, row.Stream)
	require.NotNil(t, row.FirstContentAt)
	require.False(t, row.FirstContentAt.Before(row.StartedAt))
	require.True(t, row.UsageKnown)
	require.Equal(t, "gpt-hall", row.ResponseModel)
}

func TestProviderHallHooks_StreamInterruptedAfter200IsFailed(t *testing.T) {
	hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI)}, func(hallUpstreamCall, int) *http.Response {
		return hallSSE(hallResponsesSSE[:3]...) // created + item + delta, no terminal event
	})
	rec := hh.do(t, hallRequestOpts{entry: "responses", stream: true})
	require.Equal(t, http.StatusOK, rec.Code, "headers already sent")
	row := hh.singleFinish(t)
	require.Equal(t, service.ProviderHallOutcomeFailed, row.Outcome)
	require.Equal(t, service.ProviderHallExclusionNone, row.ExclusionReason)
	require.GreaterOrEqual(t, row.Submissions, 1)
	// The passthrough path discards the partial result on a missing terminal
	// event (openai_gateway_passthrough.go: `return nil, errors.New(...)`), so
	// the first-content time is not retained here. Recorded as a B2 limitation.
	require.False(t, row.UsageKnown)
}

func TestProviderHallHooks_RateLimitedThenSwitchedAccount(t *testing.T) {
	hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI), hallAccount(9902, service.PlatformOpenAI)}, func(call hallUpstreamCall, n int) *http.Response {
		if call.accountID == 9901 {
			return hallJSON(http.StatusTooManyRequests, `{"error":{"type":"rate_limit_error","message":"slow down"}}`)
		}
		return hallJSON(http.StatusOK, hallResponsesOK)
	})
	rec := hh.do(t, hallRequestOpts{entry: "responses"})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	calls := hh.upstream.snapshot()
	require.GreaterOrEqual(t, len(calls), 2)
	require.Equal(t, int64(9902), calls[len(calls)-1].accountID)
	row := hh.singleFinish(t)
	require.Equal(t, service.ProviderHallOutcomeSuccess, row.Outcome)
	require.Equal(t, len(calls), row.Submissions, "every real model submission counts, including the 429")
	require.Equal(t, service.PlatformOpenAI, row.Platform)
}

func TestProviderHallHooks_AllAccountsFailIsFailed(t *testing.T) {
	hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI), hallAccount(9902, service.PlatformOpenAI)}, func(hallUpstreamCall, int) *http.Response {
		return hallJSON(http.StatusBadGateway, `{"error":{"message":"temporary upstream failure"}}`)
	})
	rec := hh.do(t, hallRequestOpts{entry: "responses"})
	require.GreaterOrEqual(t, rec.Code, 500)
	row := hh.singleFinish(t)
	require.Equal(t, service.ProviderHallOutcomeFailed, row.Outcome)
	require.Equal(t, len(hh.upstream.snapshot()), row.Submissions)
	require.GreaterOrEqual(t, row.Submissions, 2, "same-account retry and account switch both count")
	require.False(t, row.UsageKnown)
}

func TestProviderHallHooks_NoAccountIsFailedWithoutSubmission(t *testing.T) {
	hh := newHallHarness(t, nil, func(hallUpstreamCall, int) *http.Response {
		t.Fatal("no upstream call expected")
		return nil
	})
	rec := hh.do(t, hallRequestOpts{entry: "responses"})
	require.NotEqual(t, http.StatusOK, rec.Code)
	row := hh.singleFinish(t)
	require.Equal(t, service.ProviderHallOutcomeFailed, row.Outcome)
	require.Zero(t, row.Submissions)
	require.Empty(t, row.Platform)
}

func TestProviderHallHooks_ClientCancelBeforeSubmitIsExcluded(t *testing.T) {
	hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI)}, func(hallUpstreamCall, int) *http.Response {
		return hallJSON(http.StatusOK, hallResponsesOK)
	})
	hh.do(t, hallRequestOpts{entry: "responses", cancel: true})
	row := hh.singleFinish(t)
	require.Equal(t, service.ProviderHallOutcomeExcluded, row.Outcome)
	require.Equal(t, service.ProviderHallExclusionClientCancel, row.ExclusionReason)
	require.Zero(t, row.Submissions)
}

func TestProviderHallHooks_UnlistedGroupIsNotTracked(t *testing.T) {
	hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI)}, func(hallUpstreamCall, int) *http.Response {
		return hallJSON(http.StatusOK, hallResponsesOK)
	})
	rec := hh.do(t, hallRequestOpts{entry: "responses", groupID: 4242})
	require.Equal(t, http.StatusOK, rec.Code)
	starts, finishes, billings := hh.facts.rows()
	require.Empty(t, starts)
	require.Empty(t, finishes)
	require.Empty(t, billings, "untracked requests never emit billing events")
}

func TestProviderHallTaskHeaderNeverForwarded(t *testing.T) {
	ref := service.ProviderHallTaskRef{Kind: service.ProviderHallSourceProbe, JobID: 5, SampleID: 9, TraceID: uuid.New()}
	header := service.FormatProviderHallTaskHeader(ref)
	for _, entry := range []string{"responses", "chat_completions", "messages"} {
		t.Run(entry+"/user_key_mismatch", func(t *testing.T) {
			hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI)}, hallAnswer)
			rec := hh.do(t, hallRequestOpts{entry: entry, header: map[string]string{service.ProviderHallTaskHeader: header}})
			require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
			calls := hh.upstream.snapshot()
			require.NotEmpty(t, calls)
			for _, call := range calls {
				for key := range call.header {
					require.NotEqual(t, strings.ToLower(service.ProviderHallTaskHeader), strings.ToLower(key), "task header leaked to upstream on %s", entry)
				}
				require.NotContains(t, string(call.body), ref.TraceID.String())
			}
			row := hh.singleFinish(t)
			require.Equal(t, service.ProviderHallSourceUser, row.Source)
			require.Equal(t, service.ProviderHallOutcomeExcluded, row.Outcome)
			require.Equal(t, service.ProviderHallExclusionTaskHeaderMismatch, row.ExclusionReason)
			require.NotEqual(t, ref.TraceID, row.TraceID)
		})
		t.Run(entry+"/probe_key", func(t *testing.T) {
			hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI)}, hallAnswer)
			rec := hh.do(t, hallRequestOpts{entry: entry, keyID: hallProbeKey, header: map[string]string{service.ProviderHallTaskHeader: header}})
			require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
			for _, call := range hh.upstream.snapshot() {
				require.Empty(t, call.header.Values(service.ProviderHallTaskHeader))
			}
			row := hh.singleFinish(t)
			require.Equal(t, service.ProviderHallSourceProbe, row.Source)
			require.Equal(t, ref.TraceID, row.TraceID, "probe rows reuse the runner's trace so reconciliation can join")
			require.NotNil(t, row.SampleID)
			require.Equal(t, int64(9), *row.SampleID)
			require.Equal(t, service.ProviderHallOutcomeSuccess, row.Outcome)
		})
	}
	t.Run("probe_key_without_header", func(t *testing.T) {
		hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI)}, func(hallUpstreamCall, int) *http.Response {
			return hallJSON(http.StatusOK, hallResponsesOK)
		})
		hh.do(t, hallRequestOpts{entry: "responses", keyID: hallProbeKey})
		row := hh.singleFinish(t)
		require.Equal(t, service.ProviderHallSourceProbe, row.Source)
		require.Nil(t, row.SampleID)
	})
}

func TestProviderHallHooks_CollectorDisabledLeavesHandlerUntouched(t *testing.T) {
	hh := newHallHarness(t, []service.Account{hallAccount(9901, service.PlatformOpenAI)}, func(hallUpstreamCall, int) *http.Response {
		return hallJSON(http.StatusOK, hallResponsesOK)
	})
	hh.h.SetProviderHallCollector(nil)
	rec := hh.do(t, hallRequestOpts{entry: "responses", header: map[string]string{service.ProviderHallTaskHeader: "probe:1:1:" + uuid.NewString()}})
	require.Equal(t, http.StatusOK, rec.Code)
	for _, call := range hh.upstream.snapshot() {
		require.Empty(t, call.header.Values(service.ProviderHallTaskHeader), "header is stripped even when collection is off")
	}
	starts, finishes, _ := hh.facts.rows()
	require.Empty(t, starts)
	require.Empty(t, finishes)
}
