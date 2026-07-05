package kingsoftantivirus

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	kingsoftantivirusOfficialWebsite = "https://www.kingsoft.com"
	kingsoftantivirusIconURL         = "https://www.kingsoft.com/favicon.ico"
)

type Kingsoftantivirus struct{}

func init() {
	plugin.Register(&Kingsoftantivirus{})
}

func (p *Kingsoftantivirus) Name() string {
	return "kingsoft-antivirus"
}

func (p *Kingsoftantivirus) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "kingsoft-antivirus",
				Name:            "金山毒霸",
				Description:     "金山出品的综合安全防护软件，提供病毒查杀、系统清理和网络安全防护。",
				Organization:    "Kingsoft",
				OfficialWebsite: kingsoftantivirusOfficialWebsite,
				Icon:            kingsoftantivirusIconURL,
				Tags:            []string{"安全防护", "杀毒软件"},
				Categories:      []string{"security"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: kingsoftantivirusOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, kingsoftantivirusOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "金山毒霸下载页", URL: "https://www.kingsoft.com/duba/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Kingsoftantivirus) Disabled() bool { return false }
