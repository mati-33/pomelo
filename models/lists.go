package models

import (
	"github.com/charmbracelet/bubbletea"
	"pomelo/lists"
	"pomelo/styles"
	"strconv"
	"strings"
)

type listsScreen struct {
	lists   []lists.List
	focused int
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

		case "enter":
			return m, PushScreen(NewListDetailsScreen(m.lists[m.focused]))
		}
	}
	return m, nil
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

	return styles.ListsStyle.Width(Width).Render(b.String())
}

func newListsScreen() listsScreen {
	return listsScreen{}
}
