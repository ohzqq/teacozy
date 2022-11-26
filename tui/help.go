package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
)

type Help struct {
	help.Model
	*info.Info
	keys     []key.Key
	ListKeys key.KeyMap
	KeyMap   key.KeyMap
}

func NewHelp() Help {
	m := help.New()
	m.ShowAll = false
	h := Help{
		Model:  m,
		KeyMap: key.NewKeyMap(),
	}
	h.Info = info.New(h.KeyMap)
	h.ListNavigation()
	return h
}

func (h *Help) ListNavigation() {
	lk := list.ListKeyMap()
	km := key.NewKeyMap()
	km.AddBind(lk.CursorUp)
	km.AddBind(lk.CursorDown)
	h.Info.AddContent("List Nav")
	h.Info.AddFields(km)
}
