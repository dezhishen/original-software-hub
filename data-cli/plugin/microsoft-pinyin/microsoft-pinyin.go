package microsoftpinyin

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Microsoftpinyin struct{}

func init() {
	plugin.Register(&Microsoftpinyin{})
}

func (p *Microsoftpinyin) Name() string {
	return "microsoft-pinyin"
}

func (p *Microsoftpinyin) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "microsoft-pinyin",
				Name:            "微软拼音",
				Description:     "中文输入法。",
				Organization:    "Microsoft",
				OfficialWebsite: "https://www.microsoft.com",
				Icon:            "",
				Tags:            []string{"输入法"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Microsoftpinyin) Disabled() bool { return true }
