package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[CProps]

	mainRouter reactea.Component[router.Props]

	Model *Model
}

type CProps struct {
	Props
	SetCursor func(int)
	SetStart  func(int)
	SetEnd    func(int)
}

func New() *Component {
	m := Component{}
	return &m
}

func (c *Component) Init(props CProps) tea.Cmd {
	c.UpdateProps(props)
	c.Model = NewModel(props.Props)
	return nil
}

func (m *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(m)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Component) AfterUpdate() tea.Cmd {
	m.Props().SetCursor(m.Model.Cursor)
	m.Props().SetStart(m.Model.Start)
	m.Props().SetEnd(m.Model.End)
	//return keys.ChangeRoute("default")
	return nil
}

func (m *Component) Render(w, h int) string {
	m.Model.SetWidth(w)
	m.Model.SetHeight(h)
	return m.Model.View()
}
