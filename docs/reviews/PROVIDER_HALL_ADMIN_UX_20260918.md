# 供应商大厅管理端完善验收记录

日期：2026-09-18。范围：本地实现与隔离环境验收，无生产部署、远程推送或线上配置变更。

## 来源与分支

本轮实际检查发现工作区存在 Git，开始于 `release/v0.1.179-provider-hall` 的 `c86a99f8a`，在 `feature/provider-hall-admin-ux` 实施。历史项目地图中“无 Git”的描述仅代表当时快照。交付时将功能提交本地合入原 release 分支并切回该分支；用户已有的未跟踪 `site/` 未修改、未纳入提交。

## 实现行为

- 默认进入分组列表，使用共享 DataTable、筛选和分页，批量返回未删除分组、平台支持状态、目标、模型、Key 状态、最近探测/验证与实际回退档案。分组展示和目标在同一编辑弹窗保存。
- 模型候选来自账号映射、混合路由和平台列表；仅主动刷新调用上游模型列表，逐账号返回结果，不改账号映射。通配规则不导入具体目标，本站请求模型和上游模型分别展示。选中模型预览后按模型＋协议复用档案；新增档案不自动确认工具能力、响应别名或参考价格，新增目标默认停用。
- 专用 Key 仅返回标识和状态。事务锁协调运营用户变更与 Key 创建，优先复用当前用户/分组的有效已登记 Key；并发点击与重试不重复创建。被改绑的旧 Key 被跳过，不自动充值或启用任务。
- 统一设置保存采用共享分组版本号；预检查复用写入校验但不写入、登记 Key 或发送模型请求。保存成功而入队失败时保留新版本，仅重试入队；入队携带刚保存的目标及档案版本，并保留幂等键直到收到成功响应。
- 全局与目标分别增加自动调度开关，任务总开关继续限制所有付费任务。手动任务允许未上架、自动开关关闭；定时任务额外要求上架及两个自动开关开启。队列清理、入队和每个样本派发均重查；已经发出的样本继续收尾/对账。
- 批量操作支持上下架、目标启停、自动调度和周期，最多 100 组，必须明确选择目标并预览。逐组事务及版本检查，返回结果和可翻译原因；失败项保留，可重新加载其版本后预览重试。批量请求单独限制为 4 MiB，其他写接口保留 64 KiB。
- 模型档案支持搜索/分页、来源分组与模型选择、自定义输入、默认和共享引用信息、参考价及有效状态。配置页展示开关依赖、运营用户状态、变更停用影响、节点选择、分钟/小时/百分比单位及预算 0 不限额。默认档案身份修改返回明确原因与配置入口。
- 网关连接检查只请求固定 `/health`，复用既有地址约束、5 秒超时及禁止重定向，不携带 Key。任务可按分组名称、模型、来源、类型、状态筛选，运行中刷新，详情和目标双向跳转。路由参数统一写入，防止切换页签覆盖定位信息。
- 运行状态优先呈现异常、对账、预算、任务和节点；队列、水位、聚合等信息置于诊断详情。中英文同步，移动端宽表在容器内横向滚动。

新增接口位于 `/api/v1/admin/provider-hall`：`GET /groups`、`GET /groups/:id/models`、`POST /groups/:id/models/refresh`、`GET/POST /groups/:id/probe-keys`、`POST /config/preflight`、`POST /groups/:id/preflight`、`POST /config/check-gateway`、`PUT /groups/:id/settings`、`POST /groups/batch`。原接口与响应主体保留，档案/任务/健康响应增加辅助字段；管理端鉴权覆盖新增入口。

## 迁移与兼容

新增前向迁移 `234_provider_hall_admin_schedule.sql`，未改写既有迁移。旧全局自动调度回填 `tasks_enabled`，旧目标回填开启；新配置/目标自动调度默认关闭。旧客户端省略自动字段时保留现值，新建目标使用关闭默认。Ent schema、生成代码、DTO、前端类型及 Wire 一并更新。

目标版本只在执行配置变化时递增，避免编辑展示信息或其他目标使任务失效。**自动调度单独变化时不递增执行版本**，由共享分组版本防止覆盖，并在队列/派发重查自动开关；这样关闭自动调度不会误取消合法手动任务。上架开关同样不取消手动测试。

预检查返回首个字段或目标级问题，属于提示；实际写入和派发重新检查。用户大厅统计、定价、权限与默认档案回退公式保持现有实现。

## 本轮验证

