package teacozy

import (
	tea "github.com/charmbracelet/bubbletea"
)

type DefaultMenus struct {
	Menus Menus
	*Menu
}

func DefaultTuiMenus() DefaultMenus {
	menu := DefaultMenus{
		Menus: make(Menus),
		Menu:  DefaultMenu().SetToggle("?", "Menu").SetLabel("Menu"),
	}
	sort := SortListMenu()
	menu.Menus.Set("sort list", sort)

	fn := func(ui *TUI) tea.Cmd {
		return UpdateStatusCmd("poot")
	}
	menu.NewKey("o", "sort list", fn)
	menu.Menus.Set("default", menu.Menu)
	return menu
}

func GoToMenuCmd(m *Menu) MenuFunc {
	return func(ui *TUI) tea.Cmd {
		//ui.CurrentMenu = m
		//ui.ShowMenu()
		//ui.HideInfo()
		//fmt.Println(ui.CurrentMenu.Label)
		//return UpdateStatusCmd(m.Label)
		//return SetFocusedViewCmd(ui.CurrentMenu.Label)
		return ChangeMenuCmd(m)
	}
}

func SortListMenu() *Menu {
	m := DefaultMenu().SetToggle("o", "sort list").SetLabel("sort list")
	fn := func(ui *TUI) tea.Cmd {
		return UpdateStatusCmd(m.Label)
	}
	m.NewKey("t", "test", fn)
	return m
}
