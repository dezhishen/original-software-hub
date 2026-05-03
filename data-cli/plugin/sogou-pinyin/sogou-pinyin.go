package sogoupinyin

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Sogoupinyin struct{}

func init() {
	plugin.Register(&Sogoupinyin{})
}

func (p *Sogoupinyin) Name() string {
	return "sogou-pinyin"
}

func (p *Sogoupinyin) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "sogou-pinyin",
				Name:            "搜狗输入法",
				Description:     "中文输入法。",
				Organization:    "搜狗",
				OfficialWebsite: "https://pinyin.sogou.com",
				Icon:            "",
				Tags:            []string{"输入法"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Sogoupinyin) Disabled() bool { return true }
