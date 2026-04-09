package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/inspiring-group/inspiring-swiss-knife/ui"
)

func main() {
	// Create and run the Bubble Tea program
	p := tea.NewProgram(
		ui.New(),
		tea.WithAltScreen(),       // use full terminal, restore on exit
		tea.WithMouseCellMotion(), // enable mouse support
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
