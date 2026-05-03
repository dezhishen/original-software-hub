package autocad

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Autocad struct{}

func init() {
	plugin.Register(&Autocad{})
}

func (p *Autocad) Name() string {
	return "autocad"
}

func (p *Autocad) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "autocad",
				Name:            "AutoCAD",
				Description:     "计算机辅助设计软件。",
				Organization:    "",
				OfficialWebsite: "https://www.autodesk.com/products/autocad",
				Icon:            "",
				Tags:            []string{"设计", "CAD"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Autocad) Disabled() bool { return true }
