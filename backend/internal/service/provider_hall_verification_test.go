//go:build unit

package service

import (
	"math/rand"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func hallResult(testID string, seq int, result, model string) ProviderHallSampleResult {
	return ProviderHallSampleResult{TestID: testID, Seq: seq, Status: ProviderHallSampleReceived, Result: result, ResponseModel: model}
}

func hallFullSuite(model string, tools bool) []ProviderHallSampleResult {
	var out []ProviderHallSampleResult
	ids := []string{ProviderHallTestArith, ProviderHallTestJSON}
	if tools {
		ids = append(ids, ProviderHallTestTool)
	}
	for _, id := range ids {
		for seq := 1; seq <= 3; seq++ {
			out = append(out, hallResult(id, seq, ProviderHallResultPassed, model))
		}
	}
	return out
}

func TestProviderHallBuildSuite(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	plain := ProviderHallBuildSuite(ProviderHallJobProfile{Model: "gpt-x", Protocol: "responses", OutputLimit: 512}, rng)
	require.Len(t, plain, 6, "no tool suite without supports_tools")
	for _, tc := range plain {
		require.NotEqual(t, ProviderHallTestTool, tc.TestID)
		require.Equal(t, 512, tc.MaxTokens, "output_limit caps max tokens")
		require.NotEmpty(t, tc.Prompt)
	}
	withTools := ProviderHallBuildSuite(ProviderHallJobProfile{Model: "gpt-x", Protocol: "responses", SupportsTools: true}, rng)
	require.Len(t, withTools, 9)
	tools := 0
	for _, tc := range withTools {
		if tc.TestID == ProviderHallTestTool {
			tools++
			require.True(t, tc.Tool)
			require.Len(t, tc.Expect.Code, 10)
		}
		if tc.TestID == ProviderHallTestArith {
			require.Equal(t, ProviderHallVerificationMaxTokens, tc.MaxTokens)
		}
	}
	require.Equal(t, 3, tools)
	probe := ProviderHallProbeSuite(ProviderHallJobProfile{OutputLimit: 100})
	require.Len(t, probe, 1)
	require.Equal(t, 100, probe[0].MaxTokens)
}

func TestProviderHallEvaluate(t *testing.T) {
	ok := func(text string) *ProviderHallProbeResponse {
		return &ProviderHallProbeResponse{Terminal: true, Text: text}
	}
	arith := ProviderHallTestCase{TestID: ProviderHallTestArith, Expect: ProviderHallExpectation{Kind: "integer", Integer: 123456}}
	jsonCase := ProviderHallTestCase{TestID: ProviderHallTestJSON, Expect: ProviderHallExpectation{Kind: "json", Token: "abc12345", Length: 8}}
	tool := ProviderHallTestCase{TestID: ProviderHallTestTool, Tool: true, Expect: ProviderHallExpectation{Kind: "tool", Code: "zz"}}
	cases := []struct {
		name   string
		tc     ProviderHallTestCase
		resp   *ProviderHallProbeResponse
		want   string
		fenced bool
	}{
		{"arith_exact", arith, ok("123456"), ProviderHallResultPassed, false},
		{"arith_thousands", arith, ok(" 123,456 \n"), ProviderHallResultPassed, false},
		{"arith_fenced", arith, ok("```\n123456\n```"), ProviderHallResultPassed, true},
		{"arith_prose", arith, ok("The answer is 123456"), ProviderHallResultFailed, false},
		{"arith_wrong", arith, ok("123457"), ProviderHallResultFailed, false},
		{"json_exact", jsonCase, ok(`{"token":"abc12345","length":8}`), ProviderHallResultPassed, false},
		{"json_fenced", jsonCase, ok("```json\n{\"length\": 8, \"token\": \"abc12345\"}\n```"), ProviderHallResultPassed, true},
		{"json_extra_key", jsonCase, ok(`{"token":"abc12345","length":8,"x":1}`), ProviderHallResultFailed, false},
		{"json_wrong_token", jsonCase, ok(`{"token":"abc12346","length":8}`), ProviderHallResultFailed, false},
		{"json_invalid", jsonCase, ok(`token abc12345`), ProviderHallResultFailed, false},
		{"tool_match", tool, &ProviderHallProbeResponse{Terminal: true, ToolCalls: []ProviderHallToolCall{{Name: "hall_echo", Arguments: `{"code":"zz"}`}}}, ProviderHallResultPassed, false},
		{"tool_wrong_code", tool, &ProviderHallProbeResponse{Terminal: true, ToolCalls: []ProviderHallToolCall{{Name: "hall_echo", Arguments: `{"code":"zy"}`}}}, ProviderHallResultFailed, false},
		{"tool_text_only", tool, ok("zz"), ProviderHallResultFailed, false},
		{"stream_incomplete", arith, &ProviderHallProbeResponse{Text: "123456"}, ProviderHallResultError, false},
		{"upstream_error", arith, &ProviderHallProbeResponse{Terminal: true, ErrorCode: "upstream_error"}, ProviderHallResultError, false},
		{"nil", arith, nil, ProviderHallResultError, false},
		{"probe_ok", ProviderHallTestCase{TestID: ProviderHallTestProbe, Expect: ProviderHallExpectation{Kind: "probe"}}, ok("OK"), ProviderHallResultPassed, false},
		{"probe_empty", ProviderHallTestCase{TestID: ProviderHallTestProbe, Expect: ProviderHallExpectation{Kind: "probe"}}, ok("  "), ProviderHallResultFailed, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, detail := ProviderHallEvaluate(tc.tc, tc.resp)
			require.Equal(t, tc.want, got)
			_, fenced := detail["fenced"]
			require.Equal(t, tc.fenced, fenced, "fence recorded but tolerated")
			if s, ok := detail["actual"].(string); ok {
				require.LessOrEqual(t, len(s), 2048)
			}
		})
	}
}

