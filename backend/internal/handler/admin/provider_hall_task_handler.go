package admin

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// Small runtime surfaces so tests can substitute fakes. The concrete runner,
// collector and aggregator satisfy them and are nil-safe.
type providerHallJobRunner interface {
	EnqueueManual(ctx context.Context, kind service.ProviderHallJobKind, groupID, profileID int64, idempotencyKey string, actorID int64) (*service.ProviderHallJob, bool, error)
	CancelJob(ctx context.Context, jobID int64) error
}

type providerHallQueueSource interface {
	QueueStats() (depth, capacity int, dropped uint64, overflowed bool)
}

type providerHallAggregatorSource interface {
	Status() (lastRun time.Time, lastErr string)
}

// providerHallJobReads is the subset of the job repository the admin API uses.
type providerHallJobReads interface {
	ListJobs(ctx context.Context, f service.ProviderHallJobFilter, page, pageSize int) ([]service.ProviderHallJob, int, error)
	GetJob(ctx context.Context, id int64) (*service.ProviderHallJob, error)
	ListSamples(ctx context.Context, jobID int64) ([]service.ProviderHallSample, error)
	GetVerification(ctx context.Context, jobID int64) (*service.ProviderHallVerification, error)
	CountJobs(ctx context.Context, now time.Time) (service.ProviderHallJobCounts, error)
	SumSpend(ctx context.Context, day string) (service.ProviderHallSpendSummary, error)
	ListVerifications(ctx context.Context, groupID int64, profileID *int64, page, pageSize int) ([]service.ProviderHallVerification, int, error)
}

func providerHallQueryInt(c *gin.Context, name string, def, min, max int) (int, bool) {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return def, true
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < min || v > max {
		response.BadRequest(c, "invalid "+name)
		return 0, false
	}
	return v, true
}

func (h *ProviderHallHandler) enqueue(c *gin.Context, kind service.ProviderHallJobKind) {
	actor, ok := providerHallAdmin(c)
	if !ok {
		return
	}
	groupID, ok := providerHallID(c)
	if !ok {
		return
	}
	var input dto.ProviderHallEnqueueInput
	if !providerHallBody(c, &input) {
		return
	}
	if input.ProfileID < 1 {
		response.BadRequest(c, "profile_id is required")
		return
	}
	key := strings.TrimSpace(input.IdempotencyKey)
	if len(key) > 128 || strings.ContainsAny(key, "\r\n\t ") {
		response.BadRequest(c, "idempotency_key must be at most 128 characters without whitespace")
		return
	}
	if h.runner == nil {
		response.ErrorFrom(c, service.ErrProviderHallNotReady.WithMetadata(map[string]string{"switch": "tasks_enabled"}))
		return
	}
	job, reused, err := h.runner.EnqueueManual(c.Request.Context(), kind, groupID, input.ProfileID, key, actor)
	if response.ErrorFrom(c, err) {
		return
	}
	response.Accepted(c, dto.ProviderHallEnqueueResponse{JobID: job.ID, Status: string(job.Status), Reused: reused})
}

func (h *ProviderHallHandler) EnqueueProbe(c *gin.Context) {
	h.enqueue(c, service.ProviderHallJobProbe)
}
func (h *ProviderHallHandler) EnqueueVerification(c *gin.Context) {
	h.enqueue(c, service.ProviderHallJobVerification)
}

func (h *ProviderHallHandler) ListJobs(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	filter := service.ProviderHallJobFilter{}
	switch status := service.ProviderHallJobStatus(strings.TrimSpace(c.Query("status"))); status {
	case "", service.ProviderHallJobQueued, service.ProviderHallJobRunning, service.ProviderHallJobSucceeded, service.ProviderHallJobFailed, service.ProviderHallJobCancelled, service.ProviderHallJobUnknown:
		filter.Status = status
	default:
		response.BadRequest(c, "invalid status")
		return
	}
	switch kind := service.ProviderHallJobKind(strings.TrimSpace(c.Query("kind"))); kind {
	case "", service.ProviderHallJobProbe, service.ProviderHallJobVerification:
		filter.Kind = kind
	default:
		response.BadRequest(c, "invalid kind")
		return
	}
	if raw := strings.TrimSpace(c.Query("group_id")); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || id < 1 {
			response.BadRequest(c, "invalid group_id")
			return
		}
		filter.GroupID = id
	}
	page, ok := providerHallQueryInt(c, "page", 1, 1, 1_000_000)
	if !ok {
		return
	}
	pageSize, ok := providerHallQueryInt(c, "page_size", 20, 1, 100)
	if !ok {
		return
	}
	if h.jobs == nil {
		response.Paginated(c, []service.ProviderHallJob{}, 0, page, pageSize)
		return
	}
	items, total, err := h.jobs.ListJobs(c.Request.Context(), filter, page, pageSize)
	if response.ErrorFrom(c, err) {
		return
	}
	if items == nil {
		items = []service.ProviderHallJob{}
	}
	response.Paginated(c, items, int64(total), page, pageSize)
}

