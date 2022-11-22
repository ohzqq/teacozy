package form

import tea "github.com/charmbracelet/bubbletea"

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
	Data FieldData
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
