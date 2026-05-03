package powerword

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Powerword struct{}

func init() {
	plugin.Register(&Powerword{})
}

func (p *Powerword) Name() string {
	return "powerword"
}

func (p *Powerword) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "powerword",
				Name:            "金山词霸",
				Description:     "翻译与词典工具。",
				Organization:    "Kingsoft",
				OfficialWebsite: "https://www.kingsoft.com",
				Icon:            "",
				Tags:            []string{"翻译", "词典"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Powerword) Disabled() bool { return true }
