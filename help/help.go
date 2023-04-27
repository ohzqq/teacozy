package help

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	KeyMap keys.KeyMap
}

func New(km keys.KeyMap) *Component {
}

func (m *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(m)
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, k := range m.KeyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	m.view, cmd = m.view.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Component) Render(w, h int) string {
	m.view.SetWidth(w)
	m.view.SetHeight(h)
	return m.view.View()
}

func DefaultKeyMap() keys.KeyMap {
	km := []*keys.Binding{
		keys.Esc().AddKeys("q").Cmd(keys.ChangeRoute("prev")),
	}
	return keys.NewKeyMap(km...)
}
