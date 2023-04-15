package view

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/style"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Style    style.List
	Viewport viewport.Model
}

type Props struct {
	Fields item.Choices
}

func New() *Component {
	m := Component{
		Style: style.ListDefaults(),
	}
	return &m
}

func (m *Component) Init(props Props) tea.Cmd {
	m.UpdateProps(props)
	m.Viewport = viewport.New(0, 0)
	return nil
}

func (m *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(m)
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Component) Render(w, h int) string {
	m.SetWidth(w)
	m.SetHeight(h)
	return m.Viewport.View()
}

func (m *Component) SetWidth(w int) {
	m.Viewport.Width = w
}

func (m *Component) SetHeight(h int) {
	m.Viewport.Height = h
}
