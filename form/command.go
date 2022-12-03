package form

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SaveFormFunc func(m *Form) tea.Cmd

type SaveAndExitFormMsg struct {
	Save SaveFormFunc
}

func SaveAndExitFormCmd() tea.Cmd {
	return func() tea.Msg {
		return SaveAndExitFormMsg{}
	}
}

type SaveFormAsHashMsg struct{}

func SaveFormAsHash(m *Form) tea.Cmd {
	fn := func() tea.Msg {
		m.Hash = make(map[string]string)
		for _, item := range m.Fields.fields {
			m.Hash[item.Key()] = item.Value()
		}
		return SaveFormAsHashMsg{}
	}
	return fn
}

type ExitFormMsg struct{}

func ExitFormCmd() tea.Cmd {
	return func() tea.Msg {
		return ExitFormMsg{}
	}
}

type ViewFormMsg struct{}

func ViewFormCmd() tea.Cmd {
	return func() tea.Msg {
		return ViewFormMsg{}
	}
}

type HideFormMsg struct{}

func HideFormCmd() tea.Cmd {
	return func() tea.Msg {
		return HideFormMsg{}
	}
}

type EditFormItemMsg struct {
	*Field
}

func EditFormItemCmd(item *Field) tea.Cmd {
	return func() tea.Msg {
		return EditFormItemMsg{Field: item}
	}
}

type FormChangedMsg struct {
	*Field
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
