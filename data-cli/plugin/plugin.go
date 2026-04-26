// Package plugin defines the Plugin interface and shared data types used
// by all data-cli plugins.
package plugin

import "sync"

// ── Shared types ──────────────────────────────────────────────────────────────

// Source describes how to fetch a remote JSON/JSONP file.
type Source struct {
	Mode          string `json:"mode"`
	Path          string `json:"path"`
	CallbackParam string `json:"callbackParam,omitempty"`
	TimeoutMs     int    `json:"timeoutMs,omitempty"`
}

// Meta holds index-level metadata.
type Meta struct {
	Version     string `json:"version"`
	GeneratedAt string `json:"generatedAt"`
	Generator   string `json:"generator"`
}

// IndexPayload is the root index.json / index.js payload.
type IndexPayload struct {
	Meta         Meta   `json:"meta"`
	SoftwareList Source `json:"softwareList"`
}

// SoftwareItem represents one entry in the software list.
type SoftwareItem struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Pinyin          string   `json:"pinyin,omitempty"`
	Icon            string   `json:"icon,omitempty"`
	Description     string   `json:"description"`
	Organization    string   `json:"organization"`
	OfficialWebsite string   `json:"officialWebsite"`
	Tags            []string `json:"tags,omitempty"`
	Source          Source   `json:"source"`
}

// SoftwareListPayload is the software-list.json payload.
type SoftwareListPayload struct {
	UpdatedAt string         `json:"updatedAt"`
	Items     []SoftwareItem `json:"items"`
}

// Link is a single download link.
type Link struct {
	Type  string `json:"type"`  // "direct"
	Label string `json:"label"` // display text
	URL   string `json:"url"`
}

// Variant is one (architecture, platform) combination.
type Variant struct {
	Architecture string `json:"architecture"`
	Platform     string `json:"platform"`
	Links        []Link `json:"links"`
}

// Version represents one release of a software.
type Version struct {
	Version     string    `json:"version"`
	ReleaseDate string    `json:"releaseDate"`
	OfficialURL string    `json:"officialUrl"`
	Variants    []Variant `json:"variants"`
}

// VersionPayload is the versions/<id>.json payload.
type VersionPayload struct {
	SoftwareID string    `json:"softwareId"`
	UpdatedAt  string    `json:"updatedAt"`
	Versions   []Version `json:"versions"`
}

// SoftwareData is the in-memory full data returned by plugins.
// The main program is responsible for transforming it into output files.
type SoftwareData struct {
	Item     SoftwareItem
	Versions []Version
}

// ── Plugin interface ──────────────────────────────────────────────────────────

// Plugin is implemented by each data source plugin package.
type Plugin interface {
	// Name returns the plugin name for logging (e.g. "chrome", "github").
	Name() string
	// Fetch returns one or more software data items.
	Fetch() ([]SoftwareData, error)
}

// ── Registry ──────────────────────────────────────────────────────────────────

var (
	mu       sync.Mutex
	registry []Plugin
)

// Register adds a plugin to the global registry. Call from plugin init().
func Register(p Plugin) {
	mu.Lock()
	defer mu.Unlock()
	registry = append(registry, p)
}

// All returns a snapshot of all registered plugins.
func All() []Plugin {
	mu.Lock()
	defer mu.Unlock()
	out := make([]Plugin, len(registry))
	copy(out, registry)
	return out
}
