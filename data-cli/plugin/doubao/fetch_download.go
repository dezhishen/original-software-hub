package doubao

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
	DownloadPageURL string
}

const downloadPageURL = "https://www.doubao.com/download"

func fetchDownloadInfo() (*downloadInfo, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(downloadPageURL)
	if err != nil {
		return nil, fmt.Errorf("fetch doubao download page: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read doubao response: %w", err)
	}
	htmlContent := string(body)

	info := &downloadInfo{DownloadPageURL: downloadPageURL}

	versionPattern := regexp.MustCompile(`ScmProductionVersion\s*:\s*["\']?([0-9.]+)`)
	if matches := versionPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.Version = matches[1]
	}
	if info.Version == "" {
		versionPattern2 := regexp.MustCompile(`version["\s:]+([0-9]+\.[0-9]+\.[0-9]+)`)
		if matches := versionPattern2.FindStringSubmatch(htmlContent); len(matches) > 1 {
			info.Version = matches[1]
		}
	}

	windowsPattern := regexp.MustCompile(`https?://[^\s"<>]+(?:doubao|byted)[^\s"<>]*\.exe`)
	if matches := windowsPattern.FindStringSubmatch(htmlContent); len(matches) > 0 {
		info.WindowsURL = matches[0]
	}

	macPattern := regexp.MustCompile(`https?://[^\s"<>]+(?:doubao|byted)[^\s"<>]*\.dmg`)
	if matches := macPattern.FindStringSubmatch(htmlContent); len(matches) > 0 {
		info.MacIntelURL = matches[0]
	}

	info.ReleaseDate = time.Now().UTC().Format("2006-01-02")
	if info.Version == "" {
		info.Version = "Latest"
	}
	return info, nil
}
