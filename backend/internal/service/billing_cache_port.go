package service

import (
	"time"
)

// SubscriptionCacheData represents cached subscription data
type SubscriptionCacheData struct {
	Status             string
	StartsAt           time.Time
	ExpiresAt          time.Time
	DailyWindowStart   *time.Time
	WeeklyWindowStart  *time.Time
	MonthlyWindowStart *time.Time
	DailyUsage         float64
	WeeklyUsage        float64
	MonthlyUsage       float64
	Version            int64
	SnapshotVersion    int64
	CacheRevision      int64
}
