# 供应商大厅：用户大厅前端（B7）实施记录

日期：2026-09-12。本批次实现用户端"供应商大厅"页面（`/providers`）的前端部分，对应实施计划 B7。后端用户 API（B5）由另一并行批次实现；本批次以已定稿的 DTO 合约 `backend/internal/handler/dto/provider_hall_user.go` 为准，用模拟数据完成开发与测试，尚未对接真实后端运行，未部署。

## 已实现的行为

- 新增用户页面 `/providers`（路由名 `ProviderHall`，需登录），侧边栏"渠道状态"之后新增"供应商大厅"入口，由公共设置 `provider_hall_enabled` 控制（`opt-in` 模式：后端未注入或为 `false` 时隐藏）。该设置是后端根据大厅配置的"用户展示"开关派生的只读值；系统设置页"功能开关"页签新增一张只读说明卡片并链接到 `/admin/provider-hall`，不提供独立 Toggle。
- 页面顶部：监测中徽章、说明、数据截至时间（`data_through`）、手动刷新和自动刷新按钮（30/60/120/300 秒，默认 60 秒，页面隐藏或加载中时暂停）。
- 概览：可用分组环形数字、公开分组、异常线路、检测通过，以及可用/异常/未知图例。数值全部来自后端 `summary`。
- 工具栏：时间维度 chips（6h/24h/7d/30d）、排序 chips（默认 / 倍率 / 真实价格 / 平均（最快95%）/ 缓存命中 / 成功率 / 自定义；单击循环 降序→升序→关闭）、自定义三级排序弹窗、模型筛选（来自 `catalog.models`，值为 `model::protocol`）、分组搜索（300 ms 防抖，≤100 字符）。所有筛选、排序、页码写入路由 query；排序偏好另存 `localStorage` 键 `providerHall.sort.<userId>`，路由 query 优先于本地偏好。
- 列表：桌面端固定列宽表格（sticky 表头，行可聚焦，Enter/Space 展开，↑/↓ 在行间移动焦点），768px 以下切换为卡片列表。每行显示分组名/说明、倍率、真实价格与预测倍率（带计算说明 tooltip）、健康状态与每档案健康点、平均（最快95%）TTFT 数值 + 快/中/慢 + 相对页面最大值的条、1h 缓存命中环、1h 成功率环（评分色）、探测 vs 最快95% 双线迷你图（缺口区带）、模型检测判定 pill 与报告入口、"使用此分组"按钮。
- 值渲染只依赖后端 `Metric.state` 与 `reason_code`：`ok` 显示值；`stale` 显示值并加"数据过期"角标；`incomplete` 有值则显示值并标注；`insufficient/disabled/not_applicable` 显示按 `reason_code` 翻译的占位文案，无翻译时回退到 state 文案。价格、倍率均为十进制字符串，前端只做字符串截断和千分位，不做浮点运算。
- 展开详情：指标网格（倍率、输入价、1h 成功率、6h 端到端可用率、1h 缓存率、首 Token、端到端、平均（最快95%）、默认档案 P90、TPS、最近监测）、三线趋势图（探测端到端 / 探测首 Token / 真实最快95%，同一 ms 轴，缺口桶以插件绘制半透明区带）、模型健康 chips、最近探测 Token 输入/输出。详情按"分组 + 时间维度"缓存，自动刷新时重新拉取已展开分组。
- 检测报告弹窗：报告列表（可按档案筛选、分页）与单份报告（判定、执行状态、原因、完成/有效期、算术/JSON/工具/模型名比对四个套件的计数与响应模型名）。
- "使用此分组"弹窗：切换现有 Key（`keysAPI.list` 分页加载，显示 原分组 → 新分组，已在该组的 Key 禁用；`keysAPI.update(id, { group_id })`）与创建新 Key（`keysAPI.create(name, groupId)`）。失败时通过 `showError(extractApiErrorMessage)` 提示并保留选择/输入。
- 刷新失败保留旧数据并显示"数据过期"横幅；首次加载失败显示错误与重试；迟到的旧请求响应被序号机制丢弃。

## 配色与图形

按 `dataviz` 技能流程执行：真实最快95%（蓝 `#2a78d6` / 暗色 `#3987e5`）、探测端到端（橙 `#eb6834` / `#d95926`）、探测首 Token（青 `#1baf7a` / `#199e70`）为分类色槽 1–3，`validate_palette.js` 明暗两种模式全部通过（明色模式青色对比度 2.74:1 给出 WARN，已用图例 + 直接标签作为补救）。成功率环、TTFT 条使用状态色（good `#0ca30c` / warning `#fab219` / critical `#d03b3b`）并始终伴随文字标签；缺口区带为 critical 色 16% 透明度。趋势图只有一个 y 轴（毫秒）。暗色主题仅通过 `dark:` 与 `.dark-theme` 选择器实现。

## 源码

新增：

