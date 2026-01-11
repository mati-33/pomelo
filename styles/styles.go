package styles

import (
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	Grey   = gloss.Color("245")
	Accent = gloss.Color("56")
	White  = gloss.Color("#fff")
	Black  = gloss.Color("#000")

	LogoStyle      = gloss.NewStyle().Align(gloss.Center)
	ListStyle      = gloss.NewStyle().Margin(1, 1).Padding(0, 1)
	ListFocusStyle = gloss.NewStyle().Margin(1, 0).Border(gloss.NormalBorder(), false, false, false, true).BorderLeftForeground(Accent)
	ListsStyle     = gloss.NewStyle().Align(gloss.Center)
	InfoStyle      = gloss.NewStyle().Foreground(Grey)
	NameStyle      = gloss.NewStyle().Foreground(White)
	NotifSuccess   = gloss.NewStyle().Background(gloss.Color("35")).Foreground(Black).Padding(0, 1).Margin(0, 1)
)
