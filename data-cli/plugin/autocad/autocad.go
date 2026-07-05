package autocad

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	autocadOfficialWebsite = "https://www.autodesk.com/products/autocad"
	autocadIconURL         = "https://www.autodesk.com/favicon.ico"
)

type Autocad struct{}

func init() {
	plugin.Register(&Autocad{})
}

func (p *Autocad) Name() string {
	return "autocad"
}

func (p *Autocad) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "autocad",
				Name:            "AutoCAD",
				Description:     "Autodesk 出品的专业计算机辅助设计（CAD）软件，广泛用于建筑、工程与制造领域。",
				Organization:    "Autodesk",
				OfficialWebsite: autocadOfficialWebsite,
				Icon:            autocadIconURL,
				Tags:            []string{"设计", "CAD"},
				Categories:      []string{"design"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: autocadOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, autocadOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "AutoCAD 免费试用下载", URL: "https://www.autodesk.com/products/autocad/free-trial"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "AutoCAD 免费试用下载", URL: "https://www.autodesk.com/products/autocad/free-trial"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Autocad) Disabled() bool { return false }
