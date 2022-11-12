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

type MenuCmd func(m *UI) tea.Cmd

// form commands
type SaveAsHashMsg struct{}

func SaveAsHashCmd() tea.Cmd {
	return func() tea.Msg {
		return SaveAsHashMsg{}
	}
}

type EditInfoMsg struct{}

func EditInfoCmd() tea.Cmd {
	return func() tea.Msg {
		return EditInfoMsg{}
	}
}

type EditItemMsg struct {
	Field
}

func EditItemCmd(field Field) tea.Cmd {
	return func() tea.Msg {
		return EditItemMsg{Field: field}
	}
}

type UpdateContentMsg struct {
	Field
}

func UpdateContentCmd(key, val string) tea.Cmd {
	return func() tea.Msg {
		return UpdateContentMsg{Field: NewDefaultField(key, val)}
	}
}

// info commands
type HideMsg struct{}

func HideCmd() tea.Cmd {
	return func() tea.Msg {
		return HideMsg{}
	}
}

type ShowMsg struct{}

func ShowCmd() tea.Cmd {
	return func() tea.Msg {
		return ShowMsg{}
	}
}

type EditMsg struct {
	Fields *Fields
}

func EditCmd(f *Fields) tea.Cmd {
	return func() tea.Msg {
		return EditMsg{Fields: f}
	}
}

//item commands

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
