package info

import (
	tea "github.com/charmbracelet/bubbletea"
)

// info commands
type ToggleVisibleMsg struct{}

func ToggleVisibleCmd() tea.Cmd {
	return func() tea.Msg {
		return ToggleVisibleMsg{}
	}
}

type HideInfoMsg struct{}

func HideInfoCmd() tea.Cmd {
	return func() tea.Msg {
		return HideInfoMsg{}
	}
}

type UpdateContentMsg struct {
	Content string
}

func UpdateContentCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return UpdateContentMsg{Content: content}
	}
}

type ShowInfoMsg struct{}

func ShowInfoCmd() tea.Cmd {
	return func() tea.Msg {
		return ShowInfoMsg{}
	}
}

type EditInfoMsg struct {
	//Data teacozy.FormData
}

func EditInfoCmd() tea.Cmd {
	return func() tea.Msg {
		return EditInfoMsg{}
	}
}
