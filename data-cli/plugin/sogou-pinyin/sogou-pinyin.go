package sogoupinyin

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	sogoupinyinOfficialWebsite = "https://pinyin.sogou.com"
	sogoupinyinIconURL         = "https://pinyin.sogou.com/favicon.ico"
	sogoupinyinDownloadURL     = "https://pinyin.sogou.com/d/?f=1&t=0"
)

type Sogoupinyin struct{}

func init() {
	plugin.Register(&Sogoupinyin{})
}

func (p *Sogoupinyin) Name() string {
	return "sogou-pinyin"
}

func (p *Sogoupinyin) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "sogou-pinyin",
				Name:            "搜狗输入法",
				Description:     "搜狗出品的免费中文输入法，拥有丰富的词库和智能输入功能。",
				Organization:    "搜狗",
				OfficialWebsite: sogoupinyinOfficialWebsite,
				Icon:            sogoupinyinIconURL,
				Tags:            []string{"输入法"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: sogoupinyinOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, sogoupinyinOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "搜狗输入法 Windows 版", URL: sogoupinyinDownloadURL},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "搜狗输入法 Mac 版下载页", URL: "https://pinyin.sogou.com/mac/"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "应用市场", URL: "https://pinyin.sogou.com/android/"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/id917710510"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Sogoupinyin) Disabled() bool { return false }
