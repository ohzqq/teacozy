package component

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/keys"
)

type TextInput struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[TextInputProps]

	input  textinput.Model
	Prefix string
	Style  lipgloss.Style
}

type TextInputProps struct {
	SetValue func(string)
}

func NewTextInput() *TextInput {
	c := &TextInput{
		input:  textinput.New(),
		Prefix: "> ",
		Style:  lipgloss.NewStyle().Foreground(color.Cyan()),
	}

	c.input.Prompt = c.Prefix
	c.input.PromptStyle = c.Style
	c.input.KeyMap = keys.TextInputDefault()

	return c
}
