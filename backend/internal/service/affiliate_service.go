package service

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

var (
	ErrAffiliateProfileNotFound   = infraerrors.NotFound("AFFILIATE_PROFILE_NOT_FOUND", "affiliate profile not found")
	ErrAffiliateCodeInvalid       = infraerrors.BadRequest("AFFILIATE_CODE_INVALID", "invalid affiliate code")
	ErrAffiliateCodeTaken         = infraerrors.Conflict("AFFILIATE_CODE_TAKEN", "affiliate code already in use")
	ErrAffiliateAlreadyBound      = infraerrors.Conflict("AFFILIATE_ALREADY_BOUND", "affiliate inviter already bound")
	ErrAffiliateQuotaEmpty        = infraerrors.BadRequest("AFFILIATE_QUOTA_EMPTY", "no affiliate quota available to transfer")
	ErrAffiliateMinInviteesNotMet = infraerrors.BadRequest("AFFILIATE_MIN_INVITEES_NOT_MET", "minimum qualified invitees threshold not met")
	ErrAffiliateDaysQuotaEmpty    = infraerrors.BadRequest("AFFILIATE_DAYS_QUOTA_EMPTY", "no subscription days available to transfer")
)

const (
	affiliateInviteesLimit = 100
	// AffiliateCodeMinLength / AffiliateCodeMaxLength bound both system-generated
	// 12-char codes and admin-customized codes (e.g. "VIP2026").
	AffiliateCodeMinLength = 4
	AffiliateCodeMaxLength = 32
)

// affiliateCodeValidChar accepts uppercase letters, digits, underscore and dash.
// All input passes through strings.ToUpper before validation, so lowercase from
// users is normalized — admins may supply mixed case in their UI.
var affiliateCodeValidChar = func() [256]bool {
	var tbl [256]bool
	for c := byte('A'); c <= 'Z'; c++ {
		tbl[c] = true
	}
	for c := byte('0'); c <= '9'; c++ {
		tbl[c] = true
	}
	tbl['_'] = true
	tbl['-'] = true
	return tbl
}()

// isValidAffiliateCodeFormat validates code format for both binding (user input)
// and admin updates. Caller is expected to upper-case the input first.
func isValidAffiliateCodeFormat(code string) bool {
	if len(code) < AffiliateCodeMinLength || len(code) > AffiliateCodeMaxLength {
		return false
	}
	for i := 0; i < len(code); i++ {
		if !affiliateCodeValidChar[code[i]] {
			return false
		}
	}
	return true
}

type AffiliateSummary struct {
	UserID               int64     `json:"user_id"`
	AffCode              string    `json:"aff_code"`
	AffCodeCustom        bool      `json:"aff_code_custom"`
	AffRebateRatePercent *float64  `json:"aff_rebate_rate_percent,omitempty"`
	InviterID            *int64    `json:"inviter_id,omitempty"`
	AffCount             int       `json:"aff_count"`
	AffQuota             float64   `json:"aff_quota"`
	AffFrozenQuota       float64   `json:"aff_frozen_quota"`
	AffHistoryQuota      float64   `json:"aff_history_quota"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type AffiliateInvitee struct {
	UserID      int64      `json:"user_id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	TotalRebate float64    `json:"total_rebate"`
	Qualified   bool       `json:"qualified"`
}

type AffiliateDetail struct {
	UserID          int64   `json:"user_id"`
	AffCode         string  `json:"aff_code"`
	InviterID       *int64  `json:"inviter_id,omitempty"`
	AffCount        int     `json:"aff_count"`
	AffQuota        float64 `json:"aff_quota"`
	AffFrozenQuota  float64 `json:"aff_frozen_quota"`
	AffHistoryQuota float64 `json:"aff_history_quota"`
	// EffectiveRebateRatePercent 是当前用户作为邀请人时实际生效的返利比例：
	// 优先用户自己的专属比例（aff_rebate_rate_percent），否则回退到全局比例。
	// 用于在用户的 /affiliate 页面直观展示「分享后能拿到多少」。
	EffectiveRebateRatePercent float64                          `json:"effective_rebate_rate_percent"`
	QualifiedInviteeCount      int                              `json:"qualified_invitee_count"`
	MinQualifiedInvitees       int                              `json:"min_qualified_invitees"`
	SubscriptionDaysQuota      []AffiliateSubscriptionDaysQuota `json:"subscription_days_quota"`
	Invitees                   []AffiliateInvitee               `json:"invitees"`
}

