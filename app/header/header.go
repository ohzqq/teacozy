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
		Style: lipgloss.NewStyle().Background(color.Purple()).Foreground(color.Black()),
	}
}

func (c Component) Render(w, h int) string {
	return c.Style.Render(c.Props().Msg)
}

type UpdateHeaderMsg struct {
	Header string
}

func UpdateHeader(h string) tea.Cmd {
	return func() tea.Msg {
		return UpdateHeaderMsg{Header: h}
	}
}
