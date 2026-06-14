//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type updateServiceCacheStub struct {
	data string
}

func (s *updateServiceCacheStub) GetUpdateInfo(context.Context) (string, error) {
	if s.data == "" {
		return "", errors.New("cache miss")
	}
	return s.data, nil
}

func (s *updateServiceCacheStub) SetUpdateInfo(_ context.Context, data string, _ time.Duration) error {
	s.data = data
	return nil
}

type updateServiceGitHubClientStub struct {
	release *GitHubRelease
	repo    string
}

func (s *updateServiceGitHubClientStub) FetchLatestRelease(_ context.Context, repo string) (*GitHubRelease, error) {
	s.repo = repo
	return s.release, nil
}

func (s *updateServiceGitHubClientStub) DownloadFile(context.Context, string, string, int64) error {
	panic("DownloadFile should not be called when no update is available")
}

func (s *updateServiceGitHubClientStub) FetchChecksumFile(context.Context, string) ([]byte, error) {
	panic("FetchChecksumFile should not be called when no update is available")
}

func TestUpdateServicePerformUpdateNoUpdateReturnsSentinel(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{
			release: &GitHubRelease{
				TagName: "v0.1.132-smartapi.1",
				Name:    "v0.1.132-smartapi.1",
			},
		},
		"0.1.132-smartapi.1",
		"release",
	)

	err := svc.PerformUpdate(context.Background())

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrNoUpdateAvailable))
	require.ErrorIs(t, err, ErrNoUpdateAvailable)
}

func TestUpdateServiceUsesSmartAPIForkRepoAndComparesPatch(t *testing.T) {
	client := &updateServiceGitHubClientStub{
		release: &GitHubRelease{
			TagName: "v0.1.136-smartapi.2",
			Name:    "v0.1.136-smartapi.2",
		},
	}
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		client,
		"0.1.136-smartapi.1",
		"release",
	)

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.Equal(t, "hua7448/sub2api", client.repo)
	require.True(t, info.HasUpdate)
	require.Equal(t, "0.1.136-smartapi.2", info.LatestVersion)
	require.Equal(t, "hua7448/sub2api", info.UpdateRepo)
	require.Equal(t, "smartapi", info.StableTagSuffix)
}

func TestUpdateServiceAllowsConfiguredForkRepo(t *testing.T) {
	client := &updateServiceGitHubClientStub{
		release: &GitHubRelease{
			TagName: "v0.1.136-smartapi.1",
			Name:    "v0.1.136-smartapi.1",
		},
	}
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		client,
		"0.1.136-smartapi.1",
		"release",
		config.UpdateConfig{GitHubRepo: "example/smartapi", StableTagSuffix: "smartapi"},
	)

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.Equal(t, "example/smartapi", client.repo)
	require.False(t, info.HasUpdate)
	require.Equal(t, "example/smartapi", info.UpdateRepo)
}

func TestUpdateServiceRejectsNonSmartAPIStableTag(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{
			release: &GitHubRelease{
				TagName: "v0.1.137",
				Name:    "v0.1.137",
			},
		},
		"0.1.136-smartapi.1",
		"release",
	)

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.False(t, info.HasUpdate)
	require.Contains(t, info.Warning, "not a stable smartapi release")
	require.Equal(t, "0.1.136-smartapi.1", info.CurrentVersion)
}

func TestCompareVersionsSmartAPISuffix(t *testing.T) {
	require.Equal(t, -1, compareVersions("0.1.136-smartapi.1", "0.1.136-smartapi.2"))
	require.Equal(t, 1, compareVersions("0.1.137-smartapi.1", "0.1.136-smartapi.9"))
	require.Equal(t, 0, compareVersions("v0.1.136-smartapi.2", "0.1.136-smartapi.2"))
}
