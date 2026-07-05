package weiyun

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	weiyunOfficialWebsite = "https://www.weiyun.com"
	weiyunIconURL         = "https://www.weiyun.com/favicon.ico"
)

type Weiyun struct{}

func init() {
	plugin.Register(&Weiyun{})
}

func (p *Weiyun) Name() string {
	return "weiyun"
}

func (p *Weiyun) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "weiyun",
				Name:            "腾讯微云",
				Description:     "腾讯出品的云存储服务，支持文件同步、备份和在线预览。",
				Organization:    "Tencent",
				OfficialWebsite: weiyunOfficialWebsite,
				Icon:            weiyunIconURL,
				Tags:            []string{"云存储"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: weiyunOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, weiyunOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "腾讯微云 Windows 下载", URL: "https://www.weiyun.com/download.html"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "腾讯微云 macOS 下载", URL: "https://www.weiyun.com/download.html"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Weiyun) Disabled() bool { return false }
