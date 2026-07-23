# SmartAPI 发布记录

## v0.1.163-smartapi.1

- 状态：GitHub Release 已发布，4146 trial 基础验证通过；尚未部署生产
- 官方基线：0.1.163
- 同步分支：`sync/upstream-2026-07-23`
- 发布分支：`release/v0.1.163-smartapi.1`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.163-smartapi.1
- 目标镜像：`ghcr.io/hua7448/sub2api:0.1.163-smartapi.1`
- 发布提交：`480927918aa374df3a2a114f87ba3051bcc342c7`
- 同步提交：`2174beafe`（merge upstream `v0.1.163`）
- 上游 tag：`d0bdd7e77`（`v0.1.163`）

### 本次变更

- 同步官方 `v0.1.163` 基线；未合并 `v0.1.163` tag 之后 `upstream/main` 的同基线散提交。
- 本地 3 个冲突：
  - `backend/cmd/server/VERSION`：同步阶段保留当前 fork 版本，发布准备阶段改为 `0.1.163-smartapi.1`。
  - `deploy/docker-compose.yml`：保留 `ghcr.io/hua7448/sub2api` 明确版本镜像，拒绝 upstream 的 `weishaw/sub2api:latest`。
  - `frontend/src/components/layout/AppHeader.vue`：保留 SmartAPI 的 `app-header` / `toolbar-button` 样式，同时合入官方移动端间距与 Docs 按钮隐藏优化。
- 官方主要改动：
  - 分组级 OpenAI 推理策略：支持设置推理力度上限和精确映射，HTTP 与 WebSocket 转发统一强制执行。
  - Grok 兼容 `/responses/compact`，支持 compact 请求调度 Grok 账号，并支持链式中继受保护视频下载。
  - Redis 连接支持 ACL username。
  - 优雅关停超时不再跳过清理流程，避免缓冲用量/计费记录丢失。
  - hosted `image_generation` 工具图片 token 合并计入 `/responses` 计费；故障转移后同步缓存计费口径对齐。
  - Grok OAuth 模型同步、策略类 403 模型级隔离、Codex client tools 往返保留、OpenCode/CodeBuddy 缓存会话保留。
  - 多处移动端布局、iOS 输入框缩放、优惠码过期时间、系统日志清理错误提示、用量筛选竞态等修复。
- 保留 SmartAPI 固定规则：
  - `backend/internal/service/update_service.go` 默认更新源仍为 `hua7448/sub2api`，稳定 tag 后缀仍为 `smartapi`。
  - `.goreleaser*.yaml` 仍为 `prerelease: false`。
  - 生产部署仍要求明确版本 tag，不使用 `latest`。
  - 4146 试运行与 4145 切换脚本默认镜像更新为 `ghcr.io/hua7448/sub2api:0.1.163-smartapi.1`。

### 验证记录

- `gofmt` 已对本次变更 Go 文件执行。
- `git diff --check` 与 `git diff --cached --check` 通过。
- `pnpm --dir frontend run build` 通过；仅有 Vite dynamic import/chunk size、Browserslist 数据较旧、pnpm ignored build scripts 等非阻断警告。
- `cd backend && go test ./...` 首次因 `proxy.golang.org` 下载 `golang.org/x/*` 依赖超时失败；改用 `GOPROXY=https://goproxy.cn,direct go test ./...` 后通过。
- 4146 trial 基础验证通过：
  - `/health` 返回 `{"status":"ok"}`。
  - `/app/sub2api --version` 输出 `0.1.163-smartapi.1`（build log commit 显示 `docker`，符合服务器本地 trial 镜像构建）。
  - trial 日志 error scan 为 `(none)`。
  - `schema_migrations` 已包含 `185_group_reasoning_effort_policy.sql`。
  - trial DB `groups` 表已包含 `max_reasoning_effort`、`reasoning_effort_mappings` 两列。
- GitHub Actions Release workflow run `29974813137` 成功：
  - `headBranch=v0.1.163-smartapi.1`
  - `headSha=480927918aa374df3a2a114f87ba3051bcc342c7`
  - `update-version`、`build-frontend`、`release`、`sync-version-file` 全部 success。
