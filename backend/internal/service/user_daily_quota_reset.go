package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

const subscriptionQuotaCacheTimeout = 5 * time.Second

// Keep retrying beyond the billing subscription cache's five-minute maximum
// TTL. If Redis stays unavailable, stale entries expire before this loop ends.
var subscriptionQuotaCacheRetryDelays = []time.Duration{
	5 * time.Second,
	12 * time.Second,
	30 * time.Second,
	1 * time.Minute,
	2 * time.Minute,
	4 * time.Minute,
	6 * time.Minute,
}

type userDailyQuotaCacheRetryState struct {
	resetAt    time.Time
	generation uint64
}

// ResetUserDailyQuota clears only the authenticated user's current daily usage
// and starts a fresh daily window. Weekly and monthly usage remain billable.
func (s *SubscriptionService) ResetUserDailyQuota(ctx context.Context, userID, subscriptionID int64, operationKeyHash string) (*UserDailyQuotaResetResult, error) {
	if s == nil || s.userSubRepo == nil {
		return nil, ErrSubscriptionNotFound
	}
	resetRepo, ok := s.userSubRepo.(UserDailyQuotaResetRepository)
	if !ok {
		return nil, fmt.Errorf("subscription repository does not support user daily quota resets")
	}

	result, err := resetRepo.ResetDailyQuotaForUser(ctx, userID, subscriptionID, operationKeyHash)
	if err != nil {
		return nil, err
	}
	result.ResetAt = result.ResetAt.UTC()

	cacheCtx, cancel := context.WithTimeout(context.Background(), subscriptionQuotaCacheTimeout)
	cacheErr := s.invalidateSubscriptionResetCaches(cacheCtx, userID, result.GroupID, result.ResetAt)
	cancel()
	if cacheErr != nil {
		log.Printf(
			"Warning: subscription quota reset cache invalidation failed for user %d group %d: %v",
			userID,
			result.GroupID,
			cacheErr,
		)
		s.scheduleSubscriptionResetCacheRetries(userID, result.GroupID, result.ResetAt)
	}

	return result, nil
}

func (s *SubscriptionService) invalidateSubscriptionResetCaches(ctx context.Context, userID, groupID int64, resetAt time.Time) error {
	s.InvalidateSubCacheSync(userID, groupID)
	return s.invalidateSubscriptionResetRemoteCaches(ctx, userID, groupID, resetAt)
}

func (s *SubscriptionService) invalidateSubscriptionResetRemoteCaches(ctx context.Context, userID, groupID int64, resetAt time.Time) error {
	if s.billingCacheService == nil {
		return nil
	}

	var errs []error
	if err := s.billingCacheService.InvalidateSubscriptionAt(ctx, userID, groupID, resetAt); err != nil {
		errs = append(errs, fmt.Errorf("invalidate billing subscription cache: %w", err))
	}
	if err := s.billingCacheService.PublishSubscriptionCacheInvalidation(ctx, subCacheKey(userID, groupID)); err != nil {
		errs = append(errs, fmt.Errorf("publish subscription cache invalidation: %w", err))
	}
	return errors.Join(errs...)
}

// Retain the original helpers for focused tests and downstream callers while
// sharing the same barrier/retry machinery with administrator quota resets.
func (s *SubscriptionService) invalidateUserDailyQuotaCaches(ctx context.Context, userID, groupID int64, resetAt time.Time) error {
	return s.invalidateSubscriptionResetCaches(ctx, userID, groupID, resetAt)
}

func (s *SubscriptionService) scheduleUserDailyQuotaCacheRetries(userID, groupID int64, resetAt time.Time) {
	s.scheduleSubscriptionResetCacheRetries(userID, groupID, resetAt)
}

func (s *SubscriptionService) scheduleSubscriptionResetCacheRetries(userID, groupID int64, resetAt time.Time) {
	if s == nil || s.billingCacheService == nil {
		return
	}

	key := subCacheKey(userID, groupID)
	s.dailyQuotaCacheRetryMu.Lock()
	if s.dailyQuotaCacheRetries == nil {
		s.dailyQuotaCacheRetries = make(map[string]*userDailyQuotaCacheRetryState)
	}
	if state, ok := s.dailyQuotaCacheRetries[key]; ok {
		if resetAt.After(state.resetAt) {
			state.resetAt = resetAt
			state.generation++
		}
		s.dailyQuotaCacheRetryMu.Unlock()
		return
	}
	state := &userDailyQuotaCacheRetryState{resetAt: resetAt, generation: 1}
	s.dailyQuotaCacheRetries[key] = state
	s.dailyQuotaCacheRetryMu.Unlock()

	retryDelays := append([]time.Duration(nil), subscriptionQuotaCacheRetryDelays...)
	go func() {
		for {
			_, cycleGeneration := s.userDailyQuotaCacheRetrySnapshot(state)
			previousDelay := time.Duration(0)
			lastAttemptedCutoff := time.Time{}
			for _, delay := range retryDelays {
				timer := time.NewTimer(delay - previousDelay)
				<-timer.C
				previousDelay = delay

				for {
					cutoff, _ := s.userDailyQuotaCacheRetrySnapshot(state)
					lastAttemptedCutoff = cutoff
					ctx, cancel := context.WithTimeout(context.Background(), subscriptionQuotaCacheTimeout)
					err := s.invalidateSubscriptionResetCaches(ctx, userID, groupID, cutoff)
					cancel()
					if err == nil {
						if s.completeUserDailyQuotaCacheRetry(key, state, cutoff) {
							return
						}
						// A newer reset joined while this attempt was in flight.
						// Retry it immediately while Redis is known to be healthy.
						continue
					}
					log.Printf(
						"Warning: subscription quota reset cache retry failed for user %d group %d at delay %s: %v",
						userID,
						groupID,
						delay,
						err,
					)
					break
				}
			}

			if s.finishUserDailyQuotaCacheRetry(key, state, cycleGeneration, lastAttemptedCutoff) {
				return
			}
			// A newer reset joined after the final failed attempt. Keep the
			// existing worker registered and retry the new cutoff.
		}
	}()
}

func (s *SubscriptionService) userDailyQuotaCacheRetrySnapshot(state *userDailyQuotaCacheRetryState) (time.Time, uint64) {
	s.dailyQuotaCacheRetryMu.Lock()
	defer s.dailyQuotaCacheRetryMu.Unlock()
	return state.resetAt, state.generation
}

func (s *SubscriptionService) completeUserDailyQuotaCacheRetry(key string, state *userDailyQuotaCacheRetryState, cutoff time.Time) bool {
	s.dailyQuotaCacheRetryMu.Lock()
	defer s.dailyQuotaCacheRetryMu.Unlock()
	if state.resetAt.After(cutoff) {
		return false
	}
	if current, ok := s.dailyQuotaCacheRetries[key]; ok && current == state {
		delete(s.dailyQuotaCacheRetries, key)
	}
	return true
}

func (s *SubscriptionService) finishUserDailyQuotaCacheRetry(key string, state *userDailyQuotaCacheRetryState, cycleGeneration uint64, lastAttemptedCutoff time.Time) bool {
	s.dailyQuotaCacheRetryMu.Lock()
	defer s.dailyQuotaCacheRetryMu.Unlock()
	current, ok := s.dailyQuotaCacheRetries[key]
	if !ok || current != state {
		return true
	}
	if state.generation != cycleGeneration || (!lastAttemptedCutoff.IsZero() && state.resetAt.After(lastAttemptedCutoff)) {
		return false
	}
	delete(s.dailyQuotaCacheRetries, key)
	return true
}
