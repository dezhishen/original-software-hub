# 软件信息总表（info）

用途：用于补充和生成软件列表（software-list.json）所需的关键信息，并持续跟踪每个插件是否完成。

状态统计（当前）：
- 已完成：22
- 待补充：8
- 待确认：42
- 合计：72

字段约定：
- 版本接口/页面 URL：用于获取版本号与发布时间的来源地址。
- 下载地址映射：平台/架构到下载链接的映射规则（如 Windows x64、macOS Apple Silicon）。
- 提取规则：版本号、发布日期、链接提取逻辑（JSON 字段 / XPath / 正则 / 重定向规则等）。
- 状态：已完成、待补充 或 待确认。
## 已完成插件（20）

| 插件ID | 软件名称 | 官网 | 版本接口/页面 URL | 下载地址映射 | 提取规则 | 状态 |
|---|---|---|---|---|---|---|
| alipan | 阿里云盘 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| baidunetdisk | 百度网盘 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| chrome | Google Chrome | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| dingtalk | 钉钉 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| doubao | 豆包 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| firefox | Mozilla Firefox | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| github | GitHub（配置驱动） | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| huorong | 火绒 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| neteasecloudmusic | 网易云音乐 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| qq | QQ | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| qqmusic | QQ音乐 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| steam | Steam | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| tencentmeeting | 腾讯会议 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| todesk | ToDesk | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| uuremote | UU远程 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| wecom | 企业微信 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| weixin | 微信 | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| wps | WPS Office | 已在插件中配置 | 已实现 | 已实现 | 已实现 | 已完成 |
| 7zip | 7-Zip | https://www.7-zip.org | https://www.7-zip.org/（首页 Download 区块） | Windows x64/x86/arm64 直链（GitHub Releases） | XPath 定位下载表格 + 正则提取版本/日期（例如 `Download 7-Zip 26.01 (2026-04-27)`）；按架构关键词归类 | 已完成 |
| foxmail | Foxmail | https://www.foxmail.com | https://www.foxmail.com/win/ 与 https://www.foxmail.com/mac/ | Windows x64/x86 直链（/win/download 302）；macOS universal 直链（/mac/download 302） | Windows 页面正则提取 `最新版本：x.x.x (yyyy-mm-dd)`；Mac 从 dmg 文件名提取版本并优先使用 Last-Modified 作为发布日期 | 已完成 |
| anydesk | AnyDesk | https://anydesk.com.cn | https://anydesk.com.cn/zhs/downloads/windows（内嵌 JS `var downloads={...}`） | Windows x64 直链、macOS universal 直链、Linux deb/rpm/tar.gz、Android Google Play、iOS App Store | 正则提取 `var downloads=` JSON 块后 `encoding/json` 解析；各平台 version 字段直接可用 | 已完成 |
| evernote | Evernote | https://evernote.com | https://evernote.com/download（页面含 `Evernote-latest.exe` / `.dmg` 直链） | Windows x64 直链、macOS universal 直链、Android Google Play、iOS App Store | 下载页提取 `latest` 直链；HEAD 请求 `Last-Modified` 响应头作为发布日期和版本标识（无版本号字段） | 已完成 |

## 待确认插件（44）

说明：以下条目已真实访问官网 URL 并抓到发行线索（版本候选、日期候选或下载链接候选），需人工确认后写入插件抓取代码。

