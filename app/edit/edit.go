package edit

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/style"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	input textarea.Model

	Placeholder string
	Prompt      string
	Style       style.List
}

type Props struct {
	Value string
	Edit  func(string)
}

func New() *Component {
	tm := &Component{
		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
		input:  textarea.New(),
	}
	return tm
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.input.Prompt = c.Prompt
	c.input.PromptStyle = c.Style.Prompt
	c.input.Placeholder = c.Placeholder
	m.input.SetValue(props.Value)
	return c.input.Focus()
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			c.input.Reset()
			return message.StopFiltering()
		}
	}

	var cmd tea.Cmd
	c.input, cmd = c.input.Update(msg)
	c.Props().Edit(c.input.Value())
	return cmd
}

func (c *Component) Render(int, int) string {
	return c.input.View()
}
