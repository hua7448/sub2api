package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/apikey"
	"github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/ent/providerhallconfig"
	"github.com/Wei-Shaw/sub2api/ent/providerhallgroup"
	"github.com/Wei-Shaw/sub2api/ent/providerhallprofile"
	"github.com/Wei-Shaw/sub2api/ent/providerhalltarget"
	"github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/ent/usersubscription"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

// providerHallJobRepository owns jobs, samples, reports and spend with raw
// SQL; the dispatch recheck reads configuration rows through Ent inside the
// same transaction.
type providerHallJobRepository struct{ db *sql.DB }

func NewProviderHallJobRepository(db *sql.DB) service.ProviderHallJobRepository {
	return &providerHallJobRepository{db: db}
}

const providerHallJobColumns = `id, kind, target_id, group_id, profile_id, config_snapshot, slot_at, not_before, idempotency_key, requested_by,
status, lease_owner, lease_until, attempts, budget_day, error_code, error_message, created_at, started_at, finished_at, updated_at`

type providerHallScanner interface {
	Scan(dest ...any) error
}

func providerHallScanJob(row providerHallScanner) (*service.ProviderHallJob, error) {
	var j service.ProviderHallJob
	var snapshot []byte
	var kind, status string
	var slotAt, notBefore, leaseUntil, startedAt, finishedAt sql.NullTime
	var idem, leaseOwner sql.NullString
	var requestedBy sql.NullInt64
	var budgetDay sql.NullString
	if err := row.Scan(&j.ID, &kind, &j.TargetID, &j.GroupID, &j.ProfileID, &snapshot, &slotAt, &notBefore, &idem, &requestedBy,
		&status, &leaseOwner, &leaseUntil, &j.Attempts, &budgetDay, &j.ErrorCode, &j.ErrorMessage, &j.CreatedAt, &startedAt, &finishedAt, &j.UpdatedAt); err != nil {
		return nil, err
	}
	j.Kind, j.Status = service.ProviderHallJobKind(kind), service.ProviderHallJobStatus(status)
	if len(snapshot) > 0 {
		_ = json.Unmarshal(snapshot, &j.Snapshot)
	}
	if j.Snapshot.Profile.ModelAliases == nil {
		j.Snapshot.Profile.ModelAliases = []string{}
	}
	nt := func(v sql.NullTime) *time.Time {
		if !v.Valid {
			return nil
		}
		t := v.Time.UTC()
		return &t
	}
	j.SlotAt, j.NotBefore, j.LeaseUntil, j.StartedAt, j.FinishedAt = nt(slotAt), nt(notBefore), nt(leaseUntil), nt(startedAt), nt(finishedAt)
	j.CreatedAt, j.UpdatedAt = j.CreatedAt.UTC(), j.UpdatedAt.UTC()
	if idem.Valid {
		j.IdempotencyKey = &idem.String
	}
	if leaseOwner.Valid {
		j.LeaseOwner = &leaseOwner.String
	}
	if requestedBy.Valid {
		j.RequestedBy = &requestedBy.Int64
	}
	if budgetDay.Valid {
		day := budgetDay.String
		if len(day) > 10 {
			day = day[:10]
		}
		j.BudgetDay = &day
	}
	return &j, nil
}

const providerHallSampleColumns = `id, job_id, test_id, seq, trace_id, status, prepared_at, dispatched_at, received_at, client_request_id, http_status,
ttft_ms, total_ms, generation_ms, input_tokens, output_tokens, response_model, result, error_code, detail, billing_status, actual_cost`

func providerHallScanSample(row providerHallScanner) (*service.ProviderHallSample, error) {
	var s service.ProviderHallSample
	var status, traceID string
	var dispatchedAt, receivedAt sql.NullTime
	var clientRequestID, responseModel, result sql.NullString
	var httpStatus, ttft, total, generation, in, out sql.NullInt64
	var detail []byte
	var cost sql.NullString
	if err := row.Scan(&s.ID, &s.JobID, &s.TestID, &s.Seq, &traceID, &status, &s.PreparedAt, &dispatchedAt, &receivedAt, &clientRequestID, &httpStatus,
		&ttft, &total, &generation, &in, &out, &responseModel, &result, &s.ErrorCode, &detail, &s.BillingStatus, &cost); err != nil {
		return nil, err
	}
	s.Status = service.ProviderHallSampleStatus(status)
	s.TraceID, _ = uuid.Parse(traceID)
	s.PreparedAt = s.PreparedAt.UTC()
	if dispatchedAt.Valid {
		t := dispatchedAt.Time.UTC()
		s.DispatchedAt = &t
	}
	if receivedAt.Valid {
		t := receivedAt.Time.UTC()
		s.ReceivedAt = &t
	}
	ns := func(v sql.NullString) *string {
		if !v.Valid {
			return nil
		}
		return &v.String
	}
	ni := func(v sql.NullInt64) *int {
		if !v.Valid {
			return nil
		}
		i := int(v.Int64)
		return &i
	}
	s.ClientRequestID, s.ResponseModel, s.Result = ns(clientRequestID), ns(responseModel), ns(result)
	s.HTTPStatus, s.TTFTMs, s.TotalMs, s.GenerationMs, s.InputTokens, s.OutputTokens = ni(httpStatus), ni(ttft), ni(total), ni(generation), ni(in), ni(out)
	if len(detail) > 0 {
		s.Detail = json.RawMessage(detail)
	} else {
		s.Detail = json.RawMessage(`{}`)
	}
	if cost.Valid {
		if d, err := decimal.NewFromString(cost.String); err == nil {
			s.ActualCost = &d
		}
	}
	return &s, nil
}

