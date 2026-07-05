//go:build integration

package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

// This integration test builds and runs the data-cli against a small
// set of plugins and validates the output JSON files.
// Run with: go test -tags=integration -run TestIntegration ./data-cli/

func TestIntegration_GenerateData(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	outDir := filepath.Join(tmpDir, "output")
	versionsDir := filepath.Join(outDir, "versions")

	// Build the CLI binary
	binaryPath := filepath.Join(tmpDir, "data-cli")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("build data-cli: %v", err)
	}

	// Run with 3 plugins that don't require network to verify basic structure.
	// Use skip-unchanged=false to force writing output.
	run := exec.Command(binaryPath,
		"-out", outDir,
		"-plugins", "steam,7zip",
		"-skip-unchanged=false",
		"-concurrency", "2",
	)
	run.Dir = "."
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	if err := run.Run(); err != nil {
		t.Fatalf("run data-cli: %v", err)
	}

	// ── Validate output files exist ──
	checkFile(t, outDir, "index.json")
	checkFile(t, outDir, "software-list.json")
	checkFile(t, versionsDir, "steam.json")
	checkFile(t, versionsDir, "7zip.json")

	// ── Validate index.json structure ──
	indexData := readJSON(t, outDir, "index.json")
	var index plugin.IndexPayload
	if err := json.Unmarshal(indexData, &index); err != nil {
		t.Fatalf("unmarshal index.json: %v", err)
	}
	if index.Meta.Version == "" {
		t.Error("index.meta.version is empty")
	}
	if index.Meta.Generator != "data-cli" {
		t.Errorf("index.meta.generator = %q, want %q", index.Meta.Generator, "data-cli")
	}
	if index.SoftwareList.Mode != "json" {
		t.Errorf("index.softwareList.mode = %q, want %q", index.SoftwareList.Mode, "json")
	}

	// ── Validate software-list.json structure ──
	listData := readJSON(t, outDir, "software-list.json")
	var list plugin.SoftwareListPayload
	if err := json.Unmarshal(listData, &list); err != nil {
		t.Fatalf("unmarshal software-list.json: %v", err)
	}
	if list.UpdatedAt == "" {
		t.Error("software-list.updatedAt is empty")
	}
	if len(list.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list.Items))
	}

	// Check each item
	itemByID := map[string]plugin.SoftwareItem{}
	for _, item := range list.Items {
		itemByID[item.ID] = item
	}

	for _, expectedID := range []string{"steam", "7zip"} {
		item, ok := itemByID[expectedID]
		if !ok {
			t.Errorf("software-list missing item %q", expectedID)
			continue
		}
		if item.Name == "" {
			t.Errorf("[%s] Name is empty", expectedID)
		}
		if item.Description == "" {
			t.Errorf("[%s] Description is empty", expectedID)
		}
		if item.Source.Path == "" {
			t.Errorf("[%s] Source.Path is empty", expectedID)
		}
		if !strings.Contains(item.Source.Path, expectedID+".json") {
			t.Errorf("[%s] Source.Path = %q, should contain %q", expectedID, item.Source.Path, expectedID+".json")
		}
	}

	// ── Validate version files ──
	for _, id := range []string{"steam", "7zip"} {
		versionData := readJSON(t, versionsDir, id+".json")
		var payload plugin.PlatformPayload
		if err := json.Unmarshal(versionData, &payload); err != nil {
			t.Fatalf("unmarshal %s.json: %v", id, err)
		}
		if payload.SoftwareID != id {
			t.Errorf("softwareId = %q, want %q", payload.SoftwareID, id)
		}
		if payload.UpdatedAt == "" {
			t.Errorf("[%s] UpdatedAt is empty", id)
		}
		if len(payload.Platforms) == 0 {
			t.Errorf("[%s] Platforms is empty", id)
		}
	}
}

func TestIntegration_NoPlugins(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	outDir := filepath.Join(tmpDir, "output")

	binaryPath := filepath.Join(tmpDir, "data-cli")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("build data-cli: %v", err)
	}

	run := exec.Command(binaryPath, "-out", outDir, "-plugins", "nonexistent")
	run.Dir = "."
	output, err := run.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for unknown plugin")
	}
	if !strings.Contains(string(output), "unknown plugins") {
		t.Errorf("output = %q, want to contain 'unknown plugins'", string(output))
	}
}

// ── helpers ────────────────────────────────────────────────────────────────

func checkFile(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("expected file %s does not exist", path)
	}
}

func readJSON(t *testing.T, dir, name string) []byte {
	t.Helper()
	path := filepath.Join(dir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
