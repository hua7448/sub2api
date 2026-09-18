//go:build unit

package fakeopenai_test

import (
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/Wei-Shaw/sub2api/internal/testutil/fakeopenai"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var protocols = []string{"responses", "chat_completions", "messages"}

func newClient(t *testing.T) *service.ProviderHallProbeClient {
	t.Helper()
	t.Setenv("PROVIDER_HALL_ALLOW_LOOPBACK", "1")
	return service.NewProviderHallProbeClient(10 * time.Second)
}

func probeReq(origin, protocol, model string, tc service.ProviderHallTestCase) service.ProviderHallProbeRequest {
	return service.ProviderHallProbeRequest{Origin: origin, Protocol: protocol, Model: model, Key: "sk-fake",
		Task: service.ProviderHallTaskRef{Kind: service.ProviderHallSourceProbe, JobID: 1, SampleID: 1, TraceID: uuid.New()}, Case: tc, Version: "test"}
}

func TestFakeOpenAI_ProbeStreamsParseOnEveryProtocol(t *testing.T) {
	srv := fakeopenai.New(fakeopenai.Options{InputTokens: 21, OutputTokens: 7, CacheReadTokens: 5, ModelAlias: map[string]string{"gpt-5": "gpt-5-2026-01-01"}})
	defer srv.Close()
	client := newClient(t)
	profile := service.ProviderHallJobProfile{Model: "gpt-5", ModelAliases: []string{"gpt-5-2026-01-01"}}
	for _, protocol := range protocols {
		t.Run(protocol, func(t *testing.T) {
			tc := service.ProviderHallProbeSuite(profile)[0]
			res := client.Send(context.Background(), probeReq(srv.URL, protocol, "gpt-5", tc))
			require.False(t, res.Uncertain, "%+v", res)
			require.Empty(t, res.ErrorCode)
			require.Equal(t, http.StatusOK, res.HTTPStatus)
			require.NotNil(t, res.Response)
			require.True(t, res.Response.Terminal)
			require.Equal(t, "gpt-5-2026-01-01", res.Response.Model)
			require.Equal(t, "OK", res.Response.Text)
			require.NotNil(t, res.Response.TTFTMs)
			require.NotNil(t, res.Response.TotalMs)
			require.NotNil(t, res.Response.InputTokens)
			require.Equal(t, 21, *res.Response.InputTokens)
			require.NotNil(t, res.Response.OutputTokens)
			require.Equal(t, 7, *res.Response.OutputTokens)
			result, _ := service.ProviderHallEvaluate(tc, res.Response)
			require.Equal(t, service.ProviderHallResultPassed, result)
			require.Equal(t, "matched", service.ProviderHallModelMatch(profile, res.Response.Model))
		})
	}
	reqs := srv.Requests()
	require.Len(t, reqs, 3)
	for _, r := range reqs {
		require.Equal(t, "Bearer sk-fake", r.Header.Get("Authorization"))
		require.Equal(t, "gpt-5", r.Body["model"])
		require.NotEmpty(t, r.Protocol)
	}
	srv.Reset()
	require.Zero(t, srv.Count())
}

func TestFakeOpenAI_VerificationSuitePassesByDefault(t *testing.T) {
	srv := fakeopenai.New(fakeopenai.Options{})
	defer srv.Close()
	client := newClient(t)
	for _, protocol := range protocols {
		t.Run(protocol, func(t *testing.T) {
			profile := service.ProviderHallJobProfile{Model: "gpt-5", Protocol: protocol, SupportsTools: true}
			cases := service.ProviderHallBuildSuite(profile, rand.New(rand.NewSource(7)))
			require.Len(t, cases, 9)
			var samples []service.ProviderHallSampleResult
			for _, tc := range cases {
				res := client.Send(context.Background(), probeReq(srv.URL, protocol, "gpt-5", tc))
				require.Empty(t, res.ErrorCode, "%s/%d: %+v", tc.TestID, tc.Seq, res)
				result, detail := service.ProviderHallEvaluate(tc, res.Response)
				require.Equal(t, service.ProviderHallResultPassed, result, "%s/%d detail=%v text=%q tools=%v", tc.TestID, tc.Seq, detail, res.Response.Text, res.Response.ToolCalls)
				if tc.Tool {
					require.Len(t, res.Response.ToolCalls, 1)
					require.Equal(t, fakeopenai.ToolName, res.Response.ToolCalls[0].Name)
				}
				samples = append(samples, service.ProviderHallSampleResult{TestID: tc.TestID, Seq: tc.Seq, Status: service.ProviderHallSampleReceived, Result: result, ResponseModel: res.Response.Model})
			}
			verdict, execution, reason, _ := service.ProviderHallVerdict(samples, profile)
			require.Equal(t, service.ProviderHallVerdictPassed, verdict, reason)
			require.Equal(t, service.ProviderHallExecutionCompleted, execution)
		})
	}
}

func TestFakeOpenAI_KnobsProduceEachVerdict(t *testing.T) {
	client := newClient(t)
	run := func(t *testing.T, opts fakeopenai.Options, protocol string) (string, string) {
		srv := fakeopenai.New(opts)
		defer srv.Close()
		profile := service.ProviderHallJobProfile{Model: "gpt-5", Protocol: protocol, SupportsTools: true}
		var samples []service.ProviderHallSampleResult
		for _, tc := range service.ProviderHallBuildSuite(profile, rand.New(rand.NewSource(3))) {
			res := client.Send(context.Background(), probeReq(srv.URL, protocol, "gpt-5", tc))
			result, _ := service.ProviderHallEvaluate(tc, res.Response)
			model := ""
			if res.Response != nil {
				model = res.Response.Model
			}
			samples = append(samples, service.ProviderHallSampleResult{TestID: tc.TestID, Seq: tc.Seq, Status: service.ProviderHallSampleReceived, Result: result, ResponseModel: model})
		}
		verdict, _, reason, _ := service.ProviderHallVerdict(samples, profile)
		return verdict, reason
	}
	for _, protocol := range protocols {
		t.Run(protocol, func(t *testing.T) {
			v, r := run(t, fakeopenai.Options{WrongArithmetic: true}, protocol)
			require.Equal(t, service.ProviderHallVerdictFailed, v)
			require.Equal(t, "arith_failed", r)
			v, r = run(t, fakeopenai.Options{WrongJSON: true}, protocol)
			require.Equal(t, service.ProviderHallVerdictFailed, v)
			require.Equal(t, "json_failed", r)
			v, r = run(t, fakeopenai.Options{WrongTool: true}, protocol)
			require.Equal(t, service.ProviderHallVerdictFailed, v)
			require.Equal(t, "tool_failed", r)
			v, r = run(t, fakeopenai.Options{MismatchModel: "other-model"}, protocol)
			require.Equal(t, service.ProviderHallVerdictSuspected, v)
			require.Equal(t, "model_mismatch", r)
		})
	}
}

func TestFakeOpenAI_FaultInjection(t *testing.T) {
	srv := fakeopenai.New(fakeopenai.Options{})
	defer srv.Close()
	client := newClient(t)
	tc := service.ProviderHallProbeSuite(service.ProviderHallJobProfile{})[0]
	post := func(protocol, fault string, stream bool) *http.Response {
		body, err := service.ProviderHallBuildRequestBody(protocol, "gpt-5", tc)
		require.NoError(t, err)
		if !stream {
			body = []byte(strings.Replace(string(body), `"stream":true`, `"stream":false`, 1))
		}
		path, _ := service.ProviderHallProtocolPath(protocol)
		req, _ := http.NewRequest(http.MethodPost, srv.URL+path, strings.NewReader(string(body)))
		req.Header.Set("Content-Type", "application/json")
		if fault != "" {
			req.Header.Set(fakeopenai.FaultHeader, fault)
		}
		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		return res
	}
	for _, protocol := range protocols {
		t.Run(protocol+"/status_header", func(t *testing.T) {
			res := post(protocol, "status=503", true)
			defer res.Body.Close()
			require.Equal(t, 503, res.StatusCode)
			sent := client.Send(context.Background(), service.ProviderHallProbeRequest{Origin: srv.URL, Protocol: protocol, Model: "gpt-5", Key: "k", Case: tc,
				Task: service.ProviderHallTaskRef{Kind: service.ProviderHallSourceProbe, JobID: 1, SampleID: 2, TraceID: uuid.New()}})
			require.Equal(t, http.StatusOK, sent.HTTPStatus, "no header → healthy")
		})
		t.Run(protocol+"/truncate_header", func(t *testing.T) {
			res := post(protocol, "truncate=1", true)
			defer res.Body.Close()
			require.Equal(t, 200, res.StatusCode)
			parsed, err := service.ProviderHallParseStream(protocol, res.Body, time.Now(), nil)
			require.NoError(t, err)
			require.False(t, parsed.Terminal)
			require.Equal(t, service.ProviderHallErrStreamIncomplete, parsed.ErrorCode)
			require.NotNil(t, parsed.TTFTMs, "first content was delivered before the cut")
		})
		t.Run(protocol+"/delay_header", func(t *testing.T) {
			start := time.Now()
			res := post(protocol, "delay=120;first_token_delay=80", true)
			defer res.Body.Close()
			parsed, err := service.ProviderHallParseStream(protocol, res.Body, start, nil)
			require.NoError(t, err)
			require.True(t, parsed.Terminal)
			require.GreaterOrEqual(t, *parsed.TTFTMs, 190)
		})
		t.Run(protocol+"/non_stream", func(t *testing.T) {
			res := post(protocol, "", false)
			defer res.Body.Close()
			require.Equal(t, 200, res.StatusCode)
			require.Contains(t, res.Header.Get("Content-Type"), "application/json")
			raw, _ := io.ReadAll(res.Body)
			var payload map[string]any
			require.NoError(t, json.Unmarshal(raw, &payload))
			require.Equal(t, "gpt-5", payload["model"])
			require.Contains(t, string(raw), `"OK"`)
			require.Contains(t, string(raw), "tokens")
		})
	}
	t.Run("429_rotation_and_reset", func(t *testing.T) {
		srv.Reset()
		srv.SetOptions(fakeopenai.Options{Status429Every: 2})
		defer srv.SetOptions(fakeopenai.Options{})
		codes := []int{}
		for i := 0; i < 4; i++ {
			res := post("chat_completions", "", true)
			codes = append(codes, res.StatusCode)
			_, _ = io.Copy(io.Discard, res.Body)
			res.Body.Close()
		}
		require.Equal(t, []int{200, 429, 200, 429}, codes)
		require.Equal(t, 4, srv.Count())
	})
	t.Run("unknown_path_and_auth", func(t *testing.T) {
		res, err := http.Get(srv.URL + "/v1/models")
		require.NoError(t, err)
		res.Body.Close()
		require.Equal(t, 404, res.StatusCode)
		srv.SetOptions(fakeopenai.Options{RequireAuth: true})
		defer srv.SetOptions(fakeopenai.Options{})
		res = post("responses", "", true)
		res.Body.Close()
		require.Equal(t, 401, res.StatusCode)
	})
}
