# 供应商大厅 B8 汇总：批次索引与整合验证

日期：2026-09-12。本文是 B2～B8 全部批次的索引和最终整合验证记录；每个批次的行为、源码范围、实际命令与未执行项以各自记录为准。未部署生产。

## 批次记录

| 批次 | 内容 | 记录 |
|---|---|---|
| B2 | 采集、计费关联、事实表（迁移 231） | [PROVIDER_HALL_B2_20260912.md](PROVIDER_HALL_B2_20260912.md) |
| B3 | 聚合、快照、对账、保留（迁移 232） | [PROVIDER_HALL_B3_20260912.md](PROVIDER_HALL_B3_20260912.md) |
| B4 | 任务 Runner、探测、模型检测、预算（迁移 233） | [PROVIDER_HALL_B4_20260912.md](PROVIDER_HALL_B4_20260912.md) |
| B5 | 用户 API、定价叠加、排序、公开设置 | [PROVIDER_HALL_B5_20260912.md](PROVIDER_HALL_B5_20260912.md) |
| B6 | 管理任务 API、jobs/health 页签 | [PROVIDER_HALL_B6_20260912.md](PROVIDER_HALL_B6_20260912.md) |
| B7 | 用户大厅前端 `/providers` | [PROVIDER_HALL_B7_20260912.md](PROVIDER_HALL_B7_20260912.md) |
| B8 端到端 | 隔离实例、真实流量、主动任务、浏览器截图 | [PROVIDER_HALL_B8_E2E_20260912.md](PROVIDER_HALL_B8_E2E_20260912.md) |
| B8 性能与回归 | 矩阵 H、网关开销、全量回归、基线失败 | [PROVIDER_HALL_B8_PERF_20260912.md](PROVIDER_HALL_B8_PERF_20260912.md) |

## 并行实施的共享接缝

各批次由并行代理实施，为避免共享文件冲突预先建立了以下接缝，后续维护需要沿用：

- `service/provider_hall_aggregator.go` 与 `provider_hall_runner.go` 以 `Ready()` 报告运行时是否就绪；`ProvideProviderHallService` 由此推导 `collection_enabled`/`tasks_enabled` 的可启用性，`display_enabled` 由 `ProvideProviderHallQueryService` 置位。`cmd/server/wire.go` 的清理顺序为 Runner → Aggregator → Collector。
- `repository/provider_hall_localdb_test.go`：`PROVIDER_HALL_PG_PORT` 指定集群端口；升级路径由迁移文件自动发现（`empty`、`upgrade_228` … `upgrade_232`）；各批次通过 `registerProviderHallLocalDBContract` 在 `init()` 登记数据库合约。
- 用户 API JSON 合约固定在 `backend/internal/handler/dto/provider_hall_user.go`，前端 `src/types/providerHall.ts` 与之对应。
- 共享模拟上游：Go `backend/internal/testutil/fakeopenai`、Node `frontend/scripts/lib/fake-openai-upstream.mjs`。

## 各批次完成后的整合验证（本轮实际执行）

所有子批次合并后，在同一代码树上执行：

1. `go generate ./cmd/server`、`go build ./...`、`go vet -tags=unit ./internal/... ./cmd/...`、`go vet -tags=providerhall_localdb ./internal/repository`、`go vet -tags=perf ./internal/repository ./internal/handler`：通过；供应商大厅相关源码 `gofmt -l` 无输出。
2. `go test -tags=unit -race ./internal/service ./internal/handler/... ./internal/repository ./internal/testutil/... ./cmd/server -run 'ProviderHall|FakeOpenAI' -count=1`：service、handler、handler/admin、handler/dto、testutil/fakeopenai 全部 ok。
3. `PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin go test -tags=providerhall_localdb ./internal/repository -run ProviderHall -count=1 -v`：6 条路径（empty、upgrade_228～232）全部通过，180 个子测试。
4. `pnpm --dir frontend run lint:check`、`pnpm --dir frontend run typecheck`：退出码 0，无类型错误。
5. `pnpm --dir frontend exec vitest run`：256 个文件中 255 通过；唯一失败 `src/components/account/__tests__/OAuthAuthorizationFlow.spec.ts` 为既有基线（该文件与其导入链 `api/admin → api/client → i18n/index.ts` 自 2026-09-07 起未改动，本轮无任何文件出现在其失败堆栈中）。
6. `node frontend/scripts/lib/fake-openai-upstream.selfcheck.mjs`：退出码 0。
7. 端到端脚本与性能测试的执行结果见两份 B8 记录；本轮抽查了 1440 像素明色列表与展开截图，列结构、状态徽章、环形指标、检测徽章与"使用此分组"按钮与目标截图一致。

