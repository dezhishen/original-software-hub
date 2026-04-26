package dingtalk

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	dingtalkOfficialWebsite = "https://www.dingtalk.com/"
	dingtalkDownloadPage    = "https://page.dingtalk.com/wow/z/dingtalk/default/download"
	dingtalkIconURL         = "https://www.dingtalk.com/favicon.ico"
)

// DingTalk implements plugin.Plugin for Alibaba DingTalk client.
type DingTalk struct{}

func init() {
	plugin.Register(&DingTalk{})
}

func (d *DingTalk) Name() string {
	return "dingtalk"
}

func (d *DingTalk) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")
	version := "Latest"

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "dingtalk",
				Name:            "钉钉",
				Icon:            dingtalkIconURL,
				Description:     "阿里巴巴旗下企业协作与即时通讯应用。",
				Organization:    "Alibaba",
				OfficialWebsite: dingtalkOfficialWebsite,
				Tags:            []string{"办公协作", "即时通讯"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: dingtalkDownloadPage,
					Variants: []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links:        []plugin.Link{{Type: "direct", Label: "钉钉 Windows 下载", URL: dingtalkDownloadPage}},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links:        []plugin.Link{{Type: "direct", Label: "钉钉 macOS 下载", URL: dingtalkDownloadPage}},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links:        []plugin.Link{{Type: "direct", Label: "钉钉 Linux 下载", URL: dingtalkDownloadPage}},
						},
					},
				},
			},
		},
	}, nil
}
