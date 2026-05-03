package potplayer

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Potplayer struct{}

func init() {
	plugin.Register(&Potplayer{})
}

func (p *Potplayer) Name() string {
	return "potplayer"
}

func (p *Potplayer) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "potplayer",
				Name:            "PotPlayer",
				Description:     "媒体播放器。",
				Organization:    "",
				OfficialWebsite: "https://potplayer.daum.net",
				Icon:            "",
				Tags:            []string{"媒体播放", "音视频"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Potplayer) Disabled() bool { return true }
