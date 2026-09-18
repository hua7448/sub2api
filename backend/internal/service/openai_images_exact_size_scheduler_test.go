package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

func TestAccountSupportsOpenAIImageCapability_ExactSizeRequiresAPIKey(t *testing.T) {
	oauth := &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuth}
	unmarkedAPIKey := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKey}
	apiKey := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Extra: map[string]any{
		OpenAIImageExactSizeSupportedExtraKey: true,
	}}
	mappedAPIKey := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Extra: map[string]any{
		OpenAIImageExactSizeSupportedExtraKey: true,
	}, Credentials: map[string]any{
		"model_mapping": map[string]any{"gpt-image-2": "gpt-image-1"},
	}}

	require.True(t, oauth.SupportsOpenAIImageCapability(OpenAIImagesCapabilityNative))
	require.False(t, oauth.SupportsOpenAIImageCapability(OpenAIImagesCapabilityExactSize))
	require.False(t, unmarkedAPIKey.SupportsOpenAIImageCapability(OpenAIImagesCapabilityExactSize))
	require.True(t, apiKey.SupportsOpenAIImageCapability(OpenAIImagesCapabilityExactSize))
	require.False(t, mappedAPIKey.SupportsOpenAIImageCapability(OpenAIImagesCapabilityExactSize))
}

func TestOpenAIGatewayService_SelectAccountForExactImageSizeNeverFallsBackToOAuth(t *testing.T) {
	groupID := int64(27)
	oauth := Account{
		ID: 9001, Platform: PlatformOpenAI, Type: AccountTypeOAuth,
		Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 100,
	}
	apiKey := Account{
		ID: 9002, Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
		Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 1,
		Extra: map[string]any{OpenAIImageExactSizeSupportedExtraKey: true},
	}
	mappedAPIKey := Account{
		ID: 9003, Platform: PlatformOpenAI, Type: AccountTypeAPIKey,
		Status: StatusActive, Schedulable: true, Concurrency: 1, Priority: 100,
		Extra: map[string]any{OpenAIImageExactSizeSupportedExtraKey: true},
		Credentials: map[string]any{
			"model_mapping": map[string]any{"gpt-image-2": "gpt-image-1"},
		},
	}
	newService := func(accounts []Account) *OpenAIGatewayService {
		cfg := &config.Config{}
		cfg.Gateway.Scheduling.LoadBatchEnabled = false
		return &OpenAIGatewayService{
			accountRepo:        schedulerTestOpenAIAccountRepo{accounts: accounts},
			cache:              &schedulerTestGatewayCache{},
			cfg:                cfg,
			rateLimitService:   newOpenAIAdvancedSchedulerRateLimitService("true"),
			concurrencyService: NewConcurrencyService(schedulerTestConcurrencyCache{}),
		}
	}

	t.Run("mixed pool selects API key", func(t *testing.T) {
		selection, _, err := newService([]Account{oauth, mappedAPIKey, apiKey}).SelectAccountWithSchedulerForImages(
			context.Background(), &groupID, "", "gpt-image-2", nil, OpenAIImagesCapabilityExactSize,
		)
		require.NoError(t, err)
		require.NotNil(t, selection)
		require.NotNil(t, selection.Account)
		require.Equal(t, apiKey.ID, selection.Account.ID)
	})

	t.Run("pool without opted-in API key fails closed", func(t *testing.T) {
		selection, _, err := newService([]Account{oauth, mappedAPIKey}).SelectAccountWithSchedulerForImages(
			context.Background(), &groupID, "", "gpt-image-2", nil, OpenAIImagesCapabilityExactSize,
		)
		require.ErrorIs(t, err, ErrNoAvailableAccounts)
		require.Nil(t, selection)
	})
}
