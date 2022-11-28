package tui

import "github.com/ohzqq/teacozy/key"

type keyMap struct {
	key.KeyMap
}

func KeyMap() key.KeyMap {
	km := key.NewKeyMap()
	km.AddBind(key.ToggleItem)
	km.AddBind(key.HelpKey)
	km.AddBind(key.Quit)
	km.AddBind(key.SaveAndExit)
	km.AddBind(key.EditField)
	km.AddBind(key.Enter)
	km.AddBind(key.FullScreen)
	km.AddBind(key.InfoKey)
	km.AddBind(key.MenuKey)
	km.AddBind(key.SortList)
	km.AddBind(key.PrevScreen)
	km.AddBind(key.ExitScreen)
	km.AddBind(key.UnToggleAllItems)
	km.AddBind(key.ToggleAllItems)
	km.AddBind(key.ToggleItemList)
	return km
}
