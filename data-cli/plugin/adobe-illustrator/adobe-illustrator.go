package adobeillustrator

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	adobeillustratorOfficialWebsite = "https://www.adobe.com/products/illustrator"
	adobeillustratorIconURL         = "https://www.adobe.com/favicon.ico"
)

type Adobeillustrator struct{}

func init() {
	plugin.Register(&Adobeillustrator{})
}

func (p *Adobeillustrator) Name() string {
	return "adobe-illustrator"
}

func (p *Adobeillustrator) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "adobe-illustrator",
				Name:            "Adobe Illustrator",
				Description:     "Adobe 出品的专业矢量图形设计软件，广泛用于插画、排版与 Logo 设计。",
				Organization:    "Adobe",
				OfficialWebsite: adobeillustratorOfficialWebsite,
				Icon:            adobeillustratorIconURL,
				Tags:            []string{"设计", "矢量图形"},
				Categories:      []string{"design"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: adobeillustratorOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, adobeillustratorOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Illustrator 免费试用", URL: "https://www.adobe.com/products/illustrator/free-trial-download.html"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Illustrator 免费试用", URL: "https://www.adobe.com/products/illustrator/free-trial-download.html"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Adobeillustrator) Disabled() bool { return false }
