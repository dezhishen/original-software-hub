package foxmail

import (
	"fmt"
	"io"
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
	foxmailWebsiteURL   = "https://www.foxmail.com/"
	foxmailWindowsURL   = "https://www.foxmail.com/win/"
	foxmailMacURL       = "https://www.foxmail.com/mac/"
	foxmailWindowsIDURL = "https://www.foxmail.com/win/download"
	foxmailMacIDURL     = "https://www.foxmail.com/mac/download"
	foxmailIconURL      = "https://www.foxmail.com/favicon.ico"
)

var (
	foxmailWindowsVersionPattern = regexp.MustCompile(`最新版本\s*[：:]\s*([0-9]+(?:\.[0-9]+)*)\s*\((20\d{2}-\d{2}-\d{2})\)`)
	foxmailMacVersionPattern     = regexp.MustCompile(`Foxmail_for_Mac_([0-9]+(?:\.[0-9]+)*)\.dmg`)
)

type Foxmail struct{}

func init() {
	plugin.Register(&Foxmail{})
}

func (p *Foxmail) Name() string {
	return "foxmail"
}

func (p *Foxmail) Fetch() ([]plugin.SoftwareData, error) {
	versions, err := fetchFoxmailVersions()
	if err != nil {
		return nil, err
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "foxmail",
				Name:            "Foxmail",
				Description:     "腾讯推出的桌面邮件客户端。",
				Organization:    "Tencent",
				OfficialWebsite: foxmailWebsiteURL,
				Icon:            foxmailIconURL,
				Tags:            []string{"邮件"},
			},
			Versions: versions,
		},
	}, nil
}

func fetchFoxmailVersions() ([]plugin.Version, error) {
	client := &http.Client{Timeout: 20 * time.Second}

	windowsVersion, windowsDate, windowsDownload, err := fetchFoxmailWindows(client)
	if err != nil {
		return nil, fmt.Errorf("foxmail windows: %w", err)
	}

	macVersion, macDate, macDownload, err := fetchFoxmailMac(client)
	if err != nil {
		return nil, fmt.Errorf("foxmail mac: %w", err)
	}

	windowsRelease := plugin.Version{
		Version:     windowsVersion,
		ReleaseDate: windowsDate,
		OfficialURL: foxmailWindowsURL,
		Platforms: plugin.PlatformsFromVariants(windowsVersion, windowsDate, foxmailWindowsURL, []plugin.Variant{
			{
				Architecture: "x64/x86",
				Platform:     "Windows",
				Links: []plugin.Link{
					{Type: "direct", Label: fileNameFromURL(macOrWinFallback(windowsDownload, foxmailWindowsIDURL)), URL: macOrWinFallback(windowsDownload, foxmailWindowsIDURL)},
				},
			},
		}),
	}

	macRelease := plugin.Version{
		Version:     macVersion,
		ReleaseDate: macDate,
		OfficialURL: foxmailMacURL,
		Platforms: plugin.PlatformsFromVariants(macVersion, macDate, foxmailMacURL, []plugin.Variant{
			{
				Architecture: "universal",
				Platform:     "macOS",
				Links: []plugin.Link{
					{Type: "direct", Label: fileNameFromURL(macOrWinFallback(macDownload, foxmailMacIDURL)), URL: macOrWinFallback(macDownload, foxmailMacIDURL)},
				},
			},
		}),
	}

	return mergeVersionsAsTabbed([]plugin.Version{windowsRelease, macRelease}), nil
}

func fetchFoxmailWindows(client *http.Client) (version, releaseDate, finalURL string, err error) {
	_, htmlText, err := fetchHTML(client, foxmailWindowsURL)
	if err != nil {
		return "", "", "", err
	}

	if m := foxmailWindowsVersionPattern.FindStringSubmatch(htmlText); len(m) >= 3 {
		version, releaseDate = m[1], m[2]
	}
	if version == "" {
		version = "latest"
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	finalURL, err = resolveDownloadURL(client, foxmailWindowsIDURL)
	if err != nil {
		return "", "", "", err
	}
	return version, releaseDate, finalURL, nil
}

func fetchFoxmailMac(client *http.Client) (version, releaseDate, finalURL string, err error) {
	doc, _, err := fetchHTML(client, foxmailMacURL)
	if err != nil {
		return "", "", "", err
	}

	downloadNode, _ := htmlquery.Query(doc, `//a[@id='download']`)
	downloadURL := foxmailMacIDURL
	if downloadNode != nil {
		if href := strings.TrimSpace(htmlquery.SelectAttr(downloadNode, "href")); href != "" {
			downloadURL = resolveURL(foxmailMacURL, href)
		}
	}

	finalURL, err = resolveDownloadURL(client, downloadURL)
	if err != nil {
		return "", "", "", err
	}

	if m := foxmailMacVersionPattern.FindStringSubmatch(path.Base(mustParseURL(finalURL).Path)); len(m) >= 2 {
		version = m[1]
	}
	if version == "" {
		version = "latest"
	}

	releaseDate = detectLastModifiedDate(client, finalURL)
	if releaseDate == "" {
		if n, _ := htmlquery.Query(doc, `(//span[contains(@class,'changelog-date')])[1]`); n != nil {
			releaseDate = strings.TrimSpace(htmlquery.InnerText(n))
		}
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return version, releaseDate, finalURL, nil
}

func fetchHTML(client *http.Client, rawURL string) (*html.Node, string, error) {
	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, "", fmt.Errorf("http get %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("http get %s: status %d", rawURL, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read body %s: %w", rawURL, err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, "", fmt.Errorf("parse html %s: %w", rawURL, err)
	}

	return doc, string(body), nil
}

func resolveDownloadURL(client *http.Client, rawURL string) (string, error) {
	headURL, err := resolveRedirect(client, rawURL, http.MethodHead)
	if err == nil && headURL != "" {
		return headURL, nil
	}
	getURL, getErr := resolveRedirect(client, rawURL, http.MethodGet)
	if getErr == nil && getURL != "" {
		return getURL, nil
	}
	if err != nil {
		return "", err
	}
	return "", getErr
}

func resolveRedirect(client *http.Client, rawURL, method string) (string, error) {
	req, err := http.NewRequest(method, rawURL, nil)
	if err != nil {
		return "", err
	}

	redirectClient := &http.Client{
		Timeout: client.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := redirectClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location := strings.TrimSpace(resp.Header.Get("Location"))
		if location == "" {
			return "", fmt.Errorf("redirect location is empty")
		}
		return resolveURL(rawURL, location), nil
	}

	if resp.StatusCode == http.StatusOK {
		return rawURL, nil
	}

	return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
}

func detectLastModifiedDate(client *http.Client, rawURL string) string {
	req, err := http.NewRequest(http.MethodHead, rawURL, nil)
	if err != nil {
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	lm := strings.TrimSpace(resp.Header.Get("Last-Modified"))
	if lm == "" {
		return ""
	}
	t, err := time.Parse(time.RFC1123, lm)
	if err != nil {
		t, err = time.Parse(time.RFC1123Z, lm)
		if err != nil {
			return ""
		}
	}
	return t.UTC().Format("2006-01-02")
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
		OfficialURL: foxmailWebsiteURL,
		Platforms:   platforms,
	}}
}

func macOrWinFallback(v, fallback string) string {
	v = strings.TrimSpace(v)
	if v != "" {
		return v
	}
	return fallback
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
		return "下载链接"
	}
	name := path.Base(u.Path)
	if name == "." || name == "/" || name == "" {
		return "下载链接"
	}
	return name
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return &url.URL{}
	}
	return u
}

func (p *Foxmail) Disabled() bool { return false }
