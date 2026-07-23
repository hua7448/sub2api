#!/usr/bin/env bash
# 4146 隔离试运行：在 trial 容器里跑指定镜像，验证启动 / 迁移 / 基本健康。
# 只替换 sub2api-image-gallery-trial，保留 trial DB/Redis，绝不触碰生产 4145。
#
# 用法：
#   bash scripts/trial-4146.sh                          # 默认跑本 fork 最新发布镜像
#   TRIAL_IMAGE=ghcr.io/hua7448/sub2api:<tag> bash scripts/trial-4146.sh
#
# 详见 docs/TRIAL_DEPLOYMENT_CN.md 与 docs/FORK_WORKFLOW_CN.md。
set -uo pipefail

IMAGE="${TRIAL_IMAGE:-ghcr.io/hua7448/sub2api:0.1.163-smartapi.1}"
TRIAL_APP="sub2api-image-gallery-trial"
TRIAL_DB="sub2api-image-gallery-postgres-trial"
TRIAL_REDIS="sub2api-image-gallery-redis-trial"
DB_USER="${DB_USER:-sub2api}"
DB_NAME="${DB_NAME:-sub2api}"

echo "=== 1. image ==="
docker image inspect "$IMAGE" >/dev/null && echo "image OK: $IMAGE" || { echo "image missing: $IMAGE  (先 docker pull)"; exit 1; }

echo "=== 2. trial network ==="
TRIAL_NET=""
for c in "$TRIAL_APP" "$TRIAL_DB" "$TRIAL_REDIS"; do
  TRIAL_NET="$(docker inspect -f '{{range $n,$_ := .NetworkSettings.Networks}}{{println $n}}{{end}}' "$c" 2>/dev/null | head -n1 || true)"
  [ -n "$TRIAL_NET" ] && { echo "network (from $c): $TRIAL_NET"; break; }
done
[ -n "$TRIAL_NET" ] || { echo "no trial network found (trial 容器是否存在？)"; exit 1; }

echo "=== 3. env ==="
docker inspect -f '{{range .Config.Env}}{{println .}}{{end}}' "$TRIAL_APP" > /tmp/trial.env
[ -s /tmp/trial.env ] || { echo "trial env empty"; exit 1; }
echo "env lines: $(wc -l < /tmp/trial.env)"

echo "=== 4. replace trial app (keep trial DB/Redis) ==="
docker rm -f "$TRIAL_APP" >/dev/null 2>&1 || true
docker run -d --name "$TRIAL_APP" --restart unless-stopped --network "$TRIAL_NET" -p 4146:8080 \
  --env-file /tmp/trial.env \
  -e SERVER_HOST=0.0.0.0 -e SERVER_PORT=8080 \
  -e DATABASE_HOST="$TRIAL_DB" -e REDIS_HOST="$TRIAL_REDIS" -e AUTO_SETUP=true \
  "$IMAGE"
echo "started."

echo "=== 5. wait ==="; sleep 10

echo "=== 6. health ==="
echo -n "image : "; docker inspect -f '{{.Config.Image}}' "$TRIAL_APP"
echo -n "health: "; (curl -fsS http://127.0.0.1:4146/health && echo " ok") || echo "FAILED"
echo -n "ver   : "; docker exec "$TRIAL_APP" /app/sub2api --version 2>&1 | head -1 || true

echo "=== 7. migrations (trial DB) ==="
for t in prompt_audit_jobs prompt_audit_events ops_ingress_reject_aggregates auth_cache_invalidation_outbox; do
  printf '%s: ' "$t"
  docker exec "$TRIAL_DB" psql -U "$DB_USER" -d "$DB_NAME" -tc "select count(*) from $t;" 2>&1 | tr -d '\n ' || true
  echo
done

echo "=== 8. error scan ==="
docker logs --tail=150 "$TRIAL_APP" 2>&1 | grep -iE 'migrat|error|panic|fatal' | head -20 || echo "(none)"
echo "=== done ==="
