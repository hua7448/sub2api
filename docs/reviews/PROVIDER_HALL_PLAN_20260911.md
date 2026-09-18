# 供应商大厅完整实施计划（用户端 + 管理端）

## Context

目标是交付 `docs/screenshots/hubScreen/*.png` 所示的用户端"供应商大厅"（`/providers`）及配套管理端。原计划 `docs/PROVIDER_HALL_DEVELOPMENT_PLAN.md` 共 5 批，目前只完成第 1 批（配置/权限/纯算法）与第 4 批中的后台配置页：迁移 229/230、管理 API config/groups/profiles/targets、`/admin/provider-hall` 页面、V2 探测流量排除。采集、聚合、任务、检测、用户 API、用户页面全部未开始，三个运行开关被 `ErrProviderHallNotReady` 一刀切锁死。

本计划把剩余工作拆成 B2～B8 七个可独立验收的批次，每批列出文件、签名、DDL、测试。已由产品确认的决策：

- 分组的模型集合 = 该组已启用的 targets（group×profile）。行级 TTFT / P90 / 成功率**合并**该组所有已启用 profile 的真实流量（快照 `profile_id=0`）；缓存率、真实价格、预测倍率按"默认档案"（config.default_model+protocol 对应 target 若启用，否则取该组 profile_id 最小的已启用 target）；每个 target 一个健康点。
- 参考价与参考缓存率由管理员在 profile 维护（已有字段），不接 OpenRouter。
- 模型检测沿用三套断言（算术 ×3、JSON 约束 ×3、`supports_tools` 时强制工具调用 ×3）+ 响应模型名比对；判定优先级固定。
- 详情面板额外指标（6h 端到端可用率、TPS、首 Token/端到端延迟、默认档案 P90、最近探测 Token 输入/输出）全部来自探测样本；真实流量绝对计数永不返回给用户。

实施中采用的假设（不阻塞，若不符合请在对应批次调整）：
1. 采集节点 `node_id` 取 `INSTANCE_ID` 环境变量，缺省为 hostname；`expected_nodes` 必须与之一致。
2. 检测 JSON 套件容忍 ``` 围栏（记录 `fenced:true` 但不判失败）。
3. 用户选择 model/protocol 筛选时只影响"哪些分组出现"，行级价格/缓存率仍按默认档案。

## 总体架构

```
网关 handler (Responses/Chat/Messages)
  └─ ProviderHallRequestTracker(ctx) ──► ProviderHallCollector(有界队列, epoch/seq/barrier) ──► provider_hall_requests + dirty_buckets
RecordUsage(worker pool) ── applied/cost ──► Collector.RecordBilling ──► requests.billing_* / provider_hall_spend
ProviderHallAggregator(60s, primary+lock) ── requests ──► metrics_1m / latency_counts_1m ──► snapshots(1h 窗口, tier1/5)
ProviderHallRunner(primary+lock) ── jobs/samples ──► 经本站网关真实请求(专用 Key) ──► verifications / spend / probe series
用户 API /api/v1/provider-hall ── 权限过滤 → 叠加用户价格 → 排序 → 分页 ──► HallRow
前端 /providers (ProviderHallView) + /admin/provider-hall 新页签 jobs/health
```

关键复用点（均已在源码中核实）：

| 用途 | 复用 |
|---|---|
| 请求 ID | `ctxkey.ClientRequestID`（每个网关请求都有；计费去重键 = `client:`+它；响应头回显） |
| 请求开始 / ctx 挂载 | 三个 handler 中 `WithOpenAIRequestPricingContext` 之后：`openai_gateway_handler.go:509-511`、`openai_chat_completions.go:152-153`、`openai_gateway_handler.go:1082-1084` |
| 提交计数 | 包裹 `h.gatewayService.Forward*`：`openai_gateway_handler.go:620`、`openai_chat_completions.go:238`、`openai_gateway_handler.go:1177` |
| TTFT / 终态 | `OpenAIForwardResult.FirstTokenMs/Duration/Usage/UpstreamResponseModel/ClientDisconnect`；err==nil=成功；`*UpstreamFailoverError`=换号/重试；`"stream usage incomplete: missing terminal event"`=HTTP200 后流中断 |
| 计费结果 | `openai_gateway_usage.go:443-460` 中 `applyUsageBilling` 返回的 `applied bool`（当前被丢弃）、`cost *CostBreakdown{InputCost,TotalCost,ActualCost}`、`usageLog.RequestID` |
| 对账账本 | `usage_billing_dedup(request_id, api_key_id)`、`usage_logs.request_id` |
| 上游头白名单 | `openaiAllowedHeaders`（`openai_gateway_service.go:74`）、`isOpenAIPassthroughAllowedRequestHeader`；自定义头默认不透传，仍显式 `Del` |
| 探测流量排除 | `ProviderHallRepository.IsProbeKey`；V2 聚合已按 `provider_hall_probe_keys.registered_at` 排除 |
| 当前报价 | `ModelPricingResolver.Resolve` → `ResolvedPricing{Mode, BasePricing{InputPricePerToken, CacheReadPricePerToken, CacheCreationPricePerToken}, Intervals}`；`ResolveUserGroupRateMultiplier` + `Group.PeakMultiplierAt(now)` |
| 后台 worker | `service/wire.go Provide*` + `shouldStartGlobalWorkers(cfg)` + `tryAcquireSingletonLeaderLock(ctx, lockCache, db, key, owner, ttl)`；`cmd/server/wire.go provideCleanup` 注册 Stop |
| 聚合模板 | `channel_monitor_v2_aggregator.go`（timer+kick 循环、55s runOnce）、`repository/channel_monitor_v2_aggregation.go`（raw SQL、DELETE+INSERT…SELECT、watermark、保留规则） |
| 批量删除 | `ops_cleanup_executor.go:81-111 deleteOldRowsByID` |
| 用户路由 | `routes/user.go` authenticated 链 + `panelRateLimiter.Heavy()`；守卫参照 `routes/admin.go:850-870 channelMonitorModeV2Guard`；当前用户 `middleware.GetAuthSubjectFromContext` |
| 菜单开关 | `setting_public.go GetPublicSettings/GetPublicSettingsForInjection` + `dto/settings.go` + `setting_handler.go` + 前端 `utils/featureFlags.ts defineFlag/makeSidebarFlag` + `AppSidebar.vue buildSelfNavItems` |
| Key 创建/切组 | `POST /api/v1/keys`、`PUT /api/v1/keys/:id {group_id}`、`GET /groups/available`、`GET /groups/rates`；前端 `keysAPI.create/update/list` |
| 页面模板 | `views/user/ChannelStatusV2View.vue`（range chips、路由 query 同步、Abort+sequence、可展开行）、`features/channel-monitor-v2/{MetricCell,MonitorTrendChart}.vue`、`components/common/{HelpTooltip,BaseDialog,Select,Pagination,AutoRefreshButton}.vue`、`composables/useAutoRefresh.ts` |
| 测试基座 | `//go:build unit`；`providerhall_localdb`（`PROVIDER_HALL_PG_BIN`，`provider_hall_localdb_test.go` 已含 empty/upgrade_228/229，需加 upgrade_230）；httptest 模拟上游参照 `openai_gateway_bound_egress_test.go:166-177`；前端 Playwright 脚本参照 `scripts/provider-hall-admin-smoke.mjs` |

---

## B2 采集与计费关联（约 5 人日）

### 迁移 `backend/migrations/231_provider_hall_facts.sql`

