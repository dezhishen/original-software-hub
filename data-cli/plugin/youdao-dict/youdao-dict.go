package youdaodict

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Youdaodict struct{}

func init() {
	plugin.Register(&Youdaodict{})
}

func (p *Youdaodict) Name() string {
	return "youdao-dict"
}

func (p *Youdaodict) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "youdao-dict",
				Name:            "有道词典（桌面版）",
				Description:     "翻译与词典工具。",
				Organization:    "有道",
				OfficialWebsite: "https://youdao.com",
				Icon:            "",
				Tags:            []string{"翻译", "词典"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Youdaodict) Disabled() bool { return true }
