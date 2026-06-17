# 生图广场开发、测试与发布 Runbook

本文沉淀 `gpt_image_playground` 集成到 sub2api 后的固定经验。后续涉及生图广场、提示词库、KEY 选择、下载、离页保护、4146 试运行、正式发布和回滚时，先按本文检查，再执行 `docs/FORK_WORKFLOW_CN.md`。

## 核心边界

- 生图广场是 `frontend/playground` 下的独立 React/Vite 子应用，不把 React 组件混入主站 Vue 应用。
- 主站入口只负责跳转 `/image-playground/index.html`，管理端仍使用 Vue 的生图广场设置页。
- 浏览器永远不保存、不展示用户真实 API Key。playground 只能保存 `api_key_id`、模型和 UI 参数。
- 前端只调用 `/api/v1/image-playground/...`，不能直接请求 `/v1/...`、外部 OpenAI 域名、自定义 Base URL 或自定义 provider。
- 后端代理必须用登录态校验当前用户、KEY 归属、KEY 状态、额度、过期时间和分组 `allow_image_generation`。
- 本地历史、收藏、图片缓存、ZIP 导出继续使用 playground 的 IndexedDB，但 IndexedDB 不是计费事实来源。
- 正式发布优先走 GitHub Actions `Release` workflow 和 SmartAPI 一键更新；服务器手动镜像只用于 4146 试运行或应急回退。

## 关键文件

- 用户入口：`frontend/src/views/user/ImageGalleryView.vue`
- 管理设置页：`frontend/src/views/admin/ImageGalleryAdminView.vue`
- 主站嵌入和路由：`backend/internal/web/embed_on.go`
- 后端路由：`backend/internal/server/routes/user.go`
- 后端代理：`backend/internal/handler/image_playground_handler.go`
- playground API 层：`frontend/playground/src/lib/sub2api.ts`
- playground 主站继承：`frontend/playground/src/lib/sub2apiHost.ts`
- playground 下载：`frontend/playground/src/lib/downloadImages.ts`
- playground 离页保护：`frontend/playground/src/components/GenerationNavigationWarning.tsx`
- 提示词库索引：`frontend/playground/src/data/promptLibrary/index.ts`
- 提示词库分片：`frontend/playground/src/data/promptLibrary/chunks/`
- 4146 试运行：`docs/TRIAL_DEPLOYMENT_CN.md`

## 开发原则

1. 先确认需求属于主站、管理端、后端代理还是 playground 子应用，尽量只改对应边界。
2. 改 playground 时优先复用原项目的状态和组件结构，但必须保留 sub2api 安全边界。
3. 所有显示文字要走已有中英文方法：主站 Vue 使用 i18n，playground 使用 `text(...)` 或 `hostText(...)`。
4. 涉及站点名称、favicon、标题、语言、返回控制台等主站体验时，要从 host settings 继承，不写死上游 `GPT Image Playground`。
5. 涉及下载、缓存、IndexedDB 时要同时考虑 `data:`、本地图片 ID、HTTP URL 三种输入。
6. 涉及生成过程时，要考虑用户刷新、关闭、跳转、切换标签页导致任务中断的体验提示。
7. 涉及提示词库时，禁止运行时远程加载；把来源、license、commit 和数据 vendor 到本地，并保留署名。

## 本地验证

常规发布前至少执行：

```bash
pnpm --dir frontend run build
cd backend && go test ./...
```

只改 playground 时也建议执行：

```bash
pnpm --dir frontend/playground run build
```

改下载逻辑时执行专项测试：

```bash
pnpm --dir frontend/playground exec vitest run src/lib/downloadImages.test.ts
```

统计提示词模板数量：

```bash
printf 'chunk prompt count: '
rg '"prompt":' frontend/playground/src/data/promptLibrary/chunks/*.ts | wc -l
```

注意：

- `pnpm --dir frontend run build` 会同时构建主站和 playground。
- `backend/internal/web/dist/` 是生成产物，默认不提交；构建后要用 `git status --short` 确认没有意外变更。
- 如果 playground 上游遗留测试因 provider 被限制为 sub2api 而失败，先判断是否与本次改动相关；发布阻断以本次改动相关测试和完整 build 为准。

## 4146 试运行

所有生产前功能验证先走 4146 隔离环境：

- 正式应用：`sub2api`，端口 `4145`，数据库 `sub2api-postgres`，Redis `sub2api-redis`
- 试运行应用：`sub2api-image-gallery-trial`，端口 `4146`，数据库 `sub2api-image-gallery-postgres-trial`，Redis `sub2api-image-gallery-redis-trial`

固定命令见 `docs/TRIAL_DEPLOYMENT_CN.md`。

常见坑：

