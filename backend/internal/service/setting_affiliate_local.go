package service

import (
	"context"
	"strconv"
	"strings"
)

// GetAffiliateMinQualifiedInvitees returns how many qualified customers must
// already exist before the current customer's recharge can earn a rebate.
// Zero disables the threshold.
func (s *SettingService) GetAffiliateMinQualifiedInvitees(ctx context.Context) int {
	raw, err := s.settingRepo.GetValue(ctx, SettingKeyAffiliateMinQualifiedInvitees)
	if err != nil {
		return AffiliateMinQualifiedInviteesDefault
	}
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value < 0 {
		return AffiliateMinQualifiedInviteesDefault
	}
	if value > AffiliateMinQualifiedInviteesMax {
		return AffiliateMinQualifiedInviteesMax
	}
	return value
}
