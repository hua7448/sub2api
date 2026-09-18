//go:build unit

package service

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestResolveRebateRatePercent_PerUserOverride verifies that per-inviter
// AffRebateRatePercent overrides the global rate, that NULL falls back to the
// global rate, and that out-of-range exclusive rates are clamped silently.
//
// SettingService is left nil here so globalRebateRatePercent returns the
// documented default (AffiliateRebateRateDefault = 20%) - this exercises the
// fallback path without spinning up a settings stub.
func TestResolveRebateRatePercent_PerUserOverride(t *testing.T) {
	t.Parallel()
	svc := &AffiliateService{}

	// nil exclusive rate falls back to global default (20%)
	require.InDelta(t, AffiliateRebateRateDefault,
		svc.resolveRebateRatePercent(context.Background(), &AffiliateSummary{}), 1e-9)

	// exclusive rate set overrides global
	rate := 50.0
	require.InDelta(t, 50.0,
		svc.resolveRebateRatePercent(context.Background(), &AffiliateSummary{AffRebateRatePercent: &rate}), 1e-9)

	// exclusive rate 0 returns 0 (no rebate, intentional)
	zero := 0.0
	require.InDelta(t, 0.0,
		svc.resolveRebateRatePercent(context.Background(), &AffiliateSummary{AffRebateRatePercent: &zero}), 1e-9)

	// exclusive rate above max is clamped to Max
	tooHigh := 250.0
	require.InDelta(t, AffiliateRebateRateMax,
		svc.resolveRebateRatePercent(context.Background(), &AffiliateSummary{AffRebateRatePercent: &tooHigh}), 1e-9)

	// exclusive rate below min is clamped to Min
	tooLow := -5.0
	require.InDelta(t, AffiliateRebateRateMin,
		svc.resolveRebateRatePercent(context.Background(), &AffiliateSummary{AffRebateRatePercent: &tooLow}), 1e-9)
}

// TestIsEnabled_NilSettingServiceReturnsDefault verifies that IsEnabled
// safely handles a nil settingService dependency by returning the default
// (off). This protects callers from nil-pointer crashes in misconfigured
// environments.
func TestIsEnabled_NilSettingServiceReturnsDefault(t *testing.T) {
	t.Parallel()
	svc := &AffiliateService{}
	require.False(t, svc.IsEnabled(context.Background()))
	require.Equal(t, AffiliateEnabledDefault, svc.IsEnabled(context.Background()))
}

// TestValidateExclusiveRate_BoundaryAndInvalid covers the validator used by
// admin-facing rate setters: nil is always valid (clear), in-range values
// are accepted, NaN/Inf and out-of-range values produce a typed BadRequest.
func TestValidateExclusiveRate_BoundaryAndInvalid(t *testing.T) {
	t.Parallel()
	require.NoError(t, validateExclusiveRate(nil))

	for _, v := range []float64{0, 0.01, 50, 99.99, 100} {
		v := v
		require.NoError(t, validateExclusiveRate(&v), "value %v should be valid", v)
	}

	for _, v := range []float64{-0.01, 100.01, -100, 200} {
		v := v
		require.Error(t, validateExclusiveRate(&v), "value %v should be rejected", v)
	}

	nan := math.NaN()
	require.Error(t, validateExclusiveRate(&nan))
	posInf := math.Inf(1)
	require.Error(t, validateExclusiveRate(&posInf))
	negInf := math.Inf(-1)
	require.Error(t, validateExclusiveRate(&negInf))
}

func TestMaskEmail(t *testing.T) {
	t.Parallel()
	require.Equal(t, "a***@g***.com", maskEmail("alice@gmail.com"))
	require.Equal(t, "x***@d***", maskEmail("x@domain"))
	require.Equal(t, "", maskEmail(""))
}

