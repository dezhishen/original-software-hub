package qqbrowser

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Qqbrowser struct{}

func init() {
	plugin.Register(&Qqbrowser{})
}

func (p *Qqbrowser) Name() string {
	return "qq-browser"
}

func (p *Qqbrowser) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "qq-browser",
				Name:            "QQ 浏览器",
				Description:     "网页浏览器。",
				Organization:    "Tencent",
				OfficialWebsite: "https://browser.qq.com",
				Icon:            "",
				Tags:            []string{"浏览器"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Qqbrowser) Disabled() bool { return true }
