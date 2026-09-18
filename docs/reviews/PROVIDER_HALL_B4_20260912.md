# 供应商大厅 B4：任务 Runner、探测、模型检测、预算与对账实施记录

日期：2026-09-12。本批在 B2 采集与 B3 聚合之上实现主动任务运行时。用户 API、后台任务页签和用户大厅页面属于 B5～B7，未在本批实现；本批未部署生产，`tasks_enabled` 仍需管理员在配置齐备后显式开启。

## 已实现的行为

- **迁移 `233_provider_hall_jobs.sql`**：新增 `provider_hall_jobs`、`provider_hall_samples`、`provider_hall_verifications`、`provider_hall_spend`，字段、约束与索引与实施计划一致。任务不对目标建外键，目标或分组删除后历史保留；样本 `(job_id,test_id,seq)` 与 `trace_id` 唯一；费用表以“计费请求 ID + Key ID”为主键。
- **排程**：primary 实例每 15 秒在领导锁 `provider-hall-scheduler`（1 分钟）下运行一次。对每个已启用、分组已上架、已绑定专用 Key 的目标，按 `now.Truncate(interval)` 生成当前时隙的探测与检测任务（`ON CONFLICT DO NOTHING`），探测 `not_before` 抖动 ±30 秒、检测 0～5 分钟；停机期间的时隙不补跑。目标停用、分组下架、Key 解绑或删除时排队任务取消为 `target_disabled`；快照中的目标版本、档案版本或 Key 与现行不一致时取消为 `config_changed`。
- **领取与租约**：`FOR UPDATE SKIP LOCKED` 领取，同一目标不会同时有两个 running 任务；租约 120 秒、每 30 秒续约；续约失败时任务协程立即停止。租约过期的任务：样本没有 dispatched/uncertain 时回到 `queued` 可被重新领取；有在途样本时标记 `unknown`，之后只核对、不重发。
- **发出前检查**（单事务 + `pg_advisory_xact_lock(hashtext('provider_hall_dispatch'))`，固定顺序）：`tasks_enabled` 与目标启用、分组上架、快照版本/Key/运营用户/网关 origin 未变 → 用当前行重跑 `ValidateProviderHallTargetBinding`（不通过则 `failed/probe_key_invalid`，不发送付费请求）→ 任一样本账单超过 120 秒仍为 pending/uncertain 则任务回队列 `billing_backlog`，60 秒后再试 → 当日已确认花费达到 `daily_budget`（预算大于 0 时）则 `cancelled/budget_exhausted` → 全局 dispatched 样本少于 2 且本目标无在途，否则等待（最多 60 秒后交还任务）。通过后样本先以 `dispatched` 提交，再发送。
- **预算日**：任务首次发出时按 Asia/Shanghai 日期固定 `budget_day`，后续样本和费用继承；修改预算不清零当日花费；超出预算后到账的在途样本账单仍记入 `provider_hall_spend`。
- **发送**：`http.Client{Timeout: 90s}`，`CheckRedirect` 拒绝跳转；目标为 `gateway_origin + /v1/responses|/v1/chat/completions|/v1/messages`；非环回地址必须为 https，环回地址只在 `PROVIDER_HALL_ALLOW_LOOPBACK=1` 时允许，`UpdateConfig` 与运行时都按该规则校验。请求头：`Authorization: Bearer <专用 Key>`、`Accept: text/event-stream`、`X-Provider-Hall-Task: <kind>:<job>:<sample>:<trace>`、`X-Client-Request-ID: <trace>`、`User-Agent: sub2api-provider-hall/<version>`；Messages 协议另附 `x-api-key` 与 `anthropic-version`。Key 明文只用于这一次请求头，不写入任何表和日志。
- **结果分类**：2xx 且流达到终态 → `received`；明确非 2xx → `received/error/http_<status>`；连接被拒、DNS 失败、跳转 → `received/error`（未准入，不属于不确定）；超时或响应头之后流中断 → `uncertain`，永不重发。每个 tick 用 `provider_hall_requests.trace_id` 对账：请求事实到达终态时补齐首 Token、Token 数、模型名和账单（探测按成功即通过，检测样本因内容已丢失记为 `error/response_lost`）；10 分钟内无事实则 `received/error/uncertain_timeout`。账单镜像：B2 采集器在 `WriteBatch` 内对 `source ∈ {probe,verification}` 的计费事件同事务写入 `provider_hall_spend`（applied→confirmed、failed→failed、uncertain→uncertain、duplicate 不写、not_applicable 只把样本标为 unbilled）并同步 `provider_hall_samples.billing_status`；对账再从事实表兜底一次。
- **探测客户端**：三协议请求体（Responses `input_text`/`max_output_tokens`/`store:false`、Chat `stream_options.include_usage`、Messages `max_tokens`）与 SSE 解析：TTFT 取首个非空文本增量或首个工具调用项；模型名与 usage 分别来自 `response.completed`、Chat 最终 usage 块和 `message_start`/`message_delta`；终态分别为 `response.completed`、`[DONE]`、`message_stop`。探测输出上限 256、检测 1024，且不超过档案 `output_limit`；`generation_ms = total_ms − ttft_ms`，`TPS = output_tokens / generation_ms × 1000`。
- **模型检测**：算术 ×3（`a,b ∈ [100,999]`，`+ − *`）、JSON 约束 ×3（8 位 token 与 length，键集合必须精确）、`supports_tools` 时强制 `hall_echo` 工具调用 ×3（三协议各自的 `tools`/`tool_choice` 形态）。围栏 ``` 被剥离并在 detail 记录 `fenced:true`，不判失败。模型名比对 `∈ {model} ∪ model_aliases`。判定优先级固定：模型不符 ≥2 → `suspected`；任一必需套件失败 ≥2 → `failed`；样本缺失/异常/模型名缺失 → `insufficient`；否则 `passed`。`execution_status` 为 completed/partial/error；报告 48 小时过期。档案的 model/protocol/supports_tools/aliases 变化时，配置服务在保存后立即把该档案的报告标为 `stale`，runner 每 tick 再按报告中记录的档案指纹核对一次，两条路径结果一致。
- **收尾**：探测任务样本通过则 `succeeded`，否则 `failed` 并记录样本错误码；检测任务写入 `provider_hall_verifications`。任务存在 uncertain 样本时先转为 `unknown`，样本对账完成后由后续 tick 评分收尾，不产生新请求。
- **管理入口（供 B6 使用）**：`ProviderHallRunner.EnqueueManual` 校验开关、目标、Key 和预算，幂等键相同的重复请求返回已有任务并标记 `reused`；`CancelJob` 对 queued 或样本全部 prepared 的 running 任务生效，有在途样本返回 `PROVIDER_HALL_JOB_IN_FLIGHT`，已结束任务返回 `PROVIDER_HALL_JOB_NOT_CANCELLABLE`。
- **停机**：`Stop()` 停止领取并最多等待 90 秒在途发送完成，之后 B3 聚合器与 B2 采集器再停止（顺序由 `cmd/server/wire.go` 的清理步骤保证）。`PROVIDER_HALL_DISABLE_RUNNER=1` 跳过启动。

## 源码

- `backend/migrations/233_provider_hall_jobs.sql`
- `backend/internal/service/provider_hall_runner.go`、`provider_hall_runner_types.go`（任务/样本/报告/费用类型与 `ProviderHallJobRepository` 接口）、`provider_hall_probe_client.go`、`provider_hall_verification.go`、`provider_hall_budget.go`
- `backend/internal/repository/provider_hall_job_repo.go`（新增 `NewProviderHallJobRepository`，已加入 `repository/wire.go`）；`provider_hall_fact_repo.go` 的 `WriteBatch` 增加探测账单 → 费用表/样本镜像
- `backend/internal/service/provider_hall.go`：`gateway_origin` 环回规则、`SetJobControl`、`SaveProfile` 后的报告 stale 标记；`provider_hall_target.go`：`SaveTargets` 后取消已停用/版本变化目标的排队任务。两个钩子都在仓库事务提交之后执行（仓库层拥有事务），runner 每 tick 的收敛逻辑保证钩子缺席或失败时结果仍一致
- `backend/internal/service/wire.go`：`ProvideProviderHallRunner(jobs, cfgRepo, db, lockCache, cfg, buildInfo)`；`Ready()` 现返回 true，`ProvideProviderHallService` 据此把 `tasks_enabled` 的就绪标志置为 collection 就绪 ∧ runner 就绪
- 测试：`provider_hall_verification_test.go`、`provider_hall_probe_client_test.go`（unit）、`provider_hall_job_db_test.go`（`integration || providerhall_localdb`，通过 `registerProviderHallLocalDBContract("b4_jobs", …)` 注册）

## 尚未接线的一行

`ProvideProviderHallService` 不在本批的修改范围内。要让配置保存时立即取消排队任务/标记报告过期，需要在该函数中加入 `svc.SetJobControl(runner.JobControl())`。缺少这一行时功能仍然正确：runner 每 15 秒会取消漂移任务并按指纹标记报告，发出前检查也会拒绝版本不一致的任务，只是延迟最多一个 tick。

## 本轮验证

Docker 在本机不可用（`docker info` 失败），`integration` 标签的用例以 `providerhall_localdb` 路径补跑。以下为 2026-09-12 实际执行结果：

1. `gofmt -l` 对本批新增与修改文件无输出；`go build ./...`、`go vet -tags=unit ./internal/service ./internal/repository`、`go generate ./cmd/server`（wire）与 `go test -tags=unit ./cmd/server` 通过。
2. `go test -tags=unit -race ./internal/service ./internal/repository -run ProviderHall -count=1`：通过。本批新增 51 个 unit 子用例：套件构造（无 `supports_tools` 时无工具套件、`output_limit` 封顶）、18 组断言评估（含围栏容忍、千分位、多余键、工具参数）、11 组判定优先级（2 次不符 + 3 次失败 → 疑似；同套件 2 次失败 → 失败；1 失败 + 1 异常 → 证据不足；缺样本/uncertain/模型名缺失 → 证据不足；失败优先于证据不足；全异常 → execution error）、别名匹配、48 小时与 stale 过期、指纹与别名顺序无关、预算/上海日/账单状态映射；三协议 6 组 SSE fixture（文本与工具调用）、流中断、`response.failed`；三协议请求体；经 httptest 网关真实发送并断言 5 个请求头与 `X-Client-Request-ID` 回显、429 → `http_429`、响应头后中断 → uncertain、连接被拒 → 非 uncertain、跳转被拒且目标未收到 Key、无 `PROVIDER_HALL_ALLOW_LOOPBACK` 时拒绝环回、非环回必须 https。
3. `PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin PROVIDER_HALL_PG_PORT=15483 go test -tags=providerhall_localdb ./internal/repository -run ProviderHall -count=1 -v`：PostgreSQL 18 隔离集群，`empty`、`upgrade_228`～`upgrade_232` 六条路径全部通过，共 156 个子用例（含 B2/B3 合约），其中 B4 每条路径 8 组：
   - 排程生成探测 + 检测两个时隙任务且重跑不重复；两个并发领取者只有一个拿到任务；探测经 httptest 网关恰好发出一次请求，任务头、模型回显、Token、TTFT、`client_request_id`、健康点与趋势序列均可读回；
   - 检测任务生成 6 个样本、生成报告（假网关一律回答 OK，判定 `failed/arith_failed`，模型全部匹配），修改档案 `supports_tools` 后报告立即 stale；
   - prepared 状态崩溃后租约过期回队列并可再领取（attempts=2）；dispatched 状态崩溃转 `unknown`，后续 tick 不发送任何请求；事实到达后对账补齐 TTFT/账单，收尾为 `succeeded` 且网关计数不变；无事实的 uncertain 样本 10 分钟后 `uncertain_timeout`；
   - 幂等键复用返回同一任务并 `reused=true`；未启用目标拒绝；管理取消的三种状态（queued 可取消、样本全 prepared 可取消、有在途 → `PROVIDER_HALL_JOB_IN_FLIGHT`）；
   - 停用目标后排队任务立即 `cancelled/target_disabled`，tick 不再排程；领取后目标版本变化在发出前 `cancelled/config_changed` 且不发送；漂移快照被 tick 取消；
   - Key 停用后不发送，任务 `failed/probe_key_invalid`，样本仍为 prepared；
   - 16:30Z 发出的样本 `budget_day=2026-09-13`（上海零点归属），任务与费用继承该日；本目标在途时第二个样本等待；预算耗尽后到账的在途账单仍写入费用表（`1.75` 中含 `0.25`）且样本 `confirmed`；预算达限时发出前检查取消、手动发起返回 `PROVIDER_HALL_BUDGET_EXHAUSTED`；账单挂起超过 120 秒回队列 `billing_backlog`，未超过时放行；全局 2 个 dispatched 样本时第三个等待；
   - `tasks_enabled=false` 时已领取任务在发出前取消；去掉 `PROVIDER_HALL_ALLOW_LOOPBACK` 后环回 origin 不发送，样本记 `request_invalid`。
4. 两个 primary 同时领取、崩溃重启不重发两类演练均在上述数据库用例中以真实 PostgreSQL 行锁与租约完成，没有依赖 mock。

未执行：`make -C backend test-unit` / `test-integration` 全量、Go 1.26.6 兼容、B6 管理 API 与前端页签（属后续批次）。没有向任何真实上游发出付费请求。
