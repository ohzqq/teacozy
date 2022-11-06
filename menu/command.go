package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

type CmdFunc func() tea.Cmd

type UpdateMenuContentMsg struct {
	Content string
}

func UpdateMenuContentCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return UpdateMenuContentMsg{Content: content}
	}
}
