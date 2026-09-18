//go:build unit

package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestProvideOpenAICodexVersionSyncServiceHonorsInstanceRole(t *testing.T) {
	tests := []struct {
		name       string
		role       string
		wantWrites []string
	}{
		{name: "api does not start", role: config.InstanceRoleAPI},
		{name: "primary starts", role: config.InstanceRolePrimary, wantWrites: []string{"0.146.0"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newCodexVersionSyncSettingRepoStub(nil)
			github := &codexVersionSyncGitHubStub{
				latest: &GitHubRelease{TagName: "rust-v0.146.0"},
			}
			cfg := &config.Config{InstanceRole: tt.role}

			svc := ProvideOpenAICodexVersionSyncService(repo, &SettingService{}, github, cfg)
			svc.Stop()

			require.Equal(t, tt.wantWrites, repo.syncedWrites())
		})
	}
}

func TestProvideOllamaCloudUsageServiceHonorsInstanceRole(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		wantStarted bool
	}{
		{name: "api does not start", role: config.InstanceRoleAPI},
		{name: "primary starts", role: config.InstanceRolePrimary, wantStarted: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{InstanceRole: tt.role}
			svc := ProvideOllamaCloudUsageService(nil, nil, nil, nil, cfg, nil, nil)
			svc.Stop()

			svc.mu.Lock()
			started := svc.started
			svc.mu.Unlock()
			require.Equal(t, tt.wantStarted, started)
		})
	}
}
