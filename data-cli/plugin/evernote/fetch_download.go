package evernote

import (
	"fmt"
	"net/http"
	"time"
)

const (
	evernoteWinURL = "https://win.desktop.evernote.com/builds/Evernote-latest.exe"
	evernoteMacURL = "https://mac.desktop.evernote.com/builds/Evernote-latest.dmg"
)

// fetchLastModified performs an HTTP HEAD request and returns the Last-Modified
// date formatted as "YYYY-MM-DD". Falls back to today's date on any error.
func fetchLastModified(rawURL string) (string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
		// Do not follow redirects — we only need response headers from the first hop.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Head(rawURL)
	if err != nil {
		return "", fmt.Errorf("head %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	lm := resp.Header.Get("Last-Modified")
	if lm == "" {
		return "", fmt.Errorf("Last-Modified header absent for %s", rawURL)
	}
	t, err := http.ParseTime(lm)
	if err != nil {
		return "", fmt.Errorf("parse Last-Modified %q: %w", lm, err)
	}
	return t.UTC().Format("2006-01-02"), nil
}
