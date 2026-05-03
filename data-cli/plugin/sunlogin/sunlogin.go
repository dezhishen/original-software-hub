package sunlogin

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Sunlogin struct{}

func init() {
	plugin.Register(&Sunlogin{})
}

func (p *Sunlogin) Name() string {
	return "sunlogin"
}

func (p *Sunlogin) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "sunlogin",
				Name:            "向日葵远程控制",
				Description:     "远程控制工具。",
				Organization:    "向日葵",
				OfficialWebsite: "https://sunlogin.oray.com",
				Icon:            "",
				Tags:            []string{"远程控制"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Sunlogin) Disabled() bool { return true }
