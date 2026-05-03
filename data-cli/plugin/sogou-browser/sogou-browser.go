package sogoubrouser

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Sogoubrouser struct{}

func init() {
	plugin.Register(&Sogoubrouser{})
}

func (p *Sogoubrouser) Name() string {
	return "sogou-browser"
}

func (p *Sogoubrouser) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "sogou-browser",
				Name:            "搜狗高速浏览器",
				Description:     "网页浏览器。",
				Organization:    "搜狗",
				OfficialWebsite: "https://www.sogou.com",
				Icon:            "",
				Tags:            []string{"浏览器"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Sogoubrouser) Disabled() bool { return true }