func TestProviderHallVerdictPriority(t *testing.T) {
	profile := ProviderHallJobProfile{Model: "gpt-x", Protocol: "responses", SupportsTools: true, ModelAliases: []string{"gpt-x-2025"}}
	cases := []struct {
		name      string
		samples   func() []ProviderHallSampleResult
		verdict   string
		execution string
		reason    string
	}{
		{"all_pass", func() []ProviderHallSampleResult { return hallFullSuite("gpt-x", true) }, ProviderHallVerdictPassed, ProviderHallExecutionCompleted, ""},
		{"alias_matches", func() []ProviderHallSampleResult { return hallFullSuite("gpt-x-2025", true) }, ProviderHallVerdictPassed, ProviderHallExecutionCompleted, ""},
		{"two_mismatches_plus_three_failures_is_suspected", func() []ProviderHallSampleResult {
			s := hallFullSuite("gpt-x", true)
			s[0].ResponseModel, s[1].ResponseModel = "other", "other"
			s[3].Result, s[4].Result, s[5].Result = ProviderHallResultFailed, ProviderHallResultFailed, ProviderHallResultFailed
			return s
		}, ProviderHallVerdictSuspected, ProviderHallExecutionCompleted, "model_mismatch"},
		{"same_suite_two_failures_is_failed", func() []ProviderHallSampleResult {
			s := hallFullSuite("gpt-x", true)
			s[3].Result, s[4].Result = ProviderHallResultFailed, ProviderHallResultFailed
			return s
		}, ProviderHallVerdictFailed, ProviderHallExecutionCompleted, "json_failed"},
		{"one_mismatch_does_not_suspect", func() []ProviderHallSampleResult {
			s := hallFullSuite("gpt-x", true)
			s[0].ResponseModel = "other"
			return s
		}, ProviderHallVerdictPassed, ProviderHallExecutionCompleted, ""},
		{"one_failure_plus_one_error_is_insufficient", func() []ProviderHallSampleResult {
			s := hallFullSuite("gpt-x", true)
			s[0].Result = ProviderHallResultFailed
			s[1].Result = ProviderHallResultError
			return s
		}, ProviderHallVerdictInsufficient, ProviderHallExecutionPartial, "sample_error"},
		{"missing_sample_is_insufficient", func() []ProviderHallSampleResult {
			s := hallFullSuite("gpt-x", true)
			return s[:8]
		}, ProviderHallVerdictInsufficient, ProviderHallExecutionPartial, "samples_missing"},
		{"uncertain_sample_is_insufficient", func() []ProviderHallSampleResult {
			s := hallFullSuite("gpt-x", true)
			s[8].Status = ProviderHallSampleUncertain
			return s
		}, ProviderHallVerdictInsufficient, ProviderHallExecutionPartial, "samples_missing"},
		{"missing_model_name_is_insufficient", func() []ProviderHallSampleResult {
			s := hallFullSuite("gpt-x", true)
			s[2].ResponseModel = ""
			return s
		}, ProviderHallVerdictInsufficient, ProviderHallExecutionCompleted, "model_missing"},
		{"failure_beats_insufficient", func() []ProviderHallSampleResult {
			s := hallFullSuite("gpt-x", true)
			s[0].Result, s[1].Result = ProviderHallResultFailed, ProviderHallResultFailed
			s[6].Result = ProviderHallResultError
			return s
		}, ProviderHallVerdictFailed, ProviderHallExecutionPartial, "arith_failed"},
		{"all_error_is_execution_error", func() []ProviderHallSampleResult {
			s := hallFullSuite("gpt-x", true)
			for i := range s {
				s[i].Result = ProviderHallResultError
			}
			return s
		}, ProviderHallVerdictInsufficient, ProviderHallExecutionError, "sample_error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			verdict, execution, reason, summary := ProviderHallVerdict(tc.samples(), profile)
			require.Equal(t, tc.verdict, verdict)
			require.Equal(t, tc.execution, execution)
			require.Equal(t, tc.reason, reason)
			require.NotNil(t, summary.Tool)
			require.NotEmpty(t, summary.ProfileFingerprint)
		})
	}
	t.Run("no_tool_suite_without_supports_tools", func(t *testing.T) {
		verdict, _, _, summary := ProviderHallVerdict(hallFullSuite("gpt-x", false), ProviderHallJobProfile{Model: "gpt-x", Protocol: "responses"})
		require.Equal(t, ProviderHallVerdictPassed, verdict)
		require.Nil(t, summary.Tool)
		// A stray tool sample is ignored rather than counted.
		verdict, _, _, _ = ProviderHallVerdict(append(hallFullSuite("gpt-x", false), hallResult(ProviderHallTestTool, 1, ProviderHallResultFailed, "gpt-x"), hallResult(ProviderHallTestTool, 2, ProviderHallResultFailed, "gpt-x")), ProviderHallJobProfile{Model: "gpt-x", Protocol: "responses"})
		require.Equal(t, ProviderHallVerdictPassed, verdict)
	})
}

