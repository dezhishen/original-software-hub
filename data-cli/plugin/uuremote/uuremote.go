package uuremote

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	uuremoteOfficialWebsite = "https://uuyc.163.com/"
	uuremoteDownloadPage    = "https://uuyc.163.com/download/"
	uuremoteIconURL         = "https://uuyc.res.netease.com/pc/gw/20241129172408/img/logo_f05fa07b.png"
	uuremoteIOSStoreURL     = "https://apps.apple.com/cn/app/uu%E8%BF%9C%E7%A8%8B/id1642306791"
)

var (
	uuWindowsPattern = regexp.MustCompile(`(?is)<div[^>]+id="js_Btn_windows"[^>]*>([^<]+)</div>`)
	uuMacPattern     = regexp.MustCompile(`(?is)<div[^>]+id="js_Btn_mac"[^>]*>([^<]+)</div>`)
	uuAndroidPattern = regexp.MustCompile(`(?is)<div[^>]+id="js_Btn_android"[^>]*>([^<]+)</div>`)
	uuIOSPattern     = regexp.MustCompile(`(?is)<div[^>]+id="js_Btn_ios"[^>]*>([^<]+)</div>`)
	uuVersionPattern = regexp.MustCompile(`当前版本：\s*V\s*([0-9]+(?:\.[0-9]+)+)`)
	uuDatePattern    = regexp.MustCompile(`更新于\s*([0-9]{4})\.([0-9]{2})\.([0-9]{2})`)
)

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

	version := findMatch(html, uuVersionPattern)
	if version == "" {
		version = "Latest"
	}

	releaseDate := parseReleaseDate(html)
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	variants := make([]plugin.Variant, 0, 5)
	if windowsURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows",
			Links:        []plugin.Link{{Type: "direct", Label: "UU远程 Windows 下载", URL: windowsURL}},
		})
	}
	if macURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "macOS",
			Links:        []plugin.Link{{Type: "direct", Label: "UU远程 macOS 下载", URL: macURL}},
		})
	}
	if androidURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "arm64",
			Platform:     "Android",
			Links:        []plugin.Link{{Type: "direct", Label: "UU远程 Android 下载", URL: androidURL}},
		})
	}
	if iosDownloadURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "universal",
			Platform:     "iOS / iPadOS",
			Links: []plugin.Link{
				{Type: "store", Label: "App Store", URL: uuremoteIOSStoreURL},
				{Type: "webpage", Label: "iOS 下载页", URL: iosDownloadURL},
			},
		})
	}
	variants = append(variants, plugin.Variant{
		Architecture: "universal",
		Platform:     "Web",
		Links:        []plugin.Link{{Type: "webpage", Label: "UU远程官方下载页", URL: uuremoteDownloadPage}},
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
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: uuremoteDownloadPage,
					Variants:    variants,
				},
			},
		},
	}, nil
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

func parseReleaseDate(s string) string {
	m := uuDatePattern.FindStringSubmatch(s)
	if len(m) < 4 {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", m[1], m[2], m[3])
}
