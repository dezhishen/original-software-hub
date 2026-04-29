---
name: "New Data Plugin"
description: "Generate a new data-cli plugin from provided clues, following this repository's existing plugin structure and registration pattern."
argument-hint: "填写软件名、线索来源、平台范围、是否允许先生成骨架。"
agent: "Data Plugin Generator"
---
基于当前仓库的 data-cli 插件模式，为目标软件生成一个新的数据插件。

请严格遵循现有目录结构、plugin 数据模型、平台分组方式和 main.go 注册方式，不要发明新接口。

按下面输入执行：

- 软件名称：{{软件名称}}
- 目标插件 ID：{{插件ID}}
- 官网或下载页：{{官网或下载页}}
- 已有线索：{{接口 / HTML / XHR / JSON / Release / 脚本片段}}
- 需要覆盖的平台：{{Windows / macOS / Linux / Android / iOS 等}}
- 是否允许先生成骨架：{{是 / 否}}
- 其他限制：{{例如不外查、只根据提供线索实现}}

输出要求：

1. 先说明参考了哪个现有插件模式。
2. 然后直接实施代码改动，包含必要的插件文件和 main.go 接入。
3. 如果线索不足，先生成骨架并明确指出仍需补充的抓取信息。
4. 最后说明做了哪些验证，或为什么暂时无法验证。