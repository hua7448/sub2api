-- 1) New table: per-group subscription days rebate pool for inviters.
CREATE TABLE IF NOT EXISTS user_affiliate_subscription_days (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    pending_days INT NOT NULL DEFAULT 0,
    frozen_days INT NOT NULL DEFAULT 0,
    history_days INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, group_id)
);

CREATE INDEX IF NOT EXISTS idx_uasd_user_id ON user_affiliate_subscription_days(user_id);

COMMENT ON TABLE user_affiliate_subscription_days IS '邀请返利：按分组累积的订阅天数池';
COMMENT ON COLUMN user_affiliate_subscription_days.pending_days IS '可提取天数';
COMMENT ON COLUMN user_affiliate_subscription_days.frozen_days IS '冻结中天数（待解冻）';
COMMENT ON COLUMN user_affiliate_subscription_days.history_days IS '历史累计天数';

-- 2) Extend affiliate ledger to support subscription-days actions.
ALTER TABLE user_affiliate_ledger
    ADD COLUMN IF NOT EXISTS group_id BIGINT NULL,
    ADD COLUMN IF NOT EXISTS days INT NULL;

COMMENT ON COLUMN user_affiliate_ledger.group_id IS '订阅天数返利关联的分组ID';
COMMENT ON COLUMN user_affiliate_ledger.days IS '订阅天数（accrue_days/transfer_days 使用）';
