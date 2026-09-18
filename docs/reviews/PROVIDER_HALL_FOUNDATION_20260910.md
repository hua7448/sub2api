# 供应商大厅第一批实施记录

日期：2026-09-10。对应 [开发计划](../PROVIDER_HALL_DEVELOPMENT_PLAN.md)。

## 当前交付范围

本次开始实施第一批中的配置、权限与固定算法部分，不是完整大厅交付，也不宣称第一批全部完成。

- 新增 `ProviderHallService`、Ent repository、管理 handler、DTO、路由及 Wire 装配。
- 新增三种 Ent schema：`ProviderHallConfig`、`ProviderHallGroup`、`ProviderHallProfile`，已运行 `go generate ./ent`；Wire 源定义更新后已运行 `go generate ./cmd/server`。
- 新增前向迁移 `229_provider_hall_foundation.sql`，只创建 `provider_hall_config/groups/profiles`。默认不创建任何上架条目和模型档案，不复制原分组。
- 三个开关默认关闭。当前实现拒绝启用任何开关，返回 `409 / PROVIDER_HALL_NOT_READY`，避免在采集、任务或快照未实现时误开放功能。
- 复用 API Key 服务的可用分组规则，建立详情/报告共享授权方法，检查上架、活跃、公开/专属权限及有效订阅。用户路由尚未开放；尚无用户可读的报告或统计数据。
- 新增精确毫秒频次的最快 95% 均值和 P90、缓存率、重试感知成功率、输入费用分摊、输入均价及上海预算日基础算法。它们尚未接入真实网关、扣费或聚合。
- 前端仅新增管理 API 类型和调用函数，以及 API 合约测试；没有页面、菜单或可运行的大厅体验，没有构建或替换嵌入 UI。

## 管理 API

前缀 `/api/v1/admin/provider-hall`，沿用管理员认证、审计及现有管理面中间件，handler 另检查管理员身份。

| 方法与路径 | 本次行为 |
| --- | --- |
| `GET /config` | 返回单例配置及版本、更新时间、修改人 |
| `PUT /config` | 校验并事务更新配置，版本不匹配返回 409 |
| `GET /groups/:id` | 获取原分组的大厅配置；尚未配置时返回 `version=0`、`listed=false` |
| `PUT /groups/:id` | 首次保存要求 `version=0`；后续保存携带读取的版本；仅支持 OpenAI/混合分组 |
| `GET /profiles` | 获取模型与协议档案列表 |
| `POST /profiles` | 携带 `version=0` 创建，模型与协议组合唯一 |
| `PUT /profiles/:id` | 携带当前版本更新档案；重复组合返回 409 |

写入使用独立 DTO，拒绝未知字段和多份 JSON，限制请求体为 64 KiB。修改人由认证上下文取得，不能由请求体指定。`PUT` 为完整配置替换；必须发送该 DTO 的所有有效配置值，不应原样回传含 `updated_at/updated_by/id/group_id` 的读取响应。

配置金额通过 decimal 处理并作为十进制字符串传输。预算为 `numeric(20,8)`，分析价格为 `numeric(24,10)`；进入数据库前拒绝超精度、负数和科学计数形式，避免数据库自动舍入。参考价格、缓存率、确认时间必须一起配置或一起清空，缓存率限制为 0～1。

本站 origin 当前只允许无路径、查询、凭据或片段的 HTTPS origin。这里只保存配置，不会发起请求；后续 runner 仍须实现本站归属验证、禁止跨 origin 重定向和独立的本地测试许可。

默认模型必须存在对应协议档案。修改正在被默认配置引用的档案模型/协议返回冲突；先改默认比较配置，再修改该档案。repository 使用事务与统一锁顺序，避免并发更改留下失效引用。

## 删除与版本语义

大厅分组主键直接引用原分组 ID，硬删除级联移除上架配置。原分组软删除时，每次读取及保存都检查原记录；保存时锁定原分组，避免在并发删除后新增配置。名称和排序更新只修改配置版本。

配置和档案没有删除接口。任务、历史报告及其保留期限尚未实现，不能用当前三张配置表推断未来事实表的删除策略。

## 本次验证

