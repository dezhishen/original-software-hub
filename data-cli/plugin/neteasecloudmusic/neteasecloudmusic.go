package neteasecloudmusic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	neteaseDownloadPage = "https://music.163.com/download"
	neteaseWebsite      = "https://music.163.com/"
	neteaseIconURL      = "https://s1.music.126.net/style/favicon.ico"
)

var (
	neteaseMacVersionPattern   = regexp.MustCompile(`NeteaseMusic_([0-9]+(?:\.[0-9]+)*)_([0-9]+)_web\.dmg`)
	neteaseLinuxVersionPattern = regexp.MustCompile(`netease-cloud-music_([0-9]+(?:\.[0-9]+)*)_amd64_([a-z0-9_]+)_(20\d{6})\.deb`)
	neteasePCEndpointPattern   = regexp.MustCompile(`https://music\.163\.com/api/pc/package/download/latest`)
	neteaseMacEndpointPattern  = regexp.MustCompile(`https://music\.163\.com/api/osx/download/latest`)
	neteaseUWPPattern          = regexp.MustCompile(`https://www\.microsoft\.com/store/apps/9nblggh6g0jf`)
	neteaseLinuxLinkPattern    = regexp.MustCompile(`http://d1\.music\.126\.net/dmusic/netease-cloud-music_[^"'\s]+\.deb`)
)

// NeteaseCloudMusic implements plugin.Plugin for NetEase Cloud Music.
type NeteaseCloudMusic struct{}

func init() {
	plugin.Register(&NeteaseCloudMusic{})
}

func (n *NeteaseCloudMusic) Name() string {
	return "neteasecloudmusic"
}

func (n *NeteaseCloudMusic) Fetch() ([]plugin.SoftwareData, error) {
	versions, err := fetchNeteaseVersions()
	if err != nil {
		return nil, fmt.Errorf("netease cloud music: %w", err)
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "netease-music",
				Name:            "网易云音乐",
				Icon:            neteaseIconURL,
				Description:     "网易云音乐官方客户端，支持多端音乐播放。",
				Organization:    "网易公司",
				OfficialWebsite: neteaseWebsite,
				Tags:            []string{"音乐", "流媒体"},
			},
			Versions: versions,
		},
	}, nil
}

type neteasePCResponse struct {
	Code int `json:"code"`
	Data struct {
		AppVer      string `json:"appVer"`
		BuildVer    string `json:"buildVer"`
		DownloadURL string `json:"downloadUrl"`
	} `json:"data"`
}

func fetchNeteaseVersions() ([]plugin.Version, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(neteaseDownloadPage)
	if err != nil {
		return nil, fmt.Errorf("http get %s: %w", neteaseDownloadPage, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read html: %w", err)
	}
	htmlText := string(body)

	pcEndpoint := firstMatch(htmlText, neteasePCEndpointPattern)
	if pcEndpoint == "" {
		return nil, fmt.Errorf("pc download endpoint not found in download page")
	}
	macEndpoint := firstMatch(htmlText, neteaseMacEndpointPattern)
	if macEndpoint == "" {
		return nil, fmt.Errorf("mac download endpoint not found in download page")
	}

	pcVersion, pcDate, pcURL, err := fetchPCLatest(client, pcEndpoint)
	if err != nil {
		return nil, err
	}

	macVersion, macDate, macURL, err := fetchMacLatest(client, macEndpoint)
	if err != nil {
		return nil, err
	}

	linuxLinks := extractLinuxLinks(htmlText)
	uwpURL := firstMatch(htmlText, neteaseUWPPattern)

	versions := make([]plugin.Version, 0, 3)
	versions = append(versions, plugin.Version{
		Version:     "Windows " + pcVersion,
		ReleaseDate: pcDate,
		OfficialURL: neteaseDownloadPage,
		Variants:    buildWindowsVariants(pcURL, uwpURL),
	})
	versions = append(versions, plugin.Version{
		Version:     "macOS " + macVersion,
		ReleaseDate: macDate,
		OfficialURL: neteaseDownloadPage,
		Variants: []plugin.Variant{
			{
				Architecture: "universal",
				Platform:     "macOS",
				Links: []plugin.Link{
					{Type: "direct", Label: "网易云音乐 macOS 安装包 (dmg)", URL: macURL},
				},
			},
		},
	})

	if lv := buildLinuxVersion(linuxLinks); lv != nil {
		versions = append(versions, *lv)
	}

	return versions, nil
}

func fetchPCLatest(client *http.Client, endpoint string) (version, releaseDate, downloadURL string, err error) {
	resp, err := client.Get(endpoint)
	if err != nil {
		return "", "", "", fmt.Errorf("fetch pc latest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("pc latest status %d", resp.StatusCode)
	}

	var payload neteasePCResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", "", "", fmt.Errorf("decode pc latest json: %w", err)
	}

	downloadURL = strings.TrimSpace(payload.Data.DownloadURL)
	if downloadURL == "" {
		return "", "", "", fmt.Errorf("pc latest json missing downloadUrl")
	}
	appVer := strings.TrimSpace(payload.Data.AppVer)
	buildVer := strings.TrimSpace(payload.Data.BuildVer)
	if appVer == "" {
		appVer = "latest"
	}
	if buildVer != "" {
		version = fmt.Sprintf("%s (%s)", appVer, buildVer)
	} else {
		version = appVer
	}

	releaseDate = detectLastModifiedDate(client, downloadURL)
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}
	return version, releaseDate, downloadURL, nil
}

