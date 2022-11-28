package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
)

type Help struct {
	*info.Info
	key.KeyMap
}

func NewHelp() Help {
	i := info.New()
	i.Toggle = key.NewKey("?", "help")
	km := KeyMap()
	//km := key.NewKeyMap()
	//km.Add(i.Toggle)
	return Help{
		Info:   i,
		KeyMap: km,
	}
}

func (h *Help) Render() {
	h.NewSection().SetTitle("Help").SetFields(h.KeyMap)
	h.NewSection().SetTitle("Navigation").SetFields(ListKeyMap())
}

func GoToHelp(m tea.Model) tea.Cmd {
	if ui, ok := m.(*TUI); ok {
		return ui.ShowHelp()
	}
	return nil
}

func ListKeyMap() key.KeyMap {
	lk := list.ListKeyMap()
	km := key.NewKeyMap()
	km.AddBind(lk.CursorUp)
	km.AddBind(lk.CursorDown)
	return km
}
