package adobephotoshop

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Adobephotoshop struct{}

func init() {
	plugin.Register(&Adobephotoshop{})
}

func (p *Adobephotoshop) Name() string {
	return "adobe-photoshop"
}

func (p *Adobephotoshop) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "adobe-photoshop",
				Name:            "Adobe Photoshop",
				Description:     "图像编辑软件。",
				Organization:    "Adobe",
				OfficialWebsite: "https://www.adobe.com/products/photoshop",
				Icon:            "",
				Tags:            []string{"图像处理", "设计"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Adobephotoshop) Disabled() bool { return true }
