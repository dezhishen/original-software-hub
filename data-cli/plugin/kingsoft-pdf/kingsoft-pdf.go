package kingsoftpdf

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	kingsoftpdfOfficialWebsite = "https://www.kingsoft.com"
	kingsoftpdfIconURL         = "https://www.kingsoft.com/favicon.ico"
)

type Kingsoftpdf struct{}

func init() {
	plugin.Register(&Kingsoftpdf{})
}

func (p *Kingsoftpdf) Name() string {
	return "kingsoft-pdf"
}

func (p *Kingsoftpdf) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "kingsoft-pdf",
				Name:            "金山 PDF",
				Description:     "金山出品的免费 PDF 阅读器，支持文档阅读、注释和格式转换。",
				Organization:    "Kingsoft",
				OfficialWebsite: kingsoftpdfOfficialWebsite,
				Icon:            kingsoftpdfIconURL,
				Tags:            []string{"PDF", "文档阅读"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: kingsoftpdfOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, kingsoftpdfOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "金山 PDF 下载页", URL: "https://www.kingsoft.com/pdf/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Kingsoftpdf) Disabled() bool { return false }
