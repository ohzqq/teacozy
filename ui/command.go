package ui

import tea "github.com/charmbracelet/bubbletea"

type SetFocusedViewMsg string

func SetFocusedViewCmd(v string) tea.Cmd {
	return func() tea.Msg {
		return SetFocusedViewMsg(v)
	}
}

type UpdateMenuContentMsg string

func UpdateMenuContentCmd(s string) tea.Cmd {
	return func() tea.Msg {
		return UpdateMenuContentMsg(s)
	}
}

type MenuCmd func(m *UI) tea.Cmd

type UpdateStatusMsg string

func UpdateStatusCmd(status string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusMsg(status)
	}
}
