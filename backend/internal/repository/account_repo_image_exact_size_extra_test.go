package repository

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestShouldEnqueueSchedulerOutboxForExtraUpdates_ImageExactSizeCapabilityIsRelevant(t *testing.T) {
	require.True(t, shouldEnqueueSchedulerOutboxForExtraUpdates(map[string]any{
		service.OpenAIImageExactSizeSupportedExtraKey: true,
	}))
}
