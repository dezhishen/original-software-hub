package libreoffice

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	libreofficeHomeURL   = "https://zh-cn.libreoffice.org/"
	libreofficeMirrorURL = "https://mirrors.cloud.tencent.com/libreoffice/libreoffice/stable/"
	libreofficeIconURL   = "https://zh-cn.libreoffice.org/themes/libreofficenew/img/logo.png"
)

// versionDirPattern matches version directories like "26.2.4/" on the stable index page.
var versionDirPattern = regexp.MustCompile(`^(\d+\.\d+\.\d+)/$`)

type LibreOffice struct{}

func init() {
	plugin.Register(&LibreOffice{})
}

func (p *LibreOffice) Name() string {
	return "libreoffice"
}

func (p *LibreOffice) Disabled() bool { return false }

func (p *LibreOffice) Fetch() ([]plugin.SoftwareData, error) {
	version, releaseDate, err := fetchLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("libreoffice: %w", err)
	}

	variants := buildVariants(version)

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "libreoffice",
				Name:            "LibreOffice",
				Description:     "自由开源的办公套件，包含文字处理、电子表格、演示文稿等组件。",
				Organization:    "The Document Foundation",
				OfficialWebsite: libreofficeHomeURL,
				Icon:            libreofficeIconURL,
				Tags:            []string{"办公", "开源"},
				Categories:      []string{"office"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: libreofficeHomeURL,
					Platforms:   plugin.PlatformsFromVariants(version, releaseDate, libreofficeHomeURL, variants),
				},
			},
		},
	}, nil
}

// fetchLatestVersion parses the stable directory listing and returns the latest version number
// and the date of that version directory on the mirror.
func fetchLatestVersion() (version, releaseDate string, err error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(libreofficeMirrorURL)
	if err != nil {
		return "", "", fmt.Errorf("http get %s: %w", libreofficeMirrorURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("parse html: %w", err)
	}

	// Collect all version directory entries.
	type versionEntry struct {
		ver  string
		date string
	}
	var entries []versionEntry

	rows, err := htmlquery.QueryAll(doc, `//a[contains(@href,'/')]`)
	if err != nil {
		return "", "", fmt.Errorf("query version links: %w", err)
	}

	for _, row := range rows {
		href := strings.TrimSpace(htmlquery.SelectAttr(row, "href"))
		m := versionDirPattern.FindStringSubmatch(href)
		if m == nil {
			continue
		}
		ver := m[1]

		// Try to extract the date from the parent row text.
		// The Apache directory listing puts the date after the link.
		parent := row.Parent
		dateText := ""
		if parent != nil {
			dateText = strings.TrimSpace(htmlquery.InnerText(parent))
		}
		// dateText looks like "26.2.4/  04-Jun-2026 17:54  -"
		date := extractDate(dateText)
		entries = append(entries, versionEntry{ver: ver, date: date})
	}

	if len(entries) == 0 {
		return "", "", fmt.Errorf("no version directories found on %s", libreofficeMirrorURL)
	}

	// Sort by version (semver-like numeric comparison).
	sort.Slice(entries, func(i, j int) bool {
		return compareVersions(entries[i].ver, entries[j].ver) > 0
	})

	latest := entries[0]
	if latest.date == "" {
		latest.date = time.Now().UTC().Format("2006-01-02")
	}
	return latest.ver, latest.date, nil
}

// extractDate pulls a date like "04-Jun-2026" from a directory listing text.
func extractDate(s string) string {
	// Match patterns like "04-Jun-2026"
	re := regexp.MustCompile(`(\d{2}-[A-Z][a-z]{2}-\d{4})`)
	m := re.FindString(s)
	if m == "" {
		return ""
	}
	t, err := time.Parse("02-Jan-2006", m)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// compareVersions compares two dotted version strings numerically.
// Returns 1 if a > b, -1 if a < b, 0 if equal.
func compareVersions(a, b string) int {
	ap := parseVersion(a)
	bp := parseVersion(b)
	for i := 0; i < len(ap) && i < len(bp); i++ {
		if ap[i] > bp[i] {
			return 1
		}
		if ap[i] < bp[i] {
			return -1
		}
	}
	if len(ap) > len(bp) {
		return 1
	}
	if len(ap) < len(bp) {
		return -1
	}
	return 0
}

func parseVersion(v string) []int {
	parts := strings.Split(v, ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			n = 0
		}
		out = append(out, n)
	}
	return out
}

// baseMirrorURL returns the version-specific base URL on the mirror.
func baseMirrorURL(version string) string {
	return fmt.Sprintf("https://mirrors.cloud.tencent.com/libreoffice/libreoffice/stable/%s", version)
}

// buildVariants constructs the download variants for all supported platforms.
func buildVariants(version string) []plugin.Variant {
	base := baseMirrorURL(version)

	return []plugin.Variant{
		// ── Windows ──
		{
			Architecture: "x64",
			Platform:     "Windows",
			Links: []plugin.Link{
				{Type: "direct", Label: "msi 安装包 (64位)", URL: fmt.Sprintf("%s/win/x86_64/LibreOffice_%s_Win_x86-64.msi", base, version)},
			},
		},
		{
			Architecture: "x86",
			Platform:     "Windows",
			Links: []plugin.Link{
				{Type: "direct", Label: "msi 安装包 (32位)", URL: fmt.Sprintf("%s/win/x86/LibreOffice_%s_Win_x86.msi", base, version)},
			},
		},
		{
			Architecture: "arm64",
			Platform:     "Windows",
			Links: []plugin.Link{
				{Type: "direct", Label: "msi 安装包 (ARM64)", URL: fmt.Sprintf("%s/win/aarch64/LibreOffice_%s_Win_aarch64.msi", base, version)},
			},
		},
		// ── macOS ──
		{
			Architecture: "Intel",
			Platform:     "macOS",
			Links: []plugin.Link{
				{Type: "direct", Label: "dmg 安装包 (Intel)", URL: fmt.Sprintf("%s/mac/x86_64/LibreOffice_%s_MacOS_x86-64.dmg", base, version)},
			},
		},
		{
			Architecture: "Apple Silicon",
			Platform:     "macOS",
			Links: []plugin.Link{
				{Type: "direct", Label: "dmg 安装包 (Apple Silicon)", URL: fmt.Sprintf("%s/mac/aarch64/LibreOffice_%s_MacOS_aarch64.dmg", base, version)},
			},
		},
		// ── Linux ──
		{
			Architecture: "x64",
			Platform:     "Linux",
			Links: []plugin.Link{
				{Type: "direct", Label: "deb 包 (64位)", URL: fmt.Sprintf("%s/deb/x86_64/LibreOffice_%s_Linux_x86-64_deb.tar.gz", base, version)},
				{Type: "direct", Label: "rpm 包 (64位)", URL: fmt.Sprintf("%s/rpm/x86_64/LibreOffice_%s_Linux_x86-64_rpm.tar.gz", base, version)},
			},
		},
		{
			Architecture: "arm64",
			Platform:     "Linux",
			Links: []plugin.Link{
				{Type: "direct", Label: "deb 包 (ARM64)", URL: fmt.Sprintf("%s/deb/aarch64/LibreOffice_%s_Linux_aarch64_deb.tar.gz", base, version)},
				{Type: "direct", Label: "rpm 包 (ARM64)", URL: fmt.Sprintf("%s/rpm/aarch64/LibreOffice_%s_Linux_aarch64_rpm.tar.gz", base, version)},
			},
		},
	}
}
