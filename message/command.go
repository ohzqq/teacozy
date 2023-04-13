package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ReturnSelectionsMsg struct {
	Quitting bool
}

func ReturnSelections() tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

func ReturnSels(limit, numSel int) tea.Cmd {
	msg := ReturnSelectionsMsg{}
	if limit == 1 {
		return ToggleItem()
	}
	if numSel == 0 {
		msg.Quitting = true
		return ToggleItem()
	}
	return func() tea.Msg {
		return msg
	}
}

type QuitMsg struct{}

func Quit() tea.Cmd {
	return func() tea.Msg {
		return QuitMsg{}
	}
}

type ChangeRouteMsg struct {
	Name string
}

type GetConfirmationMsg struct {
	Question string
}

func GetConfirmation(q string) tea.Cmd {
	return func() tea.Msg {
		return GetConfirmationMsg{Question: q}
	}
}

type ConfirmMsg struct {
	Confirmed bool
}

func Confirm(confirm bool) tea.Cmd {
	return func() tea.Msg {
		return ConfirmMsg{Confirmed: confirm}
	}
}

func ChangeRoute(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeRouteMsg{Name: name}
	}
}

type ToggleItemMsg struct{}

func ToggleItem() tea.Cmd {
	return func() tea.Msg {
		return ToggleItemMsg{}
	}
}

type ToggleAllItemsMsg struct{}

func ToggleAllItems() tea.Cmd {
	return func() tea.Msg {
		return ToggleAllItemsMsg{}
	}
}

type ShowHelpMsg struct{}

func ShowHelp() tea.Cmd {
	return func() tea.Msg {
		return ShowHelpMsg{}
	}
}

type HideHelpMsg struct{}

func HideHelp() tea.Cmd {
	return func() tea.Msg {
		return HideHelpMsg{}
	}
}

type LineUpMsg struct {
	Lines int
}

type LineDownMsg struct {
	Lines int
}

type PageDownMsg struct{}
type HalfPageDownMsg struct{}
type PageUpMsg struct{}
type HalfPageUpMsg struct{}

func PageDown() tea.Msg {
	return PageDownMsg{}
}

func PageUp() tea.Msg {
	return PageUpMsg{}
}

func HalfPageDown() tea.Msg {
	return HalfPageDownMsg{}
}

func HalfPageUp() tea.Msg {
	return HalfPageUpMsg{}
}

type NextMsg struct{}
type PrevMsg struct{}
type TopMsg struct{}
type BottomMsg struct{}

func LineUp() tea.Msg {
	return LineUpMsg{}
}

func LineDown() tea.Msg {
	return LineDownMsg{}
}

func NextPage() tea.Cmd {
	return func() tea.Msg {
		return NextMsg{}
	}
}

func PrevPage() tea.Cmd {
	return func() tea.Msg {
		return PrevMsg{}
	}
}

func Top() tea.Cmd {
	return func() tea.Msg {
		return TopMsg{}
	}
}

func Bottom() tea.Cmd {
	return func() tea.Msg {
		return BottomMsg{}
	}
}
