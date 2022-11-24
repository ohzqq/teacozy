package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	keybind "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
)

func itemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		//var title string
		var curItem *Field
		if i, ok := m.SelectedItem().(*Field); ok {
			curItem = i
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keybind.EditField):
				return EditFormItemCmd(curItem)
			}
		}
		return nil
	}
	d.SetSpacing(0)
	help := []key.Binding{keybind.EditField}
	d.ShortHelpFunc = func() []key.Binding {
		return help
	}
	d.Styles = style.NewDefaultItemStyles()

	return d
}
