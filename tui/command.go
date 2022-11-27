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

func (ui *TUI) ToggleFullScreenCmd() tea.Cmd {
	return func() tea.Msg {
		if ui.fullScreen {
			return tea.EnterAltScreen()
		}
		return tea.ExitAltScreen()
	}
}

// menu commands
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
	*teacozy.FormData
}

func EditInfoCmd(f *teacozy.FormData) tea.Cmd {
	return func() tea.Msg {
		return EditInfoMsg{
			FormData: f,
		}
	}
}

type ShowItemInfoMsg struct{ *list.Item }

func ShowItemInfoCmd(item *list.Item) tea.Cmd {

	return func() tea.Msg {
		return ShowItemInfoMsg{Item: item}
	}
}