func TestIsValidAffiliateCodeFormat(t *testing.T) {
	t.Parallel()

	// 邀请码格式校验同时服务于：
	// 1) 系统自动生成的 12 位随机码（A-Z 去 I/O，2-9 去 0/1）
	// 2) 管理员设置的自定义专属码（如 "VIP2026"、"NEW_USER-1"）
	// 因此校验放宽到 [A-Z0-9_-]{4,32}（要求调用方先 ToUpper）。
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"valid canonical 12-char", "ABCDEFGHJKLM", true},
		{"valid all digits 2-9", "234567892345", true},
		{"valid mixed", "A2B3C4D5E6F7", true},
		{"valid admin custom short", "VIP1", true},
		{"valid admin custom with hyphen", "NEW-USER", true},
		{"valid admin custom with underscore", "VIP_2026", true},
		{"valid 32-char max", "ABCDEFGHIJKLMNOPQRSTUVWXYZ012345", true},
		// Previously-excluded chars (I/O/0/1) are now allowed since admins may use them.
		{"letter I now allowed", "IBCDEFGHJKLM", true},
		{"letter O now allowed", "OBCDEFGHJKLM", true},
		{"digit 0 now allowed", "0BCDEFGHJKLM", true},
		{"digit 1 now allowed", "1BCDEFGHJKLM", true},
		{"too short (3 chars)", "ABC", false},
		{"too long (33 chars)", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456", false},
		{"lowercase rejected (caller must ToUpper first)", "abcdefghjklm", false},
		{"empty", "", false},
		{"utf8 non-ascii", "ÄÄÄÄÄÄ", false}, // bytes out of charset
		{"ascii punctuation .", "ABCDEFGHJK.M", false},
		{"whitespace", "ABCDEFGHJK M", false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, isValidAffiliateCodeFormat(tc.in))
		})
	}
}

type affiliateRepoNoop struct{}

func (affiliateRepoNoop) RunInTransaction(context.Context, func(context.Context) error) error {
	panic("unexpected RunInTransaction call")
}

func (affiliateRepoNoop) EnsureUserAffiliate(context.Context, int64) (*AffiliateSummary, error) {
	panic("unexpected EnsureUserAffiliate call")
}

func (affiliateRepoNoop) GetAffiliateByCode(context.Context, string) (*AffiliateSummary, error) {
	panic("unexpected GetAffiliateByCode call")
}

func (affiliateRepoNoop) BindInviter(context.Context, int64, int64) (bool, error) {
	panic("unexpected BindInviter call")
}

func (affiliateRepoNoop) AccrueQuota(context.Context, int64, int64, float64, int, *int64, *int64) (bool, error) {
	panic("unexpected AccrueQuota call")
}

func (affiliateRepoNoop) GetAccruedRebateFromInvitee(context.Context, int64, int64) (float64, error) {
	panic("unexpected GetAccruedRebateFromInvitee call")
}

func (affiliateRepoNoop) ThawFrozenQuota(context.Context, int64) (float64, error) {
	panic("unexpected ThawFrozenQuota call")
}

func (affiliateRepoNoop) TransferQuotaToBalance(context.Context, int64) (float64, float64, error) {
	panic("unexpected TransferQuotaToBalance call")
}

func (affiliateRepoNoop) ListInvitees(context.Context, int64, int) ([]AffiliateInvitee, error) {
	panic("unexpected ListInvitees call")
}

func (affiliateRepoNoop) CountQualifiedInvitees(context.Context, int64) (int, error) {
	panic("unexpected CountQualifiedInvitees call")
}

func (affiliateRepoNoop) CountQualifiedInviteesExcluding(context.Context, int64, int64) (int, error) {
	panic("unexpected CountQualifiedInviteesExcluding call")
}

func (affiliateRepoNoop) AccrueSubscriptionDays(context.Context, int64, int64, int64, int, int, *int64) error {
	panic("unexpected AccrueSubscriptionDays call")
}

func (affiliateRepoNoop) TransferSubscriptionDays(context.Context, int64, int64) (int, error) {
	panic("unexpected TransferSubscriptionDays call")
}

func (affiliateRepoNoop) ListSubscriptionDaysQuota(context.Context, int64) ([]AffiliateSubscriptionDaysQuota, error) {
	panic("unexpected ListSubscriptionDaysQuota call")
}