func (r *providerHallJobRepository) ListSchedulableTargets(ctx context.Context) ([]service.ProviderHallSchedulableTarget, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT t.id, t.group_id, t.profile_id, t.probe_key_id, t.version, t.probe_interval_seconds, t.verification_interval_seconds, t.enabled,
       COALESCE(hg.listed, false) AND g.id IS NOT NULL AND g.status = 'active',
       p.id, p.version, p.model, p.protocol, p.supports_tools, p.output_limit, p.model_aliases
FROM provider_hall_targets t
JOIN provider_hall_profiles p ON p.id = t.profile_id
LEFT JOIN provider_hall_groups hg ON hg.group_id = t.group_id
LEFT JOIN groups g ON g.id = t.group_id AND g.deleted_at IS NULL
ORDER BY t.group_id, t.profile_id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallSchedulableTarget
	for rows.Next() {
		var t service.ProviderHallSchedulableTarget
		var probeKey sql.NullInt64
		var aliases []byte
		if err := rows.Scan(&t.TargetID, &t.GroupID, &t.ProfileID, &probeKey, &t.TargetVersion, &t.ProbeIntervalSeconds, &t.VerificationIntervalSeconds, &t.Enabled, &t.Listed,
			&t.Profile.ID, &t.Profile.Version, &t.Profile.Model, &t.Profile.Protocol, &t.Profile.SupportsTools, &t.Profile.OutputLimit, &aliases); err != nil {
			return nil, err
		}
		if probeKey.Valid {
			id := probeKey.Int64
			t.ProbeKeyID = &id
		}
		if len(aliases) > 0 {
			_ = json.Unmarshal(aliases, &t.Profile.ModelAliases)
		}
		if t.Profile.ModelAliases == nil {
			t.Profile.ModelAliases = []string{}
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func providerHallSnapshotJSON(s service.ProviderHallJobSnapshot) ([]byte, error) {
	if s.Profile.ModelAliases == nil {
		s.Profile.ModelAliases = []string{}
	}
	return json.Marshal(s)
}

func (r *providerHallJobRepository) EnqueueSlot(ctx context.Context, in service.ProviderHallEnqueueInput) (int64, bool, error) {
	if in.SlotAt == nil {
		return 0, false, errors.New("provider hall: slot_at required")
	}
	snapshot, err := providerHallSnapshotJSON(in.Snapshot)
	if err != nil {
		return 0, false, err
	}
	var id int64
	err = r.db.QueryRowContext(ctx, `
INSERT INTO provider_hall_jobs (kind, target_id, group_id, profile_id, config_snapshot, slot_at, not_before)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (target_id, kind, slot_at) WHERE slot_at IS NOT NULL DO NOTHING
RETURNING id`, string(in.Kind), in.TargetID, in.GroupID, in.ProfileID, snapshot, in.SlotAt.UTC(), providerHallNullTime(in.NotBefore)).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return id, true, nil
}

func providerHallNullTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC()
}

func (r *providerHallJobRepository) EnqueueManual(ctx context.Context, in service.ProviderHallEnqueueInput) (*service.ProviderHallJob, bool, error) {
	snapshot, err := providerHallSnapshotJSON(in.Snapshot)
	if err != nil {
		return nil, false, err
	}
	var idem any
	if in.IdempotencyKey != nil && strings.TrimSpace(*in.IdempotencyKey) != "" {
		idem = strings.TrimSpace(*in.IdempotencyKey)
	}
	var requestedBy any
	if in.RequestedBy != nil {
		requestedBy = *in.RequestedBy
	}
	row := r.db.QueryRowContext(ctx, `
INSERT INTO provider_hall_jobs (kind, target_id, group_id, profile_id, config_snapshot, not_before, idempotency_key, requested_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (idempotency_key) WHERE idempotency_key IS NOT NULL DO NOTHING
RETURNING `+providerHallJobColumns, string(in.Kind), in.TargetID, in.GroupID, in.ProfileID, snapshot, providerHallNullTime(in.NotBefore), idem, requestedBy)
	job, err := providerHallScanJob(row)
	if err == nil {
		return job, false, nil
	}
	if !errors.Is(err, sql.ErrNoRows) || idem == nil {
		return nil, false, err
	}
	existing, err := providerHallScanJob(r.db.QueryRowContext(ctx, `SELECT `+providerHallJobColumns+` FROM provider_hall_jobs WHERE idempotency_key = $1`, idem))
	if err != nil {
		return nil, false, err
	}
	return existing, true, nil
}

func (r *providerHallJobRepository) CancelQueuedNotIn(ctx context.Context, activeTargetIDs []int64, reason string, now time.Time) (int64, error) {
	if activeTargetIDs == nil {
		activeTargetIDs = []int64{}
	}
	res, err := r.db.ExecContext(ctx, `
UPDATE provider_hall_jobs SET status = 'cancelled', error_code = $2, finished_at = $3, updated_at = $3, lease_owner = NULL, lease_until = NULL
WHERE status = 'queued' AND target_id <> ALL($1::bigint[])`, pq.Array(activeTargetIDs), reason, now.UTC())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *providerHallJobRepository) CancelQueuedForTarget(ctx context.Context, targetID int64, reason string, now time.Time) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
UPDATE provider_hall_jobs SET status = 'cancelled', error_code = $2, finished_at = $3, updated_at = $3, lease_owner = NULL, lease_until = NULL
WHERE status = 'queued' AND target_id = $1`, targetID, reason, now.UTC())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// CancelQueuedDrifted cancels queued jobs whose snapshot no longer matches the
// live target/profile versions or probe key.
func (r *providerHallJobRepository) CancelQueuedDrifted(ctx context.Context, now time.Time) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
UPDATE provider_hall_jobs j SET status = 'cancelled', error_code = 'config_changed', finished_at = $1, updated_at = $1
FROM provider_hall_targets t, provider_hall_profiles p
WHERE j.status = 'queued' AND t.id = j.target_id AND p.id = j.profile_id AND (
      (j.config_snapshot->>'target_version')::bigint <> t.version
   OR (j.config_snapshot->'profile'->>'version')::bigint <> p.version
   OR (j.config_snapshot->>'probe_key_id')::bigint IS DISTINCT FROM t.probe_key_id)`, now.UTC())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *providerHallJobRepository) RecoverExpiredLeases(ctx context.Context, now time.Time) (int64, int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer func() { _ = tx.Rollback() }()
	// Jobs with an in-flight sample are never resent: they only get checked.
	unknownRes, err := tx.ExecContext(ctx, `
UPDATE provider_hall_jobs j SET status = 'unknown', error_code = 'lease_lost', lease_owner = NULL, lease_until = NULL, updated_at = $1
WHERE j.status = 'running' AND j.lease_until < $1
  AND EXISTS (SELECT 1 FROM provider_hall_samples s WHERE s.job_id = j.id AND s.status IN ('dispatched', 'uncertain'))`, now.UTC())
	if err != nil {
		return 0, 0, err
	}
	requeueRes, err := tx.ExecContext(ctx, `
UPDATE provider_hall_jobs j SET status = 'queued', lease_owner = NULL, lease_until = NULL, updated_at = $1
WHERE j.status = 'running' AND j.lease_until < $1`, now.UTC())
	if err != nil {
		return 0, 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}
	unknown, _ := unknownRes.RowsAffected()
	requeued, _ := requeueRes.RowsAffected()
	return requeued, unknown, nil
}

