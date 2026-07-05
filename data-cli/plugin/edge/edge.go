package edge

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	edgeOfficialWebsite = "https://www.microsoft.com/edge"
	edgeIconURL         = "https://www.microsoft.com/favicon.ico"
)

type Edge struct{}

func init() {
	plugin.Register(&Edge{})
}

func (p *Edge) Name() string {
	return "edge"
}

func (p *Edge) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "edge",
				Name:            "Microsoft Edge",
				Description:     "Microsoft 出品的 Chromium 内核浏览器，支持跨平台同步。",
				Organization:    "Microsoft",
				OfficialWebsite: edgeOfficialWebsite,
				Icon:            edgeIconURL,
				Tags:            []string{"浏览器"},
				Categories:      []string{"browser"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: edgeOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, edgeOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "Edge 安装包 (x64)", URL: "https://go.microsoft.com/fwlink/?linkid=2108834&Channel=Stable&language=zh-CN"},
							},
						},
						{
							Architecture: "x86",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "Edge 安装包 (x86)", URL: "https://go.microsoft.com/fwlink/?linkid=2108834&Channel=Stable&language=zh-CN&arch=x86"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "Edge 安装包 (arm64)", URL: "https://go.microsoft.com/fwlink/?linkid=2108834&Channel=Stable&language=zh-CN&arch=arm64"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: "Edge 安装包 (macOS)", URL: "https://go.microsoft.com/fwlink/?linkid=2108834&Channel=Stable&language=zh-CN&platform=mac"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: "Edge 安装包 (macOS ARM)", URL: "https://go.microsoft.com/fwlink/?linkid=2108834&Channel=Stable&language=zh-CN&platform=mac&arch=arm64"},
							},
						},
						{
							Architecture: "x64 (deb)",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "direct", Label: "Edge deb 安装包", URL: "https://packages.microsoft.com/repos/edge/pool/main/m/microsoft-edge-stable/microsoft-edge-stable_1.0.0_amd64.deb"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Edge) Disabled() bool { return false }
