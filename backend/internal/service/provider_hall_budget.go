package service

import (
	"time"

	"github.com/shopspring/decimal"
)

// Spend statuses in provider_hall_spend and the sample mirror.
const (
	ProviderHallSpendConfirmed = "confirmed"
	ProviderHallSpendUncertain = "uncertain"
	ProviderHallSpendFailed    = "failed"

	ProviderHallSampleBillingPending   = "pending"
	ProviderHallSampleBillingConfirmed = "confirmed"
	ProviderHallSampleBillingUncertain = "uncertain"
	ProviderHallSampleBillingFailed    = "failed"
	ProviderHallSampleBillingUnbilled  = "unbilled"
)

// ProviderHallSpendStatusForBilling maps a billing outcome to the spend ledger
// status and the sample mirror. ok=false means "write nothing" (duplicate).
func ProviderHallSpendStatusForBilling(billingStatus string) (spendStatus, sampleStatus string, writeSpend, ok bool) {
	switch billingStatus {
	case "applied":
		return ProviderHallSpendConfirmed, ProviderHallSampleBillingConfirmed, true, true
	case "failed":
		return ProviderHallSpendFailed, ProviderHallSampleBillingFailed, true, true
	case "uncertain":
		return ProviderHallSpendUncertain, ProviderHallSampleBillingUncertain, true, true
	case "not_applicable":
		return "", ProviderHallSampleBillingUnbilled, false, true
	default: // duplicate: the first bill already counted
		return "", "", false, false
	}
}

// ProviderHallBudgetExhausted applies the daily budget rule: a zero budget
// means unlimited; otherwise confirmed spend at or above the budget stops
// new dispatches. Uncertain spend is reported but does not block.
func ProviderHallBudgetExhausted(dailyBudget string, confirmed decimal.Decimal) bool {
	budget, err := decimal.NewFromString(dailyBudget)
	if err != nil || !budget.IsPositive() {
		return false
	}
	return confirmed.GreaterThanOrEqual(budget)
}

// ProviderHallBudgetDayFor returns the budget day a dispatch belongs to. The
// job keeps the day of its first dispatch so samples sent after midnight are
// still charged to the day the job started.
func ProviderHallBudgetDayFor(existing *string, dispatchedAt time.Time) string {
	if existing != nil && *existing != "" {
		return *existing
	}
	return ProviderHallBudgetDay(dispatchedAt)
}
