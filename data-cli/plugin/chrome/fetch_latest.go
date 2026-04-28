package chrome

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const releasesAPI = "https://chromiumdash.appspot.com/fetch_releases?channel=Stable&platform=Windows&num=1&offset=0"

type release struct {
	Version string `json:"version"`
	Time    int64  `json:"time"`
}

func fetchLatestStable() (string, string, string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(releasesAPI)
	if err != nil {
		return "", "", "", fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var releases []release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", "", "", fmt.Errorf("decode: %w", err)
	}
	if len(releases) == 0 {
		return "", "", "", fmt.Errorf("empty release list")
	}

	r := releases[0]
	date := time.UnixMilli(r.Time).UTC().Format("2006-01-02")
	officialURL := "https://chromereleases.googleblog.com/"
	return r.Version, date, officialURL, nil
}
