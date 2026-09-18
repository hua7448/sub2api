# 供应商大厅 B5：用户 API 与定价实施记录

日期：2026-09-12。本批在 B2～B4（采集、聚合、任务）之上交付用户端只读 API、用户价格叠加、预测倍率、排序分页和公共设置开关。B7 前端已按同一 DTO 合约先行完成；本批未做浏览器联调（属 B8），未部署。

## 已实现的行为

- 新增用户接口，全部挂在 `authenticated` 链、`panelRateLimiter.Heavy()` 与展示守卫之下；`display_enabled=false` 时守卫直接 404，接口内部再次重查开关：
  - `GET /api/v1/provider-hall?range&model&protocol&search&sort&page&page_size[&snapshot_id]`
  - `GET /api/v1/provider-hall/groups/:id?range`
  - `GET /api/v1/provider-hall/groups/:id/verifications?profile_id&page&page_size`
- 列表固定顺序：重查 `display_enabled` → `apiKeyService.GetAvailableGroups`（已排除过期订阅、撤销专属）且 active、平台为 openai/composite → 已上架分组 → `search`（展示名/说明/原名，不区分大小写）→ model/protocol 筛选（该组有对应已启用 target，只影响出现与否）→ 解析 T（`snapshot_id` 或 `aggregator_state.last_window_end`）→ 一次性批量读取 → 叠加用户价格 → 排序 → 内存分页 → 汇总。
- 行级指标：TTFT 最快 95% 均值、P90、成功率取合并快照（`profile_id=0`）；缓存率、历史均价、预测倍率取默认档案（配置默认 model+protocol 对应 target 若启用，否则该组最小 profile_id 的已启用 target）。状态机沿用 B3 的 `disabled → incomplete → insufficient → stale → ok`。
- 当前报价：`ResolveUserGroupRateMultiplier × Group.PeakMultiplierAt(pricing_at)`，报价可用当且仅当 token 模式、无区间阶梯且输入价 > 0；价格为十进制字符串（10 位），订阅分组单位 `quota_per_million`，其余 `usd_per_million`。
- 预测倍率：`(P×(1−h)+C×h) / (Pr×(1−hr)+Cr×hr)`，倍率已折入 P/C。`not_applicable` 原因码：`quote_unavailable`、`cache_creation_present`、`cache_rate_unavailable`、`reference_missing`、`reference_zero`；参考价确认时间超过 7 天 → `stale/reference_expired`；缓存率本身 stale → 预测倍率随之 stale（与计划“缓存率非 ok → n/a”相比，stale 视为可用但标记过期，见“偏差”）。
- 历史均价单位：`subscription_share=1` → quota，`0` → usd，混合 → `not_applicable/mixed_billing`；历史值不按当前倍率重算。
- 健康点：每个已启用 target 取最近一次 received 探测，且距今不超过 3×probe_interval；passed → up，否则 down；无新鲜探测 → unknown。分组状态：任一 down → down；否则任一 unknown → unknown；否则 up。汇总 `available`=up、`abnormal`=down、`verified`=检测 passed 且未过期，均基于筛选后的全集。
- 排序：最多三级 `field[:asc|desc]`，字段限 `display_order,rate,historical_price,predicted_rate,ttft_fast95,cache_rate,success_rate`；未发布值（state ∉ {ok,stale}）无论方向始终排末尾；`group_id asc` 兜底。重复字段、未知字段、非法方向、超过三级均 400。
- 迷你趋势 / 详情趋势：按 range 取桶（6h/300s×72、24h/900s×96、7d/3600s×168、30d/21600s×120），真实曲线读 tier5 合并快照，缺失或覆盖不完整 → `null + gap:true`；探测曲线按 `date_bin` 聚合 probe 样本平均耗时。
- 详情面板：6h 端到端可用率（received 探测 passed/total，<12 个样本 → `insufficient/samples_below_12`）、TPS / 首 Token / 端到端（最近一次 passed 探测，超过 5 分钟标 stale）、P90（默认档案快照）、最近样本 Token 输入/输出、各档案健康与检测徽章。用户响应不含账号、凭据、上游地址、请求体或绝对流量计数。
- 缓存：只缓存与用户无关的基底（上架分组、target、T 时刻快照、tier5 序列、探测健康、最新报告、探测桶），按 (T, range, 算法版本) 键 30 秒，并发缺失串行加载一次；权限过滤、价格、排序、分页全部按请求计算，个性化响应不缓存。
- 公共设置新增 `provider_hall_enabled`（`GetPublicSettings`、注入 payload、`dto` 与 `setting_handler`），来源为大厅配置 `display_enabled`，通过 `SettingService.SetProviderHallDisplayReader` 注入并缓存 30 秒；未注入时恒为 false（opt-in 默认）。
- 就绪标志：`ProvideProviderHallQueryService` 接线时调用 `ProviderHallService.SetDisplayReady(true)`，管理员从此可以打开 `display_enabled`（仍要求 `collection_enabled`）。

