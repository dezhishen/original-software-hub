package p360antivirus

import (
	"net/http"
	"strings"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	p360antivirusOfficialWebsite = "https://sd.360.cn"
	p360antivirusIconURL         = "https://sd.360.cn/favicon.ico"
	p360antivirusDownloadURL     = "https://dl.360.cn/sd/360sd_inst.exe"
)

type P360antivirus struct{}

func init() {
	plugin.Register(&P360antivirus{})
}

func (p *P360antivirus) Name() string {
	return "360-antivirus"
}

func (p *P360antivirus) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := detectP360antivirusReleaseDate()
	if releaseDate == "" {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "360-antivirus",
				Name:            "360 杀毒",
				Description:     "360 出品的安全与防护工具，提供病毒查杀、系统清理等功能。",
				Organization:    "360",
				OfficialWebsite: p360antivirusOfficialWebsite,
				Icon:            p360antivirusIconURL,
				Tags:            []string{"安全防护", "杀毒软件"},
				Categories:      []string{"security"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: p360antivirusOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, p360antivirusOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "360 杀毒安装包", URL: p360antivirusDownloadURL},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func detectP360antivirusReleaseDate() string {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodHead, p360antivirusDownloadURL, nil)
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

func (p *P360antivirus) Disabled() bool { return false }