- Release 校验：
  - `isDraft=false`
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.163-smartapi.1`
  - assets 包含 `checksums.txt`、linux amd64/arm64、darwin amd64/arm64、windows amd64。
  - 本地环境没有 `docker`/`crane`/`skopeo`/`oras`，且 GitHub token 缺少 `read:packages`，未在本机独立执行 GHCR manifest inspect；生产切换前需在服务器执行 `docker pull ghcr.io/hua7448/sub2api:0.1.163-smartapi.1` 补齐镜像可拉取验证。

### 部署状态

- 本次变更包含数据库迁移：
  - `backend/migrations/185_group_reasoning_effort_policy.sql`
- 迁移内容为 `groups` 表新增字段：
  - `max_reasoning_effort VARCHAR(20) NOT NULL DEFAULT ''`
  - `reasoning_effort_mappings JSONB NOT NULL DEFAULT '[]'::jsonb`
- 生产替换前必须先备份数据库；本次已完成 4146 基础健康与迁移验证，建议生产前继续人工验证登录、API Key、分组推理策略、账号调度、计费、OpenAI/Codex/Grok/Anthropic 转发、异步生图任务、用量统计和后台设置。
- 尚未执行生产 4145 切换。

### 回滚提示

- 本次包含数据库迁移。虽然是加字段且带默认值，旧应用通常会忽略新字段，但生产回滚前仍必须评估数据库向后兼容性。
- 如已发布并部署，回滚优先使用后台一键更新退回上一稳定 SmartAPI tag；容器级回滚使用 `deploy/switch-4145.sh` 指定上一明确镜像 tag。
- 若迁移后出现数据级问题，使用生产发布前备份恢复数据库；不要直接删除 volume 或数据目录。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.162-smartapi.1

- 状态：已被 `v0.1.163-smartapi.1` 发布准备替代（未推送 tag、未执行 4146 试运行、未部署生产）
- 官方基线：0.1.162
- 同步分支：`sync/upstream-2026-07-22`
- 发布分支：`release/v0.1.162-smartapi.1`
- 目标镜像：`ghcr.io/hua7448/sub2api:0.1.162-smartapi.1`
- 同步提交：`07eba13cd`（merge upstream `v0.1.162`）
- 上游 tag：`27f094e09`（`v0.1.162`）

### 本次变更

- 同步官方 `v0.1.162` 基线；未合并 `v0.1.162` tag 之后 `upstream/main` 的同基线散提交。
- 本地仅 1 个冲突：`frontend/src/components/layout/AppHeader.vue`。处理结果：保留 SmartAPI 的 `toolbar-button` 统一样式，同时采用官方新增的 `common.userMenu` 本地化无障碍文案。
- 官方主要改动：
  - 客户端真实 IP 解析改为可配置，支持显式可信代理与自定义客户端 IP 请求头，并记录相关审计日志。
  - 异步生图对象存储改为后台配置，保存即生效，并修复 `image_storage` 等凭证环境变量不可达问题。
  - 更新检查支持 `UPDATE_GITHUB_TOKEN`，规避 GitHub 未认证 API 限流。
  - API Key 部分更新时不再静默清空 IP 白/黑名单。
  - 备份设置拒绝用自动生成临时密钥持久化 S3 `SecretAccessKey`，避免重启后无法解密。
  - OpenAI/Codex/Grok/Anthropic 兼容链路包含多项 failover、缓存、计费、模型发现和内容类型修复。
  - 前端暗色模式、批量生图指引多语言、可用渠道滚动、订阅到期显示等体验修复。
- 保留 SmartAPI 固定规则：
  - `backend/internal/service/update_service.go` 默认更新源仍为 `hua7448/sub2api`，稳定 tag 后缀仍为 `smartapi`。
  - `.goreleaser*.yaml` 仍为 `prerelease: false`。
  - 生产部署仍要求明确版本 tag，不使用 `latest`。
  - 4146 试运行与 4145 切换脚本默认镜像更新为 `ghcr.io/hua7448/sub2api:0.1.162-smartapi.1`。

### 验证记录

- `git diff --check` 通过。
- `cd backend && go test ./...` 通过。
- `pnpm --dir frontend run build` 通过；仅有 Vite dynamic import/chunk size、Browserslist 数据较旧、pnpm ignored build scripts 等非阻断警告。

### 部署状态

- 本次同步未新增 `backend/migrations`，也未修改 ent schema/migrate 文件；不属于数据库迁移发布。
- 仍需在正式发布前执行 GitHub Actions Release workflow，并确认：
  - GitHub Release `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.162-smartapi.1`
  - assets 包含 `checksums.txt` 与各平台二进制包
  - GHCR 镜像 `ghcr.io/hua7448/sub2api:0.1.162-smartapi.1` 可拉取
- 生产替换前必须先用独立 trial 容器和非生产端口 `4146` 验证登录、API Key、账号调度、计费、备份对象存储配置、客户端 IP 设置、OpenAI/Codex/Grok/Anthropic 转发、异步生图任务、更新检查和用量统计。

### 回滚提示

- 本次无数据库迁移，常规应用回滚风险低于含迁移版本；但新版本可能写入新的后台设置项（客户端 IP 解析、对象存储配置、GitHub token 等）。
- 如已发布并部署，回滚优先使用后台一键更新退回上一稳定 SmartAPI tag；容器级回滚使用 `deploy/switch-4145.sh` 指定上一明确镜像 tag。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.161-smartapi.1

- 发布时间：2026-07-20
- 官方基线：0.1.161（含 0.1.160）
- 同步分支：`sync/upstream-2026-07-20`（ff 合入 main）
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.161-smartapi.1
- 镜像：`ghcr.io/hua7448/sub2api:0.1.161-smartapi.1`（多架构 linux/amd64 + linux/arm64）
- 发布提交：`333d06f3fa6167107a34f009c28a03e66b00e7fa`
- 上游 HEAD：`19149ca19`（upstream v0.1.161）

### 本次变更

- 同步官方 `v0.1.160` + `v0.1.161` 两个基线（约 70 个 upstream 提交）。
- 6 个冲突解决，保留 SmartAPI 固定规则：
  - `backend/cmd/server/VERSION`：sync commit 用基线 `0.1.161`，发布前再设 `0.1.161-smartapi.1`。
  - `backend/internal/service/wire.go`：同时保留 fork 的 `wire.Bind(modelPricingBoardGroupProvider, *APIKeyService)` 与 upstream 新增的 `ProvideAuthCacheInvalidationWorker`。
  - `backend/cmd/server/wire_gen.go`：`ProvideAdminHandlers` 调用按签名顺序同时保留 fork 的 `imageGalleryHandler` 与 upstream 新增的 `promptAdminHandler`。
  - `backend/internal/web/embed_on.go` / `embed_test.go`：fork 的 `injectPlaygroundSettings` / `injectSiteIcon` 与 upstream 新增的 `injectSiteFavicon` / `safeImageURL` 并存（补回被冲突吃掉的函数闭合括号）。
  - `deploy/docker-compose.yml`：保留 `ghcr.io/hua7448/sub2api` 镜像，拒绝 upstream 的 `sub2api:latest`。
  - `.goreleaser*.yaml` / `.github/workflows/release.yml`：`prerelease: false`、SmartAPI 文案、`-smartapi.N` tag 示例保留。
- 上游主要改动：安全（敏感操作 step-up 2FA 默认关闭、OpenAI 兼容 prompt 审计 + 审计事件持久化完整 prompt、入口拒绝聚合）、API key 鉴权（durable auth-cache invalidation outbox）、Grok（OAuth media 付费资格、签名/受保护视频代理、media model mapping、free probe 加密恢复）、路由/池（模型级临时冷却隔离、池模式遵循临时不可调度规则）、计费（上游 Sub2API 计费倍率探测刷新 + 账号展示）、网关（Claude Code 1m 模型后缀归一化、Responses 流 content_part 事件、瞬态账号耗尽归类 503）、部署（Dockerfile 交叉编译去除 QEMU、redis compose 命令参数）。

### 验证记录

- 本地验证：
  - `gofmt` 通过；`cd backend && go build ./...` 通过。
  - `cd backend && go test ./...` 全绿（含 `service` 101s、`securityaudit`、`repository`、`migrations`，无 FAIL/skip）。
  - `pnpm --dir frontend run build` 通过。
- GitHub Actions：Release workflow run `29701978580` 成功（build-frontend / update-version / release / sync-version-file 全 success），以 tag ref 触发（`headBranch=v0.1.161-smartapi.1`）。
- Release 校验：`isDraft=false`、`isPrerelease=false`、`/releases/latest` 指向 `v0.1.161-smartapi.1`、assets 含 `checksums.txt` + linux amd64/arm64 + darwin/windows；ghcr 镜像多架构。
- 4146 试运行（脚本 `deploy/trial-4146.sh`）：trial 容器替换为 `0.1.161-smartapi.1`，`/health` ok，`--version` = `0.1.161-smartapi.1`（commit `333d06f3f`），4 个迁移表均已建（count=0 正常），日志无 error/panic。
- 生产部署：后台一键热更新到 `0.1.161-smartapi.1`。热更新路径成立依据——迁移 SQL 与前端 dist 均 `go:embed` 进二进制，且本次 `backend/resources/`、`openai/instructions*.txt` 资源文件相对 0.1.159 零改动。生产 `--version` = `0.1.161-smartapi.1`，`/health` ok，4 个迁移表在生产库生效，日志无报错。

### 部署状态

- 本次变更包含数据库迁移（均为加法型 `CREATE TABLE IF NOT EXISTS` / `ADD COLUMN IF NOT EXISTS`，可逆、向后兼容）：
  - `181_prompt_audit.sql`
  - `182_prompt_audit_full_prompt.sql`
  - `183_ops_ingress_reject_aggregates.sql`
  - `184_auth_cache_invalidation_outbox.sql`
- 注意：upstream 复用了编号 181——目录里同时存在 `181_group_duplicate_operation_id.sql`（旧）和 `181_prompt_audit.sql`（新）。迁移 runner 按**文件名**（`schema_migrations.filename` 主键 + `sort.Strings`）追踪，两条是不同迁移、各自应用，安全。详见 `docs/FORK_WORKFLOW_CN.md`「经验补充」。
- 生产已切换：镜像仍为 `ghcr.io/hua7448/sub2api:0.1.136-smartapi.1`（热更新不动镜像 tag），二进制为 `0.1.161-smartapi.1`。
- 生产发布使用明确版本 tag，不使用 `latest`。
- 备份：`/root/sub2api-deploy/sub2api-2026-07-19-154640.dump`（41.8MB，迁移前快照）。

### 回滚提示

- 切换前生产状态：镜像 `ghcr.io/hua7448/sub2api:0.1.136-smartapi.1`，二进制 `0.1.159-smartapi.2`（commit `03249e45f`），记录于服务器 `/root/sub2api-deploy/PREV_VERSION.txt`。
- 热更新回滚：后台一键更新退回 `0.1.159-smartapi.2`。迁移为加法型，退回旧二进制安全（旧版本忽略新表）。
- 数据级回滚（仅必要时）：用 `sub2api-2026-07-19-154640.dump` 通过 `pg_restore` 恢复生产库。
- 如需把镜像基线刷新到与二进制一致（消除 0.1.136 → 0.1.161 的镜像漂移），用 `deploy/switch-4145.sh` 做 docker 换镜像；本次未做，留作后续可选。

## v0.1.159-smartapi.2

- 发布时间：2026-07-17
- 官方基线：0.1.159
- 同步分支：`sync/upstream-2026-07-17`
- 发布分支：`release/v0.1.159-smartapi.2`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.159-smartapi.2
- 镜像：`ghcr.io/hua7448/sub2api:0.1.159-smartapi.2`
- 发布提交：`03249e45f07bd2261bb3ff2770f6e9bd6d80c2e2`

### 本次变更

- 基于 `v0.1.159-smartapi.1` 追加 CI 修复并重新发布；`v0.1.159-smartapi.1` 已被本版本替代，不建议部署 `.1`。
- 保留 `.1` 的全部变更：同步官方 `v0.1.159`、保留 SmartAPI 定制、加入 `kimi-k3` 定价补缺。
- 本次同步不是用官方代码覆盖本地版本；已保留本地 SmartAPI 定制功能、发布规则和既有差异，并在冲突处理中合入官方 `v0.1.159` 新功能。
- 修复后端 unit-tag 测试：Image Gallery 的 `mustJSON` helper 与 upstream unit 测试 helper 冲突，改为专用命名；GLM-5.2 回归测试更新为当前专属价格预期。
- 修复后端 golangci-lint：Image Gallery、repository、update service 中的 `errcheck` 问题。
- 修复前端 playground lint：文件名清理正则、类型断言前导分号、`prefer-const` 等问题。

### 验证记录

- 本地验证：
  - `git diff --check` 通过。
  - `cd backend && go test -tags=unit ./...` 通过。
  - `cd backend && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.9.0 run ./...` 通过。
  - `make test-frontend` 通过；仅有既有 playground unused warning。
  - `pnpm --dir frontend run build` 通过；仅有既有 Vite dynamic import/chunk size、Browserslist 数据较旧、pnpm ignored build scripts 等警告。
- GitHub Actions：
  - `main` CI run `29551889934` 成功。
  - `release/v0.1.159-smartapi.2` CI run `29551889844` 成功。
  - tag `v0.1.159-smartapi.2` CI run `29551889793` 成功。
  - Release workflow run `29551889895` 成功。
- Release 校验：
  - `isDraft=false`
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.159-smartapi.2`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`、`darwin_amd64`、`darwin_arm64`、`windows_amd64`
  - workflow 以 tag ref 触发：`headBranch=v0.1.159-smartapi.2`，`headSha=03249e45f07bd2261bb3ff2770f6e9bd6d80c2e2`
- 4146 试运行：本轮未在服务器实际替换 4146 trial 容器；正式切换生产前必须按 `docs/FORK_WORKFLOW_CN.md` 和 `docs/TRIAL_DEPLOYMENT_CN.md` 完成试运行。

### 部署状态

- 本次变更包含数据库迁移：
  - `177_add_subscription_plan_currency.sql`
  - `178_channel_image_input_price.sql`
  - `179_usage_log_image_input_tokens.sql`
  - `180_audit_logs.sql`
  - `181_group_duplicate_operation_id.sql`
- 生产发布必须使用明确版本 tag：`ghcr.io/hua7448/sub2api:0.1.159-smartapi.2`，不得使用 `latest`。
- 生产替换前必须备份数据库，并先使用独立 trial 容器和非生产端口 `4146` 验证登录、API Key、账号调度、计费、审计日志、step-up 2FA、异步图片任务、对象存储、OpenAI/Codex/Grok 转发、用量统计和后台设置。
- 只替换应用容器，不删除或重建正式 PostgreSQL、Redis、`data/` 或 Docker volume。
- 本轮只完成代码同步、CI 修复、tag、GitHub Release 和 GHCR 镜像发布；未执行生产容器替换。

### 回滚提示

- 上一生产稳定版本仍按实际生产部署记录确认；若生产尚未部署 `v0.1.159-smartapi.*`，上一稳定版本通常为 `v0.1.155-smartapi.1`。
- 因本版包含数据库迁移，生产回滚前必须评估迁移是否向后兼容。
- 新版本可能写入订阅套餐币种、渠道图片输入 token 单价、usage logs 图片输入 token/费用、审计日志和重复 operation id 防护相关数据；若回滚到旧应用，需要确认旧版本对这些字段、表和约束的兼容性。
- 若已执行生产迁移，不要直接回退应用后忽略数据库状态；必须先确认迁移是否向后兼容，必要时从发布前备份恢复。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.159-smartapi.1

- 发布时间：2026-07-17
- 官方基线：0.1.159
- 同步分支：`sync/upstream-2026-07-17`
- 发布分支：`release/v0.1.159-smartapi.1`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.159-smartapi.1
- 镜像：`ghcr.io/hua7448/sub2api:0.1.159-smartapi.1`
- 发布提交：`c6644d0e19c1a2633c06f06e84fcfd3aac4a5e55`

### 本次变更

- 同步官方 `v0.1.159` 基线更新；未合入 `v0.1.159` tag 之后 `upstream/main` 上的同基线 `VERSION` chore 提交，fork 版本文件使用 `0.1.159-smartapi.1`。
- 吸收官方安全增强：操作审计日志、会话 IP/UA 绑定、敏感操作 step-up 2FA、管理员角色提升保护、API Key ACL 信任开关下的 client IP 归一化，以及相关后台路由和设置项。
- 吸收官方异步图片任务、图片结果对象存储、图片输入 token/费用拆分、上游计费倍率探测、自定义图片输入 token 单价、OpenAI 账号调度倍率、API Key 计费倍率自省等计费和图片链路更新。
- 吸收官方 OpenAI/Codex/Grok 修复：Agent Identity、Responses/WebSocket v2、Codex image tool、alpha/search、Grok OAuth/SSO、自定义上游端点、Free cache、调度缓存、failover、passthrough 错误语义和前端配置体验。
- 保留 SmartAPI 定制：模型价格广场、生图广场/Image Gallery、外部充值入口、SmartAPI 更新源、`-smartapi.N` 稳定 tag 规则、完整 release assets 和 `prerelease=false`。
- 新增 `kimi-k3` 定价补缺：输入缓存命中 ¥2.00/1M tokens，输入缓存未命中 ¥20.00/1M tokens，输出 ¥100.00/1M tokens，上下文窗口 1,048,576 tokens，并补充 billing/pricing 回归测试。

### 验证记录

- 后端验证：
  - `cd backend && go test ./...` 通过。
- 前端验证：
  - `pnpm --dir frontend run build` 通过；仅有既有 Vite dynamic import/chunk size、pnpm ignored build scripts 等警告。
- Release workflow：GitHub Actions `Release` run `29550537819` 成功。
- Release 校验：
  - `isDraft=false`
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.159-smartapi.1`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`、`darwin_amd64`、`darwin_arm64`、`windows_amd64`
  - workflow 以 tag ref 触发：`headBranch=v0.1.159-smartapi.1`，`headSha=c6644d0e19c1a2633c06f06e84fcfd3aac4a5e55`
- 4146 试运行：本轮未在服务器实际替换 4146 trial 容器；正式切换生产前必须按 `docs/FORK_WORKFLOW_CN.md` 和 `docs/TRIAL_DEPLOYMENT_CN.md` 完成试运行。

### 部署状态

- 本次变更包含数据库迁移：
  - `177_add_subscription_plan_currency.sql`
  - `178_channel_image_input_price.sql`
  - `179_usage_log_image_input_tokens.sql`
  - `180_audit_logs.sql`
  - `181_group_duplicate_operation_id.sql`
- 生产发布必须使用明确版本 tag：`ghcr.io/hua7448/sub2api:0.1.159-smartapi.1`，不得使用 `latest`。
- 生产替换前必须备份数据库，并先使用独立 trial 容器和非生产端口 `4146` 验证登录、API Key、账号调度、计费、审计日志、step-up 2FA、异步图片任务、对象存储、OpenAI/Codex/Grok 转发、用量统计和后台设置。
- 只替换应用容器，不删除或重建正式 PostgreSQL、Redis、`data/` 或 Docker volume。
- 本轮只完成代码同步、tag、GitHub Release 和 GHCR 镜像发布；未执行生产容器替换。

### 回滚提示

- 上一稳定版本：`v0.1.155-smartapi.1`
- 因本版包含数据库迁移，生产回滚前必须评估迁移是否向后兼容。
- 新版本可能写入订阅套餐币种、渠道图片输入 token 单价、usage logs 图片输入 token/费用、审计日志和重复 operation id 防护相关数据；若回滚到旧应用，需要确认旧版本对这些字段、表和约束的兼容性。
- 若已执行生产迁移，不要直接回退应用后忽略数据库状态；必须先确认迁移是否向后兼容，必要时从发布前备份恢复。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.155-smartapi.1

- 发布时间：2026-07-14
- 官方基线：0.1.155
- 同步分支：`sync/upstream-2026-07-14`
- 发布分支：未单独创建 release 分支；同步分支验证后 fast-forward 合入 `main` 并打 tag
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.155-smartapi.1
- 镜像：`ghcr.io/hua7448/sub2api:0.1.155-smartapi.1`
- 发布提交：`cbf39e044196d1d582d379977235f48e2b451ce6`

### 本次变更

- 同步官方 `v0.1.155` 基线更新；未合入 `v0.1.155` tag 之后 `upstream/main` 上的同基线零散 chore 提交。
- 吸收官方 Grok 渠道健康监控、Grok Web SSO 批量导入、Grok 免费账号滚动 24 小时配额估算、OAuth 新导入账号自动探活、Free 计划徽标、系统日志 host 过滤和可选 server timing 指标采集。
- 吸收官方网关与 OpenAI/Codex 修复：HTTP/2 keep-alive PING、Responses Lite 图像工具保留、Responses namespace 原样透传、WSv2 namespace 转发、非流式图片请求 keepalive、流式图片最终状态补齐、请求体重复扫描优化、Codex manifest 刷新稳定性与 API Key 上游代理模型列表。
- 吸收官方调度与计费修复：调度器全量重建并发合并、账号到期暂停和代理到期改投不再触发全量重建、事件延迟计算修复、OpenAI 长上下文计费改为账号级开关且默认关闭、reset credits 检测改进。
- 吸收官方 Apple container 部署支持和相关 shell 测试；Go 基线保持 `go1.26.5`，release workflow 与 `backend/go.mod` 一致。
- 冲突解决保留 SmartAPI 定制：`backend/cmd/server/VERSION` 使用 `0.1.155-smartapi.1`，嵌入前端继续支持 `data/public` 本地覆盖和 image playground 路由，同时接入官方静态资源长缓存 header。

### 验证记录

- 本地检查：
  - `git diff --check` 通过。
  - `go version` 为 `go1.26.5 darwin/arm64`。
- 后端验证：
  - `cd backend && go test ./...` 通过。
- 前端验证：
  - `pnpm --dir frontend run build` 通过，包含主前端和 `frontend/playground` 构建；仅有既有 Vite chunk size / dynamic import、Browserslist 数据较旧、playground 依赖 build script approval、pnpm 版本提示等警告。
- Release workflow：GitHub Actions `Release` run `29341993375` 成功。
- Release 校验：
  - `isDraft=false`
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.155-smartapi.1`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`、`darwin_amd64`、`darwin_arm64`、`windows_amd64`
  - workflow 以 tag ref 触发：`headBranch=v0.1.155-smartapi.1`，`headSha=cbf39e044196d1d582d379977235f48e2b451ce6`
- 4146 试运行：本轮未在服务器实际替换 4146 trial 容器；正式切换生产前必须按 `docs/FORK_WORKFLOW_CN.md` 和 `docs/TRIAL_DEPLOYMENT_CN.md` 完成试运行。

### 部署状态

- 本次变更包含数据库迁移：
  - `174_add_usage_log_long_context_billing.sql`
  - `174_add_usage_logs_api_key_latest_ip_index_notx.sql`
  - `175_add_ops_system_logs_host.sql`
  - `175_default_openai_long_context_billing.sql`
  - `175a_add_ops_system_logs_host_index_notx.sql`
  - `176_channel_monitor_grok_provider.sql`
- 生产发布必须使用明确版本 tag：`ghcr.io/hua7448/sub2api:0.1.155-smartapi.1`，不得使用 `latest`。
- 生产替换前必须备份数据库，并先使用独立 trial 容器和非生产端口 `4146` 验证登录、API Key、账号调度、计费、Grok OAuth/SSO、Grok 监控、OpenAI/Codex 转发、Responses namespace、图片生成、用量统计、系统日志和后台设置。
- 只替换应用容器，不删除或重建正式 PostgreSQL、Redis、`data/` 或 Docker volume。
- 本轮只完成代码同步、tag、GitHub Release 和 GHCR 镜像发布；未执行生产容器替换。

### 回滚提示

- 上一稳定版本：`v0.1.152-smartapi.1`
- 因本版包含数据库迁移，生产回滚前必须评估迁移是否向后兼容。
- 新版本可能写入长上下文计费字段、ops system logs host 字段、Grok channel monitor provider 配置及相关索引；若回滚到旧应用，需要确认旧版本对这些字段、索引和统计逻辑的兼容性。
- 若已执行生产迁移，不要直接回退应用后忽略数据库状态；必须先确认迁移是否向后兼容，必要时从发布前备份恢复。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.151-smartapi.1

- 发布时间：2026-07-11
- 官方基线：0.1.151
- 同步分支：`sync/upstream-2026-07-11-v0.1.151`
- 发布分支：`release/v0.1.151-smartapi.1`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.151-smartapi.1
- 镜像：`ghcr.io/hua7448/sub2api:0.1.151-smartapi.1`

### 本次变更

- 同步官方 `v0.1.151` 基线更新，包含 `v0.1.150` 与 `v0.1.151` 的 upstream 修复。
- 吸收官方 Codex 上游身份配对修复、GPT-5.6 计费/用量完整性修复、用户级 Fast/Flex 策略、compact SSE/stream bridge 加固、Grok reasoning effort 兼容、setup-token 自动刷新、parallel tool calls 兼容映射、后台用户用量拆分和英文语言包补齐。
- 吸收官方支付和订阅相关修复：账单并发与支付恢复加固、订阅额度重置、兑换调整、用户余额与支付履约回归测试补充。
- 新增 usage log `request_type=4` 的数据库约束兼容，用于记录 cyber-policy blocked 请求，避免和 legacy `request_type=0` 混淆。
- 保留 SmartAPI 定制：fork 更新源 `hua7448/sub2api`、`-smartapi.N` 稳定 tag 过滤、`prerelease=false`、完整 release assets、模型价格 fallback 与现有 SmartAPI 功能边界。

### 验证记录

- 前端验证：
  - `pnpm --dir frontend run build` 通过，包含主前端和 `frontend/playground` 构建；仅有既有 Vite chunk size / dynamic import、Browserslist 数据较旧、playground 依赖 build script approval 等警告。
- 后端验证：
  - 首次 `cd backend && go test ./...` 因 Go 自动下载 `go1.26.5` 工具链访问 `proxy.golang.org` 超时未进入测试执行。
  - 改用文档代理环境后，`ALL_PROXY=http://127.0.0.1:7897 HTTPS_PROXY=http://127.0.0.1:7897 HTTP_PROXY=http://127.0.0.1:7897 GOPROXY=https://proxy.golang.org,direct go test ./...` 通过。
