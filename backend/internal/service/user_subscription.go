package service

import "time"

const subscriptionDayDuration = 24 * time.Hour
const subscriptionWeekDuration = 7 * subscriptionDayDuration
const UserDailyQuotaResetIdempotencyTTL = 8 * subscriptionDayDuration

type UserSubscription struct {
	ID      int64
	UserID  int64
	GroupID int64

	StartsAt  time.Time
	ExpiresAt time.Time
	Status    string

	DailyWindowStart           *time.Time
	WeeklyWindowStart          *time.Time
	MonthlyWindowStart         *time.Time
	DailyQuotaResetWeekStart   *time.Time
	DailyQuotaResetAvailable   bool
	DailyQuotaResetAvailableAt *time.Time

	DailyUsageUSD   float64
	WeeklyUsageUSD  float64
	MonthlyUsageUSD float64

	AssignedBy *int64
	AssignedAt time.Time
	Notes      string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	User           *User
	Group          *Group
	AssignedByUser *User
}

// UserDailyQuotaResetResult is the user-visible snapshot produced by a manual
// daily quota reset. GroupID is internal cache-routing metadata.
type UserDailyQuotaResetResult struct {
	SubscriptionID             int64      `json:"subscription_id"`
	ResetAt                    time.Time  `json:"reset_at"`
	DailyUsageUSD              float64    `json:"daily_usage_usd"`
	WeeklyUsageUSD             float64    `json:"weekly_usage_usd"`
	MonthlyUsageUSD            float64    `json:"monthly_usage_usd"`
	DailyQuotaResetAvailable   bool       `json:"daily_quota_reset_available"`
	DailyQuotaResetAvailableAt *time.Time `json:"daily_quota_reset_available_at"`
	GroupID                    int64      `json:"-"`
	OperationReplayed          bool       `json:"-"`
}

func (s *UserSubscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive && time.Now().Before(s.ExpiresAt)
}

func (s *UserSubscription) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *UserSubscription) DaysRemaining() int {
	return s.daysRemainingAt(time.Now())
}

func (s *UserSubscription) daysRemainingAt(now time.Time) int {
	remaining := s.ExpiresAt.Sub(now)
	if remaining <= 0 {
		return 0
	}

	days := int(remaining / subscriptionDayDuration)
	if remaining%subscriptionDayDuration != 0 {
		days++
	}
	return days
}

func (s *UserSubscription) IsWindowActivated() bool {
	return s.DailyWindowStart != nil || s.WeeklyWindowStart != nil || s.MonthlyWindowStart != nil
}

func (s *UserSubscription) HasOneTimeDailyQuota() bool {
	if s == nil || s.StartsAt.IsZero() || s.ExpiresAt.IsZero() {
		return false
	}
	return !s.ExpiresAt.After(s.StartsAt.AddDate(0, 0, 1))
}

func (s *UserSubscription) NeedsDailyReset() bool {
	return s.NeedsDailyResetAt(time.Now())
}

func (s *UserSubscription) NeedsDailyResetAt(now time.Time) bool {
	_, ok := s.automaticDailyWindowStartAt(now)
	return ok
}

func (s *UserSubscription) NeedsWeeklyReset() bool {
	return s.NeedsWeeklyResetAt(time.Now())
}

func (s *UserSubscription) NeedsWeeklyResetAt(now time.Time) bool {
	if s.WeeklyWindowStart == nil {
		return false
	}
	return !now.Before(s.WeeklyWindowStart.Add(7 * 24 * time.Hour))
}

func (s *UserSubscription) NeedsMonthlyReset() bool {
	return s.NeedsMonthlyResetAt(time.Now())
}

func (s *UserSubscription) NeedsMonthlyResetAt(now time.Time) bool {
	if s.MonthlyWindowStart == nil {
		return false
	}
	return !now.Before(s.MonthlyWindowStart.Add(30 * 24 * time.Hour))
}

func (s *UserSubscription) canAutomaticallyResetDailyAt(now time.Time) bool {
	_, ok := s.automaticDailyWindowStartAt(now)
	return !s.HasOneTimeDailyQuota() && ok
}

func (s *UserSubscription) automaticDailyWindowStartAt(now time.Time) (time.Time, bool) {
	if s == nil || s.DailyWindowStart == nil || s.HasOneTimeDailyQuota() || s.ExpiresAt.IsZero() || !now.Before(s.ExpiresAt) {
		return time.Time{}, false
	}
	windowStart := subscriptionDailyCalendarStart(now)
	if !s.DailyWindowStart.Before(windowStart) || !windowStart.Before(s.ExpiresAt) {
		return time.Time{}, false
	}
	return windowStart, true
}

func (s *UserSubscription) canAutomaticallyResetWeeklyAt(now time.Time) bool {
	_, ok := s.automaticWindowStartAt(s.WeeklyWindowStart, 7*24*time.Hour, now)
	return ok
}

func (s *UserSubscription) canAutomaticallyResetMonthlyAt(now time.Time) bool {
	_, ok := s.automaticWindowStartAt(s.MonthlyWindowStart, 30*24*time.Hour, now)
	return ok
}

