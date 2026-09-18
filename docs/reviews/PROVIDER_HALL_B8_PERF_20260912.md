# 供应商大厅 B8：性能验收与回归记录

日期：2026-09-12。本记录只覆盖 B8 的性能（矩阵 H）与回归两部分；端到端脚本、浏览器截图与灰度演练见同日的 `PROVIDER_HALL_B8_E2E_20260912.md`。未部署生产。

环境：macOS 26.6.2 / Apple M5，Go 1.26.5（另用 `GOTOOLCHAIN=go1.26.6` 复跑聚焦集），PostgreSQL 18.6（Homebrew，`/opt/homebrew/opt/postgresql@18/bin`），Node 24.15.0 / pnpm 10.33.2。本机没有 Docker，`integration` 标签无法执行。

## 新增文件

- `backend/internal/repository/provider_hall_perf_test.go`（`//go:build perf`）：自建隔离 PostgreSQL 集群（`PROVIDER_HALL_PG_BIN`，端口 `PROVIDER_HALL_PG_PORT` 默认 15495，`shared_buffers=512MB`、`synchronous_commit=off`），执行全部迁移后造数：2 个档案、100 个上架分组、每组 2 个启用目标；tier-5 快照 30 天（每 5 分钟一个窗口）+ tier-1 快照最近 48 小时，合并行与默认档案行各一份，共 2,189,000 行；通过 `COPY` 写入 1,000,000 条 `provider_hall_requests` 用户事实，均匀分布在最近 8 天（多出的一天用于验证保留清理）。随后：
  1. 用记录语句的 pq 连接器跑一次冷列表与一次详情，对每条 `SELECT` 执行 `EXPLAIN (FORMAT JSON)`，断言语句文本与计划都不含 `usage_logs` / `provider_hall_requests`，且冷列表语句数 ≤ 12。
  2. 20 个并发用户各调一次 `ProviderHallQueryService.List`（冷基底缓存）取 P95 ≤ 2 s；同一服务预热后 20 并发 × 5 轮取 P95 ≤ 500 ms；另测 20 次详情。
  3. 把聚合水位拨回 240 分钟前，直接计时 `Reconcile`、240 分钟 `Recompute`、`Prune`，再重置水位用 `ProviderHallAggregator.RunOnce` 跑一整 tick，断言 < 55 s 且无错误。
- `backend/internal/handler/provider_hall_overhead_test.go`（`//go:build perf`）：真实 `OpenAIGatewayService` + 进程内桩上游，非流式 Responses 请求；采集器缺席 vs 启用（目标索引已加载，后台写线程排空到空实现的事实仓库）各 200 次预热 + 2000 次测量，交错跑两轮共 4000/4000 个样本；断言启用后 P95 增量 ≤ 5 ms，并核对采集器确实产生了事实行、队列无丢弃。

没有修改任何非测试源码；两个性能用例首次运行即通过，未发现需要修复的供应商大厅缺陷。

## 性能结果（矩阵 H）

命令：

```
PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin PROVIDER_HALL_PG_PORT=15495 \
  go test -tags=perf ./internal/repository ./internal/handler -run ProviderHall -count=1 -v -timeout 30m
```

结果：`internal/repository` ok（43.9 s），`internal/handler` ok（1.3 s）。

| 项目 | 实测 | 门限 |
|---|---|---|
| 造数：2,189,000 快照行 / 1,000,000 请求行 | 8.9 s / 7.3 s（COPY） | — |
| 冷列表语句数（100 组可见） | 9 条 | ≤ 12 |
| 详情语句数 | 14 条 | — |
| EXPLAIN 检查的 SELECT | 23 条，均不触及 `usage_logs`、`provider_hall_requests` | 0 条触及 |
| 列表 P95，20 并发冷缓存 | 43 ms（n=20，最大 43 ms） | ≤ 2 s |
| 列表 P95，20 并发热缓存 | 14 ms（n=100，最大 15 ms） | ≤ 500 ms |
| 详情 P95（顺序 20 组） | 1 ms | — |
| 对账（无待对账行） | < 1 ms | — |
| 240 分钟重算 + 发布 72,000 个快照（100 组 × 3 档 × 240 窗口） | 8.83 s | < 55 s |
| 保留清理（删除 125,029 条第 8 天请求、800 条过期 tier-1 快照） | 5.34 s | — |
| 聚合器整 tick（240 分钟积压，含锁、节点维护、对账、重算） | 8.81 s | < 55 s |
| 网关 P50 / P95，采集器缺席 | 25 µs / 34 µs | — |
| 网关 P50 / P95，采集器启用 | 25 µs / 32 µs | P95 增量 ≤ 5 ms |
| 启用变体采集器写入 | 6,500 行 / 13 批，队列 0/8192，丢弃 0 | 丢弃 0 |