```sql
CREATE TABLE provider_hall_collector_epochs (
  id bigserial PRIMARY KEY, node_id varchar(128) NOT NULL, build_version varchar(64) NOT NULL DEFAULT '',
  algorithm_version smallint NOT NULL, registered_at timestamptz NOT NULL DEFAULT now(),
  heartbeat_at timestamptz NOT NULL DEFAULT now(), exited_at timestamptz, exit_reason varchar(32) NOT NULL DEFAULT '',
  persisted_seq bigint NOT NULL DEFAULT 0, confirmed_at timestamptz, overflowed boolean NOT NULL DEFAULT false);
CREATE INDEX provider_hall_collector_epochs_node_idx ON provider_hall_collector_epochs(node_id, registered_at DESC);

CREATE TABLE provider_hall_coverage_gaps (
  id bigserial PRIMARY KEY, node_id varchar(128) NOT NULL, epoch_id bigint NOT NULL,
  scope varchar(16) NOT NULL CHECK (scope IN ('collection','billing')),
  started_at timestamptz NOT NULL, ended_at timestamptz, reason varchar(32) NOT NULL,   -- queue_overflow|db_unavailable|node_exit|missed_heartbeat
  created_at timestamptz NOT NULL DEFAULT now(), CHECK (ended_at IS NULL OR ended_at >= started_at));
CREATE INDEX provider_hall_coverage_gaps_range_idx ON provider_hall_coverage_gaps(started_at, ended_at);

CREATE TABLE provider_hall_requests (
  trace_id uuid PRIMARY KEY, node_id varchar(128) NOT NULL, epoch_id bigint NOT NULL, last_seq bigint NOT NULL,
  group_id bigint NOT NULL, profile_id bigint, api_key_id bigint NOT NULL,
  protocol varchar(32) NOT NULL CHECK (protocol IN ('responses','chat_completions','messages')),
  requested_model varchar(200) NOT NULL, upstream_model varchar(200) NOT NULL DEFAULT '', platform varchar(32) NOT NULL DEFAULT '',
  source varchar(16) NOT NULL CHECK (source IN ('user','probe','verification')), sample_id bigint,
  stream boolean NOT NULL DEFAULT false,
  started_at timestamptz NOT NULL, first_content_at timestamptz, ended_at timestamptz,
  ttft_ms integer CHECK (ttft_ms > 0), submissions integer NOT NULL DEFAULT 0 CHECK (submissions >= 0),
  outcome varchar(16) NOT NULL DEFAULT 'pending' CHECK (outcome IN ('pending','success','failed','excluded')),
  exclusion_reason varchar(32) NOT NULL DEFAULT '',   -- client_cancel|policy_reject|invalid_request|task_header_mismatch|early_exit
  usage_known boolean NOT NULL DEFAULT false,
  input_tokens bigint, output_tokens bigint, cache_read_tokens bigint, cache_creation_tokens bigint,
  billing_request_id varchar(255), billing_api_key_id bigint, request_fingerprint varchar(64),
  billing_status varchar(16) NOT NULL DEFAULT 'pending' CHECK (billing_status IN ('pending','applied','duplicate','failed','uncertain','not_applicable')),
  billing_mode varchar(16), is_subscription boolean, rate_multiplier numeric(12,6),
  input_base_cost numeric(24,10), total_base_cost numeric(24,10), actual_cost numeric(20,8), billed_at timestamptz,
  algorithm_version smallint NOT NULL DEFAULT 1,
  created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now());
CREATE INDEX provider_hall_requests_window_idx ON provider_hall_requests(group_id, profile_id, ended_at) WHERE ended_at IS NOT NULL AND source = 'user';
CREATE INDEX provider_hall_requests_reconcile_idx ON provider_hall_requests(started_at) WHERE outcome = 'success' AND billing_status IN ('pending','uncertain');
CREATE INDEX provider_hall_requests_started_idx ON provider_hall_requests(started_at);
CREATE INDEX provider_hall_requests_sample_idx ON provider_hall_requests(sample_id) WHERE sample_id IS NOT NULL;

CREATE TABLE provider_hall_dirty_buckets (
  id bigserial PRIMARY KEY, group_id bigint NOT NULL, profile_id bigint NOT NULL, minute timestamptz NOT NULL,
  reason varchar(32) NOT NULL DEFAULT '', created_at timestamptz NOT NULL DEFAULT now(), UNIQUE (group_id, profile_id, minute));
```

### 新文件

- `backend/internal/service/provider_hall_tracker.go`
  ```go
  const ProviderHallTaskHeader = "X-Provider-Hall-Task"   // "<job_id>:<sample_id>:<trace_uuid>"
  type ProviderHallRequestTracker struct { /* traceID, startedAt, group/profile/protocol/model/apiKey, source, sampleID, platform, submissions, firstContentAt, outcome, exclusion, usage, finished atomic.Bool, sink */ }
  func (c *ProviderHallCollector) Begin(ctx, in ProviderHallBeginInput) (*ProviderHallRequestTracker, context.Context) // 采集关闭或 (group,protocol,model) 不是已启用 target 且非探测 Key 时返回 (nil, ctx)；两次 map 查找，永不出错
  func ProviderHallTrackerFromContext(ctx) *ProviderHallRequestTracker
  func (t *ProviderHallRequestTracker) TraceID() uuid.UUID
  func (t *ProviderHallRequestTracker) Submit(account *Account)      // 每次 Forward* 前：submissions++，首次记 platform
  func (t *ProviderHallRequestTracker) Attempt(forwardStart time.Time, res *OpenAIForwardResult, err error)
  //   res.FirstTokenMs!=nil 且首次 → firstContentAt = forwardStart+FirstTokenMs（TTFT = firstContentAt-startedAt，含本机鉴权/排队/重试）
  //   err==nil → success + usage（只记一次）；res.ClientDisconnect → client_cancel；其余 err 记为最后失败
  func (t *ProviderHallRequestTracker) Reject(reason)                 // Begin 之后的策略/利润/cyber 拒绝
  func (t *ProviderHallRequestTracker) NoAccount()                    // 账号池不足：F=1, A=0
  func (t *ProviderHallRequestTracker) Finish()                       // 幂等；Begin 后立即 defer
  ```
  所有方法对 nil receiver 安全。`Finish` 判定：success → `success`；否则 client_cancel → `excluded`；否则 Reject → `excluded/<reason>`；否则有失败尝试或 NoAccount → `failed`；否则 `excluded/early_exit`。流中断错误带部分 result → `failed`。
  任务头：有 `X-Provider-Hall-Task` 且 Key 是探测 Key（registered_at ≤ startedAt）→ `source=probe|verification, sample_id`；有头但 Key 不匹配 → `source=user, exclusion=task_header_mismatch`；探测 Key 无头 → `source=probe, sample_id=nil`。
- `backend/internal/service/provider_hall_collector.go`
  ```go
  type ProviderHallCollector struct { repo ProviderHallFactRepository; cfgRepo ProviderHallRepository; nodeID string; epochID int64; queue chan providerHallEvent /* cap 8192 */; seq atomic.Uint64; overflowed atomic.Bool; dropped atomic.Uint64; targets atomic.Pointer[providerHallTargetIndex] }
  func (c *ProviderHallCollector) Start()/Stop()          // 在所有处理网关流量的实例启动（不受 shouldStartGlobalWorkers 限制）
  func (c *ProviderHallCollector) RecordBilling(ProviderHallBillingEvent)   // 实现 ProviderHallBillingSink
  ```
  - 事件 `start / finish / billing / barrier`，每个带 `seq`；入队非阻塞，队满则 `dropped++`、`overflowed=true`、事件丢弃。
  - 写线程：每 250ms 或 500 条一批，单事务：`UpsertRequests`（`ON CONFLICT (trace_id) DO UPDATE` 各列 `COALESCE`，先到后到都能合并）→ `MarkDirty(group, profile, minute(ended_at))`（同事务）→ 探测计费写 `provider_hall_spend`（B4 表存在时）→ 批内含 barrier 且未溢出时 `UPDATE collector_epochs SET persisted_seq, confirmed_at=barrier.at`。
  - barrier 每 1s 入队（`at = now-1s`）。
  - 溢出：首次 flush 时 `OpenGap(scope=collection, started_at=上次 confirmed_at, reason=queue_overflow)`，停止推进水位；队列低于 50% 且新 barrier 落库后 `CloseGap`，从该 barrier 恢复。DB 写失败退避重试保留批次，期间溢出走同一路径（reason `db_unavailable`）。
  - 30s 刷新 config + targets + probe keys 快照；10s 心跳。
  - `Stop()`：5s 内 flush，`exited_at=now(), exit_reason=shutdown`；未 flush 事件写 gap `node_exit`。
  - 崩溃/失联由聚合器判定：`heartbeat_at < now-45s` 且未退出 → 写 gap `missed_heartbeat`、标 `exit_reason=lost`。
- `backend/internal/repository/provider_hall_fact_repo.go`（raw `*sql.DB`）：
  ```go
  type ProviderHallFactRepository interface {
    RegisterEpoch(ctx, node, build string, algo int) (int64, error); Heartbeat(ctx, epochID) error; MarkEpochExited(ctx, epochID, reason) error
    WriteBatch(ctx, b ProviderHallFactBatch) error       // upserts + dirty + spend + watermark，一个事务
    OpenGap(ctx, g ProviderHallGap) (int64, error); CloseGap(ctx, id, endedAt) error
    ListEnabledTargets(ctx) ([]ProviderHallTargetRef, error); ListProbeKeys(ctx) (map[int64]time.Time, error) }
  ```
- `backend/internal/handler/provider_hall_gateway_hooks.go`：`providerHallBegin(c, apiKey, protocol, reqModel, stream, requestStart)`、`providerHallAttempt(...)` 小助手，三个 handler 共用。

