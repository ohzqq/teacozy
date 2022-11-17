package teacozy

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func SortListMenu() *Menu {
	t := key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "sort list"),
	)
	testHelpKeys := []MenuItem{
		NewMenuItem("t", "select item", UiTestKeyAction),
		NewMenuItem("o", "deselect item", UiTestKeyAction),
	}
	m := NewMenu("test", t, testHelpKeys...)
	return m
}

func UiTestKeyAction(m *TUI) tea.Cmd {
	return UpdateStatusCmd(fmt.Sprintf("%v", "poot"))
}
