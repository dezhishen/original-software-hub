package snipaste

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Snipaste struct{}

func init() {
	plugin.Register(&Snipaste{})
}

func (p *Snipaste) Name() string {
	return "snipaste"
}

func (p *Snipaste) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "snipaste",
				Name:            "Snipaste",
				Description:     "截图与贴图工具。",
				Organization:    "Snipaste",
				OfficialWebsite: "https://www.snipaste.com",
				Icon:            "",
				Tags:            []string{"截图", "效率工具"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Snipaste) Disabled() bool { return true }