func (r *providerHallJobRepository) Claim(ctx context.Context, owner string, now time.Time, lease time.Duration) (*service.ProviderHallJob, error) {
	row := r.db.QueryRowContext(ctx, `
UPDATE provider_hall_jobs SET status = 'running', lease_owner = $1, lease_until = $2, started_at = COALESCE(started_at, $3), attempts = attempts + 1, updated_at = $3
WHERE id = (
    SELECT id FROM provider_hall_jobs q
    WHERE q.status = 'queued' AND (q.not_before IS NULL OR q.not_before <= $3)
      AND NOT EXISTS (SELECT 1 FROM provider_hall_jobs r WHERE r.target_id = q.target_id AND r.status = 'running')
    ORDER BY q.not_before NULLS FIRST, q.created_at
    LIMIT 1 FOR UPDATE SKIP LOCKED)
RETURNING `+providerHallJobColumns, owner, now.Add(lease).UTC(), now.UTC())
	job, err := providerHallScanJob(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return job, err
}

func (r *providerHallJobRepository) RenewLease(ctx context.Context, jobID int64, owner string, until time.Time) (bool, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE provider_hall_jobs SET lease_until = $3, updated_at = now() WHERE id = $1 AND lease_owner = $2 AND status = 'running'`, jobID, owner, until.UTC())
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
}

func (r *providerHallJobRepository) Requeue(ctx context.Context, jobID int64, notBefore time.Time, code string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE provider_hall_jobs SET status = 'queued', lease_owner = NULL, lease_until = NULL, not_before = $2, error_code = $3, updated_at = now() WHERE id = $1 AND status = 'running'`, jobID, notBefore.UTC(), code)
	return err
}

func (r *providerHallJobRepository) Complete(ctx context.Context, jobID int64, status service.ProviderHallJobStatus, code, message string, now time.Time) error {
	if len(message) > 2000 {
		message = message[:2000]
	}
	_, err := r.db.ExecContext(ctx, `
UPDATE provider_hall_jobs SET status = $2, error_code = $3, error_message = $4, finished_at = $5, updated_at = $5, lease_owner = NULL, lease_until = NULL
WHERE id = $1 AND status IN ('running', 'unknown', 'queued')`, jobID, string(status), code, message, now.UTC())
	return err
}

func (r *providerHallJobRepository) MarkUnknown(ctx context.Context, jobID int64, code string, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE provider_hall_jobs SET status = 'unknown', error_code = $2, lease_owner = NULL, lease_until = NULL, updated_at = $3 WHERE id = $1 AND status = 'running'`, jobID, code, now.UTC())
	return err
}

func (r *providerHallJobRepository) Cancel(ctx context.Context, jobID int64, reason string, now time.Time) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	var status string
	if err := tx.QueryRowContext(ctx, `SELECT status FROM provider_hall_jobs WHERE id = $1 FOR UPDATE`, jobID).Scan(&status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return service.ErrProviderHallNotFound
		}
		return err
	}
	switch service.ProviderHallJobStatus(status) {
	case service.ProviderHallJobQueued:
	case service.ProviderHallJobRunning, service.ProviderHallJobUnknown:
		var inflight int
		if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_samples WHERE job_id = $1 AND status IN ('dispatched', 'uncertain', 'received')`, jobID).Scan(&inflight); err != nil {
			return err
		}
		if inflight > 0 {
			return service.ErrProviderHallJobInFlight
		}
	default:
		return service.ErrProviderHallJobNotCancellable
	}
	if _, err := tx.ExecContext(ctx, `UPDATE provider_hall_jobs SET status = 'cancelled', error_code = $2, finished_at = $3, updated_at = $3, lease_owner = NULL, lease_until = NULL WHERE id = $1`, jobID, reason, now.UTC()); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *providerHallJobRepository) CreateSamples(ctx context.Context, jobID int64, cases []service.ProviderHallTestCase) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	for _, tc := range cases {
		detail, _ := json.Marshal(map[string]any{"case": tc})
		if _, err := tx.ExecContext(ctx, `
INSERT INTO provider_hall_samples (job_id, test_id, seq, trace_id, detail) VALUES ($1, $2, $3, gen_random_uuid(), $4)
ON CONFLICT (job_id, test_id, seq) DO NOTHING`, jobID, tc.TestID, tc.Seq, detail); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *providerHallJobRepository) ListSamples(ctx context.Context, jobID int64) ([]service.ProviderHallSample, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+providerHallSampleColumns+` FROM provider_hall_samples WHERE job_id = $1 ORDER BY test_id, seq`, jobID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallSample
	for rows.Next() {
		s, err := providerHallScanSample(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *s)
	}
	return out, rows.Err()
}

