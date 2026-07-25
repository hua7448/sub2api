#!/usr/bin/env bash
# 4146 trial 只读验证：健康 / 版本 / migration 完整性 / 错误扫描 / (可选)内嵌模型抽样。
# 通用：自动扫描当前仓库 backend/migrations/*.sql，逐一对比 trial DB 的 schema_migrations，
# 不随版本写死 migration 清单。每次 trial 启动后、合入 main 前，必须执行并通过。
#
# 用法（在 /root/sub2api-src，checkout 到待验证的分支后）：
#   bash deploy/verify-4146.sh                            # 基础检查
#   CHECK_MODEL=claude-opus-5 bash deploy/verify-4146.sh  # 额外确认某模型定价已内嵌进二进制
#   MIGRATIONS_DIR=/path/to/migrations bash deploy/verify-4146.sh
#
# 退出码：0 = 全绿，可继续发布；1 = 有失败，禁止继续。
# 详见 docs/TRIAL_DEPLOYMENT_CN.md 与 docs/FORK_WORKFLOW_CN.md。
set -uo pipefail

TRIAL_APP="${TRIAL_APP:-sub2api-image-gallery-trial}"
TRIAL_DB="${TRIAL_DB:-sub2api-image-gallery-postgres-trial}"
DB_USER="${DB_USER:-sub2api}"
DB_NAME="${DB_NAME:-sub2api}"
PORT="${PORT:-4146}"
CHECK_MODEL="${CHECK_MODEL:-}"

# 默认从仓库根取 migrations 目录（trial 镜像必须由同一 checkout 构建才能对齐）
if [ -z "${MIGRATIONS_DIR:-}" ]; then
  root="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  [ -n "$root" ] && MIGRATIONS_DIR="$root/backend/migrations"
fi
MIGRATIONS_DIR="${MIGRATIONS_DIR:-./backend/migrations}"
VERSION_FILE="${MIGRATIONS_DIR%/migrations}/cmd/server/VERSION"

PASS=0; FAIL=0
green() { printf '  \033[32mPASS\033[0m %s\n' "$1"; PASS=$((PASS+1)); }
red()   { printf '  \033[31mFAIL\033[0m %s\n' "$1"; FAIL=$((FAIL+1)); }
note()  { printf '   .     %s\n' "$1"; }

echo "=== verify 4146 trial ==="
note "app=$TRIAL_APP db=$TRIAL_DB port=$PORT migrations=$MIGRATIONS_DIR"

echo "=== 1. container & image ==="
img="$(docker inspect -f '{{.Config.Image}}' "$TRIAL_APP" 2>/dev/null || true)"
[ -n "$img" ] && note "image: $img" || red "无法 inspect $TRIAL_APP（容器不存在？）"

echo "=== 2. health & binary version ==="
if curl -fsS "http://127.0.0.1:${PORT}/health" >/dev/null 2>&1; then green "health ok (/health)"; else red "health FAILED"; fi
ver="$(docker exec "$TRIAL_APP" /app/sub2api --version 2>&1 | grep -oE 'Sub2API [0-9][0-9.]*' | head -1 || true)"
[ -n "$ver" ] && green "binary version: $ver" || red "无法读取二进制版本"
cv="$(cat "$VERSION_FILE" 2>/dev/null | tr -d '\n ' || true)"
if [ -n "$cv" ] && [ -n "$ver" ]; then
  case "$ver" in
    *"$cv"*) green "checkout VERSION ($cv) == binary" ;;
    *) red "checkout VERSION ($cv) != binary ($ver) — trial 镜像可能不是当前 checkout 构建" ;;
  esac
fi

echo "=== 3. migration completeness (expected from backend/migrations) ==="
if [ ! -d "$MIGRATIONS_DIR" ]; then
  red "MIGRATIONS_DIR 不存在: $MIGRATIONS_DIR"
else
  for f in $(ls -1 "$MIGRATIONS_DIR"/*.sql 2>/dev/null | xargs -n1 basename 2>/dev/null | sort); do
    ap="$(docker exec "$TRIAL_DB" psql -U "$DB_USER" -d "$DB_NAME" -tAc "select 1 from schema_migrations where filename='$f';" 2>/dev/null | tr -d '\n ' || true)"
    if [ "$ap" = "1" ]; then green "$f"; else red "$f 未应用"; fi
  done
fi

echo "=== 4. error scan (last 200 log lines) ==="
errs="$(docker logs --tail=200 "$TRIAL_APP" 2>&1 | grep -iE 'panic|fatal|level=error|migrat.*fail' | head -10 || true)"
if [ -z "$errs" ]; then green "日志无 panic/fatal/error/migration-fail"; else red "日志含错误:"; printf '%s\n' "$errs"; fi

echo "=== 5. pricing embed sample ==="
if [ -n "$CHECK_MODEL" ]; then
  cnt="$(docker exec "$TRIAL_APP" sh -c "strings /app/sub2api 2>/dev/null | grep -F -c -- '$CHECK_MODEL'" 2>/dev/null | tr -d '\n ' || true)"
  [ "${cnt:-0}" -ge 1 ] && green "二进制内嵌模型 '$CHECK_MODEL' (count=$cnt)" || red "二进制未内嵌模型 '$CHECK_MODEL'"
else
  note "跳过（设 CHECK_MODEL=<模型名> 启用，如 claude-opus-5）"
fi

echo "=== summary: $PASS PASS / $FAIL FAIL ==="
if [ "$FAIL" -eq 0 ]; then
  echo "ALL GREEN — 可继续合入 main / 发布"
  exit 0
else
  echo "HAS FAILURES — 禁止继续发布"
  exit 1
fi
