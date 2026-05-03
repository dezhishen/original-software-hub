package adobepremierepro

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Adobepremierepro struct{}

func init() {
	plugin.Register(&Adobepremierepro{})
}

func (p *Adobepremierepro) Name() string {
	return "adobe-premiere-pro"
}

func (p *Adobepremierepro) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "adobe-premiere-pro",
				Name:            "Adobe Premiere Pro",
				Description:     "专业视频剪辑软件。",
				Organization:    "Adobe",
				OfficialWebsite: "https://www.adobe.com/products/premiere",
				Icon:            "",
				Tags:            []string{"视频剪辑"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Adobepremierepro) Disabled() bool { return true }
