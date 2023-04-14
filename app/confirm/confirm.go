package confirm

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/app/status"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/keys"
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

type ConfirmedMsg struct {
	Action tea.Cmd
}

func New() *Component {
	return &Component{}
}

func Action(q string, a tea.Cmd) tea.Cmd {
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
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, k := range c.KeyMap() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}
	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	q := fmt.Sprintf("%s (y/n)", c.Props().Question)
	return lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()).Render(q)
}

func (c *Component) Confirmed() tea.Msg {
	return status.StatusMsg{Status: fmt.Sprintf("%s %v", c.Props().Question, c.confirmed)}
}

func Confirmed(a tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return ConfirmedMsg{
			Action: a,
		}
	}
}

func (c *Component) KeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Quit(),
		keys.Yes().Cmd(Confirmed(c.Props().Action)),
		keys.No().Cmd(Confirmed(nil)),
	}
	return km
}
