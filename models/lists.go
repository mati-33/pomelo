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
	modifyMode
)

type item struct {
	data.List
}

func (i item) Title() string { return i.Name }
func (i item) Description() string {
	return fmt.Sprintf("Created: %s :: tasks: %d", i.Created.Format("02-01-2006 15:04"), i.TaskCount)
}
func (i item) FilterValue() string { return i.Name }

type listsScreen struct {
	list   list.Model
	input  textinput.Model
	mode   mode
	err    error
	db     *sql.DB
	width  int
	height int
}

func newListsScreen(db *sql.DB) listsScreen {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "type / to search"
	l.FilterInput.Prompt = "/"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()

	i := textinput.New()
	i.Prompt = ""

	return listsScreen{
		list:  l,
		input: i,
		mode:  listMode,
		db:    db,
	}
}

type LoadResult struct {
	lists []list.Item
	err   error
}

type ListAdded struct{}

type ListDeleted struct{}

type ListModified struct{}

func (m listsScreen) Init() tea.Cmd {
	return GetAllListsCmd(m.db)
}

func (m listsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case ListAdded, ListDeleted, ListModified:
		return m, GetAllListsCmd(m.db)

	case LoadResult:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		cmd := m.list.SetItems(msg.lists)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-3)

	case tea.KeyMsg:
		switch m.mode {
		case listMode:
			if m.list.FilterState() != list.Filtering {
				switch msg.String() {
				case "a":
					m.mode = addMode
					m.list.SetShowTitle(false)
					m.input.Prompt = "name: "
					return m, m.input.Focus()
				case "d":
					m.mode = deleteMode
					m.list.SetShowTitle(false)
					m.input.Prompt = "delete? (y/n): "
					return m, m.input.Focus()
				case "r":
					m.mode = modifyMode
					m.list.SetShowTitle(false)
					m.input.Prompt = "rename: "
					name := m.list.SelectedItem().(item).Name
					m.input.SetValue(name)
					return m, m.input.Focus()
				case "enter":
					id := m.list.SelectedItem().(item).ID
					return m, func() tea.Msg {
						screen := newListScreen(id, m.db, m.width, m.height-3)
						return PushScreenMsg{screen, screen.Init()}
					}
				}
			}

		case deleteMode:
			switch msg.String() {
			case "enter":
				value := m.input.Value()
				if value == "n" {
					m.mode = listMode
					m.list.SetShowTitle(true)
					m.input.SetValue("")
					m.input.Blur()
				}
				if value == "y" {
					m.mode = listMode
					m.list.SetShowTitle(true)
					m.input.SetValue("")
					m.input.Blur()
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
				m.input.SetValue("")
				m.input.Blur()
			}

		case addMode:
			switch msg.String() {
			case "enter":
				name := m.input.Value()
				if len(name) > 0 {
					m.mode = listMode
					m.list.SetShowTitle(true)
					m.input.SetValue("")
					m.input.Blur()
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
				m.input.SetValue("")
				m.input.Blur()
			}

		case modifyMode:
			switch msg.String() {
			case "enter":
				renamed := m.input.Value()
				if len(renamed) > 0 {
					id := m.list.SelectedItem().(item).ID
					m.mode = listMode
					m.list.SetShowTitle(true)
					m.input.SetValue("")
					m.input.Blur()
					return m, func() tea.Msg {
						err := data.ModifyList(m.db, id, renamed)
						if err != nil {
							return nil
						}
						return ListModified{}
					}
				}
			case "esc":
				m.mode = listMode
				m.list.SetShowTitle(true)
				m.input.SetValue("")
				m.input.Blur()
			}
		}
	}

	var cmd tea.Cmd
	switch m.mode {
	case listMode:
		m.list, cmd = m.list.Update(msg)
	case addMode, deleteMode, modifyMode:
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m listsScreen) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error!!: %v", m.err)
	}

	ret := ""

	switch m.mode {
	case addMode, deleteMode, modifyMode:
		ret += m.input.View() + "\n"
	}

	ret += m.list.View() + "\n"
	ret += "j   k   a add  d delete  r rename  enter details  ctrl+c exit"
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
