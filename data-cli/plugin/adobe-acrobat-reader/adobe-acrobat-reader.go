package adobeacrobatreader

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	adobeacrobatreaderOfficialWebsite = "https://www.adobe.com/acrobat/pdf-reader"
	adobeacrobatreaderIconURL         = "https://www.adobe.com/favicon.ico"
	adobeacrobatreaderDownloadURL     = "https://get.adobe.com/reader/"
)

type Adobeacrobatreader struct{}

func init() {
	plugin.Register(&Adobeacrobatreader{})
}

func (p *Adobeacrobatreader) Name() string {
	return "adobe-acrobat-reader"
}

func (p *Adobeacrobatreader) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "adobe-acrobat-reader",
				Name:            "Adobe Acrobat Reader",
				Description:     "Adobe 出品的免费 PDF 文档阅读、编辑与签名工具。",
				Organization:    "Adobe",
				OfficialWebsite: adobeacrobatreaderOfficialWebsite,
				Icon:            adobeacrobatreaderIconURL,
				Tags:            []string{"PDF", "文档阅读"},
				Categories:      []string{"office"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: adobeacrobatreaderOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, adobeacrobatreaderOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Adobe Acrobat Reader 下载页", URL: adobeacrobatreaderDownloadURL},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Adobe Acrobat Reader 下载页", URL: adobeacrobatreaderDownloadURL},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "Google Play", URL: "https://play.google.com/store/apps/details?id=com.adobe.reader"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/adobe-acrobat-reader/id469337564"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Adobeacrobatreader) Disabled() bool { return false }
