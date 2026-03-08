package provider

import (
	"fmt"

	"github.com/SssHhhAaaDddOoWww/cinex_/internal/types"
)

func GetTorrents(query string, mediaType string) ([]types.Torrent, error) {

	leet := NewLeetProvider()
	results, err := leet.Search(query)
	if err == nil && len(results) > 0 {
		return results, nil
	}

	return nil, fmt.Errorf("no torrents found for %q on any source", query)
}

func GetTorrentsTV(title string, season int, episode int) ([]types.Torrent, error) {

	query := fmt.Sprintf("%s S%02dE%02d", title, season, episode)
	return GetTorrents(query, "tv")
}
