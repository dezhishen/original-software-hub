package sogoubrouser

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	sogoubrowserOfficialWebsite = "https://www.sogou.com"
	sogoubrowserIconURL         = "https://www.sogou.com/favicon.ico"
)

type Sogoubrouser struct{}

func init() {
	plugin.Register(&Sogoubrouser{})
}

func (p *Sogoubrouser) Name() string {
	return "sogou-browser"
}

func (p *Sogoubrouser) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "sogou-browser",
				Name:            "搜狗高速浏览器",
				Description:     "搜狗出品的双核浏览器，结合 Trident 和 Chromium 内核，提供高速上网体验。",
				Organization:    "搜狗",
				OfficialWebsite: sogoubrowserOfficialWebsite,
				Icon:            sogoubrowserIconURL,
				Tags:            []string{"浏览器"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: "https://browser.sogou.com",
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, "https://browser.sogou.com", []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "搜狗浏览器下载", URL: "https://browser.sogou.com/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Sogoubrouser) Disabled() bool { return false }
