package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorPrimary   = lipgloss.Color("#")
	ColorSecondary = lipgloss.Color("111")
	ColorText      = lipgloss.Color("15")
	ColorMuted1    = lipgloss.Color("245")
	ColorMuted2    = lipgloss.Color("235")
)

var (
	List     = lipgloss.NewStyle().Margin(0, 1)
	HelpKey  = lipgloss.NewStyle().Foreground(ColorSecondary)
	HelpDesc = lipgloss.NewStyle().Foreground(ColorMuted1)
	Help     = lipgloss.NewStyle().Margin(0, 1, 0, 3)
	Input    = lipgloss.NewStyle().MarginLeft(3).MarginBottom(1)
	Header   = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, false, true, false).
			BorderForeground(ColorMuted2).
			MarginLeft(3).
			MarginRight(3).
			MarginBottom(1)
)
