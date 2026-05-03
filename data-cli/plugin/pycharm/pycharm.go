package pycharm

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Pycharm struct{}

func init() {
	plugin.Register(&Pycharm{})
}

func (p *Pycharm) Name() string {
	return "pycharm"
}

func (p *Pycharm) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "pycharm",
				Name:            "PyCharm",
				Description:     "JetBrains 系列 IDE。",
				Organization:    "JetBrains",
				OfficialWebsite: "https://www.jetbrains.com/pycharm",
				Icon:            "",
				Tags:            []string{"开发工具", "IDE"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Pycharm) Disabled() bool { return true }
