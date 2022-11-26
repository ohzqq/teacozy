package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// list commands
type ReturnSelectionsMsg struct{}

func ReturnSelectionsCmd() tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

type ExitSelectionsListMsg struct{}

func (m *List) ExitSelectionsListCmd() tea.Cmd {
	return func() tea.Msg {
		m.SelectionList = false
		return ExitSelectionsListMsg{}
	}
}

func ToggleAllItemsCmd(l *List) {
	l.Items.ToggleAllSelectedItems()
}

type UpdateVisibleItemsMsg string

func UpdateVisibleItemsCmd(opt string) tea.Cmd {
	return func() tea.Msg {
		return UpdateVisibleItemsMsg(opt)
	}
}

func (m *List) ShowVisibleItemsCmd() tea.Cmd {
	return func() tea.Msg {
		return UpdateVisibleItemsMsg("visible")
	}
}

func (m *List) ShowSelectedItemsCmd() tea.Cmd {
	return func() tea.Msg {
		return UpdateVisibleItemsMsg("selected")
	}
}

type UpdateStatusMsg struct{ Msg string }

func UpdateStatusCmd(status string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusMsg{Msg: status}
	}
}

type SortItemsMsg struct{ Items []*Item }

func SortItemsCmd(items []*Item) tea.Cmd {
	return func() tea.Msg {
		return SortItemsMsg{Items: items}
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

type SetItemMsg struct{ *Item }

func SetItemCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return SetItemMsg{Item: item}
	}
}

type SetItemsMsg struct{ Items []list.Item }

func SetItemsCmd(items []list.Item) tea.Cmd {
	return func() tea.Msg {
		return SetItemsMsg{Items: items}
	}
}

// item commands
type ToggleItemChildrenMsg struct{ *Item }

func ToggleItemChildrenCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return ToggleItemChildrenMsg{Item: item}
	}
}

type ToggleSelectedItemMsg struct{ *Item }

func ToggleSelectedItemCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return ToggleSelectedItemMsg{Item: item}
	}
}

type ShowItemInfoMsg struct{ *Item }

func ShowItemInfoCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return ShowItemInfoMsg{Item: item}
	}
}

type EditItemValueMsg struct{ *Item }

func EditItemValueCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return EditItemValueMsg{Item: item}
	}
}
