package evernote

import (
	"fmt"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	evernoteOfficialWebsite = "https://evernote.com"
	evernoteDownloadPage    = "https://evernote.com/download"
	evernoteIconURL         = "https://evernote.com/favicon.ico"
	evernoteAndroidURL      = "https://play.google.com/store/apps/details?id=com.evernote"
	evernoteIOSURL          = "https://apps.apple.com/app/evernote-notes-organizer/id281796108"
)

// Evernote implements plugin.Plugin for Evernote / 印象笔记.
type Evernote struct{}

func init() {
	plugin.Register(&Evernote{})
}

func (p *Evernote) Name() string {
	return "evernote"
}

func (p *Evernote) Fetch() ([]plugin.SoftwareData, error) {
	// Use the Windows package's Last-Modified date as the canonical release date.
	releaseDate, err := fetchLastModified(evernoteWinURL)
	if err != nil {
		// Graceful fallback: use today.
		releaseDate = time.Now().UTC().Format("2006-01-02")
		fmt.Printf("[evernote] warn: %v, using today as release date\n", err)
	}

	// Version is derived from release date since the download endpoint only
	// provides "latest" redirects with no embedded version number.
	version := releaseDate

	variants := []plugin.Variant{
		{
			Architecture: "x64",
			Platform:     "Windows",
			Links: []plugin.Link{
				{Type: "direct", Label: "Evernote Windows 安装包", URL: evernoteWinURL},
			},
		},
		{
			Architecture: "universal",
			Platform:     "macOS",
			Links: []plugin.Link{
				{Type: "direct", Label: "Evernote macOS 安装包", URL: evernoteMacURL},
			},
		},
		{
			Architecture: "all",
			Platform:     "Android",
			Links: []plugin.Link{
				{Type: "store", Label: "Google Play", URL: evernoteAndroidURL},
			},
		},
		{
			Architecture: "all",
			Platform:     "iOS / iPadOS",
			Links: []plugin.Link{
				{Type: "store", Label: "App Store", URL: evernoteIOSURL},
			},
		},
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "evernote",
				Name:            "Evernote",
				Description:     "跨平台云端笔记软件，支持记录、整理与同步。",
				Organization:    "Evernote Corporation",
				OfficialWebsite: evernoteOfficialWebsite,
				Icon:            evernoteIconURL,
				Tags:            []string{"笔记", "效率工具"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: evernoteDownloadPage,
					Platforms:   plugin.PlatformsFromVariants(version, releaseDate, evernoteDownloadPage, variants),
				},
			},
		},
	}, nil
}

func (p *Evernote) Disabled() bool { return false }
