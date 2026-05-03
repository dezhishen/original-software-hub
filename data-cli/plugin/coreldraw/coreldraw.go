package coreldraw

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Coreldraw struct{}

func init() {
	plugin.Register(&Coreldraw{})
}

func (p *Coreldraw) Name() string {
	return "coreldraw"
}

func (p *Coreldraw) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "coreldraw",
				Name:            "CorelDRAW",
				Description:     "平面设计软件。",
				Organization:    "Corel",
				OfficialWebsite: "https://www.coreldraw.com",
				Icon:            "",
				Tags:            []string{"设计"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Coreldraw) Disabled() bool { return true }
