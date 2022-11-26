//go:build ignore

package tui

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/menu"
)

func SortListMenu() *Menu {
	kAsc := key.NewKey("k", "key (asc)").SetCmd(SortListByKey("asc"))
	kDesc := key.NewKey("K", "key (desc)").SetCmd(SortListByKey("desc"))
	vAsc := key.NewKey("v", "value (asc)").SetCmd(SortListByValue("asc"))
	vDesc := key.NewKey("V", "value (desc)").SetCmd(SortListByValue("desc"))
	km := key.NewKeyMap()
	km.Add(kAsc)
	km.Add(kDesc)
	km.Add(vAsc)
	km.Add(vDesc)
	m := menu.New("o", "sort by", km)
	return m
}

func SortListByValue(order string) teacozy.MenuFunc {
	return func(m tea.Model) tea.Cmd {
		ui := m.(*TUI)
		main := ui.Main.(*List)
		items := main.Items.All()
		sort.SliceStable(items,
			func(i, j int) bool {
				if order == "asc" {
					return items[i].Content() < items[j].Content()
				}
				return items[j].Content() < items[i].Content()
			},
		)
		return SortItemsCmd(items)
	}
}

func SortListByKey(order string) teacozy.MenuFunc {
	return func(m tea.Model) tea.Cmd {
		ui := m.(*TUI)
		main := ui.Main.(*List)
		items := main.Items.All()
		sort.SliceStable(items,
			func(i, j int) bool {
				if order == "asc" {
					return items[i].Name() < items[j].Name()
				}
				return items[j].Name() < items[i].Name()
			},
		)
		return SortItemsCmd(items)
	}
}
