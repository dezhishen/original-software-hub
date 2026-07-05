package kugoumusic

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	kugoumusicOfficialWebsite = "https://www.kugou.com"
	kugoumusicIconURL         = "https://www.kugou.com/favicon.ico"
)

type Kugoumusic struct{}

func init() {
	plugin.Register(&Kugoumusic{})
}

func (p *Kugoumusic) Name() string {
	return "kugou-music"
}

func (p *Kugoumusic) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "kugou-music",
				Name:            "酷狗音乐",
				Description:     "酷狗出品的音乐播放客户端，提供在线音乐、MV、电台等丰富的音娱服务。",
				Organization:    "酷狗",
				OfficialWebsite: kugoumusicOfficialWebsite,
				Icon:            kugoumusicIconURL,
				Tags:            []string{"音乐"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: kugoumusicOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, kugoumusicOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "酷狗音乐 Windows 下载", URL: "https://www.kugou.com/download/"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "应用市场", URL: "https://www.kugou.com/android/"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/id472966188"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Kugoumusic) Disabled() bool { return false }
