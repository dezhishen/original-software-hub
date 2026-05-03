package outlook

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Outlook struct{}

func init() {
	plugin.Register(&Outlook{})
}

func (p *Outlook) Name() string {
	return "outlook"
}

func (p *Outlook) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "outlook",
				Name:            "Microsoft Outlook",
				Description:     "邮件与日历客户端。",
				Organization:    "Microsoft",
				OfficialWebsite: "https://www.microsoft.com/outlook",
				Icon:            "https://www.microsoft.com/favicon.ico",
				Tags:            []string{"邮件", "日历"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Outlook) Disabled() bool { return true }
