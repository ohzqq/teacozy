package message

import tea "github.com/charmbracelet/bubbletea"

type UpdateItem struct {
	Cmd func(int) tea.Cmd
}

type EditItemMsg struct {
	Index int
}

func EditItem() tea.Msg {
	fn := func(idx int) tea.Cmd {
		return func() tea.Msg {
			return EditItemMsg{
				Index: idx,
			}
		}
	}
	return UpdateItem{
		Cmd: fn,
	}
}