- Release workflow：GitHub Actions `Release` run `29133433458` 成功。
- Release 校验：
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.151-smartapi.1`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`、`darwin_amd64`、`darwin_arm64`、`windows_amd64`
  - workflow 以 tag ref 触发：`headBranch=v0.1.151-smartapi.1`，`headSha=c27aa5fa388468e1cf2262b1675efaccddf49ae8`
- GHCR 镜像 tag 按 GoReleaser 配置为 `ghcr.io/hua7448/sub2api:0.1.151-smartapi.1`；本地环境无 `docker` CLI，且当前 `gh` token 缺少 `read:packages`，未额外读取 registry manifest。
- 4146 试运行：本轮未在服务器实际替换 4146 trial 容器；正式切换生产前必须按 `docs/FORK_WORKFLOW_CN.md` 和 `docs/TRIAL_DEPLOYMENT_CN.md` 完成试运行。

### 部署状态

- 本次变更包含数据库迁移：
  - `173_allow_cyber_blocked_usage_request_type.sql`
- 生产发布必须使用明确版本 tag：`ghcr.io/hua7448/sub2api:0.1.151-smartapi.1`，不得使用 `latest`。
- 生产替换前必须备份数据库，并先使用独立 trial 容器和非生产端口 `4146` 验证登录、API Key、账号调度、计费、支付恢复、订阅额度、OpenAI/Codex 转发、compact SSE/stream bridge、Grok reasoning effort、用量统计和后台设置。
- 只替换应用容器，不删除或重建正式 PostgreSQL、Redis、`data/` 或 Docker volume。

