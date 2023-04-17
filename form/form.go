package form

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/view"
)

type Component struct {
	*view.Component

	KeyMap keys.KeyMap
}

func New() *Component {
	return &Component{
		Component: view.New(),
		KeyMap:    DefaultKeyMap(),
	}
}

func (m *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(m)
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Props().Editable {
			if k := keys.Edit(); key.Matches(msg, k.Binding) {
				return nil
			}
		}
		for _, k := range m.KeyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	cmd = m.Component.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func DefaultKeyMap() keys.KeyMap {
	return keys.KeyMap{
		keys.Esc().AddKeys("q").Cmd(keys.ChangeRoute("prev")),
	}
}
