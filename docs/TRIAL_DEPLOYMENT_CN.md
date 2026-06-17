# 4146 隔离试运行环境

本文记录当前服务器上固定使用的试运行方式。所有功能验证先走 4146 trial 环境，确认通过后再按 `docs/FORK_WORKFLOW_CN.md` 进入正式发布流程。

## 固定隔离边界

正式环境：

- 应用容器：`sub2api`
- 端口：`4145:8080`
- 数据库：`sub2api-postgres`
- Redis：`sub2api-redis`

试运行环境：

- 应用容器：`sub2api-image-gallery-trial`
- 端口：`4146:8080`
- 数据库：`sub2api-image-gallery-postgres-trial`
- Redis：`sub2api-image-gallery-redis-trial`

试运行只允许替换 `sub2api-image-gallery-trial`。不要停止、删除、重建或迁移正式环境容器和数据。

## 禁止操作

不要执行：

```bash
docker compose down -v
rm -rf data postgres_data redis_data
docker volume rm ...
```

不要删除或重建：

```text
sub2api
sub2api-postgres
sub2api-redis
sub2api-image-gallery-postgres-trial
sub2api-image-gallery-redis-trial
```

## 更新 trial 镜像

在服务器执行：

```bash
cd /root/sub2api-src

git fetch origin
git checkout feature/image-gallery
git pull --ff-only origin feature/image-gallery

TRIAL_TAG="$(git rev-parse --short HEAD)"
TRIAL_IMAGE="sub2api-test:image-playground-${TRIAL_TAG}"
echo "$TRIAL_IMAGE"

docker build -t "$TRIAL_IMAGE" .
```

如果只是复用已经构建好的镜像，可以跳过 `docker build`。

替换 trial 容器前必须确认镜像真实存在，避免 `docker run` 尝试从远端拉取本地测试镜像：

```bash
docker image inspect "$TRIAL_IMAGE" >/dev/null || {
  echo "本地不存在 $TRIAL_IMAGE，请先 docker build -t \"$TRIAL_IMAGE\" ."
  exit 1
}
```

## 替换 4146 trial 应用

只替换应用容器，保留 trial 数据库和 Redis：

```bash
TRIAL_NET="$(docker inspect -f '{{range $name, $_ := .NetworkSettings.Networks}}{{println $name}}{{end}}' sub2api-image-gallery-trial 2>/dev/null | head -n1)"
test -n "$TRIAL_NET" || TRIAL_NET="$(docker inspect -f '{{range $name, $_ := .NetworkSettings.Networks}}{{println $name}}{{end}}' sub2api-image-gallery-postgres-trial | head -n1)"
test -n "$TRIAL_NET" || { echo "TRIAL_NET 为空，先不要继续"; exit 1; }

docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' sub2api-image-gallery-trial > /tmp/sub2api-image-gallery-trial.env
test -s /tmp/sub2api-image-gallery-trial.env || { echo "/tmp/sub2api-image-gallery-trial.env 不存在，先不要继续"; exit 1; }

docker rm -f sub2api-image-gallery-trial

TRIAL_TAG="$(git rev-parse --short HEAD)"
TRIAL_IMAGE="sub2api-test:image-playground-${TRIAL_TAG}"
docker image inspect "$TRIAL_IMAGE" >/dev/null || {
  echo "本地不存在 $TRIAL_IMAGE，请先 docker build -t \"$TRIAL_IMAGE\" ."
  exit 1
}

docker run -d \
  --name sub2api-image-gallery-trial \
  --restart unless-stopped \
  --network "$TRIAL_NET" \
  -p 4146:8080 \
  --env-file /tmp/sub2api-image-gallery-trial.env \
  -e SERVER_HOST=0.0.0.0 \
  -e SERVER_PORT=8080 \
  -e DATABASE_HOST=sub2api-image-gallery-postgres-trial \
  -e REDIS_HOST=sub2api-image-gallery-redis-trial \
  -e AUTO_SETUP=true \
  "$TRIAL_IMAGE"
```

## 健康检查

```bash
docker ps | grep sub2api
docker inspect -f '{{.Config.Image}}' sub2api-image-gallery-trial
curl -fsS http://127.0.0.1:4146/health && echo
docker logs --tail=100 sub2api-image-gallery-trial
```

`docker inspect` 输出必须等于刚才的 `TRIAL_IMAGE`，例如 `sub2api-test:image-playground-ac8950aa`。如果仍显示旧 tag，说明 4146 trial 容器没有用新镜像启动。

## 生图广场测试清单

浏览器打开：

```text
http://服务器IP:4146
```

后台设置：

```text
管理后台 -> 生图广场
```

确认：

- 生图广场已启用。
- allowed models 包含实际可用模型，例如 `gpt-image-1`。
- default model 是可用模型。
- Agent 可按测试需要开启。
- Web Search 默认关闭，除非明确测试联网能力。
- 测试用户的 KEY 所属分组开启 `allow_image_generation`。

用户侧入口：

```text
http://服务器IP:4146/image-playground/index.html
```

验证：

- 页面能打开，不返回 404。
- 输入框旁有提示词参考按钮，点击后出现提示词库。
- 设置里没有真实 API Key 输入框。
- 设置里没有 Base URL 输入框。
- KEY 下拉只显示当前登录用户自己的可用 KEY。
- 文生图成功。
- 参考图/蒙版编辑成功。
- Agent 生图成功。
- 浏览器 Network 只访问 `/api/v1/image-playground/proxy/...`，不直接访问 `/v1/...` 或外部 OpenAI 域名。
- 管理后台 usage logs 有记录。
- 余额或额度扣减正常。

## 故障排查

如果 trial 应用启动失败，收集：

```bash
docker logs --tail=200 sub2api-image-gallery-trial
docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' sub2api-image-gallery-trial | grep -E 'DATABASE_|REDIS_|SERVER_|AUTO_SETUP'
docker inspect -f '{{range $name, $_ := .NetworkSettings.Networks}}{{println $name}}{{end}}' sub2api-image-gallery-trial
```

如果 `docker images` 里出现较多 `<none>` 悬空镜像，可以只清理 dangling 镜像：

```bash
docker image prune -f
```

不要使用带 volume 删除语义的清理命令。
