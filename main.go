package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

const pomeloASCII = `
 ▄▄▄  ▄▄▄  ▄   ▄  ▄▄▄  ▄    ▄▄▄ 
 █▄█  █ █  █▀▄▀█  █▄▄  █    █ █ 
 █    ███  █   █  █▄▄  █▄▄  ███ 
                                     `

var (
	grey   = gloss.Color("245")
	accent = gloss.Color("56")
	white  = gloss.Color("#fff")

	logoStyle      = gloss.NewStyle().Align(gloss.Center)
	listStyle      = gloss.NewStyle().Margin(1, 1).Padding(0, 1)
	listFocusStyle = gloss.NewStyle().Margin(1, 0).Border(gloss.NormalBorder(), false, false, false, true).BorderLeftForeground(accent)
	listsStyle     = gloss.NewStyle().Align(gloss.Center)
	infoStyle      = gloss.NewStyle().Foreground(grey)
	nameStyle      = gloss.NewStyle().Foreground(white)
)

func main() {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
}

type task struct {
	name   string
	isDone bool
}

type list struct {
	name     string
	created  string
	modified string
	tasks    []task
}

type model struct {
	lists   []list
	focused int

	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "j":
			if m.focused < len(m.lists)-1 {
				m.focused++
			}

		case "k":
			if m.focused > 0 {
				m.focused--
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	b := strings.Builder{}
	b.WriteString(logoStyle.Width(m.width).Render(pomeloASCII))
	b.WriteString("\n\n")

	for i, l := range m.lists {
		bold := false
		style := listStyle.Width(m.width - 4)
		if i == m.focused {
			bold = true
			style = listFocusStyle.Width(m.width - 6)
		}

		taskStr := nameStyle.Bold(bold).Render(l.name)
		taskStr += "\n"
		taskStr += infoStyle.Render("created:", l.created, "tasks:", strconv.Itoa(len(l.tasks)))
		b.WriteString(style.Render(taskStr))
		b.WriteString("\n")
	}

	return listsStyle.Width(m.width).Render(b.String())
}

func initialModel() model {
	return model{lists: []list{
		{
			name:     "pomelo project",
			created:  "12:17 AM",
			modified: "-",
			tasks: []task{
				{name: "tui"},
				{name: "backend"},
			},
		},
		{
			name:     "06-10-2026",
			created:  "11:11 AM",
			modified: "12:12",
			tasks:    []task{},
		},
		{
			name:     "07-10-2026",
			created:  "10:10 PM",
			modified: "-",
			tasks:    []task{},
		},
		{
			name:     "terminal typing app",
			created:  "04:27 AM",
			modified: "-",
			tasks:    []task{},
		},
		{
			name:     "terminal chat app",
			created:  "08:00 PM",
			modified: "09:30",
			tasks:    []task{},
		},
	}}
}