func (affiliateRepoNoop) ListInviters(context.Context, AffiliateAdminFilter) ([]AffiliateInviterEntry, int64, error) {
	panic("unexpected ListInviters call")
}

func (affiliateRepoNoop) AdjustUserQuota(context.Context, int64, float64) error {
	panic("unexpected AdjustUserQuota call")
}

func (affiliateRepoNoop) UpdateUserAffCode(context.Context, int64, string) error {
	panic("unexpected UpdateUserAffCode call")
}

func (affiliateRepoNoop) ResetUserAffCode(context.Context, int64) (string, error) {
	panic("unexpected ResetUserAffCode call")
}

func (affiliateRepoNoop) SetUserRebateRate(context.Context, int64, *float64) error {
	panic("unexpected SetUserRebateRate call")
}

func (affiliateRepoNoop) BatchSetUserRebateRate(context.Context, []int64, *float64) error {
	panic("unexpected BatchSetUserRebateRate call")
}

func (affiliateRepoNoop) ListUsersWithCustomSettings(context.Context, AffiliateAdminFilter) ([]AffiliateAdminEntry, int64, error) {
	panic("unexpected ListUsersWithCustomSettings call")
}

func (affiliateRepoNoop) ListAffiliateInviteRecords(context.Context, AffiliateRecordFilter) ([]AffiliateInviteRecord, int64, error) {
	panic("unexpected ListAffiliateInviteRecords call")
}

func (affiliateRepoNoop) ListAffiliateRebateRecords(context.Context, AffiliateRecordFilter) ([]AffiliateRebateRecord, int64, error) {
	panic("unexpected ListAffiliateRebateRecords call")
}

func (affiliateRepoNoop) ListAffiliateTransferRecords(context.Context, AffiliateRecordFilter) ([]AffiliateTransferRecord, int64, error) {
	panic("unexpected ListAffiliateTransferRecords call")
}

func (affiliateRepoNoop) GetAffiliateUserOverview(context.Context, int64) (*AffiliateUserOverview, error) {
	panic("unexpected GetAffiliateUserOverview call")
}

type affiliateAccrualThresholdRepoStub struct {
	affiliateRepoNoop
	inviteeUserID      int64
	inviterID          int64
	qualifiedExcluding int
	countCalls         int
	accrued            bool
	accruedGroupID     int64
	accruedDays        int
	sourceRedeemCodeID int64
	inviteeCreatedAt   time.Time
}

func (r *affiliateAccrualThresholdRepoStub) EnsureUserAffiliate(_ context.Context, userID int64) (*AffiliateSummary, error) {
	switch userID {
	case r.inviteeUserID:
		inviterID := r.inviterID
		return &AffiliateSummary{UserID: userID, InviterID: &inviterID, CreatedAt: r.inviteeCreatedAt}, nil
	case r.inviterID:
		return &AffiliateSummary{UserID: userID}, nil
	default:
		return nil, errors.New("unexpected affiliate user")
	}
}

func (r *affiliateAccrualThresholdRepoStub) CountQualifiedInviteesExcluding(_ context.Context, inviterID, excludeInviteeUserID int64) (int, error) {
	if inviterID != r.inviterID || excludeInviteeUserID != r.inviteeUserID {
		return 0, errors.New("unexpected threshold query")
	}
	r.countCalls++
	return r.qualifiedExcluding, nil
}

func (r *affiliateAccrualThresholdRepoStub) AccrueQuota(_ context.Context, inviterID, inviteeUserID int64, _ float64, _ int, _ *int64, sourceRedeemCodeID *int64) (bool, error) {
	if inviterID != r.inviterID || inviteeUserID != r.inviteeUserID {
		return false, errors.New("unexpected rebate accrual")
	}
	r.accrued = true
	if sourceRedeemCodeID != nil {
		r.sourceRedeemCodeID = *sourceRedeemCodeID
	}
	return true, nil
}

