package visualstudio

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Visualstudio struct{}

func init() {
	plugin.Register(&Visualstudio{})
}

func (p *Visualstudio) Name() string {
	return "visual-studio"
}

func (p *Visualstudio) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "visual-studio",
				Name:            "Visual Studio",
				Description:     "集成开发环境。",
				Organization:    "Microsoft",
				OfficialWebsite: "https://visualstudio.microsoft.com",
				Icon:            "https://visualstudio.microsoft.com/favicon.ico",
				Tags:            []string{"开发工具", "IDE"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Visualstudio) Disabled() bool { return true }
