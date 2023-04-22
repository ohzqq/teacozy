package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/frame"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	input textarea.Model

	KeyMap  keys.KeyMap
	Prefix  string
	Style   lipgloss.Style
	help    keys.KeyMap
	current int
}

type Props struct {
	teacozy.Props
	ShowHelp func([]map[string]string)
}

func New() *Component {
	c := &Component{
		input:  textinput.New(),
		Prefix: "> ",
		Style:  lipgloss.NewStyle().Foreground(color.Cyan()),
		KeyMap: DefaultKeyMap(),
	}

	return c
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.input.Prompt = c.Prefix
	c.input.PromptStyle = c.Style
	c.input.KeyMap = keys.TextInputDefaultKeyMap
	return c.input.Focus()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case keys.ToggleItemMsg:
		c.Props().ToggleItems(c.current)
		cmds = append(cmds, keys.LineDown)
	case tea.KeyMsg:
		for _, k := range c.KeyMap {
			if key.Matches(msg, k.Binding) {
				c.input.Reset()
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
	props := c.Props().Props
	props.SetPerPage(h - 1)
	props.Filter(c.input.Value())
	props.Selectable = true
	props.SetCurrent = c.setCurrent
	return lipgloss.JoinVertical(lipgloss.Left, view, teacozy.Renderer(props, w, h))
}

func (c *Component) Initialize(a *frame.App) {
	a.Routes["filter"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		comp := New()
		p := Props{
			Props:       a.ItemProps(),
			ToggleItems: a.ToggleItems,
		}
		a.SetKeyMap(pagy.DefaultKeyMap())
		return comp, comp.Init(p)
	}
}

func (c *Component) setCurrent(i int) {
	c.current = i
}

func DefaultKeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Quit(),
		keys.Help(),
		keys.Toggle(),
		keys.Enter().WithHelp("stop filtering").Cmd(StopFiltering),
		keys.Esc().Cmd(StopFiltering),
	}
	return km
}

func StopFiltering() tea.Msg {
	return keys.ReturnToListMsg{}
}
