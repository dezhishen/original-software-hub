package notepadplusplus

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Notepadplusplus struct{}

func init() {
	plugin.Register(&Notepadplusplus{})
}

func (p *Notepadplusplus) Name() string {
	return "notepad-plus-plus"
}

func (p *Notepadplusplus) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "notepad-plus-plus",
				Name:            "Notepad++",
				Description:     "文本编辑器。",
				Organization:    "",
				OfficialWebsite: "https://notepad-plus-plus.org",
				Icon:            "",
				Tags:            []string{"开发工具", "文本编辑器"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Notepadplusplus) Disabled() bool { return true }
