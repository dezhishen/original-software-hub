package potplayer

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	potplayerOfficialWebsite = "https://potplayer.daum.net"
)

type Potplayer struct{}

func init() {
	plugin.Register(&Potplayer{})
}

func (p *Potplayer) Name() string {
	return "potplayer"
}

func (p *Potplayer) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "potplayer",
				Name:            "PotPlayer",
				Description:     "多功能多媒体播放器，支持广泛的视频格式，界面简洁高效。",
				Organization:    "Daum Communications",
				OfficialWebsite: potplayerOfficialWebsite,
				Icon:            "",
				Tags:            []string{"媒体播放", "音视频"},
				Categories:      []string{"media"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: potplayerOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, potplayerOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "PotPlayer 下载页", URL: "https://potplayer.daum.net/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Potplayer) Disabled() bool { return false }
