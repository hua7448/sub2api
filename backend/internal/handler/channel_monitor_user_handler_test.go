//go:build unit

package handler

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

type channelMonitorUserSettingRepoStub struct {
	values map[string]string
}

func (s *channelMonitorUserSettingRepoStub) Get(_ context.Context, key string) (*service.Setting, error) {
	return &service.Setting{Key: key, Value: s.values[key]}, nil
}

func (s *channelMonitorUserSettingRepoStub) GetValue(_ context.Context, key string) (string, error) {
	return s.values[key], nil
}

func (s *channelMonitorUserSettingRepoStub) Set(context.Context, string, string) error {
	return nil
}

func (s *channelMonitorUserSettingRepoStub) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		out[key] = s.values[key]
	}
	return out, nil
}

func (s *channelMonitorUserSettingRepoStub) SetMultiple(context.Context, map[string]string) error {
	return nil
}

func (s *channelMonitorUserSettingRepoStub) GetAll(context.Context) (map[string]string, error) {
	return s.values, nil
}

func (s *channelMonitorUserSettingRepoStub) Delete(context.Context, string) error {
	return nil
}

func TestChannelMonitorUserFeatureUsesPublicStatusSwitch(t *testing.T) {
	repo := &channelMonitorUserSettingRepoStub{values: map[string]string{
		service.SettingKeyChannelMonitorEnabled:      "true",
		service.SettingKeyChannelStatusPublicEnabled: "false",
	}}
	settingSvc := service.NewSettingService(repo, &config.Config{})
	h := NewChannelMonitorUserHandler(nil, settingSvc)

	require.False(t, h.settingService.GetChannelStatusPublicRuntime(context.Background()).Enabled)
	require.True(t, h.settingService.GetChannelMonitorRuntime(context.Background()).Enabled)
}
