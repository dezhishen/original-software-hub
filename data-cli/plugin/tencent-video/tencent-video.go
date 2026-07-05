package tencentvideo

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	tencentvideoOfficialWebsite = "https://v.qq.com"
	tencentvideoIconURL         = "https://v.qq.com/favicon.ico"
)

type Tencentvideo struct{}

func init() {
	plugin.Register(&Tencentvideo{})
}

func (p *Tencentvideo) Name() string {
	return "tencent-video"
}

func (p *Tencentvideo) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "tencent-video",
				Name:            "腾讯视频",
				Description:     "腾讯出品的在线视频平台客户端，提供电视剧、电影、综艺、动漫等内容。",
				Organization:    "Tencent",
				OfficialWebsite: tencentvideoOfficialWebsite,
				Icon:            tencentvideoIconURL,
				Tags:            []string{"视频"},
				Categories:      []string{"media"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: tencentvideoOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, tencentvideoOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "腾讯视频 Windows 客户端下载", URL: "https://v.qq.com/download.html"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "腾讯视频 macOS 客户端下载", URL: "https://v.qq.com/download.html"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Tencentvideo) Disabled() bool { return false }
