package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/pager"
)

type KeyMap struct {
	ToggleFocus key.Binding
	Quit        key.Binding
	List        list.KeyMap
	Pager       pager.KeyMap
	Input       textinput.KeyMap
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
	}
}
