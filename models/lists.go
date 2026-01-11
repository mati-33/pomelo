package models

import (
	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"pomelo/components"
	"pomelo/lists"
	"pomelo/styles"
	"strconv"
	"strings"
)

type listsScreen struct {
	lists   []lists.List
	focused int
	notif   components.Notif
}

func (m listsScreen) Init() tea.Cmd {
	return func() tea.Msg {
		return lists.GetAllLists()
	}
}

func (m listsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case []lists.List:
		m.lists = msg
		return m, nil

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

		case "ctrl+a":
			lists.AddList(lists.List{
				Name:     "added from tui",
				Created:  "todo",
				Modified: "nil",
				Tasks:    []lists.Task{},
			})
			return m, tea.Batch(func() tea.Msg { return lists.GetAllLists() }, m.notif.Notify("hello world"))

		case "enter":
			return m, PushScreen(NewListDetailsScreen(m.lists[m.focused]))
		}

	}
	notif, cmd := m.notif.Update(msg)
	m.notif = notif

	return m, cmd
}

func (m listsScreen) View() string {
	b := strings.Builder{}
	b.WriteString(styles.LogoStyle.Width(Width).Render(pomeloASCII))
	b.WriteString("\n\n")

	if len(m.lists) == 0 {
		b.WriteString("no lists defined\n")
		return b.String()
	}

	for i, l := range m.lists {
		bold := false
		style := styles.ListStyle.Width(Width - 4)
		if i == m.focused {
			bold = true
			style = styles.ListFocusStyle.Width(Width - 6)
		}

		taskStr := styles.NameStyle.Bold(bold).Render(l.Name)
		taskStr += "\n"
		taskStr += styles.InfoStyle.Render("created:", l.Created, "tasks:", strconv.Itoa(len(l.Tasks)))
		b.WriteString(style.Render(taskStr))
		b.WriteString("\n")
	}

	h := gloss.Height(b.String())
	notif := gloss.Place(Width, Height+1-h, gloss.Right, gloss.Bottom, m.notif.View())
	b.WriteString(notif)

	return styles.ListsStyle.Width(Width).Render(b.String())
}

func newListsScreen() listsScreen {
	return listsScreen{notif: components.Notif{}}
}
