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
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/360-antivirus"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/360-browser"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/360-safe-guard"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/360-zip"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/7zip"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/adobe-acrobat-reader"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/adobe-after-effects"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/adobe-illustrator"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/adobe-photoshop"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/adobe-premiere-pro"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/alipan"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/anydesk"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/autocad"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/baidu-pinyin"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/baidunetdisk"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/chrome"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/coreldraw"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/dingtalk"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/doubao"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/eclipse"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/edge"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/evernote"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/firefox"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/foxit-reader"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/foxmail"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/git"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/github"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/huorong"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/iflytek-pinyin"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/intellij-idea"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/iqiyi"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/jianying"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/kingsoft-antivirus"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/kingsoft-pdf"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/kugou-music"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/microsoft-office"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/microsoft-pinyin"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/neteasecloudmusic"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/notepad-plus-plus"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/onedrive"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/onenote"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/outlook"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/potplayer"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/powerword"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/pycharm"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/qq"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/qq-browser"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/qqmusic"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/snipaste"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/sogou-browser"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/sogou-pinyin"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/steam"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/sunlogin"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/teamviewer"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/tencent-pc-manager"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/tencent-video"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/tencentmeeting"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/todesk"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/uuremote"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/virtualbox"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/visual-studio"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/vlc"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/vmware-workstation"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/vscode"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/wecom"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/weixin"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/weiyun"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/winrar"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/wps"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/youdao-dict"
	_ "github.com/dezhishen/original-software-hub/data-cli/plugin/youku"

	"github.com/dezhishen/original-software-hub/data-cli/plugin"
	"github.com/mozillazg/go-pinyin"
)

