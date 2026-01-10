package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"pomelo/models"
)

func main() {
	m := models.InitialPomeloModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
}
