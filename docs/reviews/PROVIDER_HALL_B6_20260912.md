# 供应商大厅 B6：管理 API 扩展与后台页签实施记录

日期：2026-09-12。本批在 B4 任务 Runner 之上补齐管理端的任务操作与运行状态视图。用户大厅（B5/B7）与全流程验收（B8）由其他批次负责；本批没有部署，也没有对真实上游发出付费请求。

## 已实现的行为

后端（均在 `/api/v1/admin/provider-hall` 之下，沿用管理员鉴权：未登录 401、非管理员 403）：

- `POST /groups/:id/probes`、`POST /groups/:id/verifications`：请求体 `{profile_id, idempotency_key}`，成功返回 202 `{job_id, status, reused}`。幂等键最长 128 字符、不含空白；同一幂等键重复提交返回既有任务并标 `reused:true`。Runner 的业务拒绝按原码透传为 400：`PROVIDER_HALL_TASKS_DISABLED`、`PROVIDER_HALL_BUDGET_EXHAUSTED`、`PROVIDER_HALL_TARGET_DISABLED`；Runner 未接线（api 实例或缺省构建）返回 409 `PROVIDER_HALL_NOT_READY`。请求体拒绝未知字段，无效输入不会触达 Runner。
- `GET /jobs?status&kind&group_id&page&page_size`：分页列表（`page_size` 1～100，默认 20），状态/类型取值校验；无任务仓库时返回空列表而非报错。
- `GET /jobs/:id`：任务 + 全部样本 + 检测报告（探测任务为 `null`）。响应只含 Key ID 等标识，不含 Key 材料、上游凭据或请求正文；不存在返回 404。
- `POST /jobs/:id/cancel`：排队任务或样本尚未发出的运行中任务取消；已有样本发出返回 409 `PROVIDER_HALL_JOB_IN_FLIGHT`；已结束返回 409 `PROVIDER_HALL_JOB_NOT_CANCELLABLE`。
- `GET /groups/:id/verifications?profile_id&page&page_size`：报告列表，包含未上架分组；`page_size` 上限 50。
- `GET /health`：采集（开关、每节点最新代次、缺席的预期节点、未关闭缺口、本机队列深度/容量/丢弃数）、聚合器（水位、滞后秒数、脏分钟数、最近运行与错误；进程内状态比持久化状态更新时优先）、对账（待对账、不确定、24 小时内失败）、预算（Asia/Shanghai 预算日、预算、已确认/不确定支出、在途样本、暂停原因 `tasks_disabled` 优先于 `budget_exhausted`）、任务计数。心跳超过 45 秒的节点即显示失联，不等聚合器关闭代次。

前端 `/admin/provider-hall`：

- 页签扩展为 `全局配置 | 分组与目标 | 模型档案 | 任务 | 运行状态`，任务与运行状态面板首次进入时才加载。
- 分组与目标：已保存且启用的目标行提供“立即探测 / 立即检测”。每次点击生成一次幂等键并保留到响应返回，双击不会创建两个任务；业务拒绝（任务未启用、预算耗尽、目标未启用、未就绪）后丢弃该键，传输失败保留该键以便重试。成功与“已存在相同任务”分别提示。
- 任务页签：类型/状态/分组筛选、分页、查看详情、取消（`ConfirmDialog` 二次确认）。409 在途冲突以明确文案呈现且不清空列表。
- 任务详情弹窗：任务字段、检测报告判定/执行状态/各套断言与模型比对结果、样本表（状态、结果、首 Token/端到端、Token、响应模型、计费）。
- 运行状态页签：`StatCard` 概览、节点表、缺口列表、聚合/预算/对账明细，`AutoRefreshButton` 15/30/60 秒（默认 30 秒），页面隐藏时暂停。
- 中英文文案同步；`helpers.ts` 新增五个业务错误码映射与幂等键生成。

## 源码