type AffiliateRepository interface {
	EnsureUserAffiliate(ctx context.Context, userID int64) (*AffiliateSummary, error)
	GetAffiliateByCode(ctx context.Context, code string) (*AffiliateSummary, error)
	BindInviter(ctx context.Context, userID, inviterID int64) (bool, error)
	AccrueQuota(ctx context.Context, inviterID, inviteeUserID int64, amount float64, freezeHours int, sourceOrderID, sourceRedeemCodeID *int64) (bool, error)
	GetAccruedRebateFromInvitee(ctx context.Context, inviterID, inviteeUserID int64) (float64, error)
	ThawFrozenQuota(ctx context.Context, userID int64) (float64, error)
	TransferQuotaToBalance(ctx context.Context, userID int64) (float64, float64, error)
	ListInvitees(ctx context.Context, inviterID int64, limit int) ([]AffiliateInvitee, error)

	// 有效邀请人数统计
	CountQualifiedInvitees(ctx context.Context, inviterID int64) (int, error)
	CountQualifiedInviteesExcluding(ctx context.Context, inviterID, excludeInviteeUserID int64) (int, error)

	// 订阅天数返利
	AccrueSubscriptionDays(ctx context.Context, inviterID, inviteeUserID int64, groupID int64, days int, freezeHours int, sourceRedeemCodeID *int64) error
	TransferSubscriptionDays(ctx context.Context, userID int64, groupID int64) (int, error)
	ListSubscriptionDaysQuota(ctx context.Context, userID int64) ([]AffiliateSubscriptionDaysQuota, error)

	// 管理端：邀请人列表 + 余额调整
	ListInviters(ctx context.Context, filter AffiliateAdminFilter) ([]AffiliateInviterEntry, int64, error)
	AdjustUserQuota(ctx context.Context, userID int64, newQuota float64) error

	// 管理端：用户级专属配置
	UpdateUserAffCode(ctx context.Context, userID int64, newCode string) error
	ResetUserAffCode(ctx context.Context, userID int64) (string, error)
	SetUserRebateRate(ctx context.Context, userID int64, ratePercent *float64) error
	BatchSetUserRebateRate(ctx context.Context, userIDs []int64, ratePercent *float64) error
	ListUsersWithCustomSettings(ctx context.Context, filter AffiliateAdminFilter) ([]AffiliateAdminEntry, int64, error)
	ListAffiliateInviteRecords(ctx context.Context, filter AffiliateRecordFilter) ([]AffiliateInviteRecord, int64, error)
	ListAffiliateRebateRecords(ctx context.Context, filter AffiliateRecordFilter) ([]AffiliateRebateRecord, int64, error)
	ListAffiliateTransferRecords(ctx context.Context, filter AffiliateRecordFilter) ([]AffiliateTransferRecord, int64, error)
	GetAffiliateUserOverview(ctx context.Context, userID int64) (*AffiliateUserOverview, error)
}

type affiliateTransactionRunner interface {
	RunInTransaction(ctx context.Context, fn func(context.Context) error) error
}

// AffiliateSubscriptionDaysQuota 按分组累积的订阅天数返利
type AffiliateSubscriptionDaysQuota struct {
	GroupID     int64  `json:"group_id"`
	GroupName   string `json:"group_name"`
	PendingDays int    `json:"pending_days"`
	FrozenDays  int    `json:"frozen_days"`
	HistoryDays int    `json:"history_days"`
}

// AffiliateInviterEntry 邀请人列表条目（管理后台）
type AffiliateInviterEntry struct {
	UserID          int64   `json:"user_id"`
	Email           string  `json:"email"`
	Username        string  `json:"username"`
	AffCode         string  `json:"aff_code"`
	AffCount        int     `json:"aff_count"`
	QualifiedCount  int     `json:"qualified_count"`
	AffQuota        float64 `json:"aff_quota"`
	AffFrozenQuota  float64 `json:"aff_frozen_quota"`
	AffHistoryQuota float64 `json:"aff_history_quota"`
}

// AffiliateAdminFilter 列表筛选条件
type AffiliateAdminFilter struct {
	Search   string
	Page     int
	PageSize int
}

// AffiliateAdminEntry 专属用户列表条目
type AffiliateAdminEntry struct {
	UserID               int64    `json:"user_id"`
	Email                string   `json:"email"`
	Username             string   `json:"username"`
	AffCode              string   `json:"aff_code"`
	AffCodeCustom        bool     `json:"aff_code_custom"`
	AffRebateRatePercent *float64 `json:"aff_rebate_rate_percent,omitempty"`
	AffCount             int      `json:"aff_count"`
}

type AffiliateRecordFilter struct {
	Search   string
	Page     int
	PageSize int
	StartAt  *time.Time
	EndAt    *time.Time
	SortBy   string
	SortDesc bool
	// ScopeInviterID 限定返利记录只返回该邀请人（user_affiliate_ledger.user_id）的记录，
	// 用于用户端「我的返利记录」。nil 表示管理端不限定。
	ScopeInviterID *int64
	// ScopeUserID 限定提取记录只返回该用户的记录，用于用户端「我的提取记录」。nil 表示不限定。
	ScopeUserID *int64
}

