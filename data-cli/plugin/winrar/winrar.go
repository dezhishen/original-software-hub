package winrar

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	winrarOfficialWebsite = "https://www.rarlab.com"
	winrarIconURL         = "https://www.rarlab.com/favicon.ico"
)

type Winrar struct{}

func init() {
	plugin.Register(&Winrar{})
}

func (p *Winrar) Name() string {
	return "winrar"
}

func (p *Winrar) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "winrar",
				Name:            "WinRAR",
				Description:     "高效的压缩与解压工具，支持 RAR 和 ZIP 等多种格式。",
				Organization:    "RARLAB",
				OfficialWebsite: winrarOfficialWebsite,
				Icon:            winrarIconURL,
				Tags:            []string{"压缩"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: winrarOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, winrarOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "WinRAR 下载页", URL: "https://www.rarlab.com/download.htm"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Winrar) Disabled() bool { return false }
