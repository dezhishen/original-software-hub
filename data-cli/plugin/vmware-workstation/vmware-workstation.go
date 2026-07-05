package vmwareworkstation

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	vmwareworkstationOfficialWebsite = "https://www.vmware.com/products/workstation"
	vmwareworkstationIconURL         = "https://www.vmware.com/favicon.ico"
)

type Vmwareworkstation struct{}

func init() {
	plugin.Register(&Vmwareworkstation{})
}

func (p *Vmwareworkstation) Name() string {
	return "vmware-workstation"
}

func (p *Vmwareworkstation) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "vmware-workstation",
				Name:            "VMware Workstation Pro",
				Description:     "VMware 出品的专业桌面虚拟化软件，支持在单台 PC 上运行多个虚拟机。",
				Organization:    "VMware",
				OfficialWebsite: vmwareworkstationOfficialWebsite,
				Icon:            vmwareworkstationIconURL,
				Tags:            []string{"虚拟机"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: vmwareworkstationOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, vmwareworkstationOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "VMware Workstation 免费试用", URL: "https://www.vmware.com/products/workstation-pro/workstation-pro-evaluation.html"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "webpage", Label: "VMware Workstation Linux 下载", URL: "https://www.vmware.com/products/workstation-pro/workstation-pro-evaluation.html"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Vmwareworkstation) Disabled() bool { return false }
