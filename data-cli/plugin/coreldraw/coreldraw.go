package coreldraw

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	coreldrawOfficialWebsite = "https://www.coreldraw.com"
	coreldrawIconURL         = "https://www.coreldraw.com/favicon.ico"
)

type Coreldraw struct{}

func init() {
	plugin.Register(&Coreldraw{})
}

func (p *Coreldraw) Name() string {
	return "coreldraw"
}

func (p *Coreldraw) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "coreldraw",
				Name:            "CorelDRAW",
				Description:     "Corel 出品的专业平面设计与矢量插图软件。",
				Organization:    "Corel",
				OfficialWebsite: coreldrawOfficialWebsite,
				Icon:            coreldrawIconURL,
				Tags:            []string{"设计"},
				Categories:      []string{"design"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: coreldrawOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, coreldrawOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "CorelDRAW 免费试用下载", URL: "https://www.coreldraw.com/cn/pages/free-trial/"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "CorelDRAW 免费试用下载", URL: "https://www.coreldraw.com/cn/pages/free-trial/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Coreldraw) Disabled() bool { return false }
