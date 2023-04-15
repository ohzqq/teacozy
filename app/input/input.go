package input

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/style"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	input textinput.Model

	Placeholder string
	Prompt      string
	Style       style.List
	KeyMap      keys.KeyMap
}

type Props struct {
	Filter func(string)
}

func New() *Component {
	tm := &Component{
		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
		input:  textinput.New(),
		KeyMap: DefaultKeyMap(),
	}
	return tm
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.input.Prompt = c.Prompt
	c.input.PromptStyle = c.Style.Prompt
	c.input.Placeholder = c.Placeholder
	return c.input.Focus()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, k := range c.KeyMap {
			if key.Matches(msg, k.Binding) {
				c.input.Reset()
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	c.input, cmd = c.input.Update(msg)
	c.Props().Filter(c.input.Value())
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *Component) Render(int, int) string {
	return c.input.View()
}

func DefaultKeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Quit(),
		keys.Enter().WithHelp("stop filtering").Cmd(StopFiltering),
		keys.Esc(),
	}
	return km
}

type StartFilteringMsg struct{}

func StartFiltering() tea.Msg {
	return StartFilteringMsg{}
}

type StopFilteringMsg struct{}

func StopFiltering() tea.Msg {
	return StopFilteringMsg{}
}
