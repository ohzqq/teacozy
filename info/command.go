package info

import tea "github.com/charmbracelet/bubbletea"

type HideMsg struct{}

func HideCmd() tea.Cmd {
	return func() tea.Msg {
		return HideMsg{}
	}
}

type ShowMsg struct{}

func ShowCmd() tea.Cmd {
	return func() tea.Msg {
		return ShowMsg{}
	}
}

type EditMsg struct {
	Fields *Fields
}

func EditCmd(f *Fields) tea.Cmd {
	return func() tea.Msg {
		return EditMsg{Fields: f}
	}
}