func (s *UserSubscription) automaticWindowStartAt(previous *time.Time, period time.Duration, now time.Time) (time.Time, bool) {
	anchor, ok := s.effectiveWindowAnchor(previous)
	if !ok {
		return time.Time{}, false
	}
	next := anchor.Add(period)
	if now.Before(next) || !next.Before(s.ExpiresAt) {
		return time.Time{}, false
	}

	periods := now.Sub(anchor) / period
	lastPeriodBeforeExpiry := (s.ExpiresAt.Sub(anchor) - 1) / period
	if periods > lastPeriodBeforeExpiry {
		periods = lastPeriodBeforeExpiry
	}
	return anchor.Add(periods * period), true
}

func (s *UserSubscription) effectiveWindowAnchor(previous *time.Time) (time.Time, bool) {
	if previous == nil {
		return time.Time{}, false
	}

	anchor := *previous
	// Older subscriptions initialized their first windows at midnight on their
	// start date. Only that initial value is unambiguous; later midnight anchors
	// may be manual resets and must remain authoritative.
	legacyAnchor := startOfDay(s.StartsAt)
	if legacyAnchor.Before(s.StartsAt) && anchor.Equal(legacyAnchor) {
		anchor = s.StartsAt
	}
	return anchor, true
}

func (s *UserSubscription) currentWeeklyQuotaWindowAt(now time.Time) (time.Time, time.Time, bool) {
	if s == nil {
		return time.Time{}, time.Time{}, false
	}
	anchor, ok := s.effectiveWindowAnchor(s.WeeklyWindowStart)
	if !ok {
		if s.DailyQuotaResetWeekStart == nil {
			return time.Time{}, time.Time{}, false
		}
		anchor = *s.DailyQuotaResetWeekStart
	}
	if now.Before(anchor) {
		return time.Time{}, time.Time{}, false
	}

	windowStart := anchor.Add((now.Sub(anchor) / subscriptionWeekDuration) * subscriptionWeekDuration)
	windowEnd := windowStart.Add(subscriptionWeekDuration)
	return windowStart, windowEnd, true
}

// SetDailyQuotaResetAvailabilityAt derives whether the user still has the one
// self-service daily reset allowed in the current subscription weekly window.
func (s *UserSubscription) SetDailyQuotaResetAvailabilityAt(now time.Time) {
	if s == nil {
		return
	}
	s.DailyQuotaResetAvailable = false
	s.DailyQuotaResetAvailableAt = nil

	if s.Status != SubscriptionStatusActive ||
		now.Before(s.StartsAt) ||
		!now.Before(s.ExpiresAt) ||
		s.Group == nil ||
		s.Group.DailyLimitUSD == nil ||
		*s.Group.DailyLimitUSD <= 0 {
		return
	}
	if s.DailyQuotaResetWeekStart == nil {
		s.DailyQuotaResetAvailable = true
		return
	}

	windowStart, _, ok := s.currentWeeklyQuotaWindowAt(now)
	if !ok {
		return
	}
	if s.DailyQuotaResetWeekStart.Before(windowStart) {
		s.DailyQuotaResetAvailable = true
		return
	}
	if s.DailyQuotaResetWeekStart.After(windowStart) {
		return
	}
	windowEnd := windowStart.Add(subscriptionWeekDuration)
	if windowEnd.Before(s.ExpiresAt) {
		availableAt := windowEnd
		s.DailyQuotaResetAvailableAt = &availableAt
	}
}

func (s *UserSubscription) windowResetTime(previous *time.Time, period time.Duration) *time.Time {
	anchor, ok := s.effectiveWindowAnchor(previous)
	if !ok {
		return nil
	}
	resetAt := anchor.Add(period)
	if !s.ExpiresAt.IsZero() && !resetAt.Before(s.ExpiresAt) {
		return nil
	}
	return &resetAt
}

func (s *UserSubscription) DailyResetTime() *time.Time {
	if s.DailyWindowStart == nil {
		return nil
	}
	if s.HasOneTimeDailyQuota() {
		t := s.ExpiresAt
		return &t
	}
	resetAt := subscriptionDailyCalendarStart(*s.DailyWindowStart).AddDate(0, 0, 1)
	if !s.ExpiresAt.IsZero() && !resetAt.Before(s.ExpiresAt) {
		return nil
	}
	return &resetAt
}

func (s *UserSubscription) WeeklyResetTime() *time.Time {
	return s.windowResetTime(s.WeeklyWindowStart, 7*24*time.Hour)
}

func (s *UserSubscription) MonthlyResetTime() *time.Time {
	return s.windowResetTime(s.MonthlyWindowStart, 30*24*time.Hour)
}

func (s *UserSubscription) CheckDailyLimit(group *Group, additionalCost float64) bool {
	if !group.HasDailyLimit() {
		return true
	}
	return s.DailyUsageUSD+additionalCost <= *group.DailyLimitUSD
}

func (s *UserSubscription) CheckWeeklyLimit(group *Group, additionalCost float64) bool {
	if !group.HasWeeklyLimit() {
		return true
	}
	return s.WeeklyUsageUSD+additionalCost <= *group.WeeklyLimitUSD
}

func (s *UserSubscription) CheckMonthlyLimit(group *Group, additionalCost float64) bool {
	if !group.HasMonthlyLimit() {
		return true
	}
	return s.MonthlyUsageUSD+additionalCost <= *group.MonthlyLimitUSD
}

func (s *UserSubscription) CheckAllLimits(group *Group, additionalCost float64) (daily, weekly, monthly bool) {
	daily = s.CheckDailyLimit(group, additionalCost)
	weekly = s.CheckWeeklyLimit(group, additionalCost)
	monthly = s.CheckMonthlyLimit(group, additionalCost)
	return
}
