package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/util"
)

type Help struct {
	*info.Info
	KeyMap  key.KeyMap
	ListNav key.KeyMap
}

func NewHelp() Help {
	i := info.New().SetSize(util.TermWidth()-1, util.TermHeight()-1)
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
	h.NewSection().SetTitle("Help").SetFields(h.KeyMap)
	h.NewSection().SetTitle("Navigation").SetFields(h.ListNav)
	return h
}

//func (m *Help) Update(msg tea.Msg) (*Help, tea.Cmd) {
//}

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
