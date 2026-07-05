package pycharm

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	pycharmOfficialWebsite = "https://www.jetbrains.com/pycharm"
	pycharmIconURL         = "https://www.jetbrains.com/favicon.ico"
)

type Pycharm struct{}

func init() {
	plugin.Register(&Pycharm{})
}

func (p *Pycharm) Name() string {
	return "pycharm"
}

func (p *Pycharm) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "pycharm",
				Name:            "PyCharm",
				Description:     "JetBrains 出品的专业 Python 集成开发环境（IDE），提供 Community（免费）和 Professional（付费）版本。",
				Organization:    "JetBrains",
				OfficialWebsite: pycharmOfficialWebsite,
				Icon:            pycharmIconURL,
				Tags:            []string{"开发工具", "IDE"},
				Categories:      []string{"development"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: pycharmOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, pycharmOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "PyCharm Community 下载", URL: "https://www.jetbrains.com/pycharm/download/#section=windows"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "PyCharm macOS 下载", URL: "https://www.jetbrains.com/pycharm/download/#section=mac"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "webpage", Label: "PyCharm Linux 下载", URL: "https://www.jetbrains.com/pycharm/download/#section=linux"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Pycharm) Disabled() bool { return false }
