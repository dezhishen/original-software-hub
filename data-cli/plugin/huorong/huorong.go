package huorong

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	huorongOfficialWebsite = "https://www.huorong.cn/person"
	huorongIconURL         = "https://cdn-www.huorong.cn/Public/Uploads/uploadfile/images/20240301/b1icon13.svg"

	huorongWindowsURL = "https://downloads.huorong.cn/setup/HuorongSetup.exe"
	huorongMacOSURL   = "https://downloads.huorong.cn/setup/Huorong.dmg"
)

// Huorong implements plugin.Plugin for Huorong antivirus client.
type Huorong struct{}

func init() {
	plugin.Register(&Huorong{})
}

func (h *Huorong) Name() string {
	return "huorong"
}

func (h *Huorong) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate, err := detectReleaseDate([]string{huorongWindowsURL, huorongMacOSURL})
	if err != nil {
		releaseDate = time.Now().UTC().Format("2006-01-02")
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "huorong",
				Name:            "火绒",
				Icon:            huorongIconURL,
				Description:     "火绒安全软件，专业的个人电脑防护工具。",
				Organization:    "Huorong",
				OfficialWebsite: huorongOfficialWebsite,
				Tags:            []string{"安全防护", "杀毒软件"},
			},
			Versions: []plugin.Version{
				{
					Version:     "Latest",
					ReleaseDate: releaseDate,
					OfficialURL: huorongOfficialWebsite,
					Variants: []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "火绒安装包 (exe)", URL: huorongWindowsURL},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: "火绒安装包 (dmg)", URL: huorongMacOSURL},
							},
						},
					},
				},
			},
		},
	}, nil
}

func detectReleaseDate(urls []string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	for _, u := range urls {
		resp, err := client.Head(u)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
				if t, err := time.Parse(time.RFC1123, lastModified); err == nil {
					return t.UTC().Format("2006-01-02"), nil
				}
			}
		}
	}
	return time.Now().UTC().Format("2006-01-02"), fmt.Errorf("detect huorong release date")
}
