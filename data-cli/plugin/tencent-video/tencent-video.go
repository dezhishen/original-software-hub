package tencentvideo

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Tencentvideo struct{}

func init() {
	plugin.Register(&Tencentvideo{})
}

func (p *Tencentvideo) Name() string {
	return "tencent-video"
}

func (p *Tencentvideo) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "tencent-video",
				Name:            "腾讯视频",
				Description:     "PC 客户端。",
				Organization:    "Tencent",
				OfficialWebsite: "https://v.qq.com",
				Icon:            "",
				Tags:            []string{"视频"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Tencentvideo) Disabled() bool { return true }
