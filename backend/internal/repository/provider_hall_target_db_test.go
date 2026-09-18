//go:build integration || providerhall_localdb

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

type providerHallTargetFixture struct {
	db       *sql.DB
	client   *ent.Client
	repo     service.ProviderHallRepository
	svc      *service.ProviderHallService
	user     *ent.User
	group    *ent.Group
	profiles []*service.ProviderHallProfile
}

func newProviderHallTargetFixture(t *testing.T, db *sql.DB) *providerHallTargetFixture {
	t.Helper()
	ctx := context.Background()
	f := &providerHallTargetFixture{db: db, client: ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))}
	f.repo = NewProviderHallRepository(f.client)
	f.svc = service.NewProviderHallService(f.repo, nil, nil, nil)
	var err error
	f.user, err = f.client.User.Create().SetEmail(fmt.Sprintf("hall-target-%d@example.test", time.Now().UnixNano())).SetPasswordHash("test-only").Save(ctx)
	require.NoError(t, err)
	f.group, err = f.client.Group.Create().SetName(fmt.Sprintf("hall-target-%d", time.Now().UnixNano())).SetPlatform(service.PlatformOpenAI).Save(ctx)
	require.NoError(t, err)
	for i := range 2 {
		p, err := f.repo.SaveProfile(ctx, service.ProviderHallProfile{Model: fmt.Sprintf("gpt-hall-target-%d-%d", f.group.ID, i), Protocol: "responses", OutputLimit: 256, ModelAliases: []string{}})
		require.NoError(t, err)
		f.profiles = append(f.profiles, p)
	}
	original, err := f.repo.GetConfig(ctx)
	require.NoError(t, err)
	cfg := *original
	cfg.OperatorUserID, cfg.UpdatedBy = &f.user.ID, &f.user.ID
	_, err = f.repo.UpdateConfig(ctx, cfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		current, err := f.repo.GetConfig(ctx)
		require.NoError(t, err)
		original.Version = current.Version
		_, err = f.repo.UpdateConfig(ctx, *original)
		require.NoError(t, err)
		for _, query := range []string{
			`DELETE FROM ops_error_logs WHERE group_id=$1`,
			`DELETE FROM groups WHERE id=$1`,
			`DELETE FROM provider_hall_probe_keys WHERE group_id=$1`,
		} {
			_, err = db.ExecContext(ctx, query, f.group.ID)
			require.NoError(t, err)
		}
		for _, p := range f.profiles {
			_, err = db.ExecContext(ctx, `DELETE FROM provider_hall_profiles WHERE id=$1`, p.ID)
			require.NoError(t, err)
		}
		_, err = db.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, f.user.ID)
		require.NoError(t, err)
	})
	return f
}

func (f *providerHallTargetFixture) key(t *testing.T) *ent.APIKey {
	t.Helper()
	k, err := f.client.APIKey.Create().SetUserID(f.user.ID).SetGroupID(f.group.ID).
		SetName("hall-test").SetKey(fmt.Sprintf("hall-test-%d", time.Now().UnixNano())).Save(context.Background())
	require.NoError(t, err)
	return k
}

func (f *providerHallTargetFixture) input(keyID int64) service.ProviderHallTargetSet {
	return service.ProviderHallTargetSet{GroupID: f.group.ID, Items: []service.ProviderHallTarget{{ProfileID: f.profiles[0].ID, ProbeKeyID: &keyID, Enabled: true}}}
}

