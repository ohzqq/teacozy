package confirm

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
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
	teacozy.Props
	Question string
	Confirm  ConfirmFunc
}

type GetConfirmationMsg struct {
	Props
}

type ConfirmFunc func(bool) tea.Cmd

func New() *Component {
	return &Component{}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	return frame.ChangeRoute(c)
}

func ConfirmThing() tea.Cmd {
	return frame.ChangeRoute(New())
}

func GetConfirmation(q string, c ConfirmFunc, props teacozy.Props) tea.Cmd {
	confirm := New()
	p := Props{
		Props:    props,
		Question: q,
		Confirm:  c,
	}
	return confirm.Init(p)
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

func (c *Component) Initializer(props teacozy.Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		p := Props{
			Props: props,
		}
		p.DisableKeys()
		return component, component.Init(p)
	}
}

func (c Component) Name() string {
	return "confirm"
}

func (c *Component) Render(w, h int) string {
	q := fmt.Sprintf("%s (y/n)", c.Props().Question)
	view := lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()).Render(q)
	props := c.Props().Props
	return lipgloss.JoinVertical(lipgloss.Left, view, teacozy.Renderer(props, w, h-1))
}

func (c *Component) Confirmed(y bool) tea.Cmd {
	cmd := c.Props().Confirm(y)
	return tea.Batch(cmd, keys.ChangeRoute("prev"))
	//return tea.Batch(keys.ChangeRoute("prev"), cmd)
}
