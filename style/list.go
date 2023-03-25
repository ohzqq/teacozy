package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

const (
	Bullet   = "•"
	Ellipsis = "…"
)

const (
	CursorPrefix     = "x"
	PromptPrefix     = "> "
	SelectedPrefix   = "x"
	UnselectedPrefix = " "
)

var (
	Prompt      = lipgloss.NewStyle().Foreground(color.Cyan())
	Cursor      = lipgloss.NewStyle().Foreground(color.Green())
	Unselected  = lipgloss.NewStyle().Foreground(color.Fg())
	Selected    = lipgloss.NewStyle().Foreground(color.Grey())
	Current     = lipgloss.NewStyle().Foreground(color.Grey())
	Subdued     = lipgloss.NewStyle().Foreground(color.White())
	VerySubdued = lipgloss.NewStyle().Foreground(color.Grey())
	Foreground  = lipgloss.NewStyle().Foreground(color.Fg())
)

type List struct {
	SelectedPrefix   lipgloss.Style
	Text             lipgloss.Style
	Match            lipgloss.Style
	Cursor           lipgloss.Style
	UnselectedPrefix lipgloss.Style
	Header           lipgloss.Style
	Prompt           lipgloss.Style
}

type ListItem struct {
	SelectedPrefix   lipgloss.Style
	UnselectedPrefix lipgloss.Style
	Match            lipgloss.Style
	Cursor           lipgloss.Style
	Text             lipgloss.Style
}