### 修改点

- 三个 handler：pricing ctx 之后 `tracker, ctx := h.providerHall.Begin(...); c.Request = c.Request.WithContext(ctx); defer tracker.Finish()`；无账号分支 `tracker.NoAccount()`；Forward 前 `Submit(account)`、后 `Attempt(forwardStart, result, err)`；cyber/利润拒绝分支 `Reject(policy_reject)`；`submitResponsesUsage` 闭包外取 `traceID := tracker.TraceID()` 并写入 `OpenAIRecordUsageInput.ProviderHallTraceID`（不扩展 `usageRecordContext`）。Chat/Messages 镜像。WS 入口不改。
- `service/openai_gateway_usage.go`：
  ```go
  type OpenAIRecordUsageInput struct { ...; ProviderHallTraceID uuid.UUID }
  type ProviderHallBillingEvent struct { TraceID uuid.UUID; BillingRequestID string; APIKeyID int64; Fingerprint string; Status string /* applied|duplicate|failed|not_applicable */; Mode string; IsSubscription bool; Multiplier float64; InputBaseCost, TotalBaseCost, ActualCost decimal.Decimal; BilledAt time.Time }
  type ProviderHallBillingSink interface { RecordBilling(ProviderHallBillingEvent) }
  ```
  443-460 处改为 `applied, err := applyUsageBilling(...)`，两个分支之后都调 `s.providerHallCollector.RecordBilling(...)`：`applied→applied`、`err!=nil→failed`、`!applied→duplicate`；`Mode != token`（图片/按次/视频）→ `not_applicable`。原有错误返回行为不变。
- `openai_gateway_service.go`：字段 `providerHallCollector ProviderHallBillingSink` + `SetProviderHallBillingSink`。
- 上游请求构造器加 `req.Header.Del(ProviderHallTaskHeader)`：`openai_gateway_forward.go`（~1140 白名单复制之后）、`openai_gateway_passthrough.go`（~528）、`openai_gateway_messages_anthropic_native.go buildNativeAnthropicUpstreamRequest`、chat-completions 各 builder。
- `service/provider_hall.go`：引入就绪标志（见"开关解锁"节）。
- `repository/wire.go`、`service/wire.go`（`ProvideProviderHallCollector(factRepo, cfgRepo, cfg)`，所有角色启动）、`handler/wire.go`、`cmd/server/wire.go provideCleanup`（Collector flush 步骤）。

### 测试（矩阵 A 补齐、B、C 采集侧、D 部分）

- `provider_hall_algorithm_test.go`：补齐 A 的固定样例（1000..19000+100000 → fast95=10000、P90=18000；N=0/1/10/20；重复边界值；800/1100；2/3；$1/M）。
- `provider_hall_tracker_test.go`、`provider_hall_gateway_hooks_test.go`（unit，httptest 上游）：三入口 × 非流式/SSE/透传/转换；同账号重试、换账号、429、超时、无账号、HTTP200 流中断、正常终态、客户端取消、纯心跳；断言 submissions / 唯一终态 / 首内容时间；混合分组 `group_id=入口组`、`platform=account.Platform`。
- `provider_hall_collector_test.go`（unit+race）：billing 先于 finish / finish 先于 billing / 重复账单 / 扣费失败不算免费 / 非 token 模式 not_applicable；队列溢出开缺口 / barrier 推进水位 / Stop flush 并标退出。
- `TestProviderHallTaskHeaderNeverForwarded`：三协议上游都收不到该头。
- `provider_hall_fact_repo_test.go`（integration）+ `provider_hall_localdb_test.go` 增 `upgrade_230`：迁移 231、唯一键、金额精度、upsert 幂等。

---

## B3 聚合、快照、对账、保留（约 5 人日）

### 迁移 `232_provider_hall_metrics.sql`

```sql
CREATE TABLE provider_hall_metrics_1m (
  group_id bigint NOT NULL, profile_id bigint NOT NULL, minute timestamptz NOT NULL, algorithm_version smallint NOT NULL,
  success_count int NOT NULL DEFAULT 0, failed_count int NOT NULL DEFAULT 0, submissions int NOT NULL DEFAULT 0, excluded_count int NOT NULL DEFAULT 0,
  ttft_sample_count int NOT NULL DEFAULT 0, usage_success_count int NOT NULL DEFAULT 0,
  input_tokens bigint NOT NULL DEFAULT 0, cache_read_tokens bigint NOT NULL DEFAULT 0, cache_creation_tokens bigint NOT NULL DEFAULT 0, output_tokens bigint NOT NULL DEFAULT 0,
  billed_count int NOT NULL DEFAULT 0, billed_input_tokens bigint NOT NULL DEFAULT 0, billed_input_cost numeric(24,10) NOT NULL DEFAULT 0,
  billing_pending_count int NOT NULL DEFAULT 0, billing_uncertain_count int NOT NULL DEFAULT 0,
  coverage varchar(16) NOT NULL DEFAULT 'complete' CHECK (coverage IN ('complete','collection_gap','node_unconfirmed','version_mismatch')),
  billing_coverage varchar(16) NOT NULL DEFAULT 'complete' CHECK (billing_coverage IN ('complete','billing_gap','pending')),
  computed_at timestamptz NOT NULL DEFAULT now(), PRIMARY KEY (group_id, profile_id, minute, algorithm_version));
CREATE INDEX provider_hall_metrics_1m_minute_idx ON provider_hall_metrics_1m(minute);

CREATE TABLE provider_hall_latency_counts_1m (
  group_id bigint NOT NULL, profile_id bigint NOT NULL, minute timestamptz NOT NULL, algorithm_version smallint NOT NULL,
  ttft_ms int NOT NULL CHECK (ttft_ms > 0), count int NOT NULL CHECK (count > 0),
  PRIMARY KEY (group_id, profile_id, minute, algorithm_version, ttft_ms));
CREATE INDEX provider_hall_latency_counts_1m_minute_idx ON provider_hall_latency_counts_1m(minute);

CREATE TABLE provider_hall_snapshots (
  id bigserial PRIMARY KEY, group_id bigint NOT NULL, profile_id bigint NOT NULL,   -- profile_id=0 为该组合并
  window_end timestamptz NOT NULL,           -- T，窗口 [T-60m, T)
  tier smallint NOT NULL CHECK (tier IN (1,5)),  -- minute(T)%5==0 时为 5
  algorithm_version smallint NOT NULL,
  ttft_sample_count int NOT NULL, ttft_fast95_mean_ms numeric(12,3), ttft_p90_ms int,
  success_count int NOT NULL, failed_count int NOT NULL, submissions int NOT NULL, success_rate numeric(11,10),
  usage_success_count int NOT NULL, input_tokens bigint NOT NULL, cache_read_tokens bigint NOT NULL, cache_creation_tokens bigint NOT NULL, cache_rate numeric(11,10),
  billed_count int NOT NULL, billed_input_tokens bigint NOT NULL, billed_input_cost numeric(24,10) NOT NULL, historical_price numeric(24,10),
  subscription_share numeric(11,10),
  coverage varchar(16) NOT NULL, billing_coverage varchar(16) NOT NULL, coverage_reason varchar(32) NOT NULL DEFAULT '',
  computed_at timestamptz NOT NULL DEFAULT now(), computed_version int NOT NULL DEFAULT 1,
  UNIQUE (group_id, profile_id, window_end, algorithm_version));
CREATE INDEX provider_hall_snapshots_window_idx ON provider_hall_snapshots(window_end, tier);

CREATE TABLE provider_hall_aggregator_state (
  id smallint PRIMARY KEY DEFAULT 1 CHECK (id = 1), algorithm_version smallint NOT NULL DEFAULT 1,
  last_window_end timestamptz, last_run_at timestamptz, last_error varchar(2000) NOT NULL DEFAULT '', updated_at timestamptz NOT NULL DEFAULT now());
INSERT INTO provider_hall_aggregator_state (id) VALUES (1);
```

### 新文件

