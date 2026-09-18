package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Test identifiers. "probe" is the single availability request; the other
// three form the verification suite.
const (
	ProviderHallTestProbe = "probe"
	ProviderHallTestArith = "arith"
	ProviderHallTestJSON  = "json"
	ProviderHallTestTool  = "tool"

	ProviderHallProbeMaxTokens        = 256
	ProviderHallVerificationMaxTokens = 1024
	ProviderHallVerificationTTL       = 48 * time.Hour
	ProviderHallSuiteRepeats          = 3
	ProviderHallToolName              = "hall_echo"
	providerHallProbePrompt           = "Reply with exactly the word OK."
	providerHallDetailLimit           = 2048
)

// Verdicts and execution statuses.
const (
	ProviderHallVerdictPassed       = "passed"
	ProviderHallVerdictFailed       = "failed"
	ProviderHallVerdictSuspected    = "suspected"
	ProviderHallVerdictInsufficient = "insufficient"

	ProviderHallExecutionCompleted = "completed"
	ProviderHallExecutionPartial   = "partial"
	ProviderHallExecutionError     = "error"
)

// ProviderHallExpectation is the deterministic answer a test case requires.
type ProviderHallExpectation struct {
	Kind    string `json:"kind"` // probe | integer | json | tool
	Integer int64  `json:"integer,omitempty"`
	Token   string `json:"token,omitempty"`
	Length  int    `json:"length,omitempty"`
	Code    string `json:"code,omitempty"`
}

type ProviderHallTestCase struct {
	TestID    string                  `json:"test_id"`
	Seq       int                     `json:"seq"`
	Prompt    string                  `json:"prompt"`
	MaxTokens int                     `json:"max_tokens"`
	Tool      bool                    `json:"tool"`
	Expect    ProviderHallExpectation `json:"expect"`
}

// ProviderHallToolSchema is the hall_echo parameter schema used by all protocols.
var ProviderHallToolSchema = json.RawMessage(`{"type":"object","properties":{"code":{"type":"string"}},"required":["code"],"additionalProperties":false}`)

// ProviderHallProbeSuite is the single availability request.
func ProviderHallProbeSuite(profile ProviderHallJobProfile) []ProviderHallTestCase {
	return []ProviderHallTestCase{{TestID: ProviderHallTestProbe, Seq: 1, Prompt: providerHallProbePrompt,
		MaxTokens: providerHallMaxTokens(profile, ProviderHallProbeMaxTokens), Expect: ProviderHallExpectation{Kind: "probe"}}}
}

func providerHallMaxTokens(profile ProviderHallJobProfile, wanted int) int {
	if profile.OutputLimit > 0 && profile.OutputLimit < wanted {
		return profile.OutputLimit
	}
	return wanted
}

// ProviderHallBuildSuite creates the verification cases: arithmetic ×3, JSON
// constraint ×3 and, when the profile declares tools, forced tool call ×3.
func ProviderHallBuildSuite(profile ProviderHallJobProfile, rng *rand.Rand) []ProviderHallTestCase {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	maxTokens := providerHallMaxTokens(profile, ProviderHallVerificationMaxTokens)
	cases := make([]ProviderHallTestCase, 0, 3*ProviderHallSuiteRepeats)
	ops := []byte{'+', '-', '*'}
	for seq := 1; seq <= ProviderHallSuiteRepeats; seq++ {
		a, b := int64(100+rng.Intn(900)), int64(100+rng.Intn(900))
		op := ops[rng.Intn(len(ops))]
		var want int64
		switch op {
		case '+':
			want = a + b
		case '-':
			want = a - b
		default:
			want = a * b
		}
		cases = append(cases, ProviderHallTestCase{TestID: ProviderHallTestArith, Seq: seq, MaxTokens: maxTokens,
			Prompt: fmt.Sprintf("Compute %d %c %d. Reply with only the integer result and nothing else.", a, op, b),
			Expect: ProviderHallExpectation{Kind: "integer", Integer: want}})
	}
	for seq := 1; seq <= ProviderHallSuiteRepeats; seq++ {
		token := providerHallRandomToken(rng, 8)
		cases = append(cases, ProviderHallTestCase{TestID: ProviderHallTestJSON, Seq: seq, MaxTokens: maxTokens,
			Prompt: fmt.Sprintf(`Return only a JSON object with exactly two keys: "token" set to the string %q and "length" set to the integer %d. No prose, no markdown, nothing else.`, token, len(token)),
			Expect: ProviderHallExpectation{Kind: "json", Token: token, Length: len(token)}})
	}
	if profile.SupportsTools {
		for seq := 1; seq <= ProviderHallSuiteRepeats; seq++ {
			code := providerHallRandomToken(rng, 10)
			cases = append(cases, ProviderHallTestCase{TestID: ProviderHallTestTool, Seq: seq, MaxTokens: maxTokens, Tool: true,
				Prompt: fmt.Sprintf("Call the %s tool with the code %q. Do not answer in text.", ProviderHallToolName, code),
				Expect: ProviderHallExpectation{Kind: "tool", Code: code}})
		}
	}
	return cases
}

