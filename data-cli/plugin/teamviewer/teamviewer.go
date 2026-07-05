package teamviewer

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	teamviewerOfficialWebsite = "https://www.teamviewer.com"
	teamviewerIconURL         = "https://www.teamviewer.com/favicon.ico"
	teamviewerDownloadURL     = "https://download.teamviewer.com/download/TeamViewer_Setup_x64.exe"
)

type Teamviewer struct{}

func init() {
	plugin.Register(&Teamviewer{})
}

func (p *Teamviewer) Name() string {
	return "teamviewer"
}

func (p *Teamviewer) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "teamviewer",
				Name:            "TeamViewer",
				Description:     "全球知名的远程控制、远程访问和远程支持软件。",
				Organization:    "TeamViewer AG",
				OfficialWebsite: teamviewerOfficialWebsite,
				Icon:            teamviewerIconURL,
				Tags:            []string{"远程控制"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: teamviewerOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, teamviewerOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "TeamViewer 安装包 (x64)", URL: teamviewerDownloadURL},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "TeamViewer macOS 下载", URL: "https://www.teamviewer.com/zh-cn/download/macos/"},
							},
						},
						{
							Architecture: "x64 (deb)",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "webpage", Label: "TeamViewer Linux 下载", URL: "https://www.teamviewer.com/zh-cn/download/linux/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Teamviewer) Disabled() bool { return false }
