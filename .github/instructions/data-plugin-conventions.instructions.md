---
name: "Data Plugin Conventions"
description: "Use when 编写、修改、评审 data-cli 数据插件，或接入新的软件抓取逻辑。涵盖 plugin 结构、平台映射、main.go 注册、版本输出和线索不足时的处理方式。"
applyTo: ["data-cli/plugin/**/*.go", "data-cli/main.go"]
---
# Data Plugin Conventions

- 复用现有接口和数据结构：plugin.Plugin、plugin.SoftwareData、plugin.SoftwareItem、plugin.Version、plugin.PlatformsFromVariants。
- 不要为单个插件发明新的返回模型、额外的全局注册机制或特殊输出格式。
- 每个插件目录至少保持一个主文件；如果抓取逻辑明显复杂，再拆出 fetch_*.go 等辅助文件。
- 插件主文件应包含：常量元数据、插件类型、init 注册、Name 方法、Fetch 方法。

## Item 元数据

- ID 必须稳定、可用于 versions/<id>.json 文件名。
- Name、Description、Organization、OfficialWebsite、Icon 应尽量在插件内显式定义。
- Tags 只在确实有助于前端检索时添加，风格应与现有插件保持一致。

## Versions 和平台映射

- 优先把抓取结果先整理成 []plugin.Variant，再调用 plugin.PlatformsFromVariants。
- Platform 命名优先复用现有写法，例如 Windows、macOS、Linux、Android、iOS / iPadOS。
- Architecture 命名保持可读且稳定，例如 x64、x86、arm64、Intel、Apple Silicon、universal。
- Link.Type 只使用仓库里已有语义，例如 direct、store、webpage。
- 版本号、发布日期、官方发布页链接必须来自明确线索；线索不完整时可以返回骨架，但不要伪造值。

## main.go 接入

- 新插件默认应补充 data-cli/main.go 中的空白导入，保持可执行入口完整。
- 不要移除或重排无关插件导入，除非任务明确要求。

## 线索不足时

- 允许先生成可继续完善的骨架。
- 骨架也应保持可读的常量、目录结构和 Fetch 流程占位。
- 明确指出缺少的是接口响应、下载页结构、版本字段还是平台映射信息。

## 验证

- 修改后优先做最小验证，例如 go build 或项目既有构建命令。
- 若无法验证，明确说明原因，不要声称插件已可运行。