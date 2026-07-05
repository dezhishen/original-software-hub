package git

import (
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

const (
	gitOfficialWebsite = "https://git-scm.com"
	gitIconURL         = "https://git-scm.com/favicon.ico"
)

type Git struct{}

func init() {
	plugin.Register(&Git{})
}

func (p *Git) Name() string {
	return "git"
}

func (p *Git) Fetch() ([]plugin.SoftwareData, error) {
	releaseDate := time.Now().UTC().Format("2006-01-02")

	return []plugin.SoftwareData{
		{
			Item: plugin.SoftwareItem{
				ID:              "git",
				Name:            "Git",
				Description:     "Linus Torvalds 创建的开源分布式版本控制系统，是目前最流行的版本控制工具。",
				Organization:    "Git Community",
				OfficialWebsite: gitOfficialWebsite,
				Icon:            gitIconURL,
				Tags:            []string{"开发工具", "版本控制"},
				Categories:      []string{"development"},
			},
			Versions: []plugin.Version{
				{
					Version:     "latest",
					ReleaseDate: releaseDate,
					OfficialURL: gitOfficialWebsite,
					Platforms: plugin.PlatformsFromVariants("latest", releaseDate, gitOfficialWebsite, []plugin.Variant{
						{
							Architecture: "x64",
							Platform:     "Windows",
							Links: []plugin.Link{
								{Type: "direct", Label: "Git for Windows", URL: "https://git-scm.com/download/win"},
							},
						},
						{
							Architecture: "universal",
							Platform:     "macOS",
							Links: []plugin.Link{
								{Type: "direct", Label: "Git for macOS", URL: "https://git-scm.com/download/mac"},
							},
						},
						{
							Architecture: "x64",
							Platform:     "Linux",
							Links: []plugin.Link{
								{Type: "webpage", Label: "Git for Linux", URL: "https://git-scm.com/download/linux"},
							},
						},
					}),
				},
			},
		},
	}, nil
}

func (p *Git) Disabled() bool { return false }
