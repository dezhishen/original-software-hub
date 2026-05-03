package vlc

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Vlc struct{}

func init() {
	plugin.Register(&Vlc{})
}

func (p *Vlc) Name() string {
	return "vlc"
}

func (p *Vlc) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "vlc",
				Name:            "VLC 媒体播放器",
				Description:     "多媒体播放器。",
				Organization:    "",
				OfficialWebsite: "https://www.videolan.org/vlc",
				Icon:            "https://www.videolan.org/favicon.ico",
				Tags:            []string{"媒体播放", "音视频"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Vlc) Disabled() bool { return true }
