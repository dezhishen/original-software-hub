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

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

func main() {
	outDir := flag.String("out", "../frontend/data", "Output directory (json/ and jsonp/ will be created inside)")
	flag.Parse()

	jsonDir := filepath.Join(*outDir, "json", "versions")
	jsonpDir := filepath.Join(*outDir, "jsonp", "versions")
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
		fmt.Printf("[%s] fetching versions...\n", p.ID())
		info, err := p.FetchSoftwareInfo()
		if err != nil {
			log.Printf("[%s] FetchSoftwareInfo error: %v", p.ID(), err)
			continue
		}
		versions, err := p.FetchVersions()
		if err != nil {
			log.Printf("[%s] FetchVersions error: %v", p.ID(), err)
			continue
		}

		versionPayload := plugin.VersionPayload{
			SoftwareID: p.ID(),
			UpdatedAt:  today(),
			Versions:   versions,
		}
		if err := writeJSON(filepath.Join(jsonDir, p.ID()+".json"), versionPayload); err != nil {
			log.Printf("[%s] write json: %v", p.ID(), err)
			continue
		}
		if err := writeJSONP(filepath.Join(jsonpDir, p.ID()+".js"), versionPayload); err != nil {
			log.Printf("[%s] write jsonp: %v", p.ID(), err)
			continue
		}

		item := *info
		item.Source = plugin.Source{
			Mode:          "jsonp",
			Path:          "versions/" + p.ID() + ".js",
			CallbackParam: "callback",
			TimeoutMs:     8000,
		}
		listItems = append(listItems, item)
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
