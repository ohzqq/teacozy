package tui

import (
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/menu"
)

func NewHelp() *menu.Menu {
	m := menu.New("?", "help")
	return m
}

func ListKeyMap() key.KeyMap {
	lk := list.ListKeyMap()
	km := key.NewKeyMap()
	km.AddBind(lk.CursorUp)
	km.AddBind(lk.CursorDown)
	return km
}
