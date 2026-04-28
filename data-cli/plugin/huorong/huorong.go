package huorong

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/dezhishen/original-software-hub/data-cli/plugin"
	"golang.org/x/net/html"
)

const (
	huorongOfficialWebsite = "https://www.huorong.cn/person"
	huorongFallbackIconURL = "https://cdn-www.huorong.cn/Public/Uploads/uploadfile/images/20240301/b1icon13.svg"
)

var (
	huorongVersionPattern = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	huorongDatePattern    = regexp.MustCompile(`(20\d{2})\.(\d{2})\.(\d{2})`)
	huorongFileDatePath   = regexp.MustCompile(`/files/(20\d{2})(\d{2})(\d{2})/`)
)

// Huorong implements plugin.Plugin for Huorong antivirus client.
type Huorong struct{}

func init() {
	plugin.Register(&Huorong{})
}

func (h *Huorong) Name() string {
	return "huorong"
}

func (h *Huorong) Fetch() ([]plugin.SoftwareData, error) {
	version, releaseDate, iconURL, variants, err := fetchHuorongDownloadInfo()
	if err != nil {
		return nil, fmt.Errorf("huorong: %w", err)
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "huorong",
				Name:            "火绒",
				Icon:            iconURL,
				Description:     "火绒安全软件，专业的个人电脑防护工具。",
				Organization:    "Huorong",
				OfficialWebsite: huorongOfficialWebsite,
				Tags:            []string{"安全防护", "杀毒软件"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: huorongOfficialWebsite,
					Platforms:   plugin.PlatformsFromVariants(version, releaseDate, huorongOfficialWebsite, variants),
				},
			},
		},
	}, nil
}

func fetchHuorongDownloadInfo() (string, string, string, []plugin.Variant, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(huorongOfficialWebsite)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("fetch official page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", nil, fmt.Errorf("fetch official page: status %d", resp.StatusCode)
	}

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("parse official page: %w", err)
	}

	iconURL, err := findHuorongIconURL(doc)
	if err != nil {
		iconURL = huorongFallbackIconURL
	}

	manualURL, err := findHuorongManualURL(doc)
	if err != nil {
		return "", "", "", nil, err
	}

	downloads, err := findHuorongDownloadLinks(doc)
	if err != nil {
		return "", "", "", nil, err
	}

	variants := make([]plugin.Variant, 0, 3)
	version := "Latest"
	releaseDate := releaseDateFromManualURL(manualURL)

	platforms := []struct {
		arch  string
		key   string
		label string
	}{
		{arch: "x64", key: "x64UrlAll", label: "火绒安装包 (x64 exe)"},
		{arch: "x86", key: "x86UrlAll", label: "火绒安装包 (x86 exe)"},
		{arch: "arm64", key: "arm64UrlAll", label: "火绒安装包 (ARM64 exe)"},
	}

	for _, platform := range platforms {
		downloadURL, ok := downloads[platform.key]
		if !ok {
			continue
		}

		resolvedURL, resolvedVersion, resolvedDate, err := resolveHuorongDownload(client, downloadURL)
		if err != nil {
			resolvedURL = downloadURL
		} else {
			if version == "Latest" && resolvedVersion != "" {
				version = resolvedVersion
			}
			if resolvedDate != "" {
				releaseDate = resolvedDate
			}
		}

		variants = append(variants, plugin.Variant{
			Architecture: platform.arch,
			Platform:     "Windows",
			Links: []plugin.Link{
				{Type: "direct", Label: platform.label, URL: resolvedURL},
			},
		})
	}

	if len(variants) == 0 {
		return "", "", "", nil, fmt.Errorf("no huorong download links found on official page")
	}

	if version == "Latest" {
		if manualVersion := extractHuorongVersion(manualURL); manualVersion != "" {
			version = manualVersion
		}
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return version, releaseDate, iconURL, variants, nil
}

