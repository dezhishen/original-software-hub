package onenote

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Onenote struct{}

func init() {
	plugin.Register(&Onenote{})
}

func (p *Onenote) Name() string {
	return "onenote"
}

func (p *Onenote) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "onenote",
				Name:            "OneNote",
				Description:     "笔记工具。",
				Organization:    "Microsoft",
				OfficialWebsite: "https://www.microsoft.com/onenote",
				Icon:            "https://www.microsoft.com/favicon.ico",
				Tags:            []string{"笔记", "效率工具"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Onenote) Disabled() bool { return true }
