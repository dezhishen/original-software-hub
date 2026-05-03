package git

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Git struct{}

func init() {
	plugin.Register(&Git{})
}

func (p *Git) Name() string {
	return "git"
}

func (p *Git) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "git",
				Name:            "Git",
				Description:     "版本控制工具。",
				Organization:    "",
				OfficialWebsite: "https://git-scm.com",
				Icon:            "https://git-scm.com/favicon.ico",
				Tags:            []string{"开发工具", "版本控制"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Git) Disabled() bool { return true }