## 未执行与已知限制

- `CI=1 make -C backend test-integration`：本机没有 Docker，未执行；集成合约全部通过 localdb 路径补跑（见第 3 条及各批次记录）。
- 端到端运行时间短，趋势与迷你图只出现 1 个真实点和 1 个探测点，6h 端到端可用率与 TPS 只出现"样本不足"状态；多点趋势线、WS/OAuth 入口、权限拒绝场景未在浏览器中验证。
- 全页截图中固定表头在页面中部重复出现，属 Playwright 全页截图的表现，不是页面缺陷；1440 像素"列表"与"展开"两张截图内容相同（列表截图在展开行之后拍摄），已保留。
- 运营用户余额为 0 时探测请求会被网关以 `INSUFFICIENT_BALANCE` 拒绝，Runner 与 health 不检查运营用户余额，这是部署前置条件而非代码修复项。
- API Key 账号的 Responses 请求会转换为上游 chat_completions，上游在终态事件前断流会被网关按零用量成功处理，因此该路径的"断流"不会产生失败样本。
- 后端全量 `make test-unit` 有 9 个既有基线失败（与供应商大厅无关，明细见性能与回归记录）；聚焦竞态集中 `…ExpiredCacheRefreshesWithoutBlocking` 为既有失败。
- Go 1.26.6 兼容仅对聚焦集复跑；未宣称 CI 通过。

## 验收后修复：全局配置页的三个开关

验收后发现后台"全局配置"里的三个运行开关仍是早期批次的占位（始终禁用且不勾选），端到端脚本直接调用了 API 所以没有暴露。修复：`GET/PUT /admin/provider-hall/config` 响应新增 `readiness {collection, tasks, display}`（`dto.ProviderHallConfigView`），表单按就绪状态启用对应开关、提交实际勾选值，并显示每个开关的前置条件说明；未就绪时仍禁用并提示"当前部署尚未提供此功能"。已补充组件测试（就绪时恰好启用对应开关、提交值不含 `readiness`）。同时：`gateway_origin` 保存前去掉末尾单个斜杠；`PROVIDER_HALL_INVALID_CONFIG` 的错误提示按返回的 `metadata.field` 追加具体字段说明（例如"用户大厅需要先开启真实请求采集并保存"），未知字段显示字段名，admin handler/dto 单测、eslint、typecheck、admin 相关 Vitest 68 文件 401 用例通过。

## 验收后修复：开启用户大厅后普通用户看不到入口

原因：首屏公开设置由后端注入 `index.html` 并缓存，只有系统设置保存时才失效；大厅配置保存在独立表，开启 `display_enabled` 后既不刷新 30 秒的展示开关缓存，也不让注入 HTML 失效，普通用户在重启前始终拿到 `provider_hall_enabled=false`。修复：`ProviderHallService.SetConfigSavedCallback` 在配置保存成功后调用 `SettingService.NotifyProviderHallConfigChanged`，清除展示开关缓存并触发既有的设置更新回调（注入 HTML 缓存失效）；在 `ProvideProviderHallQueryService` 中接线。新增单测：保存成功触发回调、版本冲突不触发、通知后公开设置立即翻转且 HTML 缓存失效计数为 1。已打开的用户页面仍需刷新一次才会重新拉取公开设置。

## 部署包

2026-09-12 本地构建 `~/Downloads/sub2api-provider-hall-0.1.179-20260912.tar.gz`：`pnpm run build` 产物内嵌，`CGO_ENABLED=0 go build -tags embed -trimpath`，版本 0.1.179、`BuildType=release`，含 `sub2api-linux-amd64`、`sub2api-linux-arm64`、SHA256 校验与 `DEPLOY_NOTES.md`。已确认二进制内含迁移 233、`provider_hall_enabled` 与用户大厅页面。尚未部署。

## 灰度建议

沿用开发计划的顺序：应用迁移 231～233 且三个开关保持关闭 → 全部节点升级并在 `expected_nodes` 登记 `INSTANCE_ID` → 仅开启采集并核对 health 中节点确认水位 → 为少量目标绑定专用 Key 并开启任务，确认运营用户余额与预算 → 人工核对一小时快照与账单 → 开启展示 → 观察 24 小时。任一阶段可以只关闭对应开关回退，迁移无需回滚。
