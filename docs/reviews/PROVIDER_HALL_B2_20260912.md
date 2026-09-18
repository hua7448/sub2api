# 供应商大厅 B2：采集与计费关联实施记录

日期：2026-09-12。对应 [完整实施计划](../PROVIDER_HALL_DEVELOPMENT_PLAN.md) 的 B2 批次。本批只交付请求事实采集、计费结果关联和任务头处理；聚合、快照、任务、用户 API 与页面属于 B3～B8，仍未开放。`collection_enabled` 开关继续由就绪标志锁定，B3 接线完成前后台无法启用。没有部署生产。

## 已实现的行为

- **请求跟踪器**（`ProviderHallRequestTracker`）：三个网关入口（`/openai/v1/responses`、`/v1/chat/completions`、`/v1/messages`）在鉴权、参数、准入和定价上下文之后登记候选样本，`defer Finish()` 保证恰好一次终态。所有方法对 nil 接收者安全，采集关闭时 handler 路径零改动。`/responses` 子路径（compact 等）不进入排行范围。
- **提交计数在真实出站边界**：`WrapProviderHallUpstream` 包装 `HTTPUpstream`，只有 POST 到 `/responses`、`/chat/completions`、`/messages` 的模型调用计一次提交；令牌刷新、模型列表、能力探针、compact 都不计。同账号重试和换号都计数，账号选择、利润否决不计。
- **终态判定**：成功 > 客户端取消（`excluded/client_cancel`）> 策略拒绝（`excluded/policy_reject`）> 有失败尝试或无可用账号（`failed`）> 未提交即退出（`excluded/early_exit`）。HTTP 200 后流中断记为 `failed`。TTFT 从请求到达计时，含本机鉴权、排队与重试；平台取首个收到提交的账号，`group_id` 恒为入口分组。
- **任务头**：`X-Provider-Hall-Task = <kind>:<job>:<sample>:<trace>` 在三个入口一进入就从入站请求删除，采集关闭或不在排行范围时同样删除；探测 Key（登记时间 ≤ 请求开始）携带合法头 → `source=probe|verification`、复用 runner 的 trace 与 sample；普通 Key 携带头 → `source=user`、`excluded/task_header_mismatch`；探测 Key 无头 → `source=probe`、无 sample。所有上游构造器都是白名单透传，WS 入口只复制固定头名，该头不可能出站。
- **采集器**（`ProviderHallCollector`）：每个处理网关流量的实例启动（不受 primary 限制），有界队列 8192、每 250ms 或 500 条一批单事务落库、每秒 barrier 推进水位、10s 心跳、30s 刷新配置/目标/探测 Key 快照。队满丢弃并开 `queue_overflow` 缺口，队列回落且新 barrier 落库后关闭；写库失败保留批次并退避，超过 4 批则按溢出处理。`Stop()` 5 秒内刷完，未刷完写 `node_exit` 缺口并标记退出；未干净退出的旧 epoch 在下次注册时标 `crash` 并从上次确认水位开缺口。
- **事实表**（迁移 231）：`provider_hall_collector_epochs`、`coverage_gaps`、`requests`、`dirty_buckets`。`requests` 以 trace 为主键，开始事件 `DO NOTHING`、终态事件只在 `outcome='pending'` 时合并，先到后到顺序无关且幂等；终态和脏桶标记同事务提交。
- **计费关联**：`OpenAIRecordUsageInput.ProviderHallTraceID` 在 handler 闭包外取值并传入异步入账；`RecordUsage` 保留 `applyUsageBilling` 的 `applied` 结果，映射为 `applied / duplicate / failed / uncertain(请求冲突) / not_applicable(简单模式、图片/视频/搜索、非 token 计费)`，连同计费键、指纹、倍率、`InputCost/TotalCost/ActualCost` 十进制快照写回 `requests`。已结算状态不会被后续事件覆盖，`uncertain` 可被后续 `applied` 补齐。扣费路径的返回行为不变。
- **就绪标志**：`ProviderHallReadiness{Collection,Tasks,Display}` 替代一刀切；`ProvideProviderHallService` 由已接线的 worker 推导，本批 `Collection=false`，`collection_enabled=true` 仍返回 `409 / PROVIDER_HALL_NOT_READY{switch:collection_enabled}`。

## 源码

- `backend/migrations/231_provider_hall_facts.sql`
- `backend/internal/service/provider_hall_tracker.go`、`provider_hall_collector.go`、`provider_hall.go`（就绪标志）、`openai_gateway_usage.go`（`emitProviderHallBilling`）、`openai_gateway_service.go`（`SetProviderHallBillingSink`、上游包装）
- `backend/internal/repository/provider_hall_fact_repo.go`
- `backend/internal/handler/provider_hall_gateway_hooks.go`；`openai_gateway_handler.go`、`openai_chat_completions.go` 中的 Begin/NoAccount/Reject/SelectAccount/Attempt 调用点
- Wire：`repository/wire.go`、`service/wire.go`（`ProvideProviderHallCollector` 所有角色启动）、`handler/wire.go`、`cmd/server/wire.go`（Runner → Aggregator → Collector 顺序停止）
- 测试：`service/provider_hall_algorithm_test.go`、`provider_hall_collector_test.go`、`provider_hall_billing_link_test.go`、`handler/provider_hall_gateway_hooks_test.go`、`repository/provider_hall_fact_db_test.go`（经 `provider_hall_contracts_registry_test.go` 注册进 localdb 套件）

