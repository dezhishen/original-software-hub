package doubao

import (
	"fmt"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	doubaoID           = "doubao"
	doubaoName         = "豆包客户端"
	doubaoDescription  = "字节跳动推出的 AI 助手，提供智能对话、文案创作、代码编程等功能的桌面客户端"
	doubaoOfficialURL  = "https://www.doubao.com"
	doubaoIconURL      = "https://lf-flow-web-cdn.doubao.com/obj/flow-doubao/doubao/desktop_online_web/static/image/logo-doubao.7d723a57.png"
	doubaoOrganization = "ByteDance"
)

var doubaoTags = []string{"AI 助手", "智能对话", "效率工具"}

// Doubao implements plugin.Plugin for ByteDance Doubao desktop client.
type Doubao struct{}

func init() {
	plugin.Register(&Doubao{})
}

func (d *Doubao) Name() string {
	return doubaoID
}

func (d *Doubao) Fetch() ([]plugin.SoftwareData, error) {
	info, err := fetchDownloadInfo()
	if err != nil {
		return nil, fmt.Errorf("fetch doubao info: %w", err)
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
					Label: fmt.Sprintf("豆包 Windows (%s)", "x64"),
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
					Label: "豆包 macOS (Intel)",
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
					Label: "豆包 macOS (Apple Silicon)",
					URL:   info.MacARM64URL,
				},
			},
		})
	}

	// If no direct links are found, keep platform granularity instead of collapsing to Desktop.
	if len(variants) == 0 {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows",
			Links: []plugin.Link{
				{
					Type:  "webpage",
					Label: "豆包 Windows 下载页",
					URL:   info.DownloadPageURL,
				},
			},
		})
		variants = append(variants, plugin.Variant{
			Architecture: "universal",
			Platform:     "macOS",
			Links: []plugin.Link{
				{
					Type:  "webpage",
					Label: "豆包 macOS 下载页",
					URL:   info.DownloadPageURL,
				},
			},
		})
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              doubaoID,
				Name:            doubaoName,
				Icon:            doubaoIconURL,
				Description:     doubaoDescription,
				Organization:    doubaoOrganization,
				OfficialWebsite: doubaoOfficialURL,
				Tags:            doubaoTags,
			},
			Versions: []plugin.Version{
				{
					Version:     info.Version,
					ReleaseDate: info.ReleaseDate,
					OfficialURL: info.DownloadPageURL,
					Platforms:   plugin.PlatformsFromVariants(info.Version, info.ReleaseDate, info.DownloadPageURL, variants),
				},
			},
		},
	}, nil
}

func (p *Doubao) Disabled() bool { return false }
