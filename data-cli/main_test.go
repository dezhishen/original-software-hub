package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
)

// ── flattenVersionsToPlatforms ─────────────────────────────────────────────

func TestFlattenVersionsToPlatforms_empty(t *testing.T) {
	result := flattenVersionsToPlatforms(nil)
	if len(result) != 0 {
		t.Errorf("expected empty, got %d", len(result))
	}
}

func TestFlattenVersionsToPlatforms_singleVersion(t *testing.T) {
	versions := []plugin.Version{
		{
			Version: "1.0", ReleaseDate: "2026-01-01", OfficialURL: "https://example.com",
			Platforms: []plugin.PlatformRelease{
				{
					Platform: "Windows", Version: "1.0", ReleaseDate: "2026-01-01",
					Packages: []plugin.PlatformPackage{
						{Architecture: "x64", Links: []plugin.Link{{Type: "direct", URL: "https://example.com/win.exe"}}},
					},
				},
			},
		},
	}
	result := flattenVersionsToPlatforms(versions)
	if len(result) != 1 {
		t.Fatalf("expected 1 platform, got %d", len(result))
	}
	if result[0].Platform != "Windows" {
		t.Errorf("platform = %q", result[0].Platform)
	}
	if result[0].Version != "1.0" {
		t.Errorf("version = %q", result[0].Version)
	}
	if result[0].ReleaseDate != "2026-01-01" {
		t.Errorf("releaseDate = %q", result[0].ReleaseDate)
	}
	if len(result[0].Packages) != 1 {
		t.Fatalf("expected 1 package, got %d", len(result[0].Packages))
	}
	if result[0].Packages[0].Architecture != "x64" {
		t.Errorf("arch = %q", result[0].Packages[0].Architecture)
	}
}

func TestFlattenVersionsToPlatforms_mergesMultipleVersions(t *testing.T) {
	versions := []plugin.Version{
		{
			Version: "1.0", ReleaseDate: "2026-01-01",
			Platforms: []plugin.PlatformRelease{
				{
					Platform: "Windows", Version: "1.0", ReleaseDate: "2026-01-01",
					Packages: []plugin.PlatformPackage{
						{Architecture: "x64", Links: []plugin.Link{{Type: "direct", URL: "https://example.com/old.exe"}}},
					},
				},
			},
		},
		{
			Version: "2.0", ReleaseDate: "2026-06-01",
			Platforms: []plugin.PlatformRelease{
				{
					Platform: "Windows", Version: "2.0", ReleaseDate: "2026-06-01",
					Packages: []plugin.PlatformPackage{
						{Architecture: "x64", Links: []plugin.Link{{Type: "direct", URL: "https://example.com/new.exe"}}},
					},
				},
				{
					Platform: "macOS", Version: "2.0", ReleaseDate: "2026-06-01",
					Packages: []plugin.PlatformPackage{
						{Architecture: "arm64", Links: []plugin.Link{{Type: "direct", URL: "https://example.com/mac.dmg"}}},
					},
				},
			},
		},
	}
	result := flattenVersionsToPlatforms(versions)
	if len(result) != 2 {
		t.Fatalf("expected 2 platforms, got %d: %+v", len(result), result)
	}
	// Should use newest version info
	for _, p := range result {
		if p.Platform == "Windows" {
			if p.Version != "2.0" {
				t.Errorf("Windows version = %q, want %q", p.Version, "2.0")
			}
		}
	}
}

func TestFlattenVersionsToPlatforms_deduplicates(t *testing.T) {
	versions := []plugin.Version{
		{
			Version: "1.0",
			Platforms: []plugin.PlatformRelease{
				{Platform: "Windows", Version: "1.0", Packages: []plugin.PlatformPackage{
					{Architecture: "x64", Links: []plugin.Link{
						{Type: "direct", URL: "https://example.com/same.exe"},
						{Type: "direct", URL: "https://example.com/same.exe"}, // duplicate
					}},
				}},
			},
		},
	}
	result := flattenVersionsToPlatforms(versions)
	if len(result) == 0 || len(result[0].Packages) == 0 {
		t.Fatal("expected at least 1 platform with 1 package")
	}
	if len(result[0].Packages[0].Links) != 1 {
		t.Errorf("expected 1 link after dedup, got %d", len(result[0].Packages[0].Links))
	}
}

