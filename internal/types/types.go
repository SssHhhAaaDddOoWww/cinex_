package types

type Show struct {
	ID   int    `json:"id"`
	Name string `json:"title"`
	Type string `json:"media_type"`
	Year string `json:"release_date"`
}
type Episode struct {
	ID     int    `json:"id"`
	Title  string `json:"name"`
	Number int    `json:"episode_number"`
}
type Season struct {
	Name   string `json:"name"`
	Number int    `json:"season_number"`
}

type Torrent struct {
	Title   string
	Magnet  string
	Seeders int
	Size    string
	Quality string
}

type Provider interface {
	Name() string
	Search(query string) ([]Torrent, error)
}
