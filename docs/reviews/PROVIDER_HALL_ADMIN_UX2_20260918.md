# 供应商大厅管理端配置界面重做验收记录

日期：2026-09-18。范围：本地实现与隔离环境验收，无生产部署、无远程推送、无线上配置变更。分支 `release/v0.1.179-provider-hall`。

## 背景

前一版（`docs/reviews/PROVIDER_HALL_ADMIN_UX_20260918.md`）功能可跑通，但管理员配置路径不顺：分组目标必须逐组进弹窗才能看到，模型档案的模型名靠手输而上游候选已在手边，全局配置字段平铺无状态说明。本轮按「参考 sub2api 其余管理页（`TablePageLayout` + `DataTable`）重做这几个页面」执行，含后端按需扩展。

## 实现行为

### 分组与目标：可展开主表，按组保存

- 页面改用共享 `TablePageLayout`（固定筛选区 + 可滚动表格 + 固定分页），与分组管理、账号管理一致。
- 分组行保留摘要列（分组、平台、上架、目标数与模型、Key 状态、最近探测/验证、操作）；展开行为即原弹窗内的目标编辑表，含模型/协议、专用 Key、探测周期、验证周期、启用、自动调度、测试与查看结果。
- **草稿按分组隔离**：`draft`/`baseline`/`targetVersion`/`keys` 均以 `group_id` 为键，同时展开多组互不覆盖；折叠有未保存修改时先确认。新增「保存本组」，只提交该组，替代原来一次提交整页。
- 展开时才并发拉取该组的目标、Key、模型候选；候选选择与上游刷新也在展开区。
- 上架、展示名、描述、排序值拆到独立的「管理」弹窗（`#hall-listing`），并提供入口说明。保存展示信息会写入该组目标，若此时有未保存的目标草稿会先确认。

### 修复：保存必然版本冲突

原 `payloadFor()` 恒发 `version: 0`，而后端以分组版本做乐观并发校验，导致任何保存都返回 `PROVIDER_HALL_VERSION_CONFLICT`。本轮改为由 `targetVersion[group_id]` 提供，来源是最近一次目标读取或保存响应。该缺陷由隔离 E2E 暴露（`PUT /groups/2/settings -> 409`），并在前端补了针对性用例。

### 分组排序：让 `display_order` 真正生效

- 后端 `ListAdminGroups` 新增 `sort=display_order`：按大厅展示顺序升序，未上架分组排最后，同序按名称与 ID；新增 `providerHallGroupLess`。
- 前端新增「调整展示顺序」拖拽弹窗（`vue-draggable-plus`），保存时用既有批量接口一次写入整批顺序；部分失败会列出未生效的分组。
- 筛选项从与标签不符的「名称升/降序」修正为三个并列选项，并新增「大厅展示顺序」。

### 模型档案：接上游候选，去掉手写来源

- 新增 `GET /api/v1/admin/provider-hall/profile-candidates`：跨全部分组按模型＋协议合并候选，去重后附加 `sources`（账号映射/平台候选/混合路由）与 `groups`（可用于哪些分组）。原按分组的 `GET /groups/:id/models` 保留不动。
- 新增档案时从候选下拉直接选定，模型名与协议**锁定只读**；确需自定义要显式勾选「手动输入模型名」。编辑既有档案不锁定。
- 新增「从上游批量建档」：搜索、隐藏已有档案、全选可见项、一次创建多个；已存在同模型同协议的档案跳过并汇总提示。
- 新增 `DELETE /profiles/:id`（`DELETE /api/v1/admin/provider-hall/profiles/:id`）：在同一事务内先锁档案行再查引用，仍被目标引用时返回 `PROVIDER_HALL_PROFILE_IN_USE` 409 并带引用分组 ID；未被引用则删除。前端删除走 `ConfirmDialog`，错误按既有 `hallError` 映射。
- 列表新增关联分组、上游可用性、基准状态（有效/待确认）及筛选。

### 全局配置：状态摘要与依赖说明

- 新增「当前状态」区：四个开关的当前值与关键细节（预期节点数、运营用户是否就绪、网关 Origin、自动调度）。
- 新增「下一步」提示，按依赖顺序给出当前最该处理的一项（未就绪的运行时 → 空预期节点 → 运营用户 → 网关 → 预算 → 展示）。
- 分区改为：当前状态 / 功能开关与依赖 / 运营用户与预算 / 用户大厅展示 / 采集节点 / 本站网关；未就绪开关在提示中说明缺失原因，不再只说「未启用」。
- 预算为 0 时明确提示不限额，并显示「今日已确认支出 / 预算」；节点区显示发现时间、提供「补齐缺席节点」；进入页面自动调用一次网关健康检查，并显示结果。