- `backend/internal/service/provider_hall_aggregator.go`：仿 V2 aggregator（Start/Stop/kick、60s、runOnce 55s、`tryAcquireSingletonLeaderLock("provider-hall-aggregator", 2m)`）。`runOnce`：
  1. `collection_enabled=false` → 只做保留清理（每 10 tick）。
  2. `T = now-45s 按分钟截断`；`W = state.last_window_end`，为空则 `W=T`（不回填历史）。
  3. 节点维护：失联 epoch 写 gap；对 `expected_nodes` 逐个取最新 epoch 的 `confirmed_at` 与 `algorithm_version`。
  4. 待算分钟 `M = [W,T) ∪ dirty.minute`（读 `id <= maxDirtyID`），每 tick 最多 240 分钟，截断则再 kick。
  5. 单 ReadCommitted 事务 `RecomputeMinutes(ctx, M, nodeStatus)`：DELETE 两表 → `INSERT … SELECT … FROM provider_hall_requests WHERE source='user' AND profile_id IS NOT NULL AND ended_at ∈ 分钟 GROUP BY`（`count FILTER`；billed 只算 `outcome=success AND billing_status=applied AND billing_mode=token`；输入费用 = `actual_cost*input_base_cost/NULLIF(total_base_cost,0)`；合法零费用计入 billed_count）→ latency_counts 按 `ttft_ms` 分组 → 每分钟 coverage（有 collection gap 重叠 → `collection_gap`；预期节点 `confirmed_at < m+1m` 或缺席 → `node_unconfirmed`；版本不一致 → `version_mismatch`）与 billing_coverage（`billing_gap` / 有 pending|uncertain → `pending`）。另对 `[T-65m,T)` 只重评 coverage 列，由非完整翻为完整时把依赖的 60 个快照加入 S。
  6. 快照集合 `S = {T} ∪ {m+k, k∈[1,60], m∈M, t≤T, 在保留期内}`；对每个 t、每个已启用 target `(g,p)` 及 `(g,0)`：一条 SQL 汇总 metrics_1m `[t-60m,t)`，一条汇总 latency_counts（profile 0 按 group 且限定该组已启用 profile），Go 内用 `ProviderHallExactLatency/SuccessRate/CacheRate/InputPrice` 计算，`INSERT … ON CONFLICT DO UPDATE … computed_version+1`。
  7. 同事务删 dirty（`id<=maxDirtyID`）、更新 state；提交即原子发布。
  8. 对账（每 tick，独立短事务，先于 5）：`outcome=success AND billing_status IN (pending,uncertain) AND started_at ∈ (now-7d, now-2m)` 取 500 条；有计费键 → 查 `usage_billing_dedup` 与 `usage_logs`：两者都有 → `applied` 并从 usage_logs 取 `actual_cost/input_cost/total_cost`；只有 dedup → 维持 `uncertain`；10 分钟后仍无 → `failed`；无计费键（RecordUsage 未运行）10 分钟后 → `failed`。每次状态变更标脏。**只读账本，不重试扣费。** 同一流程按 `client_request_id` 补齐 B4 探测样本与 spend。
  9. 保留（每 10 tick，`deleteOldRowsByID` 分批 5000）：requests 7d（排除 `billing_status IN (pending,uncertain)` 与被未完成任务引用者）；metrics_1m/latency_counts 7d；snapshots tier1 48h、tier5 35d；dirty 7d；已关闭 gaps 35d；已退出 epochs 35d；jobs/samples/verifications/spend 90d（排除 dispatched/uncertain 样本所属任务）。
- `backend/internal/service/provider_hall_snapshot.go`：快照行 → `dto.ProviderHallMetric[T]`：
  ```go
  func ProviderHallTTFTMetric(in) (fast95 Metric[float64], p90 Metric[int64]); ProviderHallSuccessMetric; ProviderHallCacheMetric; ProviderHallPriceMetric(in) Metric[string]
  ```
  优先级：`disabled`（collection_disabled/target_disabled/group_unlisted）→ `incomplete`（collection_gap/node_unconfirmed/version_mismatch；价格另加 billing_gap/billing_pending）→ `insufficient`（`samples_below_20`：TTFT 用 ttft_sample_count、成功率用 success+failed、缓存用 usage_success_count；`no_usage` 分母 0；`bills_below_200`）→ `stale`（`now - window_end > 5m`；快照缺失 → `insufficient/snapshot_missing`）→ `ok`。仅 ok/stale 带 value。常量 `ProviderHallMinQualitySamples=20`、`ProviderHallMinBilledSamples=200`。
- `backend/internal/repository/provider_hall_aggregation.go`：`ListDirty`、`RecomputeMinutes`、`PublishSnapshots`、`Reconcile`、`Prune`、`MarkLostEpochs`。
- 趋势桶：6h→300s(72 点)、24h→900s(96)、7d→3600s(168)、30d→21600s(120)。真实曲线取 tier5 快照 `window_end = ANY(bucketEnds)`，缺失或非完整 → `null + gap:true`。探测曲线（B4）：`date_bin($bucket, received_at, $rangeStart)` 对 `provider_hall_samples(kind=probe, status=received)` 求 `avg(total_ms)`、`avg(ttft_ms)`、失败数。

### 修改点

- `service/wire.go`：`ProvideProviderHallAggregator(aggRepo, cfgRepo, db, lockCache, cfg)`，`shouldStartGlobalWorkers` 门控；`cmd/server/wire.go` cleanup。
- `service/provider_hall.go`：就绪标志 `Collection=true`（B2+B3 完成后在 Provide 中设置）。

### 测试（矩阵 A 端到端、D）

- `provider_hall_snapshot_test.go`：状态优先级表驱动（缺口 > 样本不足 > 过期）。
- `provider_hall_aggregator_test.go`（unit）：T/W/M 计算；锁未获得跳过；240 分钟截断再 kick。
- `provider_hall_aggregation_test.go`（integration/localdb）：迁移 231/232 空库、228/229/230 升级、重复启动；迟到终态标脏并重算 60 个快照；溢出缺口使窗口 incomplete；预期节点缺席阻止 complete；对账 applied/uncertain/failed；保留不删待对账；乱序事件；A 样例入库 → 重算 → 断言快照值。

---

## B4 任务 Runner、探测、模型检测、预算、对账（约 7 人日）

### 迁移 `233_provider_hall_jobs.sql`

```sql
CREATE TABLE provider_hall_jobs (
  id bigserial PRIMARY KEY, kind varchar(16) NOT NULL CHECK (kind IN ('probe','verification')),
  target_id bigint NOT NULL,   -- 无 FK：历史在目标/分组删除后保留
  group_id bigint NOT NULL, profile_id bigint NOT NULL,
  config_snapshot jsonb NOT NULL,   -- {profile:{id,version,model,protocol,supports_tools,output_limit,model_aliases}, target_version, probe_key_id, operator_user_id, gateway_origin}
  slot_at timestamptz, not_before timestamptz, idempotency_key varchar(128), requested_by bigint,
  status varchar(16) NOT NULL DEFAULT 'queued' CHECK (status IN ('queued','running','succeeded','failed','cancelled','unknown')),
  lease_owner varchar(128), lease_until timestamptz, attempts int NOT NULL DEFAULT 0, budget_day date,
  error_code varchar(64) NOT NULL DEFAULT '', error_message varchar(2000) NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now(), started_at timestamptz, finished_at timestamptz, updated_at timestamptz NOT NULL DEFAULT now());
CREATE UNIQUE INDEX provider_hall_jobs_slot_uq ON provider_hall_jobs(target_id, kind, slot_at) WHERE slot_at IS NOT NULL;
CREATE UNIQUE INDEX provider_hall_jobs_idempotency_uq ON provider_hall_jobs(idempotency_key) WHERE idempotency_key IS NOT NULL;
CREATE INDEX provider_hall_jobs_claim_idx ON provider_hall_jobs(status, not_before, created_at) WHERE status IN ('queued','running');
CREATE INDEX provider_hall_jobs_target_idx ON provider_hall_jobs(target_id, kind, created_at DESC);
CREATE INDEX provider_hall_jobs_group_idx ON provider_hall_jobs(group_id, created_at DESC);

CREATE TABLE provider_hall_samples (
  id bigserial PRIMARY KEY, job_id bigint NOT NULL REFERENCES provider_hall_jobs(id) ON DELETE CASCADE,
  test_id varchar(32) NOT NULL,   -- probe|arith|json|tool
  seq smallint NOT NULL, trace_id uuid NOT NULL,
  status varchar(16) NOT NULL DEFAULT 'prepared' CHECK (status IN ('prepared','dispatched','received','uncertain')),
  prepared_at timestamptz NOT NULL DEFAULT now(), dispatched_at timestamptz, received_at timestamptz,
  client_request_id varchar(64), http_status int, ttft_ms int, total_ms int, generation_ms int, input_tokens int, output_tokens int,
  response_model varchar(200), result varchar(16) CHECK (result IN ('passed','failed','error','model_mismatch')),
  error_code varchar(64) NOT NULL DEFAULT '', detail jsonb NOT NULL DEFAULT '{}',   -- expected/actual 截断 2KB，不含用户内容
  billing_status varchar(16) NOT NULL DEFAULT 'pending' CHECK (billing_status IN ('pending','confirmed','uncertain','failed','unbilled')),
  actual_cost numeric(20,8), UNIQUE (job_id, test_id, seq), UNIQUE (trace_id));
CREATE INDEX provider_hall_samples_inflight_idx ON provider_hall_samples(status, dispatched_at) WHERE status IN ('dispatched','uncertain');
CREATE INDEX provider_hall_samples_received_idx ON provider_hall_samples(received_at) WHERE status = 'received';

CREATE TABLE provider_hall_verifications (
  job_id bigint PRIMARY KEY REFERENCES provider_hall_jobs(id) ON DELETE CASCADE,
  group_id bigint NOT NULL, profile_id bigint NOT NULL, target_id bigint NOT NULL,
  verdict varchar(16) NOT NULL CHECK (verdict IN ('passed','failed','suspected','insufficient')),
  execution_status varchar(16) NOT NULL CHECK (execution_status IN ('completed','partial','error')),
  reason_code varchar(64) NOT NULL DEFAULT '',
  summary jsonb NOT NULL,   -- {arithmetic:{passed,failed,error}, json:{...}, tool:{...}|null, model:{matched,mismatched,missing,seen:[]}}
  profile_version bigint NOT NULL, target_version bigint NOT NULL,
  completed_at timestamptz NOT NULL, expires_at timestamptz NOT NULL, stale boolean NOT NULL DEFAULT false, created_at timestamptz NOT NULL DEFAULT now());
CREATE INDEX provider_hall_verifications_lookup_idx ON provider_hall_verifications(group_id, profile_id, completed_at DESC);

CREATE TABLE provider_hall_spend (
  billing_request_id varchar(255) NOT NULL, api_key_id bigint NOT NULL, job_id bigint, sample_id bigint,
  budget_day date NOT NULL, actual_cost numeric(20,8) NOT NULL DEFAULT 0,
  status varchar(16) NOT NULL CHECK (status IN ('confirmed','uncertain','failed')), confirmed_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(), PRIMARY KEY (billing_request_id, api_key_id));
CREATE INDEX provider_hall_spend_day_idx ON provider_hall_spend(budget_day, status);
```

