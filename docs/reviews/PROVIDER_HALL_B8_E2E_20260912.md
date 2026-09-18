# 供应商大厅 B8 全流程验收记录（端到端部分）

日期：2026-09-12。本记录只覆盖 B8 的端到端验收：隔离实例、真实网关流量、主动任务、用户页面浏览器操作与截图。性能与全量回归（`provider_hall_perf_test.go`、`TestProviderHallGatewayOverhead`、`make test-unit` 等）另见 `PROVIDER_HALL_B8_PERF_20260912.md`。未部署生产。

## 交付物

- `frontend/scripts/provider-hall-e2e.mjs`，npm 脚本 `test:e2e:provider-hall`（`pnpm --dir frontend run test:e2e:provider-hall`）。脚本自建并销毁隔离实例，不依赖已有服务；只绑定 127.0.0.1。
- `frontend/scripts/lib/fake-openai-upstream.mjs`：故障注入除请求头外新增提示词内嵌形式 `[[fake:status=500]]`，因为本站网关不透传自定义请求头。自检脚本仍通过。
- 结果目录 `docs/reviews/artifacts/provider-hall-e2e-20260912/`：`summary.json`（步骤、计数、指标快照、任务详情、截图清单）、26 张用户页截图、`admin-smoke/` 36 张后台截图与 `summary.json`、`db-diagnostics.txt`（实例关闭前导出的大厅各表汇总）、后端/前端/Redis/PostgreSQL 日志。

## 实例与流程（脚本实际执行）

环境：Go 1.26.5、Node 24.15.0、pnpm 10.33.2（项目声明 Node 20/pnpm 9）、PostgreSQL 18（`/opt/homebrew/opt/postgresql@18/bin`，端口 15490）、Homebrew Redis（端口 16390）、后端 127.0.0.1:8091、Vite `--mode provider-hall` 127.0.0.1:3100、模拟上游 127.0.0.1:9450、Playwright 1.58.2 + 本机 Chrome（`PROVIDER_HALL_BROWSER_CHANNEL=chrome`）。后端从源码 `go build` 得到，`RUN_MODE=standard`，`AUTO_SETUP=1` 正常初始化（没有手写 `.installed`），`INSTANCE_ID=hall-e2e-node`，`PROVIDER_HALL_ALLOW_LOOPBACK=1`。

1. 管理员：给运营用户（管理员本人）充值 100；创建 `gpt-4.1-mini` 的 responses（`supports_tools`）、chat_completions、messages 三个档案并填写参考价（0.40/0.10/缓存率 0.30）；创建两个 openai 分组（倍率 1 与 1.5，开启 Messages 入口），各绑定一个指向模拟上游的 API Key 账号，各创建一个专用探测 Key；写入配置（默认档案、`gateway_origin` 指向本实例、`expected_nodes=[hall-e2e-node]`、预算 5），上架两组并启用 6 个目标；开关顺序 collection → tasks → display。
2. 普通用户：每个分组各一个 Key。等待采集器 30 秒目标索引刷新后，通过本实例网关发送流量：三个协议各 320 条（Alpha 250 + Beta 60 + 10 条上游断流；每 10 条 1 条非流式），随后再发 12 条上游 HTTP 500 注入请求，并等待账号从网关冷却中恢复。
3. 管理 API 手动发起探测与检测（含幂等键复用断言），轮询到任务结束。
4. 在 display 仍关闭时运行 B6 扩展后的后台冒烟脚本（`PROVIDER_HALL_ENABLE_TASKS=1`），随后开启 display。
5. 轮询用户 API 直到 Alpha 的 TTFT/历史价格状态为 `ok`，再等到第一个完整 5 分钟桶发布（迷你图有点）。断言权限、指标、排序、详情、无效参数 400、响应中不含 `account_id/api_key/base_url/upstream/credential` 键。
6. Playwright 以普通用户打开 `/providers` 执行页面操作与截图，最后清理进程和临时目录。

单次完整运行约 4.5 分钟（`summary.json` 的 `started_at`/`finished_at`：04:51:43Z → 04:56:16Z）。

## 结果（最终一次运行，`summary.json`）

| 项目 | 结果 |
|---|---|
| 网关流量 | responses/chat_completions/messages 各 320 条，均 200；注入断流 30 条均 200（见下文说明）；注入上游 500 的 12 条：2 条 502（换号失败）、10 条 503（账号冷却，无可用账号） |
| 事实表 | 用户成功行 961，失败行 13（`db-diagnostics.txt`）；缓存读 token 精确等于输入的一半 |
| 手动探测 | 任务 `succeeded`，样本 passed，TTFT 3 ms，计费 `confirmed` |
| 手动检测 | 任务 `succeeded`，判定 `passed/completed`，算术 3/3、JSON 3/3、工具 3/3、模型名匹配 9/9 |
| 用户列表 | 3 个上架分组可见（含冒烟脚本上架的无账号分组）；Alpha：TTFT 最快95% 1.397 ms、P90 3 ms、缓存率 0.5、成功率 0.9836、历史价 $0.20/M、预测倍率 0.8065，全部 `ok`；Beta 倍率 1.5、成功率 1；`rate:desc,ttft_fast95:asc` 排序 Beta 在前；`range=99h` 返回 400 |
| 详情 | 探测首 Token/端到端 3 ms、默认档案 P90 7 ms、最近探测 Token 输入 12/输出 5；6h 端到端可用率与 TPS 为 `insufficient`（样本不足 12 / 无足够探测），趋势 72 点中真实 1 点、探测 1 点 |
| 浏览器 | 19 项检查全部通过：行内真实指标（TTFT 数值、价格、环形值、三档案健康点、检测徽章）；倍率芯片 off→desc→asc 与 URL/首行同步；自定义三级排序对话框写入 `sort=success_rate,...` 并持久化到 `localStorage`；24h/6h 时间维度与搜索过滤；展开行显示指标网格、趋势 canvas、模型健康 chips；检测报告弹窗；"使用此分组"创建新 Key（归属 Alpha）与切换到 Beta（`GET /api/v1/keys` 校验）；键盘 Enter 展开、↑↓ 移动焦点、Space 收起；无浏览器异常 |
| 截图 | 390/1440/1920 × 明暗 × 列表/展开详情/使用分组弹窗/报告弹窗 24 张 + 排序对话框、展开、报告、创建 Key 4 张；每张断言 `scrollWidth <= innerWidth` 且无按钮文字截断 |
| 后台冒烟 | 11 项检查通过，含主动任务开启后"立即探测"返回 202、任务/运行状态页签与 390/1440/1920 明暗截图 |

