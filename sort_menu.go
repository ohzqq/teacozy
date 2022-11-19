package teacozy

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
)

type DefaultMenus struct {
	Menus Menus
	*Menu
}

func GoToMenuCmd(m *Menu) MenuFunc {
	return func(ui *TUI) tea.Cmd {
		return ChangeMenuCmd(m)
	}
}

func SortListMenu() *Menu {
	m := DefaultMenu().SetToggle("o", "sort by").SetLabel("sort by")
	m.NewKey("k", "key (asc)", SortListByKey("asc"))
	m.NewKey("K", "key (desc)", SortListByKey("desc"))
	m.NewKey("v", "value (asc)", SortListByValue("asc"))
	m.NewKey("V", "value (desc)", SortListByValue("desc"))
	return m
}

func SortListByValue(order string) MenuFunc {
	return func(ui *TUI) tea.Cmd {
		items := ui.Main.Items.All()
		sort.SliceStable(items,
			func(i, j int) bool {
				if order == "asc" {
					return items[i].Value() < items[j].Value()
				}
				return items[j].Value() < items[i].Value()
			},
		)
		return SortItemsCmd(items)
	}
}
func SortListByKey(order string) MenuFunc {
	return func(ui *TUI) tea.Cmd {
		items := ui.Main.Items.All()
		sort.SliceStable(items,
			func(i, j int) bool {
				if order == "asc" {
					return items[i].Key() < items[j].Key()
				}
				return items[j].Key() < items[i].Key()
			},
		)
		return SortItemsCmd(items)
	}
}
