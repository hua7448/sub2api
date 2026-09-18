package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/stretchr/testify/require"
)

type leaderboardRepoTestDouble struct {
	UsageLogRepository

	mu             sync.Mutex
	listCalls      int
	rankCalls      int
	requestedLimit int
	entries        []usagestats.LeaderboardEntry
	myRank         *usagestats.LeaderboardMyRank
}

func (r *leaderboardRepoTestDouble) GetLeaderboard(_ context.Context, _ usagestats.LeaderboardType, _, _ time.Time, limit int) ([]usagestats.LeaderboardEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listCalls++
	r.requestedLimit = limit
	return r.entries, nil
}

func (r *leaderboardRepoTestDouble) GetUserLeaderboardRank(_ context.Context, _ int64, _ usagestats.LeaderboardType, _, _ time.Time) (*usagestats.LeaderboardMyRank, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rankCalls++
	return r.myRank, nil
}

func TestUsageServiceGetLeaderboardCachesSharedRowsWithoutLimitPollution(t *testing.T) {
	repo := &leaderboardRepoTestDouble{
		entries: []usagestats.LeaderboardEntry{
			{Rank: 1, UserID: 10, MaskedEmail: "al**@exa...", Value: 12},
			{Rank: 2, UserID: 20, MaskedEmail: "bo**@exa...", Value: 8},
			{Rank: 3, UserID: 30, MaskedEmail: "ca**@exa...", Value: 4},
		},
		myRank: &usagestats.LeaderboardMyRank{Rank: 2, Value: 8, NextRankGap: 4},
	}
	svc := NewUsageService(repo, nil, nil, nil)

	publicResult, err := svc.GetLeaderboard(context.Background(), usagestats.LeaderboardTypeCost, usagestats.LeaderboardPeriodToday, 0, 1)
	require.NoError(t, err)
	require.Len(t, publicResult.Items, 1)
	require.Nil(t, publicResult.MyRank)
	require.Equal(t, "\U0001f451 消费之王", publicResult.Items[0].Title)

	authenticatedResult, err := svc.GetLeaderboard(context.Background(), usagestats.LeaderboardTypeCost, usagestats.LeaderboardPeriodToday, 99, 2)
	require.NoError(t, err)
	require.Len(t, authenticatedResult.Items, 2)
	require.Equal(t, "\U0001f3c6 榜上有名", authenticatedResult.Items[1].Title)
	require.Equal(t, repo.myRank, authenticatedResult.MyRank)

	repo.mu.Lock()
	require.Equal(t, 1, repo.listCalls, "same type/period should reuse the process cache")
	require.Equal(t, leaderboardCacheMaxItems, repo.requestedLimit, "cache should always load the maximum supported row count")
	require.Equal(t, 1, repo.rankCalls, "public access must skip the personalized rank query")
	repo.mu.Unlock()
	require.Empty(t, repo.entries[0].Title, "title decoration must not mutate repository or cached rows")
}

func TestLeaderboardPeriodToTimeRangeAt(t *testing.T) {
	require.NoError(t, timezone.Init("UTC"))
	now := time.Date(2026, time.April, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		period    usagestats.LeaderboardPeriod
		wantStart time.Time
		wantEnd   time.Time
	}{
		{usagestats.LeaderboardPeriodLast24h, time.Date(2026, time.April, 14, 10, 30, 0, 0, time.UTC), now},
		{usagestats.LeaderboardPeriodToday, time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC), now},
		{usagestats.LeaderboardPeriodYesterday, time.Date(2026, time.April, 14, 0, 0, 0, 0, time.UTC), time.Date(2026, time.April, 15, 0, 0, 0, 0, time.UTC)},
		{usagestats.LeaderboardPeriodLast7d, time.Date(2026, time.April, 9, 0, 0, 0, 0, time.UTC), now},
		{usagestats.LeaderboardPeriodMonth, time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), now},
		{usagestats.LeaderboardPeriodLast30d, time.Date(2026, time.March, 17, 0, 0, 0, 0, time.UTC), now},
		{usagestats.LeaderboardPeriodLastMonth, time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range tests {
		t.Run(string(tc.period), func(t *testing.T) {
			start, end := leaderboardPeriodToTimeRangeAt(tc.period, now)
			require.Equal(t, tc.wantStart, start)
			require.Equal(t, tc.wantEnd, end)
		})
	}
}
