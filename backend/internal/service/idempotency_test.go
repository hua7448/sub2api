package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unicode/utf8"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

type inMemoryIdempotencyRepo struct {
	mu     sync.Mutex
	nextID int64
	data   map[string]*IdempotencyRecord
}

func newInMemoryIdempotencyRepo() *inMemoryIdempotencyRepo {
	return &inMemoryIdempotencyRepo{
		nextID: 1,
		data:   make(map[string]*IdempotencyRecord),
	}
}

func (r *inMemoryIdempotencyRepo) key(scope, hash string) string {
	return scope + "|" + hash
}

func cloneRecord(in *IdempotencyRecord) *IdempotencyRecord {
	if in == nil {
		return nil
	}
	out := *in
	if in.ResponseStatus != nil {
		v := *in.ResponseStatus
		out.ResponseStatus = &v
	}
	if in.ResponseBody != nil {
		v := *in.ResponseBody
		out.ResponseBody = &v
	}
	if in.ErrorReason != nil {
		v := *in.ErrorReason
		out.ErrorReason = &v
	}
	if in.LockedUntil != nil {
		v := *in.LockedUntil
		out.LockedUntil = &v
	}
	return &out
}

func (r *inMemoryIdempotencyRepo) CreateProcessing(_ context.Context, record *IdempotencyRecord) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := r.key(record.Scope, record.IdempotencyKeyHash)
	if _, ok := r.data[k]; ok {
		return false, nil
	}
	rec := cloneRecord(record)
	rec.ID = r.nextID
	rec.CreatedAt = time.Now()
	rec.UpdatedAt = rec.CreatedAt
	r.nextID++
	r.data[k] = rec
	record.ID = rec.ID
	record.CreatedAt = rec.CreatedAt
	record.UpdatedAt = rec.UpdatedAt
	return true, nil
}

func (r *inMemoryIdempotencyRepo) GetByScopeAndKeyHash(_ context.Context, scope, keyHash string) (*IdempotencyRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return cloneRecord(r.data[r.key(scope, keyHash)]), nil
}

func (r *inMemoryIdempotencyRepo) TryReclaim(_ context.Context, id int64, fromStatus string, now, newLockedUntil, newExpiresAt time.Time) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rec := range r.data {
		if rec.ID != id {
			continue
		}
		if rec.Status != fromStatus {
			return false, nil
		}
		if rec.LockedUntil != nil && rec.LockedUntil.After(now) {
			return false, nil
		}
		rec.Status = IdempotencyStatusProcessing
		rec.LockedUntil = &newLockedUntil
		rec.ExpiresAt = newExpiresAt
		rec.ErrorReason = nil
		rec.UpdatedAt = time.Now()
		return true, nil
	}
	return false, nil
}

func (r *inMemoryIdempotencyRepo) ExtendProcessingLock(_ context.Context, id int64, requestFingerprint string, newLockedUntil, newExpiresAt time.Time) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, rec := range r.data {
		if rec.ID != id {
			continue
		}
		if rec.Status != IdempotencyStatusProcessing || rec.RequestFingerprint != requestFingerprint {
			return false, nil
		}
		rec.LockedUntil = &newLockedUntil
		rec.ExpiresAt = newExpiresAt
		rec.UpdatedAt = time.Now()
		return true, nil
	}
	return false, nil
}

func (r *inMemoryIdempotencyRepo) MarkSucceeded(_ context.Context, id int64, responseStatus int, responseBody string, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rec := range r.data {
		if rec.ID != id {
			continue
		}
		rec.Status = IdempotencyStatusSucceeded
		rec.LockedUntil = nil
		rec.ExpiresAt = expiresAt
		rec.UpdatedAt = time.Now()
		rec.ErrorReason = nil
		rec.ResponseStatus = &responseStatus
		rec.ResponseBody = &responseBody
		return nil
	}
	return errors.New("record not found")
}

func (r *inMemoryIdempotencyRepo) MarkFailedRetryable(_ context.Context, id int64, errorReason string, lockedUntil, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rec := range r.data {
		if rec.ID != id {
			continue
		}
		rec.Status = IdempotencyStatusFailedRetryable
		rec.LockedUntil = &lockedUntil
		rec.ExpiresAt = expiresAt
		rec.UpdatedAt = time.Now()
		rec.ErrorReason = &errorReason
		return nil
	}
	return errors.New("record not found")
}

func (r *inMemoryIdempotencyRepo) DeleteExpired(_ context.Context, now time.Time, _ int) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var deleted int64
	for k, rec := range r.data {
		if !rec.ExpiresAt.After(now) {
			delete(r.data, k)
			deleted++
		}
	}
	return deleted, nil
}

func TestIdempotencyCoordinator_RequireKey(t *testing.T) {
	resetIdempotencyMetricsForTest()
	repo := newInMemoryIdempotencyRepo()
	cfg := DefaultIdempotencyConfig()
	cfg.ObserveOnly = false
	coordinator := NewIdempotencyCoordinator(repo, cfg)

	_, err := coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:      "test.scope",
		Method:     "POST",
		Route:      "/test",
		ActorScope: "admin:1",
		RequireKey: true,
		Payload:    map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(err), infraerrors.Code(ErrIdempotencyKeyRequired))
}

