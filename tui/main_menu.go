package tui

import "github.com/ohzqq/teacozy/key"

type MainMenu struct {
	*Menu
}

func NewMainMenu() MainMenu {
	mk := key.NewKey("m", "menu")
	m := MainMenu{
		Menu: NewMenu(mk),
	}

	return m
}
