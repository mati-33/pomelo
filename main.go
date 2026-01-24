package main

import (
	"errors"
	"fmt"
	"os"

	"pomelo/data"
	"pomelo/models"

	"github.com/charmbracelet/bubbletea"
)

func main() {
	path := os.Getenv("HOME") + "/.local/state/pomelo.db"

	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(path)
		if err != nil {
			quit(err)
		}
	} else if err != nil {
		quit(err)
	}

	db, err := data.SetupDB(path)
	if err != nil {
		quit(err)
	}
	defer db.Close()

	err = data.SetupTables(db)
	if err != nil {
		quit(err)
	}

	// lr := data.ListsRepo{Db: db}
	// mocks := []struct{ name, desc string }{
	// 	{name: "pomelo", desc: "current projext"},
	// 	{name: "06-10-2026", desc: ""},
	// 	{name: "07-10-2026", desc: ""},
	// 	{name: "terminal typing app", desc: "something like keybr but for terminal"},
	// 	{name: "terminal chat app", desc: "something like sack but better"},
	// }
	//
	// for _, m := range mocks {
	// 	_, _ = lr.Add(m.name, m.desc)
	// }

	m := models.InitialPomeloModel(db)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		quit(err)
	}
}

func quit(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
	os.Exit(1)
}
