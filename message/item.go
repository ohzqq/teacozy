package message

import tea "github.com/charmbracelet/bubbletea"

type UpdateItem struct {
	Cmd func(int) tea.Cmd
}

func EditItem() tea.Msg {
	fn := func(int) tea.Cmd {
		return func() tea.Msg {
			return EditItemMsg{
				Index: idx,
			}
		}
	}
	return UpdateItemMsg{
		Cmd: fn,
	}
}