func TestIdempotencyCoordinator_ReplaySucceededResult(t *testing.T) {
	resetIdempotencyMetricsForTest()
	repo := newInMemoryIdempotencyRepo()
	cfg := DefaultIdempotencyConfig()
	coordinator := NewIdempotencyCoordinator(repo, cfg)

	execCount := 0
	exec := func(ctx context.Context) (any, error) {
		execCount++
		return map[string]any{"count": execCount}, nil
	}

	opts := IdempotencyExecuteOptions{
		Scope:          "test.scope",
		Method:         "POST",
		Route:          "/test",
		ActorScope:     "user:1",
		RequireKey:     true,
		IdempotencyKey: "case-1",
		Payload:        map[string]any{"a": 1},
	}

	first, err := coordinator.Execute(context.Background(), opts, exec)
	require.NoError(t, err)
	require.False(t, first.Replayed)

	second, err := coordinator.Execute(context.Background(), opts, exec)
	require.NoError(t, err)
	require.True(t, second.Replayed)
	require.Equal(t, 1, execCount, "second request should replay without executing business logic")

	metrics := GetIdempotencyMetricsSnapshot()
	require.Equal(t, uint64(1), metrics.ClaimTotal)
	require.Equal(t, uint64(1), metrics.ReplayTotal)
}

func TestIdempotencyCoordinator_ReclaimExpiredSucceededRecord(t *testing.T) {
	resetIdempotencyMetricsForTest()
	repo := newInMemoryIdempotencyRepo()
	coordinator := NewIdempotencyCoordinator(repo, DefaultIdempotencyConfig())

	opts := IdempotencyExecuteOptions{
		Scope:          "test.scope.expired",
		Method:         "POST",
		Route:          "/test/expired",
		ActorScope:     "user:99",
		RequireKey:     true,
		IdempotencyKey: "expired-case",
		Payload:        map[string]any{"k": "v"},
	}

	execCount := 0
	exec := func(ctx context.Context) (any, error) {
		execCount++
		return map[string]any{"count": execCount}, nil
	}

	first, err := coordinator.Execute(context.Background(), opts, exec)
	require.NoError(t, err)
	require.NotNil(t, first)
	require.False(t, first.Replayed)
	require.Equal(t, 1, execCount)

	keyHash := HashIdempotencyKey(opts.IdempotencyKey)
	repo.mu.Lock()
	existing := repo.data[repo.key(opts.Scope, keyHash)]
	require.NotNil(t, existing)
	existing.ExpiresAt = time.Now().Add(-time.Second)
	repo.mu.Unlock()

	second, err := coordinator.Execute(context.Background(), opts, exec)
	require.NoError(t, err)
	require.NotNil(t, second)
	require.False(t, second.Replayed, "expired record should be reclaimed and execute business logic again")
	require.Equal(t, 2, execCount)

	third, err := coordinator.Execute(context.Background(), opts, exec)
	require.NoError(t, err)
	require.NotNil(t, third)
	require.True(t, third.Replayed)
	payload, ok := third.Data.(map[string]any)
	require.True(t, ok)
	require.Equal(t, float64(2), payload["count"])

	metrics := GetIdempotencyMetricsSnapshot()
	require.GreaterOrEqual(t, metrics.ClaimTotal, uint64(2))
	require.GreaterOrEqual(t, metrics.ReplayTotal, uint64(1))
}

func TestIdempotencyCoordinator_SameKeyDifferentPayloadConflict(t *testing.T) {
	resetIdempotencyMetricsForTest()
	repo := newInMemoryIdempotencyRepo()
	cfg := DefaultIdempotencyConfig()
	coordinator := NewIdempotencyCoordinator(repo, cfg)

	_, err := coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "test.scope",
		Method:         "POST",
		Route:          "/test",
		ActorScope:     "user:1",
		RequireKey:     true,
		IdempotencyKey: "case-2",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.NoError(t, err)

	_, err = coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "test.scope",
		Method:         "POST",
		Route:          "/test",
		ActorScope:     "user:1",
		RequireKey:     true,
		IdempotencyKey: "case-2",
		Payload:        map[string]any{"a": 2},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(err), infraerrors.Code(ErrIdempotencyKeyConflict))

	metrics := GetIdempotencyMetricsSnapshot()
	require.Equal(t, uint64(1), metrics.ConflictTotal)
}

func TestIdempotencyCoordinator_BackoffAfterRetryableFailure(t *testing.T) {
	resetIdempotencyMetricsForTest()
	repo := newInMemoryIdempotencyRepo()
	cfg := DefaultIdempotencyConfig()
	cfg.FailedRetryBackoff = 2 * time.Second
	coordinator := NewIdempotencyCoordinator(repo, cfg)

	opts := IdempotencyExecuteOptions{
		Scope:          "test.scope",
		Method:         "POST",
		Route:          "/test",
		ActorScope:     "user:1",
		RequireKey:     true,
		IdempotencyKey: "case-3",
		Payload:        map[string]any{"a": 1},
	}

	_, err := coordinator.Execute(context.Background(), opts, func(ctx context.Context) (any, error) {
		return nil, infraerrors.InternalServer("UPSTREAM_ERROR", "upstream error")
	})
	require.Error(t, err)

	_, err = coordinator.Execute(context.Background(), opts, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(err), infraerrors.Code(ErrIdempotencyRetryBackoff))
	require.Greater(t, RetryAfterSecondsFromError(err), 0)

	metrics := GetIdempotencyMetricsSnapshot()
	require.GreaterOrEqual(t, metrics.RetryBackoffTotal, uint64(2))
	require.GreaterOrEqual(t, metrics.ConflictTotal, uint64(1))
	require.GreaterOrEqual(t, metrics.ProcessingDurationCount, uint64(1))
}

