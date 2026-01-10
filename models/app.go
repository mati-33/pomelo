package models

import (
	"github.com/charmbracelet/bubbletea"
)

const pomeloASCII = `
 ▄▄▄  ▄▄▄  ▄   ▄  ▄▄▄  ▄    ▄▄▄ 
 █▄█  █ █  █▀▄▀█  █▄▄  █    █ █ 
 █    ███  █   █  █▄▄  █▄▄  ███ 
                                     `

var (
	Width  int
	Height int
)

type pomeloModel struct {
	stack []tea.Model
}

func InitialPomeloModel() pomeloModel {
	stack := []tea.Model{}
	listsScreen := newListsScreen(0)
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
		return m, nil

	case PushScreenMsg:
		m.stack = append(m.stack, msg.Screen)
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

func (m pomeloModel) View() string {
	return m.stack[len(m.stack)-1].View()
}

type PopScreenMsg struct{}

type PushScreenMsg struct {
	Screen tea.Model
}

func PopScreen() tea.Msg {
	return PopScreenMsg{}
}

func PushScreen(screen tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushScreenMsg{screen}
	}
}
