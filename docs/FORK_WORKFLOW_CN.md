# Sub2API Fork 开发与发布规范

本文档是本 fork 的硬规则，用于保证本地开发、服务器试运行、正式发布、同步官方更新和回滚都有固定路径。除非明确决定调整流程，否则所有改动按本文执行。

## 远端仓库

- `origin`: `git@github.com:hua7448/sub2api.git`
- `upstream`: `https://github.com/Wei-Shaw/sub2api.git`

本地不得把 `origin` 指回官方仓库。官方仓库只作为同步来源。

## 分支规范

- `main`: 正式发布分支。只有通过验证、准备发布的代码才能合入。
- `develop`: 迭代集成分支。多个功能可先合到这里验证。
- `feature/<name>`: 新功能分支，例如 `feature/custom-branding`。
- `fix/<name>`: 缺陷修复分支，例如 `fix/docker-legal-docs`。
- `theme/<name>`: 视觉主题分支，例如 `theme/anthropic-colors`。
- `sync/upstream-<date>`: 同步官方更新的临时分支，例如 `sync/upstream-2026-06-14`。
- `release/<version>`: 发布准备分支，例如 `release/v0.1.136-smartapi.1`。
- `hotfix/<name>`: 生产紧急修复分支。

禁止直接在 `main` 上做日常开发。所有业务改动先进入功能分支，经验证后再合并。

## 版本命名

发布 tag 使用：

```text
v<官方基线版本>-smartapi.<序号>
```

示例：

```text
v0.1.136-smartapi.1
v0.1.136-smartapi.2
v0.1.137-smartapi.1
```

规则：

- `<官方基线版本>` 取当前合并过的 upstream 版本。
- `<序号>` 从 1 开始递增。
- 管理员后台“一键更新”只检查本 fork 的稳定 tag，即 `-smartapi.<序号>` 后缀。
- 不使用 `latest` 作为生产部署依据；生产部署必须写明确版本 tag。
- `latest` 只能作为辅助标签，不能作为回滚依据。

## 本地开发流程

本地 Go 环境：

- Go 版本必须与 `backend/go.mod` 和 release workflow 保持一致；当前为 `go1.26.4`。
- macOS 本地推荐使用 Homebrew 安装固定稳定版：

```bash
brew install go
go version
```

- 网络不稳定时可以临时使用本地代理执行安装或测试：

```bash
HOMEBREW_NO_AUTO_UPDATE=1 \
ALL_PROXY=http://127.0.0.1:7897 \
HTTPS_PROXY=http://127.0.0.1:7897 \
HTTP_PROXY=http://127.0.0.1:7897 \
brew install go

ALL_PROXY=http://127.0.0.1:7897 \
HTTPS_PROXY=http://127.0.0.1:7897 \
HTTP_PROXY=http://127.0.0.1:7897 \
GOPROXY=https://proxy.golang.org,direct \
go test ./...
```

- 提交 Go 代码前必须对改过的 Go 文件执行 `gofmt`。不要提交因 `go mod tidy` 引入的无关依赖变更。

开始新需求：

```bash
git checkout main
git pull origin main
git fetch upstream
git checkout -b feature/<name>
```

完成开发后至少执行相关验证：

```bash
pnpm --dir frontend run build
```

后端改动必须执行：

```bash
cd backend
go test ./...
```

提交信息使用简洁前缀：

- `feat:` 新功能
- `fix:` 修复
- `style:` 视觉或样式调整
- `docs:` 文档
- `chore:` 构建、配置、维护
- `sync:` 同步 upstream

示例：

```bash
git commit -m "style: apply anthropic inspired theme"
```

## 合并规则

- 功能分支先合入 `develop` 做集成验证。
- 准备发布时，从 `develop` 创建 `release/<version>`。
- `release/<version>` 验证通过后合入 `main`。
- `main` 只允许包含已准备发布或已发布的提交。

小范围、低风险改动可以跳过 `develop`，但仍必须走功能分支到 `main` 的合并路径，并完成验证。

