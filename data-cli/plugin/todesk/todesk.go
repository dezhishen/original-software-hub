package todesk

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
	todeskOfficialWebsite = "https://www.todesk.com/"
	todeskDownloadPage    = "https://www.todesk.com/download.html"
	todeskIconURL         = "https://www.todesk.com/favicon.ico"
)

var (
	todeskWindowsLinkPattern = regexp.MustCompile(`https://dl\.todesk\.com/irrigation/ToDesk_[0-9]+(?:\.[0-9]+)+\.exe`)
	todeskMacLinkPattern     = regexp.MustCompile(`https://dl\.todesk\.com/macos/ToDesk_[0-9]+(?:\.[0-9]+)+\.pkg`)
	todeskAndroidPattern     = regexp.MustCompile(`https://dl\.todesk\.com/android/ToDesk_[0-9]+(?:\.[0-9]+)+\.apk`)
	todeskTVPattern          = regexp.MustCompile(`https://dl\.todesk\.com/android/ToDesk_TV_[0-9]+(?:\.[0-9]+)+\.apk`)
	todeskIOSPattern         = regexp.MustCompile(`https://apps\.apple\.com/cn/app/todesk/id[0-9]+`)
	todeskLinuxPattern       = regexp.MustCompile(`https://dl\.todesk\.com/linux/[a-zA-Z0-9._-]*amd64\.deb`)
	todeskVersionPattern     = regexp.MustCompile(`ToDesk_([0-9]+(?:\.[0-9]+)+)\.`)
	todeskReleaseDatePattern = regexp.MustCompile(`win_release_date:"([0-9]+(?:\.[0-9]+){2})"`)
)

// ToDesk implements plugin.Plugin for ToDesk remote control client.
type ToDesk struct{}

func init() {
	plugin.Register(&ToDesk{})
}

func (t *ToDesk) Name() string {
	return "todesk"
}

func (t *ToDesk) Fetch() ([]plugin.SoftwareData, error) {
	html, err := fetchToDeskHTML()
	if err != nil {
		return nil, fmt.Errorf("fetch todesk download page: %w", err)
	}

	decodedHTML := decodeNuxtString(html)
	windowsURL := findFullMatch(decodedHTML, todeskWindowsLinkPattern)
	macURL := findFullMatch(decodedHTML, todeskMacLinkPattern)
	androidURL := findFullMatch(decodedHTML, todeskAndroidPattern)
	tvURL := findFullMatch(decodedHTML, todeskTVPattern)
	iosURL := findFullMatch(decodedHTML, todeskIOSPattern)
	linuxURL := findFullMatch(decodedHTML, todeskLinuxPattern)

	version := firstNonEmpty(
		extractVersion(windowsURL),
		extractVersion(macURL),
		extractVersion(androidURL),
	)
	if version == "" {
		version = "Latest"
	}

	releaseDate := parseDotDate(findMatch(html, todeskReleaseDatePattern))
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}
	variants := make([]plugin.Variant, 0, 6)
	if windowsURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows",
			Links:        []plugin.Link{{Type: "direct", Label: "ToDesk Windows 安装包", URL: windowsURL}},
		})
	}
	if macURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "universal",
			Platform:     "macOS",
			Links:        []plugin.Link{{Type: "direct", Label: "ToDesk macOS 安装包", URL: macURL}},
		})
	}
	if linuxURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Linux",
			Links:        []plugin.Link{{Type: "webpage", Label: "ToDesk Linux 下载页", URL: linuxURL}},
		})
	}
	if androidURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "arm64",
			Platform:     "Android",
			Links:        []plugin.Link{{Type: "direct", Label: "ToDesk Android 安装包", URL: androidURL}},
		})
	}
	if tvURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "arm64",
			Platform:     "Android TV",
			Links:        []plugin.Link{{Type: "direct", Label: "ToDesk TV 安装包", URL: tvURL}},
		})
	}
	if iosURL != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "universal",
			Platform:     "iOS / iPadOS",
			Links:        []plugin.Link{{Type: "store", Label: "App Store", URL: iosURL}},
		})
	}

	if len(variants) == 0 {
		return nil, fmt.Errorf("no downloadable variants found")
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "todesk",
				Name:            "ToDesk",
				Icon:            todeskIconURL,
				Description:     "ToDesk 远程控制客户端，支持多端远程协作与设备管理。",
				Organization:    "ToDesk",
				OfficialWebsite: todeskOfficialWebsite,
				Tags:            []string{"远程控制", "远程办公", "协作"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: todeskDownloadPage,
					Variants:    variants,
				},
			},
		},
	}, nil
}

func fetchToDeskHTML() (string, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(todeskDownloadPage)
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

func extractVersion(raw string) string {
	m := todeskVersionPattern.FindStringSubmatch(strings.TrimSpace(raw))
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

func decodeNuxtString(v string) string {
	v = strings.ReplaceAll(v, `\u002F`, "/")
	v = strings.ReplaceAll(v, `\u0026`, "&")
	v = strings.ReplaceAll(v, `\/`, "/")
	return strings.TrimSpace(v)
}

func parseDotDate(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	v = strings.ReplaceAll(v, ".", "-")
	parts := strings.Split(v, "-")
	if len(parts) != 3 {
		return ""
	}
	if len(parts[1]) == 1 {
		parts[1] = "0" + parts[1]
	}
	if len(parts[2]) == 1 {
		parts[2] = "0" + parts[2]
	}
	return strings.Join(parts, "-")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
