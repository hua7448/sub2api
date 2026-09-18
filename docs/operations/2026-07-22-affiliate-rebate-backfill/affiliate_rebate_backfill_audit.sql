\set ON_ERROR_STOP on

WITH cfg AS (
    SELECT
        COALESCE(MAX(value::numeric) FILTER (WHERE key = 'affiliate_rebate_rate'), 20) AS global_rate,
        COALESCE(MAX(value::int) FILTER (WHERE key = 'affiliate_min_qualified_invitees'), 1) AS min_qualified,
        COALESCE(MAX(value::int) FILTER (WHERE key = 'affiliate_rebate_duration_days'), 0) AS duration_days
    FROM settings
    WHERE key IN (
        'affiliate_rebate_rate',
        'affiliate_min_qualified_invitees',
        'affiliate_rebate_duration_days'
    )
), direct_redeems AS (
    SELECT
        rc.id AS redeem_code_id,
        rc.used_by AS invitee_id,
        invitee_aff.inviter_id,
        rc.type AS reward_source,
        rc.value AS source_amount,
        rc.group_id,
        CASE
            WHEN rc.type = 'subscription' AND rc.validity_days = 0 THEN 30
            ELSE rc.validity_days
        END AS source_days,
        rc.used_at,
        invitee_aff.created_at AS invited_at,
        LEAST(100::numeric, GREATEST(0::numeric,
            COALESCE(inviter_aff.aff_rebate_rate_percent, cfg.global_rate)
        )) AS rebate_rate,
        cfg.min_qualified,
        cfg.duration_days
    FROM redeem_codes rc
    JOIN user_affiliates invitee_aff ON invitee_aff.user_id = rc.used_by
    JOIN user_affiliates inviter_aff ON inviter_aff.user_id = invitee_aff.inviter_id
    CROSS JOIN cfg
    WHERE rc.status = 'used'
      AND rc.used_by IS NOT NULL
      AND rc.used_at IS NOT NULL
      AND invitee_aff.inviter_id IS NOT NULL
      AND invitee_aff.created_at <= rc.used_at
      AND (
          (rc.type = 'balance' AND rc.value > 0)
          OR
          (rc.type = 'subscription' AND rc.group_id IS NOT NULL AND rc.validity_days >= 0)
      )
      -- Payment fulfillment deliberately suppresses redeem-code rewards and
      -- writes an order-sourced monetary rebate instead.
      AND NOT EXISTS (
          SELECT 1
          FROM payment_orders po
          WHERE po.recharge_code = rc.code
      )
), event_eligibility AS (
    SELECT
        event.*,
        (
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
        ) AS prior_qualified
    FROM direct_redeems event
), expected AS (
    SELECT
        eligible.*,
        CASE
            WHEN eligible.reward_source = 'balance'
                THEN ROUND(eligible.source_amount * eligible.rebate_rate / 100, 8)
            ELSE 0::numeric
        END AS expected_amount,
        CASE
            WHEN eligible.reward_source = 'subscription'
                THEN FLOOR(eligible.source_days * eligible.rebate_rate / 100)::int
            ELSE 0
        END AS expected_days,
        CASE
            WHEN eligible.reward_source = 'balance' THEN 'accrue'
            ELSE 'accrue_days'
        END AS expected_action
    FROM event_eligibility eligible
    WHERE eligible.prior_qualified >= eligible.min_qualified
      AND (
          eligible.duration_days = 0
          OR eligible.used_at <= eligible.invited_at + make_interval(days => eligible.duration_days)
      )
), auditable AS (
    SELECT
        expected.*,
        source_ledger.id AS source_ledger_id,
        legacy_ledger.id AS legacy_ledger_id
    FROM expected
    LEFT JOIN LATERAL (
        SELECT ledger.id
        FROM user_affiliate_ledger ledger
        WHERE ledger.action = expected.expected_action
          AND ledger.source_redeem_code_id = expected.redeem_code_id
        ORDER BY ledger.id
        LIMIT 1
    ) source_ledger ON true
    LEFT JOIN LATERAL (
        SELECT ledger.id
        FROM user_affiliate_ledger ledger
        WHERE source_ledger.id IS NULL
          AND ledger.source_redeem_code_id IS NULL
          AND ledger.user_id = expected.inviter_id
          AND ledger.source_user_id = expected.invitee_id
          AND ledger.action = expected.expected_action
          AND ledger.created_at BETWEEN expected.used_at - interval '5 seconds'
                                    AND expected.used_at + interval '5 seconds'
          AND (
              (expected.expected_action = 'accrue'
               AND ABS(ledger.amount - expected.expected_amount) < 0.00000001)
              OR
              (expected.expected_action = 'accrue_days'
               AND ledger.group_id = expected.group_id
               AND ledger.days = expected.expected_days)
          )
        ORDER BY ABS(EXTRACT(epoch FROM (ledger.created_at - expected.used_at))), ledger.id
        LIMIT 1
    ) legacy_ledger ON true
    WHERE expected.expected_amount > 0 OR expected.expected_days > 0
)
SELECT
    reward_source,
    COUNT(*) AS eligible_events,
    COUNT(*) FILTER (WHERE source_ledger_id IS NOT NULL OR legacy_ledger_id IS NOT NULL) AS already_rewarded,
    COUNT(*) FILTER (WHERE source_ledger_id IS NULL AND legacy_ledger_id IS NULL) AS missing_events,
    COALESCE(SUM(expected_amount) FILTER (WHERE source_ledger_id IS NULL AND legacy_ledger_id IS NULL), 0) AS missing_amount,
    COALESCE(SUM(expected_days) FILTER (WHERE source_ledger_id IS NULL AND legacy_ledger_id IS NULL), 0) AS missing_days
