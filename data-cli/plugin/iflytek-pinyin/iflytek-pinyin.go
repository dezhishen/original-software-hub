package iflytekpinyin

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Iflytekpinyin struct{}

func init() {
	plugin.Register(&Iflytekpinyin{})
}

func (p *Iflytekpinyin) Name() string {
	return "iflytek-pinyin"
}

func (p *Iflytekpinyin) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "iflytek-pinyin",
				Name:            "讯飞输入法",
				Description:     "中文输入法。",
				Organization:    "讯飞",
				OfficialWebsite: "https://shurufa.iflytek.com",
				Icon:            "",
				Tags:            []string{"输入法"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Iflytekpinyin) Disabled() bool { return true }
