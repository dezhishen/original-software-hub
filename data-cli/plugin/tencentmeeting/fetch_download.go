package tencentmeeting

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const downloadAPI = "https://meeting.tencent.com/web-service/query-download-info"

type downloadItem struct {
	Channel  string `json:"channel"`
	Platform string `json:"platform"`
	URL      string `json:"url"`
	Version  string `json:"version"`
	SubDate  string `json:"sub-date"`
}

type downloadResp struct {
	Code     int            `json:"code"`
	InfoList []downloadItem `json:"info-list"`
}

func fetchDownloadInfo() ([]downloadItem, error) {
	downloadConfig := []map[string]any{
		{"package-type": "app", "channel": "0300000000", "platform": "mac", "arch": "x86_64"},
		{"package-type": "app", "channel": "0300000000", "platform": "mac", "arch": "arm64"},
		{"package-type": "app", "channel": "0300000000", "platform": "windows"},
		{"package-type": "app", "channel": "0300000000", "platform": "windows", "arch": "x86_64"},
		{"package-type": "app", "channel": "1410000001", "platform": "ios"},
		{"package-type": "app", "channel": "0300000000", "platform": "android"},
		{"package-type": "app", "channel": "0300000000", "platform": "linux", "arch": "x86_64", "decorators": []string{"deb"}},
		{"package-type": "app", "channel": "0300000000", "platform": "linux", "arch": "arm64", "decorators": []string{"deb"}},
		{"package-type": "app", "channel": "0300000000", "platform": "linux", "arch": "loongarch64", "decorators": []string{"deb"}},
	}

	qRaw, err := json.Marshal(downloadConfig)
	if err != nil {
		return nil, fmt.Errorf("marshal download config: %w", err)
	}

	nonce := randomNonce(16)
	apiURL := fmt.Sprintf("%s?q=%s&nonce=%s", downloadAPI, url.QueryEscape(string(qRaw)), nonce)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var payload downloadResp
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if payload.Code != 0 {
		return nil, fmt.Errorf("api code %d", payload.Code)
	}
	if len(payload.InfoList) == 0 {
		return nil, fmt.Errorf("empty info-list")
	}
	return payload.InfoList, nil
}

func randomNonce(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if n <= 0 {
		return "nonce"
	}
	b := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}