### 新文件

- `backend/internal/service/provider_hall_runner.go`（primary，15s tick，Stop 等待在途派发 ≤90s）
  1. **排程**（leader lock `provider-hall-scheduler` 1m）：`tasks_enabled` 时对每个 enabled target `slot = now.Truncate(interval)`，`INSERT … ON CONFLICT DO NOTHING`，`not_before = slot + jitter`（探测 ±30s、检测 0～5m），只排当前时隙不补跑；目标停用/下架/删除 → queued 任务 `cancelled/target_disabled`；租约过期的 running：样本全 prepared → 回 `queued`，有 dispatched/uncertain → `unknown`（只核对不重发）。
  2. **领取**（DB 租约，`FOR UPDATE SKIP LOCKED`，同 target 无 running）：
     ```sql
     UPDATE provider_hall_jobs SET status='running', lease_owner=$1, lease_until=now()+interval '120 seconds', started_at=COALESCE(started_at,now()), attempts=attempts+1
     WHERE id = (SELECT id FROM provider_hall_jobs WHERE status='queued' AND (not_before IS NULL OR not_before<=now())
                 AND NOT EXISTS (SELECT 1 FROM provider_hall_jobs r WHERE r.target_id=provider_hall_jobs.target_id AND r.status='running')
                 ORDER BY not_before NULLS FIRST, created_at LIMIT 1 FOR UPDATE SKIP LOCKED) RETURNING *;
     ```
     租约每 30s 续；首次领取时按套件创建 `prepared` 样本（`ON CONFLICT (job_id,test_id,seq) DO NOTHING`）。
  3. **发出前检查（固定顺序，事务内 `pg_advisory_xact_lock(hashtext('provider_hall_dispatch'))`）**：tasks_enabled 且目标启用、分组上架、`config_snapshot` 中 target/profile 版本未变（否则 `cancelled/config_changed`）→ 用最新行重跑 `ValidateProviderHallTargetBinding`（否则 `failed/probe_key_invalid`，不发付费请求）→ 账单积压：任一样本 dispatched/received 且 `billing_status IN (pending,uncertain)` 超过 120s → 任务回 queued、`not_before=now+60s`、`billing_backlog` → 预算：`Σ spend(budget_day=今天, confirmed) ≥ daily_budget`（预算>0 时）→ `cancelled/budget_exhausted`；`budget_day` 在首次派发时定、后续样本继承（跨零点归发出日；改预算不清零）→ 并发：全局 `dispatched` 样本 <2 且本 target 无在途，否则本 tick 等待。通过后 `samples SET status='dispatched', dispatched_at=now()` **先提交再发送**。
  4. **发送**：`http.Client{Timeout:90s}`，`CheckRedirect` 拒绝；目标 = `gateway_origin + path`，仅 `PROVIDER_HALL_ALLOW_LOOPBACK=1` 时允许环回（`UpdateConfig` 同步校验）；头 `Authorization: Bearer <专用 Key>`、`Accept: text/event-stream`、`X-Provider-Hall-Task: <job>:<sample>:<trace>`、`X-Client-Request-ID: <trace>`、`User-Agent: sub2api-provider-hall/<version>`。结果：2xx 且收到终态 → `received`；明确非 2xx → `received/error(http_<status>)`；连接拒绝/DNS → `received/error(transport_refused)`；超时/头后 EOF/可能已准入的错误 → `uncertain`（永不重发，由对账按 `requests.trace_id` 补齐；10 分钟无结果 → `received/error/uncertain_timeout`）。记录响应头 `X-Client-Request-ID` 用于 spend 对账。
  5. **收尾**：探测 → 样本 passed 则 `succeeded` 否则 `failed`，并 upsert `provider_hall_probe_health`（放 samples 表的派生视图即可：最新 received 样本）；检测 → 计算判定，写 `verifications`（`expires_at = completed_at+48h`）。`SaveProfile/SaveTargets` 改动 model/protocol/supports_tools/aliases 时在既有事务中把受影响报告 `stale=true`。
- `backend/internal/service/provider_hall_probe_client.go`：三协议请求体与 SSE 解析
  - Responses：`{"model":M,"input":[{"role":"user","content":[{"type":"input_text","text":P}]}],"stream":true,"store":false,"max_output_tokens":N}`；TTFT = 首个非空 `response.output_text.delta`（或 `function_call` 的 `output_item.added`）；模型与 usage 来自 `response.completed`。
  - Chat：`{"model":M,"messages":[...],"stream":true,"max_tokens":N,"stream_options":{"include_usage":true}}`；TTFT = 首个非空 `delta.content` 或 `delta.tool_calls`；终态 `[DONE]`。
  - Messages：`{"model":M,"max_tokens":N,"stream":true,"messages":[...]}`；模型与输入 usage 来自 `message_start`，TTFT = 首个 `content_block_delta`，输出 usage 来自 `message_delta`，终态 `message_stop`。
  - N：探测 256、检测 1024；探测 prompt `"Reply with exactly the word OK."`；`generation_ms = total_ms - ttft_ms`，`TPS = output_tokens/generation_ms*1000`。
- `backend/internal/service/provider_hall_verification.go`
  ```go
  func ProviderHallBuildSuite(profile ProviderHallProfile, rng *rand.Rand) []ProviderHallTestCase
  func ProviderHallVerdict(samples []ProviderHallSampleResult, profile ProviderHallProfile) (verdict, reason string, summary ProviderHallVerificationSummary)
  ```
  - `arith`×3：`a,b∈[100,999]`，op∈{+,-,*}，prompt "Compute {a} {op} {b}. Reply with only the integer result and nothing else."；断言去空白/千分位后整数相等。
  - `json`×3：`token`=8 位 `[a-z0-9]`，prompt 要求只返回 `{"token":"…","length":N}`；断言剥离可选围栏后解析，键集合精确为 {token,length} 且值相等。
  - `tool`×3（仅 supports_tools）：工具 `hall_echo`，schema `{"type":"object","properties":{"code":{"type":"string"}},"required":["code"],"additionalProperties":false}`；Responses `tools:[{"type":"function","name":"hall_echo","parameters":S}], tool_choice:{"type":"function","name":"hall_echo"}`；Chat `tools:[{"type":"function","function":{...}}], tool_choice:{"type":"function","function":{"name":"hall_echo"}}`；Messages `tools:[{"name":"hall_echo","input_schema":S}], tool_choice:{"type":"tool","name":"hall_echo"}`；断言出现 `hall_echo` 调用且参数为 `{"code": c}`。
  - 模型比对：`response_model ∈ {profile.Model} ∪ model_aliases` → matched；空 → missing；否则 mismatched。
  - 判定：`mismatched≥2 → suspected`；否则任一必需套件 `failed≥2 → failed`；否则有样本缺失/error/模型名 missing → `insufficient`；否则 `passed`。`execution_status`：completed/partial/error。超过 48h 或 stale → 展示为 expired。
