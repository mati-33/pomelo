package models

import (
	"github.com/charmbracelet/bubbletea"
	"pomelo/styles"
	"strconv"
	"strings"
)

type task struct {
	name   string
	isDone bool
}

type list struct {
	name     string
	created  string
	modified string
	tasks    []task
}

type listsScreen struct {
	lists   []list
	focused int
}

func (m listsScreen) Init() tea.Cmd {
	return nil
}

func (m listsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

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

	for i, l := range m.lists {
		bold := false
		style := styles.ListStyle.Width(Width - 4)
		if i == m.focused {
			bold = true
			style = styles.ListFocusStyle.Width(Width - 6)
		}

		taskStr := styles.NameStyle.Bold(bold).Render(l.name)
		taskStr += "\n"
		taskStr += styles.InfoStyle.Render("created:", l.created, "tasks:", strconv.Itoa(len(l.tasks)))
		b.WriteString(style.Render(taskStr))
		b.WriteString("\n")
	}

	return styles.ListsStyle.Width(Width).Render(b.String())
}

func newListsScreen(focused int) listsScreen {
	return listsScreen{
		focused: focused,
		lists: []list{
			{
				name:     "pomelo project",
				created:  "12:17 AM",
				modified: "-",
				tasks: []task{
					{name: "tui"},
					{name: "backend"},
				},
			},
			{
				name:     "06-10-2026",
				created:  "11:11 AM",
				modified: "12:12",
				tasks:    []task{},
			},
			{
				name:     "07-10-2026",
				created:  "10:10 PM",
				modified: "-",
				tasks:    []task{},
			},
			{
				name:     "terminal typing app",
				created:  "04:27 AM",
				modified: "-",
				tasks:    []task{},
			},
			{
				name:     "terminal chat app",
				created:  "08:00 PM",
				modified: "09:30",
				tasks:    []task{},
			},
		}}
}
