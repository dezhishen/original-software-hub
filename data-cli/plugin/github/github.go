package github

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
	"github.com/dezhishen/original-software-hub/data-cli/util"
	"gopkg.in/yaml.v3"
)

const (
	envConfigInline = "OSH_GITHUB_CONFIG"
	envConfigPath   = "OSH_GITHUB_CONFIG_PATH"
	envToken        = "OSH_GITHUB_TOKEN"
)

type pluginConfig struct {
	Token string       `yaml:"token"`
	Repos []repoConfig `yaml:"repos"`
}

type repoConfig struct {
	ID              string      `yaml:"id"`
	Disabled        bool        `yaml:"disabled"`
	Owner           string      `yaml:"owner"`
	Repo            string      `yaml:"repo"`
	Name            string      `yaml:"name"`
	Description     string      `yaml:"description"`
	Organization    string      `yaml:"organization"`
	OfficialWebsite string      `yaml:"officialWebsite"`
	Icon            string      `yaml:"icon"`
	Tags            []string    `yaml:"tags"`
	Assets          []assetRule `yaml:"assets"`
}

type assetRule struct {
	// Keywords 中的每个词都必须出现在资源文件名或下载 URL 中（大小写不敏感，AND 逻辑）。
	// 好处：与具体文件命名无关，不同仓库使用相同声明即可匹配。
	Keywords     []string `yaml:"keywords"`
	Platform     string   `yaml:"platform"`
	Architecture string   `yaml:"architecture"`
	Package      string   `yaml:"package"`
}

type githubPlugin struct {
	repos []repoConfig
	token string
}

func init() {
	cfg, err := loadConfig()
	if err != nil {
		log.Printf("[github] config load error: %v", err)
		return
	}
	if cfg == nil || len(cfg.Repos) == 0 {
		return
	}

	token := strings.TrimSpace(cfg.Token)
	if envValue := strings.TrimSpace(os.Getenv(envToken)); envValue != "" {
		token = envValue
	} else if envValue := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); envValue != "" {
		token = envValue
	}

	repos := make([]repoConfig, 0, len(cfg.Repos))
	for _, repo := range cfg.Repos {
		if repo.Disabled {
			log.Printf("[github] skip disabled repo: %s/%s", repo.Owner, repo.Repo)
			continue
		}
		normalized, ok := normalizeRepoConfig(repo)
		if !ok {
			log.Printf("[github] skip invalid repo config: owner=%q repo=%q", repo.Owner, repo.Repo)
			continue
		}
		repos = append(repos, normalized)
	}

	if len(repos) == 0 {
		return
	}

	plugin.Register(&githubPlugin{repos: repos, token: token})
}

func (p *githubPlugin) Name() string {
	return "github"
}

func (p *githubPlugin) Fetch() ([]plugin.SoftwareData, error) {
	results := make([]plugin.SoftwareData, 0, len(p.repos))

	for _, repo := range p.repos {
		release, err := util.FetchGitHubLatestReleaseWithToken(repo.Owner, repo.Repo, p.token)
		if err != nil {
			return nil, fmt.Errorf("fetch latest release for %s/%s: %w", repo.Owner, repo.Repo, err)
		}

		results = append(results, plugin.SoftwareData{
			Item: plugin.SoftwareItem{
				ID:              repo.ID,
				Name:            repo.Name,
				Icon:            repo.Icon,
				Description:     repo.Description,
				Organization:    repo.Organization,
				OfficialWebsite: repo.OfficialWebsite,
				Tags:            repo.Tags,
			},
			Versions: []plugin.Version{
				{
					Version:     strings.TrimSpace(release.TagName),
					ReleaseDate: util.FormatGitHubDate(strings.TrimSpace(release.PublishedAt)),
					OfficialURL: strings.TrimSpace(release.HTMLURL),
					Variants:    buildVariants(release.Assets, repo.Assets),
				},
			},
		})
	}

	return results, nil
}

func loadConfig() (*pluginConfig, error) {
	if inline := strings.TrimSpace(os.Getenv(envConfigInline)); inline != "" {
		return parseConfig([]byte(inline))
	}

	candidates := []string{}
	if path := strings.TrimSpace(os.Getenv(envConfigPath)); path != "" {
		candidates = append(candidates, path)
	}
	candidates = append(candidates,
		"config/github.yaml",
		"config/github.yml",
		"../config/github.yaml",
		"../config/github.yml",
	)

	for _, candidate := range candidates {
		cleanPath := filepath.Clean(candidate)
		content, err := os.ReadFile(cleanPath)
		if err == nil {
			return parseConfig(content)
		}
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read %s: %w", cleanPath, err)
		}
	}

	return nil, nil
}

