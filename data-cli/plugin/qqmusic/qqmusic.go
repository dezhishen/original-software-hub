package qqmusic

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/dezhishen/original-software-hub/data-cli/plugin"
	"golang.org/x/net/html"
)

const (
	qqMusicDownloadURL = "https://y.qq.com/download/download.html"
	qqMusicWebsiteURL  = "https://y.qq.com/"
	qqMusicIconURL     = "https://y.qq.com/favicon.ico"
	qqMusicWindowsPage = "https://y.qq.com/download/welcome_pc_v15/index.html?ADTAG=YQQ"
	qqMusicMacPage     = "https://y.qq.com/download/mac.html?part=1&ADTAG=YQQ"
)

const (
	qqMusicIOSStoreURL    = "https://itunes.apple.com/cn/app/qq-yin-le/id414603431?mt=8"
	qqMusicIPadStoreURL   = "https://itunes.apple.com/cn/app/qq-yin-lehd-du-bo-zhong-guo/id429885089?l=en&mt=8"
	qqMusicAndroidPageURL = "https://y.qq.com/download/download.html"
)

var (
	qqMusicVersionPattern = regexp.MustCompile(`最新\s*版\s*:\s*([0-9]+(?:\.[0-9]+)*)`)
	qqMusicDatePattern    = regexp.MustCompile(`发布时间\s*[：:]\s*(20\d{2}-\d{2}-\d{2})`)
)

// QQMusic implements plugin.Plugin for Tencent QQ Music.
type QQMusic struct{}

func init() {
	plugin.Register(&QQMusic{})
}

func (q *QQMusic) Name() string {
	return "qqmusic"
}

func (q *QQMusic) Fetch() ([]plugin.SoftwareData, error) {
	versions, err := fetchQQMusicDesktopVersions()
	if err != nil {
		return nil, err
	}
	// Append mobile version
	today := time.Now().UTC().Format("2006-01-02")
	mobileVersion := plugin.Version{
		Version:     "latest",
		ReleaseDate: today,
		OfficialURL: qqMusicDownloadURL,
		Platforms: plugin.PlatformsFromVariants("latest", today, qqMusicDownloadURL, []plugin.Variant{
			{
				Architecture: "universal",
				Platform:     "iOS / iPadOS",
				Links: []plugin.Link{
					{Type: "store", Label: "App Store (iPhone)", URL: qqMusicIOSStoreURL},
					{Type: "store", Label: "App Store (iPad)", URL: qqMusicIPadStoreURL},
				},
			},
			{
				Architecture: "arm64",
				Platform:     "Android",
				Links: []plugin.Link{
					{Type: "webpage", Label: "QQ音乐 Android 下载页", URL: qqMusicAndroidPageURL},
				},
			},
		}),
	}
	versions = append(versions, mobileVersion)
	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "qqmusic",
				Name:            "QQ音乐",
				Icon:            qqMusicIconURL,
				Description:     "腾讯旗下在线音乐平台客户端。",
				Organization:    "Tencent Music",
				OfficialWebsite: qqMusicWebsiteURL,
				Tags:            []string{"音乐", "流媒体"},
			},
			Versions: versions,
		},
	}, nil
}

