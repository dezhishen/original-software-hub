package adobephotoshop

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	adobephotoshopOfficialWebsite = "https://www.adobe.com/products/photoshop"
	adobephotoshopIconURL         = "https://www.adobe.com/favicon.ico"
)

type Adobephotoshop struct{}

func init() {
	plugin.Register(&Adobephotoshop{})
}

func (p *Adobephotoshop) Name() string {
	return "adobe-photoshop"
}

func (p *Adobephotoshop) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "adobe-photoshop",
				Name:            "Adobe Photoshop",
				Description:     "Adobe 出品的专业图像编辑与设计软件，广泛用于摄影、设计与数字艺术。",
				Organization:    "Adobe",
				OfficialWebsite: adobephotoshopOfficialWebsite,
				Icon:            adobephotoshopIconURL,
				Tags:            []string{"图像处理", "设计"},
				Categories:      []string{"design"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: adobephotoshopOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, adobephotoshopOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Photoshop 免费试用", URL: "https://www.adobe.com/products/photoshop/free-trial-download.html"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Photoshop 免费试用", URL: "https://www.adobe.com/products/photoshop/free-trial-download.html"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Adobephotoshop) Disabled() bool { return false }
