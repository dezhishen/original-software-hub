#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$ROOT_DIR/frontend"
PORT="${PORT:-5173}"

cd "$FRONTEND_DIR"

if ! command -v npm >/dev/null 2>&1; then
  echo "Error: npm 未安装，请先安装 Node.js。" >&2
  exit 1
fi

if ! command -v node >/dev/null 2>&1; then
  echo "Error: node 未安装，请先安装 Node.js。" >&2
  exit 1
fi

echo "[frontend] 安装依赖..."
npm install

echo "[frontend] 构建资源..."
npm run build

echo "[frontend] 启动开发预览服务: http://127.0.0.1:${PORT}"
exec node ./scripts/dev-server.mjs
