package tui

import (
	"github.com/charmbracelet/bubbles/help"
	bubblekey "github.com/charmbracelet/bubbles/key"
	"github.com/ohzqq/teacozy/key"
)

type Help struct {
	help.Model
	keys []key.Key
}

func NewHelp(keys ...key.Key) Help {
	m := help.New()
	m.ShowAll = false
	return Help{
		Model: m,
		keys:  keys,
	}
}

func (h Help) ShortHelp() []bubblekey.Binding {
	var keys []bubblekey.Binding
	for _, key := range h.keys {
		keys = append(keys, bubblekey.Binding)
	}
	return keys
}

func (h Help) FullHelp() [][]bubblekey.Binding {
	var keys [][]bubblekey.Binding
	keys = append(keys, h.ShortHelp())
	return keys
}

func (h Help) View() string {
	return h.Model.View(h)
}
