//go:build integration || providerhall_localdb

package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func init() { registerProviderHallLocalDBContract("admin_ux", providerHallManagementContracts) }

func providerHallManagementContracts(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	t.Run("admin_settings_preflight_atomic_save_and_unchanged_versions", func(t *testing.T) {
		f := newProviderHallTargetFixture(t, db)
		key := f.key(t)
		input := f.input(key.ID)
		input.Listing = &service.ProviderHallGroup{DisplayName: "Draft", Listed: true}
		input.Preflight = true
		_, err := f.svc.SaveTargets(ctx, input, f.user.ID)
		require.NoError(t, err)
		before, err := f.repo.GetTargets(ctx, f.group.ID)
		require.NoError(t, err)
		require.Zero(t, before.Version)
		require.Empty(t, before.Items)
		registered, err := f.repo.IsProbeKey(ctx, key.ID, time.Now())
		require.NoError(t, err)
		require.False(t, registered)
		input.Preflight = false
		saved, err := f.svc.SaveTargets(ctx, input, f.user.ID)
		require.NoError(t, err)
		require.True(t, saved.Listing.Listed)
		require.Equal(t, "Draft", saved.Listing.DisplayName)
		input.Version = saved.Version
		input.Listing.DisplayName = "Renamed"
		saved2, err := f.svc.SaveTargets(ctx, input, f.user.ID)
		require.NoError(t, err)
		require.Equal(t, saved.Items[0].Version, saved2.Items[0].Version)
		require.False(t, *saved2.Items[0].AutoScheduleEnabled, "new targets are manual-only")
		auto := true
		input.Version = saved2.Version
		input.Items[0].AutoScheduleEnabled = &auto
		saved3, err := f.svc.SaveTargets(ctx, input, f.user.ID)
		require.NoError(t, err)
		input.Version = saved3.Version
		input.Items[0].AutoScheduleEnabled = nil
		saved4, err := f.svc.SaveTargets(ctx, input, f.user.ID)
		require.NoError(t, err)
		require.True(t, *saved4.Items[0].AutoScheduleEnabled, "old clients preserve scheduling")
		input.Version = saved4.Version
		input.Listing.DisplayName = "Must roll back"
		input.Items[0].ProbeKeyID = new(int64(999999999))
		_, err = f.svc.SaveTargets(ctx, input, f.user.ID)
		require.Error(t, err)
		listing, err := f.repo.GetGroup(ctx, f.group.ID)
		require.NoError(t, err)
		require.Equal(t, "Renamed", listing.DisplayName)
	})
	t.Run("dedicated_key_retries_are_serialized_and_never_expose_secrets", func(t *testing.T) {
		f := newProviderHallTargetFixture(t, db)
		management := f.repo.(service.ProviderHallManagementRepository)
		var wg sync.WaitGroup
		ids := make(chan int64, 8)
		errs := make(chan error, 8)
		for range 8 {
			wg.Go(func() {
				key, err := management.EnsureProbeKey(ctx, f.group.ID, f.profiles[0].ID, f.user.ID)
				errs <- err
				if key != nil {
					ids <- key.ID
					raw, e := json.Marshal(key)
					if e == nil {
						require.NotContains(t, string(raw), "sk-")
					}
				}
			})
		}
		wg.Wait()
		close(ids)
		close(errs)
		for err := range errs {
			require.NoError(t, err)
		}
		unique := map[int64]bool{}
		for id := range ids {
			unique[id] = true
		}
		require.Len(t, unique, 1)
		keys, err := management.ListProbeKeys(ctx, f.group.ID)
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, "registered", keys[0].Status)
		page, err := management.ListAdminGroups(ctx, service.ProviderHallGroupFilter{Search: f.group.Name, Page: 1, PageSize: 20})
		require.NoError(t, err)
		require.Equal(t, 1, page.Total)
	})
	t.Run("dedicated_key_reuse_skips_changed_binding", func(t *testing.T) {
		f := newProviderHallTargetFixture(t, db)
		management := f.repo.(service.ProviderHallManagementRepository)
		first, err := management.EnsureProbeKey(ctx, f.group.ID, f.profiles[0].ID, f.user.ID)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, `UPDATE api_keys SET group_id=NULL WHERE id=$1`, first.ID)
		require.NoError(t, err)
		replacement, err := management.EnsureProbeKey(ctx, f.group.ID, f.profiles[0].ID, f.user.ID)
		require.NoError(t, err)
		require.NotEqual(t, first.ID, replacement.ID)
		retry, err := management.EnsureProbeKey(ctx, f.group.ID, f.profiles[0].ID, f.user.ID)
		require.NoError(t, err)
		require.Equal(t, replacement.ID, retry.ID)
	})
	t.Run("manual_unlisted_survives_auto_disable_and_scheduled_jobs_stop", func(t *testing.T) {
		f := newProviderHallJobFixture(t, db)
		f.blockSlots(t)
		manual := f.enqueue(t, service.ProviderHallJobProbe, "")
		_, err := db.ExecContext(ctx, `UPDATE provider_hall_groups SET listed=false WHERE group_id=$1`, f.group.ID)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, `UPDATE provider_hall_config SET auto_schedule_enabled=false`)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, `UPDATE provider_hall_targets SET auto_schedule_enabled=false WHERE id=$1`, f.target.ID)
		require.NoError(t, err)
		r := f.runner(t)
		r.RunOnce(ctx)
		r.WaitInflight()
		require.Equal(t, service.ProviderHallJobSucceeded, f.job(t, manual.ID).Status)
		require.EqualValues(t, 1, f.gateway.hits.Load())
		targets, err := f.jobs.ListSchedulableTargets(ctx)
		require.NoError(t, err)
		for _, target := range targets {
			if target.TargetID == f.target.ID {
				require.True(t, target.Enabled)
				require.False(t, target.Listed)
				require.False(t, target.AutoScheduleEnabled)
			}
		}
		_, _, err = r.EnqueueManual(ctx, service.ProviderHallJobProbe, f.group.ID, f.profiles[0].ID, "manual-unlisted", f.user.ID)
		require.NoError(t, err)
		_, _, err = r.EnqueueManual(service.ProviderHallWithExpectedVersions(ctx, 999, 999), service.ProviderHallJobProbe, f.group.ID, f.profiles[0].ID, "conflict", f.user.ID)
		require.ErrorIs(t, err, service.ErrProviderHallConflict)
	})
	t.Run("scheduled_admission_requires_all_auto_switches_and_listing", func(t *testing.T) {
		f := newProviderHallJobFixture(t, db)
		snapshot := f.snapshot(t)
		for _, global := range []bool{false, true} {
			for _, target := range []bool{false, true} {
				for _, listed := range []bool{false, true} {
					_, err := db.ExecContext(ctx, `UPDATE provider_hall_config SET auto_schedule_enabled=$1`, global)
					require.NoError(t, err)
					_, err = db.ExecContext(ctx, `UPDATE provider_hall_targets SET auto_schedule_enabled=$1 WHERE id=$2`, target, f.target.ID)
					require.NoError(t, err)
					_, err = db.ExecContext(ctx, `UPDATE provider_hall_groups SET listed=$1 WHERE group_id=$2`, listed, f.group.ID)
					require.NoError(t, err)
					now := time.Now()
					_, created, err := f.jobs.EnqueueSlot(ctx, service.ProviderHallEnqueueInput{Kind: service.ProviderHallJobProbe, TargetID: f.target.ID, GroupID: f.group.ID, ProfileID: f.profiles[0].ID, Snapshot: snapshot, SlotAt: &now})
					if global && target && listed {
						require.NoError(t, err)
						require.True(t, created)
					} else {
						require.ErrorIs(t, err, service.ErrProviderHallTargetDisabled)
					}
				}
			}
		}
	})
}
