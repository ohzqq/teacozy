package teacozy

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func ConfirmMenu() *Menu {
	m := DefaultMenu().SetToggle("a", "action").SetLabel("action")
	m.NewKey("y", "yes", SaveChanges(true))
	m.NewKey("n", "no", SaveChanges(false))
	content := "save changes? y/n"
	style := lipgloss.NewStyle().
		Foreground(DefaultColors().Black).
		Background(DefaultColors().Red)
	m.SetContent(style.Render(content))
	return m
}

func SaveChanges(confirm bool) MenuFunc {
	return func(m tea.Model) tea.Cmd {
		return func() tea.Msg {
			return ConfirmMenuMsg(confirm)
		}
	}
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
