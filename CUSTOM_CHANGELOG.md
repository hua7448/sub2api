# Sub2API 定制改动记录

本文档记录 `/home/ubuntu/sub2api-source` 中保留的本地功能、生产修复及源码归档。条目按时间倒序维护；以后每次更新必须同时记录需求、行为、源码范围、迁移/配置、验证和部署状态。

## 2026-09-12：供应商大厅 B2～B8 全部批次交付、验收后修复与部署包

状态：本地实施与验收完成，已打出 Linux 部署包（`sub2api-provider-hall-0.1.179-20260912.tar.gz`，含 amd64/arm64 二进制、校验值与 `DEPLOY_NOTES.md`），尚未部署生产；本机无 Docker，`integration` 标签以 PostgreSQL 18 localdb 路径补跑。

- 需求与行为：按 `docs/PROVIDER_HALL_DEVELOPMENT_PLAN.md` 与已确认的产品决策，完成真实请求采集与计费关联（迁移 231）、分钟聚合/快照/对账/保留（迁移 232）、探测与模型检测任务 Runner（迁移 233）、用户 API 与定价叠加、管理任务 API 与 jobs/health 页签、用户大厅 `/providers`、共享模拟上游、端到端与性能验收。三个运行开关由各自的就绪标志解锁，按采集 → 主动任务 → 用户大厅顺序开启。
- 验收后修复：全局配置页三个开关原为占位禁用，现按后端返回的 `readiness` 启用并提交实际值；`gateway_origin` 去掉末尾斜杠；`PROVIDER_HALL_INVALID_CONFIG` 按 `metadata.field` 给出具体字段说明；大厅配置保存后清除公开设置的展示开关缓存并让注入的 `index.html` 缓存失效（否则用户在重启前看不到入口）。端到端过程中另修复样本账单积压误暂停派发、两处 390px 溢出、表格列宽裁剪按钮。
- 源码范围：`backend/migrations/231..233`；`backend/internal/service/provider_hall_{tracker,collector,aggregator,snapshot,runner,probe_client,verification,budget,query,pricing,sort,admin}*.go` 及 `provider_hall.go`、`setting_public.go`、`openai_gateway_usage.go`、三个网关 handler 的采集挂点；`repository/provider_hall_{fact,aggregation,job,read,admin}_repo*.go`；用户与管理 handler/DTO/路由；`internal/testutil/fakeopenai`；前端 `views/user/ProviderHallView.vue`、`components/provider-hall/*`、`composables/useProviderHall.ts`、`api/providerHall.ts`、后台 jobs/health 组件、i18n、路由与侧栏；`frontend/scripts/provider-hall-e2e.mjs` 与 `lib/fake-openai-upstream.mjs`。
- 迁移/配置：新增表均为附加，不改旧表；新增环境变量 `INSTANCE_ID`（节点名）、`PROVIDER_HALL_ALLOW_LOOPBACK`（仅本地）、`PROVIDER_HALL_DISABLE_{COLLECTOR,AGGREGATOR,RUNNER}`；公开设置新增 `provider_hall_enabled`。
- 验证：Go 聚焦单测（含 race）全部通过；localdb 六条升级路径 180 子测试通过；perf 标签列表 P95 冷 43 ms / 热 14 ms，网关采集开销无可测增量；前端 lint/typecheck 零错误，Vitest 255/256（唯一失败为既有 `OAuthAuthorizationFlow.spec.ts`）；端到端脚本最终一轮 19/19 检查、28 张截图通过。后端全量 test-unit 的 9 个失败为既有基线。明细见 `docs/reviews/PROVIDER_HALL_B8_20260912.md` 及 B2～B8 各记录。
- 部署要点：运营用户需有余额；多实例每台唯一 `INSTANCE_ID` 并登记到预期采集节点；线上网关 Origin 填 `https://域名` 不带路径。

## 2026-09-11：恢复大厅实施基线的构建与回归

