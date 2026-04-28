package firefox

import (
	"fmt"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

func init() {
	plugin.Register(&Firefox{})
}

// Firefox implements plugin.Plugin for Mozilla Firefox.
type Firefox struct{}

func (f *Firefox) Name() string { return "firefox" }

func (f *Firefox) Fetch() ([]plugin.SoftwareData, error) {
	version, releaseDate, officialURL, err := fetchLatestStable()
	if err != nil {
		return nil, fmt.Errorf("firefox: %w", err)
	}

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "firefox",
				Name:            "Mozilla Firefox",
				Icon:            "https://www.mozilla.org/media/img/favicons/firefox/browser/favicon-196x196.59e3822720be.png",
				Description:     "Mozilla 出品的快速、安全、开源浏览器。",
				Organization:    "Mozilla Foundation",
				OfficialWebsite: "https://www.mozilla.org/firefox/",
				Tags:            []string{"浏览器", "网络", "开源"},
			},
			Versions: []plugin.Version{
				{
					Version:     version,
					ReleaseDate: releaseDate,
					OfficialURL: officialURL,
					Platforms: plugin.PlatformsFromVariants(version, releaseDate, officialURL, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "安装包 (exe)", URL: "https://download.mozilla.org/?product=firefox-latest-ssl&os=win64&lang=zh-CN"},
								{Type: "direct", Label: "安装包 (msi)", URL: "https://download.mozilla.org/?product=firefox-msi-latest-ssl&os=win64&lang=zh-CN"},
							},
						},
						{
							Architecture: "x86",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "安装包 (exe)", URL: "https://download.mozilla.org/?product=firefox-latest-ssl&os=win&lang=zh-CN"},
								{Type: "direct", Label: "安装包 (msi)", URL: "https://download.mozilla.org/?product=firefox-msi-latest-ssl&os=win&lang=zh-CN"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: "dmg 安装包", URL: "https://download.mozilla.org/?product=firefox-latest-ssl&os=osx&lang=zh-CN"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "direct", Label: "tar.bz2", URL: "https://download.mozilla.org/?product=firefox-latest-ssl&os=linux64&lang=zh-CN"},
							},
						},
					}),
				},
			},
		},
	}, nil
}
