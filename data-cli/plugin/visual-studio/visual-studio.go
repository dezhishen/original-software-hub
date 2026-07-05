package visualstudio

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	visualstudioOfficialWebsite = "https://visualstudio.microsoft.com"
	visualstudioIconURL         = "https://visualstudio.microsoft.com/favicon.ico"
)

type Visualstudio struct{}

func init() {
	plugin.Register(&Visualstudio{})
}

func (p *Visualstudio) Name() string {
	return "visual-studio"
}

func (p *Visualstudio) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "visual-studio",
				Name:            "Visual Studio",
				Description:     "Microsoft 出品的专业集成开发环境（IDE），支持 C#、C++、Python 等多种语言。",
				Organization:    "Microsoft",
				OfficialWebsite: visualstudioOfficialWebsite,
				Icon:            visualstudioIconURL,
				Tags:            []string{"开发工具", "IDE"},
				Categories:      []string{"development"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: visualstudioOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, visualstudioOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "Visual Studio Community 2022", URL: "https://visualstudio.microsoft.com/thank-you-downloading-visual-studio/?sku=Community&rel=17"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Visual Studio for Mac", URL: "https://visualstudio.microsoft.com/vs/mac/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Visualstudio) Disabled() bool { return false }
