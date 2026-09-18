package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/google/uuid"
)

// ProviderHallTaskHeader carries the runner's task correlation for probe and
// verification requests. Handlers remove it from the inbound request before
// any upstream builder runs; it never leaves this process.
const ProviderHallTaskHeader = "X-Provider-Hall-Task"

// ProviderHallAlgorithmVersion is stamped on every fact and aggregate so a
// formula change can never mix with rows computed under the old one.
const ProviderHallAlgorithmVersion = 1

type ProviderHallOutcome string

const (
	ProviderHallOutcomePending  ProviderHallOutcome = "pending"
	ProviderHallOutcomeSuccess  ProviderHallOutcome = "success"
	ProviderHallOutcomeFailed   ProviderHallOutcome = "failed"
	ProviderHallOutcomeExcluded ProviderHallOutcome = "excluded"
)

type ProviderHallExclusion string

const (
	ProviderHallExclusionNone               ProviderHallExclusion = ""
	ProviderHallExclusionClientCancel       ProviderHallExclusion = "client_cancel"
	ProviderHallExclusionPolicyReject       ProviderHallExclusion = "policy_reject"
	ProviderHallExclusionInvalidRequest     ProviderHallExclusion = "invalid_request"
	ProviderHallExclusionTaskHeaderMismatch ProviderHallExclusion = "task_header_mismatch"
	ProviderHallExclusionEarlyExit          ProviderHallExclusion = "early_exit"
)

type ProviderHallSource string

const (
	ProviderHallSourceUser         ProviderHallSource = "user"
	ProviderHallSourceProbe        ProviderHallSource = "probe"
	ProviderHallSourceVerification ProviderHallSource = "verification"
)

// ProviderHallTaskRef is the parsed task header. The runner writes the same
// trace ID into its sample row so reconciliation can join on requests.trace_id.
type ProviderHallTaskRef struct {
	Kind     ProviderHallSource
	JobID    int64
	SampleID int64
	TraceID  uuid.UUID
}

func FormatProviderHallTaskHeader(ref ProviderHallTaskRef) string {
	return string(ref.Kind) + ":" + strconv.FormatInt(ref.JobID, 10) + ":" + strconv.FormatInt(ref.SampleID, 10) + ":" + ref.TraceID.String()
}

func ParseProviderHallTaskHeader(value string) (ProviderHallTaskRef, bool) {
	parts := strings.Split(strings.TrimSpace(value), ":")
	if len(parts) != 4 {
		return ProviderHallTaskRef{}, false
	}
	kind := ProviderHallSource(parts[0])
	if kind != ProviderHallSourceProbe && kind != ProviderHallSourceVerification {
		return ProviderHallTaskRef{}, false
	}
	jobID, err1 := strconv.ParseInt(parts[1], 10, 64)
	sampleID, err2 := strconv.ParseInt(parts[2], 10, 64)
	traceID, err3 := uuid.Parse(parts[3])
	if err1 != nil || err2 != nil || err3 != nil || jobID < 1 || sampleID < 1 {
		return ProviderHallTaskRef{}, false
	}
	return ProviderHallTaskRef{Kind: kind, JobID: jobID, SampleID: sampleID, TraceID: traceID}, true
}

// ProviderHallRequestRow is one fact row. Start events carry identity only;
// finish events carry the full outcome. Billing arrives separately.
type ProviderHallRequestRow struct {
	TraceID          uuid.UUID
	NodeID           string
	EpochID          int64
	Seq              uint64
	GroupID          int64
	ProfileID        *int64
	APIKeyID         int64
	Protocol         string
	RequestedModel   string
	UpstreamModel    string
	ResponseModel    string
	Platform         string
	Source           ProviderHallSource
	SampleID         *int64
	Stream           bool
	StartedAt        time.Time
	FirstContentAt   *time.Time
	EndedAt          *time.Time
	TTFTMs           *int
	Submissions      int
	Outcome          ProviderHallOutcome
	ExclusionReason  ProviderHallExclusion
	UsageKnown       bool
	InputTokens      *int64
	OutputTokens     *int64
	CacheReadTokens  *int64
	CacheCreation    *int64
	AlgorithmVersion int
}

type ProviderHallBeginInput struct {
	APIKeyID       int64
	GroupID        int64
	Protocol       string
	RequestedModel string
	Stream         bool
	StartedAt      time.Time
	TaskHeader     string
}

type providerHallTrackerCtxKey struct{}

// ProviderHallRequestTracker is created per gateway request. Every method is
// safe on a nil receiver so handlers never branch on collection being enabled.
type ProviderHallRequestTracker struct {
	mu   sync.Mutex
	sink *ProviderHallCollector
	ctx  context.Context
	row  ProviderHallRequestRow

	submissions    int
	firstContentAt *time.Time
	success        bool
	failed         bool
	noAccount      bool
	clientCancel   bool
	reject         ProviderHallExclusion
	usage          *OpenAIUsage
	finished       bool
}

func ProviderHallTrackerFromContext(ctx context.Context) *ProviderHallRequestTracker {
	if ctx == nil {
		return nil
	}
	t, _ := ctx.Value(providerHallTrackerCtxKey{}).(*ProviderHallRequestTracker)
	return t
}

func (t *ProviderHallRequestTracker) TraceID() uuid.UUID {
	if t == nil {
		return uuid.Nil
	}
	return t.row.TraceID
}

func (t *ProviderHallRequestTracker) Source() ProviderHallSource {
	if t == nil {
		return ""
	}
	return t.row.Source
}

