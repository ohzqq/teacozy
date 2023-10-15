package input

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

type Style struct {
	Prompt      lipgloss.Style
	Text        lipgloss.Style
	Placeholder lipgloss.Style
}

func DefaultStyle() Style {
	s := Style{}

	s.Prompt = lipgloss.NewStyle().
		Foreground(color.Green)

	s.Text = lipgloss.NewStyle().
		Foreground(color.Fg)

	s.Placeholder = lipgloss.NewStyle().
		Foreground(color.Grey)

	return s
}
