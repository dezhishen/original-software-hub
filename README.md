# Original Software Hub

一个帮助用户快速找到常用软件**官方最新版本下载链接**的静态站点。所有数据均来源于各软件官网，定时自动抓取，无第三方镜像。

## 功能

- 聚合主流软件的官方下载地址与版本信息
- 支持 Windows / macOS / Linux / Android / iOS 多平台
- 数据每日自动更新（通过 GitHub Actions）
- 纯静态前端，无后端服务依赖

## 已收录软件

| 软件 | 分类 |
|------|------|
| 微信 | 即时通讯 |
| QQ | 即时通讯 |
| 企业微信 | 办公协作 |
| 钉钉 | 办公协作 |
| 腾讯会议 | 会议 |
| WPS Office | 办公套件 |
| 豆包 | AI 工具 |
| Chrome | 浏览器 |
| Firefox | 浏览器 |
| Steam | 游戏平台 |
| 阿里云盘 | 网盘 |
| 百度网盘 | 网盘 |
| 网易云音乐 | 音乐 |
| QQ 音乐 | 音乐 |
| 火绒安全 | 安全 |
| GitHub Desktop | 开发工具 |
| ToDesk | 远程控制 |
| UU 远程 | 远程控制 |

## 项目结构

```
├── data-cli/          # Go 数据抓取 CLI
│   ├── main.go        # 入口，插件调度与输出
│   └── plugin/        # 各软件抓取插件
├── frontend/          # 静态前端
│   ├── data/json/     # 生成的版本数据（JSON）
│   └── assets/        # 图标等静态资源
├── templates/         # 数据模板
└── .github/workflows/ # CI/CD 自动化流程
```

## 本地运行

**抓取数据**

```bash
cd data-cli
go run . -concurrency 5 -schedule-order priority -skip-unchanged
```

可选参数：

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-out` | `../frontend/data/json` | 输出目录 |
| `-plugins` | `all` | 指定插件，如 `weixin,qq` |
| `-concurrency` | `3` | 并发数 |
| `-schedule-order` | `priority` | 排序方式：`input` / `alpha` / `priority` |
| `-skip-unchanged` | `true` | 版本未变化时跳过写入与图标下载 |

**启动前端开发服务器**

```bash
cd frontend
npm install
npm run dev
```

## 新增软件插件

1. 在 `data-cli/plugin/<name>/` 下创建 `<name>.go`
2. 实现 `plugin.Plugin` 接口（`Name()` + `Fetch()`），并在 `init()` 中注册
3. 在 `data-cli/main.go` 中添加对应的 `_ ".../<name>"` 空导入

## CI/CD

通过 GitHub Actions（`.github/workflows/deploy-data.yml`）：

1. 从 `page` 分支恢复历史版本数据（用于 skip-unchanged 比对）
2. 运行 `data-cli` 抓取最新版本
3. 构建前端静态文件
4. 部署至 `page` 分支（GitHub Pages）

## License

MIT
