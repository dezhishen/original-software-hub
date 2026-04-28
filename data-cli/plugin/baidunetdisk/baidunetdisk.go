package baidunetdisk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	baiduNetdiskOfficialWebsite = "https://pan.baidu.com/"
	baiduNetdiskVersionPage     = "https://yun.baidu.com/disk/version"
	baiduNetdiskChangeLogAPI    = "https://yun.baidu.com/disk/cmsdata?do=changelog&platform=%s&page=1&num=1"
	baiduNetdiskIconURL         = "https://nd-static.bdstatic.com/box-static/lite-clouddisk-ui/res/static/images/favicon.ico"
)

var versionPattern = regexp.MustCompile(`V([0-9]+(?:\.[0-9]+)+)`)

// BaiduNetdisk implements plugin.Plugin for Baidu Netdisk client.
type BaiduNetdisk struct{}

func init() {
	plugin.Register(&BaiduNetdisk{})
}

func (b *BaiduNetdisk) Name() string {
	return "baidunetdisk"
}

func (b *BaiduNetdisk) Fetch() ([]plugin.SoftwareData, error) {
	winEntry, err := fetchLatestEntry("guanjia")
	if err != nil {
		return nil, fmt.Errorf("fetch windows changelog: %w", err)
	}
	macEntry, err := fetchLatestEntry("mac")
	if err != nil {
		return nil, fmt.Errorf("fetch mac changelog: %w", err)
	}

	version := extractVersion(winEntry.Version)
	if version == "" {
		version = extractVersion(macEntry.Version)
	}
	if version == "" {
		version = "Latest"
	}

	releaseDate := parsePublishDate(winEntry.Publish)
	if releaseDate == "" {
		releaseDate = parsePublishDate(macEntry.Publish)
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	variants := []plugin.Variant{}
	if winEntry.URL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows",
			Links:        []plugin.Link{{Type: "direct", Label: "百度网盘 Windows 下载", URL: winEntry.URL}},
		})
	}
	if winEntry.URLLegacy != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x86",
			Platform:     "Windows",
			Links:        []plugin.Link{{Type: "direct", Label: "百度网盘 Windows 经典版", URL: winEntry.URLLegacy}},
		})
	}
	if macEntry.URL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "macOS",
			Links:        []plugin.Link{{Type: "direct", Label: "百度网盘 macOS 下载", URL: macEntry.URL}},
		})
	}
	if macEntry.URLLegacy != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "arm64",
			Platform:     "macOS",
			Links:        []plugin.Link{{Type: "direct", Label: "百度网盘 macOS ARM64 下载", URL: macEntry.URLLegacy}},
		})
	}
	if len(variants) == 0 {
		variants = append(variants, plugin.Variant{
			Architecture: "通用",
			Platform:     "Web",
			Links:        []plugin.Link{{Type: "webpage", Label: "百度网盘版本页", URL: baiduNetdiskVersionPage}},
		})
	}

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
					OfficialURL: baiduNetdiskVersionPage,
					Platforms:   plugin.PlatformsFromVariants(version, releaseDate, baiduNetdiskVersionPage, variants),
				},
			},
		},
	}, nil
}

type changelogResponse struct {
	ErrorNo int              `json:"errorno"`
	List    []changelogEntry `json:"list"`
}

type changelogEntry struct {
	Publish   string `json:"publish"`
	Version   string `json:"version"`
	URL       string `json:"url"`
	URLLegacy string `json:"url_1"`
}

func fetchLatestEntry(platform string) (*changelogEntry, error) {
	api := fmt.Sprintf(baiduNetdiskChangeLogAPI, platform)
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(api)
	if err != nil {
		return nil, fmt.Errorf("http get %s: %w", api, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}

	var payload changelogResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode changelog: %w", err)
	}
	if payload.ErrorNo != 0 || len(payload.List) == 0 {
		return nil, fmt.Errorf("empty changelog for platform %s", platform)
	}
	entry := payload.List[0]
	entry.URL = strings.TrimSpace(entry.URL)
	entry.URLLegacy = strings.TrimSpace(entry.URLLegacy)
	entry.Version = strings.TrimSpace(entry.Version)
	entry.Publish = strings.TrimSpace(entry.Publish)
	return &entry, nil
}

func extractVersion(raw string) string {
	raw = strings.TrimSpace(raw)
	if m := versionPattern.FindStringSubmatch(raw); len(m) >= 2 {
		return m[1]
	}
	return ""
}

func parsePublishDate(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02 15:04:05", raw)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}
