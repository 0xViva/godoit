package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Add    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Toggle key.Binding
	Quit   key.Binding
	Enter  key.Binding
	Esc    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Add, k.Edit, k.Delete, k.Toggle, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Add, k.Edit, k.Delete, k.Toggle, k.Enter, k.Esc, k.Quit},
	}
}

var keys = keyMap{
	Up:     key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/↑", "")),
	Down:   key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/↓", "")),
	Add:    key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
	Edit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Delete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	Toggle: key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "toggle")),
	Enter:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
	Esc:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	Quit:   key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q/ctrl+c", "quit")),
}
