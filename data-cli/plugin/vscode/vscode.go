package vscode

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	vscodeOfficialWebsite = "https://code.visualstudio.com"
	vscodeIconURL         = "https://code.visualstudio.com/favicon.ico"
	vscodeUpdateAPI       = "https://update.code.visualstudio.com/api/update/win32-x64/stable/latest"
)

type Vscode struct{}

func init() {
	plugin.Register(&Vscode{})
}

func (p *Vscode) Name() string {
	return "vscode"
}

func (p *Vscode) Fetch() ([]plugin.SoftwareData, error) {
	version := "latest"
	releaseDate := time.Now().UTC().Format("2006-01-02")

	// Try to detect version from update API
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(vscodeUpdateAPI)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			var result struct {
				Version string `json:"productVersion"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && result.Version != "" {
				version = result.Version
			}
		}
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "vscode",
				Name:            "Visual Studio Code",
				Description:     "Microsoft 出品的轻量级开源代码编辑器，支持丰富的扩展生态。",
				Organization:    "Microsoft",
				OfficialWebsite: vscodeOfficialWebsite,
				Icon:            vscodeIconURL,
				Tags:            []string{"开发工具", "代码编辑器"},
				Categories:      []string{"development"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: vscodeOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants(version, releaseDate, vscodeOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "VSCode 用户安装包 (x64)", URL: "https://code.visualstudio.com/sha/download?build=stable&os=win32-x64-user"},
								{Type: "direct", Label: "VSCode 系统安装包 (x64)", URL: "https://code.visualstudio.com/sha/download?build=stable&os=win32-x64"},
							},
						},
						{
							Architecture: "x86",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "VSCode 用户安装包 (x86)", URL: "https://code.visualstudio.com/sha/download?build=stable&os=win32-x86-user"},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "VSCode 用户安装包 (arm64)", URL: "https://code.visualstudio.com/sha/download?build=stable&os=win32-arm64-user"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: "VSCode for macOS (Universal)", URL: "https://code.visualstudio.com/sha/download?build=stable&os=darwin-universal"},
							},
						},
						{
							Architecture: "x64 (deb)",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "direct", Label: "VSCode .deb (amd64)", URL: "https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64"},
							},
						},
						{
							Architecture: "x64 (rpm)",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "direct", Label: "VSCode .rpm (amd64)", URL: "https://code.visualstudio.com/sha/download?build=stable&os=linux-rpm-x64"},
							},
						},
						{
							Architecture: "arm64 (deb)",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "direct", Label: "VSCode .deb (arm64)", URL: "https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-arm64"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Vscode) Disabled() bool { return false }
