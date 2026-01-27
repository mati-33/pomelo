package models

import (
	"database/sql"

	"github.com/charmbracelet/bubbletea"
)

type pomeloModel struct {
	stack  []tea.Model
	db     *sql.DB
	loaded bool
}

func InitialPomeloModel(db *sql.DB) pomeloModel {
	return pomeloModel{
		stack:  []tea.Model{},
		db:     db,
		loaded: false,
	}
}

func (m pomeloModel) Init() tea.Cmd {
	return nil
}

func (m pomeloModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		Width = msg.Width
		Height = msg.Height
		if !m.loaded {
			firstScreen := newListsScreen(m.db)
			m.stack = append(m.stack, firstScreen)
			m.loaded = true
			return m, firstScreen.Init()
		}

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

	if len(m.stack) > 0 {
		var cmd tea.Cmd

		currentScreen := m.stack[len(m.stack)-1]
		currentScreen, cmd = currentScreen.Update(msg)
		m.stack[len(m.stack)-1] = currentScreen

		return m, cmd
	}

	return m, nil
}

func (m pomeloModel) View() string {
	if len(m.stack) > 0 {
		return m.stack[len(m.stack)-1].View()
	}

	return ""
}
