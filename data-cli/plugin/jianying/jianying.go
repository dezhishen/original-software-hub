package jianying

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Jianying struct{}

func init() {
	plugin.Register(&Jianying{})
}

func (p *Jianying) Name() string {
	return "jianying"
}

func (p *Jianying) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "jianying",
				Name:            "剪映（桌面版）",
				Description:     "视频剪辑工具。",
				Organization:    "剪映",
				OfficialWebsite: "https://www.capcut.com",
				Icon:            "",
				Tags:            []string{"视频剪辑"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Jianying) Disabled() bool { return true }