func (h *ProviderHallHandler) GetJob(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	if h.jobs == nil {
		response.ErrorFrom(c, service.ErrProviderHallNotFound)
		return
	}
	ctx := c.Request.Context()
	job, err := h.jobs.GetJob(ctx, id)
	if response.ErrorFrom(c, err) {
		return
	}
	samples, err := h.jobs.ListSamples(ctx, id)
	if response.ErrorFrom(c, err) {
		return
	}
	if samples == nil {
		samples = []service.ProviderHallSample{}
	}
	var report *service.ProviderHallVerification
	if job.Kind == service.ProviderHallJobVerification {
		if report, err = h.jobs.GetVerification(ctx, id); response.ErrorFrom(c, err) {
			return
		}
	}
	response.Success(c, gin.H{"job": job, "samples": samples, "verification": report})
}

func (h *ProviderHallHandler) CancelJob(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	id, ok := providerHallID(c)
	if !ok {
		return
	}
	if h.runner == nil {
		response.ErrorFrom(c, service.ErrProviderHallNotReady.WithMetadata(map[string]string{"switch": "tasks_enabled"}))
		return
	}
	if response.ErrorFrom(c, h.runner.CancelJob(c.Request.Context(), id)) {
		return
	}
	response.Success(c, gin.H{"job_id": id, "status": string(service.ProviderHallJobCancelled)})
}

// ListGroupVerifications returns reports for any group, listed or not.
func (h *ProviderHallHandler) ListGroupVerifications(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	groupID, ok := providerHallID(c)
	if !ok {
		return
	}
	var profileID *int64
	if raw := strings.TrimSpace(c.Query("profile_id")); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || id < 1 {
			response.BadRequest(c, "invalid profile_id")
			return
		}
		profileID = &id
	}
	page, ok := providerHallQueryInt(c, "page", 1, 1, 1_000_000)
	if !ok {
		return
	}
	pageSize, ok := providerHallQueryInt(c, "page_size", 20, 1, 50)
	if !ok {
		return
	}
	if h.jobs == nil {
		response.Paginated(c, []service.ProviderHallVerification{}, 0, page, pageSize)
		return
	}
	items, total, err := h.jobs.ListVerifications(c.Request.Context(), groupID, profileID, page, pageSize)
	if response.ErrorFrom(c, err) {
		return
	}
	if items == nil {
		items = []service.ProviderHallVerification{}
	}
	response.Paginated(c, items, int64(total), page, pageSize)
}

