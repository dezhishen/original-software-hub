package youku

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	youkuOfficialWebsite = "https://www.youku.com"
	youkuIconURL         = "https://www.youku.com/favicon.ico"
)

type Youku struct{}

func init() {
	plugin.Register(&Youku{})
}

func (p *Youku) Name() string {
	return "youku"
}

func (p *Youku) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "youku",
				Name:            "优酷",
				Description:     "优酷出品的视频平台客户端，提供电视剧、电影、综艺、动漫等丰富内容。",
				Organization:    "优酷",
				OfficialWebsite: youkuOfficialWebsite,
				Icon:            youkuIconURL,
				Tags:            []string{"视频"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: youkuOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, youkuOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "优酷 Windows 客户端下载", URL: "https://www.youku.com/download/"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "优酷 macOS 客户端下载", URL: "https://www.youku.com/download/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Youku) Disabled() bool { return false }
