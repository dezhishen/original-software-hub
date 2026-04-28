package todesk

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	todeskOfficialWebsite = "https://www.todesk.com/"
	todeskDownloadPage    = "https://www.todesk.com/download.html"
	todeskLinuxPage       = "https://www.todesk.com/linux.html"
	todeskIconURL         = "https://www.todesk.com/favicon.ico"
)

var (
	todeskWindowsLinkPattern = regexp.MustCompile(`https://dl\.todesk\.com/irrigation/ToDesk_[0-9]+(?:\.[0-9]+)+\.exe`)
	todeskMacLinkPattern     = regexp.MustCompile(`https://dl\.todesk\.com/macos/ToDesk_[0-9]+(?:\.[0-9]+)+\.pkg`)
	todeskAndroidPattern     = regexp.MustCompile(`https://dl\.todesk\.com/android/ToDesk_[0-9]+(?:\.[0-9]+)+\.apk`)
	todeskTVPattern          = regexp.MustCompile(`https://dl\.todesk\.com/android/ToDesk_TV_[0-9]+(?:\.[0-9]+)+\.apk`)
	todeskIOSPattern         = regexp.MustCompile(`https://apps\.apple\.com/cn/app/todesk/id[0-9]+`)
	todeskLinuxPattern       = regexp.MustCompile(`https://dl\.todesk\.com/linux/[a-zA-Z0-9._-]+`)
	todeskVersionPattern     = regexp.MustCompile(`ToDesk_([0-9]+(?:\.[0-9]+)+)\.`)
	todeskTVVersionPattern   = regexp.MustCompile(`ToDesk_TV_([0-9]+(?:\.[0-9]+)+)\.apk`)
	todeskAndroidDatePattern = regexp.MustCompile(`android_release_date:"([0-9]+(?:\.[0-9]+){2})"`)
	todeskIOSVersionPattern  = regexp.MustCompile(`ios_version:"([0-9]+(?:\.[0-9]+)+)"`)
	todeskIOSDatePattern     = regexp.MustCompile(`ios_release_date:"([0-9]+(?:\.[0-9]+){2})"`)
	todeskLinuxDatePattern   = regexp.MustCompile(`linux_release_date:"([0-9]+(?:\.[0-9]+){2})"`)
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
	linuxURLs := extractLinuxDirectLinks(decodedHTML)

	versions := make([]plugin.Version, 0, 6)

	if windowsURL != "" {
		versions = append(versions, plugin.Version{
			Version:     normalizeVersion(extractVersion(windowsURL)),
			ReleaseDate: fallbackReleaseDate(""),
			OfficialURL: todeskDownloadPage,
			Variants: []plugin.Variant{{
				Architecture: "x64",
				Platform:     "Windows",
				Links:        []plugin.Link{{Type: "direct", Label: "ToDesk Windows 安装包", URL: windowsURL}},
			}},
		})
	}

	if macURL != "" {
		versions = append(versions, plugin.Version{
			Version:     normalizeVersion(extractVersion(macURL)),
			ReleaseDate: fallbackReleaseDate(""),
			OfficialURL: todeskDownloadPage,
			Variants: []plugin.Variant{{
				Architecture: "universal",
				Platform:     "macOS",
				Links:        []plugin.Link{{Type: "direct", Label: "ToDesk macOS 安装包", URL: macURL}},
			}},
		})
	}

	if len(linuxURLs) > 0 {
		linuxVariants := make([]plugin.Variant, 0, len(linuxURLs)+1)
		for _, u := range linuxURLs {
			linuxVariants = append(linuxVariants, buildLinuxVariant(u))
		}
		linuxVariants = append(linuxVariants, plugin.Variant{
			Architecture: "all",
			Platform:     "Linux",
			Links:        []plugin.Link{{Type: "webpage", Label: "ToDesk Linux 发行版下载页", URL: todeskLinuxPage}},
		})

		versions = append(versions, plugin.Version{
			Version:     normalizeVersion(latestVersionFromURLs(linuxURLs)),
			ReleaseDate: fallbackReleaseDate(parseDotDate(findMatch(html, todeskLinuxDatePattern))),
			OfficialURL: todeskLinuxPage,
			Variants:    linuxVariants,
		})
	}

	if androidURL != "" {
		versions = append(versions, plugin.Version{
			Version:     normalizeVersion(extractVersion(androidURL)),
			ReleaseDate: fallbackReleaseDate(parseDotDate(findMatch(html, todeskAndroidDatePattern))),
			OfficialURL: todeskDownloadPage,
			Variants: []plugin.Variant{{
				Architecture: "arm64",
				Platform:     "Android",
				Links:        []plugin.Link{{Type: "direct", Label: "ToDesk Android 安装包", URL: androidURL}},
			}},
		})
	}

	if tvURL != "" {
		versions = append(versions, plugin.Version{
			Version:     normalizeVersion(extractTVVersion(tvURL)),
			ReleaseDate: fallbackReleaseDate(""),
			OfficialURL: todeskDownloadPage,
			Variants: []plugin.Variant{{
				Architecture: "arm64",
				Platform:     "Android TV",
				Links:        []plugin.Link{{Type: "direct", Label: "ToDesk TV 安装包", URL: tvURL}},
			}},
		})
	}

	if iosURL != "" {
		versions = append(versions, plugin.Version{
			Version:     normalizeVersion(findMatch(html, todeskIOSVersionPattern)),
			ReleaseDate: fallbackReleaseDate(parseDotDate(findMatch(html, todeskIOSDatePattern))),
			OfficialURL: todeskDownloadPage,
			Variants: []plugin.Variant{{
				Architecture: "universal",
				Platform:     "iOS / iPadOS",
				Links:        []plugin.Link{{Type: "store", Label: "App Store", URL: iosURL}},
			}},
		})
	}

	if len(versions) == 0 {
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
			Versions: versions,
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

func extractTVVersion(raw string) string {
	m := todeskTVVersionPattern.FindStringSubmatch(strings.TrimSpace(raw))
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

func extractLinuxDirectLinks(decodedHTML string) []string {
	matches := todeskLinuxPattern.FindAllString(decodedHTML, -1)
	if len(matches) == 0 {
		return nil
	}
	uniq := make(map[string]struct{}, len(matches))
	for _, m := range matches {
		m = strings.TrimSpace(m)
		if m == "" {
			continue
		}
		uniq[m] = struct{}{}
	}
	out := make([]string, 0, len(uniq))
	for u := range uniq {
		out = append(out, u)
	}
	sort.Strings(out)
	return out
}

func buildLinuxVariant(u string) plugin.Variant {
	name := strings.ToLower(path.Base(u))
	arch := "x64"
	switch {
	case strings.Contains(name, "arm64") || strings.Contains(name, "armv7"):
		arch = "arm64"
	case strings.Contains(name, "amd64") || strings.Contains(name, "x86_64"):
		arch = "x64"
	}

	label := "ToDesk Linux 下载包"
	switch {
	case strings.Contains(name, "uos"):
		label = "ToDesk Linux UOS 包"
	case strings.Contains(name, "kylin"):
		label = "ToDesk Linux 麒麟包"
	case strings.Contains(name, "nfschina"):
		label = "ToDesk Linux 方德包"
	case strings.HasSuffix(name, ".rpm"):
		label = "ToDesk Linux RPM 包"
	case strings.HasSuffix(name, ".deb"):
		label = "ToDesk Linux DEB 包"
	case strings.HasSuffix(name, ".pkg.tar.zst"):
		label = "ToDesk Linux Arch 包"
	}

	return plugin.Variant{
		Architecture: arch,
		Platform:     "Linux",
		Links:        []plugin.Link{{Type: "direct", Label: label, URL: u}},
	}
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

func fallbackReleaseDate(v string) string {
	v = strings.TrimSpace(v)
	if v != "" {
		return v
	}
	return time.Now().UTC().Format("2006-01-02")
}

func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "Latest"
	}
	return v
}

func latestVersionFromURLs(urls []string) string {
	best := ""
	for _, u := range urls {
		v := extractVersionFromFileName(path.Base(u))
		if compareVersion(v, best) > 0 {
			best = v
		}
	}
	return best
}

func extractVersionFromFileName(name string) string {
	re := regexp.MustCompile(`([0-9]+(?:\.[0-9]+)+)`)
	m := re.FindStringSubmatch(strings.ToLower(strings.TrimSpace(name)))
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

func compareVersion(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return -1
	}
	if b == "" {
		return 1
	}
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		ai := 0
		bi := 0
		if i < len(as) {
			fmt.Sscanf(as[i], "%d", &ai)
		}
		if i < len(bs) {
			fmt.Sscanf(bs[i], "%d", &bi)
		}
		if ai > bi {
			return 1
		}
		if ai < bi {
			return -1
		}
	}
	return 0
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func (x *ToDesk) CompareWithPrevious(previous plugin.PreviousState) ([]plugin.FetchResult, error) {
	items, err := x.Fetch()
	if err != nil {
		return nil, err
	}
	return plugin.BuildCompareResults(items, previous), nil
}
