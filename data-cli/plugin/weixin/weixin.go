package weixin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	weixinOfficialWebsite = "https://weixin.qq.com/"
	weixinIconURL         = "https://res.wx.qq.com/a/wx_fed/web2/static/img/site_icon/favicon32.ico"

	weixinWindowsURL = "https://dldir1.qq.com/weixin/Windows/WeChatSetup.exe"
	weixinMacOSURL   = "https://dldir1.qq.com/weixin/mac/WeChat.dmg"
	weixinWebURL     = "https://weixin.qq.com/"
)

// WeChat implements plugin.Plugin for Tencent WeChat messenger.
type WeChat struct{}

func init() {
	plugin.Register(&WeChat{})
}

func (w *WeChat) Name() string {
	return "weixin"
}

func (w *WeChat) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate, err := detectReleaseDate([]string{weixinWindowsURL, weixinMacOSURL})
	if err != nil {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "weixin",
				Name:            "微信",
				Icon:            weixinIconURL,
				Description:     "腾讯微信，提供消息、通话、朋友圈等功能的通讯软件。",
				Organization:    "Tencent",
				OfficialWebsite: weixinOfficialWebsite,
				Tags:            []string{"即时通讯", "社交"},
			},
			Versions: []plugin.Version{
				{
					Version:     "Latest",
					ReleaseDate: releaseDate,
					OfficialURL: weixinOfficialWebsite,
					Variants: []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "微信安装包 (exe)", URL: weixinWindowsURL},
								{Type: "direct", Label: "在线更新版", URL: weixinWebURL},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: "微信安装包 (dmg)", URL: weixinMacOSURL},
							},
						},
						{
							Architecture: "通用",
							Platform:     "Web",
							Links: []plugin.Link{
								{Type: "direct", Label: "网页版微信", URL: weixinWebURL},
							},
						},
					},
				},
			},
		},
	}, nil
}

func detectReleaseDate(urls []string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	for _, u := range urls {
		resp, err := client.Head(u)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
				if t, err := time.Parse(time.RFC1123, lastModified); err == nil {
					return t.UTC().Format("2006-01-02"), nil
				}
			}
		}
	}
	return time.Now().UTC().Format("2006-01-02"), fmt.Errorf("detect weixin release date")
}
