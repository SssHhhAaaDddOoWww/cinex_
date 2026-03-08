package tui

import (
	"fmt"

	"github.com/SssHhhAaaDddOoWww/cinex_/internal/player"
	"github.com/SssHhhAaaDddOoWww/cinex_/internal/search"
	"github.com/SssHhhAaaDddOoWww/cinex_/internal/torrent/provider"
	"github.com/SssHhhAaaDddOoWww/cinex_/internal/types"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

type model struct {
	screen   screen
	query    string
	shows    []types.Show
	selected types.Show 
	seasons  []types.Season
	season   types.Season 
	episodes []types.Episode
	torrents []types.Torrent
	cursor   int
	err      error
	loading  bool
}

const (
	screenSearch screen = iota
	screenResults
	screenSeasons
	screenEpisodes
	screenTorrents
)

func SearchCMD(query string) tea.Cmd {
	return func() tea.Msg {
		shows, err := search.Search(query)
		return SearchResultsMsg{shows: shows, err: err}
	}
}

func SearchSeason(seriesID int) tea.Cmd {
	return func() tea.Msg {
		seasons, err := search.GetSeasons(seriesID)
		return SeasonResultsMsg{seasons: seasons, err: err}
	}
}

func SearchEp(seriesID, seasonNumber int) tea.Cmd {
	return func() tea.Msg {
		episodes, err := search.GetEpisodes(seriesID, seasonNumber)
		return EpisodeResultsMsg{episodes: episodes, err: err}
	}
}

func fetchTorrentsMovie(title string) tea.Cmd {
	return func() tea.Msg {
		torrents, err := provider.GetTorrents(title, "movie")
		return TorrentResultMsg{torrents, err}
	}
}

func fetchTorrentsEpisode(title string, season int, episode int) tea.Cmd {
	return func() tea.Msg {
		torrents, err := provider.GetTorrentsTV(title, season, episode)
		return TorrentResultMsg{torrents, err}
	}
}

type SearchResultsMsg struct {
	shows []types.Show
	err   error
}

type SeasonResultsMsg struct {
	seasons []types.Season
	err     error
}

type EpisodeResultsMsg struct {
	episodes []types.Episode
	err      error
}

type TorrentResultMsg struct {
	torrents []types.Torrent
	err      error
}



func InitialModel() model {
	return model{screen: screenSearch}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch m.screen {

		case screenSearch:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "backspace":
				if len(m.query) > 0 {
					m.query = m.query[:len(m.query)-1]
				}
			case "enter":
				if m.query != "" {
					m.loading = true
					m.err = nil
					return m, SearchCMD(m.query)
				}
			default:
				if len(msg.String()) == 1 {
					m.query += msg.String()
				}
			}

		case screenResults:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.shows)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.shows) > 0 {
					m.selected = m.shows[m.cursor]
					m.loading = true
					m.err = nil
					
					if m.selected.Type == "movie" {
						return m, fetchTorrentsMovie(m.selected.Name)
					}
					return m, SearchSeason(m.selected.ID)
				}
			case "esc":
				m.screen = screenSearch
				m.shows = nil
				m.cursor = 0
				m.err = nil
			}

		case screenSeasons:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.seasons)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.seasons) > 0 {
					m.season = m.seasons[m.cursor]
					m.loading = true
					m.err = nil
					return m, SearchEp(m.selected.ID, m.season.Number)
				}
			case "esc":
				m.screen = screenResults
				m.seasons = nil
				m.cursor = 0
				m.err = nil
			}

		case screenEpisodes:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.episodes)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.episodes) > 0 {
					episode := m.episodes[m.cursor]
					m.loading = true
					m.err = nil
					return m, fetchTorrentsEpisode(m.selected.Name, m.season.Number, episode.Number)
				}
			case "esc":
				m.screen = screenSeasons
				m.episodes = nil
				m.cursor = 0
				m.err = nil
			}

		case screenTorrents:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.torrents)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.torrents) > 0 {
					player.Start(m.torrents[m.cursor].Magnet)
					return m, nil
				}
			case "esc":
				m.cursor = 0
				m.err = nil
				if m.selected.Type == "movie" {
					m.screen = screenResults
				} else {
					m.screen = screenEpisodes
				}
			}
		}

	case SearchResultsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.shows = msg.shows
		m.cursor = 0
		m.err = nil
		m.screen = screenResults

	case SeasonResultsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.seasons = msg.seasons
		m.cursor = 0
		m.err = nil
		m.screen = screenSeasons

	case EpisodeResultsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.episodes = msg.episodes
		m.cursor = 0
		m.err = nil
		m.screen = screenEpisodes

	case TorrentResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.torrents = msg.torrents
		m.cursor = 0
		m.screen = screenTorrents
	}

	return m, nil
}

func (m model) View() string {
	
	if m.err != nil {
		return fmt.Sprintf("cinex_\n\nerror: %v\n\npress esc to go back\n", m.err)
	}

	
	if m.loading {
		return "cinex_\n\nsearching...\n"
	}

	s := "cinex_\n\n"

	switch m.screen {

	case screenSearch:
		s += "Search: " + m.query + "_\n\n"
		s += "press enter to search  ctrl+c to quit\n"

	case screenResults:
		s += "Results:\n\n"
		for i, show := range m.shows {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			year := ""
			if len(show.Year) >= 4 {
				year = show.Year[:4]
			}
			s += fmt.Sprintf("%s %-35s %s [%s]\n", cursor, show.Name, year, show.Type)
		}
		s += "\n↑/↓ to move  enter to select  esc to go back\n"

	case screenSeasons:
		s += fmt.Sprintf("Seasons — %s\n\n", m.selected.Name)
		for i, season := range m.seasons {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, season.Name)
		}
		s += "\n↑/↓ to move  enter to select  esc to go back\n"

	case screenEpisodes:
		s += fmt.Sprintf("Episodes — %s %s\n\n", m.selected.Name, m.season.Name)
		for i, ep := range m.episodes {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s E%02d — %s\n", cursor, ep.Number, ep.Title)
		}
		s += "\n↑/↓ to move  enter to select  esc to go back\n"

	case screenTorrents:
		s += fmt.Sprintf("Torrents — %s\n\n", m.selected.Name)
		for i, t := range m.torrents {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			title := t.Title
			if len(title) > 45 {
				title = title[:45] + "..."
			}

			s += fmt.Sprintf("%s [%d] %-48s %s seeds:%d\n",
				cursor,
				i,
				title,
				t.Size,
				t.Seeders,
			)
		}
		s += "\n↑/↓ to move  enter to stream  esc to go back\n"
	}

	return s
}
