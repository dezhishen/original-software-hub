package weixin

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
	weixinOfficialWebsite = "https://weixin.qq.com/"
	weixinUpdatesURL      = "https://weixin.qq.com/updates"
	weixinIconURL         = "https://sgnewres.wechat.com/t/ofs-wechat/newsroom-web/res/_next/static/media/wechat-with-color.7a890de8.svg"
)

var (
	reNuxtData = regexp.MustCompile(`(?s)<script[^>]+id="__NUXT_DATA__"[^>]*>(.*?)</script>`)
	reVersion  = regexp.MustCompile(`_(\d+\.\d+[\d.]*)\.(?:exe|dmg)`)
)

// WeChat implements plugin.Plugin for Tencent WeChat messenger.
type WeChat struct{}

func init() {
	plugin.Register(&WeChat{})
}

func (w *WeChat) Name() string {
	return "weixin"
}

func (w *WeChat) Fetch() ([]plugin.SoftwareData, error) {
	info, err := fetchUpdatesInfo()
	if err != nil {
		return nil, fmt.Errorf("fetch weixin updates: %w", err)
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "weixin",
				Name:            "微信",
				Icon:            weixinIconURL,
				Description:     "腾讯微信，提供消息、通话、朋友圈等功能的通讯软件。",
				Organization:    "Tencent",
				OfficialWebsite: weixinOfficialWebsite,
				Tags:            []string{"即时通讯", "社交"},
			},
			Versions: []plugin.Version{
				{
					Version:     info.version,
					ReleaseDate: info.releaseDate,
					OfficialURL: weixinUpdatesURL,
					Platforms:   plugin.PlatformsFromVariants(info.version, info.releaseDate, weixinUpdatesURL, info.variants),
				},
			},
		},
	}, nil
}

type updatesInfo struct {
	version     string
	releaseDate string
	variants    []plugin.Variant
}

func fetchUpdatesInfo() (*updatesInfo, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(weixinUpdatesURL)
	if err != nil {
		return nil, fmt.Errorf("http get %s: %w", weixinUpdatesURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	m := reNuxtData.FindSubmatch(body)
	if len(m) < 2 {
		return nil, fmt.Errorf("__NUXT_DATA__ script not found")
	}

	var flat []json.RawMessage
	if err := json.Unmarshal(m[1], &flat); err != nil {
		return nil, fmt.Errorf("parse __NUXT_DATA__: %w", err)
	}

	// str resolves a flat-array index to a string.
	str := func(idx int) string {
		if idx < 0 || idx >= len(flat) {
			return ""
		}
		var s string
		if err := json.Unmarshal(flat[idx], &s); err != nil {
			return ""
		}
		return s
	}

	intVal := func(raw json.RawMessage) int {
		var n int
		if json.Unmarshal(raw, &n) == nil {
			return n
		}
		return -1
	}

	// Find the downloadConf object: a dict with both "windows" and "winx86" keys.
	var conf rawConf
	confFound := false
	for i := range flat {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(flat[i], &obj); err != nil {
			continue
		}
		if _, ok := obj["windows"]; !ok {
			continue
		}
		if _, ok := obj["winx86"]; !ok {
			continue
		}
		conf = rawConf{
			windows: str(intVal(obj["windows"])),
			mac:     str(intVal(obj["mac"])),
			winx86:  str(intVal(obj["winx86"])),
			linux:   str(intVal(obj["linux"])),
			win10:   str(intVal(obj["win10"])),
			web:     str(intVal(obj["web"])),
		}
		confFound = true
		break
	}
	if !confFound {
		return nil, fmt.Errorf("downloadConf not found in __NUXT_DATA__")
	}

	// Extract version from the Windows x64 URL filename (e.g. WeChatWin_4.1.9.exe → 4.1.9).
	version := versionFromURL(conf.windows)
	if version == "" {
		version = versionFromURL(conf.mac)
	}
	if version == "" {
		return nil, fmt.Errorf("cannot determine version from download URLs")
	}

	// Find the latest Windows release date from the updates records.
	releaseDate := latestWindowsDate(flat, str)
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return &updatesInfo{
		version:     version,
		releaseDate: releaseDate,
		variants:    buildVariants(version, conf),
	}, nil
}

func versionFromURL(u string) string {
	if m := reVersion.FindStringSubmatch(u); len(m) >= 2 {
		return m[1]
	}
	return ""
}

// latestWindowsDate finds the publishDate of the first (latest) windows version record.
func latestWindowsDate(flat []json.RawMessage, str func(int) string) string {
	for i := range flat {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(flat[i], &obj); err != nil {
			continue
		}
		platRaw, ok1 := obj["platform"]
		verRaw, ok2 := obj["version"]
		dateRaw, ok3 := obj["publishDate"]
		if !ok1 || !ok2 || !ok3 {
			continue
		}
		var platIdx, verIdx, dateIdx int
		if json.Unmarshal(platRaw, &platIdx) != nil {
			continue
		}
		if json.Unmarshal(verRaw, &verIdx) != nil {
			continue
		}
		if json.Unmarshal(dateRaw, &dateIdx) != nil {
			continue
		}
		if str(platIdx) != "windows" || str(verIdx) == "" {
			continue
		}
		if date := str(dateIdx); date != "" {
			return date
		}
	}
	return ""
}

type rawConf struct {
	windows, mac, winx86, linux, win10, web string
}

func buildVariants(version string, conf rawConf) []plugin.Variant {
	var variants []plugin.Variant

	if conf.windows != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows",
			Links: []plugin.Link{
				{Type: "direct", Label: fmt.Sprintf("微信 %s 安装包 (exe)", version), URL: conf.windows},
			},
		})
	}
	if conf.winx86 != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x86",
			Platform:     "Windows",
			Links: []plugin.Link{
				{Type: "direct", Label: fmt.Sprintf("微信 %s 安装包 x86 (exe)", version), URL: conf.winx86},
			},
		})
	}
	if conf.win10 != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows (Store)",
			Links: []plugin.Link{
				{Type: "store", Label: "Microsoft Store", URL: conf.win10},
			},
		})
	}
	if conf.mac != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "universal",
			Platform:     "macOS",
			Links: []plugin.Link{
				{Type: "direct", Label: fmt.Sprintf("微信 %s 安装包 (dmg)", version), URL: conf.mac},
			},
		})
	}
	if conf.linux != "" && !strings.HasPrefix(conf.linux, "https://weixin.qq.com") {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Linux",
			Links: []plugin.Link{
				{Type: "webpage", Label: "Linux 版微信", URL: conf.linux},
			},
		})
	}
	if conf.web != "" {
		variants = append(variants, plugin.Variant{
			Architecture: "通用",
			Platform:     "Web",
			Links: []plugin.Link{
				{Type: "webpage", Label: "网页版微信", URL: conf.web},
			},
		})
	}

	return variants
}