// TryDispatch runs the pre-dispatch checks in the fixed order under a global
// advisory lock and, when everything passes, marks the sample dispatched. The
// caller sends only after this commits.
func (r *providerHallJobRepository) TryDispatch(ctx context.Context, in service.ProviderHallDispatchInput) (service.ProviderHallDispatchOutcome, error) {
	job := in.Job
	if job == nil {
		return service.ProviderHallDispatchOutcome{}, errors.New("provider hall: job required")
	}
	if in.MaxInflight <= 0 {
		in.MaxInflight = 2
	}
	if in.BacklogAfter <= 0 {
		in.BacklogAfter = 120 * time.Second
	}
	now := in.Now.UTC()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtext('provider_hall_dispatch'))`); err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	cancel := func(code, msg string) (service.ProviderHallDispatchOutcome, error) {
		return service.ProviderHallDispatchOutcome{Decision: service.ProviderHallDispatchCancel, Code: code, Message: msg}, nil
	}
	client := ent.NewClient(ent.Driver(entsql.NewDriver(dialect.Postgres, entsql.Conn{ExecQuerier: tx})))

	// 1. Switches, target, listing and snapshot versions.
	cfg, err := client.ProviderHallConfig.Query().Where(providerhallconfig.IDEQ(1)).Only(ctx)
	if err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	if !cfg.TasksEnabled || !in.TasksEnabled {
		return cancel(service.ProviderHallJobCodeTasksDisabled, "tasks are disabled")
	}
	target, err := client.ProviderHallTarget.Query().Where(providerhalltarget.IDEQ(job.TargetID)).Only(ctx)
	if ent.IsNotFound(err) || err == nil && !target.Enabled {
		return cancel(service.ProviderHallJobCodeTargetDisabled, "target disabled or removed")
	}
	if err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	listing, err := client.ProviderHallGroup.Query().Where(providerhallgroup.IDEQ(job.GroupID)).Only(ctx)
	if ent.IsNotFound(err) || err == nil && !listing.Listed {
		return cancel(service.ProviderHallJobCodeTargetDisabled, "group not listed")
	}
	if err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	profile, err := client.ProviderHallProfile.Query().Where(providerhallprofile.IDEQ(job.ProfileID)).Only(ctx)
	if ent.IsNotFound(err) {
		return cancel(service.ProviderHallJobCodeConfigChanged, "profile removed")
	}
	if err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	if target.Version != job.Snapshot.TargetVersion || profile.Version != job.Snapshot.Profile.Version ||
		target.ProbeKeyID == nil || *target.ProbeKeyID != job.Snapshot.ProbeKeyID ||
		cfg.OperatorUserID == nil || *cfg.OperatorUserID != job.Snapshot.OperatorUserID || cfg.GatewayOrigin != job.Snapshot.GatewayOrigin {
		return cancel(service.ProviderHallJobCodeConfigChanged, "target, profile or configuration changed since the job was created")
	}

	// 2. Binding recheck with live rows: an invalid key never sends a paid request.
	g, err := client.Group.Query().Where(group.IDEQ(job.GroupID), group.DeletedAtIsNil()).Only(ctx)
	if ent.IsNotFound(err) {
		return cancel(service.ProviderHallJobCodeTargetDisabled, "group removed")
	}
	if err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	fail := func(code, msg string) (service.ProviderHallDispatchOutcome, error) {
		return service.ProviderHallDispatchOutcome{Decision: service.ProviderHallDispatchFail, Code: code, Message: msg}, nil
	}
	u, err := client.User.Query().Where(user.IDEQ(*cfg.OperatorUserID), user.DeletedAtIsNil()).WithAllowedGroups().Only(ctx)
	if ent.IsNotFound(err) {
		return fail(service.ProviderHallJobCodeProbeKeyInvalid, "operator missing")
	}
	if err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	operator := userEntityToService(u)
	for _, allowed := range u.Edges.AllowedGroups {
		operator.AllowedGroups = append(operator.AllowedGroups, allowed.ID)
	}
	hasSubscription := false
	if g.SubscriptionType == service.SubscriptionTypeSubscription {
		hasSubscription, err = client.UserSubscription.Query().Where(
			usersubscription.UserIDEQ(u.ID), usersubscription.GroupIDEQ(g.ID), usersubscription.DeletedAtIsNil(),
			usersubscription.StatusEQ(service.SubscriptionStatusActive), usersubscription.ExpiresAtGT(now)).Exist(ctx)
		if err != nil {
			return service.ProviderHallDispatchOutcome{}, err
		}
	}
	k, err := client.APIKey.Query().Where(apikey.IDEQ(*target.ProbeKeyID), apikey.DeletedAtIsNil()).Only(ctx)
	if ent.IsNotFound(err) {
		return fail(service.ProviderHallJobCodeProbeKeyInvalid, "probe key missing")
	}
	if err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	var route service.CompositeRouteDecision
	if g.Platform == service.PlatformComposite {
		route, err = service.NewCompositeRouteResolver(NewCompositeModelRouteRepository(client)).Resolve(ctx, g.ID, profile.Model, string(profile.Protocol))
		if err != nil {
			return service.ProviderHallDispatchOutcome{}, err
		}
	}
	if err := service.ValidateProviderHallTargetBinding(groupEntityToService(g), providerHallProfileFromEnt(profile), apiKeyEntityToService(k), operator, hasSubscription, route, now); err != nil {
		return fail(service.ProviderHallJobCodeProbeKeyInvalid, err.Error())
	}

	// 3. Billing backlog: an old unconfirmed bill pauses dispatch.
	var backlog int
	if err := tx.QueryRowContext(ctx, `
