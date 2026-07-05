package jianying

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	jianyingOfficialWebsite = "https://www.capcut.com"
	jianyingIconURL         = "https://www.capcut.com/favicon.ico"
)

type Jianying struct{}

func init() {
	plugin.Register(&Jianying{})
}

func (p *Jianying) Name() string {
	return "jianying"
}

func (p *Jianying) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "jianying",
				Name:            "剪映（桌面版）",
				Description:     "字节跳动出品的免费视频剪辑工具，支持桌面端和移动端，提供丰富的剪辑和特效功能。",
				Organization:    "字节跳动",
				OfficialWebsite: jianyingOfficialWebsite,
				Icon:            jianyingIconURL,
				Tags:            []string{"视频剪辑"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: jianyingOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, jianyingOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "剪映 Windows 版下载", URL: "https://www.capcut.com/download"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "剪映 macOS 版下载", URL: "https://www.capcut.com/download"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "应用市场", URL: "https://www.capcut.com/download"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/id1500855883"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Jianying) Disabled() bool { return false }
