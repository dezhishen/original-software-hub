package onenote

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	onenoteOfficialWebsite = "https://www.microsoft.com/onenote"
	onenoteIconURL         = "https://www.microsoft.com/favicon.ico"
)

type Onenote struct{}

func init() {
	plugin.Register(&Onenote{})
}

func (p *Onenote) Name() string {
	return "onenote"
}

func (p *Onenote) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "onenote",
				Name:            "OneNote",
				Description:     "Microsoft 出品的数字笔记应用，支持多设备同步与协作。",
				Organization:    "Microsoft",
				OfficialWebsite: onenoteOfficialWebsite,
				Icon:            onenoteIconURL,
				Tags:            []string{"笔记", "效率工具"},
				Categories:      []string{"note-taking"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: onenoteOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, onenoteOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "OneNote 官方页", URL: onenoteOfficialWebsite},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "store", Label: "Mac App Store", URL: "https://apps.apple.com/app/microsoft-onenote/id410395246"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "Google Play", URL: "https://play.google.com/store/apps/details?id=com.microsoft.office.onenote"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/microsoft-onenote/id410395246"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Onenote) Disabled() bool { return false }
