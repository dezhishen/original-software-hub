package microsoftoffice

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Microsoftoffice struct{}

func init() {
	plugin.Register(&Microsoftoffice{})
}

func (p *Microsoftoffice) Name() string {
	return "microsoft-office"
}

func (p *Microsoftoffice) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "microsoft-office",
				Name:            "Microsoft Office",
				Description:     "办公套件（Word、Excel、PowerPoint 等）。",
				Organization:    "Microsoft",
				OfficialWebsite: "https://www.microsoft.com/office",
				Icon:            "",
				Tags:            []string{"办公软件", "文档处理"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Microsoftoffice) Disabled() bool { return true }
