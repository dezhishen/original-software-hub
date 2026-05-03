package adobeillustrator

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Adobeillustrator struct{}

func init() {
	plugin.Register(&Adobeillustrator{})
}

func (p *Adobeillustrator) Name() string {
	return "adobe-illustrator"
}

func (p *Adobeillustrator) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "adobe-illustrator",
				Name:            "Adobe Illustrator",
				Description:     "矢量图形软件。",
				Organization:    "Adobe",
				OfficialWebsite: "https://www.adobe.com/products/illustrator",
				Icon:            "",
				Tags:            []string{"设计", "矢量图形"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Adobeillustrator) Disabled() bool { return true }
