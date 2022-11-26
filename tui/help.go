package tui

import (
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
)

type Help struct {
	info     *info.Info
	ListKeys key.KeyMap
	key.KeyMap
}

func NewHelp() Help {
	h := Help{
		KeyMap: key.NewKeyMap(),
	}
	h.info = info.New(h.KeyMap)
	return h
}

func (h Help) Info() *info.Info {
	h.ListNavigation()
	return h.info
}

func (h *Help) ListNavigation() {
	lk := list.ListKeyMap()
	km := key.NewKeyMap()
	km.AddBind(lk.CursorUp)
	km.AddBind(lk.CursorDown)
	h.info.AddContent("List Nav")
	h.info.AddFields(km)
}