- `backend/internal/service/provider_hall_budget.go`：`RecordSpend`（由 Collector 的 billing 事件在 `source∈{probe,verification}` 时 upsert：applied→confirmed、failed→failed、duplicate→no-op；样本 `billing_status` 镜像）、`SumSpend(day)`。
- `backend/internal/repository/provider_hall_job_repo.go`：`EnqueueSlot`、`EnqueueManual`、`Claim`、`Heartbeat`、`Complete`、`Cancel`、`CreateSamples`、`MarkSample*`、`ListInflight`、`SumSpend`、`ListJobs(filter,page)`、`GetJob`、`LatestProbePerTarget`、`ProbeSeries`、`LatestVerification`、`ListVerifications`。

### 修改点

- `service/provider_hall.go`：就绪标志 `Tasks=true`；`UpdateConfig` 中 `tasks_enabled` 要求 `collection_enabled`、`gateway_origin`、`operator_user_id` 齐备。
- `service/provider_hall_target.go`：目标停用/版本变更 → 取消 queued 任务。
- `service/wire.go`：`ProvideProviderHallRunner(jobRepo, cfgRepo, apiKeyService, userRepo, groupRepo, db, lockCache, cfg)`；`cmd/server/wire.go` cleanup（Runner 最先停）。

### 测试（矩阵 E）

- `provider_hall_runner_test.go`（integration/localdb + race）：两 primary 同时领取只发一次；prepared 崩溃可重领；dispatched 崩溃不重发；幂等键复用；停用取消未发样本；并发 2；预算达限停发；超额在途仍记账；账单未知 120s 暂停；上海零点归属；Key 失效不发送。
- `provider_hall_verification_test.go`：判定优先级表驱动（2 次不符+3 次失败→疑似；同一套件 2 次失败→失败；1 失败+1 异常→证据不足；全过→通过；未声明工具则无 tool 套件）；别名匹配；过期规则。
- `provider_hall_probe_client_test.go`：三协议 SSE fixture 解析；拒绝重定向；头已发送且上游收不到。

---

## B5 用户 API 与定价（约 4 人日）

### 新文件

- `backend/internal/handler/dto/provider_hall_user.go`（字段名即 JSON）：
  ```go
  type ProviderHallQuote struct { InputPrice, CachePrice string; Unit string /* usd_per_million|quota_per_million */; Applicable bool; ReasonCode string }
  type ProviderHallModelHealth struct { ProfileID int64; Model, Protocol, Status string /* up|down|unknown */; CheckedAt *time.Time }
  type ProviderHallHealth struct { Status string; CheckedAt *time.Time; Models []ProviderHallModelHealth }
  type ProviderHallRowMetrics struct { TTFTFast95Ms Metric[float64]; TTFTP90Ms Metric[int64]; CacheRate Metric[float64]; SuccessRate Metric[float64] }
  type ProviderHallSparkPoint struct { T time.Time; RealTTFTMs *float64; ProbeMs *float64; Gap bool }
  type ProviderHallSparkline struct { Range string; BucketSeconds int; Points []ProviderHallSparkPoint }
  type ProviderHallVerificationBadge struct { Verdict *string; ReasonCode string; CompletedAt *time.Time; Expired bool; ReportID *int64 }
  type ProviderHallRow struct { GroupID int64; Name, Description string; DisplayOrder int; Rate string; Quote ProviderHallQuote; HistoricalPrice, PredictedRate Metric[string]; Health ProviderHallHealth; Metrics ProviderHallRowMetrics; Sparkline ProviderHallSparkline; Verification ProviderHallVerificationBadge; DefaultProfile ProviderHallProfileRef }
  type ProviderHallListResponse struct { SnapshotID string; DataThrough *time.Time; MetricVersion int; PricingAt time.Time; Range string; Catalog{Models []ProfileRef; Default ProfileRef}; Summary{Available, Listed, Abnormal, Verified int}; Items []ProviderHallRow; Pagination }
  type ProviderHallDetailResponse struct { ProviderHallRow; Detail{E2EAvailability6h Metric[float64]; TPS Metric[float64]; ProbeTTFTMs, ProbeTotalMs, P90Ms Metric[int64]; ProbeTokens *{Input, Output int; At time.Time}; LastProbeAt *time.Time}; Trend{Range string; BucketSeconds int; Points []{T; RealFast95Ms, ProbeTotalMs, ProbeTTFTMs *float64; Gap bool; ProbeFailed int}}; Profiles []{ProfileRef; Health; Verification} }
  type ProviderHallVerificationReport struct { ReportID int64; ProfileID int64; Model, Protocol, Verdict, ExecutionStatus, ReasonCode string; CompletedAt, ExpiresAt time.Time; Expired bool; Summary json.RawMessage }
  ```
  不含账号 ID、凭据、上游地址、请求体、绝对流量计数（样本量只通过 state/reason_code 表达）。
- `backend/internal/service/provider_hall_query.go`（`ProviderHallQueryService`）、`provider_hall_pricing.go`（`ProviderHallPredictedRate(m,P,C,h,Pr,Cr,hr)`）、`provider_hall_sort.go`；`backend/internal/repository/provider_hall_read_repo.go`。
  `List(ctx, userID, q)` 严格顺序：
  1. 再查 `display_enabled`（守卫之外每接口重查）。
  2. 权限过滤：`apiKeyService.GetAvailableGroups(ctx, userID)`（已排除过期订阅与撤销专属）→ active 且 openai|composite → `ListListedGroups(ids)` → `search`（display_name/description/原名，不区分大小写）→ model/protocol 筛选（该组有对应 enabled target）。
  3. 解析 `T`（`snapshot_id` 或 `aggregator_state.last_window_end`）、各组默认档案；**批量**读：快照 `(g,0)` 与 `(g,default)`、每 target 最新探测（`DISTINCT ON`，3×probe_interval 内）、默认 target 最新报告、range 对应 sparkline（tier5 快照 + 探测均值）。
  4. 叠加用户价格：`m = ResolveUserGroupRateMultiplier(...) × PeakMultiplierAt(pricingAt)`；`rp = pricing.Resolve(model=默认档案, group)`；可报价 iff `Mode==token && len(Intervals)==0 && InputPricePerToken>0`；`P = InputPricePerToken×1e6×m`、`C = CacheReadPricePerToken×1e6×m`（decimal 10 位）。预测倍率 `not_applicable`：报价不可用 / 快照 `cache_creation_tokens>0` / 缓存率非 ok / 参考价缺失 / 分母 0；`stale`：`reference_confirmed_at < now-7d`。历史均价单位：`subscription_share==1` → `quota_per_million`，0 → `usd_per_million`，混合 → `not_applicable/mixed_billing`。历史值不按当前倍率重算。
  5. 排序：最多三级 `field[:asc|desc]`，允许 `display_order,rate,historical_price,predicted_rate,ttft_fast95,cache_rate,success_rate`；state ∉ {ok,stale} 视为缺失、始终末尾；`group_id asc` 兜底。
  6. 内存分页（`total=len(filtered)`）；summary 基于过滤后全集：available=up、listed=total、abnormal=down、verified=passed 且未过期。
  7. 仅可缓存与用户无关的基底（按 T,range 键 30s）；个性化响应不缓存。
  `GetGroup`：`AuthorizeUserGroup` → 行 + detail（`e2e_availability_6h` = 6h 内 received 探测 passed/total，<12 个样本 → insufficient；TPS/首 Token/端到端取最近一次 passed 探测；`p90_ms` 取默认档案快照；`probe_tokens` 取最近样本）+ trend + 各档案健康/报告。`ListVerifications`：`AuthorizeUserGroup` → 分页（默认档案，可传 profile_id，page_size≤50）。
- `backend/internal/handler/provider_hall_handler.go`（`ProviderHallUserHandler`：`List/GetGroup/ListVerifications`）；参数校验 400 `PROVIDER_HALL_INVALID_QUERY`；`range` 默认 `cfg.DefaultRange`；`page_size∈[1,100]` 默认 50；`search≤100` rune。
- `backend/internal/server/routes/provider_hall_user.go`：
  ```go
  hall := authenticated.Group("/provider-hall"); hall.Use(panelRateLimiter.Heavy(), providerHallDisplayGuard(svc)) // display 关闭 → 404
  hall.GET("", h.ProviderHall.List); hall.GET("/groups/:id", h.ProviderHall.GetGroup); hall.GET("/groups/:id/verifications", h.ProviderHall.ListVerifications)
  ```