状态：本地源码修复与验证完成，未部署生产；供应商大厅采集、聚合、任务和用户展示仍按开发计划保持关闭，未宣称完整闭环已交付。

- 修复前端账户用量批量刷新辅助函数、测试弹窗音视频状态、注册域名配额设置读取，以及分组创建请求中未使用的定价辅助代码，恢复正式类型检查和构建。
- 保留 Provider Hall 现有配置/权限/算法实现及 `PROVIDER_HALL_NOT_READY` 启用保护；未修改 229/230 迁移、真实网关采集、计费或后台 worker 生命周期。
- 验证：Node 20/pnpm 9 前端 `vue-tsc`、正式 `vite build`、受影响 ESLint 及 10 个测试文件 87 项通过；Go 服务构建、ProviderHall 单元聚焦测试和 PostgreSQL 本地迁移测试通过。Docker 集成与完整 A～H 尚未执行。

## 维护规则

1. 应用源码的统一归档目录是 `/home/ubuntu/sub2api-source`。
2. 生产 UI 的锁定构建产物仍以 `/home/ubuntu/sub2api-production-ui-source/backend/internal/web/dist` 为准；不得把开发目录的 `dist` 直接用于生产。
3. 本仓库只保存可维护源码、通用 migration 和必要的操作留档，不保存生产二进制、fallback binary、systemd drop-in、UI guard、`node_modules`、构建缓存或部署备份。
4. 一次性数据修复必须放在 `docs/operations/<date>-<topic>/`，写明快照条件、是否可重跑和校验值，不能混入自动 migration。
5. 每次发布后补充测试结果、生产版本和回滚点；未验证或未部署的内容必须明确标记。

## 2026-09-11：修复供应商大厅后台配置的两项 P2 问题

状态：两项修复及聚焦回归完成，仅更新独立本地测试实例，未部署生产；完整大厅的后续分期仍未完成。

- 分组展示名称、说明的 Ent 校验改为 `MaxRuneLen(100/2000)`，与 service 的字符限制一致；前端说明输入上限由 1000 调整为 2000。执行 `go generate ./ent`；现有 PostgreSQL varchar 已按字符计数，无新增迁移，未修改 229/230。
- 全局配置加载期间禁用编辑、保存及重复加载，提交处理函数也检查忙碌状态；收到保存结果时取消旧配置请求。配置和模型档案改为独立加载，避免配置重载覆盖刚保存的档案。
- 增加 service 中文长度边界测试、真实数据库中文创建/更新/超限回滚检查，以及实际父子组件的 4 项延迟响应测试；扩展现有 Playwright 脚本验证桌面/移动端加载互斥和中文保存边界。
- 本轮通过：大厅 Go 聚焦竞态测试、后端构建；PostgreSQL 空库及从 228/229 升级的三条路径；前端 17 项测试和涉及文件 ESLint；真实本地浏览器验收及 24 张截图。
- 全量前端类型检查仍被既有账户、分组和注册页面错误阻断，本次修改文件未报错。未重跑全量后端、完整 CI 或完整大厅 A～H 验收；未构建或替换嵌入 UI。

详细变更、测试命令和截图见 `docs/reviews/PROVIDER_HALL_P2_FIXES_20260911.md`。2026-09-10 条目保留为历史验证记录。

## 2026-09-10：供应商大厅目标绑定与后台配置页面

状态：本轮源码和配置闭环验证完成；完整采集、调度、报告和用户大厅尚未完成，未部署。

