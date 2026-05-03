package vmwareworkstation

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Vmwareworkstation struct{}

func init() {
	plugin.Register(&Vmwareworkstation{})
}

func (p *Vmwareworkstation) Name() string {
	return "vmware-workstation"
}

func (p *Vmwareworkstation) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "vmware-workstation",
				Name:            "VMware Workstation",
				Description:     "虚拟机软件。",
				Organization:    "VMware",
				OfficialWebsite: "https://www.vmware.com/products/workstation",
				Icon:            "",
				Tags:            []string{"虚拟机"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Vmwareworkstation) Disabled() bool { return true }
