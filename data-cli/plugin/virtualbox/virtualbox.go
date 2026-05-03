package virtualbox

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Virtualbox struct{}

func init() {
	plugin.Register(&Virtualbox{})
}

func (p *Virtualbox) Name() string {
	return "virtualbox"
}

func (p *Virtualbox) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "virtualbox",
				Name:            "Oracle VM VirtualBox",
				Description:     "虚拟机软件。",
				Organization:    "Oracle",
				OfficialWebsite: "https://www.virtualbox.org",
				Icon:            "",
				Tags:            []string{"虚拟机"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Virtualbox) Disabled() bool { return true }
