package powerword

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	powerwordOfficialWebsite = "https://www.kingsoft.com"
	powerwordIconURL         = "https://www.kingsoft.com/favicon.ico"
)

type Powerword struct{}

func init() {
	plugin.Register(&Powerword{})
}

func (p *Powerword) Name() string {
	return "powerword"
}

func (p *Powerword) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "powerword",
				Name:            "金山词霸",
				Description:     "金山出品的翻译与词典工具，支持英语查词、翻译和语音朗读等功能。",
				Organization:    "Kingsoft",
				OfficialWebsite: powerwordOfficialWebsite,
				Icon:            powerwordIconURL,
				Tags:            []string{"翻译", "词典"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: powerwordOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, powerwordOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "金山词霸 PC 版下载", URL: "https://cp.iciba.com/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Powerword) Disabled() bool { return false }