SELECT count(*) FROM provider_hall_samples
WHERE status IN ('dispatched', 'received', 'uncertain') AND billing_status IN ('pending', 'uncertain') AND dispatched_at < $1`, now.Add(-in.BacklogAfter)).Scan(&backlog); err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	if backlog > 0 {
		return service.ProviderHallDispatchOutcome{Decision: service.ProviderHallDispatchRequeue, Code: service.ProviderHallJobCodeBillingBacklog, Message: fmt.Sprintf("%d bills pending longer than %s", backlog, in.BacklogAfter)}, nil
	}

	// 4. Budget: the job's day is fixed at first dispatch.
	budgetDay := service.ProviderHallBudgetDayFor(job.BudgetDay, now)
	var confirmed string
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(sum(actual_cost), 0)::text FROM provider_hall_spend WHERE budget_day = $1::date AND status = 'confirmed'`, budgetDay).Scan(&confirmed); err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	confirmedSpend, _ := decimal.NewFromString(confirmed)
	if service.ProviderHallBudgetExhausted(cfg.DailyBudget.String(), confirmedSpend) {
		return cancel(service.ProviderHallJobCodeBudgetExhausted, "daily budget reached")
	}

	// 5. Concurrency: at most MaxInflight dispatched samples globally, one per target.
	var globalInflight, targetInflight int
	if err := tx.QueryRowContext(ctx, `
SELECT count(*) FILTER (WHERE true), count(*) FILTER (WHERE j.target_id = $1)
FROM provider_hall_samples s JOIN provider_hall_jobs j ON j.id = s.job_id
WHERE s.status = 'dispatched'`, job.TargetID).Scan(&globalInflight, &targetInflight); err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	if globalInflight >= in.MaxInflight || targetInflight > 0 {
		return service.ProviderHallDispatchOutcome{Decision: service.ProviderHallDispatchWait, Code: "concurrency"}, nil
	}
	res, err := tx.ExecContext(ctx, `UPDATE provider_hall_samples SET status = 'dispatched', dispatched_at = $2 WHERE id = $1 AND job_id = $3 AND status = 'prepared'`, in.SampleID, now, job.ID)
	if err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	if n, _ := res.RowsAffected(); n != 1 {
		return service.ProviderHallDispatchOutcome{}, fmt.Errorf("provider hall: sample %d is not prepared", in.SampleID)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE provider_hall_jobs SET budget_day = COALESCE(budget_day, $2::date), updated_at = $3 WHERE id = $1`, job.ID, budgetDay, now); err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	if err := tx.Commit(); err != nil {
		return service.ProviderHallDispatchOutcome{}, err
	}
	return service.ProviderHallDispatchOutcome{Decision: service.ProviderHallDispatchGo, Key: k.Key, BudgetDay: budgetDay, Origin: cfg.GatewayOrigin}, nil
}

func providerHallDetailJSON(detail map[string]any) []byte {
	if detail == nil {
		detail = map[string]any{}
	}
	b, err := json.Marshal(detail)
	if err != nil {
		return []byte(`{}`)
	}
	if len(b) > 4096 {
		// Detail is diagnostic only; keep the row small.
		b, _ = json.Marshal(map[string]any{"truncated": true, "test_id": detail["test_id"], "seq": detail["seq"], "error": detail["error"]})
	}
	return b
}

func (r *providerHallJobRepository) MarkSampleReceived(ctx context.Context, rec service.ProviderHallSampleReceipt) error {
	toNull := func(v *int) any {
		if v == nil {
			return nil
		}
		return *v
	}
	var model, result any
	if rec.ResponseModel != "" {
		model = rec.ResponseModel
	}
	if rec.Result != "" {
		result = rec.Result
	}
	var clientID any
	if rec.ClientRequestID != "" {
		clientID = rec.ClientRequestID
	}
	_, err := r.db.ExecContext(ctx, `
UPDATE provider_hall_samples SET status = 'received', received_at = $2, client_request_id = COALESCE($3, client_request_id), http_status = $4,
    ttft_ms = $5, total_ms = $6, generation_ms = $7, input_tokens = $8, output_tokens = $9, response_model = $10, result = $11, error_code = $12, detail = $13
WHERE id = $1 AND status IN ('dispatched', 'uncertain')`,
		rec.SampleID, rec.ReceivedAt.UTC(), clientID, rec.HTTPStatus, toNull(rec.TTFTMs), toNull(rec.TotalMs), toNull(rec.GenerationMs),
		toNull(rec.InputTokens), toNull(rec.OutputTokens), model, result, rec.ErrorCode, providerHallDetailJSON(rec.Detail))
	return err
}

func (r *providerHallJobRepository) MarkSampleUncertain(ctx context.Context, sampleID int64, clientRequestID, code string, now time.Time) error {
	var clientID any
	if clientRequestID != "" {
		clientID = clientRequestID
	}
	_, err := r.db.ExecContext(ctx, `UPDATE provider_hall_samples SET status = 'uncertain', client_request_id = COALESCE($2, client_request_id), error_code = $3 WHERE id = $1 AND status = 'dispatched'`, sampleID, clientID, code)
	_ = now
	return err
}

// ReconcileSamples fills uncertain samples from the collector's request facts
// (joined on trace_id), mirrors billing outcomes, and times out samples that
// never produced a fact.
// providerHallDefinitiveBillingGrace is how long a sample with a definitive
// unbillable outcome keeps billing_status=pending so a late asynchronous bill
// can still land before it is closed as unbilled.
const providerHallDefinitiveBillingGrace = 90 * time.Second

func (r *providerHallJobRepository) ReconcileSamples(ctx context.Context, now time.Time, uncertainTimeout time.Duration) (int64, error) {
	if uncertainTimeout <= 0 {
		uncertainTimeout = 10 * time.Minute
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()
	var total int64
	// Uncertain samples whose request reached a terminal state: the response
	// content is lost, so the case is an error, but timing and tokens survive.
	res, err := tx.ExecContext(ctx, `
UPDATE provider_hall_samples s SET status = 'received', received_at = COALESCE(r.ended_at, $1),
    ttft_ms = COALESCE(s.ttft_ms, r.ttft_ms), input_tokens = COALESCE(s.input_tokens, r.input_tokens::int), output_tokens = COALESCE(s.output_tokens, r.output_tokens::int),
    response_model = COALESCE(s.response_model, NULLIF(r.response_model, '')),
    result = CASE WHEN r.outcome = 'success' AND s.test_id = 'probe' THEN 'passed' WHEN r.outcome = 'success' THEN 'error' ELSE 'error' END,
    error_code = CASE WHEN r.outcome = 'success' AND s.test_id = 'probe' THEN '' WHEN r.outcome = 'success' THEN 'response_lost' ELSE 'request_' || r.outcome END
FROM provider_hall_requests r
WHERE r.trace_id = s.trace_id AND s.status = 'uncertain' AND r.outcome <> 'pending' AND r.ended_at IS NOT NULL`, now.UTC())
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	total += n
	// Billing mirror from facts (safety net for events the collector could not
	// attach because the sample row raced its own dispatch).
	res, err = tx.ExecContext(ctx, `
UPDATE provider_hall_samples s SET
    billing_status = CASE r.billing_status WHEN 'applied' THEN 'confirmed' WHEN 'failed' THEN 'failed' WHEN 'uncertain' THEN 'uncertain' WHEN 'not_applicable' THEN 'unbilled' ELSE s.billing_status END,
    actual_cost = CASE WHEN r.billing_status = 'applied' THEN r.actual_cost ELSE s.actual_cost END
FROM provider_hall_requests r
WHERE r.trace_id = s.trace_id AND s.billing_status IN ('pending', 'uncertain') AND r.billing_status IN ('applied', 'failed', 'uncertain', 'not_applicable')
  AND r.billing_status::text <> CASE s.billing_status WHEN 'uncertain' THEN 'uncertain' ELSE '' END`)
	if err != nil {
		return 0, err
	}
	n, _ = res.RowsAffected()
	total += n
	res, err = tx.ExecContext(ctx, `
INSERT INTO provider_hall_spend (billing_request_id, api_key_id, job_id, sample_id, budget_day, actual_cost, status, confirmed_at)
SELECT r.billing_request_id, COALESCE(r.billing_api_key_id, r.api_key_id), s.job_id, s.id,
       COALESCE(j.budget_day, (COALESCE(r.billed_at, $1::timestamptz) AT TIME ZONE 'Asia/Shanghai')::date),
       COALESCE(r.actual_cost, 0),
       CASE r.billing_status WHEN 'applied' THEN 'confirmed' WHEN 'failed' THEN 'failed' ELSE 'uncertain' END,
       CASE WHEN r.billing_status = 'applied' THEN r.billed_at END
FROM provider_hall_requests r JOIN provider_hall_samples s ON s.trace_id = r.trace_id JOIN provider_hall_jobs j ON j.id = s.job_id
WHERE r.source IN ('probe', 'verification') AND r.billing_request_id IS NOT NULL AND r.billing_status IN ('applied', 'failed', 'uncertain')
ON CONFLICT (billing_request_id, api_key_id) DO UPDATE SET status = EXCLUDED.status, actual_cost = EXCLUDED.actual_cost, confirmed_at = EXCLUDED.confirmed_at
WHERE provider_hall_spend.status <> 'confirmed' AND provider_hall_spend.status <> EXCLUDED.status`, now.UTC())
	if err != nil {
		return 0, err
	}
	n, _ = res.RowsAffected()
	total += n
	// Uncertain samples with no fact after the timeout are closed as errors.
	res, err = tx.ExecContext(ctx, `
UPDATE provider_hall_samples s SET status = 'received', received_at = $1, result = 'error', error_code = 'uncertain_timeout',
    billing_status = CASE WHEN s.billing_status = 'pending' THEN 'unbilled' ELSE s.billing_status END
WHERE s.status = 'uncertain' AND s.dispatched_at < $2
  AND NOT EXISTS (SELECT 1 FROM provider_hall_requests r WHERE r.trace_id = s.trace_id AND (r.outcome <> 'pending' OR r.billing_status NOT IN ('pending')))`, now.UTC(), now.Add(-uncertainTimeout).UTC())
	if err != nil {
		return 0, err
	}
	n, _ = res.RowsAffected()
	total += n
	// Received samples that will never be billed (definitive non-2xx answers,
	// excluded/failed requests) stop blocking the backlog check after a short
	// grace for the asynchronous billing worker. Waiting the full uncertain
	// timeout here would pause every dispatch for ten minutes after one 4xx.
	res, err = tx.ExecContext(ctx, `
UPDATE provider_hall_samples s SET billing_status = 'unbilled'
WHERE s.status = 'received' AND s.billing_status = 'pending' AND s.received_at < $1
  AND (s.result = 'error' AND s.http_status IS NOT NULL AND (s.http_status < 200 OR s.http_status >= 300)
       OR EXISTS (SELECT 1 FROM provider_hall_requests r WHERE r.trace_id = s.trace_id AND r.outcome IN ('excluded', 'failed') AND r.billing_status = 'pending'))`, now.Add(-providerHallDefinitiveBillingGrace).UTC())
	if err != nil {
		return 0, err
	}
	n, _ = res.RowsAffected()
	total += n
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *providerHallJobRepository) ListInflight(ctx context.Context) ([]service.ProviderHallSample, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+providerHallSampleColumns+` FROM provider_hall_samples WHERE status IN ('dispatched', 'uncertain') ORDER BY dispatched_at`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallSample
	for rows.Next() {
		s, err := providerHallScanSample(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *s)
	}
	return out, rows.Err()
}

// ListFinalizableJobs returns unknown jobs (and running jobs whose lease has
// lapsed) whose samples are all terminal, so the runner can score them.
func (r *providerHallJobRepository) ListFinalizableJobs(ctx context.Context, now time.Time) ([]service.ProviderHallJob, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT `+providerHallJobColumns+` FROM provider_hall_jobs j
WHERE (j.status = 'unknown' OR j.status = 'running' AND j.lease_until < $1)
  AND EXISTS (SELECT 1 FROM provider_hall_samples s WHERE s.job_id = j.id)
  AND NOT EXISTS (SELECT 1 FROM provider_hall_samples s WHERE s.job_id = j.id AND s.status IN ('prepared', 'dispatched', 'uncertain'))
ORDER BY j.id LIMIT 100`, now.UTC())
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallJob
	for rows.Next() {
		j, err := providerHallScanJob(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *j)
	}
	return out, rows.Err()
}