func fetchMacLatest(client *http.Client, endpoint string) (version, releaseDate, downloadURL string, err error) {
	downloadURL, err = resolveRedirectURL(client, endpoint)
	if err != nil {
		return "", "", "", fmt.Errorf("resolve mac latest: %w", err)
	}

	base := path.Base(mustParseURL(downloadURL).Path)
	if m := neteaseMacVersionPattern.FindStringSubmatch(base); len(m) >= 3 {
		version = fmt.Sprintf("%s (%s)", m[1], m[2])
	} else {
		version = "latest"
	}

	releaseDate = detectLastModifiedDate(client, downloadURL)
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}
	return version, releaseDate, downloadURL, nil
}

func buildWindowsVariants(pcURL, uwpURL string) []plugin.Variant {
	variants := []plugin.Variant{
		{
			Architecture: "x64/x86",
			Platform:     "Windows",
			Links: []plugin.Link{
				{Type: "direct", Label: "网易云音乐 Windows 安装包", URL: pcURL},
			},
		},
	}
	if strings.TrimSpace(uwpURL) != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows (Store)",
			Links: []plugin.Link{
				{Type: "store", Label: "Microsoft Store", URL: uwpURL},
			},
		})
	}
	return variants
}

func buildLinuxVersion(links []plugin.Link) *plugin.Version {
	if len(links) == 0 {
		return nil
	}

	version := "latest"
	releaseDate := ""
	for _, link := range links {
		base := path.Base(mustParseURL(link.URL).Path)
		if m := neteaseLinuxVersionPattern.FindStringSubmatch(base); len(m) >= 4 {
			if m[1] > version || version == "latest" {
				version = m[1]
			}
			if d := compactDateToISO(m[3]); d > releaseDate {
				releaseDate = d
			}
		}
	}
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &plugin.Version{
		Version:     "Linux " + version,
		ReleaseDate: releaseDate,
		OfficialURL: neteaseDownloadPage,
		Variants: []plugin.Variant{
			{
				Architecture: "x64",
				Platform:     "Linux",
				Links:        links,
			},
		},
	}
}

func extractLinuxLinks(htmlText string) []plugin.Link {
	matches := neteaseLinuxLinkPattern.FindAllString(htmlText, -1)
	if len(matches) == 0 {
		return nil
	}

	links := make([]plugin.Link, 0, len(matches))
	seen := map[string]struct{}{}
	for _, raw := range matches {
		u := strings.TrimSpace(raw)
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}

		links = append(links, plugin.Link{Type: "direct", Label: fileNameFromURL(u), URL: u})
	}
	sort.SliceStable(links, func(i, j int) bool { return links[i].Label < links[j].Label })
	return links
}

func firstMatch(s string, pattern *regexp.Regexp) string {
	if pattern == nil {
		return ""
	}
	return strings.TrimSpace(pattern.FindString(s))
}

func resolveRedirectURL(client *http.Client, rawURL string) (string, error) {
	req, err := http.NewRequest(http.MethodHead, rawURL, nil)
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
		if location != "" {
			return location, nil
		}
	}
	return rawURL, nil
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

func compactDateToISO(raw string) string {
	if len(raw) != 8 {
		return ""
	}
	return raw[0:4] + "-" + raw[4:6] + "-" + raw[6:8]
}

func fileNameFromURL(raw string) string {
	u := mustParseURL(raw)
	name := strings.TrimSpace(path.Base(u.Path))
	if name == "" || name == "." || name == "/" {
		return strings.TrimSpace(raw)
	}
	return name
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return &url.URL{Path: strings.TrimSpace(raw)}
	}
	return u
}

func (x *NeteaseCloudMusic) FetchWithPrevious(previous plugin.PreviousState) ([]plugin.FetchResult, error) {
	items, err := x.Fetch()
	if err != nil {
		return nil, err
	}
	return plugin.BuildFetchResults(items, previous), nil
}
