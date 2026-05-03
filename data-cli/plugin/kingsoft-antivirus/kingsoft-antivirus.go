package kingsoftantivirus

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Kingsoftantivirus struct{}

func init() {
	plugin.Register(&Kingsoftantivirus{})
}

func (p *Kingsoftantivirus) Name() string {
	return "kingsoft-antivirus"
}

func (p *Kingsoftantivirus) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "kingsoft-antivirus",
				Name:            "金山毒霸",
				Description:     "安全与防护工具。",
				Organization:    "Kingsoft",
				OfficialWebsite: "https://www.kingsoft.com",
				Icon:            "",
				Tags:            []string{"安全防护", "杀毒软件"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Kingsoftantivirus) Disabled() bool { return true }