## 发现并修复的问题

1. `backend/internal/repository/provider_hall_job_repo.go`：收到明确非 2xx 应答的样本此前要等 10 分钟才被记为 `unbilled`，期间账单积压规则会让所有派发暂停 10 分钟；改为 90 秒宽限（`providerHallDefinitiveBillingGrace`），`provider_hall_job_db_test.go` 增加对应用例。
2. `frontend/src/components/admin/provider-hall/ProviderHallHealth.vue`、`ProviderHallConfigForm.vue`：网格缺少基础 `grid-cols-1`，390 像素下运行状态页横向溢出（冒烟脚本捕获）。
3. `frontend/scripts/provider-hall-admin-smoke.mjs`：B6 扩展页签后档案弹窗循环没有切回"模型档案"页签，导致找不到编辑按钮。
4. `frontend/src/components/provider-hall/HallTable.vue`：列宽在 1080px 最小宽度下最后两列（模型检测、使用分组）放不下按钮，"使用此分组"被裁切；重新分配列宽。
5. `frontend/scripts/lib/fake-openai-upstream.mjs`：`setOptions` 是整体替换，脚本每次切换故障都重述基础选项（否则缓存 token 归零）；新增提示词内嵌故障标记。

## 观察到但未改动的行为（写入记录供产品/后续判断）

- 运营用户余额为 0 时探测请求被网关以 403 `INSUFFICIENT_BALANCE` 拒绝，任务标 `http_403`。Runner 的发出前检查与运行状态页都不校验运营用户余额，部署前必须给运营用户充值。
- API Key 账号的 Responses 请求经网关转换为上游 chat_completions；上游在终态事件前断流时，转换层会补发 `response.completed` 与 `[DONE]`，网关返回 200 且无错误，大厅按成功、零用量记录。因此本轮"断流"注入没有产生失败样本，失败样本改由上游 500 注入产生。
- 上游连续 500 后网关把该账号短暂冷却（本轮约 1～2 分钟），同组无其他账号时后续请求 503；大厅按"无账号"记为失败（提交 0 次）。
- 迷你图与趋势只包含 T 之前的完整 5 分钟桶，最近不足 5 分钟的数据不显示；短时运行只能得到 1 个真实点与 1 个探测点，折线需要至少 2 点才绘制。
- Playwright 整页截图会把吸顶头部与侧栏重复渲染到页面中部，是截图方式的产物，非页面缺陷。深色截图通过切换 `.dark-theme` 类完成，应用主题 store 仍锁定浅色。

## 实际执行的命令

```sh
pnpm --dir frontend run test:e2e:provider-hall        # PROVIDER_HALL_BROWSER_CHANNEL=chrome，共运行 19 次，最终 2 次通过
node frontend/scripts/lib/fake-openai-upstream.selfcheck.mjs
pnpm --dir frontend exec eslint scripts/provider-hall-e2e.mjs scripts/provider-hall-admin-smoke.mjs scripts/lib/fake-openai-upstream.mjs src/components/admin/provider-hall/ProviderHallHealth.vue src/components/admin/provider-hall/ProviderHallConfigForm.vue src/components/provider-hall/HallTable.vue
pnpm --dir frontend exec vitest run src/components/provider-hall src/components/admin/provider-hall src/views/admin/__tests__   # 4 文件 23 用例、36 文件 205 用例通过
PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin PROVIDER_HALL_PG_PORT=15496 go -C backend test -tags=providerhall_localdb ./internal/repository -run ProviderHall -count=1 -v   # 六条迁移路径 180 个子测试通过
go -C backend test -tags=unit ./internal/repository ./internal/service -run ProviderHall -count=1
go -C backend vet -tags=unit ./internal/repository && go -C backend vet -tags=providerhall_localdb ./internal/repository
```

## 未执行或未验证

- Docker 不可用，`integration` 标签未运行；数据库合约以 localdb 路径补跑。
- 6h 端到端可用率与 TPS 因短时运行样本不足，只验证了 `insufficient` 状态与占位文案，没有验证有值时的渲染。
- 迷你图/趋势只出现 1 个真实点、1 个探测点，折线本身未绘制（缺口区渲染已验证）；多点折线仍需更长时间运行验证。
- 未验证 WS 入口、OAuth/Codex 账号路径、订阅计费分组、专属分组不可见等权限场景在浏览器中的表现（用户 API 层已由 B5 单测覆盖）。
- 未在真实上游或多实例（多个 `expected_nodes`）下运行。