## 同步官方更新

同步 upstream 必须单独开分支：

```bash
git checkout main
git pull origin main
git fetch upstream
git checkout -b sync/upstream-YYYY-MM-DD
git merge upstream/main
```

如有冲突，解决后执行前后端验证。同步分支验证通过后合入 `develop` 或 `main`。

同步官方更新时禁止混入自定义功能开发。同步提交使用：

```bash
git commit -m "sync: merge upstream main YYYY-MM-DD"
```

## 官方迭代跟踪准则

当需要“跟踪官方更新”时，先判断官方基线版本和变更类型，再决定是否同步和发布。生产永远只更新到本 fork 验证过的 SmartAPI 稳定 tag，不直接跟随 upstream、origin/main 或 `latest`。

默认只有官方发布新的基线版本 tag（例如从 `v0.1.138` 到 `v0.1.139`）时，才认为需要进入常规 upstream 同步评估。同一官方基线内的 upstream 零散提交默认不合并；如确认为安全修复、严重 bugfix 或影响核心可用性，必须明确说明这是同基线紧急例外，并仍按同步分支、验证和 4146 试运行流程执行。

检查官方差异：

```bash
git checkout main
git fetch origin
git fetch upstream
git log --oneline main..upstream/main
git diff --stat main..upstream/main
```

按风险分级处理：

- 低风险：文档、样式、非核心前端展示、小范围 bugfix。可以按周同步，验证通过后正常发 `v<官方基线版本>-smartapi.N`。
- 中风险：账号调度、计费统计、渠道管理、缓存、并发、后台设置、Docker/CI 构建。必须开 `sync/upstream-YYYY-MM-DD`，跑相关后端测试和前端 build，先 4146 试运行。
- 高风险：数据库迁移、支付、OAuth 登录/授权、注册登录、API key 创建/鉴权、订阅/余额扣费、权限判断、数据备份/恢复、安全策略。必须单独列出风险点、验证回滚路径、备份数据库，并在发布前说明是否存在不可逆迁移。

同步执行准则：

- 同步分支只做 upstream merge 和冲突解决，不混入 SmartAPI 新功能、主题调整或业务定制。
- 冲突解决优先保留本 fork 的固定规则：SmartAPI 更新源、`-smartapi.N` tag 规则、`prerelease: false`、生产不用 `latest`、4146 先试运行。
- 如果 upstream 修改了同一块更新逻辑、发布 workflow、Dockerfile、配置默认值，必须重新确认一键更新链路：`/releases/latest`、release assets、checksum、容器内 `/app` 可写、重启策略。
- 如果 upstream 新增或修改数据库 migration，发布说明必须标注“包含数据库迁移”，生产替换前必须备份，回滚前必须评估迁移是否向后兼容。
- 如果官方更新只是发布 tag 但 main 无新提交，仍以代码差异为准；不要为了追 tag 直接发布。

同步完成后的发布命名：

```text
官方基线 0.1.137 第一次 SmartAPI 发布：v0.1.137-smartapi.1
同一官方基线上的 SmartAPI 修复：v0.1.137-smartapi.2
下一次官方基线 0.1.138：v0.1.138-smartapi.1
```

建议节奏：

- 普通官方更新：每周检查一次。
- 安全修复、严重 bugfix、影响核心可用性的 upstream 更新：可作为同基线紧急例外，当天同步、当天 4146 试运行。
- 无明确收益的 upstream 大改：先记录差异和风险，不急于发布到生产。

## 镜像发布

正式发布优先使用 GitHub Actions `Release` workflow。需要管理员后台一键更新时，必须使用完整 release：

