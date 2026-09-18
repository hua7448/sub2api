# 供应商大厅：目标配置与后台页面实施记录

日期：2026-09-10。本轮在第一批配置、权限和算法基础上继续实施，未完成完整大厅，未部署。

## 已实现的行为

- 新增 `GET/PUT /api/v1/admin/provider-hall/groups/:id/targets`。目标集合与分组展示配置共用版本，事务内校验分组、模型档案、运营用户、专属/订阅权限、Key 归属和实际混合路由；Messages 仍受原分组入口开关约束。
- 整组保存只接受可写字段，前端 API 层剔除服务端行 ID 和审计元数据，HTTP 层拒绝未知字段；未提交的既有目标停用并保留 ID。每组最多 100 个档案，探测间隔默认 300 秒，检测间隔默认 86400 秒。失效 Key 不妨碍停用目标。
- 专用 Key 首次登记必须未使用；已经登记的 Key 仅允许原运营用户和原分组复用。来源登记只含 ID 和时间，停用目标或删除 Key、用户、分组后仍保留。变更/清空运营用户会在同一事务内停用已启用目标并增加相关分组版本。
- V2 的请求指标、用户指标、延迟直方图和错误聚合排除登记时间之后的探测 Key 请求，保留登记之前的普通流量。原有聚合公式和保留策略不变。
- 后台新增 `/admin/provider-hall`，入口位于“渠道管理”。包含全局配置、模型档案、分组展示信息和目标绑定；管理员可以通过页面完成保存，金额始终使用十进制字符串。
- 页面具备保存冲突提示、失败保留输入、未保存离开提醒、分组切换取消旧请求、运营用户搜索、Key 分页加载和模型基准确认；同步中英文。既有目标保留停用状态，新建未保存行可移除。
- 采集、用户展示、主动任务仍默认关闭，后端继续拒绝启用尚未接入的运行功能。没有修改真实网关转发、扣费路径、V1 探测器或现有后台任务角色规则。

## 源码和迁移

- `backend/migrations/230_provider_hall_targets.sql`：新增 targets 和永久 probe_keys；没有修改已应用的 229 迁移。
- `backend/ent/schema/provider_hall.go` 及 Ent 生成输出。
- `backend/internal/service/provider_hall_target.go`、`backend/internal/repository/provider_hall_target_repo.go`，以及已有 config repository 的运营用户变更事务。
- 管理 handler、DTO、路由和对应测试；`channel_monitor_v2_aggregation.go` 的四处来源过滤。
- `frontend/src/views/admin/ProviderHallView.vue`、`components/admin/provider-hall/`、管理 API、路由、侧栏和中英文 locale。用户 Key 查询 API 增加可选分页和 AbortSignal，原调用保持兼容。
- 用 pnpm 9 添加 Playwright 1.58.2 和管理页真实接口验收脚本 `frontend/scripts/provider-hall-admin-smoke.mjs`，锁文件仍为 9.0。

## 本轮验证

已执行并通过：

1. Go 1.26.5 聚焦 `ProviderHall|ChannelMonitorV2` 单元和竞态检查，包含 service、admin handler、DTO 和 repository 包。
2. PostgreSQL 18.6 隔离集群：空库、从 228 迁移状态升级、从 229 迁移状态升级，三条路径均执行重复迁移和基础合约。每条路径新增 5 组真实目标/来源合约，共 15 组，覆盖原子回滚、CAS 竞争、停用和重新启用、权限/订阅、混合路由、删除后的来源保留和 V2 四类聚合。
3. Go 1.26.6 的大厅 service/admin 聚焦兼容检查；repository 包在该 `unit` + `ProviderHall` 筛选下没有匹配用例。
4. `make -C backend build`。
5. Node 20/pnpm 9：管理 API 6 项、后台组件 7 项测试通过；涉及文件 ESLint 通过；新增依赖后的 frozen-lockfile 安装通过。
6. Playwright 1.58.2 + 本机 Chrome 152：真实本地接口完成预算精度往返、创建档案、上架分组和绑定专用 Key；三个运行开关最终仍关闭，没有发出付费探测。390/1440/1920 像素的三个页签、390/1440 像素的编辑弹窗均完成明暗样式检查，共 22 张最终截图；检查整页横向溢出、按钮内容截断、键盘页签和浏览器异常。修复了表头隐藏文字导致的移动端溢出，以及本页和弹窗的深色背景对比问题。

