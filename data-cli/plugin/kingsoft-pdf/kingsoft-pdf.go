package kingsoftpdf

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Kingsoftpdf struct{}

func init() {
	plugin.Register(&Kingsoftpdf{})
}

func (p *Kingsoftpdf) Name() string {
	return "kingsoft-pdf"
}

func (p *Kingsoftpdf) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "kingsoft-pdf",
				Name:            "金山 PDF",
				Description:     "文档阅读与 PDF 工具。",
				Organization:    "Kingsoft",
				OfficialWebsite: "https://www.kingsoft.com",
				Icon:            "",
				Tags:            []string{"PDF", "文档阅读"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Kingsoftpdf) Disabled() bool { return true }
