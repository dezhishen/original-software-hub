package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const githubReleasesURL = "https://api.github.com/repos/%s/%s/releases/latest"

// GitHubRelease holds the fields we care about from the GitHub releases API.
type GitHubRelease struct {
	TagName     string        `json:"tag_name"`
	PublishedAt string        `json:"published_at"`
	HTMLURL     string        `json:"html_url"`
	Assets      []GitHubAsset `json:"assets"`
}

// GitHubAsset is a single release asset.
type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	ContentType        string `json:"content_type"`
}

// FetchGitHubLatestRelease fetches the latest release from a GitHub repo.
// owner and repo must be the GitHub user/org and repository name respectively.
func FetchGitHubLatestRelease(owner, repo string) (*GitHubRelease, error) {
	return FetchGitHubLatestReleaseWithToken(owner, repo, "")
}

// FetchGitHubLatestReleaseWithToken fetches the latest release using an optional token.
// When token is present, it is sent as a Bearer token to avoid low anonymous rate limits.
func FetchGitHubLatestReleaseWithToken(owner, repo, token string) (*GitHubRelease, error) {
	url := fmt.Sprintf(githubReleasesURL, owner, repo)
	client := &http.Client{Timeout: 20 * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if strings.TrimSpace(token) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api status %d for %s/%s", resp.StatusCode, owner, repo)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &release, nil
}

// FormatGitHubDate converts a GitHub ISO 8601 timestamp to YYYY-MM-DD.
func FormatGitHubDate(iso8601 string) string {
	t, err := time.Parse(time.RFC3339, iso8601)
	if err != nil {
		return iso8601
	}
	return t.Format("2006-01-02")
}
