package microsoftpinyin

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	microsoftpinyinOfficialWebsite = "https://www.microsoft.com"
	microsoftpinyinIconURL         = "https://www.microsoft.com/favicon.ico"
)

type Microsoftpinyin struct{}

func init() {
	plugin.Register(&Microsoftpinyin{})
}

func (p *Microsoftpinyin) Name() string {
	return "microsoft-pinyin"
}

func (p *Microsoftpinyin) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "microsoft-pinyin",
				Name:            "微软拼音",
				Description:     "Windows 内置的中文拼音输入法，由 Microsoft 开发维护，支持智能词组与云输入。",
				Organization:    "Microsoft",
				OfficialWebsite: microsoftpinyinOfficialWebsite,
				Icon:            microsoftpinyinIconURL,
				Tags:            []string{"输入法"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: "https://support.microsoft.com/zh-cn/windows/%E5%9C%A8-windows-%E4%B8%AD%E7%AE%A1%E7%90%86%E8%BE%93%E5%85%A5%E6%B3%95-1b7af6fc-e66f-4bbe-bec0-52e5e33bbe93",
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, microsoftpinyinOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Windows 输入法管理帮助", URL: "https://support.microsoft.com/zh-cn/windows/%E5%9C%A8-windows-%E4%B8%AD%E7%AE%A1%E7%90%86%E8%BE%93%E5%85%A5%E6%B3%95-1b7af6fc-e66f-4bbe-bec0-52e5e33bbe93"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Microsoftpinyin) Disabled() bool { return false }
