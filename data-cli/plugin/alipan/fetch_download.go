package alipan

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

type downloadInfo struct {
	Version         string
	ReleaseDate     string
	WindowsURL      string
	MacIntelURL     string
	MacARM64URL     string
	AndroidURL      string
	IOSURL          string
	DownloadPageURL string
}

const downloadPageURL = "https://www.alipan.com/download"

func fetchDownloadInfo() (*downloadInfo, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(downloadPageURL)
	if err != nil {
		return nil, fmt.Errorf("fetch alipan download page: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read alipan response: %w", err)
	}
	htmlContent := string(body)

	info := &downloadInfo{DownloadPageURL: downloadPageURL}

	versionPattern := regexp.MustCompile(`latest_version\s*:\s*['\"]?(V[\d.]+)`)
	if matches := versionPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.Version = matches[1]
	}

	windowsPattern := regexp.MustCompile(`app_windows_download_link\s*:\s*['\"]([^'\"]+\.exe)['\"]`)
	if matches := windowsPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.WindowsURL = matches[1]
	}

	macIntelPattern := regexp.MustCompile(`app_mac_download_link\s*:\s*['\"]([^'\"]+\.dmg)['\"]`)
	if matches := macIntelPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.MacIntelURL = matches[1]
	}

	macARM64Pattern := regexp.MustCompile(`app_mac_arm64_download_link\s*:\s*['\"]([^'\"]+arm64[^'\"]+\.dmg)['\"]`)
	if matches := macARM64Pattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.MacARM64URL = matches[1]
	}

	androidPattern := regexp.MustCompile(`app_android_download_link\s*:\s*['\"]([^'\"]+\.apk)['\"]`)
	if matches := androidPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.AndroidURL = matches[1]
	}

	iosPattern := regexp.MustCompile(`app_ios_download_link\s*:\s*['\"]([^'\"]+)['\"]`)
	if matches := iosPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.IOSURL = matches[1]
	}

	info.ReleaseDate = time.Now().UTC().Format("2006-01-02")
	if info.Version == "" {
		info.Version = "Latest"
	}

	return info, nil
}
