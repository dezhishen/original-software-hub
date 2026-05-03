package vscode

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Vscode struct{}

func init() {
	plugin.Register(&Vscode{})
}

func (p *Vscode) Name() string {
	return "vscode"
}

func (p *Vscode) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "vscode",
				Name:            "Visual Studio Code",
				Description:     "代码编辑器。",
				Organization:    "Microsoft",
				OfficialWebsite: "https://code.visualstudio.com",
				Icon:            "https://code.visualstudio.com/favicon.ico",
				Tags:            []string{"开发工具", "代码编辑器"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Vscode) Disabled() bool { return true }
