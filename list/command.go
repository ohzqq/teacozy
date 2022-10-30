package list

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

func MultiSelectAction(m List) func(List) tea.Cmd {
	fn := func(m List) tea.Cmd {
		for _, item := range m.Items.All {
			i := item.(Item)
			m.Items.Selected[i.Idx] = item
		}
		if m.Items.HasSelections() {
			return tea.Quit
		}
		return nil
	}
	return fn
}

func SingleSelectAction(m List) func(List) tea.Cmd {
	fn := func(m List) tea.Cmd {
		cur := m.Model.SelectedItem()
		curItem := cur.(Item)
		curItem.IsSelected = true
		m.Items.Selected[curItem.Idx] = cur
		if m.Items.HasSelections() {
			return tea.Quit
		}
		return nil
	}
	return fn
}
