package edge

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Edge struct{}

func init() {
	plugin.Register(&Edge{})
}

func (p *Edge) Name() string {
	return "edge"
}

func (p *Edge) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "edge",
				Name:            "Microsoft Edge",
				Description:     "微软浏览器。",
				Organization:    "Microsoft",
				OfficialWebsite: "https://www.microsoft.com/edge",
				Icon:            "https://www.microsoft.com/favicon.ico",
				Tags:            []string{"浏览器"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Edge) Disabled() bool { return true }
