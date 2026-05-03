package p360browser

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type P360browser struct{}

func init() {
	plugin.Register(&P360browser{})
}

func (p *P360browser) Name() string {
	return "360-browser"
}

func (p *P360browser) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "360-browser",
				Name:            "360 浏览器",
				Description:     "网页浏览器。",
				Organization:    "",
				OfficialWebsite: "https://browser.360.cn",
				Icon:            "",
				Tags:            []string{"浏览器"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *P360browser) Disabled() bool { return true }