func providerHallTargetDatabaseContracts(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	t.Run("target_atomic_registration_and_cas", func(t *testing.T) {
		f := newProviderHallTargetFixture(t, db)
		good, used := f.key(t), f.key(t)
		_, err := f.client.APIKey.UpdateOneID(used.ID).SetLastUsedAt(time.Now()).Save(ctx)
		require.NoError(t, err)
		input := f.input(good.ID)
		input.Items = append(input.Items, service.ProviderHallTarget{ProfileID: f.profiles[1].ID, ProbeKeyID: &used.ID, Enabled: true})
		_, err = f.svc.SaveTargets(ctx, input, f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallProbeKey)
		set, err := f.repo.GetTargets(ctx, f.group.ID)
		require.NoError(t, err)
		require.Zero(t, set.Version)
		require.Empty(t, set.Items)
		probe, err := f.repo.IsProbeKey(ctx, good.ID, time.Now())
		require.NoError(t, err)
		require.False(t, probe, "an invalid second target must roll back the first registration")
		var wg sync.WaitGroup
		results := make(chan error, 2)
		for range 2 {
			wg.Go(func() { _, err := f.svc.SaveTargets(ctx, f.input(good.ID), f.user.ID); results <- err })
		}
		wg.Wait()
		close(results)
		successes, conflicts := 0, 0
		for err := range results {
			if err == nil {
				successes++
			} else {
				require.ErrorIs(t, err, service.ErrProviderHallConflict)
				conflicts++
			}
		}
		require.Equal(t, 1, successes)
		require.Equal(t, 1, conflicts)
		set, err = f.repo.GetTargets(ctx, f.group.ID)
		require.NoError(t, err)
		require.Equal(t, int64(1), set.Version)
		require.Len(t, set.Items, 1)
		_, err = db.ExecContext(ctx, `INSERT INTO provider_hall_targets(group_id,profile_id) VALUES($1,$2)`, f.group.ID, f.profiles[0].ID)
		require.Error(t, err, "group/profile uniqueness must also hold in SQL")
		_, err = db.ExecContext(ctx, `UPDATE provider_hall_targets SET probe_interval_seconds=59 WHERE group_id=$1`, f.group.ID)
		require.Error(t, err)
	})

	t.Run("target_stop_reenable_and_permanent_provenance", func(t *testing.T) {
		f := newProviderHallTargetFixture(t, db)
		key := f.key(t)
		before := time.Now().Add(-time.Second)
		set, err := f.svc.SaveTargets(ctx, f.input(key.ID), f.user.ID)
		require.NoError(t, err)
		targetID := set.Items[0].ID
		probe, err := f.repo.IsProbeKey(ctx, key.ID, before)
		require.NoError(t, err)
		require.False(t, probe)
		registration, err := f.client.ProviderHallProbeKey.Get(ctx, key.ID)
		require.NoError(t, err)
		probe, err = f.repo.IsProbeKey(ctx, key.ID, registration.RegisteredAt)
		require.NoError(t, err)
		require.True(t, probe)
		_, err = f.client.APIKey.UpdateOneID(key.ID).SetExpiresAt(time.Now().Add(-time.Minute)).Save(ctx)
		require.NoError(t, err)
		set.Items = nil
		set, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.NoError(t, err, "stopping a target must not require a healthy key")
		require.False(t, set.Items[0].Enabled)
		require.Equal(t, targetID, set.Items[0].ID)
		set.Items[0].Enabled = true
		_, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallProbeKey)
		_, err = f.client.APIKey.UpdateOneID(key.ID).ClearExpiresAt().SetLastUsedAt(time.Now()).SetQuotaUsed(1).Save(ctx)
		require.NoError(t, err)
		set, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.NoError(t, err, "registered key reuse is allowed for the same owner and group")
		require.Equal(t, targetID, set.Items[0].ID)
		cfg, err := f.repo.GetConfig(ctx)
		require.NoError(t, err)
		cfg.OperatorUserID = nil
		_, err = f.repo.UpdateConfig(ctx, *cfg)
		require.NoError(t, err)
		current, err := f.repo.GetTargets(ctx, f.group.ID)
		require.NoError(t, err)
		require.False(t, current.Items[0].Enabled)
		require.Equal(t, set.Version+1, current.Version)
		_, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallConflict)
		_, err = f.client.Group.UpdateOneID(f.group.ID).SetDeletedAt(time.Now()).Save(ctx)
		require.NoError(t, err)
		_, err = f.repo.GetTargets(ctx, f.group.ID)
		require.ErrorIs(t, err, service.ErrProviderHallNotFound)
		_, err = f.svc.SaveTargets(ctx, *current, f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallNotFound)
		_, err = db.ExecContext(ctx, `DELETE FROM groups WHERE id=$1`, f.group.ID)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, f.user.ID)
		require.NoError(t, err)
		probe, err = f.repo.IsProbeKey(ctx, key.ID, time.Now())
		require.NoError(t, err)
		require.True(t, probe, "key, user and group deletion must not erase provenance")
	})

	t.Run("target_binding_permissions", func(t *testing.T) {
		f := newProviderHallTargetFixture(t, db)
		key := f.key(t)
		input := f.input(key.ID)
		_, err := f.client.APIKey.UpdateOneID(key.ID).ClearGroupID().Save(ctx)
		require.NoError(t, err)
		_, err = f.svc.SaveTargets(ctx, input, f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallProbeKey)
		_, err = f.client.APIKey.UpdateOneID(key.ID).SetGroupID(f.group.ID).SetStatus(service.StatusAPIKeyDisabled).Save(ctx)
		require.NoError(t, err)
		_, err = f.svc.SaveTargets(ctx, input, f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallProbeKey)
		_, err = f.client.APIKey.UpdateOneID(key.ID).SetStatus(service.StatusActive).Save(ctx)
		require.NoError(t, err)
		_, err = f.client.Group.UpdateOneID(f.group.ID).SetIsExclusive(true).Save(ctx)
		require.NoError(t, err)
		_, err = f.svc.SaveTargets(ctx, input, f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallOperator)
		_, err = f.client.User.UpdateOneID(f.user.ID).AddAllowedGroupIDs(f.group.ID).Save(ctx)
		require.NoError(t, err)
		set, err := f.svc.SaveTargets(ctx, input, f.user.ID)
		require.NoError(t, err)
		_, err = f.client.Group.UpdateOneID(f.group.ID).SetSubscriptionType(service.SubscriptionTypeSubscription).Save(ctx)
		require.NoError(t, err)
		_, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallOperator)
		sub, err := f.client.UserSubscription.Create().SetUserID(f.user.ID).SetGroupID(f.group.ID).SetStartsAt(time.Now().Add(-time.Hour)).SetExpiresAt(time.Now().Add(time.Hour)).Save(ctx)
		require.NoError(t, err)
		set, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.NoError(t, err)
		_, err = f.client.UserSubscription.UpdateOneID(sub.ID).SetExpiresAt(time.Now().Add(-time.Second)).Save(ctx)
		require.NoError(t, err)
		_, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallOperator)
	})

	t.Run("target_composite_route_recheck", func(t *testing.T) {
		f := newProviderHallTargetFixture(t, db)
		_, err := f.client.Group.UpdateOneID(f.group.ID).SetPlatform(service.PlatformComposite).Save(ctx)
		require.NoError(t, err)
		route, err := f.client.CompositeModelRoute.Create().SetGroupID(f.group.ID).SetPublicModel(f.profiles[0].Model).
			SetTargetPlatform(service.PlatformDeepseek).SetEndpoint("responses").Save(ctx)
		require.NoError(t, err)
		input := f.input(f.key(t).ID)
		_, err = f.svc.SaveTargets(ctx, input, f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallTargetRoute)
		_, err = f.client.CompositeModelRoute.UpdateOneID(route.ID).SetTargetPlatform(service.PlatformOpenAI).Save(ctx)
		require.NoError(t, err)
		_, err = f.svc.SaveTargets(ctx, input, f.user.ID)
		require.NoError(t, err)
	})

	t.Run("target_probe_excluded_from_v2", func(t *testing.T) {
		f := newProviderHallTargetFixture(t, db)
		probeKey, ordinaryKey := f.key(t), f.key(t)
		set, err := f.svc.SaveTargets(ctx, f.input(probeKey.ID), f.user.ID)
		require.NoError(t, err)
		registration, err := f.client.ProviderHallProbeKey.Get(ctx, probeKey.ID)
		require.NoError(t, err)
		account, err := f.client.Account.Create().SetName("hall-v2-test").SetPlatform(service.PlatformOpenAI).SetType(service.AccountTypeAPIKey).Save(ctx)
		require.NoError(t, err)
		t.Cleanup(func() { _, _ = db.ExecContext(ctx, `DELETE FROM accounts WHERE id=$1`, account.ID) })
		before, after := registration.RegisteredAt.Add(-time.Second), registration.RegisteredAt.Add(time.Second)
		for i, sample := range []struct {
			key    int64
			at     time.Time
			tokens int
		}{
			{probeKey.ID, before, 100}, {probeKey.ID, after, 9000}, {ordinaryKey.ID, after, 200},
		} {
			requestID := fmt.Sprintf("hall-v2-%d-%d", f.group.ID, i)
			_, err = f.client.UsageLog.Create().SetUserID(f.user.ID).SetAPIKeyID(sample.key).SetAccountID(account.ID).
				SetRequestID(requestID).SetModel(f.profiles[0].Model).SetGroupID(f.group.ID).SetInputTokens(sample.tokens).
				SetActualCost(0.01).SetFirstTokenMs(100).SetDurationMs(200).SetCreatedAt(sample.at).Save(ctx)
			require.NoError(t, err)
			_, err = db.ExecContext(ctx, `INSERT INTO ops_error_logs(request_id,user_id,api_key_id,group_id,platform,model,error_phase,error_type,error_owner,status_code,created_at) VALUES($1,$2,$3,$4,'openai',$5,'upstream','upstream','provider',500,$6)`, requestID+"-err", f.user.ID, sample.key, f.group.ID, f.profiles[0].Model, sample.at)
			require.NoError(t, err)
		}
		// Stopping the target cannot put already recorded probes back into V2.
		set.Items = nil
		_, err = f.svc.SaveTargets(ctx, *set, f.user.ID)
		require.NoError(t, err)
		v2 := NewChannelMonitorV2Repository(db)
		for range 2 {
			require.NoError(t, v2.RecomputeRange(ctx, before.Truncate(time.Minute), after.Truncate(time.Minute).Add(time.Minute)))
			for _, table := range []string{"channel_monitor_v2_metrics_1m", "channel_monitor_v2_user_metrics_1m"} {
				var success, failures, tokens int
				require.NoError(t, db.QueryRowContext(ctx, `SELECT sum(success_requests),sum(error_requests),sum(input_tokens) FROM `+table+` WHERE group_id=$1`, f.group.ID).Scan(&success, &failures, &tokens))
				require.Equal(t, 2, success)
				require.Equal(t, 2, failures)
				require.Equal(t, 300, tokens)
			}
			var count int
			require.NoError(t, db.QueryRowContext(ctx, `SELECT sum(sample_count) FROM channel_monitor_v2_latency_histograms_1m WHERE group_id=$1 AND user_id=0 AND metric='ttft'`, f.group.ID).Scan(&count))
			require.Equal(t, 2, count)
			require.NoError(t, db.QueryRowContext(ctx, `SELECT sum(error_requests) FROM channel_monitor_v2_error_metrics_1m WHERE group_id=$1`, f.group.ID).Scan(&count))
			require.Equal(t, 2, count)
		}
	})
}
