# 供应商大厅 (Provider Hall) 分支知识库

> 更新: 2026-09-18,由 Claude Code 从开发过程记忆整理。
> 本文件是本分支(供应商大厅版本)的权威背景文档;dev-memory/ 下其他文件是早期 bmai 项目遗留,已过时,仅供历史参考。

## 1. 这个分支是什么

- **仓库**: `/Users/simple/smartq/sub2api`,origin = `github.com/hua7448/sub2api`(另有 upstream = Wei-Shaw/sub2api)
- **分支**: `release/v0.1.179-provider-hall`,**孤儿分支(orphan)**,与 main 无共同历史
- **首个提交**: `fede1b498`(2026-09-18 从 `/Users/simple/smartq/sub2api-clean-source-baibai` 全量导入)
- **基线**: 上游 Wei-Shaw/sub2api v0.1.179 + 供应商大厅全套改动(B2–B8)
- **与 main 的关系**: main 是 smartapi 版(上游 v0.1.169-smartapi.1,2026-08-02 同步),两边上游基线已岔开。**约定不合并**;确需合并时用 `git merge --allow-unrelated-histories`(预期冲突量大)
- **检出开发**(与 main 目录并存):
  ```sh
  cd /Users/simple/smartq/sub2api
  git worktree add ../sub2api-aihub release/v0.1.179-provider-hall
  ```
- **导入排除项**: `frontend/node_modules`、`backend/bin`(编译产物)、`docs/reviews/artifacts`(20MB e2e 测试产物,留在 `/Users/simple/smartq/sub2api-clean-source-baibai/docs/reviews/artifacts/`)
- **源目录不是 git 仓库**: `/Users/simple/smartq/sub2api-clean-source-baibai` 是原始开发现场(无版本记录),以后开发请在本分支的检出里做

## 2. 供应商大厅功能状态(截至 2026-09-12)

- **批次 B2–B8 全部完成并本地验证**: B2 采集 → B3 聚合 → B4 调度 → B5 用户 API → B6 管理端 → B7 用户页 → B8 e2e/性能
- **迁移**: 229–233(231 facts / 232 metrics / 233 jobs;测试环境验证过 17 张 provider-hall 表)
- **用户页**: `/providers`(设计稿在 `docs/screenshots/hubScreen/`)
- **用户 API DTO 契约**: `backend/internal/handler/dto/provider_hall_user.go`(snake_case tag,**前端依赖此契约,勿破坏**)
- **验收后修复(9/12)**: 管理端配置表单开关由 `readiness` 驱动;保存触发 `SettingService.NotifyProviderHallConfigChanged`(清显示缓存 + 失效注入 index.html);无效配置报错带字段名
- **关键文件**:
  - `backend/internal/service/provider_hall_aggregator.go`、`provider_hall_runner.go`(有 `Ready() bool` 就绪缝)
  - `backend/internal/repository/provider_hall_localdb_test.go`(localdb 测试 harness)
  - `cmd/server/wire.go` 的 ProvideProviderHallService:Collection = aggregator.Ready(),Tasks = collection && runner.Ready();清理顺序 Runner → Aggregator → Collector(同一 "ProviderHall" step)
- **每批记录**: `docs/reviews/PROVIDER_HALL_B{2..8}_20260912.md`(+ `B8_E2E`、`B8_PERF`、`FOUNDATION`、`P2_FIXES`、`TARGETS_ADMIN`)
- **总计划**: `docs/reviews/PROVIDER_HALL_PLAN_20260911.md`(B2→B8 批次计划 + 产品契约,2026-09-11 定稿)

## 3. 产品契约(2026-09-11 已与需求方确认,勿重新询问)

