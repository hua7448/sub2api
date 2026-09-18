package service

import (
	"context"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)

type schedulerFastRecheckRepo struct {
	schedulerTestOpenAIAccountRepo
	fastCalls int
}

func (r *schedulerFastRecheckRepo) GetByIDForSchedulerRecheck(ctx context.Context, id int64, groupID *int64) (*Account, error) {
	r.fastCalls++
	return r.GetByID(ctx, id)
}

func TestOpenAIGatewayService_RecheckUsesOneQueryRepositoryFastPath(t *testing.T) {
	groupID := int64(7041)
	snapshot := Account{
		ID: 70401, Name: "snapshot", Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
		Status: StatusActive, Schedulable: true, Concurrency: 2, GroupIDs: []int64{groupID},
	}
	dbAccount := snapshot
	dbAccount.Name = "authoritative"
	repo := &schedulerFastRecheckRepo{schedulerTestOpenAIAccountRepo: schedulerTestOpenAIAccountRepo{accounts: []Account{dbAccount}}}
	svc := &OpenAIGatewayService{
		accountRepo:       repo,
		cfg:               &config.Config{RunMode: config.RunModeStandard},
		schedulerSnapshot: &SchedulerSnapshotService{cache: &openAISnapshotCacheStub{}},
	}

	got := svc.recheckSelectedOpenAIAccountFromDB(context.Background(), &snapshot, &groupID, PlatformOpenAI, "gpt-5.1", false, "")
	require.NotNil(t, got)
	require.Equal(t, "authoritative", got.Name)
	require.Equal(t, 1, repo.fastCalls)
	require.Equal(t, []int64{groupID}, got.GroupIDs)
}

func TestOpenAIGatewayService_RecheckKeepsProxyForUngroupedAccount(t *testing.T) {
	proxyID := int64(8041)
	snapshot := Account{ID: 80401, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Status: StatusActive, Schedulable: true, ProxyID: &proxyID, Proxy: &Proxy{ID: proxyID}}
	dbAccount := snapshot
	dbAccount.Proxy = nil
	repo := &schedulerFastRecheckRepo{schedulerTestOpenAIAccountRepo: schedulerTestOpenAIAccountRepo{accounts: []Account{dbAccount}}}
	svc := &OpenAIGatewayService{accountRepo: repo, cfg: &config.Config{RunMode: config.RunModeStandard}, schedulerSnapshot: &SchedulerSnapshotService{cache: &openAISnapshotCacheStub{}}}
	got, err := svc.loadSelectedOpenAIAccountForRecheck(context.Background(), &snapshot, nil)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Same(t, snapshot.Proxy, got.Proxy)
	require.Empty(t, got.GroupIDs)
	require.Empty(t, got.AccountGroups)
}