func parseConfig(content []byte) (*pluginConfig, error) {
	cfg := &pluginConfig{}
	if err := yaml.Unmarshal(content, cfg); err != nil {
		return nil, fmt.Errorf("parse github config: %w", err)
	}
	return cfg, nil
}

func normalizeRepoConfig(repo repoConfig) (repoConfig, bool) {
	repo.ID = strings.TrimSpace(repo.ID)
	repo.Owner = strings.TrimSpace(repo.Owner)
	repo.Repo = strings.TrimSpace(repo.Repo)
	repo.Name = strings.TrimSpace(repo.Name)
	repo.Description = strings.TrimSpace(repo.Description)
	repo.Organization = strings.TrimSpace(repo.Organization)
	repo.OfficialWebsite = strings.TrimSpace(repo.OfficialWebsite)
	repo.Icon = strings.TrimSpace(repo.Icon)

	if repo.Owner == "" || repo.Repo == "" {
		return repo, false
	}
	if repo.ID == "" {
		repo.ID = strings.ToLower(strings.ReplaceAll(repo.Repo, " ", "-"))
	}
	if repo.Name == "" {
		repo.Name = repo.Repo
	}
	if repo.Organization == "" {
		repo.Organization = repo.Owner
	}
	if repo.OfficialWebsite == "" {
		repo.OfficialWebsite = fmt.Sprintf("https://github.com/%s/%s", repo.Owner, repo.Repo)
	}
	if repo.Description == "" {
		repo.Description = fmt.Sprintf("%s/%s 的 GitHub Release 下载入口。", repo.Owner, repo.Repo)
	}

	for index := range repo.Assets {
		for ki, kw := range repo.Assets[index].Keywords {
			repo.Assets[index].Keywords[ki] = strings.ToLower(strings.TrimSpace(kw))
		}
		repo.Assets[index].Platform = strings.TrimSpace(repo.Assets[index].Platform)
		repo.Assets[index].Architecture = strings.TrimSpace(repo.Assets[index].Architecture)
		repo.Assets[index].Package = strings.TrimSpace(repo.Assets[index].Package)
	}

	return repo, true
}

func buildVariants(assets []util.GitHubAsset, rules []assetRule) []plugin.Variant {
	type variantKey struct {
		platform     string
		architecture string
		pkg          string
	}

	variantMap := map[variantKey][]plugin.Link{}

	for _, asset := range assets {
		assetName := strings.TrimSpace(asset.Name)
		assetURL := strings.TrimSpace(asset.BrowserDownloadURL)
		if assetName == "" || assetURL == "" {
			continue
		}

		matched := false
		for _, rule := range rules {
			if !matchesAssetRule(asset, rule) {
				continue
			}
			matched = true
			key := variantKey{
				platform:     defaultString(rule.Platform, "通用"),
				architecture: defaultString(rule.Architecture, "通用"),
				pkg:          strings.ToLower(strings.TrimSpace(rule.Package)),
			}
			// 始终使用真实文件名作为链接标签，与 GitHub Releases 页面保持一致
			variantMap[key] = append(variantMap[key], plugin.Link{Type: "direct", Label: assetName, URL: assetURL})
		}

		if len(rules) == 0 && !matched {
			key := variantKey{platform: "通用", architecture: "通用", pkg: ""}
			variantMap[key] = append(variantMap[key], plugin.Link{Type: "direct", Label: assetName, URL: assetURL})
		}
	}

	keys := make([]variantKey, 0, len(variantMap))
	for key := range variantMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].platform == keys[j].platform {
			if keys[i].architecture == keys[j].architecture {
				return keys[i].pkg < keys[j].pkg
			}
			return keys[i].architecture < keys[j].architecture
		}
		return keys[i].platform < keys[j].platform
	})

	variants := make([]plugin.Variant, 0, len(keys))
	for _, key := range keys {
		archLabel := key.architecture
		if key.pkg != "" {
			archLabel = fmt.Sprintf("%s (%s)", archLabel, key.pkg)
		}
		variants = append(variants, plugin.Variant{
			Architecture: archLabel,
			Platform:     key.platform,
			Links:        variantMap[key],
		})
	}
	return variants
}

func matchesAssetRule(asset util.GitHubAsset, rule assetRule) bool {
	if len(rule.Keywords) == 0 {
		return false
	}

	// 使用文件名和下载 URL 共同组成搜索空间（小写）
	haystack := strings.ToLower(strings.TrimSpace(asset.Name)) +
		" " + strings.ToLower(strings.TrimSpace(asset.BrowserDownloadURL))

	// 所有关键词都必须出现（AND 逻辑）
	for _, kw := range rule.Keywords {
		if kw == "" {
			continue
		}
		if !strings.Contains(haystack, kw) {
			return false
		}
	}
	return true
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