- 新增前向迁移 `230_provider_hall_targets.sql`，创建目标表和永久探测 Key 来源登记表；保留已应用的 229 迁移。更新 Ent schema、生成代码及依赖装配。
- 新增管理目标集合读写接口，使用分组级版本 CAS；在事务内检查 Key、运营用户、分组权限、订阅和混合路由。未提交的已有目标停用保留 ID，失效 Key 仍可停用；变更运营用户会停用相关目标并使旧配置版本失效。
- 专用 Key 首次登记仅接受未使用 Key，登记记录不保存凭据，并在删除 Key/用户/分组后保留。V2 四类聚合排除登记后的探测请求，保持普通流量原有公式。
- 后台新增供应商大厅配置页，包含全局配置、模型档案、分组展示与目标绑定。复用当前组件、API 客户端和权限边界；同步中英文，保留冲突输入，取消过期请求。采集、展示、主动任务仍关闭且不可启用。
- 用 Node 20/pnpm 9 加入 Playwright 1.58.2 及管理页真实接口验收脚本，锁文件仍为 9.0；没有构建或替换嵌入 UI。
- 验证：Go 1.26.5 大厅/V2 聚焦竞态检查、Go 1.26.6 service/admin 大厅检查、后端构建通过；PostgreSQL 18.6 空库/从 228 升级/从 229 升级及每条路径的 5 组目标/来源合约通过；前端 13 项 API/组件测试、涉及文件 ESLint、frozen-lockfile 安装通过。
- 浏览器：独立 DATA_DIR、PostgreSQL、Redis，经正常 AUTO_SETUP 初始化且 RUN_MODE=standard；真实接口完成配置保存、建档、上架和 Key 绑定。390/1440/1920 像素布局与明暗样式、弹窗、键盘和溢出检查通过，保存 22 张最终截图。新增页面与弹窗已修复移动端溢出和深色对比问题；没有发出付费探测。
- 限制：全量前端类型检查仍被原有账户/分组/注册模块错误阻断；完整 A～H、全链路网关和性能验收未完成，不宣称 CI 或生产验证通过。

源码范围、复现命令、截图和剩余工作见 `docs/reviews/PROVIDER_HALL_TARGETS_ADMIN_20260910.md`。

## 2026-09-10：开始供应商大厅第一批实施

状态：已完成配置、权限和算法基础的部分源码及本地验证；大厅完整功能和第一批剩余项目尚未完成，未部署。

- 新增 ProviderHall service/repository/admin handler/DTO、管理路由及 Wire 装配；配置、原分组上架信息及模型协议档案通过 Ent schema 生成。
- 新增前向迁移 `229_provider_hall_foundation.sql`，创建三个配置表。金额保持 decimal 字符串与数据库 numeric 精度，写入执行版本 CAS，原分组软删除撤销读取、硬删除级联删除上架配置。
- 采集、展示、任务开关默认全部关闭；当前代码拒绝启用未接入的运行功能，返回 `PROVIDER_HALL_NOT_READY`。没有改动真实网关、扣费路径或后台 worker 生命周期。
- 建立复用 API Key 可用分组规则的授权方法，以及 TTFT、缓存率、重试感知成功率、输入费用/均价及上海预算日基础算法。用户接口和页面尚未开放。
- 前端新增管理 API 类型/调用与 3 项 Vitest 测试，未修改锁文件、页面、中英文文案或嵌入 UI。
- 验证：Go 1.26.5 大厅竞态测试含子用例 32 项通过；Go 1.26.6 聚焦兼容检查通过；本地 PostgreSQL 18.6 空库/升级库、重复迁移及 repository 合约通过；后端构建、前端新增文件 ESLint 和 3 项 API 测试通过。
- 限制：后端全量单元测试、前端全量类型检查存在未改动模块失败；Docker 集成聚焦命令被原有测试编译错误阻断，本机也没有 Docker。未完成 A～H 全量验收，不宣称 CI 或生产验证通过。

详细源码范围、API、复现命令及未完成清单见 `docs/reviews/PROVIDER_HALL_FOUNDATION_20260910.md`。

## 2026-07-24：源码集中归档

状态：已把当前生产前端和后端中缺失的可维护代码合并回本仓库；本次归档本身不执行生产部署。

来源：