| 插件ID | 软件名称 | 官网 | 版本接口/页面 URL | 下载地址映射 | 提取规则 | 状态 |
|---|---|---|---|---|---|---|
| 360-antivirus | 360 杀毒 | https://sd.360.cn/ | https://sd.360.cn/（版本候选: 7.0; 日期候选: 2020-12-11） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| 360-browser | 360 浏览器 | https://browser.360.cn | https://browser.360.cn（版本候选: 1.360; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| 360-zip | 360 压缩 | https://yasuo.360.cn | https://yasuo.360.cn（版本候选: v8.0.1; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| adobe-after-effects | After Effects | https://www.adobe.com/products/aftereffects | https://www.adobe.com/products/aftereffects（版本候选: 未识别; 日期候选: 2026-01-12） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| adobe-illustrator | Adobe Illustrator | https://www.adobe.com/products/illustrator | https://www.adobe.com/products/illustrator（版本候选: 未识别; 日期候选: 2026-01-12） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| adobe-photoshop | Adobe Photoshop | https://www.adobe.com/products/photoshop | https://www.adobe.com/products/photoshop（版本候选: 未识别; 日期候选: 2026-01-12） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| adobe-premiere-pro | Adobe Premiere Pro | https://www.adobe.com/products/premiere | https://www.adobe.com/products/premiere（版本候选: 未识别; 日期候选: 2026-01-12） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| coreldraw | CorelDRAW | https://www.coreldraw.com | https://www.coreldraw.com（版本候选: 0.7; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| eclipse | Eclipse | https://www.eclipse.org | https://www.eclipse.org（版本候选: 0.144.2; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| edge | Microsoft Edge | https://www.microsoft.com/edge | https://www.microsoft.com/edge（版本候选: 7.25; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| foxit-reader | 福昕阅读器（Foxit Reader） | https://www.foxitsoftware.com/pdf-reader | https://www.foxitsoftware.com/pdf-reader（版本候选: 1.0; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| git | Git | https://git-scm.com | https://git-scm.com（版本候选: 0.155.3; 日期候选: 2026-04-20） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| intellij-idea | IntelliJ IDEA | https://www.jetbrains.com/idea | https://www.jetbrains.com/idea（版本候选: 2.47; 日期候选: 2021-06-01） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| iqiyi | 爱奇艺 | https://www.iqiyi.com | https://www.iqiyi.com（版本候选: 17.051.25240; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| jianying | 剪映（桌面版） | https://www.capcut.com | https://www.capcut.com（版本候选: 1.3125; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| kingsoft-antivirus | 金山毒霸 | https://www.kingsoft.com | https://www.kingsoft.com（版本候选: 1.0.0; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| kingsoft-pdf | 金山 PDF | https://www.kingsoft.com | https://www.kingsoft.com（版本候选: 1.0.0; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| kugou-music | 酷狗音乐 | https://www.kugou.com | https://www.kugou.com（版本候选: 0.0.8; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| microsoft-office | Microsoft Office | https://www.microsoft.com/office | https://www.microsoft.com/office（版本候选: v1.25.0; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| microsoft-pinyin | 微软拼音 | https://www.microsoft.com | https://www.microsoft.com（版本候选: v1.25.0; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| notepad-plus-plus | Notepad++ | https://notepad-plus-plus.org | https://notepad-plus-plus.org（版本候选: 1.0; 日期候选: 2097-11.13） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| onenote | OneNote | https://www.microsoft.com/onenote | https://www.microsoft.com/onenote（版本候选: v1.25.0; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| outlook | Microsoft Outlook | https://www.microsoft.com/outlook | https://www.microsoft.com/outlook（版本候选: v1.25.0; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| powerword | 金山词霸 | https://www.kingsoft.com | https://www.kingsoft.com（版本候选: 1.0.0; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| pycharm | PyCharm | https://www.jetbrains.com/pycharm | https://www.jetbrains.com/pycharm（版本候选: 2.53; 日期候选: 2024.3.2） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| qq-browser | QQ 浏览器 | https://browser.qq.com | https://browser.qq.com（版本候选: 1.0; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| snipaste | Snipaste | https://www.snipaste.com | https://www.snipaste.com（版本候选: 3.4.1; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| sogou-browser | 搜狗高速浏览器 | https://www.sogou.com | https://www.sogou.com（版本候选: 113.246.105.143; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| sogou-pinyin | 搜狗输入法 | https://pinyin.sogou.com | https://pinyin.sogou.com（版本候选: 1.0; 日期候选: 未识别） | 下载候选: http://ime.gtimg.com/pc/sogou_wubi_5.5f.exe（待平台归类） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| sunlogin | 向日葵远程控制 | https://sunlogin.oray.com | https://sunlogin.oray.com（版本候选: 1.0; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| teamviewer | TeamViewer | https://www.teamviewer.com | https://www.teamviewer.com（版本候选: 05.292; 日期候选: 2026-04-30） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| tencent-pc-manager | 腾讯电脑管家 | https://guanjia.qq.com | https://guanjia.qq.com（版本候选: 0.5; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| tencent-video | 腾讯视频 | https://v.qq.com | https://v.qq.com（版本候选: v2.9.6; 日期候选: 2026-04-27） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| virtualbox | Oracle VM VirtualBox | https://www.virtualbox.org | https://www.virtualbox.org（版本候选: 7.1; 日期候选: 2026-04-21） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| visual-studio | Visual Studio | https://visualstudio.microsoft.com | https://visualstudio.microsoft.com（版本候选: 0.3.3; 日期候选: 2019-08-19） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| vlc | VLC 媒体播放器 | https://www.videolan.org/vlc | https://www.videolan.org/vlc（版本候选: 3.7.2; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| vmware-workstation | VMware Workstation | https://www.vmware.com/products/desktop-hypervisor/workstation-and-fusion | https://www.vmware.com/products/desktop-hypervisor/workstation-and-fusion（版本候选: 0.100104; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| vscode | Visual Studio Code | https://code.visualstudio.com | https://code.visualstudio.com（版本候选: 1.0; 日期候选: 2026-03-17） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| weiyun | 腾讯微云 | https://www.weiyun.com | https://www.weiyun.com（版本候选: 2.1.27262305; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| winrar | WinRAR | https://www.rarlab.com | https://www.rarlab.com（版本候选: 4.01; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| youdao-dict | 有道词典（桌面版） | https://youdao.com | https://youdao.com（版本候选: 11.2; 日期候选: 未识别） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |
| youku | 优酷 | https://www.youku.com | https://www.youku.com（版本候选: 0.8; 日期候选: 2025-05-09） | 未直接抓到安装包链接（需二次页面跟进） | HTML 正则首轮提取；下一步补充接口探测与 DOM 定位 | 待确认 |

## 待补充插件（8）

说明：以下条目已真实访问官网 URL，但当前未抓到稳定发行线索，暂保留待补充。

| 插件ID | 软件名称 | 官网 | 版本接口/页面 URL | 下载地址映射 | 提取规则 | 状态 |
|---|---|---|---|---|---|---|
| 360-safe-guard | 360 安全卫士 | https://www.360.com | 待补充 | 待补充 | 待补充 | 待补充 |
| adobe-acrobat-reader | Adobe Acrobat Reader | https://www.adobe.com/acrobat/pdf-reader | 待补充 | 待补充 | 待补充 | 待补充 |
| autocad | AutoCAD | https://www.autodesk.com/products/autocad | 待补充 | 待补充 | 待补充 | 待补充 |
| baidu-pinyin | 百度输入法 | https://shurufa.baidu.com | 待补充 | 待补充 | 待补充 | 待补充 |
| iflytek-pinyin | 讯飞输入法 | https://shurufa.iflytek.com | 待补充 | 待补充 | 待补充 | 待补充 |
| onedrive | OneDrive | https://www.microsoft.com/onedrive | 待补充 | 待补充 | 待补充 | 待补充 |
| potplayer | PotPlayer | https://potplayer.daum.net | 待补充 | 待补充 | 待补充 | 待补充 |
| windows-defender | Windows Defender | https://www.microsoft.com | 待补充 | 待补充 | 待补充 | 待补充 |