FROM auditable
GROUP BY reward_source
ORDER BY reward_source;

WITH cfg AS (
    SELECT
        COALESCE(MAX(value::numeric) FILTER (WHERE key = 'affiliate_rebate_rate'), 20) AS global_rate,
        COALESCE(MAX(value::int) FILTER (WHERE key = 'affiliate_min_qualified_invitees'), 1) AS min_qualified,
        COALESCE(MAX(value::int) FILTER (WHERE key = 'affiliate_rebate_duration_days'), 0) AS duration_days
    FROM settings
    WHERE key IN ('affiliate_rebate_rate', 'affiliate_min_qualified_invitees', 'affiliate_rebate_duration_days')
), direct_redeems AS (
    SELECT rc.id AS redeem_code_id, rc.used_by AS invitee_id, ia.inviter_id,
           rc.type AS reward_source, rc.value AS source_amount, rc.group_id,
           CASE WHEN rc.type = 'subscription' AND rc.validity_days = 0 THEN 30 ELSE rc.validity_days END AS source_days,
           rc.used_at, ia.created_at AS invited_at,
           LEAST(100::numeric, GREATEST(0::numeric, COALESCE(inviter.aff_rebate_rate_percent, cfg.global_rate))) AS rebate_rate,
           cfg.min_qualified, cfg.duration_days
    FROM redeem_codes rc
    JOIN user_affiliates ia ON ia.user_id = rc.used_by
    JOIN user_affiliates inviter ON inviter.user_id = ia.inviter_id
    CROSS JOIN cfg
    WHERE rc.status = 'used' AND rc.used_by IS NOT NULL AND rc.used_at IS NOT NULL
      AND ia.inviter_id IS NOT NULL AND ia.created_at <= rc.used_at
      AND ((rc.type = 'balance' AND rc.value > 0)
           OR (rc.type = 'subscription' AND rc.group_id IS NOT NULL AND rc.validity_days >= 0))
      AND NOT EXISTS (SELECT 1 FROM payment_orders po WHERE po.recharge_code = rc.code)
), expected AS (
    SELECT event.*,
           CASE WHEN event.reward_source = 'balance' THEN ROUND(event.source_amount * event.rebate_rate / 100, 8) ELSE 0::numeric END AS expected_amount,
           CASE WHEN event.reward_source = 'subscription' THEN FLOOR(event.source_days * event.rebate_rate / 100)::int ELSE 0 END AS expected_days,
           CASE WHEN event.reward_source = 'balance' THEN 'accrue' ELSE 'accrue_days' END AS expected_action
    FROM direct_redeems event
    WHERE (
        SELECT COUNT(DISTINCT other_aff.user_id)
        FROM user_affiliates other_aff
        WHERE other_aff.inviter_id = event.inviter_id
          AND other_aff.user_id <> event.invitee_id
          AND EXISTS (
              SELECT 1 FROM redeem_codes qc
              WHERE qc.used_by = other_aff.user_id AND qc.status = 'used' AND qc.used_at IS NOT NULL
                AND other_aff.created_at <= qc.used_at AND qc.used_at < event.used_at
                AND ((qc.type = 'balance' AND qc.value > 0)
                     OR (qc.type = 'subscription' AND qc.validity_days >= 0))
          )
    ) >= event.min_qualified
      AND (event.duration_days = 0 OR event.used_at <= event.invited_at + make_interval(days => event.duration_days))
), missing AS (
    SELECT expected.*
    FROM expected
    WHERE (expected.expected_amount > 0 OR expected.expected_days > 0)
      AND NOT EXISTS (
          SELECT 1 FROM user_affiliate_ledger ledger
          WHERE ledger.action = expected.expected_action
            AND ledger.source_redeem_code_id = expected.redeem_code_id
      )
      AND NOT EXISTS (
          SELECT 1 FROM user_affiliate_ledger ledger
          WHERE ledger.source_redeem_code_id IS NULL
            AND ledger.user_id = expected.inviter_id
            AND ledger.source_user_id = expected.invitee_id
            AND ledger.action = expected.expected_action
            AND ledger.created_at BETWEEN expected.used_at - interval '5 seconds' AND expected.used_at + interval '5 seconds'
            AND ((expected.expected_action = 'accrue' AND ABS(ledger.amount - expected.expected_amount) < 0.00000001)
                 OR (expected.expected_action = 'accrue_days' AND ledger.group_id = expected.group_id AND ledger.days = expected.expected_days))
      )
)
SELECT m.inviter_id, inviter.email AS inviter_email,
       m.invitee_id, invitee.email AS invitee_email,
       m.redeem_code_id, m.reward_source, m.group_id,
       m.source_amount, m.source_days, m.rebate_rate,
       m.expected_amount, m.expected_days, m.used_at
FROM missing m
JOIN users inviter ON inviter.id = m.inviter_id
JOIN users invitee ON invitee.id = m.invitee_id
ORDER BY m.used_at, m.redeem_code_id;
