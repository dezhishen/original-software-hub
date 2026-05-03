package p360safeguard

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type P360safeguard struct{}

func init() {
	plugin.Register(&P360safeguard{})
}

func (p *P360safeguard) Name() string {
	return "360-safe-guard"
}

func (p *P360safeguard) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "360-safe-guard",
				Name:            "360 安全卫士",
				Description:     "安全与防护工具。",
				Organization:    "",
				OfficialWebsite: "https://www.360.com",
				Icon:            "",
				Tags:            []string{"安全防护"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *P360safeguard) Disabled() bool { return true }
