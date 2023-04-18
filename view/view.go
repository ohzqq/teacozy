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
}

type Props struct {
	Fields []item.Item
}

func New() *Component {
	m := Component{}
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

	var s []string
	for _, i := range m.Props().Fields {
		f := i.Render(m.Viewport.Width, m.Viewport.Height)
		s = append(s, lipgloss.NewStyle().Width(m.Width()).Render(f))
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

// Height returns the viewport height of the list.
func (m Component) Height() int {
	return m.Viewport.Height
}

// Width returns the viewport width of the list.
func (m Component) Width() int {
	return m.Viewport.Width
}