func TestIdempotencyCoordinator_ConcurrentSameKeySingleSideEffect(t *testing.T) {
	resetIdempotencyMetricsForTest()
	repo := newInMemoryIdempotencyRepo()
	cfg := DefaultIdempotencyConfig()
	cfg.ProcessingTimeout = 2 * time.Second
	coordinator := NewIdempotencyCoordinator(repo, cfg)

	opts := IdempotencyExecuteOptions{
		Scope:          "test.scope.concurrent",
		Method:         "POST",
		Route:          "/test/concurrent",
		ActorScope:     "user:7",
		RequireKey:     true,
		IdempotencyKey: "concurrent-case",
		Payload:        map[string]any{"v": 1},
	}

	var execCount int32
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = coordinator.Execute(context.Background(), opts, func(ctx context.Context) (any, error) {
				atomic.AddInt32(&execCount, 1)
				time.Sleep(80 * time.Millisecond)
				return map[string]any{"ok": true}, nil
			})
		}()
	}
	wg.Wait()

	replayed, err := coordinator.Execute(context.Background(), opts, func(ctx context.Context) (any, error) {
		atomic.AddInt32(&execCount, 1)
		return map[string]any{"ok": true}, nil
	})
	require.NoError(t, err)
	require.True(t, replayed.Replayed)
	require.Equal(t, int32(1), atomic.LoadInt32(&execCount), "concurrent same-key requests should execute business side-effect once")

	metrics := GetIdempotencyMetricsSnapshot()
	require.Equal(t, uint64(1), metrics.ClaimTotal)
	require.Equal(t, uint64(1), metrics.ReplayTotal)
	require.GreaterOrEqual(t, metrics.ConflictTotal, uint64(1))
}

func seedProcessingIdempotencyRecord(t *testing.T, repo *inMemoryIdempotencyRepo, opts IdempotencyExecuteOptions, lockedUntil *time.Time) {
	t.Helper()
	fingerprint, err := BuildIdempotencyFingerprint(opts.Method, opts.Route, opts.ActorScope, opts.Payload)
	require.NoError(t, err)
	record := &IdempotencyRecord{
		Scope:              opts.Scope,
		IdempotencyKeyHash: HashIdempotencyKey(opts.IdempotencyKey),
		RequestFingerprint: fingerprint,
		Status:             IdempotencyStatusProcessing,
		LockedUntil:        lockedUntil,
		ExpiresAt:          time.Now().Add(time.Hour),
	}
	owner, err := repo.CreateProcessing(context.Background(), record)
	require.NoError(t, err)
	require.True(t, owner)
}

func TestIdempotencyCoordinator_ReclaimsProcessingWithoutLock(t *testing.T) {
	repo := newInMemoryIdempotencyRepo()
	cfg := DefaultIdempotencyConfig()
	coordinator := NewIdempotencyCoordinator(repo, cfg)
	opts := IdempotencyExecuteOptions{
		Scope:          "test.scope.processing.nil-lock",
		Method:         "POST",
		Route:          "/test/processing/nil-lock",
		ActorScope:     "user:17",
		RequireKey:     true,
		IdempotencyKey: "processing-nil-lock",
		Payload:        map[string]any{"v": 1},
	}
	seedProcessingIdempotencyRecord(t, repo, opts, nil)

	var executions atomic.Int32
	result, err := coordinator.Execute(context.Background(), opts, func(context.Context) (any, error) {
		executions.Add(1)
		return map[string]any{"ok": true}, nil
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.Replayed)
	require.Equal(t, int32(1), executions.Load())
	stored, err := repo.GetByScopeAndKeyHash(context.Background(), opts.Scope, HashIdempotencyKey(opts.IdempotencyKey))
	require.NoError(t, err)
	require.Equal(t, IdempotencyStatusSucceeded, stored.Status)
}

func TestIdempotencyCoordinator_DoesNotReclaimActiveProcessing(t *testing.T) {
	repo := newInMemoryIdempotencyRepo()
	cfg := DefaultIdempotencyConfig()
	coordinator := NewIdempotencyCoordinator(repo, cfg)
	opts := IdempotencyExecuteOptions{
		Scope:          "test.scope.processing.active",
		Method:         "POST",
		Route:          "/test/processing/active",
		ActorScope:     "user:18",
		RequireKey:     true,
		IdempotencyKey: "processing-active",
		Payload:        map[string]any{"v": 1},
	}
	lockedUntil := time.Now().Add(time.Minute)
	seedProcessingIdempotencyRecord(t, repo, opts, &lockedUntil)

	var executions atomic.Int32
	_, err := coordinator.Execute(context.Background(), opts, func(context.Context) (any, error) {
		executions.Add(1)
		return map[string]any{"ok": true}, nil
	})

	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyInProgress), infraerrors.Code(err))
	require.Equal(t, int32(0), executions.Load())
	require.Greater(t, RetryAfterSecondsFromError(err), 0)
}