- 群组模型集 = 启用的 targets(group×profile);行级 TTFT/P90/成功率**合并所有启用 profile 的真实流量**(snapshot `profile_id=0`);缓存率/真实价/预测倍率用群组默认 profile;每个 target 一个健康点
- 参考价 + 参考缓存率:管理端维护在 profile 上;**不做** OpenRouter 集成
- 模型验证 = 三断言套件(算术×3、JSON×3、`supports_tools` 时强制工具调用×3)+ 响应模型名检查(不是截图 tooltip 写的 "Juice fingerprint")
- 细节指标(6h e2e 可用率、TPS、探针 TTFT vs e2e、默认 profile P90、最近探针 token 出入)全部来自探针样本;**绝不向用户暴露真实流量绝对数**
- 既定假设:node_id = `INSTANCE_ID` 环境变量否则 hostname;JSON 套件容忍 ``` 代码栅栏;模型/协议过滤只影响哪些群组出现

## 4. 本机(Mac)开发/测试环境

- **Docker 在这台 Mac 上不可用** → `integration` build tag 跑不了,一律走 localdb 路径
- **localdb harness**(2026-09-12 约定):
  - `PROVIDER_HALL_PG_PORT`(默认 15479)让多套件并行
  - 升级模式从迁移文件自动发现(`upgrade_228..upgrade_<newest-1>`)
  - 批次契约用 `registerProviderHallLocalDBContract(name, fn)` 在 `*_db_test.go` 的 `init()` 注册(tag `integration || providerhall_localdb`)
  - `PROVIDER_HALL_LOCALDB_MODES=empty` 收窄模式
- **e2e**: `pnpm --dir frontend run test:e2e:provider-hall`(即 `node scripts/provider-hall-e2e.mjs`);管理端冒烟 `test:e2e:provider-hall-admin`;性能测试 `-tags=perf`
- **版本注意**: PostgreSQL 18 在 `/opt/homebrew/opt/postgresql@18/bin`;本机 Node 24 / pnpm 10,项目目标 Node 20 / pnpm 9(勿随手重新生成 lockfile)
- **工具链冲突待解**:`backend/go.mod` 声明 Go 1.26.5,CI 断言 1.26.6,Dockerfile 用 1.26.6,生产契约脚本钉 Linux Go 1.26.5——无关工作时不要顺手统一

## 5. 已知基线问题(与大厅无关,勿浪费时间排查)

- 前端 vitest 基线失败 1 个:`OAuthAuthorizationFlow.spec.ts`
- Go 单测基线失败 9 个(上游/环境原因)
- 短 e2e 只会给单点趋势;探针/验证任务需要操作员账户有余额

## 6. 服务器测试部署(179.255.156.111,SSH root 免密)

- **部署目录** `/opt/sub2api-provider-hall-test`:app + postgres + redis 三容器(compose),仅 app 发布端口 `127.0.0.1:18082 → 容器 8080`;postgres/redis 不发布宿主端口
- **服务器上没有源码**,只有二进制部署;镜像 `sub2api-provider-hall:0.1.179-20260912`;原始包 `/root/sub2apiAIHUB_TEST/sub2api-provider-hall-0.1.179-20260912.tar.gz`
- **compose 里有 `PROVIDER_HALL_ALLOW_LOOPBACK: "1"`**(2026-09-15 加,允许测试网关走 loopback)
- **大坑**:provider-hall 配置 `gateway_origin` 必须是 `http://127.0.0.1:8080` —— 容器内访问不到宿主发布的 18082(症状:探针全部 `transport_refused`、HTTP status 0)
- **测试库**:PostgreSQL 库/用户 `sub2api_test`,volume `sub2api-provider-hall-test-pgdata`;应用配置 `data/config.yaml`;**真实密钥只在服务器**(`data/config.yaml`、`credentials/` 均 root-only),仓库里是脱敏的 `deploy/provider-hall-test/config.example.yaml`
- **测试资金**(2026-09-15 注):操作员 user 1 + 探针 API key 1 各充值 100(此前 403 INSUFFICIENT_BALANCE);ordinary key 2 用于刷真实用户指标
- **当前开关**:`collection_enabled=true`、`display_enabled=true`、`tasks_enabled=false`(9/16 关定时探测防继续花钱);手动探测/验证需临时打开 tasks
- **踩坑清单**(9/17 整理自运维记录):
  - 改容器运行时环境变量必须**重建容器**,`docker restart` 不会生效
  - PostgreSQL 18 的 volume 必须挂 `/var/lib/postgresql`(挂旧路径 `/var/lib/postgresql/data` 起不来)
  - 拷贝安装标记/配置文件会**错误地跳过初始化**;要走真实的 AUTO_SETUP/向导 + API 登录验证
  - loopback 网关 origin 的校验需要 `PROVIDER_HALL_ALLOW_LOOPBACK=1`;生产形态 origin 只应是 `https://域名`(不带路径)
  - 用户侧 TTFT/平均延迟卡片需要**普通用户网关流量且记录了 ttft_ms**,只跑探针/验证任务刷不出来
  - 后续定时验证任务会独立更新"最近验证"状态,可能覆盖已填充的真实用户指标卡的 verified 展示
  - 手动采样被中断会留下未派发的 pending 样本,造成 billing backlog,需标记 unbilled 清理
