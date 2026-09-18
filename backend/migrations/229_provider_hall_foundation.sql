-- Configuration only. Runtime workers and public display remain disabled.
CREATE TABLE provider_hall_config (
    id bigint PRIMARY KEY DEFAULT 1 CONSTRAINT provider_hall_config_singleton CHECK (id = 1),
    version bigint NOT NULL DEFAULT 1 CHECK (version > 0),
    collection_enabled boolean NOT NULL DEFAULT false,
    display_enabled boolean NOT NULL DEFAULT false,
    tasks_enabled boolean NOT NULL DEFAULT false,
    default_model varchar(200) NOT NULL DEFAULT '',
    default_protocol varchar(32) NOT NULL DEFAULT 'responses' CHECK (default_protocol IN ('responses', 'chat_completions', 'messages')),
    default_range varchar(8) NOT NULL DEFAULT '6h' CHECK (default_range IN ('6h', '24h', '7d', '30d')),
    gateway_origin varchar(2048) NOT NULL DEFAULT '',
    operator_user_id bigint REFERENCES users(id) ON DELETE SET NULL,
    daily_budget numeric(20,8) NOT NULL DEFAULT 0 CONSTRAINT provider_hall_config_budget_nonnegative CHECK (daily_budget >= 0),
    expected_nodes jsonb NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(expected_nodes) = 'array'),
    updated_at timestamptz NOT NULL DEFAULT now(),
    updated_by bigint REFERENCES users(id) ON DELETE SET NULL
);

INSERT INTO provider_hall_config (id) VALUES (1);

CREATE TABLE provider_hall_groups (
    group_id bigint PRIMARY KEY REFERENCES groups(id) ON DELETE CASCADE,
    version bigint NOT NULL DEFAULT 1 CHECK (version > 0),
    listed boolean NOT NULL DEFAULT false,
    display_name varchar(100) NOT NULL DEFAULT '',
    description varchar(2000) NOT NULL DEFAULT '',
    display_order integer NOT NULL DEFAULT 0,
    updated_at timestamptz NOT NULL DEFAULT now(),
    updated_by bigint REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE provider_hall_profiles (
    id bigserial PRIMARY KEY,
    version bigint NOT NULL DEFAULT 1 CHECK (version > 0),
    model varchar(200) NOT NULL CHECK (length(btrim(model)) > 0),
    protocol varchar(32) NOT NULL CHECK (protocol IN ('responses', 'chat_completions', 'messages')),
    supports_tools boolean NOT NULL DEFAULT false,
    output_limit integer NOT NULL DEFAULT 256 CONSTRAINT provider_hall_profile_output_limit CHECK (output_limit BETWEEN 1 AND 1024),
    model_aliases jsonb NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(model_aliases) = 'array'),
    reference_input_price numeric(24,10) CONSTRAINT provider_hall_profile_input_price CHECK (reference_input_price >= 0),
    reference_cache_price numeric(24,10) CONSTRAINT provider_hall_profile_cache_price CHECK (reference_cache_price >= 0),
    reference_cache_rate numeric(11,10) CONSTRAINT provider_hall_profile_cache_rate CHECK (reference_cache_rate BETWEEN 0 AND 1),
    reference_confirmed_at timestamptz,
    updated_at timestamptz NOT NULL DEFAULT now(),
    updated_by bigint REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE (model, protocol),
    CHECK ((reference_input_price IS NULL AND reference_cache_price IS NULL AND reference_cache_rate IS NULL AND reference_confirmed_at IS NULL)
        OR (reference_input_price IS NOT NULL AND reference_cache_price IS NOT NULL AND reference_cache_rate IS NOT NULL AND reference_confirmed_at IS NOT NULL))
);
