package anydesk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

const (
	anydeskDownloadPage = "https://anydesk.com.cn/zhs/downloads/windows"
)

// downloadsEntry represents one platform block inside the page-embedded JS object.
type downloadsEntry struct {
	Version  string             `json:"version"`
	URL      string             `json:"url"`
	Packages []downloadsPackage `json:"packages"`
}

type downloadsPackage struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	Version string `json:"version"`
	URL     string `json:"url"`
}

type downloadsMap map[string]downloadsEntry

// varDownloadsRe matches the first occurrence of `var downloads={...};` in the page.
var varDownloadsRe = regexp.MustCompile(`(?s)var downloads=(\{.+?\});`)

func fetchAnyDeskDownloads() (downloadsMap, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest(http.MethodGet, anydeskDownloadPage, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	m := varDownloadsRe.FindSubmatch(body)
	if len(m) < 2 {
		return nil, fmt.Errorf("var downloads not found in page")
	}

	var dm downloadsMap
	if err := json.Unmarshal(m[1], &dm); err != nil {
		return nil, fmt.Errorf("unmarshal downloads JSON: %w", err)
	}
	return dm, nil
}
