package teacozy

import tea "github.com/charmbracelet/bubbletea"

func ConfirmMenu() *Menu {
	m := DefaultMenu().SetToggle("a", "action").SetLabel("action")
	m.NewKey("y", "yes", ConfirmOrDeny(true))
	m.NewKey("n", "no", ConfirmOrDeny(false))
	return m
}

type ConfirmOrDenyMsg bool

func ConfirmOrDeny(confirm bool) MenuFunc {
	return func(m tea.Model) tea.Cmd {
		ui := m.(*TUI)
		if confirm {
			ui.actionConfirmed = true
		}
		return nil
	}
}