- UI v10 源码：`/home/ubuntu/sub2api-production-ui-source/frontend`
- 当前后端源码：`/home/ubuntu/sub2api-bnode-source/backend`
- 生产 UI 版本：`registration-alias-toggle-20260724-v10`
- v10 UI tree SHA-256：`5b0b7168e5ecb64dba4f749b35e1783e45b2a62ec1d053cd82c98fbf981fbf40`
- v10 后端候选 binary SHA-256：`6945a71d21a50bfa5eebdc5a36fec49d4de9c40b3e42bfbed193111251b3ae7f`
- 回滚版本：v9，binary SHA-256 `d12077e18c4765151983cdebd1348a9c46e809883c4d094a5e6716aae97f362c`

明确排除：生产 `dist`、锁标记、guard、installer、systemd 配置、二进制、`node_modules`、TypeScript 构建缓存和 `*.backup-deploy-*`。

本次归档验证：

- 前端 scoped ESLint 和 typecheck 通过；7 个聚焦测试文件共 59 项全部通过。
- 前端完整 Vitest 共 1,238 项，通过 1,230 项；8 项既有失败位于本次未修改的账户弹窗、倍率单元格、Key 使用弹窗、Spark Shadow、分组列设置、Stripe lazy loading 和系统日志测试。
- 后端 migration 185、别名邮箱、邀请返利、repository SQL、handler 和订阅履约测试全部通过。
- `git diff --check` 通过；没有生成或复制前端 `dist`、后端 binary 或生产部署文件。

## 2026-07-24：禁止服务商别名邮箱注册（UI v10）

需求：恢复后台“禁止别名邮箱注册”按钮，并保留原有深色主题和生产 UI。

行为：

- 管理入口：系统设置 -> 安全与认证 -> 注册设置 -> 禁止别名邮箱注册。
- 设置键：`registration_email_provider_alias_disabled`，升级默认值为 `false`；生产数据库当前值为 `true`。
- 开启后拒绝 Gmail 点号地址、`googlemail.com`、QQ 非数字别名和 Foxmail 别名。
- 覆盖普通注册、验证码发送、邮箱绑定、OAuth 首次注册和 OAuth 补邮箱；已有账号登录不因该策略受阻。
- 后端返回结构化原因 `EMAIL_ALIAS_NOT_ALLOWED`，注册页和邮箱验证页显示本地化错误。

主要源码：

- 后端：`backend/internal/service/registration_email_policy.go`、`auth_service.go`、`auth_email_oauth_auto.go`、`auth_oauth_email_flow.go`、设置 service/handler/DTO 及对应测试。
- 前端：`frontend/src/views/admin/SettingsView.vue`、`frontend/src/api/admin/settings.ts`、`frontend/src/utils/authError.ts`、注册/验证页面、中英文文案及测试。
- 嵌入 UI revision：`backend/internal/web/embed_on.go` 更新为 `registration-alias-toggle-20260724-v10`。

验证：前端 27 个聚焦测试通过，`vue-tsc`、相关 ESLint 通过；Playwright 已检查桌面/移动端的明暗主题、开关状态、溢出和文本重叠；后端带 `unit` 标签的别名策略测试通过。

部署记录：v10 已发布到 5 个生产节点，并验证 `panel.sharesai.xyz`、`api.sharesai.xyz`、`bnode.sharesai.xyz`、`cn2-api.sharesai.xyz` 返回相同 revision 和关键静态资源。

## 2026-07-22 至 2026-07-24：禁止 +tag 邮箱注册

该策略与服务商别名策略独立：

- 设置键：`registration_email_plus_alias_disabled`。
- 缺少设置时按 `true` 处理，默认拒绝本地部分包含 `+` 的注册邮箱。
- 后端错误原因：`EMAIL_PLUS_ALIAS_NOT_ALLOWED`。
- 公开设置、管理员设置 API、审计日志、普通注册、OAuth 首次注册和邮箱绑定均接入该策略。
- 当前生产数据库值为 `true`。
- 当前 UI 不增加第二个可见开关；服务商别名按钮仍只控制 provider-specific alias，避免把两个策略混称。

