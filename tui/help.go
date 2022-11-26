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
		Model:    m,
		KeyMap:   key.NewKeyMap(),
		ListKeys: key.NewKeyMap(),
	}
	h.Info = info.New(h.KeyMap)
	h.ListNavigation()
	h.Info.AddContent("List Nav")
	h.Info.SetData(h.ListKeys)
	return h
}

func (h *Help) ListNavigation() {
	lk := list.ListKeyMap()
	h.ListKeys.AddBind(lk.CursorUp)
	h.ListKeys.AddBind(lk.CursorDown)
}
