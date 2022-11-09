package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
)

type SaveAsHashMsg struct{}

func SaveAsHashCmd() tea.Cmd {
	return func() tea.Msg {
		return SaveAsHashMsg{}
	}
}

type EditInfoMsg struct{}

func EditInfoCmd() tea.Cmd {
	return func() tea.Msg {
		return EditInfoMsg{}
	}
}

type EditItemMsg struct {
	info.Field
}

func EditItemCmd(field info.Field) tea.Cmd {
	return func() tea.Msg {
		return EditItemMsg{Field: field}
	}
}

type UpdateContentMsg struct {
	info.Field
}

func UpdateContentCmd(key, val string) tea.Cmd {
	return func() tea.Msg {
		return UpdateContentMsg{Field: info.NewDefaultField(key, val)}
	}
}
