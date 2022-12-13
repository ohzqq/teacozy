package form

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateStatusMsg struct{ Msg string }

func UpdateStatusCmd(status string) tea.Cmd {
	return func() tea.Msg {
		return UpdateStatusMsg{Msg: status}
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

type SaveFormFunc func(m *Form) tea.Cmd

type SaveAndExitFormMsg struct {
	Save SaveFormFunc
}

func (f *Form) SaveForm(funcs ...SaveFormFunc) tea.Cmd {
	fn := SaveChangesAsHash
	if len(funcs) > 0 {
		fn = funcs[0]
	}
	return fn(f)
}

func SaveAndExitFormCmd() tea.Cmd {
	return func() tea.Msg {
		return SaveAndExitFormMsg{
			Save: SaveChangesAsHash,
		}
	}
}

type SaveFormAsHashMsg struct {
	Hash map[string]string
}

func SaveFormAsHash(m *Form) tea.Cmd {
	fn := func() tea.Msg {
		m.Hash = m.Fields.StringMap()
		return SaveFormAsHashMsg{Hash: m.Hash}
	}
	return fn
}

func SaveChangesAsHash(m *Form) tea.Cmd {
	fn := func() tea.Msg {
		m.Hash = m.Fields.StringMapChanges()
		return SaveFormAsHashMsg{Hash: m.Hash}
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