- 后端新增：`internal/service/provider_hall_admin.go`（健康读模型与 `ProviderHallAdminRepository` 接口）、`internal/repository/provider_hall_admin_repo.go`、`internal/handler/admin/provider_hall_task_handler.go`；`internal/repository/wire.go` 增一行 `NewProviderHallAdminRepository`。
- 后端修改：`internal/handler/admin/provider_hall_handler.go`（构造函数接入 Runner、任务仓库、Collector、Aggregator、管理仓库，均可为 nil）、`internal/handler/dto/provider_hall.go`（入队与健康 DTO，`ProviderHallAdminHealth`）、`internal/server/routes/provider_hall.go`（七条新路由）、`cmd/server/wire_gen.go`（`go generate ./cmd/server`）。
- 测试：`internal/handler/admin/provider_hall_task_handler_test.go`（`TestProviderHallAdminHandler_{TaskAuthentication,Probes,Verifications,Jobs,Cancel,Health}`）、`internal/repository/provider_hall_admin_db_test.go`（localdb 合约 `b6_admin`）；既有 handler 测试更新构造函数调用。
- 前端新增：`components/admin/provider-hall/ProviderHallJobs.vue`、`ProviderHallJobDetailDialog.vue`、`ProviderHallHealth.vue`、`__tests__/ProviderHallTasks.spec.ts`。
- 前端修改：`api/admin/providerHall.ts`（任务/样本/报告/健康类型与七个方法）、`views/admin/ProviderHallView.vue`、`components/admin/provider-hall/ProviderHallGroups.vue`、`helpers.ts`、`i18n/locales/{zh,en}/admin/providerHall.ts`、`api/__tests__/admin.providerHall.spec.ts`、`scripts/provider-hall-admin-smoke.mjs`。

## 本轮验证

以下为 2026-09-12 实际执行结果（Go 1.26.5、Node 24.15.0、pnpm 10.33.2）：

1. `gofmt -l`（本批文件）无输出；`go build ./...`、`go vet -tags=unit ./internal/... ./cmd/...`、`go generate ./cmd/server` 通过。
2. `go test -tags=unit -race ./internal/handler/admin ./internal/handler/dto ./internal/service ./cmd/server -run ProviderHall -count=1`：四个包全部 ok；admin + dto 包 ProviderHall 用例共 39 个通过（其中本批新增 6 个顶级用例、22 个子用例）。`go test -tags=unit ./cmd/server` 通过。
3. PostgreSQL 18 隔离集群（`PROVIDER_HALL_PG_PORT=15484`）：`go test -tags=providerhall_localdb ./internal/repository -run ProviderHall -count=1 -v` 通过，六条路径（empty、upgrade_228～upgrade_232）各执行 `admin_health_reads` 合约（最新代次去重、只列未关闭缺口、聚合状态、脏桶与对账计数、24 小时窗口）。
4. 前端：`pnpm run lint:check`、`pnpm run typecheck` 均无错误；`vitest run src/components/admin/provider-hall src/views/admin/__tests__/ProviderHallView.spec.ts src/api/__tests__/admin.providerHall.spec.ts src/i18n`：12 个文件 56 个用例通过（本批新增 9 个组件用例、3 个 API 契约用例，覆盖双击复用幂等键、业务拒绝后换键、禁用目标不可点、列表/筛选/取消、409 呈现、详情报告摘要、健康页节点与暂停原因、中英文键一致）。
5. `node --check` 与 ESLint 检查通过 `scripts/provider-hall-admin-smoke.mjs`。

## 未执行项与限制

- Docker 不可用，`integration` 标签未运行；数据库合约以 localdb 路径覆盖。
- Playwright 冒烟脚本已扩展（任务/运行状态页签截图、主动任务关闭时“立即探测”返回 400 且不创建任务；设置 `PROVIDER_HALL_ENABLE_TASKS=1` 时才开启开关并断言 202），但本批没有对隔离实例实际运行，留待 B8 全流程验收。
- 健康视图中的“缺席的预期节点”只比对当前存活代次；预期节点名单为空时不报缺席。
- 手动任务的实际派发、预算扣减与报告生成由 B4 Runner 负责，本批只暴露入队、查询与取消。
