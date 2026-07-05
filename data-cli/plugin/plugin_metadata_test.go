package plugin

import (
	"strings"
	"testing"
)

// AllPluginTypes returns all registered plugins for testing.
// Tests in this file verify every plugin's metadata and data structure.
// Note: when running `go test ./plugin/`, the init() functions in individual
// plugin subdirectories (e.g., 7zip, chrome) do NOT execute because they
// are separate Go packages. Only plugins registered within the plugin package
// itself will appear. Use `go test ./...` from the data-cli root to test all.
func TestAllPlugins_nameNotEmpty(t *testing.T) {
	for _, p := range All() {
		name := strings.TrimSpace(p.Name())
		if name == "" {
			t.Error("plugin with empty Name() found")
		}
	}
}

func TestAllPlugins_fetchReturnsNonNil(t *testing.T) {
	for _, p := range All() {
		data, err := p.Fetch()
		if err != nil {
			// Skip network-dependent plugins in offline test.
			// If the plugin is disabled, skip the fetch verification.
			if p.Disabled() {
				continue
			}
			t.Logf("[%s] Fetch() returned error (may be network-dependent): %v", p.Name(), err)
			continue
		}
		if data == nil {
			t.Errorf("[%s] Fetch() returned nil data", p.Name())
		}
	}
}

// TestAllPlugins_metadataStructure validates that each plugin's returned data
// has the required metadata fields populated.
func TestAllPlugins_metadataStructure(t *testing.T) {
	for _, p := range All() {
		data, err := p.Fetch()
		if err != nil || data == nil {
			continue
		}

		for _, item := range data {
			sw := item.Item
			id := strings.TrimSpace(sw.ID)
			name := strings.TrimSpace(sw.Name)

			if id == "" {
				t.Errorf("[%s] Item.ID is empty", p.Name())
			}
			if name == "" {
				t.Errorf("[%s] Item.Name is empty", p.Name())
			}
			if !isValidID(id) {
				t.Errorf("[%s] Item.ID %q contains invalid characters", p.Name(), id)
			}

			// ID should match the expected pattern: lowercase letters, digits, hyphens
			for _, r := range id {
				if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
					t.Errorf("[%s] Item.ID %q contains invalid char %q", p.Name(), id, string(r))
				}
			}
		}
	}
}

// TestAllPlugins_versionStructure validates the version data structure.
func TestAllPlugins_versionStructure(t *testing.T) {
	for _, p := range All() {
		data, err := p.Fetch()
		if err != nil || data == nil {
			continue
		}

		for _, item := range data {
			for vi, v := range item.Versions {
				ver := strings.TrimSpace(v.Version)
				if ver == "" {
					t.Errorf("[%s] versions[%d].Version is empty", p.Name(), vi)
				}

				for pi, plat := range v.Platforms {
					platform := strings.TrimSpace(plat.Platform)
					if platform == "" {
						t.Errorf("[%s] versions[%d].Platforms[%d].Platform is empty", p.Name(), vi, pi)
					}

					for pki, pkg := range plat.Packages {
						arch := strings.TrimSpace(pkg.Architecture)
						if arch == "" {
							t.Errorf("[%s] versions[%d].Platforms[%d].Packages[%d].Architecture is empty", p.Name(), vi, pi, pki)
						}

						for li, link := range pkg.Links {
							url := strings.TrimSpace(link.URL)
							linkType := strings.TrimSpace(link.Type)
							if url == "" {
								t.Errorf("[%s] versions[%d].Platforms[%d].Packages[%d].Links[%d].URL is empty", p.Name(), vi, pi, pki, li)
							}
							if linkType == "" {
								t.Errorf("[%s] versions[%d].Platforms[%d].Packages[%d].Links[%d].Type is empty", p.Name(), vi, pi, pki, li)
							}
							if linkType != "" && !isValidLinkType(linkType) {
								t.Errorf("[%s] versions[%d].Platforms[%d].Packages[%d].Links[%d].Type=%q is not a known type", p.Name(), vi, pi, pki, li, linkType)
							}
							if linkType == "direct" && !isValidURL(url) {
								t.Errorf("[%s] versions[%d].Platforms[%d].Packages[%d].Links[%d].URL=%q is invalid", p.Name(), vi, pi, pki, li, url)
							}
						}
					}
				}
			}
		}
	}
}

// TestAllPlugins_metadataConsistency checks that plugin name and ID are consistent.
func TestAllPlugins_metadataConsistency(t *testing.T) {
	for _, p := range All() {
		data, err := p.Fetch()
		if err != nil || data == nil {
			continue
		}
		for _, item := range data {
			// Plugins's Name() should match or relate to Item.ID
			pName := p.Name()
			swID := item.Item.ID
			if pName != swID && !strings.HasPrefix(swID, pName) && !strings.HasPrefix(pName, swID) {
				t.Logf("[%s] plugin Name()=%q vs Item.ID=%q differ", pName, pName, swID)
			}

			// Description should not be empty
			if strings.TrimSpace(item.Item.Description) == "" {
				t.Errorf("[%s/%s] Description is empty", p.Name(), swID)
			}

			// OfficialWebsite should be a valid URL if non-empty
			website := strings.TrimSpace(item.Item.OfficialWebsite)
			if website != "" && !isValidURL(website) {
				t.Errorf("[%s/%s] OfficialWebsite=%q is not a valid URL", p.Name(), swID, website)
			}
		}
	}
}

// TestAllPlugins_platformsConsistency checks that version/platform metadata is consistent.
func TestAllPlugins_platformsConsistency(t *testing.T) {
	validPlatforms := map[string]bool{
		"Windows": true, "macOS": true, "Linux": true,
		"Android": true, "iOS / iPadOS": true,
		"Web": true, "Windows (Store)": true,
	}
	for _, p := range All() {
		data, err := p.Fetch()
		if err != nil || data == nil {
			continue
		}
		for _, item := range data {
			for _, v := range item.Versions {
				for _, plat := range v.Platforms {
					if !validPlatforms[plat.Platform] {
						t.Logf("[%s] uses non-standard platform %q", p.Name(), plat.Platform)
					}
				}
			}
		}
	}
}

// ── helpers ────────────────────────────────────────────────────────────────

func isValidID(id string) bool {
	if id == "" {
		return false
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}
	}
	return true
}

func isValidLinkType(typ string) bool {
	switch typ {
	case "direct", "store", "webpage":
		return true
	default:
		return false
	}
}

func isValidURL(raw string) bool {
	if raw == "" {
		return false
	}
	return strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://")
}
