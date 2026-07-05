package notepadplusplus

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	notepadplusplusOfficialWebsite = "https://notepad-plus-plus.org"
	notepadplusplusIconURL         = "https://notepad-plus-plus.org/favicon.ico"
)

type Notepadplusplus struct{}

func init() {
	plugin.Register(&Notepadplusplus{})
}

func (p *Notepadplusplus) Name() string {
	return "notepad-plus-plus"
}

func (p *Notepadplusplus) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "notepad-plus-plus",
				Name:            "Notepad++",
				Description:     "开源免费的源代码编辑器，支持多种编程语言，提供丰富的插件扩展。",
				Organization:    "Notepad++ Team",
				OfficialWebsite: notepadplusplusOfficialWebsite,
				Icon:            notepadplusplusIconURL,
				Tags:            []string{"开发工具", "文本编辑器"},
				Categories:      []string{"development"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: notepadplusplusOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, notepadplusplusOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Notepad++ 下载页", URL: "https://notepad-plus-plus.org/downloads/"},
							},
						},
						{
							Architecture: "x86",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Notepad++ 下载页 (x86)", URL: "https://notepad-plus-plus.org/downloads/"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Notepadplusplus) Disabled() bool { return false }
