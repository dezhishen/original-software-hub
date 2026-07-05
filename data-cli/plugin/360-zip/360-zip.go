package p360zip

import (
	"net/http"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	p360zipOfficialWebsite = "https://yasuo.360.cn"

	p360zipDownloadURL = "https://yasuo.360.cn/360zip_setup.exe"
)

type P360zip struct{}

func init() {
	plugin.Register(&P360zip{})
}

func (p *P360zip) Name() string {
	return "360-zip"
}

func (p *P360zip) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := detectP360zipReleaseDate()
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "360-zip",
				Name:            "360 压缩",
				Description:     "360 出品的免费压缩与解压工具，支持多种格式。",
				Organization:    "360",
				OfficialWebsite: p360zipOfficialWebsite,
				Icon:            "",
				Tags:            []string{"压缩"},
				Categories:      []string{"compression"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: p360zipOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, p360zipOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "360 压缩安装包", URL: p360zipDownloadURL},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func detectP360zipReleaseDate() string {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodHead, p360zipDownloadURL, nil)
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

func (p *P360zip) Disabled() bool { return false }
