package item

import tea "github.com/charmbracelet/bubbletea"

type ToggleItemListMsg int

func ToggleItemListCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleItemListMsg(idx)
	}
}

type ToggleSelectedItemMsg int

func ToggleSelectedItemCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleSelectedItemMsg(idx)
	}
}

type UpdateInfoContentMsg string

func UpdateInfoContentCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return UpdateInfoContentMsg(content)
	}
}

type EditItemMsg string

func EditItemCmd() tea.Cmd {
	return func() tea.Msg {
		return EditItemMsg("")
	}
}