### DataTable 扩展

新增可选 `expandable`：桌面渲染 `tr[data-expanded-for]` 详情行，移动端在卡片内展开，支持 `expandedKeys` 受控模式与 `expand` 事件。未开启时行为与之前完全一致。

## 源码范围

- 后端：`repository/provider_hall_management_repo.go`（排序、跨组候选合并、档案删除）、`service/provider_hall_management.go` 与 `service/provider_hall.go`（新错误码 `PROVIDER_HALL_PROFILE_IN_USE`）、`handler/admin/provider_hall_management_handler.go`、`server/routes/provider_hall.go`。
- 前端：`components/admin/provider-hall/{ProviderHallGroups,ProviderHallProfiles,ProviderHallConfigForm}.vue` 重写，`helpers.ts`、`api/admin/providerHall.ts`、`components/common/DataTable.vue` 扩展，中英文案，测试两处，`scripts/provider-hall-admin-smoke.mjs` 改走展开区与展示弹窗。

## 兼容与迁移

无新增迁移，无环境变量变更。新增两个管理接口与一个 `sort` 取值；`ProviderHallModelCandidate` 新增 `sources`/`groups` 为附加字段，原 `source` 保留（合并列表取首个来源），旧接口响应主体不变。分组页的保存语义由「整页提交」变为「按组提交」，但请求体与后端契约一致。

## 本轮验证

环境 macOS arm64、Go 1.26.4、Node 24.15、pnpm 9、PostgreSQL 18、本机 Chrome。

| 检查 | 结果与范围 |
| --- | --- |
| Go 单元 | `go test -tags=unit ./internal/service ./internal/handler/admin ./internal/handler/dto -run ProviderHall` 通过；管理 handler 鉴权用例已纳入 `DeleteProfile` 与 `ListAllProfileCandidates` |
| PostgreSQL localdb | 空库及 228～233 升级路径全部通过；新增档案删除拒绝→解除引用→删除→再删 404、跨分组候选合并去重、展示顺序排序三组契约 |
| 前端静态 | `typecheck`、`lint:check`、`build` 通过 |
| 前端测试 | provider-hall 相关 10 文件 70 项通过（新增：版本号随保存提交、按组独立草稿、展示顺序、档案引用拒绝等）；全量 1763 项通过，唯一失败为既有基线 `OAuthAuthorizationFlow.spec.ts` |
| 隔离端到端 | 自建 PostgreSQL/Redis/模拟上游，经 AUTO_SETUP 启动，`RUN_MODE=standard`；覆盖展开目标区、创建/复用专用 Key、保存本组、立即探测与任务往返、展示信息保存并回读；390/1440/1920 × 明暗共 36 张截图，无页面越界、无浏览器异常 |
| 统计来源 | 探针事实仍为 `source='probe'` 且未进入用户指标，末端 SQL 断言通过 |

过程中的一次真实失败记录：E2E 首次运行在 `PUT /groups/2/settings` 返回 409，定位为前端未提交分组版本号；修复后复跑全绿，并补了对应单测。

复现命令（仓库根目录）：

```sh
(cd backend && go test -tags=unit ./internal/service ./internal/handler/admin ./internal/handler/dto -run ProviderHall -count=1)
(cd backend && PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin go test -race -tags=providerhall_localdb ./internal/repository -run '^TestProviderHallLocalDatabase$' -count=1 -timeout=10m)
npx --yes pnpm@9 --dir frontend run typecheck
npx --yes pnpm@9 --dir frontend run lint:check
npx --yes pnpm@9 --dir frontend run test:run
PROVIDER_HALL_BROWSER_CHANNEL=chrome npx --yes pnpm@9 --dir frontend run test:e2e:provider-hall
```

## 验证边界与后续

- 未跑 Docker integration、未跑全量 CI、未做生产部署。前端构建保留了既有的 Browserslist/PostCSS 与分块体积告警。
- 管理端宽表在窄屏下仍是容器内横向滚动，未做移动端专用的精简表单；分组展开区在 390px 下以横向滚动为主。
- **降智探测尚未实现**：仓库内目前没有相关实现。可复用 `service/provider_hall_verification.go` 的断言套件框架（`ProviderHallBuildSuite` + `ProviderHallVerdict`）扩展为目标模型的质量阶梯用例，沿用探针任务与预算约束；需单独一轮设计与验收，因为每次探测都产生真实费用。
- 分组排序弹窗一次最多处理 100 组（沿用批量接口上限），超出时会提示；如有更多分组需要另行分批。
