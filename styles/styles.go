package styles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorPrimary     = lipgloss.Color("33")  // #0087ff
	ColorSecondary   = lipgloss.Color("111") // #87afff
	ColorText        = lipgloss.Color("15")  // #ffffff
	ColorMuted1      = lipgloss.Color("245") // #8a8a8a
	ColorMuted2      = lipgloss.Color("235") // #262626
	ColorAdd         = lipgloss.Color("36")  // #00af87
	ColorRename      = lipgloss.Color("136") // #af8700
	ColorDelete      = lipgloss.Color("204") // #ff5f87
	ColorFilter      = lipgloss.Color("135") // #af5fff
	ColorTaskDone    = ColorAdd
	ColorTaskNotDone = ColorDelete

	List        = lipgloss.NewStyle().Margin(0, 3, 0, 1)
	HelpKey     = lipgloss.NewStyle().Foreground(ColorSecondary)
	HelpDesc    = lipgloss.NewStyle().Foreground(ColorMuted1)
	Help        = lipgloss.NewStyle().Margin(0, 3, 0, 3)
	InputBase   = lipgloss.NewStyle().Foreground(ColorText).Padding(0, 0, 0, 1).Margin(0, 1, 0, 3)
	InputAdd    = InputBase.Background(ColorAdd)
	InputRename = InputBase.Background(ColorRename)
	InputDelete = InputBase.Background(ColorDelete)
	TaskDone    = lipgloss.NewStyle().Foreground(ColorTaskDone)
	TaskNotDone = lipgloss.NewStyle().Foreground(ColorTaskNotDone)
	Header      = lipgloss.NewStyle().
			Margin(1, 1, 2, 1).
			Padding(0, 1).
			Border(lipgloss.ThickBorder(), false, true, false, true).
			BorderForeground(ColorPrimary)
)

func ListStyles() list.Styles {
	s := list.DefaultStyles()

	s.Title = lipgloss.NewStyle().Foreground(ColorText).Background(ColorPrimary).Padding(0, 1)
	s.FilterPrompt = lipgloss.NewStyle().Foreground(ColorText).Background(ColorFilter).Padding(0, 1).Margin(0, 1, 0, 0)
	s.FilterCursor = lipgloss.NewStyle().Foreground(ColorText)
	s.NoItems = lipgloss.NewStyle().Foreground(ColorText).Margin(0, 2)

	return s
}

func ListItemStyles() list.DefaultItemStyles {
	s := list.NewDefaultItemStyles()

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(ColorPrimary).
		Foreground(ColorText).
		Bold(true).
		Padding(0, 0, 0, 1)
	s.SelectedDesc = s.SelectedTitle.UnsetBold().UnsetUnderline()

	return s
}

func HelpStyles() help.Styles {
	keyStyle := lipgloss.NewStyle().Foreground(ColorSecondary)
	descStyle := lipgloss.NewStyle().Foreground(ColorMuted1)
	sepStyle := lipgloss.NewStyle().Foreground(ColorMuted1)

	return help.Styles{
		ShortKey:       keyStyle,
		ShortDesc:      descStyle,
		ShortSeparator: sepStyle,
		Ellipsis:       sepStyle,
		FullKey:        keyStyle,
		FullDesc:       descStyle,
		FullSeparator:  sepStyle,
	}
}
