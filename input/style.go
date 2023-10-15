package input

import "github.com/charmbracelet/lipgloss"

type Style struct {
	Prompt      lipgloss.Style
	Text        lipgloss.Style
	Placeholder lipgloss.Style
	Cursor      lipgloss.Style
}

func DefaultStyle() Style {
	s := Style{}
}
