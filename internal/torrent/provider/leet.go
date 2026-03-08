package provider

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/SssHhhAaaDddOoWww/cinex_/internal/types"
)

type LeetProvider struct{}

func NewLeetProvider() *LeetProvider {
	return &LeetProvider{}
}

func (l *LeetProvider) Name() string {
	return "1337x"
}

func (l *LeetProvider) Search(query string) ([]types.Torrent, error) {

	key := os.Getenv("LeetBase")
	if key == "" {
		key = "https://1337x.to"
	}
	searchURL := fmt.Sprintf("%s/search/%s/1/", key, url.QueryEscape(query))

	body, err := fetchPage(searchURL)
	if err != nil {
		return nil, err
	}

	// extract torrent page links from search results
	links := extractLinks(body)
	if len(links) == 0 {
		return nil, fmt.Errorf("no results found on 1337x")
	}

	// limit to first 10 results
	if len(links) > 10 {
		links = links[:10]
	}

	var torrents []types.Torrent
	for _, link := range links {
		torrent, err := fetchTorrentPage(link)
		if err != nil {
			continue
		}
		torrents = append(torrents, torrent)
	}

	return torrents, nil
}

// fetchPage fetches a URL and returns the body as string
func fetchPage(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// 1337x needs a user agent or it blocks the request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// extractLinks pulls torrent detail page URLs from the search results HTML
func extractLinks(html string) []string {
	re := regexp.MustCompile(`href="(/torrent/\d+/[^"]+)"`)
	key := os.Getenv("LeetBase")
	if key == "" {
		key = "https://1337x.to"
	}
	matches := re.FindAllStringSubmatch(html, -1)

	seen := map[string]bool{}
	var links []string
	for _, match := range matches {
		path := match[1]
		if !seen[path] {
			seen[path] = true
			links = append(links, key+path)
		}
	}

	return links
}

// fetchTorrentPage visits a torrent detail page and extracts magnet + metadata
func fetchTorrentPage(pageURL string) (types.Torrent, error) {
	body, err := fetchPage(pageURL)
	if err != nil {
		return types.Torrent{}, err
	}

	magnet := extractMagnet(body)
	if magnet == "" {
		return types.Torrent{}, fmt.Errorf("no magnet found on page")
	}

	return types.Torrent{
		Title:   extractTitle(body),
		Magnet:  magnet,
		Seeders: extractSeeders(body),
		Size:    extractSize(body),
	}, nil
}

func extractMagnet(html string) string {
	re := regexp.MustCompile(`href="(magnet:\?[^"]+)"`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func extractTitle(html string) string {
	re := regexp.MustCompile(`<title>([^<]+)</title>`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		// 1337x titles look like "Suits S01 Download - 1337x"
		title := strings.TrimSuffix(match[1], " Download - 1337x")
		return strings.TrimSpace(title)
	}
	return "Unknown"
}

func extractSeeders(html string) int {
	re := regexp.MustCompile(`<span class="seeds">(\d+)</span>`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		var seeds int
		fmt.Sscanf(match[1], "%d", &seeds)
		return seeds
	}
	return 0
}

func extractSize(html string) string {
	re := regexp.MustCompile(`<span class="size"[^>]*>([^<]+)<`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return "Unknown"
}
