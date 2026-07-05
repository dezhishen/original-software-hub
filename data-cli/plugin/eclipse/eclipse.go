package eclipse

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	eclipseOfficialWebsite = "https://www.eclipse.org"
)

type Eclipse struct{}

func init() {
	plugin.Register(&Eclipse{})
}

func (p *Eclipse) Name() string {
	return "eclipse"
}

func (p *Eclipse) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "eclipse",
				Name:            "Eclipse IDE",
				Description:     "Eclipse 基金会出品的开源集成开发环境（IDE），支持 Java、C/C++、Python 等语言。",
				Organization:    "Eclipse Foundation",
				OfficialWebsite: eclipseOfficialWebsite,
				Icon:            "",
				Tags:            []string{"开发工具", "IDE"},
				Categories:      []string{"development"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: eclipseOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, eclipseOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Eclipse IDE 下载页", URL: "https://www.eclipse.org/downloads/"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Eclipse IDE 下载页", URL: "https://www.eclipse.org/downloads/"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Eclipse IDE 下载页", URL: "https://www.eclipse.org/downloads/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Eclipse) Disabled() bool { return false }
