package plugin

import (
	"testing"
)

func TestPlatformsFromVariants_basic(t *testing.T) {
	variants := []Variant{
		{Architecture: "x64", Platform: "Windows", Links: []Link{{Type: "direct", URL: "https://example.com/win-x64.exe", Label: "Windows x64"}}},
		{Architecture: "x86", Platform: "Windows", Links: []Link{{Type: "direct", URL: "https://example.com/win-x86.exe", Label: "Windows x86"}}},
		{Architecture: "arm64", Platform: "macOS", Links: []Link{{Type: "direct", URL: "https://example.com/mac-arm64.dmg", Label: "macOS ARM"}}},
		{Architecture: "x64", Platform: "Linux", Links: []Link{{Type: "direct", URL: "https://example.com/linux-x64.deb", Label: "Linux x64"}}},
	}

	platforms := PlatformsFromVariants("1.0.0", "2026-01-15", "https://example.com", variants)
	if len(platforms) != 3 {
		t.Fatalf("expected 3 platforms, got %d", len(platforms))
	}

	// Check platform order (sorted alphabetically)
	expectedOrder := []string{"Linux", "Windows", "macOS"}
	for i, p := range platforms {
		if p.Platform != expectedOrder[i] {
			t.Errorf("platform[%d] = %q, want %q", i, p.Platform, expectedOrder[i])
		}
		if p.Version != "1.0.0" {
			t.Errorf("platform[%d].Version = %q, want %q", i, p.Version, "1.0.0")
		}
		if p.ReleaseDate != "2026-01-15" {
			t.Errorf("platform[%d].ReleaseDate = %q, want %q", i, p.ReleaseDate, "2026-01-15")
		}
	}
}

func TestPlatformsFromVariants_empty(t *testing.T) {
	platforms := PlatformsFromVariants("1.0.0", "2026-01-15", "https://example.com", nil)
	if platforms != nil {
		t.Errorf("expected nil for empty variants, got %v", platforms)
	}
}

func TestPlatformsFromVariants_single(t *testing.T) {
	variants := []Variant{
		{Architecture: "x64", Platform: "Windows", Links: []Link{{Type: "direct", URL: "https://example.com/setup.exe", Label: "安装包"}}},
	}
	platforms := PlatformsFromVariants("2.0.0", "", "https://example.com", variants)
	if len(platforms) != 1 {
		t.Fatalf("expected 1 platform, got %d", len(platforms))
	}
	if platforms[0].Platform != "Windows" {
		t.Errorf("platform = %q, want %q", platforms[0].Platform, "Windows")
	}
	if platforms[0].ReleaseDate != "" {
		t.Errorf("expected empty ReleaseDate, got %q", platforms[0].ReleaseDate)
	}
	if len(platforms[0].Packages) != 1 {
		t.Fatalf("expected 1 package, got %d", len(platforms[0].Packages))
	}
	if platforms[0].Packages[0].Architecture != "x64" {
		t.Errorf("arch = %q, want %q", platforms[0].Packages[0].Architecture, "x64")
	}
}

func TestPlatformsFromVariants_unknownPlatform(t *testing.T) {
	variants := []Variant{
		{Architecture: "x64", Platform: "", Links: []Link{{Type: "direct", URL: "https://example.com/pkg", Label: "Package"}}},
	}
	platforms := PlatformsFromVariants("1.0", "2026-06-01", "https://example.com", variants)
	if len(platforms) != 1 {
		t.Fatalf("expected 1 platform, got %d", len(platforms))
	}
	if platforms[0].Platform != "Unknown" {
		t.Errorf("platform = %q, want %q", platforms[0].Platform, "Unknown")
	}
}

func TestPlatformsFromVariants_multipleArchs(t *testing.T) {
	variants := []Variant{
		{Architecture: "arm64", Platform: "Windows", Links: []Link{{Type: "direct", URL: "https://example.com/arm64.exe", Label: "ARM64"}}},
		{Architecture: "x64", Platform: "Windows", Links: []Link{{Type: "direct", URL: "https://example.com/x64.exe", Label: "x64"}}},
		{Architecture: "x86", Platform: "Windows", Links: []Link{{Type: "direct", URL: "https://example.com/x86.exe", Label: "x86"}}},
	}
	platforms := PlatformsFromVariants("3.0", "2026-07-01", "https://example.com", variants)
	if len(platforms) != 1 {
		t.Fatalf("expected 1 platform, got %d", len(platforms))
	}
	// Should be sorted alphabetically by architecture
	archs := make([]string, len(platforms[0].Packages))
	for i, pkg := range platforms[0].Packages {
		archs[i] = pkg.Architecture
	}
	expected := []string{"arm64", "x64", "x86"}
	for i, a := range archs {
		if a != expected[i] {
			t.Errorf("arch[%d] = %q, want %q", i, a, expected[i])
		}
	}
}
