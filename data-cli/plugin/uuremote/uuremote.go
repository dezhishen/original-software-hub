package uuremote

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	uuremoteOfficialWebsite = "https://uuyc.163.com/"
	uuremoteDownloadPage    = "https://uuyc.163.com/download/"
	uuremoteIconURL         = "https://uuyc.res.netease.com/pc/gw/20241129172408/img/logo_f05fa07b.png"
	uuremoteIOSStoreURL     = "https://apps.apple.com/cn/app/uu%E8%BF%9C%E7%A8%8B/id1642306791"
	uuAndroidTVDownloadURL  = "https://adl.netease.com/d/g/uuremote/c/adtv"
)

var (
	uuDownloadJSURLPattern = regexp.MustCompile(`https://uuyc\.res\.netease\.com/pc/gw/[0-9]+/js/download/index_[a-z0-9]+\.js`)
	uuVersionPageMapRe     = regexp.MustCompile(`(?s)e="(https://uuyc\.163\.com/download/page/[^"]+\.html)",n="(https://uuyc\.163\.com/download/page/[^"]+\.html)",o="(https://uuyc\.163\.com/download/page/[^"]+\.html)",r="(https://uuyc\.163\.com/download/page/[^"]+\.html)".*?p\("(https://uuyc\.163\.com/download/page/[^"]+\.html)"\)`)
	uuMetaKeywordsPattern  = regexp.MustCompile(`(?is)<meta[^>]+name=["']keywords["'][^>]+content=["']([^"']+)["']`)
	uuMetaDescPattern      = regexp.MustCompile(`(?is)<meta[^>]+name=["']description["'][^>]+content=["']([^"']+)["']`)
	uuWindowsPattern       = regexp.MustCompile(`(?is)<div[^>]+id="js_Btn_windows"[^>]*>([^<]+)</div>`)
	uuMacPattern           = regexp.MustCompile(`(?is)<div[^>]+id="js_Btn_mac"[^>]*>([^<]+)</div>`)
	uuAndroidPattern       = regexp.MustCompile(`(?is)<div[^>]+id="js_Btn_android"[^>]*>([^<]+)</div>`)
	uuIOSPattern           = regexp.MustCompile(`(?is)<div[^>]+id="js_Btn_ios"[^>]*>([^<]+)</div>`)
	uuVersionGeneric       = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)+)`)
	uuVersionPattern       = regexp.MustCompile(`当前版本：\s*V\s*([0-9]+(?:\.[0-9]+)+)`)
	uuDatePattern          = regexp.MustCompile(`更新于\s*([0-9]{4})\.([0-9]{2})\.([0-9]{2})`)
)

type uuPlatformMeta struct {
	Version     string
	ReleaseDate string
}

type uuPlatformConfig struct {
	Key          string
	PlatformName string
	Arch         string
	URL          string
	StoreURL     string
}

// UURemote implements plugin.Plugin for NetEase UU Remote client.
type UURemote struct{}

func init() {
	plugin.Register(&UURemote{})
}

func (u *UURemote) Name() string {
	return "uuremote"
}

func (u *UURemote) Fetch() ([]plugin.SoftwareData, error) {
	html, err := fetchPageHTML(uuremoteDownloadPage)
	if err != nil {
		return nil, fmt.Errorf("fetch uu remote page: %w", err)
	}

	windowsURL := findMatch(html, uuWindowsPattern)
	macURL := findMatch(html, uuMacPattern)
	androidURL := findMatch(html, uuAndroidPattern)
	iosDownloadURL := findMatch(html, uuIOSPattern)
	platformMeta := map[string]uuPlatformMeta{}
	if jsURL := findFullMatch(html, uuDownloadJSURLPattern); jsURL != "" {
		jsCode, err := fetchPageHTML(jsURL)
		if err == nil {
			for k, pageURL := range extractPlatformVersionPages(jsCode) {
				meta, err := fetchPlatformMeta(pageURL)
				if err == nil {
					platformMeta[k] = meta
				}
			}
		}
	}

	platforms := []uuPlatformConfig{
		{Key: "windows", PlatformName: "Windows", Arch: "x64", URL: windowsURL},
		{Key: "macos", PlatformName: "macOS", Arch: "x64", URL: macURL},
		{Key: "ios", PlatformName: "iOS / iPadOS", Arch: "universal", URL: iosDownloadURL, StoreURL: uuremoteIOSStoreURL},
		{Key: "android", PlatformName: "Android", Arch: "arm64", URL: androidURL},
		{Key: "androidtv", PlatformName: "Android TV", Arch: "arm64", URL: uuAndroidTVDownloadURL},
	}

	versions := make([]plugin.Version, 0, len(platforms))
	for _, cfg := range platforms {
		if strings.TrimSpace(cfg.URL) == "" {
			continue
		}
		meta := platformMeta[cfg.Key]
		version := strings.TrimSpace(meta.Version)
		if version == "" {
			version = "Latest"
		}
		releaseDate := strings.TrimSpace(meta.ReleaseDate)
		if releaseDate == "" {
			releaseDate = time.Now().UTC().Format("2006-01-02")
		}

		links := []plugin.Link{{Type: "direct", Label: fmt.Sprintf("UU远程 %s 下载", cfg.PlatformName), URL: cfg.URL}}
		if cfg.StoreURL != "" {
			links = append([]plugin.Link{{Type: "store", Label: "App Store", URL: cfg.StoreURL}}, links...)
		}

		versions = append(versions, plugin.Version{
			Version:     version,
			ReleaseDate: releaseDate,
			OfficialURL: uuremoteDownloadPage,
			Variants: []plugin.Variant{
				{
					Architecture: cfg.Arch,
					Platform:     cfg.PlatformName,
					Links:        links,
				},
			},
		})
	}

	if len(versions) == 0 {
		fallbackVersion := findMatch(html, uuVersionPattern)
		if fallbackVersion == "" {
			fallbackVersion = "Latest"
		}
		fallbackDate := parseReleaseDate(html)
		if fallbackDate == "" {
			fallbackDate = time.Now().UTC().Format("2006-01-02")
		}
		versions = append(versions, plugin.Version{
			Version:     fallbackVersion,
			ReleaseDate: fallbackDate,
			OfficialURL: uuremoteDownloadPage,
			Variants: []plugin.Variant{
				{Architecture: "x64", Platform: "Windows", Links: []plugin.Link{{Type: "direct", Label: "UU远程 Windows 下载", URL: windowsURL}}},
			},
		})
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return versions[i].Variants[0].Platform < versions[j].Variants[0].Platform
	})

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "uu-remote",
				Name:            "UU远程",
				Icon:            uuremoteIconURL,
				Description:     "网易 UU 远程控制客户端，支持低延迟远控与多端协作。",
				Organization:    "NetEase",
				OfficialWebsite: uuremoteOfficialWebsite,
				Tags:            []string{"远程控制", "远程办公", "网易"},
			},
			Versions: versions,
		},
	}, nil
}

func extractPlatformVersionPages(jsCode string) map[string]string {
	out := map[string]string{}
	m := uuVersionPageMapRe.FindStringSubmatch(jsCode)
	if len(m) < 6 {
		return out
	}
	out["macos"] = strings.TrimSpace(m[1])
	out["ios"] = strings.TrimSpace(m[2])
	out["android"] = strings.TrimSpace(m[3])
	out["androidtv"] = strings.TrimSpace(m[4])
	out["windows"] = strings.TrimSpace(m[5])
	return out
}

func fetchPlatformMeta(pageURL string) (uuPlatformMeta, error) {
	html, err := fetchPageHTML(strings.TrimSpace(pageURL))
	if err != nil {
		return uuPlatformMeta{}, err
	}
	keywords := findMatch(html, uuMetaKeywordsPattern)
	desc := findMatch(html, uuMetaDescPattern)

	version := extractVersionText(keywords)
	if version == "" {
		version = extractVersionText(desc)
	}
	date := parseReleaseDate(desc)
	if date == "" {
		date = parseReleaseDate(keywords)
	}
	return uuPlatformMeta{Version: version, ReleaseDate: date}, nil
}

func extractVersionText(s string) string {
	m := uuVersionPattern.FindStringSubmatch(s)
	if len(m) >= 2 {
		return strings.TrimSpace(m[1])
	}
	m = uuVersionGeneric.FindStringSubmatch(s)
	if len(m) >= 2 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

func fetchPageHTML(pageURL string) (string, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(pageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func findMatch(s string, re *regexp.Regexp) string {
	m := re.FindStringSubmatch(s)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}

func findFullMatch(s string, re *regexp.Regexp) string {
	m := re.FindString(s)
	return strings.TrimSpace(m)
}

func parseReleaseDate(s string) string {
	m := uuDatePattern.FindStringSubmatch(s)
	if len(m) < 4 {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", m[1], m[2], m[3])
}

func (x *UURemote) CompareWithPrevious(previous plugin.PreviousState) ([]plugin.FetchResult, error) {
	items, err := x.Fetch()
	if err != nil {
		return nil, err
	}
	return plugin.BuildCompareResults(items, previous), nil
}
