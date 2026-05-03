package winrar

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Winrar struct{}

func init() {
	plugin.Register(&Winrar{})
}

func (p *Winrar) Name() string {
	return "winrar"
}

func (p *Winrar) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "winrar",
				Name:            "WinRAR",
				Description:     "压缩与解压工具。",
				Organization:    "",
				OfficialWebsite: "https://www.rarlab.com",
				Icon:            "",
				Tags:            []string{"压缩"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Winrar) Disabled() bool { return true }