环境为 macOS arm64、Go 模块工具链 1.26.5、Node 24.15、pnpm 9.15.9、PostgreSQL 18、隔离 Redis、Playwright＋本机 Chrome。使用冻结锁文件安装，未修改 pnpm 锁文件。

| 检查 | 结果与范围 |
| --- | --- |
| Go 大厅单元竞态 | service、admin handler、DTO 的 ProviderHall 聚焦测试通过；包括固定健康接口、不跟随重定向、新管理接口鉴权 |
| PostgreSQL localdb 竞态 | 空库及从 228、229、230、231、232、233 升级的 7 条路径通过；每条覆盖重复迁移、数据库合约、并发保存、8 路并发 Key 创建、改绑 Key 复用、无写预检查、原子保存、旧字段省略、自动开关全部 8 种组合 |
| 前端 | typecheck、完整只读 lint、构建通过；10 个相关 Vitest 文件共 69 项通过，含用户大厅回归、冲突保留草稿、入队失败只重试、Key 防双击、批量部分失败重试、跨页路由定位 |
| 隔离端到端 | 自建 PostgreSQL、Redis、DATA_DIR、模拟上游，经正常 AUTO_SETUP，RUN_MODE=standard；完成模型候选/刷新、Key 并发复用、未上架手动测试、版本冲突、入队幂等、批量预览/执行/部分失败、浏览器保存并测试及任务往返 |
| 统计来源 | 直接验证本轮真实探测落库为 probe 且关联 sample，未出现对应 user 请求，未污染用户成功率/延迟聚合 |
| 布局 | 390/1440/1920 × 明暗模式，五页签与分组弹窗共 36 张截图；无页面或弹窗越界，无浏览器 JS 异常；人工抽查移动编辑、桌面编辑、移动配置和桌面状态 |

最终 E2E 的模拟请求计数：2 次大厅推理（API 手动测试及浏览器保存并测试各 1 次）、1 次主动模型列表刷新、1 次原有账号能力检查。账号能力检查是创建测试账号时的既有行为；大厅候选查询、配置刷新、预检查与健康检查均不触发推理。

截图和结果保留于本地 `docs/reviews/artifacts/provider-hall-admin-ux-20260918/`（按仓库规则忽略，不纳入源码提交）。浏览器深色检查直接切换 `dark-theme` 类，验证组件的深色样式；既有全局主题策略没有改动。测试引导教程的已读标记仅注入隔离浏览器。

复现命令（仓库根目录）：

```sh
(cd backend && go test -race -tags=unit ./internal/service ./internal/handler/admin ./internal/handler/dto -run ProviderHall -count=1)
(cd backend && PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin go test -race -tags=providerhall_localdb ./internal/repository -run '^TestProviderHallLocalDatabase$' -count=1 -timeout=8m)
npx --yes pnpm@9 --dir frontend run typecheck
npx --yes pnpm@9 --dir frontend run lint:check
npx --yes pnpm@9 --dir frontend run test:run src/components/admin/provider-hall/__tests__ src/views/admin/__tests__/ProviderHallView.spec.ts src/components/provider-hall/__tests__ src/views/user/__tests__/ProviderHallView.spec.ts src/composables/__tests__/useProviderHall.spec.ts src/api/__tests__/admin.providerHall.spec.ts
npx --yes pnpm@9 --dir frontend run build
PROVIDER_HALL_BROWSER_CHANNEL=chrome npx --yes pnpm@9 --dir frontend run test:e2e:provider-hall
```

测试脚本可用 `PROVIDER_HALL_PG_BIN`、`PROVIDER_HALL_REDIS_BIN`、`PROVIDER_HALL_BROWSER_CHANNEL`、`PROVIDER_HALL_ARTIFACTS` 覆盖本机位置。脚本结束清理自建服务和数据，只留下验收证据。

## 验证边界与发布

本轮没有重跑全仓库 Go/Vitest、Docker integration 或生产合约脚本，不把聚焦通过视为完整 CI 兼容。Node 与 CI 的 Node 20 不同，Go 1.26.5/1.26.6 的既有 CI、Docker、生产脚本冲突未修改。前端构建有现有 Browserslist、PostCSS、混合导入等告警，但退出成功。构建仅写本地嵌入目录，没有替换生产 UI。

后续发布仍需明确版本标签、数据库备份和隔离试运行。**回退不识别自动调度限制的旧程序前，必须先关闭 `tasks_enabled` 总开关。** 前向迁移保留，无 Down SQL。本轮未部署。