- 公共设置：`PublicSettings.ProviderHallEnabled` (`provider_hall_enabled`) 贯通 service/injection payload/dto/handler；来源 `SettingService.SetProviderHallDisplayReader(interface{ DisplayEnabled(ctx) bool })`，由 `ProvideProviderHallService` 注入，30s 缓存。

### 修改点

- `routes/user.go`、`handler/handler.go`（`Handlers.ProviderHall`）、`handler/wire.go`、`service/wire.go`（`ProvideProviderHallQueryService(readRepo, cfgRepo, apiKeyService, gateway *OpenAIGatewayService, pricing *ModelPricingResolver, groupRepo)`）。
- `service/provider_hall.go`：就绪标志 `Display`（B7 前端完成后置 true）。

### 测试（矩阵 F、H 查询部分）

- `provider_hall_query_test.go`：权限（普通/专属/订阅/过期/撤销/下架 → 不可见；详情与报告 404）；两个专属倍率用户交替访问不串价（race）；JSON 反射遍历无 `account/token/key/upstream/count` 键；三级排序缺失末尾；预测倍率固定样例（阶梯→n/a、写缓存→n/a、参考过期→stale）。
- `provider_hall_handler_test.go`：参数校验；管理路由普通用户 403。
- localdb：100 组 × 30 天快照的列表查询 `EXPLAIN` 不扫 `usage_logs`/`provider_hall_requests`，语句数固定（无逐行查询）。

---

## B6 管理 API 扩展与后台页签（约 3 人日）

### 后端（`routes/provider_hall.go`、`handler/admin/provider_hall_handler.go`、`dto/provider_hall.go`、`service/provider_hall_admin.go`）

```
POST /admin/provider-hall/groups/:id/probes          {profile_id, idempotency_key} → 202 {job_id, status, reused}
POST /admin/provider-hall/groups/:id/verifications   同上
GET  /admin/provider-hall/jobs?status&kind&group_id&page&page_size
GET  /admin/provider-hall/jobs/:id                   任务 + 样本（detail 不含 Key 材料）+ 报告
POST /admin/provider-hall/jobs/:id/cancel            queued→cancelled；running 且样本全 prepared→cancelled；有 dispatched → 409 PROVIDER_HALL_JOB_IN_FLIGHT
GET  /admin/provider-hall/groups/:id/verifications   含未上架组
GET  /admin/provider-hall/health
```
`EnqueueManualJob(ctx, kind, groupID, profileID, idemKey, actorID)`：校验目标启用/运营用户/Key（`PROVIDER_HALL_TASKS_DISABLED`、`PROVIDER_HALL_BUDGET_EXHAUSTED` 等为 BadRequest 业务码）；`ON CONFLICT (idempotency_key) DO NOTHING RETURNING id`，冲突返回已有任务 `reused:true`。
Health DTO：`Collection{Enabled; Nodes[]{node_id, epoch_id, version, heartbeat_at, confirmed_at, persisted_seq, overflowed, lost}; MissingExpected[]; OpenGaps[]; LocalQueue{Depth,Capacity,Dropped}}`、`Aggregator{Watermark, LagSeconds, DirtyCount, LastRunAt, LastError}`、`Reconciliation{Pending, Uncertain, Failed24h}`、`Budget{Day, Budget, ConfirmedSpend, UncertainSpend, InFlight, PausedReason}`、`Jobs{Queued, Running, Unknown, Failed24h}`。

### 前端

- `api/admin/providerHall.ts`：`ProviderHallJob/Sample/Health` 类型 + `enqueueProbe/enqueueVerification/listJobs/getJob/cancelJob/getHealth/listGroupVerifications`。
- `views/admin/ProviderHallView.vue`：tabs `config | groups | profiles | jobs | health`（沿用 `hall-tab`、`dirty`、i18n `admin.providerHall.<tab>`）。
- 新组件 `components/admin/provider-hall/ProviderHallJobs.vue`（筛选、`Pagination`、`ConfirmDialog` 取消）、`ProviderHallJobDetailDialog.vue`（`BaseDialog`，样本表、报告摘要）、`ProviderHallHealth.vue`（`StatCard` 网格、节点表、缺口、`AutoRefreshButton` 30s）；`ProviderHallGroups.vue` 目标行加"立即探测 / 立即检测"（幂等键每次点击生成一次并保留到响应返回）。
- i18n `admin/providerHall.ts` zh/en；扩展 `scripts/provider-hall-admin-smoke.mjs`（jobs/health 页签截图、发起探测 → 202）。
- 测试：`TestProviderHallAdminHandler_{Probes,Verifications,Jobs,Cancel,Health}`（非管理员 403）；Vitest `ProviderHallAdmin.spec.ts` 扩展（列表、取消、409 呈现、双击复用键）。

---

## B7 用户大厅前端（约 6 人日）

### 文件

```
src/api/providerHall.ts                 list(params, signal) / getGroup(id, range, signal) / listVerifications(id, params, signal)；类型对应 B5 DTO；ProviderHallMetric<T> 抽到 src/types/providerHall.ts 与 admin 共用
src/composables/useProviderHall.ts     状态 range/model/protocol/search/sort[]/page 写入路由 query；排序偏好 localStorage `providerHall.sort.<userId>`；
                                        fetch() AbortController + sequence 丢弃迟到响应；失败保留旧数据并 stale=true；
                                        useAutoRefresh 60s，shouldPause: () => document.hidden || loading；expanded Set<number>；detail 按 group+range 缓存
src/views/user/ProviderHallView.vue    AppLayout 卡片：HallSummary、HallToolbar、HallTable（桌面）/ 卡片列表（<768px）；顶部"监测中"徽章、说明、更新时间、图例
src/components/provider-hall/
  HallSummary.vue          可用/公开分组环形数字 + 异常线路/检测通过 + 图例（可用/异常/未知）+ AutoRefreshButton     复用 MetricCell 样式、HallRing
  HallToolbar.vue          时间维度 chips（.tabs/.tab/.tab-active）、排序 chips（默认/倍率/真实价格/平均最快95%/缓存命中/成功率/自定义，点击 asc→desc→off）、模型 Select、搜索
  HallSortDialog.vue       "自定义"三级排序编辑                                                                      BaseDialog、Select
  HallTable.vue            sticky 表头、固定列宽、键盘 Enter/Space 展开、方向键移动；emits expand/use-group          ChannelStatusV2View 展开行模式
  HallRow.vue              分组名/说明、倍率、真实价格/预测倍率（HallMetricTooltip）、状态+HallModelDots、HallTtftBar、HallRing×2、HallSparkline、HallVerificationBadge、"使用此分组"
  HallDetailPanel.vue      指标网格（倍率、输入价、1h 成功率、6h 端到端可用率、1h 缓存率、首 Token、端到端、平均最快95%、P90、TPS、最近监测）+ 双线趋势（探测 vs 最快95%，缺口区带 plugin）+ 模型健康 chips + Token 输入/输出 + 报告入口   MonitorTrendChart 改造
  HallRing.vue             纯 SVG 环（r=18, stroke-dasharray）；成功率评分色 clamp((rate-0.5)/0.5)；state 感知（灰 + 样本不足/缺口/过期）   新建
  HallSparkline.vue        纯 SVG 双折线（null 断开）+ 缺口矩形，约 200×40，aria-label                              新建
  HallTtftBar.vue          数值 + 快/中/慢 + 相对页面最大值的条                                                       V2 error-rate bar 模式
  HallModelDots.vue        每档案点（up 绿 / down 红 / unknown 灰）+ tooltip（模型、协议、checked_at）                 status-dot
  HallVerificationBadge.vue 判定 pill（passed/failed/suspected/insufficient/expired/none）+ "查看检测报告"           .badge-*
  HallVerificationDialog.vue 报告列表 + 单份报告（判定、方式、每套断言结果、模型比对）                                BaseDialog width="wide"、Pagination
  HallMetricTooltip.vue    包装 HelpTooltip：公式文案 + 窗口/计算时间 + state 原因                                    HelpTooltip
  UseGroupDialog.vue       BaseDialog：Tab"切换现有 Key"（keysAPI.list，显示 原分组→新分组，keysAPI.update(id,{group_id})）/ Tab"创建新 Key"（keysAPI.create(name, groupId)）；错误 showError(extractApiErrorMessage)、保留输入
src/utils/providerHallFormat.ts        formatMs / formatRate / formatPrice（十进制字符串，不做浮点运算）/ stateLabel(reason_code)
src/i18n/locales/zh/providerHall.ts, en/providerHall.ts   顶级 key providerHall.*；common.ts 增 nav.providerHall；花括号用 {'{'} 转义；在 zh/index.ts、en/index.ts 展开
```

