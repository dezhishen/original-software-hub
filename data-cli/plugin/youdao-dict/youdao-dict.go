package youdaodict

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	youdaodictOfficialWebsite = "https://youdao.com"
	youdaodictIconURL         = "https://youdao.com/favicon.ico"
)

type Youdaodict struct{}

func init() {
	plugin.Register(&Youdaodict{})
}

func (p *Youdaodict) Name() string {
	return "youdao-dict"
}

func (p *Youdaodict) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "youdao-dict",
				Name:            "有道词典（桌面版）",
				Description:     "网易有道出品的免费翻译与词典工具，支持多语种翻译、查词和语音朗读。",
				Organization:    "网易有道",
				OfficialWebsite: youdaodictOfficialWebsite,
				Icon:            youdaodictIconURL,
				Tags:            []string{"翻译", "词典"},
				Categories:      []string{"translation"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: youdaodictOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, youdaodictOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "有道词典 PC 版下载", URL: "https://youdao.com/dict/download"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "有道词典 Mac 版下载", URL: "https://youdao.com/dict/download"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Youdaodict) Disabled() bool { return false }
