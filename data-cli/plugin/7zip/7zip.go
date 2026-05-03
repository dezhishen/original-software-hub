package sevenzip

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Sevenzip struct{}

func init() {
	plugin.Register(&Sevenzip{})
}

func (p *Sevenzip) Name() string {
	return "7zip"
}

func (p *Sevenzip) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "7zip",
				Name:            "7-Zip",
				Description:     "压缩与解压工具。",
				Organization:    "",
				OfficialWebsite: "https://www.7-zip.org",
				Icon:            "",
				Tags:            []string{"压缩"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Sevenzip) Disabled() bool { return true }