### 回滚提示

- 上一稳定版本：`v0.1.149-smartapi.1`
- 因本版包含数据库迁移，生产回滚前必须评估迁移是否向后兼容。
- 新版本可能写入 `usage_logs.request_type=4`；若回滚到旧应用，需要确认旧版本对该 request type 的统计、展示和约束处理是否兼容。
- 若已执行生产迁移，不要直接回退应用后忽略数据库状态；必须先确认迁移是否向后兼容，必要时从发布前备份恢复。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.149-smartapi.1

- 发布时间：2026-07-10
- 官方基线：0.1.149
- 同步分支：`sync/upstream-2026-07-10-v0.1.149`
- 发布分支：`release/v0.1.149-smartapi.1`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.149-smartapi.1
- 镜像：`ghcr.io/hua7448/sub2api:0.1.149-smartapi.1`

### 本次变更

- 同步官方 `v0.1.149` 基线更新。
- 吸收官方版本回滚 UI、用户角色/token 排行、用量布局与延迟展示、compact SSE bridge 修复、上游错误透传、Grok 修复和批量生图相关能力。
- 同步官方 Go 版本要求到 `go1.26.5`，release workflow 也校验 `go1.26.5`。
- 保留 SmartAPI 定制：fork 更新源 `hua7448/sub2api`、`-smartapi.N` 稳定 tag 过滤、`prerelease=false`、模型价格广场、生图广场/Image Gallery、外部充值入口和相关 feature flags。
- 定价 fallback 继续合并磁盘 fallback 与二进制内嵌 fallback，避免后台一键热更新只替换二进制时丢失内嵌价格补缺。

