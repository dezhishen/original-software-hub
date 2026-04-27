package wps

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	wpsHomepageURL = "https://www.wps.cn/"
	wpsWebsiteURL  = "https://www.wps.cn/"
	wpsLinuxURL    = "https://www.wps.cn/product/wpslinux"
	wpsIconURL     = "https://volcengine-kdocs-cache.wpscdn.cn/kdocs/img/logo.5c78b00f.svg"
)

var (
	reOfficialJS = regexp.MustCompile(`src="(//volcengine-kdocs-cache\.wpscdn\.cn/kdocs/js/official[^\" ]*\.js)"`)
	reWindowsEXE = regexp.MustCompile(`https://official-package\.wpscdn\.cn/wps/download/WPS_Setup_(\d+)\.exe`)
	reMacDMG     = regexp.MustCompile(`https://package\.mac\.wpscdn\.cn/mac_wps_pkg/([0-9.]+)/[^\"' ]+\.dmg`)
	reLinuxMeta  = regexp.MustCompile(`版本\s*([0-9.]+)\s*([0-9]{4}\.[0-9]{2}\.[0-9]{2})`)
)

// WPS implements plugin.Plugin for WPS Office.
type WPS struct{}

func init() {
	plugin.Register(&WPS{})
}

func (w *WPS) Name() string {
	return "wps"
}

func (w *WPS) Fetch() ([]plugin.SoftwareData, error) {
	bundleURL, err := fetchOfficialBundleURL()
	if err != nil {
		return nil, fmt.Errorf("resolve wps official bundle: %w", err)
	}

	bundleText, err := fetchText(bundleURL)
	if err != nil {
		return nil, fmt.Errorf("fetch wps official bundle: %w", err)
	}

	linuxPageText, err := fetchText(wpsLinuxURL)
	if err != nil {
		return nil, fmt.Errorf("fetch wps linux page: %w", err)
	}

	versions, err := buildVersions(bundleText, linuxPageText)
	if err != nil {
		return nil, err
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "wps",
				Name:            "WPS Office",
				Icon:            wpsIconURL,
				Description:     "金山办公旗下 WPS Office 办公软件。",
				Organization:    "Kingsoft Office",
				OfficialWebsite: wpsWebsiteURL,
				Tags:            []string{"办公软件", "文档处理"},
			},
			Versions: versions,
		},
	}, nil
}

func fetchOfficialBundleURL() (string, error) {
	html, err := fetchText(wpsHomepageURL)
	if err != nil {
		return "", err
	}

	m := reOfficialJS.FindStringSubmatch(html)
	if len(m) < 2 {
		return "", fmt.Errorf("official js bundle not found")
	}

	bundle := strings.TrimSpace(m[1])
	if strings.HasPrefix(bundle, "//") {
		return "https:" + bundle, nil
	}
	if strings.HasPrefix(bundle, "http://") || strings.HasPrefix(bundle, "https://") {
		return bundle, nil
	}
	return "", fmt.Errorf("unsupported bundle url: %s", bundle)
}

func buildVersions(bundleText, linuxPageText string) ([]plugin.Version, error) {
	versions := make([]plugin.Version, 0, 3)

	if v := buildWindowsVersion(bundleText); v != nil {
		versions = append(versions, *v)
	}
	if v := buildMacVersion(bundleText); v != nil {
		versions = append(versions, *v)
	}
	if v := buildLinuxVersion(linuxPageText); v != nil {
		versions = append(versions, *v)
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no wps download info extracted from official bundle")
	}
	return versions, nil
}

