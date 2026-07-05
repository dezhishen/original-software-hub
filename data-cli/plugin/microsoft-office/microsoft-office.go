package microsoftoffice

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	microsoftofficeOfficialWebsite = "https://www.microsoft.com/office"
	microsoftofficeIconURL         = "https://www.microsoft.com/favicon.ico"
)

type Microsoftoffice struct{}

func init() {
	plugin.Register(&Microsoftoffice{})
}

func (p *Microsoftoffice) Name() string {
	return "microsoft-office"
}

func (p *Microsoftoffice) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "microsoft-office",
				Name:            "Microsoft Office",
				Description:     "Microsoft 出品的办公套件，包含 Word、Excel、PowerPoint、Outlook 等经典应用。",
				Organization:    "Microsoft",
				OfficialWebsite: microsoftofficeOfficialWebsite,
				Icon:            microsoftofficeIconURL,
				Tags:            []string{"办公软件", "文档处理"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: microsoftofficeOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, microsoftofficeOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Microsoft Office 官方页", URL: microsoftofficeOfficialWebsite},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Microsoft Office for Mac 官方页", URL: microsoftofficeOfficialWebsite},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "Google Play", URL: "https://play.google.com/store/apps/details?id=com.microsoft.office.officehubrow"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/microsoft-office/id541164041"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Microsoftoffice) Disabled() bool { return false }