### 验证记录

- 冲突解决后检查：
  - `git diff --cached --check` 通过。
  - 未发现残留冲突标记。
- 后端验证：
  - 首次 `go version` 使用默认 `proxy.golang.org` 拉取 `go1.26.5` 超时；改用 `GOPROXY=https://goproxy.cn,direct` 后工具链下载成功。
  - `GOPROXY=https://goproxy.cn,direct go test ./internal/repository ./internal/service ./internal/handler ./internal/handler/dto ./internal/server ./internal/server/routes ./internal/setup ./cmd/server ./cmd/jwtgen` 通过。
  - `GOPROXY=https://goproxy.cn,direct go test ./...` 通过。
- 前端验证：
  - `pnpm --dir frontend run build` 通过，包含主前端和 `frontend/playground` 构建；仅有既有 Vite chunk size / dynamic import、Browserslist 数据较旧、playground 依赖 build script approval、pnpm 版本提示等警告。
- Release workflow：GitHub Actions `Release` run `29062963299` 成功。
- Release 校验：
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.149-smartapi.1`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`、`darwin_amd64`、`darwin_arm64`、`windows_amd64`
  - workflow 以 tag ref 触发：`headBranch=v0.1.149-smartapi.1`，`headSha=42f49665fd9a24ae590e5a7abefdd56f56e5a572`
- 4146 试运行：本轮未在服务器实际替换 4146 trial 容器；正式切换生产前必须按 `docs/FORK_WORKFLOW_CN.md` 和 `docs/TRIAL_DEPLOYMENT_CN.md` 完成试运行。

### 部署状态

- 本次变更包含数据库迁移：
  - `159_batch_image_foundation.sql`
  - `160_add_user_frozen_balance.sql`
  - `160_batch_image_provider_refs.sql`
  - `161_batch_image_pricing_snapshot.sql`
  - `162_add_group_batch_image_generation_gate.sql`
  - `163_batch_image_default_discount_and_hold_ratio.sql`
  - `164_batch_image_download_and_user_delete.sql`
  - `165_hide_pre_upstream_batch_image_failures.sql`
  - `166_batch_image_task_name.sql`
  - `167_clear_auto_batch_image_task_names.sql`
  - `168_restore_empty_batch_image_task_names.sql`
  - `169_batch_image_parent_batch.sql`
  - `170_add_grok_video_pricing_controls.sql`
  - `171_allow_video_usage_without_image_size.sql`
  - `172_video_per_second_billing_metadata.sql`
- 生产发布必须使用明确版本 tag：`ghcr.io/hua7448/sub2api:0.1.149-smartapi.1`，不得使用 `latest`。
- 生产替换前必须备份数据库，并先使用独立 trial 容器和非生产端口 `4146` 验证登录、API Key、账号调度、计费、批量生图、Grok 视频计费、用量统计和后台回滚 UI。
- 只替换应用容器，不删除或重建正式 PostgreSQL、Redis、`data/` 或 Docker volume。

### 回滚提示

- 上一稳定版本：`v0.1.146-smartapi.1`
- 因本版包含数据库迁移，回滚前必须评估旧版本对新增批量生图表、用户冻结余额、分组批量生图开关、Grok 视频价格控制和视频计费元数据字段的兼容性。
- 若已执行生产迁移，不要直接回退应用后忽略数据库状态；必须先确认迁移是否向后兼容，必要时从发布前备份恢复。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.146-smartapi.1

- 发布时间：2026-07-09
- 官方基线：0.1.146
- 同步分支：`sync/upstream-2026-07-09-v0.1.146`
- 发布分支：`release/v0.1.146-smartapi.1`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.146-smartapi.1
- 镜像：`ghcr.io/hua7448/sub2api:0.1.146-smartapi.1`

### 本次变更

- 同步官方 `v0.1.146` 基线更新。
- 吸收官方 API Key 并发统计、账号请求头覆写、账号导入拖拽批量处理、OpenAI 新模型/价格、Grok 图片价格控制、订阅人民币预览、Redis scan/index cleanup hardening、endpoint/responses compact 修复、Codex/OAuth header 与 User-Agent 修复、非 v1 模型同步修复等更新。
- 合入本 fork 的生图广场发布改动：按 `gpt_image_playground v0.6.12` 对齐普通生图、编辑、streaming partial image、结果处理、重试、ZIP 导出和 Agent 生图逻辑。
- 保留 sub2api 安全边界：前端不保存真实 API Key，不开放 Base URL/自定义 provider；所有生图请求经 `/api/v1/image-playground/proxy/...` 后端鉴权代理，并校验登录态、KEY 归属、状态、额度、过期时间和分组 `allow_image_generation`。
- 增强生图后端代理 guard：Images/Responses 请求按管理端 allowed models、allowed agent models、尺寸、质量、输出格式、Agent/Web Search 开关做服务端白名单校验。
- 更新 `docs/IMAGE_PLAYGROUND_RUNBOOK_CN.md`，明确默认生图链路沿用上游 `gpt_image_playground v0.6.12`，sub2api 只做同源代理、安全适配和计费权限边界。

### 验证记录

- 纯 upstream 同步分支验证：
  - `git diff --check` 通过。
  - `pnpm --dir frontend run build` 通过。
  - `cd backend && go test ./internal/repository ./internal/service ./internal/handler ./internal/server/routes` 通过。
