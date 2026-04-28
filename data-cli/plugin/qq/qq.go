package qq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const pcConfigURL = "https://cdn-go.cn/qq-web/im.qq.com_new/latest/rainbow/pcConfig.json"

// QQ implements plugin.Plugin for Tencent QQ.
type QQ struct{}

func init() {
	plugin.Register(&QQ{})
}

func (q *QQ) Name() string {
	return "qq"
}

func (q *QQ) Fetch() ([]plugin.SoftwareData, error) {
	cfg, err := fetchPCConfig()
	if err != nil {
		return nil, fmt.Errorf("fetch qq pc config: %w", err)
	}

	versions := make([]plugin.Version, 0, 3)
	if version := buildWindowsVersion(cfg.Windows); version != nil {
		versions = append(versions, *version)
	}
	if version := buildLinuxVersion(cfg.Linux); version != nil {
		versions = append(versions, *version)
	}
	if version := buildMacOSVersion(cfg.MacOS); version != nil {
		versions = append(versions, *version)
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("qq pc config has no versions")
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "qq",
				Name:            "QQ",
				Icon:            "https://static-res.qq.com/static-res/imqq/qq-logo.png",
				Description:     "腾讯即时通信软件。",
				Organization:    "Tencent",
				OfficialWebsite: "https://im.qq.com/",
				Tags:            []string{"即时通讯", "社交"},
			},
			Versions: versions,
		},
	}, nil
}

type pcConfig struct {
	Windows windowsConfig `json:"Windows"`
	Linux   linuxConfig   `json:"Linux"`
	MacOS   macOSConfig   `json:"macOS"`
}

type windowsConfig struct {
	Version          string `json:"version"`
	UpdateDate       string `json:"updateDate"`
	DownloadURL      string `json:"downloadUrl"`
	NTDownloadURL    string `json:"ntDownloadUrl"`
	NTDownloadX64URL string `json:"ntDownloadX64Url"`
	NTDownloadARMURL string `json:"ntDownloadARMUrl"`
}

type linuxConfig struct {
	Version              string            `json:"version"`
	UpdateDate           string            `json:"updateDate"`
	X64DownloadURL       map[string]string `json:"x64DownloadUrl"`
	ARMDownloadURL       map[string]string `json:"armDownloadUrl"`
	LoongarchDownloadURL string            `json:"loongarchDownloadUrl"`
	MIPSDownloadURL      string            `json:"mipsDownloadUrl"`
}

type macOSConfig struct {
	Version     string `json:"version"`
	UpdateDate  string `json:"updateDate"`
	DownloadURL string `json:"downloadUrl"`
}

func fetchPCConfig() (*pcConfig, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(pcConfigURL)
	if err != nil {
		return nil, fmt.Errorf("http get %s: %w", pcConfigURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var cfg pcConfig
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode pc config: %w", err)
	}
	return &cfg, nil
}

func buildWindowsVersion(win windowsConfig) *plugin.Version {
	version := strings.TrimSpace(win.Version)
	if version == "" {
		return nil
	}

	variants := []plugin.Variant{
		buildVariant("x64", "Windows", []string{win.NTDownloadX64URL}),
		buildVariant("x86", "Windows", []string{win.NTDownloadURL, win.DownloadURL}),
		buildVariant("arm64", "Windows", []string{win.NTDownloadARMURL}),
	}
	variants = compactVariants(variants)
	if len(variants) == 0 {
		return nil
	}

	return &plugin.Version{
		Version:     "Windows " + version,
		ReleaseDate: strings.TrimSpace(win.UpdateDate),
		OfficialURL: "https://im.qq.com/pcqq/index.shtml",
		Variants:    variants,
	}
}

func buildLinuxVersion(linux linuxConfig) *plugin.Version {
	version := strings.TrimSpace(linux.Version)
	if version == "" {
		return nil
	}

	variants := []plugin.Variant{
		buildVariantFromMap("x64", "Linux", linux.X64DownloadURL),
		buildVariantFromMap("arm64", "Linux", linux.ARMDownloadURL),
		buildVariant("loongarch64", "Linux", []string{linux.LoongarchDownloadURL}),
		buildVariant("mips64el", "Linux", []string{linux.MIPSDownloadURL}),
	}
	variants = compactVariants(variants)
	if len(variants) == 0 {
		return nil
	}

	return &plugin.Version{
		Version:     "Linux " + version,
		ReleaseDate: strings.TrimSpace(linux.UpdateDate),
		OfficialURL: "https://im.qq.com/linuxqq/index.shtml",
		Variants:    variants,
	}
}

func buildMacOSVersion(mac macOSConfig) *plugin.Version {
	version := strings.TrimSpace(mac.Version)
	if version == "" {
		return nil
	}

	variants := compactVariants([]plugin.Variant{
		buildVariant("universal", "macOS", []string{mac.DownloadURL}),
	})
	if len(variants) == 0 {
		return nil
	}

	return &plugin.Version{
		Version:     "macOS " + version,
		ReleaseDate: strings.TrimSpace(mac.UpdateDate),
		OfficialURL: "https://im.qq.com/macqq/index.shtml",
		Variants:    variants,
	}
}

func buildVariant(arch, platform string, urls []string) plugin.Variant {
	seen := map[string]struct{}{}
	links := make([]plugin.Link, 0, len(urls))

	for _, raw := range urls {
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

	return plugin.Variant{
		Architecture: arch,
		Platform:     platform,
		Links:        links,
	}
}

func buildVariantFromMap(arch, platform string, packageURLs map[string]string) plugin.Variant {
	if len(packageURLs) == 0 {
		return plugin.Variant{Architecture: arch, Platform: platform}
	}

	keys := make([]string, 0, len(packageURLs))
	for packageType := range packageURLs {
		keys = append(keys, packageType)
	}
	sort.Strings(keys)

	seen := map[string]struct{}{}
	links := make([]plugin.Link, 0, len(keys))
	for _, packageType := range keys {
		u := strings.TrimSpace(packageURLs[packageType])
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		links = append(links, plugin.Link{
			Type:  "direct",
			Label: strings.ToLower(strings.TrimSpace(packageType)),
			URL:   u,
		})
	}

	return plugin.Variant{
		Architecture: arch,
		Platform:     platform,
		Links:        links,
	}
}

func compactVariants(variants []plugin.Variant) []plugin.Variant {
	out := make([]plugin.Variant, 0, len(variants))
	for _, variant := range variants {
		if len(variant.Links) == 0 {
			continue
		}
		out = append(out, variant)
	}
	return out
}

func fileNameFromURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return strings.TrimSpace(raw)
	}
	name := path.Base(parsed.Path)
	name = strings.TrimSpace(name)
	if name == "" || name == "." || name == "/" {
		return strings.TrimSpace(raw)
	}
	return name
}

func (x *QQ) FetchWithPrevious(previous plugin.PreviousState) ([]plugin.FetchResult, error) {
	items, err := x.Fetch()
	if err != nil {
		return nil, err
	}
	return plugin.BuildFetchResults(items, previous), nil
}
