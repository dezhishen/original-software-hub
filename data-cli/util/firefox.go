package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const firefoxVersionsAPI = "https://product-details.mozilla.org/1.0/firefox_versions.json"

type firefoxVersions struct {
	LatestVersion   string `json:"LATEST_FIREFOX_VERSION"`
	LastReleaseDate string `json:"LAST_RELEASE_DATE"` // "YYYY-MM-DD"
}

// FetchFirefoxLatestStable returns (version, releaseDate, officialURL, error).
// Version and release date are fetched directly from Mozilla's product-details API.
func FetchFirefoxLatestStable() (string, string, string, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(firefoxVersionsAPI)
	if err != nil {
		return "", "", "", fmt.Errorf("http get versions: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	var fv firefoxVersions
	if err := json.NewDecoder(resp.Body).Decode(&fv); err != nil {
		return "", "", "", fmt.Errorf("decode versions: %w", err)
	}
	if fv.LatestVersion == "" {
		return "", "", "", fmt.Errorf("empty latest version")
	}

	officialURL := fmt.Sprintf("https://www.mozilla.org/en-US/firefox/%s/releasenotes/", fv.LatestVersion)
	return fv.LatestVersion, fv.LastReleaseDate, officialURL, nil
}
