# 供应商大厅 B3：聚合、快照、对账与保留实施记录

日期：2026-09-12。本批在 B2 采集事实之上实现主实例聚合器；用户 API、任务 Runner、前端大厅仍属后续批次，未部署。

## 已实现的行为

- 新增迁移 `232_provider_hall_metrics.sql`：分钟指标 `provider_hall_metrics_1m`、精确毫秒频次 `provider_hall_latency_counts_1m`、一小时窗口快照 `provider_hall_snapshots`、单行状态 `provider_hall_aggregator_state`。相对计划 DDL 增加一列 `metrics_1m.billed_subscription_count`，否则快照的 `subscription_share` 无法从分钟行推导。
- `ProviderHallAggregator` 在 `shouldStartGlobalWorkers` 实例启动，60 秒一轮，单轮 55 秒超时，领导锁 `provider-hall-aggregator`（Redis 优先、DB advisory 兜底）。每轮顺序：读配置 → 失联 epoch（心跳落后 45 秒）标 `lost` 并写 `missed_heartbeat` 缺口 → 只读账本对账 → `T = now-45s` 截到分钟，`W = state.last_window_end`（空则 `W=T`，不回填；最早不早于 `T-7d`）→ `M = [W,T) ∪ dirty`，每轮至多 240 分钟，截断后立即 kick 再跑 → 单个 ReadCommitted 事务完成重算、覆盖率复评、快照发布、脏桶删除和水位推进 → 每 10 轮做保留清理。采集关闭时仅做保留清理。
- 分钟重算：删除后按 `provider_hall_requests`（`source='user'`、`profile_id` 非空、当前算法版本、`ended_at` 落在该分钟）重建；成功/失败/排除计数、提交数、`ttft_sample_count` 只取成功且流式且有 TTFT 的样本；计费只认 `outcome=success AND billing_status=applied AND billing_mode=token`，输入费用 = `actual_cost×input_base_cost/total_base_cost`（合法零费用计入 `billed_count`）。每个已启用 target 在重算分钟都有一行（零计数占位），使无流量分钟也带覆盖状态。
- 覆盖率固定优先级 `collection_gap > node_unconfirmed > version_mismatch > complete`：collection 缺口与分钟重叠；`expected_nodes` 中任一节点的最新 epoch 缺席或 `confirmed_at < m+1m`；算法版本不一致。`expected_nodes` 为空时以当前存活（未退出）的 epoch 节点为名单。计费覆盖率 `billing_gap > pending > complete`。`[T-65m,T)` 内非本轮重算的分钟每轮复评，覆盖率任一方向变化都把依赖的快照重新发布（计划只要求"翻为完整"方向；两个方向都发布可避免名单修正后快照仍显示完整）。
- 快照集合 `S = {T} ∪ {m+k, k∈[1,60]}`，限定 `t ≤ T` 且在 tier5 保留期内；tier1 早于 48 小时的不再发布。每个已启用 target `(g,p)` 与合并行 `(g,0)` 各一条，合并行只合并该组当前已启用的 profile。每个分组只读两次（分钟行、延迟频次）后在内存内滚动 60 分钟窗口，用 `ProviderHallExactLatency/SuccessRate/CacheRate/InputPrice` 计算；`ON CONFLICT … computed_version+1`。
- 对账每轮最多 500 条：`outcome=success` 且 `billing_status ∈ {pending,uncertain}`、`started_at ∈ (now-7d, now-2m)`。有计费键：`usage_billing_dedup` 与 `usage_logs` 都有 → `applied` 并取 `actual_cost/input_cost/total_cost`；只有 dedup → `uncertain`，超过 24 小时仍无日志 → `failed`；两者皆无且超过 10 分钟 → `failed`。无计费键超过 10 分钟 → `failed`。状态变更同事务写脏桶（`reason=reconcile`）。只读账本，不重试扣费。探测/检测来源的请求行同样对账，B4 的样本与 spend 按 `requests.trace_id/billing_status` 自行补齐。
- 保留（批量 5000）：requests 7 天（排除待对账及被 `dispatched/uncertain` 样本引用者，样本表存在时才检查）；metrics/latency 7 天；snapshots tier1 48 小时、tier5 35 天；dirty 7 天；已关闭缺口 35 天；已退出 epoch 35 天；jobs/spend 90 天（表存在且无未完成样本时）。
- 指标映射 `provider_hall_snapshot.go`：状态优先级 `disabled(collection_disabled/target_disabled/group_unlisted) → incomplete(coverage；价格另含 billing_gap/billing_pending) → insufficient(snapshot_missing/samples_below_20/no_usage/bills_below_200) → stale(window_end 早于 5 分钟) → ok`，仅 ok/stale 带值；常量 `ProviderHallMinQualitySamples=20`、`ProviderHallMinBilledSamples=200`；趋势桶 6h→300s/72、24h→900s/96、7d→3600s/168、30d→21600s/120。因 `handler/dto` 已导入 `service`，service 内定义了同构的 `ProviderHallMetric[T]`，B5 在 handler 层转换为 `dto.ProviderHallMetric[T]`。
- 就绪标志：`ProvideProviderHallService` 由 `aggregator.Ready()` 推导 `Collection`，本批之后管理员可以打开 `collection_enabled`。

