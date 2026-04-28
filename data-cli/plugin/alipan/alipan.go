package alipan

import (
	"fmt"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
	"github.com/dezhishen/original-software-hub/data-cli/util"
)

const (
	alipanID           = "alipan"
	alipanName         = "阿里云盘"
	alipanDescription  = "阿里巴巴推出的高速、可靠的个人网盘，支持文件备份、同步和分享"
	alipanOfficialURL  = "https://www.alipan.com"
	alipanIconURL      = "https://img.alicdn.com/imgextra/i2/O1CN01DOYcs71v3B6bOemVM_!!6000000006116-2-tps-512-512.png"
	alipanOrganization = "Alibaba"
)

// Alipan implements plugin.Plugin for Alibaba Aliyun Drive desktop client.
type Alipan struct{}

func init() {
	plugin.Register(&Alipan{})
}

func (a *Alipan) Name() string {
	return alipanID
}

func (a *Alipan) Fetch() ([]plugin.SoftwareData, error) {
	info, err := util.FetchAlipanDownloadInfo()
	if err != nil {
		return nil, fmt.Errorf("fetch alipan info: %w", err)
	}

	variants := []plugin.Variant{}

	// Windows x64
	if info.WindowsURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows",
			Links: []plugin.Link{
				{
					Type:  "direct",
					Label: fmt.Sprintf("阿里云盘 Windows (%s)", "x64"),
					URL:   info.WindowsURL,
				},
			},
		})
	}

	// macOS Intel
	if info.MacIntelURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "Intel",
			Platform:     "macOS",
			Links: []plugin.Link{
				{
					Type:  "direct",
					Label: fmt.Sprintf("阿里云盘 macOS (Intel)"),
					URL:   info.MacIntelURL,
				},
			},
		})
	}

	// macOS Apple Silicon (ARM64)
	if info.MacARM64URL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "Apple Silicon",
			Platform:     "macOS",
			Links: []plugin.Link{
				{
					Type:  "direct",
					Label: fmt.Sprintf("阿里云盘 macOS (Apple Silicon)"),
					URL:   info.MacARM64URL,
				},
			},
		})
	}

	// Android
	if info.AndroidURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "arm64",
			Platform:     "Android",
			Links: []plugin.Link{
				{
					Type:  "direct",
					Label: "阿里云盘 Android",
					URL:   info.AndroidURL,
				},
			},
		})
	}

	// iOS
	if info.IOSUrl != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "arm64",
			Platform:     "iOS",
			Links: []plugin.Link{
				{
					Type:  "store",
					Label: "阿里云盘 App Store",
					URL:   info.IOSUrl,
				},
			},
		})
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              alipanID,
				Name:            alipanName,
				Icon:            alipanIconURL,
				Description:     alipanDescription,
				Organization:    alipanOrganization,
				OfficialWebsite: alipanOfficialURL,
			},
			Versions: []plugin.Version{
				{
					Version:     info.Version,
					ReleaseDate: info.ReleaseDate,
					OfficialURL: info.DownloadPageURL,
					Variants:    variants,
				},
			},
		},
	}, nil
}

func (x *Alipan) CompareWithPrevious(previous plugin.PreviousState) ([]plugin.FetchResult, error) {
	items, err := x.Fetch()
	if err != nil {
		return nil, err
	}
	return plugin.BuildCompareResults(items, previous), nil
}
