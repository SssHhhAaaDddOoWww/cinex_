package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(InitialModel())
	_, err := p.Run()
	if err != nil {
		fmt.Println("Error occured while starting !!!")
		os.Exit(1)

	}
}