type AffiliateInviteRecord struct {
	InviterID       int64     `json:"inviter_id"`
	InviterEmail    string    `json:"inviter_email"`
	InviterUsername string    `json:"inviter_username"`
	InviteeID       int64     `json:"invitee_id"`
	InviteeEmail    string    `json:"invitee_email"`
	InviteeUsername string    `json:"invitee_username"`
	AffCode         string    `json:"aff_code"`
	TotalRebate     float64   `json:"total_rebate"`
	CreatedAt       time.Time `json:"created_at"`
}

type AffiliateRebateRecord struct {
	LedgerID        int64     `json:"ledger_id"`
	RewardType      string    `json:"reward_type"`
	SourceType      string    `json:"source_type"`
	RedeemCodeID    int64     `json:"redeem_code_id"`
	GroupID         int64     `json:"group_id"`
	GroupName       string    `json:"group_name"`
	RebateDays      int       `json:"rebate_days"`
	OrderID         int64     `json:"order_id"`
	OutTradeNo      string    `json:"out_trade_no"`
	InviterID       int64     `json:"inviter_id"`
	InviterEmail    string    `json:"inviter_email"`
	InviterUsername string    `json:"inviter_username"`
	InviteeID       int64     `json:"invitee_id"`
	InviteeEmail    string    `json:"invitee_email"`
	InviteeUsername string    `json:"invitee_username"`
	OrderAmount     float64   `json:"order_amount"`
	PayAmount       float64   `json:"pay_amount"`
	RebateAmount    float64   `json:"rebate_amount"`
	PaymentType     string    `json:"payment_type"`
	OrderStatus     string    `json:"order_status"`
	CreatedAt       time.Time `json:"created_at"`
}

type AffiliateTransferRecord struct {
	LedgerID            int64     `json:"ledger_id"`
	UserID              int64     `json:"user_id"`
	UserEmail           string    `json:"user_email"`
	Username            string    `json:"username"`
	Amount              float64   `json:"amount"`
	BalanceAfter        *float64  `json:"balance_after,omitempty"`
	AvailableQuotaAfter *float64  `json:"available_quota_after,omitempty"`
	FrozenQuotaAfter    *float64  `json:"frozen_quota_after,omitempty"`
	HistoryQuotaAfter   *float64  `json:"history_quota_after,omitempty"`
	SnapshotAvailable   bool      `json:"snapshot_available"`
	CurrentBalance      float64   `json:"-"`
	RemainingQuota      float64   `json:"-"`
	FrozenQuota         float64   `json:"-"`
	HistoryQuota        float64   `json:"-"`
	CreatedAt           time.Time `json:"created_at"`
}

type AffiliateUserOverview struct {
	UserID              int64   `json:"user_id"`
	Email               string  `json:"email"`
	Username            string  `json:"username"`
	AffCode             string  `json:"aff_code"`
	RebateRatePercent   float64 `json:"rebate_rate_percent"`
	RebateRateCustom    bool    `json:"-"`
	InvitedCount        int     `json:"invited_count"`
	RebatedInviteeCount int     `json:"rebated_invitee_count"`
	AvailableQuota      float64 `json:"available_quota"`
	HistoryQuota        float64 `json:"history_quota"`
}

type AffiliateService struct {
	repo                 AffiliateRepository
	settingService       *SettingService
	authCacheInvalidator APIKeyAuthCacheInvalidator
	billingCacheService  *BillingCacheService
	subscriptionAssigner DefaultSubscriptionAssigner
}

func NewAffiliateService(repo AffiliateRepository, settingService *SettingService, authCacheInvalidator APIKeyAuthCacheInvalidator, billingCacheService *BillingCacheService, subscriptionAssigner DefaultSubscriptionAssigner) *AffiliateService {
	return &AffiliateService{
		repo:                 repo,
		settingService:       settingService,
		authCacheInvalidator: authCacheInvalidator,
		billingCacheService:  billingCacheService,
		subscriptionAssigner: subscriptionAssigner,
	}
}

// IsEnabled reports whether the affiliate (邀请返利) feature is turned on.
func (s *AffiliateService) IsEnabled(ctx context.Context) bool {
	if s == nil || s.settingService == nil {
		return AffiliateEnabledDefault
	}
	return s.settingService.IsAffiliateEnabled(ctx)
}