- release 分支验证：
  - `git diff --check` 通过。
  - `pnpm --dir frontend/playground exec vitest run src/lib/api.test.ts` 通过。
  - `pnpm --dir frontend/playground exec vitest run src/lib/sub2apiDirect.test.ts` 通过。
  - `pnpm --dir frontend/playground exec vitest run src/lib/downloadImages.test.ts src/lib/exportFileName.test.ts src/lib/dataUrl.test.ts src/lib/exportZip.test.ts` 通过。
  - `pnpm --dir frontend/playground run build` 通过。
  - `pnpm --dir frontend run build` 通过。
  - `cd backend && go test ./internal/repository ./internal/service ./internal/handler ./internal/server/routes` 通过。
- Release workflow：GitHub Actions `Release` run `28991661313` 成功。
- Release 校验：
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.146-smartapi.1`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`、`darwin_amd64`、`darwin_arm64`、`windows_amd64`
  - GHCR 镜像 tag 按 GoReleaser 配置为 `ghcr.io/hua7448/sub2api:0.1.146-smartapi.1`
- 4146 试运行：本轮已输出固定隔离 trial 环境命令；未在服务器实际替换 4146 容器，正式切换生产前仍需按 `docs/TRIAL_DEPLOYMENT_CN.md` 完成登录、生图、Agent、Network、usage logs 和额度扣减检查。

### 部署状态

- 本次变更不包含数据库迁移。
- `main` 已推送到 `4ccb107c`，tag `v0.1.146-smartapi.1` 指向同一提交。
- 生产发布必须使用明确版本 tag：`ghcr.io/hua7448/sub2api:0.1.146-smartapi.1`，不得使用 `latest`。
- 生产替换前必须先使用 4146 trial 环境验证；只替换应用容器，不删除或重建正式 PostgreSQL、Redis、`data/` 或 Docker volume。

### 回滚提示

- 上一稳定版本：`v0.1.145-smartapi.2`
- 因本版无数据库迁移，回滚以替换应用镜像 tag 或二进制为主。
- 回滚时不得执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.145-smartapi.2

- 发布时间：2026-07-06
- 官方基线：0.1.145
- 发布分支：`release/v0.1.145-smartapi.2`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.145-smartapi.2
- 镜像：`ghcr.io/hua7448/sub2api:0.1.145-smartapi.2`

### 本次变更

- 修复后台一键更新/二进制热更新后 Kimi K2.7 模型广场官方价仍显示 `¥4/M` 的问题。
- 将 `backend/resources/model-pricing/model_prices_and_context_window.json` 通过 `go:embed` 编入 release 二进制；即使热更新只替换 `/app/sub2api`、没有同步 `/app/resources` 目录，也能拿到 SmartAPI bundled fallback 价格。
- 价格服务现在会依次合并磁盘 fallback 与二进制内嵌 fallback 中缺失的模型条目；远程/持久化价格源仍优先，内嵌 fallback 只补缺失项。
- 部署示例中的应用镜像改为明确 SmartAPI tag：`ghcr.io/hua7448/sub2api:0.1.145-smartapi.2`。

### 验证记录

- 后端相关测试：`go test ./resources ./internal/service -run 'TestMergeMissingFallbackPricing|TestDefaultPricingIncludesCodexAutoReview|TestModelPricingBoard|TestGetModelPricing'` 通过。
- 后端价格测试：`go test ./internal/service -run 'Test.*Pricing'` 通过。
- 后端完整测试：`go test ./...` 通过。
- 前端完整构建：`pnpm --dir frontend run build` 通过（仅有既有 Vite chunk size / dynamic import、Browserslist 数据较旧、playground 依赖 build script approval 警告）。
- Release workflow：GitHub Actions `Release` run `28804653488` 成功。
- Release 校验：
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.145-smartapi.2`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`、`darwin_amd64`、`darwin_arm64`、`windows_amd64`
- 4146 试运行：正式替换生产 4145 前必须完成试运行健康检查。

### 问题复盘

- `v0.1.144-smartapi.2` 已把 Kimi K2.7 官方人民币价写入 bundled fallback JSON，但生产运行时优先读取持久化价格文件；旧持久化文件没有 Kimi K2.7，因此回退到代码兜底价，输出价被展示成 `¥4/M`。
- `v0.1.145-smartapi.1` 增加了“从 bundled fallback JSON 补缺失条目”的逻辑，但该修复仍假设 `/app/resources/model-pricing/model_prices_and_context_window.json` 会随版本更新。
- 实际后台一键更新/二进制热更新只替换 `/app/sub2api`，不会同步 `/app/resources`。因此 `.1` 在正式环境中仍可能读到旧 resources，Kimi K2.7 价格继续不生效。
- `.2` 的最终修复是把 fallback JSON 用 `go:embed` 编进二进制；以后涉及一键更新必须假设“只有二进制会更新”，不能依赖 release archive 或容器内资源目录同步。
- 后续凡是影响热更新行为的静态资源、模型价格、默认模板、规则表或内置配置，如果需要随二进制版本生效，必须内嵌到 Go binary，或提供明确的数据迁移/复制机制，并在发布说明中写明验证方法。

### 部署状态

- 本次变更不包含数据库迁移。
- 生产发布必须使用明确版本 tag：`ghcr.io/hua7448/sub2api:0.1.145-smartapi.2`，不得使用 `latest`。
- 如果使用后台一键更新，更新后需确认运行版本为 `0.1.145-smartapi.2`，并刷新模型广场验证 `kimi-k2.7-code` 官方价显示为输入 `¥6.5/M`、输出 `¥27/M`、缓存读 `¥1.3/M`。

### 回滚提示

- 上一稳定版本：`v0.1.145-smartapi.1`
- 回滚时只替换应用镜像 tag 或应用二进制，不删除数据库、Redis、`data/` 或 Docker volume。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.145-smartapi.1

- 发布时间：2026-07-06
- 官方基线：0.1.145
- 发布分支：`release/v0.1.145-smartapi.1`
- 同步分支：`sync/upstream-2026-07-06`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.145-smartapi.1
- 镜像：`ghcr.io/hua7448/sub2api:0.1.145-smartapi.1`

### 本次变更

- 同步官方 `v0.1.145` 基线更新。
- 新增 API Key 账号请求头覆写能力，并包含官方后续审计修复：补齐禁止名单、beta 头对称校验与批量清空防护。
- 适配 OpenAI 新模型 `gpt-5.6-sol`、`gpt-5.6-terra`、`gpt-5.6-luna`。
- 引入 OpenAI 高级调度器相关控制与审计修复：支持粘性加权、订阅优先、调度权值配置，并补充调度缓存与账号列表相关测试。
- 支付侧同步 EasyPay 自定义支付方式、内置支付方式精确匹配、订阅 CNY 换算显式 opt-in 以及支付配置校验修复。
- 继续保留 SmartAPI 定制：fork 更新源、明确 SmartAPI release tag、模型价格广场、用户侧渠道页面开关、外部充值入口、站点文档等。
- 修复 Kimi K2.7 模型广场官方价显示：当生产持久化价格文件缺少镜像内新增模型条目时，运行时会把 bundled fallback JSON 中缺失的条目补入内存价格表，避免 `kimi-k2.7-code` 回退到旧的 `¥4/M` 展示口径。
- 部署示例中的应用镜像改为明确 SmartAPI tag：`ghcr.io/hua7448/sub2api:0.1.145-smartapi.1`，避免生产环境误用 upstream `latest`。

### 验证记录

- 冲突解决：`backend/cmd/server/VERSION`、`frontend/src/components/layout/AppSidebar.vue`、`frontend/src/views/admin/__tests__/SettingsView.spec.ts` 已处理；保留 fork 侧边栏 logo 视觉风格，并吸收 upstream logo 回首页行为。
- JSON 结构校验：`jq empty backend/resources/model-pricing/model_prices_and_context_window.json` 通过。
- 后端相关测试：`go test ./internal/service -run 'TestMergeMissingFallbackPricing|TestDefaultPricingIncludesCodexAutoReview|TestModelPricingBoard|TestGetModelPricing|Test.*Pricing'` 通过。
- 后端完整测试：`go test ./...` 通过。
- 前端完整构建：`pnpm --dir frontend run build` 通过（仅有既有 Vite chunk size / dynamic import、Browserslist 数据较旧、playground 依赖 build script approval 警告）。
- 4146 试运行：正式替换生产 4145 前必须完成试运行健康检查。

