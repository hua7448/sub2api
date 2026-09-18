-- Provider hall request facts and collector bookkeeping. Facts are written by
-- every gateway instance; aggregation (232) reads them. No user-visible data.
CREATE TABLE provider_hall_collector_epochs (
    id bigserial PRIMARY KEY,
    node_id varchar(128) NOT NULL,
    build_version varchar(64) NOT NULL DEFAULT '',
    algorithm_version smallint NOT NULL,
    registered_at timestamptz NOT NULL DEFAULT now(),
    heartbeat_at timestamptz NOT NULL DEFAULT now(),
    exited_at timestamptz,
    exit_reason varchar(32) NOT NULL DEFAULT '',
    persisted_seq bigint NOT NULL DEFAULT 0,
    confirmed_at timestamptz,
    overflowed boolean NOT NULL DEFAULT false
);
CREATE INDEX provider_hall_collector_epochs_node_idx ON provider_hall_collector_epochs(node_id, registered_at DESC);

CREATE TABLE provider_hall_coverage_gaps (
    id bigserial PRIMARY KEY,
    node_id varchar(128) NOT NULL,
    epoch_id bigint NOT NULL,
    scope varchar(16) NOT NULL CHECK (scope IN ('collection', 'billing')),
    started_at timestamptz NOT NULL,
    ended_at timestamptz,
    reason varchar(32) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CHECK (ended_at IS NULL OR ended_at >= started_at)
);
CREATE INDEX provider_hall_coverage_gaps_range_idx ON provider_hall_coverage_gaps(started_at, ended_at);

CREATE TABLE provider_hall_requests (
    trace_id uuid PRIMARY KEY,
    node_id varchar(128) NOT NULL,
    epoch_id bigint NOT NULL,
    last_seq bigint NOT NULL,
    group_id bigint NOT NULL,
    profile_id bigint,
    api_key_id bigint NOT NULL,
    protocol varchar(32) NOT NULL CHECK (protocol IN ('responses', 'chat_completions', 'messages')),
    requested_model varchar(200) NOT NULL,
    upstream_model varchar(200) NOT NULL DEFAULT '',
    response_model varchar(200) NOT NULL DEFAULT '',
    platform varchar(32) NOT NULL DEFAULT '',
    source varchar(16) NOT NULL CHECK (source IN ('user', 'probe', 'verification')),
    sample_id bigint,
    stream boolean NOT NULL DEFAULT false,
    started_at timestamptz NOT NULL,
    first_content_at timestamptz,
    ended_at timestamptz,
    ttft_ms integer CHECK (ttft_ms > 0),
    submissions integer NOT NULL DEFAULT 0 CHECK (submissions >= 0),
    outcome varchar(16) NOT NULL DEFAULT 'pending' CHECK (outcome IN ('pending', 'success', 'failed', 'excluded')),
    exclusion_reason varchar(32) NOT NULL DEFAULT '',
    usage_known boolean NOT NULL DEFAULT false,
    input_tokens bigint,
    output_tokens bigint,
    cache_read_tokens bigint,
    cache_creation_tokens bigint,
    billing_request_id varchar(255),
    billing_api_key_id bigint,
    request_fingerprint varchar(64),
    billing_status varchar(16) NOT NULL DEFAULT 'pending'
        CHECK (billing_status IN ('pending', 'applied', 'duplicate', 'failed', 'uncertain', 'not_applicable')),
    billing_mode varchar(16),
    is_subscription boolean,
    rate_multiplier numeric(12,6),
    input_base_cost numeric(24,10),
    total_base_cost numeric(24,10),
    actual_cost numeric(20,8),
    billed_at timestamptz,
    algorithm_version smallint NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX provider_hall_requests_window_idx ON provider_hall_requests(group_id, profile_id, ended_at)
    WHERE ended_at IS NOT NULL AND source = 'user';
CREATE INDEX provider_hall_requests_reconcile_idx ON provider_hall_requests(started_at)
    WHERE outcome = 'success' AND billing_status IN ('pending', 'uncertain');
CREATE INDEX provider_hall_requests_started_idx ON provider_hall_requests(started_at);
CREATE INDEX provider_hall_requests_sample_idx ON provider_hall_requests(sample_id) WHERE sample_id IS NOT NULL;

CREATE TABLE provider_hall_dirty_buckets (
    id bigserial PRIMARY KEY,
    group_id bigint NOT NULL,
    profile_id bigint NOT NULL,
    minute timestamptz NOT NULL,
    reason varchar(32) NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (group_id, profile_id, minute)
);
