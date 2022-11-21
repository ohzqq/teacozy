package teacozy

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func ActionMenu() *Menu {
	m := DefaultMenu().SetToggle("a", "action").SetLabel("action")
	m.NewKey("P", "print to stdout", PrintItemsMenuFunc)
	return m
}

func PrintItems(items ...*Item) tea.Cmd {
	for _, i := range items {
		fmt.Println(i.String())
	}
	//return nil
	return tea.Quit
}

func PrintItemsMenuFunc(m tea.Model) tea.Cmd {
	ui := m.(*TUI)
	main := ui.Main.(*List)
	main.SetAction(PrintItems)
	return ReturnSelectionsCmd()
}
