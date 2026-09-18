-- Provider hall minute metrics, exact latency frequencies, hourly snapshots and
-- aggregator state. Written only by the primary aggregator; read by the user
-- and admin APIs. Facts stay in 231; nothing here is user-editable.
CREATE TABLE provider_hall_metrics_1m (
    group_id bigint NOT NULL,
    profile_id bigint NOT NULL,
    minute timestamptz NOT NULL,
    algorithm_version smallint NOT NULL,
    success_count int NOT NULL DEFAULT 0,
    failed_count int NOT NULL DEFAULT 0,
    submissions int NOT NULL DEFAULT 0,
    excluded_count int NOT NULL DEFAULT 0,
    ttft_sample_count int NOT NULL DEFAULT 0,
    usage_success_count int NOT NULL DEFAULT 0,
    input_tokens bigint NOT NULL DEFAULT 0,
    cache_read_tokens bigint NOT NULL DEFAULT 0,
    cache_creation_tokens bigint NOT NULL DEFAULT 0,
    output_tokens bigint NOT NULL DEFAULT 0,
    billed_count int NOT NULL DEFAULT 0,
    billed_subscription_count int NOT NULL DEFAULT 0,
    billed_input_tokens bigint NOT NULL DEFAULT 0,
    billed_input_cost numeric(24,10) NOT NULL DEFAULT 0,
    billing_pending_count int NOT NULL DEFAULT 0,
    billing_uncertain_count int NOT NULL DEFAULT 0,
    coverage varchar(16) NOT NULL DEFAULT 'complete'
        CHECK (coverage IN ('complete', 'collection_gap', 'node_unconfirmed', 'version_mismatch')),
    billing_coverage varchar(16) NOT NULL DEFAULT 'complete'
        CHECK (billing_coverage IN ('complete', 'billing_gap', 'pending')),
    computed_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (group_id, profile_id, minute, algorithm_version)
);
CREATE INDEX provider_hall_metrics_1m_minute_idx ON provider_hall_metrics_1m(minute);

CREATE TABLE provider_hall_latency_counts_1m (
    group_id bigint NOT NULL,
    profile_id bigint NOT NULL,
    minute timestamptz NOT NULL,
    algorithm_version smallint NOT NULL,
    ttft_ms int NOT NULL CHECK (ttft_ms > 0),
    count int NOT NULL CHECK (count > 0),
    PRIMARY KEY (group_id, profile_id, minute, algorithm_version, ttft_ms)
);
CREATE INDEX provider_hall_latency_counts_1m_minute_idx ON provider_hall_latency_counts_1m(minute);

CREATE TABLE provider_hall_snapshots (
    id bigserial PRIMARY KEY,
    group_id bigint NOT NULL,
    profile_id bigint NOT NULL,            -- 0 = merged over the group's enabled profiles
    window_end timestamptz NOT NULL,       -- T, window [T-60m, T)
    tier smallint NOT NULL CHECK (tier IN (1, 5)),
    algorithm_version smallint NOT NULL,
    ttft_sample_count int NOT NULL,
    ttft_fast95_mean_ms numeric(12,3),
    ttft_p90_ms int,
    success_count int NOT NULL,
    failed_count int NOT NULL,
    submissions int NOT NULL,
    success_rate numeric(11,10),
    usage_success_count int NOT NULL,
    input_tokens bigint NOT NULL,
    cache_read_tokens bigint NOT NULL,
    cache_creation_tokens bigint NOT NULL,
    cache_rate numeric(11,10),
    billed_count int NOT NULL,
    billed_input_tokens bigint NOT NULL,
    billed_input_cost numeric(24,10) NOT NULL,
    historical_price numeric(24,10),
    subscription_share numeric(11,10),
    coverage varchar(16) NOT NULL,
    billing_coverage varchar(16) NOT NULL,
    coverage_reason varchar(32) NOT NULL DEFAULT '',
    computed_at timestamptz NOT NULL DEFAULT now(),
    computed_version int NOT NULL DEFAULT 1,
    UNIQUE (group_id, profile_id, window_end, algorithm_version)
);
CREATE INDEX provider_hall_snapshots_window_idx ON provider_hall_snapshots(window_end, tier);

CREATE TABLE provider_hall_aggregator_state (
    id smallint PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    algorithm_version smallint NOT NULL DEFAULT 1,
    last_window_end timestamptz,
    last_run_at timestamptz,
    last_error varchar(2000) NOT NULL DEFAULT '',
    updated_at timestamptz NOT NULL DEFAULT now()
);
INSERT INTO provider_hall_aggregator_state (id) VALUES (1);
