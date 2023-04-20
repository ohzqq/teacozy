package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	mainRouter reactea.Component[router.Props]

	*Model
}

func New() *Component {
	m := Component{}
	return &m
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.Model = NewModel(props)
	return nil
}

func (m *Component) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Component) Render(w, h int) string {
	m.SetWidth(w)
	m.SetHeight(h)
	return m.Model.View()
}