func buildWindowsVersion(bundleText string) *plugin.Version {
	matches := reWindowsEXE.FindAllStringSubmatch(bundleText, -1)
	if len(matches) == 0 {
		return nil
	}

	type candidate struct {
		url     string
		version int
	}
	cands := make([]candidate, 0, len(matches))
	seen := map[string]struct{}{}
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		url := strings.TrimSpace(m[0])
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}

		v, err := strconv.Atoi(strings.TrimSpace(m[1]))
		if err != nil {
			continue
		}
		cands = append(cands, candidate{url: url, version: v})
	}
	if len(cands) == 0 {
		return nil
	}

	sort.Slice(cands, func(i, j int) bool {
		return cands[i].version > cands[j].version
	})
	latest := cands[0]

	releaseDate := detectLastModifiedDate(latest.url)
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &plugin.Version{
		Version:     fmt.Sprintf("Windows %d", latest.version),
		ReleaseDate: releaseDate,
		OfficialURL: wpsHomepageURL,
		Variants: []plugin.Variant{
			{
				Architecture: "x64/x86",
				Platform:     "Windows",
				Links: []plugin.Link{
					{Type: "direct", Label: fmt.Sprintf("WPS Setup %d (exe)", latest.version), URL: latest.url},
				},
			},
		},
	}
}

func buildMacVersion(bundleText string) *plugin.Version {
	matches := reMacDMG.FindAllStringSubmatch(bundleText, -1)
	if len(matches) == 0 {
		return nil
	}

	type candidate struct {
		url     string
		version string
	}
	cands := make([]candidate, 0, len(matches))
	seen := map[string]struct{}{}
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		url := strings.TrimSpace(m[0])
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}

		cands = append(cands, candidate{
			url:     url,
			version: strings.TrimSpace(m[1]),
		})
	}
	if len(cands) == 0 {
		return nil
	}

	sort.Slice(cands, func(i, j int) bool {
		return compareSemver(cands[i].version, cands[j].version) > 0
	})
	latest := cands[0]

	releaseDate := detectLastModifiedDate(latest.url)
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &plugin.Version{
		Version:     fmt.Sprintf("macOS %s", latest.version),
		ReleaseDate: releaseDate,
		OfficialURL: wpsHomepageURL,
		Variants: []plugin.Variant{
			{
				Architecture: "x64",
				Platform:     "macOS",
				Links: []plugin.Link{
					{Type: "direct", Label: fmt.Sprintf("WPS Office %s (dmg)", latest.version), URL: latest.url},
				},
			},
		},
	}
}

func buildLinuxVersion(linuxPageText string) *plugin.Version {
	version, releaseDate := extractLinuxVersionMeta(linuxPageText)
	if version == "" {
		version = "latest"
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &plugin.Version{
		Version:     fmt.Sprintf("Linux %s", version),
		ReleaseDate: releaseDate,
		OfficialURL: wpsLinuxURL,
		Variants: []plugin.Variant{
			{
				Architecture: "x64",
				Platform:     "Linux",
				Links: []plugin.Link{
					{Type: "webpage", Label: "WPS Linux 详情页", URL: wpsLinuxURL},
				},
			},
		},
	}
}

func extractLinuxVersionMeta(pageText string) (string, string) {
	m := reLinuxMeta.FindStringSubmatch(pageText)
	if len(m) < 3 {
		return "", ""
	}
	return strings.TrimSpace(m[1]), strings.ReplaceAll(strings.TrimSpace(m[2]), ".", "-")
}

func fetchText(rawURL string) (string, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(rawURL)
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

func detectLastModifiedDate(rawURL string) string {
	client := &http.Client{Timeout: 12 * time.Second}
	req, err := http.NewRequest(http.MethodHead, rawURL, nil)
	if err != nil {
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	lastModified := strings.TrimSpace(resp.Header.Get("Last-Modified"))
	if lastModified == "" {
		return ""
	}
	t, err := time.Parse(time.RFC1123, lastModified)
	if err != nil {
		return ""
	}
	return t.UTC().Format("2006-01-02")
}

func compareSemver(a, b string) int {
	ap := parseSemverParts(a)
	bp := parseSemverParts(b)
	maxLen := len(ap)
	if len(bp) > maxLen {
		maxLen = len(bp)
	}
	for i := 0; i < maxLen; i++ {
		av, bv := 0, 0
		if i < len(ap) {
			av = ap[i]
		}
		if i < len(bp) {
			bv = bp[i]
		}
		if av > bv {
			return 1
		}
		if av < bv {
			return -1
		}
	}
	return 0
}

func parseSemverParts(v string) []int {
	parts := strings.Split(strings.TrimSpace(v), ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			out = append(out, 0)
			continue
		}
		out = append(out, n)
	}
	return out
}
