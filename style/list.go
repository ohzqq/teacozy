package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

const (
	CursorPrefix     = "x"
	PromptPrefix     = "> "
	SelectedPrefix   = "x"
	UnselectedPrefix = " "
)

var (
	Prompt     = lipgloss.NewStyle().Foreground(color.Cyan())
	Foreground = lipgloss.NewStyle().Foreground(color.Fg())
	Label      = lipgloss.NewStyle().Foreground(color.Purple())
)
