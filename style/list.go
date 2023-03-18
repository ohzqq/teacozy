package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

const (
	Bullet   = "•"
	Ellipsis = "…"
)

const (
	Cursor           = "> "
	Prompt           = "> "
	SelectedPrefix   = "◉ "
	UnselectedPrefix = "○ "
	CursorPrefix     = "○ "
)

var PromptStyle = lipgloss.NewStyle().Foreground(color.Cyan)
var CursorStyle = lipgloss.NewStyle().Foreground(color.Green)
var UnselectedStyle = lipgloss.NewStyle().Foreground(color.Foreground)
var SelectedStyle = lipgloss.NewStyle().Foreground(color.Grey)
var CurrentStyle = lipgloss.NewStyle().Foreground(color.Grey)

func ListStyles() list.Styles {
	verySubduedColor := color.Grey
	subduedColor := color.White

	var s list.Styles

	s.TitleBar = lipgloss.NewStyle().
		Padding(0, 0, 0, 0)

	s.Title = lipgloss.NewStyle().
		Background(color.Purple).
		Foreground(color.Black).
		Padding(0, 1)

	s.Spinner = lipgloss.NewStyle().
		Foreground(color.Cyan)

	s.FilterPrompt = lipgloss.NewStyle().
		Foreground(color.Pink)

	s.FilterCursor = lipgloss.NewStyle().
		Foreground(color.Yellow)

	s.DefaultFilterCharacterMatch = lipgloss.NewStyle().
		Underline(true)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(color.Blue).
		Padding(0, 0, 1, 2)

	s.StatusEmpty = lipgloss.NewStyle().
		Foreground(subduedColor)

	s.StatusBarActiveFilter = lipgloss.NewStyle().
		Foreground(color.Purple)

	s.StatusBarFilterCount = lipgloss.NewStyle().
		Foreground(verySubduedColor)

	s.NoItems = lipgloss.NewStyle().
		Foreground(color.Grey)

	s.ArabicPagination = lipgloss.NewStyle().
		Foreground(subduedColor)

	s.PaginationStyle = lipgloss.NewStyle().
		PaddingLeft(2) //nolint:gomnd

	s.HelpStyle = lipgloss.NewStyle().
		Padding(1, 0, 0, 2)

	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(color.Pink).
		SetString(Bullet)

	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(Bullet)

	s.DividerDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(" " + Bullet + " ")

	return s
}
