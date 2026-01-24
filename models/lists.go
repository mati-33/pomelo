package models

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"

	"pomelo/data"
)

type mode int

const (
	listMode mode = iota
	addMode
	deleteMode
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
	list        list.Model
	addInput    textinput.Model
	deleteInput textinput.Model
	mode        mode
	err         error
	db          *sql.DB
}

func newListsScreen(db *sql.DB) listsScreen {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "type / to search"
	l.FilterInput.Prompt = "/"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()

	ai := textinput.New()
	ai.Prompt = ""
	di := textinput.New()
	di.Prompt = ""
	di.CharLimit = 1

	return listsScreen{
		list:        l,
		addInput:    ai,
		deleteInput: di,
		mode:        listMode,
		db:          db,
	}
}

type LoadResult struct {
	lists []list.Item
	err   error
}

type ListAdded struct{}

type ListDeleted struct{}

func (m listsScreen) Init() tea.Cmd {
	return GetAllListsCmd(m.db)
}

func (m listsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case ListAdded, ListDeleted:
		return m, GetAllListsCmd(m.db)

	case LoadResult:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		cmd := m.list.SetItems(msg.lists)
		return m, cmd

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-3)

	case tea.KeyMsg:
		switch m.mode {
		case listMode:
			if m.list.FilterState() != list.Filtering {
				switch msg.String() {
				case "a":
					m.mode = addMode
					m.list.SetShowTitle(false)
					return m, m.addInput.Focus()
				case "d":
					m.mode = deleteMode
					m.list.SetShowTitle(false)
					return m, m.deleteInput.Focus()
				}
			}

		case deleteMode:
			switch msg.String() {
			case "enter":
				value := m.deleteInput.Value()
				if value == "n" {
					m.mode = listMode
					m.list.SetShowTitle(true)
					m.deleteInput.SetValue("")
					m.deleteInput.Blur()
				}
				if value == "y" {
					m.mode = listMode
					m.list.SetShowTitle(true)
					m.deleteInput.SetValue("")
					m.deleteInput.Blur()
					id := m.list.SelectedItem().(item).ID
					return m, func() tea.Msg {
						err := data.DeleteList(m.db, id)
						if err != nil {
							return nil
						}
						return ListDeleted{}
					}

				}
			case "esc":
				m.mode = listMode
				m.list.SetShowTitle(true)
				m.deleteInput.SetValue("")
				m.deleteInput.Blur()
			}

		case addMode:
			switch msg.String() {
			case "enter":
				name := m.addInput.Value()
				if len(name) > 0 {
					m.mode = listMode
					m.list.SetShowTitle(true)
					m.addInput.SetValue("")
					m.addInput.Blur()
					return m, func() tea.Msg {
						_, err := data.AddList(m.db, name)
						if err != nil {
							return nil
						}
						return ListAdded{}
					}
				}
			case "esc":
				m.mode = listMode
				m.list.SetShowTitle(true)
				m.addInput.SetValue("")
				m.addInput.Blur()
			}
		}
	}

	var cmd tea.Cmd
	switch m.mode {
	case listMode:
		m.list, cmd = m.list.Update(msg)
	case addMode:
		m.addInput, cmd = m.addInput.Update(msg)
	case deleteMode:
		m.deleteInput, cmd = m.deleteInput.Update(msg)
	}
	return m, cmd

}

func (m listsScreen) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error!!: %v", m.err)
	}

	ret := ""

	switch m.mode {

	case deleteMode:
		ret += "delete? (y/n): " + m.deleteInput.View() + "\n"

	case addMode:
		ret += "name: " + m.addInput.View() + "\n"
	}

	ret += m.list.View()
	return ret
}

func GetAllListsCmd(db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		lists, err := data.GetAllLists(db)
		if err != nil {
			return LoadResult{[]list.Item{}, err}
		}

		items := make([]list.Item, 0, len(lists))
		for _, l := range lists {
			items = append(items, item{l})
		}

		return LoadResult{items, err}
	}
}
