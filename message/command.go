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

type LineUpMsg struct{}
type HalfPageUpMsg struct{}
type PageUpMsg struct{}

func LineUp() tea.Msg     { return LineUpMsg{} }
func HalfPageUp() tea.Msg { return HalfPageUpMsg{} }
func PageUp() tea.Msg     { return PageUpMsg{} }

type LineDownMsg struct{}
type HalfPageDownMsg struct{}
type PageDownMsg struct{}

func LineDown() tea.Msg     { return LineDownMsg{} }
func HalfPageDown() tea.Msg { return HalfPageDownMsg{} }
func PageDown() tea.Msg     { return PageDownMsg{} }

type NextMsg struct{}
type PrevMsg struct{}

func NextPage() tea.Msg { return NextMsg{} }
func PrevPage() tea.Msg { return PrevMsg{} }

type TopMsg struct{}
type BottomMsg struct{}

func Top() tea.Msg    { return TopMsg{} }
func Bottom() tea.Msg { return BottomMsg{} }

type StartFilteringMsg struct{}

func StartFiltering() tea.Cmd {
	return func() tea.Msg {
		return StartFilteringMsg{}
	}
}

type StopFilteringMsg struct{}

func StopFiltering() tea.Cmd {
	return func() tea.Msg {
		return StopFilteringMsg{}
	}
}
