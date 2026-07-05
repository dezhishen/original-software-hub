package baidupinyin

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	baidupinyinOfficialWebsite = "https://shurufa.baidu.com"
	baidupinyinIconURL         = "https://shurufa.baidu.com/favicon.ico"
	baidupinyinDownloadURL     = "https://srf.baidu.com/?c=1&s=down&a=pc"
)

type Baidupinyin struct{}

func init() {
	plugin.Register(&Baidupinyin{})
}

func (p *Baidupinyin) Name() string {
	return "baidu-pinyin"
}

func (p *Baidupinyin) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "baidu-pinyin",
				Name:            "百度输入法",
				Description:     "百度出品的免费中文输入法，支持拼音、手写、语音等多种输入方式。",
				Organization:    "百度",
				OfficialWebsite: baidupinyinOfficialWebsite,
				Icon:            baidupinyinIconURL,
				Tags:            []string{"输入法"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: baidupinyinOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, baidupinyinOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "百度输入法 Windows 版", URL: baidupinyinDownloadURL},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "百度输入法 Mac 版", URL: "https://srf.baidu.com/input/mac.html"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "应用市场", URL: "https://shurufa.baidu.com/android"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/id414542983"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Baidupinyin) Disabled() bool { return false }
