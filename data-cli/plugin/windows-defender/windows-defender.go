package windowsdefender

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Windowsdefender struct{}

func init() {
	plugin.Register(&Windowsdefender{})
}

func (p *Windowsdefender) Name() string {
	return "windows-defender"
}

func (p *Windowsdefender) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "windows-defender",
				Name:            "Windows Defender",
				Description:     "Windows 自带安全功能。",
				Organization:    "Microsoft",
				OfficialWebsite: "https://www.microsoft.com",
				Icon:            "",
				Tags:            []string{"安全防护"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Windowsdefender) Disabled() bool { return true }
