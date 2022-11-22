package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func itemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string
		if i, ok := m.SelectedItem().(*Field); ok {
			title = i.Title()
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, Edit):
				return EditFormItemCmd(curItem)
			}
		}
	}
}

var Edit = key.NewBinding(
	key.WithKeys("e"),
	key.WithHelp("e", "edit meta"),
)
