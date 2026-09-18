\set ON_ERROR_STOP on

BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '60s';

CREATE TEMP TABLE affiliate_rebate_missing_subscription_rewards ON COMMIT DROP AS
WITH cfg AS (
    SELECT
        COALESCE(MAX(value::numeric) FILTER (WHERE key = 'affiliate_rebate_rate'), 20) AS global_rate,
        COALESCE(MAX(value::int) FILTER (WHERE key = 'affiliate_min_qualified_invitees'), 1) AS min_qualified,
        COALESCE(MAX(value::int) FILTER (WHERE key = 'affiliate_rebate_duration_days'), 0) AS duration_days,
        COALESCE(MAX(value::int) FILTER (WHERE key = 'affiliate_rebate_freeze_hours'), 0) AS freeze_hours
    FROM settings
    WHERE key IN (
        'affiliate_rebate_rate',
        'affiliate_min_qualified_invitees',
        'affiliate_rebate_duration_days',
        'affiliate_rebate_freeze_hours'
    )
), direct_subscription_redeems AS (
    SELECT
        rc.id AS redeem_code_id,
        rc.used_by AS invitee_id,
        invitee_aff.inviter_id,
        rc.group_id,
        CASE WHEN rc.validity_days = 0 THEN 30 ELSE rc.validity_days END AS source_days,
        rc.used_at,
        invitee_aff.created_at AS invited_at,
        LEAST(100::numeric, GREATEST(0::numeric,
            COALESCE(inviter_aff.aff_rebate_rate_percent, cfg.global_rate)
        )) AS rebate_rate,
        cfg.min_qualified,
        cfg.duration_days,
        cfg.freeze_hours
    FROM redeem_codes rc
    JOIN user_affiliates invitee_aff ON invitee_aff.user_id = rc.used_by
    JOIN user_affiliates inviter_aff ON inviter_aff.user_id = invitee_aff.inviter_id
    CROSS JOIN cfg
    WHERE rc.status = 'used'
      AND rc.type = 'subscription'
      AND rc.used_by IS NOT NULL
      AND rc.used_at IS NOT NULL
      AND rc.group_id IS NOT NULL
      AND rc.validity_days >= 0
      AND invitee_aff.inviter_id IS NOT NULL
      AND invitee_aff.created_at <= rc.used_at
      AND NOT EXISTS (
          SELECT 1
          FROM payment_orders po
          WHERE po.recharge_code = rc.code
      )
), expected AS (
    SELECT
        event.redeem_code_id,
        event.invitee_id,
        event.inviter_id,
        event.group_id,
        FLOOR(event.source_days * event.rebate_rate / 100)::int AS expected_days,
        event.used_at,
        event.freeze_hours
    FROM direct_subscription_redeems event
    WHERE (
        SELECT COUNT(DISTINCT other_aff.user_id)
        FROM user_affiliates other_aff
        WHERE other_aff.inviter_id = event.inviter_id
          AND other_aff.user_id <> event.invitee_id
          AND EXISTS (
              SELECT 1
              FROM redeem_codes qualifying_code
              WHERE qualifying_code.used_by = other_aff.user_id
                AND qualifying_code.status = 'used'
                AND qualifying_code.used_at IS NOT NULL
                AND other_aff.created_at <= qualifying_code.used_at
                AND qualifying_code.used_at < event.used_at
                AND (
                    (qualifying_code.type = 'balance' AND qualifying_code.value > 0)
                    OR
                    (qualifying_code.type = 'subscription' AND qualifying_code.validity_days >= 0)
                )
          )
    ) >= event.min_qualified
      AND (
          event.duration_days = 0
          OR event.used_at <= event.invited_at + make_interval(days => event.duration_days)
      )
), missing AS (
    SELECT expected.*
    FROM expected
    WHERE expected.expected_days > 0
      AND NOT EXISTS (
          SELECT 1
          FROM user_affiliate_ledger ledger
          WHERE ledger.action = 'accrue_days'
            AND ledger.source_redeem_code_id = expected.redeem_code_id
      )
      AND NOT EXISTS (
          SELECT 1
          FROM user_affiliate_ledger ledger
          WHERE ledger.source_redeem_code_id IS NULL
            AND ledger.user_id = expected.inviter_id
            AND ledger.source_user_id = expected.invitee_id
            AND ledger.action = 'accrue_days'
            AND ledger.group_id = expected.group_id
            AND ledger.days = expected.expected_days
            AND ledger.created_at BETWEEN expected.used_at - interval '5 seconds'
                                      AND expected.used_at + interval '5 seconds'
      )
)
SELECT * FROM missing;

