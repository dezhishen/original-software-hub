package p360antivirus

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type P360antivirus struct{}

func init() {
	plugin.Register(&P360antivirus{})
}

func (p *P360antivirus) Name() string {
	return "360-antivirus"
}

func (p *P360antivirus) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "360-antivirus",
				Name:            "360 杀毒",
				Description:     "安全与防护工具。",
				Organization:    "",
				OfficialWebsite: "https://www.360.com",
				Icon:            "",
				Tags:            []string{"安全防护", "杀毒软件"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *P360antivirus) Disabled() bool { return true }
