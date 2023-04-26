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
	reactea.BasicPropfulComponent[teacozy.Props]

	confirmed bool
	Question  string
	Confirm   ConfirmFunc
	Style     lipgloss.Style
}

type Props struct {
	teacozy.Props
}

type GetConfirmationMsg struct {
	Props
}

type ConfirmFunc func(bool) tea.Cmd

func New() *Component {
	return &Component{
		Style: lipgloss.NewStyle().
			Background(color.Red()).
			Foreground(color.Black()),
	}
}

func (c *Component) Init(props teacozy.Props) tea.Cmd {
	c.UpdateProps(props)
	return frame.ChangeRoute(c)
}

func GetConfirmation(q string, c ConfirmFunc, props teacozy.Props) tea.Cmd {
	confirm := New()
	confirm.Confirm = c
	confirm.Question = q
	return confirm.Init(props)
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
		keys.Esc().AddKeys("q").Cmd(c.Confirmed(false)),
		keys.Yes().Cmd(c.Confirmed(true)),
		keys.No().Cmd(c.Confirmed(false)),
	}
	return keys.NewKeyMap(km...)
}

func (c *Component) Initializer(props teacozy.Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		props.DisableKeys()
		return c, c.Init(props)
	}
}

func (c Component) Name() string {
	return "confirm"
}

func (c *Component) Render(w, h int) string {
	view := c.Style.Render(fmt.Sprintf("%s (y/n)", c.Question))
	return lipgloss.JoinVertical(
		lipgloss.Left,
		view,
		teacozy.Renderer(c.Props(), w, h-1),
	)
}

func (c *Component) Confirmed(y bool) tea.Cmd {
	cmd := c.Confirm(y)
	return tea.Batch(cmd, keys.ChangeRoute("prev"))
}
