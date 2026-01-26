package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorPrimary   = lipgloss.Color("#")
	ColorSecondary = lipgloss.Color("111")
	ColorText      = lipgloss.Color("15")
	ColorMuted1    = lipgloss.Color("245")
	ColorMuted2    = lipgloss.Color("#")
)

var (
	List     = lipgloss.NewStyle().Margin(0, 1)
	HelpKey  = lipgloss.NewStyle().Foreground(ColorSecondary)
	HelpDesc = lipgloss.NewStyle().Foreground(ColorMuted1)
	Help     = lipgloss.NewStyle().Margin(0, 1, 0, 3)
)
