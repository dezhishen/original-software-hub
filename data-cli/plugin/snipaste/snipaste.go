package snipaste

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	snipasteOfficialWebsite = "https://www.snipaste.com"
	snipasteIconURL         = "https://www.snipaste.com/favicon.ico"
)

type Snipaste struct{}

func init() {
	plugin.Register(&Snipaste{})
}

func (p *Snipaste) Name() string {
	return "snipaste"
}

func (p *Snipaste) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "snipaste",
				Name:            "Snipaste",
				Description:     "简单高效的截图工具，支持截图后贴图到屏幕上，方便高效办公。",
				Organization:    "Snipaste",
				OfficialWebsite: snipasteOfficialWebsite,
				Icon:            snipasteIconURL,
				Tags:            []string{"截图", "效率工具"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: snipasteOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, snipasteOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Snipaste 下载", URL: "https://www.snipaste.com/download"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Snipaste macOS 下载", URL: "https://www.snipaste.com/download"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Snipaste) Disabled() bool { return false }