// SelectAccount records the platform of the first account that receives a
// forward attempt. It does not count a submission.
func (t *ProviderHallRequestTracker) SelectAccount(account *Account) {
	if t == nil || account == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.finished && t.row.Platform == "" {
		t.row.Platform = account.Platform
	}
}

// RecordSubmission counts one real model submission: an upstream HTTP model
// call or a WebSocket response.create. Account selection, profit vetoes, token
// refreshes and capability probes never reach it.
func (t *ProviderHallRequestTracker) RecordSubmission() {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.finished {
		t.submissions++
	}
}

// providerHallUpstream counts submissions at the real send boundary. Only
// model endpoints count; token, quota and capability calls pass through.
type providerHallUpstream struct{ base HTTPUpstream }

// WrapProviderHallUpstream decorates an HTTPUpstream so tracked requests count
// each real model submission exactly where it leaves the process.
func WrapProviderHallUpstream(base HTTPUpstream) HTTPUpstream {
	if base == nil {
		return nil
	}
	if _, already := base.(providerHallUpstream); already {
		return base
	}
	return providerHallUpstream{base: base}
}

func providerHallIsModelSubmission(req *http.Request) bool {
	if req == nil || req.URL == nil || req.Method != http.MethodPost {
		return false
	}
	path := strings.ToLower(strings.TrimRight(req.URL.Path, "/"))
	return strings.HasSuffix(path, "/responses") || strings.HasSuffix(path, "/chat/completions") || strings.HasSuffix(path, "/messages")
}

func (u providerHallUpstream) Do(req *http.Request, proxyURL string, accountID int64, accountConcurrency int) (*http.Response, error) {
	if providerHallIsModelSubmission(req) {
		ProviderHallTrackerFromContext(req.Context()).RecordSubmission()
	}
	return u.base.Do(req, proxyURL, accountID, accountConcurrency)
}

func (u providerHallUpstream) DoWithTLS(req *http.Request, proxyURL string, accountID int64, accountConcurrency int, profile *tlsfingerprint.Profile) (*http.Response, error) {
	if providerHallIsModelSubmission(req) {
		ProviderHallTrackerFromContext(req.Context()).RecordSubmission()
	}
	return u.base.DoWithTLS(req, proxyURL, accountID, accountConcurrency, profile)
}

// Attempt records one attempt's result. First content time is captured once;
// success is sticky. TTFT is measured from request arrival, so it includes
// local auth, queueing and retries, as the specification requires.
func (t *ProviderHallRequestTracker) Attempt(forwardStart time.Time, res *OpenAIForwardResult, err error) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.finished {
		return
	}
	if res != nil {
		if res.FirstTokenMs != nil && *res.FirstTokenMs >= 0 && t.firstContentAt == nil && !forwardStart.IsZero() {
			at := forwardStart.Add(time.Duration(*res.FirstTokenMs) * time.Millisecond)
			t.firstContentAt = &at
		}
		if model := strings.TrimSpace(res.UpstreamModel); model != "" {
			t.row.UpstreamModel = model
		}
		if model := strings.TrimSpace(res.UpstreamResponseModel); model != "" {
			t.row.ResponseModel = model
		}
		if res.ClientDisconnect {
			t.clientCancel = true
		}
		if err == nil || t.usage == nil {
			usage := res.Usage
			t.usage = &usage
		}
	}
	if err == nil {
		t.success = true
		return
	}
	t.failed = true
}

func (t *ProviderHallRequestTracker) Reject(reason ProviderHallExclusion) {
	if t == nil || reason == "" {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.finished && t.reject == "" {
		t.reject = reason
	}
}

// NoAccount records "no account available before any submission" (F=1, A=0).
func (t *ProviderHallRequestTracker) NoAccount() {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.finished {
		t.noAccount = true
	}
}

// Finish resolves the terminal state exactly once and hands the row to the
// collector. Deferred right after Begin.
func (t *ProviderHallRequestTracker) Finish() {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.finished {
		return
	}
	t.finished = true
	now := t.sink.clock()
	row := t.row
	row.EndedAt = &now
	row.Submissions = t.submissions
	if t.firstContentAt != nil {
		at := *t.firstContentAt
		row.FirstContentAt = &at
		if ms := int(at.Sub(row.StartedAt) / time.Millisecond); ms > 0 {
			row.TTFTMs = &ms
		}
	}
	if t.usage != nil {
		row.UsageKnown = true
		in, out := int64(t.usage.InputTokens), int64(t.usage.OutputTokens)
		read, create := int64(t.usage.CacheReadInputTokens), int64(t.usage.CacheCreationInputTokens)
		row.InputTokens, row.OutputTokens, row.CacheReadTokens, row.CacheCreation = &in, &out, &read, &create
	}
	clientGone := t.clientCancel || (t.ctx != nil && t.ctx.Err() != nil)
	switch {
	case row.ExclusionReason == ProviderHallExclusionTaskHeaderMismatch:
		row.Outcome = ProviderHallOutcomeExcluded
	case t.success:
		row.Outcome = ProviderHallOutcomeSuccess
		row.ExclusionReason = ProviderHallExclusionNone
	case clientGone:
		row.Outcome, row.ExclusionReason = ProviderHallOutcomeExcluded, ProviderHallExclusionClientCancel
	case t.reject != "":
		row.Outcome, row.ExclusionReason = ProviderHallOutcomeExcluded, t.reject
	case t.failed || t.noAccount:
		row.Outcome = ProviderHallOutcomeFailed
		row.ExclusionReason = ProviderHallExclusionNone
	default:
		row.Outcome, row.ExclusionReason = ProviderHallOutcomeExcluded, ProviderHallExclusionEarlyExit
	}
	t.sink.enqueue(providerHallEvent{kind: providerHallEventFinish, row: row})
}