func TestIdempotencyCoordinator_ConcurrentExpiredProcessingSingleReclaimer(t *testing.T) {
	repo := newInMemoryIdempotencyRepo()
	cfg := DefaultIdempotencyConfig()
	cfg.ProcessingTimeout = 2 * time.Second
	coordinator := NewIdempotencyCoordinator(repo, cfg)
	opts := IdempotencyExecuteOptions{
		Scope:          "test.scope.processing.expired.concurrent",
		Method:         "POST",
		Route:          "/test/processing/expired/concurrent",
		ActorScope:     "user:19",
		RequireKey:     true,
		IdempotencyKey: "processing-expired-concurrent",
		Payload:        map[string]any{"v": 1},
	}
	lockedUntil := time.Now().Add(-time.Second)
	seedProcessingIdempotencyRecord(t, repo, opts, &lockedUntil)

	const callers = 8
	start := make(chan struct{})
	release := make(chan struct{})
	entered := make(chan struct{}, 1)
	results := make(chan error, callers)
	var executions atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := coordinator.Execute(context.Background(), opts, func(context.Context) (any, error) {
				if executions.Add(1) == 1 {
					entered <- struct{}{}
				}
				<-release
				return map[string]any{"ok": true}, nil
			})
			results <- err
		}()
	}
	close(start)
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for reclaimed executor")
	}
	time.Sleep(20 * time.Millisecond)
	close(release)
	wg.Wait()
	close(results)

	successes := 0
	for err := range results {
		if err == nil {
			successes++
			continue
		}
		require.Equal(t, infraerrors.Code(ErrIdempotencyInProgress), infraerrors.Code(err))
	}
	require.GreaterOrEqual(t, successes, 1)
	require.Equal(t, int32(1), executions.Load(), "only the atomic processing reclaimer may execute")
	stored, err := repo.GetByScopeAndKeyHash(context.Background(), opts.Scope, HashIdempotencyKey(opts.IdempotencyKey))
	require.NoError(t, err)
	require.Equal(t, IdempotencyStatusSucceeded, stored.Status)
}

type failingIdempotencyRepo struct{}

func (failingIdempotencyRepo) CreateProcessing(context.Context, *IdempotencyRecord) (bool, error) {
	return false, errors.New("store unavailable")
}
func (failingIdempotencyRepo) GetByScopeAndKeyHash(context.Context, string, string) (*IdempotencyRecord, error) {
	return nil, errors.New("store unavailable")
}
func (failingIdempotencyRepo) TryReclaim(context.Context, int64, string, time.Time, time.Time, time.Time) (bool, error) {
	return false, errors.New("store unavailable")
}
func (failingIdempotencyRepo) ExtendProcessingLock(context.Context, int64, string, time.Time, time.Time) (bool, error) {
	return false, errors.New("store unavailable")
}
func (failingIdempotencyRepo) MarkSucceeded(context.Context, int64, int, string, time.Time) error {
	return errors.New("store unavailable")
}
func (failingIdempotencyRepo) MarkFailedRetryable(context.Context, int64, string, time.Time, time.Time) error {
	return errors.New("store unavailable")
}
func (failingIdempotencyRepo) DeleteExpired(context.Context, time.Time, int) (int64, error) {
	return 0, errors.New("store unavailable")
}

func TestIdempotencyCoordinator_StoreUnavailableMetrics(t *testing.T) {
	resetIdempotencyMetricsForTest()
	coordinator := NewIdempotencyCoordinator(failingIdempotencyRepo{}, DefaultIdempotencyConfig())

	_, err := coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "test.scope.unavailable",
		Method:         "POST",
		Route:          "/test/unavailable",
		ActorScope:     "admin:1",
		RequireKey:     true,
		IdempotencyKey: "case-unavailable",
		Payload:        map[string]any{"v": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyStoreUnavail), infraerrors.Code(err))
	require.GreaterOrEqual(t, GetIdempotencyMetricsSnapshot().StoreUnavailableTotal, uint64(1))
}

type utf8RejectingIdempotencyRepo struct {
	inMemoryIdempotencyRepo
}

func newUTF8RejectingIdempotencyRepo() *utf8RejectingIdempotencyRepo {
	return &utf8RejectingIdempotencyRepo{inMemoryIdempotencyRepo: *newInMemoryIdempotencyRepo()}
}

func (r *utf8RejectingIdempotencyRepo) MarkSucceeded(ctx context.Context, id int64, responseStatus int, responseBody string, expiresAt time.Time) error {
	if !utf8.ValidString(responseBody) {
		return errors.New(`pq: invalid byte sequence for encoding "UTF8": 0xe8 0xb4 0x2e`)
	}
	return r.inMemoryIdempotencyRepo.MarkSucceeded(ctx, id, responseStatus, responseBody, expiresAt)
}

