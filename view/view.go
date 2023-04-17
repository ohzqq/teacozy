package view

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/item"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Viewport viewport.Model
	fields   []item.Item
}

type Props struct {
	Fields item.Choices
}

func New() *Component {
	m := Component{}
	return &m
}

func (m *Component) Init(props Props) tea.Cmd {
	m.UpdateProps(props)
	m.fields = item.ChoiceMapToItems(m.Props().Fields)
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

	var s []string
	for _, i := range m.fields {
		s = append(s, i.Render(m.Viewport.Width, m.Viewport.Height))
	}

	m.Viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, s...),
	)
	return m.Viewport.View()
}

func (m *Component) SetWidth(w int) {
	m.Viewport.Width = w
}

func (m *Component) SetHeight(h int) {
	m.Viewport.Height = h
}
