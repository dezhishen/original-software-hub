package qqbrowser

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	qqbrowserOfficialWebsite = "https://browser.qq.com"
	qqbrowserIconURL         = "https://browser.qq.com/favicon.ico"
)

type Qqbrowser struct{}

func init() {
	plugin.Register(&Qqbrowser{})
}

func (p *Qqbrowser) Name() string {
	return "qq-browser"
}

func (p *Qqbrowser) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "qq-browser",
				Name:            "QQ 浏览器",
				Description:     "腾讯出品的 Chromium 内核浏览器，提供高速、安全的网页浏览体验。",
				Organization:    "Tencent",
				OfficialWebsite: qqbrowserOfficialWebsite,
				Icon:            qqbrowserIconURL,
				Tags:            []string{"浏览器"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: qqbrowserOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, qqbrowserOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "QQ 浏览器下载", URL: "https://browser.qq.com/"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "QQ 浏览器 Mac 版下载", URL: "https://browser.qq.com/mac/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Qqbrowser) Disabled() bool { return false }
