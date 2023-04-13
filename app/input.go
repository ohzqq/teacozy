package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Input struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[InputProps]

	textinput textinput.Model
}

type InputProps struct {
	Filter func(string)
}

func NewSearch() *Input {
	return &Input{textinput: textinput.New()}
}

func (c *Input) Init(props InputProps) tea.Cmd {
	c.UpdateProps(props)
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
	return fmt.Sprintf("Enter your name: %s\nAnd press [ Enter ]", c.textinput.View())
}