func findHuorongIconURL(doc *html.Node) (string, error) {
	nodes, err := htmlquery.QueryAll(doc, `//img[contains(@src, 'b1icon13.svg')]`)
	if err != nil {
		return "", fmt.Errorf("query icon url: %w", err)
	}
	for _, node := range nodes {
		src := strings.TrimSpace(htmlquery.SelectAttr(node, "src"))
		if src != "" {
			return absoluteHuorongURL(src), nil
		}
	}

	nodes, err = htmlquery.QueryAll(doc, `//img[contains(@src, 'icon') and contains(@src, '.svg')]`)
	if err != nil {
		return "", fmt.Errorf("query fallback icon url: %w", err)
	}
	for _, node := range nodes {
		src := strings.TrimSpace(htmlquery.SelectAttr(node, "src"))
		if src != "" {
			return absoluteHuorongURL(src), nil
		}
	}

	return "", fmt.Errorf("huorong icon url not found on official page")
}

func findHuorongManualURL(doc *html.Node) (string, error) {
	nodes, err := htmlquery.QueryAll(doc, `//a[contains(@href, '.pdf')]`)
	if err != nil {
		return "", fmt.Errorf("query manual url: %w", err)
	}
	for _, node := range nodes {
		href := strings.TrimSpace(htmlquery.SelectAttr(node, "href"))
		text := strings.TrimSpace(htmlquery.InnerText(node))
		if strings.Contains(text, "用户使用手册") || strings.Contains(href, "用户操作手册") {
			return absoluteHuorongURL(href), nil
		}
	}
	return "", fmt.Errorf("huorong manual url not found on official page")
}

func findHuorongDownloadLinks(doc *html.Node) (map[string]string, error) {
	nodes, err := htmlquery.QueryAll(doc, `//a[@data-url]`)
	if err != nil {
		return nil, fmt.Errorf("query download links: %w", err)
	}

	results := map[string]string{}
	for _, node := range nodes {
		rawURL := strings.TrimSpace(htmlquery.SelectAttr(node, "data-url"))
		if !strings.Contains(rawURL, "downloadHr60.php") {
			continue
		}

		resolvedURL := absoluteHuorongURL(rawURL)
		parsedURL, err := url.Parse(resolvedURL)
		if err != nil {
			continue
		}

		plat := strings.TrimSpace(parsedURL.Query().Get("plat"))
		if plat == "" {
			continue
		}
		results[plat] = resolvedURL
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("huorong download links not found on official page")
	}
	return results, nil
}

func resolveHuorongDownload(client *http.Client, downloadURL string) (string, string, string, error) {
	request, err := http.NewRequest(http.MethodHead, downloadURL, nil)
	if err != nil {
		return "", "", "", fmt.Errorf("build head request: %w", err)
	}

	redirectClient := &http.Client{
		Timeout: client.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := redirectClient.Do(request)
	if err != nil {
		return "", "", "", fmt.Errorf("head download url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently && resp.StatusCode != http.StatusSeeOther && resp.StatusCode != http.StatusTemporaryRedirect && resp.StatusCode != http.StatusPermanentRedirect {
		return downloadURL, "", "", nil
	}

	location := strings.TrimSpace(resp.Header.Get("Location"))
	if location == "" {
		return downloadURL, "", "", nil
	}

	version := extractHuorongVersion(location)
	date := ""
	if match := huorongDatePattern.FindStringSubmatch(location); len(match) == 4 {
		date = fmt.Sprintf("%s-%s-%s", match[1], match[2], match[3])
	}

	return location, version, date, nil
}

func extractHuorongVersion(raw string) string {
	return huorongVersionPattern.FindString(raw)
}

func releaseDateFromManualURL(manualURL string) string {
	match := huorongFileDatePath.FindStringSubmatch(manualURL)
	if len(match) != 4 {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", match[1], match[2], match[3])
}

func absoluteHuorongURL(raw string) string {
	raw = strings.TrimSpace(strings.ReplaceAll(raw, "&amp;", "&"))
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if parsed.IsAbs() {
		return parsed.String()
	}
	base, err := url.Parse(huorongOfficialWebsite)
	if err != nil {
		return raw
	}
	return base.ResolveReference(parsed).String()
}