| 检查 | 实际结果 |
| --- | --- |
| Go 1.26.5 聚焦 `-tags=unit -race` | 7 个顶层测试，含子用例共 32 项通过，失败 0，跳过 0；repository 此命令无单元用例，其实际数据库检查见下行 |
| PostgreSQL 18.6 本地临时数据库 | `empty` 和 `upgrade` 两个子场景通过；原有完整迁移链及 229 迁移、重复启动、单例约束、版本 CAS 并发、唯一模型/协议、金额精度、清空参考价格、软删除及硬删除级联均通过 |
| Go 1.26.6 聚焦单元兼容检查 | service、admin handler、DTO 通过；repository 编译通过，该命令没有 repository 单元用例 |
| `make -C backend build` | 通过，普通非嵌入 UI 构建；不是生产来源或部署校验 |
| Node 20 + pnpm 9 冻结安装 | 通过，未修改 `package.json` 或 `pnpm-lock.yaml` |
| 新增前端文件 ESLint | 通过 |
| 前端大厅 API Vitest | 1 个文件、3 项通过，覆盖金额字符串、版本冲突传播和取消信号 |
| 全量 `make -C backend test-unit` | 失败；观察到未修改模块中的账户/监控复制、余额文案、Antigravity 模型映射等失败，不计为全量通过 |
| 全量前端 `typecheck` | 失败；错误位于本次未修改的 `AccountUsageCell.vue`、`AccountTestModal.vue`、`GroupsView.vue`、`EmailVerifyView.vue`、`RegisterView.vue` |
| `CI=1` repository 集成聚焦命令 | 编译失败；原有 `account_repo_integration_test.go` 的 require、`account_repo_sort_integration_test.go` 的 fmt/time 未使用，尚未执行测试；本机 `docker info` 也确认 Docker 不可用 |

全量失败是本次检查实际观察到的结果；未将历史文档中的测试结果当作基线，也未恢复旧版本运行对照。本次没有修改这些无关模块来使全量检查通过。

本地数据库测试使用新建临时目录、独立用户和仅 Unix socket 监听的 PostgreSQL 实例，结束后停止进程并删除目录。不读取应用配置或生产凭据，不通过 `.installed` 绕过安装。这是 repository 验证，不是完整应用 setup/E2E 验证，也不替代 Docker/PostgreSQL/Redis 集成套件。

本机安装了 PostgreSQL 18.6 程序及其运行目录链接，未启动 Homebrew 常驻服务或创建应用数据库。Go 模块可自动选择 1.26.5，另下载并验证了 1.26.6；原 CI 与生产工具链版本冲突保持原状，不能宣称 CI 通过。

### 复现命令

从仓库根目录执行；网络需要时可在本地显式选择可访问的 `GOPROXY`，不修改仓库或全局 Go 配置。

```sh
go -C backend test -tags=unit -race ./internal/service ./internal/handler/admin ./internal/handler/dto ./internal/repository -run ProviderHall -count=1
GOTOOLCHAIN=go1.26.6 go -C backend test -tags=unit ./internal/service ./internal/handler/admin ./internal/handler/dto ./internal/repository -run ProviderHall -count=1
PROVIDER_HALL_PG_BIN=/opt/homebrew/opt/postgresql@18/bin go -C backend test -tags=providerhall_localdb ./internal/repository -run ProviderHall -count=1 -v
npm exec --yes --package=node@20 --package=pnpm@9 -- pnpm --dir frontend exec eslint src/api/admin/providerHall.ts src/api/__tests__/admin.providerHall.spec.ts
npm exec --yes --package=node@20 --package=pnpm@9 -- pnpm --dir frontend exec vitest run src/api/__tests__/admin.providerHall.spec.ts
```

`PROVIDER_HALL_PG_BIN` 是本地可选测试入口，其他机器改为自己的 PostgreSQL bin 目录。正常 CI 仍使用计划中的 `integration` 标签；没有降低原 harness 对 Docker 的要求。

## 下一步与验收边界

第一批还需完善剩余 DTO/接口合约、目标与专用 Key 配置实体，以及三协议模拟上游基座。后续按计划继续实现 collector、三入口观察点、异步计费关联、水位/缺口、聚合快照，再实现调度、预算、检测报告和页面。

未实现或未执行：用户列表/详情/报告接口，用户报价与排序，目标/Key 注册，采集与账单对账，任务/检测/预算运行逻辑，公共功能设置投影及中英文页面文案，Playwright，完整应用 setup，全量 frontend lint/Vitest，`make test-frontend`，全量 integration，H 性能验收，生产合约及灰度部署。

只有采集、任务及展示的各自实现和验收完成后，才可逐项移除 `PROVIDER_HALL_NOT_READY` 限制。当前状态不满足发布条件，未部署。
