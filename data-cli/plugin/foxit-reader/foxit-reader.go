package foxitreader

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Foxitreader struct{}

func init() {
	plugin.Register(&Foxitreader{})
}

func (p *Foxitreader) Name() string {
	return "foxit-reader"
}

func (p *Foxitreader) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "foxit-reader",
				Name:            "福昕阅读器（Foxit Reader）",
				Description:     "文档阅读与 PDF 工具。",
				Organization:    "福昕",
				OfficialWebsite: "https://www.foxitsoftware.com/pdf-reader",
				Icon:            "",
				Tags:            []string{"PDF", "文档阅读"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Foxitreader) Disabled() bool { return true }
