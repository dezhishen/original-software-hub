package outlook

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	outlookOfficialWebsite = "https://www.microsoft.com/outlook"
	outlookIconURL         = "https://www.microsoft.com/favicon.ico"
)

type Outlook struct{}

func init() {
	plugin.Register(&Outlook{})
}

func (p *Outlook) Name() string {
	return "outlook"
}

func (p *Outlook) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "outlook",
				Name:            "Microsoft Outlook",
				Description:     "Microsoft 出品的邮件与日历管理客户端，支持日程管理和联系人管理。",
				Organization:    "Microsoft",
				OfficialWebsite: outlookOfficialWebsite,
				Icon:            outlookIconURL,
				Tags:            []string{"邮件", "日历"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: outlookOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, outlookOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Outlook 官方页", URL: outlookOfficialWebsite},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "store", Label: "Mac App Store", URL: "https://apps.apple.com/app/microsoft-outlook/id951937596"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "Google Play", URL: "https://play.google.com/store/apps/details?id=com.microsoft.office.outlook"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/microsoft-outlook/id951937596"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Outlook) Disabled() bool { return false }
