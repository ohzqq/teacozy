package app

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	ToggleFocus key.Binding
	Quit        key.Binding
	Command     key.Binding
	FullHelp    key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		ToggleFocus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch panes"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
		Command: key.NewBinding(
			key.WithKeys(":", ";"),
			key.WithHelp(":/;", "enter command"),
		),
		FullHelp: key.NewBinding(
			key.WithKeys("?", "f1"),
			key.WithHelp("?/f1", "help"),
		),
	}
}
