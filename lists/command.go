package lists

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ListAction func(m List) tea.Cmd

type UpdateVisibleItemsMsg string

func UpdateVisibleItemsCmd(m string) tea.Cmd {
	return func() tea.Msg {
		return UpdateVisibleItemsMsg(m)
	}
}

type ToggleItemMsg struct{}

func ToggleItemCmd() tea.Cmd {
	return func() tea.Msg {
		return ToggleItemMsg{}
	}
}

type SetFocusedViewMsg string

func SetFocusedViewCmd(v string) tea.Cmd {
	return func() tea.Msg {
		return SetFocusedViewMsg(v)
	}
}

func FocusListCmd() tea.Cmd {
	return func() tea.Msg {
		return SetFocusedViewMsg("list")
	}
}

type EditItemMsg string

func EditItemCmd() tea.Cmd {
	return func() tea.Msg {
		return EditItemMsg("")
	}
}