## 源码

- 新增：`backend/internal/service/provider_hall_query.go`、`provider_hall_pricing.go`、`provider_hall_sort.go`；`backend/internal/repository/provider_hall_read_repo.go`；`backend/internal/handler/provider_hall_handler.go`；`backend/internal/server/routes/provider_hall_user.go`。
- 修改：`service/provider_hall.go`（`SetDisplayReady`、`DisplayEnabled`、测试接缝 `SetGroupAccessForTest`）、`service/wire.go`（`ProvideProviderHallQueryService` 及 ProviderSet 条目）、`repository/wire.go`（`NewProviderHallReadRepository`）、`service/setting_public.go` / `setting_service.go` / `settings_view.go`、`handler/dto/settings.go`、`handler/setting_handler.go`、`handler/handler.go`、`handler/wire.go`、`routes/user.go`；`cmd/server/wire_gen.go` 由 `go generate ./cmd/server` 重生成。
- 未改动：`handler/dto/provider_hall_user.go` 的字段与 JSON 标签（B7 合约），未新增字段。
- 测试：`service/provider_hall_query_test.go`、`handler/provider_hall_handler_test.go`、`repository/provider_hall_read_db_test.go`（以 `b5_read` 注册到 localdb 合约表）。

## 本轮验证（2026-09-12 实际执行）

1. `gofmt -l` 本批文件为空；`go build ./...` 通过；`go vet -tags=unit ./internal/... ./cmd/...` 通过；`go generate ./cmd/server` 成功。
2. `go test -tags=unit -race ./internal/service ./internal/handler/... ./internal/repository ./cmd/server -run ProviderHall -count=1`：全部包通过。本批新增 service 用例 16 个（子用例计入）：权限与可见性（普通/专属/订阅/未上架/停用/非 OpenAI 平台、详情与报告 404、开关关闭 404）、行内容（报价、单位、各指标、健康、徽章、迷你趋势）、两个专属倍率用户 100 次交替并发不串价且基底只加载一次（race 通过）、搜索/筛选/分页/参数校验、三级排序与缺失末尾、详情与报告、无水位时的空快照、预测倍率固定样例（阶梯 → n/a、写缓存 → n/a、缓存率不足 → n/a、参考缺失/为零 → n/a、参考过期 → stale、缓存率 stale → stale、倍率折算）、历史均价混合计费 → n/a、snapshot_id 往返。handler 用例 3 个：11 组非法参数均返回 400 `PROVIDER_HALL_INVALID_QUERY`、range 默认取配置、非法分组 404、展示守卫 404/放行、未登录 401、DTO 反射遍历三类响应共 16+ 个嵌套类型，JSON 键不含 `account|token|key|upstream|count`（白名单 `probe_tokens`、`bucket_seconds`）。
3. PostgreSQL 18 隔离集群（`PROVIDER_HALL_PG_BIN`，端口 15484/15485/15486）：`go test -tags=providerhall_localdb ./internal/repository -run ProviderHall -count=1 -v` 六条路径（empty、upgrade_228…upgrade_232）全部通过，共 180 个子用例，其中本批 2 个 × 6：造 100 个上架分组 × 2 个已启用 target × 121 个 tier5 快照窗口（30 天）；断言列表值与排序、30d/6h 迷你趋势缺口、详情 P90/趋势/档案、报告为空；通过记录语句的连接器证明冷缓存列表固定 9 条语句（100 组与 5 组相同），每条 SELECT 的 `EXPLAIN (FORMAT JSON)` 与语句文本均不含 `usage_logs` / `provider_hall_requests`；热缓存下第二个用户的请求不再读取 `provider_hall_snapshots`。
4. 未执行：Docker 不可用，`integration` 标签未跑（以 localdb 路径替代）；未做浏览器联调与截图（B8）；未跑全量 `make test-unit` / `test-integration`。

## 偏差与限制

- 预测倍率在缓存率 `stale` 时仍计算并标记 `stale`，而不是直接 `not_applicable`；`insufficient/incomplete/disabled` 仍为 n/a。理由：聚合每分钟运行，快照过期 5 分钟即 stale，若一律 n/a 会频繁闪烁。
- 详情面板的探测指标取该组所有已启用档案的探测样本（非仅默认档案）；P90 取默认档案快照。
- 分组健康在“部分 up、部分 unknown”时报 unknown，不报 up。
- 列表读取的基底覆盖全部上架分组（非仅当前用户可见的分组），随后再按权限裁剪；这使语句数固定并可共享缓存，但当上架分组数远大于用户可见数时会多读少量行。
- `ProviderHallQueryService.List` 每次请求仍要调用 `GetAvailableGroups`（Ent 多条语句）与 `GetConfig`，这两项不在固定语句计数内；用户倍率通过网关服务的既有缓存解析。
