package item

import tea "github.com/charmbracelet/bubbletea"

type itemIndexMsg struct {
	Index int
}

type ToggleListMsg struct{ *Item }

func ToggleListCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return ToggleListMsg{Item: item}
	}
}

type ToggleSelectedMsg struct{ *Item }

func ToggleSelectedCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return ToggleSelectedMsg{Item: item}
	}
}

type UpdateInfoMsg struct{ Info string }

func UpdateInfoCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return UpdateInfoMsg{Info: content}
	}
}

type EditContentMsg struct{ *Item }

func EditContentCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return EditContentMsg{Item: item}
	}
}
