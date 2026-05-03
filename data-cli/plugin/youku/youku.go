package youku

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Youku struct{}

func init() {
	plugin.Register(&Youku{})
}

func (p *Youku) Name() string {
	return "youku"
}

func (p *Youku) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "youku",
				Name:            "优酷",
				Description:     "PC 客户端。",
				Organization:    "优酷",
				OfficialWebsite: "https://www.youku.com",
				Icon:            "",
				Tags:            []string{"视频"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Youku) Disabled() bool { return true }
