package util

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// DoubaoDownloadInfo holds Doubao client download information
type DoubaoDownloadInfo struct {
	Version         string
	ReleaseDate     string
	WindowsURL      string
	MacIntelURL     string
	MacARM64URL     string
	OfficialURL     string
	DownloadPageURL string
}

const (
	doubaoDownloadURL = "https://www.doubao.com/download"
)

// FetchDoubaoDownloadInfo fetches Doubao client download information
func FetchDoubaoDownloadInfo() (*DoubaoDownloadInfo, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(doubaoDownloadURL)
	if err != nil {
		return nil, fmt.Errorf("fetch doubao download page: %w", err)
	}
	defer resp.Body.Close()

	// Read HTML body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read doubao response: %w", err)
	}
	htmlContent := string(body)

	info := &DoubaoDownloadInfo{
		OfficialURL:     "https://www.doubao.com",
		DownloadPageURL: doubaoDownloadURL,
	}

	// Extract version from buildEnv: ScmProductionVersion:"1.0.4.6288"
	versionPattern := regexp.MustCompile(`ScmProductionVersion\s*:\s*["\']?([0-9.]+)`)
	if matches := versionPattern.FindStringSubmatch(htmlContent); len(matches) > 1 {
		info.Version = matches[1]
	}

	// If buildEnv version not found, try other patterns
	if info.Version == "" {
		versionPattern2 := regexp.MustCompile(`version["\s:]+([0-9]+\.[0-9]+\.[0-9]+)`)
		if matches := versionPattern2.FindStringSubmatch(htmlContent); len(matches) > 1 {
			info.Version = matches[1]
		}
	}

	// Extract download links from JavaScript (may not be present in static HTML)
	// Try to find direct download links
	windowsPattern := regexp.MustCompile(`https?://[^\s"<>]+(?:doubao|byted)[^\s"<>]*\.exe`)
	if matches := windowsPattern.FindStringSubmatch(htmlContent); len(matches) > 0 {
		info.WindowsURL = matches[0]
	}

	macPattern := regexp.MustCompile(`https?://[^\s"<>]+(?:doubao|byted)[^\s"<>]*\.dmg`)
	if matches := macPattern.FindStringSubmatch(htmlContent); len(matches) > 0 {
		info.MacIntelURL = matches[0]
	}

	// If download links not found in HTML, construct them from known CDN
	// Doubao uses Bytedance CDN
	if info.WindowsURL == "" {
		// Placeholder - would need actual download link discovery
		info.WindowsURL = "" // Leave empty if not found
	}

	// Release date - use current date or extract from version
	info.ReleaseDate = time.Now().UTC().Format("2006-01-02")

	if info.Version == "" {
		info.Version = "Latest"
	}

	return info, nil
}
