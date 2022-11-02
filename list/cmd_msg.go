package list

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type EnterKeyMsg string

func EnterKeyCmd() tea.Cmd {
	return func() tea.Msg {
		return ""
	}
}

type EditItemMsg string

func EditItemCmd() tea.Cmd {
	return func() tea.Msg {
		return EditItemMsg("")
	}
}

type OSExecCmdMsg struct {
	cmd ExecuteCmdFunc
}

type ExecuteCmdFunc func(Items) *exec.Cmd

func ExecuteCmd(cmd ExecuteCmdFunc) tea.Cmd {
	return func() tea.Msg {
		return OSExecCmdMsg{cmd: cmd}
	}
}

type SetFocusedViewMsg string

func SetFocusedViewCmd(v string) tea.Cmd {
	return func() tea.Msg {
		return SetFocusedViewMsg(v)
	}
}

type ReturnSelectionsMsg struct{}

func ReturnSelectionsCmd() tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

type UpdateItemsMsg Items

func UpdateItemsCmd(i Items) tea.Cmd {
	return func() tea.Msg {
		return UpdateItemsMsg(i)
	}
}

type ToggleItemListMsg Item

func ToggleItemListCmd(i Item) tea.Cmd {
	return func() tea.Msg {
		return ToggleItemListMsg(i)
	}
}

type toggleItemMsg Item

func toggleItemCmd(idx Item) tea.Cmd {
	return func() tea.Msg {
		return toggleItemMsg(idx)
	}
}

func ToggleAllItemsCmd(l *Model) {
	for _, it := range l.Model.Items() {
		i := it.(Item)
		l.AllItems.Toggle(i.id)
	}
}

type UpdateMenuContentMsg string

func UpdateMenuContentCmd(s string) tea.Cmd {
	return func() tea.Msg {
		return UpdateMenuContentMsg(s)
	}
}

type MenuCmd func(m *Model) tea.Cmd

func DeselectAllItemsCmd(l *Model) {
	for idx, it := range l.Model.Items() {
		i := it.(Item)
		i.isSelected = false
		l.Model.SetItem(idx, i)
	}
}

func SelectAllItemsCmd(l *Model) {
	for idx, it := range l.Model.Items() {
		i := it.(Item)
		i.isSelected = true
		l.Model.SetItem(idx, i)
	}
}

type SelectedItemsMsg Selections

func GetSelectedItemsCmd(i Items) tea.Cmd {
	return func() tea.Msg {
		return SelectedItemsMsg{items: i.GetSelected()}
	}
}

type UpdateDisplayedItemsMsg string

func UpdateDisplayedItemsCmd(opt string) tea.Cmd {
	return func() tea.Msg {
		return UpdateDisplayedItemsMsg(opt)
	}
}

//type ShowNestedListMsg Items