最终浏览器结果见 [summary.json](artifacts/provider-hall-admin-20260910/summary.json)，截图位于同目录，例如 [移动端目标配置](artifacts/provider-hall-admin-20260910/390-light-targets.png)、[桌面深色弹窗](artifacts/provider-hall-admin-20260910/1440-dark-profile-dialog.png)。当前应用主题 store 仍锁定浅色，本轮通过现有 `.dark-theme` 选择器验证新增页面的深色样式，没有改动全站主题策略。

全量 `vue-tsc --noEmit` 已执行，仍报告本轮未改动的 `AccountUsageCell.vue`、`AccountTestModal.vue`、`GroupsView.vue`、`EmailVerifyView.vue` 和 `RegisterView.vue` 错误；新增页面、组件和测试没有类型错误。没有重新将前次全量后端失败作为本轮新结果，也没有宣称全量检查或 CI 通过。

本轮本地独立 Redis 为新安装的 Homebrew Redis 8.10.1，未设置开机服务。最初 AUTO_SETUP 尝试因空数据库密码的既有 DSN 行为和尚未创建 DATA_DIR 失败；补齐本地连接配置和目录后，通过同一正常初始化入口完成安装。没有手工写 `.installed` 或绕过数据库迁移。

### 复现

```sh
GOPROXY=https://goproxy.cn,direct go -C backend test -tags=unit -race ./internal/service ./internal/handler/admin ./internal/handler/dto ./internal/repository -run 'ProviderHall|ChannelMonitorV2' -count=1
PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin go -C backend test -tags=providerhall_localdb ./internal/repository -run ProviderHall -count=1 -v
pnpm --dir frontend exec vitest run src/components/admin/provider-hall/__tests__/ProviderHallAdmin.spec.ts src/api/__tests__/admin.providerHall.spec.ts
```

管理页浏览器脚本要求已有正常初始化的独立本地实例、Node 20/pnpm 9，以及 `pnpm --dir frontend exec playwright install chromium` 安装的浏览器；也可通过 `PROVIDER_HALL_BROWSER_CHANNEL=chrome` 使用本机 Chrome。只允许 loopback URL，测试会修改配置并保留测试记录以供检查，不能用于生产或共用实例。

```sh
# 凭据由新建的本地实例提供，不使用历史配置。
PROVIDER_HALL_ADMIN_EMAIL=... PROVIDER_HALL_ADMIN_PASSWORD=... \
PROVIDER_HALL_FRONTEND_URL=http://127.0.0.1:3000 \
PROVIDER_HALL_BACKEND_URL=http://127.0.0.1:8081 \
pnpm --dir frontend run test:e2e:provider-hall-admin
```

该脚本只覆盖本轮管理配置闭环，不等同于计划中尚未实现的“模拟上游请求→指标→报告→用户使用分组”全链路验收脚本。

## 未完成范围

- 三协议真实请求 tracker、采集队列、计费关联、节点确认水位和缺口处理。
- 请求事实、分钟聚合、快照、迟到数据重算、历史查询与保留清理。
- primary 任务领取、租约、发出状态机、预算、对账、探测及模型一致性报告。
- 用户大厅、趋势、报告、Key 创建/切换闭环和相应鉴权、定价、浏览器与性能验收。
- 计划 A～H 全量验收、正式构建来源校验和发布。没有生成嵌入 UI 或改动生产 UI revision/cache。

本地开发实例按正常 AUTO_SETUP 流程初始化，使用独立 DATA_DIR、PostgreSQL 和 Redis，RUN_MODE 为 standard。具体本地地址、测试账号和停服命令保存在 `/tmp/provider-hall-dev.uNxBDT/README.md`；这些不是生产配置。