func (r *affiliateAccrualThresholdRepoStub) AccrueSubscriptionDays(_ context.Context, inviterID, inviteeUserID, groupID int64, days, _ int, sourceRedeemCodeID *int64) error {
	if inviterID != r.inviterID || inviteeUserID != r.inviteeUserID {
		return errors.New("unexpected subscription rebate accrual")
	}
	r.accruedGroupID = groupID
	r.accruedDays = days
	if sourceRedeemCodeID != nil {
		r.sourceRedeemCodeID = *sourceRedeemCodeID
	}
	return nil
}

func TestAccrueInviteRebateStartsAfterConfiguredQualifiedCustomers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		threshold          string
		qualifiedExcluding int
		wantRebate         float64
		wantAccrued        bool
		wantCountCalls     int
	}{
		{name: "one existing customer skips first", threshold: "1", qualifiedExcluding: 0, wantCountCalls: 1},
		{name: "one existing customer starts with second", threshold: "1", qualifiedExcluding: 1, wantRebate: 20, wantAccrued: true, wantCountCalls: 1},
		{name: "two existing customers skip second", threshold: "2", qualifiedExcluding: 1, wantCountCalls: 1},
		{name: "two existing customers start with third", threshold: "2", qualifiedExcluding: 2, wantRebate: 20, wantAccrued: true, wantCountCalls: 1},
		{name: "zero disables threshold", threshold: "0", wantRebate: 20, wantAccrued: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &affiliateAccrualThresholdRepoStub{inviteeUserID: 200, inviterID: 100, qualifiedExcluding: tt.qualifiedExcluding}
			settings := NewSettingService(&settingRepoStub{values: map[string]string{
				SettingKeyAffiliateEnabled:              "true",
				SettingKeyAffiliateRebateRate:           "20",
				SettingKeyAffiliateMinQualifiedInvitees: tt.threshold,
			}}, nil)
			svc := NewAffiliateService(repo, settings, nil, nil, nil)

			rebate, err := svc.AccrueInviteRebate(context.Background(), repo.inviteeUserID, 100)

			require.NoError(t, err)
			require.InDelta(t, tt.wantRebate, rebate, 1e-9)
			require.Equal(t, tt.wantAccrued, repo.accrued)
			require.Equal(t, tt.wantCountCalls, repo.countCalls)
		})
	}
}

func TestGetAffiliateMinQualifiedInviteesDefaultsToOne(t *testing.T) {
	t.Parallel()

	settings := NewSettingService(&settingRepoStub{values: map[string]string{}}, nil)
	require.Equal(t, 1, settings.GetAffiliateMinQualifiedInvitees(context.Background()))
}

func TestRedeemAffiliateRewardsRecordSourceAndSupportEverySubscriptionGroup(t *testing.T) {
	t.Parallel()

	settings := NewSettingService(&settingRepoStub{values: map[string]string{
		SettingKeyAffiliateEnabled:              "true",
		SettingKeyAffiliateRebateRate:           "10",
		SettingKeyAffiliateMinQualifiedInvitees: "0",
	}}, nil)

	t.Run("balance redeem", func(t *testing.T) {
		repo := &affiliateAccrualThresholdRepoStub{inviteeUserID: 200, inviterID: 100}
		svc := NewAffiliateService(repo, settings, nil, nil, nil)

		rebate, err := svc.AccrueInviteRebateForRedeem(context.Background(), repo.inviteeUserID, 50, 7001)

		require.NoError(t, err)
		require.InDelta(t, 5, rebate, 1e-9)
		require.True(t, repo.accrued)
		require.Equal(t, int64(7001), repo.sourceRedeemCodeID)
	})

	for _, groupID := range []int64{3, 11, 17, 26} {
		groupID := groupID
		t.Run("subscription group", func(t *testing.T) {
			repo := &affiliateAccrualThresholdRepoStub{inviteeUserID: 200, inviterID: 100}
			svc := NewAffiliateService(repo, settings, nil, nil, nil)
			redeemCodeID := int64(8000) + groupID

			err := svc.AccrueSubscriptionDaysRebateForRedeem(context.Background(), repo.inviteeUserID, groupID, 30, redeemCodeID)

			require.NoError(t, err)
			require.Equal(t, groupID, repo.accruedGroupID)
			require.Equal(t, 3, repo.accruedDays)
			require.Equal(t, redeemCodeID, repo.sourceRedeemCodeID)
		})
	}
}

