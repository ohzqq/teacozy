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

type ToggleItemListMsg int

func ToggleItemListCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleItemListMsg(idx)
	}
}

type ToggleSelectedItemMsg int

func ToggleSelectedItemCmd(idx int) tea.Cmd {
	return func() tea.Msg {
		return ToggleSelectedItemMsg(idx)
	}
}

func ToggleAllItemsCmd(l *Model) {
	for _, it := range l.Items {
		i := it.(Item)
		i.ToggleSelected()
		//l.Items.ToggleSelected(i.id)
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
	for idx, it := range l.List.Items() {
		i := it.(Item)
		i.isSelected = false
		l.List.SetItem(idx, i)
	}
}

func SelectAllItemsCmd(l *Model) {
	for idx, it := range l.List.Items() {
		i := it.(Item)
		i.isSelected = true
		l.List.SetItem(idx, i)
	}
}

type SelectedItemsMsg Selections

func GetSelectedItemsCmd(i Items) tea.Cmd {
	return func() tea.Msg {
		return SelectedItemsMsg{items: i.Selected()}
	}
}

type UpdateVisibleItemsMsg string

func UpdateVisibleItemsCmd(opt string) tea.Cmd {
	return func() tea.Msg {
		return UpdateVisibleItemsMsg(opt)
	}
}

//type ShowNestedListMsg Items
