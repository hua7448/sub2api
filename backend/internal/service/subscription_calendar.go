package service

import "time"

const subscriptionDailyCalendarTimezone = "Asia/Shanghai"

var subscriptionDailyCalendarLocation = mustLoadSubscriptionDailyCalendarLocation()

func mustLoadSubscriptionDailyCalendarLocation() *time.Location {
	location, err := time.LoadLocation(subscriptionDailyCalendarTimezone)
	if err != nil {
		panic("load subscription daily calendar timezone: " + err.Error())
	}
	return location
}

// subscriptionDailyCalendarStart returns the natural-day boundary used by all
// subscription daily quotas, independently of the configurable server timezone.
func subscriptionDailyCalendarStart(value time.Time) time.Time {
	localValue := value.In(subscriptionDailyCalendarLocation)
	return time.Date(localValue.Year(), localValue.Month(), localValue.Day(), 0, 0, 0, 0, subscriptionDailyCalendarLocation)
}