func (r *providerHallJobRepository) SaveVerification(ctx context.Context, v service.ProviderHallVerification) error {
	summary, err := json.Marshal(v.Summary)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
INSERT INTO provider_hall_verifications (job_id, group_id, profile_id, target_id, verdict, execution_status, reason_code, summary, profile_version, target_version, completed_at, expires_at, stale)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (job_id) DO UPDATE SET verdict = EXCLUDED.verdict, execution_status = EXCLUDED.execution_status, reason_code = EXCLUDED.reason_code,
    summary = EXCLUDED.summary, completed_at = EXCLUDED.completed_at, expires_at = EXCLUDED.expires_at, stale = EXCLUDED.stale`,
		v.JobID, v.GroupID, v.ProfileID, v.TargetID, v.Verdict, v.ExecutionStatus, v.ReasonCode, summary, v.ProfileVersion, v.TargetVersion, v.CompletedAt.UTC(), v.ExpiresAt.UTC(), v.Stale)
	return err
}

func (r *providerHallJobRepository) MarkVerificationsStale(ctx context.Context, profileID int64) (int64, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE provider_hall_verifications SET stale = true WHERE profile_id = $1 AND NOT stale`, profileID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// MarkDriftedVerificationsStale stales reports whose profile identity
// fingerprint (model/protocol/tools/aliases) no longer matches the profile.
func (r *providerHallJobRepository) MarkDriftedVerificationsStale(ctx context.Context) (int64, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT v.job_id, v.summary->>'profile_fingerprint', p.model, p.protocol, p.supports_tools, p.model_aliases
FROM provider_hall_verifications v JOIN provider_hall_profiles p ON p.id = v.profile_id
WHERE NOT v.stale AND v.expires_at > now() AND v.profile_version <> p.version`)
	if err != nil {
		return 0, err
	}
	var stale []int64
	for rows.Next() {
		var jobID int64
		var fingerprint sql.NullString
		var model, protocol string
		var tools bool
		var aliasesRaw []byte
		if err := rows.Scan(&jobID, &fingerprint, &model, &protocol, &tools, &aliasesRaw); err != nil {
			_ = rows.Close()
			return 0, err
		}
		var aliases []string
		_ = json.Unmarshal(aliasesRaw, &aliases)
		if !fingerprint.Valid || fingerprint.String != service.ProviderHallProfileFingerprint(model, protocol, tools, aliases) {
			stale = append(stale, jobID)
		}
	}
	_ = rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	if len(stale) == 0 {
		return 0, nil
	}
	res, err := r.db.ExecContext(ctx, `UPDATE provider_hall_verifications SET stale = true WHERE job_id = ANY($1::bigint[])`, pq.Array(stale))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *providerHallJobRepository) SumSpend(ctx context.Context, day string) (service.ProviderHallSpendSummary, error) {
	out := service.ProviderHallSpendSummary{Day: day}
	var confirmed, uncertain string
	if err := r.db.QueryRowContext(ctx, `
SELECT COALESCE(sum(actual_cost) FILTER (WHERE status = 'confirmed'), 0)::text, COALESCE(sum(actual_cost) FILTER (WHERE status = 'uncertain'), 0)::text
FROM provider_hall_spend WHERE budget_day = $1::date`, day).Scan(&confirmed, &uncertain); err != nil {
		return out, err
	}
	out.Confirmed, _ = decimal.NewFromString(confirmed)
	out.Uncertain, _ = decimal.NewFromString(uncertain)
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_samples WHERE status IN ('dispatched', 'uncertain')`).Scan(&out.InFlight); err != nil {
		return out, err
	}
	return out, nil
}

