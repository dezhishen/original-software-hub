package iqiyi

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	qiqiyiOfficialWebsite = "https://www.iqiyi.com"
	qiqiyiIconURL         = "https://www.iqiyi.com/favicon.ico"
)

type Iqiyi struct{}

func init() {
	plugin.Register(&Iqiyi{})
}

func (p *Iqiyi) Name() string {
	return "iqiyi"
}

func (p *Iqiyi) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "iqiyi",
				Name:            "爱奇艺",
				Description:     "爱奇艺出品的视频客户端，提供电影、电视剧、综艺等在线视频播放。",
				Organization:    "爱奇艺",
				OfficialWebsite: qiqiyiOfficialWebsite,
				Icon:            qiqiyiIconURL,
				Tags:            []string{"视频"},
				Categories:      []string{"media"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: qiqiyiOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, qiqiyiOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "爱奇艺 Windows 客户端下载", URL: "https://www.iqiyi.com/download/"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "爱奇艺 macOS 客户端下载", URL: "https://www.iqiyi.com/download/"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "应用市场", URL: "https://www.iqiyi.com/download/"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/id393765873"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Iqiyi) Disabled() bool { return false }
