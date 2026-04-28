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
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	// Blank-import each plugin to trigger its init() registration.
	// Uncomment or add plugins here to include them in the build.
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/alipan"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/baidunetdisk"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/chrome"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/dingtalk"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/doubao"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/firefox"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/github"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/huorong"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/neteasecloudmusic"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/qq"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/qqmusic"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/steam"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/tencentmeeting"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/todesk"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/uuremote"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/wecom"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/weixin"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/wps"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
	"github.com/mozillazg/go-pinyin"
)

func main() {
	runStartedAt := time.Now()
	outDir := flag.String("out", "../frontend/data/json", "Output directory for index.json, software-list.json and versions/")
	pluginsArg := flag.String("plugins", "all", "Plugins to run: all or comma-separated names (e.g. weixin,qq)")
	concurrency := flag.Int("concurrency", 3, "Maximum number of plugins to run concurrently")
	scheduleOrder := flag.String("schedule-order", "priority", "Plugin scheduling order: input,alpha,priority")
	skipUnchanged := flag.Bool("skip-unchanged", true, "Skip write/icon steps when software versions are unchanged from previous output")
	flag.Parse()
	log.Printf("[run] start outDir=%s plugins=%s concurrency=%d order=%s skipUnchanged=%t", *outDir, *pluginsArg, *concurrency, *scheduleOrder, *skipUnchanged)

	jsonDir := filepath.Join(*outDir, "versions")
	frontendRootDir := filepath.Dir(filepath.Dir(*outDir))
	iconsDir := filepath.Join(frontendRootDir, "assets", "software-icons")
	prevVersions, prevList, prevIndex := loadPreviousState(*outDir, jsonDir)
	prevItems := make(map[string]plugin.SoftwareItem, len(prevList.Items))
	for _, item := range prevList.Items {
		if strings.TrimSpace(item.ID) == "" {
			continue
		}
		prevItems[item.ID] = item
	}
	previousState := plugin.PreviousState{
		Versions: prevVersions,
		Items:    prevItems,
	}
	log.Printf("[state] loaded previous versions=%d listItems=%d hasGeneratedAt=%t", len(prevVersions), len(prevItems), strings.TrimSpace(prevIndex.Meta.GeneratedAt) != "")

	if err := ensureDir(jsonDir); err != nil {
		log.Fatalf("ensure json dir: %v", err)
	}
	if err := ensureDir(iconsDir); err != nil {
		log.Fatalf("ensure icons dir: %v", err)
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
	log.Printf("[scheduler] plugins=%d concurrency=%d order=%s skipUnchanged=%t", len(plugins), *concurrency, strings.ToLower(strings.TrimSpace(*scheduleOrder)), *skipUnchanged)

	nowUTC := time.Now().UTC()
	updatedAt := nowUTC.Format(time.RFC3339)

	listItems := make([]plugin.SoftwareItem, 0, len(plugins))
	fetchResults := fetchPluginsConcurrently(plugins, *concurrency, *scheduleOrder)
	changedCount := 0
	unchangedCount := 0
	skipUnsupportedCount := 0
	writtenCount := 0
	totalEntries := 0
	pluginFetchErrors := 0
	versionWriteErrors := 0
	iconDownloadSuccess := 0
	iconDownloadErrors := 0
	for _, result := range fetchResults {
		p := result.Plugin
		if result.Err != nil {
			pluginFetchErrors++
			log.Printf("[%s] Fetch error: %v", p.Name(), result.Err)
			continue
		}

		pluginItemCount := 0
		pluginChangedCount := 0
		pluginUnchangedCount := 0
		pluginSkipUnsupportedCount := 0
		pluginWrittenCount := 0
		pluginVersionWriteErrors := 0
		pluginIconDownloadSuccess := 0
		pluginIconDownloadErrors := 0

		for _, fetched := range result.Items {
			decision, err := resolveDataByVersionDecision(previousState, fetched, *skipUnchanged)
			if err != nil {
				log.Printf("[%s] resolve data by version: %v", p.Name(), err)
				continue
			}
			changed := decision.Changed
			entry := decision.Result

			softwareID := entry.Item.ID
			prevItem, hasPrevItem := prevItems[softwareID]
			if !decision.SkipSupported {
				skipUnsupportedCount++
				pluginSkipUnsupportedCount++
			}
			if !changed {
				unchangedCount++
				pluginUnchangedCount++
				log.Printf("[%s/%s] decision=skipped reason=%s skipSupport=%t", p.Name(), softwareID, decision.Reason, decision.SkipSupported)
			} else {
				changedCount++
				pluginChangedCount++
				log.Printf("[%s/%s] decision=changed reason=%s skipSupport=%t", p.Name(), softwareID, decision.Reason, decision.SkipSupported)
				versionPayload := plugin.VersionPayload{
					SoftwareID: softwareID,
					UpdatedAt:  updatedAt,
					Versions:   entry.Versions,
				}
				if err := writeJSON(filepath.Join(jsonDir, softwareID+".json"), versionPayload); err != nil {
					versionWriteErrors++
					pluginVersionWriteErrors++
					log.Printf("[%s/%s] write json: %v", p.Name(), softwareID, err)
					continue
				}
				writtenCount++
				pluginWrittenCount++
			}

			item := entry.Item
			if !changed {
				if hasPrevItem && strings.TrimSpace(prevItem.Icon) != "" {
					item.Icon = prevItem.Icon
				}
			} else {
				if localIcon, err := iconDownloader.Download(softwareID, item.Icon); err != nil {
					iconDownloadErrors++
					pluginIconDownloadErrors++
					log.Printf("[%s/%s] download icon: %v", p.Name(), softwareID, err)
					if hasPrevItem && strings.TrimSpace(prevItem.Icon) != "" {
						item.Icon = prevItem.Icon
					}
				} else {
					item.Icon = localIcon
					iconDownloadSuccess++
					pluginIconDownloadSuccess++
				}
			}

			item.Pinyin = buildSearchPinyin(item.Name)
			item.Source = plugin.Source{
				Mode:      "json",
				Path:      "versions/" + softwareID + ".json",
				TimeoutMs: 8000,
			}
			listItems = append(listItems, item)
			pluginItemCount++
			totalEntries++
		}

		log.Printf("[plugin-summary] name=%s mode=fetch items=%d changed=%d skipped=%d skipUnsupported=%d written=%d writeErrors=%d iconDownloadSuccess=%d iconDownloadErrors=%d", p.Name(), pluginItemCount, pluginChangedCount, pluginUnchangedCount, pluginSkipUnsupportedCount, pluginWrittenCount, pluginVersionWriteErrors, pluginIconDownloadSuccess, pluginIconDownloadErrors)
	}

	sort.SliceStable(listItems, func(i, j int) bool {
		li := sortKeyForItem(listItems[i])
		lj := sortKeyForItem(listItems[j])
		if li == lj {
			return listItems[i].ID < listItems[j].ID
		}
		return li < lj
	})

	listChanged := !reflect.DeepEqual(prevList.Items, listItems)
	dataChanged := writtenCount > 0 || listChanged
	if !dataChanged && strings.TrimSpace(prevList.UpdatedAt) != "" {
		updatedAt = strings.TrimSpace(prevList.UpdatedAt)
	}
	generatedAt := updatedAt
	if !dataChanged && strings.TrimSpace(prevIndex.Meta.GeneratedAt) != "" {
		generatedAt = strings.TrimSpace(prevIndex.Meta.GeneratedAt)
	}

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
			GeneratedAt: generatedAt,
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

	runDuration := time.Since(runStartedAt)
	log.Printf("[summary] duration=%s selectedPlugins=%d pluginFetchErrors=%d processedEntries=%d listItems=%d versionsChanged=%d versionsSkipped=%d skipUnsupported=%d versionsWritten=%d versionWriteErrors=%d iconDownloadSuccess=%d iconDownloadErrors=%d listChanged=%t dataChanged=%t", runDuration, len(plugins), pluginFetchErrors, totalEntries, len(listItems), changedCount, unchangedCount, skipUnsupportedCount, writtenCount, versionWriteErrors, iconDownloadSuccess, iconDownloadErrors, listChanged, dataChanged)
	fmt.Println("Done.")
}

type pluginJob struct {
	Index  int
	Plugin plugin.Plugin
}

type pluginFetchResult struct {
	Index    int
	Plugin   plugin.Plugin
	Items    []plugin.SoftwareData
	Err      error
	Attempts int
	Duration time.Duration
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

func fetchPluginsConcurrently(plugins []plugin.Plugin, maxConcurrency int, order string) []pluginFetchResult {
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
	jobList := buildPluginJobs(plugins, order)
	queued := make([]string, 0, len(jobList))
	for _, j := range jobList {
		queued = append(queued, j.Plugin.Name())
	}
	log.Printf("[scheduler] queue=%s", strings.Join(queued, ","))

	var wg sync.WaitGroup
	worker := func(workerID int) {
		defer wg.Done()
		for job := range jobs {
			start := time.Now()
			log.Printf("[scheduler] start plugin=%s order=%d worker=%d", job.Plugin.Name(), job.Index, workerID)
			outcome, attempts, err := fetchWithRetry(job.Plugin, defaultFetchRetries, defaultRetryBaseDelay)
			duration := time.Since(start)
			if err != nil {
				log.Printf("[scheduler] done plugin=%s status=failed attempts=%d duration=%s err=%v", job.Plugin.Name(), attempts, duration, err)
			} else {
				log.Printf("[scheduler] done plugin=%s status=ok attempts=%d duration=%s items=%d", job.Plugin.Name(), attempts, duration, len(outcome.Items))
			}
			results <- pluginFetchResult{
				Index:    job.Index,
				Plugin:   job.Plugin,
				Items:    outcome.Items,
				Err:      err,
				Attempts: attempts,
				Duration: duration,
			}
		}
	}

	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go worker(i + 1)
	}

	for _, job := range jobList {
		jobs <- job
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

type pluginFetchOutcome struct {
	Items []plugin.SoftwareData
}

func fetchWithRetry(p plugin.Plugin, maxAttempts int, baseDelay time.Duration) (pluginFetchOutcome, int, error) {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	if baseDelay <= 0 {
		baseDelay = time.Second
	}

	var lastErr error
	log.Printf("[%s] fetch mode=fetch", p.Name())
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		items, err := p.Fetch()
		if err == nil {
			return pluginFetchOutcome{Items: items}, attempt, nil
		}

		lastErr = err
		if attempt == maxAttempts {
			break
		}

		delay := retryDelay(baseDelay, attempt)
		log.Printf("[%s] fetch failed (attempt %d/%d): %v; retry in %s", p.Name(), attempt, maxAttempts, err, delay)
		time.Sleep(delay)
	}

	return pluginFetchOutcome{}, maxAttempts, fmt.Errorf("after %d attempts: %w", maxAttempts, lastErr)
}

func buildPluginJobs(plugins []plugin.Plugin, order string) []pluginJob {
	jobs := make([]pluginJob, 0, len(plugins))
	for i, p := range plugins {
		jobs = append(jobs, pluginJob{Index: i, Plugin: p})
	}

	normalized := strings.ToLower(strings.TrimSpace(order))
	switch normalized {
	case "alpha":
		sort.SliceStable(jobs, func(i, j int) bool {
			return strings.ToLower(jobs[i].Plugin.Name()) < strings.ToLower(jobs[j].Plugin.Name())
		})
	case "priority":
		sort.SliceStable(jobs, func(i, j int) bool {
			ri := pluginPriorityRank(jobs[i].Plugin.Name())
			rj := pluginPriorityRank(jobs[j].Plugin.Name())
			if ri == rj {
				return strings.ToLower(jobs[i].Plugin.Name()) < strings.ToLower(jobs[j].Plugin.Name())
			}
			return ri < rj
		})
	default:
		// input order
	}

	for i := range jobs {
		jobs[i].Index = i
	}
	return jobs
}

func pluginPriorityRank(name string) int {
	name = strings.ToLower(strings.TrimSpace(name))
	// Network-heavy or historically flaky sources are scheduled first
	// so retries can overlap with other workers and reduce total wall time.
	priority := map[string]int{
		"github":         10,
		"chrome":         20,
		"firefox":        30,
		"wps":            40,
		"todesk":         50,
		"uuremote":       60,
		"tencentmeeting": 70,
	}
	if rank, ok := priority[name]; ok {
		return rank
	}
	return 100
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

func ensureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func loadPreviousState(outDir, versionsDir string) (map[string]plugin.VersionPayload, plugin.SoftwareListPayload, plugin.IndexPayload) {
	versionMap := map[string]plugin.VersionPayload{}
	listPayload := plugin.SoftwareListPayload{}
	indexPayload := plugin.IndexPayload{}

	listPath := filepath.Join(outDir, "software-list.json")
	if data, err := os.ReadFile(listPath); err == nil {
		if err := json.Unmarshal(data, &listPayload); err != nil {
			listPayload = plugin.SoftwareListPayload{}
		}
	}

	indexPath := filepath.Join(outDir, "index.json")
	if data, err := os.ReadFile(indexPath); err == nil {
		if err := json.Unmarshal(data, &indexPayload); err != nil {
			indexPayload = plugin.IndexPayload{}
		}
	}

	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		return versionMap, listPayload, indexPayload
	}
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(strings.ToLower(ent.Name()), ".json") {
			continue
		}
		filePath := filepath.Join(versionsDir, ent.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		var payload plugin.VersionPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			continue
		}
		if strings.TrimSpace(payload.SoftwareID) == "" {
			continue
		}
		versionMap[payload.SoftwareID] = payload
	}

	return versionMap, listPayload, indexPayload
}

func versionsEqual(a, b []plugin.Version) bool {
	return reflect.DeepEqual(a, b)
}

type versionDecision struct {
	Changed       bool
	Result        plugin.SoftwareData
	SkipSupported bool
	Reason        string
}

func resolveDataByVersionDecision(previous plugin.PreviousState, fetched plugin.SoftwareData, skipUnchanged bool) (versionDecision, error) {
	entry := fetched
	softwareID := strings.TrimSpace(entry.Item.ID)
	if softwareID == "" {
		return versionDecision{}, fmt.Errorf("empty software id")
	}

	if !skipUnchanged {
		return versionDecision{Changed: true, Result: entry, SkipSupported: false, Reason: "skip-disabled"}, nil
	}

	oldPayload, hasOldPayload := previous.Versions[softwareID]
	oldVersions := oldPayload.Versions
	entryVersions := entry.Versions
	if !hasOldPayload {
		return versionDecision{Changed: true, Result: entry, SkipSupported: false, Reason: "no-previous-version"}, nil
	}

	if !versionsEqual(oldVersions, entryVersions) {
		return versionDecision{Changed: true, Result: entry, SkipSupported: true, Reason: "version-changed"}, nil
	}

	oldItem, hasOldItem := previous.Items[softwareID]
	if hasOldItem {
		return versionDecision{Changed: false, Result: plugin.SoftwareData{Item: oldItem, Versions: oldVersions}, SkipSupported: true, Reason: "version-unchanged"}, nil
	}

	return versionDecision{Changed: true, Result: entry, SkipSupported: false, Reason: "previous-item-missing"}, nil
}

// resolveDataByVersion centralizes unchanged/changed decision and previous-data reuse.
// Return values follow: changed, result, err.
func resolveDataByVersion(previous plugin.PreviousState, fetched plugin.SoftwareData, skipUnchanged bool) (bool, plugin.SoftwareData, error) {
	decision, err := resolveDataByVersionDecision(previous, fetched, skipUnchanged)
	if err != nil {
		return false, plugin.SoftwareData{}, err
	}
	return decision.Changed, decision.Result, nil
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
