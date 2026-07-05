package sunlogin

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	sunloginOfficialWebsite = "https://sunlogin.oray.com"
	sunloginIconURL         = "https://sunlogin.oray.com/favicon.ico"
)

type Sunlogin struct{}

func init() {
	plugin.Register(&Sunlogin{})
}

func (p *Sunlogin) Name() string {
	return "sunlogin"
}

func (p *Sunlogin) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "sunlogin",
				Name:            "向日葵远程控制",
				Description:     "上海贝锐出品的远程控制软件，支持远程桌面、远程文件传输和远程运维。",
				Organization:    "上海贝锐",
				OfficialWebsite: sunloginOfficialWebsite,
				Icon:            sunloginIconURL,
				Tags:            []string{"远程控制"},
				Categories:      []string{"remote"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: sunloginOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, sunloginOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "向日葵 Windows 客户端下载", URL: "https://sunlogin.oray.com/download/"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "向日葵 macOS 客户端下载", URL: "https://sunlogin.oray.com/download/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Sunlogin) Disabled() bool { return false }
