package teacozy

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func ActionMenu() *Menu {
	m := DefaultMenu().SetToggle("a", "action").SetLabel("action")
	m.NewKey("P", "print to stdout", SortListByValue("desc"))
	return m
}

func PrintItems(items ...*Item) tea.Cmd {
	for _, i := range items {
		fmt.Println(i.String())
	}
	return tea.Quit
}
