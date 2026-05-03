package eclipse

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Eclipse struct{}

func init() {
	plugin.Register(&Eclipse{})
}

func (p *Eclipse) Name() string {
	return "eclipse"
}

func (p *Eclipse) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "eclipse",
				Name:            "Eclipse",
				Description:     "集成开发环境。",
				Organization:    "",
				OfficialWebsite: "https://www.eclipse.org",
				Icon:            "",
				Tags:            []string{"开发工具", "IDE"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Eclipse) Disabled() bool { return true }
