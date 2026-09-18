//go:build integration || providerhall_localdb

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/providerhallprofile"
	_ "github.com/Wei-Shaw/sub2api/ent/runtime"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func providerHallDatabaseContracts(t *testing.T, db *sql.DB) {
	ctx := context.Background()
	client := ent.NewClient(ent.Driver(entsql.OpenDB(dialect.Postgres, db)))
	repo := NewProviderHallRepository(client)
	cfg, err := repo.GetConfig(ctx)
	require.NoError(t, err)
	require.False(t, cfg.CollectionEnabled)
	require.False(t, cfg.DisplayEnabled)
	require.False(t, cfg.TasksEnabled)
	require.Equal(t, "0.00000000", cfg.DailyBudget)
	require.Equal(t, "responses", cfg.DefaultProtocol)
	require.Equal(t, []string{}, cfg.ExpectedNodes)
	_, err = db.ExecContext(ctx, `INSERT INTO provider_hall_config(id) VALUES(2)`)
	require.Error(t, err)

	// A real CAS race: exactly one transaction may consume a config version.
	var wg sync.WaitGroup
	errors := make(chan error, 2)
	for range 2 {
		wg.Go(func() { _, err := repo.UpdateConfig(ctx, *cfg); errors <- err })
	}
	wg.Wait()
	close(errors)
	successes, conflicts := 0, 0
	for err := range errors {
		if err == nil {
			successes++
		} else {
			require.ErrorIs(t, err, service.ErrProviderHallConflict)
			conflicts++
		}
	}
	require.Equal(t, 1, successes)
	require.Equal(t, 1, conflicts)

	u, err := client.User.Create().SetEmail(fmt.Sprintf("hall-%d@example.test", time.Now().UnixNano())).SetPasswordHash("test-only").Save(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM users WHERE id=$1`, u.ID) })
	g, err := client.Group.Create().SetName(fmt.Sprintf("hall-%d", time.Now().UnixNano())).SetPlatform(service.PlatformOpenAI).Save(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM groups WHERE id=$1`, g.ID) })
	listing, err := repo.GetGroup(ctx, g.ID)
	require.NoError(t, err)
	require.Zero(t, listing.Version)
	require.False(t, listing.Listed)
	listing.Listed, listing.UpdatedBy = true, &u.ID
	listing.DisplayName = strings.Repeat("\u4e2d", 34)
	listing.Description = strings.Repeat("\u6587", 667)
	created, err := repo.SaveGroup(ctx, *listing)
	require.NoError(t, err)
	require.Equal(t, int64(1), created.Version)
	require.Equal(t, listing.DisplayName, created.DisplayName)
	require.Equal(t, listing.Description, created.Description)
	_, err = repo.SaveGroup(ctx, *listing)
	require.ErrorIs(t, err, service.ErrProviderHallConflict)
	created.DisplayName = strings.Repeat("\u4e2d", 100)
	created.Description = strings.Repeat("\u6587", 2000)
	updated, err := repo.SaveGroup(ctx, *created)
	require.NoError(t, err)
	require.Equal(t, int64(2), updated.Version)
	stored, err := repo.GetGroup(ctx, g.ID)
	require.NoError(t, err)
	require.Equal(t, created.DisplayName, stored.DisplayName)
	require.Equal(t, created.Description, stored.Description)
	for _, field := range []string{"name", "description"} {
		invalid := *updated
		if field == "name" {
			invalid.DisplayName += "\u4e2d"
		} else {
			invalid.Description += "\u6587"
		}
		_, err = repo.SaveGroup(ctx, invalid)
		require.Error(t, err, "Ent must reject values exceeding the character limit")
	}
	stored, err = repo.GetGroup(ctx, g.ID)
	require.NoError(t, err)
	require.Equal(t, updated.Version, stored.Version, "invalid writes must not consume a config version")
	_, err = client.Group.UpdateOneID(g.ID).SetDeletedAt(time.Now()).Save(ctx)
	require.NoError(t, err)
	_, err = repo.GetGroup(ctx, g.ID)
	require.ErrorIs(t, err, service.ErrProviderHallNotFound)
	_, err = repo.SaveGroup(ctx, *updated)
	require.ErrorIs(t, err, service.ErrProviderHallNotFound)

	price, zero, rate := "12345678901234.1234567890", "0.0000000000", "0.1234567890"
	confirmed := time.Now().Add(-time.Hour).UTC().Truncate(time.Microsecond)
	p := service.ProviderHallProfile{
		ProviderHallVersion: service.ProviderHallVersion{UpdatedBy: &u.ID},
		Model:               fmt.Sprintf("gpt-hall-%d", time.Now().UnixNano()), Protocol: "responses", OutputLimit: 256,
		ModelAliases: []string{"gpt-test-version"}, ReferenceInputPrice: &price, ReferenceCachePrice: &zero,
		ReferenceCacheRate: &rate, ReferenceConfirmedAt: &confirmed,
	}
	profile, err := repo.SaveProfile(ctx, p)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM provider_hall_profiles WHERE id=$1`, profile.ID) })
	require.Equal(t, price, *profile.ReferenceInputPrice)
	require.Equal(t, zero, *profile.ReferenceCachePrice)
	require.Equal(t, rate, *profile.ReferenceCacheRate)
	require.Equal(t, confirmed, *profile.ReferenceConfirmedAt)
	_, err = repo.SaveProfile(ctx, p)
	require.ErrorIs(t, err, service.ErrProviderHallProfileExists)
	profile.ReferenceInputPrice, profile.ReferenceCachePrice, profile.ReferenceCacheRate, profile.ReferenceConfirmedAt = nil, nil, nil, nil
	cleared, err := repo.SaveProfile(ctx, *profile)
	require.NoError(t, err)
	require.Nil(t, cleared.ReferenceInputPrice)
	require.Nil(t, cleared.ReferenceConfirmedAt)
	_, err = repo.SaveProfile(ctx, *profile)
	require.ErrorIs(t, err, service.ErrProviderHallConflict)
	_, err = db.ExecContext(ctx, `UPDATE provider_hall_profiles SET reference_input_price=1 WHERE id=$1`, profile.ID)
	require.Error(t, err, "partial reference must be rejected by the database")
	profiles, err := repo.ListProfiles(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, profiles)

	cfg, err = repo.GetConfig(ctx)
	require.NoError(t, err)
	cfg.DefaultModel = cleared.Model
	cfg, err = repo.UpdateConfig(ctx, *cfg)
	require.NoError(t, err)
	cleared.Model += "-renamed"
	_, err = repo.SaveProfile(ctx, *cleared)
	require.ErrorIs(t, err, service.ErrProviderHallConflict)
	cfg.DefaultModel = ""
	_, err = repo.UpdateConfig(ctx, *cfg)
	require.NoError(t, err)
	_, err = client.ProviderHallProfile.Delete().Where(providerhallprofile.IDEQ(profile.ID)).Exec(ctx)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `DELETE FROM groups WHERE id=$1`, g.ID)
	require.NoError(t, err)
	var count int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM provider_hall_groups WHERE group_id=$1`, g.ID).Scan(&count))
	require.Zero(t, count, "hard deletion must cascade the hall listing")
}
