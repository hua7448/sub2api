# SmartAPI 发布记录

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

### 部署状态

- 本次变更包含数据库迁移：
  - `154_add_ops_system_logs_api_key_id.sql`
  - `155_add_ops_system_logs_api_key_id_index_notx.sql`
  - `156_content_moderation_matched_keyword.sql`
  - `157_user_platform_quotas_add_grok.sql`
- 生产发布前必须备份数据库，并确认 4146 试运行健康检查通过。
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