- 路由：`{ path: '/providers', name: 'ProviderHall', component: () => import('@/views/user/ProviderHallView.vue'), meta: { requiresAuth: true, titleKey: 'nav.providerHall' } }`。
- 侧栏：`buildSelfNavItems` 在 `/monitor` 之后加 `{ path: '/providers', label: t('nav.providerHall'), icon: ChannelIcon, featureFlag: flagProviderHall }`；`utils/featureFlags.ts` 增 `providerHall = defineFlag({ key: 'provider_hall_enabled', mode: 'opt-in' })`；`types/index.ts PublicSettings.provider_hall_enabled`；`SettingsView.vue` 按该文件头部清单同步。
- 值渲染只依赖后端 `Metric.state`：`ok` 值；`stale` 值 + "过期"角标；`incomplete` 值（若有）+ "数据不完整"；`insufficient/disabled/not_applicable` 按 reason_code 显示占位文案。缺失值排序由后端保证。
- 图表配色遵循 `dataviz` skill；深色主题只用 `dark:` 或 `.dark-theme`。

### 测试（矩阵 G Vitest）

- `components/provider-hall/__tests__/`：`HallRing`（角度/颜色阈值）、`HallSparkline`（缺口矩形生成）、`HallTable`（状态渲染、移动端列、键盘）、`UseGroupDialog`（创建/切换参数、失败保留输入）。
- `composables/__tests__/useProviderHall.spec.ts`：筛选写路由、三级排序持久化、迟到响应丢弃、刷新失败保留旧值+stale、隐藏页暂停。
- `views/user/__tests__/ProviderHallView.spec.ts`：整页挂载 + mock API。
- 现有 `localesMessageCompile`、`localesNoKeyCollision` 自动覆盖。

---

## B8 全流程验收、性能与发布准备（约 4 人日）

- **共享模拟上游** `backend/internal/testutil/fakeopenai`（Responses/Chat/Messages SSE，可注入延迟/429/流中断，支持 `hall_echo` 工具并回显模型），Go 测试与 e2e 共用；e2e 侧另有 Node 版 `frontend/scripts/lib/fake-openai-upstream.mjs`。
- **e2e** `frontend/scripts/provider-hall-e2e.mjs`（npm `test:e2e:provider-hall`，pnpm 9）：启动隔离实例（临时 DATA_DIR/PG/Redis、`RUN_MODE=standard`、正常 AUTO_SETUP，禁止写 `.installed`）→ 启动模拟上游 → 管理员建档案、上架 2 组、绑定专用 Key、开启三个开关 → 普通 Key 每协议 ≥250 条模拟请求（含失败）→ 等两轮聚合 → 手动探测与检测 → 普通用户打开 `/providers`：断言行数据、排序、展开详情、报告弹窗、"使用此分组"创建与切换 Key → 390/1440/1920 × 明暗截图，检查 `scrollWidth<=innerWidth`、按钮截断、键盘操作 → 失败截图到 `PROVIDER_HALL_ARTIFACT_DIR` → 清理。
- **性能（H）** `backend/internal/repository/provider_hall_perf_test.go`（`//go:build perf`）：`COPY` 造 100 组 × 30 天快照 + 100 万 requests；`EXPLAIN (FORMAT JSON)` 列表/详情无 `provider_hall_requests`/`usage_logs` 顺序扫描；20 并发列表 P95 ≤500ms 热 / ≤2s 冷；`TestProviderHallGatewayOverhead` 对比采集开关前后 P95 增量 ≤5ms。
- **回归**：`make -C backend test-unit`、`CI=1 make -C backend test-integration`、`make -C backend build`、`make test-frontend`；聚焦 `-run 'ProviderHall|ChannelMonitorV2|OpenAIGateway|UsageBilling|WS'`。基线既有失败（`AccountUsageCell.vue` 等 typecheck、既有 test-unit 失败）单独记录不计入。
- **证据**：每批一份 `docs/reviews/PROVIDER_HALL_B{2..8}_<date>.md`（实际命令、用例数、`docker info` 结果、截图路径、未执行项，"跳过"不写"通过"）；更新 `docs/PROVIDER_HALL_DEVELOPMENT_PLAN.md` 状态段。

---

## 开关解锁：用就绪标志替代一刀切

`service/provider_hall.go:160-162` 改为：
```go
type ProviderHallReadiness struct{ Collection, Tasks, Display bool }
func (s *ProviderHallService) SetReadiness(r ProviderHallReadiness)   // ProvideProviderHallService 中设置
// UpdateConfig
if cfg.CollectionEnabled && !s.ready.Collection { return nil, ErrProviderHallNotReady.WithMetadata(map[string]string{"switch":"collection_enabled"}) }
if cfg.TasksEnabled && (!s.ready.Tasks || !cfg.CollectionEnabled || cfg.GatewayOrigin=="" || cfg.OperatorUserID==nil) { ... }
if cfg.DisplayEnabled && (!s.ready.Display || !cfg.CollectionEnabled) { ... }
```

| 开关 | 置 true 时机 | 前置 |
|---|---|---|
| `collection_enabled` | B3 结束（collector+aggregator 已 wire） | 节点名单已登记 |
| `tasks_enabled` | B4 结束（runner 已 wire） | collection 开、运营用户、gateway_origin、预算、至少一个有效 target |
| `display_enabled` | B7 结束（用户 API + 前端已 wire） | collection 开、至少一个 listed 分组 |

现有"任一开关→NotReady"单测拆成三例。灰度顺序沿用计划：迁移且关闭 → 全节点升级并登记 → 仅采集 → 少量目标探测 → 一小时核对 → 用户开放 → 观察 24h。

## Wiring 汇总

- `repository/wire.go`：`NewProviderHallFactRepository`、`NewProviderHallAggregationRepository`、`NewProviderHallJobRepository`、`NewProviderHallReadRepository`（均 `*sql.DB`，经 `ProvideSQLDB`）。
- `service/wire.go`：`ProvideProviderHallCollector`（所有实例）、`ProvideProviderHallAggregator`、`ProvideProviderHallRunner`（`shouldStartGlobalWorkers`）、`ProvideProviderHallService`（就绪标志 + settings reader）、`ProvideProviderHallQueryService`；`OpenAIGatewayService.SetProviderHallBillingSink(collector)`。
- `cmd/server/wire.go provideCleanup`：Runner → Aggregator → Collector(flush) 顺序停。
- `handler/wire.go`、`handler/handler.go`：`Handlers.ProviderHall *ProviderHallUserHandler`。
- 每批结束 `make -C backend generate`（wire + ent）。

## 验证命令

```sh
# 后端聚焦（每批）
GOPROXY=https://goproxy.cn,direct go -C backend test -tags=unit -race ./internal/service ./internal/handler/... ./internal/repository -run ProviderHall -count=1
PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin go -C backend test -tags=providerhall_localdb ./internal/repository -run ProviderHall -count=1 -v
CI=1 go -C backend test -tags=integration ./internal/repository -run ProviderHall -count=1   # 有 Docker 时
make -C backend generate && make -C backend build
# 前端
pnpm --dir frontend run lint:check && pnpm --dir frontend run typecheck
pnpm --dir frontend exec vitest run src/components/provider-hall src/components/admin/provider-hall src/composables/__tests__/useProviderHall.spec.ts src/views/user/__tests__/ProviderHallView.spec.ts src/views/admin/__tests__ src/api/__tests__
pnpm --dir frontend run test:e2e:provider-hall-admin
pnpm --dir frontend run test:e2e:provider-hall            # B8 新增
# 交付前
make -C backend test-unit ; CI=1 make -C backend test-integration ; make test-frontend
go -C backend test -tags=perf ./internal/repository -run ProviderHall -count=1   # H
```

Node 20 + pnpm 9、Go 1.26.5（另做 1.26.6 兼容）；`docker info` 不可用时集成用例必须以 localdb 路径补跑并在记录中注明。

## 批次顺序与工作量

B2(5) → B3(5) → B4(7) → B5(4) → B6(3) → B7(6) → B8(4)，合计约 34 人日。B5 可与 B4 并行（B5 读快照/样本表，B4 写）；B7 可在 B5 DTO 定稿后用 mock 数据先行。每批结束写实施记录、跑聚焦命令，不宣称全量通过。
