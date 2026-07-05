.PHONY: all data build preview dev clean ci-data ci-build ads-sync

DATA_DIR     := data-cli
FRONTEND_DIR := frontend
PORT         ?= 5173
CONCURRENCY  ?= 5

# ── 默认 ──────────────────────────────────────────────
# 抓取数据 → 安装依赖 → 构建 → 预览
all: data build
	cd $(FRONTEND_DIR) && pnpm preview --port $(PORT) --host 127.0.0.1

# ── 数据抓取 ──────────────────────────────────────────
data:
	cd $(DATA_DIR) && go run . -concurrency $(CONCURRENCY) -schedule-order priority -skip-unchanged

# ── 前端依赖 ──────────────────────────────────────────
install:
	cd $(FRONTEND_DIR) && pnpm install

install-ci:
	cd $(FRONTEND_DIR) && pnpm install --frozen-lockfile

# ── 前端构建 ──────────────────────────────────────────
build: install
	cd $(FRONTEND_DIR) && pnpm build

build-ci: install-ci
	cd $(FRONTEND_DIR) && pnpm build

# ── 预览 / 开发 ───────────────────────────────────────
preview:
	cd $(FRONTEND_DIR) && pnpm preview --port $(PORT) --host 127.0.0.1

dev: install
	cd $(FRONTEND_DIR) && pnpm dev --port $(PORT) --host 127.0.0.1

# ── 广告同步 ──────────────────────────────────────────
ads-sync:
	cd $(FRONTEND_DIR) && pnpm ads:sync

# ── 跳过数据抓取，直接构建 + 预览 ─────────────────────
build-only: install
	cd $(FRONTEND_DIR) && pnpm build && pnpm preview --port $(PORT) --host 127.0.0.1

# ── 清理 ──────────────────────────────────────────────
clean:
	rm -rf $(FRONTEND_DIR)/dist $(FRONTEND_DIR)/node_modules
