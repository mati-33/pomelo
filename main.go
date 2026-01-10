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
	width  int
	height int

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
	m := initialMainModel()
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

type PopScreen struct{}

type PushScreen struct {
	screen tea.Model
}

type mainModel struct {
	stack []tea.Model
}

func initialMainModel() mainModel {
	stack := []tea.Model{}
	listsScreen := newListsScreen(0)
	stack = append(stack, listsScreen)
	return mainModel{stack: stack}
}

func (m mainModel) Init() tea.Cmd {
	return m.stack[len(m.stack)-1].Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		width = msg.Width
		height = msg.Height

	case PopScreen:
		m.stack = m.stack[:len(m.stack)-1]
		return m, nil

	case PushScreen:
		m.stack = append(m.stack, msg.screen)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd

	currentScreen := m.stack[len(m.stack)-1]
	currentScreen, cmd = currentScreen.Update(msg)
	m.stack[len(m.stack)-1] = currentScreen

	return m, cmd
}

func (m mainModel) View() string {
	return m.stack[len(m.stack)-1].View()
}

type listsScreen struct {
	lists   []list
	focused int
}

func (m listsScreen) Init() tea.Cmd {
	return nil
}

func (m listsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "j":
			if m.focused < len(m.lists)-1 {
				m.focused++
			}

		case "k":
			if m.focused > 0 {
				m.focused--
			}

		case "enter":
			return m, func() tea.Msg {
				return PushScreen{newListDetailsScreen(m.lists[m.focused])}
			}
		}
	}
	return m, nil
}

func (m listsScreen) View() string {
	b := strings.Builder{}
	b.WriteString(logoStyle.Width(width).Render(pomeloASCII))
	b.WriteString("\n\n")

	for i, l := range m.lists {
		bold := false
		style := listStyle.Width(width - 4)
		if i == m.focused {
			bold = true
			style = listFocusStyle.Width(width - 6)
		}

		taskStr := nameStyle.Bold(bold).Render(l.name)
		taskStr += "\n"
		taskStr += infoStyle.Render("created:", l.created, "tasks:", strconv.Itoa(len(l.tasks)))
		b.WriteString(style.Render(taskStr))
		b.WriteString("\n")
	}

	return listsStyle.Width(width).Render(b.String())
}

func newListsScreen(focused int) listsScreen {
	return listsScreen{
		focused: focused,
		lists: []list{
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

type listDetailsScreen struct {
	list list
}

func newListDetailsScreen(list list) listDetailsScreen {
	return listDetailsScreen{list}
}

func (m listDetailsScreen) Init() tea.Cmd {
	return nil
}

func (m listDetailsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			return m, func() tea.Msg { return PopScreen{} }

		}
	}
	return m, nil
}

func (m listDetailsScreen) View() string {
	b := strings.Builder{}
	b.WriteString(logoStyle.Width(width).Render(pomeloASCII))
	b.WriteString("\n\n")

	b.WriteString("this is list details view\n")
	b.WriteString(m.list.name)

	return b.String()
}
