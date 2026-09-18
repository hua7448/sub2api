//go:build integration || providerhall_localdb

package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func init() {
	registerProviderHallLocalDBContract("b2_facts", providerHallFactDatabaseContracts)
}

func providerHallFactRow(trace uuid.UUID, group int64, profile *int64, started time.Time) service.ProviderHallRequestRow {
	return service.ProviderHallRequestRow{
		TraceID: trace, NodeID: "node-fact", EpochID: 1, GroupID: group, ProfileID: profile, APIKeyID: 77,
		Protocol: "responses", RequestedModel: "gpt-fact", Source: service.ProviderHallSourceUser,
		StartedAt: started, Outcome: service.ProviderHallOutcomePending, AlgorithmVersion: service.ProviderHallAlgorithmVersion,
	}
}

func providerHallFinishRow(start service.ProviderHallRequestRow, endedAt time.Time) service.ProviderHallRequestRow {
	row := start
	first := start.StartedAt.Add(120 * time.Millisecond)
	ttft := 120
	in, out, read, create := int64(100), int64(5), int64(40), int64(0)
	row.FirstContentAt, row.EndedAt, row.TTFTMs = &first, &endedAt, &ttft
	row.Submissions, row.Outcome, row.UsageKnown = 2, service.ProviderHallOutcomeSuccess, true
	row.InputTokens, row.OutputTokens, row.CacheReadTokens, row.CacheCreation = &in, &out, &read, &create
	row.UpstreamModel, row.ResponseModel, row.Platform = "gpt-fact-up", "gpt-fact-2026", "openai"
	return row
}

