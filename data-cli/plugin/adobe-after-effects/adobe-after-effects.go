package adobeaftereffects

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	adobeaftereffectsOfficialWebsite = "https://www.adobe.com/products/aftereffects"
	adobeaftereffectsIconURL         = "https://www.adobe.com/favicon.ico"
)

type Adobeaftereffects struct{}

func init() {
	plugin.Register(&Adobeaftereffects{})
}

func (p *Adobeaftereffects) Name() string {
	return "adobe-after-effects"
}

func (p *Adobeaftereffects) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "adobe-after-effects",
				Name:            "After Effects",
				Description:     "Adobe 出品的专业视觉特效与动态图形软件，广泛用于影视后期制作。",
				Organization:    "Adobe",
				OfficialWebsite: adobeaftereffectsOfficialWebsite,
				Icon:            adobeaftereffectsIconURL,
				Tags:            []string{"视频后期", "特效"},
				Categories:      []string{"design"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: adobeaftereffectsOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, adobeaftereffectsOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "After Effects 免费试用", URL: "https://www.adobe.com/products/aftereffects/free-trial-download.html"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "After Effects 免费试用", URL: "https://www.adobe.com/products/aftereffects/free-trial-download.html"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Adobeaftereffects) Disabled() bool { return false }
