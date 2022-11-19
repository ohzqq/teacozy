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
	menu.Menus.Set("sort by", sort)

	fn := func(ui *TUI) tea.Cmd {
		return UpdateStatusCmd("poot")
	}
	menu.NewKey("o", "sort by", fn)
	menu.Menus.Set("default", menu.Menu)
	return menu
}

func GoToMenuCmd(m *Menu) MenuFunc {
	return func(ui *TUI) tea.Cmd {
		return ChangeMenuCmd(m)
	}
}

func SortListMenu() *Menu {
	m := DefaultMenu().SetToggle("o", "sort by").SetLabel("sort by")
	//fn := func(ui *TUI) tea.Cmd {
	//  return UpdateStatusCmd("poot")
	//}
	return m
}

//func SortListByKey(ui *TUI) tea.Cmd {
//  items := ui.Main.List.Items.All()
//  sort.SliceStable(items,
//    func(i, j int) bool {
//      return items[i].Key() < items[j].Key()
//    },
//  )
//}
