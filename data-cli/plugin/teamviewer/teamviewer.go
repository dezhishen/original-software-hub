package teamviewer

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Teamviewer struct{}

func init() {
	plugin.Register(&Teamviewer{})
}

func (p *Teamviewer) Name() string {
	return "teamviewer"
}

func (p *Teamviewer) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "teamviewer",
				Name:            "TeamViewer",
				Description:     "远程控制工具。",
				Organization:    "TeamViewer",
				OfficialWebsite: "https://www.teamviewer.com",
				Icon:            "",
				Tags:            []string{"远程控制"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Teamviewer) Disabled() bool { return true }
