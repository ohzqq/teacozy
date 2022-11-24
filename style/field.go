package style

import "github.com/charmbracelet/lipgloss"

type Field struct {
	Key   lipgloss.Style
	Value lipgloss.Style
}

func DefaultFieldStyles() Field {
	return Field{
		Key:   lipgloss.NewStyle().Foreground(Color.Blue),
		Value: lipgloss.NewStyle().Foreground(Color.DefaultFg),
	}
}
