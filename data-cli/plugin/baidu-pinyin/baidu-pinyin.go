package baidupinyin

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Baidupinyin struct{}

func init() {
	plugin.Register(&Baidupinyin{})
}

func (p *Baidupinyin) Name() string {
	return "baidu-pinyin"
}

func (p *Baidupinyin) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "baidu-pinyin",
				Name:            "百度输入法",
				Description:     "中文输入法。",
				Organization:    "百度",
				OfficialWebsite: "https://shurufa.baidu.com",
				Icon:            "",
				Tags:            []string{"输入法"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Baidupinyin) Disabled() bool { return true }
