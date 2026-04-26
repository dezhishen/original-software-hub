package weixin

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	weixinOfficialWebsite = "https://weixin.qq.com/"
	weixinUpdatesURL      = "https://weixin.qq.com/updates"
	weixinIconURL         = "https://res.wx.qq.com/a/wx_fed/web2/static/img/site_icon/favicon32.ico"
	weixinWebURL          = "https://weixin.qq.com/"
)

var (
	reWinURL = regexp.MustCompile(`https://dldir[^"]+/WeChatWin_([\d.]+)\.exe`)
	reMacURL = regexp.MustCompile(`https://dldir[^"]+/WeChatMac_([\d.]+)\.dmg`)
	reDate   = regexp.MustCompile(`发布日期[：: ]*(\d{4}-\d{2}-\d{2})`)
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
	info, err := fetchUpdatesInfo()
	if err != nil {
		return nil, fmt.Errorf("fetch weixin updates: %w", err)
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
					Version:     info.version,
					ReleaseDate: info.releaseDate,
					OfficialURL: weixinUpdatesURL,
					Variants: []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: fmt.Sprintf("微信 %s 安装包 (exe)", info.version), URL: info.windowsURL},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: fmt.Sprintf("微信 %s 安装包 (dmg)", info.version), URL: info.macURL},
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

type updatesInfo struct {
	version     string
	releaseDate string
	windowsURL  string
	macURL      string
}

func fetchUpdatesInfo() (*updatesInfo, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(weixinUpdatesURL)
	if err != nil {
		return nil, fmt.Errorf("http get %s: %w", weixinUpdatesURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	html := string(body)

	winMatch := reWinURL.FindStringSubmatch(html)
	if len(winMatch) < 2 {
		return nil, fmt.Errorf("windows download URL not found on updates page")
	}
	macMatch := reMacURL.FindStringSubmatch(html)
	if len(macMatch) < 2 {
		return nil, fmt.Errorf("macos download URL not found on updates page")
	}

	version := winMatch[1]
	windowsURL := winMatch[0]
	macURL := macMatch[0]

	releaseDate := time.Now().UTC().Format("2006-01-02")
	if dateMatch := reDate.FindStringSubmatch(html); len(dateMatch) >= 2 {
		releaseDate = dateMatch[1]
	}

	return &updatesInfo{
		version:     version,
		releaseDate: releaseDate,
		windowsURL:  windowsURL,
		macURL:      macURL,
	}, nil
}