func (s *AffiliateService) EnsureUserAffiliate(ctx context.Context, userID int64) (*AffiliateSummary, error) {
	if userID <= 0 {
		return nil, infraerrors.BadRequest("INVALID_USER", "invalid user")
	}
	if s == nil || s.repo == nil {
		return nil, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	return s.repo.EnsureUserAffiliate(ctx, userID)
}

func (s *AffiliateService) GetAffiliateDetail(ctx context.Context, userID int64) (*AffiliateDetail, error) {
	// Lazy thaw: move any matured frozen quota to available before reading.
	if s != nil && s.repo != nil {
		// best-effort: thaw failure is non-fatal
		_, _ = s.repo.ThawFrozenQuota(ctx, userID)
	}

	summary, err := s.EnsureUserAffiliate(ctx, userID)
	if err != nil {
		return nil, err
	}
	invitees, err := s.listInvitees(ctx, userID)
	if err != nil {
		return nil, err
	}

	var qualifiedCount int
	if s.repo != nil {
		qualifiedCount, _ = s.repo.CountQualifiedInvitees(ctx, userID)
	}
	var minInvitees int
	if s.settingService != nil {
		minInvitees = s.settingService.GetAffiliateMinQualifiedInvitees(ctx)
	}
	var daysQuota []AffiliateSubscriptionDaysQuota
	if s.repo != nil {
		daysQuota, _ = s.repo.ListSubscriptionDaysQuota(ctx, userID)
	}

	return &AffiliateDetail{
		UserID:                     summary.UserID,
		AffCode:                    summary.AffCode,
		InviterID:                  summary.InviterID,
		AffCount:                   summary.AffCount,
		AffQuota:                   summary.AffQuota,
		AffFrozenQuota:             summary.AffFrozenQuota,
		AffHistoryQuota:            summary.AffHistoryQuota,
		EffectiveRebateRatePercent: s.resolveRebateRatePercent(ctx, summary),
		QualifiedInviteeCount:      qualifiedCount,
		MinQualifiedInvitees:       minInvitees,
		SubscriptionDaysQuota:      daysQuota,
		Invitees:                   invitees,
	}, nil
}

func (s *AffiliateService) BindInviterByCode(ctx context.Context, userID int64, rawCode string) error {
	code := strings.ToUpper(strings.TrimSpace(rawCode))
	if code == "" {
		return nil
	}
	if s == nil || s.repo == nil {
		return infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	// 总开关关闭时，注册阶段静默忽略 aff 参数（不报错，避免阻断注册流程）
	if !s.IsEnabled(ctx) {
		return nil
	}
	if !isValidAffiliateCodeFormat(code) {
		return ErrAffiliateCodeInvalid
	}

	selfSummary, err := s.repo.EnsureUserAffiliate(ctx, userID)
	if err != nil {
		return err
	}
	if selfSummary.InviterID != nil {
		return nil
	}

	inviterSummary, err := s.repo.GetAffiliateByCode(ctx, code)
	if err != nil {
		if errors.Is(err, ErrAffiliateProfileNotFound) {
			return ErrAffiliateCodeInvalid
		}
		return err
	}
	if inviterSummary == nil || inviterSummary.UserID <= 0 || inviterSummary.UserID == userID {
		return ErrAffiliateCodeInvalid
	}

	bound, err := s.repo.BindInviter(ctx, userID, inviterSummary.UserID)
	if err != nil {
		return err
	}
	if !bound {
		return ErrAffiliateAlreadyBound
	}
	return nil
}

func (s *AffiliateService) AccrueInviteRebate(ctx context.Context, inviteeUserID int64, baseRechargeAmount float64) (float64, error) {
	return s.accrueInviteRebate(ctx, inviteeUserID, baseRechargeAmount, nil, nil)
}

func (s *AffiliateService) AccrueInviteRebateForOrder(ctx context.Context, inviteeUserID int64, baseRechargeAmount float64, sourceOrderID *int64) (float64, error) {
	return s.accrueInviteRebate(ctx, inviteeUserID, baseRechargeAmount, sourceOrderID, nil)
}

// AccrueInviteRebateForRedeem accrues a balance-code rebate and records the
// exact redeem code for auditability and idempotency.
func (s *AffiliateService) AccrueInviteRebateForRedeem(ctx context.Context, inviteeUserID int64, baseRechargeAmount float64, redeemCodeID int64) (float64, error) {
	if redeemCodeID <= 0 {
		return s.AccrueInviteRebate(ctx, inviteeUserID, baseRechargeAmount)
	}
	return s.accrueInviteRebate(ctx, inviteeUserID, baseRechargeAmount, nil, &redeemCodeID)
}

func (s *AffiliateService) accrueInviteRebate(ctx context.Context, inviteeUserID int64, baseRechargeAmount float64, sourceOrderID, sourceRedeemCodeID *int64) (float64, error) {
	if s == nil || s.repo == nil {
		return 0, nil
	}
	if inviteeUserID <= 0 || baseRechargeAmount <= 0 || math.IsNaN(baseRechargeAmount) || math.IsInf(baseRechargeAmount, 0) {
		return 0, nil
	}
	// 总开关关闭时，新充值不再产生返利
	if !s.IsEnabled(ctx) {
		return 0, nil
	}

	inviteeSummary, err := s.repo.EnsureUserAffiliate(ctx, inviteeUserID)
	if err != nil {
		return 0, err
	}
	if inviteeSummary.InviterID == nil || *inviteeSummary.InviterID <= 0 {
		return 0, nil
	}

	// 加载邀请人 profile，优先使用专属比例（覆盖全局）
	inviterSummary, err := s.repo.EnsureUserAffiliate(ctx, *inviteeSummary.InviterID)
	if err != nil {
		return 0, err
	}

	// minInvitees 表示当前客户之前必须已有的有效客户数。排除当前客户，
	// 才能保证 N=1 时从第 2 个客户开始、N=2 时从第 3 个客户开始返利。
	if s.settingService != nil {
		if minInvitees := s.settingService.GetAffiliateMinQualifiedInvitees(ctx); minInvitees > 0 {
			qualifiedCount, qErr := s.repo.CountQualifiedInviteesExcluding(ctx, *inviteeSummary.InviterID, inviteeUserID)
			if qErr != nil {
				return 0, qErr
			}
			if qualifiedCount < minInvitees {
				return 0, nil
			}
		}
	}

	// 有效期检查：超过返利有效期后不再产生返利
	if s.settingService != nil {
		if durationDays := s.settingService.GetAffiliateRebateDurationDays(ctx); durationDays > 0 {
			if time.Now().After(inviteeSummary.CreatedAt.AddDate(0, 0, durationDays)) {
				return 0, nil
			}
		}
	}

	rebateRatePercent := s.resolveRebateRatePercent(ctx, inviterSummary)
	rebate := roundTo(baseRechargeAmount*(rebateRatePercent/100), 8)
	if rebate <= 0 {
		return 0, nil
	}

	// 单人上限检查：精确截断到剩余额度
	if s.settingService != nil {
		if perInviteeCap := s.settingService.GetAffiliateRebatePerInviteeCap(ctx); perInviteeCap > 0 {
			existing, err := s.repo.GetAccruedRebateFromInvitee(ctx, *inviteeSummary.InviterID, inviteeUserID)
			if err != nil {
				return 0, err
			}
			if existing >= perInviteeCap {
				return 0, nil
			}
			if remaining := perInviteeCap - existing; rebate > remaining {
				rebate = roundTo(remaining, 8)
			}
		}
	}

	var freezeHours int
	if s.settingService != nil {
		freezeHours = s.settingService.GetAffiliateRebateFreezeHours(ctx)
	}

	applied, err := s.repo.AccrueQuota(ctx, *inviteeSummary.InviterID, inviteeUserID, rebate, freezeHours, sourceOrderID, sourceRedeemCodeID)
	if err != nil {
		return 0, err
	}
	if !applied {
		return 0, nil
	}
	return rebate, nil
}

// resolveRebateRatePercent returns the inviter's exclusive rate when set,
// otherwise the global setting value (clamped to [Min, Max]).
func (s *AffiliateService) resolveRebateRatePercent(ctx context.Context, inviter *AffiliateSummary) float64 {
	if inviter != nil && inviter.AffRebateRatePercent != nil {
		v := *inviter.AffRebateRatePercent
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return s.globalRebateRatePercent(ctx)
		}
		return clampAffiliateRebateRate(v)
	}
	return s.globalRebateRatePercent(ctx)
}

// globalRebateRatePercent reads the system-wide rebate rate via SettingService,
// returning the documented default when SettingService is unavailable.
func (s *AffiliateService) globalRebateRatePercent(ctx context.Context) float64 {
	if s == nil || s.settingService == nil {
		return AffiliateRebateRateDefault
	}
	return s.settingService.GetAffiliateRebateRatePercent(ctx)
}

func (s *AffiliateService) TransferAffiliateQuota(ctx context.Context, userID int64) (float64, float64, error) {
	if s == nil || s.repo == nil {
		return 0, 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}

	// 最低有效邀请人数门槛检查
	if s.settingService != nil {
		if minInvitees := s.settingService.GetAffiliateMinQualifiedInvitees(ctx); minInvitees > 0 {
			qualifiedCount, err := s.repo.CountQualifiedInvitees(ctx, userID)
			if err != nil {
				return 0, 0, err
			}
			if qualifiedCount < minInvitees {
				return 0, 0, ErrAffiliateMinInviteesNotMet
			}
		}
	}

	transferred, balance, err := s.repo.TransferQuotaToBalance(ctx, userID)
	if err != nil {
		return 0, 0, err
	}
	if transferred > 0 {
		s.invalidateAffiliateCaches(ctx, userID)
	}
	return transferred, balance, nil
}

func (s *AffiliateService) listInvitees(ctx context.Context, inviterID int64) ([]AffiliateInvitee, error) {
	if s == nil || s.repo == nil {
		return nil, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	invitees, err := s.repo.ListInvitees(ctx, inviterID, affiliateInviteesLimit)
	if err != nil {
		return nil, err
	}
	for i := range invitees {
		invitees[i].Email = maskEmail(invitees[i].Email)
	}
	return invitees, nil
}

func roundTo(v float64, scale int) float64 {
	factor := math.Pow10(scale)
	return math.Round(v*factor) / factor
}

func maskEmail(email string) string {
	email = strings.TrimSpace(email)
	if email == "" {
		return ""
	}
	at := strings.Index(email, "@")
	if at <= 0 || at >= len(email)-1 {
		return "***"
	}

	local := email[:at]
	domain := email[at+1:]
	dot := strings.LastIndex(domain, ".")

	maskedLocal := maskSegment(local)
	if dot <= 0 || dot >= len(domain)-1 {
		return maskedLocal + "@" + maskSegment(domain)
	}

	domainName := domain[:dot]
	tld := domain[dot:]
	return maskedLocal + "@" + maskSegment(domainName) + tld
}

func maskSegment(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return "***"
	}
	if len(r) == 1 {
		return string(r[0]) + "***"
	}
	return string(r[0]) + "***"
}

func (s *AffiliateService) invalidateAffiliateCaches(ctx context.Context, userID int64) {
	if s.authCacheInvalidator != nil {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, userID)
	}
	if s.billingCacheService != nil {
		if err := s.billingCacheService.InvalidateUserBalance(ctx, userID); err != nil {
			logger.LegacyPrintf("service.affiliate", "[Affiliate] Failed to invalidate billing cache for user %d: %v", userID, err)
		}
	}
}

