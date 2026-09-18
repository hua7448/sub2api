-- Track the redeem code that produced each affiliate reward. The site uses
-- redeem codes as its primary paid-entitlement flow, so payment-order-only
-- audit fields are not sufficient for a complete affiliate history.

ALTER TABLE user_affiliate_ledger
    ADD COLUMN IF NOT EXISTS source_redeem_code_id BIGINT NULL REFERENCES redeem_codes(id) ON DELETE SET NULL;

COMMENT ON COLUMN user_affiliate_ledger.source_redeem_code_id IS
    '产生该返利流水的余额或订阅兑换码；历史无法可靠匹配时为 NULL';

CREATE INDEX IF NOT EXISTS idx_user_affiliate_ledger_source_redeem_code_id
    ON user_affiliate_ledger(source_redeem_code_id)
    WHERE source_redeem_code_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_affiliate_ledger_redeem_reward_unique
    ON user_affiliate_ledger(action, source_redeem_code_id)
    WHERE source_redeem_code_id IS NOT NULL
      AND action IN ('accrue', 'accrue_days');

-- Best-effort historical backfill. Only accept a match when exactly one code
-- fits the user, reward kind/group and five-second redemption window.
WITH candidates AS (
    SELECT ual.id AS ledger_id,
           rc.id AS redeem_code_id,
           COUNT(*) OVER (PARTITION BY ual.id) AS ledger_match_count,
           COUNT(*) OVER (PARTITION BY rc.id, ual.action) AS code_match_count,
           ROW_NUMBER() OVER (
               PARTITION BY ual.id
               ORDER BY ABS(EXTRACT(EPOCH FROM (ual.created_at - rc.used_at))), rc.id
           ) AS match_rank
    FROM user_affiliate_ledger ual
    JOIN redeem_codes rc
      ON rc.used_by = ual.source_user_id
     AND rc.used_at IS NOT NULL
     AND ual.created_at BETWEEN rc.used_at - INTERVAL '5 seconds'
                            AND rc.used_at + INTERVAL '5 seconds'
     AND (
         (ual.action = 'accrue' AND rc.type = 'balance')
         OR
         (ual.action = 'accrue_days' AND rc.type = 'subscription' AND rc.group_id = ual.group_id)
     )
    WHERE ual.source_redeem_code_id IS NULL
      AND ual.action IN ('accrue', 'accrue_days')
)
UPDATE user_affiliate_ledger ual
SET source_redeem_code_id = candidates.redeem_code_id,
    updated_at = NOW()
FROM candidates
WHERE ual.id = candidates.ledger_id
  AND candidates.ledger_match_count = 1
  AND candidates.code_match_count = 1
  AND candidates.match_rank = 1
  AND NOT EXISTS (
      SELECT 1
      FROM user_affiliate_ledger existing
      WHERE existing.source_redeem_code_id = candidates.redeem_code_id
        AND existing.action = ual.action
  );
