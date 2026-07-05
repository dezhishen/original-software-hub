# Original Software Hub

一个帮助用户快速找到常用软件**官方最新版本下载链接**的静态站点。所有数据均来源于各软件官网，定时自动抓取，无第三方镜像。
这是一个纯 vibecoding 完成的实验项目。

## 功能

- 聚合主流软件的官方下载地址与版本信息
- 支持 Windows / macOS / Linux / Android / iOS 多平台
- 数据每日自动更新（通过 GitHub Actions）
- 纯静态前端，无后端服务依赖

## 项目结构

```
├── data-cli/            # Go 数据抓取 CLI
│   ├── main.go          # 入口，插件调度与输出
│   └── plugin/          # 各软件抓取插件
├── frontend/            # Vue 3 + Vite 前端
│   ├── src/             # 源码（组件、视图、composables）
│   ├── public/          # 静态资源（图标、PWA manifest）
│   ├── data/json/       # 生成的版本数据（JSON）
│   └── assets/          # 软件图标等静态资源
├── templates/           # 数据模板
└── .github/workflows/   # CI/CD 自动化流程
```

## 本地运行

**环境要求**

- [Go](https://go.dev/) 1.22+
- [Node.js](https://nodejs.org/) 22+
- [pnpm](https://pnpm.io)（推荐通过 `corepack enable` 启用）

**常用命令**

```bash
# 一键：抓取数据 → 安装依赖 → 构建 → 预览
make all

# 仅抓取最新软件数据
make data

# 跳过数据抓取，直接构建 + 预览
make build-only

# 启动 Vite 热重载开发服务器
make dev

# 清理构建产物
make clean
```

## 插件位置说明

- 插件目录：`data-cli/plugin/`
- 每个软件插件建议使用独立目录：`data-cli/plugin/<name>/<name>.go`
- 插件注册入口：`data-cli/main.go` 中的空导入列表（`_ ".../plugin/<name>"`）

新增插件时：

1. 在 `data-cli/plugin/<name>/` 下创建 `<name>.go`
2. 实现 `plugin.Plugin` 接口（`Name()` + `Fetch()`），并在 `init()` 中注册
3. 在 `data-cli/main.go` 中添加对应的空导入

## CI/CD

通过 GitHub Actions（`.github/workflows/deploy-data.yml`），每4小时自动执行：

1. 从 `page` 分支恢复历史版本数据（用于 skip-unchanged 比对）
2. `make data` 抓取最新版本
3. `make install-ci` 安装前端依赖
4. `make ads-sync` 同步公共广告
5. `make build-ci` 构建前端静态文件
6. 部署至 `page` 分支（GitHub Pages）

## License

MIT
