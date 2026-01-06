package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbletea"
)

const pomeloASCII = `
                           ▄▄       
                           ██       
████▄ ▄███▄ ███▄███▄ ▄█▀█▄ ██ ▄███▄ 
██ ██ ██ ██ ██ ██ ██ ██▄█▀ ██ ██ ██ 
████▀ ▀███▀ ██ ██ ██ ▀█▄▄▄ ██ ▀███▀ 
██                                  
▀▀                                  
	`

func main() {
	m := model{}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
}

type model struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	b := strings.Builder{}
	b.WriteString(pomeloASCII)

	return b.String()
}
