package thunderbird

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	versionsAPI     = "https://product-details.mozilla.org/1.0/thunderbird_versions.json"
	historyMajorAPI = "https://product-details.mozilla.org/1.0/thunderbird_history_major_releases.json"
)

type versionsResp struct {
	LatestVersion string `json:"LATEST_THUNDERBIRD_VERSION"`
}

type majorHistoryResp map[string]string

type Thunderbird struct{}

func init() {
	plugin.Register(&Thunderbird{})
}

func (p *Thunderbird) Name() string {
	return "thunderbird"
}

func (p *Thunderbird) Fetch() ([]plugin.SoftwareData, error) {
	version, releaseDate, officialURL, err := fetchLatestStable()
	if err != nil {
		return nil, fmt.Errorf("thunderbird: %w", err)
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "thunderbird",
				Name:            "Mozilla Thunderbird",
				Icon:            "https://www.thunderbird.net/media/img/thunderbird/thunderbird-256.png",
				Description:     "Mozilla 出品的开源邮件客户端与日历应用。",
				Organization:    "MZLA Technologies Corporation (Mozilla Foundation)",
				OfficialWebsite: "https://www.thunderbird.net/zh-CN/",
				Tags:            []string{"邮件", "日历", "开源"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: officialURL,
					Platforms: plugin.PlatformsFromVariants(version, releaseDate, officialURL, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "安装包 (exe)", URL: fmt.Sprintf("https://download.mozilla.org/?product=thunderbird-%s-SSL&os=win64&lang=zh-CN", version)},
							},
						},
						{
							Architecture: "x64",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: "安装包 (dmg)", URL: fmt.Sprintf("https://download.mozilla.org/?product=thunderbird-%s-SSL&os=osx&lang=zh-CN", version)},
							},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "direct", Label: "tar.bz2", URL: fmt.Sprintf("https://download.mozilla.org/?product=thunderbird-%s-SSL&os=linux64&lang=zh-CN", version)},
							},
						},
						{
							Architecture: "arm64",
							Platform:     "Android",
							Links: []plugin.Link{
								{Type: "store", Label: "Google Play", URL: "https://play.google.com/store/apps/details?id=net.thunderbird.android"},
								{Type: "store", Label: "F-Droid", URL: "https://f-droid.org/packages/net.thunderbird.android"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func fetchLatestStable() (string, string, string, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(versionsAPI)
	if err != nil {
		return "", "", "", fmt.Errorf("http get versions: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("versions api unexpected status %d", resp.StatusCode)
	}

	var vr versionsResp
	if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
		return "", "", "", fmt.Errorf("decode versions: %w", err)
	}
	version := strings.TrimSpace(vr.LatestVersion)
	if version == "" {
		return "", "", "", fmt.Errorf("empty latest thunderbird version")
	}

	resp2, err := client.Get(historyMajorAPI)
	if err != nil {
		return "", "", "", fmt.Errorf("http get major history: %w", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("major history api unexpected status %d", resp2.StatusCode)
	}

	var history majorHistoryResp
	if err := json.NewDecoder(resp2.Body).Decode(&history); err != nil {
		return "", "", "", fmt.Errorf("decode major history: %w", err)
	}

	releaseDate := strings.TrimSpace(history[version])
	if releaseDate == "" {
		return "", "", "", fmt.Errorf("release date not found for version %s", version)
	}

	officialURL := "https://www.thunderbird.net/notes/"
	return version, releaseDate, officialURL, nil
}

func (p *Thunderbird) Disabled() bool { return false }
