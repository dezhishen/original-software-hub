package p360browser

import (
	"net/http"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	p360browserOfficialWebsite = "https://browser.360.cn"
	p360browserIconURL         = "https://browser.360.cn/favicon.ico"
	p360browserDownloadURL     = "https://browser.360.cn/ee/360browser_ee.exe"
)

type P360browser struct{}

func init() {
	plugin.Register(&P360browser{})
}

func (p *P360browser) Name() string {
	return "360-browser"
}

func (p *P360browser) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := detectP360browserReleaseDate()
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "360-browser",
				Name:            "360 浏览器",
				Description:     "360 出品的高速、安全网页浏览器，支持 Chromium 内核。",
				Organization:    "360",
				OfficialWebsite: p360browserOfficialWebsite,
				Icon:            p360browserIconURL,
				Tags:            []string{"浏览器"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: p360browserOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, p360browserOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "360 浏览器安装包", URL: p360browserDownloadURL},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func detectP360browserReleaseDate() string {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodHead, p360browserDownloadURL, nil)
	if err != nil {
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	lastModified := strings.TrimSpace(resp.Header.Get("Last-Modified"))
	if lastModified == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC1123, lastModified)
	if err != nil {
		return ""
	}
	return parsed.UTC().Format("2006-01-02")
}

func (p *P360browser) Disabled() bool { return false }