DO $preflight$
DECLARE
    candidate_count bigint;
    candidate_days bigint;
BEGIN
    SELECT COUNT(*), COALESCE(SUM(expected_days), 0)
    INTO candidate_count, candidate_days
    FROM affiliate_rebate_missing_subscription_rewards;

    IF candidate_count <> 65 OR candidate_days <> 195 THEN
        RAISE EXCEPTION 'affiliate backfill preflight changed: count=%, days=%; expected count=65, days=195',
            candidate_count, candidate_days;
    END IF;
END
$preflight$;

SELECT
    group_id,
    COUNT(*) AS reward_count,
    SUM(expected_days) AS reward_days
FROM affiliate_rebate_missing_subscription_rewards
GROUP BY group_id
ORDER BY group_id;

CREATE TEMP TABLE affiliate_rebate_inserted_rewards ON COMMIT DROP AS
WITH inserted AS (
    INSERT INTO user_affiliate_ledger (
        user_id,
        action,
        amount,
        source_user_id,
        source_redeem_code_id,
        group_id,
        days,
        frozen_until,
        created_at,
        updated_at
    )
    SELECT
        candidate.inviter_id,
        'accrue_days',
        0,
        candidate.invitee_id,
        candidate.redeem_code_id,
        candidate.group_id,
        candidate.expected_days,
        CASE
            WHEN candidate.freeze_hours > 0
                THEN NOW() + make_interval(hours => candidate.freeze_hours)
            ELSE NULL
        END,
        NOW(),
        NOW()
    FROM affiliate_rebate_missing_subscription_rewards candidate
    ORDER BY candidate.used_at, candidate.redeem_code_id
    ON CONFLICT DO NOTHING
    RETURNING id, user_id, group_id, days, frozen_until, source_redeem_code_id
)
SELECT * FROM inserted;

DO $insert_check$
DECLARE
    inserted_count bigint;
    inserted_days bigint;
BEGIN
    SELECT COUNT(*), COALESCE(SUM(days), 0)
    INTO inserted_count, inserted_days
    FROM affiliate_rebate_inserted_rewards;

    IF inserted_count <> 65 OR inserted_days <> 195 THEN
        RAISE EXCEPTION 'affiliate backfill insert changed: count=%, days=%; expected count=65, days=195',
            inserted_count, inserted_days;
    END IF;
END
$insert_check$;

INSERT INTO user_affiliate_subscription_days (
    user_id,
    group_id,
    pending_days,
    frozen_days,
    history_days,
    created_at,
    updated_at
)
SELECT
    inserted.user_id,
    inserted.group_id,
    SUM(CASE WHEN inserted.frozen_until IS NULL THEN inserted.days ELSE 0 END)::int,
    SUM(CASE WHEN inserted.frozen_until IS NOT NULL THEN inserted.days ELSE 0 END)::int,
    SUM(inserted.days)::int,
    NOW(),
    NOW()
FROM affiliate_rebate_inserted_rewards inserted
GROUP BY inserted.user_id, inserted.group_id
ON CONFLICT (user_id, group_id) DO UPDATE
SET pending_days = user_affiliate_subscription_days.pending_days + EXCLUDED.pending_days,
    frozen_days = user_affiliate_subscription_days.frozen_days + EXCLUDED.frozen_days,
    history_days = user_affiliate_subscription_days.history_days + EXCLUDED.history_days,
    updated_at = NOW();

SELECT
    COUNT(*) AS inserted_rewards,
    COUNT(DISTINCT user_id) AS rewarded_inviters,
    SUM(days) AS inserted_days
FROM affiliate_rebate_inserted_rewards;

COMMIT;
