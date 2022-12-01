package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
)

type Help struct {
	*info.Info
	view    viewport.Model
	KeyMap  key.KeyMap
	ListNav key.KeyMap
}

func NewHelp() Help {
	i := info.New()
	i.Show()
	i.Toggle = key.NewKey("?", "help")
	km := KeyMap()
	//km := key.NewKeyMap()
	//km.Add(i.Toggle)
	h := Help{
		Info:    i,
		KeyMap:  km,
		ListNav: ListKeyMap(),
	}
	return h
}

func (h *Help) Render() string {
	h.NewSection().SetTitle("Help").SetFields(h.KeyMap)
	h.NewSection().SetTitle("Navigation").SetFields(h.ListNav)
	return h.Info.Render()
}

//func (m *Help) Update(msg tea.Msg) (*Help, tea.Cmd) {
//}

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
