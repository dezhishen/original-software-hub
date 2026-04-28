package firefox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const versionsAPI = "https://product-details.mozilla.org/1.0/firefox_versions.json"

type versionsResp struct {
	LatestVersion   string `json:"LATEST_FIREFOX_VERSION"`
	LastReleaseDate string `json:"LAST_RELEASE_DATE"`
}

func fetchLatestStable() (string, string, string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(versionsAPI)
	if err != nil {
		return "", "", "", fmt.Errorf("http get versions: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	var fv versionsResp
	if err := json.NewDecoder(resp.Body).Decode(&fv); err != nil {
		return "", "", "", fmt.Errorf("decode versions: %w", err)
	}
	if fv.LatestVersion == "" {
		return "", "", "", fmt.Errorf("empty latest version")
	}
	officialURL := fmt.Sprintf("https://www.mozilla.org/en-US/firefox/%s/releasenotes/", fv.LatestVersion)
	return fv.LatestVersion, fv.LastReleaseDate, officialURL, nil
}
