package teacozy

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ui commands
type SetFocusedViewMsg string

func SetFocusedViewCmd(v string) tea.Cmd {
	return func() tea.Msg {
		return SetFocusedViewMsg(v)
	}
}

type UpdateMenuContentMsg string

func UpdateMenuContentCmd(s string) tea.Cmd {
	return func() tea.Msg {
		return UpdateMenuContentMsg(s)
	}
}

type MenuFunc func(m *TUI) tea.Cmd

// form commands
type SaveFormAsHashMsg struct{}

func SaveFormAsHashCmd() tea.Cmd {
	return func() tea.Msg {
		return SaveFormAsHashMsg{}
	}
}

type EditFormMsg struct{}

func EditFormCmd() tea.Cmd {
	return func() tea.Msg {
		return EditFormMsg{}
	}
}

type EditFormItemMsg struct {
	Data FieldData
	*Item
}

func EditFormItemCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return EditFormItemMsg{Data: item.Data, Item: item}
	}
}

type UpdateFormContentMsg struct {
	*Field
}

func UpdateFormContentCmd(key, val string) tea.Cmd {
	return func() tea.Msg {
		return UpdateFormContentMsg{Field: NewField(key, val)}
	}
}

// info commands
type HideInfoMsg struct{}

func HideInfoCmd() tea.Cmd {
	return func() tea.Msg {
		return HideInfoMsg{}
	}
}

type ShowInfoMsg struct{}

func ShowInfoCmd() tea.Cmd {
	return func() tea.Msg {
		return ShowInfoMsg{}
	}
}

type EditInfoMsg struct {
	Fields *Fields
}

func EditInfoCmd(f *Fields) tea.Cmd {
	return func() tea.Msg {
		return EditInfoMsg{Fields: f}
	}
}

//item commands

type ToggleItemListMsg struct{ *Item }

func ToggleItemListCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return ToggleItemListMsg{Item: item}
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

// list commands

type ReturnSelectionsMsg struct{}

func ReturnSelectionsCmd() tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
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

type UpdateStatusMsg struct{ Msg string }

func UpdateStatusCmd(status string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusMsg{Msg: status}
	}
}

type SetSizeMsg []int

func SetSizeCmd(size []int) tea.Cmd {
	return func() tea.Msg {
		return SetSizeMsg(size)
	}
}
