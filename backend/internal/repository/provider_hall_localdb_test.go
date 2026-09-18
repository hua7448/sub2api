//go:build providerhall_localdb && !integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/Wei-Shaw/sub2api/migrations"
	"github.com/stretchr/testify/require"
)

// This opt-in suite owns an isolated Unix-socket-only PostgreSQL cluster. It
// does not replace the integration tag's Docker/PostgreSQL/Redis harness.
func TestProviderHallLocalDatabase(t *testing.T) {
	bin := os.Getenv("PROVIDER_HALL_PG_BIN")
	require.NotEmpty(t, bin, "set PROVIDER_HALL_PG_BIN to the PostgreSQL bin directory")
	dir, err := os.MkdirTemp("/tmp", "provider-hall-")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	data := filepath.Join(dir, "data")
	run := func(name string, args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		out, err := exec.CommandContext(ctx, filepath.Join(bin, name), args...).CombinedOutput()
		require.NoError(t, err, "%s: %s", name, out)
	}
	// Several suites may run at once on one machine; PROVIDER_HALL_PG_PORT keeps
	// their Unix-socket clusters apart.
	port := os.Getenv("PROVIDER_HALL_PG_PORT")
	if port == "" {
		port = "15479"
	}
	run("initdb", "-D", data, "-U", "hall_test", "--auth=trust", "--no-locale", "-E", "UTF8")
	run("pg_ctl", "-D", data, "-l", filepath.Join(dir, "postgres.log"), "-o", fmt.Sprintf("-h '' -k %s -p %s", dir, port), "-w", "start")
	t.Cleanup(func() { run("pg_ctl", "-D", data, "-m", "immediate", "-w", "stop") })
	connect := func(name string) *sql.DB {
		db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=hall_test dbname=%s sslmode=disable", dir, port, name))
		require.NoError(t, err)
		t.Cleanup(func() { _ = db.Close() })
		return db
	}
	admin := connect("postgres")
	// Every provider hall migration (229 onwards) must be applied exactly once
	// on each path. The list is discovered so later batches are covered as soon
	// as their migration file exists.
	allFiles, err := fs.Glob(migrations.FS, "*.sql")
	require.NoError(t, err)
	var hallMigrations []string
	for _, name := range allFiles {
		if name >= "229_" && strings.Contains(name, "provider_hall") {
			hallMigrations = append(hallMigrations, name)
		}
	}
	sort.Strings(hallMigrations)
	require.NotEmpty(t, hallMigrations)
	// upgrade_NNN starts from a database where every migration below NNN+1 is
	// already applied, i.e. the state of a deployment that stopped at NNN.
	modes := []string{"empty", "upgrade_228"}
	for _, name := range hallMigrations {
		modes = append(modes, "upgrade_"+name[:3])
	}
	modes = modes[:len(modes)-1] // the newest migration has no "upgrade from itself" path
	if os.Getenv("PROVIDER_HALL_LOCALDB_MODES") != "" {
		modes = strings.Split(os.Getenv("PROVIDER_HALL_LOCALDB_MODES"), ",")
	}
	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			_, err := admin.Exec("CREATE DATABASE hall_" + mode)
			require.NoError(t, err)
			db := connect("hall_" + mode)
			ctx := context.Background()
			var upgradeGroupID, upgradeProfileID int64
			if mode != "empty" {
				cutoff := fmt.Sprintf("%03d_", atoiMigration(t, strings.TrimPrefix(mode, "upgrade_"))+1)
				previous := fstest.MapFS{}
				files, err := fs.Glob(migrations.FS, "*.sql")
				require.NoError(t, err)
				for _, name := range files {
					if name >= cutoff {
						continue
					}
					content, err := migrations.FS.ReadFile(name)
					require.NoError(t, err)
					previous[name] = &fstest.MapFile{Data: content}
				}
				require.NoError(t, applyMigrationsFS(ctx, db, previous))
				if mode == "upgrade_233" {
					_, err := db.ExecContext(ctx, `UPDATE provider_hall_config SET tasks_enabled=true WHERE id=1`)
					require.NoError(t, err)
					require.NoError(t, db.QueryRowContext(ctx, `INSERT INTO groups (name,platform) VALUES ('migration schedule','openai') RETURNING id`).Scan(&upgradeGroupID))
					_, err = db.ExecContext(ctx, `INSERT INTO provider_hall_groups (group_id) VALUES ($1)`, upgradeGroupID)
					require.NoError(t, err)
					require.NoError(t, db.QueryRowContext(ctx, `INSERT INTO provider_hall_profiles (model,protocol) VALUES ('migration-test','responses') RETURNING id`).Scan(&upgradeProfileID))
					_, err = db.ExecContext(ctx, `INSERT INTO provider_hall_targets (group_id,profile_id) VALUES ($1,$2)`, upgradeGroupID, upgradeProfileID)
					require.NoError(t, err)
				}
			}
			require.NoError(t, ApplyMigrations(ctx, db))
			require.NoError(t, ApplyMigrations(ctx, db), "restart must not apply a migration twice")
			if mode == "upgrade_233" {
				var globalAuto, targetAuto bool
				require.NoError(t, db.QueryRowContext(ctx, `SELECT auto_schedule_enabled FROM provider_hall_config WHERE id=1`).Scan(&globalAuto))
				require.True(t, globalAuto)
				require.NoError(t, db.QueryRowContext(ctx, `SELECT auto_schedule_enabled FROM provider_hall_targets WHERE group_id=$1`, upgradeGroupID).Scan(&targetAuto))
				require.True(t, targetAuto)
				_, err = db.ExecContext(ctx, `DELETE FROM groups WHERE id=$1`, upgradeGroupID)
				require.NoError(t, err)
				_, err = db.ExecContext(ctx, `DELETE FROM provider_hall_profiles WHERE id=$1`, upgradeProfileID)
				require.NoError(t, err)
				_, err = db.ExecContext(ctx, `UPDATE provider_hall_config SET tasks_enabled=false,auto_schedule_enabled=false`)
				require.NoError(t, err)
			}
			for _, name := range hallMigrations {
				var count int
				require.NoError(t, db.QueryRow(`SELECT count(*) FROM schema_migrations WHERE filename=$1`, name).Scan(&count))
				require.Equal(t, 1, count, name)
			}
			providerHallDatabaseContracts(t, db)
			providerHallTargetDatabaseContracts(t, db)
			for _, contract := range providerHallLocalDBContracts {
				contract.fn(t, db)
			}
		})
	}
}

func atoiMigration(t *testing.T, s string) int {
	t.Helper()
	n := 0
	for _, r := range s {
		require.True(t, r >= '0' && r <= '9', "bad migration prefix %q", s)
		n = n*10 + int(r-'0')
	}
	return n
}
