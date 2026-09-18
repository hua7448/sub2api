//go:build perf

package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// TestProviderHallGatewayOverhead compares the Responses handler latency with
// the provider hall collector absent versus enabled (target index loaded,
// background writer draining into a no-op fact repository). Plan matrix H:
// the P95 increase must stay within 5 ms.
//
//	go test -tags=perf ./internal/handler -run ProviderHallGatewayOverhead -v

const (
	overheadWarmup   = 200
	overheadRequests = 2000
	overheadGroupID  = int64(7401)
	overheadKeyID    = int64(8401)
)

type overheadFactRepo struct {
	mu      sync.Mutex
	batches int
	rows    int
	targets []service.ProviderHallTargetRef
}

func (f *overheadFactRepo) RegisterEpoch(context.Context, string, string, int, time.Time) (int64, error) {
	return 1, nil
}
func (f *overheadFactRepo) Heartbeat(context.Context, int64, time.Time) error { return nil }
func (f *overheadFactRepo) MarkEpochExited(context.Context, int64, string, time.Time) error {
	return nil
}
func (f *overheadFactRepo) WriteBatch(_ context.Context, b service.ProviderHallFactBatch) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.batches++
	f.rows += len(b.Starts) + len(b.Finishes) + len(b.Billings)
	return nil
}
func (f *overheadFactRepo) OpenGap(context.Context, service.ProviderHallGap) (int64, error) {
	return 1, nil
}
func (f *overheadFactRepo) CloseGap(context.Context, int64, time.Time) error { return nil }
func (f *overheadFactRepo) ListEnabledTargets(context.Context) ([]service.ProviderHallTargetRef, error) {
	return f.targets, nil
}
func (f *overheadFactRepo) ListProbeKeys(context.Context) (map[int64]time.Time, error) {
	return map[int64]time.Time{}, nil
}
func (f *overheadFactRepo) stats() (int, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.batches, f.rows
}

type overheadCfgRepo struct {
	service.ProviderHallRepository
}

func (overheadCfgRepo) GetConfig(context.Context) (*service.ProviderHallConfig, error) {
	return &service.ProviderHallConfig{CollectionEnabled: true, DefaultProtocol: "responses", DefaultRange: "6h"}, nil
}

type overheadAccountRepo struct {
	service.AccountRepository
	accounts []service.Account
}

func (s *overheadAccountRepo) ListSchedulableByPlatform(_ context.Context, platform string) ([]service.Account, error) {
	out := make([]service.Account, 0, len(s.accounts))
	for _, account := range s.accounts {
		if account.Platform == platform && account.IsSchedulable() {
			out = append(out, account)
		}
	}
	return out, nil
}
func (s *overheadAccountRepo) ListSchedulableByGroupIDAndPlatform(ctx context.Context, _ int64, platform string) ([]service.Account, error) {
	return s.ListSchedulableByPlatform(ctx, platform)
}
func (s *overheadAccountRepo) ListSchedulableUngroupedByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	return s.ListSchedulableByPlatform(ctx, platform)
}
func (s *overheadAccountRepo) GetByID(_ context.Context, id int64) (*service.Account, error) {
	for _, account := range s.accounts {
		if account.ID == id {
			acc := account
			return &acc, nil
		}
	}
	return nil, nil
}
func (s *overheadAccountRepo) SetRateLimited(context.Context, int64, time.Time) error { return nil }

const overheadResponsesOK = `{"id":"resp_ok","object":"response","model":"gpt-hall","status":"completed","output":[{"type":"message","id":"msg_1","status":"completed","role":"assistant","content":[{"type":"output_text","text":"OK","annotations":[]}]}],"usage":{"input_tokens":11,"output_tokens":2,"total_tokens":13}}`

type overheadUpstream struct{}

func (overheadUpstream) respond(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		_, _ = io.Copy(io.Discard, req.Body)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(overheadResponsesOK)),
	}, nil
}
func (u overheadUpstream) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	return u.respond(req)
}
func (u overheadUpstream) DoWithTLS(req *http.Request, _ string, _ int64, _ int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.respond(req)
}

func overheadPercentile(d []time.Duration, p float64) time.Duration {
	c := append([]time.Duration(nil), d...)
	sort.Slice(c, func(i, j int) bool { return c[i] < c[j] })
	idx := int(float64(len(c)-1)*p + 0.5)
	return c[idx]
}