说明：网关开销用例的上游在进程内应答，所以绝对延迟远低于真实链路；对比的是同一路径开关采集器前后的差值，实测差值为 -2 µs，落在测量噪声内。240 分钟积压是计划规定的单 tick 上限，对应线上停机 4 小时后的首个 tick；正常每分钟 tick 只重算 1 分钟加脏桶。

## 回归结果

| 命令 | 结果 |
|---|---|
| `make -C backend build` | 通过（`bin/server`，版本 0.1.179） |
| `go test -tags=unit ./...`（即 `make test-unit`） | 52 个包通过，3 个包失败共 9 个用例，全部为基线失败，见下节 |
| `go test -tags=unit -race ./... -run 'ProviderHall\|ChannelMonitorV2\|OpenAIGateway\|UsageBilling\|WS' -count=1` | 54 个包通过，`internal/service` 1 个用例失败（基线，见下节）；`ProviderHall` 相关用例全部通过 |
| 聚焦集 `-run 'ProviderHall\|FakeOpenAI'`（service、handler、admin、dto、repository、testutil/fakeopenai、cmd/server） | 79 个顶层用例、218 个子用例通过，0 失败 |
| `PROVIDER_HALL_PG_PORT=15496 go test -tags=providerhall_localdb ./internal/repository -run ProviderHall -count=1 -v` | 通过，6 条迁移路径（empty、upgrade_228…upgrade_232）× 全部批次合约，180 个子用例通过 |
| `GOTOOLCHAIN=go1.26.6`：`go build ./...` 与上述聚焦集 `-race` | 通过（9 个包 ok） |
| `make test-frontend` | 失败于第一步 `lint:check`：`frontend/scripts/provider-hall-e2e.mjs` 报 `no-inner-declarations`（该文件由 B8 端到端部分并行编写中，不属于本记录范围）；随后单独执行 `typecheck` 0 错误，`make test-frontend-critical` 13 个文件 166 个用例通过 |
| `pnpm --dir frontend exec vitest run`（全量） | 255 个文件通过、1 个文件失败（基线，见下节） |
| `CI=1 make -C backend test-integration` | 未执行：本机没有 Docker。同一批次合约已通过 `providerhall_localdb` 路径在真实 PostgreSQL 18.6 上执行 |

## 基线失败（本轮未修改、与供应商大厅无关）

以下用例的失败断言与堆栈均不涉及 `provider_hall*` 源码；`PROVIDER_HALL_FOUNDATION_20260910.md` 在实施前已记录同类失败（账户/监控复制、余额文案、Antigravity 模型映射）。

`make test-unit`（9 个）：

- `internal/handler/admin`：`TestDuplicateAccountHandlerRecoversAfterMarkSucceededFailure`、`TestDuplicateAccountHandlerPreservesIdempotencyErrorWhenRecoveryLookupFails`、`TestDuplicateChannelMonitorHandlerRecoversAfterMarkSucceededFailure`、`TestDuplicateGroupHandlerRecoversAfterMarkSucceededFailure`（幂等恢复响应头期望 `"true"` 实际为空）。
- `internal/server/middleware`：`TestAPIKeyAuthRejectsExhaustedBalance`（余额不足文案与期望英文不一致）。
- `internal/service`：`TestAntigravityGatewayService_GetMappedModel`（模型映射期望 `claude-opus-4-7`）、`TestHandle403_CNProviderStructured403TempUnschedulableFirstHit`、`TestForwardResponses_DeepSeekReasoningOnlyStreamProducesVisibleText`、`TestHandleUpstreamError_OpenAIStructured403StillPenalizes`。

聚焦 `-race`（1 个）：`TestOpenAIGatewayService_OpenAIAdvancedSchedulerRuntimeSettings_ExpiredCacheRefreshesWithoutBlocking`，竞态报告位于 `openai_account_scheduler.go:2094` 的 `resetOpenAIAdvancedSchedulerSettingCacheForTest` 与 singleflight 之间，单独复跑 3 次均稳定失败，与采集器无关。

前端（1 个文件）：`src/components/account/__tests__/OAuthAuthorizationFlow.spec.ts` 因 `vue-i18n` mock 缺少 `createI18n`，导入链 `OAuthAuthorizationFlow.vue → api/admin → api/client → i18n/index.ts` 早于本次工作存在。

## 未执行项

- Docker 集成标签（无 Docker）。
- 真实上游、多实例部署下的开销与聚合耗时；本记录的数据来自单机隔离集群与进程内桩上游。
- 端到端浏览器验收与截图（由 B8 端到端部分负责）。
