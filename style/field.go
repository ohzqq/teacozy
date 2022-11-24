package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

type Field struct {
	Key   lipgloss.Style
	Value lipgloss.Style
}

func DefaultFieldStyles() Field {
	return Field{
		Key:   lipgloss.NewStyle().Foreground(color.Blue),
		Value: lipgloss.NewStyle().Foreground(color.Foreground),
	}
}
