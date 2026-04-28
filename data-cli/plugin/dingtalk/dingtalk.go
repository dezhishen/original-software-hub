package dingtalk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	dingtalkOfficialWebsite = "https://www.dingtalk.com/"
	dingtalkDownloadPage    = "https://www.dingtalk.com/download"
	dingtalkIconURL         = "https://gw.alicdn.com/imgextra/i3/O1CN01eMicSg1GVD4uXMWGv_!!6000000000627-73-tps-32-32.ico"
)

var dingtalkVersionPattern = regexp.MustCompile(`DingTalk_v([0-9]+(?:\.[0-9]+)+)`)

// DingTalk implements plugin.Plugin for Alibaba DingTalk client.
type DingTalk struct{}

func init() {
	plugin.Register(&DingTalk{})
}

func (d *DingTalk) Name() string {
	return "dingtalk"
}

func (d *DingTalk) Fetch() ([]plugin.SoftwareData, error) {
	meta, err := fetchDownloadMeta()
	if err != nil {
		return nil, fmt.Errorf("fetch dingtalk download meta: %w", err)
	}

	version := extractDingTalkVersion(meta.WinAccessibilityDownloadLink)
	if version == "" {
		version = strings.TrimSpace(meta.Version)
	}
	if version == "" {
		version = "Latest"
	}

	releaseDate := parseDingTalkReleaseDate(meta.WinAccessibilityDownloadLink, meta.Time)
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	macStoreLink, windowsStoreLink := extractAppMarketLinks(meta)
	variants := []plugin.Variant{}
	if meta.WinAccessibilityDownloadLink != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows",
			Links:        []plugin.Link{{Type: "direct", Label: "钉钉 Windows 安装包", URL: meta.WinAccessibilityDownloadLink}},
		})
	}
	if windowsStoreLink != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows (Store)",
			Links:        []plugin.Link{{Type: "store", Label: "Microsoft Store", URL: windowsStoreLink}},
		})
	}
	if macStoreLink != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "universal",
			Platform:     "macOS",
			Links:        []plugin.Link{{Type: "store", Label: "Mac App Store", URL: macStoreLink}},
		})
	}
	variants = append(variants, plugin.Variant{
		Architecture: "通用",
		Platform:     "Web",
		Links:        []plugin.Link{{Type: "webpage", Label: "钉钉下载页", URL: dingtalkDownloadPage}},
	})

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "dingtalk",
				Name:            "钉钉",
				Icon:            dingtalkIconURL,
				Description:     "阿里巴巴旗下企业协作与即时通讯应用。",
				Organization:    "Alibaba",
				OfficialWebsite: dingtalkOfficialWebsite,
				Tags:            []string{"办公协作", "即时通讯"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: dingtalkDownloadPage,
					Variants:    variants,
				},
			},
		},
	}, nil
}

type dingtalkPayload struct {
	Version                      string `json:"version"`
	Time                         int64  `json:"time"`
	WinAccessibilityDownloadLink string `json:"winAccessibilityDownloadLink"`
	AppMarketData                struct {
		List []struct {
			Key      string `json:"key"`
			JumpLink string `json:"jumpLink"`
		} `json:"list"`
	} `json:"appMarketData"`
}

func fetchDownloadMeta() (*dingtalkPayload, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(dingtalkDownloadPage)
	if err != nil {
		return nil, fmt.Errorf("http get %s: %w", dingtalkDownloadPage, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	rawJSON, err := extractDataJSON(string(body))
	if err != nil {
		return nil, err
	}

	var payload dingtalkPayload
	if err := json.Unmarshal([]byte(rawJSON), &payload); err != nil {
		return nil, fmt.Errorf("decode __DATA__: %w", err)
	}
	payload.WinAccessibilityDownloadLink = strings.TrimSpace(payload.WinAccessibilityDownloadLink)
	payload.Version = strings.TrimSpace(payload.Version)
	return &payload, nil
}

func extractDataJSON(html string) (string, error) {
	anchor := "window.__DATA__ ="
	idx := strings.Index(html, anchor)
	if idx < 0 {
		return "", fmt.Errorf("window.__DATA__ not found")
	}

	s := html[idx+len(anchor):]
	start := strings.Index(s, "{")
	if start < 0 {
		return "", fmt.Errorf("__DATA__ json start not found")
	}

	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(s); i++ {
		ch := s[i]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			continue
		}
		if ch == '{' {
			depth++
			continue
		}
		if ch == '}' {
			depth--
			if depth == 0 {
				return s[start : i+1], nil
			}
		}
	}

	return "", fmt.Errorf("__DATA__ json end not found")
}

func extractDingTalkVersion(winURL string) string {
	if m := dingtalkVersionPattern.FindStringSubmatch(strings.TrimSpace(winURL)); len(m) >= 2 {
		return m[1]
	}
	return ""
}

func parseDingTalkReleaseDate(winURL string, ts int64) string {
	winURL = strings.TrimSpace(winURL)
	if parts := strings.Split(winURL, "/"); len(parts) > 3 {
		for _, p := range parts {
			if len(p) == 12 && isDigits(p) {
				if t, err := time.Parse("200601021504", p); err == nil {
					return t.Format("2006-01-02")
				}
			}
		}
	}
	if ts > 0 {
		return time.UnixMilli(ts).UTC().Format("2006-01-02")
	}
	return ""
}

func extractAppMarketLinks(meta *dingtalkPayload) (macStoreLink, windowsStoreLink string) {
	for _, item := range meta.AppMarketData.List {
		key := strings.TrimSpace(item.Key)
		link := strings.TrimSpace(item.JumpLink)
		if link == "" {
			continue
		}
		switch key {
		case "MacApplicationMarket":
			macStoreLink = link
		case "MicrosoftApplicationMarket":
			windowsStoreLink = link
		}
	}
	return
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func (x *DingTalk) CompareWithPrevious(previous plugin.PreviousState) ([]plugin.FetchResult, error) {
	items, err := x.Fetch()
	if err != nil {
		return nil, err
	}
	return plugin.BuildCompareResults(items, previous), nil
}