### 部署状态

- 本次变更不包含数据库迁移。
- 生产发布必须使用明确版本 tag：`ghcr.io/hua7448/sub2api:0.1.145-smartapi.1`，不得使用 `latest`。
- 生产替换前需先按 `docs/FORK_WORKFLOW_CN.md` 使用独立试运行容器和非生产端口 `4146` 验证登录、模型广场、Kimi K2.7 官方价、支付配置与 OpenAI/Anthropic 基本转发。

### 回滚提示

- 上一稳定版本：`v0.1.144-smartapi.2`
- 回滚时只替换应用镜像 tag，不删除数据库、Redis、`data/` 或 Docker volume。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.144-smartapi.2

- 发布时间：2026-07-06
- 官方基线：0.1.144
- 发布分支：`main`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.144-smartapi.2
- 镜像：`ghcr.io/hua7448/sub2api:0.1.144-smartapi.2`

### 本次变更

- 补充本地 fallback 模型定价 JSON：新增 `kimi-k2.7-code`、`kimi-k2.7-code-highspeed`、`kimi-for-coding`。
- Kimi K2.7 Code 标准版按官方人民币价配置：输入缓存未命中 `¥6.50/M`、缓存命中 `¥1.30/M`、输出 `¥27.00/M`。
- Kimi K2.7 Code HighSpeed 按标准版 2 倍配置：输入缓存未命中 `¥13.00/M`、缓存命中 `¥2.60/M`、输出 `¥54.00/M`。
- `kimi-for-coding` 作为 Kimi Code 统一模型 ID，按 Kimi K2.7 Code 标准版价格兜底；Kimi 官方后续若在公开定价源中发布该模型或 K2.7 Code 的新定价，以官方最新定价为准。
- 修正 Kimi fallback 注释，避免把 `kimi-for-coding` 误读为 K2.6 档位。
- 包含 `v0.1.144-smartapi.1` 后已合入 `main` 的站点文档维护：新增 SmartQ 文档站点，更新 Kimi 下载链接和充值说明，并补充 `site/` 维护规范。

### 验证记录

- JSON 结构校验：`jq empty backend/resources/model-pricing/model_prices_and_context_window.json` 通过。
- 后端相关测试：`go test ./internal/service -run 'Test.*Pricing|TestModelPricingBoard|TestCalculateCost|TestGetModelPricing'` 通过。
- 前端完整构建：`pnpm --dir frontend run build` 通过（仅有既有 Vite chunk size / dynamic import 警告）。
- 4146 试运行：正式替换生产 4145 前必须完成试运行健康检查。

### 部署状态

- 本次变更不包含数据库迁移。
- 生产发布前仍需按 `docs/FORK_WORKFLOW_CN.md` 使用明确版本 tag，先部署 4146 试运行容器验证，再切换生产服务。

### 回滚提示

- 上一稳定版本：`v0.1.144-smartapi.1`
- 回滚时只替换应用镜像 tag，不删除数据库、Redis、`data/` 或 Docker volume。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。

## v0.1.143-smartapi.1

- 发布时间：2026-07-02
- 官方基线：0.1.143
- 发布分支：`release/v0.1.143-smartapi.1`
- 同步分支：`sync/upstream-2026-07-02-v0.1.143`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.143-smartapi.1
- 镜像：`ghcr.io/hua7448/sub2api:0.1.143-smartapi.1`

### 本次变更

- 同步官方 `v0.1.143` 基线更新。
- 新增订阅分组高峰时段倍率：支持为订阅分组配置高峰时段、倍率和服务器时区展示，高峰倍率会透传到可用渠道、支付计划、结算与分组展示链路。
- 修正计费术语：将「文本倍率」调整为「token 倍率」，并明确高峰倍率作用于 token 计费，图片按次计费不受高峰倍率影响。
- 新增 OpenAI WebSocket `http_bridge` ingress 模式和账号级 WS 选择器。
- 支持恢复已撤销订阅，显示 OpenAI 账号额度重置到期时间，并持久化 OpenAI 订阅到期时间。
- 用量记录新增 IP 地理位置查询与展示，管理端分组列表支持自定义列显示。
- Anthropic 渠道支持 API Key Bearer 认证。
- 修复 Claude Code 流式响应 keepalive 卡顿、OpenAI OAuth `count_tokens` scope 错误、Codex 图片桥接 `/responses/compact` 注入、Gemini 推理模型无效参数、用户模型统计聚合等官方问题。
- 保留 SmartAPI 定制：fork 更新源、生图广场、Image Gallery、模型价格广场、国内模型/GLM 定价等。

### 验证记录

- 后端完整测试：首次 `go test ./...` 因下载 `github.com/tiktoken-go/tokenizer@v0.8.0` 超时失败；按规范使用本地代理重跑 `ALL_PROXY=http://127.0.0.1:7897 HTTPS_PROXY=http://127.0.0.1:7897 HTTP_PROXY=http://127.0.0.1:7897 GOPROXY=https://proxy.golang.org,direct go test ./...` 通过。
- 前端完整构建：`pnpm --dir frontend run build` 通过。
- Release workflow：GitHub Actions `Release` run `28595423046` 成功。
- Release 校验：
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.143-smartapi.1`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`、`darwin_amd64`、`darwin_arm64`、`windows_amd64`。
  - GoReleaser 步骤成功；本机无 `docker` / `skopeo`，且当前 GitHub token 缺少 `read:packages` scope，未做本地 GHCR manifest/API 校验。
- 4146 试运行：正式替换生产 4145 前必须完成试运行健康检查。

### 部署状态

- 本次变更包含数据库迁移：
  - `158_add_group_peak_rate_multiplier.sql`
  - `158_enable_grok_media_generation_groups.sql`
- 迁移会修改 `groups` 表并回填既有 Grok 分组的图片生成开关；生产发布前必须备份数据库，并先使用非生产端口 4146 试运行。
- 因包含数据库迁移，回滚前必须评估旧版本对新增分组高峰倍率字段、Grok 媒体生成分组配置的兼容性。
- Agent 未直接执行服务器 4146 替换和正式 4145 切换；需按 `docs/FORK_WORKFLOW_CN.md` 和 `docs/TRIAL_DEPLOYMENT_CN.md` 先试运行再正式切换。

### 回滚提示

- 上一稳定版本：`v0.1.142-smartapi.1`
- 回滚时只替换应用镜像 tag，不删除数据库、Redis、`data/` 或 Docker volume。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。
- 因包含数据库迁移，回滚前必须评估 `groups` 高峰倍率相关字段、Grok 媒体生成分组配置是否影响旧版本运行。

## v0.1.142-smartapi.1

- 发布时间：2026-07-02
- 官方基线：0.1.142
- 发布分支：`release/v0.1.142-smartapi.1`
- 同步分支：`sync/upstream-2026-07-02`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.142-smartapi.1
- 镜像：`ghcr.io/hua7448/sub2api:0.1.142-smartapi.1`

### 本次变更

- 同步官方 `v0.1.142` 基线更新，未合并官方 tag 之后 `upstream/main` 的同基线散提交。
- 新增 OpenAI Spark 链接型影子账号：影子账号通过 `parent_account_id` 复用母账号凭据和代理，独立走 `spark` 配额维度与用量窗口。
- 新增 Grok 媒体路由：支持 Grok 图像生成、图像编辑、视频生成和视频状态查询，并接入调度、鉴权、内容审核、错误透传和用量记录。
- 适配 Claude Sonnet 5，并新增 Anthropic OAuth dateline 指纹归一化开关，默认抹除客户端日期句式隐写位。
- 修复 OpenAI Codex OAuth 多轮 encrypted reasoning 保留问题，修复 `gpt-5.5-pro` 被错误归一化问题，并调整默认模型列表。
- 修复账号列表分页 total 不准、订阅撤销软删除失效、五平台配额更新等官方问题。
- 前端同步 Grok 图标与配色、用户使用记录“推理强度”默认列、登录协议提示 i18n 和多处 en/zh 文案修复。
- 保留 SmartAPI 定制：fork 更新源、生图广场、Image Gallery、模型价格广场、国内模型/GLM 定价等。

