package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GeminiTokenCache stores short-lived access tokens and coordinates refresh to avoid stampedes.
type GeminiTokenCache interface {
	// cacheKey should be stable for the token scope; for GeminiCli OAuth we primarily use project_id.
	GetAccessToken(ctx context.Context, cacheKey string) (string, error)
	SetAccessToken(ctx context.Context, cacheKey string, token string, ttl time.Duration) error
	DeleteAccessToken(ctx context.Context, cacheKey string) error

	// owner must be unique to this acquisition and retained for its release.
	AcquireRefreshLock(ctx context.Context, cacheKey, owner string, ttl time.Duration) (bool, error)
	ReleaseRefreshLock(ctx context.Context, cacheKey, owner string) error
}

// Each acquisition needs a distinct owner even when the account and process match.
func newOAuthRefreshLockOwner() string { return uuid.NewString() }
