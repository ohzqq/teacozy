package teacozy

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type Help struct {
	help.Model
	keys []Key
}

func NewHelp(keys ...Key) Help {
	m := help.New()
	m.ShowAll = false
	return Help{
		Model: m,
		keys:  keys,
	}
}

func (h Help) ShortHelp() []key.Binding {
	var keys []key.Binding
	for _, key := range h.keys {
		keys = append(keys, key.Bind)
	}
	return keys
}

func (h Help) FullHelp() [][]key.Binding {
	var keys [][]key.Binding
	keys = append(keys, h.ShortHelp())
	return keys
}

func (h Help) View() string {
	return h.Model.View(h)
}
