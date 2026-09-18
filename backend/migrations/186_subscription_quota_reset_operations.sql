-- Durable idempotency and operation audit for the global subscription quota reset.
-- Raw idempotency keys are never persisted; only their SHA-256 hashes are stored.
CREATE TABLE IF NOT EXISTS subscription_quota_reset_operations (
    id BIGSERIAL PRIMARY KEY,
    actor_user_id BIGINT NOT NULL,
    idempotency_key_hash VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'processing',
    cutoff_at TIMESTAMPTZ,
    captured_at TIMESTAMPTZ,
    summary JSONB NOT NULL DEFAULT '{}'::jsonb,
    target_pairs JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT subscription_quota_reset_operations_status_check
        CHECK (status IN ('processing', 'succeeded'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_subscription_quota_reset_operations_actor_key
    ON subscription_quota_reset_operations (actor_user_id, idempotency_key_hash);

CREATE INDEX IF NOT EXISTS idx_subscription_quota_reset_operations_created_at
    ON subscription_quota_reset_operations (created_at DESC);