func (r *providerHallJobRepository) ListJobs(ctx context.Context, f service.ProviderHallJobFilter, page, pageSize int) ([]service.ProviderHallJob, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}
	where := []string{"true"}
	args := []any{}
	if f.Status != "" {
		args = append(args, string(f.Status))
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if f.Kind != "" {
		args = append(args, string(f.Kind))
		where = append(where, fmt.Sprintf("kind = $%d", len(args)))
	}
	if f.GroupID > 0 {
		args = append(args, f.GroupID)
		where = append(where, fmt.Sprintf("group_id = $%d", len(args)))
	}
	clause := strings.Join(where, " AND ")
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_jobs WHERE `+clause, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`SELECT %s FROM provider_hall_jobs WHERE %s ORDER BY created_at DESC, id DESC LIMIT $%d OFFSET $%d`, providerHallJobColumns, clause, len(args)-1, len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()
	out := []service.ProviderHallJob{}
	for rows.Next() {
		j, err := providerHallScanJob(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *j)
	}
	return out, total, rows.Err()
}

func (r *providerHallJobRepository) GetJob(ctx context.Context, id int64) (*service.ProviderHallJob, error) {
	job, err := providerHallScanJob(r.db.QueryRowContext(ctx, `SELECT `+providerHallJobColumns+` FROM provider_hall_jobs WHERE id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, service.ErrProviderHallNotFound
	}
	return job, err
}

const providerHallVerificationColumns = `job_id, group_id, profile_id, target_id, verdict, execution_status, reason_code, summary, profile_version, target_version, completed_at, expires_at, stale, created_at`

func providerHallScanVerification(row providerHallScanner) (*service.ProviderHallVerification, error) {
	var v service.ProviderHallVerification
	var summary []byte
	if err := row.Scan(&v.JobID, &v.GroupID, &v.ProfileID, &v.TargetID, &v.Verdict, &v.ExecutionStatus, &v.ReasonCode, &summary, &v.ProfileVersion, &v.TargetVersion, &v.CompletedAt, &v.ExpiresAt, &v.Stale, &v.CreatedAt); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(summary, &v.Summary)
	if v.Summary.Model.Seen == nil {
		v.Summary.Model.Seen = []string{}
	}
	v.CompletedAt, v.ExpiresAt, v.CreatedAt = v.CompletedAt.UTC(), v.ExpiresAt.UTC(), v.CreatedAt.UTC()
	return &v, nil
}