func TestIdempotencyCoordinator_TruncatedStoredResponseRemainsUTF8(t *testing.T) {
	repo := newUTF8RejectingIdempotencyRepo()
	cfg := DefaultIdempotencyConfig()
	cfg.MaxStoredResponseLen = len(`{"message":"`) + 2
	coordinator := NewIdempotencyCoordinator(repo, cfg)

	opts := IdempotencyExecuteOptions{
		Scope:          "test.scope.truncate_utf8",
		Method:         "POST",
		Route:          "/api/v1/accounts/import/codex-session",
		ActorScope:     "admin:1",
		RequireKey:     true,
		IdempotencyKey: "truncate-utf8",
		Payload:        map[string]any{"content": "codex-session"},
	}

	result, err := coordinator.Execute(context.Background(), opts, func(ctx context.Context) (any, error) {
		return map[string]any{"message": strings.Repeat("\u8d26", 8)}, nil
	})
	require.NoError(t, err)
	require.NotNil(t, result)

	stored, err := repo.GetByScopeAndKeyHash(context.Background(), opts.Scope, HashIdempotencyKey(opts.IdempotencyKey))
	require.NoError(t, err)
	require.NotNil(t, stored)
	require.NotNil(t, stored.ResponseBody)
	require.True(t, utf8.ValidString(*stored.ResponseBody))
	require.Contains(t, *stored.ResponseBody, "...(truncated)")
}

func TestDefaultIdempotencyCoordinatorAndTTLs(t *testing.T) {
	SetDefaultIdempotencyCoordinator(nil)
	require.Nil(t, DefaultIdempotencyCoordinator())
	require.Equal(t, DefaultIdempotencyConfig().DefaultTTL, DefaultWriteIdempotencyTTL())
	require.Equal(t, DefaultIdempotencyConfig().SystemOperationTTL, DefaultSystemOperationIdempotencyTTL())

	coordinator := NewIdempotencyCoordinator(newInMemoryIdempotencyRepo(), IdempotencyConfig{
		DefaultTTL:         2 * time.Hour,
		SystemOperationTTL: 15 * time.Minute,
		ProcessingTimeout:  10 * time.Second,
		FailedRetryBackoff: 3 * time.Second,
		ObserveOnly:        false,
	})
	SetDefaultIdempotencyCoordinator(coordinator)
	t.Cleanup(func() {
		SetDefaultIdempotencyCoordinator(nil)
	})

	require.Same(t, coordinator, DefaultIdempotencyCoordinator())
	require.Equal(t, 2*time.Hour, DefaultWriteIdempotencyTTL())
	require.Equal(t, 15*time.Minute, DefaultSystemOperationIdempotencyTTL())
}

func TestNormalizeIdempotencyKeyAndFingerprint(t *testing.T) {
	key, err := NormalizeIdempotencyKey("  abc-123  ")
	require.NoError(t, err)
	require.Equal(t, "abc-123", key)

	key, err = NormalizeIdempotencyKey("")
	require.NoError(t, err)
	require.Equal(t, "", key)

	_, err = NormalizeIdempotencyKey(string(make([]byte, 129)))
	require.Error(t, err)

	_, err = NormalizeIdempotencyKey("bad\nkey")
	require.Error(t, err)

	fp1, err := BuildIdempotencyFingerprint("", "", "", map[string]any{"a": 1})
	require.NoError(t, err)
	require.NotEmpty(t, fp1)
	fp2, err := BuildIdempotencyFingerprint("POST", "/", "anonymous", map[string]any{"a": 1})
	require.NoError(t, err)
	require.Equal(t, fp1, fp2)

	_, err = BuildIdempotencyFingerprint("POST", "/x", "u:1", map[string]any{"bad": make(chan int)})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyInvalidPayload), infraerrors.Code(err))
}

func TestRetryAfterSecondsFromErrorBranches(t *testing.T) {
	require.Equal(t, 0, RetryAfterSecondsFromError(nil))
	require.Equal(t, 0, RetryAfterSecondsFromError(errors.New("plain")))

	err := ErrIdempotencyInProgress.WithMetadata(map[string]string{"retry_after": "12"})
	require.Equal(t, 12, RetryAfterSecondsFromError(err))

	err = ErrIdempotencyInProgress.WithMetadata(map[string]string{"retry_after": "bad"})
	require.Equal(t, 0, RetryAfterSecondsFromError(err))
}

func TestIdempotencyCoordinator_ExecuteNilExecutorAndNoKeyPassThrough(t *testing.T) {
	repo := newInMemoryIdempotencyRepo()
	coordinator := NewIdempotencyCoordinator(repo, DefaultIdempotencyConfig())

	_, err := coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope",
		IdempotencyKey: "k",
		Payload:        map[string]any{"a": 1},
	}, nil)
	require.Error(t, err)
	require.Equal(t, "IDEMPOTENCY_EXECUTOR_NIL", infraerrors.Reason(err))

	called := 0
	result, err := coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:      "scope",
		RequireKey: true,
		Payload:    map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		called++
		return map[string]any{"ok": true}, nil
	})
	require.NoError(t, err)
	require.Equal(t, 1, called)
	require.NotNil(t, result)
	require.False(t, result.Replayed)
}

type noIDOwnerRepo struct{}