- `TRIAL_NET` 为空：通常是 shell 断开后变量丢失，或 inspect 了不存在的容器。重新从 trial 应用容器或 trial Postgres 容器读取网络名。
- `Unable to find image 'sub2api-test:...' locally`：当前 `TRIAL_IMAGE` 没有 build 出来。先 `docker image inspect "$TRIAL_IMAGE"`，不存在就 `docker build -t "$TRIAL_IMAGE" .`。
- 容器启动后页面仍是旧版本：检查 `docker inspect -f '{{.Config.Image}}' sub2api-image-gallery-trial`，必须等于刚构建的测试镜像。
- 访问 `/image-playground/` 异常：优先访问 `/image-playground/index.html`，并确认后端 embed 路由已包含 image-playground 静态资源。

试运行浏览器检查：

- `http://服务器IP:4146` 可以登录主站。
- `/image-playground/index.html` 可刷新打开。
- 未登录直接进入 playground 时，会提示先登录主站并提供跳转按钮。
- 页面 title、favicon、站点名、语言跟随主站。
- 有返回控制台入口。
- 没有真实 API Key 输入框、Base URL 输入框、自定义 provider、自定义代理或远程导入入口。
- KEY 下拉只显示当前登录用户自己的 eligible keys。
- 提示词库按钮可打开，列表加载不明显卡顿。
- 文生图、参考图/蒙版编辑、Agent 生图按测试要求可用。
- 下载单张、批量下载、ZIP 导出可用。
- 生成中刷新/关闭/跳转会出现中英文离页提示。
- Network 只访问 `/api/v1/image-playground/proxy/...`，不直连 OpenAI 或外部 API 域名。
- usage logs 和余额扣减以后端记录为准。

## KEY 下拉为空排查

先看接口：

```text
GET /api/v1/image-playground/eligible-keys
```

如果返回 `data: []`，按顺序检查：

1. 当前浏览器是否是登录态。
2. KEY 是否属于当前用户。
3. KEY 状态是否 `active`。
4. KEY 是否未过期。
5. KEY 额度是否未耗尽；`quota = 0` 通常表示不限额，需结合后端逻辑确认。
6. KEY 所属分组是否 active。
7. KEY 所属分组是否开启 `allow_image_generation`。
8. 管理端生图广场是否启用。
9. allowed models/default model 是否包含实际要用的模型。

看到用户在 `/api/v1/keys` 有 KEY，但 eligible keys 为空，最常见原因是分组没有开启 `allow_image_generation`。这不是 KEY 读取 bug。

## 下载问题排查

现象：点击图片详情下载失败，Network 没有任何请求。

判断：

- 生成图可能已经在 IndexedDB 里缓存成 `data:` URL。
- 下载逻辑如果直接 `fetch(dataUrl)`，某些环境会被限制或表现不稳定，而且本来就不应该产生网络请求。

固定处理：

- 本地图片 ID 先走 `ensureImageCached(id)`。
- 如果得到 `data:` URL，直接本地解码为 `Blob`。
- 只有 `http://` 或 `https://` URL 才走 `fetch`。
- 下载测试要断言 `data:` URL 不调用 `fetch`，而是创建 object URL 后触发 `<a download>`。

验证命令：

```bash
pnpm --dir frontend/playground exec vitest run src/lib/downloadImages.test.ts
```

## 生成中离页保护

问题：图片或 Agent 生图未完成时，用户切换、刷新或关闭页面，任务会中断。

固定处理：

- 只要 `tasks` 中存在 `running`，或 Agent conversation round 存在 `running`，页面显示顶部提示。
- 同时注册 `beforeunload`，触发浏览器原生离页确认。
- 提示必须中英文齐全。
- 不要阻止用户主动离开，只做风险提醒和确认。

推荐中文文案：

```text
图片生成中，请不要切换页面或关闭窗口，完成后再离开。
```

推荐英文文案：

```text
Image generation is running. Do not switch pages or close this window until it finishes.
```

## 提示词库集成

固定原则：

- 来源数据 vendor 到 `frontend/playground/src/data/promptLibrary/`。
- `index.ts` 只放摘要和 source metadata，大 prompt 放进 `chunks/` 分片。
- UI 首屏不要一次性加载所有 prompt 正文，避免拖慢 playground 初始化。
- 保留 source name、URL、license、commit。
- 可署名，但不要在核心生成流程中跳外链。

更新后统计数量：

```bash
rg '"prompt":' frontend/playground/src/data/promptLibrary/chunks/*.ts | wc -l
```

公告里可以写准确数量，例如：

```text
生图广场已上线，效果炸裂。内置 885 个提示词模板，可直接参考、套用和改写。使用按量扣费的 Codex 分组 API Key 即可生成图片，欢迎体验。
```

## 热修发布流程

适用于下载失败、权限判断、登录态、计费链路、页面无法打开等生产紧急问题。