### 验证记录

- 后端完整测试：`go test ./...` 通过。
- 前端完整构建：`pnpm --dir frontend run build` 通过。
- Release workflow：GitHub Actions `Release` run `28559013060` 成功。
- Release 校验：
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.142-smartapi.1`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`、`darwin_amd64`、`darwin_arm64`、`windows_amd64`。
  - GHCR 镜像 `ghcr.io/hua7448/sub2api:0.1.142-smartapi.1` 已由 GoReleaser 成功推送；本机无 `docker` / `skopeo`，未做本地 manifest inspect。
- 4146 试运行命令已给出，正式切换前必须完成试运行检查。

### 部署状态

- 本次变更包含数据库迁移：
  - `154_account_spark_shadow.sql`
  - `154a_account_spark_shadow_indexes_notx.sql`
- 迁移会修改 `accounts` 表：
  - 新增 `parent_account_id`
  - 新增 `quota_dimension`
  - 增加账号父子关系、维度合法性、禁止自指等约束
  - 增加 `parent_account_id` 索引和每个母账号一个 active spark 影子账号的唯一索引
- 生产发布前必须备份数据库，并确认 4146 试运行健康检查通过。
- 因包含数据库约束和并发索引迁移，回滚前必须评估旧版本对新增列和约束的兼容性。
- Agent 未直接执行服务器 4146 替换和正式 4145 切换；需按 `docs/FORK_WORKFLOW_CN.md` 和 `docs/TRIAL_DEPLOYMENT_CN.md` 先试运行再正式切换。

### 回滚提示

- 上一稳定版本：`v0.1.141-smartapi.1`
- 回滚时只替换应用镜像 tag，不删除数据库、Redis、`data/` 或 Docker volume。
- 严禁执行 `docker compose down -v`、删除生产数据目录或删除生产 volume。
- 因包含数据库迁移，回滚前必须评估 `accounts.parent_account_id`、`accounts.quota_dimension`、相关 CHECK/FK 约束和并发索引是否影响旧版本运行。

## v0.1.141-smartapi.1

- 发布时间：2026-07-01
- 官方基线：0.1.141
- 发布分支：`release/v0.1.141-smartapi.1`
- 同步分支：`sync/upstream-2026-07-01`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.141-smartapi.1
- 镜像：`ghcr.io/hua7448/sub2api:0.1.141-smartapi.1`

### 本次变更

- 同步官方 `v0.1.139`、`v0.1.140`、`v0.1.141` 基线更新。
- 新增 Grok 订阅、Grok OAuth、Grok quota probe、Grok CLI 兼容路由和相关前后端能力。
- 新增 Codex PAT 上游认证、GPT-5.5 Codex 指令支持，并加强 `codex_cli_only` / engine fingerprint 检测。
- 新增 API Key 列设置、Ops 系统日志 API Key 过滤、OpenAI quota headroom 调度权重。
- 用户端用量分析对齐管理员视角，补充分组维度、端点/模型分布和图表展示。
- 修复多项 OpenAI 网关、图片计费、支付订阅金额、退款 pending、注册/OAuth、余额扣费和前端直连 API base 问题。
- 保留 SmartAPI 定制：fork 更新源、生图广场、Image Gallery、模型价格广场、国内模型/GLM 定价等。

### 验证记录

- 后端完整测试：`go test ./...` 通过。
- 前端完整构建：`pnpm --dir frontend run build` 通过。
- 4146 试运行命令已给出，正式切换前必须完成试运行检查。
- Release workflow：GitHub Actions `Release` run `28488231851` 成功。
- Release 校验：
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.141-smartapi.1`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`。
  - GHCR API 查询因本地 GitHub token 缺少 `read:packages` scope 返回 403；workflow 的 GoReleaser 步骤已成功。

### 部署状态

- 本次变更包含数据库迁移：
  - `154_add_ops_system_logs_api_key_id.sql`
  - `155_add_ops_system_logs_api_key_id_index_notx.sql`
  - `156_content_moderation_matched_keyword.sql`
  - `157_user_platform_quotas_add_grok.sql`
- 生产发布前必须备份数据库，并确认 4146 试运行健康检查通过。
- Agent 未直接执行服务器 4146 替换和正式 4145 切换；已给出按 `docs/FORK_WORKFLOW_CN.md` 和 `docs/TRIAL_DEPLOYMENT_CN.md` 执行的 4146 试运行命令。
- 官方 `v0.1.141` 之后的 `upstream/main` 散提交未合并；其中 `fix(openai): preserve encrypted reasoning across turns on codex OAuth path` 可后续单独评估是否作为同基线紧急例外。

### 回滚提示

- 上一稳定版本：`v0.1.138-smartapi.5`
- 回滚时只替换应用镜像 tag，不删除数据库、Redis、`data/` 或 Docker volume。
- 因包含数据库迁移，回滚前必须评估迁移向后兼容性；尤其关注 `ops_system_logs.api_key_id`、`content_moderation_logs.matched_keyword` 和 `user_platform_quotas` 平台 CHECK 约束。

## v0.1.138-smartapi.5

- 发布时间：2026-06-29
- 官方基线：0.1.138
- 发布分支：`main`
- 功能分支：`feature/domestic-model-pricing`
- Release URL：https://github.com/hua7448/sub2api/releases/tag/v0.1.138-smartapi.5
- 镜像：`ghcr.io/hua7448/sub2api:0.1.138-smartapi.5`

### 本次变更

- 新增国产模型价格兜底与模型价格广场展示支持，覆盖 DeepSeek V4、GLM、Kimi、MiniMax 等模型族。
- 模型价格广场新增国产模型分类，并按国产模型使用人民币等值展示。
- 调整国产模型在价格页的排序和厂商展示，避免统一显示为 Anthropic。
- 补充 GLM-5 系列前端白名单：`glm-5.2`、`glm-5.1`、`glm-5`、`glm-5-turbo`。
- 联网核对并记录 GLM-5-Turbo 官方 Z.AI 定价口径：输入 `$1.20/M`、缓存命中 `$0.24/M`、输出 `$4.00/M`。

### 验证记录

- 后端完整测试：`go test ./...` 通过。
- 前端完整构建：`pnpm --dir frontend run build` 通过。
- 定价页前端单测：`pnpm --dir frontend exec vitest run src/views/user/__tests__/ModelPricingBoardView.spec.ts` 通过。
- Release workflow：GitHub Actions `Release` run `28353150685` 成功。
- Release 校验：
  - `isPrerelease=false`
  - `/releases/latest` 指向 `v0.1.138-smartapi.5`
  - assets 包含 `checksums.txt`、`linux_amd64`、`linux_arm64`。

### 部署状态

- 本次变更不包含数据库迁移。
- 生产服务发布前健康检查：`http://18.217.117.151:4145/health` 返回 `{"status":"ok"}`。
- 4146 试运行健康检查：`http://18.217.117.151:4146/health` 返回 `{"status":"ok"}`。
- Agent 未直接执行服务器 4146 替换和正式 4145 切换：本地无可用服务器 SSH 私钥，SSH 返回 `Permission denied (publickey...)`。
- 已给出按 `docs/FORK_WORKFLOW_CN.md` 和 `docs/TRIAL_DEPLOYMENT_CN.md` 执行的 4146 试运行、数据库备份和正式切换命令。

### 回滚提示

- 上一稳定版本：`v0.1.138-smartapi.4`
- 回滚时只替换应用镜像 tag，不删除数据库、Redis、`data/` 或 Docker volume。
