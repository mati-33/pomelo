package models

import (
	"github.com/charmbracelet/bubbletea"
	"pomelo/lists"
	"pomelo/styles"
	"strings"
)

type listDetailsScreen struct {
	list lists.List
}

func (m listDetailsScreen) Init() tea.Cmd {
	return nil
}

func (m listDetailsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			return m, PopScreen

		}
	}
	return m, nil
}

func (m listDetailsScreen) View() string {
	b := strings.Builder{}
	b.WriteString(styles.LogoStyle.Width(Width).Render(pomeloASCII))
	b.WriteString("\n\n")

	b.WriteString("this is list details view\n")
	b.WriteString(m.list.Name)

	return b.String()
}

func NewListDetailsScreen(list lists.List) listDetailsScreen {
	return listDetailsScreen{list}
}
