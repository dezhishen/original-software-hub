package foxitreader

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	foxitreaderOfficialWebsite = "https://www.foxitsoftware.com/pdf-reader"
	foxitreaderIconURL         = "https://www.foxitsoftware.com/favicon.ico"
)

type Foxitreader struct{}

func init() {
	plugin.Register(&Foxitreader{})
}

func (p *Foxitreader) Name() string {
	return "foxit-reader"
}

func (p *Foxitreader) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "foxit-reader",
				Name:            "福昕阅读器（Foxit Reader）",
				Description:     "福昕出品的免费 PDF 阅读与编辑工具，支持注释、表单填写和文档协作。",
				Organization:    "福昕",
				OfficialWebsite: foxitreaderOfficialWebsite,
				Icon:            foxitreaderIconURL,
				Tags:            []string{"PDF", "文档阅读"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: foxitreaderOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, foxitreaderOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Foxit Reader 下载页", URL: "https://www.foxitsoftware.com/downloads/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Foxitreader) Disabled() bool { return false }
