package onedrive

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	onedriveOfficialWebsite = "https://www.microsoft.com/onedrive"
	onedriveIconURL         = "https://www.microsoft.com/favicon.ico"
)

type Onedrive struct{}

func init() {
	plugin.Register(&Onedrive{})
}

func (p *Onedrive) Name() string {
	return "onedrive"
}

func (p *Onedrive) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "onedrive",
				Name:            "OneDrive",
				Description:     "Microsoft 出品的云存储服务，支持文件同步、备份与共享。",
				Organization:    "Microsoft",
				OfficialWebsite: onedriveOfficialWebsite,
				Icon:            onedriveIconURL,
				Tags:            []string{"云存储"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: onedriveOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, onedriveOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "OneDrive 下载页", URL: onedriveOfficialWebsite},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "OneDrive 下载页", URL: onedriveOfficialWebsite},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "Google Play", URL: "https://play.google.com/store/apps/details?id=com.microsoft.skydrive"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/onedrive/id477537958"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Onedrive) Disabled() bool { return false }
