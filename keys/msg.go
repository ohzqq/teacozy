package keys

import tea "github.com/charmbracelet/bubbletea"

type ReturnSelectionsMsg struct {
	Quitting bool
}

func ReturnSelections() tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

type QuitMsg struct{}

func QuitTea() tea.Cmd {
	return func() tea.Msg {
		return QuitMsg{}
	}
}

type ChangeRouteMsg struct {
	Name string
}

func ChangeRoute(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeRouteMsg{Name: name}
	}
}

//type ToggleItemMsg struct{}

//func ToggleItem() tea.Cmd {
//return func() tea.Msg {
//return ToggleItemMsg{}
//}
//}

type ToggleAllItemsMsg struct{}

func ToggleAllItems() tea.Msg {
	return ToggleAllItemsMsg{}
}

type DeselectAllItemsMsg struct{}

func DeselectAllItems() tea.Msg {
	return DeselectAllItemsMsg{}
}

type SelectAllItemsMsg struct{}

func SelectAllItems() tea.Msg {
	return SelectAllItemsMsg{}
}

//type ShowHelpMsg struct{}

//func ShowHelp() tea.Cmd {
//return func() tea.Msg {
//return ShowHelpMsg{}
//}
//}

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

type ToggleMsg struct {
	Index int
}

func Toggle(idx int) tea.Msg {
	return func() tea.Msg {
		return ToggleMsg{Index: idx}
	}
}
