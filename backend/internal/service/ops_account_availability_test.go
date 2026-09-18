//go:build unit

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAccountTempUnschedulableForOpsKeepActive(t *testing.T) {
	now := time.Now().UTC()
	future := now.Add(5 * time.Minute)
	past := now.Add(-5 * time.Minute)

	tests := []struct {
		name string
		acc  *Account
		want bool
	}{
		{
			name: "ordinary account with active cooldown",
			acc:  &Account{TempUnschedulableUntil: &future},
			want: true,
		},
		{
			name: "keep active ignores historical cooldown",
			acc: &Account{
				TempUnschedulableUntil: &future,
				Extra:                  map[string]any{KeepStatusActiveExtraKey: true},
			},
			want: false,
		},
		{
			name: "expired cooldown",
			acc:  &Account{TempUnschedulableUntil: &past},
			want: false,
		},
		{
			name: "nil account",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, accountTempUnschedulableForOps(tt.acc, now))
		})
	}
}

func TestAccountAutomaticAvailabilityForOpsKeepActive(t *testing.T) {
	now := time.Now().UTC()
	future := now.Add(5 * time.Minute)
	past := now.Add(-5 * time.Minute)

	tests := []struct {
		name        string
		acc         *Account
		rateLimited bool
		overloaded  bool
	}{
		{
			name: "ordinary account with active automatic barriers",
			acc: &Account{
				RateLimitResetAt: &future,
				OverloadUntil:    &future,
			},
			rateLimited: true,
			overloaded:  true,
		},
		{
			name: "keep active ignores historical automatic barriers",
			acc: &Account{
				RateLimitResetAt: &future,
				OverloadUntil:    &future,
				Extra:            map[string]any{KeepStatusActiveExtraKey: true},
			},
		},
		{
			name: "expired automatic barriers",
			acc: &Account{
				RateLimitResetAt: &past,
				OverloadUntil:    &past,
			},
		},
		{name: "nil account"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.rateLimited, accountRateLimitedForOps(tt.acc, now))
			require.Equal(t, tt.overloaded, accountOverloadedForOps(tt.acc, now))
		})
	}
}

func TestAccountAvailabilityStateForOpsPreservesManualGates(t *testing.T) {
	now := time.Now().UTC()
	future := now.Add(5 * time.Minute)
	protectedExtra := map[string]any{KeepStatusActiveExtraKey: true}

	tests := []struct {
		name string
		acc  *Account
		want accountOpsAvailabilityState
	}{
		{
			name: "protected active schedulable account ignores automatic barriers",
			acc: &Account{
				Status:                 StatusActive,
				Schedulable:            true,
				RateLimitResetAt:       &future,
				OverloadUntil:          &future,
				TempUnschedulableUntil: &future,
				Extra:                  protectedExtra,
			},
			want: accountOpsAvailabilityState{available: true},
		},
		{
			name: "protected account still obeys manual schedulable false",
			acc: &Account{
				Status:      StatusActive,
				Schedulable: false,
				Extra:       protectedExtra,
			},
			want: accountOpsAvailabilityState{},
		},
		{
			name: "protected account still obeys disabled status",
			acc: &Account{
				Status:      StatusDisabled,
				Schedulable: true,
				Extra:       protectedExtra,
			},
			want: accountOpsAvailabilityState{},
		},
		{
			name: "ordinary account reports automatic barriers",
			acc: &Account{
				Status:                 StatusActive,
				Schedulable:            true,
				RateLimitResetAt:       &future,
				OverloadUntil:          &future,
				TempUnschedulableUntil: &future,
			},
			want: accountOpsAvailabilityState{
				rateLimited:       true,
				overloaded:        true,
				tempUnschedulable: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, accountAvailabilityStateForOps(tt.acc, now))
		})
	}
}