- `frontend/src/types/providerHall.ts`（与后端 DTO 一一对应的类型，`ProviderHallMetric<T>` 与管理端共用）
- `frontend/src/api/providerHall.ts`（`list / getGroup / listVerifications / serializeSort`）
- `frontend/src/composables/useProviderHall.ts`
- `frontend/src/utils/providerHallFormat.ts`
- `frontend/src/views/user/ProviderHallView.vue`
- `frontend/src/components/provider-hall/`：`HallSummary`、`HallToolbar`、`HallSortDialog`、`HallTable`、`HallRow`、`HallDetailPanel`、`HallRing`、`HallSparkline`、`HallTtftBar`、`HallModelDots`、`HallVerificationBadge`、`HallVerificationDialog`、`HallMetricTooltip`、`UseGroupDialog`、`hallGeometry.ts`
- `frontend/src/i18n/locales/zh/providerHall.ts`、`en/providerHall.ts`
- 测试：`components/provider-hall/__tests__/{testUtils.ts,HallRing,HallSparkline,HallTable,UseGroupDialog}.spec.ts`、`composables/__tests__/useProviderHall.spec.ts`、`views/user/__tests__/ProviderHallView.spec.ts`

修改：

- `frontend/src/router/index.ts`（`/providers`）、`components/layout/AppSidebar.vue`（入口）、`utils/featureFlags.ts`（`providerHall` opt-in 标志）、`types/index.ts`（`PublicSettings.provider_hall_enabled?`）
- `i18n/locales/{zh,en}/index.ts`（展开 `providerHall`）、`common.ts`（`nav.providerHall`）、`admin/settings.ts`（功能开关页只读卡片文案）
- `views/admin/SettingsView.vue`（只读说明卡片）、`api/admin/settings.ts`（`SystemSettings.provider_hall_enabled?` 只读字段）
- `api/admin/providerHall.ts`：删除本地重复的 `ProviderHallMetric` 定义，改为从 `@/types/providerHall` 重新导出（管理端测试全部通过）

未改动 `backend/` 及管理端组件。

## 本轮验证

环境：Node v24.15.0、pnpm 10.33.2（项目目标为 Node 20 / pnpm 9，本机未安装，本轮用已安装版本执行）。以下为 2026-09-12 实际执行结果：

1. `pnpm --dir frontend run lint:check`：通过（exit 0）。
2. `pnpm --dir frontend run typecheck`：通过，0 个错误。前次记录中列出的 `AccountUsageCell.vue` 等基线错误本轮未出现（全量 `vue-tsc --noEmit` exit 0）。
3. `pnpm --dir frontend exec vitest run src/components/provider-hall src/composables/__tests__/useProviderHall.spec.ts src/views/user/__tests__/ProviderHallView.spec.ts src/i18n`：14 个文件 63 个用例通过，其中新增 35 个：
   - HallRing 7（弧长/clamp、评分色阈值、ok/stale/insufficient/未知 reason 渲染）
   - HallSparkline 6（缺口合并、null 断线、共用 y 轴、空数据、组件矩形与路径）
   - HallTable 6（state 占位、相对最大值条、移动端卡片、键盘 Enter/Space/方向键、按钮不触发展开、详情按 range 隔离）
   - UseGroupDialog 4（切换只发 `group_id`、创建参数与失败保留输入、切换失败保留选择、无 Key）
   - useProviderHall 7（筛选写路由并重置页码、三级排序持久化与循环、路由优先、迟到响应丢弃、失败保留旧值+stale、隐藏页暂停、详情缓存）
   - ProviderHallView 5（整页挂载、工具栏切换、展开与报告弹窗、使用分组弹窗、stale 横幅与错误态）
   - 既有 `localesMessageCompile`、`localesNoKeyCollision` 等 i18n 测试通过。
4. 回归：`src/api/__tests__/admin.providerHall.spec.ts`、`src/components/admin/provider-hall`、`src/views/admin/__tests__/ProviderHallView.spec.ts`、`src/utils/__tests__`、`src/components/layout`：31 个文件 169 个用例通过。
5. `dataviz` 技能 `validate_palette.js` 明/暗模式各一次，全部 PASS（明色模式一项 WARN，见上）。

## 未验证与限制

- 未与真实后端联调：B5 用户 API 由并行批次实现，本批次全部基于 mock 数据；未启动应用、未截图、未做 390/1440/1920 明暗浏览器检查，这些属于 B8。
- 侧边栏入口依赖后端把 `provider_hall_enabled` 加入 `PublicSettings` 与 SSR 注入 payload（`PublicSettingsInjectionPayload`），否则 opt-in 标志在刷新后会隐藏入口直到异步设置返回；后端侧属 B5 范围。
- 趋势图在 jsdom 中以 stub 替代 `vue-chartjs`，缺口区带插件的绘制逻辑未在测试中执行，需在 B8 浏览器验收中目视检查。
- 真实价格货币前缀依据 `quote.unit`（`quota_per_million` → ¥，否则 $）；若后端历史均价单位与当前报价单位不一致（`mixed_billing`），行内按 `not_applicable` 占位显示。
- 未提供中文以外的排序偏好迁移；用户切换账号后按 `userId` 隔离偏好。
