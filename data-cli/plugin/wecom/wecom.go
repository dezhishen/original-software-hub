package wecom

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	wecomOfficialWebsite = "https://work.weixin.qq.com/"
	wecomDownloadPage    = "https://work.weixin.qq.com/#indexDownload"
	wecomIconURL         = "https://work.weixin.qq.com/favicon.ico"
)

// WeCom implements plugin.Plugin for Tencent enterprise messenger.
type WeCom struct{}

func init() {
	plugin.Register(&WeCom{})
}

func (w *WeCom) Name() string {
	return "wecom"
}

func (w *WeCom) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")
	version := "Latest"

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "wecom",
				Name:            "企业微信",
				Icon:            wecomIconURL,
				Description:     "腾讯企业微信，面向企业的办公沟通与协作平台。",
				Organization:    "Tencent",
				OfficialWebsite: wecomOfficialWebsite,
				Tags:            []string{"办公协作", "企业通讯"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: wecomDownloadPage,
					Variants: []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links:        []plugin.Link{{Type: "direct", Label: "企业微信 Windows 下载", URL: wecomDownloadPage}},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links:        []plugin.Link{{Type: "direct", Label: "企业微信 macOS 下载", URL: wecomDownloadPage}},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links:        []plugin.Link{{Type: "direct", Label: "企业微信 Linux 下载", URL: wecomDownloadPage}},
						},
					},
				},
			},
		},
	}, nil
}
