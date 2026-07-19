#!/usr/bin/env bash
# 4145 生产只读验证：版本 / 健康 / 4 个迁移表 / 错误扫描。不修改任何东西。
# 用于一键热更新（或镜像切换）之后确认生产状态。
set -uo pipefail

APP="${APP:-sub2api}"
PROD_DB="${PROD_DB:-sub2api-postgres}"
DB_USER="${DB_USER:-sub2api}"
DB_NAME="${DB_NAME:-sub2api}"
PORT="${PORT:-4145}"

echo "=== version ==="
echo -n "image : "; docker inspect -f '{{.Config.Image}}' "$APP" 2>&1
echo -n "binary: "; docker exec "$APP" /app/sub2api --version 2>&1 | head -1 || true

echo "=== health ==="
echo -n "health: "; (curl -fsS "http://127.0.0.1:${PORT}/health" && echo " ok") || echo "FAILED"

echo "=== migrations (prod DB) ==="
for t in prompt_audit_jobs prompt_audit_events ops_ingress_reject_aggregates auth_cache_invalidation_outbox; do
  printf '%s: ' "$t"
  docker exec "$PROD_DB" psql -U "$DB_USER" -d "$DB_NAME" -tc "select count(*) from $t;" 2>&1 | tr -d '\n ' || true
  echo
done

echo "=== error scan (last 150) ==="
docker logs --tail=150 "$APP" 2>&1 | grep -iE 'error|panic|fatal' | head -20 || echo "(none)"
echo "=== done ==="
