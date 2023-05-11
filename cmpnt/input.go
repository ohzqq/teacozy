package cmpnt

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type TextInput struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[InputProps]

	Input textinput.Model
}

type TextArea struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[InputProps]

	area textarea.Model
}

type InputProps struct {
	SetValue func(string)
}

func NewTextInput() *TextInput {
	input := textinput.New()
	return &TextInput{
		Input: input,
	}
}

func (c *TextInput) Init(props InputProps) tea.Cmd {
	c.UpdateProps(props)
	return c.Input.Focus()
}

func (c *TextInput) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			c.Props().SetValue(c.Input.Value())
			c.Input.Blur()
			return nil
		}
	}

	var cmd tea.Cmd
	c.Input, cmd = c.Input.Update(msg)
	return cmd
}

func (c TextInput) Render(w, h int) string {
	return c.Input.View()
}

func NewTextArea() *TextArea {
	input := textarea.New()
	return &TextArea{
		area: input,
	}
}

func (c *TextArea) Init(props InputProps) tea.Cmd {
	c.UpdateProps(props)
	return c.area.Focus()
}