func (noIDOwnerRepo) CreateProcessing(context.Context, *IdempotencyRecord) (bool, error) {
	return true, nil
}
func (noIDOwnerRepo) GetByScopeAndKeyHash(context.Context, string, string) (*IdempotencyRecord, error) {
	return nil, nil
}
func (noIDOwnerRepo) TryReclaim(context.Context, int64, string, time.Time, time.Time, time.Time) (bool, error) {
	return false, nil
}
func (noIDOwnerRepo) ExtendProcessingLock(context.Context, int64, string, time.Time, time.Time) (bool, error) {
	return false, nil
}
func (noIDOwnerRepo) MarkSucceeded(context.Context, int64, int, string, time.Time) error { return nil }
func (noIDOwnerRepo) MarkFailedRetryable(context.Context, int64, string, time.Time, time.Time) error {
	return nil
}
func (noIDOwnerRepo) DeleteExpired(context.Context, time.Time, int) (int64, error) { return 0, nil }

func TestIdempotencyCoordinator_RepoNilScopeRequiredAndRecordIDMissing(t *testing.T) {
	cfg := DefaultIdempotencyConfig()
	coordinator := NewIdempotencyCoordinator(nil, cfg)

	_, err := coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope",
		IdempotencyKey: "k",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyStoreUnavail), infraerrors.Code(err))

	coordinator = NewIdempotencyCoordinator(newInMemoryIdempotencyRepo(), cfg)
	_, err = coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		IdempotencyKey: "k2",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, "IDEMPOTENCY_SCOPE_REQUIRED", infraerrors.Reason(err))

	coordinator = NewIdempotencyCoordinator(noIDOwnerRepo{}, cfg)
	_, err = coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope-no-id",
		IdempotencyKey: "k3",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyStoreUnavail), infraerrors.Code(err))
}

type conflictBranchRepo struct {
	existing      *IdempotencyRecord
	tryReclaimErr error
	tryReclaimOK  bool
}

func (r *conflictBranchRepo) CreateProcessing(context.Context, *IdempotencyRecord) (bool, error) {
	return false, nil
}
func (r *conflictBranchRepo) GetByScopeAndKeyHash(context.Context, string, string) (*IdempotencyRecord, error) {
	return cloneRecord(r.existing), nil
}
func (r *conflictBranchRepo) TryReclaim(context.Context, int64, string, time.Time, time.Time, time.Time) (bool, error) {
	if r.tryReclaimErr != nil {
		return false, r.tryReclaimErr
	}
	return r.tryReclaimOK, nil
}
func (r *conflictBranchRepo) ExtendProcessingLock(context.Context, int64, string, time.Time, time.Time) (bool, error) {
	return false, nil
}
func (r *conflictBranchRepo) MarkSucceeded(context.Context, int64, int, string, time.Time) error {
	return nil
}
func (r *conflictBranchRepo) MarkFailedRetryable(context.Context, int64, string, time.Time, time.Time) error {
	return nil
}
func (r *conflictBranchRepo) DeleteExpired(context.Context, time.Time, int) (int64, error) {
	return 0, nil
}

func TestIdempotencyCoordinator_ConflictBranchesAndDecodeError(t *testing.T) {
	now := time.Now()
	fp, err := BuildIdempotencyFingerprint("POST", "/x", "u:1", map[string]any{"a": 1})
	require.NoError(t, err)
	badBody := "{bad-json"
	repo := &conflictBranchRepo{
		existing: &IdempotencyRecord{
			ID:                 1,
			Scope:              "scope",
			IdempotencyKeyHash: HashIdempotencyKey("k"),
			RequestFingerprint: fp,
			Status:             IdempotencyStatusSucceeded,
			ResponseBody:       &badBody,
			ExpiresAt:          now.Add(time.Hour),
		},
	}
	coordinator := NewIdempotencyCoordinator(repo, DefaultIdempotencyConfig())
	_, err = coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope",
		IdempotencyKey: "k",
		Method:         "POST",
		Route:          "/x",
		ActorScope:     "u:1",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyStoreUnavail), infraerrors.Code(err))

	repo.existing = &IdempotencyRecord{
		ID:                 2,
		Scope:              "scope",
		IdempotencyKeyHash: HashIdempotencyKey("k"),
		RequestFingerprint: fp,
		Status:             "unknown",
		ExpiresAt:          now.Add(time.Hour),
	}
	_, err = coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope",
		IdempotencyKey: "k",
		Method:         "POST",
		Route:          "/x",
		ActorScope:     "u:1",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyKeyConflict), infraerrors.Code(err))

	repo.existing = &IdempotencyRecord{
		ID:                 3,
		Scope:              "scope",
		IdempotencyKeyHash: HashIdempotencyKey("k"),
		RequestFingerprint: fp,
		Status:             IdempotencyStatusFailedRetryable,
		LockedUntil:        ptrTime(now.Add(-time.Second)),
		ExpiresAt:          now.Add(time.Hour),
	}
	repo.tryReclaimErr = errors.New("reclaim down")
	_, err = coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope",
		IdempotencyKey: "k",
		Method:         "POST",
		Route:          "/x",
		ActorScope:     "u:1",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyStoreUnavail), infraerrors.Code(err))

	repo.tryReclaimErr = nil
	repo.tryReclaimOK = false
	_, err = coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope",
		IdempotencyKey: "k",
		Method:         "POST",
		Route:          "/x",
		ActorScope:     "u:1",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyInProgress), infraerrors.Code(err))
}

