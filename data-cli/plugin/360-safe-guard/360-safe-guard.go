package p360safeguard

import (
	"net/http"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	p360safeguardOfficialWebsite = "https://www.360.com"
	p360safeguardIconURL         = "https://www.360.com/favicon.ico"
	p360safeguardDownloadURL     = "https://dl.360.cn/setup.exe"
)

type P360safeguard struct{}

func init() {
	plugin.Register(&P360safeguard{})
}

func (p *P360safeguard) Name() string {
	return "360-safe-guard"
}

func (p *P360safeguard) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := detectP360safeguardReleaseDate()
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "360-safe-guard",
				Name:            "360 安全卫士",
				Description:     "360 出品的综合安全防护工具，提供病毒查杀、系统优化、软件管理等功能。",
				Organization:    "360",
				OfficialWebsite: p360safeguardOfficialWebsite,
				Icon:            p360safeguardIconURL,
				Tags:            []string{"安全防护"},
				Categories:      []string{"security"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: p360safeguardOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, p360safeguardOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "360 安全卫士安装包", URL: p360safeguardDownloadURL},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func detectP360safeguardReleaseDate() string {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodHead, p360safeguardDownloadURL, nil)
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

func (p *P360safeguard) Disabled() bool { return false }
