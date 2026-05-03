package tencentpcmanager

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Tencentpcmanager struct{}

func init() {
	plugin.Register(&Tencentpcmanager{})
}

func (p *Tencentpcmanager) Name() string {
	return "tencent-pc-manager"
}

func (p *Tencentpcmanager) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "tencent-pc-manager",
				Name:            "腾讯电脑管家",
				Description:     "安全与防护工具。",
				Organization:    "Tencent",
				OfficialWebsite: "https://guanjia.qq.com",
				Icon:            "",
				Tags:            []string{"安全防护"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Tencentpcmanager) Disabled() bool { return true }
