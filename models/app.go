package models

import (
	"database/sql"
	"strings"

	"github.com/charmbracelet/bubbletea"
)

var (
	Width  int
	Height int
)

type pomeloModel struct {
	stack []tea.Model
}

func InitialPomeloModel(db *sql.DB) pomeloModel {
	stack := []tea.Model{}
	listsScreen := newListsScreen(db)
	stack = append(stack, listsScreen)
	return pomeloModel{stack: stack}
}

func (m pomeloModel) Init() tea.Cmd {
	return m.stack[len(m.stack)-1].Init()
}

func (m pomeloModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		Width = msg.Width
		Height = msg.Height

	case PopScreenMsg:
		m.stack = m.stack[:len(m.stack)-1]
		return m, msg.cmd

	case PushScreenMsg:
		m.stack = append(m.stack, msg.Screen)
		return m, msg.Cmd

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

func (m pomeloModel) View() string {
	b := strings.Builder{}
	b.WriteString("  pomelo v0.1.0\n\n")
	b.WriteString(m.stack[len(m.stack)-1].View())
	return b.String()
}

type PopScreenMsg struct {
	cmd tea.Cmd
}

type PushScreenMsg struct {
	Screen tea.Model
	Cmd    tea.Cmd
}

// todo: WithCommand(cmd tea.Cmd)
func PopScreen(cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return PopScreenMsg{cmd}
	}
}

func PushScreen(screen tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushScreenMsg{screen, nil}
	}
}