- `simple_release`: `false`
- 产物：`ghcr.io/hua7448/sub2api:<version>`
- Release assets 必须包含各平台二进制压缩包和 `checksums.txt`，否则后台只能检查版本，不能完成热更新。
- GoReleaser 配置必须保持 `prerelease: false`。如果 `-smartapi.<序号>` tag 被 GitHub 标成 prerelease，后台 `/releases/latest` 将看不到它。
- 管理员后台一键更新/热更新只应假设会替换应用二进制 `/app/sub2api`；不要假设 release 压缩包中的 `deploy/*`、容器镜像中的 `/app/resources` 或其他静态资源会同步到现有生产容器。
- 如果本次修复依赖模型价格、默认模板、规则表、内置 JSON/YAML 等静态资源随版本生效，必须满足至少一个条件：
  - 通过 `go:embed` 或等效方式编进 Go 二进制；
  - 提供明确的数据迁移/复制逻辑，把资源写入持久化目录；
  - 明确要求生产重新拉取镜像并重建容器，而不是后台热更新，并在发布说明中写清楚。
- 任何“fallback 文件已更新”的修复，都必须额外验证“只有二进制更新、磁盘 resources 仍旧”的场景。`v0.1.145-smartapi.1` 的 Kimi K2.7 价格修复曾因此在正式热更新后不生效，最终由 `v0.1.145-smartapi.2` 通过内嵌 fallback JSON 修复。

本地或服务器手动构建只用于试运行，不作为正式发布来源。

首次使用 fork Actions 时，GitHub 可能会禁用 workflow。需要在仓库 Actions 页面手动启用：

```text
https://github.com/hua7448/sub2api/actions
```

如需本地查看或触发 workflow，安装并登录 GitHub CLI：

```bash
HOMEBREW_NO_AUTO_UPDATE=1 \
ALL_PROXY=http://127.0.0.1:7897 \
HTTPS_PROXY=http://127.0.0.1:7897 \
HTTP_PROXY=http://127.0.0.1:7897 \
brew install gh

ALL_PROXY=http://127.0.0.1:7897 \
HTTPS_PROXY=http://127.0.0.1:7897 \
HTTP_PROXY=http://127.0.0.1:7897 \
gh auth login -h github.com -p ssh -w
```

正式发布步骤：

```bash
git checkout main
git pull origin main
git tag -a v0.1.136-smartapi.1 -m "Release v0.1.136-smartapi.1"
git push origin v0.1.136-smartapi.1
```

如果 tag 没有自动触发，或需要手动触发完整 release：

```bash
gh workflow run release.yml \
  --repo hua7448/sub2api \
  --ref v0.1.136-smartapi.1 \
  -f tag=v0.1.136-smartapi.1 \
  -f simple_release=false

gh run watch --repo hua7448/sub2api <run-id> --interval 15 --exit-status
```

注意：workflow 必须以 tag ref 触发，即 `headBranch` 应为 `v0.1.136-smartapi.1`，`headSha` 应为该 tag 指向的提交。不要在 `main` ref 上手动发布同一个 tag，否则可能用错代码构建前端产物。

Release 成功后必须确认：

```bash
gh release view v0.1.136-smartapi.1 \
  --repo hua7448/sub2api \
  --json isDraft,isPrerelease,name,tagName,url,assets

gh api repos/hua7448/sub2api/releases/latest \
  --jq '.tag_name, .prerelease, .name'
```

期望：

- `isPrerelease` 为 `false`。
- `/releases/latest` 返回当前 `v...-smartapi.N` tag。
- assets 中包含 `checksums.txt`、`linux_amd64.tar.gz`、`linux_arm64.tar.gz`。

GitHub Actions 成功后，生产 compose 使用明确镜像：

```yaml
services:
  sub2api:
    image: ghcr.io/hua7448/sub2api:0.1.136-smartapi.1
```

首个支持 SmartAPI 一键更新的生产镜像部署后，后续小版本可以继续使用管理员后台“检查更新”下载本 fork 的完整 release 二进制并重启服务。容器必须带 `restart: unless-stopped` 或等效重启策略；如果以后重新创建容器，仍应使用最新明确镜像 tag，避免回到旧镜像内置版本。

## 服务器试运行

