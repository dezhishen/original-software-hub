package iflytekpinyin

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	iflytekpinyinOfficialWebsite = "https://shurufa.iflytek.com"
)

type Iflytekpinyin struct{}

func init() {
	plugin.Register(&Iflytekpinyin{})
}

func (p *Iflytekpinyin) Name() string {
	return "iflytek-pinyin"
}

func (p *Iflytekpinyin) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "iflytek-pinyin",
				Name:            "讯飞输入法",
				Description:     "科大讯飞出品的免费中文输入法，支持语音、拼音、手写等多种输入方式。",
				Organization:    "科大讯飞",
				OfficialWebsite: iflytekpinyinOfficialWebsite,
				Icon:            "",
				Tags:            []string{"输入法"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: iflytekpinyinOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, iflytekpinyinOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "讯飞输入法 Windows 下载页", URL: "https://shurufa.iflytek.com/pc"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "webpage", Label: "讯飞输入法 macOS 下载页", URL: "https://shurufa.iflytek.com/pc"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "应用市场", URL: "https://shurufa.iflytek.com/android"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "iOS / iPadOS",
							Links: []plugin.Link{
								{Type: "store", Label: "App Store", URL: "https://apps.apple.com/app/id917710518"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Iflytekpinyin) Disabled() bool { return false }