// =========================
// Admin: 专属配置管理
// =========================

// validateExclusiveRate ensures a per-user override is finite and within
// [Min, Max]. nil is always valid (means "clear / fall back to global").
func validateExclusiveRate(ratePercent *float64) error {
	if ratePercent == nil {
		return nil
	}
	v := *ratePercent
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return infraerrors.BadRequest("INVALID_RATE", "invalid rebate rate")
	}
	if v < AffiliateRebateRateMin || v > AffiliateRebateRateMax {
		return infraerrors.BadRequest("INVALID_RATE", "rebate rate out of range")
	}
	return nil
}

// AdminUpdateUserAffCode 管理员改写用户的邀请码（专属邀请码）。
func (s *AffiliateService) AdminUpdateUserAffCode(ctx context.Context, userID int64, rawCode string) error {
	if s == nil || s.repo == nil {
		return infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	code := strings.ToUpper(strings.TrimSpace(rawCode))
	if !isValidAffiliateCodeFormat(code) {
		return ErrAffiliateCodeInvalid
	}
	return s.repo.UpdateUserAffCode(ctx, userID, code)
}

// AdminResetUserAffCode 重置用户邀请码为系统随机码。
func (s *AffiliateService) AdminResetUserAffCode(ctx context.Context, userID int64) (string, error) {
	if s == nil || s.repo == nil {
		return "", infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	return s.repo.ResetUserAffCode(ctx, userID)
}

// AdminSetUserRebateRate 设置/清除用户专属返利比例。ratePercent==nil 表示清除。
func (s *AffiliateService) AdminSetUserRebateRate(ctx context.Context, userID int64, ratePercent *float64) error {
	if s == nil || s.repo == nil {
		return infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	if err := validateExclusiveRate(ratePercent); err != nil {
		return err
	}
	return s.repo.SetUserRebateRate(ctx, userID, ratePercent)
}

// AdminBatchSetUserRebateRate 批量设置/清除用户专属返利比例。
func (s *AffiliateService) AdminBatchSetUserRebateRate(ctx context.Context, userIDs []int64, ratePercent *float64) error {
	if s == nil || s.repo == nil {
		return infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	if err := validateExclusiveRate(ratePercent); err != nil {
		return err
	}
	cleaned := make([]int64, 0, len(userIDs))
	for _, uid := range userIDs {
		if uid > 0 {
			cleaned = append(cleaned, uid)
		}
	}
	if len(cleaned) == 0 {
		return nil
	}
	return s.repo.BatchSetUserRebateRate(ctx, cleaned, ratePercent)
}

// AdminListCustomUsers 列出有专属配置的用户。
func (s *AffiliateService) AdminListCustomUsers(ctx context.Context, filter AffiliateAdminFilter) ([]AffiliateAdminEntry, int64, error) {
	if s == nil || s.repo == nil {
		return nil, 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	return s.repo.ListUsersWithCustomSettings(ctx, filter)
}

func (s *AffiliateService) AdminListInviteRecords(ctx context.Context, filter AffiliateRecordFilter) ([]AffiliateInviteRecord, int64, error) {
	if s == nil || s.repo == nil {
		return nil, 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	return s.repo.ListAffiliateInviteRecords(ctx, normalizeAffiliateRecordFilter(filter))
}

func (s *AffiliateService) AdminListRebateRecords(ctx context.Context, filter AffiliateRecordFilter) ([]AffiliateRebateRecord, int64, error) {
	if s == nil || s.repo == nil {
		return nil, 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	return s.repo.ListAffiliateRebateRecords(ctx, normalizeAffiliateRecordFilter(filter))
}

func (s *AffiliateService) AdminListTransferRecords(ctx context.Context, filter AffiliateRecordFilter) ([]AffiliateTransferRecord, int64, error) {
	if s == nil || s.repo == nil {
		return nil, 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	return s.repo.ListAffiliateTransferRecords(ctx, normalizeAffiliateRecordFilter(filter))
}

// UserListRebateRecords 返回指定邀请人（当前登录用户）名下的返利明细。
func (s *AffiliateService) UserListRebateRecords(ctx context.Context, userID int64, filter AffiliateRecordFilter) ([]AffiliateRebateRecord, int64, error) {
	if s == nil || s.repo == nil {
		return nil, 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	if userID <= 0 {
		return nil, 0, infraerrors.BadRequest("INVALID_USER", "invalid user")
	}
	filter.ScopeInviterID = &userID
	return s.repo.ListAffiliateRebateRecords(ctx, normalizeAffiliateRecordFilter(filter))
}

// UserListTransferRecords 返回当前登录用户的返利额度提取（转入余额）明细。
func (s *AffiliateService) UserListTransferRecords(ctx context.Context, userID int64, filter AffiliateRecordFilter) ([]AffiliateTransferRecord, int64, error) {
	if s == nil || s.repo == nil {
		return nil, 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	if userID <= 0 {
		return nil, 0, infraerrors.BadRequest("INVALID_USER", "invalid user")
	}
	filter.ScopeUserID = &userID
	return s.repo.ListAffiliateTransferRecords(ctx, normalizeAffiliateRecordFilter(filter))
}

func (s *AffiliateService) AdminGetUserOverview(ctx context.Context, userID int64) (*AffiliateUserOverview, error) {
	if userID <= 0 {
		return nil, infraerrors.BadRequest("INVALID_USER", "invalid user")
	}
	if s == nil || s.repo == nil {
		return nil, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	overview, err := s.repo.GetAffiliateUserOverview(ctx, userID)
	if err != nil {
		return nil, err
	}
	if overview != nil {
		if !overview.RebateRateCustom {
			overview.RebateRatePercent = s.globalRebateRatePercent(ctx)
		}
		overview.RebateRatePercent = clampAffiliateRebateRate(overview.RebateRatePercent)
	}
	return overview, nil
}

func normalizeAffiliateRecordFilter(filter AffiliateRecordFilter) AffiliateRecordFilter {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	filter.Search = strings.TrimSpace(filter.Search)
	filter.SortBy = strings.TrimSpace(filter.SortBy)
	return filter
}

// AccrueSubscriptionDaysRebate 被邀请人续费订阅时，给邀请人按比例累积天数。
func (s *AffiliateService) AccrueSubscriptionDaysRebate(ctx context.Context, inviteeUserID int64, groupID int64, days int) error {
	return s.accrueSubscriptionDaysRebate(ctx, inviteeUserID, groupID, days, nil)
}

// AccrueSubscriptionDaysRebateForRedeem handles every positive subscription
// redeem code, regardless of group/tier, and records the source code.
func (s *AffiliateService) AccrueSubscriptionDaysRebateForRedeem(ctx context.Context, inviteeUserID int64, groupID int64, days int, redeemCodeID int64) error {
	if redeemCodeID <= 0 {
		return s.AccrueSubscriptionDaysRebate(ctx, inviteeUserID, groupID, days)
	}
	return s.accrueSubscriptionDaysRebate(ctx, inviteeUserID, groupID, days, &redeemCodeID)
}

func (s *AffiliateService) accrueSubscriptionDaysRebate(ctx context.Context, inviteeUserID int64, groupID int64, days int, sourceRedeemCodeID *int64) error {
	if s == nil || s.repo == nil {
		return nil
	}
	if inviteeUserID <= 0 || groupID <= 0 || days <= 0 {
		return nil
	}
	if !s.IsEnabled(ctx) {
		return nil
	}

	inviteeSummary, err := s.repo.EnsureUserAffiliate(ctx, inviteeUserID)
	if err != nil {
		return err
	}
	if inviteeSummary.InviterID == nil || *inviteeSummary.InviterID <= 0 {
		return nil
	}

	inviterSummary, err := s.repo.EnsureUserAffiliate(ctx, *inviteeSummary.InviterID)
	if err != nil {
		return err
	}

	// 当前续费用户不计入返利生效前必须已有的有效客户数。
	if s.settingService != nil {
		if minInvitees := s.settingService.GetAffiliateMinQualifiedInvitees(ctx); minInvitees > 0 {
			qualifiedCount, qErr := s.repo.CountQualifiedInviteesExcluding(ctx, *inviteeSummary.InviterID, inviteeUserID)
			if qErr != nil {
				return qErr
			}
			if qualifiedCount < minInvitees {
				return nil
			}
		}
	}

	// Keep subscription-day rewards under the same invite relationship expiry
	// setting as balance-code rewards.
	if s.settingService != nil {
		if durationDays := s.settingService.GetAffiliateRebateDurationDays(ctx); durationDays > 0 {
			if time.Now().After(inviteeSummary.CreatedAt.AddDate(0, 0, durationDays)) {
				return nil
			}
		}
	}

	rebateRatePercent := s.resolveRebateRatePercent(ctx, inviterSummary)
	rebateDays := int(math.Floor(float64(days) * rebateRatePercent / 100))
	if rebateDays < 1 {
		return nil
	}

	var freezeHours int
	if s.settingService != nil {
		freezeHours = s.settingService.GetAffiliateRebateFreezeHours(ctx)
	}

	return s.repo.AccrueSubscriptionDays(ctx, *inviteeSummary.InviterID, inviteeUserID, groupID, rebateDays, freezeHours, sourceRedeemCodeID)
}

// TransferAffiliateSubscriptionDays 提取指定分组的订阅天数返利，为用户创建/续期订阅。
func (s *AffiliateService) TransferAffiliateSubscriptionDays(ctx context.Context, userID int64, groupID int64) (int, error) {
	if s == nil || s.repo == nil {
		return 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	if s.subscriptionAssigner == nil {
		return 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "subscription service unavailable")
	}

	// 最低有效邀请人数门槛检查
	if s.settingService != nil {
		if minInvitees := s.settingService.GetAffiliateMinQualifiedInvitees(ctx); minInvitees > 0 {
			qualifiedCount, err := s.repo.CountQualifiedInvitees(ctx, userID)
			if err != nil {
				return 0, err
			}
			if qualifiedCount < minInvitees {
				return 0, ErrAffiliateMinInviteesNotMet
			}
		}
	}

	var transferred int
	transferAndAssign := func(txCtx context.Context) error {
		var err error
		transferred, err = s.repo.TransferSubscriptionDays(txCtx, userID, groupID)
		if err != nil {
			return err
		}
		_, _, err = s.subscriptionAssigner.AssignOrExtendSubscription(txCtx, &AssignSubscriptionInput{
			UserID:       userID,
			GroupID:      groupID,
			ValidityDays: transferred,
			AssignedBy:   0,
			Notes:        "邀请返利订阅天数转入",
		})
		return err
	}

	var err error
	if txRunner, ok := s.repo.(affiliateTransactionRunner); ok {
		err = txRunner.RunInTransaction(ctx, transferAndAssign)
	} else {
		err = transferAndAssign(ctx)
	}
	if err != nil {
		return 0, err
	}
	return transferred, nil
}

// AdminListInviters 列出所有有邀请记录的用户（管理后台）。
func (s *AffiliateService) AdminListInviters(ctx context.Context, filter AffiliateAdminFilter) ([]AffiliateInviterEntry, int64, error) {
	if s == nil || s.repo == nil {
		return nil, 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	return s.repo.ListInviters(ctx, filter)
}

// AdminAdjustUserQuota 管理员调整用户返利余额。
func (s *AffiliateService) AdminAdjustUserQuota(ctx context.Context, userID int64, action string, amount float64) error {
	if s == nil || s.repo == nil {
		return infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
	}
	if userID <= 0 {
		return infraerrors.BadRequest("INVALID_USER", "invalid user")
	}

	switch action {
	case "set":
		if amount < 0 {
			return infraerrors.BadRequest("INVALID_AMOUNT", "amount cannot be negative for set action")
		}
		return s.repo.AdjustUserQuota(ctx, userID, amount)
	case "add":
		summary, err := s.repo.EnsureUserAffiliate(ctx, userID)
		if err != nil {
			return err
		}
		newQuota := summary.AffQuota + amount
		if newQuota < 0 {
			newQuota = 0
		}
		return s.repo.AdjustUserQuota(ctx, userID, newQuota)
	default:
		return infraerrors.BadRequest("INVALID_ACTION", "action must be 'set' or 'add'")
	}
}