func TestProviderHallVerificationExpiry(t *testing.T) {
	now := time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	v := &ProviderHallVerification{CompletedAt: now, ExpiresAt: now.Add(ProviderHallVerificationTTL)}
	require.False(t, v.Expired(now.Add(47*time.Hour)))
	require.True(t, v.Expired(now.Add(48*time.Hour)))
	v.Stale = true
	require.True(t, v.Expired(now))
	var nilV *ProviderHallVerification
	require.True(t, nilV.Expired(now))
}

func TestProviderHallProfileFingerprint(t *testing.T) {
	a := ProviderHallProfileFingerprint("m", "responses", true, []string{"b", "a"})
	b := ProviderHallProfileFingerprint("m", "responses", true, []string{"a", "b"})
	require.Equal(t, a, b, "alias order does not matter")
	require.NotEqual(t, a, ProviderHallProfileFingerprint("m", "responses", false, []string{"a", "b"}))
	require.NotEqual(t, a, ProviderHallProfileFingerprint("m", "chat_completions", true, []string{"a", "b"}))
}

func TestProviderHallBudgetHelpers(t *testing.T) {
	require.False(t, ProviderHallBudgetExhausted("0", decimalFromString(t, "100")), "zero budget is unlimited")
	require.False(t, ProviderHallBudgetExhausted("10", decimalFromString(t, "9.99999999")))
	require.True(t, ProviderHallBudgetExhausted("10", decimalFromString(t, "10")))
	require.True(t, ProviderHallBudgetExhausted("10", decimalFromString(t, "10.5")))
	// Shanghai midnight: 15:59:59Z is still the same day; 16:00:00Z is the next.
	require.Equal(t, "2026-09-12", ProviderHallBudgetDay(time.Date(2026, 9, 12, 15, 59, 59, 0, time.UTC)))
	require.Equal(t, "2026-09-13", ProviderHallBudgetDay(time.Date(2026, 9, 12, 16, 0, 0, 0, time.UTC)))
	day := "2026-09-01"
	require.Equal(t, day, ProviderHallBudgetDayFor(&day, time.Date(2026, 9, 12, 16, 0, 0, 0, time.UTC)), "job keeps its first dispatch day")
	for _, tc := range []struct {
		in            string
		spend, sample string
		write, ok     bool
	}{
		{"applied", ProviderHallSpendConfirmed, ProviderHallSampleBillingConfirmed, true, true},
		{"failed", ProviderHallSpendFailed, ProviderHallSampleBillingFailed, true, true},
		{"uncertain", ProviderHallSpendUncertain, ProviderHallSampleBillingUncertain, true, true},
		{"not_applicable", "", ProviderHallSampleBillingUnbilled, false, true},
		{"duplicate", "", "", false, false},
	} {
		spend, sample, write, ok := ProviderHallSpendStatusForBilling(tc.in)
		require.Equal(t, tc.spend, spend, tc.in)
		require.Equal(t, tc.sample, sample, tc.in)
		require.Equal(t, tc.write, write, tc.in)
		require.Equal(t, tc.ok, ok, tc.in)
	}
}

func decimalFromString(t *testing.T, s string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(s)
	require.NoError(t, err)
	return d
}
