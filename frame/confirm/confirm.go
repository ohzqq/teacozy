package confirm

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/frame"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	confirmed bool
}

type Props struct {
	Question string
	Confirm  Confirm
}

type GetConfirmationMsg struct {
	Props
}

type Confirm func(bool) tea.Cmd

func New() *Component {
	return &Component{}
}

func GetConfirmation(q string, c Confirm) tea.Cmd {
	return func() tea.Msg {
		return GetConfirmationMsg{
			Props: Props{
				Question: q,
				Confirm:  c,
			},
		}
	}
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, k := range c.KeyMap().Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}
	return tea.Batch(cmds...)
}

func (c *Component) KeyMap() keys.KeyMap {
	km := []*keys.Binding{
		keys.Quit(),
		keys.Yes().Cmd(c.Confirmed(true)),
		keys.No().Cmd(c.Confirmed(false)),
	}
	return keys.NewKeyMap(km...)
}

func (c *Component) Initialize(a *frame.App) {
	a.Routes["confirm"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		p := a.Confirm
		return component, component.Init(p)
	}
}

func (c *Component) Render(w, h int) string {
	q := fmt.Sprintf("%s (y/n)", c.Props().Question)
	return lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()).Render(q)
}

func (c *Component) Confirmed(y bool) tea.Cmd {
	cmd := c.Props().Confirm(y)
	return tea.Batch(cmd, keys.ChangeRoute("prev"))
	//return tea.Batch(keys.ChangeRoute("prev"), cmd)
}
