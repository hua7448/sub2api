#!/usr/bin/env sh
set -eu

APP_NAME="${APP_NAME:-sub2api-site}"
IMAGE_NAME="${IMAGE_NAME:-sub2api-site}"
PORT="${PORT:-8080}"

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

if [ ! -f "out/index.html" ]; then
  echo "out/index.html not found. Run: npm run build"
  exit 1
fi

docker build -t "$IMAGE_NAME" -f deploy/Dockerfile .

if docker ps -a --format '{{.Names}}' | grep -Fxq "$APP_NAME"; then
  docker rm -f "$APP_NAME"
fi

docker run -d \
  --name "$APP_NAME" \
  --restart unless-stopped \
  -p "$PORT:80" \
  "$IMAGE_NAME"

echo "Started $APP_NAME at http://localhost:$PORT"