```bash
git fetch origin --tags
git checkout main
git pull --ff-only origin main
git checkout -b hotfix/<name>
```

如果本地已有用户未提交改动，先保护：

```bash
git stash push -m "preserve local changes before <name>"
```

修复后提交到 hotfix 分支，并在分支上完成相关验证。随后合入 `main`：

```bash
git checkout main
git pull --ff-only origin main
git merge --no-ff hotfix/<name> -m "fix: release <name>"
pnpm --dir frontend run build
cd backend && go test ./...
cd ..
```

打下一个 SmartAPI tag：

```bash
VERSION=v0.1.136-smartapi.3
git push origin main
git tag -a "$VERSION" -m "Release $VERSION" -m "本版说明..."
git push origin "$VERSION"
```

等待 release workflow：

```bash
gh run list --repo hua7448/sub2api --limit 10 \
  --json databaseId,workflowName,headBranch,headSha,status,conclusion,createdAt,url

gh run watch --repo hua7448/sub2api <release-run-id> --interval 15 --exit-status
```

必须确认 release run 的：

- `workflowName` 是 `Release`。
- `headBranch` 是 tag，例如 `v0.1.136-smartapi.3`。
- `headSha` 是 tag 指向的提交。

Release 完成后确认：

```bash
gh release view "$VERSION" \
  --repo hua7448/sub2api \
  --json isDraft,isPrerelease,name,tagName,url,assets

gh api repos/hua7448/sub2api/releases/latest \
  --jq '.tag_name, .prerelease, .name, .assets[].name'
```

期望：

- `isDraft=false`
- `isPrerelease=false`
- `/releases/latest` 返回当前 tag
- assets 包含 `checksums.txt`、`linux_amd64.tar.gz`、`linux_arm64.tar.gz`

如果 workflow 自动提交了 `backend/cmd/server/VERSION` 同步提交，发布后拉回本地：

```bash
git pull --ff-only origin main
```

如果之前 stash 了用户改动，最后恢复：

```bash
git stash pop
```

恢复后要确认只剩用户原有改动，不混入热修文件。

## 生产更新与回滚

如果生产已部署支持 SmartAPI 完整 release 的版本，普通小版本更新优先使用管理员后台一键更新，不需要服务器重新 build 本地镜像。

仍需人工确认：

- GitHub latest 已是目标 `v...-smartapi.N`。
- Release 不是 prerelease。
- Release assets 和 checksums 存在。
- 本版是否包含数据库迁移。

如果必须手动替换生产容器，按 `docs/FORK_WORKFLOW_CN.md` 的正式部署章节执行，必须先备份数据库，并使用明确镜像 tag：

```text
ghcr.io/hua7448/sub2api:<version>
```

回滚规则：

- 无数据库迁移的小版本，应用回滚到上一个稳定 tag，风险较低。
- 有数据库迁移的版本，回滚前必须评估迁移是否向后兼容，不能只回滚应用。
- 每次发布说明里记录上一个稳定 tag。

## 常见故障速查

- `Image playground asset not found`：镜像没有包含 `backend/internal/web/dist/image-playground/`，或容器仍在旧镜像。重新 build 并确认 `docker inspect` 镜像 tag。
- `ERR_TOO_MANY_REDIRECTS`：检查 `/image-playground`、`/image-playground/`、`/image-playground/index.html` 的后端 rewrite，避免目录路径和 fallback 互相重定向。
- KEY 下拉为空：优先查分组 `allow_image_generation`。
- 页面仍显示上游名称或 favicon：检查 `sub2apiHost.ts`、`Header.tsx`、`SettingsModal.tsx`、`embed_on.go` 的 host settings 注入。
- 中英文混用：检查新文案是否使用 `text(...)`、`hostText(...)` 或 Vue i18n key。
- 下载没有请求且失败：优先查 `data:` URL 到 Blob 的本地转换。
- 提示词按钮看不到：确认容器镜像确实是最新构建，不只停留在旧 `image-playground-YYYYMMDD` tag。
- 试运行端口不是新代码：`docker ps` 只说明容器重启了，不说明用了新镜像；必须看 `docker inspect -f '{{.Config.Image}}'`。

## 发布完成后的简报模板

```text
已发布 vX.Y.Z-smartapi.N。

包含：
- 修复/新增 ...
- 本版无数据库迁移 / 本版包含数据库迁移，回滚需评估 ...

验证：
- pnpm --dir frontend run build
- cd backend && go test ./...
- 4146 试运行：通过 ...

Release：
- latest 已指向 vX.Y.Z-smartapi.N
- isPrerelease=false
- checksums、linux_amd64、linux_arm64 assets 已存在

回滚：
- 上一个稳定版本：vX.Y.Z-smartapi.N-1
```
