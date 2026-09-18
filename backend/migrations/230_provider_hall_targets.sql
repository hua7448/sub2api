-- Permanent ID-only provenance must survive key/user/group deletion.
CREATE TABLE provider_hall_probe_keys (
    api_key_id bigint PRIMARY KEY CHECK (api_key_id > 0),
    operator_user_id bigint NOT NULL CHECK (operator_user_id > 0),
    group_id bigint NOT NULL CHECK (group_id > 0),
    registered_by bigint NOT NULL CHECK (registered_by > 0),
    registered_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE provider_hall_targets (
    id bigserial PRIMARY KEY,
    group_id bigint NOT NULL REFERENCES provider_hall_groups(group_id) ON DELETE CASCADE,
    profile_id bigint NOT NULL REFERENCES provider_hall_profiles(id) ON DELETE RESTRICT,
    probe_key_id bigint REFERENCES provider_hall_probe_keys(api_key_id) ON DELETE RESTRICT,
    enabled boolean NOT NULL DEFAULT false,
    probe_interval_seconds integer NOT NULL DEFAULT 300
        CONSTRAINT provider_hall_target_probe_interval CHECK (probe_interval_seconds BETWEEN 60 AND 86400),
    verification_interval_seconds integer NOT NULL DEFAULT 86400
        CONSTRAINT provider_hall_target_verification_interval CHECK (verification_interval_seconds BETWEEN 3600 AND 604800),
    version bigint NOT NULL DEFAULT 1 CHECK (version > 0),
    updated_at timestamptz NOT NULL DEFAULT now(),
    updated_by bigint REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE (group_id, profile_id),
    CONSTRAINT provider_hall_target_enabled_key CHECK (NOT enabled OR probe_key_id IS NOT NULL)
);
CREATE INDEX provider_hall_targets_probe_key_id_idx ON provider_hall_targets(probe_key_id);
