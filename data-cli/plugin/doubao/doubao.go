package doubao

import (
	"fmt"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
	"github.com/dezhishen/original-software-hub/data-cli/util"
)

const (
	doubaoID           = "doubao"
	doubaoName         = "豆包客户端"
	doubaoDescription  = "字节跳动推出的 AI 助手，提供智能对话、文案创作、代码编程等功能的桌面客户端"
	doubaoOfficialURL  = "https://www.doubao.com"
	doubaoIconURL      = "https://lf-flow-web-cdn.doubao.com/obj/flow-doubao/doubao/desktop_online_web/static/image/logo-doubao.7d723a57.png"
	doubaoOrganization = "ByteDance"
)

// Doubao implements plugin.Plugin for ByteDance Doubao desktop client.
type Doubao struct{}

func init() {
	plugin.Register(&Doubao{})
}

func (d *Doubao) Name() string {
	return doubaoID
}

func (d *Doubao) Fetch() ([]plugin.SoftwareData, error) {
	info, err := util.FetchDoubaoDownloadInfo()
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

	// If no download variants found, provide download page link
	if len(variants) == 0 {
		variants = append(variants, plugin.Variant{
			Architecture: "Universal",
			Platform:     "Desktop",
			Links: []plugin.Link{
				{
					Type:  "direct",
					Label: "豆包官方下载页",
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

func (x *Doubao) FetchWithPrevious(previous plugin.PreviousState) ([]plugin.FetchResult, error) {
	items, err := x.Fetch()
	if err != nil {
		return nil, err
	}
	return plugin.BuildFetchResults(items, previous), nil
}