func providerHallRandomToken(rng *rand.Rand, n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rng.Intn(len(alphabet))]
	}
	return string(b)
}

// ProviderHallToolCall is one tool invocation seen in a response.
type ProviderHallToolCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ProviderHallProbeResponse is the protocol-neutral parse of one streamed reply.
type ProviderHallProbeResponse struct {
	HTTPStatus   int
	Terminal     bool
	Model        string
	Text         string
	ToolCalls    []ProviderHallToolCall
	InputTokens  *int
	OutputTokens *int
	TTFTMs       *int
	TotalMs      *int
	ErrorCode    string
}

var providerHallFencePattern = regexp.MustCompile("(?s)^```[a-zA-Z0-9_-]*\\s*\\n?(.*?)\\n?```$")
var providerHallIntegerPattern = regexp.MustCompile(`^[-+]?\d+$`)

func providerHallTruncate(s string) string {
	if len(s) > providerHallDetailLimit {
		return s[:providerHallDetailLimit]
	}
	return s
}

// ProviderHallEvaluate scores a parsed response against its test case. The
// result is passed/failed/error; detail never includes user content.
func ProviderHallEvaluate(tc ProviderHallTestCase, resp *ProviderHallProbeResponse) (string, map[string]any) {
	detail := map[string]any{"test_id": tc.TestID, "seq": tc.Seq}
	if resp == nil {
		detail["error"] = "no_response"
		return ProviderHallResultError, detail
	}
	if resp.ErrorCode != "" {
		detail["error"] = resp.ErrorCode
		return ProviderHallResultError, detail
	}
	if !resp.Terminal {
		detail["error"] = "stream_incomplete"
		return ProviderHallResultError, detail
	}
	switch tc.Expect.Kind {
	case "probe":
		if strings.TrimSpace(resp.Text) == "" && len(resp.ToolCalls) == 0 {
			detail["error"] = "empty_output"
			return ProviderHallResultFailed, detail
		}
		return ProviderHallResultPassed, detail
	case "integer":
		raw := strings.TrimSpace(providerHallStripFence(resp.Text, detail))
		cleaned := strings.NewReplacer(",", "", " ", "", "\u00a0", "", "_", "").Replace(raw)
		cleaned = strings.TrimSuffix(cleaned, ".")
		detail["expected"] = tc.Expect.Integer
		detail["actual"] = providerHallTruncate(raw)
		if !providerHallIntegerPattern.MatchString(cleaned) {
			return ProviderHallResultFailed, detail
		}
		got, err := strconv.ParseInt(cleaned, 10, 64)
		if err != nil || got != tc.Expect.Integer {
			return ProviderHallResultFailed, detail
		}
		return ProviderHallResultPassed, detail
	case "json":
		raw := strings.TrimSpace(providerHallStripFence(resp.Text, detail))
		detail["expected"] = map[string]any{"token": tc.Expect.Token, "length": tc.Expect.Length}
		detail["actual"] = providerHallTruncate(raw)
		var obj map[string]json.RawMessage
		if err := json.Unmarshal([]byte(raw), &obj); err != nil || len(obj) != 2 {
			return ProviderHallResultFailed, detail
		}
		var token string
		var length int
		if err := json.Unmarshal(obj["token"], &token); err != nil {
			return ProviderHallResultFailed, detail
		}
		if err := json.Unmarshal(obj["length"], &length); err != nil {
			return ProviderHallResultFailed, detail
		}
		if token != tc.Expect.Token || length != tc.Expect.Length {
			return ProviderHallResultFailed, detail
		}
		return ProviderHallResultPassed, detail
	case "tool":
		detail["expected"] = map[string]any{"name": ProviderHallToolName, "code": tc.Expect.Code}
		names := make([]string, 0, len(resp.ToolCalls))
		for _, call := range resp.ToolCalls {
			names = append(names, call.Name)
			if call.Name != ProviderHallToolName {
				continue
			}
			var args map[string]json.RawMessage
			if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil || len(args) != 1 {
				continue
			}
			var code string
			if err := json.Unmarshal(args["code"], &code); err == nil && code == tc.Expect.Code {
				return ProviderHallResultPassed, detail
			}
		}
		detail["actual"] = map[string]any{"tool_calls": names, "text": providerHallTruncate(resp.Text)}
		return ProviderHallResultFailed, detail
	}
	detail["error"] = "unknown_expectation"
	return ProviderHallResultError, detail
}

// providerHallStripFence tolerates a Markdown code fence around the answer and
// records that it was present without failing the case.
func providerHallStripFence(text string, detail map[string]any) string {
	trimmed := strings.TrimSpace(text)
	if m := providerHallFencePattern.FindStringSubmatch(trimmed); m != nil {
		detail["fenced"] = true
		return strings.TrimSpace(m[1])
	}
	return trimmed
}

// ProviderHallModelMatch compares the response model with the profile model
// and its aliases: matched / mismatched / missing.
func ProviderHallModelMatch(profile ProviderHallJobProfile, responseModel string) string {
	got := strings.TrimSpace(responseModel)
	if got == "" {
		return "missing"
	}
	if got == profile.Model {
		return "matched"
	}
	for _, alias := range profile.ModelAliases {
		if got == alias {
			return "matched"
		}
	}
	return "mismatched"
}

