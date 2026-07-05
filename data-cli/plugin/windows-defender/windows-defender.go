package windowsdefender

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	windowsdefenderOfficialWebsite = "https://www.microsoft.com"
	windowsdefenderIconURL         = "https://www.microsoft.com/favicon.ico"
)

type Windowsdefender struct{}

func init() {
	plugin.Register(&Windowsdefender{})
}

func (p *Windowsdefender) Name() string {
	return "windows-defender"
}

func (p *Windowsdefender) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "windows-defender",
				Name:            "Microsoft Defender",
				Description:     "Windows 内置的安全防护功能，现升级为 Microsoft Defender，提供实时病毒和威胁防护。",
				Organization:    "Microsoft",
				OfficialWebsite: windowsdefenderOfficialWebsite,
				Icon:            windowsdefenderIconURL,
				Tags:            []string{"安全防护"},
				Categories:      []string{"security"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: "https://www.microsoft.com/windows/comprehensive-security",
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, windowsdefenderOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Microsoft Defender 安全中心", URL: "https://www.microsoft.com/windows/comprehensive-security"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Windowsdefender) Disabled() bool { return false }
