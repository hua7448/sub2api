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

## 镜像发布

正式发布优先使用 GitHub Actions `Release` workflow。需要管理员后台一键更新时，必须使用完整 release：

- `simple_release`: `false`
- 产物：`ghcr.io/hua7448/sub2api:<version>`
- Release assets 必须包含各平台二进制压缩包和 `checksums.txt`，否则后台只能检查版本，不能完成热更新。

本地或服务器手动构建只用于试运行，不作为正式发布来源。

正式发布步骤：

```bash
git checkout main
git pull origin main
git tag -a v0.1.136-smartapi.1 -m "Release v0.1.136-smartapi.1"
git push origin v0.1.136-smartapi.1
```

GitHub Actions 成功后，生产 compose 使用明确镜像：

```yaml
services:
  sub2api:
    image: ghcr.io/hua7448/sub2api:0.1.136-smartapi.1
```

首个支持 SmartAPI 一键更新的生产镜像部署后，后续小版本可以继续使用管理员后台“检查更新”下载本 fork 的完整 release 二进制并重启服务。容器必须带 `restart: unless-stopped` 或等效重启策略；如果以后重新创建容器，仍应使用最新明确镜像 tag，避免回到旧镜像内置版本。

## 服务器试运行

正式切换前必须先起新端口试运行实例，例如 `4146`。

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
  --name sub2api-test \
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
  sub2api-test:<name>
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
