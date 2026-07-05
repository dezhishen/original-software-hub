package vlc

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	vlcOfficialWebsite = "https://www.videolan.org/vlc"
	vlcIconURL         = "https://www.videolan.org/favicon.ico"
)

type Vlc struct{}

func init() {
	plugin.Register(&Vlc{})
}

func (p *Vlc) Name() string {
	return "vlc"
}

func (p *Vlc) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "vlc",
				Name:            "VLC 媒体播放器",
				Description:     "VideoLAN 出品的开源跨平台多媒体播放器，支持几乎所有音频和视频格式。",
				Organization:    "VideoLAN",
				OfficialWebsite: vlcOfficialWebsite,
				Icon:            vlcIconURL,
				Tags:            []string{"媒体播放", "音视频"},
				Categories:      []string{"media"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: vlcOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, vlcOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "VLC Windows 安装包 (x64)", URL: "https://get.videolan.org/vlc/latest/win64/vlc-latest-win64.exe"},
							},
						},
						{
							Architecture: "x86",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "VLC Windows 安装包 (x86)", URL: "https://get.videolan.org/vlc/latest/win32/vlc-latest-win32.exe"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: "VLC macOS 安装包", URL: "https://get.videolan.org/vlc/latest/macosx/vlc-latest-universal.dmg"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "webpage", Label: "VLC for Linux", URL: "https://www.videolan.org/vlc/#download"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "Google Play", URL: "https://play.google.com/store/apps/details?id=org.videolan.vlc"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/vlc-for-mobile/id650377962"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Vlc) Disabled() bool { return false }
