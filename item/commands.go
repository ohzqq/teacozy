package item

import tea "github.com/charmbracelet/bubbletea"

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

type ShowInfoMsg struct{ *Item }

func ShowInfoCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return ShowInfoMsg{Item: item}
	}
}

type EditContentMsg struct{ *Item }

func EditContentCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return EditContentMsg{Item: item}
	}
}
