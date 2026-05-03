package iqiyi

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Iqiyi struct{}

func init() {
	plugin.Register(&Iqiyi{})
}

func (p *Iqiyi) Name() string {
	return "iqiyi"
}

func (p *Iqiyi) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "iqiyi",
				Name:            "爱奇艺",
				Description:     "PC 客户端。",
				Organization:    "爱奇艺",
				OfficialWebsite: "https://www.iqiyi.com",
				Icon:            "",
				Tags:            []string{"视频"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Iqiyi) Disabled() bool { return true }
