package info

import tea "github.com/charmbracelet/bubbletea"

type EditInfoMsg struct{}

func UpdateInfoCmd() tea.Cmd {
	return func() tea.Msg {
		return EditInfoMsg{}
	}
}

type EditItemMsg struct {
	Field
}

func EditItemCmd(field Field) tea.Cmd {
	return func() tea.Msg {
		return EditItemMsg{Field: field}
	}
}

type UpdateContentMsg struct {
	Field
}

func UpdateContentCmd(key, val string) tea.Cmd {
	return func() tea.Msg {
		return UpdateContentMsg{Field: NewDefaultField(key, val)}
	}
}
