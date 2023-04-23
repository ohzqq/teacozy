package keys

import tea "github.com/charmbracelet/bubbletea"

// Route keys
type ReturnToListMsg struct{}

func ReturnToList() tea.Msg { return ReturnToListMsg{} }

type ChangeRouteMsg struct {
	Name string
}

func ChangeRoute(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeRouteMsg{Name: name}
	}
}

type UpdateItemMsg struct {
	Cmd func(int) tea.Cmd
}

func UpdateItem(cmd func(int) tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return UpdateItemMsg{
			Cmd: cmd,
		}
	}
}

type ToggleItemsMsg struct {
	Index int
}

func ToggleItems(idx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleItemsMsg{
			Index: idx,
		}
	}
}

// List Msg
type ReturnSelectionsMsg struct{}

func ReturnSelections() tea.Msg { return ReturnSelectionsMsg{} }

type ToggleItemMsg struct{}
type ToggleAllItemsMsg struct{}

func ToggleAllItems() tea.Msg { return ToggleAllItemsMsg{} }
func ToggleItem() tea.Msg     { return ToggleItemMsg{} }

// Help msg
type ShowHelpMsg struct{}

func ShowHelp() tea.Msg { return ShowHelpMsg{} }

// Nav msg
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

// formj

type EditItemMsg struct {
	Index int
}
type StartEditingMsg struct{}
type StopEditingMsg struct{}
type SaveEditMsg struct{}
type ConfirmEditMsg struct{}

func SaveChanges() tea.Msg {
	return SaveEditMsg{}
}

func SaveEdit(save bool) tea.Cmd {
	if save {
		return SaveChanges
	}
	return ReturnToList
}

func ConfirmEdit() tea.Msg {
	return ConfirmEditMsg{}
}

func StopEditing() tea.Msg {
	return StopEditingMsg{}
}

func EditItem(idx int) tea.Cmd {
	return func() tea.Msg {
		return EditItemMsg{
			Index: idx,
		}
	}
}

func StartEditing() tea.Msg {
	return StartEditingMsg{}
}