## 测试覆盖

- **矩阵 A 固定样例**（`provider_hall_algorithm_test.go`）：1000..19000 + 100000 → fast95 = 10000、P90 = 18000；N = 0/1/10/20；800/1100 重复边界值；输入 map 不被改写；成功率 2/3、重试分母、无提交分母；缓存率含写缓存稀释；输入费用分摊、$1/M 均价、零账单合法、负值/零分母/输入大于总额拒绝；上海预算日跨零点。
- **入口级 hook**（`provider_hall_gateway_hooks_test.go`，真实 `OpenAIGatewayService` + 记录请求头的 stub 上游 + 真实 `ProviderHallCollector`）：三入口非流式成功（提交 1、平台、分组、Key、profile 映射、usage、计费事件 `not_applicable` 回链同一 trace）；Responses SSE 成功记录首内容时间与响应模型；HTTP 200 后流中断 → `failed`；429 换号 → 成功且提交数等于真实上游调用数；全部失败 → `failed`、提交数含同账号重试与换号；无账号 → `failed`、0 提交、无平台；转发前客户端取消 → `excluded/client_cancel`；未上架分组不采集也不发计费事件；采集器为 nil 时头仍被剥离。
- **`TestProviderHallTaskHeaderNeverForwarded`**：三入口 × 普通 Key（头不出站、`task_header_mismatch`、trace 不复用）× 探测 Key（头不出站、`source=probe`、trace/sample 复用）+ 探测 Key 无头。
- **跟踪器/采集器单测**（既有 `provider_hall_collector_test.go`）：九种终态、出站边界计数、nil 安全、计费先于/晚于终态、溢出开关缺口、写失败保留与有界、启停压测。
- **数据库合约**（`provider_hall_fact_db_test.go`，5 组，四条迁移路径各跑一遍）：epoch 注册/心跳/干净退出/崩溃缺口从确认水位起算；开始→终态→重复终态幂等、终态→迟到开始不回退、第二终态不覆盖、混合批次只有 user+profile 行标脏、脏桶唯一、`ttft_ms>0` 与枚举约束；计费先于终态保留、终态先于计费再标脏、已结算不被覆盖、`uncertain→applied`、`not_applicable` 终态、未知 trace 静默、非法状态整批回滚、`numeric(24,10)/(20,8)/(12,6)` 精度原样读回；无 barrier 不推进水位、水位与 seq 不回退、缺口开/关/夹取/幂等；已启用目标与探测 Key 索引、停用后目标消失而 Key 永久保留。

## 本轮验证

以下为 2026-09-12 实际执行结果，Go 1.26.5，`GOPROXY=https://goproxy.cn,direct`。

1. `go test -tags=unit -race ./internal/service ./internal/handler/... ./internal/repository -run ProviderHall -count=1`：五个包全部通过，共 58 个顶层用例（含同一目录内其他批次已落地的用例）。
2. `PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin PROVIDER_HALL_PG_PORT=15481 go test -tags=providerhall_localdb ./internal/repository -run ProviderHall -count=1 -v`：在仅含 B2 源码的隔离副本（其他批次同时改动中的文件替换为接线 stub）上，`empty / upgrade_228 / upgrade_229 / upgrade_230` 四条路径各执行迁移两次、基础合约、目标合约和本批 5 组事实合约，45 个子用例全部通过。在共享工作树上，同一时刻 repository 测试包因其他批次尚在编写的 `provider_hall_job_db_test.go` 编译失败，B2 合约本身在共享树上此前一次可编译的运行中亦全部通过。
3. `go vet -tags=unit ./internal/...` 通过；`gofmt -l internal` 对本批文件无输出（顺手格式化了两份既有 B2 测试文件）。
4. `go build ./...` 通过。

未执行：`CI=1 go test -tags=integration`（本机 `docker info` 不可用，Docker 路径的 PostgreSQL/Redis 集成套件没有运行，事实表合约以 localdb 路径替代并已注明）；`make -C backend generate`（本批不改 Ent/Wire 定义）。

## 已知限制

- 透传路径在收到 HTTP 200 后缺少终态事件时丢弃部分结果（`openai_gateway_passthrough.go` 中 `return nil, errors.New("stream usage incomplete: missing terminal event")`），因此该场景的 `first_content_at`/usage 不会进入事实行，只保留 `failed` 终态与提交数。改动会牵连 #5148 的部分结果入账语义，本批未修改，留待后续评估。
- 入口级测试通过 `openai_passthrough` API Key 账号驱动；Chat/Messages 入口在此配置下经 CC/Messages→Responses 转换后以流式访问上游，非流式客户端的 TTFT 在这两条入口上不作断言。OAuth/Codex 账号、WS 入口未在本批 hook 测试中覆盖。
- 客户端取消只覆盖“转发前取消”；流中途断开依赖上游 `ClientDisconnect` 标记，已在跟踪器单测中以直接调用覆盖，未在 handler 级模拟。
- 采集器目前只写 `requests`，探测计费写入 `provider_hall_spend` 由 B4 在 `WriteBatch` 中扩展。
