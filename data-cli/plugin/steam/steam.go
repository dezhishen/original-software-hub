package steam

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	steamOfficialWebsite = "https://store.steampowered.com/about/"
	steamIconURL         = "https://store.fastly.steamstatic.com/public/shared/images/responsive/share_steam_logo.png"

	steamWindowsURL = "https://media.st.dl.eccdnx.com/client/installer/SteamSetup.exe"
	steamMacOSURL   = "https://media.st.dl.eccdnx.com/client/installer/steam.dmg"
	steamLinuxURL   = "https://media.st.dl.eccdnx.com/client/installer/steam.deb"
)

// Steam implements plugin.Plugin for Valve Steam client.
type Steam struct{}

func init() {
	plugin.Register(&Steam{})
}

func (s *Steam) Name() string {
	return "steam"
}

func (s *Steam) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate, err := detectReleaseDate([]string{steamWindowsURL, steamMacOSURL, steamLinuxURL})
	if err != nil {
		return nil, fmt.Errorf("detect steam release date: %w", err)
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "steam",
				Name:            "Steam",
				Icon:            steamIconURL,
				Description:     "Valve 旗下数字游戏平台客户端。",
				Organization:    "Valve Corporation",
				OfficialWebsite: steamOfficialWebsite,
				Tags:            []string{"游戏平台", "娱乐"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: steamOfficialWebsite,
					Variants: []plugin.Variant{
						buildVariant("x86/x64", "Windows", steamWindowsURL),
						buildVariant("x64", "macOS", steamMacOSURL),
						buildVariant("x64", "Linux", steamLinuxURL),
					},
				},
			},
		},
	}, nil
}

func buildVariant(arch, platform, downloadURL string) plugin.Variant {
	url := strings.TrimSpace(downloadURL)
	return plugin.Variant{
		Architecture: arch,
		Platform:     platform,
		Links: []plugin.Link{
			{
				Type:  "direct",
				Label: fileNameFromURL(url),
				URL:   url,
			},
		},
	}
}

func detectReleaseDate(urls []string) (string, error) {
	client := &http.Client{Timeout: 12 * time.Second}
	latest := time.Time{}

	for _, raw := range urls {
		if strings.TrimSpace(raw) == "" {
			continue
		}

		req, err := http.NewRequest(http.MethodHead, raw, nil)
		if err != nil {
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		lastModified := strings.TrimSpace(resp.Header.Get("Last-Modified"))
		resp.Body.Close()
		if lastModified == "" {
			continue
		}

		parsed, err := time.Parse(time.RFC1123, lastModified)
		if err != nil {
			continue
		}
		if parsed.After(latest) {
			latest = parsed
		}
	}

	if latest.IsZero() {
		// Last-Modified 不是强保证字段，缺失时回退为当前日期，避免插件失败。
		return time.Now().Format("2006-01-02"), nil
	}

	return latest.UTC().Format("2006-01-02"), nil
}

func fileNameFromURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return strings.TrimSpace(raw)
	}
	name := strings.TrimSpace(path.Base(parsed.Path))
	if name == "" || name == "." || name == "/" {
		return strings.TrimSpace(raw)
	}
	return name
}
