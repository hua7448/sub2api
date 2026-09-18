# Sub2API 项目梳理

梳理日期：2026-09-10。范围为当前源码导出目录的结构、关键入口、功能实现与随附文档；不是逐行代码审计，也没有连接生产环境。

## 当前版本的性质

这是在 Sub2API 上持续演进的定制版本。历史资料指向上游 `Wei-Shaw/sub2api` 和 fork `bayma888/sub2api-bmai`，随后又合入了生产前后端修复及其他功能实现。

当前没有 `.git`，因此无法确认实际分支、提交祖先或计算相对上游的新增/修改行数。本文的“扩展”指当前需要关注和保留的功能，并不代表每一项都由当前维护者原创或仍未合入上游。

- 后端嵌入版本文件为 `0.1.179`，不能据此认定代码等同于官方该版本。
- `CUSTOM_CHANGELOG.md` 最近条目为 2026-07-24 的源码归档，记录了前后端来自不同工作目录。
- `docs/reviews/CPA_ACCOUNT_IDENTITY_20260907.md` 记录了更晚的身份兼容改动，明确说明当次尚未部署。
- `dev-memory/` 主要停留在 2026 年 3 月，里面的版本、分支及环境资料明显早于当前代码。
- 静态文件统计：2,511 个 Go 文件，其中 1,240 个测试文件；307 个 Vue 文件；456 个 TypeScript 文件，其中 246 个 `.spec.ts` / `.test.ts`；273 个 SQL migration 文件。Go 数量包含 Ent 生成代码，不能当作定制代码量。

## 系统结构

```text
Vue 页面 / Pinia / Vue Router
    -> frontend/src/api 的统一 HTTP 客户端
    -> /api/v1 用户、管理、支付和设置接口
    -> handler / DTO -> service -> repository -> PostgreSQL / Redis

外部 API 客户端
    -> /v1、/v1beta 及兼容别名
    -> 请求限制、API Key 验证、混合分组解析
    -> 各协议 handler + service 中的审计、额度、并发、调度和转发
    -> 上游模型服务
    -> 用量结算、日志、缓存和监控
```

上图为模块关系，具体鉴权、审计、并发和计费时序因端点而异，应以对应 handler 和回归测试为准。

后端是 Go + Gin + Ent + Wire。数据库以 PostgreSQL 保存业务事实，Redis 用于缓存、队列、锁及运行状态。前端是 Vue 3 + TypeScript + Vite + Pinia + Tailwind + vue-i18n，支付 SDK 按页面加载。

主要入口：

| 位置 | 作用 |
| --- | --- |
| `backend/cmd/server/main.go` | 安装检查、服务启动、Prompt Audit 启动与退出 |
| `backend/cmd/server/wire.go` | 依赖装配源定义；`wire_gen.go` 为生成结果 |
| `backend/internal/server/router.go` | 全局中间件、管理面与网关路由装配 |
| `backend/internal/server/routes/gateway.go` | 多协议入口、平台分派、混合分组解析及兼容别名 |
| `backend/internal/service/wire.go` | 服务构造、后台任务与实例角色判断 |
| `backend/internal/repository/ent.go` | 数据库初始化、迁移和角色边界 |
| `frontend/src/router/index.ts` | 页面入口、懒加载、权限及功能可见性 |
| `frontend/src/api/client.ts` | 鉴权、刷新、语言、时区和统一响应处理 |
| `frontend/vite.config.ts` | 开发代理、品牌设置注入、分包及 UI 构建目录 |

## 定制与扩展地图

以下模块均能在当前源码中找到实现；生产启用状态未核实。

