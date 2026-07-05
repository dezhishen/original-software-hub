package tencentpcmanager

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	tencentpcmanagerOfficialWebsite = "https://guanjia.qq.com"
	tencentpcmanagerIconURL         = "https://guanjia.qq.com/favicon.ico"
)

type Tencentpcmanager struct{}

func init() {
	plugin.Register(&Tencentpcmanager{})
}

func (p *Tencentpcmanager) Name() string {
	return "tencent-pc-manager"
}

func (p *Tencentpcmanager) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "tencent-pc-manager",
				Name:            "腾讯电脑管家",
				Description:     "腾讯出品的综合安全与管理软件，提供病毒查杀、系统加速和软件管理功能。",
				Organization:    "Tencent",
				OfficialWebsite: tencentpcmanagerOfficialWebsite,
				Icon:            tencentpcmanagerIconURL,
				Tags:            []string{"安全防护"},
				Categories:      []string{"security"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: tencentpcmanagerOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, tencentpcmanagerOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "腾讯电脑管家下载", URL: "https://guanjia.qq.com/download.html"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Tencentpcmanager) Disabled() bool { return false }