// Health assembles the operational view. Each section degrades independently:
// a failing read leaves its section empty and is reported in last_error-style
// fields rather than failing the whole endpoint.
func (h *ProviderHallHandler) Health(c *gin.Context) {
	if _, ok := providerHallAdmin(c); !ok {
		return
	}
	ctx := c.Request.Context()
	now := h.now().UTC()
	out := dto.ProviderHallAdminHealth{GeneratedAt: now}
	out.Collection.Nodes = []dto.ProviderHallHealthNode{}
	out.Collection.MissingExpected = []string{}
	out.Collection.OpenGaps = []dto.ProviderHallHealthGap{}

	cfg, err := h.service.GetConfig(ctx)
	if response.ErrorFrom(c, err) {
		return
	}
	out.Collection.Enabled = cfg.CollectionEnabled
	out.Budget.Day = service.ProviderHallBudgetDay(now)
	out.Budget.Budget = cfg.DailyBudget
	out.Budget.ConfirmedSpend, out.Budget.UncertainSpend = "0", "0"

	if h.queue != nil {
		depth, capacity, dropped, _ := h.queue.QueueStats()
		out.Collection.LocalQueue = dto.ProviderHallHealthQueue{Depth: depth, Capacity: capacity, Dropped: dropped}
	}
	alive := map[string]bool{}
	if h.admin != nil {
		epochs, err := h.admin.ListLatestEpochs(ctx)
		if response.ErrorFrom(c, err) {
			return
		}
		for _, e := range epochs {
			lost := e.ExitReason == "lost" || (e.ExitedAt == nil && e.HeartbeatAt.Before(now.Add(-service.ProviderHallNodeLostAfter)))
			if e.ExitedAt == nil && !lost {
				alive[e.NodeID] = true
			}
			out.Collection.Nodes = append(out.Collection.Nodes, dto.ProviderHallHealthNode{
				NodeID: e.NodeID, EpochID: e.EpochID, Version: e.BuildVersion, HeartbeatAt: e.HeartbeatAt, ConfirmedAt: e.ConfirmedAt,
				PersistedSeq: e.PersistedSeq, Overflowed: e.Overflowed, Lost: lost,
			})
		}
		gaps, err := h.admin.ListOpenGaps(ctx)
		if response.ErrorFrom(c, err) {
			return
		}
		for _, g := range gaps {
			out.Collection.OpenGaps = append(out.Collection.OpenGaps, dto.ProviderHallHealthGap{ID: g.ID, NodeID: g.NodeID, EpochID: g.EpochID, Scope: g.Scope, StartedAt: g.StartedAt, Reason: g.Reason})
		}
		state, err := h.admin.GetAggregatorState(ctx)
		if response.ErrorFrom(c, err) {
			return
		}
		out.Aggregator.Watermark = state.LastWindowEnd
		if state.LastWindowEnd != nil {
			lag := int64(now.Sub(*state.LastWindowEnd) / time.Second)
			out.Aggregator.LagSeconds = &lag
		}
		out.Aggregator.LastRunAt = state.LastRunAt
		out.Aggregator.LastError = state.LastError
		if out.Aggregator.DirtyCount, err = h.admin.CountDirty(ctx); response.ErrorFrom(c, err) {
			return
		}
		rec, err := h.admin.CountReconciliation(ctx, now)
		if response.ErrorFrom(c, err) {
			return
		}
		out.Reconciliation = dto.ProviderHallHealthReconciliation{Pending: rec.Pending, Uncertain: rec.Uncertain, Failed24h: rec.Failed24h}
	}
	if h.aggregator != nil {
		// In-process status wins when it is fresher than the durable state
		// (for example a lock miss recorded locally but never persisted).
		if lastRun, lastErr := h.aggregator.Status(); !lastRun.IsZero() && (out.Aggregator.LastRunAt == nil || lastRun.After(*out.Aggregator.LastRunAt)) {
			t := lastRun.UTC()
			out.Aggregator.LastRunAt = &t
			out.Aggregator.LastError = lastErr
		}
	}
	for _, node := range cfg.ExpectedNodes {
		if !alive[node] {
			out.Collection.MissingExpected = append(out.Collection.MissingExpected, node)
		}
	}
	if h.jobs != nil {
		spend, err := h.jobs.SumSpend(ctx, out.Budget.Day)
		if response.ErrorFrom(c, err) {
			return
		}
		out.Budget.ConfirmedSpend = spend.Confirmed.String()
		out.Budget.UncertainSpend = spend.Uncertain.String()
		out.Budget.InFlight = spend.InFlight
		switch {
		case !cfg.TasksEnabled:
			out.Budget.PausedReason = service.ProviderHallJobCodeTasksDisabled
		case service.ProviderHallBudgetExhausted(cfg.DailyBudget, spend.Confirmed):
			out.Budget.PausedReason = service.ProviderHallJobCodeBudgetExhausted
		}
		counts, err := h.jobs.CountJobs(ctx, now)
		if response.ErrorFrom(c, err) {
			return
		}
		out.Jobs = dto.ProviderHallHealthJobs{Queued: counts.Queued, Running: counts.Running, Unknown: counts.Unknown, Failed24h: counts.Failed24h}
	} else if !cfg.TasksEnabled {
		out.Budget.PausedReason = service.ProviderHallJobCodeTasksDisabled
	}
	response.Success(c, out)
}
