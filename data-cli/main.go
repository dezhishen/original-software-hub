package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	// Blank-import each plugin to trigger its init() registration.
	// Uncomment or add plugins here to include them in the build.
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/chrome"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/github"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/qq"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/steam"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

func main() {
	outDir := flag.String("out", "../frontend/data", "Output directory (json/ and jsonp/ will be created inside)")
	flag.Parse()

	jsonDir := filepath.Join(*outDir, "json", "versions")
	jsonpDir := filepath.Join(*outDir, "jsonp", "versions")
	if err := resetDir(jsonDir); err != nil {
		log.Fatalf("reset json dir: %v", err)
	}
	if err := resetDir(jsonpDir); err != nil {
		log.Fatalf("reset jsonp dir: %v", err)
	}
	if err := os.MkdirAll(jsonDir, 0o755); err != nil {
		log.Fatalf("create json dir: %v", err)
	}
	if err := os.MkdirAll(jsonpDir, 0o755); err != nil {
		log.Fatalf("create jsonp dir: %v", err)
	}

	plugins := plugin.All()
	if len(plugins) == 0 {
		log.Println("no plugins registered, nothing to do")
		return
	}

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
				UpdatedAt:  today(),
				Versions:   entry.Versions,
			}
			if err := writeJSON(filepath.Join(jsonDir, softwareID+".json"), versionPayload); err != nil {
				log.Printf("[%s/%s] write json: %v", p.Name(), softwareID, err)
				continue
			}
			if err := writeJSONP(filepath.Join(jsonpDir, softwareID+".js"), versionPayload); err != nil {
				log.Printf("[%s/%s] write jsonp: %v", p.Name(), softwareID, err)
				continue
			}

			item := entry.Item
			item.Source = plugin.Source{
				Mode:          "jsonp",
				Path:          "versions/" + softwareID + ".js",
				CallbackParam: "callback",
				TimeoutMs:     8000,
			}
			listItems = append(listItems, item)
		}
	}

	softwareList := plugin.SoftwareListPayload{
		UpdatedAt: today(),
		Items:     listItems,
	}
	if err := writeJSON(filepath.Join(*outDir, "json", "software-list.json"), softwareList); err != nil {
		log.Fatalf("write software-list.json: %v", err)
	}
	if err := writeJSONP(filepath.Join(*outDir, "jsonp", "software-list.js"), softwareList); err != nil {
		log.Fatalf("write software-list.js: %v", err)
	}

	indexJSON := plugin.IndexPayload{
		Meta: plugin.Meta{
			Version:     "1.0.0",
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			Generator:   "data-cli",
		},
		SoftwareList: plugin.Source{
			Mode:          "json",
			Path:          "software-list.json",
			CallbackParam: "callback",
			TimeoutMs:     8000,
		},
	}
	indexJSONP := plugin.IndexPayload{
		Meta: indexJSON.Meta,
		SoftwareList: plugin.Source{
			Mode:          "jsonp",
			Path:          "software-list.js",
			CallbackParam: "callback",
			TimeoutMs:     8000,
		},
	}
	if err := writeJSON(filepath.Join(*outDir, "json", "index.json"), indexJSON); err != nil {
		log.Fatalf("write index.json: %v", err)
	}
	if err := writeJSONP(filepath.Join(*outDir, "jsonp", "index.js"), indexJSONP); err != nil {
		log.Fatalf("write index.js: %v", err)
	}

	fmt.Println("Done.")
}

func today() string {
	return time.Now().Format("2006-01-02")
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

func writeJSONP(path string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	content := fmt.Sprintf("callback(%s);\n", data)
	return os.WriteFile(path, []byte(content), 0o644)
}