// ── sortKeyForItem ─────────────────────────────────────────────────────────

func TestSortKeyForItem_usesPinyin(t *testing.T) {
	item := plugin.SoftwareItem{Name: "微信", Pinyin: "weixin w x"}
	key := sortKeyForItem(item)
	if key != "weixin" {
		t.Errorf("key = %q, want %q", key, "weixin")
	}
}

func TestSortKeyForItem_fallsBackToName(t *testing.T) {
	item := plugin.SoftwareItem{Name: "7-Zip"}
	key := sortKeyForItem(item)
	if key != "7-zip" {
		t.Errorf("key = %q, want %q", key, "7-zip")
	}
}

func TestSortKeyForItem_respectsOrder(t *testing.T) {
	items := []plugin.SoftwareItem{
		{Name: "360 安全卫士", Pinyin: "360 anquan weishi 3"},
		{Name: "7-Zip"},
		{Name: "Chrome"},
		{Name: "微信", Pinyin: "weixin w x"},
	}
	sort.SliceStable(items, func(i, j int) bool {
		return sortKeyForItem(items[i]) < sortKeyForItem(items[j])
	})
	expected := []string{"360 安全卫士", "7-Zip", "Chrome", "微信"}
	for i, item := range items {
		if item.Name != expected[i] {
			t.Errorf("order[%d] = %q, want %q", i, item.Name, expected[i])
		}
	}
}

// ── sanitizeFileStem ───────────────────────────────────────────────────────