| 模块 | 当前行为与维护重点 | 主要源码入口 |
| --- | --- | --- |
| 邀请返利 | 余额和订阅天数两类返利；支付/兑换来源；有效客户门槛、冻结、提取及幂等 | `service/affiliate_service.go`、`repository/affiliate_repo.go`、`views/user/AffiliateView.vue`、`views/admin/affiliates/` |
| 注册与身份 | 服务商别名和 `+tag` 独立策略、邮箱去重、OAuth 身份绑定、Passkey、验证码服务 | `service/registration_email_policy.go`、`service/registration_email_alias.go`、`service/auth_*`、`routes/auth.go` |
| 排行榜 | 消费、充值、Token、请求、勤劳榜；个人和公开查询入口 | `service/usage_service.go`、`repository/usage_log_repo.go`、`views/user/LeaderboardView.vue` |
| 混合分组 | 同一分组和 Key 按模型/端点解析实际平台，支持明确模型路由 | `service/composite_route_resolver.go`、`repository/composite_model_route_repo.go`、`routes/gateway.go`、`views/admin/GroupsView.vue` |
| 上游扩展 | OpenAI、Claude、Gemini/Antigravity、Grok，以及 Kimi、智谱、DeepSeek 等兼容平台；能力按平台分派 | `routes/gateway.go`、`service/`、`internal/pkg/xai/`、`api/admin/cnProviders.ts` |
| 利润和定价 | 账号选择的利润门槛、请求价格快照、重试前复核；模型、时段、图像、音视频、搜索等价格配置 | `service/openai_profit_control.go`、`service/gateway_profit_control.go`、`service/profit_preview.go`、`views/admin/groupsProfitControl.ts` |
| 订阅额度 | 当日额度重置、用户重置周限额、上海时区日边界、数据库事务和缓存写入屏障 | `service/subscription_reset_today.go`、`repository/user_subscription_repo.go`、`repository/usage_billing_repo.go` |
| Codex 账号身份 | 持久随机种子，四种身份模式，HTTP/WS 会话和 metadata 一致性，多回合隔离 | `service/openai_codex_fingerprint.go`、`service/openai_codex_identity.go`、`service/openai_agent_identity.go` |
| 图片工作流 | OpenAI/Grok 异步任务通过 Redis 和对象存储返回结果；Gemini/Vertex 批量任务有持久状态、队列、预扣与结算恢复 | `handler/image_task_handler.go`、`service/batch_image*.go`、`repository/batch_image*.go`、`views/user/AiStudioView.vue` |
| 渠道监控 | V1 主动探测与 V2 被动日志聚合；聚合回填、隐私、缓存和前端指标矩阵 | `service/channel_monitor_v2*.go`、`repository/channel_monitor_v2*.go`、`features/channel-monitor-v2/` |
| Prompt Audit | 独立审计节点池、配置快照、异步任务、同步阻断、协议快照和管理控制台 | `backend/internal/securityaudit/`、`frontend/src/features/prompt-audit/`、`openspec/changes/add-openai-compatible-prompt-audit/` |
| 支付与运营 | 支付宝、微信、易支付、Stripe、Airwallex；订单、套餐、兑换码、审计和运营面板 | `backend/internal/payment/`、`routes/payment.go`、`views/admin/orders/`、`views/admin/ops/` |
| 多节点和 UI | Primary/API 分工、缓存失效与后台任务协调；嵌入 UI revision、防旧缓存白屏、深色和移动布局 | `config/config.go`、`repository/ent.go`、`service/wire.go`、`web/embed_on.go` |

表中未带项目根的 `service/`、`repository/`、`handler/`、`routes/` 分别位于 `backend/internal/service/`、`backend/internal/repository/`、`backend/internal/handler/`、`backend/internal/server/routes/`；`views/`、`features/`、`api/` 位于 `frontend/src/`。

## 已发现的资料和构建差异

1. **Go 版本冲突。** `backend/go.mod` 为 1.26.5；CI 按该文件安装，却用字符串断言要求 1.26.6；三个 Dockerfile 使用 1.26.6。生产合约脚本则固定 Linux Go 1.26.5。后续应区分开发、CI 和生产构建约束，再做针对性统一。本次没有修改这些文件。
2. **历史开发指南不能直接套用。** `LOCAL_DEV_GUIDE.md` 面向 Windows，Go 版本仍写 1.25.7+，并声称已有 `backend/config.yaml`，但当前导出目录没有此文件。当前开发 Compose 的 PostgreSQL/Redis 镜像也已是 18/8，与旧文档的 16/7 不同。
3. **迁移说明存在旧命令。** `backend/migrations/README.md` 中的 `make migrate-up` 和 `make migrate-down` 在现有 Makefile 中不存在。实际是应用启动调用自定义 runner；API 角色跳过迁移。
4. **迁移编号有重复前缀和跳号。** 例如 181、185、186、194、225、226 均有不同文件。runner 按完整文件名记录和排序，因此不能直接认定冲突，更不能为整齐而重编号。最高前缀 228 也不等于实际迁移文件数。
5. **源码导出不能完成生产来源校验。** `verify-production-contracts.sh` 检查规范分支、提交祖先、特定迁移哈希和固定工具链。缺少 Git 历史时，`--source-only` 也不能证明满足发布条件。
6. **前端构建会清空目标目录。** Vite 的输出是 `backend/internal/web/dist`，且 `emptyOutDir: true`。历史生产记录要求后端专用更新保留锁定 UI；本地新构建不能直接代表那份生产产物。
7. **生图子应用只有内嵌产物。** `frontend/public/image-studio/` 包含打包后的 JS/CSS 和独立 HTML。修改相关体验时，需先确认子应用源码与构建来源。
8. **开发 Compose 带有环境相关默认值。** 包含宿主代理端口 7897 和已有 Vertex 项目/存储桶名称。首次本地启动应明确覆盖为自己的测试配置，避免把旧环境默认值当作通用配置。

## 本次验证边界

已核对文件结构、关键路由、服务和依赖装配、迁移执行机制、构建命令、历史功能记录及本机工具版本。

本机可用工具为 Go 1.26.4（darwin/arm64）、pnpm 10.33.2、Node 24.15.0；未在 PATH 中找到 Docker 和 golangci-lint。当前目录没有 `frontend/node_modules`、`backend/config.yaml` 和 `backend/internal/web/dist`。

本次仅新增项目指引和本图谱，没有安装依赖、运行应用测试、启动服务、执行迁移或部署。历史文档中的测试通过和生产状态只作为历史记录引用。

后续开展具体开发时，先按模块补齐所需本地环境和聚焦测试；要精确盘点“相对官方改了什么”，还需要对应 Git 历史或明确的上游基线源码。
