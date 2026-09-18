//go:build unit

package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type providerHallSinkFake struct{ events []ProviderHallBillingEvent }

func (f *providerHallSinkFake) RecordBilling(ev ProviderHallBillingEvent) {
	f.events = append(f.events, ev)
}

func TestProviderHallBillingLinkStatuses(t *testing.T) {
	trace := uuid.New()
	tokenCost := &CostBreakdown{InputCost: 0.2, TotalCost: 0.5, ActualCost: 0.25, BillingMode: string(BillingModeToken)}
	for _, tc := range []struct {
		name          string
		input         *OpenAIRecordUsageInput
		applied       bool
		err           error
		cost          *CostBreakdown
		notApplicable bool
		want          string
		wantActual    string
	}{
		{"applied", &OpenAIRecordUsageInput{ProviderHallTraceID: trace, APIKey: &APIKey{ID: 3}, Result: &OpenAIForwardResult{}}, true, nil, tokenCost, false, "applied", "0.25"},
		{"duplicate_is_free_of_charge", &OpenAIRecordUsageInput{ProviderHallTraceID: trace, APIKey: &APIKey{ID: 3}, Result: &OpenAIForwardResult{}}, false, nil, tokenCost, false, "duplicate", "0"},
		{"fingerprint_conflict_uncertain", &OpenAIRecordUsageInput{ProviderHallTraceID: trace, APIKey: &APIKey{ID: 3}, Result: &OpenAIForwardResult{}}, false, ErrUsageBillingRequestConflict, tokenCost, false, "uncertain", "0"},
		{"failed_is_not_free", &OpenAIRecordUsageInput{ProviderHallTraceID: trace, APIKey: &APIKey{ID: 3}, Result: &OpenAIForwardResult{}}, false, errors.New("insufficient balance"), tokenCost, false, "failed", "0"},
		{"pricing_unavailable", &OpenAIRecordUsageInput{ProviderHallTraceID: trace, APIKey: &APIKey{ID: 3}, Result: &OpenAIForwardResult{}}, true, nil, &CostBreakdown{BillingMode: string(BillingModeToken)}, true, "not_applicable", "0"},
		{"per_request_mode", &OpenAIRecordUsageInput{ProviderHallTraceID: trace, APIKey: &APIKey{ID: 3}, Result: &OpenAIForwardResult{}}, true, nil, &CostBreakdown{ActualCost: 1, BillingMode: string(BillingModePerRequest)}, false, "not_applicable", "0"},
		{"media_bill", &OpenAIRecordUsageInput{ProviderHallTraceID: trace, APIKey: &APIKey{ID: 3}, Result: &OpenAIForwardResult{ImageCount: 1}}, true, nil, tokenCost, false, "not_applicable", "0"},
		{"zero_cost_applied_is_legitimately_free", &OpenAIRecordUsageInput{ProviderHallTraceID: trace, APIKey: &APIKey{ID: 3}, Result: &OpenAIForwardResult{}}, true, nil, &CostBreakdown{BillingMode: string(BillingModeToken)}, false, "applied", "0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sink := &providerHallSinkFake{}
			svc := &OpenAIGatewayService{}
			svc.SetProviderHallBillingSink(sink)
			svc.emitProviderHallBilling(tc.input, "client:req", "fp", tc.applied, tc.err, tc.cost, 1.5, true, tc.notApplicable)
			require.Len(t, sink.events, 1)
			ev := sink.events[0]
			require.Equal(t, tc.want, ev.Status)
			require.Equal(t, trace, ev.TraceID)
			require.Equal(t, "client:req", ev.BillingRequestID)
			require.Equal(t, int64(3), ev.APIKeyID)
			require.Equal(t, "fp", ev.Fingerprint)
			require.Equal(t, 1.5, ev.Multiplier)
			require.True(t, ev.IsSubscription)
			require.Equal(t, tc.wantActual, ev.ActualCost.String())
			require.False(t, ev.BilledAt.IsZero())
		})
	}
	t.Run("untracked_or_unwired_is_silent", func(t *testing.T) {
		sink := &providerHallSinkFake{}
		svc := &OpenAIGatewayService{}
		svc.emitProviderHallBilling(&OpenAIRecordUsageInput{ProviderHallTraceID: trace, APIKey: &APIKey{ID: 3}}, "r", "f", true, nil, tokenCost, 1, false, false)
		svc.SetProviderHallBillingSink(sink)
		svc.emitProviderHallBilling(&OpenAIRecordUsageInput{APIKey: &APIKey{ID: 3}}, "r", "f", true, nil, tokenCost, 1, false, false)
		svc.emitProviderHallBilling(nil, "r", "f", true, nil, tokenCost, 1, false, false)
		var nilSvc *OpenAIGatewayService
		nilSvc.emitProviderHallBilling(&OpenAIRecordUsageInput{ProviderHallTraceID: trace, APIKey: &APIKey{ID: 3}}, "r", "f", true, nil, tokenCost, 1, false, false)
		nilSvc.SetProviderHallBillingSink(sink)
		require.Empty(t, sink.events)
	})
}
