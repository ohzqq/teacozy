package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
)

type Help struct {
	*info.Info
	KeyMap  key.KeyMap
	ListNav key.KeyMap
}

func NewHelp() Help {
	h := Help{
		Info:    info.New(),
		KeyMap:  KeyMap(),
		ListNav: ListKeyMap(),
	}
	h.Show()
	h.Toggle = KeyMap().GetKey("?")
	h.NewSection().SetTitle("Help").SetFields(h.KeyMap)
	h.NewSection().SetTitle("Navigation").SetFields(h.ListNav)
	return h
}

func GoToHelpView(m tea.Model) tea.Cmd {
	if ui, ok := m.(*Tui); ok {
		ui.state = helpModel
		ui.showFullHelp = true
		ui.view = ui.Help.Info.Model
		return info.UpdateContentCmd(ui.Help.Render())
	}
	return nil
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
