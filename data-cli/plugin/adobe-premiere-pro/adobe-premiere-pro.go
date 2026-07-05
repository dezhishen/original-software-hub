package adobepremierepro

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	adobepremiereproOfficialWebsite = "https://www.adobe.com/products/premiere"
	adobepremiereproIconURL         = "https://www.adobe.com/favicon.ico"
)

type Adobepremierepro struct{}

func init() {
	plugin.Register(&Adobepremierepro{})
}

func (p *Adobepremierepro) Name() string {
	return "adobe-premiere-pro"
}

func (p *Adobepremierepro) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "adobe-premiere-pro",
				Name:            "Adobe Premiere Pro",
				Description:     "Adobe 出品的专业视频剪辑与后期制作软件，广泛用于影视与视频创作。",
				Organization:    "Adobe",
				OfficialWebsite: adobepremiereproOfficialWebsite,
				Icon:            adobepremiereproIconURL,
				Tags:            []string{"视频剪辑"},
				Categories:      []string{"design"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: adobepremiereproOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, adobepremiereproOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Premiere Pro 免费试用", URL: "https://www.adobe.com/products/premiere/free-trial-download.html"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Premiere Pro 免费试用", URL: "https://www.adobe.com/products/premiere/free-trial-download.html"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Adobepremierepro) Disabled() bool { return false }
