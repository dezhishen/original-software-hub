package anydesk

import (
	"fmt"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	anydeskOfficialWebsite = "https://anydesk.com.cn"
	anydeskIconURL         = "https://anydesk.com.cn/favicon.ico"
)

// Anydesk implements plugin.Plugin for AnyDesk remote control client.
type Anydesk struct{}

func init() {
	plugin.Register(&Anydesk{})
}

func (p *Anydesk) Name() string {
	return "anydesk"
}

func (p *Anydesk) Fetch() ([]plugin.SoftwareData, error) {
	dm, err := fetchAnyDeskDownloads()
	if err != nil {
		return nil, fmt.Errorf("fetch anydesk downloads: %w", err)
	}

	today := time.Now().UTC().Format("2006-01-02")
	versions := make([]plugin.Version, 0, 5)

	// Windows
	if win, ok := dm["windows"]; ok && win.URL != "" {
		v := strings.TrimSpace(win.Version)
		versions = append(versions, plugin.Version{
			Version:     v,
			ReleaseDate: today,
			OfficialURL: anydeskDownloadPage,
			Platforms: plugin.PlatformsFromVariants(v, today, anydeskDownloadPage, []plugin.Variant{
				{Architecture: "x64", Platform: "Windows", Links: []plugin.Link{
					{Type: "direct", Label: "AnyDesk Windows 安装包", URL: win.URL},
				}},
			}),
		})
	}

	// macOS
	if mac, ok := dm["mac"]; ok && mac.URL != "" {
		v := strings.TrimSpace(mac.Version)
		versions = append(versions, plugin.Version{
			Version:     v,
			ReleaseDate: today,
			OfficialURL: anydeskDownloadPage,
			Platforms: plugin.PlatformsFromVariants(v, today, anydeskDownloadPage, []plugin.Variant{
				{Architecture: "universal", Platform: "macOS", Links: []plugin.Link{
					{Type: "direct", Label: "AnyDesk macOS 安装包", URL: mac.URL},
				}},
			}),
		})
	}

	// Linux — emit each package as a separate variant
	if lnx, ok := dm["linux"]; ok && len(lnx.Packages) > 0 {
		v := strings.TrimSpace(lnx.Version)
		variants := make([]plugin.Variant, 0, len(lnx.Packages))
		for _, pkg := range lnx.Packages {
			if pkg.URL == "" {
				continue
			}
			variants = append(variants, plugin.Variant{
				Architecture: linuxArch(pkg.ID),
				Platform:     "Linux",
				Links: []plugin.Link{
					{Type: "direct", Label: pkg.Name, URL: pkg.URL},
				},
			})
		}
		if len(variants) > 0 {
			officialLinux := "https://anydesk.com.cn/zhs/downloads/linux"
			versions = append(versions, plugin.Version{
				Version:     v,
				ReleaseDate: today,
				OfficialURL: officialLinux,
				Platforms:   plugin.PlatformsFromVariants(v, today, officialLinux, variants),
			})
		}
	}

	// Android (store link)
	if and, ok := dm["android"]; ok && and.URL != "" {
		v := strings.TrimSpace(and.Version)
		versions = append(versions, plugin.Version{
			Version:     v,
			ReleaseDate: today,
			OfficialURL: and.URL,
			Platforms: plugin.PlatformsFromVariants(v, today, and.URL, []plugin.Variant{
				{Architecture: "all", Platform: "Android", Links: []plugin.Link{
					{Type: "store", Label: "Google Play", URL: and.URL},
				}},
			}),
		})
	}

	// iOS (store link)
	if ios, ok := dm["ios"]; ok && ios.URL != "" {
		v := strings.TrimSpace(ios.Version)
		versions = append(versions, plugin.Version{
			Version:     v,
			ReleaseDate: today,
			OfficialURL: ios.URL,
			Platforms: plugin.PlatformsFromVariants(v, today, ios.URL, []plugin.Variant{
				{Architecture: "all", Platform: "iOS / iPadOS", Links: []plugin.Link{
					{Type: "store", Label: "App Store", URL: ios.URL},
				}},
			}),
		})
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "anydesk",
				Name:            "AnyDesk",
				Description:     "轻量高速的远程控制软件，支持跨平台访问与控制。",
				Organization:    "AnyDesk Software GmbH",
				OfficialWebsite: anydeskOfficialWebsite,
				Icon:            anydeskIconURL,
				Tags:            []string{"远程控制"},
			},
			Versions: versions,
		},
	}, nil
}

// linuxArch maps the package id field to a human-readable architecture label.
func linuxArch(id string) string {
	switch {
	case strings.Contains(id, "arm64") || strings.HasSuffix(id, "_arm64"):
		return "arm64"
	case strings.Contains(id, "arm") || strings.HasSuffix(id, "_arm"):
		return "arm"
	default:
		return "x64"
	}
}

func (p *Anydesk) Disabled() bool { return false }
