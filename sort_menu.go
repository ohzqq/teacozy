package teacozy

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func DefaultTuiMenus() Menus {
	menus := make(Menus)
	menus.Set("sort list", SortListMenu())
	return menus
}

func SortListMenu() *Menu {
	m := DefaultMenu().
		SetToggle("o", "sort list")
	return m
}

func UiTestKeyAction(m *TUI) tea.Cmd {
	return UpdateStatusCmd(fmt.Sprintf("%v", "poot"))
}
