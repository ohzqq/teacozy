package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/list"
)

// ui commands
type SetFocusedViewMsg string

func SetFocusedViewCmd(v string) tea.Cmd {
	return func() tea.Msg {
		return SetFocusedViewMsg(v)
	}
}

type ActionMenuMsg struct{}

func ActionMenuCmd() tea.Cmd {
	return func() tea.Msg {
		return ActionMenuMsg{}
	}
}

type ItemChangedMsg struct {
	*list.Item
}

func ItemChangedCmd(item *list.Item) tea.Cmd {
	return func() tea.Msg {
		return ItemChangedMsg{Item: item}
	}
}

func (ui *Tui) ToggleFullScreenCmd() tea.Cmd {
	return func() tea.Msg {
		if ui.fullScreen {
			return tea.EnterAltScreen()
		}
		return tea.ExitAltScreen()
	}
}

type ToggleHelpMsg struct{}

func ToggleHelpCmd() tea.Cmd {
	return func() tea.Msg {
		return ToggleHelpMsg{}
	}
}

// menu commands
// info commands
type EditInfoMsg struct {
	Fields *teacozy.Fields
}

func EditInfoCmd(f *teacozy.Fields) tea.Cmd {
	return func() tea.Msg {
		return EditInfoMsg{
			Fields: f,
		}
	}
}

type ShowItemInfoMsg struct{ *list.Item }

func ShowItemInfoCmd(item *list.Item) tea.Cmd {
	return func() tea.Msg {
		return ShowItemInfoMsg{Item: item}
	}
}