正式切换前必须先起新端口试运行实例，例如 `4146`。

当前服务器固定的 4146 隔离试运行环境、容器名、更新命令和测试清单见 `docs/TRIAL_DEPLOYMENT_CN.md`。后续 agent 工作应优先按该文档执行，确保测试环境与正式 4145、正式 PostgreSQL、正式 Redis 隔离。

如果在 Codex App 等隔离环境中执行试运行时遇到 `Permission denied (publickey...)`，说明该环境没有服务器 SSH 凭据；这不是 GitHub Release 权限问题。必须先按 `docs/TRIAL_DEPLOYMENT_CN.md` 恢复 SSH 权限，或改由有 SSH 权限的本机 CLI 完成 4146 trial。不得因自动化环境缺少 SSH key 而跳过试运行。

生图广场相关开发、测试、热修发布、提示词库、下载/离页保护、KEY 下拉排查和上线公告模板见 `docs/IMAGE_PLAYGROUND_RUNBOOK_CN.md`。涉及该模块时，先读专项 runbook，再按本发布规范执行。

试运行可以使用服务器手动构建镜像：

```bash
cd /root/sub2api-src
git fetch origin
git checkout <branch-or-tag>
git pull
docker build -t sub2api-test:<name> .
```

试运行容器必须：

- 使用 `/root/sub2api-deploy/.env`
- 连接现有 `sub2api-postgres` 和 `sub2api-redis`
- 使用独立容器名
- 使用非主服务端口，例如 `4146`
- 不执行破坏性后台操作

示例：

```bash
docker run -d \
  --name sub2api-smartapi-test \
  --restart unless-stopped \
  --network sub2api-deploy_sub2api-network \
  -p 4146:8080 \
  --env-file /root/sub2api-deploy/.env \
  -e SERVER_HOST=0.0.0.0 \
  -e SERVER_PORT=8080 \
  -e DATABASE_HOST=sub2api-postgres \
  -e REDIS_HOST=sub2api-redis \
  -e AUTO_SETUP=true \
  -v /root/sub2api-deploy/data:/app/data \
  ghcr.io/hua7448/sub2api:<version>
```

如果 `4146` 被旧测试容器占用，先查看并清理旧测试容器，或换临时端口：

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Ports}}' | grep 4146 || true
ss -ltnp | grep ':4146' || true
docker rm -f sub2api-smartapi-test 2>/dev/null || true
```

验证通过后才能进入正式切换。

## 正式部署

生产环境只改应用镜像，不动数据库和 Redis 数据目录。

正式部署前备份数据库：

```bash
cd /root/sub2api-deploy
docker compose exec postgres pg_dump -U sub2api -d sub2api -Fc -f /tmp/sub2api.dump
docker compose cp sub2api-postgres:/tmp/sub2api.dump ./sub2api-$(date +%F-%H%M%S).dump
```

更新 `docker-compose.override.yml` 的镜像 tag：

```yaml
services:
  sub2api:
    image: ghcr.io/hua7448/sub2api:<version>
```

部署：

```bash
docker compose pull sub2api
docker compose up -d sub2api
docker compose logs -f sub2api
```

当前生产环境如直接使用 `docker run` 管理 4145 主服务，替换命令为：

```bash
docker rm -f sub2api

docker run -d \
  --name sub2api \
  --restart unless-stopped \
  --network sub2api-deploy_sub2api-network \
  -p 4145:8080 \
  --env-file /root/sub2api-deploy/.env \
  -e SERVER_HOST=0.0.0.0 \
  -e SERVER_PORT=8080 \
  -e DATABASE_HOST=sub2api-postgres \
  -e REDIS_HOST=sub2api-redis \
  -e AUTO_SETUP=true \
  -v /root/sub2api-deploy/data:/app/data \
  ghcr.io/hua7448/sub2api:<version>
```

部署后检查：

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' | grep sub2api
docker logs -f sub2api
docker exec sub2api /app/sub2api --version
curl -fsS http://127.0.0.1:4145/health && echo
```

