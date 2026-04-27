package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
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
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/firefox"
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
	frontendRootDir := filepath.Dir(filepath.Dir(*outDir))
	iconsDir := filepath.Join(frontendRootDir, "assets", "software-icons")
	if err := resetDir(jsonDir); err != nil {
		log.Fatalf("reset json dir: %v", err)
	}
	if err := os.MkdirAll(jsonDir, 0o755); err != nil {
		log.Fatalf("create json dir: %v", err)
	}
	if err := resetDir(iconsDir); err != nil {
		log.Fatalf("reset icons dir: %v", err)
	}
	iconDownloader := newIconDownloader(iconsDir, "./assets/software-icons")

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
			if localIcon, err := iconDownloader.Download(softwareID, item.Icon); err != nil {
				log.Printf("[%s/%s] download icon: %v", p.Name(), softwareID, err)
			} else {
				item.Icon = localIcon
			}
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

const (
	defaultFetchRetries   = 3
	defaultRetryBaseDelay = time.Second
	maxIconFileSizeBytes  = 5 << 20
)

type iconDownloader struct {
	client    *http.Client
	outputDir string
	relPrefix string
	cache     map[string]string
}

func newIconDownloader(outputDir, relPrefix string) *iconDownloader {
	return &iconDownloader{
		client:    &http.Client{Timeout: 12 * time.Second},
		outputDir: outputDir,
		relPrefix: strings.TrimRight(relPrefix, "/"),
		cache:     make(map[string]string),
	}
}

func (d *iconDownloader) Download(softwareID, iconRef string) (string, error) {
	iconRef = strings.TrimSpace(iconRef)
	if iconRef == "" {
		return "", nil
	}
	if !isRemoteHTTPURL(iconRef) {
		return iconRef, nil
	}

	if cached, ok := d.cache[iconRef]; ok {
		return cached, nil
	}

	req, err := http.NewRequest(http.MethodGet, iconRef, nil)
	if err != nil {
		return iconRef, err
	}
	req.Header.Set("User-Agent", "original-software-hub-data-cli/1.0")
	req.Header.Set("Accept", "image/*,*/*;q=0.8")

	resp, err := d.client.Do(req)
	if err != nil {
		return iconRef, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return iconRef, fmt.Errorf("http %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxIconFileSizeBytes+1))
	if err != nil {
		return iconRef, err
	}
	if len(body) > maxIconFileSizeBytes {
		return iconRef, fmt.Errorf("icon too large (> %d bytes)", maxIconFileSizeBytes)
	}

	ext := chooseIconExtension(iconRef, resp.Header.Get("Content-Type"), body)
	filename := sanitizeFileStem(softwareID)
	if filename == "" {
		filename = fmt.Sprintf("icon-%d", time.Now().UnixNano())
	}
	fileNameWithExt := filename + ext
	filePath := filepath.Join(d.outputDir, fileNameWithExt)
	if err := os.WriteFile(filePath, body, 0o644); err != nil {
		return iconRef, err
	}

	localPath := d.relPrefix + "/" + fileNameWithExt
	d.cache[iconRef] = localPath
	return localPath, nil
}

func chooseIconExtension(iconURL, contentType string, body []byte) string {
	if ext := iconExtensionFromContentType(contentType); ext != "" {
		return ext
	}
	if looksLikeSVG(body) {
		return ".svg"
	}
	if ext := iconExtensionFromURL(iconURL); ext != "" {
		return ext
	}
	return ".png"
}

func iconExtensionFromContentType(raw string) string {
	mediaType, _, err := mime.ParseMediaType(strings.TrimSpace(raw))
	if err != nil {
		return ""
	}
	mediaType = strings.ToLower(mediaType)
	switch mediaType {
	case "image/svg+xml":
		return ".svg"
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	case "image/x-icon", "image/vnd.microsoft.icon":
		return ".ico"
	default:
		return ""
	}
}

func iconExtensionFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	ext := strings.ToLower(path.Ext(u.Path))
	switch ext {
	case ".svg", ".png", ".jpg", ".webp", ".gif", ".ico":
		return ext
	case ".jpeg":
		return ".jpg"
	default:
		return ""
	}
}

func looksLikeSVG(body []byte) bool {
	limit := len(body)
	if limit > 512 {
		limit = 512
	}
	head := strings.ToLower(string(body[:limit]))
	return strings.Contains(head, "<svg")
}

func sanitizeFileStem(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}

	var b strings.Builder
	lastDash := false
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func isRemoteHTTPURL(raw string) bool {
	if raw == "" {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	return scheme == "http" || scheme == "https"
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
			items, err := fetchWithRetry(job.Plugin, defaultFetchRetries, defaultRetryBaseDelay)
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

func fetchWithRetry(p plugin.Plugin, maxAttempts int, baseDelay time.Duration) ([]plugin.SoftwareData, error) {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	if baseDelay <= 0 {
		baseDelay = time.Second
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		items, err := p.Fetch()
		if err == nil {
			return items, nil
		}

		lastErr = err
		if attempt == maxAttempts {
			break
		}

		delay := retryDelay(baseDelay, attempt)
		log.Printf("[%s] fetch failed (attempt %d/%d): %v; retry in %s", p.Name(), attempt, maxAttempts, err, delay)
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("after %d attempts: %w", maxAttempts, lastErr)
}

func retryDelay(base time.Duration, attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	factor := math.Pow(2, float64(attempt-1))
	return time.Duration(float64(base) * factor)
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
