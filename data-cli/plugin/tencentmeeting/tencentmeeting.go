package tencentmeeting

import (
	"fmt"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
	"github.com/dezhishen/original-software-hub/data-cli/util"
)

const (
	tencentMeetingOfficialWebsite = "https://meeting.tencent.com/"
	tencentMeetingDownloadPage    = "https://meeting.tencent.com/download/"
	tencentMeetingIconURL         = "https://meeting.tencent.com/favicon.ico"
)

// TencentMeeting implements plugin.Plugin for Tencent Meeting.
type TencentMeeting struct{}

func init() {
	plugin.Register(&TencentMeeting{})
}

func (t *TencentMeeting) Name() string {
	return "tencentmeeting"
}

func (t *TencentMeeting) Fetch() ([]plugin.SoftwareData, error) {
	items, err := util.FetchTencentMeetingDownloadInfo()
	if err != nil {
		return nil, fmt.Errorf("fetch tencent meeting download info: %w", err)
	}

	byPlatform := map[string]util.TencentMeetingDownloadItem{}
	for _, it := range items {
		if it.Platform == "" || it.URL == "" {
			continue
		}
		byPlatform[it.Platform] = it
	}

	version := firstNonEmpty(
		byPlatform["windows_x86_64"].Version,
		byPlatform["windows"].Version,
		byPlatform["mac_arm64"].Version,
		byPlatform["mac"].Version,
	)
	if version == "" {
		version = "Latest"
	}

	releaseDate := firstNonEmpty(
		byPlatform["windows_x86_64"].SubDate,
		byPlatform["windows"].SubDate,
		byPlatform["mac_arm64"].SubDate,
		byPlatform["mac"].SubDate,
	)
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	variants := make([]plugin.Variant, 0, 8)
	if it, ok := byPlatform["windows_x86_64"]; ok {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Windows",
			Links: []plugin.Link{
				{Type: "direct", Label: "腾讯会议 x64 安装包 (exe)", URL: it.URL},
			},
		})
	}
	if it, ok := byPlatform["windows"]; ok {
		variants = append(variants, plugin.Variant{
			Architecture: "x86",
			Platform:     "Windows",
			Links: []plugin.Link{
				{Type: "direct", Label: "腾讯会议 安装包 (exe)", URL: it.URL},
			},
		})
	}
	if it, ok := byPlatform["mac"]; ok {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "macOS",
			Links: []plugin.Link{
				{Type: "direct", Label: "腾讯会议 macOS Intel (dmg)", URL: it.URL},
			},
		})
	}
	if it, ok := byPlatform["mac_arm64"]; ok {
		variants = append(variants, plugin.Variant{
			Architecture: "arm64",
			Platform:     "macOS",
			Links: []plugin.Link{
				{Type: "direct", Label: "腾讯会议 macOS Apple 芯片 (dmg)", URL: it.URL},
			},
		})
	}
	if it, ok := byPlatform["linux"]; ok {
		variants = append(variants, plugin.Variant{
			Architecture: "x64",
			Platform:     "Linux",
			Links: []plugin.Link{
				{Type: "direct", Label: "腾讯会议 Linux x64 (deb)", URL: it.URL},
			},
		})
	}
	if it, ok := byPlatform["linux_arm64"]; ok {
		variants = append(variants, plugin.Variant{
			Architecture: "arm64",
			Platform:     "Linux",
			Links: []plugin.Link{
				{Type: "direct", Label: "腾讯会议 Linux arm64 (deb)", URL: it.URL},
			},
		})
	}
	if it, ok := byPlatform["linux_deb_loongarch64"]; ok {
		variants = append(variants, plugin.Variant{
			Architecture: "loongarch64",
			Platform:     "Linux",
			Links: []plugin.Link{
				{Type: "direct", Label: "腾讯会议 Linux loongarch64 (deb)", URL: it.URL},
			},
		})
	}
	if it, ok := byPlatform["android"]; ok {
		variants = append(variants, plugin.Variant{
			Architecture: "arm64",
			Platform:     "Android",
			Links: []plugin.Link{
				{Type: "direct", Label: "腾讯会议 Android (apk)", URL: it.URL},
			},
		})
	}
	if it, ok := byPlatform["ios"]; ok {
		storeURL := normalizeIOSURL(it.URL)
		variants = append(variants, plugin.Variant{
			Architecture: "universal",
			Platform:     "iOS / iPadOS",
			Links: []plugin.Link{
				{Type: "store", Label: "App Store", URL: storeURL},
			},
		})
	}

	if len(variants) == 0 {
		return nil, fmt.Errorf("no download variants found")
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "tencent-meeting",
				Name:            "腾讯会议",
				Icon:            tencentMeetingIconURL,
				Description:     "腾讯会议，支持多端高清视频会议与在线协作。",
				Organization:    "Tencent",
				OfficialWebsite: tencentMeetingOfficialWebsite,
				Tags:            []string{"视频会议", "办公协作", "腾讯"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: tencentMeetingDownloadPage,
					Variants:    variants,
				},
			},
		},
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func normalizeIOSURL(u string) string {
	u = strings.TrimSpace(u)
	if strings.HasPrefix(u, "itms-apps://") {
		return "https://itunes.apple.com/cn/app/id1484048379"
	}
	if u == "" {
		return "https://itunes.apple.com/cn/app/id1484048379"
	}
	return u
}

func (x *TencentMeeting) FetchWithPrevious(previous plugin.PreviousState) ([]plugin.FetchResult, error) {
	items, err := x.Fetch()
	if err != nil {
		return nil, err
	}
	return plugin.BuildFetchResults(items, previous), nil
}