## 回滚规则

每次正式发布必须记录上一个生产镜像 tag。

回滚时只改回旧 tag：

```yaml
services:
  sub2api:
    image: ghcr.io/hua7448/sub2api:<previous-version>
```

然后执行：

```bash
cd /root/sub2api-deploy
docker compose pull sub2api
docker compose up -d sub2api
docker compose logs -f sub2api
```

如果新版本包含数据库迁移，回滚前必须先判断迁移是否向后兼容。无法确认时，不允许直接回滚应用，需要先备份并评估数据库状态。

## 数据安全红线

禁止执行：

```bash
docker compose down -v
rm -rf data postgres_data redis_data
docker volume rm <production-volume>
```

禁止把以下内容提交到 Git：

- `.env`
- 数据库 dump
- `data/`
- `postgres_data/`
- `redis_data/`
- 用户使用记录和私有报表

## 发布检查清单

发布前必须确认：

- 当前分支已合入 `main`
- `git status` 干净
- 前端构建通过
- 后端相关测试通过，若未运行需说明原因
- 服务器试运行端口验证通过
- 已完成数据库备份
- 已记录当前生产镜像 tag，具备回滚路径
- 正式部署使用明确版本 tag，不使用 `latest`

## 站点文档维护（site/）

`site/` 是 SmartQ 客户使用文档站点（Next.js 静态站点）。

### 目录说明

| 路径 | 说明 |
|------|------|
| `site/lib/docs-data.ts` | 站点全局配置：主站、Base URL、QQ 群、工作时间、价格表、FAQ、导航 |
| `site/components/docs/hero.tsx` | 首页首屏文案和终端展示 |
| `site/components/docs/quick-start.tsx` | 快速开始四步卡片，含充值入口链接 |
| `site/components/docs/pricing.tsx` | 充值套餐表和充值按钮 |
| `site/components/docs/billing.tsx` | 计费说明 |
| `site/components/docs/api-key-guide.tsx` | API Key 使用教程 |
| `site/components/docs/codex-config.tsx` | Codex CLI/App 配置教程 |
| `site/components/docs/claude-config.tsx` | Claude Code 配置教程，含 GLM 配置 |
| `site/components/docs/kimi-config.tsx` | Kimi CLI 配置教程 |
| `site/components/docs/faq-section.tsx` / `faq-accordion.tsx` | 常见问题（数据来自 `docs-data.ts`） |
| `site/components/docs/contact.tsx` | 联系支持 |
| `site/deploy/` | Dockerfile、nginx.conf、start.sh |

### 修改流程

1. 编辑对应组件或 `site/lib/docs-data.ts`
2. 本地验证：
   ```bash
   cd site
   pnpm install
   pnpm build
   pnpm dev --port 3002
   ```
3. 浏览器打开 `http://localhost:3002` 确认
4. 提交：`feat(site): ...`

### 常见修改映射

- 改主站/Base URL/QQ 群/工作时间 → `site/lib/docs-data.ts`
- 改价格 → `site/lib/docs-data.ts` 中的 `PRICING`
- 改常见问题 → `site/lib/docs-data.ts` 中的 `FAQ`
- 改导航 → `site/lib/docs-data.ts` 中的 `NAV_ITEMS`
- 改 Kimi 下载链接 → `site/components/docs/kimi-config.tsx`
- 改充值跳转链接 → `site/components/docs/pricing.tsx` 和 `site/components/docs/quick-start.tsx`
- 改 Claude/Codex/GLM 配置示例 → `site/components/docs/claude-config.tsx` / `codex-config.tsx`

### 注意事项

- 站点使用 Next.js 15 + React 19 + Tailwind CSS 4
- 构建产物输出到 `site/out/`，已加入根目录 `.gitignore`
- API Key 等敏感信息必须脱敏，示例中使用 `sk-********************************`
- 新增页面组件后需要在 `site/app/page.tsx` 中引入并渲染
