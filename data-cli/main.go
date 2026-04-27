package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	// Blank-import each plugin to trigger its init() registration.
	// Uncomment or add plugins here to include them in the build.
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/baidunetdisk"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/chrome"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/dingtalk"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/github"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/huorong"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/neteasecloudmusic"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/qq"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/qqmusic"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/steam"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/wecom"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/weixin"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/wps"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
	"github.com/mozillazg/go-pinyin"
)

func main() {
	outDir := flag.String("out", "../frontend/data/json", "Output directory for index.json, software-list.json and versions/")
	pluginsArg := flag.String("plugins", "all", "Plugins to run: all or comma-separated names (e.g. weixin,qq)")
	concurrency := flag.Int("concurrency", 3, "Maximum number of plugins to run concurrently")
	flag.Parse()

	jsonDir := filepath.Join(*outDir, "versions")
	if err := resetDir(jsonDir); err != nil {
		log.Fatalf("reset json dir: %v", err)
	}
	if err := os.MkdirAll(jsonDir, 0o755); err != nil {
		log.Fatalf("create json dir: %v", err)
	}

	plugins, err := selectPlugins(plugin.All(), *pluginsArg)
	if err != nil {
		log.Fatalf("select plugins: %v", err)
	}
	if len(plugins) == 0 {
		log.Println("no plugins registered, nothing to do")
		return
	}
	nowUTC := time.Now().UTC()
	updatedAt := nowUTC.Format(time.RFC3339)

	listItems := make([]plugin.SoftwareItem, 0, len(plugins))
	fetchResults := fetchPluginsConcurrently(plugins, *concurrency)
	for _, result := range fetchResults {
		p := result.Plugin
		if result.Err != nil {
			log.Printf("[%s] Fetch error: %v", p.Name(), result.Err)
			continue
		}

		for _, entry := range result.Items {
			softwareID := entry.Item.ID
			if softwareID == "" {
				log.Printf("[%s] skip item with empty id", p.Name())
				continue
			}

			versionPayload := plugin.VersionPayload{
				SoftwareID: softwareID,
				UpdatedAt:  updatedAt,
				Versions:   entry.Versions,
			}
			if err := writeJSON(filepath.Join(jsonDir, softwareID+".json"), versionPayload); err != nil {
				log.Printf("[%s/%s] write json: %v", p.Name(), softwareID, err)
				continue
			}

			item := entry.Item
			item.Pinyin = buildSearchPinyin(item.Name)
			item.Source = plugin.Source{
				Mode:      "json",
				Path:      "versions/" + softwareID + ".json",
				TimeoutMs: 8000,
			}
			listItems = append(listItems, item)
		}
	}

	sort.SliceStable(listItems, func(i, j int) bool {
		li := sortKeyForItem(listItems[i])
		lj := sortKeyForItem(listItems[j])
		if li == lj {
			return listItems[i].ID < listItems[j].ID
		}
		return li < lj
	})

	softwareList := plugin.SoftwareListPayload{
		UpdatedAt: updatedAt,
		Items:     listItems,
	}
	if err := writeJSON(filepath.Join(*outDir, "software-list.json"), softwareList); err != nil {
		log.Fatalf("write software-list.json: %v", err)
	}

	indexJSON := plugin.IndexPayload{
		Meta: plugin.Meta{
			Version:     "1.0.0",
			GeneratedAt: updatedAt,
			Generator:   "data-cli",
		},
		SoftwareList: plugin.Source{
			Mode:      "json",
			Path:      "software-list.json",
			TimeoutMs: 8000,
		},
	}
	if err := writeJSON(filepath.Join(*outDir, "index.json"), indexJSON); err != nil {
		log.Fatalf("write index.json: %v", err)
	}

	fmt.Println("Done.")
}

type pluginJob struct {
	Index  int
	Plugin plugin.Plugin
}

type pluginFetchResult struct {
	Index  int
	Plugin plugin.Plugin
	Items  []plugin.SoftwareData
	Err    error
}

func fetchPluginsConcurrently(plugins []plugin.Plugin, maxConcurrency int) []pluginFetchResult {
	if len(plugins) == 0 {
		return nil
	}
	if maxConcurrency < 1 {
		maxConcurrency = 1
	}
	if maxConcurrency > len(plugins) {
		maxConcurrency = len(plugins)
	}

	jobs := make(chan pluginJob)
	results := make(chan pluginFetchResult, len(plugins))

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for job := range jobs {
			fmt.Printf("[%s] fetching data...\n", job.Plugin.Name())
			items, err := job.Plugin.Fetch()
			results <- pluginFetchResult{
				Index:  job.Index,
				Plugin: job.Plugin,
				Items:  items,
				Err:    err,
			}
		}
	}

	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go worker()
	}

	for i, p := range plugins {
		jobs <- pluginJob{Index: i, Plugin: p}
	}
	close(jobs)

	wg.Wait()
	close(results)

	out := make([]pluginFetchResult, 0, len(plugins))
	for result := range results {
		out = append(out, result)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Index < out[j].Index
	})

	return out
}

func selectPlugins(all []plugin.Plugin, pluginsArg string) ([]plugin.Plugin, error) {
	raw := strings.TrimSpace(strings.ToLower(pluginsArg))
	if raw == "" || raw == "all" {
		return all, nil
	}

	selectedNames := map[string]struct{}{}
	for _, name := range strings.Split(raw, ",") {
		n := strings.TrimSpace(strings.ToLower(name))
		if n == "" {
			continue
		}
		selectedNames[n] = struct{}{}
	}
	if len(selectedNames) == 0 {
		return nil, fmt.Errorf("invalid -plugins value: %q", pluginsArg)
	}

	filtered := make([]plugin.Plugin, 0, len(selectedNames))
	found := map[string]struct{}{}
	available := make([]string, 0, len(all))
	for _, p := range all {
		name := strings.ToLower(strings.TrimSpace(p.Name()))
		if name == "" {
			continue
		}
		available = append(available, name)
		if _, ok := selectedNames[name]; ok {
			filtered = append(filtered, p)
			found[name] = struct{}{}
		}
	}

	missing := make([]string, 0)
	for name := range selectedNames {
		if _, ok := found[name]; !ok {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		sort.Strings(available)
		return nil, fmt.Errorf("unknown plugins: %s (available: %s)", strings.Join(missing, ","), strings.Join(available, ","))
	}

	return filtered, nil
}

func resetDir(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return os.MkdirAll(path, 0o755)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func sortKeyForItem(item plugin.SoftwareItem) string {
	if item.Pinyin != "" {
		return strings.ToLower(strings.Fields(item.Pinyin)[0])
	}
	return strings.ToLower(strings.TrimSpace(item.Name))
}

func buildSearchPinyin(name string) string {
	clean := strings.TrimSpace(name)
	if clean == "" || !containsHan(clean) {
		return ""
	}

	args := pinyin.NewArgs()
	args.Style = pinyin.Normal
	pyGroups := pinyin.Pinyin(clean, args)

	parts := make([]string, 0, len(pyGroups))
	abbr := make([]rune, 0, len(pyGroups))
	for _, g := range pyGroups {
		if len(g) == 0 {
			continue
		}
		s := strings.ToLower(strings.TrimSpace(g[0]))
		if s == "" {
			continue
		}
		parts = append(parts, s)
		abbr = append(abbr, []rune(s)[0])
	}

	full := strings.Join(parts, "")
	if full == "" {
		return ""
	}
	short := string(abbr)
	if short != "" && short != full {
		return full + " " + short
	}
	return full
}

func containsHan(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}