// ProviderHallSampleResult is the verdict input for one sample.
type ProviderHallSampleResult struct {
	TestID        string
	Seq           int
	Status        ProviderHallSampleStatus
	Result        string // passed | failed | error | ""
	ResponseModel string
}

type ProviderHallSuiteSummary struct {
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Error   int `json:"error"`
	Missing int `json:"missing"`
}

type ProviderHallModelSummary struct {
	Matched    int      `json:"matched"`
	Mismatched int      `json:"mismatched"`
	Missing    int      `json:"missing"`
	Seen       []string `json:"seen"`
}

type ProviderHallVerificationSummary struct {
	Arithmetic ProviderHallSuiteSummary  `json:"arithmetic"`
	JSON       ProviderHallSuiteSummary  `json:"json"`
	Tool       *ProviderHallSuiteSummary `json:"tool"`
	Model      ProviderHallModelSummary  `json:"model"`
	// ProfileFingerprint lets the runner mark a report stale when the identity
	// fields of the profile change, independent of the version counter.
	ProfileFingerprint string `json:"profile_fingerprint"`
}

// ProviderHallProfileFingerprint covers the fields a verification depends on.
func ProviderHallProfileFingerprint(model, protocol string, supportsTools bool, aliases []string) string {
	sorted := append([]string(nil), aliases...)
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j] < sorted[j-1]; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}
	return fmt.Sprintf("%s|%s|%t|%s", model, protocol, supportsTools, strings.Join(sorted, ","))
}

// ProviderHallVerdict applies the fixed priority: repeated model mismatch →
// suspected; any required suite failing twice → failed; a missing sample,
// error or missing model name → insufficient; otherwise passed.
func ProviderHallVerdict(samples []ProviderHallSampleResult, profile ProviderHallJobProfile) (verdict, execution, reason string, summary ProviderHallVerificationSummary) {
	summary.ProfileFingerprint = ProviderHallProfileFingerprint(profile.Model, profile.Protocol, profile.SupportsTools, profile.ModelAliases)
	if profile.SupportsTools {
		summary.Tool = &ProviderHallSuiteSummary{}
	}
	suites := map[string]*ProviderHallSuiteSummary{ProviderHallTestArith: &summary.Arithmetic, ProviderHallTestJSON: &summary.JSON}
	if summary.Tool != nil {
		suites[ProviderHallTestTool] = summary.Tool
	}
	seenModels := map[string]bool{}
	present := map[string]int{}
	for _, s := range samples {
		suite := suites[s.TestID]
		if suite == nil {
			continue
		}
		present[s.TestID]++
		switch {
		case s.Status != ProviderHallSampleReceived:
			suite.Missing++
		case s.Result == ProviderHallResultPassed:
			suite.Passed++
		case s.Result == ProviderHallResultFailed:
			suite.Failed++
		default:
			suite.Error++
		}
		if s.Status == ProviderHallSampleReceived && s.Result != ProviderHallResultError {
			switch ProviderHallModelMatch(profile, s.ResponseModel) {
			case "matched":
				summary.Model.Matched++
			case "mismatched":
				summary.Model.Mismatched++
			default:
				summary.Model.Missing++
			}
			if m := strings.TrimSpace(s.ResponseModel); m != "" && !seenModels[m] {
				seenModels[m] = true
				summary.Model.Seen = append(summary.Model.Seen, m)
			}
		}
	}
	if summary.Model.Seen == nil {
		summary.Model.Seen = []string{}
	}
	for id, suite := range suites {
		if present[id] < ProviderHallSuiteRepeats {
			suite.Missing += ProviderHallSuiteRepeats - present[id]
		}
	}
	incomplete := false
	errored := false
	for _, suite := range suites {
		if suite.Missing > 0 {
			incomplete = true
		}
		if suite.Error > 0 {
			errored = true
		}
	}
	switch {
	case !incomplete && !errored:
		execution = ProviderHallExecutionCompleted
	case summary.Arithmetic.Passed+summary.Arithmetic.Failed+summary.JSON.Passed+summary.JSON.Failed == 0:
		execution = ProviderHallExecutionError
	default:
		execution = ProviderHallExecutionPartial
	}
	if summary.Model.Mismatched >= 2 {
		return ProviderHallVerdictSuspected, execution, "model_mismatch", summary
	}
	for _, id := range []string{ProviderHallTestArith, ProviderHallTestJSON, ProviderHallTestTool} {
		if suite := suites[id]; suite != nil && suite.Failed >= 2 {
			return ProviderHallVerdictFailed, execution, id + "_failed", summary
		}
	}
	if incomplete {
		return ProviderHallVerdictInsufficient, execution, "samples_missing", summary
	}
	if errored {
		return ProviderHallVerdictInsufficient, execution, "sample_error", summary
	}
	if summary.Model.Missing > 0 {
		return ProviderHallVerdictInsufficient, execution, "model_missing", summary
	}
	return ProviderHallVerdictPassed, execution, "", summary
}
