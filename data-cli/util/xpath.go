package util

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

// FetchXPath fetches a URL and evaluates an XPath expression,
// returning all matching node inner texts.
func FetchXPath(url, xpath string) ([]string, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: status %d", url, resp.StatusCode)
	}

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	nodes, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		return nil, fmt.Errorf("xpath %q: %w", xpath, err)
	}

	results := make([]string, 0, len(nodes))
	for _, node := range nodes {
		results = append(results, strings.TrimSpace(htmlquery.InnerText(node)))
	}
	return results, nil
}

// FetchXPathAttr fetches a URL, evaluates an XPath expression, and returns
// a specific attribute from all matching nodes.
func FetchXPathAttr(url, xpath, attr string) ([]string, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: status %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	nodes, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		return nil, fmt.Errorf("xpath %q: %w", xpath, err)
	}

	results := make([]string, 0, len(nodes))
	for _, node := range nodes {
		val := htmlquery.SelectAttr(node, attr)
		results = append(results, strings.TrimSpace(val))
	}
	return results, nil
}
