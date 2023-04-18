package help

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/view"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	view *view.Model

	KeyMap keys.KeyMap
}

type Props struct {
	view.Props
}

func New() *Component {
	return &Component{
		KeyMap: DefaultKeyMap(),
	}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.view = view.New(props.Props)
	return nil
}

func (m *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(m)
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, k := range m.KeyMap {
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
	return keys.KeyMap{
		keys.Esc().AddKeys("q").Cmd(keys.ChangeRoute("prev")),
	}
}
