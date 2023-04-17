package item

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

const (
	Cursor     = "x"
	Prompt     = "> "
	Selected   = "x"
	Unselected = " "
)

type Style struct {
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	Cursor     lipgloss.Style
	Match      lipgloss.Style
	Label      lipgloss.Style
}

func DefaultStyle() Style {
	return Style{
		Selected:   lipgloss.NewStyle().Foreground(color.Grey()),
		Unselected: lipgloss.NewStyle().Foreground(color.Fg()),
		Cursor:     lipgloss.NewStyle().Foreground(color.Green()),
		Match:      lipgloss.NewStyle().Foreground(color.Cyan()),
		Label:      lipgloss.NewStyle().Foreground(color.Purple()),
	}
}
