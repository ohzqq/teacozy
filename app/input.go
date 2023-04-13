package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/style"
)

type Input struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[InputProps]

	textinput textinput.Model

	Placeholder string
	Prompt      string
	Style       style.List
}

type InputProps struct {
	Filter func(string)
}

func NewSearch() *Input {
	tm := &Input{
		Style:     style.ListDefaults(),
		Prompt:    style.PromptPrefix,
		textinput: textinput.New(),
	}
	return tm
}

func (c *Input) Init(props InputProps) tea.Cmd {
	c.UpdateProps(props)
	c.textinput.Prompt = c.Prompt
	c.textinput.PromptStyle = c.Style.Prompt
	c.textinput.Placeholder = c.Placeholder
	return c.textinput.Focus()
}

func (c *Input) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			// Lifted state power! Woohooo
			c.Props().Filter(c.textinput.Value())

			reactea.SetCurrentRoute("default")
			return nil
		}
	}

	var cmd tea.Cmd
	c.textinput, cmd = c.textinput.Update(msg)
	return cmd
}

// Here we are not using width and height, but you can!
func (c *Input) Render(int, int) string {
	return c.textinput.View()
}