func fetchQQMusicDesktopVersions() ([]plugin.Version, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(qqMusicDownloadURL)
	if err != nil {
		return nil, fmt.Errorf("http get %s: %w", qqMusicDownloadURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	pcSection, err := htmlquery.Query(doc, `//div[contains(@class,'product--pc')]`)
	if err != nil || pcSection == nil {
		return nil, fmt.Errorf("desktop section not found")
	}

	items, err := htmlquery.QueryAll(pcSection, `.//li[contains(@class,'product_list__item')]`)
	if err != nil {
		return nil, fmt.Errorf("query desktop items: %w", err)
	}

	versions := make([]plugin.Version, 0, 3)
	for _, item := range items {
		title := normalizeSpace(nodeText(item, `.//h3[contains(@class,'product_list__tit')]`))
		if title == "" {
			continue
		}
		versionText := normalizeSpace(nodeText(item, `.//span[contains(@class,'product_list__version')]`))
		version := versionFromText(versionText)
		releaseDate := releaseDateFromItem(item)

		switch {
		case strings.Contains(strings.ToLower(title), "windows"):
			if v := buildWindowsVersion(item, version, releaseDate); v != nil {
				versions = append(versions, *v)
			}
		case strings.Contains(strings.ToLower(title), "mac"):
			if v := buildMacVersion(item, version, releaseDate); v != nil {
				versions = append(versions, *v)
			}
		case strings.Contains(strings.ToLower(title), "linux"):
			if v := buildLinuxVersion(item, version, releaseDate); v != nil {
				versions = append(versions, *v)
			}
		}
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no desktop versions found")
	}
	versions = mergeVersionsAsTabbed(versions)

	return versions, nil
}

func mergeVersionsAsTabbed(versions []plugin.Version) []plugin.Version {
	if len(versions) <= 1 {
		return versions
	}

	platforms := make([]plugin.PlatformRelease, 0, len(versions))
	latestDate := ""
	for _, version := range versions {
		platforms = append(platforms, version.Platforms...)
		if strings.TrimSpace(version.ReleaseDate) > latestDate {
			latestDate = strings.TrimSpace(version.ReleaseDate)
		}
	}
	if latestDate == "" {
		latestDate = time.Now().UTC().Format("2006-01-02")
	}

	return []plugin.Version{{
		Version:     "latest",
		ReleaseDate: latestDate,
		OfficialURL: qqMusicDownloadURL,
		Platforms:   platforms,
	}}
}

func buildWindowsVersion(item *html.Node, version, releaseDate string) *plugin.Version {
	if version == "" {
		version = "latest"
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &plugin.Version{
		Version:     version,
		ReleaseDate: releaseDate,
		OfficialURL: qqMusicDownloadURL,
		Platforms: plugin.PlatformsFromVariants(version, releaseDate, qqMusicDownloadURL, []plugin.Variant{
			{
				Architecture: "x64/x86",
				Platform:     "Windows",
				Links: []plugin.Link{
					{Type: "webpage", Label: "QQ音乐 Windows 下载页", URL: qqMusicWindowsPage},
				},
			},
		}),
	}
}

func buildMacVersion(item *html.Node, version, releaseDate string) *plugin.Version {
	if version == "" {
		version = "latest"
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &plugin.Version{
		Version:     version,
		ReleaseDate: releaseDate,
		OfficialURL: qqMusicDownloadURL,
		Platforms: plugin.PlatformsFromVariants(version, releaseDate, qqMusicDownloadURL, []plugin.Variant{
			{
				Architecture: "universal",
				Platform:     "macOS",
				Links: []plugin.Link{
					{Type: "webpage", Label: "QQ音乐 macOS 下载页", URL: qqMusicMacPage},
				},
			},
		}),
	}
}

func buildLinuxVersion(item *html.Node, version, releaseDate string) *plugin.Version {
	if version == "" {
		version = "latest"
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &plugin.Version{
		Version:     version,
		ReleaseDate: releaseDate,
		OfficialURL: qqMusicDownloadURL,
		Platforms: plugin.PlatformsFromVariants(version, releaseDate, qqMusicDownloadURL, []plugin.Variant{
			{
				Architecture: "x64",
				Platform:     "Linux",
				Links: []plugin.Link{
					{Type: "webpage", Label: "QQ音乐 Linux 下载页", URL: qqMusicDownloadURL},
				},
			},
		}),
	}
}

func nodeText(root *html.Node, xpath string) string {
	n, err := htmlquery.Query(root, xpath)
	if err != nil || n == nil {
		return ""
	}
	return htmlquery.InnerText(n)
}

func versionFromText(s string) string {
	m := qqMusicVersionPattern.FindStringSubmatch(normalizeSpace(s))
	if len(m) >= 2 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func releaseDateFromItem(item *html.Node) string {
	nodes, _ := htmlquery.QueryAll(item, `.//li[contains(@class,'version_list__item')]`)
	for _, n := range nodes {
		text := normalizeSpace(htmlquery.InnerText(n))
		if m := qqMusicDatePattern.FindStringSubmatch(text); len(m) >= 2 {
			return m[1]
		}
	}
	return ""
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
