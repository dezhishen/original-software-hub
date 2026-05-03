package adobeacrobatreader

import "github.com/dezhishen/original-software-hub/data-cli/plugin"

type Adobeacrobatreader struct{}

func init() {
	plugin.Register(&Adobeacrobatreader{})
}

func (p *Adobeacrobatreader) Name() string {
	return "adobe-acrobat-reader"
}

func (p *Adobeacrobatreader) Fetch() ([]plugin.SoftwareData, error) {
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "adobe-acrobat-reader",
				Name:            "Adobe Acrobat Reader",
				Description:     "文档阅读与 PDF 工具。",
				Organization:    "Adobe",
				OfficialWebsite: "https://www.adobe.com/acrobat/pdf-reader",
				Icon:            "",
				Tags:            []string{"PDF", "文档阅读"},
			},
			Versions: []plugin.Version{},
		},
	}, nil
}

func (p *Adobeacrobatreader) Disabled() bool { return true }
