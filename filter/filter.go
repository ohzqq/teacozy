package filter

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	reactea.BasicComponent
	//reactea.BasicPropfulComponent[Props]

	input textinput.Model

	KeyMap keys.KeyMap
	Prefix string
	Style  lipgloss.Style
	help   keys.KeyMap

	props Props
}

type Props struct {
	SetKeyMap func(keys.KeyMap)
	Filter    func(string, int, int) string
}

func New() *Component {
	c := &Component{
		input:  textinput.New(),
		Prefix: "> ",
		Style:  lipgloss.NewStyle().Foreground(color.Cyan()),
		KeyMap: DefaultKeyMap(),
	}

	c.input.Prompt = c.Prefix
	c.input.PromptStyle = c.Style
	c.input.KeyMap = keys.TextInputDefault()

	return c
}

func (c Component) Props() Props {
	return c.props
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case StopFilteringMsg:
		c.input.Blur()
		c.Props().SetKeyMap(keys.VimKeyMap())
		return keys.ChangeRoute("prev")
	case tea.KeyMsg:
		for _, k := range c.KeyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	if c.input.Focused() {
		c.input, cmd = c.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	view := c.input.View()
	f := c.Props().Filter(c.input.Value(), w, h-1)
	if c.input.Focused() {
		return lipgloss.JoinVertical(lipgloss.Left, view, f)
	}
	return ""
}

func (c *Component) Initializer(props teacozy.Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		comp := New()
		p := Props{
			SetKeyMap: props.SetKeyMap,
			Filter:    props.Filter,
		}
		p.SetKeyMap(keys.DefaultKeyMap())
		comp.props = p
		return comp, comp.input.Focus()
	}
}

func (c Component) Name() string {
	return "filter"
}

func DefaultKeyMap() keys.KeyMap {
	km := []*keys.Binding{
		keys.Quit(),
		keys.Help(),
		keys.Enter().WithHelp("stop filtering").Cmd(StopFiltering),
		keys.Esc().Cmd(StopFiltering),
	}
	return keys.NewKeyMap(km...)
}

type StopFilteringMsg struct{}

func StopFiltering() tea.Msg {
	return StopFilteringMsg{}
}
