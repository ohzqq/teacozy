package teacozy

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

type TextinputStyle struct {
	Prompt      lipgloss.Style
	Text        lipgloss.Style
	Placeholder lipgloss.Style
}

func TextinputDefaultStyle() TextinputStyle {
	s := TextinputStyle{}

	s.Prompt = lipgloss.NewStyle().
		Foreground(color.Green)

	s.Text = lipgloss.NewStyle().
		Foreground(color.Fg)

	s.Placeholder = lipgloss.NewStyle().
		Foreground(color.Grey)

	return s
}

func SetTextinputStyle(input *textinput.Model, style TextinputStyle) {
	input.PromptStyle = style.Prompt
	input.TextStyle = style.Text
	input.PlaceholderStyle = style.Placeholder
}