func providerHallFactDatabaseContracts(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	repo := NewProviderHallFactRepository(db)
	base := time.Date(2026, 9, 12, 8, 30, 15, 0, time.UTC)
	t.Cleanup(func() {
		for _, q := range []string{
			`DELETE FROM provider_hall_dirty_buckets`, `DELETE FROM provider_hall_requests`,
			`DELETE FROM provider_hall_coverage_gaps`, `DELETE FROM provider_hall_collector_epochs`,
		} {
			_, _ = db.ExecContext(ctx, q)
		}
	})

	t.Run("fact_epoch_lifecycle_and_crash_gap", func(t *testing.T) {
		first, err := repo.RegisterEpoch(ctx, "node-epoch", "v1", 1, base)
		require.NoError(t, err)
		require.NoError(t, repo.Heartbeat(ctx, first, base.Add(10*time.Second)))
		var heartbeat time.Time
		require.NoError(t, db.QueryRowContext(ctx, `SELECT heartbeat_at FROM provider_hall_collector_epochs WHERE id=$1`, first).Scan(&heartbeat))
		require.True(t, heartbeat.Equal(base.Add(10*time.Second)))
		// A confirmed watermark bounds the crash gap.
		confirmed := base.Add(20 * time.Second)
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: first, MaxSeq: 5, ConfirmedAt: &confirmed}))

		// Same node registers again without a clean exit: the old epoch is a
		// crash and everything after its confirmation point is a gap.
		second, err := repo.RegisterEpoch(ctx, "node-epoch", "v2", 1, base.Add(time.Minute))
		require.NoError(t, err)
		require.NotEqual(t, first, second)
		var exitReason string
		var exitedAt sql.NullTime
		require.NoError(t, db.QueryRowContext(ctx, `SELECT exit_reason, exited_at FROM provider_hall_collector_epochs WHERE id=$1`, first).Scan(&exitReason, &exitedAt))
		require.Equal(t, "crash", exitReason)
		require.True(t, exitedAt.Valid)
		var gapStart, gapEnd time.Time
		var gapReason, gapScope string
		require.NoError(t, db.QueryRowContext(ctx, `SELECT started_at, ended_at, reason, scope FROM provider_hall_coverage_gaps WHERE epoch_id=$1`, first).Scan(&gapStart, &gapEnd, &gapReason, &gapScope))
		require.True(t, gapStart.Equal(confirmed), "crash gap starts at the last confirmed watermark")
		require.True(t, gapEnd.Equal(base.Add(time.Minute)))
		require.Equal(t, "crash", gapReason)
		require.Equal(t, "collection", gapScope)

		// Clean exit: no gap, heartbeat after exit is ignored.
		require.NoError(t, repo.MarkEpochExited(ctx, second, "shutdown", base.Add(2*time.Minute)))
		require.NoError(t, repo.Heartbeat(ctx, second, base.Add(3*time.Minute)))
		require.NoError(t, db.QueryRowContext(ctx, `SELECT exit_reason, heartbeat_at FROM provider_hall_collector_epochs WHERE id=$1`, second).Scan(&exitReason, &heartbeat))
		require.Equal(t, "shutdown", exitReason)
		require.True(t, heartbeat.Equal(base.Add(2*time.Minute)))
		var gaps int
		require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_coverage_gaps WHERE epoch_id=$1`, second).Scan(&gaps))
		require.Zero(t, gaps)
		third, err := repo.RegisterEpoch(ctx, "node-epoch", "v3", 1, base.Add(4*time.Minute))
		require.NoError(t, err)
		require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_coverage_gaps WHERE epoch_id=$1`, second).Scan(&gaps))
		require.Zero(t, gaps, "a cleanly exited epoch never becomes a crash")
		require.NoError(t, repo.MarkEpochExited(ctx, third, "shutdown", base.Add(5*time.Minute)))
	})

	t.Run("fact_upsert_order_independent_and_idempotent", func(t *testing.T) {
		epoch, err := repo.RegisterEpoch(ctx, "node-upsert", "v1", 1, base)
		require.NoError(t, err)
		profile := int64(500)
		// start -> finish -> finish again
		a := providerHallFactRow(uuid.New(), 900, &profile, base)
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Starts: []service.ProviderHallRequestRow{a}, MaxSeq: 1}))
		var outcome string
		var ended sql.NullTime
		var dirty int
		require.NoError(t, db.QueryRowContext(ctx, `SELECT outcome, ended_at FROM provider_hall_requests WHERE trace_id=$1`, a.TraceID.String()).Scan(&outcome, &ended))
		require.Equal(t, "pending", outcome)
		require.False(t, ended.Valid)
		require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_dirty_buckets WHERE group_id=900`).Scan(&dirty))
		require.Zero(t, dirty, "a start never dirties a minute")

		finishA := providerHallFinishRow(a, base.Add(2*time.Second))
		finishA.Seq = 2
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Finishes: []service.ProviderHallRequestRow{finishA}, MaxSeq: 2}))
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Finishes: []service.ProviderHallRequestRow{finishA}, MaxSeq: 2}))
		var submissions, ttft int
		var upstreamModel, responseModel, platform string
		var inputTokens, cacheRead int64
		var lastSeq int64
		require.NoError(t, db.QueryRowContext(ctx, `SELECT outcome, ended_at, submissions, ttft_ms, upstream_model, response_model, platform, input_tokens, cache_read_tokens, last_seq FROM provider_hall_requests WHERE trace_id=$1`, a.TraceID.String()).
			Scan(&outcome, &ended, &submissions, &ttft, &upstreamModel, &responseModel, &platform, &inputTokens, &cacheRead, &lastSeq))
		require.Equal(t, "success", outcome)
		require.True(t, ended.Time.Equal(base.Add(2*time.Second)))
		require.Equal(t, 2, submissions)
		require.Equal(t, 120, ttft)
		require.Equal(t, "gpt-fact-up", upstreamModel)
		require.Equal(t, "gpt-fact-2026", responseModel)
		require.Equal(t, "openai", platform)
		require.Equal(t, int64(100), inputTokens)
		require.Equal(t, int64(40), cacheRead)
		require.Equal(t, int64(2), lastSeq)
		var minute time.Time
		require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*), min(minute) FROM provider_hall_dirty_buckets WHERE group_id=900 AND profile_id=500`).Scan(&dirty, &minute))
		require.Equal(t, 1, dirty, "dirty minute is unique per (group, profile, minute)")
		require.True(t, minute.Equal(base.Add(2*time.Second).Truncate(time.Minute)))

		// finish -> late start: the start must not regress the terminal row.
		b := providerHallFactRow(uuid.New(), 900, &profile, base.Add(time.Minute))
		finishB := providerHallFinishRow(b, base.Add(time.Minute+time.Second))
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Finishes: []service.ProviderHallRequestRow{finishB}, MaxSeq: 3}))
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Starts: []service.ProviderHallRequestRow{b}, MaxSeq: 4}))
		require.NoError(t, db.QueryRowContext(ctx, `SELECT outcome, submissions FROM provider_hall_requests WHERE trace_id=$1`, b.TraceID.String()).Scan(&outcome, &submissions))
		require.Equal(t, "success", outcome)
		require.Equal(t, 2, submissions)
		// A second terminal for the same trace never overwrites the first.
		failedB := finishB
		failedB.Outcome, failedB.Submissions = service.ProviderHallOutcomeFailed, 9
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Finishes: []service.ProviderHallRequestRow{failedB}, MaxSeq: 5}))
		require.NoError(t, db.QueryRowContext(ctx, `SELECT outcome, submissions FROM provider_hall_requests WHERE trace_id=$1`, b.TraceID.String()).Scan(&outcome, &submissions))
		require.Equal(t, "success", outcome)
		require.Equal(t, 2, submissions)
		require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_dirty_buckets WHERE group_id=900`).Scan(&dirty))
		require.Equal(t, 2, dirty)

		// Mixed batch (start+finish+ a probe finish + a profile-less finish):
		// only user rows with a profile dirty a minute.
		c := providerHallFactRow(uuid.New(), 901, &profile, base.Add(2*time.Minute))
		finishC := providerHallFinishRow(c, base.Add(2*time.Minute+time.Second))
		probe := providerHallFactRow(uuid.New(), 901, &profile, base.Add(2*time.Minute))
		probe.Source = service.ProviderHallSourceProbe
		sample := int64(31)
		probe.SampleID = &sample
		finishProbe := providerHallFinishRow(probe, base.Add(2*time.Minute+time.Second))
		noProfile := providerHallFactRow(uuid.New(), 901, nil, base.Add(2*time.Minute))
		finishNoProfile := providerHallFinishRow(noProfile, base.Add(2*time.Minute+time.Second))
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Starts: []service.ProviderHallRequestRow{c, probe, noProfile},
			Finishes: []service.ProviderHallRequestRow{finishC, finishProbe, finishNoProfile}, MaxSeq: 6}))
		require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_dirty_buckets WHERE group_id=901`).Scan(&dirty))
		require.Equal(t, 1, dirty)
		var source string
		var sampleID sql.NullInt64
		require.NoError(t, db.QueryRowContext(ctx, `SELECT source, sample_id FROM provider_hall_requests WHERE trace_id=$1`, probe.TraceID.String()).Scan(&source, &sampleID))
		require.Equal(t, "probe", source)
		require.Equal(t, int64(31), sampleID.Int64)

		// Database constraints: ttft must be positive, outcome/source enumerated.
		bad := providerHallFinishRow(providerHallFactRow(uuid.New(), 901, &profile, base), base.Add(time.Second))
		zero := 0
		bad.TTFTMs = &zero
		require.Error(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Finishes: []service.ProviderHallRequestRow{bad}, MaxSeq: 7}))
		bad = providerHallFactRow(uuid.New(), 901, &profile, base)
		bad.Source = "bogus"
		require.Error(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Starts: []service.ProviderHallRequestRow{bad}, MaxSeq: 8}))
		require.NoError(t, repo.MarkEpochExited(ctx, epoch, "shutdown", base.Add(time.Hour)))
	})

	t.Run("fact_billing_link_before_and_after_finish", func(t *testing.T) {
		epoch, err := repo.RegisterEpoch(ctx, "node-billing", "v1", 1, base)
		require.NoError(t, err)
		profile := int64(510)
		mkBilling := func(trace uuid.UUID, status string) service.ProviderHallBillingEvent {
			return service.ProviderHallBillingEvent{
				TraceID: trace, BillingRequestID: "client:" + trace.String(), APIKeyID: 77, Fingerprint: "fp-" + status,
				Status: status, Mode: "token", IsSubscription: false, Multiplier: 1.5,
				InputBaseCost: decimal.RequireFromString("0.0000000123"), TotalBaseCost: decimal.RequireFromString("12345678901234.1234567890"),
				ActualCost: decimal.RequireFromString("0.12345678"), BilledAt: base.Add(3 * time.Second),
			}
		}
		readBilling := func(trace uuid.UUID) (status, inputBase, totalBase, actual string, mult string) {
			require.NoError(t, db.QueryRowContext(ctx, `SELECT billing_status, input_base_cost::text, total_base_cost::text, actual_cost::text, rate_multiplier::text FROM provider_hall_requests WHERE trace_id=$1`, trace.String()).
				Scan(&status, &inputBase, &totalBase, &actual, &mult))
			return
		}

		// Billing first (same batch as the start), finish later keeps the bill.
		a := providerHallFactRow(uuid.New(), 910, &profile, base)
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Starts: []service.ProviderHallRequestRow{a}, Billings: []service.ProviderHallBillingEvent{mkBilling(a.TraceID, "applied")}, MaxSeq: 1}))
		status, inputBase, totalBase, actual, mult := readBilling(a.TraceID)
		require.Equal(t, "applied", status)
		require.Equal(t, "0.0000000123", inputBase)
		require.Equal(t, "12345678901234.1234567890", totalBase)
		require.Equal(t, "0.12345678", actual)
		require.Equal(t, "1.500000", mult)
		var dirty int
		require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_dirty_buckets WHERE group_id=910`).Scan(&dirty))
		require.Zero(t, dirty, "a bill on an unfinished row has no minute to dirty yet")
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Finishes: []service.ProviderHallRequestRow{providerHallFinishRow(a, base.Add(2*time.Second))}, MaxSeq: 2}))
		status, _, _, actual, _ = readBilling(a.TraceID)
		require.Equal(t, "applied", status)
		require.Equal(t, "0.12345678", actual)
		var outcome string
		require.NoError(t, db.QueryRowContext(ctx, `SELECT outcome FROM provider_hall_requests WHERE trace_id=$1`, a.TraceID.String()).Scan(&outcome))
		require.Equal(t, "success", outcome)

		// Finish first, bill later: dirty minute is re-marked (unique) and the
		// applied bill is final; later duplicate/failed events are ignored.
		b := providerHallFactRow(uuid.New(), 910, &profile, base)
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Starts: []service.ProviderHallRequestRow{b}, Finishes: []service.ProviderHallRequestRow{providerHallFinishRow(b, base.Add(2*time.Second))}, MaxSeq: 3}))
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Billings: []service.ProviderHallBillingEvent{mkBilling(b.TraceID, "applied")}, MaxSeq: 4}))
		status, _, _, _, _ = readBilling(b.TraceID)
		require.Equal(t, "applied", status)
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Billings: []service.ProviderHallBillingEvent{mkBilling(b.TraceID, "failed")}, MaxSeq: 5}))
		status, _, _, actual, _ = readBilling(b.TraceID)
		require.Equal(t, "applied", status, "a settled bill is never overwritten")
		require.Equal(t, "0.12345678", actual)
		require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_dirty_buckets WHERE group_id=910`).Scan(&dirty))
		require.Equal(t, 1, dirty)

		// uncertain may still be resolved later; failed and not_applicable are terminal.
		c := providerHallFactRow(uuid.New(), 910, &profile, base)
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Starts: []service.ProviderHallRequestRow{c}, Billings: []service.ProviderHallBillingEvent{mkBilling(c.TraceID, "uncertain")}, MaxSeq: 6}))
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Billings: []service.ProviderHallBillingEvent{mkBilling(c.TraceID, "applied")}, MaxSeq: 7}))
		status, _, _, _, _ = readBilling(c.TraceID)
		require.Equal(t, "applied", status)
		d := providerHallFactRow(uuid.New(), 910, &profile, base)
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Starts: []service.ProviderHallRequestRow{d}, Billings: []service.ProviderHallBillingEvent{mkBilling(d.TraceID, "not_applicable")}, MaxSeq: 8}))
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Billings: []service.ProviderHallBillingEvent{mkBilling(d.TraceID, "applied")}, MaxSeq: 9}))
		status, _, _, _, _ = readBilling(d.TraceID)
		require.Equal(t, "not_applicable", status)
		// Billing for an unknown trace is a silent no-op inside the transaction.
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Billings: []service.ProviderHallBillingEvent{mkBilling(uuid.New(), "applied")}, MaxSeq: 10}))
		bogus := mkBilling(uuid.New(), "bogus")
		e := providerHallFactRow(bogus.TraceID, 910, &profile, base)
		require.Error(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, Starts: []service.ProviderHallRequestRow{e}, Billings: []service.ProviderHallBillingEvent{bogus}, MaxSeq: 11}), "unknown billing status is rejected and the batch rolls back")
		var count int
		require.NoError(t, db.QueryRowContext(ctx, `SELECT count(*) FROM provider_hall_requests WHERE trace_id=$1`, e.TraceID.String()).Scan(&count))
		require.Zero(t, count)
		require.NoError(t, repo.MarkEpochExited(ctx, epoch, "shutdown", base.Add(time.Hour)))
	})

	t.Run("fact_watermark_and_gaps", func(t *testing.T) {
		epoch, err := repo.RegisterEpoch(ctx, "node-wm", "v1", 1, base)
		require.NoError(t, err)
		read := func() (seq int64, confirmed sql.NullTime, overflowed bool) {
			require.NoError(t, db.QueryRowContext(ctx, `SELECT persisted_seq, confirmed_at, overflowed FROM provider_hall_collector_epochs WHERE id=$1`, epoch).Scan(&seq, &confirmed, &overflowed))
			return
		}
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, MaxSeq: 10, Overflowed: true}))
		seq, confirmed, overflowed := read()
		require.Equal(t, int64(10), seq)
		require.False(t, confirmed.Valid, "no barrier: the watermark must not advance")
		require.True(t, overflowed)
		at := base.Add(30 * time.Second)
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, MaxSeq: 12, ConfirmedAt: &at}))
		seq, confirmed, overflowed = read()
		require.Equal(t, int64(12), seq)
		require.True(t, confirmed.Time.Equal(at))
		require.False(t, overflowed)
		earlier := base.Add(10 * time.Second)
		require.NoError(t, repo.WriteBatch(ctx, service.ProviderHallFactBatch{EpochID: epoch, MaxSeq: 11, ConfirmedAt: &earlier}))
		seq, confirmed, _ = read()
		require.Equal(t, int64(12), seq, "seq never regresses")
		require.True(t, confirmed.Time.Equal(at), "watermark never regresses")

		gapID, err := repo.OpenGap(ctx, service.ProviderHallGap{NodeID: "node-wm", EpochID: epoch, Scope: "collection", StartedAt: at, Reason: "queue_overflow"})
		require.NoError(t, err)
		_, _, overflowed = read()
		require.True(t, overflowed, "an open gap flags the epoch")
		var ended sql.NullTime
		require.NoError(t, db.QueryRowContext(ctx, `SELECT ended_at FROM provider_hall_coverage_gaps WHERE id=$1`, gapID).Scan(&ended))
		require.False(t, ended.Valid)
		// Closing before the start clamps to the start; closing twice is a no-op.
		require.NoError(t, repo.CloseGap(ctx, gapID, at.Add(-time.Minute)))
		require.NoError(t, db.QueryRowContext(ctx, `SELECT ended_at FROM provider_hall_coverage_gaps WHERE id=$1`, gapID).Scan(&ended))
		require.True(t, ended.Time.Equal(at))
		require.NoError(t, repo.CloseGap(ctx, gapID, at.Add(time.Hour)))
		require.NoError(t, db.QueryRowContext(ctx, `SELECT ended_at FROM provider_hall_coverage_gaps WHERE id=$1`, gapID).Scan(&ended))
		require.True(t, ended.Time.Equal(at))
		_, err = repo.OpenGap(ctx, service.ProviderHallGap{NodeID: "node-wm", EpochID: epoch, Scope: "bogus", StartedAt: at, Reason: "x"})
		require.Error(t, err)
		require.NoError(t, repo.MarkEpochExited(ctx, epoch, "shutdown", base.Add(time.Hour)))
	})

	t.Run("fact_target_and_probe_key_index", func(t *testing.T) {
		f := newProviderHallTargetFixture(t, db)
		key := f.key(t)
		_, err := f.svc.SaveTargets(ctx, f.input(key.ID), f.user.ID)
		require.NoError(t, err)
		targets, err := repo.ListEnabledTargets(ctx)
		require.NoError(t, err)
		var found *service.ProviderHallTargetRef
		for i := range targets {
			if targets[i].GroupID == f.group.ID {
				found = &targets[i]
			}
		}
		require.NotNil(t, found)
		require.Equal(t, f.profiles[0].ID, found.ProfileID)
		require.Equal(t, f.profiles[0].Model, found.Model)
		require.Equal(t, "responses", found.Protocol)
		require.True(t, found.Enabled)
		require.NotNil(t, found.ProbeKeyID)
		require.Equal(t, key.ID, *found.ProbeKeyID)
		require.Equal(t, 300, found.ProbeIntervalSeconds)
		require.Equal(t, []string{}, found.ModelAliases)
		require.Equal(t, int64(1), found.TargetVersion)
		keys, err := repo.ListProbeKeys(ctx)
		require.NoError(t, err)
		registered, ok := keys[key.ID]
		require.True(t, ok)
		require.WithinDuration(t, time.Now(), registered, time.Minute)
		// Disabled targets drop out of the index; the probe key stays registered forever.
		set, err := f.repo.GetTargets(ctx, f.group.ID)
		require.NoError(t, err)
		set.Items = nil
		_, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.NoError(t, err)
		targets, err = repo.ListEnabledTargets(ctx)
		require.NoError(t, err)
		for _, target := range targets {
			require.NotEqual(t, f.group.ID, target.GroupID)
		}
		keys, err = repo.ListProbeKeys(ctx)
		require.NoError(t, err)
		_, ok = keys[key.ID]
		require.True(t, ok)
	})
}
