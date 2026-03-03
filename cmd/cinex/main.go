package main

import (
	"fmt"
	"os"

	"github.com/SssHhhAaaDddOoWww/cinex_/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	p := tea.NewProgram(tui.InitialModel())
	_, err := p.Run()
	if err != nil {
		fmt.Println("Error occured while starting !!!")
		os.Exit(1)

	}

	// torrents, err := provider.GetTorrent("inception", "movie")
	// if err != nil {
	// 	fmt.Println("torrent error:", err)
	// 	os.Exit(1)
	// }

	// for i, t := range torrents {
	// 	fmt.Printf("[%d] %s | %s | seeds: %d\n", i+1, t.Title, t.Size, t.Seeders)
	// }

}
