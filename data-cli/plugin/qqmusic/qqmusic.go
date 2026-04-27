package qqmusic

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
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
		return nil, fmt.Errorf("qq music: %w", err)
	}

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

	return versions, nil
}

func buildWindowsVersion(item *html.Node, version, releaseDate string) *plugin.Version {
	downloadURL := normalizeSpace(attr(item, `.//a[@data-type='1']`, "data-url"))
	if downloadURL == "" {
		return nil
	}
	if version == "" {
		version = "latest"
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &plugin.Version{
		Version:     "Windows " + version,
		ReleaseDate: releaseDate,
		OfficialURL: qqMusicDownloadURL,
		Variants: []plugin.Variant{
			{
				Architecture: "x64/x86",
				Platform:     "Windows",
				Links: []plugin.Link{
					{Type: "direct", Label: "QQ音乐 Windows 安装包", URL: downloadURL},
				},
			},
		},
	}
}

func buildMacVersion(item *html.Node, version, releaseDate string) *plugin.Version {
	downloadURL := normalizeSpace(attr(item, `.//a[@data-type='2']`, "data-url"))
	if downloadURL == "" {
		return nil
	}
	if version == "" {
		version = "latest"
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &plugin.Version{
		Version:     "macOS " + version,
		ReleaseDate: releaseDate,
		OfficialURL: qqMusicDownloadURL,
		Variants: []plugin.Variant{
			{
				Architecture: "universal",
				Platform:     "macOS",
				Links: []plugin.Link{
					{Type: "direct", Label: "QQ音乐 macOS 安装包 (dmg)", URL: downloadURL},
				},
			},
		},
	}
}

func buildLinuxVersion(item *html.Node, version, releaseDate string) *plugin.Version {
	nodes, _ := htmlquery.QueryAll(item, `.//a[contains(@class,'popup_list__link') and @data-url]`)
	if len(nodes) == 0 {
		return nil
	}
	links := make([]plugin.Link, 0, len(nodes))
	seen := map[string]struct{}{}

	for _, n := range nodes {
		u := normalizeSpace(htmlquery.SelectAttr(n, "data-url"))
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}

		label := normalizeSpace(htmlquery.InnerText(n))
		if label == "" {
			label = fileNameFromURL(u)
		}
		links = append(links, plugin.Link{Type: "direct", Label: label, URL: u})
	}

	if len(links) == 0 {
		return nil
	}
	if version == "" {
		version = "latest"
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &plugin.Version{
		Version:     "Linux " + version,
		ReleaseDate: releaseDate,
		OfficialURL: qqMusicDownloadURL,
		Variants: []plugin.Variant{
			{
				Architecture: "x64",
				Platform:     "Linux",
				Links:        links,
			},
		},
	}
}

func nodeText(root *html.Node, xpath string) string {
	n, err := htmlquery.Query(root, xpath)
	if err != nil || n == nil {
		return ""
	}
	return htmlquery.InnerText(n)
}

func attr(root *html.Node, xpath, key string) string {
	n, err := htmlquery.Query(root, xpath)
	if err != nil || n == nil {
		return ""
	}
	return htmlquery.SelectAttr(n, key)
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

func fileNameFromURL(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return strings.TrimSpace(raw)
	}
	name := strings.TrimSpace(path.Base(parsed.Path))
	if name == "" || name == "." || name == "/" {
		return strings.TrimSpace(raw)
	}
	return name
}
