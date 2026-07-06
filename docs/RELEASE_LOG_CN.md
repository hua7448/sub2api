# SmartAPI 发布记录

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
- 待执行：Release workflow。
- 4146 试运行：正式替换生产 4145 前必须完成试运行健康检查。

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
