package header

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/color"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]
	Style lipgloss.Style
}

type Props struct {
	Msg string
}

func New() *Component {
	return &Component{
		Style: lipgloss.NewStyle().Foreground(color.Green()),
	}
}

func (c Component) Render(w, h int) string {
	return c.Style.Render(c.Props().Msg)
}

type StatusMsg struct {
	Status string
}

func StatusUpdate(h string) tea.Cmd {
	return func() tea.Msg {
		return StatusMsg{Status: h}
	}
}
