package intellijidea

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Intellijidea struct{}

func init() {
	plugin.Register(&Intellijidea{})
}

func (p *Intellijidea) Name() string {
	return "intellij-idea"
}

func (p *Intellijidea) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "intellij-idea",
				Name:            "IntelliJ IDEA",
				Description:     "JetBrains 系列 IDE。",
				Organization:    "JetBrains",
				OfficialWebsite: "https://www.jetbrains.com/idea",
				Icon:            "",
				Tags:            []string{"开发工具", "IDE"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Intellijidea) Disabled() bool { return true }
