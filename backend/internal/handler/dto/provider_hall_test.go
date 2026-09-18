//go:build unit

package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProviderHallUnknownMetricIsNotZero(t *testing.T) {
	encoded, err := json.Marshal(ProviderHallMetric[string]{State: ProviderHallMetricIncomplete, ReasonCode: "coverage_gap"})
	require.NoError(t, err)
	require.JSONEq(t, `{"value":null,"state":"incomplete","reason_code":"coverage_gap","window_start":null,"window_end":null,"computed_at":null}`, string(encoded))
}
