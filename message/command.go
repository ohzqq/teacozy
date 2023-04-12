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

type UpMsg struct {
	Lines int
}

type DownMsg struct {
	Lines int
}

type NextMsg struct{}
type PrevMsg struct{}
type TopMsg struct{}
type BottomMsg struct{}

func Up(l ...int) tea.Cmd {
	return func() tea.Msg {
		lines := 1
		if len(l) > 0 {
			lines = l[0]
		}
		return UpMsg{
			Lines: lines,
		}
	}
}

func Down(l ...int) tea.Cmd {
	return func() tea.Msg {
		lines := 1
		if len(l) > 0 {
			lines = l[0]
		}
		return DownMsg{
			Lines: lines,
		}
	}
}

func Next() tea.Cmd {
	return func() tea.Msg {
		return NextMsg{}
	}
}

func Prev() tea.Cmd {
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
