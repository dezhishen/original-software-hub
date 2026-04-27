package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// AlipanDownloadInfo holds Aliyun Drive desktop download information
type AlipanDownloadInfo struct {
	Version           string
	ReleaseDate       string
	WindowsURL        string
	MacIntelURL       string
	MacARM64URL       string
	LinuxURL          string
	AndroidURL        string
	IOSUrl            string
	OfficialURL       string
	DownloadPageURL   string
}

const (
	alipanDownloadURL = "https://www.alipan.com/download"
)

// FetchAlipanDownloadInfo fetches Aliyun Drive desktop client download information
func FetchAlipanDownloadInfo() (*AlipanDownloadInfo, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(alipanDownloadURL)
	if err != nil {
		return nil, fmt.Errorf("fetch alipan download page: %w", err)
	}
	defer resp.Body.Close()

	// Read HTML body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read alipan response: %w", err)
	}
	htmlContent := string(body)

	info := &AlipanDownloadInfo{
		OfficialURL:     "https://www.alipan.com",
		DownloadPageURL: alipanDownloadURL,
	}

	// Extract latest_version: 'V2.1.7'
	versionPattern := regexp.MustCompile(`latest_version\s*:\s*['\"]?(V[\d.]+)`)
	if matches := versionPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.Version = matches[1]
	}

	// Extract download links from JavaScript/HTML
	// app_windows_download_link: 'https://...'
	windowsPattern := regexp.MustCompile(`app_windows_download_link\s*:\s*['\"]([^'\"]+\.exe)['\"]`)
	if matches := windowsPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.WindowsURL = matches[1]
	}

	// app_mac_download_link: 'https://...'
	macIntelPattern := regexp.MustCompile(`app_mac_download_link\s*:\s*['\"]([^'\"]+\.dmg)['\"]`)
	if matches := macIntelPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.MacIntelURL = matches[1]
	}

	// app_mac_arm64_download_link: 'https://...'
	macARM64Pattern := regexp.MustCompile(`app_mac_arm64_download_link\s*:\s*['\"]([^'\"]+arm64[^'\"]+\.dmg)['\"]`)
	if matches := macARM64Pattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.MacARM64URL = matches[1]
	}

	// app_android_download_link
	androidPattern := regexp.MustCompile(`app_android_download_link\s*:\s*['\"]([^'\"]+\.apk)['\"]`)
	if matches := androidPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.AndroidURL = matches[1]
	}

	// app_ios_download_link
	iosPattern := regexp.MustCompile(`app_ios_download_link\s*:\s*['\"]([^'\"]+)['\"]`)
	if matches := iosPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.IOSUrl = matches[1]
	}

	// Extract release date from version info or use current date
	// Aliyun Drive doesn't provide explicit release date in HTML, use current UTC date
	info.ReleaseDate = time.Now().UTC().Format("2006-01-02")

	if info.Version == "" {
		info.Version = "Latest"
	}

	return info, nil
}

// DownloadInfoToJSON converts AlipanDownloadInfo to JSON string for debugging
func (a *AlipanDownloadInfo) ToJSON() string {
	data, _ := json.MarshalIndent(a, "", "  ")
	return string(data)
}