## 源码和迁移

- `backend/migrations/232_provider_hall_metrics.sql`
- `backend/internal/service/provider_hall_aggregator.go`（替换 stub：接口 `ProviderHallAggregationRepository`、分钟规划、覆盖率规则、循环）、`provider_hall_snapshot.go`（快照行、分钟行、指标映射、快照折叠）
- `backend/internal/repository/provider_hall_aggregation.go`（raw SQL：`GetState/RecordRun/MarkLostEpochs/ListDirty/Recompute/Reconcile/Prune`）；`repository/wire.go` 增加 `NewProviderHallAggregationRepository`
- `backend/internal/service/wire.go` 的 `ProvideProviderHallAggregator(aggRepo, facts, cfgRepo, db, lockCache, cfg)`；`cmd/server/wire_gen.go` 由 wire 重新生成
- 测试：`service/provider_hall_snapshot_test.go`、`service/provider_hall_aggregator_test.go`、`repository/provider_hall_aggregation_db_test.go`（通过 `registerProviderHallLocalDBContract("b3_aggregation", …)` 挂到共享 localdb 集群）

## 本轮验证

以下为 2026-09-12 实际执行结果（Go 1.26.5，`GOPROXY=https://goproxy.cn,direct`，工作目录 `backend/`）：

1. `gofmt -l` 本批文件：无输出。`go build ./...`：通过。`go vet -tags=unit ./internal/service ./internal/repository`、`go vet -tags=providerhall_localdb ./internal/repository`：通过。
2. `go test -tags=unit -race ./internal/service ./internal/repository -run ProviderHall -count=1`：service 包通过（含本批 11 个顶层用例：状态优先级表驱动 11 例、矩阵 A 折叠、趋势桶、分钟规划 4 例、T/W 计算与两轮推进、无锁跳过、关闭仅清理、240 分钟截断再 kick、Start/Stop、分钟覆盖率 10 例、快照时间集合）；repository 包在 unit 标签下无用例。
3. `PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin PROVIDER_HALL_PG_PORT=15482 go test -tags=providerhall_localdb ./internal/repository -run ProviderHallLocalDatabase -count=1 -v`：隔离 PostgreSQL 18 集群，六条路径（空库、从 228/229/230/231/232 升级）各执行重复迁移与全部已登记合约。本批 8 组合约 × 6 路径 = 48 组全部通过：矩阵 A 20 样本入库→重算→快照 `fast95=10000.000`、`P90=18000`、成功率、历史均价 2.0/M、订阅占比、合并行；迟到终态标脏→重算 61 个窗口且版本 +1；采集缺口使窗口 `collection_gap`，关闭后复评翻回并重发布；预期节点缺席→`node_unconfirmed`，名单修正后连同复评窗口内的旧分钟一起翻回，空名单按存活节点判定；失联 epoch→`lost`+`missed_heartbeat` 缺口且幂等；对账 applied/uncertain/failed/pending 各分支与脏桶；保留保住待对账行、删除 tier1 48h/tier5 35d/7d 分钟行；状态往返。
   同一次运行中，`fact_billing_link_before_and_after_finish`（B2/B4 事实仓库的 spend upsert，`pq: inconsistent types deduced for parameter $5`）在六条路径均失败，不属于本批文件，未处理。
4. `go generate ./cmd/server`（wire）：重新生成 `wire_gen.go`，随后 `go build ./...` 通过。

未执行：`docker info` 不可用，`CI=1 go test -tags=integration` 未跑，集成合约以 localdb 路径替代并已注明；`make -C backend generate`（ent）未执行，本批无 Ent schema 变更；前端未涉及。

## 限制与后续

- 缺失分钟行（目标启用前、聚合器启动前）视为无预期，不降低覆盖率；新目标在头一小时内因样本不足显示 `insufficient`。
- `expected_nodes` 为空时按存活 epoch 判定，管理员登记名单后才有"缺席节点阻止完整"的保护。
- 账单 24 小时仍只有 dedup 的记录会转 `failed` 以便保留清理释放，属于计划之外的收口规则。
- B4 需在 Runner 内按 `provider_hall_requests.trace_id/billing_status/actual_cost` 补齐样本与 spend；B5 读取 `provider_hall_snapshots`（`profile_id=0` 合并行）与 `ProviderHallSnapshotRow` 并转换为 DTO。