func main() {
	runStartedAt := time.Now()

	// ── CLI flags ──────────────────────────────────────────────────────────────
	outDir := flag.String("out", "../frontend/data/json", "Output directory for index.json, software-list.json and versions/")
	pluginsArg := flag.String("plugins", "all", "Plugins to run: all or comma-separated names (e.g. weixin,qq)")
	concurrency := flag.Int("concurrency", 3, "Maximum number of plugins to run concurrently")
	scheduleOrder := flag.String("schedule-order", "priority", "Plugin scheduling order: input,alpha,priority")
	skipUnchanged := flag.Bool("skip-unchanged", true, "Skip write/icon steps when software versions are unchanged from previous output")
	flag.Parse()
	log.Printf("[run] start outDir=%s plugins=%s concurrency=%d order=%s skipUnchanged=%t",
		*outDir, *pluginsArg, *concurrency, *scheduleOrder, *skipUnchanged)

	// ── Directories ────────────────────────────────────────────────────────────
	jsonDir := filepath.Join(*outDir, "versions")
	iconsDir := filepath.Join(filepath.Dir(filepath.Dir(*outDir)), "assets", "software-icons")
	for _, dir := range []string{jsonDir, iconsDir} {
		if err := ensureDir(dir); err != nil {
			log.Fatalf("ensure dir %s: %v", dir, err)
		}
	}

	// ── Previous state ─────────────────────────────────────────────────────────
	prevVersions, prevList, prevIndex := loadPreviousState(*outDir, jsonDir)
	prevItems := make(map[string]plugin.SoftwareItem, len(prevList.Items))
	for _, item := range prevList.Items {
		if strings.TrimSpace(item.ID) != "" {
			prevItems[item.ID] = item
		}
	}
	previousState := plugin.PreviousState{Versions: prevVersions, Items: prevItems}
	log.Printf("[state] loaded previous versions=%d listItems=%d hasGeneratedAt=%t",
		len(prevVersions), len(prevItems), strings.TrimSpace(prevIndex.Meta.GeneratedAt) != "")

	// ── Plugin selection ───────────────────────────────────────────────────────
	plugins, err := selectPlugins(plugin.All(), *pluginsArg)
	if err != nil {
		log.Fatalf("select plugins: %v", err)
	}
	if len(plugins) == 0 {
		log.Println("no plugins registered, nothing to do")
		return
	}
	log.Printf("[scheduler] plugins=%d concurrency=%d order=%s skipUnchanged=%t",
		len(plugins), *concurrency, strings.ToLower(strings.TrimSpace(*scheduleOrder)), *skipUnchanged)

	// ── Fetch & process ────────────────────────────────────────────────────────
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	iconDL := newIconDownloader(iconsDir, "./assets/software-icons")
	fetchResults := fetchPluginsConcurrently(plugins, *concurrency, *scheduleOrder)
	listItems, stats := processFetchResults(fetchResults, previousState, *skipUnchanged, jsonDir, updatedAt, iconDL)

	// ── Sort list items ────────────────────────────────────────────────────────
	sort.SliceStable(listItems, func(i, j int) bool {
		li, lj := sortKeyForItem(listItems[i]), sortKeyForItem(listItems[j])
		if li == lj {
			return listItems[i].ID < listItems[j].ID
		}
		return li < lj
	})

	// ── Timestamps: preserve previous values when nothing changed ──────────────
	listChanged := !reflect.DeepEqual(prevList.Items, listItems)
	dataChanged := stats.writtenCount > 0 || listChanged
	if !dataChanged && strings.TrimSpace(prevList.UpdatedAt) != "" {
		updatedAt = strings.TrimSpace(prevList.UpdatedAt)
	}
	generatedAt := updatedAt
	if !dataChanged && strings.TrimSpace(prevIndex.Meta.GeneratedAt) != "" {
		generatedAt = strings.TrimSpace(prevIndex.Meta.GeneratedAt)
	}

	// ── Write output files ─────────────────────────────────────────────────────
	if err := writeJSON(filepath.Join(*outDir, "software-list.json"), plugin.SoftwareListPayload{
		UpdatedAt: updatedAt,
		Items:     listItems,
	}); err != nil {
		log.Fatalf("write software-list.json: %v", err)
	}
	if err := writeJSON(filepath.Join(*outDir, "index.json"), plugin.IndexPayload{
		Meta:         plugin.Meta{Version: "1.0.0", GeneratedAt: generatedAt, Generator: "data-cli"},
		SoftwareList: plugin.Source{Mode: "json", Path: "software-list.json", TimeoutMs: 8000},
	}); err != nil {
		log.Fatalf("write index.json: %v", err)
	}

	// ── Summary ────────────────────────────────────────────────────────────────
	log.Printf("[summary] duration=%s selectedPlugins=%d pluginFetchErrors=%d processedEntries=%d listItems=%d versionsChanged=%d versionsSkipped=%d skipUnsupported=%d versionsWritten=%d versionWriteErrors=%d iconDownloadSuccess=%d iconDownloadErrors=%d listChanged=%t dataChanged=%t",
		time.Since(runStartedAt), len(plugins),
		stats.pluginFetchErrors, stats.totalEntries, len(listItems),
		stats.changedCount, stats.unchangedCount, stats.skipUnsupportedCount,
		stats.writtenCount, stats.versionWriteErrors,
		stats.iconDownloadSuccess, stats.iconDownloadErrors,
		listChanged, dataChanged)
	fmt.Println("Done.")
}

// runStats aggregates counters from all plugin fetch results.
type runStats struct {
	pluginFetchErrors    int
	changedCount         int
	unchangedCount       int
	skipUnsupportedCount int
	writtenCount         int
	versionWriteErrors   int
	iconDownloadSuccess  int
	iconDownloadErrors   int
	totalEntries         int
}

