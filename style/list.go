package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

const (
	Bullet   = "•"
	Ellipsis = "…"
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
	Label       = lipgloss.NewStyle().Foreground(color.Purple())
	Footer      = lipgloss.NewStyle().Foreground(color.Blue())
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
	Match lipgloss.Style
	Text  lipgloss.Style
	Label lipgloss.Style
	ItemPrefix
}

type ItemPrefix struct {
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	Cursor     lipgloss.Style
}

func ListDefaults() List {
	var s List
	s.Cursor = Cursor
	s.SelectedPrefix = Selected
	s.UnselectedPrefix = Unselected
	s.Text = Foreground
	s.Match = lipgloss.NewStyle().Foreground(color.Cyan())
	s.Header = lipgloss.NewStyle().Foreground(color.Purple())
	s.Prompt = Prompt
	return s
}
