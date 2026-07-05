package virtualbox

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	virtualboxOfficialWebsite = "https://www.virtualbox.org"
	virtualboxIconURL         = "https://www.virtualbox.org/favicon.ico"
)

type Virtualbox struct{}

func init() {
	plugin.Register(&Virtualbox{})
}

func (p *Virtualbox) Name() string {
	return "virtualbox"
}

func (p *Virtualbox) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "virtualbox",
				Name:            "Oracle VM VirtualBox",
				Description:     "Oracle 出品的开源跨平台虚拟机软件，支持在单台电脑上运行多个操作系统。",
				Organization:    "Oracle",
				OfficialWebsite: virtualboxOfficialWebsite,
				Icon:            virtualboxIconURL,
				Tags:            []string{"虚拟机"},
				Categories:      []string{"virtualization"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: virtualboxOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, virtualboxOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "VirtualBox Windows 下载", URL: "https://www.virtualbox.org/wiki/Downloads"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "VirtualBox macOS 下载", URL: "https://www.virtualbox.org/wiki/Downloads"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "webpage", Label: "VirtualBox Linux 下载", URL: "https://www.virtualbox.org/wiki/Downloads"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Virtualbox) Disabled() bool { return false }