// processFetchResults processes each plugin's fetch outcome: persists changed version JSON,
// resolves icons, and builds the flat SoftwareItem list for software-list.json.
func processFetchResults(
	results []pluginFetchResult,
	previousState plugin.PreviousState,
	skipUnchanged bool,
	jsonDir, updatedAt string,
	iconDL *iconDownloader,
) ([]plugin.SoftwareItem, runStats) {
	var stats runStats
	listItems := make([]plugin.SoftwareItem, 0, len(results))

	for _, result := range results {
		p := result.Plugin
		if result.Err != nil {
			stats.pluginFetchErrors++
			log.Printf("[%s] Fetch error: %v", p.Name(), result.Err)
			continue
		}

		var pItems, pChanged, pUnchanged, pSkipUnsupported, pWritten, pWriteErrors, pIconOK, pIconErr int

		for _, fetched := range result.Items {
			decision, err := resolveDataByVersionDecision(previousState, fetched, skipUnchanged)
			if err != nil {
				log.Printf("[%s] resolve: %v", p.Name(), err)
				continue
			}
			entry := decision.Result
			softwareID := entry.Item.ID
			prevItem, hasPrevItem := previousState.Items[softwareID]

			if !decision.SkipSupported {
				stats.skipUnsupportedCount++
				pSkipUnsupported++
			}

			if decision.Changed {
				stats.changedCount++
				pChanged++
				log.Printf("[%s/%s] decision=changed reason=%s skipSupport=%t", p.Name(), softwareID, decision.Reason, decision.SkipSupported)
				payload := plugin.PlatformPayload{
					SoftwareID: softwareID,
					UpdatedAt:  updatedAt,
					Platforms:  flattenVersionsToPlatforms(entry.Versions),
				}
				if err := writeJSON(filepath.Join(jsonDir, softwareID+".json"), payload); err != nil {
					stats.versionWriteErrors++
					pWriteErrors++
					log.Printf("[%s/%s] write json: %v", p.Name(), softwareID, err)
					continue
				}
				stats.writtenCount++
				pWritten++
			} else {
				stats.unchangedCount++
				pUnchanged++
				log.Printf("[%s/%s] decision=skipped reason=%s skipSupport=%t", p.Name(), softwareID, decision.Reason, decision.SkipSupported)
			}

			item := entry.Item
			if decision.Changed {
				if localIcon, err := iconDL.Download(softwareID, item.Icon); err != nil {
					stats.iconDownloadErrors++
					pIconErr++
					log.Printf("[%s/%s] download icon: %v", p.Name(), softwareID, err)
					if hasPrevItem && strings.TrimSpace(prevItem.Icon) != "" {
						item.Icon = prevItem.Icon
					}
				} else {
					item.Icon = localIcon
					stats.iconDownloadSuccess++
					pIconOK++
				}
			} else if hasPrevItem && strings.TrimSpace(prevItem.Icon) != "" {
				item.Icon = prevItem.Icon
			}

			item.Pinyin = buildSearchPinyin(item.Name)
			item.Source = plugin.Source{Mode: "json", Path: "versions/" + softwareID + ".json", TimeoutMs: 8000}
			listItems = append(listItems, item)
			pItems++
			stats.totalEntries++
		}

		log.Printf("[plugin-summary] name=%s mode=fetch items=%d changed=%d skipped=%d skipUnsupported=%d written=%d writeErrors=%d iconDownloadSuccess=%d iconDownloadErrors=%d",
			p.Name(), pItems, pChanged, pUnchanged, pSkipUnsupported, pWritten, pWriteErrors, pIconOK, pIconErr)
	}

	return listItems, stats
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
		enabled := make([]plugin.Plugin, 0, len(all))
		for _, p := range all {
			if p.Disabled() {
				log.Printf("[scheduler] skip disabled plugin: %s", p.Name())
				continue
			}
			enabled = append(enabled, p)
		}
		return enabled, nil
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

func loadPreviousState(outDir, versionsDir string) (map[string]plugin.PlatformPayload, plugin.SoftwareListPayload, plugin.IndexPayload) {
	versionMap := map[string]plugin.PlatformPayload{}
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
		var payload plugin.PlatformPayload
		if err := json.Unmarshal(data, &payload); err == nil && strings.TrimSpace(payload.SoftwareID) != "" {
			versionMap[payload.SoftwareID] = payload
			continue
		}

		// Backward compatibility for legacy payload with versions[]
		var legacy struct {
			SoftwareID string           `json:"softwareId"`
			UpdatedAt  string           `json:"updatedAt"`
			Versions   []plugin.Version `json:"versions"`
		}
		if err := json.Unmarshal(data, &legacy); err != nil {
			continue
		}
		if strings.TrimSpace(legacy.SoftwareID) == "" {
			continue
		}
		versionMap[legacy.SoftwareID] = plugin.PlatformPayload{
			SoftwareID: legacy.SoftwareID,
			UpdatedAt:  legacy.UpdatedAt,
			Platforms:  flattenVersionsToPlatforms(legacy.Versions),
		}
	}

	return versionMap, listPayload, indexPayload
}

