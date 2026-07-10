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
	release        *GitHubRelease
	repo           string
	recentReleases []*GitHubRelease
	recentErr      error
}

func (s *updateServiceGitHubClientStub) FetchLatestRelease(_ context.Context, repo string) (*GitHubRelease, error) {
	s.repo = repo
	return s.release, nil
}

func (s *updateServiceGitHubClientStub) FetchRecentReleases(_ context.Context, repo string, _ int) ([]*GitHubRelease, error) {
	s.repo = repo
	return s.recentReleases, s.recentErr
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

func newRollbackTestService(current string, releases []*GitHubRelease) *UpdateService {
	return NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{recentReleases: releases},
		current,
		"release",
	)
}

func TestUpdateServiceListRollbackVersionsFiltersAndCaps(t *testing.T) {
	releases := []*GitHubRelease{
		{TagName: "v0.1.148-smartapi.1", PublishedAt: "2026-07-09T00:00:00Z"},                       // newer than current: excluded
		{TagName: "v0.1.147-smartapi.1", PublishedAt: "2026-07-08T00:00:00Z"},                       // current: excluded
		{TagName: "v0.1.146-smartapi.1-rc1", PublishedAt: "2026-07-07T12:00:00Z", Prerelease: true}, // prerelease: excluded
		{TagName: "v0.1.146-smartapi.1", PublishedAt: "2026-07-07T00:00:00Z"},
		{TagName: "v0.1.145-smartapi.1", PublishedAt: "2026-07-06T00:00:00Z", Draft: true}, // draft: excluded
		{TagName: "v0.1.144-smartapi.2", PublishedAt: "2026-07-05T00:00:00Z"},
		{TagName: "v0.1.144-smartapi.2", PublishedAt: "2026-07-05T00:00:00Z"}, // duplicate: excluded
		{TagName: "v0.1.143-smartapi.1", PublishedAt: "2026-07-04T00:00:00Z"},
		{TagName: "v0.1.142-smartapi.1", PublishedAt: "2026-07-03T00:00:00Z"}, // beyond cap of 3: excluded
		{TagName: "v0.1.141", PublishedAt: "2026-07-02T00:00:00Z"},            // upstream tag: excluded
	}
	svc := newRollbackTestService("0.1.147-smartapi.1", releases)

	versions, err := svc.ListRollbackVersions(context.Background())

	require.NoError(t, err)
	require.Len(t, versions, 3)
	require.Equal(t, "0.1.146-smartapi.1", versions[0].Version)
	require.Equal(t, "0.1.144-smartapi.2", versions[1].Version)
	require.Equal(t, "0.1.143-smartapi.1", versions[2].Version)
}

func TestUpdateServiceListRollbackVersionsSortsUnorderedInput(t *testing.T) {
	releases := []*GitHubRelease{
		{TagName: "v0.1.144-smartapi.1"},
		{TagName: "v0.1.146-smartapi.1"},
		{TagName: "v0.1.145-smartapi.1"},
	}
	svc := newRollbackTestService("0.1.147-smartapi.1", releases)

	versions, err := svc.ListRollbackVersions(context.Background())

	require.NoError(t, err)
	require.Len(t, versions, 3)
	require.Equal(t, "0.1.146-smartapi.1", versions[0].Version)
	require.Equal(t, "0.1.145-smartapi.1", versions[1].Version)
	require.Equal(t, "0.1.144-smartapi.1", versions[2].Version)
}

func TestUpdateServiceListRollbackVersionsEmptyWhenNoneOlder(t *testing.T) {
	releases := []*GitHubRelease{
		{TagName: "v0.1.147-smartapi.1"},
		{TagName: "v0.1.148-smartapi.1"},
	}
	svc := newRollbackTestService("0.1.147-smartapi.1", releases)

	versions, err := svc.ListRollbackVersions(context.Background())

	require.NoError(t, err)
	require.Empty(t, versions)
}

func TestUpdateServiceListRollbackVersionsPropagatesFetchError(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{recentErr: errors.New("github unavailable")},
		"0.1.147-smartapi.1",
		"release",
	)

	_, err := svc.ListRollbackVersions(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "github unavailable")
}

func TestUpdateServiceRollbackToVersionRejectsDisallowedTargets(t *testing.T) {
	releases := []*GitHubRelease{
		{TagName: "v0.1.148-smartapi.1"},
		{TagName: "v0.1.147-smartapi.1"},
		{TagName: "v0.1.146-smartapi.1"},
		{TagName: "v0.1.145-smartapi.1"},
		{TagName: "v0.1.144-smartapi.1"},
		{TagName: "v0.1.143-smartapi.1"},
		{TagName: "v0.1.142-smartapi.1"},
	}
	svc := newRollbackTestService("0.1.147-smartapi.1", releases)

	for _, target := range []string{
		"",                    // empty
		"0.1.147-smartapi.1",  // current version
		"v0.1.147-smartapi.1", // current version with prefix
		"0.1.148-smartapi.1",  // newer than current
		"0.1.142-smartapi.1",  // older than the 3 most recent
		"9.9.9-smartapi.1",    // nonexistent
		"0.1.146",             // upstream tag format
	} {
		err := svc.RollbackToVersion(context.Background(), target)
		require.ErrorIs(t, err, ErrRollbackVersionNotAllowed, "target %q should be rejected", target)
	}
}

func TestUpdateServiceRollbackToVersionAcceptsVPrefix(t *testing.T) {
	// No platform asset in the release: the target passes the allowlist check
	// and fails later at asset lookup, proving the version itself was accepted.
	releases := []*GitHubRelease{
		{TagName: "v0.1.147-smartapi.1"},
		{TagName: "v0.1.146-smartapi.1"},
	}
	svc := newRollbackTestService("0.1.147-smartapi.1", releases)

	err := svc.RollbackToVersion(context.Background(), "v0.1.146-smartapi.1")

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrRollbackVersionNotAllowed)
	require.Contains(t, err.Error(), "no compatible release found")
}
