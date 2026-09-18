-- Track the weekly quota window in which a user used the self-service
-- daily-quota reset. Existing subscriptions start with one available reset.
ALTER TABLE user_subscriptions
    ADD COLUMN IF NOT EXISTS daily_quota_reset_week_start TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS daily_quota_reset_operations JSONB NOT NULL DEFAULT '[]'::jsonb;

