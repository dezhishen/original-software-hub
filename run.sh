#!/usr/bin/env bash
set -euo pipefail

# 用法:
#   ./run.sh              — 抓取数据 + 构建 + 预览
#   ./run.sh skip-data    — 跳过数据抓取，直接构建 + 预览
#   ./run.sh dev          — 跳过数据抓取，启动 Vite 热重载开发服务器

MODE="${1:-}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="$ROOT_DIR/frontend"
DATA_DIR="$ROOT_DIR/data-cli"
PORT="${PORT:-5173}"

# ── 依赖检查 ──────────────────────────────────────────
if ! command -v pnpm >/dev/null 2>&1; then
  echo "Error: pnpm 未找到。请安装 Node.js 后运行 corepack enable。" >&2
  exit 1
fi

if [[ "$MODE" != "skip-data" && "$MODE" != "dev" ]]; then
  if ! command -v go >/dev/null 2>&1; then
    echo "Error: go 未找到。请先安装 Go。" >&2
    exit 1
  fi
fi

# ── 数据抓取 ──────────────────────────────────────────
if [[ "$MODE" != "skip-data" && "$MODE" != "dev" ]]; then
  echo "[data-cli] 生成数据..."
  cd "$DATA_DIR"
  go run . -concurrency 5 -schedule-order priority -skip-unchanged
  cd "$ROOT_DIR"
fi

# ── 前端 ──────────────────────────────────────────────
cd "$FRONTEND_DIR"

echo "[frontend] 安装依赖..."
pnpm install

if [[ "$MODE" == "dev" ]]; then
  echo "[frontend] 启动开发服务器 (热重载): http://127.0.0.1:${PORT}"
  exec pnpm dev --port "$PORT" --host 127.0.0.1
else
  echo "[frontend] 构建资源..."
  pnpm build

  echo "[frontend] 启动预览服务: http://127.0.0.1:${PORT}"
  exec pnpm preview --port "$PORT" --host 127.0.0.1
fi
