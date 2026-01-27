package models

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"pomelo/components"
	"pomelo/data"
	"pomelo/styles"
)

type item struct {
	data.List
}

func (i item) Title() string { return i.Name }
func (i item) Description() string {
	return fmt.Sprintf("added: %s  tasks: %d", i.Created.Format("02-01-2006 15:04"), i.TaskCount)
}
func (i item) FilterValue() string { return i.Name }

type listsKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Navigate key.Binding
	Add      key.Binding
	Rename   key.Binding
	Delete   key.Binding
	Tasks    key.Binding
	Filter   key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func (k listsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Navigate, k.Help, k.Quit}
}

func (k listsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Add, k.Rename, k.Delete, k.Tasks},
		{k.Filter, k.Help, k.Quit},
	}
}

var listsKeys = listsKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("l", "move right"),
	),
	Navigate: key.NewBinding(
		key.WithKeys("h", "j", "k", "l"),
		key.WithHelp("hjkl", "navigate"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Tasks: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "see tasks"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type listsScreen struct {
	header components.Header
	list   list.Model
	input  textinput.Model
	keys   listsKeyMap
	help   help.Model
	mode   mode
	err    error
	db     *sql.DB
}

func newListsScreen(db *sql.DB) listsScreen {
	delegate := list.NewDefaultDelegate()
	delegate.Styles = styles.ListItemStyles()
	l := list.New([]list.Item{}, delegate, Width, Height-5)
	l.Title = "your lists"
	l.FilterInput.Prompt = "filter:"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()
	l.Styles = styles.ListStyles()
	l.FilterInput.PromptStyle = l.Styles.FilterPrompt
	l.FilterInput.Cursor.Style = l.Styles.FilterCursor

	i := textinput.New()
	i.Prompt = ""

	h := help.New()
	h.Styles = styles.HelpStyles()

	return listsScreen{
		header: components.NewHeader(pomeloASCI, " task manager", "v0.1.0", Width),
		list:   l,
		input:  i,
		keys:   listsKeys,
		help:   h,
		mode:   listMode,
		db:     db,
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
		m.list.SetSize(msg.Width, msg.Height-5)
		m.help.Width = msg.Width
		m.header.SetWidth(msg.Width)

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
						screen := newTasksScreen(id, m.db)
						return PushScreenMsg{screen, screen.Init()}
					}
				case "?":
					m.help.ShowAll = !m.help.ShowAll
					return m, nil
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

	header := m.header.View()
	help := m.help.View(m.keys)

	listHeight := Height - lipgloss.Height(header) - lipgloss.Height(help)
	input := ""

	var inputStyle lipgloss.Style

	switch m.mode {
	case addMode:
		inputStyle = styles.InputAdd
	case modifyMode:
		inputStyle = styles.InputRename
	case deleteMode:
		inputStyle = styles.InputDelete
	}

	switch m.mode {
	case addMode, deleteMode, modifyMode:
		m.input.PromptStyle = inputStyle
		input = m.input.View()
		listHeight = listHeight - lipgloss.Height(input)
		m.list.SetHeight(listHeight)
		list := m.list.View()
		return lipgloss.JoinVertical(lipgloss.Left, header, input, styles.List.Render(list), styles.Help.Render(help))
	}

	// todo: is there a way to do that in Update?
	m.list.SetHeight(listHeight)
	list := m.list.View()

	return lipgloss.JoinVertical(lipgloss.Left, header, styles.List.Render(list), styles.Help.Render(help))
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