本次归档同步了相关 handler、DTO、service、公开设置和带 `unit` 标签的测试，不复制生产构建产物。

## 2026-07-22 至 2026-07-24：邀请返利与最新版 UI

后端行为：

- 新增 migration `185_affiliate_redeem_code_sources.sql`，给返利账本增加 `source_redeem_code_id` 以及余额/订阅兑换来源的唯一防重复索引，并尽力回填可识别的历史来源。
- 余额兑换码和所有有效订阅兑换码均可产生对应返利；账本记录 reward type、source type、兑换码、分组和返利天数。
- `qualified` 有效客户以使用正额余额码或非退款订阅码为准，免费/管理员赠送不计入。
- 门槛统计排除正在发生返利的当前客户。默认门槛为 1，含义是先积累 1 个有效客户，从第 2 个客户开始返利。
- 订阅天数返利遵守邀请关系有效期，并通过来源兑换码保证幂等。
- 返利查询同时展示余额返利和订阅天数返利，不再只依赖支付订单。

前端 UI：

- 邀请返利页升级为 v9 最新版：显示有效客户/总邀请、门槛进度、可用与冻结余额、订阅天数池、返利记录和提取记录。
- 返利记录显示兑换/支付来源、分组、余额金额或订阅天数。
- 邀请客户显示资格状态，支持搜索；桌面和移动端均有独立布局并支持深色主题。
- API 类型增加 `qualified`、`ledger_id`、`reward_type`、`source_type`、`redeem_code_id`、`group_id`、`group_name`、`rebate_days`。

一次性历史补发：

- 只读核对和一次性执行脚本归档在 `docs/operations/2026-07-22-affiliate-rebate-backfill/`。
- `apply` 脚本绑定当时快照的 65 条缺失订阅返利、合计 195 天；不能作为通用 migration 或直接重跑。
- 生产执行结果：2026-07-22 10:25:13 UTC 写入 65 条 `accrue_days` 记录、合计 195 天，涉及 24 位邀请人；migration 185 仅执行一次。
- 通用 schema 与幂等逻辑只由 migration 185 和应用代码维护。

## 2026-07-22：Primary/API 节点角色隔离

- 新增 `INSTANCE_ROLE=primary|api` 配置和校验。
- API 节点跳过 schema migration、默认分组 seed 和仅应运行一次的全局后台 worker。
- Primary 节点负责数据库迁移、种子数据和全局任务，避免多节点重复执行。
- 相关源码位于 `backend/internal/config`、repository 初始化、service wiring 和启动测试；本仓库此前已经保留，集中归档时未覆盖。

## 2026-07-18：Codex/OpenAI 账号可靠性

- 导入或创建 OpenAI 账号时强制接入图片生成 bridge。
- 兼容 legacy PAT 标识及没有 `expires_at` 的 PAT。
- 普通、可能瞬时的 401 不再立即把账号标记为永久错误；明确不可恢复的 Unauthorized 仍按错误处理。
- 包含账号导入、OAuth、图片 bridge、限流和 401 判定的对应测试。

## 2026-07-17 至 2026-07-24：UI 与嵌入资源修复

- 保留最新版邀请返利 UI、深色主题和响应式布局。
- 修复浅色主题中 danger 确认按钮的可见性。
- 修复 K12 计划筛选和平台徽标显示，并增加组件测试。
- 嵌入子应用避免 301 循环。
- 使用 revision cookie、`no-store` 和必要的 `Clear-Site-Data` 处理旧静态资源，降低更新后白屏或错版概率。
- 生产 UI 已从 v9 升级并锁定到 v10；后端专用更新不得重建或替换生产前端。

## 后续条目模板

```text
## YYYY-MM-DD：改动标题

状态：开发中 / 已验证 / 已部署
需求：
行为：
源码：
迁移或配置：
验证：
部署与回滚：
```
