package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	Bullet   = "•"
	Ellipsis = "…"
)

func ListStyles() list.Styles {
	verySubduedColor := Color.Grey
	subduedColor := Color.White

	var s list.Styles

	s.TitleBar = lipgloss.NewStyle().
		Padding(0, 0, 0, 0)

	s.Title = lipgloss.NewStyle().
		Background(Color.Purple).
		Foreground(Color.Black).
		Padding(0, 1)

	s.Spinner = lipgloss.NewStyle().
		Foreground(Color.Cyan)

	s.FilterPrompt = lipgloss.NewStyle().
		Foreground(Color.Pink)

	s.FilterCursor = lipgloss.NewStyle().
		Foreground(Color.Yellow)

	s.DefaultFilterCharacterMatch = lipgloss.NewStyle().
		Underline(true)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(Color.Blue).
		Padding(0, 0, 1, 2)

	s.StatusEmpty = lipgloss.NewStyle().
		Foreground(subduedColor)

	s.StatusBarActiveFilter = lipgloss.NewStyle().
		Foreground(Color.Purple)

	s.StatusBarFilterCount = lipgloss.NewStyle().
		Foreground(verySubduedColor)

	s.NoItems = lipgloss.NewStyle().
		Foreground(Color.Grey)

	s.ArabicPagination = lipgloss.NewStyle().
		Foreground(subduedColor)

	s.PaginationStyle = lipgloss.NewStyle().
		PaddingLeft(2) //nolint:gomnd

	s.HelpStyle = lipgloss.NewStyle().
		Padding(1, 0, 0, 2)

	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(Color.Pink).
		SetString(Bullet)

	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(Bullet)

	s.DividerDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(" " + Bullet + " ")

	return s
}
