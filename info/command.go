package info

import tea "github.com/charmbracelet/bubbletea"

type EditInfoMsg struct{}

func UpdateInfoCmd() tea.Cmd {
	return func() tea.Msg {
		return EditInfoMsg{}
	}
}

type EditItemMsg struct {
	FormField
}

func EditItemCmd(field FormField) tea.Cmd {
	return func() tea.Msg {
		return EditItemMsg{FormField: field}
	}
}

type UpdateContentMsg struct {
	Field
}

func UpdateContentCmd(key, val string) tea.Cmd {
	return func() tea.Msg {
		return UpdateContentMsg{Field: NewField(key, val)}
	}
}