type markBehaviorRepo struct {
	inMemoryIdempotencyRepo
	failMarkSucceeded bool
	failMarkFailed    bool
}

type contextAwareMarkRepo struct {
	inMemoryIdempotencyRepo

	successFailures int32
	failedFailures  int32
	successCalls    atomic.Int32
	failedCalls     atomic.Int32
	sawCanceled     atomic.Bool
	missingDeadline atomic.Bool
}

func (r *contextAwareMarkRepo) inspectPersistenceContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		r.sawCanceled.Store(true)
		return err
	}
	if _, ok := ctx.Deadline(); !ok {
		r.missingDeadline.Store(true)
	}
	return nil
}

func (r *contextAwareMarkRepo) MarkSucceeded(ctx context.Context, id int64, responseStatus int, responseBody string, expiresAt time.Time) error {
	call := r.successCalls.Add(1)
	if err := r.inspectPersistenceContext(ctx); err != nil {
		return err
	}
	if call <= r.successFailures {
		return errors.New("transient mark succeeded failure")
	}
	return r.inMemoryIdempotencyRepo.MarkSucceeded(ctx, id, responseStatus, responseBody, expiresAt)
}

func (r *contextAwareMarkRepo) MarkFailedRetryable(ctx context.Context, id int64, errorReason string, lockedUntil, expiresAt time.Time) error {
	call := r.failedCalls.Add(1)
	if err := r.inspectPersistenceContext(ctx); err != nil {
		return err
	}
	if call <= r.failedFailures {
		return errors.New("transient mark failed retryable failure")
	}
	return r.inMemoryIdempotencyRepo.MarkFailedRetryable(ctx, id, errorReason, lockedUntil, expiresAt)
}

func TestIdempotencyCoordinator_MarkSucceededRetriesAfterRequestCancellation(t *testing.T) {
	repo := &contextAwareMarkRepo{
		inMemoryIdempotencyRepo: *newInMemoryIdempotencyRepo(),
		successFailures:         2,
	}
	coordinator := NewIdempotencyCoordinator(repo, DefaultIdempotencyConfig())
	opts := IdempotencyExecuteOptions{
		Scope:          "scope-canceled-success",
		IdempotencyKey: "canceled-success",
		Method:         "POST",
		Route:          "/canceled-success",
		ActorScope:     "u:21",
		Payload:        map[string]any{"a": 1},
	}
	requestCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := coordinator.Execute(requestCtx, opts, func(context.Context) (any, error) {
		cancel()
		return map[string]any{"ok": true}, nil
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, int32(3), repo.successCalls.Load())
	require.False(t, repo.sawCanceled.Load())
	require.False(t, repo.missingDeadline.Load())
	stored, err := repo.GetByScopeAndKeyHash(context.Background(), opts.Scope, HashIdempotencyKey(opts.IdempotencyKey))
	require.NoError(t, err)
	require.Equal(t, IdempotencyStatusSucceeded, stored.Status)
}

func TestIdempotencyCoordinator_MarkFailedRetriesAfterRequestCancellationWithShortTTL(t *testing.T) {
	repo := &contextAwareMarkRepo{
		inMemoryIdempotencyRepo: *newInMemoryIdempotencyRepo(),
		failedFailures:          1,
	}
	cfg := DefaultIdempotencyConfig()
	cfg.FailedRetryBackoff = 5 * time.Second
	coordinator := NewIdempotencyCoordinator(repo, cfg)
	opts := IdempotencyExecuteOptions{
		Scope:          "scope-canceled-failure",
		IdempotencyKey: "canceled-failure",
		Method:         "POST",
		Route:          "/canceled-failure",
		ActorScope:     "u:22",
		Payload:        map[string]any{"a": 1},
		TTL:            8 * 24 * time.Hour,
	}
	requestCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	businessErr := errors.New("business failure")
	startedAt := time.Now()

	_, err := coordinator.Execute(requestCtx, opts, func(context.Context) (any, error) {
		cancel()
		return nil, businessErr
	})

	require.ErrorIs(t, err, businessErr)
	require.Equal(t, int32(2), repo.failedCalls.Load())
	require.False(t, repo.sawCanceled.Load())
	require.False(t, repo.missingDeadline.Load())
	stored, err := repo.GetByScopeAndKeyHash(context.Background(), opts.Scope, HashIdempotencyKey(opts.IdempotencyKey))
	require.NoError(t, err)
	require.Equal(t, IdempotencyStatusFailedRetryable, stored.Status)
	require.NotNil(t, stored.LockedUntil)
	require.False(t, stored.ExpiresAt.Before(*stored.LockedUntil))
	require.WithinDuration(t, startedAt.Add(time.Minute), stored.ExpiresAt, 2*time.Second)
	require.True(t, stored.ExpiresAt.Before(startedAt.Add(2*time.Minute)), "failed record must not inherit the 8-day operation TTL")
}

func TestIdempotencyCoordinator_MarkSucceededRetryIsBounded(t *testing.T) {
	repo := &contextAwareMarkRepo{
		inMemoryIdempotencyRepo: *newInMemoryIdempotencyRepo(),
		successFailures:         100,
	}
	coordinator := NewIdempotencyCoordinator(repo, DefaultIdempotencyConfig())
	startedAt := time.Now()

	_, err := coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope-bounded-success",
		IdempotencyKey: "bounded-success",
		Method:         "POST",
		Route:          "/bounded-success",
		ActorScope:     "u:23",
		Payload:        map[string]any{"a": 1},
	}, func(context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})

	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyStoreUnavail), infraerrors.Code(err))
	require.Equal(t, int32(idempotencyFinalizeAttempts), repo.successCalls.Load())
	require.Less(t, time.Since(startedAt), time.Second)
	require.False(t, repo.missingDeadline.Load())
}

