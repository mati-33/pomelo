package models

import (
	"github.com/charmbracelet/bubbletea"
)

type mode int

const (
	listMode mode = iota
	addMode
	deleteMode
	modifyMode
)

type PopScreenMsg struct {
	cmd tea.Cmd
}

type PushScreenMsg struct {
	Screen tea.Model
	Cmd    tea.Cmd
}
