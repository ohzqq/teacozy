package item

import tea "github.com/charmbracelet/bubbletea"

type itemIndexMsg struct {
	Index int
}

type ToggleItemListMsg struct{ itemIndexMsg }

func ToggleItemListCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleItemListMsg{
			itemIndexMsg: itemIndexMsg{Index: idx},
		}
	}
}

type ToggleSelectedItemMsg struct{ itemIndexMsg }

func ToggleSelectedItemCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleSelectedItemMsg{
			itemIndexMsg: itemIndexMsg{Index: idx},
		}
	}
}

type UpdateItemInfoMsg struct{ Info string }

func UpdateItemInfoCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return UpdateItemInfoMsg{Info: content}
	}
}

type EditItemMsg struct{ itemIndexMsg }

func EditItemCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return EditItemMsg{
			itemIndexMsg: itemIndexMsg{Index: idx},
		}
	}
}
