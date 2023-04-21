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
	"github.com/ohzqq/teacozy/frame"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	input textinput.Model

	KeyMap keys.KeyMap
	Prefix string
	Style  lipgloss.Style
	help   keys.KeyMap
}

type Props struct {
	teacozy.Props
	ShowHelp    func([]map[string]string)
	ToggleItems func(...int)
}

func New() *Component {
	c := &Component{
		input: textinput.New(),
		//KeyMap: DefaultKeyMap(),
		Prefix: "> ",
		Style:  lipgloss.NewStyle().Foreground(color.Cyan()),
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
		c.Props().ToggleItems(c.Props().Cursor())
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
	return lipgloss.JoinVertical(lipgloss.Left, view, teacozy.Renderer(props, w, h))
}

func (c *Component) Initialize(a *frame.App) {
	a.Routes["filter"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		comp := New()
		p := Props{
			Props:       a.ItemProps(),
			ToggleItems: a.ToggleItems,
		}
		a.SetKeyMap(DefaultKeyMap())
		return comp, comp.Init(p)
	}
}

func DefaultKeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Quit(),
		keys.Help(),
		keys.Toggle(),
		keys.Enter().WithHelp("stop filtering").Cmd(StopFiltering),
		keys.Esc().Cmd(StopFiltering),
	}
	km = append(km, pagy.DefaultKeyMap()...)
	return km
}

func StopFiltering() tea.Msg {
	return keys.ReturnToListMsg{}
}