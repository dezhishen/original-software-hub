package baidunetdisk

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	baiduNetdiskOfficialWebsite = "https://pan.baidu.com/"
	baiduNetdiskDownloadPage    = "https://pan.baidu.com/download"
	baiduNetdiskIconURL         = "https://pan.baidu.com/favicon.ico"
)

// BaiduNetdisk implements plugin.Plugin for Baidu Netdisk client.
type BaiduNetdisk struct{}

func init() {
	plugin.Register(&BaiduNetdisk{})
}

func (b *BaiduNetdisk) Name() string {
	return "baidunetdisk"
}

func (b *BaiduNetdisk) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")
	version := "Latest"

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "baidunetdisk",
				Name:            "百度网盘",
				Icon:            baiduNetdiskIconURL,
				Description:     "百度网盘客户端，提供文件同步、分享与云存储服务。",
				Organization:    "Baidu",
				OfficialWebsite: baiduNetdiskOfficialWebsite,
				Tags:            []string{"云存储", "文件同步"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: baiduNetdiskDownloadPage,
					Variants: []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links:        []plugin.Link{{Type: "direct", Label: "百度网盘 Windows 下载", URL: baiduNetdiskDownloadPage}},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links:        []plugin.Link{{Type: "direct", Label: "百度网盘 macOS 下载", URL: baiduNetdiskDownloadPage}},
						},
					},
				},
			},
		},
	}, nil
}
