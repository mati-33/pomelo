package models

import (
	"database/sql"
	"fmt"
	"pomelo/data"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type mode2 int

const (
	listMode2 mode2 = iota
	addMode2
	deleteMode2
	modifyMode2
)

type taskItem struct {
	data.Task
}

func (i taskItem) Title() string {
	if i.IsDone {
		return "[x] " + i.Name
	}
	return "[ ] " + i.Name
}
func (i taskItem) Description() string {
	return fmt.Sprintf("Created: %s", i.Created.Format("02-01-2006 15:04"))
}
func (i taskItem) FilterValue() string { return i.Name }

type LoadTasks struct {
	tasks []list.Item
	err   error
}

type listScreen struct {
	list   list.Model
	input  textinput.Model
	listId int64
	mode   mode2
	err    error
	db     *sql.DB
	width  int
	height int
}

func newListScreen(listId int64, db *sql.DB, width, height int) listScreen {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	l.Title = "type / to search"
	l.FilterInput.Prompt = "/"
	l.SetShowHelp(false)
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()

	i := textinput.New()
	i.Prompt = ""

	return listScreen{
		list:   l,
		input:  i,
		mode:   listMode2,
		listId: listId,
		db:     db,
	}
}

type TaskAdded struct{}

type TaskDeleted struct{}

type TaskModified struct{}

type TaskToggled struct{}

func (m listScreen) Init() tea.Cmd {
	return GetAllTasksCmd(m.db, m.listId)
}

func (m listScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case error:
		m.err = msg
		return m, nil

	case TaskAdded, TaskDeleted, TaskModified, TaskToggled:
		return m, GetAllTasksCmd(m.db, m.listId)

	case LoadTasks:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		cmd := m.list.SetItems(msg.tasks)
		return m, cmd

	case tea.KeyMsg:
		switch m.mode {
		case listMode2:
			if m.list.FilterState() != list.Filtering {
				switch msg.String() {
				case "a":
					m.mode = addMode2
					m.list.SetShowTitle(false)
					m.input.Prompt = "name: "
					return m, m.input.Focus()
				case "d":
					m.mode = deleteMode2
					m.list.SetShowTitle(false)
					m.input.Prompt = "delete? (y/n): "
					return m, m.input.Focus()
				case "r":
					m.mode = modifyMode2
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
				case "esc":
					return m, func() tea.Msg { return PopScreenMsg{GetAllListsCmd(m.db)} }
				}
			}
		case deleteMode2:
			switch msg.String() {
			case "enter":
				value := m.input.Value()
				if value == "y" {
					id := m.list.SelectedItem().(taskItem).ID
					m.mode = listMode2
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
					m.mode = listMode2
					m.list.SetShowTitle(true)
					m.input.SetValue("")
					m.input.Blur()
				}
			case "esc":
				m.mode = listMode2
				m.list.SetShowTitle(true)
				m.input.SetValue("")
				m.input.Blur()
			}
		case addMode2:
			switch msg.String() {
			case "enter":
				name := m.input.Value()
				if len(name) > 0 {
					m.mode = listMode2
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
				m.mode = listMode2
				m.list.SetShowTitle(true)
				m.input.SetValue("")
				m.input.Blur()
			}
		case modifyMode2:
			switch msg.String() {
			case "enter":
				rename := m.input.Value()
				if len(rename) > 0 {
					task := m.list.SelectedItem().(taskItem)
					m.mode = listMode2
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
				m.mode = listMode2
				m.list.SetShowTitle(true)
				m.input.SetValue("")
				m.input.Blur()
			}
		}
	}

	var cmd tea.Cmd
	switch m.mode {
	case listMode2:
		m.list, cmd = m.list.Update(msg)
	case addMode2, deleteMode2, modifyMode2:
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m listScreen) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error!!: %v", m.err)
	}

	ret := ""

	switch m.mode {

	case deleteMode2, addMode2, modifyMode2:
		ret += m.input.View() + "\n"

	}

	ret += m.list.View()
	ret += "\nj   k   a add  d delete  r rename  enter done  esc   ctrl+c exit"
	return ret
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