func TestSubscriptionRedeemAffiliateRewardHonorsRelationshipDuration(t *testing.T) {
	t.Parallel()

	settings := NewSettingService(&settingRepoStub{values: map[string]string{
		SettingKeyAffiliateEnabled:            "true",
		SettingKeyAffiliateRebateRate:         "10",
		SettingKeyAffiliateRebateDurationDays: "30",
	}}, nil)
	repo := &affiliateAccrualThresholdRepoStub{
		inviteeUserID:    200,
		inviterID:        100,
		inviteeCreatedAt: time.Now().AddDate(0, 0, -31),
	}
	svc := NewAffiliateService(repo, settings, nil, nil, nil)

	err := svc.AccrueSubscriptionDaysRebateForRedeem(context.Background(), repo.inviteeUserID, 17, 30, 8017)

	require.NoError(t, err)
	require.Zero(t, repo.accruedDays)
	require.Zero(t, repo.sourceRedeemCodeID)
}

type affiliateTransferTxKey struct{}

type affiliateTransferRepoStub struct {
	affiliateRepoNoop

	transferred     int
	runTxCalled     bool
	transferSawTx   bool
	transferUserID  int64
	transferGroupID int64
}

func (r *affiliateTransferRepoStub) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	r.runTxCalled = true
	return fn(context.WithValue(ctx, affiliateTransferTxKey{}, true))
}

func (r *affiliateTransferRepoStub) TransferSubscriptionDays(ctx context.Context, userID int64, groupID int64) (int, error) {
	r.transferSawTx = ctx.Value(affiliateTransferTxKey{}) == true
	r.transferUserID = userID
	r.transferGroupID = groupID
	return r.transferred, nil
}

type affiliateTransferAssignerStub struct {
	calls []AssignSubscriptionInput
	err   error
}

func (s *affiliateTransferAssignerStub) AssignOrExtendSubscription(_ context.Context, input *AssignSubscriptionInput) (*UserSubscription, bool, error) {
	if input != nil {
		s.calls = append(s.calls, *input)
	}
	if s.err != nil {
		return nil, false, s.err
	}
	return &UserSubscription{UserID: input.UserID, GroupID: input.GroupID}, false, nil
}

func TestTransferAffiliateSubscriptionDaysAssignsTransferredDays(t *testing.T) {
	repo := &affiliateTransferRepoStub{transferred: 9}
	assigner := &affiliateTransferAssignerStub{}
	svc := NewAffiliateService(repo, nil, nil, nil, assigner)

	days, err := svc.TransferAffiliateSubscriptionDays(context.Background(), 1001, 20)

	require.NoError(t, err)
	require.Equal(t, 9, days)
	require.True(t, repo.runTxCalled)
	require.True(t, repo.transferSawTx)
	require.Equal(t, int64(1001), repo.transferUserID)
	require.Equal(t, int64(20), repo.transferGroupID)
	require.Len(t, assigner.calls, 1)
	require.Equal(t, AssignSubscriptionInput{
		UserID:       1001,
		GroupID:      20,
		ValidityDays: 9,
		AssignedBy:   0,
		Notes:        "邀请返利订阅天数转入",
	}, assigner.calls[0])
}

func TestTransferAffiliateSubscriptionDaysReturnsAssignError(t *testing.T) {
	assignErr := errors.New("assign failed")
	repo := &affiliateTransferRepoStub{transferred: 7}
	assigner := &affiliateTransferAssignerStub{err: assignErr}
	svc := NewAffiliateService(repo, nil, nil, nil, assigner)

	days, err := svc.TransferAffiliateSubscriptionDays(context.Background(), 1001, 20)

	require.ErrorIs(t, err, assignErr)
	require.Zero(t, days)
	require.True(t, repo.runTxCalled)
	require.Len(t, assigner.calls, 1)
}
