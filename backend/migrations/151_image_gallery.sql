CREATE TABLE IF NOT EXISTS image_templates (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    prompt TEXT NOT NULL,
    category VARCHAR(80) NOT NULL DEFAULT '',
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    source VARCHAR(120) NOT NULL DEFAULT 'builtin',
    source_url TEXT NOT NULL DEFAULT '',
    author VARCHAR(120) NOT NULL DEFAULT '',
    license VARCHAR(120) NOT NULL DEFAULT '',
    variables JSONB NOT NULL DEFAULT '{}'::jsonb,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    featured BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS image_template_import_runs (
    id BIGSERIAL PRIMARY KEY,
    source VARCHAR(120) NOT NULL,
    source_url TEXT NOT NULL DEFAULT '',
    status VARCHAR(40) NOT NULL DEFAULT 'completed',
    imported_count INTEGER NOT NULL DEFAULT 0,
    skipped_count INTEGER NOT NULL DEFAULT 0,
    error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS image_generation_jobs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    api_key_id BIGINT REFERENCES api_keys(id) ON DELETE SET NULL,
    status VARCHAR(40) NOT NULL DEFAULT 'pending',
    model VARCHAR(120) NOT NULL,
    prompt TEXT NOT NULL,
    params JSONB NOT NULL DEFAULT '{}'::jsonb,
    error TEXT NOT NULL DEFAULT '',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS image_assets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_id BIGINT REFERENCES image_generation_jobs(id) ON DELETE SET NULL,
    usage_log_id BIGINT REFERENCES usage_logs(id) ON DELETE SET NULL,
    api_key_id BIGINT REFERENCES api_keys(id) ON DELETE SET NULL,
    path TEXT NOT NULL,
    thumb_path TEXT NOT NULL DEFAULT '',
    sha256 VARCHAR(64) NOT NULL,
    mime_type VARCHAR(80) NOT NULL DEFAULT 'image/png',
    width INTEGER NOT NULL DEFAULT 0,
    height INTEGER NOT NULL DEFAULT 0,
    size_bytes BIGINT NOT NULL DEFAULT 0,
    visibility VARCHAR(20) NOT NULL DEFAULT 'private',
    review_status VARCHAR(20) NOT NULL DEFAULT 'none',
    hidden BOOLEAN NOT NULL DEFAULT FALSE,
    prompt_snapshot TEXT NOT NULL DEFAULT '',
    model VARCHAR(120) NOT NULL DEFAULT '',
    params JSONB NOT NULL DEFAULT '{}'::jsonb,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS image_gallery_favorites (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset_id BIGINT NOT NULL REFERENCES image_assets(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, asset_id)
);

CREATE INDEX IF NOT EXISTS idx_image_templates_enabled_category ON image_templates(enabled, category);
CREATE INDEX IF NOT EXISTS idx_image_generation_jobs_user_created ON image_generation_jobs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_image_assets_user_created ON image_assets(user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_image_assets_public_created ON image_assets(created_at DESC) WHERE visibility = 'public' AND review_status = 'approved' AND hidden = FALSE AND deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_image_assets_review ON image_assets(review_status, created_at DESC) WHERE deleted_at IS NULL;

INSERT INTO image_templates (title, prompt, category, tags, source, featured)
VALUES
    ('产品海报', '为{主体}生成一张电商产品海报，干净背景，柔和棚拍光，突出材质和卖点，留出标题排版空间。', '商业设计', '["海报","产品","电商"]'::jsonb, 'builtin', TRUE),
    ('角色设定', '创建{主体}的角色设定图，包含清晰轮廓、服装细节、配色方案和半身姿态，风格为{风格}。', '角色', '["角色","设定","插画"]'::jsonb, 'builtin', TRUE),
    ('室内空间', '生成一个{场景}室内空间，真实摄影质感，自然光，材料细节丰富，构图适合建筑杂志。', '空间', '["室内","建筑","摄影"]'::jsonb, 'builtin', TRUE),
    ('社媒封面', '为{主题}设计一张社交媒体封面，视觉中心明确，高对比度，现代排版，适合移动端浏览。', '内容运营', '["社媒","封面","运营"]'::jsonb, 'builtin', TRUE)
ON CONFLICT DO NOTHING;
