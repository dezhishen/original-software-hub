package p360zip

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type P360zip struct{}

func init() {
	plugin.Register(&P360zip{})
}

func (p *P360zip) Name() string {
	return "360-zip"
}

func (p *P360zip) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "360-zip",
				Name:            "360 压缩",
				Description:     "压缩与解压工具。",
				Organization:    "",
				OfficialWebsite: "https://yasuo.360.cn",
				Icon:            "",
				Tags:            []string{"压缩"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *P360zip) Disabled() bool { return true }
