package adobeaftereffects

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Adobeaftereffects struct{}

func init() {
	plugin.Register(&Adobeaftereffects{})
}

func (p *Adobeaftereffects) Name() string {
	return "adobe-after-effects"
}

func (p *Adobeaftereffects) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "adobe-after-effects",
				Name:            "After Effects",
				Description:     "视觉特效软件。",
				Organization:    "Adobe",
				OfficialWebsite: "https://www.adobe.com/products/aftereffects",
				Icon:            "",
				Tags:            []string{"视频后期", "特效"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Adobeaftereffects) Disabled() bool { return true }
