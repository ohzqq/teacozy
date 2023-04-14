package confirm

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/app/status"
	"github.com/ohzqq/teacozy/color"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	confirmed bool
}

type Props struct {
	Question string
	Action   tea.Cmd
}

type GetConfirmationMsg struct {
	Props
}

func New() *Component {
	return &Component{}
}

func Confirm(q string, a tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return GetConfirmationMsg{
			Props: Props{
				Question: q,
				Action:   a,
			},
		}
	}
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			c.confirmed = true
			return c.Props().Action
		case "n":
			return NotConfirmed
		}
	}
	return nil
}

func (c *Component) Render(w, h int) string {
	q := fmt.Sprintf("%s (y/n)", c.Props().Question)
	return lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()).Render(q)
}

func (c *Component) Confirmed() tea.Msg {
	return status.StatusMsg{Status: fmt.Sprintf("%v", c.confirmed)}
}

type NotConfirmedMsg struct{}

func NotConfirmed() tea.Msg {
	return NotConfirmedMsg{}
}
