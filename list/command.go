package list

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

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

type SetItemsMsg Items

func SetItemsCmd(i Items) tea.Cmd {
	return func() tea.Msg {
		return SetItemsMsg(i)
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

type UpdateMenuContentMsg string

func UpdateMenuContentCmd(s string) tea.Cmd {
	return func() tea.Msg {
		return UpdateMenuContentMsg(s)
	}
}

type MenuCmd func(m *Model) tea.Cmd

func ToggleAllItemsCmd(l *Model) {
	for _, it := range l.Items {
		i := it.(Item)
		i.ToggleSelected()
		//l.Items.ToggleSelected(i.id)
	}
}

type UpdateVisibleItemsMsg string

func UpdateVisibleItemsCmd(opt string) tea.Cmd {
	return func() tea.Msg {
		return UpdateVisibleItemsMsg(opt)
	}
}

type UpdateStatusMsg string

func UpdateStatusCmd(status string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusMsg(status)
	}
}

type SetSizeMsg []int

func SetSizeCmd(size []int) tea.Cmd {
	return func() tea.Msg {
		return SetSizeMsg(size)
	}
}

type UpdateInfoContentMsg string

func UpdateInfoContentCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return UpdateInfoContentMsg(content)
	}
}