func newOverheadHandler(t *testing.T, withCollector bool) (*OpenAIGatewayHandler, *service.ProviderHallCollector, *overheadFactRepo) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{RunMode: config.RunModeSimple}
	cfg.Default.RateMultiplier = 1
	cfg.Security.URLAllowlist.Enabled = false
	cfg.Gateway.MaxAccountSwitches = 3
	accounts := []service.Account{{
		ID: 9901, Name: "acct", Platform: service.PlatformOpenAI,
		Type: service.AccountTypeAPIKey, Status: service.StatusActive, Schedulable: true, Priority: 1,
		Credentials: map[string]any{"api_key": "sk-test", "base_url": "https://api.example.test"},
		Extra:       map[string]any{"openai_passthrough": true},
	}}
	billingCacheSvc := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	t.Cleanup(billingCacheSvc.Stop)
	gatewaySvc := service.NewOpenAIGatewayService(
		&overheadAccountRepo{accounts: accounts},
		nil, nil, nil, nil, nil, nil, cfg, nil, nil,
		service.NewBillingService(cfg, nil), nil, billingCacheSvc, overheadUpstream{}, &service.DeferredService{},
		nil, nil, nil, nil, nil, nil, nil,
	)
	h := NewOpenAIGatewayHandler(gatewaySvc, service.NewConcurrencyService(nil), billingCacheSvc,
		service.NewAPIKeyService(nil, nil, nil, nil, nil, nil, cfg), nil, nil, nil, nil, cfg)
	if !withCollector {
		return h, nil, nil
	}
	facts := &overheadFactRepo{targets: []service.ProviderHallTargetRef{
		{TargetID: 1, GroupID: overheadGroupID, ProfileID: 100, Model: "gpt-hall", Protocol: "responses", Enabled: true},
	}}
	collector := service.NewProviderHallCollector(facts, overheadCfgRepo{}, "node-perf", "perf")
	collector.Start()
	t.Cleanup(collector.Stop)
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		collector.RefreshNow()
		if tr, _ := collector.Begin(context.Background(), service.ProviderHallBeginInput{APIKeyID: overheadKeyID, GroupID: overheadGroupID, Protocol: "responses", RequestedModel: "gpt-hall"}); tr != nil {
			tr.Finish()
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	gatewaySvc.SetProviderHallBillingSink(collector)
	h.SetProviderHallCollector(collector)
	return h, collector, facts
}

func overheadRequest(t *testing.T, h *OpenAIGatewayHandler) time.Duration {
	t.Helper()
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", strings.NewReader(`{"model":"gpt-hall","input":"hi","stream":false}`))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	groupID := overheadGroupID
	c.Set(string(middleware.ContextKeyAPIKey), &service.APIKey{
		ID: overheadKeyID, GroupID: &groupID,
		User:  &service.User{ID: 1703, Status: service.StatusActive},
		Group: &service.Group{ID: groupID, Platform: service.PlatformOpenAI, Status: service.StatusActive, AllowMessagesDispatch: true},
	})
	c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: 1703, Concurrency: 0})
	t0 := time.Now()
	h.Responses(c)
	d := time.Since(t0)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	return d
}

func TestProviderHallGatewayOverhead(t *testing.T) {
	measure := func(withCollector bool) ([]time.Duration, *service.ProviderHallCollector, *overheadFactRepo) {
		h, collector, facts := newOverheadHandler(t, withCollector)
		for i := 0; i < overheadWarmup; i++ {
			overheadRequest(t, h)
		}
		samples := make([]time.Duration, 0, overheadRequests)
		for i := 0; i < overheadRequests; i++ {
			samples = append(samples, overheadRequest(t, h))
		}
		return samples, collector, facts
	}
	// Interleave two rounds each so scheduler noise hits both variants alike.
	var off, on []time.Duration
	var facts *overheadFactRepo
	var collector *service.ProviderHallCollector
	for round := 0; round < 2; round++ {
		s, _, _ := measure(false)
		off = append(off, s...)
		s, collector, facts = measure(true)
		on = append(on, s...)
	}
	offP95, onP95 := overheadPercentile(off, 0.95), overheadPercentile(on, 0.95)
	offP50, onP50 := overheadPercentile(off, 0.50), overheadPercentile(on, 0.50)
	// Let the writer drain so the "enabled" path is proven to have produced facts.
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if _, rows := facts.stats(); rows >= overheadRequests {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	batches, rows := facts.stats()
	depth, capacity, dropped, overflowed := collector.QueueStats()
	t.Logf("requests per variant: %d (after %d warm-up); collector off P50=%s P95=%s; collector on P50=%s P95=%s; delta P95=%s",
		len(off), overheadWarmup, offP50.Round(time.Microsecond), offP95.Round(time.Microsecond), onP50.Round(time.Microsecond), onP95.Round(time.Microsecond), (onP95 - offP95).Round(time.Microsecond))
	t.Logf("collector wrote %d rows in %d batches; queue depth=%d/%d dropped=%d overflowed=%v", rows, batches, depth, capacity, dropped, overflowed)
	require.Greater(t, rows, 0, "enabled variant must actually track requests")
	require.Zero(t, dropped)
	require.LessOrEqual(t, onP95-offP95, 5*time.Millisecond, "collector adds at most 5 ms at P95")
}