func TestSanitizeFileStem(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"7-Zip", "7-zip"},
		{"Visual Studio Code", "visual-studio-code"},
		{"360安全卫士", "360"},
		{"  Chrome  ", "chrome"},
		{"", ""},
		{"---abc---", "abc"},
	}
	for _, tt := range tests {
		got := sanitizeFileStem(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeFileStem(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ── selectPlugins ──────────────────────────────────────────────────────────

type mockPlugin struct {
	name     string
	disabled bool
}

func (m *mockPlugin) Name() string                          { return m.name }
func (m *mockPlugin) Fetch() ([]plugin.SoftwareData, error) { return nil, nil }
func (m *mockPlugin) Disabled() bool                        { return m.disabled }

func TestSelectPlugins_all(t *testing.T) {
	all := []plugin.Plugin{
		&mockPlugin{name: "foo", disabled: false},
		&mockPlugin{name: "bar", disabled: true},
		&mockPlugin{name: "baz", disabled: false},
	}
	selected, err := selectPlugins(all, "all")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(selected) != 2 {
		t.Fatalf("expected 2 enabled plugins, got %d", len(selected))
	}
}

func TestSelectPlugins_commaSeparated(t *testing.T) {
	all := []plugin.Plugin{
		&mockPlugin{name: "foo", disabled: false},
		&mockPlugin{name: "bar", disabled: true},
		&mockPlugin{name: "baz", disabled: false},
	}
	selected, err := selectPlugins(all, "foo,baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(selected) != 2 {
		t.Fatalf("expected 2 plugins, got %d", len(selected))
	}
}

func TestSelectPlugins_unknown(t *testing.T) {
	all := []plugin.Plugin{&mockPlugin{name: "foo", disabled: false}}
	_, err := selectPlugins(all, "nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown plugin")
	}
	if !strings.Contains(err.Error(), "unknown plugins") {
		t.Errorf("error = %q, want to contain 'unknown plugins'", err.Error())
	}
}

// ── buildSearchPinyin ──────────────────────────────────────────────────────

func TestBuildSearchPinyin(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"微信", "weixin wx"},
		{"Chrome", ""},
		{"", ""},
		{"   ", ""},
		// "360" is dropped by pinyin lib since it only handles Chinese characters
		{"360安全卫士", "anquanweishi aqws"},
	}
	for _, tt := range tests {
		got := buildSearchPinyin(tt.input)
		if got != tt.want {
			t.Errorf("buildSearchPinyin(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// ── resolveDataByVersionDecision ───────────────────────────────────────────

func TestResolveDataByVersionDecision_noPrevious(t *testing.T) {
	fetched := plugin.SoftwareData{
		Item: plugin.SoftwareItem{ID: "test-software", Name: "Test"},
		Versions: []plugin.Version{
			{Version: "1.0", Platforms: []plugin.PlatformRelease{
				{Platform: "Windows", Packages: []plugin.PlatformPackage{
					{Architecture: "x64", Links: []plugin.Link{{Type: "direct", URL: "https://example.com/test.exe"}}},
				}},
			}},
		},
	}
	previous := plugin.PreviousState{Versions: map[string]plugin.PlatformPayload{}, Items: map[string]plugin.SoftwareItem{}}

	decision, err := resolveDataByVersionDecision(previous, fetched, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !decision.Changed {
		t.Error("expected Changed=true for new software")
	}
	if decision.Reason != "no-previous-version" {
		t.Errorf("reason = %q, want %q", decision.Reason, "no-previous-version")
	}
}

func TestResolveDataByVersionDecision_unchanged(t *testing.T) {
	item := plugin.SoftwareItem{ID: "test", Name: "Test"}
	platforms := []plugin.PlatformRelease{
		{Platform: "Windows", Version: "1.0", Packages: []plugin.PlatformPackage{
			{Architecture: "x64", Links: []plugin.Link{{Type: "direct", URL: "https://example.com/test.exe"}}},
		}},
	}
	fetched := plugin.SoftwareData{Item: item, Versions: []plugin.Version{
		{Version: "1.0", Platforms: platforms},
	}}
	previous := plugin.PreviousState{
		Versions: map[string]plugin.PlatformPayload{
			"test": {SoftwareID: "test", Platforms: platforms},
		},
		Items: map[string]plugin.SoftwareItem{"test": item},
	}

	decision, err := resolveDataByVersionDecision(previous, fetched, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decision.Changed {
		t.Error("expected Changed=false for unchanged software")
	}
	if decision.Reason != "version-unchanged" {
		t.Errorf("reason = %q, want %q", decision.Reason, "version-unchanged")
	}
}

func TestResolveDataByVersionDecision_changedVersion(t *testing.T) {
	oldPlatforms := []plugin.PlatformRelease{
		{Platform: "Windows", Version: "1.0", Packages: []plugin.PlatformPackage{
			{Architecture: "x64", Links: []plugin.Link{{Type: "direct", URL: "https://example.com/old.exe"}}},
		}},
	}
	newPlatforms := []plugin.PlatformRelease{
		{Platform: "Windows", Version: "2.0", Packages: []plugin.PlatformPackage{
			{Architecture: "x64", Links: []plugin.Link{{Type: "direct", URL: "https://example.com/new.exe"}}},
		}},
	}
	fetched := plugin.SoftwareData{
		Item:     plugin.SoftwareItem{ID: "test", Name: "Test"},
		Versions: []plugin.Version{{Version: "2.0", Platforms: newPlatforms}},
	}
	previous := plugin.PreviousState{
		Versions: map[string]plugin.PlatformPayload{
			"test": {SoftwareID: "test", Platforms: oldPlatforms},
		},
		Items: map[string]plugin.SoftwareItem{"test": {ID: "test"}},
	}

	decision, err := resolveDataByVersionDecision(previous, fetched, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !decision.Changed {
		t.Error("expected Changed=true for changed version")
	}
	if decision.Reason != "version-changed" {
		t.Errorf("reason = %q, want %q", decision.Reason, "version-changed")
	}
}

// ── pluginPriorityRank ─────────────────────────────────────────────────────

func TestPluginPriorityRank_default(t *testing.T) {
	rank := pluginPriorityRank("unknown-plugin")
	if rank != 100 {
		t.Errorf("default rank = %d, want 100", rank)
	}
}

func TestPluginPriorityRank_known(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"github", 10},
		{"chrome", 20},
		{"firefox", 30},
		{"wps", 40},
		{"todesk", 50},
	}
	for _, tt := range tests {
		rank := pluginPriorityRank(tt.name)
		if rank != tt.want {
			t.Errorf("pluginPriorityRank(%q) = %d, want %d", tt.name, rank, tt.want)
		}
	}
}

// ── ensureDir ──────────────────────────────────────────────────────────────

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "a", "b", "c")
	if err := ensureDir(testDir); err != nil {
		t.Fatalf("ensureDir: %v", err)
	}
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("directory was not created")
	}
}

// ── writeJSON / loadPreviousState ──────────────────────────────────────────

func TestWriteJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.json")
	data := map[string]string{"hello": "world"}
	if err := writeJSON(path, data); err != nil {
		t.Fatalf("writeJSON: %v", err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var decoded map[string]string
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded["hello"] != "world" {
		t.Errorf("got %q, want %q", decoded["hello"], "world")
	}
}

func TestLoadPreviousState_empty(t *testing.T) {
	tmpDir := t.TempDir()
	versionsDir := filepath.Join(tmpDir, "versions")
	os.MkdirAll(versionsDir, 0755)

	versions, list, index := loadPreviousState(tmpDir, versionsDir)
	if len(versions) != 0 {
		t.Errorf("expected empty versions, got %d", len(versions))
	}
	if len(list.Items) != 0 {
		t.Errorf("expected empty list, got %d", len(list.Items))
	}
	if index.Meta.Version != "" {
		t.Errorf("expected empty meta version, got %q", index.Meta.Version)
	}
}

func TestLoadPreviousState_withFiles(t *testing.T) {
	tmpDir := t.TempDir()
	versionsDir := filepath.Join(tmpDir, "versions")
	os.MkdirAll(versionsDir, 0755)

	// Write software-list.json
	listPayload := plugin.SoftwareListPayload{
		UpdatedAt: "2026-07-05T00:00:00Z",
		Items: []plugin.SoftwareItem{
			{ID: "test", Name: "Test Software", Description: "A test"},
		},
	}
	writeJSON(filepath.Join(tmpDir, "software-list.json"), listPayload)

	// Write a version file
	platformPayload := plugin.PlatformPayload{
		SoftwareID: "test",
		UpdatedAt:  "2026-07-05T00:00:00Z",
		Platforms: []plugin.PlatformRelease{
			{Platform: "Windows", Version: "1.0"},
		},
	}
	writeJSON(filepath.Join(versionsDir, "test.json"), platformPayload)

	// Write index.json
	indexPayload := plugin.IndexPayload{
		Meta: plugin.Meta{Version: "1.0.0", GeneratedAt: "2026-07-05T00:00:00Z", Generator: "data-cli"},
	}
	writeJSON(filepath.Join(tmpDir, "index.json"), indexPayload)

	versions, list, index := loadPreviousState(tmpDir, versionsDir)
	if len(versions) != 1 {
		t.Errorf("expected 1 version, got %d", len(versions))
	}
	if len(list.Items) != 1 {
		t.Errorf("expected 1 list item, got %d", len(list.Items))
	}
	if index.Meta.Version != "1.0.0" {
		t.Errorf("meta version = %q", index.Meta.Version)
	}

	v, ok := versions["test"]
	if !ok {
		t.Fatal("test software not found in versions map")
	}
	if len(v.Platforms) != 1 {
		t.Errorf("expected 1 platform, got %d", len(v.Platforms))
	}
}

// ── containsHan ────────────────────────────────────────────────────────────

func TestContainsHan(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"微信", true},
		{"Chrome", false},
		{"", false},
		{"360安全卫士", true},
		{"7-Zip", false},
	}
	for _, tt := range tests {
		got := containsHan(tt.input)
		if got != tt.want {
			t.Errorf("containsHan(%q) = %t, want %t", tt.input, got, tt.want)
		}
	}
}

// ── isRemoteHTTPURL ────────────────────────────────────────────────────────

func TestIsRemoteHTTPURL(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"https://example.com/file.exe", true},
		{"http://example.com/file.exe", true},
		{"ftp://example.com/file.exe", false},
		{"", false},
		{"/local/path/file.exe", false},
		{"file:///local/file.exe", false},
	}
	for _, tt := range tests {
		got := isRemoteHTTPURL(tt.input)
		if got != tt.want {
			t.Errorf("isRemoteHTTPURL(%q) = %t, want %t", tt.input, got, tt.want)
		}
	}
}

// ── retryDelay ─────────────────────────────────────────────────────────────

func TestRetryDelay(t *testing.T) {
	base := time.Second
	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{1, time.Second},
		{2, 2 * time.Second},
		{3, 4 * time.Second},
		{4, 8 * time.Second},
	}
	for _, tt := range tests {
		got := retryDelay(base, tt.attempt)
		if got != tt.want {
			t.Errorf("retryDelay(%d) = %v, want %v", tt.attempt, got, tt.want)
		}
	}
}
