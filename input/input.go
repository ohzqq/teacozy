package input

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/keys"
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
	Filter   func(string)
	ShowHelp func([]map[string]string)
}

func New() *Component {
	c := &Component{
		input:  textinput.New(),
		KeyMap: DefaultKeyMap(),
		Prefix: "> ",
		Style:  lipgloss.NewStyle().Foreground(color.Cyan()),
	}

	c.help = append(c.help, c.KeyMap...)
	c.help = append(c.help, keys.TextInput()...)

	return c
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.input.Prompt = c.Prefix
	c.input.PromptStyle = c.Style
	return c.input.Focus()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case keys.ShowHelpMsg:
		c.Props().ShowHelp(c.help.Map())
		cmds = append(cmds, keys.ChangeRoute("help"))
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
		c.Props().Filter(c.input.Value())
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (c *Component) Render(int, int) string {
	return c.input.View()
}

func DefaultKeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Quit(),
		keys.Help(),
		keys.Enter().WithHelp("stop filtering").Cmd(StopFiltering),
		keys.Esc().Cmd(StopFiltering),
	}
	return km
}

func StopFiltering() tea.Msg {
	return keys.ReturnToListMsg{}
}
