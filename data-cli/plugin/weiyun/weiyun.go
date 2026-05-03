package weiyun

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Weiyun struct{}

func init() {
	plugin.Register(&Weiyun{})
}

func (p *Weiyun) Name() string {
	return "weiyun"
}

func (p *Weiyun) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "weiyun",
				Name:            "腾讯微云",
				Description:     "云存储与同步客户端。",
				Organization:    "Tencent",
				OfficialWebsite: "https://www.weiyun.com",
				Icon:            "",
				Tags:            []string{"云存储"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Weiyun) Disabled() bool { return true }