func (r *providerHallJobRepository) GetVerification(ctx context.Context, jobID int64) (*service.ProviderHallVerification, error) {
	v, err := providerHallScanVerification(r.db.QueryRowContext(ctx, `SELECT `+providerHallVerificationColumns+` FROM provider_hall_verifications WHERE job_id = $1`, jobID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return v, err
}

func (r *providerHallJobRepository) CountJobs(ctx context.Context, now time.Time) (service.ProviderHallJobCounts, error) {
	var c service.ProviderHallJobCounts
	err := r.db.QueryRowContext(ctx, `
SELECT count(*) FILTER (WHERE status = 'queued'), count(*) FILTER (WHERE status = 'running'), count(*) FILTER (WHERE status = 'unknown'),
       count(*) FILTER (WHERE status = 'failed' AND finished_at >= $1)
FROM provider_hall_jobs`, now.Add(-24*time.Hour).UTC()).Scan(&c.Queued, &c.Running, &c.Unknown, &c.Failed24h)
	return c, err
}

const providerHallProbeHealthSQL = `
SELECT j.target_id, j.group_id, j.profile_id, j.id, s.id, COALESCE(s.result, ''), s.error_code, s.received_at, s.ttft_ms, s.total_ms, s.generation_ms, s.input_tokens, s.output_tokens, s.dispatched_at
FROM provider_hall_samples s JOIN provider_hall_jobs j ON j.id = s.job_id
WHERE j.kind = 'probe' AND s.status = 'received' AND s.received_at >= $1`

func providerHallScanProbeHealth(rows *sql.Rows) (*service.ProviderHallProbeHealth, error) {
	var h service.ProviderHallProbeHealth
	var ttft, total, generation, in, out sql.NullInt64
	var dispatched sql.NullTime
	if err := rows.Scan(&h.TargetID, &h.GroupID, &h.ProfileID, &h.JobID, &h.SampleID, &h.Result, &h.ErrorCode, &h.ReceivedAt, &ttft, &total, &generation, &in, &out, &dispatched); err != nil {
		return nil, err
	}
	ni := func(v sql.NullInt64) *int {
		if !v.Valid {
			return nil
		}
		i := int(v.Int64)
		return &i
	}
	h.TTFTMs, h.TotalMs, h.GenerationMs, h.InputTokens, h.OutputTokens = ni(ttft), ni(total), ni(generation), ni(in), ni(out)
	h.ReceivedAt = h.ReceivedAt.UTC()
	if dispatched.Valid {
		t := dispatched.Time.UTC()
		h.DispatchedAt = &t
	}
	return &h, nil
}

func (r *providerHallJobRepository) LatestProbePerTarget(ctx context.Context, since time.Time) (map[int64]service.ProviderHallProbeHealth, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT DISTINCT ON (j.target_id) * FROM (`+providerHallProbeHealthSQL+`) AS j(target_id, group_id, profile_id, id, sid, result, error_code, received_at, ttft_ms, total_ms, generation_ms, input_tokens, output_tokens, dispatched_at) ORDER BY j.target_id, j.received_at DESC, j.sid DESC`, since.UTC())
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := map[int64]service.ProviderHallProbeHealth{}
	for rows.Next() {
		h, err := providerHallScanProbeHealth(rows)
		if err != nil {
			return nil, err
		}
		out[h.TargetID] = *h
	}
	return out, rows.Err()
}

func (r *providerHallJobRepository) RecentProbeSamples(ctx context.Context, groupID int64, profileID *int64, since time.Time, limit int) ([]service.ProviderHallProbeHealth, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	var pid any
	if profileID != nil {
		pid = *profileID
	}
	rows, err := r.db.QueryContext(ctx, providerHallProbeHealthSQL+` AND j.group_id = $2 AND ($3::bigint IS NULL OR j.profile_id = $3) ORDER BY s.received_at DESC, s.id DESC LIMIT $4`, since.UTC(), groupID, pid, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallProbeHealth
	for rows.Next() {
		h, err := providerHallScanProbeHealth(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *h)
	}
	return out, rows.Err()
}

func (r *providerHallJobRepository) ProbeSeries(ctx context.Context, q service.ProviderHallProbeSeriesQuery) ([]service.ProviderHallProbePoint, error) {
	if q.BucketSeconds <= 0 {
		return nil, errors.New("provider hall: bucket_seconds required")
	}
	var pid any
	if q.ProfileID != nil {
		pid = *q.ProfileID
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT date_bin(make_interval(secs => $1), s.received_at, $2::timestamptz) AS bucket,
       count(*), count(*) FILTER (WHERE s.result = 'passed'), count(*) FILTER (WHERE s.result <> 'passed' OR s.result IS NULL),
       avg(s.total_ms) FILTER (WHERE s.result = 'passed'), avg(s.ttft_ms) FILTER (WHERE s.result = 'passed')
FROM provider_hall_samples s JOIN provider_hall_jobs j ON j.id = s.job_id
WHERE j.kind = 'probe' AND s.status = 'received' AND s.received_at >= $2 AND s.received_at < $3 AND j.group_id = $4 AND ($5::bigint IS NULL OR j.profile_id = $5)
GROUP BY bucket ORDER BY bucket`, q.BucketSeconds, q.Start.UTC(), q.End.UTC(), q.GroupID, pid)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []service.ProviderHallProbePoint
	for rows.Next() {
		var p service.ProviderHallProbePoint
		var avgTotal, avgTTFT sql.NullFloat64
		if err := rows.Scan(&p.BucketStart, &p.Total, &p.Passed, &p.Failed, &avgTotal, &avgTTFT); err != nil {
			return nil, err
		}
		p.BucketStart = p.BucketStart.UTC()
		if avgTotal.Valid {
			p.AvgTotalMs = &avgTotal.Float64
		}
		if avgTTFT.Valid {
			p.AvgTTFTMs = &avgTTFT.Float64
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *providerHallJobRepository) LatestVerification(ctx context.Context, groupID, profileID int64) (*service.ProviderHallVerification, error) {
	v, err := providerHallScanVerification(r.db.QueryRowContext(ctx, `SELECT `+providerHallVerificationColumns+` FROM provider_hall_verifications WHERE group_id = $1 AND profile_id = $2 ORDER BY completed_at DESC, job_id DESC LIMIT 1`, groupID, profileID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return v, err
}

func (r *providerHallJobRepository) ListVerifications(ctx context.Context, groupID int64, profileID *int64, page, pageSize int) ([]service.ProviderHallVerification, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	var pid any
	if profileID != nil {
		pid = *profileID
	}
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_verifications WHERE group_id = $1 AND ($2::bigint IS NULL OR profile_id = $2)`, groupID, pid).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT `+providerHallVerificationColumns+` FROM provider_hall_verifications WHERE group_id = $1 AND ($2::bigint IS NULL OR profile_id = $2) ORDER BY completed_at DESC, job_id DESC LIMIT $3 OFFSET $4`, groupID, pid, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()
	out := []service.ProviderHallVerification{}
	for rows.Next() {
		v, err := providerHallScanVerification(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *v)
	}
	return out, total, rows.Err()
}
