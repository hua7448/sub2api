#!/usr/bin/env bash
# 4145 生产镜像切换：把 sub2api 容器整体换到指定镜像（二进制 + 磁盘资源一起刷新）。
# 常规正式发布默认走后台页面热更新；仅在需要刷新磁盘资源、镜像基线或容器级回滚时使用本脚本。
# 复用当前生产容器的 env，保留生产 DB/Redis 数据卷，不碰数据。
#
# 前置：已备份生产库（pg_dump）。详见 docs/FORK_WORKFLOW_CN.md「正式部署」「回滚规则」。
#
# 用法：
#   bash deploy/switch-4145.sh                                   # 默认切到本次 release
#   SWITCH_IMAGE=ghcr.io/hua7448/sub2api:<tag> bash deploy/switch-4145.sh
#   SWITCH_IMAGE=ghcr.io/hua7448/sub2api:<old-tag> bash deploy/switch-4145.sh # 回滚
set -uo pipefail

IMAGE="${SWITCH_IMAGE:-ghcr.io/hua7448/sub2api:0.1.163-smartapi.1}"
APP="sub2api"
PROD_DB="sub2api-postgres"
PROD_REDIS="sub2api-redis"
DB_USER="${DB_USER:-sub2api}"
DB_NAME="${DB_NAME:-sub2api}"

echo "=== 0. 前置检查 ==="
ls /root/sub2api-deploy/sub2api-*.dump >/dev/null 2>&1 && echo "DB dump 存在 ✓" || echo "!! 警告：没看到 /root/sub2api-deploy/sub2api-*.dump，确认已备份再继续"
docker image inspect "$IMAGE" >/dev/null && echo "image OK: $IMAGE" || { echo "image missing: $IMAGE  (先 docker pull)"; exit 1; }

echo "=== 1. 记录切换前生产状态（回滚依据）==="
PREV_IMG="$(docker inspect -f '{{.Config.Image}}' "$APP" 2>/dev/null || true)"
PREV_VER="$(docker exec "$APP" /app/sub2api --version 2>&1 | head -1 || true)"
{
  echo "switched_at: $(date -u +%FT%TZ)"
  echo "prev_image: $PREV_IMG"
  echo "prev_version: $PREV_VER"
} | tee /root/sub2api-deploy/PREV_STATE_4145.txt

echo "=== 2. 复用当前生产 env（保证配置一致）==="
docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' "$APP" > /tmp/prod.env
[ -s /tmp/prod.env ] || { echo "无法读取 $APP env，停止"; exit 1; }
echo "env lines: $(wc -l < /tmp/prod.env)"

echo "=== 3. 探测生产网络 ==="
PROD_NET="$(docker inspect -f '{{range $n,$_ := .NetworkSettings.Networks}}{{println $n}}{{end}}' "$APP" 2>/dev/null | head -n1 || true)"
[ -n "$PROD_NET" ] || PROD_NET="$(docker inspect -f '{{range $n,$_ := .NetworkSettings.Networks}}{{println $n}}{{end}}' "$PROD_DB" 2>/dev/null | head -n1 || true)"
[ -n "$PROD_NET" ] || { echo "探测不到生产网络，停止"; exit 1; }
echo "PROD_NET=$PROD_NET"

echo "=== 4. 替换生产容器（保留数据卷；连生产 DB/Redis）==="
docker rm -f "$APP" >/dev/null 2>&1 || true
docker run -d --name "$APP" --restart unless-stopped --network "$PROD_NET" -p 4145:8080 \
  --env-file /tmp/prod.env \
  -e SERVER_HOST=0.0.0.0 -e SERVER_PORT=8080 \
  -e DATABASE_HOST="$PROD_DB" -e REDIS_HOST="$PROD_REDIS" -e AUTO_SETUP=true \
  "$IMAGE"
echo "started."

echo "=== 5. wait ==="; sleep 12

echo "=== 6. 验证 ==="
echo -n "image : "; docker inspect -f '{{.Config.Image}}' "$APP"
echo -n "health: "; (curl -fsS http://127.0.0.1:4145/health && echo " ok") || echo "FAILED"
echo -n "ver   : "; docker exec "$APP" /app/sub2api --version 2>&1 | head -1 || true
echo "-- ps --"; docker ps --format '{{.Names}}\t{{.Image}}\t{{.Status}}' | grep -E 'sub2api'

echo "=== 7. 生产库迁移表 ==="
for t in prompt_audit_jobs prompt_audit_events ops_ingress_reject_aggregates auth_cache_invalidation_outbox; do
  printf '%s: ' "$t"
  docker exec "$PROD_DB" psql -U "$DB_USER" -d "$DB_NAME" -tc "select count(*) from $t;" 2>&1 | tr -d '\n ' || true
  echo
done
echo -n "185_group_reasoning_effort_policy.sql: "
docker exec "$PROD_DB" psql -U "$DB_USER" -d "$DB_NAME" -tc "select filename from schema_migrations where filename='185_group_reasoning_effort_policy.sql';" 2>&1 | tr -d '\n ' || true
echo
echo "groups reasoning columns:"
docker exec "$PROD_DB" psql -U "$DB_USER" -d "$DB_NAME" -tc "select column_name from information_schema.columns where table_name='groups' and column_name in ('max_reasoning_effort','reasoning_effort_mappings') order by column_name;" 2>&1 || true

echo "=== 8. 错误扫描 ==="
docker logs --tail=150 "$APP" 2>&1 | grep -iE 'error|panic|fatal' | head -20 || echo "(none)"

echo "=== done ==="
echo "回滚（如需）：SWITCH_IMAGE=$PREV_IMG bash deploy/switch-4145.sh"
