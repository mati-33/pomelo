package models

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"

	"pomelo/data"
)

type item struct {
	data.List
}

func (i item) Title() string { return i.Name }
func (i item) Description() string {
	return fmt.Sprintf("Created: %s :: tasks: 3", i.Created.Format("02-01-2006 15:04"))
}
func (i item) FilterValue() string { return i.Name }

type listsScreen struct {
	list list.Model
	err  error
	db   *sql.DB
}

func newListsScreen(db *sql.DB) listsScreen {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "type / to search"
	l.FilterInput.Prompt = "/"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)

	return listsScreen{
		list: l,
		db:   db,
	}
}

type LoadResult struct {
	lists []data.List
	err   error
}

func (m listsScreen) Init() tea.Cmd {
	return func() tea.Msg {
		lists, err := data.GetAllLists(m.db)
		return LoadResult{lists, err}
	}
}

func (m listsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoadResult:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		items := make([]list.Item, 0, len(msg.lists))
		for _, l := range msg.lists {
			items = append(items, item{l})
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-3)
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
		case "d":
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listsScreen) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error!!: %v", m.err)
	}
	return m.list.View()
}
