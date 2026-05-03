package kugoumusic

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Kugoumusic struct{}

func init() {
	plugin.Register(&Kugoumusic{})
}

func (p *Kugoumusic) Name() string {
	return "kugou-music"
}

func (p *Kugoumusic) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "kugou-music",
				Name:            "酷狗音乐",
				Description:     "音乐客户端。",
				Organization:    "酷狗",
				OfficialWebsite: "https://www.kugou.com",
				Icon:            "",
				Tags:            []string{"音乐"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Kugoumusic) Disabled() bool { return true }
