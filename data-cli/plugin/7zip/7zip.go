package sevenzip

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
	sevenzipHomeURL = "https://www.7-zip.org/"
	sevenzipIconURL = "https://www.7-zip.org/7ziplogo.png"
)

var (
	sevenzipMainPattern = regexp.MustCompile(`Download\s+7-Zip\s+([0-9]+(?:\.[0-9]+)*)\s*\((20\d{2}-\d{2}-\d{2})\)`)
	sevenzipVerPattern  = regexp.MustCompile(`Download\s+7-Zip\s+([0-9]+(?:\.[0-9]+)*)`)
	sevenzipPathPattern = regexp.MustCompile(`/download/([0-9]+(?:\.[0-9]+)*)/`)
)

type Sevenzip struct{}

func init() {
	plugin.Register(&Sevenzip{})
}

func (p *Sevenzip) Name() string {
	return "7zip"
}

func (p *Sevenzip) Fetch() ([]plugin.SoftwareData, error) {
	version, releaseDate, variants, err := fetchSevenzipWindowsVariants()
	if err != nil {
		return nil, err
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "7zip",
				Name:            "7-Zip",
				Description:     "高压缩率的开源压缩与解压工具。",
				Organization:    "Igor Pavlov",
				OfficialWebsite: sevenzipHomeURL,
				Icon:            sevenzipIconURL,
				Tags:            []string{"压缩"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: sevenzipHomeURL,
					Platforms:   plugin.PlatformsFromVariants(version, releaseDate, sevenzipHomeURL, variants),
				},
			},
		},
	}, nil
}

func fetchSevenzipWindowsVariants() (version, releaseDate string, variants []plugin.Variant, err error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(sevenzipHomeURL)
	if err != nil {
		return "", "", nil, fmt.Errorf("7zip: http get %s: %w", sevenzipHomeURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", nil, fmt.Errorf("7zip: unexpected status %d", resp.StatusCode)
	}

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return "", "", nil, fmt.Errorf("7zip: parse html: %w", err)
	}

	headline, _ := htmlquery.Query(doc, `//p/b[contains(.,'Download 7-Zip')]`)
	if headline != nil {
		hText := normalizeSpace(htmlquery.InnerText(headline))
		if m := sevenzipMainPattern.FindStringSubmatch(hText); len(m) >= 3 {
			version, releaseDate = m[1], m[2]
		} else if m := sevenzipVerPattern.FindStringSubmatch(hText); len(m) >= 2 {
			version = m[1]
		}
	}

	rows, err := htmlquery.QueryAll(doc, `//table[.//th[contains(.,'Windows')]]//tr[td[@class='Item']]`)
	if err != nil {
		return "", "", nil, fmt.Errorf("7zip: query download rows: %w", err)
	}

	byArch := map[string][]plugin.Link{}
	for _, row := range rows {
		archText := normalizeSpace(nodeText(row, `./td[3]`))
		arch := detectSevenzipArchitecture(archText)
		if arch == "" {
			continue
		}

		linkNode, _ := htmlquery.Query(row, `./td[1]/a`)
		if linkNode == nil {
			continue
		}
		href := strings.TrimSpace(htmlquery.SelectAttr(linkNode, "href"))
		if href == "" {
			continue
		}
		fullURL := resolveURL(sevenzipHomeURL, href)
		label := fileNameFromURL(fullURL)
		if label == "" {
			label = "7-Zip 安装包"
		}
		byArch[arch] = appendUniqueLink(byArch[arch], plugin.Link{Type: "direct", Label: label, URL: fullURL})

		if version == "" {
			if m := sevenzipPathPattern.FindStringSubmatch(fullURL); len(m) >= 2 {
				version = m[1]
			}
		}
	}

	archOrder := []string{"x64", "x86", "arm64"}
	for _, arch := range archOrder {
		links := byArch[arch]
		if len(links) == 0 {
			continue
		}
		variants = append(variants, plugin.Variant{Architecture: arch, Platform: "Windows", Links: links})
	}

	if len(variants) == 0 {
		return "", "", nil, fmt.Errorf("7zip: no windows download links found")
	}
	if version == "" {
		version = "latest"
	}
	return version, releaseDate, variants, nil
}

func detectSevenzipArchitecture(s string) string {
	lower := strings.ToLower(s)
	switch {
	case strings.Contains(lower, "arm64"):
		return "arm64"
	case strings.Contains(lower, "x86") || strings.Contains(lower, "32-bit"):
		return "x86"
	case strings.Contains(lower, "x64") || strings.Contains(lower, "64-bit"):
		return "x64"
	default:
		return ""
	}
}

func nodeText(root *html.Node, xpath string) string {
	n, err := htmlquery.Query(root, xpath)
	if err != nil || n == nil {
		return ""
	}
	return htmlquery.InnerText(n)
}

func appendUniqueLink(links []plugin.Link, link plugin.Link) []plugin.Link {
	for _, existing := range links {
		if existing.URL == link.URL {
			return links
		}
	}
	return append(links, link)
}

func resolveURL(baseURL, ref string) string {
	u, err := url.Parse(strings.TrimSpace(ref))
	if err == nil && u.IsAbs() {
		return u.String()
	}
	b, err := url.Parse(baseURL)
	if err != nil {
		return strings.TrimSpace(ref)
	}
	if u == nil {
		u, _ = url.Parse(strings.TrimSpace(ref))
	}
	if u == nil {
		return strings.TrimSpace(ref)
	}
	return b.ResolveReference(u).String()
}

func fileNameFromURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return ""
	}
	name := path.Base(u.Path)
	if name == "." || name == "/" {
		return ""
	}
	return name
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func (p *Sevenzip) Disabled() bool { return false }
