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
	"time"

	// Blank-import each plugin to trigger its init() registration.
	// Uncomment or add plugins here to include them in the build.
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/chrome"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/github"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/huorong"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/qq"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/steam"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/weixin"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

func main() {
	outDir := flag.String("out", "../frontend/data", "Output directory (json/ will be created inside)")
	pluginsArg := flag.String("plugins", "all", "Plugins to run: all or comma-separated names (e.g. weixin,qq)")
	flag.Parse()

	jsonDir := filepath.Join(*outDir, "json", "versions")
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
	for _, p := range plugins {
		fmt.Printf("[%s] fetching data...\n", p.Name())
		items, err := p.Fetch()
		if err != nil {
			log.Printf("[%s] Fetch error: %v", p.Name(), err)
			continue
		}

		for _, entry := range items {
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
			item.Source = plugin.Source{
				Mode:      "json",
				Path:      "versions/" + softwareID + ".json",
				TimeoutMs: 8000,
			}
			listItems = append(listItems, item)
		}
	}

	softwareList := plugin.SoftwareListPayload{
		UpdatedAt: updatedAt,
		Items:     listItems,
	}
	if err := writeJSON(filepath.Join(*outDir, "json", "software-list.json"), softwareList); err != nil {
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
	if err := writeJSON(filepath.Join(*outDir, "json", "index.json"), indexJSON); err != nil {
		log.Fatalf("write index.json: %v", err)
	}

	fmt.Println("Done.")
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
