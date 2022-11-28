package teacozy

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateStatusMsg struct{ Msg string }

func UpdateStatusCmd(status string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusMsg{Msg: status}
	}
}

type SetListItemMsg struct {
	Item list.Item
}

func SetListItemCmd(item list.Item) tea.Cmd {
	return func() tea.Msg {
		return SetListItemMsg{Item: item}
	}
}