- **完整运维史**: `deploy/provider-hall-test/DEPLOYMENT-RECORD.md`(9/12 隔离初始化 → 9/15 loopback 修正+注资 → 9/16 展示样本+TTFT 补样 → 9/17 坑点整理收尾)
- **过时快照**: `/root/simple/sub2api-clean-source-baibai`(服务器上 9/10 的旧拷贝),别当源码用
- 本地配套:`deploy/provider-hall-test/`(Dockerfile、compose.yaml、README、backup.sh、prepare-bootstrap.cjs、verify-login.cjs)

## 7. 维护约定

- 新批次/部署/产品决策 → 更新本文件 + 在 `docs/reviews/` 加 `PROVIDER_HALL_<主题>_<日期>.md` 记录
- 与服务器同步:改完代码在本分支提交;服务器只吃发布包(参照 `deploy/provider-hall-test/README.md`)

## 8. 管理端配置界面重做(2026-09-18)

**动的三个页面**: `components/admin/provider-hall/` 下 `ProviderHallGroups.vue`(改为 `TablePageLayout` + DataTable 可展开主表)、`ProviderHallProfiles.vue`(接上游候选批量建档)、`ProviderHallConfigForm.vue`(加状态摘要与分区)。

- **分组目标改为展开行**,不再是逐组弹窗;草稿按 `group_id` 隔离,保存粒度是「保存本组」。上架/展示名/描述/排序值拆到 `#hall-listing` 弹窗。
- **`payloadFor` 曾恒发 `version: 0`** → 后端乐观并发校验必然 409。分组版本现在存 `targetVersion[group_id]`,由最近一次 `getTargets`/`saveSettings` 响应提供。改这块时别再把版本丢掉。
- **`display_order` 现在真的生效**: `ListAdminGroups` 支持 `sort=display_order`(未上架排最后),前端有拖拽排序弹窗,保存走批量接口。
- **档案候选取全局**: 新增 `GET /profile-candidates`(跨组按模型+协议合并,带 `sources`/`groups`);原按组的 `GET /groups/:id/models` 保留。新增 `DELETE /profiles/:id`,被目标引用时 409 `PROVIDER_HALL_PROFILE_IN_USE` 并带引用分组。注意:省略 target 只是置 `enabled=false`,行仍在,所以引用不会因「移除目标」自动解除。
- **DataTable 新增 `expandable`**(可选):桌面 `tr[data-expanded-for]`、移动卡片内展开,支持受控 `expandedKeys`。不动其他页面行为。
- **降智探测还没做**,仓库里没有任何实现。可复用的底子是 `service/provider_hall_verification.go` 的 `ProviderHallBuildSuite`/`ProviderHallVerdict` 断言框架;每次探测真花钱,需要单独一轮设计。生产 `tasks_enabled` 目前关闭。
- 记录: `docs/reviews/PROVIDER_HALL_ADMIN_UX2_20260918.md`、`CUSTOM_CHANGELOG.md` 同日条目。
