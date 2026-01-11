package components

import (
	"time"

	"github.com/charmbracelet/bubbletea"
	"pomelo/styles"
)

type Notif struct {
	Text      string
	IsVisible bool
}

func (n Notif) Init() tea.Cmd {
	return nil
}

func (n Notif) Update(msg tea.Msg) (Notif, tea.Cmd) {
	switch msg := msg.(type) {

	case notifyMsg:
		n.Text = string(msg)
		n.IsVisible = true
		return n, hideNotif

	case hideNotifMsg:
		n.Text = ""
		n.IsVisible = false
	}
	return n, nil
}

func (n Notif) View() string {
	if !n.IsVisible {
		return ""
	}

	return styles.NotifSuccess.Render(n.Text)
}

func (n Notif) Notify(text string) tea.Cmd {
	return func() tea.Msg {
		return notifyMsg(text)
	}
}

type hideNotifMsg struct{}

type notifyMsg string

func hideNotif() tea.Msg {
	time.Sleep(3 * time.Second)
	return hideNotifMsg{}
}
