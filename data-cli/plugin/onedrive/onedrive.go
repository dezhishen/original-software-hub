package onedrive

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Onedrive struct{}

func init() {
	plugin.Register(&Onedrive{})
}

func (p *Onedrive) Name() string {
	return "onedrive"
}

func (p *Onedrive) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "onedrive",
				Name:            "OneDrive",
				Description:     "微软云盘。",
				Organization:    "Microsoft",
				OfficialWebsite: "https://www.microsoft.com/onedrive",
				Icon:            "https://www.microsoft.com/favicon.ico",
				Tags:            []string{"云存储"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Onedrive) Disabled() bool { return true }
