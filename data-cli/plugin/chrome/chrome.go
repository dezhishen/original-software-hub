package chrome

import (
	"fmt"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
	"github.com/dezhishen/original-software-hub/data-cli/util"
)

func init() {
	plugin.Register(&Chrome{})
}

// Chrome implements plugin.Plugin for Google Chrome.
type Chrome struct{}

func (c *Chrome) Name() string { return "chrome" }

func (c *Chrome) Fetch() ([]plugin.SoftwareData, error) {
	// Chrome publishes direct download URLs at well-known paths.
	// For the stable channel, the version can be fetched from the
	// Chrome release endpoint; here we use a static known version
	// as a starting placeholder.
	version, releaseDate, officialURL, err := util.FetchChromeLatestStable()
	if err != nil {
		return nil, fmt.Errorf("chrome: %w", err)
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "chrome",
				Name:            "Google Chrome",
				Icon:            "https://www.google.com/chrome/static/images/chrome-logo.svg",
				Description:     "Google 出品的高速、安全浏览器。",
				Organization:    "Google LLC",
				OfficialWebsite: "https://www.google.com/chrome/", Tags: []string{"浏览器", "网络"}},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: officialURL,
					Variants: []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "在线安装包 (exe)", URL: "https://dl.google.com/chrome/install/ChromeStandaloneSetup64.exe"},
								{Type: "direct", Label: "企业 MSI (x64)", URL: "https://dl.google.com/tag/s/dl/chrome/install/googlechromestandaloneenterprise64.msi"},
							},
						},
						{
							Architecture: "x86",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "在线安装包 (exe)", URL: "https://dl.google.com/chrome/install/ChromeStandaloneSetup.exe"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: "dmg 安装包", URL: "https://dl.google.com/chrome/mac/stable/GGRO/googlechrome.dmg"},
							},
						},
					},
				},
			},
		},
	}, nil
}

func (x *Chrome) FetchWithPrevious(previous plugin.PreviousState) ([]plugin.FetchResult, error) {
	items, err := x.Fetch()
	if err != nil {
		return nil, err
	}
	return plugin.BuildFetchResults(items, previous), nil
}
