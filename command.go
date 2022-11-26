package teacozy

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	*Item
}

func ItemChangedCmd(item *Item) tea.Cmd {
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
type CmdFunc func(m tea.Model) tea.Cmd

type MenuFunc func(m tea.Model) tea.Cmd

type UpdateMenuContentMsg string

func UpdateMenuContentCmd(s string) tea.Cmd {
	return func() tea.Msg {
		return UpdateMenuContentMsg(s)
	}
}

type HideMenuMsg struct{}

func HideMenuCmd() tea.Cmd {
	return func() tea.Msg {
		return HideMenuMsg{}
	}
}

type ShowMenuMsg struct{ *Menu }

func ShowMenuCmd(menu *Menu) tea.Cmd {
	return func() tea.Msg {
		return ShowMenuMsg{Menu: menu}
	}
}

type ChangeMenuMsg struct{ *Menu }

func GoToMenuCmd(menu *Menu) MenuFunc {
	return func(m tea.Model) tea.Cmd {
		return ChangeMenuCmd(menu)
	}
}

func ChangeMenuCmd(menu *Menu) tea.Cmd {
	return func() tea.Msg {
		return ChangeMenuMsg{Menu: menu}
	}
}

// form commands
type SaveFormFunc func(m *List) tea.Cmd

type SaveForm func(m *Form) tea.Cmd

type SaveAndExitFormMsg struct {
	Exit SaveFormFunc
	Save SaveForm
}

func SaveAndExitFormCmd(fn SaveFormFunc) tea.Cmd {
	return func() tea.Msg {
		return SaveAndExitFormMsg{Exit: fn}
	}
}

func SaveFormCmd(fn SaveForm) tea.Cmd {
	return func() tea.Msg {
		return SaveAndExitFormMsg{Save: fn}
	}
}

type ExitFormMsg struct{}

func ExitFormCmd() tea.Cmd {
	return func() tea.Msg {
		return ExitFormMsg{}
	}
}

type SaveFormAsHashMsg struct{}

func SaveFormAsHashCmd(m *List) tea.Cmd {
	fn := func() tea.Msg {
		m.Hash = make(map[string]string)
		for _, item := range m.Items.Flat() {
			m.Hash[item.Key()] = item.Value()
		}
		return SaveFormAsHashMsg{}
	}
	return fn
}

func SaveFormAsHash(m *Form) tea.Cmd {
	fn := func() tea.Msg {
		m.Hash = make(map[string]string)
		for _, item := range m.Items.Flat() {
			m.Hash[item.Key()] = item.Value()
		}
		return SaveFormAsHashMsg{}
	}
	return fn
}

type EditFormItemMsg struct {
	Data Field
	*Item
}

func EditFormItemCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return EditFormItemMsg{Data: item.Data, Item: item}
	}
}

type FormChangedMsg struct {
	*Item
}

func FormChangedCmd() tea.Cmd {
	return func() tea.Msg {
		return FormChangedMsg{}
	}
}

type ConfirmMenuMsg bool

func ConfirmMenuCmd(confirm bool) tea.Cmd {
	return func() tea.Msg {
		return ConfirmMenuMsg(confirm)
	}
}

type ConfirmFormSaveMsg struct{}

func ConfirmFormSaveCmd() tea.Cmd {
	return func() tea.Msg {
		return ConfirmFormSaveMsg{}
	}
}

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
	*FormData
}

func EditInfoCmd(f *FormData) tea.Cmd {
	return func() tea.Msg {
		return EditInfoMsg{
			FormData: f,
		}
	}
}

// item commands
type ToggleItemChildrenMsg struct{ *Item }

func ToggleItemChildrenCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return ToggleItemChildrenMsg{Item: item}
	}
}

type ToggleSelectedItemMsg struct{ *Item }

func ToggleSelectedItemCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return ToggleSelectedItemMsg{Item: item}
	}
}

type ShowItemInfoMsg struct{ *Item }

func ShowItemInfoCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return ShowItemInfoMsg{Item: item}
	}
}

type EditItemValueMsg struct{ *Item }

func EditItemValueCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return EditItemValueMsg{Item: item}
	}
}

// list commands
type ReturnSelectionsMsg struct{}

func ReturnSelectionsCmd() tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

type ExitSelectionsListMsg struct{}

func (m *List) ExitSelectionsListCmd() tea.Cmd {
	return func() tea.Msg {
		m.SelectionList = false
		return ExitSelectionsListMsg{}
	}
}

func ToggleAllItemsCmd(l *List) {
	l.Items.ToggleAllSelectedItems()
}

type UpdateVisibleItemsMsg string

func UpdateVisibleItemsCmd(opt string) tea.Cmd {
	return func() tea.Msg {
		return UpdateVisibleItemsMsg(opt)
	}
}

func (m *List) ShowVisibleItemsCmd() tea.Cmd {
	return func() tea.Msg {
		return UpdateVisibleItemsMsg("visible")
	}
}

func (m *List) ShowSelectedItemsCmd() tea.Cmd {
	return func() tea.Msg {
		return UpdateVisibleItemsMsg("selected")
	}
}

type UpdateStatusMsg struct{ Msg string }

func UpdateStatusCmd(status string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusMsg{Msg: status}
	}
}

type SortItemsMsg struct{ Items []*Item }

func SortItemsCmd(items []*Item) tea.Cmd {
	return func() tea.Msg {
		return SortItemsMsg{Items: items}
	}
}

type SetListItemMsg struct {
	Item list.Item
}

func SetListItemCmd(item list.Item) tea.Cmd {
	return func() tea.Msg {
		return SetListItemMsg{Item: item}
	}
}

type SetItemMsg struct{ *Item }

func SetItemCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return SetItemMsg{Item: item}
	}
}

type SetItemsMsg struct{ Items []list.Item }

func SetItemsCmd(items []list.Item) tea.Cmd {
	return func() tea.Msg {
		return SetItemsMsg{Items: items}
	}
}