func TestIdempotencyFailedRecordExpiresAtBounds(t *testing.T) {
	now := time.Date(2026, 8, 5, 12, 50, 0, 0, time.UTC)

	require.Equal(
		t,
		now.Add(time.Minute),
		idempotencyFailedRecordExpiresAt(now, now.Add(8*24*time.Hour), now.Add(5*time.Second)),
	)
	require.Equal(
		t,
		now.Add(20*time.Second),
		idempotencyFailedRecordExpiresAt(now, now.Add(20*time.Second), now.Add(5*time.Second)),
	)
	require.Equal(
		t,
		now.Add(2*time.Minute),
		idempotencyFailedRecordExpiresAt(now, now.Add(20*time.Second), now.Add(2*time.Minute)),
	)
}

func (r *markBehaviorRepo) MarkSucceeded(ctx context.Context, id int64, responseStatus int, responseBody string, expiresAt time.Time) error {
	if r.failMarkSucceeded {
		return errors.New("mark succeeded failed")
	}
	return r.inMemoryIdempotencyRepo.MarkSucceeded(ctx, id, responseStatus, responseBody, expiresAt)
}

func (r *markBehaviorRepo) MarkFailedRetryable(ctx context.Context, id int64, errorReason string, lockedUntil, expiresAt time.Time) error {
	if r.failMarkFailed {
		return errors.New("mark failed retryable failed")
	}
	return r.inMemoryIdempotencyRepo.MarkFailedRetryable(ctx, id, errorReason, lockedUntil, expiresAt)
}

func TestIdempotencyCoordinator_MarkAndMarshalBranches(t *testing.T) {
	repo := &markBehaviorRepo{inMemoryIdempotencyRepo: *newInMemoryIdempotencyRepo()}
	coordinator := NewIdempotencyCoordinator(repo, DefaultIdempotencyConfig())

	repo.failMarkSucceeded = true
	_, err := coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope-success",
		IdempotencyKey: "k1",
		Method:         "POST",
		Route:          "/ok",
		ActorScope:     "u:1",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"ok": true}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyStoreUnavail), infraerrors.Code(err))

	repo.failMarkSucceeded = false
	_, err = coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope-marshal",
		IdempotencyKey: "k2",
		Method:         "POST",
		Route:          "/bad",
		ActorScope:     "u:1",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return map[string]any{"bad": make(chan int)}, nil
	})
	require.Error(t, err)
	require.Equal(t, infraerrors.Code(ErrIdempotencyStoreUnavail), infraerrors.Code(err))

	repo.failMarkFailed = true
	_, err = coordinator.Execute(context.Background(), IdempotencyExecuteOptions{
		Scope:          "scope-fail",
		IdempotencyKey: "k3",
		Method:         "POST",
		Route:          "/fail",
		ActorScope:     "u:1",
		Payload:        map[string]any{"a": 1},
	}, func(ctx context.Context) (any, error) {
		return nil, errors.New("plain failure")
	})
	require.Error(t, err)
	require.Equal(t, "plain failure", err.Error())
}

func TestIdempotencyCoordinator_HelperBranches(t *testing.T) {
	c := NewIdempotencyCoordinator(newInMemoryIdempotencyRepo(), IdempotencyConfig{
		DefaultTTL:           time.Hour,
		SystemOperationTTL:   time.Hour,
		ProcessingTimeout:    time.Second,
		FailedRetryBackoff:   time.Second,
		MaxStoredResponseLen: 12,
		ObserveOnly:          false,
	})

	// conflictWithRetryAfter without locked_until should return base error.
	base := ErrIdempotencyInProgress
	err := c.conflictWithRetryAfter(base, nil, time.Now())
	require.Equal(t, infraerrors.Code(base), infraerrors.Code(err))

	// marshalStoredResponse should truncate.
	body, err := c.marshalStoredResponse(map[string]any{"long": "abcdefghijklmnopqrstuvwxyz"})
	require.NoError(t, err)
	require.Contains(t, body, "...(truncated)")

	// decodeStoredResponse empty and invalid json.
	out, err := c.decodeStoredResponse(nil)
	require.NoError(t, err)
	_, ok := out.(map[string]any)
	require.True(t, ok)

	invalid := "{invalid"
	_, err = c.decodeStoredResponse(&invalid)
	require.Error(t, err)
}
