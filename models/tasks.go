package models

import (
	"database/sql"
	"fmt"
	"pomelo/components"
	"pomelo/data"
	"pomelo/styles"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type taskItem struct {
	data.Task
}

func (i taskItem) Title() string { return i.Name }
func (i taskItem) Description() string {
	done := styles.TaskNotDone.Render("")
	if i.IsDone {
		done = styles.TaskDone.Render("󰸞")
	}
	return fmt.Sprintf("added: %s  done: %s", i.Created.Format("02-01-2006 15:04"), done)
}
func (i taskItem) FilterValue() string { return i.Name }

type tasksKeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Navigate   key.Binding
	Add        key.Binding
	Rename     key.Binding
	Delete     key.Binding
	ToggleDone key.Binding
	Back       key.Binding
	Filter     key.Binding
	Help       key.Binding
	Quit       key.Binding
}

func (k tasksKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Navigate, k.Help, k.Quit}
}

func (k tasksKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Add, k.Rename, k.Delete, k.ToggleDone},
		{k.Filter, k.Back, k.Help, k.Quit},
	}
}

var tasksKeys = tasksKeyMap{
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
	ToggleDone: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "toggle done"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "move back"),
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

type LoadTasks struct {
	tasks []list.Item
	err   error
}

type tasksScreen struct {
	header components.Header
	list   list.Model
	input  textinput.Model
	keys   tasksKeyMap
	help   help.Model
	listId int64
	mode   mode
	err    error
	db     *sql.DB
}

func newTasksScreen(listId int64, db *sql.DB) tasksScreen {
	delegate := list.NewDefaultDelegate()
	delegate.Styles = styles.ListItemStyles()
	l := list.New([]list.Item{}, delegate, Width, Height-5)
	l.Title = ""
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

	return tasksScreen{
		header: components.NewHeader(pomeloASCI, " task manager", "v0.1.0", Width),
		list:   l,
		input:  i,
		keys:   tasksKeys,
		help:   h,
		mode:   listMode,
		listId: listId,
		db:     db,
	}
}

type TaskAdded struct{}

type TaskDeleted struct{}

type TaskModified struct{}

type TaskToggled struct{}

func (m tasksScreen) Init() tea.Cmd {
	return tea.Batch(GetAllTasksCmd(m.db, m.listId), func() tea.Msg {
		list, err := data.GetList(m.db, m.listId)
		if err != nil {
			return nil
		}
		return list
	})
}

func (m tasksScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case error:
		m.err = msg
		return m, nil

	case TaskAdded, TaskDeleted, TaskModified, TaskToggled:
		return m, GetAllTasksCmd(m.db, m.listId)

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-3)
		m.help.Width = msg.Width
		m.header.SetWidth(msg.Width)

	case LoadTasks:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		cmd := m.list.SetItems(msg.tasks)
		return m, cmd

	case data.List:
		m.list.Title = msg.Name
		return m, nil

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
					name := m.list.SelectedItem().(taskItem).Name
					m.input.SetValue(name)
					return m, m.input.Focus()
				case "enter":
					task := m.list.SelectedItem().(taskItem)
					return m, func() tea.Msg {
						err := data.ModifyTask(m.db, task.ID, task.Name, !task.IsDone)
						if err != nil {
							return err
						}
						return TaskToggled{}
					}
				case "?":
					m.help.ShowAll = !m.help.ShowAll
					return m, nil
				case "esc":
					return m, func() tea.Msg { return PopScreenMsg{GetAllListsCmd(m.db)} }
				}
			}
		case deleteMode:
			switch msg.String() {
			case "enter":
				value := m.input.Value()
				if value == "y" {
					id := m.list.SelectedItem().(taskItem).ID
					m.mode = listMode
					m.list.SetShowTitle(true)
					m.input.SetValue("")
					m.input.Blur()
					return m, func() tea.Msg {
						err := data.DeleteTask(m.db, id)
						if err != nil {
							return err
						}
						return TaskDeleted{}
					}
				}
				if value == "n" {
					m.mode = listMode
					m.list.SetShowTitle(true)
					m.input.SetValue("")
					m.input.Blur()
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
						_, err := data.AddTask(m.db, m.listId, name)
						if err != nil {
							return nil
						}
						return TaskAdded{}
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
				rename := m.input.Value()
				if len(rename) > 0 {
					task := m.list.SelectedItem().(taskItem)
					m.mode = listMode
					m.list.SetShowTitle(true)
					m.input.SetValue("")
					m.input.Blur()
					return m, func() tea.Msg {
						err := data.ModifyTask(m.db, task.ID, rename, task.IsDone)
						if err != nil {
							return err
						}
						return TaskAdded{}
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

func (m tasksScreen) View() string {
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

func GetAllTasksCmd(db *sql.DB, listId int64) tea.Cmd {
	return func() tea.Msg {
		tasks, err := data.GetAllListTasks(db, listId)
		if err != nil {
			return LoadTasks{nil, err}
		}

		taskItems := make([]list.Item, 0, len(tasks))
		for _, t := range tasks {
			taskItems = append(taskItems, taskItem{t})
		}

		return LoadTasks{taskItems, nil}
	}
}
