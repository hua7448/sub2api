-- Provider hall probe/verification jobs, samples, reports and spend ledger.
-- Written by the primary runner; the collector mirrors billing outcomes into
-- provider_hall_spend / provider_hall_samples.billing_status.
CREATE TABLE provider_hall_jobs (
    id bigserial PRIMARY KEY,
    kind varchar(16) NOT NULL CHECK (kind IN ('probe', 'verification')),
    target_id bigint NOT NULL,   -- no FK: history survives target/group deletion
    group_id bigint NOT NULL,
    profile_id bigint NOT NULL,
    config_snapshot jsonb NOT NULL,   -- {profile:{id,version,model,protocol,supports_tools,output_limit,model_aliases}, target_version, probe_key_id, operator_user_id, gateway_origin}
    slot_at timestamptz,
    not_before timestamptz,
    idempotency_key varchar(128),
    requested_by bigint,
    status varchar(16) NOT NULL DEFAULT 'queued' CHECK (status IN ('queued', 'running', 'succeeded', 'failed', 'cancelled', 'unknown')),
    lease_owner varchar(128),
    lease_until timestamptz,
    attempts int NOT NULL DEFAULT 0,
    budget_day date,
    error_code varchar(64) NOT NULL DEFAULT '',
    error_message varchar(2000) NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    started_at timestamptz,
    finished_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX provider_hall_jobs_slot_uq ON provider_hall_jobs(target_id, kind, slot_at) WHERE slot_at IS NOT NULL;
CREATE UNIQUE INDEX provider_hall_jobs_idempotency_uq ON provider_hall_jobs(idempotency_key) WHERE idempotency_key IS NOT NULL;
CREATE INDEX provider_hall_jobs_claim_idx ON provider_hall_jobs(status, not_before, created_at) WHERE status IN ('queued', 'running');
CREATE INDEX provider_hall_jobs_target_idx ON provider_hall_jobs(target_id, kind, created_at DESC);
CREATE INDEX provider_hall_jobs_group_idx ON provider_hall_jobs(group_id, created_at DESC);

CREATE TABLE provider_hall_samples (
    id bigserial PRIMARY KEY,
    job_id bigint NOT NULL REFERENCES provider_hall_jobs(id) ON DELETE CASCADE,
    test_id varchar(32) NOT NULL,   -- probe|arith|json|tool
    seq smallint NOT NULL,
    trace_id uuid NOT NULL,
    status varchar(16) NOT NULL DEFAULT 'prepared' CHECK (status IN ('prepared', 'dispatched', 'received', 'uncertain')),
    prepared_at timestamptz NOT NULL DEFAULT now(),
    dispatched_at timestamptz,
    received_at timestamptz,
    client_request_id varchar(64),
    http_status int,
    ttft_ms int,
    total_ms int,
    generation_ms int,
    input_tokens int,
    output_tokens int,
    response_model varchar(200),
    result varchar(16) CHECK (result IN ('passed', 'failed', 'error', 'model_mismatch')),
    error_code varchar(64) NOT NULL DEFAULT '',
    detail jsonb NOT NULL DEFAULT '{}',   -- expected/actual truncated to 2KB, never user content
    billing_status varchar(16) NOT NULL DEFAULT 'pending' CHECK (billing_status IN ('pending', 'confirmed', 'uncertain', 'failed', 'unbilled')),
    actual_cost numeric(20,8),
    UNIQUE (job_id, test_id, seq),
    UNIQUE (trace_id)
);
CREATE INDEX provider_hall_samples_inflight_idx ON provider_hall_samples(status, dispatched_at) WHERE status IN ('dispatched', 'uncertain');
CREATE INDEX provider_hall_samples_received_idx ON provider_hall_samples(received_at) WHERE status = 'received';

CREATE TABLE provider_hall_verifications (
    job_id bigint PRIMARY KEY REFERENCES provider_hall_jobs(id) ON DELETE CASCADE,
    group_id bigint NOT NULL,
    profile_id bigint NOT NULL,
    target_id bigint NOT NULL,
    verdict varchar(16) NOT NULL CHECK (verdict IN ('passed', 'failed', 'suspected', 'insufficient')),
    execution_status varchar(16) NOT NULL CHECK (execution_status IN ('completed', 'partial', 'error')),
    reason_code varchar(64) NOT NULL DEFAULT '',
    summary jsonb NOT NULL,   -- {arithmetic:{passed,failed,error}, json:{...}, tool:{...}|null, model:{matched,mismatched,missing,seen:[]}}
    profile_version bigint NOT NULL,
    target_version bigint NOT NULL,
    completed_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL,
    stale boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX provider_hall_verifications_lookup_idx ON provider_hall_verifications(group_id, profile_id, completed_at DESC);

CREATE TABLE provider_hall_spend (
    billing_request_id varchar(255) NOT NULL,
    api_key_id bigint NOT NULL,
    job_id bigint,
    sample_id bigint,
    budget_day date NOT NULL,
    actual_cost numeric(20,8) NOT NULL DEFAULT 0,
    status varchar(16) NOT NULL CHECK (status IN ('confirmed', 'uncertain', 'failed')),
    confirmed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (billing_request_id, api_key_id)
);
CREATE INDEX provider_hall_spend_day_idx ON provider_hall_spend(budget_day, status);