func versionsEqual(a, b []plugin.Version) bool {
	return reflect.DeepEqual(a, b)
}

func platformsEqual(a, b []plugin.PlatformRelease) bool {
	return reflect.DeepEqual(a, b)
}

func flattenVersionsToPlatforms(versions []plugin.Version) []plugin.PlatformRelease {
	merged := map[string]plugin.PlatformRelease{}
	order := make([]string, 0, len(versions))

	for _, v := range versions {
		for _, p := range v.Platforms {
			platformName := strings.TrimSpace(p.Platform)
			if platformName == "" {
				platformName = "其他"
			}

			entry, exists := merged[platformName]
			if !exists {
				entry = plugin.PlatformRelease{Platform: platformName}
				order = append(order, platformName)
			}

			if strings.TrimSpace(entry.Version) == "" || strings.TrimSpace(p.ReleaseDate) >= strings.TrimSpace(entry.ReleaseDate) {
				if strings.TrimSpace(p.Version) != "" {
					entry.Version = strings.TrimSpace(p.Version)
				}
				if strings.TrimSpace(p.ReleaseDate) != "" {
					entry.ReleaseDate = strings.TrimSpace(p.ReleaseDate)
				}
				if strings.TrimSpace(p.OfficialURL) != "" {
					entry.OfficialURL = strings.TrimSpace(p.OfficialURL)
				}
			}

			for _, pkg := range p.Packages {
				if len(pkg.Links) == 0 {
					continue
				}
				entry.Packages = append(entry.Packages, pkg)
			}

			merged[platformName] = entry
		}
	}

	result := make([]plugin.PlatformRelease, 0, len(order))
	for _, platformName := range order {
		entry := merged[platformName]
		result = append(result, dedupePlatformRelease(entry))
	}
	return result
}

func dedupePlatformRelease(entry plugin.PlatformRelease) plugin.PlatformRelease {
	type packageKey struct {
		arch string
		url  string
	}
	seen := map[packageKey]struct{}{}
	packages := make([]plugin.PlatformPackage, 0, len(entry.Packages))

	for _, pkg := range entry.Packages {
		arch := strings.TrimSpace(pkg.Architecture)
		links := make([]plugin.Link, 0, len(pkg.Links))
		for _, link := range pkg.Links {
			u := strings.TrimSpace(link.URL)
			if u == "" {
				continue
			}
			k := packageKey{arch: arch, url: u}
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			links = append(links, link)
		}
		if len(links) == 0 {
			continue
		}
		packages = append(packages, plugin.PlatformPackage{Architecture: pkg.Architecture, Links: links})
	}

	entry.Packages = packages
	return entry
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
	oldPlatforms := oldPayload.Platforms
	entryPlatforms := flattenVersionsToPlatforms(entry.Versions)
	if !hasOldPayload {
		return versionDecision{Changed: true, Result: entry, SkipSupported: false, Reason: "no-previous-version"}, nil
	}

	if !platformsEqual(oldPlatforms, entryPlatforms) {
		return versionDecision{Changed: true, Result: entry, SkipSupported: true, Reason: "version-changed"}, nil
	}

	oldItem, hasOldItem := previous.Items[softwareID]
	if hasOldItem {
		return versionDecision{Changed: false, Result: plugin.SoftwareData{Item: oldItem, Versions: entry.Versions}, SkipSupported: true, Reason: "version-unchanged"}, nil
	}

	return versionDecision{Changed: true, Result: entry, SkipSupported: false, Reason: "previous-item-missing"}, nil
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
