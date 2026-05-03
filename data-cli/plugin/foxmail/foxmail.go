package foxmail

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Foxmail struct{}

func init() {
	plugin.Register(&Foxmail{})
}

func (p *Foxmail) Name() string {
	return "foxmail"
}

func (p *Foxmail) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "foxmail",
				Name:            "Foxmail",
				Description:     "邮件客户端。",
				Organization:    "",
				OfficialWebsite: "https://www.foxmail.com",
				Icon:            "",
				Tags:            []string{"邮件"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Foxmail) Disabled() bool { return true }
