package intellijidea

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	intellijideaOfficialWebsite = "https://www.jetbrains.com/idea"
	intellijideaIconURL         = "https://www.jetbrains.com/favicon.ico"
)

type Intellijidea struct{}

func init() {
	plugin.Register(&Intellijidea{})
}

func (p *Intellijidea) Name() string {
	return "intellij-idea"
}

func (p *Intellijidea) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "intellij-idea",
				Name:            "IntelliJ IDEA",
				Description:     "JetBrains 出品的专业 Java/Kotlin 集成开发环境（IDE），提供 Community（免费）和 Ultimate（付费）版本。",
				Organization:    "JetBrains",
				OfficialWebsite: intellijideaOfficialWebsite,
				Icon:            intellijideaIconURL,
				Tags:            []string{"开发工具", "IDE"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: intellijideaOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, intellijideaOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "IntelliJ IDEA Community 下载", URL: "https://www.jetbrains.com/idea/download/#section=windows"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "IntelliJ IDEA macOS 下载", URL: "https://www.jetbrains.com/idea/download/#section=mac"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "webpage", Label: "IntelliJ IDEA Linux 下载", URL: "https://www.jetbrains.com/idea/download/#section=linux"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Intellijidea) Disabled() bool { return false }
